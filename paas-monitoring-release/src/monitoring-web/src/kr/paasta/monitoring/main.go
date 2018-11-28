package main

import (
	"bufio"
	"io"
	"os"
	"log"
	"strings"
	"strconv"
	"net/http"
	"github.com/influxdata/influxdb/client/v2"
	_ "github.com/go-sql-driver/mysql"
	"kr/paasta/monitoring/iaas/model"
	bm "kr/paasta/monitoring/paas/model"
	"kr/paasta/monitoring/handlers"
	"kr/paasta/monitoring/utils"
	"fmt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"github.com/cihub/seelog"
	"gopkg.in/olivere/elastic.v3"
	"github.com/monasca/golang-monascaclient/monascaclient"
	"github.com/gophercloud/gophercloud"
	"github.com/alexedwards/scs"
	"time"
	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/go-redis/redis"
	"github.com/cloudfoundry-community/gogobosh"
)

type Config map[string]string

type DBConfig struct {
	DbType string
	UserName string
	UserPassword string
	Host string
	Port string
	DbName string
}

type MemberInfo struct{
	UserId       		string     `gorm:"type:varchar(50);primary_key"`
	UserPw       		string     `gorm:"type:varchar(500);null;"`
	UserEmail       	string     `gorm:"type:varchar(100);null;"`
	UserNm       		string     `gorm:"type:varchar(100);null;"`
	IaasUserId       	string     `gorm:"type:varchar(100);null;"`
	IaasUserPw       	string     `gorm:"type:varchar(100);null;"`
	PaasUserId       	string     `gorm:"type:varchar(100);null;"`
	PaasUserPw       	string     `gorm:"type:varchar(100);null;"`
	IaasUserUseYn    string     `gorm:"type:varchar(1);null;"`
	PaasUserUseYn    string     `gorm:"type:varchar(1);null;"`
	UpdatedAt       	time.Time  `gorm:"type:datetime;null;DEFAULT:null"`
	CreatedAt       	time.Time  `gorm:"type:datetime;null;DEFAULT:CURRENT_TIMESTAMP"`
}


func main() {

	sessionCookie, _ := utils.GenerateRandomString(32)
	model.SessionManager = *scs.NewCookieManager(sessionCookie) //("u46IpCV9y5Vlur8YvODJEhgOY8m9JVE4")
	model.SessionManager.Lifetime(time.Minute * 30)

	//model.SessionManager.Secure()
	//============================================
	// 기본적인 프로퍼티 설정 정보 읽어오기
	config, err := ReadConfig(`config.ini`)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	xmlFile, err := ReadXmlConfig(`log_config.xml`)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	logger, err := seelog.LoggerFromConfigAsBytes([]byte(xmlFile))

	if err != nil {
		fmt.Println(err)
		return
	}
	model.MonitLogger = logger
	UseLogger(logger)

	timeGap, _ := strconv.Atoi(config["gmt.time.gap"])
	model.GmtTimeGap = timeGap
	bm.GmtTimeGap = timeGap

	apiPort, _ := strconv.Atoi(config["server.port"])

	sysType := config["system.monitoring.type"]

	// iaas client
	var iaasDbAccessObj *gorm.DB
	var iaaSInfluxServerClient client.Client
	var iaasElasticClient *elastic.Client
	var openstackProvider model.OpenstackProvider
	var monClient *monascaclient.Client
	var auth gophercloud.AuthOptions

	// paas client
	var paaSInfluxServerClient client.Client
	var paasElasticClient *elastic.Client
	var databases bm.Databases
	var cfProvider cfclient.Config
	var boshClient *gogobosh.Client

	// Common MysqlDB
	paasConfigDbCon := new(DBConfig)
	paasConfigDbCon.DbType = config["paas.monitoring.db.type"]
	paasConfigDbCon.DbName = config["paas.monitoring.db.dbname"]
	paasConfigDbCon.UserName = config["paas.monitoring.db.username"]
	paasConfigDbCon.UserPassword = config["paas.monitoring.db.password"]
	paasConfigDbCon.Host = config["paas.monitoring.db.host"]
	paasConfigDbCon.Port = config["paas.monitoring.db.port"]

	paasConnectionString := utils.GetConnectionString(paasConfigDbCon.Host, paasConfigDbCon.Port, paasConfigDbCon.UserName, paasConfigDbCon.UserPassword, paasConfigDbCon.DbName)
	fmt.Println("String:", paasConnectionString)
	paasDbAccessObj, paasDbErr := gorm.Open(paasConfigDbCon.DbType, paasConnectionString+"?charset=utf8&parseTime=true")
	if paasDbErr != nil{
		fmt.Println("err::",paasDbErr)
		return
	}

	// memberInfo table (use common database table)
	createTable(paasDbAccessObj)

	// Redis Client
	rdClient := redis.NewClient(&redis.Options{
		Addr:     config["redis.addr"],
		Password: config["redis.password"],
	})

	// IaaS Connection Info
	if sysType == utils.SYS_TYPE_ALL || sysType == utils.SYS_TYPE_IAAS {
		iaasDbAccessObj,iaaSInfluxServerClient, iaasElasticClient, openstackProvider, monClient, auth, err = getIaasClients(config)
		if err != nil {
			log.Println(err)
			os.Exit(-1)
		}
	}


	if sysType == utils.SYS_TYPE_ALL || sysType == utils.SYS_TYPE_PAAS {
		paaSInfluxServerClient, paasElasticClient, databases, cfProvider, boshClient, err = getPaasClients(config)
		if err != nil {
			log.Println(err)
			os.Exit(-1)
		}
	}

	// Route Path 정보와 처리 서비스 연결
	var handler http.Handler
	if sysType == utils.SYS_TYPE_IAAS {
		handler = handlers.NewHandler(openstackProvider, iaaSInfluxServerClient, nil,
			iaasDbAccessObj,  paasDbAccessObj, iaasElasticClient, nil, *monClient, auth, bm.Databases{},
			cfclient.Config{}, rdClient, sysType, nil)
	}else if sysType == utils.SYS_TYPE_PAAS {
		handler = handlers.NewHandler(model.OpenstackProvider{}, nil, paaSInfluxServerClient,
			nil,  paasDbAccessObj, nil, paasElasticClient, monascaclient.Client{}, gophercloud.AuthOptions{}, databases,
			cfProvider, rdClient, sysType, boshClient)
	}else{
		handler = handlers.NewHandler(openstackProvider, iaaSInfluxServerClient, paaSInfluxServerClient,
			iaasDbAccessObj,  paasDbAccessObj, iaasElasticClient, paasElasticClient, *monClient, auth, databases,
			cfProvider, rdClient, sysType, boshClient)
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%v", apiPort), handler); err != nil {
		log.Fatalln(err)
	}

}


func UseLogger(newLogger seelog.LoggerInterface) {
	utils.Logger = newLogger
}

func ReadXmlConfig (filename string) (string, error) {
	xmlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "", err
	}
	return string(xmlFile),  nil
}

// Config 파일 읽어 오기
func ReadConfig(filename string) (Config, error) {
	// init with some bogus data
	config := Config{
		"server.ip":     "127.0.0.1",
		"server.port":   "8888",
	}

	if len(filename) == 0 {
		return config, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		// check if the line has = sign
		// and process the line. Ignore the rest.
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				// assign the config map
				config[key] = value
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}

func createTable(dbClient *gorm.DB){

	dbClient.Debug().AutoMigrate(&MemberInfo{})

	/*memberInfo := &MemberInfo {
		UserId       		:"admin",
		UserPw       		:"03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4",
		UserEmail       	:"amedio77@gmail.com",
		UserNm       		:"admin",
		IaasUserId       	:"admin",
		IaasUserPw       	:"cfmonit",
		PaasUserId       	:"admin",
		PaasUserPw       	:"admin",
		IaasUserUseYn       :"N",
		PaasUserUseYn       :"Y",
	}

	db.FirstOrCreate(memberInfo)*/
}

func getIaasClients(config Config) (iaasDbAccessObj *gorm.DB, iaaSInfluxServerClient client.Client, iaasElasticClient *elastic.Client, openstackProvider model.OpenstackProvider, monClient *monascaclient.Client, auth gophercloud.AuthOptions, err error){

	// Mysql
	iaasConfigDbCon := new(DBConfig)
	iaasConfigDbCon.DbType = config["iaas.monitoring.db.type"]
	iaasConfigDbCon.DbName = config["iaas.monitoring.db.dbname"]
	iaasConfigDbCon.UserName = config["iaas.monitoring.db.username"]
	iaasConfigDbCon.UserPassword = config["iaas.monitoring.db.password"]
	iaasConfigDbCon.Host = config["iaas.monitoring.db.host"]
	iaasConfigDbCon.Port = config["iaas.monitoring.db.port"]

	iaasConnectionString := utils.GetConnectionString(iaasConfigDbCon.Host, iaasConfigDbCon.Port, iaasConfigDbCon.UserName, iaasConfigDbCon.UserPassword, iaasConfigDbCon.DbName)
	fmt.Println("String:", iaasConnectionString)
	iaasDbAccessObj, err = gorm.Open(iaasConfigDbCon.DbType, iaasConnectionString+"?charset=utf8&parseTime=true")

	//Alarm 처리 내역 정보 Table 생성
	iaasDbAccessObj.Debug().AutoMigrate(&model.AlarmActionHistory{})


	// InfluxDB
	iaasUrl, _ := config["iaas.metric.db.url"]
	iaasUserName, _ := config["iaas.metric.db.username"]
	iaasPassword, _ := config["iaas.metric.db.password"]

	iaaSInfluxServerClient, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     iaasUrl,
		Username: iaasUserName,
		Password: iaasPassword,
	})

	// ElasticSearch
	iaasElasticUrl, _ := config["iaas.elastic.url"]
	iaasElasticClient, err = elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("http://%s", iaasElasticUrl)),
		elastic.SetSniff(false),
	)

	// Openstack Info
	openstackProvider.Region, _ 			= config["default.region"]
	openstackProvider.Username, _ 		    = config["default.username"]
	openstackProvider.Password, _ 			= config["default.password"]
	openstackProvider.Domain, _ 			= config["default.domain"]
	openstackProvider.TenantName, _ 		= config["default.tenant_name"]
	openstackProvider.AdminTenantId, _	 	= config["default.tenant_id"]
	openstackProvider.KeystoneUrl, _ 		= config["keystone.url"]
	openstackProvider.IdentityEndpoint, _ 	= config["identity.endpoint"]
	openstackProvider.RabbitmqUser, _ 		= config["rabbitmq.user"]
	openstackProvider.RabbitmqPass, _		= config["rabbitmq.pass"]
	openstackProvider.RabbitmqTargetNode, _	= config["rabbitmq.target.node"]

	model.MetricDBName, _ 		= config["iaas.metric.db.name"]
	model.NovaUrl, _ 			= config["nova.target.url"]
	model.NovaVersion, _ 		= config["nova.target.version"]
	model.NeutronUrl, _ 		= config["neutron.target.url"]
	model.NeutronVersion, _ 	= config["neutron.target.version"]
	model.KeystoneUrl, _ 		= config["keystone.target.url"]
	model.KeystoneVersion, _ 	= config["keystone.target.version"]
	model.CinderUrl, _ 			= config["cinder.target.url"]
	model.CinderVersion, _ 		= config["cinder.target.version"]
	model.GlanceUrl, _ 			= config["glance.target.url"]
	model.GlanceVersion,_ 		= config["glance.target.version"]
	model.DefaultTenantId, _	= config["default.tenant_id"]
	model.RabbitMqIp, _ 		= config["rabbitmq.ip"]
	model.RabbitMqPort, _ 	    = config["rabbitmq.port"]
	model.GMTTimeGap, _ 	    = strconv.ParseInt(config["gmt.time.gap"], 10, 64)

	monClient = monascaclient.New()
	monClient.SetBaseURL(config["monasca.url"])
	timeOut, _ := strconv.Atoi(config["monasca.connect.timeout"])
	monClient.SetTimeout(timeOut)

	tls, _ := strconv.ParseBool(config["monasca.secure.tls"])
	monClient.SetInsecure(tls)

	auth = gophercloud.AuthOptions{
		DomainName : config["default.domain"],
		IdentityEndpoint : config["keystone.url"],
		Username : config["default.username"],
		Password : config["default.password"],
		TenantID : config["default.tenant_id"],
	}

	return
}

func getPaasClients(config Config) (paaSInfluxServerClient client.Client, paasElasticClient *elastic.Client, databases bm.Databases,cfProvider cfclient.Config, boshClient *gogobosh.Client, err error ){

	// InfluxDB
	paasUrl, _ := config["paas.metric.db.url"]
	paasuserName, _ := config["paas.metric.db.username"]
	paasPassword, _ := config["paas.metric.db.password"]

	paaSInfluxServerClient, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     paasUrl,
		Username: paasuserName,
		Password: paasPassword,
	})

	// ElasticSearch
	paasElasticUrl, _ := config["paas.elastic.url"]
	paasElasticClient, err = elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("http://%s", paasElasticUrl)),
		elastic.SetSniff(false),
	)

	// PaaS Database
	bosh_database, _ := config["paas.metric.db.name.bosh"]
	paasta_database, _ := config["paas.metric.db.name.paasta"]
	container_database, _ := config["paas.metric.db.name.container"]


	databases.BoshDatabase = bosh_database
	databases.PaastaDatabase = paasta_database
	databases.ContainerDatabase = container_database


	// Cloud Foundry Client
	cfProvider = cfclient.Config{
		ApiAddress: config["paas.cf.client.apiaddress"],
		//Username:     "admin",
		//Password:     "admin",
		SkipSslValidation : true,
	}

	// BOSH Client Config
	boshConfig := &gogobosh.Config{
		BOSHAddress:       config["bosh.client.api.address"],
		Username:          config["bosh.client.api.username"],
		Password:          config["bosh.client.api.password"],
		HttpClient:        http.DefaultClient,
		SkipSslValidation: true,

	}
	boshClient, err = gogobosh.NewClient(boshConfig)
	if err != nil {
		log.Fatalln("Failed to create connection to the bosh server. err=", err)
	}

	return
}
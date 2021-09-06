package main

import (
	"crypto/tls"
	"fmt"

	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cihub/seelog"
	"github.com/cloudfoundry-community/gogobosh"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gophercloud/gophercloud"
	"github.com/influxdata/influxdb1-client/v2"
	"github.com/jinzhu/gorm"
	"github.com/monasca/golang-monascaclient/monascaclient"

	"kr/paasta/monitoring/common/config"
	commonModel "kr/paasta/monitoring/common/model"
	"kr/paasta/monitoring/handlers"
	"kr/paasta/monitoring/iaas/model"
	"kr/paasta/monitoring/iaas_new"
	bm "kr/paasta/monitoring/paas/model"
	"kr/paasta/monitoring/utils"
)

type Config map[string]string


type MemberInfo struct {
	UserId        string    `gorm:"type:varchar(50);primary_key"`
	UserPw        string    `gorm:"type:varchar(500);null;"`
	UserEmail     string    `gorm:"type:varchar(100);null;"`
	UserNm        string    `gorm:"type:varchar(100);null;"`
	IaasUserId    string    `gorm:"type:varchar(100);null;"`
	IaasUserPw    string    `gorm:"type:varchar(100);null;"`
	CaasUserId    string    `gorm:"type:varchar(100);null;"`
	CaasUserPw    string    `gorm:"type:varchar(100);null;"`
	PaasUserId    string    `gorm:"type:varchar(100);null;"`
	PaasUserPw    string    `gorm:"type:varchar(100);null;"`
	IaasUserUseYn string    `gorm:"type:varchar(1);null;"`
	PaasUserUseYn string    `gorm:"type:varchar(1);null;"`
	CaasUserUseYn string    `gorm:"type:varchar(1);null;"`
	UpdatedAt     time.Time `gorm:"type:datetime;null;DEFAULT:null"`
	CreatedAt     time.Time `gorm:"type:datetime;null;DEFAULT:CURRENT_TIMESTAMP"`
}

func main() {

	// 기본적인 프로퍼티 설정 정보 읽어오기
	configMap, err := config.ImportConfig(`config.ini`)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	xmlFile, err := config.ConvertXmlToString(`log_config.xml`)
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

	timeGap, _ := strconv.Atoi(configMap["gmt.time.gap"])
	model.GmtTimeGap = timeGap
	bm.GmtTimeGap = timeGap

	apiPort, _ := strconv.Atoi(configMap["server.port"])

	sysType := configMap["system.monitoring.type"]

	// iaas client
	var iaasDbAccessObj *gorm.DB
	var iaaSInfluxServerClient client.Client
	var iaasElasticClient *elasticsearch.Client
	var openstackProvider model.OpenstackProvider
	var monClient *monascaclient.Client
	var auth gophercloud.AuthOptions

	// paas client
	var paaSInfluxServerClient client.Client
	var paasElasticClient *elasticsearch.Client
	var databases bm.Databases
	//var cfProvider cfclient.Config
	var boshClient *gogobosh.Client

	// Common MysqlDB
	paasConfigDbCon := new(commonModel.DBConfig)
	paasConfigDbCon.DbType = configMap["paas.monitoring.db.type"]
	paasConfigDbCon.DbName = configMap["paas.monitoring.db.dbname"]
	paasConfigDbCon.UserName = configMap["paas.monitoring.db.username"]
	paasConfigDbCon.UserPassword = configMap["paas.monitoring.db.password"]
	paasConfigDbCon.Host = configMap["paas.monitoring.db.host"]
	paasConfigDbCon.Port = configMap["paas.monitoring.db.port"]

	paasConnectionString := utils.GetConnectionString(paasConfigDbCon.Host, paasConfigDbCon.Port, paasConfigDbCon.UserName, paasConfigDbCon.UserPassword, paasConfigDbCon.DbName)
	fmt.Println("String:", paasConnectionString)
	paasDbAccessObj, paasDbErr := gorm.Open(paasConfigDbCon.DbType, paasConnectionString+"?charset=utf8&parseTime=true")
	if paasDbErr != nil {
		fmt.Println("err::", paasDbErr)
		return
	}

	// 2021.09.06 - 왜 있는건지??
	// memberInfo table (use common database table)
	//createTable(paasDbAccessObj)

	// Redis Client
	rdClient := redis.NewClient(&redis.Options{
		Addr:     configMap["redis.addr"],
		Password: configMap["redis.password"],
	})
	cfConfig := bm.CFConfig{
		Host:           configMap["paas.monitoring.cf.host"],
		CaasBrokerHost: configMap["caas.monitoring.broker.host"],
	}
	//IaaS Connection Info
	if strings.Contains(sysType, utils.SYS_TYPE_ALL) || strings.Contains(sysType, utils.SYS_TYPE_IAAS) {
		iaasDbAccessObj, iaaSInfluxServerClient, iaasElasticClient, openstackProvider, monClient, auth, err = iaas_new.GetIaasClients(configMap)
		if err != nil {
			log.Println(err)
			os.Exit(-1)
		}
	}
	//
	if strings.Contains(sysType, utils.SYS_TYPE_ALL) || strings.Contains(sysType, utils.SYS_TYPE_PAAS) {
		fmt.Println("sysType == utils.SYS_TYPE_ALL || sysType == utils.SYS_TYPE_PAAS")
		paaSInfluxServerClient, paasElasticClient, databases, boshClient, err = getPaasClients(configMap)
		if err != nil {
			log.Println(err)
			os.Exit(-1)
		}
	}

	// Route Path 정보와 처리 서비스 연결
	var handler http.Handler

	if strings.Contains(sysType, utils.SYS_TYPE_ALL) || strings.Contains(sysType, utils.SYS_TYPE_IAAS) {
		handler = handlers.NewHandler(openstackProvider, iaaSInfluxServerClient, paaSInfluxServerClient,
			iaasDbAccessObj, paasDbAccessObj, iaasElasticClient, paasElasticClient, *monClient, auth, databases,
			rdClient, sysType, boshClient, cfConfig)
		if err := http.ListenAndServe(fmt.Sprintf(":%v", apiPort), handler); err != nil {
			log.Fatalln(err)
		}
	} else {
		handler = handlers.NewHandler(openstackProvider, iaaSInfluxServerClient, paaSInfluxServerClient,
			iaasDbAccessObj, paasDbAccessObj, iaasElasticClient, paasElasticClient, monascaclient.Client{}, auth, databases,
			rdClient, sysType, boshClient, cfConfig)
		if err := http.ListenAndServe(fmt.Sprintf(":%v", apiPort), handler); err != nil {
			log.Fatalln(err)
		}
	}

}

func UseLogger(newLogger seelog.LoggerInterface) {
	utils.Logger = newLogger
}





// 2021.09.06 - 이거 왜 있는지??
//func createTable(dbClient *gorm.DB) {
//	dbClient.Debug().AutoMigrate(&MemberInfo{})
//}



func getPaasClients(config map[string]string) (paaSInfluxServerClient client.Client, paasElasticClient *elasticsearch.Client, databases bm.Databases, boshClient *gogobosh.Client, err error) {

	// InfluxDB
	paasUrl, _ := config["paas.metric.db.url"]
	paasuserName, _ := config["paas.metric.db.username"]
	paasPassword, _ := config["paas.metric.db.password"]

	paaSInfluxServerClient, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     paasUrl,
		Username: paasuserName,
		Password: paasPassword,
		InsecureSkipVerify: true,
	})

	fmt.Printf("paaSInfluxServerClient : %v\n", paaSInfluxServerClient)
	fmt.Printf("err : %v\n", err)

	elasticsearchUsername, _ := config["paas.elasticsearch.username"]
	elasticsearchPassword, _ := config["paas.elasticsearch.password"]
	elasticsearchUrl, _ := config["paas.elasticsearch.url"]
	elasticsearchHttpsEnabled, _ := strconv.ParseBool(config["paas.elasticsearch.https_enabled"])

	cfg := elasticsearch.Config{
		Username: elasticsearchUsername,
		Password: elasticsearchPassword,
		Addresses: []string{
			elasticsearchUrl,
		},
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MaxVersion:         tls.VersionTLS11,
				InsecureSkipVerify: elasticsearchHttpsEnabled,
			},
		},
	}
	paasElasticClient, err = elasticsearch.NewClient(cfg)
	fmt.Println("paasElasticClient::", paasElasticClient)
	fmt.Println("err::", err)

	// ElasticSearch
	/*paasElasticUrl, _ := config["paas.elastic.url"]
	paasElasticClient, err = elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("http://%s", paasElasticUrl)),
		elastic.SetSniff(false),
	)*/

	// PaaS Database
	bosh_database, _ := config["paas.metric.db.name.bosh"]
	paasta_database, _ := config["paas.metric.db.name.paasta"]
	container_database, _ := config["paas.metric.db.name.container"]

	databases.BoshDatabase = bosh_database
	databases.PaastaDatabase = paasta_database
	databases.ContainerDatabase = container_database

	// Cloud Foundry Client
	//cfProvider = cfclient.Config{
	//	ApiAddress: config["paas.cf.client.apiaddress"],
	//	//Username:     "admin",
	//	//Password:     "admin",
	//	SkipSslValidation: true,
	//}

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

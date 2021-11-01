package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/rday/zabbix"
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
	"github.com/influxdata/influxdb1-client/v2"
	"github.com/jinzhu/gorm"

	"kr/paasta/monitoring/common/config"
	"kr/paasta/monitoring/handlers"
	"kr/paasta/monitoring/iaas_new"
	"kr/paasta/monitoring/iaas_new/model"
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

var logger seelog.LoggerInterface

func main() {

	xmlFile, err := config.ConvertXmlToString("log_config.xml")
	if err != nil {
		os.Exit(-1)
	}

	logger, _ = seelog.LoggerFromConfigAsBytes([]byte(xmlFile))
	model.MonitLogger = logger
	utils.Logger = logger

	// 기본적인 프로퍼티 설정 정보 읽어오기
	configMap, err := config.ImportConfig("config.ini")
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}

	timeGap, _ := strconv.Atoi(configMap["gmt.time.gap"])
	model.GmtTimeGap = timeGap
	bm.GmtTimeGap = timeGap

	apiPort, _ := strconv.Atoi(configMap["server.port"])

	sysType := configMap["system.monitoring.type"]

	// paas client
	var paaSInfluxServerClient client.Client
	var paasElasticClient *elasticsearch.Client
	var databases bm.Databases
	var boshClient *gogobosh.Client

	// config.ini 파일에서 MySQL 접속정보를 추출
	dbConnInfo := config.InitDBConnectionConfig(configMap)

	paasConnectionString := utils.GetConnectionString(dbConnInfo.Host, dbConnInfo.Port, dbConnInfo.UserName, dbConnInfo.UserPassword, dbConnInfo.DbName)
	logger.Infof("DB Connection Info : %v\n", paasConnectionString)

	paasDbAccessObj, paasDbErr := gorm.Open(dbConnInfo.DbType, paasConnectionString+"?charset=utf8&parseTime=true")
	if paasDbErr != nil {
		logger.Errorf("%v\n", paasDbErr)
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
	logger.Info(rdClient)

	cfConfig := bm.CFConfig{
		Host:           configMap["paas.monitoring.cf.host"],
		CaasBrokerHost: configMap["caas.monitoring.broker.host"],
	}
	//IaaS Connection Info
	iaasClient := iaas_new.IaasClient{}

	if strings.Contains(sysType, utils.SYS_TYPE_ALL) || strings.Contains(sysType, utils.SYS_TYPE_IAAS) {
		iaasClient, err = iaas_new.GetIaasClients(configMap)
		if err != nil {
			logger.Error(err)
			os.Exit(-1)
		}
	}

	if strings.Contains(sysType, utils.SYS_TYPE_ALL) || strings.Contains(sysType, utils.SYS_TYPE_PAAS) {
		paaSInfluxServerClient, paasElasticClient, databases, boshClient, err = getPaasClients(configMap)
		if err != nil {
			logger.Error(err)
			os.Exit(-1)
		}
	}

	// Route Path 정보와 처리 서비스 연결
	var handler http.Handler

	if strings.Contains(sysType, utils.SYS_TYPE_ALL) || strings.Contains(sysType, utils.SYS_TYPE_IAAS) {
		handler = handlers.NewHandler(iaasClient.Provider, iaasClient.InfluxClient, paaSInfluxServerClient,
			iaasClient.ConnectionPool, paasDbAccessObj, iaasClient.ElasticClient, paasElasticClient, iaasClient.AuthOpts, databases,
			rdClient, sysType, boshClient, cfConfig)
		if err := http.ListenAndServe(fmt.Sprintf(":%v", apiPort), handler); err != nil {
			logger.Error(err)
		}
	} else {
		handler = handlers.NewHandler(iaasClient.Provider, iaasClient.InfluxClient, paaSInfluxServerClient,
			iaasClient.ConnectionPool, paasDbAccessObj, iaasClient.ElasticClient, paasElasticClient, iaasClient.AuthOpts, databases,
			rdClient, sysType, boshClient, cfConfig)
		if err := http.ListenAndServe(fmt.Sprintf(":%v", apiPort), handler); err != nil {
			logger.Error(err)
		}
	}
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

	logger.Infof("paaSInfluxServerClient : %v\n", paaSInfluxServerClient)
	if err != nil {
		logger.Errorf("err : %v\n", err)
	}


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
	logger.Infof("paasElasticClient : %v\n", paasElasticClient)
	if err != nil {
		logger.Errorf("err : %v\n", err)
	}


	// ElasticSearch
	/*paasElasticUrl, _ := config["paas.elastic.url"]
	paasElasticClient, err = elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("http://%s", paasElasticUrl)),
		elastic.SetSniff(false),
	)*/

	// PaaS Database
	boshDatabase, _ := config["paas.metric.db.name.bosh"]
	paastaDatabase, _ := config["paas.metric.db.name.paasta"]
	containerDatabase, _ := config["paas.metric.db.name.container"]

	databases.BoshDatabase = boshDatabase
	databases.PaastaDatabase = paastaDatabase
	databases.ContainerDatabase = containerDatabase

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
		logger.Errorf("Failed to create connection to the bosh server. err=", err)
	}

	// Zabbix API에 연결하기
	// Zabbix API 객체 생성
	api, err := zabbix.NewAPI("http://203.255.255.101:8080/zabbix/api_jsonrpc.php", "Admin", "zabbix")
	if err != nil {
		fmt.Println(err)
		return
	}


	// Zabbix API의 버전 정보 가져오기
	ApiVersionInfo, err := api.Version()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print("*Zabbix API Version: ")
	fmt.Println(ApiVersionInfo)


	// Zabbix API에 로그인 하기
	// 로그인 성공 여부 확인
	_, err = api.Login()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("*Login Successful")


	// Zabbix 토큰 값 가져오기
	ApiToken := api.GetAuth()
	fmt.Print("*Zabbix API Token: ")
	fmt.Println(ApiToken)


	// <<< *** hostgroup.get 메서드 사용하는 영역 *** >>>
	// Zabbix에 등록된 HostGroup 전체 리스트 가져오기
	paramsHostGroup := make(map[string]interface{})
	hostGroup, err := api.HostGroup("get", paramsHostGroup)
	if err != nil {
		fmt.Println(err)
		return
	}
	hostGroupResultForJson, _ := json.MarshalIndent(hostGroup, "", "  ")
	fmt.Println(string(hostGroupResultForJson))


	// Zabbix HostGroup 중 "PaaS-TA Group"에 속한 호스트 갯수 가져오기
	// "name"이 "PaaS-TA Group"인 호스트를 검색하기 위해 "search" 파라미터 사용
	// 갯수를 가져오기 위해 "selectHosts" 파라미터의 "count" 속성 필요
	hostGroupName := make(map[string]interface{})
	hostGroupName["name"] = "PaaS-TA Group"
	paramsHostGroup["search"] = hostGroupName
	paramsHostGroup["selectHosts"] = "count"
	hostGroup, err = api.HostGroup("get", paramsHostGroup)
	if err != nil {
		fmt.Println(err)
		return
	}
	hostGroupResultForJson, _ = json.MarshalIndent(hostGroup, "", "  ")
	fmt.Println(string(hostGroupResultForJson))


	// Zabbix HostGroup 중 "PaaS-TA Group"에 속한 호스트 전체 리스트 가져오기
	// 호스트의 "name" 리스트를 가져오기 위해 "selectHosts" 파라미터가 사용할 수 있는 "hosts"의 속성 배열 값 중 "name"을 사용
	// 따라서 hosts 속성 사용 시에는 string 배열로 사용되어야 함
	hostProp := []string{"name"}
	paramsHostGroup["selectHosts"] = hostProp
	hostGroup, err = api.HostGroup("get", paramsHostGroup)
	if err != nil {
		fmt.Println(err)
		return
	}
	hostGroupResultForJson, _ = json.MarshalIndent(hostGroup, "", "  ")
	fmt.Println(string(hostGroupResultForJson))


	// <<< *** item.get 메서드 사용하는 영역 *** >>>
	// 특정 호스트에 대한 기본 시스템 정보(CPU/Memory/Disk Utilization, Interface Address) 가져오기
	paramsItem := make(map[string]interface{})
	itemFilter := make(map[string]interface{})
	nameList := []string{}
	outputList := []string{}
	interfaceProp := []string{}
	paramsItem["group"] = "PaaS-TA Group"
	paramsItem["host"] = "ebcbef8b-cf4d-409d-ab58-c7ee352b6604"
	nameList = []string{"CPU utilization", "Memory utilization", "/: Space utilization"}
	itemFilter["name"] = nameList
	paramsItem["filter"] = itemFilter
	outputList = []string{"name", "lastvalue", "units"}
	paramsItem["output"] = outputList
	interfaceProp = []string{"ip"}
	paramsItem["selectInterfaces"] = interfaceProp
	itemInfo, err := api.Item("get", paramsItem)
	if err != nil {
		fmt.Println(err)
		return
	}
	itemInfoForJson, _ := json.MarshalIndent(itemInfo, "", "  ")
	fmt.Println(string(itemInfoForJson))

	return
}

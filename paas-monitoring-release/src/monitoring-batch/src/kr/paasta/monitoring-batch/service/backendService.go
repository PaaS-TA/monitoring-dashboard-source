package service

import (
	"os"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	client "github.com/influxdata/influxdb/client/v2"
	"github.com/cloudfoundry-community/gogobosh"
	"net/http"
	//"github.com/cloudfoundry-community/go-cfclient"
	"kr/paasta/monitoring-batch/util"
	"sync"
	"github.com/go-redis/redis"
	"github.com/cloudfoundry-community/go-cfclient"
)

type BackendServices struct {
	GmtTimeGapHour   int64
	Influxclient     client.Client
	InfluxConfig     *InfluxConfig
	MonitoringDbClient         *gorm.DB
	//PortalDbClient   *gorm.DB
	BoshClient       *gogobosh.Client
	BoshConfig       *BoshConfig
	RedisClient	 	 *redis.Client
	RedisConfig	 	 *RedisConfig
	MailConfig       *MailConfig
	CfClient		 *cfclient.Client
	//ThresholdConfig *ThresholdConfig
	retry            bool
	StopChan         chan bool
	config			map[string]string
}

type InfluxConfig struct {
	InfluxUrl		string
	InfluxUser 		string
	InfluxPass 		string
	InfraDatabase  	 	string
	PaastaDatabase   	string
	ContainerDatabase   	string
	DefaultTimeRange        string
}


type DBConfig struct {
	DbType string
	UserName string
	UserPassword string
	Host string
	Port string
	DbName string
}


type BoshConfig struct {
	BoshUrl  string
	BoshIp   string
	BoshId   string
	BoshPass string
	CfDeploymentName    string
	DiegoDeploymentName string
	CellNamePrefix      string
	ServiceName         string
}

type RedisConfig struct {
	RedisAddr  		string
	RedisPassword   string
}

type ThresholdConfig struct {
	BoshCpuCritical          int
	BoshCpuWarning           int
	BoshMemoryCritical       int
	BoshMemoryWarning        int
	BoshDiskCritical         int
	BoshDiskWarning          int

	PaasTaCpuCritical        int
	PaasTaCpuWarning         int
	PaasTaMemoryCritical     int
	PaasTaMemoryWarning      int
	PaasTaDiskCritical       int
	PaasTaDiskWarning        int

	ContainerCpuCritical     int
	ContainerCpuWarning      int
	ContainerMemoryCritical  int
	ContainerMemoryWarning   int
	ContainerDiskCritical    int
	ContainerDiskWarning     int

	AlarmResendInterval      int
}


type MailConfig struct {
	SmtpHost       string
	Port           string
	MailSender     string
	SenderPwd      string
	ResouceUrl     string
	MailReceiver   string
	AlarmSend      bool
	MailTlsSend    bool
}

func NewBackendServices(gmtTimeGapHour int64, influx *InfluxConfig, configDB *DBConfig, portalDB *DBConfig, boshConfig *BoshConfig,
	redisConfig *RedisConfig, mailConfig *MailConfig, config map[string]string ) *BackendServices {

	connectionString := util.GetConnectionString(configDB.Host , configDB.Port, configDB.UserName, configDB.UserPassword, configDB.DbName )
	fmt.Println("String:",connectionString)

	dbAccessObj, err := gorm.Open(configDB.DbType, connectionString + "?charset=utf8&parseTime=true")

	//portalConnectionString := util.GetConnectionString(portalDB.Host , portalDB.Port, portalDB.UserName, portalDB.UserPassword, portalDB.DbName )
	//portalDbAccessObj, err := gorm.Open(portalDB.DbType, portalConnectionString  + "?charset=utf8&parseTime=true")
	//fmt.Println("portalCon:", portalConnectionString)

	CreateTable(dbAccessObj)
	CreateAlarmPolicyInitialData(dbAccessObj)

	InfluxServerClient, _ := client.NewHTTPClient(client.HTTPConfig{
		Addr: influx.InfluxUrl,
		Username: influx.InfluxUser,
		Password: influx.InfluxPass,
	})

	fmt.Println("xxxxxxx2", InfluxServerClient)
	boshClientConfig := &gogobosh.Config{
		BOSHAddress: 		fmt.Sprintf("https://%s", boshConfig.BoshUrl),
		Username:    		boshConfig.BoshId,
		Password:    		boshConfig.BoshPass,
		HttpClient:        	http.DefaultClient,
		SkipSslValidation: 	true,
	}
	boshClient, err := gogobosh.NewClient(boshClientConfig)
	if err != nil {
		fmt.Println("Error:", err.Error())
		panic( err )
	}
	fmt.Println("xxxxxxx3")
	redisOptions := redis.Options{
		Addr: redisConfig.RedisAddr,
		Password: redisConfig.RedisPassword,
		//DB:       0,  // use default DB
		//Addr:     "localhost:6379",
		//Password: "", // no password set
		//DB:       0,  // use default DB
		//DialTimeout:  10 * time.Second,
		//ReadTimeout:  30 * time.Second,
		//WriteTimeout: 30 * time.Second,
		//PoolSize:     10,
		//PoolTimeout:  30 * time.Second,
	}
	redisClient := redis.NewClient(&redisOptions)
	fmt.Println("xxxxxxx4")
	//cf-client
	cfclientConfig := cfclient.Config{
		ApiAddress: config["cf.client.api_address"],
		Username: config["cf.client.username"],
		Password: config["cf.client.password"],
		SkipSslValidation: true,
	}
	fmt.Println("xxxxxxx5")
	cfClient, cfErr:= cfclient.NewClient(&cfclientConfig)
	if cfErr != nil {
		fmt.Println("ssssssssssssss", cfErr)
		fmt.Errorf(">>>>> cfclient error:%v", cfErr)
		panic(cfErr)
	}

	stop := make(chan bool)
	fmt.Println("xxxxxxx6")
	return &BackendServices{
		GmtTimeGapHour: gmtTimeGapHour,
		Influxclient:	InfluxServerClient,	//MetricDB
		InfluxConfig:   influx,           	//InfluxDB 설정정보
		MonitoringDbClient:       dbAccessObj,		//Monitoring Configuration DB
		//PortalDbClient: portalDbAccessObj,
		BoshClient:     boshClient,
		BoshConfig:     boshConfig,
		RedisClient:	redisClient,
		RedisConfig:	redisConfig,
		CfClient:		cfClient,
		retry:		false,
		MailConfig:      mailConfig,
		StopChan: 	stop,
		config:			config,
	}
}


func (f *BackendServices) StopProcess(){
	go func(){
		for{
			select{
			case <- f.StopChan:
				os.Exit(1)	//Bosh Monit start Batch program automatically if the process is down.
			}
		}
	}()
}


func (f *BackendServices) CreateUpdateBoshData(boshConfig BoshConfig) error {

	err := C()
	if err != nil{
		f.StopChan <- true
	}

	CreteUpdateBoshVms(f, boshConfig, f.MonitoringDbClient)
	fmt.Println("# Bosh function end ...")

	return nil
}

func (f *BackendServices) UserPortalService() error {

	err := C()
	if err != nil{
		f.StopChan <- true
	}

	var wg sync.WaitGroup
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func(wg *sync.WaitGroup, index int) {
			defer wg.Done()
			switch index {
			case 0:
				fmt.Println(">>>>> START - AutoScaler")
				AutoScaler(f).AutoScale()
			case 1:
				fmt.Println(">>>>> START - PortalAppAlarm")
				PortalAppAlarm(f).PortalAppAlarmCollect()
			}
		}(&wg, i)
	}
	wg.Wait()

	fmt.Println(">>>>> END - UserPortalService")

	return nil
}

func (f *BackendServices) CreateAlarmData() error {

	err := C()
	if err != nil{
		f.StopChan <- true
	}

	//Update SNS Alarm target.
	UpdateSnsAlarmTarget(f)

	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func(wg *sync.WaitGroup, index int) {
			switch index {
			case 0 :
				fmt.Println("# BoshAlarmCollect start !")
				BoshAlarmCollect(f)
			case 1 :
				fmt.Println("# PaasTaAlarmCollect start !")
				PaasTaAlarmCollect(f)
			case 2 :
				fmt.Println("# ContainerAlarmCollect start !")
				ContainerAlarmCollect(f)
			}
			wg.Done()
		}(&wg, i)
	}
	wg.Wait()

	fmt.Println("# CreateAlarmData function end ...")
	return nil
}


func C () error{
	return nil
}

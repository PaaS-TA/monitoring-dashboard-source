package services

import (
	"os"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	client "github.com/influxdata/influxdb/client/v2"
	"github.com/cloudfoundry-community/gogobosh"
	"net/http"
	//"github.com/cloudfoundry-community/go-cfclient"
	"kr/paasta/monitoring/monit-batch/util"
	/*"kr/paasta/monitoring/monit-batch/dao"
	cb "kr/paasta/monitoring/monit-batch/models/base"*/
	"sync"
)

type BackendServices struct {
	GmtTimeGapHour   int64
	Influxclient     client.Client
	InfluxConfig     *InfluxConfig
	MonitoringDbClient         *gorm.DB
	PortalDbClient   *gorm.DB
	BoshClient       *gogobosh.Client
	BoshConfig       *BoshConfig
	MailConfig       *MailConfig
	//ThresholdConfig *ThresholdConfig
	retry            bool
	StopChan         chan bool
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
}


func NewBackendServices(gmtTimeGapHour int64, influx *InfluxConfig, configDB *DBConfig, portalDB *DBConfig,boshConfig *BoshConfig,
			mailConfig *MailConfig , config map[string]string ) *BackendServices{

	connectionString := util.GetConnectionString(configDB.Host , configDB.Port, configDB.UserName, configDB.UserPassword, configDB.DbName )
	fmt.Println("String:",connectionString)

	dbAccessObj, err := gorm.Open(configDB.DbType, connectionString + "?charset=utf8&parseTime=true")

	portalConnectionString := util.GetConnectionString(portalDB.Host , portalDB.Port, portalDB.UserName, portalDB.UserPassword, portalDB.DbName )

	portalDbAccessObj, err := gorm.Open(portalDB.DbType, portalConnectionString  + "?charset=utf8&parseTime=true")

	CreateTable(dbAccessObj)
	CreateAlarmPolicyInitialData(dbAccessObj)

	//CreateTablePortal(portalDbAccessObj)
	//CreatePortalInitialData(portalDbAccessObj)
	InfluxServerClient, _ := client.NewHTTPClient(client.HTTPConfig{
		Addr: influx.InfluxUrl,
		Username: influx.InfluxUser,
		Password: influx.InfluxPass,
	})


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

	stop := make(chan bool)

	return &BackendServices{
		GmtTimeGapHour: gmtTimeGapHour,
		Influxclient:	InfluxServerClient,	//MetricDB
		InfluxConfig:   influx,           	//InfluxDB 설정정보
		MonitoringDbClient:       dbAccessObj,		//Monitoring Configuration DB
		PortalDbClient: portalDbAccessObj,
		BoshClient:     boshClient,
		BoshConfig:     boshConfig,
		//ThresholdConfig: thresholdConfig,
		retry:		false,
		MailConfig:      mailConfig,
		StopChan: 	stop,
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

func (f *BackendServices) AutoScale() error {
	err := C()
	if err != nil{
		f.StopChan <- true
	}

	PortalAutoScale(f)
	return nil
}

func (f *BackendServices) CreateAlarmData() error {

	err := C()
	if err != nil{
		f.StopChan <- true
	}

	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func(wg *sync.WaitGroup, index int) {
			switch index {
				case 0 :
					BoshAlarmCollect(f)
				case 1 :
					PaasTaAlarmCollect(f)
				case 2 :
					ContainerAlarmCollect(f)
			}
			wg.Done()
		}(&wg, i)
	}
	wg.Wait()

	fmt.Println("# CreateAlarmData function end ...")
	return nil
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

func C () error{
	return nil
}


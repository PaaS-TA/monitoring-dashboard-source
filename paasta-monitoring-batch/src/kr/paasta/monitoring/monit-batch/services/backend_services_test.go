package services_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"kr/paasta/monitoring/monit-batch/util"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"kr/paasta/monitoring/monit-batch/services"
	"kr/paasta/monitoring/monit-batch/models"
	"kr/paasta/monitoring/monit-batch/dao"
	"os"
	"bufio"
	"strings"
	"io"
)

type DBConfig struct {
	DbType string
	UserName string
	UserPassword string
	Host string
	Port string
	DbName string
}

//#Top-level 단위테스트 묶음.
var _ = Describe("backendService", func() {

	var (				//#단위테스트 수행을 위해 필요한 변수 선언
	fakeDbAccessObj *gorm.DB
	runDone chan struct{}
	stopChan chan bool
	)

	//#각각의 단위테스트 수행 전 실행되는 함수
	BeforeEach(func() {

		config, _ := ReadConfig(`../config_test.ini`)
		_, configDbCon, _, _, _, config := GetObject(config)

		fmt.Println("########### Alarm Test Start #########")
		fakeDatabase := new(DBConfig)
		fakeDatabase.DbType = configDbCon.DbType
		fakeDatabase.DbName = configDbCon.DbName
		fakeDatabase.Host = configDbCon.Host
		fakeDatabase.Port = configDbCon.Port
		fakeDatabase.UserName = configDbCon.UserName
		fakeDatabase.UserPassword = configDbCon.UserPassword

		connectionString := util.GetConnectionString(fakeDatabase.Host , fakeDatabase.Port, fakeDatabase.UserName, fakeDatabase.UserPassword, fakeDatabase.DbName )

		//fakeLogger.Debug("Database ConnectString ::", connectionString)
		dbAccessObj, err := gorm.Open(fakeDatabase.DbType, connectionString + "?charset=utf8&parseTime=true")

		if err != nil{
			fmt.Errorf("Error===>",err)
		}

		fakeDbAccessObj = dbAccessObj

		//backend_service := services.NewBackendServices( n.influxCon, n.configDbCon, n.portalDbCon, n.boshCon,  n.mailConfig, n.thresholdConfig, n.config)

		stopChan = make(chan bool)
		runDone = make(chan struct{})
	})

	//#각각의 단위테스트 수행 후 실행되는 함수
	AfterEach(func() {
		close(stopChan)
		//Eventually(runDone).Should(BeClosed())
	})

	It("Get BackendService Handler", func() {
		config, _ := ReadConfig(`../config_test.ini`)
		influxCon, configDbCon, portalDbCon, boshCon, mailConfig, config := GetObject(config)

		backendService := services.NewBackendServices( -9, influxCon, configDbCon, portalDbCon, boshCon,  mailConfig, config)
		Eventually(func() *services.BackendServices {
			return backendService
		}).ShouldNot(Equal(nil))
	})

	It("Create AlarmData", func() {
		config, _ := ReadConfig(`../config_test.ini`)
		influxCon, configDbCon, portalDbCon, boshCon, mailConfig, config := GetObject(config)

		backendService := services.NewBackendServices( -9, influxCon, configDbCon, portalDbCon, boshCon,  mailConfig, config)
		err := backendService.CreateAlarmData()
		Eventually(func() error {
			return err
		}).Should(BeNil())
	})

	//AutoScale Test
	It("Auto Scale", func() {
		config, _ := ReadConfig(`../config_test.ini`)
		influxCon, configDbCon, portalDbCon, boshCon, mailConfig, config := GetObject(config)

		backendService := services.NewBackendServices( -9, influxCon, configDbCon, portalDbCon, boshCon,  mailConfig, config)
		backendService.AutoScale()


	})

	//Bosh Vm 정보 동기화 Test
	It("Bosh Info Sync", func() {
		config, _ := ReadConfig(`../config_test.ini`)
		influxCon, configDbCon, portalDbCon, boshCon, mailConfig, config := GetObject(config)

		backendService := services.NewBackendServices( -9, influxCon, configDbCon, portalDbCon, boshCon,  mailConfig, config)

		//Vm, Zone정보 초기화
		backendService.MonitoringDbClient.Where("").Delete(models.Vm{})
		backendService.MonitoringDbClient.Where("").Delete(models.Zone{})


		backendService.CreateUpdateBoshData(*boshCon)

		//Bosh VM정보가 정상 수집되었는지 체크
		Eventually(func() bool {
			zoneList, _ := dao.GetBoshVmsDao(backendService.BoshClient).GetZoneInfosList(backendService.MonitoringDbClient)
			vmList, _ := dao.GetBoshVmsDao(backendService.BoshClient).GetJobInfoList(backendService.MonitoringDbClient)
			if len(zoneList) > 0 && len(vmList) > 0{
				return true
			}
			return false
		}).Should(BeTrue())

	})

	It("Stop Process", func() {
		config, _ := ReadConfig(`../config_test.ini`)
		influxCon, configDbCon, portalDbCon, boshCon, mailConfig, config := GetObject(config)

		backendService := services.NewBackendServices( -9, influxCon, configDbCon, portalDbCon, boshCon,  mailConfig, config)
		backendService.StopProcess()
	})

})

type Config map[string]string

func GetObject(config Config) (*services.InfluxConfig, *services.DBConfig, *services.DBConfig, *services.BoshConfig, *services.MailConfig, map[string]string){

	//==============================================
	//Influx Configuration
	influxCon := new(services.InfluxConfig)
	influxCon.InfluxUrl  = config["influx.url"]
	influxCon.InfluxUser = config["influx.user"]
	influxCon.InfluxPass = config["influx.pass"]
	influxCon.InfraDatabase     = config["influx.bosh.db_name"]
	influxCon.PaastaDatabase    = config["influx.paasta.db_name"]
	influxCon.ContainerDatabase = config["influx.container.db_name"]
	influxCon.DefaultTimeRange = config["influx.defaultTimeRange"]
	//==============================================

	//==============================================
	//Monitoring configDB Configuration
	configDbCon := new(services.DBConfig)
	configDbCon.DbType = config["monitoring.db.type"]
	configDbCon.DbName = config["monitoring.db.dbname"]
	configDbCon.UserName      = config["monitoring.db.username"]
	configDbCon.UserPassword  = config["monitoring.db.password"]
	configDbCon.Host          = config["monitoring.db.host"]
	configDbCon.Port          = config["monitoring.db.port"]
	//==============================================

	//==============================================
	//Monitoring configDB Configuration
	portalDbCon := new(services.DBConfig)
	portalDbCon.DbType = config["portal.db.type"]
	portalDbCon.DbName = config["portal.db.dbname"]
	portalDbCon.UserName      = config["portal.db.username"]
	portalDbCon.UserPassword  = config["portal.db.password"]
	portalDbCon.Host          = config["portal.db.host"]
	portalDbCon.Port          = config["portal.db.port"]
	//==============================================

	fmt.Println("portalDbCon.Host====",portalDbCon.Host)

	//==============================================
	//configDB Configuration
	boshCon := new(services.BoshConfig)
	boshCon.BoshUrl  = config["bosh.api.url"]
	boshCon.BoshIp  = config["bosh.ip"]
	boshCon.BoshId   = config["bosh.admin"]
	boshCon.BoshPass = config["bosh.password"]
	boshCon.CfDeploymentName    = config["bosh.cf.deployment.name"]
	boshCon.DiegoDeploymentName = config["bosh.diego.deployment.name"]
	boshCon.CellNamePrefix      = config["bosh.diego.cell.name.prefix"]
	boshCon.ServiceName         = config["bosh.service.name"]
	//==============================================


	mailConfig := new(services.MailConfig)
	mailConfig.SmtpHost   = config["mail.smtp.host"]
	mailConfig.Port       = config["mail.smtp.port"]
	mailConfig.MailSender = config["mail.sender"]
	mailConfig.SenderPwd  = config["mail.sender.password"]
	mailConfig.ResouceUrl = config["mail.resource.url"]
	mailConfig.MailReceiver = config["mail.receiver.email"]

	return influxCon, configDbCon, portalDbCon, boshCon, mailConfig, config
}


func ReadConfig(filename string) (Config, error) {
	// init with some bogus data
	config := Config{
		"server.port": "9999",
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
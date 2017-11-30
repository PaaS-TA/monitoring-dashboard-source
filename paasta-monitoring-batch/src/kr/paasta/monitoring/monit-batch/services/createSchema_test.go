package services_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	_ "github.com/go-sql-driver/mysql"
	"kr/paasta/monitoring/monit-batch/util"
	"kr/paasta/monitoring/monit-batch/services"
	"kr/paasta/monitoring/monit-batch/models"
	"fmt"
	"github.com/jinzhu/gorm"
)

//#Top-level 단위테스트 묶음.
var _ = Describe("Bosh Service Test", func() {

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

	It("Create Table and InitalData", func() {

		services.CreateTable(fakeDbAccessObj)
		//services.CreateTablePortal(fakeDbAccessObj)
		services.CreatePortalInitialData(fakeDbAccessObj)
		alarmData1 := models.Alarm{Id: 1, OriginType:"bos", OriginId: 9999, AlarmType: "cpu", Level: "critical",
			Ip:"127.0.0.1", AppYn: "N", AlarmTitle: "alarm Occur", AlarmMessage: "alarm Message", ResolveStatus: "1", AlarmCnt: 1, CompleteUser: "test"}
		fakeDbAccessObj.Debug().FirstOrCreate(&alarmData1)

		var alarmList []models.Alarm

		fakeDbAccessObj.Debug().Model("alarms").Find(&alarmList)
		Eventually(func() int {
			return len(alarmList)
		}).ShouldNot(Equal(0))

		var alarm models.Alarm
		fakeDbAccessObj.Delete(&alarm)
	})

})


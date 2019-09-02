package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/tidwall/gjson"
	"kr/paasta/monitoring/saas/model"
	"kr/paasta/monitoring/saas/util"
	"strconv"

	//"github.com/thoas/go-funk"
	"log"
	"os"
)

var dbType string
var connectionString string

type alarmId struct {
	AlarmId string
}

func init() {
	// 기본적인 프로퍼티 설정 정보 읽어오기
	config, err := util.ReadConfig(`config.ini`)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	dbType = config["paas.monitoring.db.type"]
	dbName := config["paas.monitoring.db.dbname"]
	userName := config["paas.monitoring.db.username"]
	userPassword := config["paas.monitoring.db.password"]
	host := config["paas.monitoring.db.host"]
	port := config["paas.monitoring.db.port"]

	connectionString = fmt.Sprintf("%s:%s@%s([%s]:%s)/%s%s", userName, userPassword, "tcp", host, port, dbName, "")
}

func GetdbAccessObj() *gorm.DB {
	dbAccessObj, paasDbErr := gorm.Open(dbType, connectionString+"?charset=utf8&parseTime=true")
	if paasDbErr != nil {
		fmt.Println("err::", paasDbErr)
		return nil
	}
	return dbAccessObj
}

func CreateTable(dbClient *gorm.DB) {
	dbClient.Debug().AutoMigrate(&model.BatchAlarmInfo{}, &model.BatchAlarmExecution{}, &model.BatchAlarmReceiver{})
}

// Alarm Info
func GetBatchAlarmInfo(dbClient *gorm.DB) ([]model.BatchAlarmInfo, model.ErrMessage) {
	var alarmInfos []model.BatchAlarmInfo
	//status := dbClient.Debug().Find(&alarmInfos)
	status := dbClient.Debug().Table("batch_alarm_infos").Where("service_type = 'SaaS'").Find(&alarmInfos)

	err := util.GetError().DbCheckError(status.Error)
	if err != nil {
		return nil, err
	}

	return alarmInfos, nil
}

//// 알람 수신자 조회
func GetBatchAlarmReceiver(dbClient *gorm.DB) ([]model.BatchAlarmReceiver, model.ErrMessage) {
	var alarmReceiver []model.BatchAlarmReceiver

	status := dbClient.Debug().Table("batch_alarm_receivers").Where("service_type = 'SaaS'").Find(&alarmReceiver)

	err := util.GetError().DbCheckError(status.Error)
	if err != nil {
		return nil, err
	}

	return alarmReceiver, nil
}

//// 알람 수신자 조회
func GetBatchAlarmLog(dbClient *gorm.DB) ([]model.BatchAlarmExecution, model.ErrMessage) {
	var alarmLog []model.BatchAlarmExecution
	status := dbClient.Debug().Table("batch_alarm_executions").Where("service_type = 'SaaS' and critical_status <> 'Success' ").Find(&alarmLog)

	err := util.GetError().DbCheckError(status.Error)
	if err != nil {
		return nil, err
	}

	return alarmLog, nil
}

func InsertAlarmInfo(dbClient *gorm.DB, udateData []gjson.Result, timeData int) model.ErrMessage {
	var err model.ErrMessage

	// Delete And Insert
	status := dbClient.Debug().Where("service_type = 'SaaS'").Delete(model.BatchAlarmInfo{})

	for _, data := range udateData {
		batchAlarmInfo := model.BatchAlarmInfo{}
		tempMap := data.Map()
		tempWaring, _ := strconv.Atoi(tempMap["Warning"].String())
		tempCritical, _ := strconv.Atoi(tempMap["Critical"].String())

		cronExpression := "*/" + tempMap["Delay"].String() + " * * * *"
		batchAlarmInfo.ServiceType = "SaaS"
		batchAlarmInfo.MetricType = tempMap["Name"].String()
		batchAlarmInfo.WarningValue = tempWaring
		batchAlarmInfo.CriticalValue = tempCritical
		batchAlarmInfo.MeasureTime = timeData
		batchAlarmInfo.CronExpression = cronExpression

		if tempMap["Name"].String() == "SYSTEM_CPU" {
			batchAlarmInfo.ExecMsg = "SaaS Application : ${AppName} System CPU 현재사용률 (${Currend_value}%)"
			batchAlarmInfo.ParamData1 = "/getAgentStat/cpuLoad/chart.pinpoint"
			batchAlarmInfo.ParamData2 = "charts.y.CPU_LOAD_SYSTEM.#.2"
			batchAlarmInfo.ParamData3 = ""
		}

		if tempMap["Name"].String() == "JVM_CPU" {
			batchAlarmInfo.ExecMsg = "SaaS Application : ${AppName} JVM CPU 현재사용률 (${Currend_value}%)"
			batchAlarmInfo.ParamData1 = "/getAgentStat/cpuLoad/chart.pinpoint"
			batchAlarmInfo.ParamData2 = "charts.y.CPU_LOAD_JVM.#.2"
			batchAlarmInfo.ParamData3 = ""
		}

		if tempMap["Name"].String() == "HEAP_MEMORY" {
			batchAlarmInfo.ExecMsg = "SaaS Application : ${AppName} Heap Memory 현재사용률 (${Currend_value}%)"
			batchAlarmInfo.ParamData1 = "/getAgentStat/jvmGc/chart.pinpoint"
			batchAlarmInfo.ParamData2 = "charts.y.JVM_MEMORY_HEAP_USED.#.2"
			batchAlarmInfo.ParamData3 = "charts.y.JVM_MEMORY_HEAP_MAX.#.2"
		}

		status = dbClient.Debug().Create(&batchAlarmInfo)
		err = util.GetError().DbCheckError(status.Error)

	}
	return err
}

func InsertAlarmReceivers(dbClient *gorm.DB, receiverId string, emailData string, snsId int64) model.ErrMessage {
	var err model.ErrMessage
	batchAlarmReceiver := model.BatchAlarmReceiver{}
	// Delete And Insert
	status := dbClient.Debug().Where("service_type = 'SaaS'").Delete(model.BatchAlarmReceiver{})
	err = util.GetError().DbCheckError(status.Error)

	//batchAlarmReceiver.ReceiverId	= ""  autoincrement
	batchAlarmReceiver.ServiceType = "SaaS"
	batchAlarmReceiver.Name = "Admin"
	batchAlarmReceiver.Email = emailData
	batchAlarmReceiver.SnsId = snsId
	batchAlarmReceiver.UseYn = "Y"

	status = dbClient.Debug().Create(&batchAlarmReceiver)
	err = util.GetError().DbCheckError(status.Error)

	return err
}
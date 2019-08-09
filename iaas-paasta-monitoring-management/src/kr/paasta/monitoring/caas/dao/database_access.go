package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring/caas/model"
	"kr/paasta/monitoring/caas/util"
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
	status := dbClient.Debug().Table("batch_alarm_infos").Where("service_type = 'CaaS'").Find(&alarmInfos)

	err := util.GetError().DbCheckError(status.Error)
	if err != nil {
		return nil, err
	}

	return alarmInfos, nil
}

//// 알람 수신자 조회
func GetBatchAlarmReceiver(dbClient *gorm.DB) ([]model.BatchAlarmReceiver, model.ErrMessage) {
	var alarmReceiver []model.BatchAlarmReceiver

	status := dbClient.Debug().Table("batch_alarm_receivers").Where("service_type = 'CaaS'").Find(&alarmReceiver)

	err := util.GetError().DbCheckError(status.Error)
	if err != nil {
		return nil, err
	}

	return alarmReceiver, nil
}

//// 알람 수신자 조회
func GetBatchAlarmLog(dbClient *gorm.DB) ([]model.BatchAlarmExecution, model.ErrMessage) {
	var alarmLog []model.BatchAlarmExecution
	var tempAlarmDiv []alarmId
	var quaryAlarmDiv string

	temp := dbClient.Debug().Table("batch_alarm_infos").
		Select("alarm_id").
		Where("service_type = 'CaaS'").Find(&tempAlarmDiv)

	err1 := util.GetError().DbCheckError(temp.Error)
	if err1 != nil {
		return nil, err1
	}

	for _, data := range tempAlarmDiv {
		quaryAlarmDiv += data.AlarmId + ","
	}

	quaryAlarmDiv += "''"

	status := dbClient.Debug().Table("batch_alarm_executions").Where("alarm_id in (" + quaryAlarmDiv + ") and critical_status <> 'Success' ").Find(&alarmLog)

	err := util.GetError().DbCheckError(status.Error)
	if err != nil {
		return nil, err
	}

	return alarmLog, nil
}

//func InsertBatchExecution(dbClient *gorm.DB, batchExection *model.BatchAlarmExecution) {
//	if err := dbClient.Debug().Create(&batchExection).Error; err != nil {
//		fmt.Printf("insert error : %v\n", dbClient.Error)
//	}
//}

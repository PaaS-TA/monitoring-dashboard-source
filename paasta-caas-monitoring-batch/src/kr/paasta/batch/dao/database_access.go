package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/thoas/go-funk"
	"kr/paasta/batch/model"
	"kr/paasta/batch/util"
	"log"
	"os"
)

var dbType string
var connectionString string

func init() {
	// 기본적인 프로퍼티 설정 정보 읽어오기
	config, err := util.ReadConfig(`config.ini`)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	dbType = config["monitoring.db.type"]
	dbName := config["monitoring.db.dbname"]
	userName := config["monitoring.db.username"]
	userPassword := config["monitoring.db.password"]
	host := config["monitoring.db.host"]
	port := config["monitoring.db.port"]

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

func GetBatchAlarmInfo(dbClient *gorm.DB) []model.BatchAlarmInfo {
	var alarmInfos []model.BatchAlarmInfo
	dbClient.Debug().Find(&alarmInfos)
	return alarmInfos
}

func InsertBatchExecution(dbClient *gorm.DB, batchExection *model.BatchAlarmExecution) {
	if err := dbClient.Debug().Create(&batchExection).Error; err != nil {
		fmt.Printf("insert error : %v\n", dbClient.Error)
	}
}

// 알람 수신자 조회
func GetBatchAlarmReceiver(serviceType string, dbClient *gorm.DB) ([]string, []int64) {
	var alarmReceiver []model.BatchAlarmReceiver
	dbClient.Debug().Find(&alarmReceiver)

	var mapedEmail []string
	var mapedSnsId []int64

	if len(alarmReceiver) > 0 {
		filterEmail := funk.Filter(alarmReceiver, func(x model.BatchAlarmReceiver) bool {
			return len(x.Email) > 0 && x.ServiceType == serviceType && x.UseYn == "Y"
		}).([]model.BatchAlarmReceiver)

		mapedEmail = funk.Map(filterEmail, func(x model.BatchAlarmReceiver) string {
			return x.Email
		}).([]string)

		filterSnsId := funk.Filter(alarmReceiver, func(x model.BatchAlarmReceiver) bool {
			return x.SnsId > 0 && x.ServiceType == serviceType && x.UseYn == "Y"
		}).([]model.BatchAlarmReceiver)

		mapedSnsId = funk.Map(filterSnsId, func(x model.BatchAlarmReceiver) int64 {
			return x.SnsId
		}).([]int64)

		return mapedEmail, mapedSnsId
	}
	return nil, nil
}

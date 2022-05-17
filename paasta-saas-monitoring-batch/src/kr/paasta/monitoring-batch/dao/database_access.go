package dao

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/thoas/go-funk"
	"log"
	"os"
	"saas-monitoring-batch/model"
	"saas-monitoring-batch/util"
	"strconv"
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
	dbClient.Debug().AutoMigrate(&model.BatchAlarmInfo{}, &model.BatchAlarmExecution{}, &model.BatchAlarmReceiver{}, &model.BatchAlarmSns{}, &model.BatchAlarmExecutionResolve{})
}

func GetBatchAlarmInfo(dbClient *gorm.DB) []model.BatchAlarmInfo {
	var alarmInfos []model.BatchAlarmInfo
	dbClient.Debug().Find(&alarmInfos)
	return alarmInfos
}

func InsertBatchExecution(dbClient *gorm.DB, batchExection *model.BatchAlarmExecution) {
	sql := "INSERT  INTO batch_alarm_executions (alarm_id, service_type, critical_status, measure_value, measure_name1, measure_name2, measure_name3, execution_time, execution_result, resolve_status) "
	sql += "VALUES"
	sql += "(%s, '%s', '%s', %s, '%s', '', '',  now(), '%s', '1')"
	insertSql := fmt.Sprintf(sql, strconv.Itoa(batchExection.AlarmId), batchExection.ServiceType, batchExection.CriticalStatus,
		strconv.FormatFloat(batchExection.MeasureValue, 'f', 6, 64), batchExection.MeasureName1, batchExection.ExecutionResult)
	if err := dbClient.Debug().Exec(insertSql).Error; err != nil {
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
			return x.ReceiveType == "EMAIL" && x.ServiceType == serviceType && x.UseYn == "Y"
		}).([]model.BatchAlarmReceiver)

		mapedEmail = funk.Map(filterEmail, func(x model.BatchAlarmReceiver) string {
			return x.TargetId
		}).([]string)

		filterSnsId := funk.Filter(alarmReceiver, func(x model.BatchAlarmReceiver) bool {
			return x.ReceiveType == "SNS" && x.ServiceType == serviceType && x.UseYn == "Y"
		}).([]model.BatchAlarmReceiver)

		mapedSnsId = funk.Map(filterSnsId, func(x model.BatchAlarmReceiver) int64 {
			id, _ := strconv.ParseInt(x.TargetId, 10, 64)
			return int64(id)
		}).([]int64)

		fmt.Printf("mapedSnsId : %v\n", mapedSnsId)
		return mapedEmail, mapedSnsId
	}
	return nil, nil
}

func GetBatchAlarmSnsToken(serviceType string, dbClient *gorm.DB) model.BatchAlarmSns {
	var alarmSns model.BatchAlarmSns
	dbClient.Debug().Table("batch_alarm_sns").Where("origin_type = '" + serviceType + "'").Find(&alarmSns)
	return alarmSns
}

func SaveBatchAlarmSnsReceiver(serviceType string, dbClient *gorm.DB, targetIds []string) {
	tx := dbClient.Begin().Debug()
	for _, targetId := range targetIds {
		alarmReceiver := model.BatchAlarmReceiver{
			ServiceType: serviceType,
			ReceiveType: "SNS",
			TargetId:    targetId,
		}

		status := tx.Table("batch_alarm_receivers").
			Set("gorm:insert_option", "on duplicate key update modi_date = now(), modi_user = 'system'").Create(&alarmReceiver)
		if err := status.Error; err != nil {
			fmt.Printf("insert error : %v\n", dbClient.Error)
			tx.Rollback()
		}
	}
	tx.Commit()
}

package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring/caas/model"
	"kr/paasta/monitoring/caas/util"
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

func CreateTable(dbClient *gorm.DB) {
	dbClient.Debug().AutoMigrate(&model.BatchAlarmInfo{}, &model.BatchAlarmExecution{}, &model.BatchAlarmReceiver{}, &model.BatchAlarmSns{}, &model.BatchAlarmExecutionResolve{})
}

// Alarm Info
func GetBatchAlarmInfo(dbClient *gorm.DB) ([]model.BatchAlarmInfo, model.ErrMessage) {
	var alarmInfos []model.BatchAlarmInfo
	status := dbClient.Debug().Table("batch_alarm_infos").Where("service_type = 'CaaS'").Find(&alarmInfos)

	err := util.GetError().DbCheckError(status.Error)

	if err != nil {
		return nil, err
	}
	return alarmInfos, nil
}

//// 알람 수신자 조회
func GetBatchAlarmReceiver(dbClient *gorm.DB, receiveType string) ([]model.AlarmrReceiverResponse, model.ErrMessage) {
	var alarmReceiver []model.AlarmrReceiverResponse
	status := dbClient.Debug().Table("batch_alarm_receivers").Where("service_type = 'CaaS' AND receive_type = '" + receiveType + "'").Find(&alarmReceiver)

	err := util.GetError().DbCheckError(status.Error)
	if err != nil {
		return nil, err
	}

	return alarmReceiver, nil
}

//// 알람 수신자 조회
func GetBatchAlarmLog(dbClient *gorm.DB, searchDateFrom string, searchDateTo string, alarmType string, alarmStatus string, resolveStatus string) ([]model.BatchAlarmExecution, model.ErrMessage) {
	var queryWhere string
	queryWhere = ""

	if len(searchDateFrom) > 0 && len(searchDateTo) > 0 {
		queryWhere += " AND execution_time BETWEEN '" + searchDateFrom + " 00:00:00' AND '" + searchDateTo + " 23:59:59' "
	}

	if len(alarmType) > 0 {
		queryWhere += " AND execution_result LIKE '%" + alarmType + "%' "
	}

	if len(alarmStatus) > 0 {
		queryWhere += " AND critical_status = '" + alarmStatus + "' "
	}
	if len(resolveStatus) > 0 {
		queryWhere += " AND resolve_status = '" + resolveStatus + "' "
	}

	var alarmLog []model.BatchAlarmExecution
	status := dbClient.Debug().Table("batch_alarm_executions").Where("service_type = 'CaaS' and critical_status <> 'Success' " + queryWhere).Find(&alarmLog)

	err := util.GetError().DbCheckError(status.Error)

	if err != nil {
		return nil, err
	}

	return alarmLog, nil
}

func GetBatchAlarmResolve(dbClient *gorm.DB, id uint64) ([]model.AlarmrRsolveResponse, model.ErrMessage) {
	var alarmRsolves []model.AlarmrRsolveResponse
	status := dbClient.Debug().Table("batch_alarm_execution_resolves").Select(" resolve_id , alarm_action_desc,  reg_date").Where("excution_id = " + strconv.Itoa(int(id))).Find(&alarmRsolves)

	err := util.GetError().DbCheckError(status.Error)

	if err != nil {
		return nil, err
	}

	return alarmRsolves, nil
}

func InsertAlarmInfo(dbClient *gorm.DB, request []model.AlarmPolicyRequest, email string, emailUseYn string) model.ErrMessage {
	var err model.ErrMessage

	tx := dbClient.Begin().Debug()
	// Delete And Insert
	status := tx.Where("service_type = 'CaaS'").Delete(model.BatchAlarmInfo{})
	for _, data := range request {
		batchAlarmInfo := model.BatchAlarmInfo{}
		repeatTime := strconv.Itoa(data.RepeatTime)

		batchAlarmInfo.ServiceType = data.OriginType
		batchAlarmInfo.MetricType = data.AlarmType
		batchAlarmInfo.WarningValue = data.WarningThreshold
		batchAlarmInfo.CriticalValue = data.CriticalThreshold
		batchAlarmInfo.MeasureTime = data.MeasureTime
		batchAlarmInfo.CronExpression = "*/" + repeatTime + " * * * *"
		batchAlarmInfo.ExecMsg = "CaaS PodName : ${PodName} 현재사용률 " + data.AlarmType + " (${Currend_value}%)"
		batchAlarmInfo.ParamData1 = ""
		batchAlarmInfo.ParamData2 = ""
		batchAlarmInfo.ParamData3 = ""

		status = tx.Create(&batchAlarmInfo)

		err = util.GetError().DbCheckError(status.Error)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Email Receiver 지정
	alarmReceiver := model.BatchAlarmReceiver{
		ServiceType: "CaaS",
		ReceiveType: "EMAIL",
		TargetId:    email,
		UseYn:       emailUseYn,
	}

	status = tx.Where("service_type = 'CaaS' AND receive_type = 'EMAIL'").Delete(model.BatchAlarmReceiver{})
	err = util.GetError().DbCheckError(status.Error)
	if err != nil {
		tx.Rollback()
		return err
	}

	status = tx.Create(&alarmReceiver)
	err = util.GetError().DbCheckError(status.Error)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return err
}

func GetSnsInfo(dbClient *gorm.DB) (interface{}, model.ErrMessage) {
	var alarmSns model.BatchAlarmSnsRequest
	status := dbClient.Debug().Table("batch_alarm_sns").Where("origin_type = 'CaaS'").Find(&alarmSns)

	err := util.GetError().DbCheckError(status.Error)

	if err != nil {
		return nil, err
	}

	return alarmSns, nil
}

func GetAlarmCount(dbClient *gorm.DB, searchDateFrom string, searchDateTo string) (model.AlarmCount, model.ErrMessage) {
	var queryWhere string
	if len(searchDateFrom) > 0 && len(searchDateTo) > 0 {
		queryWhere = " AND execution_time BETWEEN '" + searchDateFrom + " 00:00:00' AND '" + searchDateTo + " 23:59:59' "
	} else {
		queryWhere = ""
	}

	var alarmCnt int
	status := dbClient.Debug().Table("batch_alarm_executions").Where("critical_status != 'Success' and service_type = 'CaaS' " + queryWhere).Count(&alarmCnt)

	err := util.GetError().DbCheckError(status.Error)

	if err != nil {
		return model.AlarmCount{}, err
	}

	alarmCount := model.AlarmCount{AlarmCnt: alarmCnt}

	return alarmCount, nil
}

func GetAlarmSnsSave(dbClient *gorm.DB, alarmSns model.BatchAlarmSnsRequest) model.ErrMessage {
	fmt.Printf("alarmSns : %v\n", alarmSns)
	status := dbClient.Debug().Table("batch_alarm_sns").
		Set("gorm:insert_option", "on duplicate key update modi_date = now(), modi_user = 'system', sns_id = '"+alarmSns.SnsId+"',  token = '"+alarmSns.Token+"', expl = '"+alarmSns.Expl+"' , sns_send_yn = '"+alarmSns.SnsSendYn+"'").Create(&alarmSns)
	err := util.GetError().DbCheckError(status.Error)

	if err != nil {
		fmt.Printf("error : %v\n", err)
		return err
	}
	return err
}

func UpdateAlarmSate(dbClient *gorm.DB, request model.AlarmrRsolveRequest) model.ErrMessage {
	var status *gorm.DB
	if request.ResolveStatus == "3" {
		status = dbClient.Debug().Table("batch_alarm_executions").Where("excution_id = ? ", request.Id).
			Updates(map[string]interface{}{"complete_date": util.GetDBCurrentTime(), "resolve_status": request.ResolveStatus})
	} else {
		status = dbClient.Debug().Table("batch_alarm_executions").Where("excution_id = ? ", request.Id).
			Updates(map[string]interface{}{"resolve_status": request.ResolveStatus})
	}
	err := util.GetError().DbCheckError(status.Error)

	if err != nil {
		fmt.Printf("error : %v\n", err)
	}

	return err
}

func CreateAlarmResolve(dbClient *gorm.DB, request model.AlarmrRsolveRequest) model.ErrMessage {
	var alarmExecutionResolve model.BatchAlarmExecutionResolve
	alarmExecutionResolve.ExcutionId = request.Id
	alarmExecutionResolve.AlarmActionDesc = request.AlarmActionDesc
	alarmExecutionResolve.RegDate = util.GetDBCurrentTime()

	status := dbClient.Debug().Table("batch_alarm_execution_resolves").Create(&alarmExecutionResolve)

	err := util.GetError().DbCheckError(status.Error)

	if err != nil {
		fmt.Printf("error : %v\n", err)
	}
	return err
}

func UpdateAlarmResolve(dbClient *gorm.DB, request model.AlarmrRsolveRequest) model.ErrMessage {
	status := dbClient.Debug().Table("batch_alarm_execution_resolves").Where("resolve_id = ? ", request.Id).
		Updates(map[string]interface{}{"alarm_action_desc": request.AlarmActionDesc})

	err := util.GetError().DbCheckError(status.Error)

	if err != nil {
		fmt.Printf("error : %v\n", err)
	}
	return err
}

func DeleteAlarmResolve(dbClient *gorm.DB, id uint64) model.ErrMessage {
	status := dbClient.Debug().Table("batch_alarm_execution_resolves").Where("resolve_id = " + strconv.Itoa(int(id))).Delete(model.BatchAlarmExecutionResolve{})
	err := util.GetError().DbCheckError(status.Error)

	if err != nil {
		fmt.Printf("error : %v\n", err)
	}
	return err
}

func DeleteAlarmSnsChannel(dbClient *gorm.DB, id int) model.ErrMessage {
	alarmSns := model.BatchAlarmSns{
		ChannelId: id,
	}
	status := dbClient.Debug().Delete(&alarmSns)

	err := util.GetError().DbCheckError(status.Error)

	if err != nil {
		fmt.Printf("error : %v\n", err)
	}
	return err
}

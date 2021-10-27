package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring-batch/model"
	"kr/paasta/monitoring-batch/model/base"
	"kr/paasta/monitoring-batch/util"
	"log"
	"time"
)

func GetAlarmPolicy(dbConn *gorm.DB) ([]model.AlarmPolicy, base.ErrMessage){
	var alarmPolicy []model.AlarmPolicy
	status := dbConn.Model(&alarmPolicy).Where("origin_type=?", "ias").Find(&alarmPolicy)

	err := util.GetError().DbCheckError(status.Error)
	if err != nil {
		log.Fatalf("Error::", err)
		return nil, err
	}

	return alarmPolicy, nil
}


func GetAlarmTarget(dbConn *gorm.DB) []model.AlarmTarget {
	var alarmTarget []model.AlarmTarget
	dbConn.Debug().Table("alarm_targets").Where("origin_type = 'ias' and mail_send_yn = 'Y'").Find(&alarmTarget).RecordNotFound()
	return alarmTarget
}

func IsExistAlarm(alarmData model.Alarm, txn *gorm.DB) (bool, model.Alarm) {
	//resolve_status
	// 1 : Alarm 발생, 2 : Alarm 접수
	// Alarm발생, 접수 상태면 Alarm을 새로 발생하지 않는다.
	var alarm model.Alarm
	isNew := txn.Model(&alarm).Where("origin_type = ? and origin_id = ? and alarm_type = ? and ( resolve_status = '1' || resolve_status = '2') and level = ? ", alarmData.OriginType, alarmData.OriginId, alarmData.AlarmType, alarmData.Level).
		Find(&alarm).RecordNotFound()
	return !isNew, alarm
}

func InsertAlarm(alarm model.Alarm, dbConn *gorm.DB) base.ErrMessage {
	alarm.AlarmCnt = 1
	alarm.RegUser = base.BAT_USER
	alarm.ModiUser = base.BAT_USER
	alarm.RegDate = time.Now()
	alarm.ModiDate = time.Now()
	alarm.AlarmSendDate = time.Now()

	status := dbConn.Create(&alarm)
	err := util.GetError().DbCheckError(status.Error)
	if err != nil {
		fmt.Println("Error::", err)
	}
	return err
}

func UpdateAlarm(alarm model.Alarm, dbConn *gorm.DB) base.ErrMessage {
	status := dbConn.Model(&alarm).Where("origin_type = ? and origin_id = ? and alarm_type = ? and level = ? and resolve_status = '1'", alarm.OriginType, alarm.OriginId, alarm.AlarmType, alarm.Level).
		Updates(map[string]interface{}{"alarm_send_date": time.Now(), "modi_date": time.Now(), "modi_user": base.BAT_USER})
	err := util.GetError().DbCheckError(status.Error)
	if err != nil {
		fmt.Println("Error::", err)
		return err
	}
	return err
}

func GetAlarm(alarm model.Alarm, dbConn *gorm.DB) model.Alarm {
	var result model.Alarm
	dbConn.Model(&alarm).
		Where("origin_type = ? and origin_id = ? and alarm_type = ? and resolve_status = '1'", alarm.OriginType, alarm.OriginId, alarm.AlarmType).
		Find(&result)

	return result
}

func GetAlarmSns(dbConn *gorm.DB) ([]model.AlarmSns, base.ErrMessage) {
	where := "sns_send_yn = 'Y'"
	where += " and (origin_type='ias' or origin_type='all')"

	var result []model.AlarmSns
	status := dbConn.Table("alarm_sns").Where(where).Find(&result)
	err := util.GetError().DbCheckError(status.Error)
	return result, err
}

func GetAlarmSnsTarget(alarmSns model.AlarmSns, txn *gorm.DB) ([]model.AlarmSnsTarget, base.ErrMessage) {

	where := fmt.Sprintf("channel_id='%d'", alarmSns.ChannelId)
	var alarmSnsTargets []model.AlarmSnsTarget
	status := txn.Debug().Table("alarm_sns_targets").Where(where).Find(&alarmSnsTargets)
	err := util.GetError().DbCheckError(status.Error)
	return alarmSnsTargets, err
}

func UpdateSnsAlarmTargets(alarmSnsTarget model.AlarmSnsTarget, txn *gorm.DB) base.ErrMessage {

	status := txn.Table("alarm_sns_targets").
		Set("gorm:insert_option", "on duplicate key update modi_date = now(), modi_user = 'system'").Create(&alarmSnsTarget)
	err := util.GetError().DbCheckError(status.Error)
	return err
}
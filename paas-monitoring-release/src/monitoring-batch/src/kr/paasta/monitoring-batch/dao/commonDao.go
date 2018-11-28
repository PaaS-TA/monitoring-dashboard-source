package dao

import (
	client "github.com/influxdata/influxdb/client/v2"
	cb "kr/paasta/monitoring-batch/model/base"
	"kr/paasta/monitoring-batch/model"
	"kr/paasta/monitoring-batch/util"
	"github.com/jinzhu/gorm"
	"fmt"
	"time"
)

type commonStruct struct {
	influxClient 	client.Client
}


func GetCommonDao(influxClient client.Client) *commonStruct{
	return &commonStruct{
		influxClient: 	influxClient,
	}
}

func (f commonStruct) GetIsNotExistAlarm( alarmData cb.Alarm,  txn *gorm.DB) (bool, cb.Alarm){

	//resolve_status
	// 1 : Alarm 발생, 2 : Alarm 접수
	// Alarm발생, 접수 상태면 Alarm을 새로 발생하지 않는다.
	var alarm cb.Alarm
	isNew := txn.Debug().Model(&alarm).Where("origin_type = ? and origin_id = ? and alarm_type = ? and ( resolve_status = '1' || resolve_status = '2') and level = ? ", alarmData.OriginType, alarmData.OriginId, alarmData.AlarmType, alarmData.Level).
		Find(&alarm).RecordNotFound()
	return isNew, alarm
}

//TestCode용
func (f commonStruct) GetIsNotExistAlarmCheck( alarmData cb.Alarm,  txn *gorm.DB) (bool, cb.Alarm){

	//resolve_status
	// 1 : Alarm 발생, 2 : Alarm 접수
	// Alarm발생, 접수 상태면 Alarm을 새로 발생하지 않는다.
	var alarm cb.Alarm
	isNew := txn.Debug().Model(&alarm).Where("origin_type = ? and resolve_status = '1' ", alarmData.OriginType).
		Find(&alarm).RecordNotFound()
	return isNew, alarm
}


func (b commonStruct) CreateAlarmData(boshAlarm cb.Alarm, txn *gorm.DB) cb.ErrMessage{

	eventData := cb.Alarm{OriginId: boshAlarm.OriginId, OriginType: boshAlarm.OriginType, AlarmType: boshAlarm.AlarmType, Level: boshAlarm.Level,
		AppYn: boshAlarm.AppYn, Ip: boshAlarm.Ip, AlarmTitle: boshAlarm.AlarmTitle, AlarmMessage: boshAlarm.AlarmMessage , ResolveStatus: boshAlarm.ResolveStatus, AlarmCnt: 1,
		RegDate: time.Now(), RegUser: "Bat" , ModiUser: cb.BAT_USER, ModiDate: time.Now(), AlarmSendDate: time.Now()}
	status := txn.Debug().Create(&eventData)
	err := util.GetError().DbCheckError(status.Error)
	if err != nil{
		fmt.Println("Error::", err )
	}
	return  err
}

func (b commonStruct) UpdateSendDate(boshAlarm cb.Alarm, txn *gorm.DB) cb.ErrMessage{

	status := txn.Debug().Model(&boshAlarm).Where("origin_type = ? and origin_id = ? and alarm_type = ? and level = ? and resolve_status = '1'", boshAlarm.OriginType, boshAlarm.OriginId, boshAlarm.AlarmType, boshAlarm.Level ).
		Updates(map[string]interface{}{ "alarm_send_date": time.Now(), "modi_date": time.Now(), "modi_user": cb.BAT_USER})
	err := util.GetError().DbCheckError(status.Error)
	if err != nil{
		fmt.Println("Error::", err )
		return   err
	}
	return  err
}

func (b commonStruct) GetAlarmData(alarm cb.Alarm, txn *gorm.DB) (bool, model.Alarm) {

	var alarmData model.Alarm
	isNew := txn.Debug().Model(&alarm).Where("origin_type = ? and origin_id = ? and alarm_type = ? and resolve_status = '1'", alarm.OriginType, alarm.OriginId, alarm.AlarmType ).
		Find(&alarmData).RecordNotFound()
	return isNew, alarmData

}

func (b commonStruct) GetAlarmSns(orginType string, txn *gorm.DB) ([]model.AlarmSns, cb.ErrMessage) {

	where := "sns_send_yn = 'Y'"
	if orginType != "" {
		where += fmt.Sprintf(" and (origin_type='%s' or origin_type='all')", orginType)
	}
	var alarmSns []model.AlarmSns
	status := txn.Debug().Table("alarm_sns").Where(where).Find(&alarmSns)
	err := util.GetError().DbCheckError(status.Error)
	return alarmSns, err
}

func (b commonStruct) GetAlarmSnsTarget(alarmSns model.AlarmSns, txn *gorm.DB) ([]model.AlarmSnsTarget, cb.ErrMessage) {

	where := fmt.Sprintf("channel_id='%d'", alarmSns.ChannelId)
	var alarmSnsTargets []model.AlarmSnsTarget
	status := txn.Debug().Table("alarm_sns_targets").Where(where).Find(&alarmSnsTargets)
	err := util.GetError().DbCheckError(status.Error)
	return alarmSnsTargets, err
}

func (b commonStruct) UpdateSnsAlarmTargets(alarmSnsTarget model.AlarmSnsTarget, txn *gorm.DB) (cb.ErrMessage) {

	status := txn.Debug().Table("alarm_sns_targets").
		Set("gorm:insert_option", "on duplicate key update modi_date = now(), modi_user = 'system'").Create(&alarmSnsTarget)
	err := util.GetError().DbCheckError(status.Error)
	return err
}


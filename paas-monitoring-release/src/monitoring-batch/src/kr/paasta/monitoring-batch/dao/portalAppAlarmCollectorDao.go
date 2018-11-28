package dao

import (
	client "github.com/influxdata/influxdb/client/v2"
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring-batch/model/base"
	"kr/paasta/monitoring-batch/model"
	"kr/paasta/monitoring-batch/util"
	"fmt"
	"time"
)

type portalAppAlarmDao struct {
	txn *gorm.DB
	influxClient client.Client
	influxDbName string
}

func GetPortalAppAlarmDao(txn *gorm.DB, influxClient client.Client, influxDbName string) *portalAppAlarmDao{
	return &portalAppAlarmDao{
		influxClient: influxClient,
		txn: txn,
		influxDbName: influxDbName,
	}
}

func (p portalAppAlarmDao) GetNotTerminatedAppAlarm() ([]model.AppAlarmHistory, base.ErrMessage) {

	var listNotTerminatedAlarm []model.AppAlarmHistory
	status := p.txn.Debug().Table("app_alarm_histories").Where("terminate_yn = 'N'").Find(&listNotTerminatedAlarm)
	err := util.GetError().DbCheckError(status.Error)
	return listNotTerminatedAlarm, err
}

func (p portalAppAlarmDao) GetAppAlarmPolicy() ([]model.AppAlarmPolicy, base.ErrMessage) {

	var listAppAlarmPolicy []model.AppAlarmPolicy
	status := p.txn.Debug().Table("app_alarm_policies").Where("alarm_use_yn = 'Y'").Find(&listAppAlarmPolicy)
	err := util.GetError().DbCheckError(status.Error)
	return listAppAlarmPolicy, err
}

func (p portalAppAlarmDao) GetSendTargetAppAlarm(interval int) ([]model.SendTargetAppAlarmHistory, base.ErrMessage) {

	var listSendTargetAppAlarm []model.SendTargetAppAlarmHistory

	query := "app_alarm_histories.alarm_id,app_alarm_histories.app_guid,app_alarm_histories.app_idx,app_alarm_histories.resource_type," +
		"app_alarm_histories.alarm_level,app_alarm_histories.app_name,app_alarm_histories.cell_ip,app_alarm_histories.container_id," +
		"app_alarm_histories.container_interface,app_alarm_histories.alarm_title,app_alarm_histories.alarm_message," +
		"app_alarm_histories.alarm_send_date,app_alarm_histories.terminate_yn," +
		"app_alarm_histories.reg_date,app_alarm_histories.reg_user,app_alarm_histories.modi_date,app_alarm_histories.modi_user," +
		"app_alarm_policies.email,app_alarm_policies.email_send_yn"

	join := "INNER JOIN app_alarm_policies ON " +
		"app_alarm_histories.app_guid = app_alarm_policies.app_guid " +
		"AND app_alarm_policies.alarm_use_yn = 'Y' " +
		"AND app_alarm_histories.terminate_yn = 'N' " +
		"AND (app_alarm_histories.alarm_send_date is null OR app_alarm_histories.alarm_send_date < DATE_SUB(?, INTERVAL ? MINUTE))"

	status := p.txn.Debug().Table("app_alarm_histories").Select(query).Joins(join, time.Now(), interval).Scan(&listSendTargetAppAlarm)
	err := util.GetError().DbCheckError(status.Error)
	return listSendTargetAppAlarm, err
}

func (p portalAppAlarmDao) UpdateTerminatedAlarm(listTerminated []model.AppAlarmHistory) base.ErrMessage {

	var listAlarmId []uint
	for _, v := range listTerminated {
		listAlarmId = append(listAlarmId, v.AlarmId)
	}

	status := p.txn.Debug().Table("app_alarm_histories").Where("alarm_id IN (?)", listAlarmId).
		Updates(map[string]interface{}{"terminate_yn": "Y", "modi_date": time.Now()})
	err := util.GetError().DbCheckError(status.Error)
	return err
}

func (p portalAppAlarmDao) UpdateContinuousAppAlarm(updated model.AppAlarmHistory) base.ErrMessage {

	where := "app_guid = ? " +
		"AND app_idx = ? " +
		"AND resource_type = ? " +
		"AND container_interface = ? " +
		"AND terminate_yn = 'N'"

	update := map[string]interface{}{
		"alarm_level": updated.AlarmLevel,
		"alarm_title": updated.AlarmTitle,
		"alarm_message": updated.AlarmMessage,
		"modi_date": time.Now(),
	}

	status := p.txn.Debug().Table("app_alarm_histories").
		Where(where, updated.AppGuid, updated.AppIdx, updated.ResourceType, updated.ContainerInterface).Update(update)
	err := util.GetError().DbCheckError(status.Error)
	return err
}

func (p portalAppAlarmDao) InsertNewAppAlarm(new model.AppAlarmHistory) base.ErrMessage {

	t := time.Now()
	new.RegDate, new.ModiDate = t, t
	status := p.txn.Debug().Table("app_alarm_histories").Create(&new)
	err := util.GetError().DbCheckError(status.Error)
	return err
}

func (p portalAppAlarmDao) UpdateAlarmSendDate(alarm model.SendTargetAppAlarmHistory) base.ErrMessage {

	t := time.Now()
	status := p.txn.Debug().Table("app_alarm_histories").Where("alarm_id = ?",	alarm.AlarmId).
		Update("alarm_send_date", t, "modi_date", t)
	err := util.GetError().DbCheckError(status.Error)
	return err
}

func (p portalAppAlarmDao) GetAppInfo(appGuid string) (_ client.Response, errMsg base.ErrMessage)  {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = base.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	sql :=	"SELECT application_id, application_name, application_index, cell_ip, container_id, container_interface, value " +
		"FROM container_metrics " +
		"WHERE time > now() - 1m AND application_id = '%s' " +
		"GROUP BY container_interface LIMIT 1;"
	q := client.Query{
		Command:  fmt.Sprintf(sql, appGuid),
		Database: p.influxDbName,
	}
	resp, err := p.influxClient.Query(q)
	if err != nil{
		errLogMsg = err.Error()
	}

	return util.GetError().CheckError(*resp, err)
}

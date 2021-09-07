package dao

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring/iaas_new/model"
	"kr/paasta/monitoring/utils"
	"strconv"
	"time"
)

type AlarmDao struct {
	influxClient client.Client
	txn          *gorm.DB
}

func GetAlarmDao(influxClient client.Client, txn *gorm.DB) *AlarmDao {
	return &AlarmDao{
		influxClient: influxClient,
		txn:          txn,
	}
}

func (a *AlarmDao) GetAlarmsActionHistoryList(request model.AlarmActionRequest) ([]model.AlarmActionHistory, error) {

	fmt.Println("alarmRequest.AlarmId::", request.AlarmId)
	t := []model.AlarmActionHistory{}
	status := a.txn.Debug().Table("alarm_action_histories").
		Select("id, alarm_id,  alarm_action_desc, reg_date + INTERVAL "+strconv.Itoa(model.GmtTimeGap)+" HOUR  as reg_date , reg_user, modi_date + INTERVAL "+strconv.Itoa(model.GmtTimeGap)+" HOUR  as modi_date  , modi_user").
		Where("alarm_id = ?", request.AlarmId).
		Order("reg_date desc").
		Find(&t)

	if status.Error != nil {
		return t, status.Error
	}
	return t, nil
}

func (a *AlarmDao) CreateAlarmAction(request model.AlarmActionRequest) error {

	actionData := model.AlarmActionHistory{AlarmId: request.AlarmId, AlarmActionDesc: request.AlarmActionDesc}

	status := a.txn.Debug().Create(&actionData)

	if status.Error != nil {
		return status.Error
	}
	return status.Error
}

func (a *AlarmDao) UpdateAlarmAction(request model.AlarmActionRequest) error {

	status := a.txn.Debug().Table("alarm_action_histories").Where("id = ? ", request.Id).
		Updates(map[string]interface{}{"alarm_action_desc": request.AlarmActionDesc, "modi_date": time.Now()})
	if status.Error != nil {
		return status.Error
	}
	return status.Error
}

func (a *AlarmDao) DeleteAlarmAction(request model.AlarmActionRequest) error {

	status := a.txn.Debug().Table("alarm_action_histories").Where("id = ? ", request.Id).Delete(&request)

	if status.Error != nil {
		return status.Error
	}
	return status.Error
}

//Instance의 현재 CPU사용률을 조회한다.
func (a AlarmDao) GetAlarmHistoryList(request model.AlarmReq) (_ client.Response, errMsg model.ErrMessage) {
	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	alarmHistorySql := "select new_state, old_state, reason from alarm_state_history where alarm_id = '%s' "

	var q client.Query
	if request.TimeRange != "" {

		alarmHistorySql += " and time > now() - %s  "

		if request.State != "" {
			alarmHistorySql += " and new_state = '%s'"
		}

		q = client.Query{
			Command: fmt.Sprintf(alarmHistorySql+" order by time desc ;",
				request.AlarmId, request.TimeRange),
			Database: model.MetricDBName,
		}
	}

	model.MonitLogger.Debug("GetAlarmHistoryList Sql==>", q)
	resp, err := a.influxClient.Query(q)

	return utils.GetError().CheckError(*resp, err)
}

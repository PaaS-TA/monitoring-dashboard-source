package services

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/jinzhu/gorm"
	"github.com/monasca/golang-monascaclient/monascaclient"
	mod "github.com/monasca/golang-monascaclient/monascaclient/models"
	"kr/paasta/monitoring/iaas_new/dao"
	"kr/paasta/monitoring/iaas_new/integration"
	"kr/paasta/monitoring/iaas_new/model"
	pm "kr/paasta/monitoring/paas/model"
	"kr/paasta/monitoring/utils"
	"reflect"
	"time"
)

type AlarmStatusService struct {
	monClient    monascaclient.Client
	influxClient client.Client
	txn          *gorm.DB
}

func GetAlarmStatusService(monClient monascaclient.Client, influxClient client.Client, txn *gorm.DB) *AlarmStatusService {
	return &AlarmStatusService{
		monClient:    monClient,
		influxClient: influxClient,
		txn:          txn,
	}
}

func (a *AlarmStatusService) GetAlarmStatusCount(query mod.AlarmQuery) (pm.AlarmStatusCountResponse, error) {
	// monasca rdb 직접 조회로 수정 요청
	totalData, dbErr := dao.GetMonascaDbDao(a.txn).GetAlarmCount(*query.State)

	if dbErr != nil {
		return totalData, dbErr
	}

	return totalData, nil

	/*totalData, errCnt := integration.GetMonasca(a.monClient).GetAlarmCount(query)

	if errCnt != nil{
		return nil, errCnt
	}
	resultData := map[string]interface{}{
		model.RESULT_CNT: totalData[0][0],
	}
	return resultData, nil*/
}

func (a *AlarmStatusService) GetAlarmStatusList(query mod.AlarmQuery) (map[string]interface{}, error) {

	var allQuery mod.AlarmQuery

	allQuery.Severity = query.Severity
	allQuery.State = query.State

	totalData, errCnt := integration.GetMonasca(a.monClient).GetAlarmList(allQuery)

	if errCnt != nil {
		return nil, errCnt
	}

	alarmStatusList, err := integration.GetMonasca(a.monClient).GetAlarmList(query)

	if err != nil {
		return nil, err
	}

	var result []model.AlarmStatus

	for _, data := range alarmStatusList {

		alarmDefinition, err := dao.GetMonascaDbDao(a.txn).GetAlarmsDefinition(data.Id)
		if err != nil {
			return nil, err
		}
		alarmStatus := data
		alarmStatus.AlarmDefinitionId = alarmDefinition.AlarmDefinitionId
		alarmStatus.AlarmDefinitionName = alarmDefinition.Name
		alarmStatus.Expression = alarmDefinition.Expression
		alarmStatus.Severity = alarmDefinition.Severity

		result = append(result, alarmStatus)
	}
	resultData := map[string]interface{}{
		model.RESULT_CNT:  len(totalData),
		model.RESULT_DATA: result,
	}

	return resultData, err
}

func (a *AlarmStatusService) GetAlarmHistoryList(alarmReq model.AlarmReq) (result []model.AlarmHistory, err model.ErrMessage) {

	alarmHistoryResp, err := dao.GetAlarmDao(a.influxClient, a.txn).GetAlarmHistoryList(alarmReq)
	if err != nil {
		model.MonitLogger.Error("Error==>", err)
		return result, err
	}

	alarmHistoryList, err := utils.GetResponseConverter().InfluxConverterToMap(alarmHistoryResp)

	if err != nil {
		return result, err
	}

	for _, data := range alarmHistoryList {
		var alarmHistory model.AlarmHistory

		alarmHistory.Id = alarmReq.AlarmId
		occurDate := time.Unix(reflect.ValueOf(data["time"]).Int(), 0)
		alarmHistory.Time = occurDate.Format("2006-01-02 15:04:05")
		alarmHistory.NewState = reflect.ValueOf(data["new_state"]).String()
		alarmHistory.OldState = reflect.ValueOf(data["old_state"]).String()
		alarmHistory.Reason = reflect.ValueOf(data["reason"]).String()
		result = append(result, alarmHistory)
	}

	return result, err

}

func (a *AlarmStatusService) GetAlarmStatus(alarmId string) (result model.AlarmStatus, err error) {

	result, err = integration.GetMonasca(a.monClient).GetAlarm(alarmId)

	if err != nil {
		return result, err
	}

	alarmDefinition, err := dao.GetMonascaDbDao(a.txn).GetAlarmsDefinition(result.Id)
	if err != nil {
		return result, err
	}

	result.AlarmDefinitionId = alarmDefinition.AlarmDefinitionId
	result.AlarmDefinitionName = alarmDefinition.Name
	result.Expression = alarmDefinition.Expression
	result.Severity = alarmDefinition.Severity

	return result, nil
}

func (a *AlarmStatusService) GetAlarmHistoryActionList(alarmId string) (result []model.AlarmActionResponse, err error) {

	var alarmRequest model.AlarmActionRequest
	alarmRequest.AlarmId = alarmId

	alarmActionList, err := dao.GetAlarmDao(a.influxClient, a.txn).GetAlarmsActionHistoryList(alarmRequest)

	for _, data := range alarmActionList {
		var alarmAction model.AlarmActionResponse
		alarmAction.Id = data.Id
		alarmAction.AlarmId = data.AlarmId
		alarmAction.AlarmActionDesc = data.AlarmActionDesc
		alarmAction.RegDate = data.RegDate.Add(time.Duration(model.GmtTimeGap) * time.Hour).Format("2006-01-02 15:04:05")
		alarmAction.RegUser = data.RegUser
		alarmAction.RegDate = data.RegDate.Add(time.Duration(model.GmtTimeGap) * time.Hour).Format("2006-01-02 15:04:05")
		alarmAction.ModiUser = data.ModiUser

		result = append(result, alarmAction)
	}
	if err != nil {
		return result, err
	}

	return result, nil
}

func (h *AlarmStatusService) CreateAlarmHistoryAction(request model.AlarmActionRequest) error {

	dbErr := dao.GetAlarmDao(h.influxClient, h.txn).CreateAlarmAction(request)
	return dbErr
}

func (h *AlarmStatusService) UpdateAlarmAction(request model.AlarmActionRequest) error {

	dbErr := dao.GetAlarmDao(h.influxClient, h.txn).UpdateAlarmAction(request)
	return dbErr
}

func (h *AlarmStatusService) DeleteAlarmAction(request model.AlarmActionRequest) error {

	dbErr := dao.GetAlarmDao(h.influxClient, h.txn).DeleteAlarmAction(request)
	return dbErr
}

func (a *AlarmStatusService) GetIaasAlarmRealTimeCount() (model.AlarmRealtimeCountResponse, error) {

	var query mod.AlarmQuery
	var result model.AlarmRealtimeCountResponse

	targetState := model.ALARM_STATE_ALARM
	query.State = &targetState

	statusListResult, err := a.GetAlarmStatusList(query)
	if err != nil {
		return result, err
	}

	list := statusListResult[model.RESULT_DATA].([]model.AlarmStatus)
	for _, v := range list {
		if model.ALARM_STATE_ALARM == v.State {
			result.TotalCnt++
			switch v.Severity {
			case model.ALARM_SEVERITY_CRITICAL:
				result.CriticalCnt++
			case model.ALARM_SEVERITY_HIGH:
				result.WarningCnt++
			}
		}
	}

	return result, nil
}

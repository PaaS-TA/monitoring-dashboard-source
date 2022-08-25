package service

import (
	"github.com/jinzhu/gorm"
	"monitoring-portal/iaas_new/dao"
	"monitoring-portal/paas/model"
)

type AlarmService struct {
	txn *gorm.DB
}

func GetAlarmService(txn *gorm.DB) *AlarmService {
	return &AlarmService{
		txn: txn,
	}
}

//Service
func (h *AlarmService) GetAlarmList(request model.AlarmRequest, txn *gorm.DB) (model.AlarmPagingResponse, model.ErrMessage) {

	var alarmListPaging model.AlarmPagingResponse
	alarms, rowCnt, dbErr := dao.GetAlarmDao(h.txn).GetAlarmList(request, txn)

	if dbErr != nil {
		return alarmListPaging, dbErr
	}

	alarmListPaging.PageIndex = request.PagingReq.PageIndex
	alarmListPaging.PageItem = request.PagingReq.PageItem
	alarmListPaging.TotalCount = rowCnt
	alarmListPaging.AlarmResponse = alarms

	return alarmListPaging, nil
}

func (h *AlarmService) GetAlarmListCount(request model.AlarmRequest, txn *gorm.DB) (model.AlarmStatusCountResponse, model.ErrMessage) {

	totalData, dbErr := dao.GetAlarmDao(h.txn).GetAlarmListCount(request, txn)
	if dbErr != nil {
		return totalData, dbErr
	}

	return totalData, nil
}

func (h *AlarmService) GetAlarmResolveStatus(request model.AlarmRequest, txn *gorm.DB) ([]model.AlarmResponse, model.ErrMessage) {

	alarms, dbErr := dao.GetAlarmDao(h.txn).GetAlarmResolveStatus(request, txn)

	if dbErr != nil {
		return alarms, dbErr
	}

	return alarms, nil
}

func (h *AlarmService) GetAlarmDetail(request model.AlarmRequest, txn *gorm.DB) (model.AlarmDetailResponse, model.ErrMessage) {

	alarm, dbErr := dao.GetAlarmDao(h.txn).GetAlarmDetail(request, txn)

	if dbErr != nil {
		return alarm, dbErr
	}

	var alarmRequest model.AlarmRequest
	alarmRequest.Id = alarm.Id

	alarmAction, err := dao.GetAlarmDao(h.txn).GetAlarmsAction(alarmRequest, txn)
	if err != nil {
		return alarm, err
	}
	alarm.AlarmActionResponse = alarmAction

	return alarm, nil
}

func (h *AlarmService) UpdateAlarm(request model.AlarmRequest, txn *gorm.DB) model.ErrMessage {

	dbErr := dao.GetAlarmDao(h.txn).UpdateAlarm(request, txn)
	return dbErr
}

func (h *AlarmService) CreateAlarmAction(request model.AlarmActionRequest, txn *gorm.DB) model.ErrMessage {

	dbErr := dao.GetAlarmDao(h.txn).CreateAlarmAction(request, txn)
	return dbErr
}

func (h *AlarmService) UpdateAlarmAction(request model.AlarmActionRequest, txn *gorm.DB) model.ErrMessage {

	dbErr := dao.GetAlarmDao(h.txn).UpdateAlarmAction(request, txn)
	return dbErr
}

func (h *AlarmService) DeleteAlarmAction(request model.AlarmActionRequest, txn *gorm.DB) model.ErrMessage {

	dbErr := dao.GetAlarmDao(h.txn).DeleteAlarmAction(request, txn)
	return dbErr
}

func (h *AlarmService) GetAlarmStat(request model.AlarmStatRequest, txn *gorm.DB) (model.AlarmStatResponse, model.ErrMessage) {

	alarms, dbErr := dao.GetAlarmDao(h.txn).GetAlarmStat(request, txn)

	if dbErr != nil {
		return alarms, dbErr
	}

	return alarms, nil
}

func (h *AlarmService) GetAlarmStatGraph(request model.AlarmStatRequest) (result []map[string]interface{}, errMsg model.ErrMessage) {

	for _, v := range request.Args.([]model.AlarmStat) {
		request.Level = v.Level
		request.Origin = v.Origin
		request.Type = v.Type
		daoResult, errMsg := dao.GetAlarmDao(h.txn).GetAlarmListByPeriod(request)
		if errMsg != nil {
			return result, errMsg
		}
		stats := map[string]interface{}{
			model.RESULT_NAME:      v.Alias,
			model.RESULT_STAT_NAME: daoResult,
		}
		result = append(result, stats)
	}

	return result, errMsg
}

func (h *AlarmService) GetPaasAlarmRealTimeCount() (model.AlarmRealtimeCountResponse, model.ErrMessage) {

	var result model.AlarmRealtimeCountResponse

	list, err := dao.GetAlarmDao(h.txn).GetPaasAlarmRealTimeList()
	if err != nil {
		return result, err
	}

	result.TotalCnt = len(list)
	for _, v := range list {
		switch v.Level {
		case model.STATE_CRITICAL:
			result.CriticalCnt++
		case model.STATE_WARNING:
			result.WarningCnt++
		}
	}

	return result, nil
}

func (h *AlarmService) GetPaasAlarmRealTimeList() (model.AlarmRealtimeListResponse, model.ErrMessage) {

	var result model.AlarmRealtimeListResponse
	var list []model.AlarmResponse

	list, err := dao.GetAlarmDao(h.txn).GetPaasAlarmRealTimeList()
	if err != nil {
		return result, err
	}

	result.TotalCount = len(list)
	result.AlarmResponse = list

	return result, nil
}

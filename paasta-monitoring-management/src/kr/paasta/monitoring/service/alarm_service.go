package service

import (
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring/domain"
	"kr/paasta/monitoring/dao"
)

type AlarmService struct {
	txn   *gorm.DB
}

func GetAlarmService(txn *gorm.DB) *AlarmService {
	return &AlarmService{
		txn:   txn,
	}
}
//Service
func (h *AlarmService) GetAlarmList(request domain.AlarmRequest, txn *gorm.DB) (domain.AlarmPagingResponse, domain.ErrMessage) {

	var alarmListPaging domain.AlarmPagingResponse
	alarms, rowCnt, dbErr := dao.GetAlarmDao(h.txn).GetAlarmList(request, txn)

	if dbErr != nil{
		return alarmListPaging , dbErr
	}

	alarmListPaging.PageIndex = request.PageIndex
	alarmListPaging.PageItem  = request.PageItem
	alarmListPaging.TotalCount = rowCnt
	alarmListPaging.AlarmResponse = alarms

	return alarmListPaging, nil
}

func (h *AlarmService) GetAlarmResolveStatus(request domain.AlarmRequest, txn *gorm.DB) ([]domain.AlarmResponse, domain.ErrMessage) {

	alarms, dbErr := dao.GetAlarmDao(h.txn).GetAlarmResolveStatus(request, txn)

	if dbErr != nil{
		return alarms , dbErr
	}

	return alarms, nil
}

func (h *AlarmService) GetAlarmDetail(request domain.AlarmRequest, txn *gorm.DB) (domain.AlarmDetailResponse, domain.ErrMessage) {

	alarm, dbErr := dao.GetAlarmDao(h.txn).GetAlarmDetail(request, txn)

	if dbErr != nil{
		return alarm , dbErr
	}

	var alarmRequest domain.AlarmRequest
	alarmRequest.Id = alarm.Id

	alarmAction , err  := dao.GetAlarmDao(h.txn).GetAlarmsAction(alarmRequest, txn)
	if err != nil {
		return alarm , err
	}
	alarm.AlarmActionResponse = alarmAction

	return alarm, nil
}

func (h *AlarmService) UpdateAlarm(request domain.AlarmRequest, txn *gorm.DB) domain.ErrMessage {

	dbErr := dao.GetAlarmDao(h.txn).UpdateAlarm(request, txn)
	return dbErr
}

func (h *AlarmService) CreateAlarmAction(request domain.AlarmActionRequest, txn *gorm.DB) domain.ErrMessage {

	dbErr := dao.GetAlarmDao(h.txn).CreateAlarmAction(request, txn)
	return dbErr
}

func (h *AlarmService) UpdateAlarmAction(request domain.AlarmActionRequest, txn *gorm.DB) domain.ErrMessage {

	dbErr := dao.GetAlarmDao(h.txn).UpdateAlarmAction(request, txn)
	return dbErr
}

func (h *AlarmService) DeleteAlarmAction(request domain.AlarmActionRequest, txn *gorm.DB) domain.ErrMessage {

	dbErr := dao.GetAlarmDao(h.txn).DeleteAlarmAction(request, txn)
	return dbErr
}


func (h *AlarmService) GetAlarmStat(request domain.AlarmStatRequest, txn *gorm.DB) (domain.AlarmStatResponse, domain.ErrMessage) {

	alarms, dbErr := dao.GetAlarmDao(h.txn).GetAlarmStat(request, txn)

	if dbErr != nil{
		return alarms , dbErr
	}

	return alarms, nil
}
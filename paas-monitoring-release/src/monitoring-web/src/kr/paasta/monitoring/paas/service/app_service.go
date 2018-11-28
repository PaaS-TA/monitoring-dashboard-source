package service

import (
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring/paas/model"
	"kr/paasta/monitoring/paas/dao"
)

type AppService struct {
	txn   *gorm.DB
}

func GetAppService(txn *gorm.DB) *AppService{
	return &AppService{
		txn: 	txn,
	}
}


func (h *AppService) UpdatePaasAppAutoScalingPolicy(request model.AppAutoscalingPolicy) (res model.ResultResponse) {

	result := dao.GetAppDao(h.txn).UpdatePaasAppAutoScalingPolicy(request)

	if result != "" {
		res.Status = model.RESULT_FAIL
		res.Message = result
		return res
	} else {
		res.Status = model.RESULT_SUCCESS
		//res.Message = model.RESULT_SUCCESS
		return res
	}
}

func (h *AppService) GetPaasAppAutoScalingPolicy(request model.AppAlarmReq) (res model.AppAutoscalingPolicy, err model.ErrMessage) {

	res, err = dao.GetAppDao(h.txn).GetPaasAppAutoScalingPolicy(request)

	return res, err
}


func (h *AppService) UpdatePaasAppPolicyInfo(request model.AppAlarmPolicy) (res model.ResultResponse) {

	result := dao.GetAppDao(h.txn).UpdatePaasAppPolicyInfo(request)

	if result != "" {
		res.Status = model.RESULT_FAIL
		res.Message = result
		return res
	} else {
		res.Status = model.RESULT_SUCCESS
		//res.Message = model.RESULT_SUCCESS
		return res
	}
}


func (h *AppService) GetPaasAppPolicyInfo(request model.AppAlarmReq) (res model.AppAlarmPolicy, err model.ErrMessage) {

	res, err = dao.GetAppDao(h.txn).GetPaasAppPolicyInfo(request)

	return res, err
}

func (h *AppService) GetPaasAppAlarmList(request model.AppAlarmReq) (model.AppAlarmPagingRes, model.ErrMessage) {
	var alarmListPaging model.AppAlarmPagingRes
	alarms, rowCnt, dbErr := dao.GetAppDao(h.txn).GetPaasAppAlarmList(request)

	if dbErr != nil{
		return alarmListPaging , dbErr
	}

	alarmListPaging.PageIndex = request.PageIndex
	alarmListPaging.PageItem  = request.PageItems
	alarmListPaging.TotalCount = rowCnt
	alarmListPaging.AppAlarmList = alarms

	return alarmListPaging, nil
}


func (h *AppService) DeletePaasAppPolicy(guid string) (res model.ResultResponse) {

	result := dao.GetAppDao(h.txn).DeletePaasAppPolicy(guid)

	if result != "" {
		res.Status = model.RESULT_FAIL
		res.Message = result
		return res
	} else {
		res.Status = model.RESULT_SUCCESS
		//res.Message = model.RESULT_SUCCESS
		return res
	}
}


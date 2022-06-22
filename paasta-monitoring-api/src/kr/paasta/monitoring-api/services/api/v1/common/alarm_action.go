package common

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	dao "paasta-monitoring-api/dao/api/v1/common"
	models "paasta-monitoring-api/models/api/v1"
	"strconv"
	"time"
)

type AlarmActionService struct {
	DbInfo *gorm.DB
}


func GetAlarmActionService(DbInfo *gorm.DB) *AlarmActionService {
	return &AlarmActionService{
		DbInfo: DbInfo,
	}
}


func (ap *AlarmActionService) CreateAlarmAction(request models.AlarmActionRequest) (string, error) {
	params := models.AlarmActions {
		AlarmId : request.AlarmId,
		AlarmActionDesc: request.AlarmActionDesc,
		RegDate: time.Now(),
		RegUser: request.RegUser,
	}
	alarmParams := models.Alarms{
		Id: request.AlarmId,
	}
	alarmResult, err := dao.GetAlarmDao(ap.DbInfo).GetAlarms(alarmParams)

	if len(alarmResult) <= 0 {
		err = errors.New("Not exist alarms data.")
		return "FAILED CREATE ALARM ACTION!", err
	}

	err = dao.GetAlarmActionDao(ap.DbInfo).CreateAlarmAction(params)
	if err != nil {
		return "FAILED CREATE ALARM ACTION!", err
	}
	return "SUCCEEDED CREATE ALARM ACTION!", nil
}


func (ap *AlarmActionService) GetAlarmAction(ctx echo.Context) ([]models.AlarmActions, error) {
	alarmId, _ := strconv.Atoi(ctx.QueryParam("alarmId"))
	params := models.AlarmActions{
		AlarmId: alarmId,
		AlarmActionDesc: ctx.QueryParam("alarmActionDesc"),
	}
	results, err := dao.GetAlarmActionDao(ap.DbInfo).GetAlarmAction(params)
	if err != nil {
		return results, err
	}
	return results, nil
}


func (ap *AlarmActionService) UpdateAlarmAction(request models.AlarmActionRequest) (string, error) {
	params := models.AlarmActions {
		Id : request.Id,
		AlarmActionDesc: request.AlarmActionDesc,
		ModiDate: time.Now(),
		ModiUser: request.RegUser,
	}

	err := dao.GetAlarmActionDao(ap.DbInfo).UpdateAlarmAction(params)
	if err != nil {
		return "FAILED UPDATE ALARM ACTION!", err
	}
	return "SUCCEEDED UPDATE ALARM ACTION!", nil
}


func (ap *AlarmActionService) DeleteAlarmAction(request models.AlarmActionRequest) (string, error) {
	err := dao.GetAlarmActionDao(ap.DbInfo).DeleteAlarmAction(request)
	if err != nil {
		return "FAILED DELETE ALARM ACTION!", err
	}
	return "SUCCEEDED DELETE ALARM ACTION!", nil
}
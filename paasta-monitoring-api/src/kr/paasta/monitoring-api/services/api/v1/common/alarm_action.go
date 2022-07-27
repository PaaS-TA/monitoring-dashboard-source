package common

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	"paasta-monitoring-api/dao/api/v1/common"
	"paasta-monitoring-api/helpers"
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


func (service *AlarmActionService) CreateAlarmAction(ctx echo.Context) (string, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	var request models.AlarmActionRequest
	err := helpers.BindJsonAndCheckValid(ctx, &request)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	request.RegUser = ctx.Get("userId").(string)

	params := models.AlarmActions {
		AlarmId : request.AlarmId,
		AlarmActionDesc: request.AlarmActionDesc,
		RegDate: time.Now(),
		RegUser: request.RegUser,
	}
	alarmParams := models.Alarms{
		Id: request.AlarmId,
	}
	alarmResult, err := common.GetAlarmDao(service.DbInfo).GetAlarms(alarmParams)

	if len(alarmResult) <= 0 {
		err = errors.New("Not exist alarms data.")
		return "FAILED CREATE ALARM ACTION!", err
	}

	err = common.GetAlarmActionDao(service.DbInfo).CreateAlarmAction(params)
	if err != nil {
		return "FAILED CREATE ALARM ACTION!", err
	}
	return "SUCCEEDED CREATE ALARM ACTION!", nil
}


func (service *AlarmActionService) GetAlarmAction(ctx echo.Context) ([]models.AlarmActions, error) {
	alarmId, _ := strconv.Atoi(ctx.QueryParam("alarmId"))
	params := models.AlarmActions{
		AlarmId: alarmId,
		AlarmActionDesc: ctx.QueryParam("alarmActionDesc"),
	}
	results, err := common.GetAlarmActionDao(service.DbInfo).GetAlarmAction(params)
	if err != nil {
		return results, err
	}
	return results, nil
}


func (service *AlarmActionService) UpdateAlarmAction(ctx echo.Context) (string, *models.ApiError) {

	var request models.AlarmActionRequest
	errValid := helpers.BindJsonAndCheckValid(ctx, &request)
	if errValid != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, models.ERR_PARAM_VALIDATION, errValid.Error())
		apiError := &models.ApiError{
			OriginError: errValid,
			Message: errValid.Error(),
		}
		return "", apiError
	}
	params := models.AlarmActions {
		Id : request.Id,
		AlarmActionDesc: request.AlarmActionDesc,
		ModiDate: time.Now(),
		ModiUser: ctx.Get("userId").(string),
	}

	err := common.GetAlarmActionDao(service.DbInfo).UpdateAlarmAction(params)
	if err != nil {
		return models.FAIL_UPD_ALARM_ACTION, err
	}
	return models.SUCC_UPD_ALARM_ACTION, nil
}


func (service *AlarmActionService) DeleteAlarmAction(ctx echo.Context) (string, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	var request models.AlarmActionRequest
	err := helpers.BindJsonAndCheckValid(ctx, &request)
	if err != nil {
		logger.Error(err)
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Invalid JSON provided, please check the request JSON", err.Error())
		return "", err
	}
	errAction := common.GetAlarmActionDao(service.DbInfo).DeleteAlarmAction(request)
	if errAction != nil {
		logger.Error(errAction)
		return "FAILED DELETE ALARM ACTION!", errAction
	}
	return "SUCCEEDED DELETE ALARM ACTION!", nil
}
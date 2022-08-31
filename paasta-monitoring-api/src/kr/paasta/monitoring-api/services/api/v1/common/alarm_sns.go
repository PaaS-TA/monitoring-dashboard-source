package common

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	errors "golang.org/x/xerrors"
	"gorm.io/gorm"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	"paasta-monitoring-api/dao/api/v1/common"
	"paasta-monitoring-api/helpers"
	models "paasta-monitoring-api/models/api/v1"
	"time"
)

type AlarmSnsService struct {
	DbInfo *gorm.DB
}

func GetAlarmSnsService(DbInfo *gorm.DB) *AlarmSnsService {
	return &AlarmSnsService{
		DbInfo: DbInfo,
	}
}

func (service *AlarmSnsService) CreateAlarmSns(ctx echo.Context) (string, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	var params []models.AlarmSns
	err := helpers.BindJsonAndCheckValid(ctx, &params)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Invalid JSON provided, please check the request JSON", err.Error())
		return "", errors.New(models.ERR_PARAM_VALIDATION)
	}
	for i, _ := range params {
		params[i].RegUser = ctx.Get("userId").(string)
		params[i].RegDate = time.Now()
	}

	errAction := common.GetAlarmSnsDao(service.DbInfo).CreateAlarmSns(params)
	if errAction != nil {
		logger.Error(err)
		return "Failed to register sns account.", errAction
	}
	return "Succeeded to register sns account.", nil
}

func (service *AlarmSnsService) GetAlarmSns(ctx echo.Context) ([]models.AlarmSns, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	params := models.AlarmSns{
		OriginType: ctx.QueryParam("originType"),
		SnsType:    ctx.QueryParam("snsType"),
		SnsSendYN:  ctx.QueryParam("snsSendYn"),
	}

	results, err := common.GetAlarmSnsDao(service.DbInfo).GetAlarmSns(params)
	if err != nil {
		logger.Error(err)
		return results, err
	}
	return results, nil
}

func (service *AlarmSnsService) UpdateAlarmSns(ctx echo.Context) (string, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	params := &models.AlarmSns{}
	err := helpers.BindJsonAndCheckValid(ctx, &params)
	if err != nil {
		return "", errors.New(models.ERR_PARAM_VALIDATION)
	}
	params.ModiUser = ctx.Get("userId").(string)
	params.ModiDate = time.Now()

	errAction := common.GetAlarmSnsDao(service.DbInfo).UpdateAlarmSns(params)
	if errAction != nil {
		logger.Error(err)
		return "FAILED UPDATE SNS ACCOUNT!", errAction
	}
	return "SUCCEEDED UPDATE SNS ACCOUNT!", nil
}

func (service *AlarmSnsService) DeleteAlarmSns(ctx echo.Context) (string, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	var params models.AlarmSns
	err := helpers.BindJsonAndCheckValid(ctx, &params)
	if err != nil {
		return "", errors.New(models.ERR_PARAM_VALIDATION)
	}

	errAction := common.GetAlarmSnsDao(service.DbInfo).DeleteAlarmSns(params)
	if errAction != nil {
		logger.Error(errAction)
		return "Failed to delete sns account.", errAction
	}
	return "Succeeded to delete sns account.", nil
}

package common

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"github.com/labstack/echo/v4"
	"paasta-monitoring-api/dao/api/v1/common"
	models "paasta-monitoring-api/models/api/v1"
)

type AlarmService struct {
	DbInfo *gorm.DB
}

func GetAlarmService(DbInfo *gorm.DB) *AlarmService {
	return &AlarmService {
		DbInfo: DbInfo,
	}
}

func (ap *AlarmService) GetAlarms(ctx echo.Context) ([]models.Alarms, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	params := models.Alarms{
		OriginType: ctx.QueryParam("originType"),
		AlarmType: ctx.QueryParam("alarmType"),
		Level: ctx.QueryParam("level"),
		ResolveStatus: ctx.QueryParam("resolveStatus"),
	}

	results, err := common.GetAlarmDao(ap.DbInfo).GetAlarms(params)
	if err != nil {
		logger.Error(err)
		return results, err
	}
	return results, nil
}
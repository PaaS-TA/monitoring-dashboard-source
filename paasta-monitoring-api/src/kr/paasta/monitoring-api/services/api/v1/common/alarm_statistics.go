package common

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"github.com/labstack/echo/v4"
	"paasta-monitoring-api/dao/api/v1/common"
	models "paasta-monitoring-api/models/api/v1"
)

type AlarmStatisticsService struct {
	DbInfo *gorm.DB
}


func GetAlarmStatisticsService(DbInfo *gorm.DB) *AlarmStatisticsService {
	return &AlarmStatisticsService{
		DbInfo: DbInfo,
	}
}

func (service *AlarmStatisticsService) GetAlarmStatistics(ctx echo.Context) ([]map[string]interface{}, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	var params models.AlarmStatisticsParam
	params.OriginType = ctx.QueryParam("originType")
	params.ResourceType = ctx.QueryParam("resourceType")
	params.AliasPrefix = ctx.QueryParam("aliasPrefix")

	switch ctx.QueryParam("period") {
	case "d":
		params.Period = "DAY"
		params.TimeCriterion = "HOUR"
		params.DateFormat = "%Y-%m-%d %H:00"
	case "w":
		params.Period = "WEEK"
		params.TimeCriterion = "DAY"
		params.DateFormat = "%Y-%m-%d"
	case "m":
		params.Period = "MONTH"
		params.TimeCriterion = "DAY"
		params.DateFormat = "%Y-%m-%d"
	case "y":
		params.Period = "YEAR"
		params.TimeCriterion = "MONTH"
		params.DateFormat = "%Y-%m-01"
	}

	request := []models.AlarmStatisticsCriteriaRequest{
		{params.AliasPrefix + "Warning", "warning", params.OriginType, params.ResourceType},
		{params.AliasPrefix + "Critical", "critical", params.OriginType, params.ResourceType},
	}
	params.ExtraParams = request

	results, err := common.GetAlarmStatisticsDao(service.DbInfo).GetAlarmStatisticsForGraphByTime(params)
	if err != nil {
		logger.Error(err)
		return results, err
	}
	return results, nil
}


func (service *AlarmStatisticsService) GetAlarmStatisticsResource(ctx echo.Context) ([]map[string]interface{}, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	var params models.AlarmStatisticsParam
	params.OriginType = ctx.Param("originType")

	switch ctx.QueryParam("period") {
	case "d":
		params.Period = "DAY"
		params.TimeCriterion = "HOUR"
		params.DateFormat = "%Y-%m-%d %H:00"
	case "w":
		params.Period = "WEEK"
		params.TimeCriterion = "DAY"
		params.DateFormat = "%Y-%m-%d"
	case "m":
		params.Period = "MONTH"
		params.TimeCriterion = "DAY"
		params.DateFormat = "%Y-%m-%d"
	case "y":
		params.Period = "YEAR"
		params.TimeCriterion = "MONTH"
		params.DateFormat = "%Y-%m-01"
	}

	request := []models.AlarmStatisticsCriteriaRequest{
		{"cpu-Warning", "warning", params.OriginType, "cpu"},
		{"cpu-Critical", "critical", params.OriginType, "cpu"},
		{"memory-Warning", "warning", params.OriginType, "memory"},
		{"memory-Critical", "critical", params.OriginType, "memory"},
		{"disk-Warning", "warning", params.OriginType, "disk"},
		{"disk-Critical", "critical", params.OriginType, "disk"},
	}
	params.ExtraParams = request

	results, err := common.GetAlarmStatisticsDao(service.DbInfo).GetAlarmStatisticsForGraphByTime(params)
	if err != nil {
		logger.Error(err)
		return results, err
	}
	return results, nil
}
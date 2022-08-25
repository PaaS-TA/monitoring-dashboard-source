package common

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
		{params.OriginType + params.ResourceType + "-Warning", "warning", params.OriginType, params.ResourceType},
		{params.OriginType + params.ResourceType + "-Critical", "critical", params.OriginType, params.ResourceType},
	}
	params.ExtraParams = request

	results, err := common.GetAlarmStatisticsDao(service.DbInfo).GetAlarmStatisticsForGraphByTime(params)
	if err != nil {
		logger.Error(err)
		return results, err
	}
	return results, nil
}

func (service *AlarmStatisticsService) GetAlarmStatisticsService(ctx echo.Context) ([]map[string]interface{}, error) {
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
		{"bos-Warning", "warning", "bos", ""},
		{"bos-Critical", "critical", "bos", ""},
		{"pas-Warning", "warning", "pas", ""},
		{"pas-Critical", "critical", "pas", ""},
		{"con-Warning", "warning", "con", ""},
		{"con-Critical", "critical", "con", ""},
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
		{"cpu-Warning", "warning", "", "cpu"},
		{"cpu-Critical", "critical", "", "cpu"},
		{"memory-Warning", "warning", "", "memory"},
		{"memory-Critical", "critical", "", "memory"},
		{"disk-Warning", "warning", "", "disk"},
		{"disk-Critical", "critical", "", "disk"},
	}
	params.ExtraParams = request

	results, err := common.GetAlarmStatisticsDao(service.DbInfo).GetAlarmStatisticsForGraphByTime(params)
	if err != nil {
		logger.Error(err)
		return results, err
	}
	return results, nil
}

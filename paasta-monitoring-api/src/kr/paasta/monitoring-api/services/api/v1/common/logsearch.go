package common

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	dao "paasta-monitoring-api/dao/api/v1/common"
	"paasta-monitoring-api/helpers"
	models "paasta-monitoring-api/models/api/v1"
)

type LogSearchService struct {
	DbInfo         *gorm.DB
	InfluxDbClient models.InfluxDbClient
	BoshInfoList   []models.Bosh
}

func GetLogSearchService(DbInfo *gorm.DB, InfluxDbClient models.InfluxDbClient, BoshInfoList []models.Bosh) *LogSearchService {
	return &LogSearchService{
		DbInfo:         DbInfo,
		InfluxDbClient: InfluxDbClient,
		BoshInfoList:   BoshInfoList,
	}
}

func (l *LogSearchService) GetLogs(ctx echo.Context) ([]models.Logs, error) {
	params := models.Logs{
		UUID:       ctx.Param("uuid"),
		Keyword:    ctx.QueryParam("keyword"),
		TargetDate: ctx.QueryParam("targetDate"),
		StartTime:  ctx.QueryParam("startTime"),
		EndTime:    ctx.QueryParam("endTime"),
		Period:     ctx.QueryParam("period"),
		LogType:    ctx.QueryParam("logType"),
	}
	params = helpers.InfluxTimeSetFormatter(params)

	var results []models.Logs
	switch params.LogType {
	case "bosh":
		for _, boshInfo := range l.BoshInfoList {
			if boshInfo.UUID == params.UUID {
				response, err := dao.GetLogSearchDao(l.DbInfo, l.InfluxDbClient).GetLogs(params)
				if err != nil {
					return results, err
				}
				responseVal, _ := helpers.InfluxConverterList(response, "")

				result := params
				result.Messages = responseVal["metric"]
				results = append(results, result)
			}
		}
	case "cf":
		response, err := dao.GetLogSearchDao(l.DbInfo, l.InfluxDbClient).GetLogs(params)
		if err != nil {
			return results, err
		}
		responseVal, _ := helpers.InfluxConverterList(response, "")
		result := params
		result.Messages = responseVal["metric"]
		results = append(results, result)
	}
	return results, nil
}

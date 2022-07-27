package common

import (
	"gorm.io/gorm"
	"github.com/labstack/echo/v4"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	"paasta-monitoring-api/connections"
	common "paasta-monitoring-api/services/api/v1/common"
)

type AlarmStatisticsController struct {
	DbInfo *gorm.DB
}


func GetAlarmStatisticsController(conn connections.Connections) *AlarmStatisticsController {
	return &AlarmStatisticsController{
		DbInfo: conn.DbInfo,
	}
}


// GetAlarmStatisticsTotal
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      알람 통계 그래프를 그리기 위한 데이터 가져오기
//  @Description  알람 통계 그래프를 그리기 위한 데이터를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/ap/alarm/statistics/total [get]
func (ap *AlarmStatisticsController) GetAlarmStatistics(ctx echo.Context) error {
	results, err := common.GetAlarmStatisticsService(ap.DbInfo).GetAlarmStatistics(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get alarm statistics.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get alarm statistics.", results)
	return nil
}



// GetAlarmStatisticsResource
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      알람 통계 그래프(자원별)를 그리기 위한 데이터 가져오기
//  @Description  알람 통계 그래프(자원별)를 그리기 위한 데이터를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/ap/alarm/statistics/resource [get]
func (ap *AlarmStatisticsController) GetAlarmStatisticsResource(ctx echo.Context) error {
	results, err := common.GetAlarmStatisticsService(ap.DbInfo).GetAlarmStatisticsResource(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get alarm statistics.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get alarm statistics.", results)
	return nil
}
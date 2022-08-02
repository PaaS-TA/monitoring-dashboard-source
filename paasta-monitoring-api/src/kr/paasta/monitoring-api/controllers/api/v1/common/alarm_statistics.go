package common

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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

// GetAlarmStatistics
//  @tags         Common
//  @Summary      알람 통계 그래프를 그리기 위한 데이터 가져오기
//  @Description  알람 통계 그래프를 그리기 위한 데이터를 가져온다.
//  @Description  필수 인자를 제외한 옵션 인자는 중복하여 사용할 수 없다.
//  @Description  즉 originType, resourceType을 각각 개별로 사용해야 한다.
//  @Accept       json
//  @Produce      json
//  @Param        period        query     string  true   "Period"         enums(d, w, m, y)
//  @Param        originType    query     string  false  "Origin Type"    enums(bos, pas, con, ias)
//  @Param        resourceType  query     string  false  "Resource Type"  enums(cpu, memory, disk)
//  @Success      200           {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/alarm/stats [get]
func (ap *AlarmStatisticsController) GetAlarmStatistics(ctx echo.Context) error {
	results, err := common.GetAlarmStatisticsService(ap.DbInfo).GetAlarmStatistics(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get alarm statistics.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get alarm statistics.", results)
	return nil
}

// GetAlarmStatisticsService
//  @tags         Common
//  @Summary      알람 통계 그래프(서비스별)를 그리기 위한 데이터 가져오기
//  @Description  알람 통계 그래프(서비스별)를 그리기 위한 데이터를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        period  query     string  true  "Period"  enums(d, w, m, y)
//  @Success      200     {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/alarm/stats/service [get]
func (ap *AlarmStatisticsController) GetAlarmStatisticsService(ctx echo.Context) error {
	results, err := common.GetAlarmStatisticsService(ap.DbInfo).GetAlarmStatisticsService(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get alarm statistics.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get alarm statistics.", results)
	return nil
}

// GetAlarmStatisticsResource
//  @tags         Common
//  @Summary      알람 통계 그래프(자원별)를 그리기 위한 데이터 가져오기
//  @Description  알람 통계 그래프(자원별)를 그리기 위한 데이터를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        period  query     string  true  "Period"  enums(d, w, m, y)
//  @Success      200     {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/alarm/stats/resource [get]
func (ap *AlarmStatisticsController) GetAlarmStatisticsResource(ctx echo.Context) error {
	results, err := common.GetAlarmStatisticsService(ap.DbInfo).GetAlarmStatisticsResource(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get alarm statistics.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get alarm statistics.", results)
	return nil
}

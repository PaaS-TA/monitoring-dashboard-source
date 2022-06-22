package common

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	"paasta-monitoring-api/connections"
	models "paasta-monitoring-api/models/api/v1"
	service "paasta-monitoring-api/services/api/v1/common"
)

type AlarmController struct {
	DbInfo *gorm.DB
}

func GetAlarmController(conn connections.Connections) *AlarmController {
	return &AlarmController {
		DbInfo: conn.DbInfo,
	}
}

// GetAlarms
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      전체 알람 현황 가져오기
//  @Description  전체 알람 현황을 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.Alarms}
//  @Router       /api/v1/ap/alarm/status [get]
func (ap *AlarmController) GetAlarms(ctx echo.Context) error {
	params := models.Alarms{
		OriginType: ctx.QueryParam("origin_type"),
		AlarmType: ctx.QueryParam("alarm_type"),
		Level: ctx.QueryParam("level"),
		ResolveStatus: ctx.QueryParam("resolve_status"),
	}

	results, err := service.GetAlarmService(ap.DbInfo).GetAlarms(params)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get alarm status.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get alarms status.", results)
	return nil
}
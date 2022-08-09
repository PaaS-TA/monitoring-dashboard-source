package common

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	"paasta-monitoring-api/connections"
	service "paasta-monitoring-api/services/api/v1/common"
)

type AlarmController struct {
	DbInfo *gorm.DB
}

func GetAlarmController(conn connections.Connections) *AlarmController {
	return &AlarmController{
		DbInfo: conn.DbInfo,
	}
}

// GetAlarms
//  @Tags         Common
//  @Summary      알람 현황 가져오기
//  @Description  알람 현황을 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        originType     query     string  false  "Origin Type"  enums(bos, pas, con, ias)
//  @Param        alarmType      query     string  false  "Alarm Type"   enums(cpu, memory, disk, fail)
//  @Param        level          query     string  false  "Level"        enums(warning, critical, fail)
//  @Param        resolveStatus  query     string  false  "Resolve Status"
//  @Success      200            {object}  apiHelpers.BasicResponseForm{responseInfo=v1.Alarms}
//  @Router       /api/v1/alarm [get]
func (ap *AlarmController) GetAlarms(ctx echo.Context) error {
	results, err := service.GetAlarmService(ap.DbInfo).GetAlarms(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get alarm status.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get alarms status.", results)
	}
	return nil
}

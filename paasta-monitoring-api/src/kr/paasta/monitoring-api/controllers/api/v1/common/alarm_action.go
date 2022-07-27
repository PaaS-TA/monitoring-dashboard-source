package common

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	"paasta-monitoring-api/connections"
	models "paasta-monitoring-api/models/api/v1"
	service "paasta-monitoring-api/services/api/v1/common"
)

type AlarmActionController struct {
	DbInfo *gorm.DB
}

func GetAlarmActionController(conn connections.Connections) *AlarmActionController {
	return &AlarmActionController{
		DbInfo: conn.DbInfo,
	}
}


// CreateAlarmAction
//  @Tags         Common
//  @Summary      알람 조치 내용 신규 작성하기
//  @Description  알람에 대한 조치 내용을 신규 작성한다.
//  @Accept       json
//  @Produce      json
//  @Param        AlarmActions  body      v1.AlarmActions  true  "신규 작성할 알람 정보를 주입한다."
//  @Success      200                 {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/alarm/action [post]
func (controller *AlarmActionController) CreateAlarmAction(ctx echo.Context) error {
	results, err := service.GetAlarmActionService(controller.DbInfo).CreateAlarmAction(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to create alarm actions.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to create alarm actions.", results)
	}
	return nil
}


// GetAlarmAction
//  @Tags         Common
//  @Summary      알람 조치 내용 가져오기
//  @Description  알람에 대한 조치 내용을 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        alarmId          query     int     false  "Alarm ID"
//  @Param        alarmActionDesc  query     string  false  "Alarm Action Desc"
//  @Success      200              {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/alarm/action [get]
func (controller *AlarmActionController) GetAlarmAction(ctx echo.Context) error {
	results, err := service.GetAlarmActionService(controller.DbInfo).GetAlarmAction(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get alarm action.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get alarm action.", results)
	}
	return nil
}


// UpdateAlarmAction
//  @Tags         Common
//  @Summary      알람 조치 내용 수정하기
//  @Description  알람에 대한 조치 내용을 수정한다.
//  @Accept       json
//  @Produce      json
//  @Param        AlarmActionRequest  body      v1.AlarmActionRequest  true  "수정할 알람 정보(ID)를  주입한다."
//  @Success      200                 {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/alarm/action [patch]
func (controller *AlarmActionController) UpdateAlarmAction(ctx echo.Context) error {
	results, err := service.GetAlarmActionService(controller.DbInfo).UpdateAlarmAction(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to update alarm action.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, results, nil)
	}
	return nil
}


// DeleteAlarmAction
//  @Tags         Common
//  @Summary      알람에 대한 조치 내용 삭제하기
//  @Description  알람에 대한 조치 내용을 삭제한다.
//  @Accept       json
//  @Produce      json
//  @Param        AlarmActionRequest  body      v1.AlarmActionRequest  true  "삭제할 알람 정보(ID)를  주입한다."
//  @Success      200                 {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/alarm/action [delete]
func (controller *AlarmActionController) DeleteAlarmAction(ctx echo.Context) error {
	results, err := service.GetAlarmActionService(controller.DbInfo).DeleteAlarmAction(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, models.FAIL_DEL_ALARM_ACTION, err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, models.SUCC_DEL_ALARM_ACTION, results)
	}
	return nil
}
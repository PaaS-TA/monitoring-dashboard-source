package common

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	"paasta-monitoring-api/connections"
	"paasta-monitoring-api/helpers"
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
//  @Summary      알람에 대한 조치 내용 작성하기
//  @Description  알람에 대한 조치 내용을 작성한다.
//  @Accept       json
//  @Produce      json
//  @Param        AlarmActionRequest  body      v1.AlarmActionRequest  true  "새로 작성할 알람 정보를 주입한다."
//  @Success      200                 {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/alarm/action [post]
func (controller *AlarmActionController) CreateAlarmAction(ctx echo.Context) error {
	var request models.AlarmActionRequest
	err := helpers.BindJsonAndCheckValid(ctx, &request)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Invalid JSON provided, please check the request JSON", err.Error())
		return err
	}

	request.RegUser = ctx.Get("userId").(string)
	results, err := service.GetAlarmActionService(controller.DbInfo).CreateAlarmAction(request)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to create alarm actions.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to create alarm actions.", results)
	return nil
}

// GetAlarmAction
//  @Tags         Common
//  @Summary      알람에 대한 조치 내용 가져오기
//  @Description  알람에 대한 조치 내용을 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/alarm/action [get]
func (ap *AlarmActionController) GetAlarmAction(ctx echo.Context) error {
	results, err := service.GetAlarmActionService(ap.DbInfo).GetAlarmAction(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get alarm action.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get alarm action.", results)
	return nil
}

// UpdateAlarmAction
//  @Tags         Common
//  @Summary      알람에 대한 조치 내용 수정하기
//  @Description  알람에 대한 조치 내용을 수정한다.
//  @Accept       json
//  @Produce      json
//  @Param        AlarmActionRequest  body      v1.AlarmActionRequest  true  "수정할 알람 정보를 주입한다."
//  @Success      200                 {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/alarm/action [patch]
func (controller *AlarmActionController) UpdateAlarmAction(ctx echo.Context) error {
	var request models.AlarmActionRequest
	err := helpers.BindJsonAndCheckValid(ctx, &request)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Invalid JSON provided, please check the request JSON", err.Error())
		return err
	}

	request.RegUser = ctx.Get("userId").(string)
	results, err := service.GetAlarmActionService(controller.DbInfo).UpdateAlarmAction(request)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to update alarm action.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to update alarm action.", results)
	return nil
}

// DeleteAlarmAction
//  @Tags         Common
//  @Summary      알람에 대한 조치 내용 삭제하기
//  @Description  알람에 대한 조치 내용을 삭제한다.
//  @Accept       json
//  @Produce      json
//  @Param        AlarmActionRequest  body      v1.AlarmActionRequest  true  "삭제할 알람 정보(Id)를  주입한다."
//  @Success      200                 {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/alarm/action [delete]
func (ap *AlarmActionController) DeleteAlarmAction(c echo.Context) error {
	var request models.AlarmActionRequest
	err := helpers.BindJsonAndCheckValid(c, &request)
	if err != nil {
		apiHelpers.Respond(c, http.StatusBadRequest, "Invalid JSON provided, please check the request JSON", err.Error())
		return err
	}

	results, err := service.GetAlarmActionService(ap.DbInfo).DeleteAlarmAction(request)
	if err != nil {
		apiHelpers.Respond(c, http.StatusBadRequest, "Failed to delete alarm action.", err.Error())
		return err
	}

	apiHelpers.Respond(c, http.StatusOK, "Succeeded to delete alarm action.", results)
	return nil
}

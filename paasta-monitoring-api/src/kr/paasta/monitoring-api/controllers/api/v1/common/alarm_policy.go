package common

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	"paasta-monitoring-api/connections"
	"paasta-monitoring-api/helpers"
	models "paasta-monitoring-api/models/api/v1"
	common "paasta-monitoring-api/services/api/v1/common"
)

type AlarmPolicyController struct {
	DbInfo *gorm.DB
}

func GetAlarmPolicyController(conn connections.Connections) *AlarmPolicyController {
	return &AlarmPolicyController{
		DbInfo: conn.DbInfo,
	}
}

// GetAlarmPolicy
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      전체 알람 정책 가져오기
//  @Description  전체 알람 정책을 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.AlarmPolicies}
//  @Router       /api/v1/alarm/policy [get]
func (controller *AlarmPolicyController) GetAlarmPolicy(ctx echo.Context) error {
	results, err := common.GetAlarmPolicyService(controller.DbInfo).GetAlarmPolicy(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get alarm policy.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get alarms policy.", results)
	return nil
}

// UpdateAlarmPolicy
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      전체 알람 정책 업데이트하기
//  @Description  전체 알람 정책을 업데이트 한다.
//  @Accept       json
//  @Produce      json
//  @Param        AlarmPolicyRequest  body      v1.AlarmPolicyRequest  true  "알람 정책을 변경하기 위한 정보를 주입한다."
//  @Success      200                {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/alarm/policy [put]
func (controller *AlarmPolicyController) UpdateAlarmPolicy(ctx echo.Context) error {
	var request []models.AlarmPolicyRequest
	err := helpers.BindJsonAndCheckValid(ctx, &request)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Invalid JSON provided, please check the request JSON", err.Error())
		return err
	}

	results, err := common.GetAlarmPolicyService(controller.DbInfo).UpdateAlarmPolicy(request)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to update alarm policy.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to update alarms policy.", results)
	return nil
}

// UpdateAlarmTarget
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      전체 알람 타겟 업데이트하기
//  @Description  전체 알람 타겟을 업데이트 한다.
//  @Accept       json
//  @Produce      json
//  @Param        AlarmTargetRequest  body      v1.AlarmPolicyRequest  true  "알람 타겟을 변경하기 위한 정보를 주입한다."
//  @Success      200                {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/alarm/target [put]
func (controller *AlarmPolicyController) UpdateAlarmTarget(ctx echo.Context) error {
	var request []models.AlarmTargetRequest
	err := helpers.BindJsonAndCheckValid(ctx, &request)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Invalid JSON provided, please check the request JSON", err.Error())
		return err
	}

	results, err := common.GetAlarmPolicyService(controller.DbInfo).UpdateAlarmTarget(request)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to update alarm target.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to update alarms target.", results)
	return nil
}

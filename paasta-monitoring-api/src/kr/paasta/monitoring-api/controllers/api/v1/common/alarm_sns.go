package common

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	"paasta-monitoring-api/connections"
	"paasta-monitoring-api/helpers"
	models "paasta-monitoring-api/models/api/v1"
	service "paasta-monitoring-api/services/api/v1/common"
)

type AlarmSnsController struct {
	DbInfo *gorm.DB
}

func GetAlarmSnsController(conn connections.Connections) *AlarmSnsController {
	return &AlarmSnsController {
		DbInfo: conn.DbInfo,
	}
}


// CreateAlarmSns
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      알람 받을 SNS 계정 등록하기
//  @Description  알람 받을 SNS 계정을 등록한다.
//  @Accept       json
//  @Produce      json
//  @Param        SnsAccountRequest  body      v1.SnsAccountRequest  true  "알람 받을 SNS 계정 정보를 주입한다."
//  @Success      200                {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/ap/alarm/sns [post]
func (controller *AlarmSnsController) CreateAlarmSns(c echo.Context) error {
	var request models.SnsAccountRequest
	err := helpers.BindJsonAndCheckValid(c, &request)
	if err != nil {
		apiHelpers.Respond(c, http.StatusBadRequest, "Invalid JSON provided, please check the request JSON", err.Error())
		return err
	}

	results, err := service.GetAlarmSnsService(controller.DbInfo).CreateAlarmSns(request)
	if err != nil {
		apiHelpers.Respond(c, http.StatusBadRequest, "Failed to register sns account.", err.Error())
		return err
	}

	apiHelpers.Respond(c, http.StatusOK, "Succeeded to register sns account.", results)
	return nil
}


// GetAlarmSns
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      알람 받는 SNS 계정 가져오기
//  @Description  알람 받는 SNS 계정 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.AlarmSns}
//  @Router       /api/v1/ap/alarm/sns [get]
func (controller *AlarmSnsController) GetAlarmSns(c echo.Context) error {
	params := models.AlarmSns{
		OriginType: c.QueryParam("origin_type"),
		SnsType: c.QueryParam("sns_type"),
		SnsSendYN: c.QueryParam("sns_send_yn"),
	}

	results, err := service.GetAlarmSnsService(controller.DbInfo).GetAlarmSns(params)
	if err != nil {
		apiHelpers.Respond(c, http.StatusBadRequest, "Failed to get sns alarm list.", err.Error())
		return err
	}

	apiHelpers.Respond(c, http.StatusOK, "Succeeded to get sns alarm list.", results)
	return nil
}


// UpdateAlarmSns
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      알람 받을 SNS 계정 수정하기
//  @Description  알람 받을 SNS 계정 정보를 수정한다.
//  @Accept       json
//  @Produce      json
//  @Param        SnsAccountRequest  body      v1.SnsAccountRequest  true  "수정할 SNS 계정 정보를 주입한다."
//  @Success      200  {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/ap/alarm/sns [put]
func (controller *AlarmSnsController) UpdateAlarmSns(c echo.Context) error {
	var request models.SnsAccountRequest
	err := helpers.BindJsonAndCheckValid(c, &request)
	if err != nil {
		apiHelpers.Respond(c, http.StatusBadRequest, "Invalid JSON provided, please check the request JSON", err.Error())
		return err
	}

	results, err := service.GetAlarmSnsService(controller.DbInfo).UpdateAlarmSns(request)
	if err != nil {
		apiHelpers.Respond(c, http.StatusBadRequest, "Failed to update sns account.", err.Error())
		return err
	}

	apiHelpers.Respond(c, http.StatusOK, "Succeeded to update sns account.", results)
	return nil
}


// DeleteAlarmSns
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      알람 받는 SNS 계정 삭제하기
//  @Description  알람 받는 SNS 계정을 삭제한다.
//  @Accept       json
//  @Produce      json
//  @Param        SnsAccountRequest  body      v1.SnsAccountRequest  true  "삭제할 SNS 계정을 정보(ChannelId)를  주입한다."
//  @Success      200  {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/ap/alarm/sns [delete]
func (controller *AlarmSnsController) DeleteAlarmSns(c echo.Context) error {
	var request models.SnsAccountRequest
	err := helpers.BindJsonAndCheckValid(c, &request)
	if err != nil {
		apiHelpers.Respond(c, http.StatusBadRequest, "Invalid JSON provided, please check the request JSON", err.Error())
		return err
	}

	results, err := service.GetAlarmSnsService(controller.DbInfo).DeleteAlarmSns(request)
	if err != nil {
		apiHelpers.Respond(c, http.StatusBadRequest, "Failed to delete sns account.", err.Error())
		return err
	}

	apiHelpers.Respond(c, http.StatusOK, "Succeeded to delete sns account.", results)
	return nil
}
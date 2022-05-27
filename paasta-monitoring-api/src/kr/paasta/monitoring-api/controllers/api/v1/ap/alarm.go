package ap

import (
	"GoEchoProject/apiHelpers"
	"GoEchoProject/connections"
	"GoEchoProject/helpers"
	models "GoEchoProject/models/api/v1"
	AP "GoEchoProject/services/api/v1/ap"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ApAlarmController struct {
	DbInfo *gorm.DB
}

func GetApAlarmController(conn connections.Connections) *ApAlarmController {
	return &ApAlarmController{
		DbInfo: conn.DbInfo,
	}
}

// GetAlarmStatus
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      전체 알람 현황 가져오기
//  @Description  전체 알람 현황을 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.Alarms}
//  @Router       /api/v1/ap/alarm/status [get]
func (ap *ApAlarmController) GetAlarmStatus(c echo.Context) error {
	results, err := AP.GetApAlarmService(ap.DbInfo).GetAlarmStatus()
	if err != nil {
		apiHelpers.Respond(c, http.StatusBadRequest, "Failed to get alarm status.", err.Error())
		return err
	}

	apiHelpers.Respond(c, http.StatusOK, "Succeeded to get alarms status.", results)
	return nil
}

// GetAlarmPolicy
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      전체 알람 정책 가져오기
//  @Description  전체 알람 정책을 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.AlarmPolicies}
//  @Router       /api/v1/ap/alarm/policy [get]
func (ap *ApAlarmController) GetAlarmPolicy(c echo.Context) error {
	results, err := AP.GetApAlarmService(ap.DbInfo).GetAlarmPolicy()
	if err != nil {
		apiHelpers.Respond(c, http.StatusBadRequest, "Failed to get alarm policy.", err.Error())
		return err
	}

	apiHelpers.Respond(c, http.StatusOK, "Succeeded to get alarms policy.", results)
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
//  @Success      200                 {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/ap/alarm/policy [put]
func (ap *ApAlarmController) UpdateAlarmPolicy(c echo.Context) error {
	var request []models.AlarmPolicyRequest
	err := helpers.BindRequestAndCheckValid(c, &request)
	if err != nil {
		apiHelpers.Respond(c, http.StatusBadRequest, "Invalid JSON provided, please check the request JSON", err.Error())
		return err
	}

	results, err := AP.GetApAlarmService(ap.DbInfo).UpdateAlarmPolicy(request)
	if err != nil {
		apiHelpers.Respond(c, http.StatusBadRequest, "Failed to update alarm policy.", err.Error())
		return err
	}

	apiHelpers.Respond(c, http.StatusOK, "Succeeded to update alarms policy.", results)
	return nil
}

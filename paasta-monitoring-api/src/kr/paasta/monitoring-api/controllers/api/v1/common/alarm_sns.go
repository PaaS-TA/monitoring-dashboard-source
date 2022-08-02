package common

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	"paasta-monitoring-api/connections"
	service "paasta-monitoring-api/services/api/v1/common"
)

type AlarmSnsController struct {
	DbInfo *gorm.DB
}

func GetAlarmSnsController(conn connections.Connections) *AlarmSnsController {
	return &AlarmSnsController{
		DbInfo: conn.DbInfo,
	}
}

// CreateAlarmSns
//  @tags         Common
//  @Summary      알람 받을 SNS 계정 등록하기
//  @Description  알람 받을 SNS 계정을 등록한다.
//  @Accept       json
//  @Produce      json
//  @Param        AlarmSns  body      []v1.AlarmSns  true  "알람 받을 SNS 계정 정보를 주입한다."
//  @Success      200       {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/alarm/sns [post]
func (controller *AlarmSnsController) CreateAlarmSns(ctx echo.Context) error {
	results, err := service.GetAlarmSnsService(controller.DbInfo).CreateAlarmSns(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to register sns account.", err.Error())
		return err
	}
	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to register sns account.", results)
	return nil
}

// GetAlarmSns
//  @tags         Common
//  @Summary      알람 받는 SNS 계정 가져오기
//  @Description  알람 받는 SNS 계정 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        originType  query     string  false  "Origin Type"  enums(all, bos, pas, con, ias)
//  @param        snsType     query     string  false  "SNS Type"     default(telegram)
//  @Param        snsSendYn   query     string  false  "SNS Send YN"  enum(Y, N)
//  @Success      200         {object}  apiHelpers.BasicResponseForm{responseInfo=v1.AlarmSns}
//  @Router       /api/v1/alarm/sns [get]
func (controller *AlarmSnsController) GetAlarmSns(ctx echo.Context) error {
	results, err := service.GetAlarmSnsService(controller.DbInfo).GetAlarmSns(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get sns alarm list.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get sns alarm list.", results)
	return nil
}

// UpdateAlarmSns
//  @tags         Common
//  @Summary      알람 받을 SNS 계정 수정하기
//  @Description  알람 받을 SNS 계정 정보를 수정한다.
//  @Accept       json
//  @Produce      json
//  @Param        AlarmSns  body      v1.AlarmSns  true  "수정할 SNS 계정 정보(channelId)를  주입한다."
//  @Success      200       {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/alarm/sns [put]
func (controller *AlarmSnsController) UpdateAlarmSns(ctx echo.Context) error {
	results, err := service.GetAlarmSnsService(controller.DbInfo).UpdateAlarmSns(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to update sns account.", err.Error())
		return err
	}
	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to update sns account.", results)
	return nil
}

// DeleteAlarmSns
//  @tags         Common
//  @Summary      알람 받는 SNS 계정 삭제하기
//  @Description  알람 받는 SNS 계정을 삭제한다.
//  @Accept       json
//  @Produce      json
//  @Param        AlarmSns  body      v1.AlarmSns  true  "삭제할 SNS 계정 정보(channelId)를  주입한다."
//  @Success      200       {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/alarm/sns [delete]
func (controller *AlarmSnsController) DeleteAlarmSns(ctx echo.Context) error {
	results, err := service.GetAlarmSnsService(controller.DbInfo).DeleteAlarmSns(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to delete sns account.", err.Error())
		return err
	}
	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to delete sns account.", results)
	return nil
}

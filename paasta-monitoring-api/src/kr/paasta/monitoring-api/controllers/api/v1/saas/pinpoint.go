package saas

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	models "paasta-monitoring-api/models/api/v1"
	service "paasta-monitoring-api/services/api/v1/saas"
)

type PinpointController struct {
	SaaS models.SaaS
}

func GetPinpointController(saas models.SaaS) *PinpointController {
	return &PinpointController{
		SaaS: saas,
	}
}

// GetAgentList
//  @Tags         SaaS
//  @Summary      핀포인트 모니터링 에이전트 리스트 가져오기
//  @Description  핀포인트 모니터링 에이전트 리스트를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/saas/pinpoint/getAgentList [get]
func (controller *PinpointController) GetAgentList(ctx echo.Context) error {
	result, err := service.GetPinpointService(controller.SaaS).GetAgentList(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Pinpoint Agent List.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Pinpoint Agent List.", result)
	}
	return nil
}

// GetAgentStat
//  @Tags         SaaS
//  @Summary      핀포인트 모니터링 에이전트 통계 가져오기
//  @Description  핀포인트 모니터링 에이전트 통계를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        chartType  path      string  true  "차트 타입을 주입한다."                         enums(cpuLoad, jvmGc, activeTrace, responseTime)
//  @Param        agentId    query     string  true  "에이전트 아이디를 주입한다."                      example(3149075187)
//  @Param        period     query     string  true  "현재를 기준으로 조회하고 싶은 기간을 주입한다. 단위는 분(m),  시(h),  일(d)  단위를  사용한다."        example(1m)
//  @Success      200        {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/saas/pinpoint/{chartType}/getAgentStat [get]
func (controller *PinpointController) GetAgentStat(ctx echo.Context) error {
	result, err := service.GetPinpointService(controller.SaaS).GetAgentStat(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Pinpoint Agent Stat.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Pinpoint Agent Stat", result)
	}
	return nil
}

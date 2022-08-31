package saas

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	models "paasta-monitoring-api/models/api/v1"
	service "paasta-monitoring-api/services/api/v1/saas"
)

type SaasController struct {
	SaaS models.SaaS
}

func GetSaasController(saas models.SaaS) *SaasController {
	return &SaasController{
		SaaS: saas,
	}
}

// GetApplicationStatus
//  @Tags         SaaS
//  @Summary      애플리케이션(핀포인트 에이전트) 스테이터스 가져오기
//  @Description  애플리케이션(핀포인트 에이전트) 스테이터스를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200     {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/saas/app/status [get]
func (controller *SaasController) GetApplicationStatus(ctx echo.Context) error {
	result, err := service.GetSaasService(controller.SaaS).GetApplicationStatus(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Application Status.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Application Status.", result)
	}
	return nil
}

// GetApplicationUsage
//  @Tags         SaaS
//  @Summary      애플리케이션(핀포인트 에이전트) 통합 사용량 정보 가져오기
//  @Description  애플리케이션(핀포인트 에이전트) 통합 사용량 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        period  query     string  false  "현재를 기준으로 조회하고 싶은 기간을 주입한다. 단위는 분(m),  시(h),  일(d)  단위를  사용한다."        example(1m)
//  @Success      200     {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/saas/app/usage [get]
func (controller *SaasController) GetApplicationUsage(ctx echo.Context) error {
	result, err := service.GetSaasService(controller.SaaS).GetApplicationUsage(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Application Usage.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Application Usage", result)
	}
	return nil
}

// GetApplicationUsageList
//  @Tags         SaaS
//  @Summary      애플리케이션(핀포인트 에이전트) 개별 사용량 정보 가져오기
//  @Description  애플리케이션(핀포인트 에이전트) 개별 사용량 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        period  query     string  false  "현재를 기준으로 조회하고 싶은 기간을 주입한다. 단위는 분(m),  시(h),  일(d)  단위를  사용한다."        example(1m)
//  @Success      200  {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/saas/app/usage/list [get]
func (controller *SaasController) GetApplicationUsageList(ctx echo.Context) error {
	result, err := service.GetSaasService(controller.SaaS).GetApplicationUsageList(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Application Usage List.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Application Usage List", result)
	}
	return nil
}

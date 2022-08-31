package iaas

import (
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	service "paasta-monitoring-api/services/api/v1/iaas"

	"github.com/gophercloud/gophercloud"
	"github.com/labstack/echo/v4"
)

type OpenstackController struct {
	OpenstackProvider *gophercloud.ProviderClient
}

func GetOpenstackController(openstackProvider *gophercloud.ProviderClient) *OpenstackController {
	return &OpenstackController{
		OpenstackProvider: openstackProvider,
	}
}

// GetHypervisorStatistics
//  @tags         IaaS
//  @Summary      오픈스택 하이퍼바이저 통계 가져오기
//  @Description  오픈스택 하이퍼바이저 통계를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/iaas/hypervisor/stats [get]
func (controller *OpenstackController) GetHypervisorStatistics(ctx echo.Context) error {
	results, err := service.GetOpenstackService(controller.OpenstackProvider).GetHypervisorStatistics(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor Statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Hypervisor Statistics.", results)
	}
	return nil
}

// GetHypervisorList
//  @tags         IaaS
//  @Summary      오픈스택 하이퍼바이저 리스트 가져오기
//  @Description  오픈스택 하이퍼바이저 리스트를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/iaas/hypervisor/list [get]
func (controller *OpenstackController) GetHypervisorList(ctx echo.Context) error {
	results, err := service.GetOpenstackService(controller.OpenstackProvider).GetHypervisorList(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor List.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Hypervisor List.", results)
	}
	return nil
}

// GetProjectList
//  @tags         IaaS
//  @Summary      오픈스택 프로젝트 리스트 가져오기
//  @Description  오픈스택 프로젝트 리스트를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/iaas/project/list [get]
func (controller *OpenstackController) GetProjectList(ctx echo.Context) error {
	results, err := service.GetOpenstackService(controller.OpenstackProvider).GetProjectList(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Project List.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Project List.", results)
	}
	return nil
}

// GetProjectUsage
//  @tags         IaaS
//  @Summary      오픈스택 프로젝트 사용량 가져오기
//  @Description  오픈스택 프로젝트 사용량을 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=[]usage.TenantUsage}
//  @Router       /api/v1/iaas/instance/usage/list [get]
func (controller *OpenstackController) GetProjectUsage(ctx echo.Context) error {
	results, err := service.GetOpenstackService(controller.OpenstackProvider).RetrieveTenantUsage(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Project Usage.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Project Usage.", results)
	}
	return nil
}

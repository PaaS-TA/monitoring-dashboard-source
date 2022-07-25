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


func (controller *OpenstackController) GetHypervisorStatistics(ctx echo.Context) error {
	results, err := service.GetOpenstackService(controller.OpenstackProvider).GetHypervisorStatistics(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}
	return nil
}


/**
	하이퍼바이저 목록 (상세 정보 포함) 조회
*/
func (controller *OpenstackController) GetHypervisorList(ctx echo.Context) error {
	results, err := service.GetOpenstackService(controller.OpenstackProvider).GetHypervisorList(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}
	return nil
}


/**
	프로젝트 목록(만) 조회
*/
func (controller *OpenstackController) GetProjectList(ctx echo.Context) error {
	results, err := service.GetOpenstackService(controller.OpenstackProvider).GetProjectList(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}
	return nil
}


/**
프로젝트(테넌트) 목록과 usage 정보를 조회
	- 프로젝트에 속한 인스턴스 목록과 usage 조회도 가능하나 현재는 비활성화 되어 있음
*/
func (controller *OpenstackController) GetProjectUsage(ctx echo.Context) error {
	results, err := service.GetOpenstackService(controller.OpenstackProvider).RetrieveTenantUsage(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}
	return nil
}
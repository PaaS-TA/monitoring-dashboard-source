package iaas

import (
	"fmt"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	service "paasta-monitoring-api/services/api/v1/iaas"

	"github.com/labstack/echo/v4"
	"github.com/gophercloud/gophercloud"
)

type (
	OpenstackController struct {
		OpenstackProvider *gophercloud.ProviderClient
	}
)


func GetOpenstackController(openstackProvider *gophercloud.ProviderClient) *OpenstackController {
	return &OpenstackController{
		OpenstackProvider: openstackProvider,
	}
}


func (controller *OpenstackController) GetHypervisorStatistics(ctx echo.Context) error {
	fmt.Println(controller.OpenstackProvider.TokenID)
	results, err := service.GetOpenstackService(controller.OpenstackProvider).GetHypervisorStatistics()
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
func (osService *OpenstackController) GetHypervisorList(ctx echo.Context) error {
	results, err := service.GetOpenstackService(osService.OpenstackProvider).GetHypervisorList()

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
func (osService *OpenstackController) GetProjectList(ctx echo.Context) error {
	serverParams := make(map[string]interface{}, 0)
	results, err := service.GetOpenstackService(osService.OpenstackProvider).GetProjectList(serverParams)

	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}
	return nil
}
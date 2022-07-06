package caas

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	models "paasta-monitoring-api/models/api/v1"
	service "paasta-monitoring-api/services/api/v1/caas"
)

type WorkloadController struct {
	CaasConfig models.CaasConfig
}

func GetWorkloadController(config models.CaasConfig) *WorkloadController {
	return &WorkloadController {
		CaasConfig: config,
	}
}

func (controller *WorkloadController) GetWorkloadStatus(ctx echo.Context) error {
	results, err := service.GetWorkloadService(controller.CaasConfig).GetWorkloadStatus()
	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}
	return nil
}


func (controller *WorkloadController) GetWorkloadList(ctx echo.Context) error {
	results, err := service.GetWorkloadService(controller.CaasConfig).GetWorkloadList()
	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}
	return nil
}


func (controller *WorkloadController) GetWorkloadDetailMetrics(ctx echo.Context) error {
	results, err := service.GetWorkloadService(controller.CaasConfig).GetWorkloadDetailMetrics(ctx.QueryParam("workload"))
	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}
	return nil
}


func (controller *WorkloadController) GetWorkloadContainerList(ctx echo.Context) error {
	results, err := service.GetWorkloadService(controller.CaasConfig).GetWorkloadContainerList(ctx.QueryParam("workload"))
	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}
	return nil
}
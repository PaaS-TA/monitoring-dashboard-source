package caas

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	models "paasta-monitoring-api/models/api/v1"
	service "paasta-monitoring-api/services/api/v1/caas"
)

type WorkloadController struct {
	CaaS models.CaaS
}

func GetWorkloadController(config models.CaaS) *WorkloadController {
	return &WorkloadController {
		CaaS: config,
	}
}

func (controller *WorkloadController) GetWorkloadStatus(ctx echo.Context) error {
	results, err := service.GetWorkloadService(controller.CaaS).GetWorkloadStatus(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}
	return nil
}


func (controller *WorkloadController) GetWorkloadList(ctx echo.Context) error {
	results, err := service.GetWorkloadService(controller.CaaS).GetWorkloadList()
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}
	return nil
}


func (controller *WorkloadController) GetWorkloadDetailMetrics(ctx echo.Context) error {
	results, err := service.GetWorkloadService(controller.CaaS).GetWorkloadDetailMetrics(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}
	return nil
}


func (controller *WorkloadController) GetWorkloadContainerList(ctx echo.Context) error {
	results, err := service.GetWorkloadService(controller.CaaS).GetWorkloadContainerList(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}
	return nil
}


func (controller *WorkloadController) GetContainerMetrics(ctx echo.Context) error {
	results, err := service.GetWorkloadService(controller.CaaS).GetContainerMetrics(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}
	return nil
}


func (controller *WorkloadController) GetContainerLog(ctx echo.Context) error {
	results, err := service.GetWorkloadService(controller.CaaS).GetContainerLog(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}
	return nil
}
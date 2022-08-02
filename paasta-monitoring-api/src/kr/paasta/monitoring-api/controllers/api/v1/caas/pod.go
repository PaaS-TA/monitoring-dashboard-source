package caas

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	models "paasta-monitoring-api/models/api/v1"
	service "paasta-monitoring-api/services/api/v1/caas"
)

type PodController struct {
	CaaS models.CaaS
}

func GetPodController(config models.CaaS) *PodController{
	return &PodController{
		CaaS: config,
	}
}

func (controller *PodController) GetPodStatus(ctx echo.Context) error {
	results, err := service.GetPodService(controller.CaaS).GetPodStatus(ctx)
	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}
	return nil
}


func (controller *PodController) GetPodList(ctx echo.Context) error {
	results, err := service.GetPodService(controller.CaaS).GetPodList()
	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}
	return nil
}


func (controller *PodController) GetPodDetailMetrics(ctx echo.Context) error {
	pod := ctx.QueryParam("pod")
	results, err := service.GetPodService(controller.CaaS).GetPodDetailMetrics(pod)
	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}
	return nil
}


func (controller *PodController) GetPodContainerList(ctx echo.Context) error {
	pod := ctx.QueryParam("pod")
	results, err := service.GetPodService(controller.CaaS).GetPodContainerList(pod)
	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}
	return nil
}
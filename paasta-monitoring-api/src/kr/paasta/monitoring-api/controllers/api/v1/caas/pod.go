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
	CaasConfig models.CaasConfig
}

func GetPodController(config models.CaasConfig) *PodController{
	return &PodController{
		CaasConfig: config,
	}
}

func (controller *PodController) GetPodStatus(ctx echo.Context) error {
	results, err := service.GetPodService(controller.CaasConfig).GetPodStatus()
	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}
	return nil
}

package saas

import (
	"github.com/labstack/echo/v4"
	"log"
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


func (controller *PinpointController) GetAgentList(ctx echo.Context) error {
	result, err := service.GetPinpointService(controller.SaaS).GetAgentList()
	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", result)
	}
	return nil
}


func (controller *PinpointController) GetAgentStat(ctx echo.Context) error {
	result, err := service.GetPinpointService(controller.SaaS).GetAgentStat(ctx)
	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", result)
	}
	return nil
}
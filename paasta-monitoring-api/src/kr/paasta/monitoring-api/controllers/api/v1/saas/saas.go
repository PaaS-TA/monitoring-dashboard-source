package saas

import (
	"github.com/labstack/echo/v4"
	"log"
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


func (controller *SaasController) GetApplicationStatus(ctx echo.Context) error {
	result, err := service.GetSaasService(controller.SaaS).GetApplicationStatus()
	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", result)
	}
	return nil
}


func (controller *SaasController) GetApplicationUsage(ctx echo.Context) error {
	period := ctx.QueryParam("period")
	result, err := service.GetSaasService(controller.SaaS).GetApplicationUsage(period)
	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", result)
	}
	return nil
}


func (controller *SaasController) GetApplicationUsageList(ctx echo.Context) error {
	period := ctx.QueryParam("period")
	result, err := service.GetSaasService(controller.SaaS).GetApplicationUsageList(period)
	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", result)
	}
	return nil
}
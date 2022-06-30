package caas

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	models "paasta-monitoring-api/models/api/v1"
	service "paasta-monitoring-api/services/api/v1/caas"
)

type ClusterController struct {
	CaasConfig models.CaasConfig
}

func GetClusterController(config models.CaasConfig) *ClusterController{
	return &ClusterController{
		CaasConfig: config,
	}
}

func (controller *ClusterController) GetClusterAverage(ctx echo.Context) error {
	typeParam := ctx.Param("type")
	results, err := service.GetClusterService(controller.CaasConfig).GetClusterAverage(typeParam)
	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}
	return nil
}


func (controller *ClusterController) GetWorkNodeList(ctx echo.Context) error {
	results, err := service.GetClusterService(controller.CaasConfig).GetWorkNodeList()
	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}
	return nil
}


func (controller *ClusterController) GetWorkNode(ctx echo.Context) error {
	nodeName := ctx.QueryParam("nodename")
	instance := ctx.QueryParam("instance")

	results, err := service.GetClusterService(controller.CaasConfig).GetWorkNode(nodeName, instance)
	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Worker Node data.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}
	return nil
}
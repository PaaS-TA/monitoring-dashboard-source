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
	CaaS models.CaaS
}

func GetClusterController(config models.CaaS) *ClusterController{
	return &ClusterController{
		CaaS: config,
	}
}

func (controller *ClusterController) GetClusterAverage(ctx echo.Context) error {

	results, err := service.GetClusterService(controller.CaaS).GetClusterAverage(ctx)
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
	results, err := service.GetClusterService(controller.CaaS).GetWorkNodeList(ctx)
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


	results, err := service.GetClusterService(controller.CaaS).GetWorkNode(ctx)
	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Worker Node data.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}
	return nil
}
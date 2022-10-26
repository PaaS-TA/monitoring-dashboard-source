package controller

import (
	"github.com/cloudfoundry-community/go-cfclient"
	"monitoring-portal/paas/service"
	"monitoring-portal/paas/util"
	"net/http"
)

type CloudFoundryController struct {
	cfClient *cfclient.Client
}

func GetCloudFoundryController(cfClient *cfclient.Client) *CloudFoundryController {
	return &CloudFoundryController{
		cfClient: cfClient,
	}
}

func (c *CloudFoundryController) GetPaasDiagram(w http.ResponseWriter, r *http.Request) {
	// TODO
	result := service.GetCloudFoundryService(c.cfClient).GetPaasDiagram()
	util.RenderJsonResponse(result, w)
}
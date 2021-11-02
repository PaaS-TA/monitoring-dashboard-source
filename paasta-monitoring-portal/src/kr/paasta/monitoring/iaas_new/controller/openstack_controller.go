package controller

import (
	client "github.com/influxdata/influxdb1-client/v2"

	"kr/paasta/monitoring/iaas_new/model"
	"kr/paasta/monitoring/iaas_new/service"
	"kr/paasta/monitoring/utils"
	"net/http"
)

type OpenstackController struct {
	OpenstackProvider model.OpenstackProvider
	influxClient      client.Client
}

func NewOpenstackController(openstackProvider model.OpenstackProvider, influxClient client.Client) *OpenstackController {
	return &OpenstackController{
		OpenstackProvider: openstackProvider,
		influxClient:      influxClient,
	}
}


func (osService *OpenstackController) GetHypervisorStatistics(w http.ResponseWriter, r *http.Request) {

	provider, userName, err := utils.GetOpenstackProvider(r)

	result, err := service.GetOpenstackService(osService.OpenstackProvider, provider, osService.influxClient).GetHypervisorStatistics(userName)

	if err != nil {
		model.MonitLogger.Error("GetServerList error :", err)
		utils.ErrRenderJsonResponse(err, w)
	} else {
		utils.RenderJsonResponse(result, w)
	}
}

func (service *OpenstackController) GetServerList(w http.ResponseWriter, r *http.Request) {

}


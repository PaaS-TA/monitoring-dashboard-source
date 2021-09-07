package controller

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"kr/paasta/monitoring/iaas_new/model"
	"kr/paasta/monitoring/iaas_new/service"
	"kr/paasta/monitoring/utils"
	"net/http"
)

//Main Page Controller
type OpenstackServices struct {
	OpenstackProvider model.OpenstackProvider
	influxClient      client.Client
}

func NewMainController(openstackProvider model.OpenstackProvider, influxClient client.Client) *OpenstackServices {
	return &OpenstackServices{
		OpenstackProvider: openstackProvider,
		influxClient:      influxClient,
	}
}

func (h *OpenstackServices) Main(w http.ResponseWriter, r *http.Request) {
	model.MonitLogger.Debug("Main API Called")

	url := "/public/index.html"
	http.Redirect(w, r, url, 302)
}

func (s *OpenstackServices) OpenstackSummary(w http.ResponseWriter, r *http.Request) {

	provider, username, err := utils.GetOpenstackProvider(r)
	projectResourceSummary, err := services.GetMainService(s.OpenstackProvider, provider, s.influxClient).GetOpenstackSummary(username)

	if err != nil {
		model.MonitLogger.Error("GetOpenstackResources error :", err)
		utils.ErrRenderJsonResponse(err, w)
	} else {
		utils.RenderJsonResponse(projectResourceSummary, w)
	}

}

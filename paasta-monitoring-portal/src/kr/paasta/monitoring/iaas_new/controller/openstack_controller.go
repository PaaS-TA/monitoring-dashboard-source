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

/**
	하이퍼바이저 통계 데이터 조회
		- 화면 : public/index.html
 */
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


/**
	서버 목록 조회
 */
func (osService *OpenstackController) GetServerList(w http.ResponseWriter, r *http.Request) {

	tenantIdParam := r.URL.Query().Get("tenantId")


	provider, _, err := utils.GetOpenstackProvider(r)

	serverParams := make(map[string]interface{}, 0)
	serverParams["allTenants"] = true
	if tenantIdParam != "" {
		serverParams["tenantId"] = tenantIdParam
	}
	result, err := service.GetOpenstackService(osService.OpenstackProvider, provider, osService.influxClient).GetServerList(serverParams)

	if err != nil {
		model.MonitLogger.Error("GetServerList error :", err)
		utils.ErrRenderJsonResponse(err, w)
	} else {
		utils.RenderJsonResponse(result, w)
	}
}


/**
	프로젝트(테넌트) 목록과 속한 인스턴스 목록 및 usage 정보를 조회
 */
func (osService *OpenstackController) GetProjectUsage(w http.ResponseWriter, r *http.Request) {
	tenantIdParam := r.URL.Query().Get("tenantId")

	provider, _, err := utils.GetOpenstackProvider(r)

	serverParams := make(map[string]interface{}, 0)
	serverParams["allTenants"] = true
	if tenantIdParam != "" {
		serverParams["tenantId"] = tenantIdParam
	}
	result := service.GetOpenstackService(osService.OpenstackProvider, provider, osService.influxClient).RetrieveTenantUsage()

	if err != nil {
		model.MonitLogger.Error("GetServerList error :", err)
		utils.ErrRenderJsonResponse(err, w)
	} else {
		utils.RenderJsonResponse(result, w)
	}

}


/**
	@Unused
	프로젝트 목록(만) 조회
 */
func (osService *OpenstackController) GetProjectList(w http.ResponseWriter, r *http.Request) {
	provider, _, err := utils.GetOpenstackProvider(r)

	serverParams := make(map[string]interface{}, 0)
	result, err := service.GetOpenstackService(osService.OpenstackProvider, provider, osService.influxClient).GetProjectList(serverParams)

	if err != nil {
		model.MonitLogger.Error("GetProjectList error :", err)
		utils.ErrRenderJsonResponse(err, w)
	} else {
		utils.RenderJsonResponse(result, w)
	}
}


/**
	하이퍼바이저 목록 (상세 정보 포함) 조회
 */
func (osService *OpenstackController) GetHypervisorList(w http.ResponseWriter, r *http.Request) {
	provider, _, err := utils.GetOpenstackProvider(r)

	result, err := service.GetOpenstackService(osService.OpenstackProvider, provider, osService.influxClient).GetHypervisorList()

	if err != nil {
		model.MonitLogger.Error("GetProjectList error :", err)
		utils.ErrRenderJsonResponse(err, w)
	} else {
		utils.RenderJsonResponse(result, w)
	}
}



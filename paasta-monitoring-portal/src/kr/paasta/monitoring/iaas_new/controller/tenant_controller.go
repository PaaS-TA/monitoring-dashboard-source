package controller

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"kr/paasta/monitoring/iaas_new/model"
	"kr/paasta/monitoring/iaas_new/service"
	"kr/paasta/monitoring/utils"
	"net/http"
)

//Tenant Controller
type OpenstackTenant struct {
	openstackProvider model.OpenstackProvider
	influxClient      client.Client
}

func NewOpenstackTenantController(openstackProvider model.OpenstackProvider, influxClient client.Client) *OpenstackTenant {
	return &OpenstackTenant{
		openstackProvider: openstackProvider,
		influxClient:      influxClient,
	}
}

func (s *OpenstackTenant) TenantSummary(w http.ResponseWriter, r *http.Request) {

	//tenantName은 조회조건 (Optional)
	var apiRequest model.TenantReq
	apiRequest.TenantName = r.URL.Query().Get("tenantName")

	provider, username, _ := utils.GetOpenstackProvider(r)
	s.openstackProvider.Username = username
	tenantSummary, err := service.GetTenantService(s.openstackProvider, provider, s.influxClient).GetTenantSummary(apiRequest)

	if err != nil {
		utils.ErrRenderJsonResponse(err, w)
	} else {
		utils.RenderJsonResponse(tenantSummary, w)
	}
}

func (s *OpenstackTenant) GetTenantInstanceList(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.TenantReq
	apiRequest.TenantId = r.FormValue(":instanceId")
	//hostname은 조회조건 (Optional)
	apiRequest.HostName = r.URL.Query().Get("hostname")
	//Paging Size (Optional)
	apiRequest.Limit = r.URL.Query().Get("limit")
	//Paging 처리시 현재 Page Limit의 마지막 Instance Id를 요청 받으면 다음 Page를 조회 할 수 있다.
	//Limit과 같이 사용되어야 함 (Optional)
	apiRequest.Marker = r.URL.Query().Get("marker")

	validation := apiRequest.TenantInstanceRequestValidate(apiRequest)
	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validation.Error()))
		return
	}

	provider, _, _ := utils.GetOpenstackProvider(r)
	tenantSummary, err := service.GetTenantService(s.openstackProvider, provider, s.influxClient).GetTenantInstanceList(apiRequest)
	if err != nil {
		utils.ErrRenderJsonResponse(err, w)
	} else {
		utils.RenderJsonResponse(tenantSummary, w)
	}
}

func (s *OpenstackTenant) GetInstanceCpuUsageList(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.DetailReq
	apiRequest.InstanceId = r.FormValue(":instanceId")
	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy = r.URL.Query().Get("groupBy")

	validation := apiRequest.InstanceMetricRequestValidate(apiRequest)
	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validation.Error()))
		return
	}
	provider, _, _ := utils.GetOpenstackProvider(r)
	cpuUsageList, err := service.GetTenantService(s.openstackProvider, provider, s.influxClient).GetInstanceCpuUsageList(apiRequest)
	if err != nil {
		utils.ErrRenderJsonResponse(err, w)
	} else {
		utils.RenderJsonResponse(cpuUsageList, w)
	}
}

func (s *OpenstackTenant) GetInstanceMemoryUsageList(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.DetailReq
	apiRequest.InstanceId = r.FormValue(":instanceId")
	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy = r.URL.Query().Get("groupBy")

	validation := apiRequest.InstanceMetricRequestValidate(apiRequest)
	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validation.Error()))
		return
	}
	provider, _, _ := utils.GetOpenstackProvider(r)
	cpuUsageList, err := service.GetTenantService(s.openstackProvider, provider, s.influxClient).GetInstanceMemoryUsageList(apiRequest)
	if err != nil {
		utils.ErrRenderJsonResponse(err, w)
	} else {
		utils.RenderJsonResponse(cpuUsageList, w)
	}
}

func (s *OpenstackTenant) GetInstanceDiskReadList(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.DetailReq
	apiRequest.InstanceId = r.FormValue(":instanceId")
	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy = r.URL.Query().Get("groupBy")

	validation := apiRequest.InstanceMetricRequestValidate(apiRequest)
	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validation.Error()))
		return
	}

	provider, _, _ := utils.GetOpenstackProvider(r)
	cpuUsageList, err := service.GetTenantService(s.openstackProvider, provider, s.influxClient).GetInstanceDiskIoKbyteList(apiRequest, "read")
	if err != nil {
		utils.ErrRenderJsonResponse(err, w)
	} else {
		utils.RenderJsonResponse(cpuUsageList, w)
	}
}

func (s *OpenstackTenant) GetInstanceDiskWriteList(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.DetailReq
	apiRequest.InstanceId = r.FormValue(":instanceId")
	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy = r.URL.Query().Get("groupBy")

	validation := apiRequest.InstanceMetricRequestValidate(apiRequest)
	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validation.Error()))
		return
	}
	provider, _, _ := utils.GetOpenstackProvider(r)
	cpuUsageList, err := service.GetTenantService(s.openstackProvider, provider, s.influxClient).GetInstanceDiskIoKbyteList(apiRequest, "write")
	if err != nil {
		utils.ErrRenderJsonResponse(err, w)
	} else {
		utils.RenderJsonResponse(cpuUsageList, w)
	}
}

func (s *OpenstackTenant) GetInstanceNetworkIoList(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.DetailReq
	apiRequest.InstanceId = r.FormValue(":instanceId")
	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy = r.URL.Query().Get("groupBy")

	validation := apiRequest.InstanceMetricRequestValidate(apiRequest)
	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validation.Error()))
		return
	}
	provider, _, _ := utils.GetOpenstackProvider(r)
	cpuUsageList, err := service.GetTenantService(s.openstackProvider, provider, s.influxClient).GetInstanceNetworkIoKbyteList(apiRequest)
	if err != nil {
		utils.ErrRenderJsonResponse(err, w)
	} else {
		utils.RenderJsonResponse(cpuUsageList, w)
	}
}

func (s *OpenstackTenant) GetInstanceNetworkPacketsList(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.DetailReq
	apiRequest.InstanceId = r.FormValue(":instanceId")
	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy = r.URL.Query().Get("groupBy")

	validation := apiRequest.InstanceMetricRequestValidate(apiRequest)
	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validation.Error()))
		return
	}
	provider, _, _ := utils.GetOpenstackProvider(r)
	cpuUsageList, err := service.GetTenantService(s.openstackProvider, provider, s.influxClient).GetInstanceNetworkPacketsList(apiRequest)
	if err != nil {
		utils.ErrRenderJsonResponse(err, w)
	} else {
		utils.RenderJsonResponse(cpuUsageList, w)
	}
}

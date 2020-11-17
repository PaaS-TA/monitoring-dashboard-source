package controller

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring/paas/model"
	"kr/paasta/monitoring/paas/service"
	"kr/paasta/monitoring/paas/util"
	"net/http"
	"strconv"
)

//Gorm Object Struct
type ContainerService struct {
	txn          *gorm.DB
	influxClient client.Client
	databases    model.Databases
}

func GetContainerController(txn *gorm.DB, influxClient client.Client, databases model.Databases) *ContainerService {
	return &ContainerService{
		txn:          txn,
		influxClient: influxClient,
		databases:    databases,
	}
}

func (h *ContainerService) GetContainerDeploy(w http.ResponseWriter, r *http.Request) {

	containerDeployList, err := service.GetContainerService(h.txn, h.influxClient, h.databases).GetContainerDeploy()
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(containerDeployList, w)
	}
}

func (h *ContainerService) GetCellOverview(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.ContainerReq

	resList, err := service.GetContainerService(h.txn, h.influxClient, h.databases).GetCellOverview(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(resList, w)
	}
}

func (h *ContainerService) GetCellOverviewStatusList(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.ContainerReq
	apiRequest.Status = r.FormValue(":status")

	resList, err := service.GetContainerService(h.txn, h.influxClient, h.databases).GetCellOverviewStatusList(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(resList, w)
	}
}

func (h *ContainerService) GetContainerOverview(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.ContainerReq

	resList, err := service.GetContainerService(h.txn, h.influxClient, h.databases).GetContainerOverview(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(resList, w)
	}
}

func (h *ContainerService) GetContainerOverviewStatusList(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.ContainerReq
	apiRequest.Status = r.FormValue(":status")

	resList, err := service.GetContainerService(h.txn, h.influxClient, h.databases).GetContainerOverviewStatusList(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(resList, w)
	}
}

func (h *ContainerService) GetContainerSummary(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.ContainerReq
	//Page 번호
	apiRequest.PageIndex, _ = strconv.Atoi(r.FormValue("pageIndex"))
	//Page당 보여주는 갯수
	apiRequest.PageItems, _ = strconv.Atoi(r.FormValue("pageItems"))

	param := r.FormValue("zoneName")

	resList, err := service.GetContainerService(h.txn, h.influxClient, h.databases).GetContainerSummary(apiRequest, param)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(resList, w)
	}
}

func (h *ContainerService) GetContainerRelationship(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.ContainerReq
	apiRequest.Name = r.FormValue(":name")

	resList, err := service.GetContainerService(h.txn, h.influxClient, h.databases).GetContainerRelationship(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(resList, w)
	}
}

func (h *ContainerService) GetPaasMainContainerView(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.ContainerReq

	resList, err := service.GetContainerService(h.txn, h.influxClient, h.databases).GetPaasMainContainerView(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(resList, w)
	}
}

func (h *ContainerService) GetPaasContainerCpuUsages(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.ContainerReq
	apiRequest.ContainerName = r.FormValue(":id")
	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy = r.URL.Query().Get("groupBy")

	apiRequest.Item = []model.ContainerDetailReq{
		{model.CON_MTR_CPU_USAGE, "cpu"}}

	resList, err := service.GetContainerService(h.txn, h.influxClient, h.databases).GetPaasContainerUsages(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(resList, w)
	}
}

func (h *ContainerService) GetPaasContainerCpuLoads(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.ContainerReq
	apiRequest.ContainerName = r.FormValue(":id")
	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy = r.URL.Query().Get("groupBy")

	apiRequest.Item = []model.ContainerDetailReq{
		{model.CON_MTR_LOAD_AVG, "1m"}}

	resList, err := service.GetContainerService(h.txn, h.influxClient, h.databases).GetPaasContainerUsages(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(resList, w)
	}
}

func (h *ContainerService) GetPaasContainerMemoryUsages(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.ContainerReq
	apiRequest.ContainerName = r.FormValue(":id")
	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy = r.URL.Query().Get("groupBy")

	apiRequest.Item = []model.ContainerDetailReq{
		{model.CON_MTR_MEM_USAGE, "mem"}}

	resList, err := service.GetContainerService(h.txn, h.influxClient, h.databases).GetPaasContainerUsages(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(resList, w)
	}
}

func (h *ContainerService) GetPaasContainerDiskUsages(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.ContainerReq
	apiRequest.ContainerName = r.FormValue(":id")
	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy = r.URL.Query().Get("groupBy")

	apiRequest.Item = []model.ContainerDetailReq{
		{model.CON_MTR_DISK_USAGE, "disk"}}

	resList, err := service.GetContainerService(h.txn, h.influxClient, h.databases).GetPaasContainerUsages(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(resList, w)
	}
}

func (h *ContainerService) GetPaasContainerNetworkBytes(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.ContainerReq
	apiRequest.ContainerName = r.FormValue(":id")
	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy = r.URL.Query().Get("groupBy")

	apiRequest.Item = []model.ContainerDetailReq{
		{model.CON_MTR_RX_BYTES, "rx"},
		{model.CON_MTR_TX_BYTES, "tx"}}

	resList, err := service.GetContainerService(h.txn, h.influxClient, h.databases).GetPaasContainerUsages(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(resList, w)
	}
}

func (h *ContainerService) GetPaasContainerNetworkDrops(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.ContainerReq
	apiRequest.ContainerName = r.FormValue(":id")
	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy = r.URL.Query().Get("groupBy")

	apiRequest.Item = []model.ContainerDetailReq{
		{model.CON_MTR_RX_DROPPED, "rx"},
		{model.CON_MTR_TX_DROPPED, "tx"}}

	resList, err := service.GetContainerService(h.txn, h.influxClient, h.databases).GetPaasContainerUsages(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(resList, w)
	}
}

func (h *ContainerService) GetPaasContainerNetworkErrors(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.ContainerReq
	apiRequest.ContainerName = r.FormValue(":id")
	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy = r.URL.Query().Get("groupBy")

	apiRequest.Item = []model.ContainerDetailReq{
		{model.CON_MTR_RX_ERRORS, "rx"},
		{model.CON_MTR_TX_ERRORS, "tx"}}

	resList, err := service.GetContainerService(h.txn, h.influxClient, h.databases).GetPaasContainerUsages(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(resList, w)
	}
}

package controller

import (
	client "github.com/influxdata/influxdb1-client/v2"
	model "kr/paasta/monitoring/paas/model"
	"kr/paasta/monitoring/paas/service"
	"kr/paasta/monitoring/paas/util"
	"net/http"
	//models "kr/paasta/monitoring/iaas/model"
	"github.com/jinzhu/gorm"
	"strconv"
)

//Gorm Object Struct
type BoshStatusService struct {
	txn          *gorm.DB
	influxClient client.Client
	databases    model.Databases
}

func GetBoshStatusController(txn *gorm.DB, influxClient client.Client, databases model.Databases) *BoshStatusService {
	return &BoshStatusService{
		txn:          txn,
		influxClient: influxClient,
		databases:    databases,
	}
}

func (h *BoshStatusService) GetBoshStatusOverview(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.BoshSummaryReq
	//Page 번호
	apiRequest.PageIndex, _ = strconv.Atoi(r.FormValue("pageIndex"))
	//Page당 보여주는 갯수
	apiRequest.PageItem, _ = strconv.Atoi(r.FormValue("pageItems"))

	boshOverview, err := service.GetBoshStatusService(h.txn, h.influxClient, h.databases).GetBoshStatusOverview(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(boshOverview, w)
	}
}

func (h *BoshStatusService) GetBoshStatusSummary(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.BoshSummaryReq
	//Page 번호
	apiRequest.PageIndex, _ = strconv.Atoi(r.FormValue("pageIndex"))
	//Page당 보여주는 갯수
	apiRequest.PageItem, _ = strconv.Atoi(r.FormValue("pageItems"))

	boshSummary, err := service.GetBoshStatusService(h.txn, h.influxClient, h.databases).GetBoshStatusSummary(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(boshSummary, w)
	}
}

func (h *BoshStatusService) GetBoshStatusTopprocess(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.BoshSummaryReq
	apiRequest.Id = r.FormValue(":id")

	topProcess, err := service.GetBoshStatusService(h.txn, h.influxClient, h.databases).GetTopProcessListByMemory(apiRequest)

	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(topProcess, w)
	}
}

func (h *BoshStatusService) GetBoshCpuUsageList(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.BoshDetailReq
	apiRequest.Id = r.FormValue(":id")
	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy = r.URL.Query().Get("groupBy")

	validation := apiRequest.MetricRequestValidate(apiRequest)
	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validation.Error()))
		return
	}
	resList, err := service.GetBoshStatusService(h.txn, h.influxClient, h.databases).GetBoshCpuUsageList(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(resList, w)
	}
}

func (h *BoshStatusService) GetBoshCpuLoadList(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.BoshDetailReq
	apiRequest.Id = r.FormValue(":id")
	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy = r.URL.Query().Get("groupBy")

	validation := apiRequest.MetricRequestValidate(apiRequest)
	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validation.Error()))
		return
	}
	resList, err := service.GetBoshStatusService(h.txn, h.influxClient, h.databases).GetBoshCpuLoadList(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(resList, w)
	}
}

func (h *BoshStatusService) GetBoshMemoryUsageList(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.BoshDetailReq
	apiRequest.Id = r.FormValue(":id")
	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy = r.URL.Query().Get("groupBy")

	validation := apiRequest.MetricRequestValidate(apiRequest)
	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validation.Error()))
		return
	}
	resList, err := service.GetBoshStatusService(h.txn, h.influxClient, h.databases).GetBoshMemoryUsageList(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(resList, w)
	}
}

func (h *BoshStatusService) GetBoshDiskUsageList(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.BoshDetailReq
	apiRequest.Id = r.FormValue(":id")
	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy = r.URL.Query().Get("groupBy")

	validation := apiRequest.MetricRequestValidate(apiRequest)
	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validation.Error()))
		return
	}
	resList, err := service.GetBoshStatusService(h.txn, h.influxClient, h.databases).GetBoshDiskUsageList(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(resList, w)
	}
}

func (h *BoshStatusService) GetBoshDiskIoList(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.BoshDetailReq
	apiRequest.Id = r.FormValue(":id")
	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy = r.URL.Query().Get("groupBy")

	validation := apiRequest.MetricRequestValidate(apiRequest)
	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validation.Error()))
		return
	}
	resList, err := service.GetBoshStatusService(h.txn, h.influxClient, h.databases).GetBoshDiskIoList(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(resList, w)
	}
}

func (h *BoshStatusService) GetBoshNetworkByteList(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.BoshDetailReq
	apiRequest.Id = r.FormValue(":id")
	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy = r.URL.Query().Get("groupBy")

	validation := apiRequest.MetricRequestValidate(apiRequest)
	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validation.Error()))
		return
	}
	resList, err := service.GetBoshStatusService(h.txn, h.influxClient, h.databases).GetBoshNetworkByteList(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(resList, w)
	}
}

func (h *BoshStatusService) GetBoshNetworkPacketList(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.BoshDetailReq
	apiRequest.Id = r.FormValue(":id")
	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy = r.URL.Query().Get("groupBy")

	validation := apiRequest.MetricRequestValidate(apiRequest)
	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validation.Error()))
		return
	}
	resList, err := service.GetBoshStatusService(h.txn, h.influxClient, h.databases).GetBoshNetworkPacketList(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(resList, w)
	}
}

func (h *BoshStatusService) GetBoshNetworkDropList(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.BoshDetailReq
	apiRequest.Id = r.FormValue(":id")
	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy = r.URL.Query().Get("groupBy")

	validation := apiRequest.MetricRequestValidate(apiRequest)
	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validation.Error()))
		return
	}
	resList, err := service.GetBoshStatusService(h.txn, h.influxClient, h.databases).GetBoshNetworkDropList(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(resList, w)
	}
}

func (h *BoshStatusService) GetBoshNetworkErrorList(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.BoshDetailReq
	apiRequest.Id = r.FormValue(":id")
	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy = r.URL.Query().Get("groupBy")

	validation := apiRequest.MetricRequestValidate(apiRequest)
	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validation.Error()))
		return
	}
	resList, err := service.GetBoshStatusService(h.txn, h.influxClient, h.databases).GetBoshNetworkErrorList(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(resList, w)
	}
}

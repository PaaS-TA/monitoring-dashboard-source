package controller

import (
	"kr/paasta/monitoring/paas/model"
	"kr/paasta/monitoring/paas/service"
	"net/http"
	"kr/paasta/monitoring/paas/util"
	"github.com/influxdata/influxdb/client/v2"
	"fmt"
)

type InfluxServerClient struct {
	client   client.Client
	databases model.Databases
}

func GetMetricsController(client client.Client, databases model.Databases) *InfluxServerClient {
	return &InfluxServerClient{
		client:   client,
		databases: databases,
	}
}

func (h InfluxServerClient) GetDiskIOList(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.MetricsRequest

	apiRequest.Origin = r.FormValue(":origin")

	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.ServiceName      = r.URL.Query().Get("serviceName")
	apiRequest.Index            = r.URL.Query().Get("index")
	apiRequest.Addr             = r.URL.Query().Get("addr")
	apiRequest.TimeRangeFrom    = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo      = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy          = r.URL.Query().Get("groupBy")

	//service호출 (Gorm Obj 매개 변수)
	result, err := service.GetMetricsService(h.client, h.databases).GetDiskIOList(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

func (h InfluxServerClient) GetNetworkIOList(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.MetricsRequest

	apiRequest.Origin = r.FormValue(":origin")

	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.ServiceName      = r.URL.Query().Get("serviceName")
	apiRequest.Index            = r.URL.Query().Get("index")
	apiRequest.Addr             = r.URL.Query().Get("addr")
	apiRequest.TimeRangeFrom    = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo      = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy          = r.URL.Query().Get("groupBy")

	//service호출 (Gorm Obj 매개 변수)
	result, err := service.GetMetricsService(h.client, h.databases).GetNetworkIOList(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

func (h InfluxServerClient) GetTopProcessList(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.MetricsRequest

	apiRequest.Origin = r.FormValue(":origin")

	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.ServiceName      = r.URL.Query().Get("serviceName")
	apiRequest.Index            = r.URL.Query().Get("index")
	apiRequest.Addr             = r.URL.Query().Get("addr")
	apiRequest.TimeRangeFrom    = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo      = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy          = r.URL.Query().Get("groupBy")

	//service호출 (Gorm Obj 매개 변수)
	result, err := service.GetMetricsService(h.client, h.databases).GetTopProcessList(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}


/**
 Date: 2017-08-16
 Description: Application guid 및 조회기간 정보를 전달받아 해당 Application의 CPU 변화량 정보를 리턴한다.
 	      Portal 시스템에서 호출한다.
 */
func (h InfluxServerClient) GetAppCpuUsage(w http.ResponseWriter, r *http.Request) {
	var apiRequest model.MetricsRequest
	apiRequest.ServiceName      = r.FormValue(":guid")
	apiRequest.Index       	    = r.FormValue(":idx")
	apiRequest.DefaultTimeRange = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom    = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo      = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy          = r.URL.Query().Get("groupBy")

	fmt.Println("Range::",apiRequest.DefaultTimeRange)
	if apiRequest.DefaultTimeRange == "" && apiRequest.TimeRangeFrom == "" && apiRequest.TimeRangeTo == ""{
		apiRequest.DefaultTimeRange = "15m"	//15 minutes
	}

	if apiRequest.GroupBy == ""{
		apiRequest.GroupBy = "1m"	//15 minutes
	}

	result, err := service.GetMetricsService(h.client, h.databases).GetAppCpuUsage(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}



/**
 Date: 2017-08-16
 Description: Application guid 및 조회기간 정보를 전달받아 해당 Application의 Memory 변화량 정보를 리턴한다.
 	      Portal 시스템에서 호출한다.
 */
func (h InfluxServerClient) GetAppMemoryUsage(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.MetricsRequest
	apiRequest.ServiceName 	     = r.FormValue(":guid")
	apiRequest.Index             = r.FormValue(":idx")
	apiRequest.DefaultTimeRange  = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom     = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo       = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy          = r.URL.Query().Get("groupBy")


	if apiRequest.DefaultTimeRange == "" && apiRequest.TimeRangeFrom == "" && apiRequest.TimeRangeTo == ""{
		apiRequest.DefaultTimeRange = "15m"	//15 minutes
	}

	if apiRequest.GroupBy == ""{
		apiRequest.GroupBy = "1m"	//15 minutes
	}

	result, err := service.GetMetricsService(h.client, h.databases).GetAppMemoryUsage(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}


//App Disk Usage
func (h InfluxServerClient) GetDiskUsage(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.MetricsRequest
	apiRequest.ServiceName 	     = r.FormValue(":guid")
	apiRequest.Index             = r.FormValue(":idx")
	apiRequest.DefaultTimeRange  = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom     = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo       = r.URL.Query().Get("timeRangeTo")
	apiRequest.GroupBy          = r.URL.Query().Get("groupBy")

	if apiRequest.DefaultTimeRange == "" && apiRequest.TimeRangeFrom == "" && apiRequest.TimeRangeTo == ""{
		apiRequest.DefaultTimeRange = "15m"	//15 minutes
	}

	if apiRequest.GroupBy == ""{
		apiRequest.GroupBy = "1m"	//15 minutes
	}

	//service호출 (Gorm Obj 매개 변수)
	result, err := service.GetMetricsService(h.client, h.databases).GetDiskUsage(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

/**
 Date: 2017-08-16
 Description: Application guid 및 조회기간 정보를 전달받아 해당 Application의 Network 변화량 정보를 리턴한다.
 	      Portal 시스템에서 호출한다.
 */
func (h InfluxServerClient) GetAppNetworkIoKByte(w http.ResponseWriter, r *http.Request) {
	fmt.Println("xxxxxxxxxxxxxxxx")
	var apiRequest model.MetricsRequest
	apiRequest.ServiceName 	     = r.FormValue(":guid")
	apiRequest.Index             = r.FormValue(":idx")
	apiRequest.DefaultTimeRange  = r.URL.Query().Get("defaultTimeRange")
	apiRequest.TimeRangeFrom     = r.URL.Query().Get("timeRangeFrom")
	apiRequest.TimeRangeTo       = r.URL.Query().Get("timeRangeTo")

	if apiRequest.DefaultTimeRange == "" && apiRequest.TimeRangeFrom == "" && apiRequest.TimeRangeTo == ""{
		apiRequest.DefaultTimeRange = "15m"	//15 minutes
	}

	if apiRequest.GroupBy == ""{
		apiRequest.GroupBy = "1m"	//15 minutes
	}

	result, err := service.GetMetricsService(h.client, h.databases).GetAppNetworkKByte(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

/**
 Date: 2017-08-14
 Description: Application guid를 전달받아 해당 Application의 Resource 정보를 리턴한다.
 	      Portal 시스템에서 호출한다. (기존 paasta-monitoring-api-release 대체)
 */
func (h InfluxServerClient) GetApplicationResources(w http.ResponseWriter, r *http.Request) {
	var apiRequest model.MetricsRequest
	apiRequest.ServiceName = r.URL.Query().Get("app_id")
	apiRequest.Index       = r.URL.Query().Get("app_index")
	apiRequest.DefaultTimeRange = "30s"

	//service호출 (Gorm Obj 매개 변수)
	result, err := service.GetMetricsService(h.client, h.databases).GetApplicationResources(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}
func (h InfluxServerClient) GetApplicationResourcesAll(w http.ResponseWriter, r *http.Request) {
	var apiRequest model.MetricsRequest
	apiRequest.Index = r.URL.Query().Get("limit")
	apiRequest.DefaultTimeRange = "30s"

	//service호출 (Gorm Obj 매개 변수)
	result, err := service.GetMetricsService(h.client, h.databases).GetApplicationResourcesAll(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

func (h InfluxServerClient) PaasMain(w http.ResponseWriter, r *http.Request) {

	url := "/public/dist/index.html"
	http.Redirect(w, r, url, 302)
}





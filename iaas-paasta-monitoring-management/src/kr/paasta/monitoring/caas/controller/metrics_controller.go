package controller

import (
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring/caas/model"
	"kr/paasta/monitoring/caas/service"
	"kr/paasta/monitoring/caas/util"
	"net/http"
)

type MetricController struct {
	txn *gorm.DB
}

func (s *MetricController) GetClusterAvg(w http.ResponseWriter, r *http.Request) {
	//service호출
	result, err := service.GetMetricsService().GetClusterAvg()
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

func (s *MetricController) GetWorkNodeList(w http.ResponseWriter, r *http.Request) {
	//service호출
	result, err := service.GetMetricsService().GetWorkNodeList()
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

func (s *MetricController) GetWorkNodeInfo(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.MetricsRequest

	apiRequest.Nodename = r.URL.Query().Get("NodeName")
	apiRequest.Instance = r.URL.Query().Get("Instance")

	//service호출
	result, err := service.GetMetricsService().GetWorkNodeInfo(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

func (s *MetricController) GetContainerList(w http.ResponseWriter, r *http.Request) {
	var apiRequest model.MetricsRequest
	apiRequest.WorkloadsName = r.URL.Query().Get("WorkloadsName")
	apiRequest.PodName = r.URL.Query().Get("PodName")

	//service호출
	result, err := service.GetMetricsService().GetContainerList(apiRequest)

	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

func (s *MetricController) GetContainerInfo(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.MetricsRequest

	apiRequest.ContainerName = r.URL.Query().Get("ContainerName")
	apiRequest.PodName = r.URL.Query().Get("PodName")
	apiRequest.NameSpace = r.URL.Query().Get("NameSpace")

	//service호출
	result, err := service.GetMetricsService().GetContainerInfo(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

func (s *MetricController) GetContainerLog(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.MetricsRequest

	apiRequest.ContainerName = r.URL.Query().Get("ContainerName")
	apiRequest.PodName = r.URL.Query().Get("PodName")
	apiRequest.NameSpace = r.URL.Query().Get("NameSpace")

	//service호출
	result, err := service.GetMetricsService().GetContainerLog(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

func (s *MetricController) GetClusterOverView(w http.ResponseWriter, r *http.Request) {
	//service호출
	result, err := service.GetMetricsService().GetClusterOverView()
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

func (s *MetricController) GetWorkloadsStatus(w http.ResponseWriter, r *http.Request) {
	//service호출
	result, err := service.GetMetricsService().GetWorkloadsStatus()
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

func (s *MetricController) GetMasterNodeUsage(w http.ResponseWriter, r *http.Request) {
	//service호출
	result, err := service.GetMetricsService().GetWorkNodeList()
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

func (s *MetricController) GetWorkloadsContiSummary(w http.ResponseWriter, r *http.Request) {
	//service호출
	result, err := service.GetMetricsService().GetWorkloadsContiSummary()
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

func (s *MetricController) GetWorkloadsUsage(w http.ResponseWriter, r *http.Request) {
	var apiRequest model.MetricsRequest

	apiRequest.WorkloadsName = r.URL.Query().Get("WorkloadsName")

	//service호출
	result, err := service.GetMetricsService().GetWorkloadsUsage(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

func (s *MetricController) GetPodStatList(w http.ResponseWriter, r *http.Request) {
	//service호출
	result, err := service.GetMetricsService().GetPodStatList()
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

func (s *MetricController) GetPodMetricList(w http.ResponseWriter, r *http.Request) {
	//service호출
	result, err := service.GetMetricsService().GetPodMetricList()
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

func (s *MetricController) GetPodInfo(w http.ResponseWriter, r *http.Request) {
	var apiRequest model.MetricsRequest

	apiRequest.PodName = r.URL.Query().Get("PodName")

	//service호출
	result, err := service.GetMetricsService().GetPodInfo(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

func (s *MetricController) GetWorkNodeInfoGraph(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.MetricsRequest

	apiRequest.Nodename = r.URL.Query().Get("NodeName")
	apiRequest.Instance = r.URL.Query().Get("Instance")

	//service호출
	result, err := service.GetMetricsService().GetWorkNodeInfoGraph(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

func (s *MetricController) GetWorkloadsInfoGraph(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.MetricsRequest

	apiRequest.WorkloadsName = r.URL.Query().Get("WorkloadsName")

	//service호출
	result, err := service.GetMetricsService().GetWorkloadsInfoGraph(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

func (s *MetricController) GetPodInfoGraph(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.MetricsRequest

	apiRequest.PodName = r.URL.Query().Get("PodName")

	//service호출
	result, err := service.GetMetricsService().GetPodInfoGraph(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

func (s *MetricController) GetContainerInfoGraph(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.MetricsRequest

	apiRequest.ContainerName = r.URL.Query().Get("ContainerName")
	apiRequest.PodName = r.URL.Query().Get("PodName")
	apiRequest.NameSpace = r.URL.Query().Get("NameSpace")

	//service호출
	result, err := service.GetMetricsService().GetContainerInfoGraph(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

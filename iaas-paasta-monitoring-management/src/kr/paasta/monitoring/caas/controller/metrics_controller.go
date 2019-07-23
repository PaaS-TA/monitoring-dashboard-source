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
	//service호출
	result, err := service.GetMetricsService().GetContainerList()

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

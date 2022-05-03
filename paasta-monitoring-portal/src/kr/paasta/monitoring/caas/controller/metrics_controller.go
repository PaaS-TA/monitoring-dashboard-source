package controller

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"monitoring-portal/caas/model"
	"monitoring-portal/caas/service"
	"monitoring-portal/caas/util"
	"net/http"
	"strconv"
)

type MetricController struct {
	txn *gorm.DB
}

func NewMetricControllerr(txn *gorm.DB) *MetricController {
	return &MetricController{
		txn: txn,
	}
}

func (s *MetricController) GetClusterAvg(w http.ResponseWriter, r *http.Request) {
	//service호출
	result, err := service.GetMetricsService().GetClusterAvg()
	if err != nil {
		util.RenderJsonResponse(err, w)
		return
	}
	util.RenderJsonResponse(result, w)
}

func (s *MetricController) GetWorkNodeList(w http.ResponseWriter, r *http.Request) {
	//service호출
	result, err := service.GetMetricsService().GetWorkNodeList()
	if err != nil {
		util.RenderJsonResponse(err, w)
		return
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
		return
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
		return
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
		return
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
		return
	}
	util.RenderJsonResponse(result, w)
}

func (s *MetricController) GetClusterOverView(w http.ResponseWriter, r *http.Request) {
	//service호출
	result, err := service.GetMetricsService().GetClusterOverView()
	if err != nil {
		util.RenderJsonResponse(err, w)
		return
	}
	util.RenderJsonResponse(result, w)
}

func (s *MetricController) GetWorkloadsStatus(w http.ResponseWriter, r *http.Request) {
	//service호출
	result, err := service.GetMetricsService().GetWorkloadsStatus()
	if err != nil {
		util.RenderJsonResponse(err, w)
		return
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

func (s *MetricController) GetWorkNodeAvg(w http.ResponseWriter, r *http.Request) {
	//service호출
	result, err := service.GetMetricsService().GetWorkNodeAvg()
	if err != nil {
		util.RenderJsonResponse(err, w)
		return
	}
	util.RenderJsonResponse(result, w)
}

func (s *MetricController) GetWorkloadsContiSummary(w http.ResponseWriter, r *http.Request) {
	//service호출
	result, err := service.GetMetricsService().GetWorkloadsContiSummary()
	if err != nil {
		util.RenderJsonResponse(err, w)
		return
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
	} else {
		util.RenderJsonResponse(result, w)
	}

}

func (s *MetricController) GetPodStatList(w http.ResponseWriter, r *http.Request) {
	//service호출
	result, err := service.GetMetricsService().GetPodStatList()
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(result, w)
	}
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
	} else {
		util.RenderJsonResponse(result, w)
	}

}

func (s *MetricController) GetWorkNodeInfoGraph(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.MetricsRequest

	apiRequest.Nodename = r.URL.Query().Get("NodeName")
	apiRequest.Instance = r.URL.Query().Get("Instance")

	//service호출
	result, err := service.GetMetricsService().GetWorkNodeInfoGraph(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(result, w)
	}

}

func (s *MetricController) GetWorkloadsInfoGraph(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.MetricsRequest
	apiRequest.WorkloadsName = r.URL.Query().Get("WorkloadsName")

	//service호출
	result, err := service.GetMetricsService().GetWorkloadsInfoGraph(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(result, w)
	}

}

func (s *MetricController) GetPodInfoGraph(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.MetricsRequest

	apiRequest.PodName = r.URL.Query().Get("PodName")

	//service호출
	result, err := service.GetMetricsService().GetPodInfoGraph(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(result, w)
	}

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
	} else {
		util.RenderJsonResponse(result, w)
	}

}

func (s *MetricController) GetWorkNodeInfoGraphList(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.MetricsRequest

	apiRequest.Nodename = r.URL.Query().Get("NodeName")
	apiRequest.Instance = r.URL.Query().Get("Instance")

	//service호출
	result, err := service.GetMetricsService().GetWorkNodeInfoGraphList(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(result, w)
	}

}

//Alarm Process
func (s *MetricController) GetAlarmInfo(w http.ResponseWriter, r *http.Request) {
	//service호출
	result, err := service.GetAlarmService(s.txn).GetAlarmInfo()
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

func (s *MetricController) GetAlarmUpdate(w http.ResponseWriter, r *http.Request) {
	var apiRequest []model.AlarmPolicyRequest
	data, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	err := json.Unmarshal(data, &apiRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	service.GetAlarmService(s.txn).GetAlarmUpdate(apiRequest)
}

func (s *MetricController) GetAlarmLog(w http.ResponseWriter, r *http.Request) {
	searchDateFrom := r.URL.Query().Get("searchDateFrom")
	searchDateTo := r.URL.Query().Get("searchDateTo")

	alarmType := r.URL.Query().Get("alarmType")
	alarmStatus := r.URL.Query().Get("alarmStatus")
	resolveStatus := r.URL.Query().Get("resolveStatus")

	result, err := service.GetAlarmService(s.txn).GetAlarmLog(searchDateFrom, searchDateTo, alarmType, alarmStatus, resolveStatus)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(result, w)
	}

}

func (p *MetricController) GetSnsInfo(w http.ResponseWriter, r *http.Request) {
	result, err := service.GetAlarmService(p.txn).GetSnsInfo()
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(result, w)
	}

}

func (p *MetricController) GetAlarmCount(w http.ResponseWriter, r *http.Request) {
	searchDateFrom := r.URL.Query().Get("searchDateFrom")
	searchDateTo := r.URL.Query().Get("searchDateTo")
	result, err := service.GetAlarmService(p.txn).GetAlarmCount(searchDateFrom, searchDateTo)

	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(result, w)
	}

}

func (p *MetricController) GetlarmSnsSave(w http.ResponseWriter, r *http.Request) {

	var alarmSns model.BatchAlarmSnsRequest
	err := json.NewDecoder(r.Body).Decode(&alarmSns)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	error := service.GetAlarmService(p.txn).GetlarmSnsSave(alarmSns)

	util.RenderJsonResponse(error, w)
}

func (h *MetricController) UpdateAlarmState(w http.ResponseWriter, r *http.Request) {

	var alarmrRsolveRequest model.AlarmrRsolveRequest
	err := json.NewDecoder(r.Body).Decode(&alarmrRsolveRequest)
	defer r.Body.Close()

	id, _ := strconv.Atoi(r.FormValue(":id"))
	alarmrRsolveRequest.Id = uint64(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	error := service.GetAlarmService(h.txn).UpdateAlarmSate(alarmrRsolveRequest)
	util.RenderJsonResponse(error, w)
}

func (h *MetricController) CreateAlarmResolve(w http.ResponseWriter, r *http.Request) {
	var alarmrRsolveRequest model.AlarmrRsolveRequest
	err := json.NewDecoder(r.Body).Decode(&alarmrRsolveRequest)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	error := service.GetAlarmService(h.txn).CreateAlarmResolve(alarmrRsolveRequest)
	util.RenderJsonResponse(error, w)
}

func (h *MetricController) DeleteAlarmResolve(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue(":id"))

	error := service.GetAlarmService(h.txn).DeleteAlarmResolve(uint64(id))
	util.RenderJsonResponse(error, w)
	return
}

func (h *MetricController) UpdateAlarmResolve(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue(":id"))

	var alarmrRsolveRequest model.AlarmrRsolveRequest
	err := json.NewDecoder(r.Body).Decode(&alarmrRsolveRequest)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	alarmrRsolveRequest.Id = uint64(id)
	error := service.GetAlarmService(h.txn).UpdateAlarmResolve(alarmrRsolveRequest)
	util.RenderJsonResponse(error, w)
	return
}

func (h *MetricController) GetAlarmSnsReceiver(w http.ResponseWriter, r *http.Request) {
	alarmReceiver, _ := service.GetAlarmService(h.txn).GetAlarmSnsReceiver()
	util.RenderJsonResponse(alarmReceiver, w)
}

func (h *MetricController) DeleteAlarmSnsChannel(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue(":id"))

	err := service.GetAlarmService(h.txn).DeleteAlarmSnsChannel(id)
	util.RenderJsonResponse(err, w)
}

func (h *MetricController) GetAlarmActionList(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(r.FormValue(":id"))

	result, err := service.GetAlarmService(h.txn).GetAlarmActionList(id)

	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(result, w)
	}
}

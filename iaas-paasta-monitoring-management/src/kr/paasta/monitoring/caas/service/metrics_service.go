package service

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"kr/paasta/monitoring/caas/model"
	"kr/paasta/monitoring/caas/util"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// PLUS(+)
	PLUS = "%2B"

	// metricUrl
	SUB_URI       = "/api/v1/query?query="
	SUB_URI_RANGE = "/api/v1/query_range?query="
	K8S_SUB_URI   = "/api/v1/"

	//Select division
	WORKLOADS = "workLoads"
	POD       = "pod"
	CPU       = "cpu"
	MEMORY    = "memory"
	DISK      = "disk"

	//Jpath Type (JSON PATH)
	VALUE_0_DATA           = "data.result.0.value.0"
	VALUE_1_DATA           = "data.result.0.value.1"
	VALUE_2_DATA           = "data.result.0.value.2"
	RESULT_0_METRIC_0_DATA = "data.result.0.metric.0"

	//(PromQl)
	//Cluster Usage Metrics
	PQ_POD_USAGE    = "sum(kube_pod_info)/sum(kube_node_status_allocatable_pods{node=~'.*'})"
	PQ_CPU_USAGE    = "cluster:node_cpu:sum_rate5m"
	PQ_MEMORY_USAGE = "(sum(node_memory_MemTotal_bytes)-sum(node_memory_MemFree_bytes" + PLUS + "node_memory_Buffers_bytes" + PLUS + "node_memory_Cached_bytes))/sum(node_memory_MemTotal_bytes)*100"
	PQ_DISK_USAGE   = "(sum(node_filesystem_files)-sum(node_filesystem_files_free))/sum(node_filesystem_files)*100"

	//Cluster Overview
	PQ_CLUSTER_ALERTS            = "sum(ALERTS)"
	PQ_CLUSTER_RUNNING_POD       = "sum(kubelet_running_pod_count)"
	PQ_CLUSTER_RUNNING_CONTAINER = "sum(kubelet_running_container_count)"
	PQ_CLUSTER_POD_RESTART       = "kube_pod_container_status_restarts_total"
	PQ_CLUSTER_NODES             = "sum(kube_node_info)"

	//Workloads Status (Deployment, Daemonset, StateFulset, PodContainer)
	PQ_DEPLOYMENT_TOTAL       = "sum(kube_deployment_status_replicas)"
	PQ_DEPLOYMENT_AVAILABLE   = "sum(kube_deployment_status_replicas_available)"
	PQ_DEPLOYMENT_UNAVAILABLE = "sum(kube_deployment_status_replicas_unavailable)"
	PQ_DEPLOYMENT_UPDATED     = "sum(kube_deployment_status_replicas_updated)"
	PQ_DAEMONSET_READY        = "sum(kube_daemonset_status_number_ready)"
	PQ_DAEMONSET_AVAILABLE    = "sum(kube_daemonset_status_number_available)"
	PQ_DAEMONSET_UNAVAILABLE  = "sum(kube_daemonset_status_number_unavailable)"
	PQ_DAEMONSET_MISSCHEDULED = "sum(kube_daemonset_status_number_misscheduled)"
	PQ_STATEFULSET_TOTAL      = "sum(kube_statefulset_status_replicas)"
	PQ_STATEFULSET_READY      = "sum(kube_statefulset_status_replicas_ready)"
	PQ_STATEFULSET_UPDATED    = "sum(kube_statefulset_status_replicas_updated)"
	PQ_STATEFULSET_REVISION   = "sum(kube_statefulset_status_update_revision)"
	PQ_PODCONTAINER_READY     = "sum(kube_pod_container_status_ready)"
	PQ_PODCONTAINER_RUNNING   = "sum(kube_pod_container_status_running)"
	PQ_PODCONTAINER_RESTATS   = "sum(kube_pod_container_status_restarts_total)"
	PQ_PODCONTAINER_TERMINATE = "sum(kube_pod_container_status_terminated)"

	//Work Node Usage Metrics
	PQ_WORK_NODE_NAME_LIST    = "count(node_uname_info)by(instance,nodename,namespace)"
	PQ_WORK_NODE_CPU_USAGE    = "sum(instance:node_cpu:ratio)by(instance)*100"
	PQ_WORK_NODE_MEMORY_USAGE = "max(((node_memory_MemTotal_bytes{job='node-exporter'}-" +
		"node_memory_MemFree_bytes{job='node-exporter'}" +
		"-node_memory_Buffers_bytes{job='node-exporter'}" +
		"-node_memory_Cached_bytes{job='node-exporter'})" +
		"/node_memory_MemTotal_bytes{job='node-exporter'})*100)by(instance)"
	//CHANGE METRIC CPU_USE => CPU ALLOC
	//PQ_WORK_NODE_CPU_USE      = "avg(node_cpu_seconds_total{job='node-exporter',mode!='idle'})by(instance)"
	PQ_WORK_NODE_CPU_ALLOC  = "sum(kube_node_status_allocatable_cpu_cores*100)by(node)"
	PQ_WORK_NODE_MEMORY_USE = "sum(node_memory_MemTotal_bytes{job='node-exporter'})by(instance)"
	PQ_WORK_NODE_DISK_USE   = "sum(node_filesystem_size_bytes{job='node-exporter'})by(instance)"
	PQ_WORK_NODE_CONDITION  = "count(kube_node_status_condition{condition='Ready',status='true'})by(node)"

	//Workloads Container Metrics
	/*
		PQ_COTAINER_NAME_LIST  = "count(kube_$workloadName_metadata_generation)by(namespace,$workloadName)"
		PQ_COTAINER_CPU_USAGE  = "sum(rate(container_cpu_usage_seconds_total{container_name!='POD'image!=''}[2m]))by(namespace,pod_name,container_name)*100"
	*/

	//Pod Usage Metrics
	PQ_POD_LIST         = "count(container_cpu_usage_seconds_total{pod_name!='',container_name!='POD',image!=''})by(pod_name)"
	PQ_POD_CPU_USAGE    = "avg(rate(container_cpu_usage_seconds_total{pod_name!=''}[5m]))by(pod_name)*100"
	PQ_POD_MEMORY_USE   = "sum(container_memory_working_set_bytes{pod_name!=''})by(pod_name)/1024/1024"
	PQ_POD_DISK_USE     = "sum(container_fs_usage_bytes{pod_name!=''})by(pod_name)/1024/1024"
	PQ_POD_DISK_USAGE   = "sum(container_fs_usage_bytes{container_name!=''})by(pod_name)/avg(container_fs_limit_bytes{container_name!=''})by(pod_name)*100"
	PQ_POD_MEMORY_USAGE = "(sum(container_memory_working_set_bytes{pod_name!=''})by(pod_name))/(sum(container_memory_usage_bytes{pod_name!=''})by(pod_name))*100"

	//Container Usage Metrics
	//	PQ_COTAINER_NAME_LIST  = "count(container_cpu_usage_seconds_total{container_name!='POD',image!=''})by(namespace,pod_name,container_name)"
	PQ_COTAINER_CPU_USE      = "sum(container_cpu_usage_seconds_total{container_name!='POD',image!=''})by(namespace,pod_name,container_name)" //(MS)
	PQ_COTAINER_CPU_USAGE    = "sum(rate(container_cpu_usage_seconds_total{container_name!=''}[5m])*100)by(namespace,pod_name,container_name)"
	PQ_COTAINER_MEMORY_USE   = "sum(container_memory_working_set_bytes{container_name!=''})by(namespace,pod_name,container_name)/1024/1024"
	PQ_COTAINER_DISK_USE     = "sum(container_fs_usage_bytes{container_name!=''})by(namespace,pod_name,container_name)/1024/1024"
	PQ_COTAINER_MEMORY_USAGE = "sum(container_memory_working_set_bytes{container_name!=''}/container_memory_usage_bytes{container_name!=''}*100)by(namespace,pod_name,container_name)"

	//Pod Phase
	PQ_POD_PHASE = "count(kube_pod_status_phase>0)by(phase)"
)

type MetricsService struct {
	promethusUrl      string
	promethusRangeUrl string
	k8sApiUrl         string
}

func GetMetricsService() *MetricsService {
	config, err := util.ReadConfig(`config.ini`)
	prometheusUrl, _ := config["prometheus.addr"]
	url := prometheusUrl + SUB_URI

	k8sApiUrl, _ := config["kubernetesApi.addr"]
	k8sUrl := k8sApiUrl + K8S_SUB_URI

	rangeUrl := prometheusUrl + SUB_URI_RANGE

	if err != nil {
		log.Println(err)
	}

	return &MetricsService{
		promethusUrl:      url,
		k8sApiUrl:         k8sUrl,
		promethusRangeUrl: rangeUrl,
	}
}

func (s *MetricsService) GetClusterAvg() (model.ClusterAvg, model.ErrMessage) {
	// Metrics Call func
	podUsageData, _ := strconv.ParseFloat(GetResourceUsage(s.promethusUrl+PQ_POD_USAGE, VALUE_1_DATA), 64)
	cpuUsageData, _ := strconv.ParseFloat(GetResourceUsage(s.promethusUrl+PQ_CPU_USAGE, VALUE_1_DATA), 64)
	memoryUsageData, _ := strconv.ParseFloat(GetResourceUsage(s.promethusUrl+PQ_MEMORY_USAGE, VALUE_1_DATA), 64)
	diskUsageData, _ := strconv.ParseFloat(GetResourceUsage(s.promethusUrl+PQ_DISK_USAGE, VALUE_1_DATA), 64)

	// Struct Metrics Values Setting
	podUsage := fmt.Sprintf("%.2f", podUsageData)
	cpuUsage := fmt.Sprintf("%.2f", cpuUsageData)
	memoryUsage := fmt.Sprintf("%.2f", memoryUsageData)
	diskUsage := fmt.Sprintf("%.2f", diskUsageData)

	var dataClusterAvg model.ClusterAvg
	dataClusterAvg.PodUsage = podUsage
	dataClusterAvg.CpuUsage = cpuUsage
	dataClusterAvg.MemoryUsage = memoryUsage
	dataClusterAvg.DiskUsage = diskUsage

	return dataClusterAvg, nil
}

func (s *MetricsService) GetWorkNodeList() ([]model.WorkNodeList, model.ErrMessage) {
	// Metrics Call func
	workNodeNameList := GetWorkNodeNameList(s.promethusUrl + PQ_WORK_NODE_NAME_LIST)
	workNodeMemUsageList := GetWorkNodeMemUsageList(s.promethusUrl + PQ_WORK_NODE_MEMORY_USAGE)
	workNodeCpuUsageList := GetWorkNodeCpuUsageList(s.promethusUrl + PQ_WORK_NODE_CPU_USAGE)
	workNodeDiskUseList := GetWorkNodeDiskUseList(s.promethusUrl + PQ_WORK_NODE_DISK_USE)
	workNodeCpuUseList := GetWorkNodeCpuUseList(s.promethusUrl + PQ_WORK_NODE_CPU_ALLOC)
	workNodeMemUseList := GetWorkNodeMemUseList(s.promethusUrl + PQ_WORK_NODE_MEMORY_USE)
	workNodeConditionList := GetWorkNodeConditionList(s.promethusUrl + PQ_WORK_NODE_CONDITION)

	// Merge Maps
	workerNodeList := WorkNodeMapMerge(
		workNodeNameList,
		workNodeMemUsageList,
		workNodeCpuUsageList,
		workNodeDiskUseList,
		workNodeCpuUseList,
		workNodeMemUseList,
		workNodeConditionList)

	return workerNodeList, nil
}

func (s *MetricsService) GetWorkNodeAvg() (model.WorkNodeInfo, model.ErrMessage) {
	/*
		Make promQl

	*/
	// Metrics Call func
	cpuUsageData, _ := strconv.ParseFloat(GetResourceUsage(s.promethusUrl+PQ_CPU_USAGE, VALUE_1_DATA), 64)
	memoryUsageData, _ := strconv.ParseFloat(GetResourceUsage(s.promethusUrl+PQ_MEMORY_USAGE, VALUE_1_DATA), 64)
	diskUsageData, _ := strconv.ParseFloat(GetResourceUsage(s.promethusUrl+PQ_DISK_USAGE, VALUE_1_DATA), 64)

	// Struct Metrics Values Setting
	cpuUsage := fmt.Sprintf("%.2f", cpuUsageData)
	memoryUsage := fmt.Sprintf("%.2f", memoryUsageData)
	diskUsage := fmt.Sprintf("%.2f", diskUsageData)

	var workNodeInfo model.WorkNodeInfo
	workNodeInfo.CpuUsage = cpuUsage
	workNodeInfo.MemoryUsage = memoryUsage
	workNodeInfo.DiskUsage = diskUsage

	return workNodeInfo, nil
}

func (s *MetricsService) GetWorkNodeInfo(request model.MetricsRequest) (model.WorkNodeInfo, model.ErrMessage) {
	nodeName := request.Nodename
	instance := request.Instance

	/*
		Make promQl

	*/
	// 1.podUsage (input:nodeName)
	pqPodUsage := "sum(kube_pod_info{node='" + nodeName + "'})" +
		"/sum(kube_node_status_allocatable_pods{node='" + nodeName + "'})"

	// 2.cpuUsage (input:Instance)
	pqCpuUsage := "instance:node_cpu:ratio{instance='" + instance + "'}*100"

	// 3.memoryUsage (input:Instance)
	pqMemoryUsage := "(sum(node_memory_MemTotal_bytes{instance='" + instance + "'})-sum(node_memory_MemFree_bytes{instance='" + instance + "'}" + PLUS + "node_memory_Buffers_bytes{instance='" + instance + "'}" + PLUS + "node_memory_Cached_bytes{instance='" + instance + "'}))/sum(node_memory_MemTotal_bytes{instance='" + instance + "'})*100"

	// 4.diskUsage (input:Instance)
	var pqDiskUsage = "(sum(node_filesystem_files{instance='" + instance + "'})-sum(node_filesystem_files_free{instance='" + instance + "'}))/sum(node_filesystem_files{instance='" + instance + "'})*100"

	podUsageData, _ := strconv.ParseFloat(GetResourceUsage(s.promethusUrl+pqPodUsage, VALUE_1_DATA), 64)
	cpuUsageData, _ := strconv.ParseFloat(GetResourceUsage(s.promethusUrl+pqCpuUsage, VALUE_1_DATA), 64)
	memoryUsageData, _ := strconv.ParseFloat(GetResourceUsage(s.promethusUrl+pqMemoryUsage, VALUE_1_DATA), 64)
	diskUsageData, _ := strconv.ParseFloat(GetResourceUsage(s.promethusUrl+pqDiskUsage, VALUE_1_DATA), 64)

	// Struct Metrics Values Setting
	podUsage := fmt.Sprintf("%.2f", podUsageData)
	cpuUsage := fmt.Sprintf("%.2f", cpuUsageData)
	memoryUsage := fmt.Sprintf("%.2f", memoryUsageData)
	diskUsage := fmt.Sprintf("%.2f", diskUsageData)

	var workNodeInfo model.WorkNodeInfo
	workNodeInfo.PodUsage = podUsage
	workNodeInfo.CpuUsage = cpuUsage
	workNodeInfo.MemoryUsage = memoryUsage
	workNodeInfo.DiskUsage = diskUsage

	return workNodeInfo, nil
}

func (s *MetricsService) GetContainerList(request model.MetricsRequest) ([]model.ContainerMetricList, model.ErrMessage) {
	// Metrics Call func
	containerNameList := GetContainerNameList(s.promethusUrl, request)
	//containerCpuUseList := GetContainerCpuUseList(s.promethusUrl + PQ_COTAINER_CPU_USE)
	containerCpuUsageList := GetContainerCpuUsageList(s.promethusUrl + PQ_COTAINER_CPU_USAGE)
	containerMemUseList := GetContainerMemUseList(s.promethusUrl + PQ_COTAINER_MEMORY_USE)
	containerMemUsageList := GetContainerMemUsageList(s.promethusUrl + PQ_COTAINER_MEMORY_USAGE)
	containerDiskUseList := GetContainerDiskUseList(s.promethusUrl + PQ_COTAINER_DISK_USE)

	contanierList := ContainerMapMerge(
		containerNameList,
		//containerCpuUseList,
		containerCpuUsageList,
		containerMemUseList,
		containerMemUsageList,
		containerDiskUseList)

	return contanierList, nil
}

func (s *MetricsService) GetContainerInfo(request model.MetricsRequest) (model.ContainerInfo, model.ErrMessage) {
	containerName := request.ContainerName
	nameSpace := request.NameSpace
	podName := request.PodName

	/*
		Make promQl

	*/
	// 1.cpuUsage (input:nodeName,nameSpace,podName)
	pqCpuUsage := "sum(rate(container_cpu_usage_seconds_total{container_name!='POD',image!='',container_name='" + containerName + "',namespace='" + nameSpace + "',pod_name='" + podName + "'}[2m]))by(namespace,pod_name,container_name)*100"

	// 2.memoryUsage (input:nodeName,nameSpace,podName)
	pqMemoryUsage := "sum(container_memory_working_set_bytes{container_name!='POD',image!='',container_name='" + containerName + "',namespace='" + nameSpace + "',pod_name='" + podName + "'})/avg(machine_memory_bytes)*100"

	// 3.diskUsage (input:nodeName,nameSpace,podName)
	pqDiskUsage :=
		"sum(container_fs_usage_bytes{container_name!='POD',image!='',container_name='" + containerName + "',namespace='" + nameSpace + "',pod_name='" + podName + "'})" +
			"/max(container_fs_limit_bytes{container_name='" + containerName + "',namespace='" + nameSpace + "',pod_name='" + podName + "'})"

	// Metrics Call func
	cpuUsageData, _ := strconv.ParseFloat(GetResourceUsage(s.promethusUrl+pqCpuUsage, VALUE_1_DATA), 64)
	memoryUsageData, _ := strconv.ParseFloat(GetResourceUsage(s.promethusUrl+pqMemoryUsage, VALUE_1_DATA), 64)
	diskUsageData, _ := strconv.ParseFloat(GetResourceUsage(s.promethusUrl+pqDiskUsage, VALUE_1_DATA), 64)

	// Struct Metrics Values Setting
	cpuUsage := fmt.Sprintf("%.2f", cpuUsageData)
	memoryUsage := fmt.Sprintf("%.2f", memoryUsageData)
	diskUsage := fmt.Sprintf("%.2f", diskUsageData)

	var containerInfo model.ContainerInfo
	containerInfo.CpuUsage = cpuUsage
	containerInfo.MemoryUsage = memoryUsage
	containerInfo.DiskUsage = diskUsage

	return containerInfo, nil
}

func (s *MetricsService) GetContainerLog(request model.MetricsRequest) (model.K8sLog, model.ErrMessage) {
	nameSpace := request.NameSpace
	podName := request.PodName

	// 1.K8S_LOG
	k8sLogUrl := "namespaces/" + nameSpace + "/pods/" + podName + "/log"

	// Metrics Call func
	k8sLog := GetContainerLog(s.k8sApiUrl + k8sLogUrl)

	var containerLog model.K8sLog
	containerLog.Log = k8sLog

	return containerLog, nil
}

func (s *MetricsService) GetClusterOverView() (model.ClusterOverview, model.ErrMessage) {
	// Metrics Call func
	alertsData := GetResourceUsage(s.promethusUrl+PQ_CLUSTER_ALERTS, VALUE_1_DATA)
	runningPodData := GetResourceUsage(s.promethusUrl+PQ_CLUSTER_RUNNING_POD, VALUE_1_DATA)
	runningContainerData := GetResourceUsage(s.promethusUrl+PQ_CLUSTER_RUNNING_CONTAINER, VALUE_1_DATA)
	podRestartData := GetResourceUsage(s.promethusUrl+PQ_CLUSTER_POD_RESTART, VALUE_1_DATA)
	nodesData := GetResourceUsage(s.promethusUrl+PQ_CLUSTER_NODES, VALUE_1_DATA)

	// Struct Metrics Values Setting
	var dataClusterOverview model.ClusterOverview
	dataClusterOverview.Alerts = alertsData
	dataClusterOverview.RunningPod = runningPodData
	dataClusterOverview.Runningcontainer = runningContainerData
	dataClusterOverview.PodRestart = podRestartData
	dataClusterOverview.Nodes = nodesData

	return dataClusterOverview, nil
}

func (s *MetricsService) GetWorkloadsStatus() ([]model.WorkloadsStatus, model.ErrMessage) {
	// Metrics Call func && Struct Metrics Values Setting
	dataWorkloadsStatus := make([]model.WorkloadsStatus, 4)
	dataWorkloadsStatus[0].Name = "Deployment"
	dataWorkloadsStatus[0].Total = GetResourceUsage(s.promethusUrl+PQ_DEPLOYMENT_TOTAL, VALUE_1_DATA)
	dataWorkloadsStatus[0].Available = GetResourceUsage(s.promethusUrl+PQ_DEPLOYMENT_AVAILABLE, VALUE_1_DATA)
	dataWorkloadsStatus[0].Unavailable = GetResourceUsage(s.promethusUrl+PQ_DEPLOYMENT_UNAVAILABLE, VALUE_1_DATA)
	dataWorkloadsStatus[0].Updated = GetResourceUsage(s.promethusUrl+PQ_DEPLOYMENT_UPDATED, VALUE_1_DATA)

	dataWorkloadsStatus[1].Name = "DaemonSet"
	dataWorkloadsStatus[1].Available = GetResourceUsage(s.promethusUrl+PQ_DAEMONSET_AVAILABLE, VALUE_1_DATA)
	dataWorkloadsStatus[1].Unavailable = GetResourceUsage(s.promethusUrl+PQ_DAEMONSET_UNAVAILABLE, VALUE_1_DATA)
	dataWorkloadsStatus[1].Ready = GetResourceUsage(s.promethusUrl+PQ_DAEMONSET_READY, VALUE_1_DATA)
	dataWorkloadsStatus[1].Misscheduled = GetResourceUsage(s.promethusUrl+PQ_DAEMONSET_MISSCHEDULED, VALUE_1_DATA)

	dataWorkloadsStatus[2].Name = "Stateful"
	dataWorkloadsStatus[2].Total = GetResourceUsage(s.promethusUrl+PQ_STATEFULSET_TOTAL, VALUE_1_DATA)
	dataWorkloadsStatus[2].Updated = GetResourceUsage(s.promethusUrl+PQ_STATEFULSET_UPDATED, VALUE_1_DATA)
	dataWorkloadsStatus[2].Ready = GetResourceUsage(s.promethusUrl+PQ_STATEFULSET_READY, VALUE_1_DATA)
	dataWorkloadsStatus[2].Revision = GetResourceUsage(s.promethusUrl+PQ_STATEFULSET_REVISION, VALUE_1_DATA)

	dataWorkloadsStatus[3].Name = "Pod"
	dataWorkloadsStatus[3].Ready = GetResourceUsage(s.promethusUrl+PQ_PODCONTAINER_READY, VALUE_1_DATA)
	dataWorkloadsStatus[3].Running = GetResourceUsage(s.promethusUrl+PQ_PODCONTAINER_RUNNING, VALUE_1_DATA)
	dataWorkloadsStatus[3].Restart = GetResourceUsage(s.promethusUrl+PQ_PODCONTAINER_RESTATS, VALUE_1_DATA)
	dataWorkloadsStatus[3].Terminated = GetResourceUsage(s.promethusUrl+PQ_PODCONTAINER_TERMINATE, VALUE_1_DATA)

	return dataWorkloadsStatus, nil
}

func (s *MetricsService) GetWorkloadsContiSummary() ([]model.WorkloadsContiSummary, model.ErrMessage) {
	// Metrics Call func && Struct Metrics Values Setting
	dataWorkloadsContiSummary := make([]model.WorkloadsContiSummary, 3)
	deplomentMetric := model.WorkloadsContiSummary{}
	statefulsetMetric := model.WorkloadsContiSummary{}
	daemonsetMetric := model.WorkloadsContiSummary{}

	//goroutine process
	var waitGroup sync.WaitGroup
	waitGroup.Add(3)

	go func(url string, workLoadName string) {
		deplomentMetric = GetWorkloadsMetrics(url, workLoadName)
		defer waitGroup.Done()
	}(s.promethusUrl, "deployment")

	go func(url string, workLoadName string) {
		statefulsetMetric = GetWorkloadsMetrics(url, workLoadName)
		defer waitGroup.Done()
	}(s.promethusUrl, "statefulset")

	go func(url string, workLoadName string) {
		daemonsetMetric = GetWorkloadsMetrics(url, workLoadName)
		defer waitGroup.Done()
	}(s.promethusUrl, "daemonset")

	waitGroup.Wait()

	dataWorkloadsContiSummary[0] = deplomentMetric
	dataWorkloadsContiSummary[1] = statefulsetMetric
	dataWorkloadsContiSummary[2] = daemonsetMetric

	return dataWorkloadsContiSummary, nil

}

func (s *MetricsService) GetWorkloadsUsage(request model.MetricsRequest) (model.ContainerInfo, model.ErrMessage) {
	workloadsName := request.WorkloadsName

	// Metrics Call func && Struct Metrics Values Setting
	dataMetric := model.WorkloadsContiSummary{}
	dataMetric = GetWorkloadsMetrics(s.promethusUrl, workloadsName)

	workloadMetric := model.ContainerInfo{}
	workloadMetric.CpuUsage = dataMetric.CpuUsage
	workloadMetric.MemoryUsage = dataMetric.MemoryUsage
	workloadMetric.DiskUsage = dataMetric.DiskUsage
	return workloadMetric, nil
}

func (s *MetricsService) GetPodStatList() (model.PodPhase, model.ErrMessage) {
	//체크해야할 POD 상태 목록
	podPhaseItem := map[int]string{
		0: "Total",
		1: "Failed",
		2: "Pending",
		3: "Running",
		4: "Succeeded",
		5: "Unknown",
	}

	dataTotal := 0
	dataFailed := "0"
	dataPending := "0"
	dataRunning := "0"
	dataSucceeded := "0"
	dataUnknown := "0"

	//// Metrics Call func
	podPhaseData := GetPodPhaseList(s.promethusUrl + PQ_POD_PHASE)
	//
	var podPhase model.PodPhase
	var tmpValue1 string
	var tmpValue2 string
	//check := true

	for _, val1 := range podPhaseData {
		for key, val2 := range podPhaseItem {
			tmpValue1 = val1["phase"]
			tmpValue2 = val1["value"]

			dataVal, err1 := strconv.Atoi(tmpValue2)
			if err1 != nil {
				log.Println(err1)
			}

			if tmpValue1 == val2 {
				if key == 1 {
					dataFailed = tmpValue2
				}

				if key == 2 {
					dataPending = tmpValue2
				}

				if key == 3 {
					dataRunning = tmpValue2
				}

				if key == 4 {
					dataSucceeded = tmpValue2
				}

				if key == 5 {
					dataUnknown = tmpValue2
				}

				dataTotal += dataVal
			}
		}
	}

	podPhase.Total = strconv.Itoa(dataTotal)
	podPhase.Failed = dataFailed
	podPhase.Pending = dataPending
	podPhase.Running = dataRunning
	podPhase.Succeeded = dataSucceeded
	podPhase.Unknown = dataUnknown

	return podPhase, nil
}

func (s *MetricsService) GetPodMetricList() ([]model.PodMetricList, model.ErrMessage) {
	// Metrics Call func
	podNameList := GetPodNameList(s.promethusUrl + PQ_POD_LIST)
	//podCpuUseList := GetPodCpuUseList(s.promethusUrl + PQ_POD_CPU_USE)
	podCpuUsageList := GetPodCpuUsageList(s.promethusUrl + PQ_POD_CPU_USAGE)
	podMemUseList := GetPodMemUseList(s.promethusUrl + PQ_POD_MEMORY_USE)
	podDiskUseList := GetPodDiskUseList(s.promethusUrl + PQ_POD_DISK_USE)
	podMemUsageList := GetPodMemUsageList(s.promethusUrl + PQ_POD_MEMORY_USAGE)
	podDiskUsageList := GetPodDiskUsageList(s.promethusUrl + PQ_POD_DISK_USAGE)

	podList := PodMapMerge(
		podNameList,
		//podCpuUseList,
		podCpuUsageList,
		podMemUseList,
		podMemUsageList,
		podDiskUseList,
		podDiskUsageList)

	return podList, nil
}

func (s *MetricsService) GetPodInfo(request model.MetricsRequest) (model.ContainerInfo, model.ErrMessage) {
	podName := request.PodName

	/*
		Make promQl

	*/
	// 1.cpuUsage (input:nodeName,nameSpace,podName)
	pqCpuUsage := "sum(rate(container_cpu_usage_seconds_total{container_name!='POD',image!='',pod_name='" + podName + "'}[2m]))by(pod_name)*100"

	// 2.memoryUsage (input:nodeName,nameSpace,podName)
	pqMemoryUsage := "sum(container_memory_working_set_bytes{container_name!='POD',image!='',pod_name='" + podName + "'})/avg(machine_memory_bytes)*100"

	// 3.diskUsage (input:nodeName,nameSpace,podName)
	pqDiskUsage :=
		"sum(container_fs_usage_bytes{container_name!='POD',image!='',pod_name='" + podName + "'})" +
			"/max(container_fs_limit_bytes{container_name!='POD',image!='',pod_name='" + podName + "'})"

	// Metrics Call func
	cpuUsageData, _ := strconv.ParseFloat(GetResourceUsage(s.promethusUrl+pqCpuUsage, VALUE_1_DATA), 64)
	memoryUsageData, _ := strconv.ParseFloat(GetResourceUsage(s.promethusUrl+pqMemoryUsage, VALUE_1_DATA), 64)
	diskUsageData, _ := strconv.ParseFloat(GetResourceUsage(s.promethusUrl+pqDiskUsage, VALUE_1_DATA), 64)

	// Struct Metrics Values Setting
	cpuUsage := fmt.Sprintf("%.2f", cpuUsageData)
	memoryUsage := fmt.Sprintf("%.2f", memoryUsageData)
	diskUsage := fmt.Sprintf("%.2f", diskUsageData)

	var podInfo model.ContainerInfo
	podInfo.CpuUsage = cpuUsage
	podInfo.MemoryUsage = memoryUsage
	podInfo.DiskUsage = diskUsage

	return podInfo, nil
}

func (s *MetricsService) GetWorkNodeInfoGraph(request model.MetricsRequest) ([]model.GraphMetric, model.ErrMessage) {
	nodeName := request.Nodename
	instance := request.Instance

	//interval : currentTime between endTime(interval Time (sec))
	//timeStep : timeSeries graph(timeStep Time (sec))
	fromToTimeParmameter := util.GetPromqlFromToParameter(3600, "600")

	/*
		Make promQl

	*/
	// 1.podUsage (input:nodeName)
	pqPodUsage := "sum(kube_pod_info{node='" + nodeName + "'})" + fromToTimeParmameter

	//// 2.cpuUsage (input:Instance)
	pqCpuUsage := "node_cpu_seconds_total{mode!='idle',job='node-exporter',instance='" + instance + "'}" + fromToTimeParmameter
	//
	//// 3.memoryUsage (input:Instance)
	pqMemoryUsage := "max(((node_memory_MemTotal_bytes{job='node-exporter',instance='" + instance + "'})))" + fromToTimeParmameter
	//
	//// 4.diskUsage (input:nodeName)
	pqDiskUsage := "sum(container_fs_usage_bytes{id='/',node='" + nodeName + "'})" + fromToTimeParmameter

	// Metrics Call func
	podGraphData := GetGraphDataMap(s.promethusRangeUrl+pqPodUsage, POD)
	cpuGraphData := GetGraphDataMap(s.promethusRangeUrl+pqCpuUsage, CPU)
	memGraphData := GetGraphDataMap(s.promethusRangeUrl+pqMemoryUsage, MEMORY)
	diskGraphData := GetGraphDataMap(s.promethusRangeUrl+pqDiskUsage, DISK)

	graphMetric := make([]model.GraphMetric, 4)
	graphMetric[0].Metric = podGraphData
	graphMetric[0].Name = "Pod"
	graphMetric[1].Metric = cpuGraphData
	graphMetric[1].Name = "Cpu"
	graphMetric[2].Metric = memGraphData
	graphMetric[2].Name = "Memory"
	graphMetric[3].Metric = diskGraphData
	graphMetric[3].Name = "Disk"

	return graphMetric, nil
}

func (s *MetricsService) GetWorkNodeInfoGraphList(request model.MetricsRequest) ([]model.GraphMetricList, model.ErrMessage) {
	nodeName := request.Nodename
	instance := request.Instance

	//interval : currentTime between endTime(interval Time (sec))
	//timeStep : timeSeries graph(timeStep Time (sec))
	fromToTimeParmameter := util.GetPromqlFromToParameter(3600, "600")

	/*
	   Make promQl

	*/
	// 1.podUsage (input:nodeName)

	pqPodUsage := "sum(kube_pod_info{node='" + nodeName + "'})" + fromToTimeParmameter

	//// 2.cpuUsage (input:Instance)
	pqCpuUsage := "node_cpu_seconds_total{mode!='idle',job='node-exporter',instance='" + instance + "'}" + fromToTimeParmameter
	//
	// 3.memoryUsage (input:Instance)
	pqMemoryUsage := "max(((node_memory_MemTotal_bytes{job='node-exporter',instance='" + instance + "'})))" + fromToTimeParmameter
	//
	// 4.diskUsage (input:nodeName)
	pqDiskUsage := "sum(container_fs_usage_bytes{id='/',node='" + nodeName + "'})" + fromToTimeParmameter

	// Metrics Call func
	podGraphData := GetGraphDataMap(s.promethusRangeUrl+pqPodUsage, POD)
	cpuGraphData := GetGraphDataMap(s.promethusRangeUrl+pqCpuUsage, CPU)
	memGraphData := GetGraphDataMap(s.promethusRangeUrl+pqMemoryUsage, MEMORY)
	diskGraphData := GetGraphDataMap(s.promethusRangeUrl+pqDiskUsage, DISK)

	//fmt.Println(podGraphData)

	len := len(podGraphData)

	graphMetricList := make([]model.GraphMetricList, len)

	for i := 0; i < len; i++ {
		time_date, _ := strconv.Atoi(podGraphData[i]["time"])
		graphMetricList[i].Time = time.Unix(int64(time_date), 0).Format("2006-01-02 15:04:05")
		graphMetricList[i].PodUsage = podGraphData[i]["usage"]
		graphMetricList[i].CpuUsage = cpuGraphData[i]["usage"]
		graphMetricList[i].MemoryUsage = memGraphData[i]["usage"]
		graphMetricList[i].DiskUsage = diskGraphData[i]["usage"]
	}

	return graphMetricList, nil
}

func (s *MetricsService) GetWorkloadsInfoGraph(request model.MetricsRequest) ([]model.GraphMetric, model.ErrMessage) {
	workloadsName := request.WorkloadsName

	//interval : currentTime between endTime(interval Time (sec))
	//timeStep : timeSeries graph(timeStep Time (sec))
	fromToTimeParmameter := util.GetPromqlFromToParameter(3600, "600")

	//goroutine setting
	runtime.GOMAXPROCS(5)
	var wm sync.WaitGroup

	//Workloads(WL) Container PromQl
	pqWLmetaDataList := "count(kube_" + workloadsName + "_metadata_generation)by(namespace," + workloadsName + ")"
	resp, err := http.Get(s.promethusUrl + pqWLmetaDataList)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	str2 := string(data)

	jsonString1 := gjson.Get(str2, "data.result.#")

	var workLoadName string
	var nameSpace string

	dataWLcpu := make([][]map[string]string, int(jsonString1.Int()))
	dataWLmemory := make([][]map[string]string, int(jsonString1.Int()))
	dataWLdisk := make([][]map[string]string, int(jsonString1.Int()))

	wm.Add(int(jsonString1.Int()))

	for i := 0; i < int(jsonString1.Int()); i++ {
		jsonData := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".metric")
		jsonDataMap := jsonData.Map()
		workLoadName = jsonDataMap[workloadsName].String()
		nameSpace = jsonDataMap["namespace"].String()

		go func(url string, nameSpace string, workLoadName string, dataWLcpu [][]map[string]string, dataWLmemory [][]map[string]string, dataWLdisk [][]map[string]string, fromToTimeParmameter string, i int) {
			defer wm.Done()
			dataWLcpu[i] = GetWorkloadsGraphCpuUse(s.promethusRangeUrl, nameSpace, workLoadName, fromToTimeParmameter)
			dataWLmemory[i] = GetWorkloadsGraphMemoryUse(s.promethusRangeUrl, nameSpace, workLoadName, fromToTimeParmameter)
			dataWLdisk[i] = GetWorkloadsGraphDiskUse(s.promethusRangeUrl, nameSpace, workLoadName, fromToTimeParmameter)
		}(s.promethusRangeUrl, nameSpace, workLoadName, dataWLcpu, dataWLmemory, dataWLdisk, fromToTimeParmameter, i)
	}
	wm.Wait()

	resultCpuData := make([]map[string]string, int(len(dataWLcpu)))
	resultMemoryData := make([]map[string]string, int(len(dataWLmemory)))
	resultDiskData := make([]map[string]string, int(len(dataWLdisk)))

	var timeData string
	var tempCpuData float64
	var tempMemoryData float64
	var tempDiskData float64

	// Metrics Map Calculation (add by index)
	//Cpu
	for i := 0; i < len(dataWLcpu); i++ {
		tempCpuData = 0
		for _, val := range dataWLcpu[i] {
			cpuData, dErr := strconv.ParseFloat((val["usage"]), 64)
			if dErr != nil {
				log.Println(dErr)
			}

			tempCpuData += cpuData
			timeData = val["time"]
		}
		tempMap := make(map[string]string)
		tempMap["usage"] = fmt.Sprintf("%.02f", tempCpuData)
		tempMap["time"] = timeData
		resultCpuData[i] = tempMap

	}
	//Memory
	for i := 0; i < len(dataWLmemory); i++ {
		tempMemoryData = 0
		for _, val := range dataWLmemory[i] {
			memData, dErr := strconv.ParseFloat((val["usage"]), 64)
			if dErr != nil {
				log.Println(dErr)
			}

			tempMemoryData += memData
			timeData = val["time"]
		}
		tempMap := make(map[string]string)
		tempMap["usage"] = fmt.Sprintf("%.02f", tempMemoryData)
		tempMap["time"] = timeData
		resultMemoryData[i] = tempMap
	}
	//Disk
	for i := 0; i < len(dataWLdisk); i++ {
		tempDiskData = 0
		for _, val := range dataWLdisk[i] {
			diskData, dErr := strconv.ParseFloat((val["usage"]), 64)
			if dErr != nil {
				log.Println(dErr)
			}

			tempDiskData += diskData
			timeData = val["time"]
		}
		tempMap := make(map[string]string)
		tempMap["usage"] = fmt.Sprintf("%.02f", tempDiskData)
		tempMap["time"] = timeData
		resultDiskData[i] = tempMap
	}

	graphMetric := make([]model.GraphMetric, 3)
	graphMetric[0].Metric = resultCpuData
	graphMetric[0].Name = "Cpu"
	graphMetric[1].Metric = resultMemoryData
	graphMetric[1].Name = "Memory"
	graphMetric[2].Metric = resultDiskData
	graphMetric[2].Name = "Disk"

	return graphMetric, nil
}

func (s *MetricsService) GetPodInfoGraph(request model.MetricsRequest) ([]model.GraphMetric, model.ErrMessage) {
	podName := request.PodName

	//interval : currentTime between endTime(interval Time (sec))
	//timeStep : timeSeries graph(timeStep Time (sec))
	fromToTimeParmameter := util.GetPromqlFromToParameter(3600, "600")

	// 1.cpuUsage (input:nodeName,nameSpace,podName)
	pqCpuUsage := "sum(container_cpu_usage_seconds_total{container_name!='POD',image!='',pod_name='" + podName + "'})" + fromToTimeParmameter

	// 2.memoryUsage (input:nodeName,nameSpace,podName)
	pqMemoryUsage := "sum(container_memory_working_set_bytes{container_name!='POD',image!='',pod_name='" + podName + "'})" + fromToTimeParmameter

	// 3.diskUsage (input:nodeName,nameSpace,podName)
	pqDiskUsage :=
		"sum(container_fs_usage_bytes{container_name!='POD',image!='',pod_name='" + podName + "'})" + fromToTimeParmameter

	// Metrics Call func
	cpuGraphData := GetGraphDataMap(s.promethusRangeUrl+pqCpuUsage, CPU)
	memGraphData := GetGraphDataMap(s.promethusRangeUrl+pqMemoryUsage, MEMORY)
	diskGraphData := GetGraphDataMap(s.promethusRangeUrl+pqDiskUsage, DISK)

	graphMetric := make([]model.GraphMetric, 3)
	graphMetric[0].Metric = cpuGraphData
	graphMetric[0].Name = "Cpu"
	graphMetric[1].Metric = memGraphData
	graphMetric[1].Name = "Memory"
	graphMetric[2].Metric = diskGraphData
	graphMetric[2].Name = "Disk"

	return graphMetric, nil
}

func (s *MetricsService) GetContainerInfoGraph(request model.MetricsRequest) ([]model.GraphMetric, model.ErrMessage) {
	containerName := request.ContainerName
	nameSpace := request.NameSpace
	podName := request.PodName

	//interval : currentTime between endTime(interval Time (sec))
	//timeStep : timeSeries graph(timeStep Time (sec))
	fromToTimeParmameter := util.GetPromqlFromToParameter(3600, "600")

	// 1.cpuUsage (input:nodeName,nameSpace,podName)
	pqCpuUsage := "sum(container_cpu_usage_seconds_total{container_name!='POD',image!='',container_name='" + containerName + "',namespace='" + nameSpace + "',pod_name='" + podName + "'})" + fromToTimeParmameter

	// 2.memoryUsage (input:nodeName,nameSpace,podName)
	pqMemoryUsage := "sum(container_memory_working_set_bytes{container_name!='POD',image!='',container_name='" + containerName + "',namespace='" + nameSpace + "',pod_name='" + podName + "'})" + fromToTimeParmameter

	// 3.diskUsage (input:nodeName,nameSpace,podName)
	pqDiskUsage :=
		"sum(container_fs_usage_bytes{container_name!='POD',image!='',container_name='" + containerName + "',namespace='" + nameSpace + "',pod_name='" + podName + "'})" + fromToTimeParmameter

	// Metrics Call func
	cpuGraphData := GetGraphDataMap(s.promethusRangeUrl+pqCpuUsage, CPU)
	memGraphData := GetGraphDataMap(s.promethusRangeUrl+pqMemoryUsage, MEMORY)
	diskGraphData := GetGraphDataMap(s.promethusRangeUrl+pqDiskUsage, DISK)

	graphMetric := make([]model.GraphMetric, 3)
	graphMetric[0].Metric = cpuGraphData
	graphMetric[0].Name = "Cpu"
	graphMetric[1].Metric = memGraphData
	graphMetric[1].Name = "Memory"
	graphMetric[2].Metric = diskGraphData
	graphMetric[2].Name = "Disk"

	return graphMetric, nil
}

//Cluster Metrics func
func GetResourceUsage(url string, jpath string) string {
	var matricValue string

	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	str2 := string(data)

	jsonString1 := gjson.Get(str2, jpath)

	matricValue = jsonString1.String()

	return matricValue
}

//WorkNode Metrics func
func GetWorkNodeNameList(url string) []map[string]string {
	//var matricValue string
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	str2 := string(data)

	jsonString1 := gjson.Get(str2, "data.result.#")

	jsonMap := make([]map[string]string, 0)

	for i := 0; i < int(jsonString1.Int()); i++ {
		tempMap := make(map[string]string)
		jsonData := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".metric")
		jsonDataMap := jsonData.Map()
		tempMap["instance"] = jsonDataMap["instance"].String()
		tempMap["namespace"] = jsonDataMap["namespace"].String()
		tempMap["nodename"] = jsonDataMap["nodename"].String()

		jsonMap = append(jsonMap, tempMap)
	}

	return jsonMap
}
func GetWorkNodeMemUsageList(url string) []map[string]string {
	//var matricValue string
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	str2 := string(data)

	jsonString1 := gjson.Get(str2, "data.result.#")

	var jsonMap []map[string]string

	for i := 0; i < int(jsonString1.Int()); i++ {
		tempMap := make(map[string]string)
		jsonData1 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".metric.instance")
		tempData1 := jsonData1.String()
		jsonData2 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".value.1")
		tempData2 := jsonData2.Float()

		tempMap["instance"] = tempData1
		tempMap["value"] = fmt.Sprintf("%.0f", tempData2)

		jsonMap = append(jsonMap, tempMap)
	}
	return jsonMap
}
func GetWorkNodeCpuUsageList(url string) []map[string]string {
	//var matricValue string
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	str2 := string(data)

	jsonString1 := gjson.Get(str2, "data.result.#")

	var jsonMap []map[string]string

	for i := 0; i < int(jsonString1.Int()); i++ {
		tempMap := make(map[string]string)
		jsonData1 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".metric.instance")
		tempData1 := jsonData1.String()
		jsonData2 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".value.1")
		tempData2 := jsonData2.Float()

		tempMap["instance"] = tempData1
		tempMap["value"] = fmt.Sprintf("%.0f", tempData2)

		jsonMap = append(jsonMap, tempMap)
	}

	return jsonMap
}
func GetWorkNodeDiskUseList(url string) []map[string]string {
	//var matricValue string
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	str2 := string(data)

	jsonString1 := gjson.Get(str2, "data.result.#")

	var jsonMap []map[string]string

	for i := 0; i < int(jsonString1.Int()); i++ {
		tempMap := make(map[string]string)
		jsonData1 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".metric.instance")
		tempData1 := jsonData1.String()
		jsonData2 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".value.1")
		tempData2 := jsonData2.String()

		tempMap["instance"] = tempData1
		tempMap["value"] = util.ConByteToGB(tempData2)

		jsonMap = append(jsonMap, tempMap)
	}

	return jsonMap
}
func GetWorkNodeCpuUseList(url string) []map[string]string {
	//var matricValue string
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	str2 := string(data)

	jsonString1 := gjson.Get(str2, "data.result.#")

	var jsonMap []map[string]string

	for i := 0; i < int(jsonString1.Int()); i++ {
		tempMap := make(map[string]string)
		jsonData1 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".metric.node")
		tempData1 := jsonData1.String()
		jsonData2 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".value.1")
		tempData2 := jsonData2.Float()

		tempMap["node"] = tempData1
		tempMap["value"] = fmt.Sprintf("%.0f", tempData2)

		jsonMap = append(jsonMap, tempMap)
	}

	return jsonMap
}
func GetWorkNodeMemUseList(url string) []map[string]string {
	//var matricValue string
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	str2 := string(data)

	jsonString1 := gjson.Get(str2, "data.result.#")

	var jsonMap []map[string]string

	for i := 0; i < int(jsonString1.Int()); i++ {
		tempMap := make(map[string]string)
		jsonData1 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".metric.instance")
		tempData1 := jsonData1.String()
		jsonData2 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".value.1")
		tempData2 := jsonData2.String()

		tempMap["instance"] = tempData1
		tempMap["value"] = util.ConByteToGB(tempData2)

		jsonMap = append(jsonMap, tempMap)
	}

	return jsonMap
}
func GetWorkNodeConditionList(url string) []map[string]string {
	//var matricValue string
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	str2 := string(data)

	jsonString1 := gjson.Get(str2, "data.result.#")

	var jsonMap []map[string]string

	for i := 0; i < int(jsonString1.Int()); i++ {
		tempMap := make(map[string]string)
		jsonData1 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".metric.node")
		tempData1 := jsonData1.String()
		jsonData2 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".value.1")
		tempData2 := jsonData2.String()

		tempMap["node"] = tempData1
		tempMap["value"] = tempData2

		jsonMap = append(jsonMap, tempMap)
	}

	return jsonMap
}
func GetGraphDataMap(url string, division string) []map[string]string {
	//var matricValue string
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	str2 := string(data)

	jsonString1 := gjson.Get(str2, "data.result.0.values.#")

	jsonMap := make([]map[string]string, 0)

	for i := 0; i < int(jsonString1.Int()); i++ {
		tempMap := make(map[string]string)
		jsonData := gjson.Get(str2, "data.result.0.values."+strconv.Itoa(i)+".0")
		jsonData1 := gjson.Get(str2, "data.result.0.values."+strconv.Itoa(i)+".1")
		tempMap["time"] = jsonData.String()
		if division == "pod" {
			tempMap["usage"] = fmt.Sprintf("%.02f", jsonData1.Float())
		} else if division == "cpu" {
			tempMap["usage"] = fmt.Sprintf("%.02f", jsonData1.Float())
		} else if division == "memory" {
			tempMap["usage"] = util.ConByteToMB(fmt.Sprintf("%.02f", jsonData1.Float()))
		} else if division == "disk" {
			tempMap["usage"] = util.ConByteToMB(fmt.Sprintf("%.02f", jsonData1.Float()))
		} else {
			tempMap["usage"] = fmt.Sprintf("%.2f", jsonData1.Float())
		}

		jsonMap = append(jsonMap, tempMap)
	}
	return jsonMap
}

//Workload Metrics func
func GetWorkloadsMetrics(url string, workloadsName string) model.WorkloadsContiSummary {
	//goroutine setting
	runtime.GOMAXPROCS(5)
	var wm sync.WaitGroup
	var workloadsMetrics model.WorkloadsContiSummary

	//Workloads(WL) Container PromQl
	pqWLmetaDataList := "count(kube_" + workloadsName + "_metadata_generation)by(namespace," + workloadsName + ")"
	resp, err := http.Get(url + pqWLmetaDataList)
	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	str2 := string(data)

	jsonString1 := gjson.Get(str2, "data.result.#")

	//var dataWLcpu float64
	var dataWLcpuUsage float64
	var dataWLmemory float64
	var dataWLmemoryUsage float64
	var dataWLdisk float64
	var workLoadName string
	var nameSpace string

	wm.Add(int(jsonString1.Int()))

	for i := 0; i < int(jsonString1.Int()); i++ {
		jsonData := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".metric")
		jsonDataMap := jsonData.Map()
		workLoadName = jsonDataMap[workloadsName].String()
		nameSpace = jsonDataMap["namespace"].String()

		go func(url string, nameSpace string, workLoadName string) {
			defer wm.Done()
			//dataWLcpu += GetDivsionContiCpuUse(url, nameSpace, workLoadName, WORKLOADS)
			dataWLcpuUsage += GetDivsionContiCpuUsage(url, nameSpace, workLoadName, WORKLOADS)
			dataWLmemory += GetDivsionContiMemoryUse(url, nameSpace, workLoadName, WORKLOADS)
			dataWLmemoryUsage += GetDivsionContiMemoryUsage(url, nameSpace, workLoadName, WORKLOADS)
			dataWLdisk += GetDivsionContiDiskUse(url, nameSpace, workLoadName, WORKLOADS)
		}(url, nameSpace, workLoadName)
	}
	wm.Wait()

	workloadsMetrics.Name = workloadsName
	//workloadsMetrics.Cpu = fmt.Sprintf("%.2f", dataWLcpu)
	workloadsMetrics.CpuUsage = fmt.Sprintf("%.2f", dataWLcpuUsage/float64(jsonString1.Int()))
	workloadsMetrics.Memory = fmt.Sprintf("%.2f", dataWLmemory)
	workloadsMetrics.Disk = fmt.Sprintf("%.2f", dataWLdisk)
	workloadsMetrics.MemoryUsage = fmt.Sprintf("%.2f", dataWLmemoryUsage/float64(jsonString1.Int()))

	dataDiskUsage := GetDivsionContiDiskUsage(url, dataWLdisk)
	workloadsMetrics.DiskUsage = fmt.Sprintf("%.2f", dataDiskUsage)

	return workloadsMetrics
}
func GetDivsionContiCpuUse(url string, namespace string, podname string, division string) float64 {
	pqUrl := "sum(container_cpu_usage_seconds_total{container_name!='POD',image!='',namespace='" + namespace + "',pod_name=~'" + podname + "-.*'})"
	dataWL, _ := strconv.ParseFloat(GetResourceUsage(url+pqUrl, VALUE_1_DATA), 64)
	return dataWL
}
func GetDivsionContiMemoryUse(url string, namespace string, podname string, division string) float64 {
	var pqUrl string
	if division == "workLoads" {
		pqUrl = "sum(container_memory_working_set_bytes{container_name!='',namespace='" + namespace + "',pod_name=~'" + podname + "-.*'})"
	} else if division == "pod" {
		pqUrl = "sum(container_memory_working_set_bytes{container_name!='',pod_name=~'" + podname + "'})"
	}

	dataWL, _ := strconv.ParseFloat(GetResourceUsage(url+pqUrl, VALUE_1_DATA), 64)
	return dataWL
}
func GetDivsionContiCpuUsage(url string, namespace string, podname string, division string) float64 {
	var pqUrl string
	if division == "workLoads" {
		pqUrl = "avg(rate(container_cpu_usage_seconds_total{container_name!='',namespace='" + namespace + "',pod_name=~'" + podname + "-.*'}[5m])*100)"
	} else if division == "pod" {
		pqUrl = "avg(rate(container_cpu_usage_seconds_total{container_name!='POD',pod_name=~'" + podname + "'}[5m])*100)"
	}
	dataWL, _ := strconv.ParseFloat(GetResourceUsage(url+pqUrl, VALUE_1_DATA), 64)
	return dataWL
}
func GetDivsionContiMemoryUsage(url string, namespace string, podname string, division string) float64 {
	var pqUrl string
	if division == "workLoads" {
		pqUrl = "container_memory_working_set_bytes{container_name!='',namespace='" + namespace + "',pod_name=~'" + podname + "-.*'}/container_memory_usage_bytes{container_name!='',namespace='" + namespace + "',pod_name=~'" + podname + "-.*'}*100"
	} else if division == "pod" {
		pqUrl = "container_memory_working_set_bytes{container_name!='',pod_name=~'" + podname + "-.*'}/container_memory_usage_bytes{container_name!='',namespace='" + namespace + "',pod_name=~'" + podname + "-.*'}*100"
	}
	dataWL, _ := strconv.ParseFloat(GetResourceUsage(url+pqUrl, VALUE_1_DATA), 64)

	return dataWL
}
func GetDivsionContiDiskUse(url string, namespace string, podname string, division string) float64 {
	var pqUrl string
	if division == "workLoads" {
		pqUrl = "sum(container_fs_usage_bytes{container_name!='',namespace='" + namespace + "',pod_name=~'" + podname + "-.*'}/1024/1024)"
		fmt.Println(pqUrl)
	} else if division == "pod" {
		pqUrl = "sum(container_fs_usage_bytes{container_name!='',pod_name=~'" + podname + "-.*'}/1024/1024)"
	}
	dataWL, _ := strconv.ParseFloat(GetResourceUsage(url+pqUrl, VALUE_1_DATA), 64)
	return dataWL
}
func GetDivsionContiDiskUsage(url string, diskUse float64) float64 {
	//1.container_fs_limit_bytes
	var contLimitBytes float64
	dataContilimit := url + "sum(container_fs_limit_bytes{id='/'})"

	//var matricValue string
	resp, err := http.Get(dataContilimit)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	str2 := string(data)

	jsonString0 := gjson.Get(str2, "data.result.0.value.1")

	contLimitBytes = jsonString0.Float()

	contLimitBytes = diskUse / contLimitBytes * 100

	return contLimitBytes
}
func GetDivsionContiNameList(url string, namespace string, podname string, division string) []map[string]string {
	var pqUrl string
	if division == "workLoads" {
		pqUrl = "count(kube_pod_container_info{namespace='" + namespace + "',pod=~'" + podname + "-.*'})by(namespace,pod,container)"
	} else if division == "pod" {
		pqUrl = "count(kube_pod_container_info{pod='" + podname + "'})by(namespace,pod,container)"
	}

	resp, err := http.Get(url + pqUrl)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	str2 := string(data)

	jsonString1 := gjson.Get(str2, "data.result.#")

	jsonMap := make([]map[string]string, 0)

	for i := 0; i < int(jsonString1.Int()); i++ {
		tempMap := make(map[string]string)
		jsonData := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".metric")
		jsonDataMap := jsonData.Map()
		tempMap["namespace"] = jsonDataMap["namespace"].String()
		tempMap["podname"] = jsonDataMap["pod"].String()
		tempMap["containername"] = jsonDataMap["container"].String()

		jsonMap = append(jsonMap, tempMap)
	}

	return jsonMap
}

func GetWorkloadsGraphCpuUse(url string, namespace string, podname string, fromToTimeParmameter string) []map[string]string {
	pqUrl := "sum(container_cpu_usage_seconds_total{container_name!='POD',image!='',namespace='" + namespace + "',pod_name=~'" + podname + "-.*'})" + fromToTimeParmameter
	dataWL := GetGraphDataMap(url+pqUrl, CPU)
	return dataWL
}
func GetWorkloadsGraphMemoryUse(url string, namespace string, podname string, fromToTimeParmameter string) []map[string]string {
	pqUrl := "sum(container_memory_working_set_bytes{container_name!='POD',image!='',namespace='" + namespace + "',pod_name=~'" + podname + "-.*'})" + fromToTimeParmameter
	dataWL := GetGraphDataMap(url+pqUrl, MEMORY)
	return dataWL
}
func GetWorkloadsGraphDiskUse(url string, namespace string, podname string, fromToTimeParmameter string) []map[string]string {
	pqUrl := "sum(container_fs_usage_bytes{container_name!='POD',image!='',namespace='" + namespace + "',pod_name=~'" + podname + "-.*'})" + fromToTimeParmameter
	dataWL := GetGraphDataMap(url+pqUrl, DISK)
	return dataWL
}

//Pod metrics func
func GetPodNameList(url string) []map[string]string {
	//var matricValue string
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}
	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	str2 := string(data)

	jsonString1 := gjson.Get(str2, "data.result.#")

	jsonMap := make([]map[string]string, 0)

	for i := 0; i < int(jsonString1.Int()); i++ {
		tempMap := make(map[string]string)
		jsonData := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".metric.pod_name")
		tempMap["podname"] = jsonData.String()
		jsonMap = append(jsonMap, tempMap)
	}

	return jsonMap
}
func GetPodCpuUseList(url string) []map[string]string {
	//var matricValue string
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	str2 := string(data)

	jsonString1 := gjson.Get(str2, "data.result.#")

	jsonMap := make([]map[string]string, 0)

	for i := 0; i < int(jsonString1.Int()); i++ {
		tempMap := make(map[string]string)
		jsonData := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".metric")
		jsonDataMap := jsonData.Map()
		jsonData1 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".value.1")
		tempMap["podname"] = jsonDataMap["pod_name"].String()
		tempMap["value"] = fmt.Sprintf("%.2f", jsonData1.Float())

		jsonMap = append(jsonMap, tempMap)
	}
	return jsonMap
}
func GetPodCpuUsageList(url string) []map[string]string {
	//var matricValue string
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	str2 := string(data)

	jsonString1 := gjson.Get(str2, "data.result.#")

	jsonMap := make([]map[string]string, 0)

	for i := 0; i < int(jsonString1.Int()); i++ {
		tempMap := make(map[string]string)
		jsonData := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".metric")
		jsonDataMap := jsonData.Map()
		jsonData1 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".value.1")
		tempMap["podname"] = jsonDataMap["pod_name"].String()
		tempMap["value"] = fmt.Sprintf("%.2f", jsonData1.Float())

		jsonMap = append(jsonMap, tempMap)
	}

	return jsonMap
}
func GetPodMemUseList(url string) []map[string]string {
	//var matricValue string
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	str2 := string(data)

	jsonString1 := gjson.Get(str2, "data.result.#")

	jsonMap := make([]map[string]string, 0)

	for i := 0; i < int(jsonString1.Int()); i++ {
		tempMap := make(map[string]string)
		jsonData := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".metric")
		jsonDataMap := jsonData.Map()
		jsonData1 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".value.1")
		tempMap["podname"] = jsonDataMap["pod_name"].String()
		tempMap["value"] = fmt.Sprintf("%.2f", jsonData1.Float())
		jsonMap = append(jsonMap, tempMap)
	}
	return jsonMap
}
func GetPodMemUsageList(url string) []map[string]string {
	//var matricValue string
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	str2 := string(data)

	jsonString1 := gjson.Get(str2, "data.result.#")

	jsonMap := make([]map[string]string, 0)

	for i := 0; i < int(jsonString1.Int()); i++ {
		tempMap := make(map[string]string)
		jsonData := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".metric")
		jsonDataMap := jsonData.Map()
		jsonData1 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".value.1")
		tempMap["podname"] = jsonDataMap["pod_name"].String()
		tempMap["value"] = fmt.Sprintf("%.2f", jsonData1.Float())
		jsonMap = append(jsonMap, tempMap)
	}
	return jsonMap
}
func GetPodDiskUseList(url string) []map[string]string {
	//var matricValue string
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	str2 := string(data)

	jsonString1 := gjson.Get(str2, "data.result.#")

	jsonMap := make([]map[string]string, 0)

	for i := 0; i < int(jsonString1.Int()); i++ {
		tempMap := make(map[string]string)
		jsonData := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".metric")
		jsonDataMap := jsonData.Map()
		jsonData1 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".value.1")
		tempMap["podname"] = jsonDataMap["pod_name"].String()
		tempMap["value"] = fmt.Sprintf("%.2f", jsonData1.Float())

		jsonMap = append(jsonMap, tempMap)
	}
	return jsonMap
}
func GetPodPhaseList(url string) []map[string]string {
	//var matricValue string
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	str2 := string(data)

	jsonString1 := gjson.Get(str2, "data.result.#")

	var jsonMap []map[string]string

	for i := 0; i < int(jsonString1.Int()); i++ {
		tempMap := make(map[string]string)
		jsonData1 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".metric.phase")
		tempData1 := jsonData1.String()
		jsonData2 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".value.1")
		tempData2 := jsonData2.Float()

		tempMap["phase"] = tempData1
		tempMap["value"] = fmt.Sprintf("%.0f", tempData2)

		jsonMap = append(jsonMap, tempMap)
	}
	return jsonMap
}
func GetPodDiskUsageList(url string) []map[string]string {
	//var matricValue string
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	str2 := string(data)

	jsonString1 := gjson.Get(str2, "data.result.#")

	jsonMap := make([]map[string]string, 0)

	for i := 0; i < int(jsonString1.Int()); i++ {
		tempMap := make(map[string]string)
		jsonData := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".metric")
		jsonDataMap := jsonData.Map()
		jsonData1 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".value.1")
		tempMap["podname"] = jsonDataMap["pod_name"].String()
		tempMap["value"] = util.ConByteToMB(jsonData1.String())
		jsonMap = append(jsonMap, tempMap)
	}
	return jsonMap
}

//Container Metrics func
func GetContainerNameList(url string, request model.MetricsRequest) []map[string]string {
	//	PQ_COTAINER_NAME_LIST  = "count(container_cpu_usage_seconds_total{container_name!='POD',image!=''})by(namespace,pod_name,container_name)"
	workloadName := request.WorkloadsName
	podName := request.PodName

	jsonMap := make([]map[string]string, 0)

	//파라메터 종류에 따라 분기(WorkloadName, PodName)
	if len(strings.TrimSpace(workloadName)) != 0 {
		//goroutine setting
		runtime.GOMAXPROCS(5)
		var wm sync.WaitGroup
		//var workloadsMetrics model.WorkloadsContiSummary

		//Workloads(WL) Container PromQl
		pqWLmetaDataList := "count(kube_" + workloadName + "_metadata_generation)by(namespace," + workloadName + ")"
		resp, err := http.Get(url + pqWLmetaDataList)

		if err != nil {
			log.Println(err)
		}
		//defer resp.Body.Close()

		// 결과 출력
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		str2 := string(data)

		jsonString1 := gjson.Get(str2, "data.result.#")

		var workLoadName string
		var nameSpace string

		wm.Add(int(jsonString1.Int()))

		for i := 0; i < int(jsonString1.Int()); i++ {
			jsonData := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".metric")
			jsonDataMap := jsonData.Map()
			workLoadName = jsonDataMap[workloadName].String()
			nameSpace = jsonDataMap["namespace"].String()

			go func(workLoadName string, nameSpace string) {
				jsonMap = append(jsonMap, GetDivsionContiNameList(url, nameSpace, workLoadName, WORKLOADS)...)
				wm.Done()
			}(workLoadName, nameSpace)
		}
		wm.Wait()
	} else if len(strings.TrimSpace(podName)) != 0 {
		jsonMap = append(jsonMap, GetDivsionContiNameList(url, "", podName, POD)...)
	}

	return jsonMap
}
func GetContainerCpuUseList(url string) []map[string]string {
	//var matricValue string
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	str2 := string(data)

	jsonString1 := gjson.Get(str2, "data.result.#")

	jsonMap := make([]map[string]string, 0)

	for i := 0; i < int(jsonString1.Int()); i++ {
		tempMap := make(map[string]string)
		jsonData := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".metric")
		jsonDataMap := jsonData.Map()
		jsonData1 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".value.1")
		tempMap["namespace"] = jsonDataMap["namespace"].String()
		tempMap["podname"] = jsonDataMap["pod_name"].String()
		tempMap["containername"] = jsonDataMap["container_name"].String()
		tempMap["value"] = fmt.Sprintf("%.2f", jsonData1.Float())

		jsonMap = append(jsonMap, tempMap)
	}
	return jsonMap
}
func GetContainerCpuUsageList(url string) []map[string]string {
	//var matricValue string
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	str2 := string(data)

	jsonString1 := gjson.Get(str2, "data.result.#")

	jsonMap := make([]map[string]string, 0)

	for i := 0; i < int(jsonString1.Int()); i++ {
		tempMap := make(map[string]string)
		jsonData := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".metric")
		jsonDataMap := jsonData.Map()
		jsonData1 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".value.1")
		tempMap["namespace"] = jsonDataMap["namespace"].String()
		tempMap["podname"] = jsonDataMap["pod_name"].String()
		tempMap["containername"] = jsonDataMap["container_name"].String()
		tempMap["value"] = fmt.Sprintf("%.2f", jsonData1.Float())

		jsonMap = append(jsonMap, tempMap)
	}

	return jsonMap
}
func GetContainerMemUseList(url string) []map[string]string {
	//var matricValue string
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	str2 := string(data)

	jsonString1 := gjson.Get(str2, "data.result.#")

	jsonMap := make([]map[string]string, 0)

	for i := 0; i < int(jsonString1.Int()); i++ {
		tempMap := make(map[string]string)
		jsonData := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".metric")
		jsonDataMap := jsonData.Map()
		jsonData1 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".value.1")
		tempMap["namespace"] = jsonDataMap["namespace"].String()
		tempMap["podname"] = jsonDataMap["pod_name"].String()
		tempMap["containername"] = jsonDataMap["container_name"].String()
		//		tempMap["value"] =	jsonData1.String()
		tempMap["value"] = fmt.Sprintf("%.2f", jsonData1.Float())
		jsonMap = append(jsonMap, tempMap)
	}
	return jsonMap
}
func GetContainerMemUsageList(url string) []map[string]string {
	//var matricValue string
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	str2 := string(data)

	jsonString1 := gjson.Get(str2, "data.result.#")

	jsonMap := make([]map[string]string, 0)

	for i := 0; i < int(jsonString1.Int()); i++ {
		tempMap := make(map[string]string)
		jsonData := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".metric")
		jsonDataMap := jsonData.Map()
		jsonData1 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".value.1")
		tempMap["namespace"] = jsonDataMap["namespace"].String()
		tempMap["podname"] = jsonDataMap["pod_name"].String()
		tempMap["containername"] = jsonDataMap["container_name"].String()
		tempMap["value"] = fmt.Sprintf("%.2f", jsonData1.Float())
		jsonMap = append(jsonMap, tempMap)
	}
	return jsonMap
}
func GetContainerDiskUseList(url string) []map[string]string {
	//var matricValue string
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	str2 := string(data)

	jsonString1 := gjson.Get(str2, "data.result.#")

	jsonMap := make([]map[string]string, 0)

	for i := 0; i < int(jsonString1.Int()); i++ {
		tempMap := make(map[string]string)
		jsonData := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".metric")
		jsonDataMap := jsonData.Map()
		jsonData1 := gjson.Get(str2, "data.result."+strconv.Itoa(i)+".value.1")
		tempMap["namespace"] = jsonDataMap["namespace"].String()
		tempMap["podname"] = jsonDataMap["pod_name"].String()
		tempMap["containername"] = jsonDataMap["container_name"].String()
		tempMap["value"] = fmt.Sprintf("%.2f", jsonData1.Float())

		jsonMap = append(jsonMap, tempMap)
	}
	return jsonMap
}
func GetContainerLog(url string) string {
	var metricLog string

	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	metricLog = string(data)

	return metricLog
}

func WorkNodeMapMerge(
	workNodeNameList []map[string]string,
	workNodeMemUsageList []map[string]string,
	workNodeCpuUsageList []map[string]string,
	workNodeDiskUseList []map[string]string,
	workNodeCpuUseList []map[string]string,
	workNodeMemUseList []map[string]string,
	workNodeConditionList []map[string]string) []model.WorkNodeList {

	//workNode := model.WorkNode{}
	var workNodeList []model.WorkNodeList
	workNodeList = make([]model.WorkNodeList, len(workNodeNameList))

	for idx, data := range workNodeNameList {
		workNodeList[idx].Instance = data["instance"]
		workNodeList[idx].NodeName = data["nodename"]
		workNodeList[idx].NameSpace = data["namespace"]
	}

	for i := 0; i < len(workNodeList); i++ {
		dataInstance := workNodeList[i].Instance
		dataNodeName := workNodeList[i].NodeName
		for _, data := range workNodeMemUsageList {
			if strings.Compare(dataInstance, data["instance"]) == 0 {
				workNodeList[i].MemoryUsage = data["value"]
			}
		}

		//NodeCpuUsage
		for _, data := range workNodeCpuUsageList {
			if strings.Compare(dataInstance, data["instance"]) == 0 {
				workNodeList[i].CpuUsage = data["value"]
			}
		}

		//NodeDiskUsage
		for _, data := range workNodeDiskUseList {
			if strings.Compare(dataInstance, data["instance"]) == 0 {
				workNodeList[i].Disk = data["value"]
			}
		}

		//NodeCpuUse
		for _, data := range workNodeCpuUseList {
			if strings.Compare(dataNodeName, data["node"]) == 0 {
				workNodeList[i].Cpu = data["value"]
			}
		}

		//NodeMemUse
		for _, data := range workNodeMemUseList {
			if strings.Compare(dataInstance, data["instance"]) == 0 {
				workNodeList[i].Memory = data["value"]
			}
		}

		//NodeConditionReady(true, false)
		for _, data := range workNodeConditionList {
			if strings.Contains(dataNodeName, data["node"]) {
				workNodeList[i].Ready = "TRUE"
			}
		}
	}
	//
	//workNode.WorkNode = make([]model.WorkNodeList, len(workNodeList))
	//for i := 0; i < len(workNodeList); i++ {
	//	workNode.WorkNode[i] = workNodeList[i]
	//}

	return workNodeList
}

func ContainerMapMerge(
	containerNameList []map[string]string,
	//	containerCpuUseList []map[string]string,
	containerCpuUsageList []map[string]string,
	containerMemUseList []map[string]string,
	containerMemUsageList []map[string]string,
	containerDiskUseList []map[string]string) []model.ContainerMetricList {
	//
	//containerMetric := model.ContainerMetric{}
	var containerMetricList []model.ContainerMetricList

	containerMetricList = make([]model.ContainerMetricList, len(containerNameList))

	for idx, data := range containerNameList {
		containerMetricList[idx].NameSpace = data["namespace"]
		containerMetricList[idx].PodName = data["podname"]
		containerMetricList[idx].ContainerName = data["containername"]
	}

	for i := 0; i < len(containerMetricList); i++ {
		nameSpace := containerMetricList[i].NameSpace
		podName := containerMetricList[i].PodName
		containerName := containerMetricList[i].ContainerName
		//
		//for _, data := range containerCpuUseList {
		//	if (strings.Compare(nameSpace, data["namespace"]) == 0) && (strings.Compare(podName, data["podname"]) == 0) && (strings.Compare(containerName, data["containername"]) == 0) {
		//		containerMetricList[i].Cpu = data["value"]
		//	}
		//}

		for _, data := range containerCpuUsageList {
			if (strings.Compare(nameSpace, data["namespace"]) == 0) && (strings.Compare(podName, data["podname"]) == 0) && (strings.Compare(containerName, data["containername"]) == 0) {
				containerMetricList[i].CpuUsage = data["value"]
			}
		}

		for _, data := range containerMemUseList {
			if (strings.Compare(nameSpace, data["namespace"]) == 0) && (strings.Compare(podName, data["podname"]) == 0) && (strings.Compare(containerName, data["containername"]) == 0) {
				containerMetricList[i].Memory = data["value"]
			}
		}

		for _, data := range containerMemUsageList {
			if (strings.Compare(nameSpace, data["namespace"]) == 0) && (strings.Compare(podName, data["podname"]) == 0) && (strings.Compare(containerName, data["containername"]) == 0) {
				containerMetricList[i].MemoryUsage = data["value"]
			}
		}

		for _, data := range containerDiskUseList {
			if (strings.Compare(nameSpace, data["namespace"]) == 0) && (strings.Compare(podName, data["podname"]) == 0) && (strings.Compare(containerName, data["containername"]) == 0) {
				containerMetricList[i].Disk = data["value"]
			}
		}
	}

	//containerMetric.ContainerMetric = make([]model.ContainerMetricList, len(containerMetricList))
	//for i := 0; i < len(containerMetricList); i++ {
	//	containerMetric.ContainerMetric[i] = containerMetricList[i]
	//}

	return containerMetricList
}

func PodMapMerge(
	podNameList []map[string]string,
	//podCpuUseList []map[string]string,
	podCpuUsageList []map[string]string,
	podMemUseList []map[string]string,
	podMemUsageList []map[string]string,
	podDiskUseList []map[string]string,
	podDiskUsageList []map[string]string) []model.PodMetricList {

	var podMetricList []model.PodMetricList

	podMetricList = make([]model.PodMetricList, len(podNameList))

	for idx, data := range podNameList {
		podMetricList[idx].PodName = data["podname"]
	}

	for i := 0; i < len(podMetricList); i++ {
		podName := podMetricList[i].PodName

		//for _, data := range podCpuUseList {
		//	if strings.Compare(podName, data["podname"]) == 0 {
		//		podMetricList[i].Cpu = data["value"]
		//	}
		//}

		for _, data := range podCpuUsageList {
			if strings.Compare(podName, data["podname"]) == 0 {
				podMetricList[i].CpuUsage = data["value"]
			}
		}

		for _, data := range podMemUseList {
			if strings.Compare(podName, data["podname"]) == 0 {
				podMetricList[i].Memory = data["value"]
			}
		}

		for _, data := range podMemUsageList {
			if strings.Compare(podName, data["podname"]) == 0 {
				podMetricList[i].MemoryUsage = data["value"]
			}
		}

		for _, data := range podDiskUseList {
			if strings.Compare(podName, data["podname"]) == 0 {
				podMetricList[i].Disk = data["value"]
			}
		}

		for _, data := range podDiskUsageList {
			if strings.Compare(podName, data["podname"]) == 0 {
				podMetricList[i].DiskUsage = data["value"]
			}
		}
	}
	return podMetricList
}

package service

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"kr/paasta/monitoring/caas/model"
	"kr/paasta/monitoring/caas/util"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	// metricUrl
	SUB_URI     = "/api/v1/query?query="
	K8S_SUB_URI = "/api/v1/"

	//Jpath Type
	VALUE_0_DATA           = "data.result.0.value.0"
	VALUE_1_DATA           = "data.result.0.value.1"
	VALUE_2_DATA           = "data.result.0.value.2"
	RESULT_0_METRIC_0_DATA = "data.result.0.metric.0"

	//(PromQl)
	//Cluster Usage Metrics
	PQ_POD_USAGE    = "sum(kube_pod_info)/sum(kube_node_status_allocatable_pods{node=~'.*'})"
	PQ_CPU_USAGE    = "sum(rate(container_cpu_usage_seconds_total{id='/'}[10m]))/sum(machine_cpu_cores)*100"
	PQ_MEMORY_USAGE = "sum(container_memory_working_set_bytes{id='/'})/sum(machine_memory_bytes)*100"
	PQ_DISK_USAGE   = "sum(container_fs_usage_bytes{id='/'})/sum(container_fs_limit_bytes{id='/'})*100"

	//Work Node Usage Metrics
	PQ_WORK_NODE_NAME_LIST    = "count(node_uname_info)by(instance,nodename,namespace)"
	PQ_WORK_NODE_CPU_USAGE    = "(sum(irate(node_cpu_seconds_total{mode!='idle',job='node-exporter'}[2m]))by(instance))*100"
	PQ_WORK_NODE_MEMORY_USAGE = "max(((node_memory_MemTotal_bytes{job='node-exporter'}-" +
		"node_memory_MemFree_bytes{job='node-exporter'}" +
		"-node_memory_Buffers_bytes{job='node-exporter'}" +
		"-node_memory_Cached_bytes{job='node-exporter'})" +
		"/node_memory_MemTotal_bytes{job='node-exporter'})*100)by(instance)"
	PQ_WORK_NODE_CPU_USE    = "avg(node_cpu_seconds_total{job='node-exporter',mode!='idle'})by(instance)"
	PQ_WORK_NODE_MEMORY_USE = "sum(node_memory_MemTotal_bytes{job='node-exporter'})by(instance)"
	PQ_WORK_NODE_DISK_USE   = "sum(node_filesystem_size_bytes{job='node-exporter'})by(instance)"
	PQ_WORK_NODE_CONDITION  = "count(kube_node_status_condition{condition='Ready',status='true'})by(node)"

	//Container Usage Metrics
	PQ_COTAINER_NAME_LIST  = "count(container_cpu_usage_seconds_total{container_name!='POD',image!=''})by(namespace,pod_name,container_name)"
	PQ_COTAINER_CPU_USAGE  = "sum(rate(container_cpu_usage_seconds_total{container_name!='POD',image!=''}[2m]))by(namespace,pod_name,container_name)*100"
	PQ_COTAINER_CPU_USE    = "sum(container_cpu_usage_seconds_total{container_name!='POD',image!=''})by(namespace,pod_name,container_name)"  //(MS)
	PQ_COTAINER_MEMORY_USE = "sum(container_memory_working_set_bytes{container_name!='POD',image!=''})by(namespace,pod_name,container_name)" //(MB)
	PQ_COTAINER_DISK_USE   = "sum(container_fs_usage_bytes{container_name!='POD',image!=''})by(namespace,pod_name,container_name)"
	//PQ_COTAINER_MEMORY_USAGE
	//sum(container_memory_working_set_bytes{container_name!='POD',image!=''}/avg(machine_memory_bytes)*100"
)

type MetricsService struct {
	promethusUrl string
	k8sApiUrl    string
}

func GetMetricsService() *MetricsService {
	config, err := util.ReadConfig(`config.ini`)
	prometheusUrl, _ := config["prometheus.addr"]
	url := prometheusUrl + SUB_URI

	k8sApiUrl, _ := config["kubernetesApi.addr"]
	k8sUrl := k8sApiUrl + K8S_SUB_URI

	if err != nil {
		log.Println(err)
	}

	return &MetricsService{
		promethusUrl: url,
		k8sApiUrl:    k8sUrl,
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

func (s *MetricsService) GetWorkNodeList() (model.WorkNode, model.ErrMessage) {
	// Metrics Call func
	workNodeNameList := GetWorkNodeNameList(s.promethusUrl + PQ_WORK_NODE_NAME_LIST)
	workNodeMemUsageList := GetWorkNodeMemUsageList(s.promethusUrl + PQ_WORK_NODE_MEMORY_USAGE)
	workNodeCpuUsageList := GetWorkNodeCpuUsageList(s.promethusUrl + PQ_WORK_NODE_CPU_USAGE)
	workNodeDiskUseList := GetWorkNodeDiskUseList(s.promethusUrl + PQ_WORK_NODE_DISK_USE)
	workNodeCpuUseList := GetWorkNodeCpuUseList(s.promethusUrl + PQ_WORK_NODE_CPU_USE)
	workNodeMemUseList := GetWorkNodeMemUseList(s.promethusUrl + PQ_WORK_NODE_MEMORY_USE)
	workNodeConditionList := GetWorkNodeConditionList(s.promethusUrl + PQ_WORK_NODE_CONDITION)

	// Merge Maps
	workerNodeList := mergeMap(
		workNodeNameList,
		workNodeMemUsageList,
		workNodeCpuUsageList,
		workNodeDiskUseList,
		workNodeCpuUseList,
		workNodeMemUseList,
		workNodeConditionList)

	return workerNodeList, nil
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
	pqCpuUsage := "(sum(irate(node_cpu_seconds_total{mode!='idle',job='node-exporter',instance='" + instance + "'}[2m])))*100"

	// 3.memoryUsage (input:Instance)
	pqMemoryUsage :=
		"max(((node_memory_MemTotal_bytes{job='node-exporter',instance='" + instance + "'}-" +
			"node_memory_MemFree_bytes{job='node-exporter',instance='" + instance + "'}" +
			"-node_memory_Buffers_bytes{job='node-exporter',instance='" + instance + "'}" +
			"-node_memory_Cached_bytes{job='node-exporter',instance='" + instance + "'})" +
			"/node_memory_MemTotal_bytes{job='node-exporter',instance='" + instance + "'})*100)"

	// 4.diskUsage (input:nodeName)
	pqDiskUsage :=
		"sum(container_fs_usage_bytes{id='/',node='" + nodeName + "'})" +
			"/sum(container_fs_limit_bytes{id='/',node='" + nodeName + "'})*100"

	// Metrics Call func
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

func (s *MetricsService) GetContainerList() (model.ContainerMetric, model.ErrMessage) {
	// Metrics Call func
	containerNameList := GetContainerNameList(s.promethusUrl + PQ_COTAINER_NAME_LIST)
	containerCpuUseList := GetContainerCpuUseList(s.promethusUrl + PQ_COTAINER_CPU_USE)
	containerCpuUsageList := GetContainerCpuUsageList(s.promethusUrl + PQ_COTAINER_CPU_USAGE)
	containerMemUseList := GetContainerMemUseList(s.promethusUrl + PQ_COTAINER_MEMORY_USE)
	containerMemUsageList := GetContainerMemUsageList(s.promethusUrl)
	containerDiskUseList := GetContainerDiskUseList(s.promethusUrl + PQ_COTAINER_DISK_USE)

	contanierList := mergeMap2(
		containerNameList,
		containerCpuUseList,
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

func GetContainerNameList(url string) []map[string]string {
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
		tempMap["namespace"] = jsonDataMap["namespace"].String()
		tempMap["podname"] = jsonDataMap["pod_name"].String()
		tempMap["containername"] = jsonDataMap["container_name"].String()

		jsonMap = append(jsonMap, tempMap)
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
		tempMap["value"] = util.ConByteToMB(jsonData1.String())
		jsonMap = append(jsonMap, tempMap)
	}
	return jsonMap
}

func GetContainerMemUsageList(url string) []map[string]string {
	// 1.machineMem
	var machineMem float64
	machineMemUri := url + "avg(machine_memory_bytes)"

	//var matricValue string
	resp, err := http.Get(machineMemUri)

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

	machineMem = jsonString0.Float()

	// 2.memUsageUri
	memUsageUri := url + "sum(container_memory_working_set_bytes{container_name!='POD',image!=''})by(namespace,pod_name,container_name)"

	//var matricValue string
	resp1, err := http.Get(memUsageUri)

	if err != nil {
		log.Println(err)
	}

	//defer resp.Body.Close()

	// 결과 출력
	data1, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		log.Println(err)
	}

	str3 := string(data1)

	jsonString1 := gjson.Get(str3, "data.result.#")

	jsonMap := make([]map[string]string, 0)

	for i := 0; i < int(jsonString1.Int()); i++ {
		tempMap := make(map[string]string)
		jsonData := gjson.Get(str3, "data.result."+strconv.Itoa(i)+".metric")
		jsonDataMap := jsonData.Map()
		jsonData1 := gjson.Get(str3, "data.result."+strconv.Itoa(i)+".value.1")
		tempMap["namespace"] = jsonDataMap["namespace"].String()
		tempMap["podname"] = jsonDataMap["pod_name"].String()
		tempMap["containername"] = jsonDataMap["container_name"].String()
		tempMap["value"] = fmt.Sprintf("%.2f", (jsonData1.Float() / machineMem * 100))
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
		tempMap["value"] = util.ConByteToMB(jsonData1.String())

		jsonMap = append(jsonMap, tempMap)
	}
	return jsonMap
}

func GetContainerLog(url string) string {
	var metricLog string

	fmt.Println(url)

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

	fmt.Println(metricLog)

	return metricLog
}

//sub_method
func GetResourceUsage(url string, jpath string) string {
	var matricValue string

	fmt.Println(url)

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

	jsonString1 := gjson.Get(str2, jpath)

	matricValue = jsonString1.String()

	return matricValue
}

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

func mergeMap(
	workNodeNameList []map[string]string,
	workNodeMemUsageList []map[string]string,
	workNodeCpuUsageList []map[string]string,
	workNodeDiskUseList []map[string]string,
	workNodeCpuUseList []map[string]string,
	workNodeMemUseList []map[string]string,
	workNodeConditionList []map[string]string) model.WorkNode {

	workNode := model.WorkNode{}
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
			if strings.Compare(dataInstance, data["instance"]) == 0 {
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

	workNode.WorkNode = make([]model.WorkNodeList, len(workNodeList))
	for i := 0; i < len(workNodeList); i++ {
		workNode.WorkNode[i] = workNodeList[i]
	}

	return workNode
}

func mergeMap2(
	containerNameList []map[string]string,
	containerCpuUseList []map[string]string,
	containerCpuUsageList []map[string]string,
	containerMemUseList []map[string]string,
	containerMemUsageList []map[string]string,
	containerDiskUseList []map[string]string) model.ContainerMetric {

	containerMetric := model.ContainerMetric{}
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

		for _, data := range containerCpuUseList {
			if (strings.Compare(nameSpace, data["namespace"]) == 0) && (strings.Compare(podName, data["podname"]) == 0) && (strings.Compare(containerName, data["containername"]) == 0) {
				containerMetricList[i].Cpu = data["value"]
			}
		}

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

	containerMetric.ContainerMetric = make([]model.ContainerMetricList, len(containerMetricList))
	for i := 0; i < len(containerMetricList); i++ {
		containerMetric.ContainerMetric[i] = containerMetricList[i]
	}

	return containerMetric
}

package caas

import (
	"github.com/tidwall/gjson"
	"log"
	"paasta-monitoring-api/helpers"
	models "paasta-monitoring-api/models/api/v1"
	"runtime"
	"strconv"
	"strings"
	"sync"
)


type WorkloadService struct {
	CaasConfig models.CaasConfig
}

func GetWorkloadService(config models.CaasConfig) *WorkloadService{
	return &WorkloadService{
		CaasConfig: config,
	}
}


func (service *WorkloadService) GetWorkloadStatus() ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	deploymentTotalResult, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query=" + models.PROMQL_WORKLOAD_DEPLOYMENT_TOTAL)
	deploymentTotal := gjson.Get(string(deploymentTotalResult), "data.result.0.value.1").String()
	deploymentActiveResult, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query=" + models.PROMQL_WORKLOAD_DEPLOYMENT_AVAILABLE)
	deploymentActive := gjson.Get(string(deploymentActiveResult), "data.result.0.value.1").String()
	deploymentInactiveResult, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query=" + models.PROMQL_WORKLOAD_DEPLOYMENT_UNAVAILABLE)
	deploymentInactive := gjson.Get(string(deploymentInactiveResult), "data.result.0.value.1").String()
	deploymentUpdatedResult, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query=" + models.PROMQL_WORKLOAD_DEPLOYMENT_UPDATED)
	deploymentUpdated := gjson.Get(string(deploymentUpdatedResult), "data.result.0.value.1").String()
	deploymentMap := make(map[string]interface{})
	deploymentMap["name"] = "Deployment"
	deploymentMap["total"] = deploymentTotal
	deploymentMap["available"] = deploymentActive
	deploymentMap["unavailable"] = deploymentInactive
	deploymentMap["updated"] = deploymentUpdated

	statefulsetTotalResult, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query=" + models.PROMQL_WORKLOAD_STATEFULSET_TOTAL)
	statefulsetTotal := gjson.Get(string(statefulsetTotalResult), "data.result.0.value.1").String()
	statefulsetActiveResult, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query=" + models.PROMQL_WORKLOAD_STATEFULSET_AVAILABLE)
	statefulsetActive := gjson.Get(string(statefulsetActiveResult), "data.result.0.value.1").String()
	statefulsetInactiveResult, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query=" + models.PROMQL_WORKLOAD_STATEFULSET_UNAVAILABLE)
	statefulsetInactive := gjson.Get(string(statefulsetInactiveResult), "data.result.0.value.1").String()
	statefulsetUpdatedResult, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query=" + models.PROMQL_WORKLOAD_STATEFULSET_UPDATED)
	statefulsetUpdated := gjson.Get(string(statefulsetUpdatedResult), "data.result.0.value.1").String()
	statefulsetMap := make(map[string]interface{})
	statefulsetMap["name"] = "Stateful"
	statefulsetMap["total"] = statefulsetTotal
	statefulsetMap["available"] = statefulsetActive
	statefulsetMap["unavailable"] = statefulsetInactive
	statefulsetMap["updated"] = statefulsetUpdated

	daemonsetReadyResult, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query=" + models.PROMQL_WORKLOAD_DAEMONSET_READY)
	daemonsetReady := gjson.Get(string(daemonsetReadyResult), "data.result.0.value.1").String()
	daemonsetActiveResult, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query=" + models.PROMQL_WORKLOAD_DAEMONSET_AVAILABLE)
	daemonsetActive := gjson.Get(string(daemonsetActiveResult), "data.result.0.value.1").String()
	daemonsetInactiveResult, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query=" + models.PROMQL_WORKLOAD_DAEMONSET_UNAVAILABLE)
	daemonsetInactive := gjson.Get(string(daemonsetInactiveResult), "data.result.0.value.1").String()
	daemonsetMisscheduleResult, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query=" + models.PROMQL_WORKLOAD_DAEMONSET_MISSCHEDULED)
	daemonsetMisschedule := gjson.Get(string(daemonsetMisscheduleResult), "data.result.0.value.1").String()
	daemonsetMap := make(map[string]interface{})
	daemonsetMap["name"] = "DaemonSet"
	daemonsetMap["ready"] = daemonsetReady
	daemonsetMap["available"] = daemonsetActive
	daemonsetMap["unavailable"] = daemonsetInactive
	daemonsetMap["misscheduled"] = daemonsetMisschedule

	podcontainerReadyResult, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query=" + models.PROMQL_WORKLOAD_PODCONTAINER_READY)
	podcontainerReady := gjson.Get(string(podcontainerReadyResult), "data.result.0.value.1").String()
	podcontainerRunningResult, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query=" + models.PROMQL_WORKLOAD_PODCONTAINER_RUNNING)
	podcontainerRunning := gjson.Get(string(podcontainerRunningResult), "data.result.0.value.1").String()
	podcontainerRestatsResult, _ := 	helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query=" + models.PROMQL_WORKLOAD_PODCONTAINER_RESTARTS)
	podcontainerRestarts := gjson.Get(string(podcontainerRestatsResult), "data.result.0.value.1").String()
	podcontainerTerminateResult, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query=" + models.PROMQL_WORKLOAD_PODCONTAINER_TERMINATE)
	podcontainerTerminate := gjson.Get(string(podcontainerTerminateResult), "data.result.0.value.1").String()
	podcontainerMap := make(map[string]interface{})
	podcontainerMap["name"] = "Pod"
	podcontainerMap["ready"] = podcontainerReady
	podcontainerMap["running"] = podcontainerRunning
	podcontainerMap["restart"] = podcontainerRestarts
	podcontainerMap["terminated"] = podcontainerTerminate

	result = append(result, deploymentMap)
	result = append(result, statefulsetMap)
	result = append(result, daemonsetMap)
	result = append(result, podcontainerMap)

	return result, nil
}


func (service *WorkloadService) GetWorkloadList() ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	var waitGroup sync.WaitGroup
	waitGroup.Add(3)

	deplomentMetric := make(map[string]interface{})
	statefulsetMetric := make(map[string]interface{})
	daemonsetMetric := make(map[string]interface{})

	go func(url string, workLoadName string) {
		deplomentMetric = getWorkloadMetrics(url + "/query", workLoadName)
		defer waitGroup.Done()
	}(service.CaasConfig.PromethusUrl, "deployment")

	go func(url string, workLoadName string) {
		statefulsetMetric = getWorkloadMetrics(url + "/query", workLoadName)
		defer waitGroup.Done()
	}(service.CaasConfig.PromethusUrl, "statefulset")

	go func(url string, workLoadName string) {
		daemonsetMetric = getWorkloadMetrics(url + "/query", workLoadName)
		defer waitGroup.Done()
	}(service.CaasConfig.PromethusUrl, "daemonset")

	waitGroup.Wait()

	result = append(result, deplomentMetric)
	result = append(result, statefulsetMetric)
	result = append(result, daemonsetMetric)

	return result, nil
}


var metricTypes = [...]string{"cpu", "memory", "disk"}
var metricValues = [len(metricTypes)]float64{}
func (service *WorkloadService) GetWorkloadDetailMetrics(workloadParam string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	runtime.GOMAXPROCS(5)
	fromToTimeParmameter := GetPromqlFromToParameter(3600, "600")

	promqlWorkloadList := "count(kube_" + workloadParam + "_metadata_generation)by(namespace," + workloadParam + ")"
	workloadsByte, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query="+promqlWorkloadList) // Retrieve workload list per type
	workloadArray := gjson.Get(string(workloadsByte), "data.result")

	// cpu, memory, disk 배열 루프
	for _, metricType := range metricTypes {

		seriesMap := make(map[string]interface{})
		var seriesDataArr []map[string]interface{}

		// 워크로드 배열 루프
		for workloadIdx, workload := range workloadArray.Array() {
			itemMap := workload.Get("metric")
			var namespace string
			var workloadOrPod string
			if strings.Compare(workloadParam, "daemonset") == 0 {
				workloadOrPod = itemMap.Get("pod").String()
				namespace = ""
			} else {
				workloadOrPod = itemMap.Get(workloadParam).String()
				namespace = itemMap.Get("namespace").String()
			}

			promQLStr := makePromQLScriptForWorkloadMetrics(metricType, namespace, workloadOrPod, fromToTimeParmameter)
			metricsBytes, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query_range", "query="+promQLStr)  // Retrieve workload's metric data
			metricsResult := gjson.Get(string(metricsBytes), "data.result.0.values")

			for metricsIdx, item := range metricsResult.Array() {
				timestamp := item.Get("0").String()
				usage := item.Get("1").Float()

				var itemMap map[string]interface{}
				if workloadIdx == 0 {
					itemMap = make(map[string]interface{})
					itemMap["time"] = timestamp
					itemMap["usage"] = usage
					seriesDataArr = append(seriesDataArr, itemMap)
				} else {
					itemMap = seriesDataArr[metricsIdx]

					prevUsage := itemMap["usage"].(float64)
					itemMap["usage"] = prevUsage + usage
				}
			}
		}
		seriesMap["name"] = metricType
		seriesMap["metric"] = seriesDataArr

		result = append(result, seriesMap)

	}
	log.Printf("result : %v\n", result)
	return result, nil
}


func (service *WorkloadService) GetWorkloadContainerList(workloadParam string) ([]map[string]interface{}, error) {
	var containerList []map[string]interface{}

	promqlWorkloadList := "count(kube_" + workloadParam + "_metadata_generation)by(namespace," + workloadParam + ")"
	workloadsByte, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query="+promqlWorkloadList) // Retrieve workload list per type
	workloadArray := gjson.Get(string(workloadsByte), "data.result")

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(workloadArray.Array()))

	for _, item := range workloadArray.Array() {
		workload := item.Get("metric." + workloadParam).String()
		namespace := item.Get("metric.namespace").String()

		// TODO: Retrieve container list in workload
		promqlContainerList := "count(kube_pod_container_info{namespace='" + namespace + "',pod=~'" + workload + "-.*'})by(namespace,pod,container)"
		containersByte, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query="+promqlContainerList) // Retrieve container list in workload
		containerArray := gjson.Get(string(containersByte), "data.result")

		for _, container := range containerArray.Array() {
			containerMap := make(map[string]interface{})
			podName := container.Get("metric.pod").String()
			containerName := container.Get("metric.container").String()
			containerMap["namespace"] = namespace
			containerMap["pod"] = podName
			containerMap["container"] = containerName
			containerList = append(containerList, containerMap)
		}
	}

	var cpuUseList []map[string]interface{}
	promqlCpuUse := "sum(container_cpu_usage_seconds_total{container!='POD',image!=''})by(namespace,pod,container)"
	cpuUseByte, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query="+promqlCpuUse) // Retrieve container list in workload
	cpuUseArray := gjson.Get(string(cpuUseByte), "data.result")
	for _, item := range cpuUseArray.Array() {
		cpuUseMap := make(map[string]interface{})
		namespace := item.Get("metric.namespace").String()
		podName := item.Get("metric.pod").String()
		containerName := item.Get("metric.container").String()
		value := item.Get("value.1").Float()
		cpuUseMap["namespace"] = namespace
		cpuUseMap["pod"] = podName
		cpuUseMap["container"] = containerName
		cpuUseMap["value"] = value
		cpuUseList = append(cpuUseList, cpuUseMap)
	}

	var cpuUsageList []map[string]interface{}
	promqlCpuUsage := "sum(rate(container_cpu_usage_seconds_total{container!='POD',image!=''}[5m])*100)by(namespace,pod,container)"
	cpuUsageByte, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query="+promqlCpuUsage) // Retrieve container list in workload
	cpuUsageArray := gjson.Get(string(cpuUsageByte), "data.result")
	for _, item := range cpuUsageArray.Array() {
		cpuUsageMap := make(map[string]interface{})
		namespace := item.Get("metric.namespace").String()
		podName := item.Get("metric.pod").String()
		containerName := item.Get("metric.container").String()
		value := item.Get("value.1").Float()
		cpuUsageMap["namespace"] = namespace
		cpuUsageMap["pod"] = podName
		cpuUsageMap["container"] = containerName
		cpuUsageMap["value"] = value
		cpuUsageList = append(cpuUsageList, cpuUsageMap)
	}

	var memoryUseList []map[string]interface{}
	promqlMemoryUse := "sum(container_memory_working_set_bytes{container!='POD',image!=''})by(namespace,pod,container)/1024/1024"
	memoryUseByte, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query="+promqlMemoryUse) // Retrieve container list in workload
	memoryUseArray := gjson.Get(string(memoryUseByte), "data.result")
	for _, item := range memoryUseArray.Array() {
		memoryUseMap := make(map[string]interface{})
		namespace := item.Get("metric.namespace").String()
		podName := item.Get("metric.pod").String()
		containerName := item.Get("metric.container").String()
		value := item.Get("value.1").Float()
		memoryUseMap["namespace"] = namespace
		memoryUseMap["pod"] = podName
		memoryUseMap["container"] = containerName
		memoryUseMap["value"] = value
		memoryUseList = append(memoryUseList, memoryUseMap)
	}

	var memoryUsageList []map[string]interface{}
	promqlMemoryUsage := "avg(container_memory_working_set_bytes{container!='POD',image!=''})by(namespace,pod,container)/scalar(sum(machine_memory_bytes))*100*scalar(count(container_memory_usage_bytes{container!='POD',image!=''}))"
	memoryUsageByte, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query="+promqlMemoryUsage) // Retrieve container list in workload
	memoryUsageArray := gjson.Get(string(memoryUsageByte), "data.result")
	for _, item := range memoryUsageArray.Array() {
		memoryUsageMap := make(map[string]interface{})
		namespace := item.Get("metric.namespace").String()
		podName := item.Get("metric.pod").String()
		containerName := item.Get("metric.container").String()
		value := item.Get("value.1").Float()
		memoryUsageMap["namespace"] = namespace
		memoryUsageMap["pod"] = podName
		memoryUsageMap["container"] = containerName
		memoryUsageMap["value"] = value
		memoryUsageList = append(memoryUsageList, memoryUsageMap)
	}

	var diskUseList []map[string]interface{}
	promqlDiskUse := "sum(container_fs_usage_bytes{container!='POD',image!=''})by(namespace,pod,container)/1024/1024"
	diskUseByte, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query="+promqlDiskUse) // Retrieve container list in workload
	diskUseArray := gjson.Get(string(diskUseByte), "data.result")
	for _, item := range diskUseArray.Array() {
		diskUseMap := make(map[string]interface{})
		namespace := item.Get("metric.namespace").String()
		podName := item.Get("metric.pod").String()
		containerName := item.Get("metric.container").String()
		value := item.Get("value.1").Float()
		diskUseMap["namespace"] = namespace
		diskUseMap["pod"] = podName
		diskUseMap["container"] = containerName
		diskUseMap["value"] = value
		diskUseList = append(diskUseList, diskUseMap)
	}

	for _, containerMap := range containerList {
		namespace := containerMap["namespace"].(string)
		pod := containerMap["pod"].(string)
		container := containerMap["container"].(string)
		for _, cpuUseMap := range cpuUseList {
			if (strings.Compare(namespace, cpuUseMap["namespace"].(string)) == 0) &&
				(strings.Compare(pod, cpuUseMap["pod"].(string)) == 0) &&
					(strings.Compare(container, cpuUseMap["container"].(string)) == 0) {
				containerMap["cpu"] = cpuUseMap["value"].(float64)
			}
		}
		for _, cpuUsage := range cpuUsageList {
			if (strings.Compare(namespace, cpuUsage["namespace"].(string)) == 0) &&
				(strings.Compare(pod, cpuUsage["pod"].(string)) == 0) &&
				(strings.Compare(container, cpuUsage["container"].(string)) == 0) {
				containerMap["cpuUsage"] = cpuUsage["value"].(float64)
			}
		}
		for _, memoryUseMap := range memoryUseList {
			if (strings.Compare(namespace, memoryUseMap["namespace"].(string)) == 0) &&
				(strings.Compare(pod, memoryUseMap["pod"].(string)) == 0) &&
				(strings.Compare(container, memoryUseMap["container"].(string)) == 0) {
				containerMap["memory"] = memoryUseMap["value"].(float64)
			}
		}
		for _, memoryUsageMap := range memoryUsageList {
			if (strings.Compare(namespace, memoryUsageMap["namespace"].(string)) == 0) &&
				(strings.Compare(pod, memoryUsageMap["pod"].(string)) == 0) &&
				(strings.Compare(container, memoryUsageMap["container"].(string)) == 0) {
				containerMap["memoryUsage"] = memoryUsageMap["value"].(float64)
			}
		}
		for _, diskUseMap := range diskUseList {
			if (strings.Compare(namespace, diskUseMap["namespace"].(string)) == 0) &&
				(strings.Compare(pod, diskUseMap["pod"].(string)) == 0) &&
				(strings.Compare(container, diskUseMap["container"].(string)) == 0) {
				containerMap["disk"] = diskUseMap["value"].(float64)
			}
		}
	}
	return containerList, nil
}


func (service *WorkloadService) GetContainerStatus(namespace string, container string, pod string) ([]map[string]interface{}, error) {
	var resultList []map[string]interface{}

	fromToTimeParmameter := GetPromqlFromToParameter(3600, "600")
	pqCpuUsage := "sum(container_cpu_usage_seconds_total{container!='POD',image!='',container='" + container + "',namespace='" + namespace + "',pod='" + pod + "'})" + fromToTimeParmameter
	pqMemoryUsage := "sum(container_memory_working_set_bytes{container!='POD',image!='',container='" + container + "',namespace='" + namespace + "',pod='" + pod + "'})" + fromToTimeParmameter
	pqDiskUsage := "sum(container_fs_usage_bytes{container!='POD',image!='',container='" + container + "',namespace='" + namespace + "',pod='" + pod + "'})" + fromToTimeParmameter

	metricsBytes, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query_range", "query="+pqCpuUsage)  // Retrieve workload's metric data
	metricsResult := gjson.Get(string(metricsBytes), "data.result.0.values")
	var seriesDataArr []map[string]interface{}
	for _, item := range metricsResult.Array() {
		timestamp := item.Get("0").String()
		usage := item.Get("1").Float()

		var itemMap map[string]interface{}
		itemMap = make(map[string]interface{})
		itemMap["time"] = timestamp
		itemMap["usage"] = usage
		seriesDataArr = append(seriesDataArr, itemMap)
	}
	cpuMetricsMap := make(map[string]interface{})
	cpuMetricsMap["name"] = "cpu"
	cpuMetricsMap["metric"] = seriesDataArr
	resultList = append(resultList, cpuMetricsMap)
	seriesDataArr = nil

	metricsBytes, _ = helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query_range", "query="+pqMemoryUsage)  // Retrieve workload's metric data
	metricsResult = gjson.Get(string(metricsBytes), "data.result.0.values")
	for _, item := range metricsResult.Array() {
		timestamp := item.Get("0").String()
		usage := item.Get("1").Float()

		var itemMap map[string]interface{}
		itemMap = make(map[string]interface{})
		itemMap["time"] = timestamp
		itemMap["usage"] = usage
		seriesDataArr = append(seriesDataArr, itemMap)
	}
	memoryMetricsMap := make(map[string]interface{})
	memoryMetricsMap["name"] = "memory"
	memoryMetricsMap["metric"] = seriesDataArr
	resultList = append(resultList, memoryMetricsMap)
	seriesDataArr = nil
	
	metricsBytes, _ = helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query_range", "query="+pqDiskUsage)  // Retrieve workload's metric data
	metricsResult = gjson.Get(string(metricsBytes), "data.result.0.values")
	for _, item := range metricsResult.Array() {
		timestamp := item.Get("0").String()
		usage := item.Get("1").Float()

		var itemMap map[string]interface{}
		itemMap = make(map[string]interface{})
		itemMap["time"] = timestamp
		itemMap["usage"] = usage
		seriesDataArr = append(seriesDataArr, itemMap)
	}
	diskMetricsMap := make(map[string]interface{})
	diskMetricsMap["name"] = "disk"
	diskMetricsMap["metric"] = seriesDataArr
	resultList = append(resultList, diskMetricsMap)

	return resultList, nil
}


func makePromQLScriptForWorkloadMetrics(metricType string, namespace string, pod string, timeCondition string) string {
	var promQLStr string
	var promQLLabel string
	var promQLCondition string

	switch metricType {
	case "cpu":
		promQLLabel = "container_cpu_usage_seconds_total"
		break
	case "memory":
		promQLLabel = "container_memory_working_set_bytes"
		break
	case "disk":
		promQLLabel = "container_fs_usage_bytes"
		break
	}

	promQLCondition = "{container!='POD',image!='',"
	if len(namespace) > 0 {
		promQLCondition += "namespace='" + namespace + "',pod=~'" + pod + "-.*'}"
	} else {
		promQLCondition += "pod='" + pod + "'}"
	}
	promQLStr = "sum(" + promQLLabel + promQLCondition + ")" + timeCondition
	//log.Printf("promQL script : %v\n", promQLStr)

	return promQLStr
}


func getWorkloadMetrics(url string, workloadNameParam string) map[string]interface{} {
	result := make(map[string]interface{})

	runtime.GOMAXPROCS(5)
	var wg sync.WaitGroup

	promqlWorkloadList := "count(kube_" + workloadNameParam + "_metadata_generation)by(namespace," + workloadNameParam + ")"
	resultBytes, _ := helpers.RequestHttpGet(url, "query="+promqlWorkloadList)
	resultArray := gjson.Get(string(resultBytes), "data.result")

	var cpuUse float64
	var cpuUsage float64
	var memoryUse float64
	var memoryUsage float64
	var diskUse float64
	var diskUsage float64

	wg.Add(len(resultArray.Array()))

	for _, item := range resultArray.Array() {
		workload := item.Get("metric." + workloadNameParam).String()
		namespace := item.Get("metric.namespace").String()

		go func(url string, workload string, namespace string) {
			defer wg.Done()

			pqUrl := "sum(container_cpu_usage_seconds_total{container!='POD',image!='',namespace='" + namespace + "',pod=~'" + workload + "-.*'})"
			resultBytes, _ := helpers.RequestHttpGet(url, "query="+pqUrl)
			cpuUseVal, _ := strconv.ParseFloat(gjson.Get(string(resultBytes), "data.result.0.value.1").String(), 64)
			cpuUse += cpuUseVal

			pqUrl = "sum(rate(container_cpu_usage_seconds_total{container!='',namespace='" + namespace + "',pod=~'" + workload + "-.*'}[5m]))*100"
			resultBytes, _ = helpers.RequestHttpGet(url, "query="+pqUrl)
			cpuUsageVal, _ := strconv.ParseFloat(gjson.Get(string(resultBytes), "data.result.0.value.1").String(), 64)
			cpuUsage += cpuUsageVal

			pqUrl = "sum(container_memory_working_set_bytes{container!='',namespace='" + namespace + "',pod=~'" + workload + "-.*'})/1024/1024"
			resultBytes, _ = helpers.RequestHttpGet(url, "query="+pqUrl)
			memoryUseVal, _ := strconv.ParseFloat(gjson.Get(string(resultBytes), "data.result.0.value.1").String(), 64)
			memoryUse += memoryUseVal

			pqUrl = "avg(container_memory_working_set_bytes{container!='',namespace='" + namespace + "',pod=~'" + workload + "-.*'})/scalar(sum(machine_memory_bytes))*100*scalar(count(container_memory_usage_bytes{container!='POD',image!=''}))"
			resultBytes, _ = helpers.RequestHttpGet(url, "query="+pqUrl)
			memoryUsageVal, _ := strconv.ParseFloat(gjson.Get(string(resultBytes), "data.result.0.value.1").String(), 64)
			memoryUsage += memoryUsageVal

			pqUrl = "sum(container_fs_usage_bytes{container!='',namespace='" + namespace + "',pod=~'" + workload + "-.*'}/1024/1024)"
			resultBytes, _ = helpers.RequestHttpGet(url, "query="+pqUrl)
			diskUseVal, _ := strconv.ParseFloat(gjson.Get(string(resultBytes), "data.result.0.value.1").String(), 64)
			diskUse += diskUseVal

			pqUrl = "sum(container_fs_usage_bytes{container!='POD',image!=''.container!='',namespace='" + namespace + "',pod=~'" + workload + "-.*'})by(pod)/max(container_fs_limit_bytes{container!='POD',image!='',container!='',namespace='" + namespace + "',pod=~'" + workload + "-.*'})by(pod)*100"
			resultBytes, _ = helpers.RequestHttpGet(url, "query="+pqUrl)
			diskUsageVal, _ := strconv.ParseFloat(gjson.Get(string(resultBytes), "data.result.0.value.1").String(), 64)
			diskUsage += diskUsageVal

		}(url, workload, namespace)
	}
	wg.Wait()

	result["name"] = workloadNameParam
	result["cpu"] = cpuUse
	result["cpuUsage"] = cpuUsage
	result["memory"] = memoryUse
	result["memoryUsage"] = memoryUsage
	result["disk"] = diskUse
	result["diskUsage"] = diskUsage

	return result
}



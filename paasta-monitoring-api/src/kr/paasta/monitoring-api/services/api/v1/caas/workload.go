package caas

import (
	"github.com/tidwall/gjson"
	"paasta-monitoring-api/helpers"
	models "paasta-monitoring-api/models/api/v1"
	"runtime"
	"strconv"
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


func (service *WorkloadService) GetWorkloadContainerList() ([]map[string]interface{}, error) {
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

	return result;
}
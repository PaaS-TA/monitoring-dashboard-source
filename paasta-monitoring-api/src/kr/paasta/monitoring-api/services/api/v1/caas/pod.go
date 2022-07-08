package caas

import (
	"github.com/tidwall/gjson"
	"paasta-monitoring-api/helpers"
	models "paasta-monitoring-api/models/api/v1"
	"strings"
)

type PodService struct {
	CaaS models.CaaS
}

func GetPodService(config models.CaaS) *PodService{
	return &PodService{
		CaaS: config,
	}
}


func (service *PodService) GetPodStatus() ([]map[string]interface{}, error) {
	var resultList []map[string]interface{}

	podStatusBytes, err := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query=" + models.PROMQL_POD_PHASE, "")
	if err != nil {
		return nil, err
	}
	podStatusArray := gjson.Get(string(podStatusBytes), "data.result")
	for _, item := range podStatusArray.Array() {
		phase := item.Get("metric.phase").String()
		value := item.Get("value.1").Float()
		itemMap := make(map[string]interface{})
		itemMap["phase"] = phase
		itemMap["value"] = value
		resultList = append(resultList, itemMap)
	}

	return resultList, nil
}

func (service *PodService) GetPodList() ([]map[string]interface{}, error) {
	resultBytes, _ := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+models.PROMQL_POD_LIST, "")
	resultArray := gjson.Get(string(resultBytes), "data.result")

	var podList []map[string]interface{}
	for _, item := range resultArray.Array() {
		itemMap := make(map[string]interface{})
		pod := item.Get("metric.pod").String()
		namespace := item.Get("metric.namespace").String()
		itemMap["pod"] = pod
		itemMap["namespace"] = namespace
		podList = append(podList, itemMap)
	}

	resultBytes, _ = helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+models.PROMQL_POD_CPU_USE, "")
	cpuUseResultArray := gjson.Get(string(resultBytes), "data.result")
	var cpuUseList []map[string]interface{}
	for _, item := range cpuUseResultArray.Array() {
		itemMap := make(map[string]interface{})
		pod := item.Get("metric.pod").String()
		value := item.Get("value.1").Float()
		itemMap["pod"] = pod
		itemMap["value"] = value
		cpuUseList = append(cpuUseList, itemMap)
	}

	resultBytes, _ = helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+models.PROMQL_POD_CPU_USAGE, "")
	cpuUsageResultArray := gjson.Get(string(resultBytes), "data.result")
	var cpuUsageList []map[string]interface{}
	for _, item := range cpuUsageResultArray.Array() {
		itemMap := make(map[string]interface{})
		pod := item.Get("metric.pod").String()
		value := item.Get("value.1").Float()
		itemMap["pod"] = pod
		itemMap["value"] = value
		cpuUsageList = append(cpuUsageList, itemMap)
	}

	resultBytes, _ = helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+models.PROMQL_POD_MEMORY_USE, "")
	memoryUseResultArray := gjson.Get(string(resultBytes), "data.result")
	var memoryUseList []map[string]interface{}
	for _, item := range memoryUseResultArray.Array() {
		itemMap := make(map[string]interface{})
		pod := item.Get("metric.pod").String()
		value := item.Get("value.1").Float()
		itemMap["pod"] = pod
		itemMap["value"] = value
		memoryUseList = append(memoryUseList, itemMap)
	}

	resultBytes, _ = helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+models.PROMQL_POD_MEMORY_USAGE, "")
	memoryUsageResultArray := gjson.Get(string(resultBytes), "data.result")
	var memoryUsageList []map[string]interface{}
	for _, item := range memoryUsageResultArray.Array() {
		itemMap := make(map[string]interface{})
		pod := item.Get("metric.pod").String()
		value := item.Get("value.1").Float()
		itemMap["pod"] = pod
		itemMap["value"] = value
		memoryUsageList = append(memoryUsageList, itemMap)
	}

	resultBytes, _ = helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+models.PROMQL_POD_DISK_USE, "")
	diskUseResultArray := gjson.Get(string(resultBytes), "data.result")
	var diskUseList []map[string]interface{}
	for _, item := range diskUseResultArray.Array() {
		itemMap := make(map[string]interface{})
		pod := item.Get("metric.pod").String()
		value := item.Get("value.1").Float()
		itemMap["pod"] = pod
		itemMap["value"] = value
		diskUseList = append(diskUseList, itemMap)
	}

	resultBytes, _ = helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+models.PROMQL_POD_DISK_USAGE, "")
	diskUsageResultArray := gjson.Get(string(resultBytes), "data.result")
	var diskUsageList []map[string]interface{}
	for _, item := range diskUsageResultArray.Array() {
		itemMap := make(map[string]interface{})
		pod := item.Get("metric.pod").String()
		value := item.Get("value.1").Float()
		itemMap["pod"] = pod
		itemMap["value"] = value
		diskUsageList = append(diskUsageList, itemMap)
	}

	for i, item := range podList {
		pod := item["pod"].(string)

		cpuUseMap := cpuUseList[i]
		if strings.Compare(pod, cpuUseMap["pod"].(string)) == 0 {
			item["cpu"] = cpuUseMap["value"]
		}

		cpuUsageMap := cpuUsageList[i]
		if strings.Compare(pod, cpuUsageMap["pod"].(string)) == 0 {
			item["cpuUsage"] = cpuUsageMap["value"]
		}

		memoryUseMap := memoryUseList[i]
		if strings.Compare(pod, memoryUseMap["pod"].(string)) == 0 {
			item["memory"] = memoryUseMap["value"]
		}

		memoryUsageMap := memoryUsageList[i]
		if strings.Compare(pod, memoryUsageMap["pod"].(string)) == 0 {
			item["memoryUsage"] = memoryUsageMap["value"]
		}

		diskUseMap := diskUseList[i]
		if strings.Compare(pod, diskUseMap["pod"].(string)) == 0 {
			item["disk"] = diskUseMap["value"]
		}

		diskUsageMap := diskUsageList[i]
		if strings.Compare(pod, diskUsageMap["pod"].(string)) == 0 {
			item["diskUsage"] = diskUsageMap["value"]
		}
	}

	return podList, nil
}


func (service *PodService) GetPodDetailMetrics(pod string) ([]map[string]interface{}, error) {
	var resultList []map[string]interface{}
	fromToTimeParmameter := GetPromqlFromToParameter(3600, "600")

	promqlCpuUsage := service.CaaS.MakePromQLScriptForWorkloadMetrics("cpu", "", pod, fromToTimeParmameter)
	promqlMemoryUsage := service.CaaS.MakePromQLScriptForWorkloadMetrics("memory", "", pod, fromToTimeParmameter)
	promqlDiskUsage := service.CaaS.MakePromQLScriptForWorkloadMetrics("disk", "", pod, fromToTimeParmameter)

	helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query_range", "query="+promqlCpuUsage, "")
	helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query_range", "query="+promqlMemoryUsage, "")
	helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query_range", "query="+promqlDiskUsage, "")

	// CPU
	metricsBytes, _ := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query_range", "query="+promqlCpuUsage, "") // Retrieve workload's metric data
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

	// Memory
	metricsBytes, _ = helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query_range", "query="+promqlMemoryUsage, "") // Retrieve workload's metric data
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

	// Disk
	metricsBytes, _ = helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query_range", "query="+promqlDiskUsage, "") // Retrieve workload's metric data
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
package cp

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"paasta-monitoring-api/helpers"
	models "paasta-monitoring-api/models/api/v1"
	"strings"
)

type PodService struct {
	CaaS models.CP
}

func GetPodService(config models.CP) *PodService {
	return &PodService{
		CaaS: config,
	}
}

func (service *PodService) GetPodStatus(ctx echo.Context) ([]map[string]interface{}, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	var resultList []map[string]interface{}
	podStatusBytes, err := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+models.PROMQL_POD_PHASE, "")
	if err != nil {
		logger.Error(err)
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

func (service *PodService) GetPodContainerList(pod string) ([]map[string]interface{}, error) {

	// Retrieve container list in workload
	promqlContainerList := "count(kube_pod_container_info{pod='" + pod + "'})by(namespace,pod,container)"
	containersByte, _ := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+promqlContainerList, "") // Retrieve container list in workload
	containerArray := gjson.Get(string(containersByte), "data.result")

	var containerList []map[string]interface{}
	for _, container := range containerArray.Array() {
		containerMap := make(map[string]interface{})
		namespace := container.Get("metric.namespace").String()
		podName := container.Get("metric.pod").String()
		containerName := container.Get("metric.container").String()
		containerMap["namespace"] = namespace
		containerMap["pod"] = podName
		containerMap["container"] = containerName
		containerList = append(containerList, containerMap)
	}

	fmt.Printf("container list : %v\n", containerList)

	var cpuUseList []map[string]interface{}
	promqlCpuUse := "sum(container_cpu_usage_seconds_total{container!='POD',image!=''})by(namespace,pod,container)"
	cpuUseByte, _ := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+promqlCpuUse, "") // Retrieve container list in workload
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
	cpuUsageByte, _ := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+promqlCpuUsage, "") // Retrieve container list in workload
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
	memoryUseByte, _ := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+promqlMemoryUse, "") // Retrieve container list in workload
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
	memoryUsageByte, _ := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+promqlMemoryUsage, "") // Retrieve container list in workload
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
	diskUseByte, _ := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+promqlDiskUse, "") // Retrieve container list in workload
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

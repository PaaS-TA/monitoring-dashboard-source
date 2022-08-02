package cp

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"paasta-monitoring-api/helpers"
	models "paasta-monitoring-api/models/api/v1"
	"strconv"
	"strings"
	"time"
)

type ClusterService struct {
	CaaS models.CP
}

func GetClusterService(config models.CP) *ClusterService {
	return &ClusterService{
		CaaS: config,
	}
}

func (service *ClusterService) GetClusterAverage(ctx echo.Context) (map[string]interface{}, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	result := make(map[string]interface{})
	typeParam := ctx.Param("type")
	switch typeParam {
	case "pod":
		podUsageBytes, err := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+models.PROMQL_POD_USAGE, "")
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		podUsage := gjson.Get(string(podUsageBytes), "data.result.0.value.1")
		podUsageFloat, _ := strconv.ParseFloat(podUsage.String(), 64)
		result["pod"] = podUsageFloat
		break

	case "cpu":
		cpuUsageBytes, err := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+models.PROMQL_CPU_USAGE, "")
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		cpuUsage := gjson.Get(string(cpuUsageBytes), "data.result.0.value.1")
		cpuUsageFloat, _ := strconv.ParseFloat(cpuUsage.String(), 64)
		result["cpu"] = cpuUsageFloat
		break

	case "disk":
		diskUsageBytes, err := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+models.PROMQL_DISK_USAGE, "")
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		diskUsage := gjson.Get(string(diskUsageBytes), "data.result.0.value.1")
		diskUsageFloat, _ := strconv.ParseFloat(diskUsage.String(), 64)
		result["disk"] = diskUsageFloat
		break

	case "memory":
		memoryUsageBytes, err := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+models.PROMQL_MEMORY_USAGE, "")
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		memoryUsage := gjson.Get(string(memoryUsageBytes), "data.result.0.value.1")
		memoryUsageFloat, _ := strconv.ParseFloat(memoryUsage.String(), 64)
		result["memory"] = memoryUsageFloat
		break

	}
	return result, nil
}

func (service *ClusterService) GetWorkNodeList(ctx echo.Context) ([]map[string]interface{}, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	var result []map[string]interface{}
	nameListBytes, err := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+models.PROMQL_WORKNODE_NAME_LIST, "")
	memoryListBytes, err := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+models.PROMQL_WORKNODE_MEMORY_USAGE, "")
	memoryUseListBytes, err := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+models.PROMQL_WORKNODE_MEMORY_USE, "")
	cpuUsageListBytes, err := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+models.PROMQL_WORKNODE_CPU_USAGE, "")
	cpuUseListBytes, err := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+models.PROMQL_WORKNODE_CPU_ALLOC, "")
	diskUseListBytes, err := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+models.PROMQL_WORKNODE_DISK_USE, "")
	conditionListBytes, err := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query", "query="+models.PROMQL_WORKNODE_CONDITION, "")
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	nameListResult := gjson.Get(string(nameListBytes), "data.result")
	memoryListResult := gjson.Get(string(memoryListBytes), "data.result")
	memoryUseListResult := gjson.Get(string(memoryUseListBytes), "data.result")
	cpuUsageListResult := gjson.Get(string(cpuUsageListBytes), "data.result")
	cpuUseListResult := gjson.Get(string(cpuUseListBytes), "data.result")
	diskUseListResult := gjson.Get(string(diskUseListBytes), "data.result")
	conditionListResult := gjson.Get(string(conditionListBytes), "data.result")

	for i, info := range nameListResult.Array() {
		instance := info.Get("metric.instance").String()
		namespace := info.Get("metric.namespace").String()
		nodename := info.Get("metric.nodename").String()

		resultMap := make(map[string]interface{})

		resultMap["instance"] = instance
		resultMap["namespace"] = namespace
		resultMap["nodename"] = nodename

		memoryMap := memoryListResult.Array()[i]
		if strings.Compare(instance, memoryMap.Get("metric.instance").String()) == 0 {
			usageFloat, _ := strconv.ParseFloat(memoryMap.Get("value.1").String(), 64)
			resultMap["memory_usage"] = usageFloat
		}

		memoryListMap := memoryUseListResult.Array()[i]
		if strings.Compare(instance, memoryListMap.Get("metric.instance").String()) == 0 {
			useCapacityInt, _ := strconv.ParseUint(memoryListMap.Get("value.1").String(), 10, 64)
			resultMap["memory"] = useCapacityInt
		}

		cpuUsageListMap := cpuUsageListResult.Array()[i]
		if strings.Compare(instance, cpuUsageListMap.Get("metric.instance").String()) == 0 {
			usageFloat, _ := strconv.ParseFloat(cpuUsageListMap.Get("value.1").String(), 64)
			resultMap["cpu_usage"] = usageFloat
		}

		cpuUseListMap := cpuUseListResult.Array()[i]
		if strings.Compare(instance, cpuUseListMap.Get("metric.instance").String()) == 0 {
			resultMap["cpu"] = cpuUseListMap.Get("value.1").String()
		}

		diskUseListMap := diskUseListResult.Array()[i]
		if strings.Compare(instance, diskUseListMap.Get("metric.instance").String()) == 0 {
			useCapacityInt, _ := strconv.ParseUint(diskUseListMap.Get("value.1").String(), 10, 64)
			resultMap["disk"] = useCapacityInt
		}

		conditionListMap := conditionListResult.Array()[i]
		fmt.Println(conditionListMap)
		if strings.Compare(nodename, conditionListMap.Get("metric.node").String()) == 0 {
			resultMap["ready"] = true
		} else {
			resultMap["ready"] = false
		}
		result = append(result, resultMap)
	}
	return result, nil
}

func (service *ClusterService) GetWorkNode(ctx echo.Context) ([]map[string]interface{}, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	var result []map[string]interface{}

	nodeName := ctx.QueryParam("nodename")
	instanceName := ctx.QueryParam("instance")
	fromToTimeParmameter := GetPromqlFromToParameter(3600, "600")

	// 1.podUsage (input:nodeName)
	pqPodUsage := "sum(kube_pod_info{node='" + nodeName + "'})" + fromToTimeParmameter
	// 2.cpuUsage (input:Instance)
	pqCpuUsage := "node_cpu_seconds_total{mode!='idle',job='node-exporter',instance='" + instanceName + "'}" + fromToTimeParmameter
	// 3.memoryUsage (input:Instance)
	pqMemoryUsage := "node_memory_Active_bytes{job='node-exporter',instance='" + instanceName + "'}" + fromToTimeParmameter
	// 4.diskUsage (input:nodeName)
	pqDiskUsage := "sum(node_filesystem_size_bytes{instance='" + instanceName + "'}-node_filesystem_free_bytes{instance='" + instanceName + "'})" + fromToTimeParmameter

	podSeriesBytes, err := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query_range", "query="+pqPodUsage, "")
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	cpuSeriesBytes, err := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query_range", "query="+pqCpuUsage, "")
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	memorySeriesBytes, err := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query_range", "query="+pqMemoryUsage, "")
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	diskSeriesBytes, err := helpers.RequestHttpGet(service.CaaS.PromethusUrl+"/query_range", "query="+pqDiskUsage, "")
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	podSeriesResult := gjson.Get(string(podSeriesBytes), "data.result.0.values")
	cpuSeriesResult := gjson.Get(string(cpuSeriesBytes), "data.result.0.values")
	memorySeriesResult := gjson.Get(string(memorySeriesBytes), "data.result.0.values")
	diskSeriesResult := gjson.Get(string(diskSeriesBytes), "data.result.0.values")

	podSeriesMap := make(map[string]interface{})
	var podSeriesDataArr []map[string]interface{}
	for _, item := range podSeriesResult.Array() {
		itemMap := make(map[string]interface{})
		timestamp := item.Get("0").String()
		usage := item.Get("1").String()
		itemMap["time"] = timestamp
		itemMap["usage"] = usage
		podSeriesDataArr = append(podSeriesDataArr, itemMap)
	}
	podSeriesMap["name"] = "pod"
	podSeriesMap["metric"] = podSeriesDataArr

	cpuSeriesMap := make(map[string]interface{})
	var cpuSeriesDataArr []map[string]interface{}
	for _, item := range cpuSeriesResult.Array() {
		itemMap := make(map[string]interface{})
		timestamp := item.Get("0").String()
		usage := item.Get("1").String()
		itemMap["time"] = timestamp
		itemMap["usage"] = usage
		cpuSeriesDataArr = append(cpuSeriesDataArr, itemMap)
	}
	cpuSeriesMap["name"] = "cpu"
	cpuSeriesMap["metric"] = cpuSeriesDataArr

	memorySeriesMap := make(map[string]interface{})
	var memorySeriesDataArr []map[string]interface{}
	for _, item := range memorySeriesResult.Array() {
		itemMap := make(map[string]interface{})
		timestamp := item.Get("0").String()
		usage := item.Get("1").String()
		itemMap["time"] = timestamp
		itemMap["usage"] = usage
		memorySeriesDataArr = append(memorySeriesDataArr, itemMap)
	}
	memorySeriesMap["name"] = "memory"
	memorySeriesMap["metric"] = memorySeriesDataArr

	diskSeriesMap := make(map[string]interface{})
	var diskSeriesDataArr []map[string]interface{}
	for _, item := range diskSeriesResult.Array() {
		itemMap := make(map[string]interface{})
		timestamp := item.Get("0").String()
		usage := item.Get("1").String()
		itemMap["time"] = timestamp
		itemMap["usage"] = usage
		diskSeriesDataArr = append(diskSeriesDataArr, itemMap)
	}
	diskSeriesMap["name"] = "disk"
	diskSeriesMap["metric"] = diskSeriesDataArr

	result = append(result, podSeriesMap)
	result = append(result, cpuSeriesMap)
	result = append(result, memorySeriesMap)
	result = append(result, diskSeriesMap)

	return result, nil
}

func GetPromqlFromToParameter(interval int64, timeStep string) string {
	currentTime := time.Now().Unix()
	previousTime := currentTime - interval
	start := strconv.FormatInt(previousTime, 10)
	end := strconv.FormatInt(currentTime, 10)

	parameter := "&start=" + start + "&end=" + end + "&step=" + timeStep
	return parameter
}

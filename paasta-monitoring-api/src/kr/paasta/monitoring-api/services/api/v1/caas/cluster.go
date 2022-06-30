package caas

import (
	"github.com/tidwall/gjson"
	"fmt"
	"paasta-monitoring-api/helpers"
	models "paasta-monitoring-api/models/api/v1"
	"strconv"
	"strings"
)

type ClusterService struct {
	CaasConfig models.CaasConfig
}

func GetClusterService(config models.CaasConfig) *ClusterService{
	return &ClusterService{
		CaasConfig: config,
	}
}

func (service *ClusterService) GetClusterAverage(typeParam string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	switch typeParam {
	case "pod":
		podUsageBytes, err := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl, "query="+models.PROMQL_POD_USAGE)
		if err != nil {
			return nil, err
		}
		podUsage := gjson.Get(string(podUsageBytes), "data.result.0.value.1")
		podUsageFloat, _ := strconv.ParseFloat(podUsage.String(), 64)
		result["pod"] = podUsageFloat
		break;

	case "cpu":
		cpuUsageBytes, err := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl, "query="+models.PROMQL_CPU_USAGE)
		if err != nil {
			return nil, err
		}
		cpuUsage := gjson.Get(string(cpuUsageBytes), "data.result.0.value.1")
		cpuUsageFloat, _ := strconv.ParseFloat(cpuUsage.String(), 64)
		result["cpu"] = cpuUsageFloat
		break;

	case "disk" :
		diskUsageBytes, err := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl, "query="+models.PROMQL_DISK_USAGE)
		if err != nil {
			return nil, err
		}
		diskUsage := gjson.Get(string(diskUsageBytes), "data.result.0.value.1")
		diskUsageFloat, _ := strconv.ParseFloat(diskUsage.String(), 64)
		result["disk"] = diskUsageFloat
		break;

	case "memory" :
		memoryUsageBytes, err := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl, "query="+models.PROMQL_MEMORY_USAGE)
		if err != nil {
			return nil, err
		}
		memoryUsage := gjson.Get(string(memoryUsageBytes), "data.result.0.value.1")
		memoryUsageFloat, _ := strconv.ParseFloat(memoryUsage.String(), 64)
		result["memory"] = memoryUsageFloat
		break;

	}
	return result, nil
}


func (service *ClusterService) GetWorkNodeList() ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	nameListBytes, err := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl, "query="+models.PROMQL_WORKNODE_NAME_LIST)
	memoryListBytes, err := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl, "query="+models.PROMQL_WORKNODE_MEMORY_USAGE)
	memoryUseListBytes, err := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl, "query="+models.PROMQL_WORKNODE_MEMORY_USE)
	cpuUsageListBytes, err := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl, "query="+models.PROMQL_WORKNODE_CPU_USAGE)
	cpuUseListBytes, err := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl, "query="+models.PROMQL_WORKNODE_CPU_ALLOC)
	diskUseListBytes, err := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl, "query="+models.PROMQL_WORKNODE_DISK_USE)
	conditionListBytes, err := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl, "query="+models.PROMQL_WORKNODE_CONDITION)
	if err != nil {
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
			useCapacityInt, _ := strconv.ParseUint(memoryListMap.Get("value.1").String(), 10,64)
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
			useCapacityInt, _ := strconv.ParseUint(diskUseListMap.Get("value.1").String(), 10,64)
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
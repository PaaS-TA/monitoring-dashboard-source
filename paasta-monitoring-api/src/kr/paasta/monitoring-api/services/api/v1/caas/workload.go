package caas

import (
	"github.com/tidwall/gjson"
	"paasta-monitoring-api/helpers"
	models "paasta-monitoring-api/models/api/v1"
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


	daemonsetActiveResult, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query=" + models.PROMQL_WORKLOAD_DAEMONSET_AVAILABLE)
	daemonsetActive := gjson.Get(string(daemonsetActiveResult), "data.result.0.value.1").String()
	daemonsetInactiveResult, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query=" + models.PROMQL_WORKLOAD_DAEMONSET_UNAVAILABLE)
	daemonsetInactive := gjson.Get(string(daemonsetInactiveResult), "data.result.0.value.1").String()
	daemonsetReadyResult, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query=" + models.PROMQL_WORKLOAD_DAEMONSET_READY)
	daemonsetReady := gjson.Get(string(daemonsetReadyResult), "data.result.0.value.1").String()
	daemonsetMisscheduleResult, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query=" + models.PROMQL_WORKLOAD_DAEMONSET_MISSCHEDULED)
	daemonsetMisschedule := gjson.Get(string(daemonsetMisscheduleResult), "data.result.0.value.1").String()
	daemonsetMap := make(map[string]interface{})
	daemonsetMap["name"] = "DaemonSet"
	daemonsetMap["ready"] = daemonsetReady
	daemonsetMap["available"] = daemonsetActive
	daemonsetMap["unavailable"] = daemonsetInactive
	daemonsetMap["misscheduled"] = daemonsetMisschedule

	statefulsetTotalResult, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query=" + models.PROMQL_WORKLOAD_STATEFULSET_TOTAL)
	statefulsetTotal := gjson.Get(string(statefulsetTotalResult), "data.result.0.value.1").String()
	statefulsetReadyResult, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query=" + models.PROMQL_WORKLOAD_STATEFULSET_READY)
	statefulsetReady := gjson.Get(string(statefulsetReadyResult), "data.result.0.value.1").String()
	statefulsetRevisionResult, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query=" + models.PROMQL_WORKLOAD_STATEFULSET_REVISION)
	statefulsetRevision := gjson.Get(string(statefulsetRevisionResult), "data.result.0.value.1").String()
	statefulsetUpdatedResult, _ := helpers.RequestHttpGet(service.CaasConfig.PromethusUrl+"/query", "query=" + models.PROMQL_WORKLOAD_STATEFULSET_UPDATED)
	statefulsetUpdated := gjson.Get(string(statefulsetUpdatedResult), "data.result.0.value.1").String()
	statefulsetMap := make(map[string]interface{})
	statefulsetMap["name"] = "Stateful"
	statefulsetMap["total"] = statefulsetTotal
	statefulsetMap["ready"] = statefulsetReady
	statefulsetMap["revision"] = statefulsetRevision
	statefulsetMap["updated"] = statefulsetUpdated

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
	result = append(result, daemonsetMap)
	result = append(result, statefulsetMap)
	result = append(result, podcontainerMap)

	return result, nil

}

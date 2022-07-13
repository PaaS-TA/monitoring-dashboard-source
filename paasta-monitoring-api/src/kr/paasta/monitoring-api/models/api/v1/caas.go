package v1


const (
	PROMQL_POD_USAGE = "sum(kube_pod_info)/sum(sum(kube_node_status_allocatable{resource='pods'})by(node))*100"
	PROMQL_CPU_USAGE = "avg(instance:node_cpu:ratio)*100"
	PROMQL_DISK_USAGE = "(sum(node_filesystem_size_bytes)-sum(node_filesystem_free_bytes))/sum(node_filesystem_size_bytes)*100"
	PROMQL_MEMORY_USAGE = "(sum(node_memory_MemTotal_bytes)-sum(node_memory_MemFree_bytes+node_memory_Buffers_bytes+node_memory_Cached_bytes))/sum(node_memory_MemTotal_bytes)*100"

	PROMQL_WORKNODE_NAME_LIST    = "count(node_uname_info)by(instance,nodename,namespace)"
	//PROMQL_WORKNODE_NAME_LIST2   = "kube_node_created"
	PROMQL_WORKNODE_CPU_USAGE    = "instance:node_cpu:ratio*100"
	PROMQL_WORKNODE_MEMORY_USAGE = "max(((node_memory_MemTotal_bytes{job='node-exporter'}-" +
		"node_memory_MemFree_bytes{job='node-exporter'}" +
		"-node_memory_Buffers_bytes{job='node-exporter'}" +
		"-node_memory_Cached_bytes{job='node-exporter'})" +
		"/node_memory_MemTotal_bytes{job='node-exporter'})*100)by(instance)"

	PROMQL_WORKNODE_CPU_ALLOC  = "sum(node_cpu_seconds_total{mode='iowait',cpu='0'})by(instance)"
	PROMQL_WORKNODE_MEMORY_USE = "node_memory_Active_bytes"
	PROMQL_WORKNODE_DISK_USE  = "sum(node_filesystem_size_bytes-node_filesystem_free_bytes)by(instance)"
	PROMQL_WORKNODE_CONDITION = "count(kube_node_status_condition{condition='Ready',status='true'})by(node)"


	PROMQL_WORKLOAD_DEPLOYMENT_TOTAL        = "sum(kube_deployment_status_replicas)"
	PROMQL_WORKLOAD_DEPLOYMENT_AVAILABLE    = "sum(kube_deployment_status_replicas_available)"
	PROMQL_WORKLOAD_DEPLOYMENT_UNAVAILABLE  = "sum(kube_deployment_status_replicas_unavailable)"
	PROMQL_WORKLOAD_DEPLOYMENT_UPDATED      = "sum(kube_deployment_status_replicas_updated)"

	PROMQL_WORKLOAD_STATEFULSET_TOTAL        = "sum(kube_statefulset_status_replicas)"
	PROMQL_WORKLOAD_STATEFULSET_AVAILABLE    = "sum(kube_statefulset_status_replicas_available)"
	PROMQL_WORKLOAD_STATEFULSET_UNAVAILABLE  = "sum(kube_statefulset_status_replicas_unavailable)"
	PROMQL_WORKLOAD_STATEFULSET_UPDATED      = "sum(kube_statefulset_status_replicas_updated)"

	PROMQL_WORKLOAD_DAEMONSET_READY        = "sum(kube_daemonset_status_number_ready)"
	PROMQL_WORKLOAD_DAEMONSET_AVAILABLE    = "sum(kube_daemonset_status_number_available)"
	PROMQL_WORKLOAD_DAEMONSET_UNAVAILABLE  = "sum(kube_daemonset_status_number_unavailable)"
	PROMQL_WORKLOAD_DAEMONSET_MISSCHEDULED = "sum(kube_daemonset_status_number_misscheduled)"

	PROMQL_WORKLOAD_PODCONTAINER_READY     = "sum(kube_pod_container_status_ready)"
	PROMQL_WORKLOAD_PODCONTAINER_RUNNING   = "sum(kube_pod_container_status_running)"
	PROMQL_WORKLOAD_PODCONTAINER_RESTARTS  = "sum(kube_pod_container_status_restarts_total)"
	PROMQL_WORKLOAD_PODCONTAINER_TERMINATE = "sum(kube_pod_container_status_terminated)"

	PROMQL_POD_PHASE = "count(kube_pod_status_phase>0)by(phase)"
	PROMQL_POD_LIST         = "sum(container_cpu_usage_seconds_total{pod!='',image!=''})by(pod,namespace)"
	PROMQL_POD_CPU_USE      = "sum(container_cpu_usage_seconds_total{pod!='',image!=''})by(pod)"
	PROMQL_POD_CPU_USAGE    = "sum(rate(container_cpu_usage_seconds_total{pod!='',image!=''}[5m]))by(pod,namespace)*100"
	PROMQL_POD_MEMORY_USE   = "sum(container_memory_working_set_bytes{pod!='',image!=''})by(pod,namespace)/1024/1024"
	PROMQL_POD_DISK_USE     = "sum(container_fs_usage_bytes{pod!='',image!=''})by(pod,namespace)/1024/1024"
	PROMQL_POD_DISK_USAGE   = "sum(container_fs_usage_bytes{pod!='',image!=''})by(pod,namespace)/max(container_fs_limit_bytes{pod!='',image!=''})by(pod,namespace)*100"
	PROMQL_POD_MEMORY_USAGE = "avg(container_memory_working_set_bytes{pod!='',image!=''})by(pod,namespace)/scalar(sum(machine_memory_bytes))*100*scalar(count(container_memory_usage_bytes{pod!='',image!=''}))"

)


var BULITIN_WORKLOAD = [...]string {"deployment", "daemonset", "statefulset"}

type CaaS struct {
	PromethusUrl string
	PromethusRangeUrl string
	K8sUrl string
	K8sAdminToken string
}


type (

	ClusterAverage struct {
		PodUsage    string `json:"PodUsage"`
		CpuUsage    string `json:"CpuUsage"`
		MemoryUsage string `json:"MemoryUsage"`
		DiskUsage   string `json:"DiskUsage"`
	}
)

func (c CaaS) MakePromQLScriptForWorkloadMetrics(metricType string, namespace string, pod string, timeCondition string) string {
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
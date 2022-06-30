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
)

type (
	CaasConfig struct {
		PromethusUrl string
		PromethusRangeUrl string
		K8sUrl string
		K8sAdminToken string
	}

	ClusterAverage struct {
		PodUsage    string `json:"PodUsage"`
		CpuUsage    string `json:"CpuUsage"`
		MemoryUsage string `json:"MemoryUsage"`
		DiskUsage   string `json:"DiskUsage"`
	}
)
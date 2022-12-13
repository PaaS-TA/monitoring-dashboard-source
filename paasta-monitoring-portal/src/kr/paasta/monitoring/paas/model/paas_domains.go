package model

import (
	"encoding/json"
	"github.com/cloudfoundry-community/gogobosh"
)

const (
	TOP_PROCESS_CNT          = 5
	MB                       = 1048576
	KB                       = 1024
	STATE_RUNNING            = "running"
	STATE_FAILED             = "failed"
	STATE_CRITICAL           = "critical"
	STATE_WARNING            = "warning"
	ORIGIN_TYPE_PAAS         = "pas"
	DISK_STATE_NORMAL        = "healthy"
	METRIC_NAME_TOTAL_MEMORY = "memoryStats.TotalMemory"
	//METRIC_NAME_MEMORY_USAGE				= "memoryStats.UsedPercent"
	METRIC_NAME_FREE_MEMORY              = "memoryStats.FreeMemory"
	METRIC_NAME_TOTAL_DISK_ROOT          = "diskStats./.Total"
	METRIC_NAME_TOTAL_DISK_VCAP          = "diskStats./var/vcap/data.Total"
	METRIC_NAME_CPU_CORE_PREFIX          = "cpuStats.core."
	METRIC_NAME_CPU_LOAD_AVG_01_MIN      = "cpuStats.LoadAvg1Stats"
	METRIC_NAME_CPU_LOAD_AVG_05_MIN      = "cpuStats.LoadAvg5Stats"
	METRIC_NAME_CPU_LOAD_AVG_15_MIN      = "cpuStats.LoadAvg15Stats"
	METRIC_NAME_DISK_ROOT_USAGE          = "diskStats./.Usage"
	METRIC_NAME_DISK_VCAP_USAGE          = "diskStats./var/vcap/data.Usage"
	METRIC_NAME_DISK_IO_ROOT_READ_BYTES  = "diskIOStats.\\/\\..*.readBytes"
	METRIC_NAME_DISK_IO_ROOT_WRITE_BYTES = "diskIOStats.\\/\\..*.writeBytes"
	METRIC_NAME_DISK_IO_VCAP_READ_BYTES  = "diskIOStats.\\/var\\/vcap\\/data\\..*.readBytes"
	METRIC_NAME_DISK_IO_VCAP_WRITE_BYTES = "diskIOStats.\\/var\\/vcap\\/data\\..*.writeBytes"
	METRIC_NETWORK_IO_BYTES_SENT         = "networkIOStats.eth0.bytesSent"
	METRIC_NETWORK_IO_BYTES_RECV         = "networkIOStats.eth0.bytesRecv"
	METRIC_NETWORK_IO_PACKET_SENT        = "networkIOStats.eth0.packetSent"
	METRIC_NETWORK_IO_PACKET_RECV        = "networkIOStats.eth0.packetRecv"
	METRIC_NETWORK_IO_DROP_IN            = "networkIOStats.eth0.dropIn"
	METRIC_NETWORK_IO_DROP_OUT           = "networkIOStats.eth0.dropOut"
	METRIC_NETWORK_IO_ERR_IN             = "networkIOStats.eth0.errIn"
	METRIC_NETWORK_IO_ERR_OUT            = "networkIOStats.eth0.errOut"
	BOSH_DEPLOYMENT_NAME_CF              = "cf"
	BOSH_NAME                            = "bosh"
	USAGE_NAME_CPU                       = "CPU"
	USAGE_NAME_MEMORY                    = "Memory"
	USAGE_NAME_DISK_ROOT                 = "Disk(/)"
	USAGE_NAME_DISK_VCAP                 = "Disk(Data)"
)

type (
	PaasRequest struct {
		PagingReq               PagingReq
		Origin                  string `json:"origin"`
		Addr                    string `json:"addr"`
		MetricName              string `json:"metricName"`
		DefaultTimeRange        string `json:"defaultTimeRange"`
		TimeRangeFrom           string `json:"timeRangeFrom"`
		TimeRangeTo             string `json:"timeRangeTo"`
		GroupBy                 string `json:"groupBy"`
		ServiceName             string `json:"serviceName"`
		Ip                      string `json:"ip"`
		Index                   string `json:"index"`
		Name                    string `json:"name"`
		Id                      string `json:"id"`
		Args                    interface{}
		IsLikeQuery             bool
		IsRespondKb             bool
		IsNonNegativeDerivative bool
		Status                  string
	}

	PaasResponse struct {
		ServiceName string      `json:"serviceName"`
		Ip          string      `json:"ip"`
		Status      string      `json:"status"`
		Core        int         `json:"core"`
		CpuUsage    float64     `json:"cpuUsage"`
		MemorySize  json.Number `json:"memorySize"`
		MemoryUsage float64     `json:"memoryUsage"`
		DiskSize    json.Number `json:"diskSize"`
		DiskStatus  string      `json:"diskStatus"`
		Person      []string    `json:"persons"`
	}

	PaasOverview struct {
		Running  string `json:"Running"`
		Failed   string `json:"Failed"`
		Critical string `json:"Critical"`
		Warning  string `json:"Warning"`
		Total    string `json:"Total"`
	}

	PaasOverviewStatus struct {
		Data []PaasVm `json:"data"`
	}

	PaasSummary struct {
		Data         []PaasVm     `json:"data"`
		TotalCount   int          `json:"totalCount"`
		PageItems    int          `json:"pageItems"`
		PaasOverview PaasOverview `json:"overview"`
	}

	PaasVm struct {
		Name           string  `json:"name"`
		Id             string  `json:"id"`
		Address        string  `json:"address"`
		State          string  `json:"state"`
		Core           string  `json:"core"`
		CpuUsage       float64 `json:"cpuUsage"`
		CpuState       string  `json:"cpuErrStat"`
		TotalMemory    int64   `json:"totalMemory"`
		MemoryUsage    float64 `json:"memoryUsage"`
		MemoryState    string  `json:"memErrStat"`
		TotalDisk      int64   `json:"totalDisk"`
		TotalDiskUsage float64 `json:"-"`
		TotalDiskState string  `json:"diskRootErrStat"`
		DataDisk       int64   `json:"dataDisk"`
		DataDiskUsage  float64 `json:"-"`
		DataDiskState  string  `json:"diskDataErrStat"`
		DiskState      string  `json:"diskStatus"`
	}
	/*
		PaasTopProcessList struct {
			Data		[]PaasProcessUsage		`json:"data"`
		}
	*/
	PaasProcessUsage struct {
		Index   int64  `json:"index"`
		Time    string `json:"time"`
		Pid     string `json:"pid"`
		Process string `json:"process"`
		Memory  int64  `json:"memory"`
	}

	ResultVm struct {
		Id     string `json:"id"`
		ZoneId string `json:"zoneId"`
		Name   string `json:"name"`
		Ip     string `json:"ip"`
		VmType string `json:"vmType"`
	}

	UsageByTime struct {
		Time  int64       `json:"time"`
		Usage json.Number `json:"usage"`
	}

	MetricInfo struct {
		Metric []UsageByTime `json:"metric"`
		Name   string        `json:"name"`
	}

	MetricArg struct {
		Name  string
		Alias string
	}

	MemoryMetricArg struct {
		NameMemoryTotal string
		NameMemoryFree  string
		Alias           string
	}

	BoshDeployment struct {
		Name   string
		Status string
		VMS    []gogobosh.VM
	}

	MonitVms struct {
		Name        string            `json:"name"` //bosh
		Deployments []MonitDeployment `json:"children"`
	}

	MonitDeployment struct {
		Name   string   `json:"name"` //paasta-container
		Status string   `json:"status"`
		VMS    []BoshVm `json:"children"`
	}

	BoshVm struct {
		Name         string        `json:"name"` //access_z1/0
		Status       string        `json:"status"`
		BoshVmUsages []BoshVmUsage `json:"children"`
	}

	BoshVmUsage struct {
		Name   string  `json:"name"` //cpu,memory,disk
		Status string  `json:"status"`
		Usages []Usage `json:"children"`
	}

	Usage struct {
		Usage string `json:"name"`
	}

	Diagram struct {
		Id string            `json:"id"`
		Name string          `json:"name"`
		Title string         `json:"title"`
		Children []Diagram   `json:"children"`
	}
)

// Default Parent State: "running"
func DetermineVmState(parentState string, childState string) string {

	if childState == STATE_FAILED {
		parentState = STATE_FAILED

	} else if childState == STATE_CRITICAL && parentState != STATE_FAILED {
		parentState = STATE_CRITICAL

	} else if childState == STATE_WARNING && parentState != STATE_FAILED && parentState != STATE_CRITICAL {
		parentState = STATE_WARNING
	}

	return parentState
}

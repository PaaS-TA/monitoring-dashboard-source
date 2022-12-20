package model

import "encoding/json"

type (
	MetricsRequest struct {
		Origin           string `json:"origin"`
		Addr             string `json:"addr"`
		MetricName       string `json:"metricName"`
		DefaultTimeRange string `json:"defaultTimeRange"`
		TimeRangeFrom    string `json:"timeRangeFrom"`
		TimeRangeTo      string `json:"timeRangeTo"`
		GroupBy          string `json:"groupBy"`
		ServiceName      string `json:"serviceName"`
		Ip               string `json:"ip"`
		Index            string `json:"index"`
	}

	MetricsResponse struct {
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

	ErrMessage map[string]interface{}

	DiskIOUsage struct {
		DiskIOUsage []DiskIOUsageList `json:"data"`
	}
	DiskIOUsageList struct {
		Name string       `json:"name"`
		Data []MultiUsage `json:"data"`
	}

	MultiUsage struct {
		Time  int64       `json:"time"`
		Usage json.Number `json:"totalUsage"`
	}

	NetworkIOUsage struct {
		Person         []string             `json:"person"`
		NetworkIOUsage []NetworkIOUsageList `json:"data"`
	}
	NetworkIOUsageList struct {
		Name string         `json:"name"`
		Unit string         `json:"unit"`
		Data []NetworkUsage `json:"data"`
	}

	NetworkUsage struct {
		Time  int64   `json:"time"`
		Usage float64 `json:"totalUsage"`
	}

	TopProcessList struct {
		ServiceName string            `json:"serviceName"`
		Data        []TopProcessUsage `json:"data"`
	}

	TopProcessUsage struct {
		Index   string      `json:"index"`
		Memory  json.Number `json:"memory"`
		Pid     json.Number `json:"pid"`
		Process string      `json:"process"`
		Time    int64       `json:"time"`
	}

	AppNetworkIOKbyte struct {
		Name string                   `json:"name"`
		Data []map[string]interface{} `json:"data"`
	}

	ApplicationResources struct {
		CpuUsage  json.Number       `json:"cpu_usage"`
		MemUsage  json.Number       `json:"mem_usage"`
		DiskUsage json.Number       `json:"disk_usage"`
		Data      []ApplicationInfo `json:"data"`
	}
	ApplicationInfo struct {
		Time  int64       `json:"time"`
		Id    string      `json:"id"`
		Index string      `json:"index"`
		Name  string      `json:"name"`
		Value json.Number `json:"value"`
	}
)

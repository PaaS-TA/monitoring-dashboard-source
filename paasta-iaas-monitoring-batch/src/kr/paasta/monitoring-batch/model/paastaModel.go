package model

type(
	VmReq struct {
		ServiceName string        	`json:"serviceName"`
		Ip          string        	`json:"ip"`
		MountPoint  string        	`json:"mountPoint"`
		Status      string        	`json:"status"`
		MetricDatabase string           `json:"database"`
		DefaultTimeRange string         `json:"defaultTimeRange"`
		MeasureTimeList	[]AlarmItemMeasureTime
	}

	PaasTaResponse struct{
		ServiceName    string        	`json:"serviceName"`
		Ip          string        	`json:"ip"`
		Status      string              `json:"status"`
		CpuUsage    float64             `json:"cpuUsage"`
		MemoryUsage float64             `json:"memoryUsage"`
		DiskUsage  float64             `json:"diskStatus"`
		DiskRootUsage  float64             `json:"diskRootStatus"`
	}
)
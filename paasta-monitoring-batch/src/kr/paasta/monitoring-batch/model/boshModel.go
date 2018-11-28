package model

import "github.com/cloudfoundry-community/gogobosh"

type(
	//Bosh Topology Struct
	BoshDeployments struct{
		Name 		string          `json:"name"`
		VMS 		[]gogobosh.VM   `json:"children"`
	}

	BoshReq struct {
		ServiceName string        	`json:"serviceName"`
		Ip          string        	`json:"ip"`
		MountPoint  string        	`json:"mountPoint"`
		Status      string        	`json:"status"`
		MetricDatabase string           `json:"database"`
		DefaultTimeRange string         `json:"defaultTimeRange"`
		MeasureTimeList	[]AlarmItemMeasureTime
	}

	BoshResponse struct{
		ServiceName    string        	`json:"serviceName"`
		Ip          string        	`json:"ip"`
		Status      string              `json:"status"`
		CpuUsage    float64             `json:"cpuUsage"`
		MemoryUsage float64             `json:"memoryUsage"`
		DiskUsage  float64             `json:"diskStatus"`
		DiskRootUsage  float64             `json:"diskRootStatus"`
	}
)



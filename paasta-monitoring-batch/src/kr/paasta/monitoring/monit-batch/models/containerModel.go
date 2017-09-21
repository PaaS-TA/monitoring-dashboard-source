package models

type(

	//Zone DetailView Model Parameters
	ZonesReq struct {
		ApplicationName   string        	`json:"applicationName"`
		ApplicationIndex  string                `json:"applicationIndex"`
		ContainerName     string                `json:"containerName"`
		Addr              string        	`json:"addr"`
		Name              string                `json:"name"`
		CellIp            string        	`json:"cellIp"`
		MetricDatabase string                   `json:"database"`
		DefaultTimeRange string                 `json:"defaultTimeRange"`
	}

	ZoneCellInfo struct {
		ZoneName string
		CellName string
		Ip       string
		Id       uint
	}

	ZoneInfo struct {
		ZoneName string
		AppInfo []AppInfo
	}

	AppInfo struct{
		AppName string
		ContainerInfo []string
	}

	CellTileView struct{
		CellName   string
		ContainerTileView   []ContainerTileView
	}

	ContainerTileView struct {
		AppName  string
		AppGuid  string
		AppIndex int
		ContainerName string
		AlarmType     string
		Status        string
		CpuUsage      float64
		MemoryUsage   float64
		DiskUsage     float64
	}

	AppInfoResponse struct{
		AppName     string
		Index       string
		Status      string
		CpuUsage    float64
		MemoryUsage float64
		DiskStatus  string
	}

	AutoScaleAction struct{
		AppName string
		AppGuid string
		CpuUsage float64
		MemoryUsage float64
		Action  string
		Cause   string
	}
)

package model

type(

	//Zone DetailView Model Parameters
	ZonesReq struct {
		ApplicationName   string
		ApplicationIndex  string
		ContainerName     string
		Addr              string
		Name              string
		CellIp            string
		MetricDatabase    string
		DefaultTimeRange  string
		CheckTIme         string
		MeasureTimeList	[]AlarmItemMeasureTime
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
		AppName     string      `json:"appName"`
		AppGuid     string    	`json:"appGuid"`
		CpuUsage    string     `json:"cpuUsage"`
		MemoryUsage string     `json:"memoryUsage"`
		Action      string      `json:"action"`
		Instance    string      `json:"instanceCnt"`
		Cause       string      `json:"cause"`
	}

	ApplicationContainerInfo struct {
		ApplicationId		string
		ApplicationIndex	string
		ApplicationName		string
		ContainerInterface	string
		ContainerId			string
		CellIp				string
		CellName			string
		CpuUsage			float64
		MemoryUsage			float64
	}

	ApplicationInfo struct {
		ApplicationId				string
		ApplicationContainerInfo	[]ApplicationContainerInfo
	}

	AutoScaleTarget struct {
		AppName			string
		AppGuid			string
		CpuUsage		string
		MemoryUsage		string
		Action			string
		InstanceCnt		string
		Cause 			string
	}
)

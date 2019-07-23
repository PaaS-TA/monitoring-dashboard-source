package model

type (
	MetricsRequest struct {
		Nodename      string `json:"nodename"`
		Instance      string `json:"instance"`
		PodName       string `json:"PodName"`
		NameSpace     string `json:"NameSpace"`
		ContainerName string `json:"ContainerName"`
	}

	ClusterAvg struct {
		PodUsage    string `json:"PodUsage"`
		CpuUsage    string `json:"CpuUsage"`
		MemoryUsage string `json:"MemoryUsage"`
		DiskUsage   string `json:"DiskUsage"`
	}

	WorkNode struct {
		WorkNode []WorkNodeList `json:"Data"`
	}

	WorkNodeList struct {
		NodeName    string `json:"NodeName"`
		NameSpace   string `json:"NameSpace"`
		Ready       string `json:"Ready"`
		Cpu         string `json:"Cpu"`
		Memory      string `json:"Memory"`
		Disk        string `json:"Disk"`
		CpuUsage    string `json:"CpuUsage"`
		MemoryUsage string `json:"MemoryUsage"`
		Instance    string `json:"Instance"`
	}

	WorkNodeInfo struct {
		PodUsage    string `json:"PodUsage"`
		CpuUsage    string `json:"CpuUsage"`
		MemoryUsage string `json:"MemoryUsage"`
		DiskUsage   string `json:"DiskUsage"`
	}

	ContainerMetric struct {
		ContainerMetric []ContainerMetricList `json:"Data"`
	}

	ContainerMetricList struct {
		ContainerName string `json:"ContainerName"`
		PodName       string `json:"PodName"`
		NameSpace     string `json:"NameSpace"`
		Cpu           string `json:"Cpu"`
		Memory        string `json:"Memory"`
		Disk          string `json:"Disk"`
		CpuUsage      string `json:"CpuUsage"`
		MemoryUsage   string `json:"MemoryUsage"`
	}

	ContainerInfo struct {
		CpuUsage    string `json:"CpuUsage"`
		MemoryUsage string `json:"MemoryUsage"`
		DiskUsage   string `json:"DiskUsage"`
	}

	K8sLog struct {
		Log string `json:"Log"`
	}
)

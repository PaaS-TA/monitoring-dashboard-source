package model

type (
	MetricsRequest struct {
		Nodename      string `json:"nodename"`
		Instance      string `json:"instance"`
		PodName       string `json:"PodName"`
		NameSpace     string `json:"NameSpace"`
		ContainerName string `json:"ContainerName"`
		WorkloadsName string `json:"WorkloadsName"`
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

	ClusterOverview struct {
		Alerts           string `json:"Alerts"`
		RunningPod       string `json:"RunningPod"`
		Runningcontainer string `json:"RunningContainer"`
		PodRestart       string `json:"PodRestart"`
		Nodes            string `json:"Nodes"`
	}

	WorkloadsStatus struct {
		Name         string `json:"Name"`
		Total        string `json:"Total"`
		Available    string `json:"Available"`
		Unavailable  string `json:"Unavailable"`
		Updated      string `json:"Updated"`
		Ready        string `json:"Ready"`
		Revision     string `json:"Revision"`
		Misscheduled string `json:"Misscheduled"`
		Restart      string `json:"Restart"`
		Terminated   string `json:"terminated"`
		Running      string `json:"Running"`
	}

	WorkloadsContiSummary struct {
		Name        string `json:"Name"`
		Cpu         string `json:"Cpu"`
		Memory      string `json:"Memory"`
		Disk        string `json:"Disk"`
		CpuUsage    string `json:"CpuUsage"`
		MemoryUsage string `json:"MemoryUsage"`
		DiskUsage   string `json:"DiskUsage"`
	}

	PodPhase struct {
		Total     string `json:"Total"`
		Failed    string `json:"Failed"`
		Pending   string `json:"Pending"`
		Running   string `json:"Running"`
		Succeeded string `json:"Succeeded"`
		Unknown   string `json:"Unknown"`
	}

	PodMetricList struct {
		PodName     string `json:"PodName"`
		NameSpace   string `json:"NameSpace"`
		Cpu         string `json:"Cpu"`
		Memory      string `json:"Memory"`
		Disk        string `json:"Disk"`
		CpuUsage    string `json:"CpuUsage"`
		MemoryUsage string `json:"MemoryUsage"`
		DiskUsage   string `json:"DiskUsage"`
	}

	GraphMetric struct {
		Metric []map[string]string `json:"metric"`
		Name   string              `json:"name"`
	}

	GraphMetricList struct {
		Time        string `json:"Time"`
		PodUsage    string `json:"PodUsage"`
		CpuUsage    string `json:"CpuUsage"`
		MemoryUsage string `json:"MemoryUsage"`
		DiskUsage   string `json:"DiskUsage"`
	}
)

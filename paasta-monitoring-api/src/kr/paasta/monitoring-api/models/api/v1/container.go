package v1

import client "github.com/influxdata/influxdb1-client/v2"

type (
	ZoneInfo struct {
		ZoneName string     `json:"zoneName"`
		CellCnt  uint       `json:"cellCnt"`
		CellInfo []CellInfo `json:"cellInfo,omitempty" gorm:"-"`
	}

	CellInfo struct {
		ZoneName     string    `json:"zoneName"`
		CellName     string    `json:"cellName"`
		Ip           string    `json:"cellIp"`
		Id           uint      `json:"cellId"`
		AppCnt       uint      `json:"appCnt"`
		ContainerCnt uint      `json:"containerCnt"`
		AppInfo      []AppInfo `json:"appInfo,omitempty" gorm:"-"`
	}

	AppInfo struct {
		CellName      string         `json:"cellName"`
		AppName       string         `json:"appName"`
		Uri           string         `json:"uri"`
		Buildpack     string         `json:"buildpack"`
		Status        string         `json:"status"`
		Instances     int            `json:"instances"`
		Memory        int            `json:"memory"`
		DiskQuota     int            `json:"diskQuota"`
		CfApi         string         `json:"cfApi"`
		CreatedAt     string         `json:"createdAt"`
		UpdatedAt     string         `json:"updatedAt"`
		ContainerInfo *ContainerInfo `json:"containerInfo,omitempty" gorm:"-"`
	}

	ContainerInfo struct {
		AppName   string      `json:"appName"`
		Container []Container `json:"container"`
	}

	Container struct {
		AppIndex    string `json:"appIndex"`
		ContainerId string `json:"containerId"`
	}

	Overview struct {
		ZoneInfo []ZoneInfo `json:"zoneInfo,omitempty"`
	}

	Status struct {
		Running  uint `json:"running"`
		Warning  uint `json:"warning"`
		Critical uint `json:"critical"`
		Failed   uint `json:"failed"`
	}

	StatusByResource struct {
		CpuStatus    string
		MemoryStatus string
		DiskStatus   string
		TotalStatus  string
	}

	InfluxDbName struct {
		BoshDatabase      string
		PaastaDatabase    string
		ContainerDatabase string
		LoggingDatabase   string
	}

	InfluxDbClient struct {
		HttpClient client.Client
		DbName     InfluxDbName
	}

	InfluxDbQueryRequest struct {
		Sql        string
		CellIp     string
		MetricName string
	}

	CellMetricData struct {
		CpuCore   []map[string]interface{}
		CpuUsage  map[string]interface{}
		MemTotal  map[string]interface{}
		MemFree   map[string]interface{}
		DiskTotal map[string]interface{}
		DiskUsage map[string]interface{}
	}

	CellMetricDataFloat64 struct {
		CpuCore   uint
		CpuUsage  float64
		MemTotal  float64
		MemFree   float64
		MemUsage  float64
		DiskTotal float64
		DiskUsage float64
	}
)

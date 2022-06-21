package v1

import client "github.com/influxdata/influxdb1-client/v2"

type (
	CellInfo struct {
		Name string `json:"name"`
		Ip   string `json:"ip"`
		Id   uint   `json:"id"`
	}

	ZoneInfo struct {
		Name string `json:"name"`
	}

	AppInfo struct {
		Name      string `json:"name"`
		Uri       string `json:"uri"`
		Buildpack string `json:"buildpack"`
		Status    string `json:"status"`
		Instances int    `json:"instances"`
		Memory    int    `json:"memory"`
		DiskQuota int    `json:"diskQuota"`
		CfApi     string `json:"cfApi"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
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
)

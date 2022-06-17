package v1

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
	}

	Databases struct {
		BoshDatabase      string
		PaastaDatabase    string
		ContainerDatabase string
		LoggingDatabase   string
	}
)

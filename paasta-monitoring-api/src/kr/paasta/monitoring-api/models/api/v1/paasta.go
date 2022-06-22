package v1

type (
	PaastaSummary struct {
		Name                string `json:"name"`
		Ip                  string `json:"ip"`
		UUID                string `json:"uuid"`
		SqlQuery            string
		Time                string
		MetricName          string
		PaastaSummaryMetric PaastaSummaryMetric
	}

	PaastaSummaryMetric struct {
		State           string  `json:"state"`
		Core            string  `json:"core"`
		CpuUsage        float64 `json:"cpuUsage"`
		TotalMemory     float64 `json:"totalMemory"`
		MemoryUsage     float64 `json:"memoryUsage"`
		TotalDisk       float64 `json:"totalDisk"`
		DataDisk        float64 `json:"dataDisk"`
		DiskStatus      string  `json:"diskStatus"`
		PaastaState     string  `json:"-"`
		CpuErrStat      string  `json:"cpuErrStat"`
		MemErrStat      string  `json:"memErrStat"`
		DiskRootErrStat string  `json:"diskRootErrStat"`
		DiskDataErrStat string  `json:"diskDataErrStat"`
	}

	PaastaOverview struct {
		Running  string `json:"running"`
		Failed   string `json:"failed"`
		Critical string `json:"critical"`
		Warning  string `json:"warning"`
		Total    string `json:"total"`
	}

	Paasta struct {
		Id       int    `json:"id"`
		ZoneId   int    `json:"zoneId"`
		Name     string `json:"name"`
		Ip       string `json:"ip"`
		VmType   string `json:"vmType"`
		RegDate  string `json:"regDate"`
		RegUser  string `json:"regUser"`
		ModiDate string `json:"modiDate"`
		ModiUser string `json:"modiUser"`
	}

	PaastaProcess struct {
		Index   int64  `json:"index"`
		Time    string `json:"time"`
		Pid     string `json:"pid"`
		Process string `json:"process"`
		Memory  int64  `json:"memory"`
		UUID    string `json:"uuid"`
	}

	PaastaChart struct {
		UUID             string `json:"uuid"`
		MetricName       string `json:"metricname"`
		SqlQuery         string `json:"sqlquery"`
		DefaultTimeRange string `json:"defaulttimerange"`
		TimeRangeFrom    string `json:"timerangefrom"`
		TimeRangeTo      string `json:"timerangeto"`
		GroupBy          string `json:"groupby"`
		IsConvertKb      bool   `json:"isconvertkb"`
		MetricData       map[string]interface{}
	}

	PaastaLog struct {
		UUID       string `json:"uuid"`
		LogType    string `json:"logType"`
		Keyword    string `json:"keyword"`
		TargetDate string `json:"targetDate"`
		Period     string `json:"period"`
		StartTime  string `json:"startTime"`
		EndTime    string `json:"endTime"`
		Messages   interface{}
	}
)

package v1

type (
	PaastaSummary struct {
		Name                string `json:"name"`
		UUID                string `json:"uuid"`
		Address             string `json:"address"`
		PaastaSummaryMetric PaastaSummaryMetric
	}

	PaastaRequest struct {
		//PagingReq               PagingReq
		Origin                  string `json:"origin"`
		Addr                    string `json:"addr"`
		MetricName              string `json:"metricName"`
		DefaultTimeRange        string `json:"defaultTimeRange"`
		TimeRangeFrom           string `json:"timeRangeFrom"`
		TimeRangeTo             string `json:"timeRangeTo"`
		GroupBy                 string `json:"groupBy"`
		ServiceName             string `json:"serviceName"`
		Ip                      string `json:"ip"`
		Index                   string `json:"index"`
		Name                    string `json:"name"`
		Id                      string `json:"id"`
		Args                    interface{}
		IsLikeQuery             bool
		IsRespondKb             bool
		IsNonNegativeDerivative bool
		Status                  string
	}

	PaastaSummaryMetric struct {
		State          string  `json:"state"`
		Core           string  `json:"core"`
		CpuUsage       float64 `json:"cpuUsage"`
		CpuState       string  `json:"cpuErrStat"`
		TotalMemory    int64   `json:"totalMemory"`
		MemoryUsage    float64 `json:"memoryUsage"`
		MemoryState    string  `json:"memErrStat"`
		TotalDisk      int64   `json:"totalDisk"`
		TotalDiskUsage float64 `json:"-"`
		TotalDiskState string  `json:"diskRootErrStat"`
		DataDisk       int64   `json:"dataDisk"`
		DataDiskUsage  float64 `json:"-"`
		DataDiskState  string  `json:"diskDataErrStat"`
		DiskState      string  `json:"diskStatus"`
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
		UUID                    string `json:"uuid"`
		MetricName              string `json:"metricname"`
		SqlQuery                string `json:"sqlquery"`
		DefaultTimeRange        string `json:"defaulttimerange"`
		TimeRangeFrom           string `json:"timerangefrom"`
		TimeRangeTo             string `json:"timerangeto"`
		GroupBy                 string `json:"groupby"`
		IsConvertKb             bool   `json:"isconvertkb"`
		MetricData              map[string]interface{}
		IsLikeQuery             bool
		IsRespondKb             bool
		IsNonNegativeDerivative bool
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

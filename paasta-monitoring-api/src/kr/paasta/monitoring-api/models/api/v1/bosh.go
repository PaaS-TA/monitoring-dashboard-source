package v1

import "github.com/cloudfoundry-community/gogobosh"

const (
	BOSH_STATE_FAIL    = "failed"
	BOSH_STATE_NORMAL  = "healthy"
	BOSH_STATE_RUNNING = "running"

	MTR_CPU_CORE            = "cpuStats.core"
	MTR_CPU_LOAD_1M         = "cpuStats.LoadAvg1Stats"
	MTR_CPU_LOAD_5M         = "cpuStats.LoadAvg5Stats"
	MTR_CPU_LOAD_15M        = "cpuStats.LoadAvg15Stats"
	MTR_MEM_TOTAL           = "memoryStats.TotalMemory"
	MTR_MEM_FREE            = "memoryStats.FreeMemory"
	MTR_MEM_USAGE           = "memoryStats.UsedPercent"
	MTR_DISK_TOTAL          = "diskStats./.Total"
	MTR_DISK_USED           = "diskStats./.Used"
	MTR_DISK_USAGE          = "diskStats./.Usage"
	MTR_DISK_DATA_TOTAL     = "diskStats./var/vcap/data.Total"
	MTR_DISK_DATA_USED      = "diskStats./var/vcap/data.Used"
	MTR_DISK_DATA_USAGE     = "diskStats./var/vcap/data.Usage"
	MTR_DISK_USAGE_STR      = "diskStats.%s.Usage"
	MTR_DISK_IO_READ_STR    = "diskIOStats.%s.readBytes"
	MTR_DISK_IO_WRITE_STR   = "diskIOStats.%s.writeBytes"
	MTR_NETWORK_BYTE_SENT   = "networkIOStats.%s.bytesSent"
	MTR_NETWORK_BYTE_RECV   = "networkIOStats.%s.bytesRecv"
	MTR_NETWORK_PACKET_SENT = "networkIOStats.%s.packetSent"
	MTR_NETWORK_PACKET_RECV = "networkIOStats.%s.packetRecv"
	MTR_NETWORK_DROP_IN     = "networkIOStats.%s.dropIn"
	MTR_NETWORK_DROP_OUT    = "networkIOStats.%s.dropOut"
	MTR_NETWORK_ERROR_IN    = "networkIOStats.%s.errIn"
	MTR_NETWORK_ERROR_OUT   = "networkIOStats.%s.errOut"

	IFX_MTR_PROC_NAME = "process_name"
	IFX_MTR_PROC_PID  = "proc_pid"
	IFX_MTR_MEM_USAGE = "mem_usage"
	IFX_MTR_TIME      = "time"

	RESP_DATA_CPU_NAME             = "cpu"
	RESP_DATA_LOAD_1M_NAME         = "1m"
	RESP_DATA_LOAD_5M_NAME         = "5m"
	RESP_DATA_LOAD_15M_NAME        = "15m"
	RESP_DATA_MEM_NAME             = "memory"
	RESP_DATA_NETWORK_IO_SENT_NAME = "sent"
	RESP_DATA_NETWORK_IO_RECV_NAME = "recv"
	RESP_DATA_NETWORK_IO_IN_NAME   = "in"
	RESP_DATA_NETWORK_IO_OUT_NAME  = "out"
)

type (
	BoshSummary struct {
		Name              string `json:"name"`
		Ip                string `json:"ip"`
		UUID              string `json:"uuid"`
		SqlQuery          string
		Time              string
		MetricName        string
		BoshSummaryMetric BoshSummaryMetric
	}

	BoshSummaryMetric struct {
		State           string  `json:"state"`
		Core            string  `json:"core"`
		CpuUsage        float64 `json:"cpuUsage"`
		TotalMemory     float64 `json:"totalMemory"`
		MemoryUsage     float64 `json:"memoryUsage"`
		TotalDisk       float64 `json:"totalDisk"`
		DataDisk        float64 `json:"dataDisk"`
		DiskStatus      string  `json:"diskStatus"`
		BoshState       string  `json:"-"`
		CpuErrStat      string  `json:"cpuErrStat"`
		MemErrStat      string  `json:"memErrStat"`
		DiskRootErrStat string  `json:"diskRootErrStat"`
		DiskDataErrStat string  `json:"diskDataErrStat"`
	}

	BoshOverview struct {
		Running  string `json:"running"`
		Failed   string `json:"failed"`
		Critical string `json:"critical"`
		Warning  string `json:"warning"`
		Total    string `json:"total"`
	}

	Bosh struct {
		UUID       string           `json:"uuid"`
		Name       string           `json:"name"`
		Ip         string           `json:"ip"`
		Deployname string           `json:"deployment"`
		Address    string           `json:"address"`
		Username   string           `json:"username"`
		Password   string           `json:"password"`
		Client     *gogobosh.Client `json:"client"`
	}

	BoshProcess struct {
		Index   string  `json:"index"`
		Time    string  `json:"time"`
		Pid     string  `json:"pid"`
		Process string  `json:"process"`
		Memory  float64 `json:"memory"`
		UUID    string  `json:"uuid"`
	}

	BoshChart struct {
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

	BoshLog struct {
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

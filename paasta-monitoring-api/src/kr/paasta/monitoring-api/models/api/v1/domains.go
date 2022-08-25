package v1

const (
	CSRF_TOKEN_NAME   = "X-XSRF-TOKEN"
	TEST_TOKEN_NAME   = "TestCase"
	TEST_TOKEN_VALUE  = "TestCase"
	USER_SESSION_NAME = "info"

	METRIC_NAME_CPU_USAGE    = "cpu"
	METRIC_NAME_CPU_LOAD_1M  = "1m"
	METRIC_NAME_CPU_LOAD_5M  = "5m"
	METRIC_NAME_CPU_LOAD_15M = "15m"
	METRIC_NAME_MEMORY_SWAP  = "swap"
	METRIC_NAME_MEMORY_USAGE = "memory"

	METRIC_NAME_NETWORK_ETH_IN  = "InEth"
	METRIC_NAME_NETWORK_VX_IN   = "InVxlan"
	METRIC_NAME_NETWORK_ETH_OUT = "OutEth"
	METRIC_NAME_NETWORK_VX_OUT  = "OutVxlan"

	METRIC_NAME_NETWORK_ETH_IN_ERROR  = "InEth"
	METRIC_NAME_NETWORK_VX_IN_ERROR   = "InVxlan"
	METRIC_NAME_NETWORK_ETH_OUT_ERROR = "OutEth"
	METRIC_NAME_NETWORK_VX_OUT_ERROR  = "OutVxlan"

	METRIC_NAME_NETWORK_ETH_IN_DROPPED_PACKET  = "InEth"
	METRIC_NAME_NETWORK_VX_IN_DROPPED_PACKET   = "InVxlan"
	METRIC_NAME_NETWORK_ETH_OUT_DROPPED_PACKET = "OutEth"
	METRIC_NAME_NETWORK_VX_OUT_DROPPED_PACKET  = "OutVxlan"

	METRIC_NAME_DISK_READ_KBYTE  = "read"
	METRIC_NAME_DISK_WRITE_KBYTE = "write"

	METRIC_NAME_NETWORK_IN  = "in"
	METRIC_NAME_NETWORK_OUT = "out"

	RESULT_CNT        = "totalCnt"
	RESULT_PROJECT_ID = "tenantId"
	RESULT_NAME       = "name"
	RESULT_DATA       = "data"
	RESULT_DATA_NAME  = "metric"

	VM_STATUS_NO        = "noStatus"
	VM_STATUS_RUNNING   = "running"
	VM_STATUS_IDLE      = "idle/blocked"
	VM_STATUS_PAUSED    = "paused"
	VM_STATUS_SHUTDOWN  = "shutDown"
	VM_STATUS_SHUTOFF   = "shutOff"
	VM_STATUS_CRASHED   = "crashed"
	VM_STATUS_POEWR_OFF = "powerOff"
)

const (
	TOP_PROCESS_CNT          = 5
	MB                       = 1048576
	KB                       = 1024
	STATE_RUNNING            = "running"
	STATE_FAILED             = "failed"
	STATE_CRITICAL           = "critical"
	STATE_WARNING            = "warning"
	ORIGIN_TYPE_PAAS         = "pas"
	DISK_STATE_NORMAL        = "healthy"
	METRIC_NAME_TOTAL_MEMORY = "memoryStats.TotalMemory"
	//METRIC_NAME_MEMORY_USAGE				= "memoryStats.UsedPercent"
	METRIC_NAME_FREE_MEMORY              = "memoryStats.FreeMemory"
	METRIC_NAME_TOTAL_DISK_ROOT          = "diskStats./.Total"
	METRIC_NAME_TOTAL_DISK_VCAP          = "diskStats./var/vcap/data.Total"
	METRIC_NAME_CPU_CORE_PREFIX          = "cpuStats.core."
	METRIC_NAME_CPU_LOAD_AVG_01_MIN      = "cpuStats.LoadAvg1Stats"
	METRIC_NAME_CPU_LOAD_AVG_05_MIN      = "cpuStats.LoadAvg5Stats"
	METRIC_NAME_CPU_LOAD_AVG_15_MIN      = "cpuStats.LoadAvg15Stats"
	METRIC_NAME_DISK_ROOT_USAGE          = "diskStats./.Usage"
	METRIC_NAME_DISK_VCAP_USAGE          = "diskStats./var/vcap/data.Usage"
	METRIC_NAME_DISK_IO_ROOT_READ_BYTES  = "diskIOStats.\\/\\..*.readBytes"
	METRIC_NAME_DISK_IO_ROOT_WRITE_BYTES = "diskIOStats.\\/\\..*.writeBytes"
	METRIC_NAME_DISK_IO_VCAP_READ_BYTES  = "diskIOStats.\\/var\\/vcap\\/data\\..*.readBytes"
	METRIC_NAME_DISK_IO_VCAP_WRITE_BYTES = "diskIOStats.\\/var\\/vcap\\/data\\..*.writeBytes"
	METRIC_NETWORK_IO_BYTES_SENT         = "networkIOStats.eth0.bytesSent"
	METRIC_NETWORK_IO_BYTES_RECV         = "networkIOStats.eth0.bytesRecv"
	METRIC_NETWORK_IO_PACKET_SENT        = "networkIOStats.eth0.packetSent"
	METRIC_NETWORK_IO_PACKET_RECV        = "networkIOStats.eth0.packetRecv"
	METRIC_NETWORK_IO_DROP_IN            = "networkIOStats.eth0.dropIn"
	METRIC_NETWORK_IO_DROP_OUT           = "networkIOStats.eth0.dropOut"
	METRIC_NETWORK_IO_ERR_IN             = "networkIOStats.eth0.errIn"
	METRIC_NETWORK_IO_ERR_OUT            = "networkIOStats.eth0.errOut"
	BOSH_DEPLOYMENT_NAME_CF              = "cf"
	BOSH_NAME                            = "bosh"
	USAGE_NAME_CPU                       = "CPU"
	USAGE_NAME_MEMORY                    = "Memory"
	USAGE_NAME_DISK_ROOT                 = "Disk(/)"
	USAGE_NAME_DISK_VCAP                 = "Disk(Data)"
)

const (
	DB_DATE_FORMAT       string = "2006-01-02T15:04:05+00:00"
	ALARM_LEVEL_FAIL     string = "fail"
	ALARM_LEVEL_CRITICAL string = "critical"
	ALARM_LEVEL_WARNING  string = "warning"

	ALARM_TYPE_CPU    string = "cpu"
	ALARM_TYPE_MEMORY string = "memory"
	ALARM_TYPE_DISK   string = "disk"

	ORIGIN_TYPE_BOSH      string = "bos"
	ORIGIN_TYPE_PAASTA    string = "pas"
	ORIGIN_TYPE_CONTAINER string = "con"

	RESULT_STAT_NAME = "stat"
)

type (
	LogInfo struct {
		Time    string `json:"time"`
		Message string `json:"message"`
	}
)

var GmtTimeGap int

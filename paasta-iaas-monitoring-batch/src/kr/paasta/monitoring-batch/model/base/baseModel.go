package base

import (
	"time"
)

const (

	ORIGIN_TYPE_IAAS string = "ias"
	ORIGIN_TYPE_BOSH string = ""
	ORIGIN_TYPE_CONTAINER string = ""
	ORIGIN_TYPE_PAASTA string = ""

	ALARM_LEVEL_CRITICAL string = "critical"
	ALARM_LEVEL_WARNING  string = "warning"
	ALARM_LEVEL_FAIL     string = "fail"
	ALARM_TYPE_CPU string = "cpu"
	ALARM_TYPE_DISK string = "disk"
	ALARM_TYPE_ROOTDISK string = "rootdisk"
	ALARM_TYPE_MEMORY string = "memory"


	INSERT_DATE_FORMAT string = "2006-01-02T15:04:05+00:00"
	DEFAULT_ORIGIN_ID uint = 9999
	BAT_USER string = "bat"
	VM_TYPE_CEL string = "cel"
	VM_TYPE     string = "vms"

	SNS_TYPE_TELEGRAM = "telegram"

	SCALE_IN = "I"
	SCALE_OUT = "O"
	SCALE_RESOURCE_CPU = "CPU"
	SCALE_RESOURCE_MEM = "MEM"
	SCALE_API_URI = "/external/app/updateApp​"
)



type (

	BaseModel struct {
		Origin     string        	`json:"origin"`
		JobName    string        	`json:"jobName"`
		MetricName string       	`json:"metricName"`
		DefaultTimeRange  string     	`json:"defaultTimeRange"`
		TimeRangeFrom     string        `json:"timeRangeFrom"`
		TimeRangeTo       string        `json:"timeRangeTo"`
		GroupBy           string        `json:"groupBy"`

	}

	Event struct{
		OriginType    string
		OriginId      uint
		AlarmType     string
		AlarmLevel    string
		Ip            string
		AppYn         string
		AppName       string
		AppIndex      int
		ContainerName string
		CellName      string
		Message       string
		CpuUsage      float64
		MemoryUsage   float64
		DiskUsage     float64
	}

	EventResponse struct{
		OriginType    string
		OriginId      uint
		AlarmType     string
		AlarmLevel    string
		AppYn         string
		AppName       string
		AppIndex      int
		ContainerName string
		CellName      string
		Message       string
		CpuUsage      float64
		MemoryUsage   float64
		DiskUsage     float64
		AlarmDefaultYn string
		DownTime      int
		Count         int
		ResolveTime   int
	}

	AlarmTypeResponse struct {
		OriginId    uint
		OriginType  string
		AlarmType   string
		Level       string
		Ip          string
		ServiceName string
		ResolveTime int
		DefaultYn   string
		RepeatTime  int
		Threshold   int
		DownTime    int
	}

	AlarmThreshold struct {
		OriginId    uint
		OriginType  string
		AlarmType   string
		Level       string
		Ip          string
		ServiceName string
		ResolveTime int
		DefaultYn   string
		RepeatTime  int
		Threshold   float64
		DownTime    int
		MeasureTime int
	}

	Alarm   struct{
		Id            uint
		OriginType    string
		OriginId      uint
		AlarmType     string
		Level         string
		Ip            string
		AlarmTitle    string
		AlarmMessage  string
		ResolveStatus string
		AppYn         string
		AppName       string
		AppIndex      int
		ContainerName string
		AlarmCnt      int
		AlarmSendDate time.Time  `gorm:"type:datetime;null;DEFAULT:null"`
		RegDate       time.Time  `gorm:"type:datetime;null;DEFAULT:null"`
		RegUser       string     `gorm:"type:varchar(36);null;"`
		ModiDate      time.Time  `gorm:"type:datetime;null;DEFAULT:null"`
		ModiUser      string     `gorm:"type:varchar(36);null;"`
	}

	AlarmData   struct{
		Id            uint
		OriginType    string
		OriginId      uint
		AlarmType     string
		Level         string
		Ip            string
		AlarmTitle    string
		AlarmMessage  string
		ResolveStatus string
		AppYn         string
		AppName       string
		AppIndex      int
		ContainerName string
		AlarmCnt      int
		AlarmSendDate time.Time
		RegDate       time.Time
		ModiDate      time.Time
		ModiUser      string
	}

	AlarmPolicy    struct{
		Id            uint
		OriginType    string
		OriginId      uint
		AlarmType     string
		WarningThreshold int
		CriticalThreshold int
		repeatTime int
	}

	ServiceFileSystemUsage struct {
		ServiceName     string
		FileSystemUsage []FileSystemUsage
	}
	//현재 FileSystem의 사용률을 받아온다.
	FileSystemUsage struct{
		FileSystemName    string        `json:"name"`
		Usage             float64       `json:"totalUsage"`
	}

)

type ErrMessage map[string]interface{}
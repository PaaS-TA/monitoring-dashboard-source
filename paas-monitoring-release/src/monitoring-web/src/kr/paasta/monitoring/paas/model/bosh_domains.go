package model

import (
	"errors"
)

const (
	BOSH_STATE_FAIL	= "failed"
	BOSH_STATE_NORMAL = "healthy"
	BOSH_STATE_RUNNING = "running"

	MTR_CPU_CORE = "cpuStats.core"
	MTR_CPU_LOAD_1M = "cpuStats.LoadAvg1Stats"
	MTR_CPU_LOAD_5M = "cpuStats.LoadAvg5Stats"
	MTR_CPU_LOAD_15M = "cpuStats.LoadAvg15Stats"
	MTR_MEM_TOTAL = "memoryStats.TotalMemory"
	MTR_MEM_FREE = "memoryStats.FreeMemory"
	MTR_MEM_USAGE = "memoryStats.UsedPercent"
	MTR_DISK_TOTAL = "diskStats./.Total"
	MTR_DISK_USED = "diskStats./.Used"
	MTR_DISK_USAGE = "diskStats./.Usage"
	MTR_DISK_DATA_TOTAL = "diskStats./var/vcap/data.Total"
	MTR_DISK_DATA_USED = "diskStats./var/vcap/data.Used"
	MTR_DISK_USAGE_STR ="diskStats.%s.Usage"
	MTR_DISK_IO_READ_STR = "diskIOStats.%s.readBytes"
	MTR_DISK_IO_WRITE_STR ="diskIOStats.%s.writeBytes"
	MTR_NETWORK_BYTE_SENT = "networkIOStats.%s.bytesSent"
	MTR_NETWORK_BYTE_RECV = "networkIOStats.%s.bytesRecv"
	MTR_NETWORK_PACKET_SENT = "networkIOStats.%s.packetSent"
	MTR_NETWORK_PACKET_RECV = "networkIOStats.%s.packetRecv"
	MTR_NETWORK_DROP_IN = "networkIOStats.%s.dropIn"
	MTR_NETWORK_DROP_OUT = "networkIOStats.%s.dropOut"
	MTR_NETWORK_ERROR_IN = "networkIOStats.%s.errIn"
	MTR_NETWORK_ERROR_OUT = "networkIOStats.%s.errOut"

	IFX_MTR_PROC_NAME = "process_name"
	IFX_MTR_PROC_PID = "proc_pid"
	IFX_MTR_MEM_USAGE = "mem_usage"
	IFX_MTR_TIME = "time"

	RESP_DATA_CPU_NAME = "cpu"
	RESP_DATA_LOAD_1M_NAME = "1m"
	RESP_DATA_LOAD_5M_NAME = "5m"
	RESP_DATA_LOAD_15M_NAME = "15m"
	RESP_DATA_MEM_NAME = "memory"
	RESP_DATA_NETWORK_IO_SENT_NAME = "sent"
	RESP_DATA_NETWORK_IO_RECV_NAME = "recv"
	RESP_DATA_NETWORK_IO_IN_NAME = "in"
	RESP_DATA_NETWORK_IO_OUT_NAME = "out"
)

type (

	BoshSummaryReq struct {
		PagingReq
		Name string      	`json:"name"`
		Ip string        	`json:"ip"`
		Id string			`json:"id"`
		SqlQuery string
		MetricName string
		Time string
	}

	BoshSummaryRes struct {
		Name 		string        		`json:"name"`
		Address		string        		`json:"address"`
		Id 			string        		`json:"id"`
		State		string              `json:"state"`
		Core		string              `json:"core"`
		CpuUsage	float64             `json:"cpuUsage"`
		TotalMemory	float64         	`json:"totalMemory"`
		MemoryUsage	float64             `json:"memoryUsage"`
		TotalDisk	float64        		`json:"totalDisk"`
		DataDisk	float64        		`json:"dataDisk"`
		DiskStatus	string              `json:"diskStatus"`
		BoshState	string				`json:"-"`
		CpuErrStat	string              `json:"cpuErrStat"`
		MemErrStat	string              `json:"memErrStat"`
		DiskRootErrStat	string          `json:"diskRootErrStat"`
		DiskDataErrStat	string          `json:"diskDataErrStat"`
	}

	BoshOverviewCntRes struct {
		Running		string        		`json:"running"`
		Failed		string        		`json:"failed"`
		Critical	string        		`json:"critical"`
		Warning		string              `json:"warning"`
		Total		string              `json:"total"`
	}

	BoshStatusOverviewRes struct {
		PagingRes
		Overview BoshOverviewCntRes		`json:"overview"`
		Data []BoshSummaryRes			`json:"data"`
	}

	BoshTopProcessUsage struct {
		Index       string 			`json:"index"`
		Time        string         	`json:"time"`
		Pid         string        	`json:"pid"`
		Process     string         	`json:"process"`
		Memory      float64        	`json:"memory"`
	}

	BoshTopprocessUsageRes struct {
		Data []BoshTopProcessUsage	`json:"data"`
	}

	BoshDetailReq struct {
		Id          	  string
		MetricName        string
		SqlQuery 		  string
		DefaultTimeRange  string
		TimeRangeFrom     string
		TimeRangeTo       string
		GroupBy           string
		IsConvertKb		  bool
	}

)

func (bm BoshDetailReq) MetricRequestValidate(req BoshDetailReq) error {

	if req.Id == "" {
		return errors.New("Required input value does not exist. [id]")
	}

	//조회 조건 Validation Check
	if req.TimeRangeFrom == "" && req.TimeRangeTo == ""{
		if req.DefaultTimeRange == ""{
			return errors.New("Required input value does not exist. [defaultTimeRange]")
		}else{
			if req.GroupBy != ""{
				return nil
			}else{
				return errors.New("Required input value does not exist. [groupBy]")
			}
		}
		return errors.New("Required input value does not exist. [timeRangeFrom, timeRangeTo]")
	}else{
		if req.TimeRangeFrom == "" || req.TimeRangeTo == ""{
			if req.TimeRangeFrom == ""{
				return errors.New("Required input value does not exist. [timeRangeFrom]")
			}else if req.TimeRangeTo == ""{
				return errors.New("Required input value does not exist. [timeRangeTo]")
			}
		}else{

			if req.GroupBy == ""{
				return errors.New("Required input value does not exist. [groupBy]")
			}else{
				return nil
			}
		}
	}
	return nil

}
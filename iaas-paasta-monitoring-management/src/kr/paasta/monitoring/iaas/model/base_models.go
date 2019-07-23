package model

import (
	"errors"
	"github.com/cihub/seelog"
	monascagopher "github.com/gophercloud/gophercloud"
	"github.com/rackspace/gophercloud"
	"time"
	"unicode"
)

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

type ErrMessage map[string]interface{}

var GmtTimeGap int
var TestUserName string
var TestPassword string
var TestTenantID string
var TestDomainName string
var TestIdentityEndpoint string

type (
	Cookie struct {
		Name       string
		Value      string
		Path       string
		Domain     string
		Expires    time.Time
		RawExpires string
		MaxAge     int
		Secure     bool
		HttpOnly   bool
		Raw        string
		Unparsed   []string // Raw text of unparsed attribute-value pairs
		// MaxAge=0 means no 'Max-Age' attribute specified.
		// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
		// MaxAge>0 means Max-Age attribute present and given in seconds
	}

	User struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Token    string `json:"token"`
	}

	UserSession struct {
		Username  string
		Password  string
		CsrfToken string
		//OpenstackToken    string
		MonAuth monascagopher.AuthOptions
	}

	DetailReq struct {
		HostName         string
		InstanceId       string
		MetricName       string
		MountPoint       string
		DefaultTimeRange string
		TimeRangeFrom    string
		TimeRangeTo      string
		GroupBy          string
	}

	NodeReq struct {
		HostName string
	}

	InstanceReq struct {
		InstanceId string
	}

	TenantReq struct {
		TenantId   string
		TenantName string
		Limit      string
		Marker     string
		HostName   string
	}

	AlarmReq struct {
		AlarmId   string
		State     string
		TimeRange string
	}

	ErrorMessageStruct struct {
		Message    string `json:"message"`
		HttpStatus int    `json:"HttpStatus"`
	}

	TopProcess struct {
		Index       int     `json:"index"`
		ProcessName string  `json:"processName"`
		Usage       float64 `json:"usage"`
	}

	LogMessage struct {
		Hostname     string    `json:"hostname"`
		PageIndex    int       `json:"pageIndex"`
		PageItems    int       `json:"pageItems"`
		LogType      string    `json:"logType"`
		Keyword      string    `json:"keyword"`
		Index        string    `json:"logstashIndex"`
		TargetDate   string    `json:"targetDate"`
		Period       int64     `json:"period"`
		StartTime    string    `json:"startTime"`
		EndTime      string    `json:"endTime"`
		CurrentItems int       `json:"currentItems"`
		TotalCount   int       `json:"totalCount"`
		Messages     []LogInfo `json:"messages"`
	}

	LogInfo struct {
		Time    string `json:"time"`
		Message string `json:"message"`
	}
)

var OpenStackClient []map[string]*gophercloud.ProviderClient

var MonitLogger seelog.LoggerInterface

func (bm TenantReq) TenantInstanceRequestValidate(req TenantReq) error {
	if req.Limit == "" {
		return errors.New("Required input value does not exist. [limit]")
	}
	if isInt(req.Limit) == false {
		return errors.New("Required input value require Number. [limit]")
	}
	if req.TenantId == "" {
		return errors.New("Required input value does not exist. [tenantId]")
	}
	return nil
}

func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

//var SessionManager scs.Manager // = scs.NewCookieManager("u46IpCV9y5Vlur8YvODJEhgOY8m9JVE4")

func (bm DetailReq) MetricRequestValidate(req DetailReq) error {

	if req.HostName == "" {
		return errors.New("Required input value does not exist. [hostname]")
	}

	//조회 조건 Validation Check
	if req.TimeRangeFrom == "" && req.TimeRangeTo == "" {
		if req.DefaultTimeRange == "" {
			return errors.New("Required input value does not exist. [defaultTimeRange]")
		} else {
			if req.GroupBy != "" {
				return nil
			} else {
				return errors.New("Required input value does not exist. [groupBy]")
			}
		}
		return errors.New("Required input value does not exist. [timeRangeFrom, timeRangeTo]")
	} else {
		if req.TimeRangeFrom == "" || req.TimeRangeTo == "" {
			if req.TimeRangeFrom == "" {
				return errors.New("Required input value does not exist. [timeRangeFrom]")
			} else if req.TimeRangeTo == "" {
				return errors.New("Required input value does not exist. [timeRangeTo]")
			}
		} else {

			if req.GroupBy == "" {
				return errors.New("Required input value does not exist. [groupBy]")
			} else {
				return nil
			}
		}
	}
	return nil

}

func (bm DetailReq) InstanceMetricRequestValidate(req DetailReq) error {

	if req.InstanceId == "" {
		return errors.New("Required input value does not exist. [instanceId]")
	}

	//조회 조건 Validation Check
	if req.TimeRangeFrom == "" && req.TimeRangeTo == "" {
		if req.DefaultTimeRange == "" {
			return errors.New("Required input value does not exist. [defaultTimeRange]")
		} else {
			if req.GroupBy != "" {
				return nil
			} else {
				return errors.New("Required input value does not exist. [groupBy]")
			}
		}
		return errors.New("Required input value does not exist. [timeRangeFrom, timeRangeTo]")
	} else {
		if req.TimeRangeFrom == "" || req.TimeRangeTo == "" {
			if req.TimeRangeFrom == "" {
				return errors.New("Required input value does not exist. [timeRangeFrom]")
			} else if req.TimeRangeTo == "" {
				return errors.New("Required input value does not exist. [timeRangeTo]")
			}
		} else {

			if req.GroupBy == "" {
				return errors.New("Required input value does not exist. [groupBy]")
			} else {
				return nil
			}
		}
	}
	return nil

}

func (bm LogMessage) DefaultLogValidate(req LogMessage) error {
	if req.Hostname == "" {
		return errors.New("Required input value does not exist. [hostname]")
	}
	if req.LogType == "" {
		return errors.New("Required input value does not exist. [logType]")
	}
	if req.PageIndex == 0 {
		return errors.New("Required input value does not exist. [pageIndex]")
	}
	if req.PageItems == 0 {
		return errors.New("Required input value does not exist. [pageItems]")
	}
	if req.Period == 0 {
		return errors.New("Required input value does not exist. [period]")
	}

	return nil
}

func (bm LogMessage) SpecificTimeRangeLogValidate(req LogMessage) error {
	if req.Hostname == "" {
		return errors.New("Required input value does not exist. [hostname]")
	}
	if req.LogType == "" {
		return errors.New("Required input value does not exist. [logType]")
	}
	if req.TargetDate == "" {
		return errors.New("Required input value does not exist. [targetDate]")
	}

	return nil
}

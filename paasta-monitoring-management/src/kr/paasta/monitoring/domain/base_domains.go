package domain

import (
	"time"
	"fmt"
)

const (
	DB_DATE_FORMAT string = "2006-01-02T15:04:05+00:00"
	ALARM_LEVEL_CRITICAL string = "critical"
	ALARM_LEVEL_WARNING  string = "warning"

	ALARM_TYPE_CPU string = "cpu"
	ALARM_TYPE_MEMORY  string = "memory"
	ALARM_TYPE_DISK  string = "disk"

	ORIGIN_TYPE_BOSH string = "bos"
	ORIGIN_TYPE_PAASTA string = "pas"
	ORIGIN_TYPE_CONTAINER string = "con"
)

type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

var GmtTimeGap int

type JSONTime time.Time

func (t JSONTime)MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05")) // 01/02 03:04:05PM '06 -0700 <<== GMT-0700
	return []byte(stamp), nil
}

type (
	BaseModel struct {
		Service     string        	`json:"service"`
		JobName    string        	`json:"jobName"`
		MetricName string       	`json:"metricName"`
		DefaultTimeRange  string     	`json:"defaultTimeRange"`
		TimeRangeFrom     string        `json:"timeRangeFrom"`
		TimeRangeTo       string        `json:"timeRangeTo"`
		GroupBy           string        `json:"groupBy"`
	}

	PagingReq struct {
		PageIndex int                   `json:"pageIndex"`
		PageItem  int                   `json:"pageItem"`
	}

	PagingRes struct {
		PageIndex int                   `json:"pageIndex"`
		PageItem  int                   `json:"pageItem"`
		TotalCount int                  `json:"totalCount"`
	}
)

type ErrMessage map[string]interface{}
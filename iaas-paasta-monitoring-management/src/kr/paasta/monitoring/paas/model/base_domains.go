package model

import (
	"errors"
	"fmt"
	"time"
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

	RESULT_NAME      = "name"
	RESULT_STAT_NAME = "stat"
)

type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

var GmtTimeGap int

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05")) // 01/02 03:04:05PM '06 -0700 <<== GMT-0700
	return []byte(stamp), nil
}

type (
	CFConfig struct {
		UserId         string
		UserPw         string
		Host           string
		Type           string
		CaasBrokerHost string
	}

	BaseModel struct {
		Service          string `json:"service"`
		JobName          string `json:"jobName"`
		MetricName       string `json:"metricName"`
		DefaultTimeRange string `json:"defaultTimeRange"`
		TimeRangeFrom    string `json:"timeRangeFrom"`
		TimeRangeTo      string `json:"timeRangeTo"`
		GroupBy          string `json:"groupBy"`
	}

	PagingReq struct {
		PageIndex int `json:"pageIndex"`
		PageItem  int `json:"pageItem"`
	}

	PagingRes struct {
		PageIndex  int `json:"pageIndex"`
		PageItem   int `json:"pageItem"`
		TotalCount int `json:"totalCount"`
	}

	LogMessage struct {
		Id           string    `json:"id"`
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

type ErrMessage map[string]interface{}

func (bm LogMessage) DefaultLogValidate(req LogMessage) error {
	if req.Id == "" {
		return errors.New("Required input value does not exist. [id]")
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
	if req.Id == "" {
		return errors.New("Required input value does not exist. [id]")
	}
	if req.LogType == "" {
		return errors.New("Required input value does not exist. [logType]")
	}
	if req.TargetDate == "" {
		return errors.New("Required input value does not exist. [targetDate]")
	}

	return nil
}

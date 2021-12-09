package model

import (
	"time"
)

const (
	ALARM_SEVERITY_CRITICAL  = "CRITICAL"
	ALARM_SEVERITY_HIGH      = "HIGH"
	ALARM_STATE_ALARM        = "ALARM"
	ALARM_STATE_OK           = "OK"
	ALARM_STATE_UNDETERMINED = "UNDETERMINED"
)

type (
	AlarmNotification struct {
		Id     string `json:"id"`
		Name   string `json:"name"`
		Period int    `json:"period"`
		Email  string `json:"email"`
	}

	AlarmDefinition struct {
		Id                  string   `json:"id"`
		Name                string   `json:"name"`
		Severity            string   `json:"severity"`
		Expression          string   `json:"expression"`
		AlarmAction         []string `json:"alarmAction"`
		OkAction            []string `json:"okAction"`
		MatchBy             []string `json:"matchBy"`
		UndeterminedActions []string `json:"undetermined_actions"`
		Description         string   `json:"description"`
		//totalCnt         	int      `json:"totalCnt"`
	}

	AlarmDefinitionDetail struct {
		Id                  string              `json:"id"`
		Name                string              `json:"name"`
		Severity            string              `json:"severity"`
		Expression          string              `json:"expression"`
		MatchBy             []string            `json:"matchBy"`
		UndeterminedActions []string            `json:"undetermined_actions"`
		Description         string              `json:"description"`
		AlarmNotification   []AlarmNotification `json:"alarmAction"`
		//OkAction   	 	[]string `json:"okAction"`

	}

	AlarmAction struct {
		Id     string `json:"id"`
		Name   string `json:"name"`
		Period int    `json:"period"`
		Email  string `json:"email"`
	}

	AlarmStatus struct {
		Id                  string `json:"id"`
		HostName            string `json:"hostname"`
		AlarmDefinitionName string `json:"alarmDefinitionName"`
		AlarmDefinitionId   string `json:"alarmDefinitionId"`
		MetricName          string `json:"metricName"`
		Expression          string `json:"expression"`
		Type                string `json:"type"`
		Zone                string `json:"zone"`
		Severity            string `json:"severity"`
		State               string `json:"state"`
		UpdateTime          string `json:"updateTime"`
	}

	Alarm struct {
		Id                string
		AlarmDefinitionId string
		Name              string
		Expression        string
		Severity          string
	}

	AlarmHistory struct {
		Id       string `json:"alarmId"`
		Time     string `json:"time"`
		NewState string `json:"newState"`
		OldState string `json:"oldState"`
		Reason   string `json:"reason"`
	}

	AlarmActionHistory struct {
		Id              uint      `gorm:"primary_key"`
		AlarmId         string    `gorm:"type:varchar(36);not null;"`
		AlarmActionDesc string    `gorm:"type:text;"`
		RegDate         time.Time `gorm:"type:datetime;DEFAULT:current_timestamp;not null;"`
		RegUser         string    `gorm:"type:varchar(36);DEFAULT:'system';not null;"`
		ModiDate        time.Time `gorm:"type:datetime;DEFAULT:current_timestamp;null;"`
		ModiUser        string    `gorm:"type:varchar(36);DEFAULT:'system';null;"`
	}

	AlarmActionRequest struct {
		Id              uint
		AlarmId         string
		AlarmActionDesc string
		RegDate         time.Time
		RegUser         string
		ModiDate        time.Time
		ModiUser        string
	}

	AlarmActionResponse struct {
		Id              uint   `json:"id"`
		AlarmId         string `json:"alarmId"`
		AlarmActionDesc string `json:"alarmActionDesc"`
		RegDate         string `json:"regDate"`
		RegUser         string `json:"regUser"`
		ModiDate        string `json:"ModiDate"`
		ModiUser        string `json:"modiUser"`
	}

	AlarmRealtimeCountResponse struct {
		TotalCnt    int `json:"totalCnt"`
		WarningCnt  int `json:"warningCnt"`
		CriticalCnt int `json:"criticalCnt"`
	}
)

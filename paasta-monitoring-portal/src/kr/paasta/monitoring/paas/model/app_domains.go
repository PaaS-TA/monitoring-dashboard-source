package model

import (
	"errors"
	"time"
)

const (
	RESULT_SUCCESS = "success"
	RESULT_FAIL    = "fail"
)

type (
	AppAutoScalingPolicies struct {
		AppGuid             string    `gorm:"type:varchar(50);not null;primary_key;"`
		InstanceMinCnt      uint      `gorm:"type:int unsigned;not null;"`
		InstanceMaxCnt      uint      `gorm:"type:int unsigned;not null;"`
		CpuMinThreshold     uint      `gorm:"type:int unsigned;not null;"`
		CpuMaxThreshold     uint      `gorm:"type:int unsigned;not null;"`
		MemoryMinThreshold  uint      `gorm:"type:int unsigned;not null;"`
		MemoryMaxThreshold  uint      `gorm:"type:int unsigned;not null;"`
		InstanceScalingUnit uint      `gorm:"type:int unsigned;not null;"`
		MeasureTimeSec      uint      `gorm:"type:int unsigned;not null;"`
		AutoScalingOutYn    string    `gorm:"type:varchar(1);not null;"`
		AutoScalingInYn     string    `gorm:"type:varchar(1);not null;"`
		AutoScalingCpuYn    string    `gorm:"type:varchar(1);not null;"`
		AutoScalingMemoryYn string    `gorm:"type:varchar(1);not null;"`
		RegDate             time.Time `gorm:"type:datetime;not null;DEFAULT:CURRENT_TIMESTAMP;"`
		RegUser             string    `gorm:"type:varchar(36);not null;DEFAULT:'system';"`
		ModiDate            time.Time `gorm:"type:datetime;not null;DEFAULT:CURRENT_TIMESTAMP;"`
		ModiUser            string    `gorm:"type:varchar(36);not null;DEFAULT:'system';"`
	}

	AppAlarmPolicies struct {
		AppGuid                 string    `gorm:"type:varchar(50);not null;primary_key;"`
		CpuWarningThreshold     uint      `gorm:"type:int unsigned;not null;"`
		CpuCriticalThreshold    uint      `gorm:"type:int unsigned;not null;"`
		MemoryWarningThreshold  uint      `gorm:"type:int unsigned;not null;"`
		MemoryCriticalThreshold uint      `gorm:"type:int unsigned;not null;"`
		MeasureTimeSec          uint      `gorm:"type:int unsigned;not null;"`
		Email                   string    `gorm:"type:varchar(100);null;DEFAULT:null;"`
		EmailSendYn             string    `gorm:"type:varchar(1);not null;"`
		AlarmUseYn              string    `gorm:"type:varchar(1);not null;"`
		RegDate                 time.Time `gorm:"type:datetime;not null;DEFAULT:CURRENT_TIMESTAMP;"`
		RegUser                 string    `gorm:"type:varchar(36);not null;DEFAULT:'system';"`
		ModiDate                time.Time `gorm:"type:datetime;not null;DEFAULT:CURRENT_TIMESTAMP;"`
		ModiUser                string    `gorm:"type:varchar(36);not null;DEFAULT:'system';"`
	}

	AppAlarmReq struct {
		AppGuid        string
		PageItems      int
		PageIndex      int
		ResourceType   string
		AlarmLevel     string
		SearchDateFrom string
		SearchDateTo   string
	}

	AppAutoscalingPolicy struct {
		AppGuid             string `json:"appGuid"`
		InstanceMinCnt      int    `json:"instanceMinCnt"`
		InstanceMaxCnt      int    `json:"instanceMaxCnt"`
		CpuMinThreshold     int    `json:"cpuMinThreshold"`
		CpuMaxThreshold     int    `json:"cpuMaxThreshold"`
		MemoryMinThreshold  int    `json:"memoryMinThreshold"`
		MemoryMaxThreshold  int    `json:"memoryMaxThreshold"`
		InstanceScalingUnit int    `json:"instanceVariationUnit"`
		MeasureTimeSec      int    `json:"measureTimeSec"`
		AutoScalingOutYn    string `json:"autoScalingOutYn"`
		AutoScalingInYn     string `json:"autoScalingInYn"`
		AutoScalingCpuYn    string `json:"autoScalingCpuYn"`
		AutoScalingMemoryYn string `json:"autoScalingMemoryYn"`
	}

	AppAlarmPolicy struct {
		AppGuid                 string `json:"appGuid"`
		CpuWarningThreshold     int    `json:"cpuWarningThreshold"`
		CpuCriticalThreshold    int    `json:"cpuCriticalThreshold"`
		MemoryWarningThreshold  int    `json:"memoryWarningThreshold"`
		MemoryCriticalThreshold int    `json:"memoryCriticalThreshold"`
		MeasureTimeSec          int    `json:"measureTimeSec"`
		Email                   string `json:"email"`
		EmailSendYn             string `json:"emailSendYn"`
		AlarmUseYn              string `json:"alarmUseYn"`
	}

	AppAlarmPagingRes struct {
		PagingRes
		AppAlarmList []AppAlarm `json:"data"`
	}

	AppAlarm struct {
		AlarmId      string   `json:"alarmId"`
		AppGuid      string   `json:"appGuid"`
		AppIdx       string   `json:"appIdx"`
		AppName      string   `json:"appName"`
		ResourceType string   `json:"resourceType"`
		AlarmLevel   string   `json:"alarmLevel"`
		AlarmTitle   string   `json:"alarmTitle"`
		AlarmMessage string   `json:"alarmMessage"`
		RegDate      JSONTime `json:"regDate"`
	}

	ResultResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
)

func (a AppAlarmReq) DefaultAlarmListValidate(req AppAlarmReq) error {
	if req.AppGuid == "" {
		return errors.New("Required input value does not exist. [AppGuid]")
	}
	if req.PageIndex == 0 {
		return errors.New("Required input value does not exist. [pageIndex]")
	}
	if req.PageItems == 0 {
		return errors.New("Required input value does not exist. [pageItems]")
	}

	return nil
}

func (a AppAlarmPolicy) DefaultAlarmPolicyValidate(req AppAlarmPolicy) error {
	if req.AppGuid == "" {
		return errors.New("Required input value does not exist. [AppGuid]")
	}

	if req.CpuWarningThreshold <= 0 || req.CpuWarningThreshold > 99 {
		return errors.New("Required input value must be bigger than zero and must be smaller than one hundred. [CpuWarningThreshold]")
	}
	if req.CpuCriticalThreshold <= 0 || req.CpuCriticalThreshold > 99 {
		return errors.New("Required input value must be bigger than zero and must be smaller than one hundred. [CpuCriticalThreshold]")
	}
	if req.MemoryWarningThreshold <= 0 || req.MemoryWarningThreshold > 99 {
		return errors.New("Required input value must be bigger than zero and must be smaller than one hundred. [MemoryWarningThreshold]")
	}
	if req.MemoryCriticalThreshold <= 0 || req.MemoryCriticalThreshold > 99 {
		return errors.New("Required input value must be bigger than zero and must be smaller than one hundred. [MemoryCriticalThreshold]")
	}
	if req.MeasureTimeSec <= 0 {
		return errors.New("Required input value must be bigger than zero. [MeasureTimeSec]")
	}
	if req.EmailSendYn == "Y" && req.Email == "" {
		return errors.New("Required input value does not exist. [Email]")
	}

	return nil
}

func (a AppAutoscalingPolicy) DefaultAutoScalingPolicyValidate(req AppAutoscalingPolicy) error {
	if req.AppGuid == "" {
		return errors.New("Required input value does not exist. [AppGuid]")
	}

	if req.InstanceMinCnt <= 0 {
		return errors.New("Required input value must be bigger than zero. [InstanceMinCnt]")
	}
	if req.CpuMinThreshold <= 0 || req.CpuMinThreshold > 99 {
		return errors.New("Required input value must be bigger than zero and must be smaller than one hundred. [CpuMinThreshold]")
	}
	if req.CpuMaxThreshold <= 0 || req.CpuMaxThreshold > 99 {
		return errors.New("Required input value must be bigger than zero and must be smaller than one hundred. [CpuMaxThreshold]")
	}
	if req.MemoryMinThreshold <= 0 || req.MemoryMinThreshold > 99 {
		return errors.New("Required input value must be bigger than zero and must be smaller than one hundred. [MemoryMinThreshold]")
	}
	if req.MemoryMaxThreshold <= 0 || req.MemoryMaxThreshold > 99 {
		return errors.New("Required input value must be bigger than zero and must be smaller than one hundred. [MemoryMaxThreshold]")
	}
	if req.MeasureTimeSec <= 0 {
		return errors.New("Required input value must be bigger than zero. [MeasureTimeSec]")
	}

	if req.InstanceMinCnt > req.InstanceMaxCnt {
		return errors.New("InstanceMaxCnt value must be bigger than InstanceMinCnt")
	}
	if req.CpuMinThreshold >= req.CpuMaxThreshold {
		return errors.New("CpuMaxThreshold value must be bigger than CpuMinThreshold")
	}
	if req.MemoryMinThreshold >= req.MemoryMaxThreshold {
		return errors.New("MemoryMaxThreshold value must be bigger than MemoryMinThreshold")
	}

	if req.AutoScalingCpuYn == "" {
		return errors.New("autoScalingCpuYn cannot be empty")
	}

	if req.AutoScalingMemoryYn == "" {
		return errors.New("autoScalingMemoryYn cannot be empty")
	}

	return nil
}

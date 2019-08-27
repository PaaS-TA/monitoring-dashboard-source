package model

import "time"

type BatchAlarmInfo struct {
	AlarmId     int    `gorm:"type:int unsigned auto_increment;not null;primary_key"`
	ServiceType string `gorm:"type:varchar(10);not null;"`
	MetricType  string `gorm:"type:varchar(50);not null;"`
	//CriticalType       	string     `gorm:"type:varchar(50);not null;"`
	WarningValue   int    `gorm:"type:int;not null;"`
	CriticalValue  int    `gorm:"type:int;not null;"`
	MeasureTime    int    `gorm:"type:int;not null;"`
	CronExpression string `gorm:"type:varchar(50);not null;"`
	ExecMsg        string `gorm:"type:varchar(1000);not null;"`
	ParamData1     string `gorm:"type:varchar(1000);null;"`
	ParamData2     string `gorm:"type:varchar(1000);null;"`
	ParamData3     string `gorm:"type:varchar(1000);null;"`
}

type InsertAlarmInfo struct {
	ServiceType string `gorm:"type:varchar(10);not null;"`
	MetricType  string `gorm:"type:varchar(50);not null;"`
	//CriticalType       	string     `gorm:"type:varchar(50);not null;"`
	WarningValue   int    `gorm:"type:int;not null;"`
	CriticalValue  int    `gorm:"type:int;not null;"`
	MeasureTime    int    `gorm:"type:int;not null;"`
	CronExpression string `gorm:"type:varchar(50);not null;"`
	ExecMsg        string `gorm:"type:varchar(1000);not null;"`
	ParamData1     string `gorm:"type:varchar(1000);null;"`
	ParamData2     string `gorm:"type:varchar(1000);null;"`
	ParamData3     string `gorm:"type:varchar(1000);null;"`
}

type BatchAlarmExecution struct {
	ExcutionId      uint64    `gorm:"type:bigint(20) unsigned auto_increment;not null;primary_key"`
	AlarmId         int       `gorm:"type:int;not null;"`
	ServiceType     string    `gorm:"type:varchar(10);not null;"`
	CriticalStatus  string    `gorm:"type:varchar(50);not null;"`
	MeasureValue    float64   `gorm:"type:float;not null;"`
	MeasureName1    string    `gorm:"type:varchar(200);not null;"`
	MeasureName2    string    `gorm:"type:varchar(200);not null;"`
	MeasureName3    string    `gorm:"type:varchar(200);not null;"`
	ExecutionTime   time.Time `gorm:"type:timestamp;not null;"`
	ExecutionResult string    `gorm:"type:varchar(1000);not null;"`
}

type BatchAlarmReceiver struct {
	ReceiverId  int    `gorm:"type:int unsigned auto_increment;not null;primary_key"`
	ServiceType string `gorm:"type:varchar(10);not null;"`
	Name        string `gorm:"type:varchar(100);not null;"`
	Email       string `gorm:"type:varchar(200);null;"`
	SnsId       int64  `gorm:"type:varchar(200);null;"`
	UseYn       string `gorm:"type:varchar(1);not null;"`
}

type MailContent struct {
	AlarmInfo      BatchAlarmInfo
	AlarmExecution BatchAlarmExecution
}

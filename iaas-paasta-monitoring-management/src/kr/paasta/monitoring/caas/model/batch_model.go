package model

import "time"

type BatchAlarmInfo struct {
	AlarmId        int    `gorm:"type:int unsigned auto_increment;not null;primary_key"`
	ServiceType    string `gorm:"type:varchar(10);not null;unique_index:uix_batch_alarm_info;"`
	MetricType     string `gorm:"type:varchar(50);not null;unique_index:uix_batch_alarm_info;"`
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
	ResolveStatus   string    `gorm:"type:varchar(1);not null;default '1';"`
	CompleteDate    time.Time `gorm:"null;"`
}

type BatchAlarmExecutionResolve struct {
	ResolveId       uint64    `gorm:"type:bigint(20) unsigned auto_increment;not null;primary_key"`
	ExcutionId      uint64    `gorm:"type:bigint(20);not null;"`
	AlarmActionDesc string    `gorm:"type:varchar(4000);null;"`
	RegDate         time.Time `gorm:"null;"`
}

type BatchAlarmReceiver struct {
	ReceiverId  int       `gorm:"type:int unsigned auto_increment;not null;primary_key"`
	ServiceType string    `gorm:"type:varchar(10);not null;unique_index:uix_batch_alarm_sns_targets;"`
	ReceiveType string    `gorm:"type:varchar(100);not null;unique_index:uix_batch_alarm_sns_targets;"`
	TargetId    string    `gorm:"type:varchar(200);not null;unique_index:uix_batch_alarm_sns_targets;"`
	UseYn       string    `gorm:"type:varchar(1);not null;DEFAULT:'Y';"`
	RegDate     time.Time `gorm:"type:datetime;not null;DEFAULT:CURRENT_TIMESTAMP;"`
	RegUser     string    `gorm:"type:varchar(36);not null;DEFAULT:'system';"`
	ModiDate    time.Time `gorm:"type:datetime;not null;DEFAULT:CURRENT_TIMESTAMP;"`
	ModiUser    string    `gorm:"type:varchar(36);not null;DEFAULT:'system';"`
}

type BatchAlarmSns struct {
	ChannelId  int       `gorm:"type:bigint(20) unsigned auto_increment;not null;primary_key"`
	OriginType string    `gorm:"type:varchar(10);not null;unique_index:uix_batch_alarm_sns"`
	SnsId      string    `gorm:"type:varchar(200);not null;"`
	Token      string    `gorm:"type:varchar(1000);not null;"`
	Expl       string    `gorm:"type:varchar(100);null;"`
	SnsSendYn  string    `gorm:"type:varchar(1);not null;DEFAULT:'Y';"`
	RegDate    time.Time `gorm:"type:datetime;not null;DEFAULT:CURRENT_TIMESTAMP;"`
	RegUser    string    `gorm:"type:varchar(36);not null;DEFAULT:'system';"`
	ModiDate   time.Time `gorm:"type:datetime;not null;DEFAULT:CURRENT_TIMESTAMP;"`
	ModiUser   string    `gorm:"type:varchar(36);not null;DEFAULT:'system';"`
}

type MailContent struct {
	AlarmInfo      BatchAlarmInfo
	AlarmExecution []BatchAlarmExecution
}

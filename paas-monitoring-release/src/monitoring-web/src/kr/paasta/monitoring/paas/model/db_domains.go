package model

import (
	"time"
)

type (
	BaseField struct{
		RegDate     time.Time   `gorm:"type:datetime;DEFAULT:current_timestamp;not null;"`
		RegUser     string      `gorm:"type:varchar(36);DEFAULT:'system';not null;"`
		ModiDate    time.Time   `gorm:"type:datetime;DEFAULT:current_timestamp;null;"`
		ModiUser    string      `gorm:"type:varchar(36);DEFAULT:'system';null;"`
	}

	Alarm   struct{
		Id           uint        `gorm:"primary_key"`
		OriginType   string	 `gorm:"type:varchar(3);not null;"`
		OriginId     uint	 `gorm:"type:int;not null;"`
		AlarmType    string      `gorm:"type:varchar(6);not null;"`
		Level        string      `gorm:"type:varchar(8);not null;"`
		Ip          string       `gorm:"type:varchar(15);null;"`
		AppYn        string      `gorm:"type:varchar(1);not null;"`
		AppName      string      `gorm:"type:varchar(500);null;"`
		AppIndex     uint        `gorm:"type:int;null;"`
		ContainerName  string    `gorm:"type:varchar(40);null;"`
		AlarmTitle     string    `gorm:"type:varchar(5000);not null;"`
		AlarmMessage   string    `gorm:"type:text;not null;"`
		ResolveStatus  string  	 `gorm:"type:varchar(1);not null;"`        //처리 여부 1: Alarm 발생, 2: 처리중, 3: 처리 완료
		AlarmCnt       int  	 `gorm:"type:int;not null;DEFAULT:1"`      //Alarm 발생 횟수
		BaseField
		AlarmSendDate   time.Time   `gorm:"type:datetime;null;DEFAULT:null"`           //Alarm 전송 시간
		CompleteDate    time.Time   `gorm:"type:datetime;null;DEFAULT:null"`           //Alarm 처리 완료 시간
		CompleteUser    string      `gorm:"type:varchar(36);null;"`        //Alarm 처리 완료자
	}

	AlarmAction struct{
		Id           uint       `gorm:"primary_key"`
		AlarmId      uint	`gorm:"type:integer;not null;"`
		AlarmActionDesc string	`gorm:"type:text;"`
		RegDate     time.Time   `gorm:"type:datetime;DEFAULT:current_timestamp;not null;"`
		RegUser     string      `gorm:"type:varchar(36);DEFAULT:'system';not null;"`
		ModiDate    time.Time   `gorm:"type:datetime;DEFAULT:current_timestamp;null;"`
		ModiUser    string      `gorm:"type:varchar(36);DEFAULT:'system';null;"`
	}

	AlarmSns struct{
		ChannelId		uint		`gorm:"type:int(10) unsigned;not null;AUTO_INCREMENT;primary_key;"`
		OriginType		string		`gorm:"type:varchar(3);not null;unique_index:ux_alarm_sns;"`
		SnsType			string		`gorm:"type:varchar(50);not null;unique_index:ux_alarm_sns;"`
		SnsId			string   	`gorm:"type:varchar(50);not null;unique_index:ux_alarm_sns;"`
		Token			string      `gorm:"type:varchar(100);not null;"`
		Expl			string      `gorm:"type:varchar(100);not null;"`
		RegDate			time.Time	`gorm:"type:datetime;not null;DEFAULT:CURRENT_TIMESTAMP;"`
		RegUser			string		`gorm:"type:varchar(36);not null;DEFAULT:'system';"`
		ModiDate		time.Time	`gorm:"type:datetime;not null;DEFAULT:CURRENT_TIMESTAMP;"`
		ModiUser		string		`gorm:"type:varchar(36);not null;DEFAULT:'system';"`
		SnsSendYn		string		`gorm:"type:varchar(1);not null;DEFAULT:'Y';"`
	}

	Databases struct{
		BoshDatabase string
		PaastaDatabase string
		ContainerDatabase string
	}
)

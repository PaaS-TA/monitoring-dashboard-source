package model

import "time"

type (
	AlarmRequest struct {
		Nodename      string `json:"nodename"`
		Instance      string `json:"instance"`
		PodName       string `json:"PodName"`
		NameSpace     string `json:"NameSpace"`
		ContainerName string `json:"ContainerName"`
		WorkloadsName string `json:"WorkloadsName"`
	}

	AlarmInfo struct {
		Name     string `json:"Name"`
		Warning  int    `json:"Warning"`
		Critical int    `json:"Critical"`
		Delay    string `json:"Delay"`
		AlarmId  int    `json:"AlarmId"`
	}

	ResultAlarmInfo struct {
		Result        []AlarmInfo `json:"Threshold"`
		MeasuringTime int         `json:"MeasuringTime"`
		AlarmMail     string      `json:"AlarmMail"`
		AlarmTelegram int64       `json:"AlarmTelegram"`
		ReceiverID    int         `json:"ReceiverID"`
	}

	AlarmLog struct {
		WorkNode  string    `json:"WorkNode"`
		NameSpace string    `json:"NameSpace"`
		Pod       string    `json:"Pod"`
		Status    string    `json:"Status"`
		Issue     string    `json:"Issue"`
		Time      time.Time `json:"Time"`
	}
)

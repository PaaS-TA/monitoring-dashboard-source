package model

type (
	BatchAlarmSnsRequest struct {
		ChannelId  int    `json:"id"`
		OriginType string `json:"originType"`
		SnsId      string `json:"snsId"`
		Token      string `json:"token"`
		Expl       string `json:"expl"`
		SnsSendYn  string `json:"snsSendYn"`
	}

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

	AlarmPolicyResponse struct {
		Id                int    `json:"id"`
		OriginType        string `json:"originType"`
		AlarmType         string `json:"alarmType"`
		WarningThreshold  int    `json:"warningThreshold"`
		CriticalThreshold int    `json:"criticalThreshold"`
		RepeatTime        int    `json:"repeatTime"`
		Comment           string `json:"comment"`
		MeasureTime       int    `json:"measureTime"`
		MailAddress       string `json:"mailAddress"`
		MailSendYn        string `json:"mailSendYn"`
	}

	AlarmPolicyRequest struct {
		Id                uint     `json:"id"`
		OriginType        string   `json:"originType"`
		AlarmType         string   `json:"alarmType"`
		WarningThreshold  int      `json:"warningThreshold"`
		CriticalThreshold int      `json:"criticalThreshold"`
		RepeatTime        int      `json:"repeatTime"`
		Comment           string   `json:"comment"`
		MeasureTime       int      `json:"measureTime"`
		MailAddress       string   `json:"mailAddress"`
		SnsType           string   `json:"snsType"`
		SnsId             string   `json:"snsId"`
		Token             string   `json:"token"`
		Expl              string   `json:"expl"`
		MailSendYn        string   `json:"mailSendYn"`
		SnsSendYn         string   `json:"snsSendYn"`
		ModiDate          JSONTime `json:"modiDate"`
		ModiUser          string   `json:"modiUser"`
	}

	ResultAlarmInfo struct {
		Result        []AlarmInfo `json:"Threshold"`
		MeasuringTime int         `json:"MeasuringTime"`
		AlarmMail     string      `json:"AlarmMail"`
		UseYn         string      `json:"UseYn"`
	}

	AlarmLog struct {
		Id                uint64                 `json:"id"`
		Application       string                 `json:"application"`
		Status            string                 `json:"status"`
		Issue             string                 `json:"issue"`
		RegDate           string                 `json:"regDate"`
		ResolveStatusName string                 `json:"resolveStatusName"`
		ResolveStatus     string                 `json:"resolveStatus"`
		CompleteDate      string                 `json:"completeDate"`
		Data              []AlarmrRsolveResponse `json:"data"`
	}

	AlarmCount struct {
		AlarmCnt int `json:"totalCnt"`
	}

	AlarmrRsolveRequest struct {
		Id              uint64 `json:"id"`
		ResolveStatus   string `json:"resolveStatus"`
		AlarmActionDesc string `json:"alarmActionDesc"`
	}

	AlarmrRsolveResponse struct {
		ResolveId       uint64   `json:"id"`
		AlarmActionDesc string   `json:"alarmActionDesc"`
		RegDate         JSONTime `json:"regDate"`
	}

	AlarmrReceiverResponse struct {
		ReceiverId int      `json:"ReceiverId"`
		TargetId   string   `json:"TargetId"`
		RegDate    JSONTime `json:"RegDate"`
		UseYn      string   `json:"UseYn"`
	}

	AlarmActionResponse struct {
		AlarmId         uint64   `json:"alarmId"`
		ResolveId       uint64   `json:"id"`
		AlarmActionDesc string   `json:"alarmActionDesc"`
		RegDate         JSONTime `json:"regDate"`
	}
)

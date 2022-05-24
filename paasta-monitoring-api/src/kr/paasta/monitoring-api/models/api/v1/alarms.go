package v1

// Struct's each field name is field name of tables in PaastaMonitoring Database.
// The name of struct is tables' name.
// For response container.
type (
	Alarms struct {
		Id            int    `json:"id"`
		OriginType    string `json:"originType"`
		OriginId      int    `json:"originId"`
		AlarmType     string `json:"alarmType"`
		Level         string `json:"level"`
		Ip            string `json:"ip"`
		AppYN         string `json:"appYN"`
		AppName       string `json:"appName"`
		AppIndex      int    `json:"appIndex"`
		ContainerName string `json:"containerName"`
		AlarmTitle    string `json:"alarmTitle"`
		AlarmMessage  string `json:"alarmMessage"`
		ResolveStatus string `json:"resolveStatus"`
		AlarmCnt      int    `json:"alarmCnt"`
		RegDate       string `json:"regDate"`
		RegUser       string `json:"regUser"`
		ModiDate      string `json:"modiDate"`
		ModiUser      string `json:"modiUser"`
		AlarmSendDate string `json:"alarmSendDate"`
		CompleteDate  string `json:"completeDate"`
		CompleteUser  string `json:"completeUser"`
	}

	AlarmPolicies struct {
		Id                int    `json:"id"`
		OriginType        string `json:"originType"`
		AlarmType         string `json:"alarmType"`
		WarningThreshold  int    `json:"warningThreshold"`
		CriticalThreshold int    `json:"criticalThreshold"`
		RepeatTime        int    `json:"repeatTime"`
		MeasureTime       int    `json:"measureTime"`
		Comment           string `json:"comment"`
		RegDate           string `json:"regDate"`
		RegUser           string `json:"regUser"`
		ModiDate          string `json:"modiDate"`
		ModiUser          string `json:"modiUser"`
	}
)

// For request container.
type (
	AlarmPolicyRequest struct {
		OriginType        string `json:"originType" validate:"required"`
		AlarmType         string `json:"alarmType"`
		WarningThreshold  int    `json:"warningThreshold"`
		CriticalThreshold int    `json:"criticalThreshold"`
		RepeatTime        int    `json:"repeatTime"`
		MeasureTime       int    `json:"measureTime"`
		MailAddress       string `json:"mailAddress" validate:"email"`
		MailSendYn        string `json:"mailSendYn"`
	}
)

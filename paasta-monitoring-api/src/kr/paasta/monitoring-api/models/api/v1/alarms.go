package v1

import "time"

// Struct's each field name is field name of tables in PaastaMonitoring Database.
// The name of struct is tables' name.
// For response container part.
type (
	Alarms struct {
		Id            int       `json:"id"            gorm:"primaryKey;autoIncrement;-:write"`
		OriginType    string    `json:"originType"    gorm:"not null"`
		OriginId      int       `json:"originId"      gorm:"not null"`
		AlarmType     string    `json:"alarmType"     gorm:"not null"`
		Level         string    `json:"level"         gorm:"not null"`
		Ip            string    `json:"ip"`
		AppYN         string    `json:"appYN"`
		AppName       string    `json:"appName"`
		AppIndex      int       `json:"appIndex"`
		ContainerName string    `json:"containerName"`
		AlarmTitle    string    `json:"alarmTitle"    gorm:"not null"`
		AlarmMessage  string    `json:"alarmMessage"  gorm:"not null"`
		ResolveStatus string    `json:"resolveStatus" gorm:"not null"`
		AlarmCnt      int       `json:"alarmCnt"      gorm:"not null"`
		RegDate       time.Time `json:"regDate"       gorm:"not null"`
		RegUser       string    `json:"regUser"       gorm:"not null"`
		ModiDate      time.Time `json:"modiDate"`
		ModiUser      string    `json:"modiUser"`
		AlarmSendDate time.Time `json:"-"`
		CompleteDate  time.Time `json:"completeDate"`
		CompleteUser  string    `json:"completeUser"`
	}

	AlarmPolicies struct {
		Id                int       `json:"id"                 gorm:"primaryKey;autoIncrement;not null;-:write"`
		OriginType        string    `json:"originType"         gorm:"not null"`
		AlarmType         string    `json:"alarmType"          gorm:"not null"`
		WarningThreshold  int       `json:"warningThreshold"   gorm:"not null"`
		CriticalThreshold int       `json:"criticalThreshold"  gorm:"not null"`
		RepeatTime        int       `json:"repeatTime"         gorm:"not null"`
		MeasureTime       int       `json:"measureTime"        gorm:"not null"`
		Comment           string    `json:"comment"`
		RegDate           time.Time `json:"regDate"            gorm:"<-:create"`
		RegUser           string    `json:"regUser"            gorm:"<-:create"`
		ModiDate          time.Time `json:"modiDate"           gorm:"<-:update"`
		ModiUser          string    `json:"modiUser"           gorm:"<-:update"`
	}

	AlarmSns struct {
		ChannelId  uint      `json:"channelId"    gorm:"primaryKey;autoIncrement;not null;-:write"`
		OriginType string    `json:"originType"   gorm:"not null;default:all"`
		SnsType    string    `json:"snsType"      gorm:"not null"`
		SnsId      string    `json:"snsId"        gorm:"not null"`
		Token      string    `json:"token"`
		Expl       string    `json:"expl"         gorm:"not null"`
		SnsSendYN  string    `json:"snsSendYN"    gorm:"not null;default:Y"`
		RegDate    time.Time `json:"regDate"      gorm:"<-:create"`
		RegUser    string    `json:"regUser"      gorm:"<-:create;default:system"`
		ModiDate   time.Time `json:"modiDate"     gorm:"<-:update"`
		ModiUser   string    `json:"modiUser"     gorm:"<-:update;default:system"`
	}

	AlarmActions struct {
		Id              int       `json:"id"               gorm:"primaryKey;autoIncrement;not null;-:write"`
		AlarmId         int       `json:"alarmId"`
		AlarmActionDesc string    `json:"alarmActionDesc"`
		RegDate         time.Time `json:"regDate"          gorm:"<-:create" swaggerignore:"true"`
		RegUser         string    `json:"regUser"          gorm:"<-:create" swaggerignore:"true"`
		ModiDate        time.Time `json:"modiDate"         gorm:"<-:update" swaggerignore:"true"`
		ModiUser        string    `json:"modiUser"         gorm:"<-:update" swaggerignore:"true"`
	}
)

// 사용자정의형 응답을 위한 구조체 정의 영역.
// JOIN 등으로 생성된 가상 또는 임시 테이블의 결과와 매치시킴.
type (
	CountByTimeline struct {
		Timeline int `json:"timeline"`
		Count    int `json:"count"`
	}
)

// For request container part.
type (
	AlarmPolicyRequest struct {
		OriginType        string `json:"originType" validate:"required"`
		AlarmType         string `json:"alarmType"`
		WarningThreshold  int    `json:"warningThreshold"`
		CriticalThreshold int    `json:"criticalThreshold"`
		RepeatTime        int    `json:"repeatTime"`
		MeasureTime       int    `json:"measureTime"`
	}

	AlarmTargetRequest struct {
		OriginType  string `json:"originType" validate:"required"`
		MailAddress string `json:"mailAddress" validate:"email"`
		MailSendYN  string `json:"mailSendYN"`
	}

	SnsAccountRequest struct {
		OriginType string `json:"originType"`
		SnsType    string `json:"snsType"`
		SnsId      string `json:"snsId"`
		Token      string `json:"token"`
		Expl       string `json:"expl"`
		SnsSendYN  string `json:"snsSendYN"`
	}

	AlarmActionRequest struct {
		Id              int    `json:"id"`
		AlarmId         int    `json:"alarmId"`
		AlarmActionDesc string `json:"alarmActionDesc"`
		RegUser         string `json:"regUser" swaggerignore:"true"`
	}

	AlarmStatisticsCriteriaRequest struct {
		Alias      string `json:"alias"`
		AlarmLevel string `json:"alarmLevel"`
		Service    string `json:"service"`
		Resource   string `json:"resource"`
	}

	AlarmStatisticsParam struct {
		OriginType    string
		ResourceType  string
		AliasPrefix   string
		Period        string
		TimeCriterion string
		DateFormat    string
		ExtraParams   interface{}
	}
)

package models

type(
	AlarmChannelInfoReq struct {
		OriginType  string
		OriginId    uint
		ChannelType string
		UserId	    string
	}

	AlarmChannelInfoResp struct {
		Name        string
		Email       string
		ChannelType string
	}
	AlarmConfig struct {
		SmtpHost       string
		Port           string
		MailSender     string
		SenderPwd      string
		RocketChannel  string
		MailResource   string
		MailReceiver   string
		AlarmSend      bool
	}

)
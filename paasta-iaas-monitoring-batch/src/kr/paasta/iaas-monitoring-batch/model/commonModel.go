package model

import (
	"net/http"
	"crypto/tls"
)

var PortalUrl    string
var PortalClient *http.Client

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
		MailTlsSend    bool
	}

	AlarmItemMeasureTime struct{
		Item	string
		MeasureTime	int
	}

	Mail struct {
		Sender  string
		To      []string
		Cc      []string
		Bcc     []string
		Subject string
		Body    string
	}

	SmtpServer struct {
		Host      string
		Port      string
		TlsConfig *tls.Config
	}

	MailContents struct {
		Status       string
		Message      string
		StatusDetail string
		ServerName   string
		AlarmDate    string
		ElapseTime   string
	}
)
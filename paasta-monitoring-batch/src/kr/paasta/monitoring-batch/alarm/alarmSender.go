package alarm

import (
	"crypto/tls"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	client "github.com/influxdata/influxdb/client/v2"
	"github.com/jinzhu/gorm"
	"monitoring-batch/dao"
	mod "monitoring-batch/model"
	models "monitoring-batch/model"
	cb "monitoring-batch/model/base"
	"monitoring-batch/util"
	"net/smtp"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
type Mail struct {
	Sender  string
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Body    string
}

type SmtpServer struct {
	Host      string
	Port      string
	TlsConfig *tls.Config
}

type MailContents struct {
	Status       string
	Message      string
	StatusDetail string
	ServerName   string
	AlarmDate    string
	ElapseTime   string
}
*/
type AlarmService struct {
	alarmConfig *models.AlarmConfig
}

func GetAlarmService(alarmConfig models.AlarmConfig) *AlarmService {
	return &AlarmService{
		alarmConfig: &alarmConfig,
	}
}

func (f AlarmService) DBAlarmMessageBuild(alarmSource cb.AlarmThreshold, currentSystemUsage float64) cb.Alarm {

	fmt.Println("Before----->", alarmSource)
	fmt.Println("Threshold----->", alarmSource.Threshold)
	var alarmData cb.Alarm
	alarmData.OriginType = alarmSource.OriginType
	alarmData.OriginId = alarmSource.OriginId
	alarmData.AlarmType = alarmSource.AlarmType
	alarmData.Level = alarmSource.Level
	alarmData.Ip = alarmSource.Ip
	alarmData.ResolveStatus = "1"

	if alarmData.OriginType != cb.ORIGIN_TYPE_CONTAINER {
		alarmData.AppYn = "N"

		switch alarmSource.Level {

		case cb.ALARM_LEVEL_WARNING:
			alarmData.AlarmTitle = "[" + alarmSource.ServiceName + "]의 " + alarmSource.AlarmType + " 상태 [" + alarmSource.Level + "]"
			alarmMsg := alarmSource.ServiceName + " " + alarmSource.AlarmType + "의 상태" + alarmSource.Level + "\n"
			alarmMsg += alarmSource.AlarmType + " 의 임계치인" + util.Floattostrwithprec(alarmSource.Threshold, 0) + "%를 초과하였습니다."
			alarmMsg += " \n 현재 사용률 [" + util.FloattostrDigit2(currentSystemUsage) + "]% 입니다. "

			alarmData.AlarmMessage = alarmMsg
		case cb.ALARM_LEVEL_CRITICAL:
			alarmData.AlarmTitle = "[" + alarmSource.ServiceName + "]의 " + alarmSource.AlarmType + " 상태 [" + alarmSource.Level + "]"
			alarmMsg := alarmSource.ServiceName + " " + alarmSource.AlarmType + " 의 상태" + alarmSource.Level + "\n"
			alarmMsg += alarmSource.AlarmType + " 의 임계치인" + util.Floattostrwithprec(alarmSource.Threshold, 0) + "%를 초과하였습니다."
			alarmMsg += " \n 현재 사용률 [" + util.FloattostrDigit2(currentSystemUsage) + "]% 입니다. "
			alarmData.AlarmMessage = alarmMsg
		case cb.ALARM_LEVEL_FAIL:
			alarmData.AlarmTitle = "[" + alarmSource.ServiceName + "]가 다운되었습니다."
			alarmData.AlarmMessage = "[" + alarmSource.ServiceName + "]가 다운되었습니다."
		}

	}

	return alarmData
}

func (f AlarmService) DBAppAlarmMessageBuild(alarmSource cb.AlarmThreshold, appInfo mod.ContainerTileView, cellName string, currentSystemUsage float64) cb.Alarm {

	var alarmData cb.Alarm
	alarmData.OriginType = alarmSource.OriginType
	alarmData.OriginId = alarmSource.OriginId
	alarmData.AlarmType = alarmSource.AlarmType
	alarmData.Level = alarmSource.Level
	alarmData.Ip = alarmSource.Ip
	alarmData.ResolveStatus = "1"
	alarmData.AppName = appInfo.AppName
	alarmData.AppIndex = appInfo.AppIndex
	alarmData.ContainerName = appInfo.ContainerName

	if alarmData.OriginType == cb.ORIGIN_TYPE_CONTAINER {
		alarmData.AppYn = "Y"

		switch alarmSource.Level {

		case cb.ALARM_LEVEL_WARNING:
			alarmMsg := "App상태:Warning. " + alarmData.AppName + "[" + strconv.Itoa(alarmData.AppIndex) + "]" + ", CellName:[" + cellName + "], " + alarmSource.AlarmType + "의 상태가 " + alarmSource.Level + " 입니다. \n"
			alarmMsg += alarmData.AppName + "의 임계치 " + util.Floattostrwithprec(alarmSource.Threshold, 0) + "% 를 초과했습니다. "
			alarmMsg += " \n 현재 사용률 [" + util.FloattostrDigit2(currentSystemUsage) + "]% 입니다. "
			alarmData.AlarmMessage = alarmMsg
			alarmData.AlarmTitle = alarmData.AppName + "[" + strconv.Itoa(alarmData.AppIndex) + "]" + " App 상태 Warning. "
		case cb.ALARM_LEVEL_CRITICAL:
			alarmMsg := "App상태:Critical. " + alarmData.AppName + "[" + strconv.Itoa(alarmData.AppIndex) + "]" + ", CellName:[" + cellName + "], " + alarmSource.AlarmType + "의 상태가 " + alarmSource.Level + " 입니다. \n"
			alarmMsg += alarmData.AppName + "의 임계치 " + util.Floattostrwithprec(alarmSource.Threshold, 0) + "% 를 초과했습니다. "
			alarmMsg += " \n 현재 사용률 [" + util.FloattostrDigit2(currentSystemUsage) + "]% 입니다. "
			alarmData.AlarmMessage = alarmMsg
			alarmData.AlarmTitle = alarmData.AppName + "[" + strconv.Itoa(alarmData.AppIndex) + "]" + " App 상태 Critical. "
		case cb.ALARM_LEVEL_FAIL:
			alarmData.AlarmTitle = "[" + alarmData.AppName + "]가 Down 되었습니다."
			alarmData.AlarmMessage = "[" + alarmData.AppName + "]가 Down 되었습니다."
		}
	}

	return alarmData
}

func (f AlarmService) MailAlarmMessageBuild(alarmSource cb.AlarmThreshold, alarmDate, elaspeTime string, currentSystemUsage float64) mod.MailContents {

	var mail mod.MailContents

	mail.Status = alarmSource.Level
	if alarmSource.Level == cb.ALARM_LEVEL_FAIL {
		mail.Message = alarmSource.ServiceName + "이 Down 되었습니다."
	} else {
		mail.Message = alarmSource.ServiceName + "의 " + alarmSource.AlarmType + " 임계치 [" + util.Floattostrwithprec(alarmSource.Threshold, 0) + "]%를 초과했습니다. \n 현재 사용률 [" + util.FloattostrDigit2(currentSystemUsage) + "]% 입니다."
	}

	mail.ServerName = alarmSource.ServiceName
	mail.Status = alarmSource.Level
	mail.AlarmDate = alarmDate
	mail.ElapseTime = elaspeTime + "(분)"
	return mail
}

func (f AlarmService) MailAppAlarmMessageBuild(alarmSource cb.AlarmThreshold, alarmData cb.Alarm, alarmDate, elaspeTime string, currentSystemUsage float64) mod.MailContents {

	var mail mod.MailContents
	mail.Status = alarmSource.Level
	if alarmSource.Level == cb.ALARM_LEVEL_FAIL {
		mail.Message = alarmData.AppName + "이 Down되었습니다."
	} else {
		mail.Message = alarmData.AppName + "[" + strconv.Itoa(alarmData.AppIndex) + "]" + "의 " + alarmSource.AlarmType + " 임계치 [" + util.Floattostrwithprec(alarmSource.Threshold, 0) + "]%를 초과했습니다. \n 현재 사용률 [" + util.FloattostrDigit2(currentSystemUsage) + "]% 입니다."
	}

	mail.ServerName = alarmData.AppName + "[" + strconv.Itoa(alarmData.AppIndex) + "]"
	mail.Status = alarmSource.Level
	mail.AlarmDate = alarmDate
	mail.ElapseTime = elaspeTime + "(분)"

	return mail
	//return f.GetMailForm(mail)
}

func (f AlarmService) AlarmSend(alarmSource cb.AlarmThreshold, alarmData cb.Alarm, txn *gorm.DB, metricClient client.Client, alarmConfig mod.AlarmConfig, currentSystemUsage float64) {

	var alarmDataDB models.Alarm
	if alarmData.OriginType != cb.ORIGIN_TYPE_CONTAINER {
		_, alarmDataDB = dao.GetCommonDao(metricClient).GetAlarmData(alarmData, txn)
	} else {
		_, alarmDataDB = dao.GetContainerAlarmDao(metricClient).GetAlarmData(alarmData, txn)
	}

	//Gmt 시간에서 현재 시간 얻어 온다.
	alarmDate := alarmDataDB.RegDate.Add(time.Duration(9) * time.Hour).Format("2006-01-02 15:04:05")
	nowDate := time.Now()

	//경과시간 계산(분)
	elaspeTime := nowDate.Unix() - alarmDataDB.RegDate.Unix()
	elaspeTimeMinute := elaspeTime / 60

	var message mod.MailContents
	if alarmData.OriginType != cb.ORIGIN_TYPE_CONTAINER {
		message = f.MailAlarmMessageBuild(alarmSource, alarmDate, strconv.FormatInt(elaspeTimeMinute, 10), currentSystemUsage)
	} else {
		message = f.MailAppAlarmMessageBuild(alarmSource, alarmData, alarmDate, strconv.FormatInt(elaspeTimeMinute, 10), currentSystemUsage)
	}

	emailTargetList := make([]models.AlarmChannelInfoResp, 1)

	mailReceiver := dao.GetContainerAlarmDao(metricClient).GetAlarmTarget(alarmData.OriginType, txn)

	emailTargetList[0].Email = mailReceiver.MailAddress //alarmConfig.MailReceiver
	fmt.Println("mailReceiver.MailAddress===>", mailReceiver.MailAddress)
	if f.alarmConfig.AlarmSend == true {
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			f.SendMail(alarmData.AlarmTitle, message, emailTargetList)
		}()
		go func() {
			defer wg.Done()
			f.SendTelegram(metricClient, txn, alarmData, message)
		}()
		wg.Wait()
	}
}

//Send Alarm Message To Telegram
func (f AlarmService) SendTelegram(metricClient client.Client, txn *gorm.DB, alarmData cb.Alarm, message mod.MailContents) {
	alarmSns, err := dao.GetCommonDao(metricClient).GetAlarmSns(alarmData.OriginType, txn)
	if len(alarmSns) < 1 {
		return
	}
	if err != nil {
		fmt.Println("Failed to get SNS alarm target! :", err)
		return
	} else {
		text := fmt.Sprintf("status: %s\nname: %s\ndate: %s\nelapsed: %s\n%s", message.Status, message.ServerName, message.AlarmDate, message.ElapseTime, message.Message)
		for _, sns := range alarmSns {
			if strings.ToUpper(sns.SnsSendYn) == "Y" {
				alarmSnsTarget, err := dao.GetCommonDao(metricClient).GetAlarmSnsTarget(sns, txn)
				if err != nil {
					fmt.Println("Failed to get target(user unique key)! :", err)
					return
				}
				for _, target := range alarmSnsTarget {
					bot, err := tgbotapi.NewBotAPI(sns.Token)
					if err != nil {
						fmt.Println("Failed to get telegram client connection! :", err)
					} else {
						bot.Debug = true
						botMsg, botErr := bot.Send(tgbotapi.NewMessage(target.TargetId, text))
						fmt.Printf(">>>>> botMsg=[%v], botErr[%v]\n", botMsg, botErr)
					}
				}
			}
		}
	}
}

func (f AlarmService) SendMail(subject string, body mod.MailContents, receiver []models.AlarmChannelInfoResp) {
	mail := mod.Mail{}
	mail.Sender = f.alarmConfig.MailSender
	var mailReceivers []string
	for _, data := range receiver {
		mailReceivers = append(mailReceivers, data.Email)
	}

	mail.To = mailReceivers
	mail.Subject = subject

	messageBody := f.buildMailMessage(subject, body)
	mail.Body = messageBody

	smtpServer := mod.SmtpServer{Host: f.alarmConfig.SmtpHost, Port: f.alarmConfig.Port}

	var client *smtp.Client
	var err error
	var conn *tls.Conn

	if f.alarmConfig.MailTlsSend == true {
		fmt.Println(">>>>>>>>>>>>>>> Sendmail STMP TLS Mode")
		smtpServer.TlsConfig = &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         smtpServer.Host,
		}
		auth := smtp.PlainAuth("", mail.Sender, f.alarmConfig.SenderPwd, smtpServer.Host)

		conn, err = tls.Dial("tcp", smtpServer.Host+":"+smtpServer.Port, smtpServer.TlsConfig)
		if err != nil {
			fmt.Println("smtp TLS connection error :", err)
			return
		}

		client, err = smtp.NewClient(conn, smtpServer.Host)
		if err != nil {
			fmt.Println("smtp TLS new clinet create error :", err)
			return
		}

		// step 1: Use Auth
		if err = client.Auth(auth); err != nil {
			fmt.Println("TLS client auth error :", err)
			return
		}
	} else {
		fmt.Println(">>>>>>>>>>>>>>> Sendmail STMP Mode")
		client, err = smtp.Dial(smtpServer.Host + ":" + smtpServer.Port)
		if err != nil {
			fmt.Println("smtp connection & client create error :", err)
			return
		}
	}

	// step 2: add all from and to
	if err = client.Mail(mail.Sender); err != nil {
		fmt.Println("client send mail error :", err)
		return
	}

	receivers := append(mail.To, mail.Cc...)
	receivers = append(receivers, mail.Bcc...)

	for _, k := range receivers {
		fmt.Println("sending to: ", k)
		if err = client.Rcpt(k); err != nil {
			fmt.Println("sending error :", err)
			return
		}
	}

	// Data
	w, err := client.Data()
	if err != nil {
		fmt.Println(" client send data error :", err)
		return
	}

	_, err = w.Write([]byte(messageBody))
	if err != nil {
		fmt.Println("write message error :", err)
		return
	}

	err = w.Close()
	if err != nil {
		fmt.Println("client close error:", err)
		return
	}

	//client.Quit()
	fmt.Println(">>>>>>>>Mail sent successfully")
	defer func() {

		if err := recover(); err != nil {
			//fmt.Println("######## Mail sent catch panic")
			fmt.Println(err)
		}

		if f.alarmConfig.MailTlsSend == true {
			if conn != nil {
				//fmt.Println("######## smtp TLS conn close")
				conn.Close()
			}
		}

		if client != nil {
			//fmt.Println("######## smtp client close")
			client.Close()
		}

	}()
}

func (f *AlarmService) buildMailMessage(subject string, mailContent mod.MailContents) string {

	mailForm := ""
	mailForm = mailForm + "Subject: " + subject + "!\n"

	mailForm = mailForm + "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	mailForm += ""
	mailForm += "<!DOCTYPE html>"
	mailForm += "<html>"
	mailForm += "<head>"
	mailForm += "<title>PaaS-TA Monitor</title>"
	mailForm += "<meta charset=\"utf-8\">"
	mailForm += "<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">"
	mailForm += "<style>"
	mailForm += "#text_area:before{display: inline;content: '';width: 158.16px;height: 6px;background-color: #124379;position: absolute;top: 0;left: 0;z-index: 999;}"
	mailForm += "#text_area p strong:before{display:inline-block;content:'';width:23px;height:24px;background:url(" + f.alarmConfig.MailResource + "/public/resources/img/email/ic_error.png) no-repeat;background-size:100%;vertical-align:bottom;margin-right:7px;}"
	mailForm += "</style>"
	mailForm += "</head>"
	mailForm += "<body style=\"*word-break:break-all;-ms-word-break:break-all;margin:0;padding:0;font-family:'Noto Sans', sans-serif;font-size:12px;color:#555;\">"
	mailForm += "<div id=\"wrap\" style=\"margin:0;padding:0;width:100%;height:100%;background-size:cover;overflow-y:auto;\">"
	mailForm += "    <div class=\"email_form\" style=\"width: 708px;height: 361px;padding: 19px 30px;overflow: hidden;margin: 0 auto;display: block;background: #f8f8f8;\"><div class=\"form_area\" style=\"margin:0;padding:0\">"
	mailForm += "        <div class=\"contents_area\" style=\"margin:0;padding:0\">"
	mailForm += "            <div class=\"contents-header\" style=\"margin:0;padding:0;padding-bottom:25px;\">"
	mailForm += "                <h2 style=\"margin:0;padding:0;font-family:'Noto Sans', sans-serif;font-size:12px;color:#555;\"><img src=\"" + f.alarmConfig.MailResource + "/public/resources/img/email/logo.png\" alt=\"monitorxpert\" style=\"border:0 none;width:279px;\"></h2>"
	mailForm += "            </div>"
	mailForm += "            <div class=\"contents-body\" style=\"margin:0;padding:0;position:relative;\">"
	mailForm += "                <div id=\"text_area\" style=\"margin:0;padding:18px 18px 24px;padding-bottom:12px;border:1px solid #ddd;background:#fff;height:auto;border-top:6px solid #007dc4;\">"
	mailForm += "                    <p class=\"text\" style=\"font-size: 15.5px;line-height: 26px;color: #404040;text-align: center;padding-bottom: 8px;font-family: 'Hanna', sans-serif;\">"
	mailForm += "                        <strong>" + mailContent.Status + " : </strong>"
	mailForm += "                        " + mailContent.Message
	mailForm += "                    </p>"
	mailForm += "                    <a href=\"" + f.alarmConfig.MailResource + "\"> <div id=\"text_area\" style=\"margin:0 auto;padding:5px 12px;font-family:'Noto Sans', sans-serif;font-size:14px;color:#fff;vertical-align:middle;text-align:center;outline:none;display:block;border-radius:3px;box-shadow:none;border:none;background:#124379;font-weight:600;cursor:pointer;\">PaaS-Ta Monitor로 이동</div></a>"
	mailForm += "                    <div class=\"notic\" style=\"margin:0;padding:0;padding-top:35px;\">"
	mailForm += "                        <h3 style=\"margin:0;padding:0;font-family:'Noto Sans', sans-serif;font-size:15px;color:#555;font-weight:600;padding-bottom:6px;\">&#50508;&#47548; &#51221;&#48372;</h3>"
	mailForm += "                        <table style=\"text-align: left;border-top: 3px solid #888;border-bottom: 1px solid #bbb;width: 100%;border-spacing: 0;border-collapse: collapse;color: #404040;\">"
	mailForm += "                            <tr>"
	mailForm += "                                <th style=\"background-color:#f5f5f5;padding: 4px;\">&#49345;&#53468;</th>"
	mailForm += "                                <th style=\"background-color:#f5f5f5;padding:4px;\">&#49436;&#48260; &#51060;&#47492;</th>"
	mailForm += "                                <th style=\"background-color:#f5f5f5;padding:4px;\">&#48156;&#49373; &#49884;&#44036;</th>"
	mailForm += "                                <th style=\"background-color:#f5f5f5;padding:4px;\">&#44221;&#44284; &#49884;&#44036;(&#48516;)</th>"
	mailForm += "                            </tr>"
	mailForm += "                            <tr>"
	mailForm += "                                <td style=\"padding: 8px 6px;\">" + mailContent.Status + "</td>"
	mailForm += "                                <td style=\"padding: 8px 6px;\">" + mailContent.ServerName + "</td>"
	mailForm += "                                <td style=\"padding: 8px 6px;\">" + mailContent.AlarmDate + "</td>"
	mailForm += "                                <td style=\"padding: 8px 6px;\">" + mailContent.ElapseTime + "</td>"
	mailForm += "                            </tr>"
	//mailForm += "                            <tr></tr>"
	mailForm += "                        </table>"
	mailForm += "                    </div>"
	mailForm += "                </div>"
	mailForm += "                <p class=\"copyright\" style=\"color: #6f6f6f;\"> Copyright &copy; PaaS-Ta All rights reserved </p>"
	//mailForm += "                <p></p>"
	mailForm += "            </div>"
	mailForm += "        </div>"
	mailForm += "    </div>"
	mailForm += "</div>"
	mailForm += "</body>"
	mailForm += "</html>"

	return mailForm
}

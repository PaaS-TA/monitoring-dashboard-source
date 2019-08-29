package notify

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	//"net/smtp"
	"gopkg.in/gomail.v2"
	"kr/paasta/batch/model"
	"kr/paasta/batch/util"
	"os"
)

var sender string
var resourceUrl string
var dialer *gomail.Dialer

//var message *gomail.Message
func init() {
	config, err := util.ReadConfig("config.ini")
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	smtpHost := config["mail.smtp.host"]
	smtpPort, _ := strconv.Atoi(config["mail.smtp.port"])
	username := config["mail.username"]
	password := config["mail.password"]
	sender = config["mail.sender"]
	resourceUrl = config["mail.resource.url"]

	dialer = gomail.NewDialer(smtpHost, smtpPort, username, password)

}

func SendMail(serviceType string, receivers []string, mailContent model.MailContent) {
	if len(receivers) > 0 {
		message := gomail.NewMessage()

		mailTile := mailContent.AlarmInfo.ServiceType + " [" + mailContent.AlarmInfo.MetricType + "] 사용률 임계치 초과"
		message.SetAddressHeader("From", sender, serviceType)
		message.SetHeader("Subject", mailTile)
		message.SetBody("text/html", buildMailMessage(mailContent))

		emails := strings.Join(receivers, ",")
		message.SetHeader("To", emails)

		if err := dialer.DialAndSend(message); err != nil {
			log.Printf("Could not send email to %s: %v", emails, err)
		}
	}

}

func buildMailMessage(mailContent model.MailContent) string {

	fmt.Printf("mail length : %v\n", len(mailContent.AlarmExecution))

	//measureValue := fmt.Sprintf("%.1f", mailContent.AlarmExecution.MeasureValue)
	mailTile := mailContent.AlarmInfo.ServiceType + " [" + mailContent.AlarmInfo.MetricType + "] 사용률 임계치 초과"
	mailForm := ""
	//mailForm = mailForm + "Subject: " + mailTile + "!\n"

	//mailForm = mailForm + "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	mailForm += ""
	mailForm += "<!DOCTYPE html>"
	mailForm += "<html>"
	mailForm += "<head>"
	mailForm += "<title>PaaS-TA Monitor</title>"
	mailForm += "<meta charset=\"utf-8\">"
	mailForm += "<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">"
	mailForm += "<style>"
	mailForm += "#text_area:before{display: inline;content: '';width: 158.16px;height: 6px;background-color: #124379;position: absolute;top: 0;left: 0;z-index: 999;}"
	mailForm += "#text_area p strong:before{display:inline-block;content:'';width:23px;height:24px;background-size:100%;vertical-align:bottom;margin-right:7px;}"
	//mailForm += "#text_area p strong:before{display:inline-block;content:'';width:23px;height:24px;background:url(" + resourceUrl + "/public/src/assets/images/email/ic_error.png) no-repeat;background-size:100%;vertical-align:bottom;margin-right:7px;}"
	mailForm += "</style>"
	mailForm += "</head>"
	mailForm += "<body style=\"*word-break:break-all;-ms-word-break:break-all;margin:0;padding:0;font-family:'Noto Sans', sans-serif;font-size:12px;color:#555;\">"
	mailForm += "<div id=\"wrap\" style=\"margin:0;padding:0;width:100%;height:100%;background-size:cover;overflow-y:auto;\">"
	mailForm += "    <div class=\"email_form\" style=\"width: 90%;height: 90%;padding: 19px 30px;overflow: hidden;margin: 0 auto;display: block;background: #f8f8f8;\"><div class=\"form_area\" style=\"margin:0;padding:0\">"
	mailForm += "        <div class=\"contents_area\" style=\"margin:0;padding:0\">"
	//mailForm += "            <div class=\"contents-header\" style=\"margin:0;padding:0;padding-bottom:25px;\">"
	//mailForm += "                <h2 style=\"margin:0;padding:0;font-family:'Noto Sans', sans-serif;font-size:12px;color:#555;\"><img src=\"" + resourceUrl + "/public/src/assets/images/logo.png\" alt=\"monitorxpert\" style=\"border:0 none;width:279px;\"></h2>"
	//mailForm += "            </div>"
	mailForm += "            <div class=\"contents-body\" style=\"margin:0;padding:0;position:relative;\">"
	mailForm += "                <div id=\"text_area\" style=\"margin:0;padding:18px 18px 24px;padding-bottom:12px;border:1px solid #ddd;background:#fff;height:auto;border-top:6px solid #007dc4;\">"
	mailForm += "                    <p class=\"text\" style=\"font-size: 15.5px;line-height: 26px;color: #404040;text-align: center;padding-bottom: 8px;font-family: 'Hanna', sans-serif;\">"
	mailForm += "                        <strong>" + mailTile + "</strong>"
	mailForm += "                    </p>"
	mailForm += "                    <a href=\"" + resourceUrl + "\"> <div id=\"text_area\" style=\"margin:0 auto;padding:5px 12px;font-family:'Noto Sans', sans-serif;font-size:14px;color:#fff;vertical-align:middle;text-align:center;outline:none;display:block;border-radius:3px;box-shadow:none;border:none;background:#124379;font-weight:600;cursor:pointer;\">PaaS-Ta Monitor로 이동</div></a>"
	mailForm += "                    <div class=\"notic\" style=\"margin:0;padding:0;padding-top:35px;\">"
	mailForm += "                        <h3 style=\"margin:0;padding:0;font-family:'Noto Sans', sans-serif;font-size:15px;color:#555;font-weight:600;padding-bottom:6px;\">&#50508;&#47548; &#51221;&#48372;</h3>"
	mailForm += "                        <table style=\"text-align: left;border-top: 3px solid #888;border-bottom: 1px solid #bbb;width: 100%;border-spacing: 0;border-collapse: collapse;color: #404040;\">"
	mailForm += "                            <tr>"

	if mailContent.AlarmInfo.ServiceType == "SaaS" {
		mailForm += "                                <th style=\"background-color:#f5f5f5;padding: 4px;\">Application Name</th>"
	} else {
		mailForm += "                                <th style=\"background-color:#f5f5f5;padding: 4px;\">Pod Name</th>"
	}

	mailForm += "                                <th style=\"background-color:#f5f5f5;padding:4px;\">Measure</th>"
	mailForm += "                                <th style=\"background-color:#f5f5f5;padding:4px;\">Critical Value</th>"
	mailForm += "                                <th style=\"background-color:#f5f5f5;padding:4px;\">Measure Value</th>"
	mailForm += "                                <th style=\"background-color:#f5f5f5;padding:4px;\">Status</th>"
	mailForm += "                                <th style=\"background-color:#f5f5f5;padding:4px;\">Date</th>"
	mailForm += "                            </tr>"

	for _, mail := range mailContent.AlarmExecution {
		var criticalValue int
		measureValue := fmt.Sprintf("%.1f", mail.MeasureValue)
		if mail.CriticalStatus == "Warning" {
			criticalValue = mailContent.AlarmInfo.WarningValue
		} else {
			criticalValue = mailContent.AlarmInfo.CriticalValue
		}

		mailForm += "                            <tr>"
		mailForm += "                                <td style=\"padding: 8px 6px;\">" + mail.MeasureName1 + "</td>"
		mailForm += "                                <td style=\"padding: 8px 6px;\">" + mailContent.AlarmInfo.MetricType + "</td>"
		mailForm += "                                <td style=\"padding: 8px 6px;\">" + strconv.Itoa(criticalValue) + "%</td>"
		mailForm += "                                <td style=\"padding: 8px 6px;\">" + measureValue + "%</td>"
		mailForm += "                                <td style=\"padding: 8px 6px;\">" + mail.CriticalStatus + "</td>"
		mailForm += "                                <td style=\"padding: 8px 6px;\">" + mail.ExecutionTime.Format("2006-01-02 15:04:05") + "</td>"
		mailForm += "                            </tr>"
	}
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

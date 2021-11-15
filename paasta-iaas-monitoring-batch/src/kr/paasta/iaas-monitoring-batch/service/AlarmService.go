package service

import (
	"crypto/tls"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"kr/paasta/iaas-monitoring-batch/config"
	"kr/paasta/iaas-monitoring-batch/model/base"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"zabbix-client/common"
	"zabbix-client/host"
	"zabbix-client/hostgroup"
	"zabbix-client/item"

	"github.com/cavaliercoder/go-zabbix"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"kr/paasta/iaas-monitoring-batch/dao"
	"kr/paasta/iaas-monitoring-batch/model"
	"kr/paasta/iaas-monitoring-batch/util"
)

type AlarmService struct {
	ConfigData *config.Config
	Session *zabbix.Session
	DbConn *gorm.DB
	StopChan chan bool
	retry bool
}

const (
	INSERT_ALARM int = 0
	UPDATE_ALARM int = 1
)

var hostList []zabbix.Host

func AlarmServiceBuilder(configData *config.Config) *AlarmService {
	dbConnStr := util.GetConnectionString(configData.DbHost, configData.DbPort, configData.DbUser, configData.DbPasswd, configData.DbName)

	//log.Printf("dbConnStr : %v\n", dbConnStr+"?charset=utf8&parseTime=true")

	dbConn, err := gorm.Open(configData.DbType, dbConnStr+"?charset=utf8&parseTime=true")
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}

	createAlarmPolicyInitialData(dbConn)

	session := createSession(configData)
	stop := make(chan bool)
	return &AlarmService{
		ConfigData: configData,
		Session: session,
		DbConn: dbConn,
		StopChan: stop,
		retry: false,
	}
}

/**
	알람 정책 정보를 조회
		CPU 임계치, 메모리 임계치, 디스크 임계치
	호스트 정보를 조회
	호스트의 CPU, Memory, Disk 사용량을 조회
 */

func (service *AlarmService) RunScheduler() error {

	/*
	err := new(error)
	if err != nil {
		service.StopChan <- true
	}
	*/
	log.Println(" Running RunScheduler()...")
	service.UpdateSnsAlarmTarget()

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	go func(waitGroup *sync.WaitGroup) {
		//TODO 알람 정책과 리소스 사용량 대조하여 true면 알람 발송하는 함수 (&waitGroup)호출
		service.alarmChecker()
		waitGroup.Done()
	}(&waitGroup)

	waitGroup.Wait()
	return nil
}


/**
	알람 발송 여부를 체크하는 함수

 */
func (service *AlarmService) alarmChecker() {
	var err error
	err = nil
	if err != nil {
		fmt.Println(err)
		service.StopChan <- true
	}

	// 알람 정책 불러오기
	var cpuPolicyInfo model.AlarmPolicy
	var memoryPolicyInfo model.AlarmPolicy
	var diskPolicyInfo model.AlarmPolicy

	alarmPolicyList, _ := dao.GetAlarmPolicy(service.DbConn)
	for _, policy := range alarmPolicyList {
		if policy.AlarmType == base.ALARM_TYPE_CPU {
			cpuPolicyInfo = policy
		}
		if policy.AlarmType == base.ALARM_TYPE_MEMORY {
			memoryPolicyInfo = policy
		}
		if policy.AlarmType == base.ALARM_TYPE_DISK {
			diskPolicyInfo = policy
		}
	}

	var alarmThreshold base.AlarmThreshold
	alarmThreshold.OriginType = base.ORIGIN_TYPE_IAAS

	// Zabbix에서 데이터 불러오기
	hostList = getHosts(service.Session)

	hostIdArr := make([]string, len(hostList))
	for idx, host := range hostList {
		hostIdArr[idx] = host.HostID
	}

	cpuUsage := getCpuUsage(service.Session, hostIdArr)
	memoryUsage := getMemoryUsage(service.Session, hostIdArr)
	diskUsage := getDiskUsage(service.Session, hostIdArr)

	// CPU 사용량 체크
	for _, item := range cpuUsage {
		valueType := item.LastValueType
		var usage float64
		if valueType == 0 {
			lastValue, _ := strconv.ParseFloat(item.LastValue, 64)
			usage = lastValue
		}

		hostInfo := findHostData(item.HostID)

		alarmThreshold.AlarmType = base.ALARM_TYPE_CPU
		alarmThreshold.OriginId = uint(item.HostID)
		alarmThreshold.ServiceName = hostInfo.Hostname
		alarmThreshold.Ip = hostInfo.Interfaces[0]["ip"]
		alarmThreshold.Threshold = float64(cpuPolicyInfo.CriticalThreshold)
		alarmThreshold.RepeatTime = cpuPolicyInfo.RepeatTime

		// Warning에 대한 알람 처리
		log.Printf("HostID : %d, CPU policy : %d, usage: %f", item.HostID, cpuPolicyInfo.WarningThreshold, usage)
		if cpuPolicyInfo.WarningThreshold <= int(usage) {  // Warning CPU
			log.Println("CPU Warning")

			alarmThreshold.Level = base.ALARM_LEVEL_WARNING
			alarmThreshold.Threshold = float64(cpuPolicyInfo.WarningThreshold)
			alarmData := service.makeAlarmData(alarmThreshold, usage)
			isExist, existAlarm := dao.IsExistAlarm(alarmData, service.DbConn)
			if !isExist {
				service.alarmSender(alarmThreshold, alarmData, usage, INSERT_ALARM)
			} else {
				if existAlarm.ResolveStatus == "1" {
					alarmSendAvailableTime := existAlarm.AlarmSendDate.Add(time.Duration(cpuPolicyInfo.RepeatTime) * time.Minute)
					//DB에서 받아온 시간은 GMT TIme 으로 받아 오기 때문에 9 시간을 뺀다. (config.ini  : gmt.time.hour.gap)
					availTime := alarmSendAvailableTime.Add(time.Duration(service.ConfigData.GmtTimeGapHour) * time.Hour).Unix()
					if time.Now().Unix() >= availTime {
						service.alarmSender(alarmThreshold, alarmData, usage, UPDATE_ALARM)
					}

				}
			}
		}

		// Critical에 대한 알람 처리
		log.Printf("HostID : %d, CPU policy : %d, usage: %f", item.HostID, cpuPolicyInfo.CriticalThreshold, usage)
		if cpuPolicyInfo.CriticalThreshold <= int(usage) {  // Critical CPU
			log.Println("CPU Critical")

			alarmThreshold.Level = base.ALARM_LEVEL_CRITICAL
			alarmThreshold.Threshold = float64(cpuPolicyInfo.CriticalThreshold)
			alarmData := service.makeAlarmData(alarmThreshold, usage)
			isExist, existAlarm := dao.IsExistAlarm(alarmData, service.DbConn)
			if !isExist {
				service.alarmSender(alarmThreshold, alarmData, usage, INSERT_ALARM)
			} else {
				if existAlarm.ResolveStatus == "1" {
					service.alarmSender(alarmThreshold, alarmData, usage, UPDATE_ALARM)
				}
			}
		}
	}

	// 메모리 사용량 체크
	for _, item := range memoryUsage {

		valueType := item.LastValueType

		var usage float64
		if valueType == 0 {
			lastValue, _ := strconv.ParseFloat(item.LastValue, 64)
			usage = lastValue
		}

		hostInfo := findHostData(item.HostID)

		alarmThreshold.AlarmType = base.ALARM_TYPE_MEMORY
		alarmThreshold.OriginId = uint(item.HostID)
		alarmThreshold.ServiceName = hostInfo.Hostname
		alarmThreshold.Ip = hostInfo.Interfaces[0]["ip"]
		alarmThreshold.Threshold = float64(memoryPolicyInfo.CriticalThreshold)
		alarmThreshold.RepeatTime = memoryPolicyInfo.RepeatTime

		log.Printf("HostID : %d, Memory policy : %d, usage: %f", item.HostID, memoryPolicyInfo.WarningThreshold, usage)
		if memoryPolicyInfo.WarningThreshold <= int(usage) {  // Warning Memory
			// TODO
			log.Println("Memory warning")

			alarmThreshold.Level = base.ALARM_LEVEL_WARNING
			alarmData := service.makeAlarmData(alarmThreshold, usage)
			isExist, existAlarm := dao.IsExistAlarm(alarmData, service.DbConn)
			if !isExist {
				service.alarmSender(alarmThreshold, alarmData, usage, INSERT_ALARM)
			} else {
				if existAlarm.ResolveStatus == "1" {
					service.alarmSender(alarmThreshold, alarmData, usage, UPDATE_ALARM)
				}
			}
		}

		log.Printf("HostID : %d, Memory policy : %d, usage: %f", item.HostID, memoryPolicyInfo.CriticalThreshold, usage)
		if memoryPolicyInfo.CriticalThreshold <= int(usage) {  // Critical Memory
			// TODO
			log.Println("Memory Critical")

			alarmThreshold.Level = base.ALARM_LEVEL_CRITICAL
			alarmThreshold.Threshold = float64(memoryPolicyInfo.CriticalThreshold)

			alarmData := service.makeAlarmData(alarmThreshold, usage)
			isExist, existAlarm := dao.IsExistAlarm(alarmData, service.DbConn)
			if !isExist {
				service.alarmSender(alarmThreshold, alarmData, usage, INSERT_ALARM)
			} else {
				if existAlarm.ResolveStatus == "1" {
					service.alarmSender(alarmThreshold, alarmData, usage, UPDATE_ALARM)
				}
			}
		}
	}

	// TODO 디스크 사용량 체크
	for _, item := range diskUsage {
		valueType := item.LastValueType

		var usage float64
		if valueType == 0 {
			lastValue, _ := strconv.ParseFloat(item.LastValue, 64)
			usage = lastValue
		}
		hostInfo := findHostData(item.HostID)

		alarmThreshold.AlarmType = base.ALARM_TYPE_DISK
		alarmThreshold.OriginId = uint(item.HostID)
		alarmThreshold.ServiceName = hostInfo.Hostname
		alarmThreshold.Ip = hostInfo.Interfaces[0]["ip"]
		alarmThreshold.RepeatTime = diskPolicyInfo.RepeatTime

		log.Printf("HostID : %d, Disk policy : %d, usage: %f", item.HostID, diskPolicyInfo.WarningThreshold, usage)
		if diskPolicyInfo.WarningThreshold <= int(usage) {  // Warning Memory
			// TODO
			log.Println("Disk warning")

			alarmThreshold.Level = base.ALARM_LEVEL_WARNING
			alarmThreshold.Threshold = float64(diskPolicyInfo.WarningThreshold)

			alarmData := service.makeAlarmData(alarmThreshold, usage)
			isExist, existAlarm := dao.IsExistAlarm(alarmData, service.DbConn)
			if !isExist {
				service.alarmSender(alarmThreshold, alarmData, usage, INSERT_ALARM)
			} else {
				if existAlarm.ResolveStatus == "1" {
					service.alarmSender(alarmThreshold, alarmData, usage, UPDATE_ALARM)
				}
			}
		}

		log.Printf("HostID : %d, Disk policy : %d, usage: %f", item.HostID, diskPolicyInfo.CriticalThreshold, usage)
		if diskPolicyInfo.CriticalThreshold <= int(usage) {  // Critical Memory
			// TODO
			log.Println("Disk Critical")

			alarmThreshold.Level = base.ALARM_LEVEL_CRITICAL
			alarmThreshold.Threshold = float64(diskPolicyInfo.CriticalThreshold)

			alarmData := service.makeAlarmData(alarmThreshold, usage)
			isExist, existAlarm := dao.IsExistAlarm(alarmData, service.DbConn)
			if !isExist {
				service.alarmSender(alarmThreshold, alarmData, usage, INSERT_ALARM)
			} else {
				if existAlarm.ResolveStatus == "1" {
					service.alarmSender(alarmThreshold, alarmData, usage, UPDATE_ALARM)
				}
			}
		}
	}
}

/**
	알람 발송을 담당하는 함수
	MODE.INSERT_ALARM = 0
	MODE.UPDATE_ALARM = 1
 */
func (service *AlarmService) alarmSender(alarmThreshold base.AlarmThreshold, alarm model.Alarm, usage float64, mode int) {
	if mode == INSERT_ALARM {
		dao.InsertAlarm(alarm, service.DbConn)
	} else if mode == UPDATE_ALARM {
		dao.UpdateAlarm(alarm, service.DbConn)
	}

	alarmData := dao.GetAlarm(alarm, service.DbConn)

	date := alarmData.RegDate.Add(time.Duration(9) * time.Hour).Format("2006-01-02 15:04:05")
	elaspeTime := time.Now().Unix() - alarmData.RegDate.Add(time.Duration(9) * time.Hour).Unix()
	elaspeMinute := strconv.FormatInt(elaspeTime / 60, 10)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		service.SendMail(alarmData, alarmThreshold, usage)
	}()
	go func() {
		defer wg.Done()
		service.SendTelegram(alarm, date, elaspeMinute, alarmThreshold, service.DbConn)
	}()
	wg.Wait()
}


func createSession(configData *config.Config) *zabbix.Session {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	cache := zabbix.NewSessionFileCache().SetFilePath("./zabbix_session")
	_session, err := zabbix.CreateClient(configData.ZabbixHost).
		WithCache(cache).
		WithHTTPClient(client).
		WithCredentials(configData.ZabbixAdminId, configData.ZabbixAdminPw).Connect()
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	version, err := _session.GetVersion()
	fmt.Printf("Connected to Zabbix API v%s\n", version)

	return _session
}


func getHosts(session *zabbix.Session) []zabbix.Host {
	groupParams := make(map[string]interface{}, 0)
	groupParams["name"] = "PaaS-TA Group"   // TODO : 차후 fix될 호스트그룹명을 명시할 것.. 'PaaS-TA Group'
	hostgroupList := hostgroup.GetHostgroup(session, groupParams)

	groupId := hostgroupList[0].GroupID

	hostParams := make(map[string]interface{}, 0)
	groupIds := make([]string, 1)
	groupIds[0] = groupId
	hostParams["groupIds"] = groupIds
	hostList := host.GetHostList(session, hostParams)

	return hostList
}


/**
	CPU 사용량 조회
 */
func getCpuUsage(session *zabbix.Session, hostIdArr []string) []zabbix.Item {
	itemParams := make(map[string]interface{}, 0)
	keywordArr := make([]string, 1)
	keywordArr[0] = common.SYSTEM_CPU_UTIL
	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIdArr
	return item.GetItemList(session, itemParams)
}


/**
	메모리 사용량 조회
 */
func getMemoryUsage(session *zabbix.Session, hostIdArr []string) []zabbix.Item {
	itemParams := make(map[string]interface{}, 0)
	keywordArr := make([]string, 1)
	keywordArr[0] = common.VM_MEMORY_UTILIZATION
	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIdArr
	return item.GetItemList(session, itemParams)
}


/**
	디스크 사용량 조회
 */
func getDiskUsage(session *zabbix.Session, hostIdArr []string) []zabbix.Item {
	itemParams := make(map[string]interface{}, 0)
	keywordArr := make([]string, 1)
	keywordArr[0] = common.SPACE_UTILIZATION
	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIdArr
	return item.GetItemList(session, itemParams)
}


func findHostData(hostId int) zabbix.Host {
	var result zabbix.Host
	for _, host := range hostList {
		if host.HostID == strconv.Itoa(hostId) {
			result = host
			break
		}
	}
	return result
}


func (service *AlarmService) StopProcess() {
	go func() {
		for {
			select {
			case <- service.StopChan:
				os.Exit(1) //Bosh Monit start Batch program automatically if the process is down.
			}
		}
	}()
}


func createAlarmPolicyInitialData(dbClient *gorm.DB) {
	log.Println("createAlarmPolicyInitialData")
	dbClient.AutoMigrate(&model.AlarmPolicy{}, &model.AlarmTarget{}, &model.AlarmSns{}, &model.AlarmSnsTarget{})

	iaasCpuData  := model.AlarmPolicy{Id:10, OriginType: "ias", AlarmType: "cpu", WarningThreshold: 85, CriticalThreshold: 90, RepeatTime: 10 , MeasureTime: 600 , Comment: "Initial Data"}
	iaasMemData  := model.AlarmPolicy{Id:11, OriginType: "ias", AlarmType: "memory", WarningThreshold: 85, CriticalThreshold: 90, RepeatTime: 10 , MeasureTime: 600 , Comment: "Initial Data"}
	iaasDiskData := model.AlarmPolicy{Id:12, OriginType: "ias", AlarmType: "disk", WarningThreshold: 85, CriticalThreshold: 90, RepeatTime: 10 , MeasureTime: 600 , Comment: "Initial Data"}
	alarmTaget   := model.AlarmTarget{Id:4, OriginType: "ias", MailAddress: "adminUser@gmail.com", MailSendYn: "Y" }

	dbClient.FirstOrCreate(&iaasCpuData)
	dbClient.FirstOrCreate(&iaasMemData)
	dbClient.FirstOrCreate(&iaasDiskData)
	dbClient.FirstOrCreate(&alarmTaget)
}


func (service AlarmService) SendMail(alarm model.Alarm, threshold base.AlarmThreshold, usage float64) {

	// alarm_target 테이블에서 알람 이메일을 받을 이메일 주소를 조회함
	mailReceiver := dao.GetAlarmTarget(service.DbConn)
	var mailReceivers []string
	for _, data := range mailReceiver {
		mailReceivers = append(mailReceivers, data.MailAddress)
	}

	mail := model.Mail{}
	mail.Sender = service.ConfigData.MailSender
	mail.To = mailReceivers
	mail.Subject = alarm.AlarmTitle

	// 이메일 본문 구성하기
	messageBody := service.makeMailHTMLContents(alarm, threshold, usage)
	mail.Body = messageBody

	smtpServer := model.SmtpServer{Host: service.ConfigData.SmtpHost, Port: service.ConfigData.Port}

	var client *smtp.Client
	var err error
	var conn *tls.Conn

	if service.ConfigData.MailTlsSend == true {
		fmt.Println(">>>>>>>>>>>>>>> Sendmail STMP TLS Mode")
		smtpServer.TlsConfig = &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         smtpServer.Host,
		}
		auth := smtp.PlainAuth("", mail.Sender, service.ConfigData.SenderPwd, smtpServer.Host)

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
			fmt.Println("smtp connection & client create error :", err.Error())
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

		if service.ConfigData.MailTlsSend == true {
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


func (f AlarmService) SendTelegram(alarmData model.Alarm, date string, elapseMinute string, threshold base.AlarmThreshold, dbConn *gorm.DB) {
	alarmSns, err := dao.GetAlarmSns(dbConn)
	if len(alarmSns) < 1 {
		return
	}
	if err != nil {
		fmt.Println("Failed to get SNS alarm target! :", err)
		return
	} else {

		text := fmt.Sprintf("status: %s\nname: %s\ndate: %s\nelapsed: %s\n%s", alarmData.Level, threshold.ServiceName, date, elapseMinute, alarmData.AlarmMessage)
		for _, sns := range alarmSns {
			if strings.ToUpper(sns.SnsSendYn) == "Y" {
				alarmSnsTarget, err := dao.GetAlarmSnsTarget(sns, dbConn)
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

func (f AlarmService) makeAlarmData(alarmSource base.AlarmThreshold, currentSystemUsage float64) model.Alarm {
	var alarmData model.Alarm
	alarmData.OriginType = alarmSource.OriginType
	alarmData.OriginId = alarmSource.OriginId
	alarmData.AlarmType = alarmSource.AlarmType
	alarmData.Level = alarmSource.Level
	alarmData.Ip = alarmSource.Ip
	alarmData.ResolveStatus = "1"

	var title string
	var message string

	alarmData.AppYn = "N"
	switch alarmSource.Level {
		case base.ALARM_LEVEL_WARNING:
			fallthrough
		case base.ALARM_LEVEL_CRITICAL:
			title = "[" + alarmSource.ServiceName + "]의 " + alarmSource.AlarmType + " 상태 [" + alarmSource.Level + "]"
			message = alarmSource.ServiceName + " " + alarmSource.AlarmType + " 의 상태" + alarmSource.Level + "\n"
			message += alarmSource.AlarmType + " 의 임계치인 " + util.Floattostrwithprec(alarmSource.Threshold, 0) + "%를 초과하였습니다."
			message += " \n 현재 사용률 [" + util.FloattostrDigit2(currentSystemUsage) + "]% 입니다. "
		case base.ALARM_LEVEL_FAIL:
			title = "[" + alarmSource.ServiceName + "]가 다운되었습니다."
			message = "[" + alarmSource.ServiceName + "]가 다운되었습니다."
	}
	alarmData.AlarmTitle = title
	alarmData.AlarmMessage = message

	return alarmData
}

func (service *AlarmService) makeMailHTMLContents(alarm model.Alarm, alarmThreshold base.AlarmThreshold, usage float64) string {

	alarmOccurrenceDate := alarm.RegDate.Add(time.Duration(9) * time.Hour).Format("2006-01-02 15:04:05")
	elaspeTime := time.Now().Unix() - alarm.RegDate.Unix()
	elapsedMinute := strconv.FormatInt(elaspeTime / 60, 10)

	mailForm := ""
	mailForm = mailForm + "Subject: " + alarm.AlarmTitle + "!\n"

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
	mailForm += "#text_area p strong:before{display:inline-block;content:'';width:23px;height:24px;background:url(" + service.ConfigData.ResouceUrl + "/public/resources/img/email/ic_error.png) no-repeat;background-size:100%;vertical-align:bottom;margin-right:7px;}"
	mailForm += "</style>"
	mailForm += "</head>"
	mailForm += "<body style=\"*word-break:break-all;-ms-word-break:break-all;margin:0;padding:0;font-family:'Noto Sans', sans-serif;font-size:12px;color:#555;\">"
	mailForm += "<div id=\"wrap\" style=\"margin:0;padding:0;width:100%;height:100%;background-size:cover;overflow-y:auto;\">"
	mailForm += "    <div class=\"email_form\" style=\"width: 708px;height: 361px;padding: 19px 30px;overflow: hidden;margin: 0 auto;display: block;background: #f8f8f8;\"><div class=\"form_area\" style=\"margin:0;padding:0\">"
	mailForm += "        <div class=\"contents_area\" style=\"margin:0;padding:0\">"
	mailForm += "            <div class=\"contents-header\" style=\"margin:0;padding:0;padding-bottom:25px;\">"
	mailForm += "                <h2 style=\"margin:0;padding:0;font-family:'Noto Sans', sans-serif;font-size:12px;color:#555;\"><img src=\"" + service.ConfigData.ResouceUrl + "/public/resources/img/email/logo.png\" alt=\"monitorxpert\" style=\"border:0 none;width:279px;\"></h2>"
	mailForm += "            </div>"
	mailForm += "            <div class=\"contents-body\" style=\"margin:0;padding:0;position:relative;\">"
	mailForm += "                <div id=\"text_area\" style=\"margin:0;padding:18px 18px 24px;padding-bottom:12px;border:1px solid #ddd;background:#fff;height:auto;border-top:6px solid #007dc4;\">"
	mailForm += "                    <p class=\"text\" style=\"font-size: 15.5px;line-height: 26px;color: #404040;text-align: center;padding-bottom: 8px;font-family: 'Hanna', sans-serif;\">"
	mailForm += "                        <strong>" + alarm.Level + " : </strong>"
	mailForm += "                        " + alarm.AlarmMessage
	mailForm += "                    </p>"
	mailForm += "                    <a href=\"" + service.ConfigData.ResouceUrl + "\"> <div id=\"text_area\" style=\"margin:0 auto;padding:5px 12px;font-family:'Noto Sans', sans-serif;font-size:14px;color:#fff;vertical-align:middle;text-align:center;outline:none;display:block;border-radius:3px;box-shadow:none;border:none;background:#124379;font-weight:600;cursor:pointer;\">PaaS-Ta Monitor로 이동</div></a>"
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
	mailForm += "                                <td style=\"padding: 8px 6px;\">" + alarm.Level + "</td>"
	mailForm += "                                <td style=\"padding: 8px 6px;\">" + alarmThreshold.ServiceName + "</td>"
	mailForm += "                                <td style=\"padding: 8px 6px;\">" + alarmOccurrenceDate + "</td>"
	mailForm += "                                <td style=\"padding: 8px 6px;\">" + elapsedMinute + "</td>"
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


func (service *AlarmService)UpdateSnsAlarmTarget() {

	alarmSns, err := dao.GetAlarmSns(service.DbConn)
	if err != nil {
		fmt.Println("Failed to get sns_id(ChatRoomId)! :", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(alarmSns))
	for _, v := range alarmSns {
		go func(wg *sync.WaitGroup, v model.AlarmSns) {
			defer wg.Done()

			if v.SnsType == base.SNS_TYPE_TELEGRAM {
				bot, err := tgbotapi.NewBotAPI(v.Token)
				if err != nil {
					fmt.Println(err)
				} else {
					bot.Debug = true
					var updateConfig tgbotapi.UpdateConfig
					updateConfig.Offset = 0
					updateConfig.Timeout = 30
					updates, err := bot.GetUpdates(updateConfig)
					if err != nil {
						fmt.Println(err)
					} else {
						var chatIdList []int64
						for _, update := range updates {
							if update.Message == nil {
								continue
							}
							chatIdList = append(chatIdList, update.Message.Chat.ID)
						}
						chatIdList = util.RemoveDuplicates(chatIdList)
						for _, chatId := range chatIdList {
							var alarmSnsTarget model.AlarmSnsTarget
							alarmSnsTarget.ChannelId = v.ChannelId
							alarmSnsTarget.TargetId = chatId
							err := dao.UpdateSnsAlarmTargets(alarmSnsTarget, service.DbConn)
							if err != nil {
								return
							}
						}
					}
				}
			}
		}(&wg, v)
	}
	wg.Wait()
}
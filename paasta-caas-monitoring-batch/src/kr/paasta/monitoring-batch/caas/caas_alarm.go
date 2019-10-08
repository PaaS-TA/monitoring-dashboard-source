package caas

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
	"github.com/mileusna/crontab"
	"github.com/thoas/go-funk"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"kr/paasta/monitoring-batch/dao"
	"kr/paasta/monitoring-batch/model"
	"kr/paasta/monitoring-batch/notify"
	"kr/paasta/monitoring-batch/util"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

var equalCheckAlarmInfos []model.BatchAlarmInfo
var procAlarmInfos []model.BatchAlarmInfoCheck

//var ctab *crontab.Crontab

func Startschedule(dbClient *gorm.DB, monitoringUrl string) {
	ctab := crontab.New() // create cron table
	// 수신자 정보 조회 및  저장
	ctab.MustAddJob("*/3 * * * *", jobCaasAlarmSnsTarget, dbClient)
	ctab.MustAddJob("*/1 * * * *", getAlarmInfos, dbClient, monitoringUrl)
	ctab.RunAll()
}

func jobCaasAlarmSnsTarget(dbClient *gorm.DB) {
	alarmSnsToken := dao.GetBatchAlarmSnsToken("CaaS", dbClient)
	if alarmSnsToken == (model.BatchAlarmSns{}) {
		return
	}

	bot, err := tgbotapi.NewBotAPI(alarmSnsToken.Token)

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
			var targerIds []string
			for _, update := range updates {
				if update.Message == nil {
					continue
				}
				if !util.ContainsTargetId(chatIdList, update.Message.Chat.ID) {
					targetId := strconv.Itoa(int(update.Message.Chat.ID))
					targerIds = append(targerIds, targetId)
				}
			}
			dao.SaveBatchAlarmSnsReceiver("CaaS", dbClient, targerIds)
		}
	}
}

func jobCaasAlarm(dbClient *gorm.DB, alarmInfo model.BatchAlarmInfoCheck, monitoringUrl string) {
	//podList, data := podList(monitoringUrl)
	podList := podList(monitoringUrl)

	var waitGroup sync.WaitGroup
	var runAppCnt = 0
	resutChan := make(chan model.BatchAlarmExecution, len(podList))

	for _, data := range podList {
		waitGroup.Add(1)
		runAppCnt++
		go podRunningStat(resutChan, &waitGroup, dbClient, data, alarmInfo)
	}

	var batchExecutions []model.BatchAlarmExecution

	for k := 0; k < runAppCnt; k++ {
		batchTempExecutions := <-resutChan
		if batchTempExecutions.MeasureValue > float64(alarmInfo.CriticalValue) || batchTempExecutions.MeasureValue > float64(alarmInfo.WarningValue) {
			batchExecutions = append(batchExecutions, batchTempExecutions)
		}
	}
	close(resutChan)

	if len(batchExecutions) > 0 {
		emails, _ := dao.GetBatchAlarmReceiver(alarmInfo.ServiceType, dbClient)
		go notify.SendMail(alarmInfo.ServiceType, emails, model.MailContent{alarmInfo, batchExecutions})
	}
}

//func podList(monitoringUrl string) ([]map[string]string, []byte) {
func podList(monitoringUrl string) []map[string]gjson.Result {
	url := monitoringUrl + "/v2/caas/monitoring/podList"

	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	str2 := string(data)

	jsonString1 := gjson.Get(str2, "#")

	jsonMap := make([]map[string]gjson.Result, 0)

	for i := 0; i < int(jsonString1.Int()); i++ {
		jsonData := gjson.Get(str2, ""+strconv.Itoa(i)+"")
		jsonDataMap := jsonData.Map()

		jsonMap = append(jsonMap, jsonDataMap)
	}

	return jsonMap
}

func podRunningStat(c chan model.BatchAlarmExecution, waitGroup *sync.WaitGroup, dbClient *gorm.DB, result map[string]gjson.Result, alarmInfo model.BatchAlarmInfoCheck) {
	defer waitGroup.Done()

	fmt.Printf("result : %v\n", result)

	var value float64
	warning := float64(alarmInfo.WarningValue)
	critical := float64(alarmInfo.CriticalValue)
	podName := result["PodName"].String()
	nameSpace := result["NameSpace"].String()

	if alarmInfo.MetricType == "CPU" {
		value = result["CpuUsage"].Float()

	}

	if alarmInfo.MetricType == "MEMORY" {
		value = result["MemoryUsage"].Float()
	}

	if alarmInfo.MetricType == "DISK" {
		value = result["DiskUsage"].Float()
	}

	executionResult := strings.Replace(alarmInfo.ExecMsg, "${PodName}", podName, 1)

	criticalStatus := "Success"
	if value > critical {
		criticalStatus = "Critical"
		executionResult = "[critical] " + strings.Replace(executionResult, "${Currend_value}", fmt.Sprintf("%.2f", value), 1) + ", 임계값 (" + strconv.Itoa(alarmInfo.CriticalValue) + ")%"
	} else if value > warning {
		criticalStatus = "Warning"
		executionResult = "[warning] " + strings.Replace(executionResult, "${Currend_value}", fmt.Sprintf("%.2f", value), 1) + ", 임계값 (" + strconv.Itoa(alarmInfo.WarningValue) + ")%"
	} else {
		executionResult = strings.Replace(executionResult, "${Currend_value}", fmt.Sprintf("%.1f", value), 1)
	}

	batchExecution := model.BatchAlarmExecution{
		CriticalStatus:  criticalStatus,
		AlarmId:         alarmInfo.AlarmId,
		ServiceType:     alarmInfo.ServiceType,
		MeasureValue:    value,
		MeasureName1:    podName,
		MeasureName2:    nameSpace,
		MeasureName3:    "",
		ExecutionTime:   time.Now().Local(),
		ExecutionResult: executionResult,
	}

	c <- batchExecution
	podStatReport(dbClient, alarmInfo, batchExecution, value)
}

// CaaS 스케쥴 DB 저장 및 임계지 이상 값 알림 처리
func podStatReport(dbClient *gorm.DB, alarmInfo model.BatchAlarmInfoCheck, batchExecution model.BatchAlarmExecution, result float64) {
	var token string = ""
	alarmSns := dao.GetBatchAlarmSnsToken(alarmInfo.ServiceType, dbClient)
	if alarmSns != (model.BatchAlarmSns{}) {
		token = alarmSns.Token
	}

	if result > float64(alarmInfo.CriticalValue) || result > float64(alarmInfo.WarningValue) {
		go dao.InsertBatchExecution(dbClient, &batchExecution)
		_, snsIds := dao.GetBatchAlarmReceiver(alarmInfo.ServiceType, dbClient)
		go notify.SendChatBot(snsIds, batchExecution.ExecutionResult, token)
	}
}
func getAlarmInfos(dbClient *gorm.DB, monitoringUrl string) {
	alarmInfos := dao.GetBatchAlarmInfo(dbClient)

	var isEqual bool = false
	if len(alarmInfos) > 0 {
		caasAlarms := funk.Filter(alarmInfos, func(x model.BatchAlarmInfo) bool {
			return x.ServiceType == "CaaS"
		}).([]model.BatchAlarmInfo)

		if len(equalCheckAlarmInfos) == 0 {
			isEqual = false
		} else {
			if len(caasAlarms) == len(equalCheckAlarmInfos) {
				isEqual = reflect.DeepEqual(caasAlarms, equalCheckAlarmInfos)
			} else {
				isEqual = false
			}
		}

		if isEqual == false {
			copyAlarmInfo(caasAlarms, &equalCheckAlarmInfos, &procAlarmInfos)
			for i := 0; i < len(procAlarmInfos); i++ {
				caasAlarm := &procAlarmInfos[i]
				caasAlarm.NextMinute = nextScheduleMinute(*caasAlarm)

				go jobCaasAlarm(dbClient, *caasAlarm, monitoringUrl)
			}
		} else {
			mninute := strconv.Itoa(time.Now().Minute())
			for i := 0; i < len(procAlarmInfos); i++ {
				caasAlarm := &procAlarmInfos[i]
				if mninute == caasAlarm.NextMinute {
					caasAlarm.NextMinute = nextScheduleMinute(*caasAlarm)
					go jobCaasAlarm(dbClient, *caasAlarm, monitoringUrl)
				}
			}
		}
	}

}

func copyAlarmInfo(dbAlarmInfos []model.BatchAlarmInfo, equalCheckAlarmInfos *[]model.BatchAlarmInfo, procAlarmInfos *[]model.BatchAlarmInfoCheck) {
	*equalCheckAlarmInfos = make([]model.BatchAlarmInfo, len(dbAlarmInfos), len(dbAlarmInfos))
	*procAlarmInfos = make([]model.BatchAlarmInfoCheck, len(dbAlarmInfos), len(dbAlarmInfos))

	for i, alarmInfo := range dbAlarmInfos {
		(*equalCheckAlarmInfos)[i].AlarmId = alarmInfo.AlarmId
		(*equalCheckAlarmInfos)[i].ServiceType = alarmInfo.ServiceType
		(*equalCheckAlarmInfos)[i].MetricType = alarmInfo.MetricType
		(*equalCheckAlarmInfos)[i].WarningValue = alarmInfo.WarningValue
		(*equalCheckAlarmInfos)[i].CriticalValue = alarmInfo.CriticalValue
		(*equalCheckAlarmInfos)[i].MeasureTime = alarmInfo.MeasureTime
		(*equalCheckAlarmInfos)[i].CronExpression = alarmInfo.CronExpression
		(*equalCheckAlarmInfos)[i].MeasureTime = alarmInfo.MeasureTime
		(*equalCheckAlarmInfos)[i].ExecMsg = alarmInfo.ExecMsg
		(*equalCheckAlarmInfos)[i].ParamData1 = alarmInfo.ParamData1
		(*equalCheckAlarmInfos)[i].ParamData2 = alarmInfo.ParamData2
		(*equalCheckAlarmInfos)[i].ParamData3 = alarmInfo.ParamData3

		(*procAlarmInfos)[i].AlarmId = alarmInfo.AlarmId
		(*procAlarmInfos)[i].ServiceType = alarmInfo.ServiceType
		(*procAlarmInfos)[i].MetricType = alarmInfo.MetricType
		(*procAlarmInfos)[i].WarningValue = alarmInfo.WarningValue
		(*procAlarmInfos)[i].CriticalValue = alarmInfo.CriticalValue
		(*procAlarmInfos)[i].MeasureTime = alarmInfo.MeasureTime
		(*procAlarmInfos)[i].CronExpression = alarmInfo.CronExpression
		(*procAlarmInfos)[i].MeasureTime = alarmInfo.MeasureTime
		(*procAlarmInfos)[i].ExecMsg = alarmInfo.ExecMsg
		(*procAlarmInfos)[i].ParamData1 = alarmInfo.ParamData1
		(*procAlarmInfos)[i].ParamData2 = alarmInfo.ParamData2
		(*procAlarmInfos)[i].ParamData3 = alarmInfo.ParamData3
	}
}

func nextScheduleMinute(alarmInfo model.BatchAlarmInfoCheck) string {
	alarmInfo.CronExpression = strings.Replace(alarmInfo.CronExpression, "/", "", -1)
	alarmInfo.CronExpression = strings.Replace(alarmInfo.CronExpression, " ", "", -1)
	alarmInfo.CronExpression = strings.Replace(alarmInfo.CronExpression, "*", "", -1)

	minute, _ := strconv.Atoi(alarmInfo.CronExpression)

	dulation := time.Duration(minute) * time.Minute
	now := time.Now().Add(dulation)
	return strconv.Itoa(now.Minute())
}

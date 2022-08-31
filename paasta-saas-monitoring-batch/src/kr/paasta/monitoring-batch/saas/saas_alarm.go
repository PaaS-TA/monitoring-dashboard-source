package saas

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
	"github.com/mileusna/crontab"
	"github.com/thoas/go-funk"
	"github.com/tidwall/gjson"

	"saas-monitoring-batch/dao"
	"saas-monitoring-batch/model"
	"saas-monitoring-batch/notify"
	"saas-monitoring-batch/util"
)

var equalCheckAlarmInfos []model.BatchAlarmInfo
var procAlarmInfos []model.BatchAlarmInfoCheck

func Startschedule(dbClient *gorm.DB, pinpointUrl string) {
	ctab := crontab.New() // create cron table

	ctab.MustAddJob("*/3 * * * *", jobSaasAlarmSnsTarget, dbClient)
	ctab.MustAddJob("*/1 * * * *", getAlarmInfos, dbClient, pinpointUrl)
	ctab.RunAll()
}

func jobSaasAlarmSnsTarget(dbClient *gorm.DB) {
	alarmSnsToken := dao.GetBatchAlarmSnsToken("SaaS", dbClient)

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
			dao.SaveBatchAlarmSnsReceiver("SaaS", dbClient, targerIds)
		}
	}
}

func jobSaasAlarm(dbClient *gorm.DB, alarmInfo model.BatchAlarmInfoCheck, pinpointUrl string) {
	applications, data := appNameList(pinpointUrl)
	var waitGroup sync.WaitGroup
	var runAppCnt = 0
	resutChan := make(chan model.BatchAlarmExecution, len(applications))

	for appName, _ := range applications {
		jpath := appName + ".#.status.state.code"
		json := gjson.Get(string(data), jpath)
		appStatus := json.Array()[0].Num

		switch appStatus {
		case 100:
			runAppCnt++
			agentId := gjson.Get(string(data), appName+".#.agentId").Array()[0]
			waitGroup.Add(1)
			go appRunningStat(resutChan, &waitGroup, dbClient, appName, agentId.String(), pinpointUrl, alarmInfo)
			break
			//case 200 :
			//	appNotRunningStat(dbClient, appName, alarmInfo, "Shutdown")
			//	break
			//case 201:
			//	appNotRunningStat(dbClient, appName, alarmInfo, "Shutdown")
			//	break
			//case 300:
			//	appNotRunningStat(dbClient, appName, alarmInfo, "Disconnect")
			//	break
		}
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

func appNameList(pinpointUrl string) (map[string]string, []byte) {
	url := pinpointUrl + "/getAgentList.pinpoint"

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

	var applications map[string]string

	json.Unmarshal(data, &applications)

	return applications, data
}

func appRunningStat(c chan model.BatchAlarmExecution, waitGroup *sync.WaitGroup, dbClient *gorm.DB, appName string, agentId string, pinpointUrl string, alarmInfo model.BatchAlarmInfoCheck) {
	defer waitGroup.Done()

	from := strconv.FormatInt(time.Now().Add(time.Duration(alarmInfo.MeasureTime*-1)*time.Second).UTC().Unix(), 10) + "000"
	to := strconv.FormatInt(time.Now().UTC().Unix(), 10) + "000"

	data, err := getRestCall(pinpointUrl + alarmInfo.ParamData1 + "?agentId=" + agentId + "&from=" + from + "&to=" + to + "")
	if err != nil {
		return
	}
	var result1 float64
	var result2 float64 = 1
	var result float64
	result1 = getAvgReustData(data, alarmInfo.ParamData2)
	//result = int(result1)
	if alarmInfo.MetricType == "HEAP_MEMORY" {
		result2 = getAvgReustData(data, alarmInfo.ParamData3)
		result = result1 / result2 * 100
	} else {
		result = result1
	}

	executionResult := strings.Replace(alarmInfo.ExecMsg, "${AppName}", appName, 1)

	criticalStatus := "Success"

	if result > float64(alarmInfo.CriticalValue) {
		criticalStatus = "Critical"
		executionResult = "[critical] " + strings.Replace(executionResult, "${Currend_value}", fmt.Sprintf("%.1f", result), 1) + ", 임계값 (" + strconv.Itoa(alarmInfo.CriticalValue) + ")%"
	} else if result > float64(alarmInfo.WarningValue) {
		criticalStatus = "Warning"
		executionResult = "[warning] " + strings.Replace(executionResult, "${Currend_value}", fmt.Sprintf("%.1f", result), 1) + ", 임계값 (" + strconv.Itoa(alarmInfo.WarningValue) + ")%"
	} else {
		executionResult = strings.Replace(executionResult, "${Currend_value}", fmt.Sprintf("%.1f", result), 1)
	}

	batchExecution := model.BatchAlarmExecution{
		CriticalStatus:  criticalStatus,
		AlarmId:         alarmInfo.AlarmId,
		ServiceType:     alarmInfo.ServiceType,
		MeasureValue:    result,
		MeasureName1:    appName,
		MeasureName2:    "",
		MeasureName3:    "",
		ExecutionTime:   time.Now().Local(),
		ExecutionResult: executionResult,
	}

	c <- batchExecution
	appStatReport(dbClient, alarmInfo, batchExecution, result)
}

// SaaS App Shutdown / Disconnect 처리
func appNotRunningStat(dbClient *gorm.DB, waitGroup *sync.WaitGroup, appName string, alarmInfo model.BatchAlarmInfoCheck, criticalStatus string) {
	defer waitGroup.Done()

	executionResult := "SaaS " + appName + " app is " + criticalStatus

	batchExecution := model.BatchAlarmExecution{
		CriticalStatus:  criticalStatus,
		AlarmId:         alarmInfo.AlarmId,
		ServiceType:     alarmInfo.ServiceType,
		MeasureValue:    0,
		MeasureName1:    appName,
		MeasureName2:    "",
		MeasureName3:    "",
		ExecutionTime:   time.Now().Local(),
		ExecutionResult: executionResult,
	}

	appStatReport(dbClient, alarmInfo, batchExecution, 1000)
}

// SaaS 스케쥴 DB 저장 및 임계지 이상 값 알림 처리
func appStatReport(dbClient *gorm.DB, alarmInfo model.BatchAlarmInfoCheck, batchExecution model.BatchAlarmExecution, result float64) {
	var token string = ""
	alarmSns := dao.GetBatchAlarmSnsToken(alarmInfo.ServiceType, dbClient)
	if alarmSns != (model.BatchAlarmSns{}) {
		token = alarmSns.Token
	}

	if result > float64(alarmInfo.CriticalValue) || result > float64(alarmInfo.WarningValue) {
		go dao.InsertBatchExecution(dbClient, &batchExecution)
		_, snsIds := dao.GetBatchAlarmReceiver(alarmInfo.ServiceType, dbClient)
		go notify.SendChatBot(alarmInfo.ServiceType, snsIds, batchExecution.ExecutionResult, token)
	}
}

func getRestCall(url string) (string, error) {
	resp, err := http.Get(url)

	fmt.Println(url)

	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	return string(data), err
}

func getAvgReustData(data string, jpath string) float64 {
	json := gjson.Get(data, jpath)

	mapData := funk.Map(json.Array(), func(x gjson.Result) float64 {
		return x.Num
	})

	filterData := funk.Filter(mapData, func(x float64) bool {
		return x > 0.0
	}).([]float64)

	var resultData float64
	if len(filterData) > 0 {
		resultData = funk.SumFloat64(filterData) / float64(len(filterData))
		resultData, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", resultData), 0)
	}
	return resultData
}

func getAlarmInfos(dbClient *gorm.DB, monitoringUrl string) {
	alarmInfos := dao.GetBatchAlarmInfo(dbClient)

	var isEqual bool = false
	if len(alarmInfos) > 0 {
		saasAlarms := funk.Filter(alarmInfos, func(x model.BatchAlarmInfo) bool {
			return x.ServiceType == "SaaS"
		}).([]model.BatchAlarmInfo)

		if len(equalCheckAlarmInfos) == 0 {
			isEqual = false
		} else {
			if len(saasAlarms) == len(equalCheckAlarmInfos) {
				isEqual = reflect.DeepEqual(saasAlarms, equalCheckAlarmInfos)
			} else {
				isEqual = false
			}
		}

		if isEqual == false {
			copyAlarmInfo(saasAlarms, &equalCheckAlarmInfos, &procAlarmInfos)
			for i := 0; i < len(procAlarmInfos); i++ {
				caasAlarm := &procAlarmInfos[i]
				caasAlarm.NextMinute = nextScheduleMinute(*caasAlarm)

				go jobSaasAlarm(dbClient, *caasAlarm, monitoringUrl)
			}
		} else {
			mninute := strconv.Itoa(time.Now().Minute())
			for i := 0; i < len(procAlarmInfos); i++ {
				caasAlarm := &procAlarmInfos[i]
				if mninute == caasAlarm.NextMinute {
					caasAlarm.NextMinute = nextScheduleMinute(*caasAlarm)
					go jobSaasAlarm(dbClient, *caasAlarm, monitoringUrl)
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

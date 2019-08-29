package saas

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/mileusna/crontab"
	"github.com/thoas/go-funk"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"kr/paasta/batch/dao"
	"kr/paasta/batch/model"
	"kr/paasta/batch/notify"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Startschedule(dbClient *gorm.DB, alarmInfos []model.BatchAlarmInfo, pinpointUrl string) {
	ctab := crontab.New() // create cron table

	for _, alarmInfo := range alarmInfos {
		ctab.MustAddJob(alarmInfo.CronExpression, jobSaasAlarm, dbClient, alarmInfo, pinpointUrl)
	}

	ctab.RunAll()
}

func jobSaasAlarm(dbClient *gorm.DB, alarmInfo model.BatchAlarmInfo, pinpointUrl string) {
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

func appRunningStat(c chan model.BatchAlarmExecution, waitGroup *sync.WaitGroup, dbClient *gorm.DB, appName string, agentId string, pinpointUrl string, alarmInfo model.BatchAlarmInfo) {
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
		executionResult = "[Critical] " + strings.Replace(executionResult, "${Currend_value}", fmt.Sprintf("%.1f", result), 1) + ", 임계값 (" + strconv.Itoa(alarmInfo.CriticalValue) + ")%"
	} else if result > float64(alarmInfo.WarningValue) {
		criticalStatus = "Warning"
		executionResult = "[Warning] " + strings.Replace(executionResult, "${Currend_value}", fmt.Sprintf("%.1f", result), 1) + ", 임계값 (" + strconv.Itoa(alarmInfo.WarningValue) + ")%"
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
func appNotRunningStat(dbClient *gorm.DB, waitGroup *sync.WaitGroup, appName string, alarmInfo model.BatchAlarmInfo, criticalStatus string) {
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
func appStatReport(dbClient *gorm.DB, alarmInfo model.BatchAlarmInfo, batchExecution model.BatchAlarmExecution, result float64) {
	// batch_alarm_executions 테이블 Insert
	go dao.InsertBatchExecution(dbClient, &batchExecution)

	if result > float64(alarmInfo.CriticalValue) || result > float64(alarmInfo.WarningValue) {
		_, snsIds := dao.GetBatchAlarmReceiver(alarmInfo.ServiceType, dbClient)
		//go notify.SendMail(alarmInfo.ServiceType, emails, model.MailContent {alarmInfo, batchExecution})
		go notify.SendChatBot(alarmInfo.ServiceType, snsIds, batchExecution.ExecutionResult)
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

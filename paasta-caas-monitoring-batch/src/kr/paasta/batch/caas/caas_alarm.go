package caas

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/mileusna/crontab"
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

func Startschedule(dbClient *gorm.DB, alarmInfos []model.BatchAlarmInfo, monitoringUrl string) {
	ctab := crontab.New() // create cron table

	for _, alarmInfo := range alarmInfos {
		ctab.MustAddJob(alarmInfo.CronExpression, jobCaasAlarm, dbClient, alarmInfo, monitoringUrl)
	}

	ctab.RunAll()
}

func jobCaasAlarm(dbClient *gorm.DB, alarmInfo model.BatchAlarmInfo, monitoringUrl string) {
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

func podRunningStat(c chan model.BatchAlarmExecution, waitGroup *sync.WaitGroup, dbClient *gorm.DB, result map[string]gjson.Result, alarmInfo model.BatchAlarmInfo) {
	defer waitGroup.Done()

	var value float64
	warning := float64(alarmInfo.WarningValue)
	critical := float64(alarmInfo.CriticalValue)
	podName := result["PodName"].String()

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
		executionResult = "[Critical] " + strings.Replace(executionResult, "${Currend_value}", fmt.Sprintf("%.2f", value), 1) + ", 임계값 (" + strconv.Itoa(alarmInfo.CriticalValue) + ")%"
	} else if value > warning {
		criticalStatus = "Warning"
		executionResult = "[Warning] " + strings.Replace(executionResult, "${Currend_value}", fmt.Sprintf("%.2f", value), 1) + ", 임계값 (" + strconv.Itoa(alarmInfo.WarningValue) + ")%"
	} else {
		executionResult = strings.Replace(executionResult, "${Currend_value}", fmt.Sprintf("%.1f", value), 1)
	}

	batchExecution := model.BatchAlarmExecution{
		CriticalStatus:  criticalStatus,
		AlarmId:         alarmInfo.AlarmId,
		ServiceType:     alarmInfo.ServiceType,
		MeasureValue:    value,
		MeasureName1:    podName,
		MeasureName2:    "",
		MeasureName3:    "",
		ExecutionTime:   time.Now().Local(),
		ExecutionResult: executionResult,
	}

	c <- batchExecution
	podStatReport(dbClient, alarmInfo, batchExecution, value)
}

// CaaS 스케쥴 DB 저장 및 임계지 이상 값 알림 처리
func podStatReport(dbClient *gorm.DB, alarmInfo model.BatchAlarmInfo, batchExecution model.BatchAlarmExecution, result float64) {
	// batch_alarm_executions 테이블 Insert
	go dao.InsertBatchExecution(dbClient, &batchExecution)

	if result > float64(alarmInfo.CriticalValue) || result > float64(alarmInfo.WarningValue) {
		_, snsIds := dao.GetBatchAlarmReceiver(alarmInfo.ServiceType, dbClient)
		//go notify.SendMail(alarmInfo.ServiceType, emails, model.MailContent {alarmInfo, batchExecution})
		go notify.SendChatBot(alarmInfo.ServiceType, snsIds, batchExecution.ExecutionResult)
	}
}

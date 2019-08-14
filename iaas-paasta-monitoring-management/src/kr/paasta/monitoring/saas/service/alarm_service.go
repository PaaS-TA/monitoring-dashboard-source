package service

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"kr/paasta/monitoring/saas/dao"
	"kr/paasta/monitoring/saas/model"
	//"kr/paasta/monitoring/caas/util"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	// metricUrl
	SUB_URI1 = "/api/v1/query?query="
)

type AlarmService struct {
	promethusUrl      string
	promethusRangeUrl string
	k8sApiUrl         string
}

func GetAlarmService() *AlarmService {
	return &AlarmService{}
}

func (s *AlarmService) GetAlarmInfo() (model.ResultAlarmInfo, model.ErrMessage) {
	var alramInfo []model.AlarmInfo
	var resultAlarmInfo model.ResultAlarmInfo
	var measuringTime int

	dbAccessObj := dao.GetdbAccessObj()

	//alarm Info
	alarmInfos, err := dao.GetBatchAlarmInfo(dbAccessObj)
	if err != nil {
		fmt.Println(err)
		return resultAlarmInfo, err
	}

	//Notification Info
	alarmNotis, err1 := dao.GetBatchAlarmReceiver(dbAccessObj)
	if err1 != nil {
		fmt.Println(err1)
		return resultAlarmInfo, err1
	}

	alramInfo = make([]model.AlarmInfo, len(alarmInfos))

	for idx, data := range alarmInfos {
		delay := data.CronExpression
		index := strings.Fields(delay)
		alramInfo[idx].Name = data.MetricType
		alramInfo[idx].Critical = data.CriticalValue
		alramInfo[idx].Warning = data.WarningValue
		alramInfo[idx].Delay = fmt.Sprintf("%s", strings.Replace(index[0], "*/", "", 1))
		alramInfo[idx].AlarmId = data.AlarmId
		measuringTime = data.MeasureTime

		resultAlarmInfo.Result = alramInfo
	}

	if len(alarmNotis) > 0 {
		resultAlarmInfo.MeasuringTime = measuringTime
		resultAlarmInfo.AlarmMail = alarmNotis[0].Email
		resultAlarmInfo.AlarmTelegram = alarmNotis[0].SnsId
		resultAlarmInfo.ReceiverID = alarmNotis[0].ReceiverId
	}

	return resultAlarmInfo, nil
}

func (s *AlarmService) GetAlarmUpdate(r *http.Request) {
	dbAccessObj := dao.GetdbAccessObj()

	// 결과 출력
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	str2 := string(data)

	jsonString1 := gjson.Get(str2, "Threshold")
	jsonString2 := gjson.Get(str2, "MeasuringTime")
	jsonString3 := gjson.Get(str2, "AlarmMail")
	jsonString4 := gjson.Get(str2, "AlarmTelegram")

	temp := jsonString1.Array()
	measuringTime, _ := strconv.Atoi(jsonString2.String())
	alarmMail := jsonString3.String()
	alarmTelegram := jsonString4.Int()

	err1 := dao.InsertAlarmInfo(dbAccessObj, temp, measuringTime)

	if err1 != nil {
		fmt.Println(err)
	}

	err2 := dao.InsertAlarmReceivers(dbAccessObj, alarmMail, alarmTelegram)

	if err2 != nil {
		fmt.Println(err2)
	}
}

func (s *AlarmService) GetAlarmLog() ([]model.AlarmLog, model.ErrMessage) {
	var alarmLog []model.AlarmLog

	dbAccessObj := dao.GetdbAccessObj()

	//alarm Info
	alarmLogs, err := dao.GetBatchAlarmLog(dbAccessObj)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	alarmLog = make([]model.AlarmLog, len(alarmLogs))

	for idx, data := range alarmLogs {
		alarmLog[idx].NameSpace = data.MeasureName2
		alarmLog[idx].WorkNode = data.MeasureName1
		alarmLog[idx].Issue = data.ExecutionResult
		alarmLog[idx].Pod = data.MeasureName3
		alarmLog[idx].Status = data.CriticalStatus
		alarmLog[idx].Time = data.ExecutionTime
	}

	return alarmLog, nil
}

package service

import (
	"fmt"
	"kr/paasta/monitoring/caas/dao"
	"kr/paasta/monitoring/caas/model"
	"kr/paasta/monitoring/caas/util"
	"log"
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
	config, err := util.ReadConfig(`config.ini`)
	prometheusUrl, _ := config["prometheus.addr"]
	url := prometheusUrl + SUB_URI

	k8sApiUrl, _ := config["kubernetesApi.addr"]
	k8sUrl := k8sApiUrl + K8S_SUB_URI

	rangeUrl := prometheusUrl + SUB_URI_RANGE

	if err != nil {
		log.Println(err)
	}

	return &AlarmService{
		promethusUrl:      url,
		k8sApiUrl:         k8sUrl,
		promethusRangeUrl: rangeUrl,
	}
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
		index := strings.LastIndex(delay, "/") + 1
		alramInfo[idx].Name = data.MetricType
		alramInfo[idx].Critical = data.CriticalValue
		alramInfo[idx].Warning = data.WarningValue
		alramInfo[idx].Delay = fmt.Sprintf("%c", delay[index])
		measuringTime = data.MeasureTime
	}

	resultAlarmInfo.MeasuringTime = measuringTime
	resultAlarmInfo.Result = alramInfo
	resultAlarmInfo.AlarmMail = alarmNotis[0].Email
	resultAlarmInfo.AlarmTelegram = alarmNotis[0].SnsId

	return resultAlarmInfo, nil
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

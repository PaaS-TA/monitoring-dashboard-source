package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/thoas/go-funk"
	"kr/paasta/monitoring/saas/dao"
	"kr/paasta/monitoring/saas/model"
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
	txn               *gorm.DB
}

func GetAlarmService(txn *gorm.DB) *AlarmService {
	return &AlarmService{
		txn: txn,
	}
}

func (s *AlarmService) GetAlarmInfo() ([]model.AlarmPolicyResponse, model.ErrMessage) {
	var alramInfo []model.AlarmPolicyResponse

	//alarm Info
	alarmInfos, err := dao.GetBatchAlarmInfo(s.txn)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	//Notification Info
	alarmNotis, err1 := dao.GetBatchAlarmReceiver(s.txn, "EMAIL")
	if err1 != nil {
		fmt.Println(err1)
		return nil, err1
	}

	alramInfo = make([]model.AlarmPolicyResponse, len(alarmInfos))

	for idx, data := range alarmInfos {
		delay := data.CronExpression
		index := strings.Fields(delay)
		alramInfo[idx].Id = data.AlarmId
		alramInfo[idx].OriginType = data.ServiceType
		alramInfo[idx].AlarmType = data.MetricType
		alramInfo[idx].WarningThreshold = data.WarningValue
		alramInfo[idx].CriticalThreshold = data.CriticalValue

		repeatTime, _ := strconv.Atoi(fmt.Sprintf("%s", strings.Replace(index[0], "*/", "", 1)))

		alramInfo[idx].RepeatTime = repeatTime
		alramInfo[idx].Comment = ""
		alramInfo[idx].MeasureTime = data.MeasureTime

		if len(alarmNotis) > 0 {
			alramInfo[idx].MailAddress = alarmNotis[0].TargetId
			alramInfo[idx].MailSendYn = alarmNotis[0].UseYn
		}

	}

	return alramInfo, nil
}

func (s *AlarmService) GetAlarmUpdate(request []model.AlarmPolicyRequest) model.ErrMessage {
	email := request[3].MailAddress
	emailUseYn := request[3].MailSendYn
	err1 := dao.InsertAlarmInfo(s.txn, request[:3], email, emailUseYn)

	if err1 != nil {
		fmt.Println(err1)
	}

	return err1
}

func (s *AlarmService) GetAlarmLog(searchDateFrom string, searchDateTo string, alarmType string, alarmStatus string, resolveStatus string) ([]model.AlarmLog, model.ErrMessage) {
	var alarmLog []model.AlarmLog

	//alarm Info
	alarmLogs, err := dao.GetBatchAlarmLog(s.txn, searchDateFrom, searchDateTo, alarmType, alarmStatus, resolveStatus)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	alarmLog = make([]model.AlarmLog, len(alarmLogs))

	for idx, data := range alarmLogs {
		alarmLog[idx].Id = data.ExcutionId
		alarmLog[idx].Issue = data.ExecutionResult
		alarmLog[idx].Application = data.MeasureName1
		alarmLog[idx].Status = data.CriticalStatus
		alarmLog[idx].ResolveStatus = data.ResolveStatus

		if data.ResolveStatus == "1" {
			alarmLog[idx].ResolveStatusName = "Alarm 발생"
		} else if data.ResolveStatus == "2" {
			alarmLog[idx].ResolveStatusName = "Alarm 처리중"
		} else if data.ResolveStatus == "3" {
			alarmLog[idx].ResolveStatusName = "Alarm 처리완료"
		}

		alarmLog[idx].RegDate = data.ExecutionTime.Format("2006-01-02 15:04:05")
		if funk.IsZero(data.CompleteDate) {
			alarmLog[idx].CompleteDate = ""
		} else {
			alarmLog[idx].CompleteDate = data.CompleteDate.Format("2006-01-02 15:04:05")
		}

		alarmResolves, err := dao.GetBatchAlarmResolve(s.txn, data.ExcutionId)
		if err == nil {
			alarmLog[idx].Data = alarmResolves
		}
	}

	return alarmLog, nil
}

func (s *AlarmService) GetSnsInfo() (interface{}, interface{}) {
	var alarmSns model.BatchAlarmSnsRequest

	alarmSns, err := dao.GetSnsInfo(s.txn)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return alarmSns, nil
}

func (s *AlarmService) GetAlarmCount() (model.AlarmCount, interface{}) {
	var alarmCount model.AlarmCount
	alarmCount, err := dao.GetAlarmCount(s.txn)

	if err != nil {
		fmt.Println(err)
		return alarmCount, err
	}

	return alarmCount, nil
}

func (s *AlarmService) GetlarmSnsSave(alarmSns model.BatchAlarmSnsRequest) interface{} {
	err := dao.GetAlarmSnsSave(s.txn, alarmSns)
	return err
}

func (s *AlarmService) UpdateAlarmSate(request model.AlarmrRsolveRequest) interface{} {
	err := dao.UpdateAlarmSate(s.txn, request)
	return err
}

func (s *AlarmService) CreateAlarmResolve(request model.AlarmrRsolveRequest) interface{} {
	err := dao.CreateAlarmResolve(s.txn, request)
	return err
}

func (s *AlarmService) UpdateAlarmResolve(request model.AlarmrRsolveRequest) interface{} {
	err := dao.UpdateAlarmResolve(s.txn, request)
	return err
}

func (s *AlarmService) DeleteAlarmResolve(id uint64) interface{} {
	err := dao.DeleteAlarmResolve(s.txn, id)
	return err

}

func (s *AlarmService) GetAlarmSnsReceiver() ([]model.BatchAlarmReceiver, interface{}) {
	var alarmReceiver []model.BatchAlarmReceiver
	alarmReceiver, err := dao.GetBatchAlarmReceiver(s.txn, "SNS")

	if err != nil {
		return nil, err
	}
	return alarmReceiver, err
}

func (s *AlarmService) DeleteAlarmSnsChannel(id int) interface{} {
	err := dao.DeleteAlarmSnsChannel(s.txn, id)
	return err
}

func (s *AlarmService) GetAlarmActionList(id int) ([]model.AlarmActionResponse, model.ErrMessage) {
	alarmResolves, err := dao.GetBatchAlarmResolve(s.txn, uint64(id))

	alarmActionResponses := make([]model.AlarmActionResponse, len(alarmResolves))

	for index, action := range alarmResolves {
		alarmActionResponses[index].AlarmId = uint64(id)
		alarmActionResponses[index].ResolveId = action.ResolveId
		alarmActionResponses[index].AlarmActionDesc = action.AlarmActionDesc
		alarmActionResponses[index].RegDate = action.RegDate
	}

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return alarmActionResponses, nil
}

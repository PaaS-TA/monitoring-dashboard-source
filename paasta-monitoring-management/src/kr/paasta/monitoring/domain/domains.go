package domain

import (
	"errors"
	"fmt"
)

type (
	AlarmRequest struct {
		Id                uint
		OriginType        string
		OriginId          uint
		AlarmType         string
		Level             string
		AlarmTitle        string
		ResolveStatus     string
		SearchDateFrom    string
		SearchDateTo      string
		PagingReq
	}

	AlarmActionRequest struct {
		Id                uint
		AlarmId           uint
		AlarmActionDesc   string
		RegDate           JSONTime
		RegUser           string
		ModiDate          JSONTime
		ModiUser          string
	}

	AlarmStatRequest struct {
		Interval          int
		Period            string
		SearchDateFrom    string
		SearchDateTo      string
	}

	AlarmPolicyRequest struct {
		OriginType         string   `json:"originType"`
		AlarmType          string   `json:"alarmType"`
		WarningThreshold   int      `json:"warningThreshold"`
		CriticalThreshold  int      `json:"criticalThreshold"`
		RepeatTime         int      `json:"repeatTime"`
		Comment            string   `json:"comment"`
	}

	AlarmPolicy struct {
		OriginType         string   `json:"originType"`
		AlarmType          string   `json:"alarmType"`
		WarningThreshold   int      `json:"warningThreshold"`
		CriticalThreshold  int      `json:"criticalThreshold"`
		RepeatTime  int             `json:"repeatTime"`
		Comment      string	    `json:"comment"`
	}

	AlarmPolicyResponse struct {
		Id                uint     `json:"id"`
		OriginType        string   `json:"originType"`
		AlarmType         string   `json:"alarmType"`
		WarningThreshold  int      `json:"warningThreshold"`
		CriticalThreshold int      `json:"criticalThreshold"`
		RepeatTime        int      `json:"repeatTime"`
		Comment           string   `json:"comment"`
	}

	AlarmResponse struct {
		Id                uint		`json:"id"`
		OriginType        string	`json:"originType"`
		OriginId          uint		`json:"originId"`
		OriginName        string	`json:"originName"`
		AlarmType         string	`json:"alarmType"`
		Level             string	`json:"level"`
		AlarmTitle        string	`json:"alarmTitle"`
		AlarmMessage      string	`json:"alarmMessage"`
		Ip                string	`json:"ip"`
		ResolveStatus     string	`json:"resolveStatus"`
		ResolveStatusName string	`json:"resolveStatusName"`
		AppYn             string	`json:"appYn"`
		AppName           string	`json:"appName"`
		AppIndex          int		`json:"appIndex"`
		ContainerName     string	`json:"containerName"`
		AlarmCnt          int		`json:"alarmCnt"`
		AlarmSendDate     JSONTime	`json:"alarmSendDate"`
		RegDate           JSONTime	`json:"regDate"`
		RegUser           string	`json:"regUser"`
		UserName          string	`json:"userName"`
	}

	AlarmPagingResponse struct {
		PageIndex         int		`json:"pageIndex"`
		PageItem          int		`json:"pageItem"`
		TotalCount        int		`json:"totalCount"`
		AlarmResponse []AlarmResponse   `json:"data"`
	}

	AlarmDetailResponse struct {
		Id                uint		`json:"id"`
		OriginType        string	`json:"originType"`
		OriginId          uint		`json:"originId"`
		AlarmType         string	`json:"alarmType"`
		Level             string	`json:"level"`
		AlarmTitle        string	`json:"alarmTitle"`
		AlarmMessage      string	`json:"alarmMessage"`
		OriginName        string	`json:"originName"`
		Ip                string	`json:"ip"`
		ResolveStatusName string	`json:"resolveStatusName"`
		ResolveStatus     string	`json:"resolveStatus"`
		AppYn             string	`json:"appYn"`
		AppName           string	`json:"appName"`
		AppIndex          int		`json:"appIndex"`
		ContainerName     string	`json:"containerName"`
		AlarmCnt          int		`json:"alarmCnt"`
		AlarmSendDate     JSONTime	`json:"alarmSendDate"`
		RegDate           JSONTime	`json:"regDate"`
		RegUser           string	`json:"regUser"`
		AlarmActionResponse []AlarmActionResponse   `json:"data"`
	}

	AlarmActionResponse struct {
		Id                 uint		`json:"id"`
		AlarmId            string	`json:"alarmId"`
		AlarmActionDesc    string	`json:"alarmActionDesc"`
		RegDate            JSONTime	`json:"regDate"`
		RegUser            string	`json:"regUser"`
		ModiDate           JSONTime	`json:"ModiDate"`
		ModiUser           string	`json:"modiUser"`
	}

	AlarmStatResponse struct {
		Id                      uint		`json:"id"`
		TotalCnt                int		`json:"totalCnt"`
		WarningCnt              int		`json:"warningCnt"`
		CriticalCnt             int		`json:"criticalCnt"`
		PaastaWarningCnt        int		`json:"paastaWarningCnt"`
		PaastaCriticalCnt       int		`json:"paastaCriticalCnt"`
		BoshWarningCnt          int		`json:"boshWarningCnt"`
		BoshCriticalCnt         int		`json:"boshCriticalCnt"`
		ContainerWarningCnt     int		`json:"containerWarningCnt"`
		ContainerCriticalCnt    int		`json:"containerCriticalCnt"`
		TotalResolveCnt         int		`json:"totalResolveCnt"`
		WarningResolveCnt       int		`json:"warningResolveCnt"`
		CriticalResolveCnt      int		`json:"criticalResolveCnt"`
		CpuWarningCnt           int		`json:"cpuWarningCnt"`
		CpuCriticalCnt          int		`json:"cpuCriticalCnt"`
		MemoryWarningCnt        int		`json:"memoryWarningCnt"`
		MemoryCriticalCnt       int		`json:"memoryCriticalCnt"`
		DiskWarningCnt          int		`json:"diskWarningCnt"`
		DiskCriticalCnt         int		`json:"diskCriticalCnt"`
	}
)

func (bm AlarmRequest) AlarmValidate(req AlarmRequest) error {
	if req.Id == 0 {
		return errors.New("Required input value does not exist. [Id]");
	}
	if req.ResolveStatus == "" {
		return errors.New("Required input value does not exist. [ResolveStatus]");
	}
	return nil
}
func (bm AlarmActionRequest) AlarmActionValidate(req AlarmActionRequest) error {
	if req.Id == 0 {
		return errors.New("Required input value does not exist. [Id]");
	}
	return nil
}

//Alarm 정책정보 유효성 체크
func (bm AlarmPolicyRequest) AlarmPolicyValidate(requests AlarmPolicyRequest) error {


	if requests.OriginType == "" {
		return errors.New("Required input value does not exist. [originType]");
	}

	if requests.AlarmType == "" {
		return errors.New("Required input value does not exist. [alarmType]");
	}

	if requests.WarningThreshold == 0 {
		return errors.New("Required input value does not exist. [warningThreshold]");
	}

	if requests.CriticalThreshold == 0 {
		return errors.New("Required input value does not exist. [criticalThreshold]");
	}

	fmt.Println("requests.RepeatTime:::::",requests.RepeatTime)
	if requests.RepeatTime <= 0 {
		return errors.New("Required input value does not exist. [repeatTime]");
	}

	if requests.WarningThreshold >= requests.CriticalThreshold {
		return errors.New("[warningThreshold] can not greater than criticalThreshold or equal");
	}

	if requests.CriticalThreshold > 100 {
		return errors.New("[CriticalThreshold] can not greater than 100");
	}

	return nil
}

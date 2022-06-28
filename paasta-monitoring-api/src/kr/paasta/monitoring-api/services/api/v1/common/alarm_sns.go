package common

import (
	"gorm.io/gorm"
	"paasta-monitoring-api/dao/api/v1/common"
	models "paasta-monitoring-api/models/api/v1"
	"time"
)


type AlarmSnsService struct {
	DbInfo *gorm.DB
}

func GetAlarmSnsService(DbInfo *gorm.DB) *AlarmSnsService {
	return &AlarmSnsService {
		DbInfo: DbInfo,
	}
}


func (service *AlarmSnsService) CreateAlarmSns(params []models.AlarmSns, regUser string) (string, error) {
	for _, param := range params {
		param.RegUser = regUser
		param.RegDate = time.Now()
	}

	err := common.GetAlarmSnsDao(service.DbInfo).CreateAlarmSns(params)
	if err != nil {
		return "FAILED REGISTER SNS ACCOUNT!", err
	}
	return "SUCCEEDED REGISTER SNS ACCOUNT!", nil
}

func (service *AlarmSnsService) GetAlarmSns(params models.AlarmSns) ([]models.AlarmSns, error) {
	results, err := common.GetAlarmSnsDao(service.DbInfo).GetAlarmSns(params)
	if err != nil {
		return results, err
	}
	return results, nil
}


func (service *AlarmSnsService) UpdateAlarmSns(request *models.AlarmSns) (string, error) {
	err := common.GetAlarmSnsDao(service.DbInfo).UpdateAlarmSns(request)
	if err != nil {
		return "FAILED UPDATE SNS ACCOUNT!", err
	}
	return "SUCCEEDED UPDATE SNS ACCOUNT!", nil
}


func (service *AlarmSnsService) DeleteAlarmSns(request models.AlarmSns) (string, error) {
	err := common.GetAlarmSnsDao(service.DbInfo).DeleteAlarmSns(request)
	if err != nil {
		return "FAILED DELETE SNS ACCOUNT!", err
	}
	return "SUCCEEDED DELETE SNS ACCOUNT!", nil
}

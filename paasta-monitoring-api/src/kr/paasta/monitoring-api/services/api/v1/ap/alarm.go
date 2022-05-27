package ap

import (
	AP "GoEchoProject/dao/api/v1/ap"
	models "GoEchoProject/models/api/v1"
	"github.com/jinzhu/gorm"
)

type ApAlarmService struct {
	DbInfo *gorm.DB
}

func GetApAlarmService(DbInfo *gorm.DB) *ApAlarmService {
	return &ApAlarmService{
		DbInfo: DbInfo,
	}
}

func (ap *ApAlarmService) GetAlarmStatus() ([]models.Alarms, error) {
	results, err := AP.GetApAlarmDao(ap.DbInfo).GetAlarmStatus()
	if err != nil {
		return results, err
	}
	return results, nil
}

func (ap *ApAlarmService) GetAlarmPolicy() ([]models.AlarmPolicies, error) {
	results, err := AP.GetApAlarmDao(ap.DbInfo).GetAlarmPolicy()
	if err != nil {
		return results, err
	}
	return results, nil
}

func (ap *ApAlarmService) UpdateAlarmPolicy(request []models.AlarmPolicyRequest) (string, error) {
	for _, request := range request {
		err := AP.GetApAlarmDao(ap.DbInfo).UpdateAlarmPolicy(request)
		if err != nil {
			return "FAILED UPDATE POLICY!", err
		}
	}

	return "SUCCEEDED UPDATE POLICY!", nil
}

func (ap *ApAlarmService) UpdateAlarmTarget(request []models.AlarmTargetRequest) (string, error) {
	for _, request := range request {
		err := AP.GetApAlarmDao(ap.DbInfo).UpdateAlarmTarget(request)
		if err != nil {
			return "FAILED UPDATE TARGET!", err
		}
	}

	return "SUCCEEDED UPDATE TARGET!", nil
}

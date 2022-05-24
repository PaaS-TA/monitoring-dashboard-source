package ap

import (
	AP "GoEchoProject/dao/api/v1/ap"
	models "GoEchoProject/models/api/v1"
	"github.com/jinzhu/gorm"
)

type ApService struct {
	DbInfo *gorm.DB
}

func GetApService(DbInfo *gorm.DB) *ApService {
	return &ApService{
		DbInfo: DbInfo,
	}
}

func (ap *ApService) GetAlarmStatus() ([]models.Alarms, error) {
	results, err := AP.GetApDao(ap.DbInfo).GetAlarmStatus()
	if err != nil {
		return results, err
	}
	return results, nil
}

func (ap *ApService) GetAlarmPolicy() ([]models.AlarmPolicies, error) {
	results, err := AP.GetApDao(ap.DbInfo).GetAlarmPolicy()
	if err != nil {
		return results, err
	}
	return results, nil
}

func (ap *ApService) UpdateAlarmPolicy(request []models.AlarmPolicyRequest) (string, error) {
	for _, request := range request {
		err := AP.GetApDao(ap.DbInfo).UpdateAlarmPolicy(request)
		if request.MailAddress != "" {
			ap.UpdateAlarmTarget(request)
		}
		if err != nil {
			return "FAILED UPDATE POLICY!", err
		}
	}

	return "SUCCEEDED!", nil
}

func (ap *ApService) UpdateAlarmTarget(request models.AlarmPolicyRequest) (string, error) {
	err := AP.GetApDao(ap.DbInfo).UpdateAlarmTarget(request)
	if err != nil {
		return "FAILED UPDATE TARGET!", err
	}

	return "SUCCEEDED!", nil
}

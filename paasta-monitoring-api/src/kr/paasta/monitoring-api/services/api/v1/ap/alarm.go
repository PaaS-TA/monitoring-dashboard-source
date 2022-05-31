package ap

import (
	"github.com/jinzhu/gorm"
	AP "paasta-monitoring-api/dao/api/v1/ap"
	models "paasta-monitoring-api/models/api/v1"
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

func (ap *ApAlarmService) RegisterSnsAccount(request models.SnsAccountRequest) (string, error) {
	err := AP.GetApAlarmDao(ap.DbInfo).RegisterSnsAccount(request)
	if err != nil {
		return "FAILED REGISTER SNS ACCOUNT!", err
	}
	return "SUCCEEDED REGISTER SNS ACCOUNT!", nil
}

func (ap *ApAlarmService) GetSnsAccount() ([]models.AlarmSns, error) {
	results, err := AP.GetApAlarmDao(ap.DbInfo).GetSnsAccount()
	if err != nil {
		return results, err
	}
	return results, nil
}

func (ap *ApAlarmService) DeleteSnsAccount(request models.SnsAccountRequest) (string, error) {
	err := AP.GetApAlarmDao(ap.DbInfo).DeleteSnsAccount(request)
	if err != nil {
		return "FAILED DELETE SNS ACCOUNT!", err
	}
	return "SUCCEEDED DELETE SNS ACCOUNT!", nil
}

func (ap *ApAlarmService) UpdateSnsAccount(request models.SnsAccountRequest) (string, error) {
	err := AP.GetApAlarmDao(ap.DbInfo).UpdateSnsAccount(request)
	if err != nil {
		return "FAILED UPDATE SNS ACCOUNT!", err
	}
	return "SUCCEEDED UPDATE SNS ACCOUNT!", nil
}

func (ap *ApAlarmService) CreateAlarmAction(request models.AlarmActionRequest) (string, error) {
	err := AP.GetApAlarmDao(ap.DbInfo).CreateAlarmAction(request)
	if err != nil {
		return "FAILED CREATE ALARM ACTION!", err
	}
	return "SUCCEEDED CREATE ALARM ACTION!", nil
}

func (ap *ApAlarmService) GetAlarmAction() ([]models.AlarmActions, error) {
	results, err := AP.GetApAlarmDao(ap.DbInfo).GetAlarmAction()
	if err != nil {
		return results, err
	}
	return results, nil
}

func (ap *ApAlarmService) UpdateAlarmAction(request models.AlarmActionRequest) (string, error) {
	err := AP.GetApAlarmDao(ap.DbInfo).UpdateAlarmAction(request)
	if err != nil {
		return "FAILED UPDATE ALARM ACTION!", err
	}
	return "SUCCEEDED UPDATE ALARM ACTION!", nil
}

func (ap *ApAlarmService) DeleteAlarmAction(request models.AlarmActionRequest) (string, error) {
	err := AP.GetApAlarmDao(ap.DbInfo).DeleteAlarmAction(request)
	if err != nil {
		return "FAILED DELETE ALARM ACTION!", err
	}
	return "SUCCEEDED DELETE ALARM ACTION!", nil
}

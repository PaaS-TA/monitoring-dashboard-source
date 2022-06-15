package ap

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
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

func (ap *ApAlarmService) GetAlarmStatisticsTotal(c echo.Context) ([]map[string]interface{}, error) {
	results, err := AP.GetApAlarmDao(ap.DbInfo).GetAlarmStatisticsTotal(c)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (ap *ApAlarmService) GetAlarmStatisticsService(c echo.Context) ([]map[string]interface{}, error) {
	results, err := AP.GetApAlarmDao(ap.DbInfo).GetAlarmStatisticsService(c)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (ap *ApAlarmService) GetAlarmStatisticsResource(c echo.Context) ([]map[string]interface{}, error) {
	results, err := AP.GetApAlarmDao(ap.DbInfo).GetAlarmStatisticsResource(c)
	if err != nil {
		return results, err
	}
	return results, nil
}

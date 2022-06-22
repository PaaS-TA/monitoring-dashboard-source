package common

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	common "paasta-monitoring-api/dao/api/v1/common"
	models "paasta-monitoring-api/models/api/v1"
)

type AlarmPolicyService struct {
	DbInfo *gorm.DB
}

func GetAlarmPolicyService(DbInfo *gorm.DB) *AlarmPolicyService {
	return &AlarmPolicyService{
		DbInfo: DbInfo,
	}
}

func (service *AlarmPolicyService) GetAlarmPolicy(c echo.Context) ([]models.AlarmPolicies, error) {
	params := models.AlarmPolicies {
		OriginType : c.QueryParam("originType"),
		AlarmType : c.QueryParam("alarmType"),
	}

	results, err := common.GetAlarmPolicyDao(service.DbInfo).GetAlarmPolicy(params)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (service *AlarmPolicyService) UpdateAlarmPolicy(request []models.AlarmPolicyRequest) (string, error) {
	for _, request := range request {
		err := common.GetAlarmPolicyDao(service.DbInfo).UpdateAlarmPolicy(request)
		if err != nil {
			return "FAILED UPDATE POLICY!", err
		}
	}
	return "SUCCEEDED UPDATE POLICY!", nil
}

func (service *AlarmPolicyService) UpdateAlarmTarget(request []models.AlarmTargetRequest) (string, error) {
	for _, request := range request {
		err := common.GetAlarmPolicyDao(service.DbInfo).UpdateAlarmTarget(request)
		if err != nil {
			return "FAILED UPDATE TARGET!", err
		}
	}
	return "SUCCEEDED UPDATE TARGET!", nil
}

package common

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	Common "paasta-monitoring-api/dao/api/v1"
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

func (common *AlarmPolicyService) GetAlarmPolicy(c echo.Context) ([]models.AlarmPolicies, error) {
	results, err := Common.GetCommonDao(common.DbInfo).GetAlarmPolicy(c)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (common *AlarmPolicyService) UpdateAlarmPolicy(request []models.AlarmPolicyRequest) (string, error) {
	for _, request := range request {
		err := Common.GetCommonDao(common.DbInfo).UpdateAlarmPolicy(request)
		if err != nil {
			return "FAILED UPDATE POLICY!", err
		}
	}
	return "SUCCEEDED UPDATE POLICY!", nil
}

func (common *AlarmPolicyService) UpdateAlarmTarget(request []models.AlarmTargetRequest) (string, error) {
	for _, request := range request {
		err := Common.GetCommonDao(common.DbInfo).UpdateAlarmTarget(request)
		if err != nil {
			return "FAILED UPDATE TARGET!", err
		}
	}
	return "SUCCEEDED UPDATE TARGET!", nil
}

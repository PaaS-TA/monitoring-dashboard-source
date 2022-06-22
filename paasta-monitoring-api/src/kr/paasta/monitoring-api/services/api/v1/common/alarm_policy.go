package common

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	common "paasta-monitoring-api/dao/api/v1/common"
	models "paasta-monitoring-api/models/api/v1"
	"time"
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


func (service *AlarmPolicyService) UpdateAlarmPolicy(ctx echo.Context, params []models.AlarmPolicyRequest) (string, error) {
	for _, policyParam := range params {
		param := models.AlarmPolicies{
			OriginType: policyParam.OriginType,
			AlarmType: policyParam.AlarmType,
			WarningThreshold: policyParam.WarningThreshold,
			CriticalThreshold: policyParam.CriticalThreshold,
			RepeatTime: policyParam.RepeatTime,
			MeasureTime: policyParam.MeasureTime,
			ModiUser: ctx.Get("userId").(string),
			ModiDate: time.Now(),
		}
		err := common.GetAlarmPolicyDao(service.DbInfo).UpdateAlarmPolicy(param)
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
package v1

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	models "paasta-monitoring-api/models/api/v1"
	"time"
)

type CommonDao struct {
	DbInfo *gorm.DB
}

func GetCommonDao(DbInfo *gorm.DB) *CommonDao {
	return &CommonDao{
		DbInfo: DbInfo,
	}
}

func (common *CommonDao) GetAlarmStatus() ([]models.Alarms, error) {
	var response []models.Alarms
	results := common.DbInfo.Debug().Table("alarms").
		Select("*").
		Find(&response)

	if results.Error != nil {
		fmt.Println(results.Error)
		return response, results.Error
	}

	return response, nil
}

func (common *CommonDao) GetAlarmPolicy(c echo.Context) ([]models.AlarmPolicies, error) {
	var request models.AlarmPolicies
	request.OriginType = c.QueryParam("originType")
	request.AlarmType = c.QueryParam("alarmType")

	var response []models.AlarmPolicies
	results := common.DbInfo.Debug().Table("alarm_policies").
		Select("*").
		Where(request).
		Find(&response)

	if results.Error != nil {
		fmt.Println(results.Error)
		return response, results.Error
	}

	return response, nil
}

func (common *CommonDao) UpdateAlarmPolicy(request models.AlarmPolicyRequest) error {
	results := common.DbInfo.Debug().Table("alarm_policies").
		Where("origin_type = ? AND alarm_type = ?", request.OriginType, request.AlarmType).
		Updates(map[string]interface{}{
			"warning_threshold":  request.WarningThreshold,
			"critical_threshold": request.CriticalThreshold,
			"repeat_time":        request.RepeatTime,
			"measure_time":       request.MeasureTime,
			"modi_date":          time.Now().UTC().Add(time.Hour * 9),
			"modi_user":          "admin"})

	if results.Error != nil {
		fmt.Println(results.Error)
		return results.Error
	}

	return nil
}

func (common *CommonDao) UpdateAlarmTarget(request models.AlarmTargetRequest) error {
	results := common.DbInfo.Debug().Table("alarm_targets").
		Where("origin_type = ?", request.OriginType).
		Updates(map[string]interface{}{
			"mail_address": request.MailAddress,
			"mail_send_yn": request.MailSendYN,
			"modi_date":    time.Now().UTC().Add(time.Hour * 9),
			"modi_user":    "admin"})

	if results.Error != nil {
		fmt.Println(results.Error)
		return results.Error
	}

	return nil
}

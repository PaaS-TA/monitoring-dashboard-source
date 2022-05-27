package ap

import (
	models "GoEchoProject/models/api/v1"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type ApAlarmDao struct {
	DbInfo *gorm.DB
}

func GetApAlarmDao(DbInfo *gorm.DB) *ApAlarmDao {
	return &ApAlarmDao{
		DbInfo: DbInfo,
	}
}

func (ap *ApAlarmDao) GetAlarmStatus() ([]models.Alarms, error) {
	var response []models.Alarms
	results := ap.DbInfo.Debug().Table("alarms").
		Select(" * ").
		Find(&response)

	if results.Error != nil {
		fmt.Println(results.Error)
		return response, results.Error
	}

	return response, nil
}

func (ap *ApAlarmDao) GetAlarmPolicy() ([]models.AlarmPolicies, error) {
	var response []models.AlarmPolicies
	results := ap.DbInfo.Debug().Table("alarm_policies").
		Select(" * ").
		Find(&response)

	if results.Error != nil {
		fmt.Println(results.Error)
		return response, results.Error
	}

	return response, nil
}

func (ap *ApAlarmDao) UpdateAlarmPolicy(request models.AlarmPolicyRequest) error {
	results := ap.DbInfo.Debug().Table("alarm_policies").
		Model(request).
		Where("origin_type = ? and alarm_type = ?", request.OriginType, request.AlarmType).
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

func (ap *ApAlarmDao) UpdateAlarmTarget(request models.AlarmTargetRequest) error {
	if request.MailAddress != "" {
		results := ap.DbInfo.Debug().Table("alarm_targets").
			Model(request).
			Where("origin_type = ?", request.OriginType).
			Updates(map[string]interface{}{
				"mail_address": request.MailAddress,
				"mail_send_yn": request.MailSendYn,
				"modi_date":    time.Now().UTC().Add(time.Hour * 9),
				"modi_user":    "admin"})
		if results.Error != nil {
			fmt.Println(results.Error)
			return results.Error
		}
	}

	return nil
}

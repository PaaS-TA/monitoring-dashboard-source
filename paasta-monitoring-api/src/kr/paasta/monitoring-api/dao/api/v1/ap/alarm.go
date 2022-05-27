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
		Select("*").
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
		Select("*").
		Find(&response)

	if results.Error != nil {
		fmt.Println(results.Error)
		return response, results.Error
	}

	return response, nil
}

func (ap *ApAlarmDao) UpdateAlarmPolicy(request models.AlarmPolicyRequest) error {
	results := ap.DbInfo.Debug().Table("alarm_policies").
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
	results := ap.DbInfo.Debug().Table("alarm_targets").
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

func (ap *ApAlarmDao) RegisterSnsAccount(request models.SnsAccountRequest) error {
	results := ap.DbInfo.Debug().Table("alarm_sns").Create(&request)

	if results.Error != nil {
		fmt.Println(results.Error)
		return results.Error
	}

	return nil
}

func (ap *ApAlarmDao) GetSnsAccount() ([]models.AlarmSns, error) {
	var response []models.AlarmSns
	results := ap.DbInfo.Debug().Table("alarm_sns").
		Select("*").
		Find(&response)

	if results.Error != nil {
		fmt.Println(results.Error)
		return response, results.Error
	}

	return response, nil
}

func (ap *ApAlarmDao) DeleteSnsAccount(request models.SnsAccountRequest) error {
	results := ap.DbInfo.Debug().Table("alarm_sns").
		Where("sns_id = ?", request.SnsId).
		Delete(&request)

	if results.Error != nil {
		fmt.Println(results.Error)
		return results.Error
	}

	return nil
}

func (ap *ApAlarmDao) UpdateSnsAccount(request models.SnsAccountRequest) error {
	results := ap.DbInfo.Debug().Table("alarm_sns").
		Where("sns_id = ?", request.SnsId).
		Updates(map[string]interface{}{
			"origin_type": request.OriginType,
			"sns_type":    request.SnsType,
			"sns_id":      request.SnsId,
			"token":       request.Token,
			"expl":        request.Expl,
			"sns_send_yn": request.SnsSendYN,
			"modi_date":   time.Now().UTC().Add(time.Hour * 9),
			"modi_user":   "admin"})

	if results.Error != nil {
		fmt.Println(results.Error)
		return results.Error
	}

	return nil
}

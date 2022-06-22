package common

import (
	"github.com/jinzhu/gorm"
	models "paasta-monitoring-api/models/api/v1"
	"time"
)


type AlarmSnsDao struct {
	DbInfo *gorm.DB
}

func GetAlarmSnsDao(DbInfo *gorm.DB) *AlarmSnsDao {
	return &AlarmSnsDao {
		DbInfo: DbInfo,
	}
}


func (ap *AlarmSnsDao) CreateAlarmSns(request models.SnsAccountRequest) error {
	results := ap.DbInfo.Debug().Table("alarm_sns").Create(&request)

	if results.Error != nil {
		return results.Error
	}

	return nil
}

func (ap *AlarmSnsDao) GetAlarmSns(params models.AlarmSns) ([]models.AlarmSns, error) {
	var response []models.AlarmSns
	results := ap.DbInfo.Debug().Where(&params).Find(&response)

	if results.Error != nil {
		return response, results.Error
	}

	return response, nil
}


func (ap *AlarmSnsDao) UpdateAlarmSns(request models.SnsAccountRequest) error {
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
		return results.Error
	}

	return nil
}


func (ap *AlarmSnsDao) DeleteAlarmSns(request models.SnsAccountRequest) error {
	results := ap.DbInfo.Debug().Table("alarm_sns").
		Where("sns_id = ?", request.SnsId).
		Delete(&request)

	if results.Error != nil {
		return results.Error
	}

	return nil
}


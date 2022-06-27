package common

import (
	"gorm.io/gorm"
	models "paasta-monitoring-api/models/api/v1"
)

type AlarmSnsDao struct {
	DbInfo *gorm.DB
}

func GetAlarmSnsDao(DbInfo *gorm.DB) *AlarmSnsDao {
	return &AlarmSnsDao {
		DbInfo: DbInfo,
	}
}


func (ap *AlarmSnsDao) CreateAlarmSns(request []models.AlarmSns) error {
	results := ap.DbInfo.Debug().CreateInBatches(&request, 100)
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


func (ap *AlarmSnsDao) UpdateAlarmSns(request *models.AlarmSns) error {
	results := ap.DbInfo.Debug().Model(&request).
		Where("channel_id = ?", request.ChannelId).
		Updates(&request)
	if results.Error != nil {
		return results.Error
	}

	return nil
}


func (ap *AlarmSnsDao) DeleteAlarmSns(request models.AlarmSns) error {
	results := ap.DbInfo.Debug().
		Where("channel_id = ?", request.ChannelId).
		Delete(&request)
	if results.Error != nil {
		return results.Error
	}

	return nil
}


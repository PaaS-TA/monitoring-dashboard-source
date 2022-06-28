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


func (dao *AlarmSnsDao) CreateAlarmSns(params []models.AlarmSns) error {
	results := dao.DbInfo.Debug().CreateInBatches(&params, 100)
	if results.Error != nil {
		return results.Error
	}

	return nil
}

func (dao *AlarmSnsDao) GetAlarmSns(params models.AlarmSns) ([]models.AlarmSns, error) {
	var response []models.AlarmSns
	results := dao.DbInfo.Debug().Where(&params).Find(&response)

	if results.Error != nil {
		return response, results.Error
	}

	return response, nil
}


func (dao *AlarmSnsDao) UpdateAlarmSns(params *models.AlarmSns) error {
	results := dao.DbInfo.Debug().Model(&params).
		Where("channel_id = ?", params.ChannelId).
		Updates(&params)
	if results.Error != nil {
		return results.Error
	}

	return nil
}


func (dao *AlarmSnsDao) DeleteAlarmSns(params models.AlarmSns) error {
	results := dao.DbInfo.Debug().
		Where("channel_id = ?", params.ChannelId).
		Delete(&params)
	if results.Error != nil {
		return results.Error
	}

	return nil
}


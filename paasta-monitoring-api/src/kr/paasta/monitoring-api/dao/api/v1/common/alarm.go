package common

import (
	"github.com/jinzhu/gorm"
	models "paasta-monitoring-api/models/api/v1"
)

type AlarmDao struct {
	DbInfo *gorm.DB
}

func GetAlarmDao(DbInfo *gorm.DB) *AlarmDao {
	return &AlarmDao {
		DbInfo: DbInfo,
	}
}


func (ap *AlarmDao) GetAlarms(params models.Alarms) ([]models.Alarms, error) {
	var response []models.Alarms
	results := ap.DbInfo.Debug().Where(&params).Find(&response)

	if results.Error != nil {
		return response, results.Error
	}

	return response, nil
}
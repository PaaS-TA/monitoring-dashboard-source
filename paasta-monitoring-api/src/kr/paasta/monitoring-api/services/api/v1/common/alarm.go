package common

import (
	"github.com/jinzhu/gorm"
	dao "paasta-monitoring-api/dao/api/v1/common"
	models "paasta-monitoring-api/models/api/v1"
)

type AlarmService struct {
	DbInfo *gorm.DB
}

func GetAlarmService(DbInfo *gorm.DB) *AlarmService {
	return &AlarmService {
		DbInfo: DbInfo,
	}
}

func (ap *AlarmService) GetAlarms(params models.Alarms) ([]models.Alarms, error) {
	results, err := dao.GetAlarmDao(ap.DbInfo).GetAlarms(params)
	if err != nil {
		return results, err
	}
	return results, nil
}
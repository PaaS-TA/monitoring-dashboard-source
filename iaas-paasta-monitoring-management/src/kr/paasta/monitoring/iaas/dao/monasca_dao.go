package dao

import (
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring/iaas/model"
	pm "kr/paasta/monitoring/paas/model"
)

type MonascaDao struct {
	txn   *gorm.DB
}

func GetMonascaDbDao(txn *gorm.DB) *MonascaDao {
	return &MonascaDao{
		txn:   txn,
	}
}

func (m *MonascaDao) GetAlarmsDefinition(alarmId string) (model.Alarm, error) {

	var alarm model.Alarm

	status := m.txn.Debug().Table("alarm").
		Select("alarm.id, alarm.alarm_definition_id, alarm_definition.name, alarm_definition.expression, alarm_definition.severity").
		Joins("inner join alarm_definition on alarm_definition.id = alarm.alarm_definition_id ").
		Where("alarm.id = ?", alarmId).
		Find(&alarm)

	if status.Error != nil {
		return alarm, status.Error
	}

	return alarm, status.Error
}

func (m *MonascaDao) GetAlarmCount(state string) (pm.AlarmStatusCountResponse, error) {

	t := pm.AlarmStatusCountResponse{}

	status := m.txn.Debug().Table("alarm").
		Select("count(id) as totalCnt").
		Where("state = '" + state + "'").
		Count(&t.TotalCnt)

	if status.Error != nil{
		return t, status.Error
	}
	return t, nil

}
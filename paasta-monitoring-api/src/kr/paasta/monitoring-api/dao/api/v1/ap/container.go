package ap

import (
	"fmt"
	"github.com/jinzhu/gorm"
	models "paasta-monitoring-api/models/api/v1"
)

type ApContainerDao struct {
	DbInfo *gorm.DB
}

func GetApContainerDao(DbInfo *gorm.DB) *ApContainerDao {
	return &ApContainerDao{
		DbInfo: DbInfo,
	}
}

func (ap *ApContainerDao) GetCellInfo() ([]models.CellInfo, error) {
	var response []models.CellInfo

	results := ap.DbInfo.Debug().Table("zones").Order("cell_name ASC").
		Select("zones.name AS zone_name, vms.name AS cell_name, vms.ip, vms.id").
		Joins("INNER JOIN vms ON zones.id = vms.zone_id AND vms.vm_type = 'cel'").
		Find(&response)

	if results.Error != nil {
		fmt.Println(results.Error)
		return response, results.Error
	}

	return response, nil
}

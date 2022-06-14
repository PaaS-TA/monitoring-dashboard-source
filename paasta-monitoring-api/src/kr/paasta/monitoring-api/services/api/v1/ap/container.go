package ap

import (
	"github.com/jinzhu/gorm"
	AP "paasta-monitoring-api/dao/api/v1/ap"
	models "paasta-monitoring-api/models/api/v1"
)

type ApContainerService struct {
	DbInfo *gorm.DB
}

func GetApContainerService(DbInfo *gorm.DB) *ApContainerService {
	return &ApContainerService{
		DbInfo: DbInfo,
	}
}

func (ap *ApContainerService) GetCellInfo() ([]models.CellInfo, error) {
	results, err := AP.GetApContainerDao(ap.DbInfo).GetCellInfo()
	if err != nil {
		return results, err
	}
	return results, nil
}

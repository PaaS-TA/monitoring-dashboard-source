package ap

import (
	"github.com/cloudfoundry-community/go-cfclient"
	influxDb "github.com/influxdata/influxdb1-client/v2"
	"github.com/jinzhu/gorm"
	AP "paasta-monitoring-api/dao/api/v1/ap"
	models "paasta-monitoring-api/models/api/v1"
)

type ApContainerService struct {
	DbInfo             *gorm.DB
	InfluxDbInfo       influxDb.Client
	Databases          models.Databases
	CloudFoundryClient *cfclient.Client
}

func GetApContainerService(DbInfo *gorm.DB, InfluxDbInfo influxDb.Client, Databases models.Databases, CloudFoundryClient *cfclient.Client) *ApContainerService {
	return &ApContainerService{
		DbInfo:             DbInfo,
		InfluxDbInfo:       InfluxDbInfo,
		Databases:          Databases,
		CloudFoundryClient: CloudFoundryClient,
	}
}

func (ap *ApContainerService) GetCellInfo() ([]models.CellInfo, error) {
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbInfo, ap.Databases, ap.CloudFoundryClient).GetCellInfo()
	if err != nil {
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetZoneInfo() ([]models.ZoneInfo, error) {
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbInfo, ap.Databases, ap.CloudFoundryClient).GetZoneInfo()
	if err != nil {
		return results, err
	}
	return results, nil
}
func (ap *ApContainerService) GetAppInfo() ([]models.AppInfo, error) {
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbInfo, ap.Databases, ap.CloudFoundryClient).GetAppInfo()
	if err != nil {
		return results, err
	}
	return results, nil
}

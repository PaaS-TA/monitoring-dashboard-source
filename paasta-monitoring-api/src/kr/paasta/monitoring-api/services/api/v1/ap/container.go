package ap

import (
	"github.com/cloudfoundry-community/go-cfclient"
	"gorm.io/gorm"
	AP "paasta-monitoring-api/dao/api/v1/ap"
	models "paasta-monitoring-api/models/api/v1"
)

type ApContainerService struct {
	DbInfo         *gorm.DB
	InfluxDbClient models.InfluxDbClient
	CfClient       *cfclient.Client
}

func GetApContainerService(DbInfo *gorm.DB, InfluxDbClient models.InfluxDbClient, CfClient *cfclient.Client) *ApContainerService {
	return &ApContainerService{
		DbInfo:         DbInfo,
		InfluxDbClient: InfluxDbClient,
		CfClient:       CfClient,
	}
}

func (ap *ApContainerService) GetCellInfo() ([]models.CellInfo, error) {
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetCellInfo()
	if err != nil {
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetZoneInfo() ([]models.ZoneInfo, error) {
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetZoneInfo()
	if err != nil {
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetAppInfo() ([]models.AppInfo, error) {
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetAppInfo()
	if err != nil {
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetContainerInfo() ([]models.ContainerInfo, error) {
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerInfo()
	if err != nil {
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetContainerPageOverview() (models.Overview, error) {
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerPageOverview()
	if err != nil {
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetContainerStatus() (models.Status, error) {
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerStatus()
	if err != nil {
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetCellStatus() (models.Status, error) {
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetCellStatus()
	if err != nil {
		return results, err
	}
	return results, nil
}

package ap

import (
	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/labstack/echo/v4"
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

func (ap *ApContainerService) GetZoneInfo() ([]models.ZoneInfo, error) {
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetZoneInfo()
	if err != nil {
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetCellInfo() ([]models.CellInfo, error) {
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetCellInfo()
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

func (ap *ApContainerService) GetContainerCpuUsages(ctx echo.Context) ([]map[string]interface{}, error) {
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerCpuUsages(ctx)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetContainerCpuLoads(ctx echo.Context) ([]map[string]interface{}, error) {
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerCpuLoads(ctx)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetContainerMemoryUsages(ctx echo.Context) ([]map[string]interface{}, error) {
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerMemoryUsages(ctx)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetContainerDiskUsages(ctx echo.Context) ([]map[string]interface{}, error) {
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerDiskUsages(ctx)
	if err != nil {
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetContainerNetworkBytes() (models.Status, error) {
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerNetworkBytes()
	if err != nil {
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetContainerNetworkDrops() (models.Status, error) {
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerNetworkDrops()
	if err != nil {
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetContainerNetworkErrors() (models.Status, error) {
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerNetworkErrors()
	if err != nil {
		return results, err
	}
	return results, nil
}

package ap

import (
	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
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

func (ap *ApContainerService) GetZoneInfo(ctx echo.Context) ([]models.ZoneInfo, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetZoneInfo()
	if err != nil {
		logger.Error(err)
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetCellInfo(ctx echo.Context) ([]models.CellInfo, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetCellInfo()
	if err != nil {
		logger.Error(err)
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetAppInfo(ctx echo.Context) ([]models.AppInfo, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetAppInfo()
	if err != nil {
		logger.Error(err)
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetContainerInfo(ctx echo.Context) ([]models.ContainerInfo, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerInfo()
	if err != nil {
		logger.Error(err)
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetContainerPageOverview(ctx echo.Context) (models.Overview, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerPageOverview()
	if err != nil {
		logger.Error(err)
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetContainerStatus(ctx echo.Context) (models.Status, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerStatus()
	if err != nil {
		logger.Error(err)
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetCellStatus(ctx echo.Context) (models.Status, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetCellStatus()
	if err != nil {
		logger.Error(err)
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetContainerCpuUsages(ctx echo.Context) ([]map[string]interface{}, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerCpuUsages(ctx)
	if err != nil {
		logger.Error(err)
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetContainerCpuLoads(ctx echo.Context) ([]map[string]interface{}, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerCpuLoads(ctx)
	if err != nil {
		logger.Error(err)
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetContainerMemoryUsages(ctx echo.Context) ([]map[string]interface{}, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerMemoryUsages(ctx)
	if err != nil {
		logger.Error(err)
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetContainerDiskUsages(ctx echo.Context) ([]map[string]interface{}, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerDiskUsages(ctx)
	if err != nil {
		logger.Error(err)
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetContainerNetworkBytes(ctx echo.Context) ([]map[string]interface{}, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerNetworkBytes(ctx)
	if err != nil {
		logger.Error(err)
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetContainerNetworkDrops(ctx echo.Context) ([]map[string]interface{}, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerNetworkDrops(ctx)
	if err != nil {
		logger.Error(err)
		return results, err
	}
	return results, nil
}

func (ap *ApContainerService) GetContainerNetworkErrors(ctx echo.Context) ([]map[string]interface{}, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)
	results, err := AP.GetApContainerDao(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerNetworkErrors(ctx)
	if err != nil {
		logger.Error(err)
		return results, err
	}
	return results, nil
}
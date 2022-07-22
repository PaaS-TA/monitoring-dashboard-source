package ap

import (
	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	"paasta-monitoring-api/connections"
	models "paasta-monitoring-api/models/api/v1"
	AP "paasta-monitoring-api/services/api/v1/ap"
)

type ApContainerController struct {
	DbInfo         *gorm.DB
	InfluxDbClient models.InfluxDbClient
	CfClient       *cfclient.Client
}

func GetApContainerController(conn connections.Connections) *ApContainerController {
	return &ApContainerController{
		DbInfo:         conn.DbInfo,
		InfluxDbClient: conn.InfluxDbClient,
		CfClient:       conn.CfClient,
	}
}

// GetZoneInfo
//  @Tags         AP
//  @Summary      Zone 정보 가져오기
//  @Description  Zone 기준 현재 배포된 Diego-Cell VM 갯수를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.ZoneInfo}
//  @Router       /api/v1/ap/container/zone [get]
func (ap *ApContainerController) GetZoneInfo(ctx echo.Context) error {
	results, err := AP.GetApContainerService(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetZoneInfo()
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get zone info.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get zone info.", results)
	return nil
}

// GetCellInfo
//  @Tags         AP
//  @Summary      Cell 정보 가져오기
//  @Description  현재 배포된 모든 Diego-Cell VM의 상세정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.CellInfo}
//  @Router       /api/v1/ap/container/cell [get]
func (ap *ApContainerController) GetCellInfo(ctx echo.Context) error {
	results, err := AP.GetApContainerService(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetCellInfo()
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get cell info.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get cell info.", results)
	return nil
}

// GetAppInfo
//  @Tags         AP
//  @Summary      App 정보 가져오기
//  @Description  현재 배포된 모든 App의 상세정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.AppInfo}
//  @Router       /api/v1/ap/container/app [get]
func (ap *ApContainerController) GetAppInfo(ctx echo.Context) error {
	results, err := AP.GetApContainerService(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetAppInfo()
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get apps info.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get apps info.", results)
	return nil
}

// GetContainerInfo
//  @Tags         AP
//  @Summary      Container 정보 가져오기
//  @Description  현재 배포된 모든 App 당 생성된 Container 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.ContainerInfo}
//  @Router       /api/v1/ap/container/container [get]
func (ap *ApContainerController) GetContainerInfo(ctx echo.Context) error {
	results, err := AP.GetApContainerService(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerInfo()
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get containers info.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get containers info.", results)
	return nil
}

// GetContainerPageOverview
//  @Tags         AP
//  @Summary      Container 페이지 Overview 정보 가져오기
//  @Description  Container 페이지를 위한 Overview 정보를 가져온다.
//  @Description  Zone - Cell - App - Container 구조를 매핑하여 보여준다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.Overview}
//  @Router       /api/v1/ap/container/overview [get]
func (ap *ApContainerController) GetContainerPageOverview(ctx echo.Context) error {
	results, err := AP.GetApContainerService(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerPageOverview()
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get container page overview.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get container page overview.", results)
	return nil
}

// GetContainerStatus
//  @Tags         AP
//  @Summary      Container Status 정보 가져오기
//  @Description  Container의 Status 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.Status}
//  @Router       /api/v1/ap/container/container/status [get]
func (ap *ApContainerController) GetContainerStatus(ctx echo.Context) error {
	results, err := AP.GetApContainerService(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerStatus()
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get container status.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get container status.", results)
	return nil
}

// GetCellStatus
//  @Tags         AP
//  @Summary      Cell Status 정보 가져오기
//  @Description  Diego-Cell VM의 Status 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.Status}
//  @Router       /api/v1/ap/container/cell/status [get]
func (ap *ApContainerController) GetCellStatus(ctx echo.Context) error {
	results, err := AP.GetApContainerService(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetCellStatus()
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get cell status.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get cell status.", results)
	return nil
}

// GetContainerCpuUsages
//  @Tags         AP
//  @Summary      Container CPU Usages 정보 가져오기
//  @Description  Container CPU Usages 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        id                path   string  true  "Container ID"        example(10.255.116.231-10.200.1.132)
//  @Param        defaultTimeRange  query  string  true  "Default Time Range"  example(15m)
//  @Param        groupBy           query  string  true  "Group By"            example(1m)
//  @Router       /api/v1/ap/container/container/cpu/{id}/usages [get]
func (ap *ApContainerController) GetContainerCpuUsages(ctx echo.Context) error {
	results, err := AP.GetApContainerService(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerCpuUsages(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Container CPU usages.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Container CPU usages.", results)
	return nil
}

// GetContainerCpuLoads
//  @Tags         AP
//  @Summary      Container CPU Loads 정보 가져오기
//  @Description  Container CPU Loads 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        id                path   string  true  "Container ID"        example(10.255.116.231-10.200.1.132)
//  @Param        defaultTimeRange  query  string  true  "Default Time Range"  example(15m)
//  @Param        groupBy           query  string  true  "Group By"            example(1m)
//  @Router       /api/v1/ap/container/container/cpu/{id}/loads [get]
func (ap *ApContainerController) GetContainerCpuLoads(ctx echo.Context) error {
	results, err := AP.GetApContainerService(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerCpuLoads(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Container CPU Loads.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Container CPU Loads.", results)
	return nil
}

// GetContainerMemoryUsages
//  @Tags         AP
//  @Summary      Container Memory Usages 정보 가져오기
//  @Description  Container Memory Usages 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        id                path   string  true  "Container ID"        example(10.255.116.231-10.200.1.132)
//  @Param        defaultTimeRange  query  string  true  "Default Time Range"  example(15m)
//  @Param        groupBy           query  string  true  "Group By"            example(1m)
//  @Router       /api/v1/ap/container/container/memory/{id}/usages [get]
func (ap *ApContainerController) GetContainerMemoryUsages(ctx echo.Context) error {
	results, err := AP.GetApContainerService(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerMemoryUsages(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Container Memory usages.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Container Memory usages.", results)
	return nil
}

// GetContainerDiskUsages
//  @Tags         AP
//  @Summary      Container Disk Usages 정보 가져오기
//  @Description  Container Disk Usages 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        id                path   string  true  "Container ID"        example(10.255.116.231-10.200.1.132)
//  @Param        defaultTimeRange  query  string  true  "Default Time Range"  example(15m)
//  @Param        groupBy           query  string  true  "Group By"            example(1m)
//  @Router       /api/v1/ap/container/container/disk/{id}/usages [get]
func (ap *ApContainerController) GetContainerDiskUsages(ctx echo.Context) error {
	results, err := AP.GetApContainerService(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerDiskUsages(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Container Disk usages.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Container Disk usages.", results)
	return nil
}

// GetContainerNetworkBytes
//  @Tags         AP
//  @Summary      Container Network Bytes 정보 가져오기
//  @Description  Container Network Bytes 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        id                path   string  true  "Container ID"        example(10.255.116.231-10.200.1.132)
//  @Param        defaultTimeRange  query  string  true  "Default Time Range"  example(15m)
//  @Param        groupBy           query  string  true  "Group By"            example(1m)
//  @Router       /api/v1/ap/container/container/network/{id}/bytes [get]
func (ap *ApContainerController) GetContainerNetworkBytes(ctx echo.Context) error {
	results, err := AP.GetApContainerService(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerNetworkBytes(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Container Network Bytes usages by time.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Container Network Bytes usages by time.", results)
	return nil
}

// GetContainerNetworkDrops
//  @Tags         AP
//  @Summary      Container Network Drops 정보 가져오기
//  @Description  Container Network Drops 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        id                path   string  true  "Container ID"        example(10.255.116.231-10.200.1.132)
//  @Param        defaultTimeRange  query  string  true  "Default Time Range"  example(15m)
//  @Param        groupBy           query  string  true  "Group By"            example(1m)
//  @Router       /api/v1/ap/container/container/network/{id}/drops [get]
func (ap *ApContainerController) GetContainerNetworkDrops(ctx echo.Context) error {
	results, err := AP.GetApContainerService(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerNetworkDrops(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Container Network Drops usages by time.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Container Network Drops usages by time.", results)
	return nil
}

// GetContainerNetworkErrors
//  @Tags         AP
//  @Summary      Container Network Errors 정보 가져오기
//  @Description  Container Network Errors 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        id                path   string  true  "Container ID"        example(10.255.116.231-10.200.1.132)
//  @Param        defaultTimeRange  query  string  true  "Default Time Range"  example(15m)
//  @Param        groupBy           query  string  true  "Group By"            example(1m)
//  @Router       /api/v1/ap/container/container/network/{id}/errors [get]
func (ap *ApContainerController) GetContainerNetworkErrors(ctx echo.Context) error {
	results, err := AP.GetApContainerService(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerNetworkErrors(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Container Network Errors usages by time.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Container Network Errors usages by time.", results)
	return nil
}

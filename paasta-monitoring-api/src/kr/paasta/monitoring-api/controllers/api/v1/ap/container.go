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
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      Zone 정보 가져오기
//  @Description  Zone 정보를 가져온다.
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
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      Cell 정보 가져오기
//  @Description  Cell 정보를 가져온다.
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
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      App 정보 가져오기
//  @Description  App 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=cfclient.App}
//  @Router       /api/v1/ap/container/zone [get]
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
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      App 정보 가져오기
//  @Description  App 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=cfclient.App}
//  @Router       /api/v1/ap/container/zone [get]
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
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      Container 페이지 Overview 정보 가져오기
//  @Description  Container 페이지의 Overview 정보를 가져온다.
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
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      Container의 Status 정보 가져오기
//  @Description  Container의 Status 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.Status}
//  @Router       /api/v1/ap/container/overview [get]
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
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      Cell(Diego-Cell VM)의 Status 정보 가져오기
//  @Description  Cell(Diego-Cell VM)의 Status 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.Overview}
//  @Router       /api/v1/ap/container/overview [get]
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
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      Container CPU Usages 정보 가져오기
//  @Description  Container CPU Usages 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.Overview}
//  @Router       /api/v1/ap/container/overview [get]
func (ap *ApContainerController) GetContainerCpuUsages(ctx echo.Context) error {
	results, err := AP.GetApContainerService(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerCpuUsages(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get cell status.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get cell status.", results)
	return nil
}

// GetContainerCpuLoads
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      Container CPU Loads 정보 가져오기
//  @Description  Container CPU Loads 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.Overview}
//  @Router       /api/v1/ap/container/overview [get]
func (ap *ApContainerController) GetContainerCpuLoads(ctx echo.Context) error {
	results, err := AP.GetApContainerService(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerCpuLoads(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get cell status.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get cell status.", results)
	return nil
}

// GetContainerMemoryUsages
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      Container Memory Usages 정보 가져오기
//  @Description  Container Memory Usages 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.Overview}
//  @Router       /api/v1/ap/container/overview [get]
func (ap *ApContainerController) GetContainerMemoryUsages(ctx echo.Context) error {
	results, err := AP.GetApContainerService(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerMemoryUsages(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get cell status.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get cell status.", results)
	return nil
}

// GetContainerDiskUsages
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      Container Disk Usages 정보 가져오기
//  @Description  Container Disk Usages 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.Overview}
//  @Router       /api/v1/ap/container/overview [get]
func (ap *ApContainerController) GetContainerDiskUsages(ctx echo.Context) error {
	results, err := AP.GetApContainerService(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerDiskUsages(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get cell status.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get cell status.", results)
	return nil
}

// GetContainerNetworkBytes
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      Container Network Bytes 정보 가져오기
//  @Description  Container Network Bytes 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.Overview}
//  @Router       /api/v1/ap/container/overview [get]
func (ap *ApContainerController) GetContainerNetworkBytes(ctx echo.Context) error {
	results, err := AP.GetApContainerService(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerNetworkBytes()
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get cell status.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get cell status.", results)
	return nil
}

// GetContainerNetworkDrops
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      Container Network Drops 정보 가져오기
//  @Description  Container Network Drops 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.Overview}
//  @Router       /api/v1/ap/container/overview [get]
func (ap *ApContainerController) GetContainerNetworkDrops(ctx echo.Context) error {
	results, err := AP.GetApContainerService(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerNetworkDrops()
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get cell status.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get cell status.", results)
	return nil
}

// GetContainerNetworkErrors
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      Container Newtwork Errors 정보 가져오기
//  @Description  Container Newtwork Errors 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.Overview}
//  @Router       /api/v1/ap/container/overview [get]
func (ap *ApContainerController) GetContainerNetworkErrors(ctx echo.Context) error {
	results, err := AP.GetApContainerService(ap.DbInfo, ap.InfluxDbClient, ap.CfClient).GetContainerNetworkErrors()
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get cell status.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get cell status.", results)
	return nil
}

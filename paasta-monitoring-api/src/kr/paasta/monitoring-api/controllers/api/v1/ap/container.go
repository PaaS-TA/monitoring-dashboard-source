package ap

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	"paasta-monitoring-api/connections"
	AP "paasta-monitoring-api/services/api/v1/ap"
)

type ApContainerController struct {
	DbInfo *gorm.DB
}

func GetApContainerController(conn connections.Connections) *ApContainerController {
	return &ApContainerController{
		DbInfo: conn.DbInfo,
	}
}

// GetCellInfo
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      셀 정보 가져오기
//  @Description  셀 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.CellInfo}
//  @Router       /api/v1/ap/alarm/status [get]
func (ap *ApContainerController) GetCellInfo(c echo.Context) error {
	results, err := AP.GetApContainerService(ap.DbInfo).GetCellInfo()
	if err != nil {
		apiHelpers.Respond(c, http.StatusBadRequest, "Failed to get cell info.", err.Error())
		return err
	}

	apiHelpers.Respond(c, http.StatusOK, "Succeeded to get cell info.", results)
	return nil
}

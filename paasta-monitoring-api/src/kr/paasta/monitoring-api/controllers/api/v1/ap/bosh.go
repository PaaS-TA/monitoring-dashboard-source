package ap

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	"paasta-monitoring-api/connections"
	models "paasta-monitoring-api/models/api/v1"
	AP "paasta-monitoring-api/services/api/v1/ap"
)

type BoshController struct {
	DbInfo         *gorm.DB
	InfluxDbClient models.InfluxDbClient
	BoshInfoList   []models.Bosh
}

func GetBoshController(conn connections.Connections) *BoshController {
	return &BoshController{
		DbInfo:         conn.DbInfo,
		InfluxDbClient: conn.InfluxDbClient,
		BoshInfoList:   conn.BoshInfoList,
	}
}

// GetBoshInfoList
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      Bosh의 목록을 가져온다.
//  @Description  Bosh의 목록을 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.Bosh}
//  @Router       /api/v1/bosh [get]
func (b *BoshController) GetBoshInfoList(c echo.Context) (err error) {
	results, err := AP.GetApBoshService(b.DbInfo, b.InfluxDbClient, b.BoshInfoList).GetBoshInfoList()
	if err != nil {
		apiHelpers.Respond(c, http.StatusInternalServerError, err.Error(), nil)
		return err
	}
	apiHelpers.Respond(c, http.StatusOK, "Success to get Bosh Info List", results)
	return nil
}

// GetBoshOverview
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      Bosh Overview 정보를 가져온다.
//  @Description  Bosh Overview 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.BoshOverview}
//  @Router       /api/v1/ap/bosh/overview [get]
func (b *BoshController) GetBoshOverview(c echo.Context) (err error) {
	results, err := AP.GetApBoshService(b.DbInfo, b.InfluxDbClient, b.BoshInfoList).GetBoshOverview(c)
	if err != nil {
		apiHelpers.Respond(c, http.StatusInternalServerError, err.Error(), nil)
		return err
	}
	apiHelpers.Respond(c, http.StatusOK, "Success to get Bosh Overview", results)
	return nil
}

// GetBoshSummary
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      Bosh Summary 정보를 가져온다.
//  @Description  Bosh Summary 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.BoshSummary}
//  @Router       /api/v1/ap/bosh/summary [get]
func (b *BoshController) GetBoshSummary(c echo.Context) (err error) {
	results, err := AP.GetApBoshService(b.DbInfo, b.InfluxDbClient, b.BoshInfoList).GetBoshSummary(c)
	if err != nil {
		apiHelpers.Respond(c, http.StatusInternalServerError, err.Error(), nil)
		return err
	}
	apiHelpers.Respond(c, http.StatusOK, "Success to get Bosh Summary", results)
	return nil
}

// GetBoshProcessByMemory
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      Bosh의 프로세스 목록을 가져온다.
//  @Description  Bosh의 프로세스 목록을 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        uuid  query     string  false  "Bosh의 프로세스 목록 조회시 Bosh ID를 주입한다."
//  @Success      200   {object}  apiHelpers.BasicResponseForm{responseInfo=v1.BoshProcess}
//  @Router       /api/v1/ap/bosh/process [get]
func (b *BoshController) GetBoshProcessByMemory(ctx echo.Context) (err error) {
	// Bosh Process 정보 조회
	results, err := AP.GetApBoshService(b.DbInfo, b.InfluxDbClient, b.BoshInfoList).GetBoshProcessByMemory(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusInternalServerError, err.Error(), nil)
		return err
	}
	apiHelpers.Respond(ctx, http.StatusOK, "Success to get Bosh Process By Memory", results)
	return nil
}

// GetBoshChart
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      Bosh의 차트 정보를 가져온다.
//  @Description  Bosh의 차트 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        id   query     string  false  "Bosh의 차트 정보 조회시 Bosh ID를 주입한다."
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.BoshChart}
//  @Router       /api/v1/bosh/Chart [get]
func (b *BoshController) GetBoshChart(ctx echo.Context) (err error) {
	results, err := AP.GetApBoshService(b.DbInfo, b.InfluxDbClient, b.BoshInfoList).GetBoshChart(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusInternalServerError, err.Error(), nil)
		return err
	}
	apiHelpers.Respond(ctx, http.StatusOK, "Success to get Bosh Chart", results)
	return nil
}

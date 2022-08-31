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
//  @tags         AP
//  @Summary      BOSH 목록 가져오기
//  @Description  BOSH 목록을 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.Bosh}
//  @Router       /api/v1/ap/bosh [get]
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
//  @tags         AP
//  @Summary      BOSH Overview 정보 가져오기
//  @Description  BOSH Overview 정보를 가져온다.
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
//  @tags         AP
//  @Summary      BOSH Summary 정보 가져오기
//  @Description  BOSH Summary 정보를 가져온다.
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
//  @tags         AP
//  @Summary      BOSH 프로세스 목록 가져오기
//  @Description  BOSH 프로세스 목록을 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        uuid  query     string  true  "BOSH의 UUID를 주입한다."  example(36dd3d08-5198-42b6-4130-d0c04479236f)
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
//  @tags         AP
//  @Summary      BOSH 차트 정보 가져오기
//  @Description  BOSH 차트 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        uuid              path      string  true   "BOSH의 UUID를 주입한다."  example(36dd3d08-5198-42b6-4130-d0c04479236f)
//  @Param        defaultTimeRange  query     string  true   "기본 시간 범위를 주입한다."    example(15m)
//  @Param        timeRangeFrom     query     string  false  "시작 시간을 주입한다."
//  @Param        timeRangeTo       query     string  false  "종료 시간을 주입한다."
//  @Param        groupBy           query     string  true   "시간 단위를 주입한다."  example(1m)
//  @Success      200               {object}  apiHelpers.BasicResponseForm{responseInfo=v1.BoshChart}
//  @Router       /api/v1/ap/bosh/chart/{uuid} [get]
func (b *BoshController) GetBoshChart(ctx echo.Context) (err error) {
	results, err := AP.GetApBoshService(b.DbInfo, b.InfluxDbClient, b.BoshInfoList).GetBoshChart(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusInternalServerError, err.Error(), nil)
		return err
	}
	apiHelpers.Respond(ctx, http.StatusOK, "Success to get Bosh Chart", results)
	return nil
}

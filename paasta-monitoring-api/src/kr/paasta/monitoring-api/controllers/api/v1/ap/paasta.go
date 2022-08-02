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

type PaastaController struct {
	DbInfo         *gorm.DB
	InfluxDbClient models.InfluxDbClient
	BoshInfoList   []models.Bosh
}

func GetPaastaController(conn connections.Connections) *PaastaController {
	return &PaastaController{
		DbInfo:         conn.DbInfo,
		InfluxDbClient: conn.InfluxDbClient,
		BoshInfoList:   conn.BoshInfoList,
	}
}

// GetPaastaInfoList
//  @tags         AP
//  @Summary      PaaS-TA VM 목록 가져오기
//  @Description  PaaS-TA VM 목록을 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.Paasta}
//  @Router       /api/v1/ap/paasta [get]
func (p *PaastaController) GetPaastaInfoList(ctx echo.Context) (err error) {
	results, err := AP.GetApPaastaService(p.DbInfo, p.InfluxDbClient).GetPaastaInfoList(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusInternalServerError, err.Error(), nil)
		return err
	}
	apiHelpers.Respond(ctx, http.StatusOK, "Success to get Paasta Info List", results)
	return nil
}

// GetPaastaOverview
//  @tags         AP
//  @Summary      PaaS-TA Overview 정보 가져오기
//  @Description  PaaS-TA Overview 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.PaastaOverview}
//  @Router       /api/v1/ap/paasta/overview [get]
func (p *PaastaController) GetPaastaOverview(ctx echo.Context) (err error) {
	results, err := AP.GetApPaastaService(p.DbInfo, p.InfluxDbClient).GetPaastaOverview(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusInternalServerError, err.Error(), nil)
		return err
	}
	apiHelpers.Respond(ctx, http.StatusOK, "Success to get Paasta Overview", results)
	return nil
}

// GetPaastaSummary
//  @tags         AP
//  @Summary      PaaS-TA Summary 정보 가져오기
//  @Description  PaaS-TA Summary 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.PaastaSummary}
//  @Router       /api/v1/ap/paasta/summary [get]
func (p *PaastaController) GetPaastaSummary(ctx echo.Context) (err error) {
	var paastaRequest models.PaastaRequest
	results, err := AP.GetApPaastaService(p.DbInfo, p.InfluxDbClient).GetPaastaSummary(paastaRequest)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusInternalServerError, err.Error(), nil)
		return err
	}
	apiHelpers.Respond(ctx, http.StatusOK, "Success to get Bosh Summary", results)
	return nil
}

// GetPaastaProcessByMemory
//  @tags         AP
//  @Summary      PaaS-TA VM 별 프로세스 목록 가져오기
//  @Description  PaaS-TA VM 별 프로세스 목록을 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        uuid  query     string  false  "Paasta의 프로세스 목록 조회시 VM ID를 주입한다."  example(f1db5cd8-3e5b-4980-966f-9fa88d8d85fd)
//  @Success      200   {object}  apiHelpers.BasicResponseForm{responseInfo=v1.PaastaProcess}
//  @Router       /api/v1/ap/paasta/process [get]
func (p *PaastaController) GetPaastaProcessByMemory(ctx echo.Context) (err error) {
	var paastaProcess models.PaastaProcess
	paastaProcess.UUID = ctx.QueryParam("uuid")

	// Paasta Process 정보 조회
	results, err := AP.GetApPaastaService(p.DbInfo, p.InfluxDbClient).GetPaastaProcessByMemory(paastaProcess)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusInternalServerError, err.Error(), nil)
		return err
	}
	apiHelpers.Respond(ctx, http.StatusOK, "Success to get Paasta Process By Memory", results)
	return nil
}

// GetPaastaChart
//  @tags         AP
//  @Summary      PaaS-TA VM 별 차트 정보 가져오기
//  @Description  PaaS-TA VM 별 차트 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        uuid              path      string  true   "PaaS-TA VM 별 차트 정보 조회시 VM UUID를 주입한다."   enums(f1db5cd8-3e5b-4980-966f-9fa88d8d85fd, 644ce3f1-c758-42ae-8d74-8193e28839fe)
//  @Param        defaultTimeRange  query     string  false  "PaaS-TA VM 별 차트 정보 조회시 기본 시간 범위를 주입한다."  example(15m)
//  @Param        timeRangeFrom     query     string  false  "PaaS-TA VM 별 차트 정보 조회시 시작 시간을 주입한다."
//  @Param        timeRangeTo       query     string  false  "PaaS-TA VM 별 차트 정보 조회시 종료 시간을 주입한다."
//  @Param        groupBy           query     string  false  "PaaS-TA VM 별 차트 정보 조회시 시간 단위를 주입한다."  example(1m)
//  @Success      200               {object}  apiHelpers.BasicResponseForm{responseInfo=v1.PaastaChart}
//  @Router       /api/v1/ap/paasta/chart/{uuid} [get]
func (p *PaastaController) GetPaastaChart(ctx echo.Context) (err error) {
	var paastaChart models.PaastaChart
	paastaChart.UUID = ctx.Param("uuid")
	paastaChart.DefaultTimeRange = ctx.QueryParam("defaultTimeRange")
	paastaChart.TimeRangeFrom = ctx.QueryParam("timeRangeFrom")
	paastaChart.TimeRangeTo = ctx.QueryParam("timeRangeTo")
	paastaChart.GroupBy = ctx.QueryParam("groupBy")

	results, err := AP.GetApPaastaService(p.DbInfo, p.InfluxDbClient).GetPaastaChart(paastaChart)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusInternalServerError, err.Error(), nil)
		return err
	}
	apiHelpers.Respond(ctx, http.StatusOK, "Success to get Paasta Chart", results)
	return nil
}

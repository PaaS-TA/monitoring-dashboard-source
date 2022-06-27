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
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      PaaS-TA Core의 목록을 가져온다.
//  @Description  PaaS-TA Core의 목록을 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.Paasta}
//  @Router       /api/v1/paasta [get]
func (p *PaastaController) GetPaastaInfoList(c echo.Context) (err error) {
	results, err := AP.GetApPaastaService(p.DbInfo, p.InfluxDbClient).GetPaastaInfoList()
	if err != nil {
		apiHelpers.Respond(c, http.StatusInternalServerError, err.Error(), nil)
		return err
	}
	apiHelpers.Respond(c, http.StatusOK, "Success to get Paasta Info List", results)
	return nil
}

// GetPaastaOverview
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      PaaS-TA Overview 정보를 가져온다.
//  @Description  PaaS-TA Overview 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.PaastaOverview}
//  @Router       /api/v1/ap/paasta/overview [get]
func (p *PaastaController) GetPaastaOverview(c echo.Context) (err error) {
	results, err := AP.GetApPaastaService(p.DbInfo, p.InfluxDbClient).GetPaastaOverview(c)
	if err != nil {
		apiHelpers.Respond(c, http.StatusInternalServerError, err.Error(), nil)
		return err
	}
	apiHelpers.Respond(c, http.StatusOK, "Success to get Bosh Overview", results)
	return nil
}

// GetPaastaSummary
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      PaaS-TA Summary 정보를 가져온다.
//  @Description  PaaS-TA Summary 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.PaastaSummary}
//  @Router       /api/v1/ap/paasta/summary [get]
func (p *PaastaController) GetPaastaSummary(c echo.Context) (err error) {
	results, err := AP.GetApPaastaService(p.DbInfo, p.InfluxDbClient).GetPaastaSummary(c)
	if err != nil {
		apiHelpers.Respond(c, http.StatusInternalServerError, err.Error(), nil)
		return err
	}
	apiHelpers.Respond(c, http.StatusOK, "Success to get Bosh Summary", results)
	return nil
}

// GetPaastaProcessByMemory
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      PaaS-TA Core 별 프로세스 목록을 가져온다.
//  @Description  PaaS-TA Core 별 프로세스 목록을 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        uuid   query     string  false  "Paasta의 프로세스 목록 조회시 VM ID를 주입한다."
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.PaastaProcess}
//  @Router       /api/v1/ap/paasta/process [get]
func (p *PaastaController) GetPaastaProcessByMemory(c echo.Context) (err error) {
	var paastaProcess models.PaastaProcess
	paastaProcess.UUID = c.QueryParam("uuid")

	// Paasta Process 정보 조회
	results, err := AP.GetApPaastaService(p.DbInfo, p.InfluxDbClient).GetPaastaProcessByMemory(paastaProcess)
	if err != nil {
		apiHelpers.Respond(c, http.StatusInternalServerError, err.Error(), nil)
		return err
	}
	apiHelpers.Respond(c, http.StatusOK, "Success to get Paasta Process By Memory", results)
	return nil
}

// GetPaastaChart
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      PaaS-TA Core 별 차트 정보를 가져온다.
//  @Description  PaaS-TA Core 별 차트 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        uuid   query     string  false  "PaaS-TA Core 별 차트 정보 조회시 VM UUID를 주입한다."
//  @Param        defaultTimeRange   query     string  false  "PaaS-TA Core 별 차트 정보 조회시 기본 시간 범위 (defaultTimeRange=15m)를 주입한다."
//  @Param        timeRangeFrom   query     string  false  "PaaS-TA Core 별 차트 정보 조회시 시간 범위 시작 (timeRangeFrom=2022-06-16T10:21:39)를 주입한다."
//  @Param        timeRangeTo   query     string  false  "PaaS-TA Core 별 차트 정보 조회시 시간 범위 종료 (timeRangeTo=2022-06-16T10:21:39)를 주입한다."
//  @Param        groupBy   query     string  false  "PaaS-TA Core 별 차트 정보 조회시 그룹 (groupBy=1m)을 주입한다."
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.PaastaChart}
//  @Router       /api/v1/paasta/Chart [get]
func (p *PaastaController) GetPaastaChart(c echo.Context) (err error) {
	var boshChart models.BoshChart
	boshChart.UUID = c.Param("uuid")
	boshChart.DefaultTimeRange = c.QueryParam("defaultTimeRange")
	boshChart.TimeRangeFrom = c.QueryParam("timeRangeFrom")
	boshChart.TimeRangeTo = c.QueryParam("timeRangeTo")
	boshChart.GroupBy = c.QueryParam("groupBy")

	results, err := AP.GetApPaastaService(p.DbInfo, p.InfluxDbClient).GetPaastaChart(boshChart)
	if err != nil {
		apiHelpers.Respond(c, http.StatusInternalServerError, err.Error(), nil)
		return err
	}
	apiHelpers.Respond(c, http.StatusOK, "Success to get Bosh Chart", results)
	return nil
}

// GetPaastaLog
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      Paasta의 로그 정보를 가져온다.
//  @Description  Paasta의 로그 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        uuid   query     string  false  "PaaS-TA Core 별 로그 정보 조회시 VM UUID를 주입한다."
//  @Param        keyword   query     string  false  "PaaS-TA Core 별 로그 정보 조회시 키워드 (keyword)를 주입한다."
//  @Param        targetDate   query     string  false  "PaaS-TA Core 별 로그 정보 조회시 대상날짜 (targetDate)를 주입한다."
//  @Param        startTime   query     string  false  "PaaS-TA Core 별 로그 정보 조회시 시작시간 (startTime)를 주입한다."
//  @Param        endTime   query     string  false  "PaaS-TA Core 별 로그 정보 조회시 종료시간 (endTime)를 주입한다."
//  @Param        period   query     string  false  "PaaS-TA Core 별 로그 정보 조회시 조회기간 (period)를 주입한다."
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.PaastaLog}
//  @Router       /api/v1/paasta/log [get]
func (p *PaastaController) GetPaastaLog(c echo.Context) (err error) {
	var paastaLog models.PaastaLog
	paastaLog.UUID = c.Param("uuid")
	paastaLog.Keyword = c.QueryParam("keyword")
	paastaLog.TargetDate = c.QueryParam("targetDate")
	paastaLog.StartTime = c.QueryParam("startTime")
	paastaLog.EndTime = c.QueryParam("endTime")
	paastaLog.Period = c.QueryParam("period")

	results, err := AP.GetApPaastaService(p.DbInfo, p.InfluxDbClient).GetPaastaLog(paastaLog)
	if err != nil {
		apiHelpers.Respond(c, http.StatusInternalServerError, err.Error(), nil)
		return err
	}
	apiHelpers.Respond(c, http.StatusOK, "Success to get Bosh Log", results)
	return nil
}

package common

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	"paasta-monitoring-api/connections"
	models "paasta-monitoring-api/models/api/v1"
	service "paasta-monitoring-api/services/api/v1/common"
)

type LogSearchController struct {
	DbInfo         *gorm.DB
	InfluxDbClient models.InfluxDbClient
	BoshInfoList   []models.Bosh
}

func GetLogSearchController(conn connections.Connections) *LogSearchController {
	return &LogSearchController{
		DbInfo:         conn.DbInfo,
		InfluxDbClient: conn.InfluxDbClient,
		BoshInfoList:   conn.BoshInfoList,
	}
}

// GetLogs
//  * Annotations for Swagger *
//  @Summary      로그 정보 가져오기
//  @Description  로그 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        uuid   query     string  false  "PaaS-TA Core 별 로그 정보 조회시 VM UUID를 주입한다."
//  @Param        keyword   query     string  false  "PaaS-TA Core 별 로그 정보 조회시 키워드 (keyword)를 주입한다."
//  @Param        targetDate   query     string  false  "PaaS-TA Core 별 로그 정보 조회시 대상날짜 (targetDate)를 주입한다."
//  @Param        startTime   query     string  false  "PaaS-TA Core 별 로그 정보 조회시 시작시간 (startTime)를 주입한다."
//  @Param        endTime   query     string  false  "PaaS-TA Core 별 로그 정보 조회시 종료시간 (endTime)를 주입한다."
//  @Param        period   query     string  false  "PaaS-TA Core 별 로그 정보 조회시 조회기간 (period)를 주입한다."
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.Logs}
//  @Router       /api/v1/log [get]
func (l *LogSearchController) GetLogs(ctx echo.Context) error {
	results, err := service.GetLogSearchService(l.DbInfo, l.InfluxDbClient, l.BoshInfoList).GetLogs(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get logs.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get logs.", results)
	return nil
}

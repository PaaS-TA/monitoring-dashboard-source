package ap

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"paasta-monitoring-api/connections"
	v1 "paasta-monitoring-api/models/api/v1"
)

type BoshController struct {
	DbInfo *gorm.DB
}

func GetBoshController(conn connections.Connections) *BoshController {
	return &BoshController{
		DbInfo: conn.DbInfo,
	}
}

// GetBoshOverview
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      Bosh의 상태 별 개수를 가져온다.
//  @Description  Bosh의 상태 별 개수를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.BoshSummary}
//  @Router       /api/v1/bosh/overview [get]
func (a *BoshController) GetBoshOverview(c echo.Context) (err error) {
	var BoshSummary []v1.BoshSummary
	fmt.Println(BoshSummary)
	return nil
}

// GetBoshList
//  * Annotations for Swagger *
//  @Summary      Bosh의 목록을 가져온다.
//  @Description  Bosh의 목록을 가져온다.
//  @tags         AP
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.Bosh}
//  @Router       /api/v1/bosh [get]
func (a *BoshController) GetBoshList(c echo.Context) (err error) {
	return nil
}

// GetBoshProcessList
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      Bosh의 프로세스 목록을 가져온다.
//  @Description  Bosh의 프로세스 목록을 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        id   query     string  false  "Bosh의 프로세스 목록 조회시 Bosh ID를 주입한다."
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.BoshProcess}
//  @Router       /api/v1/bosh/process [get]
func (a *BoshController) GetBoshProcessList(c echo.Context) (err error) {
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
func (a *BoshController) GetBoshChart(c echo.Context) (err error) {
	return nil
}

// GetBoshLog
//  * Annotations for Swagger *
//  @tags         AP
//  @Summary      Bosh의 로그 정보를 가져온다.
//  @Description  Bosh의 로그 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        id   query     string  false  "Bosh의 로그 정보 조회시 Bosh ID를 주입한다."
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=v1.BoshLog}
//  @Router       /api/v1/bosh/log [get]
func (a *BoshController) GetBoshLog(c echo.Context) (err error) {
	return nil
}

package v1

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	"paasta-monitoring-api/connections"
	v1service "paasta-monitoring-api/services/api/v1"
)

type UserController struct {
	DbInfo *gorm.DB
}

func GetUserController(conn connections.Connections) *UserController {
	return &UserController{
		DbInfo: conn.DbInfo,
	}
}

// GetMember
//  @Summary      전체 또는 단일 유저 정보 가져오기
//  @Description  전체 또는 단일 유저 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        userId  query     string  false  "단일 사용자 정보 조회시 유저명을 주입한다."
//  @Success      200     {object}  apiHelpers.BasicResponseForm{responseInfo=MemberInfos}
//  @Router       /api/v1/members [get]
func (controller *UserController) GetMember(ctx echo.Context) error {
	results, err := v1service.GetUserService(controller.DbInfo).GetMember(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get user.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get user.", results)
	return nil
}

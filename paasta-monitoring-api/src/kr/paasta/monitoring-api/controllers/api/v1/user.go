package v1

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	"paasta-monitoring-api/connections"
	"paasta-monitoring-api/models/api/v1"
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

// GetUsers
//  * Annotations for Swagger *
//  @Summary      전체 또는 단일 유저 정보 가져오기
//  @Description  전체 또는 단일 유저 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        username  query     string  false  "단일 사용자 정보 조회시 유저명을 주입한다."
//  @Success      200       {object}  apiHelpers.BasicResponseForm{responseInfo=UserInfo}
//  @Router       /api/v1/users [get]
func (a *UserController) GetUsers(c echo.Context) (err error) {
	var users []v1.UserInfo
	var request v1.UserInfo
	qParams := c.QueryParams()
	if val, ok := qParams["username"]; ok {
		request.Username = val[0]
	}

	if request.Username != "" {
		users, err = v1service.GetUserService(a.DbInfo).GetUser(request, c)
		if err != nil {
			apiHelpers.Respond(c, http.StatusInternalServerError, err.Error(), nil)
			return err
		}
		// 단일 사용자 정보를 전달한다.
		apiHelpers.Respond(c, http.StatusOK, "Success to get user", users)
	} else {
		users, err = v1service.GetUserService(a.DbInfo).GetUsers(request, c)
		if err != nil {
			apiHelpers.Respond(c, http.StatusInternalServerError, err.Error(), nil)
			return err
		}
		// 전체 사용자 정보를 전달한다.
		apiHelpers.Respond(c, http.StatusOK, "Success to get all users", users)
	}
	return nil
}


func (controller *UserController) GetMember(ctx echo.Context) error {
	results, err := v1service.GetUserService(controller.DbInfo).GetMember(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to update alarm policy.", err.Error())
		return err
	}

	apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to update alarms policy.", results)
	return nil
}
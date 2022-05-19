package v1

import (
	"GoEchoProject/apiHelpers"
	"GoEchoProject/connections"
	"GoEchoProject/helpers"
	"GoEchoProject/models/api/v1"
	v1service "GoEchoProject/services/api/v1"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"net/http"
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
//  @Summary      [테스트] 전체 유저 정보 가져오기
//  @Description  [테스트] 전체 유저 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm{responseInfo=UserInfo}
//  @Router       /api/test/users [get]
func (a *UserController) GetUsers(c echo.Context) (err error) {
	var apiRequest v1.UserInfo
	err = helpers.BindRequestAndCheckValid(c, &apiRequest)
	if err != nil {
		apiHelpers.Respond(c, http.StatusBadRequest, "Invalid JSON provided, please check the request JSON", nil)
		return err
	}

	// User의 GetUsers를 호출한다.
	users, err := v1service.GetUserService(a.DbInfo).GetUsers(apiRequest, c)
	if err != nil {
		apiHelpers.Respond(c, http.StatusInternalServerError, err.Error(), nil)
		return err
	}

	// 사용자 정보를 전달한다.
	apiHelpers.Respond(c, http.StatusOK, "Success to get all users", users)
	return nil
}

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

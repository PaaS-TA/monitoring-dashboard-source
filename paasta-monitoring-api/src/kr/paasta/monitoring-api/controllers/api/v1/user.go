package v1

import (
	"GoEchoProject/apiHelpers"
	"GoEchoProject/connections"
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
	/* Request Body Data Mapping */
	var apiRequest v1.UserInfo // -> &추가
	//apiRequest := new(models.UserInfo) // &제거
	if err = c.Bind(&apiRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return nil
	}

	// User의 GetUsers를 호출한다.
	users, err := v1service.GetUserService(a.DbInfo).GetUsers(apiRequest, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return nil
	}

	// 사용자 정보를 전달한다.
	apiHelpers.Respond(c, http.StatusOK, "Success to get all users", users)
	return nil
}

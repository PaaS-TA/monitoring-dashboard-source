package v1

import (
	"GoEchoProject/connections"
	"GoEchoProject/models"
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
	var apiRequest models.UserInfo // -> &추가
	//apiRequest := new(models.UserInfo) // &제거
	if err = c.Bind(&apiRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return nil
	}

	// Authentication의 CreateToken을 호출한다.
	users, err := v1service.GetUserService(a.DbInfo).GetUsers(apiRequest, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return nil
	}

	// 토근을 발급한다.
	c.JSON(http.StatusOK, users)
	return nil
}

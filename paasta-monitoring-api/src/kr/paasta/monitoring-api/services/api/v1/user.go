package v1service

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"github.com/labstack/echo/v4"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	"paasta-monitoring-api/helpers"
	"paasta-monitoring-api/models/api/v1"

	dao "paasta-monitoring-api/dao/api/v1"
)

//Gorm Object Struct
type UserService struct {
	db *gorm.DB
}

func GetUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}


func (service *UserService) GetMember(ctx echo.Context) ([]v1.MemberInfos, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	userId := ctx.QueryParam("userId")
	params := v1.MemberInfos {
		UserId : userId,
	}

	validationErr := helpers.CheckValid(params)
	if validationErr != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Invalid JSON provided, please check the request JSON", validationErr.Error())
		return nil, validationErr
	}

	members, err := dao.GetUserDao(service.db).GetMember(params)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return members, nil
}


func (h *UserService) GetUsers(request v1.UserInfo, c echo.Context) ([]v1.UserInfo, error) {
	users, err := dao.GetUserDao(h.db).GetUsers(request, c)
	if err != nil {
		fmt.Println(err.Error())
		return users, err
	}
	return users, nil
}

func (h *UserService) GetUser(request v1.UserInfo, c echo.Context) ([]v1.UserInfo, error) {
	users, err := dao.GetUserDao(h.db).GetUser(request, c)
	if err != nil {
		fmt.Println(err.Error())
		return users, err
	}
	return users, nil
}

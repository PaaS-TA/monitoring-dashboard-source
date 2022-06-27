package v1service

import (
	"fmt"
	"gorm.io/gorm"
	"github.com/labstack/echo/v4"
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

func (service *UserService) GetMember(params v1.MemberInfos) ([]v1.MemberInfos, error) {
	members, err := dao.GetUserDao(service.db).GetMember(params)
	if err != nil {
		fmt.Println(err.Error())
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

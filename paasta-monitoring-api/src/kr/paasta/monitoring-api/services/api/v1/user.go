package v1service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"paasta-monitoring-api/models/api/v1"

	dao "paasta-monitoring-api/dao/api/v1"
)

//Gorm Object Struct
type UserService struct {
	txn *gorm.DB
}

func GetUserService(txn *gorm.DB) *UserService {
	return &UserService{
		txn: txn,
	}
}

func (h *UserService) GetUsers(request v1.UserInfo, c echo.Context) ([]v1.UserInfo, error) {
	users, err := dao.GetUserDao(h.txn).GetUsers(request, c)
	if err != nil {
		fmt.Println(err.Error())
		return users, err
	}
	return users, nil
}

func (h *UserService) GetUser(request v1.UserInfo, c echo.Context) ([]v1.UserInfo, error) {
	users, err := dao.GetUserDao(h.txn).GetUser(request, c)
	if err != nil {
		fmt.Println(err.Error())
		return users, err
	}
	return users, nil
}

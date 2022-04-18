package v1service

import (
	"GoEchoProject/models"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"

	dao "GoEchoProject/dao/api/v1"
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

func (h *UserService) GetUsers(apiRequest models.UserInfo, c echo.Context) ([]models.UserInfo, error) {
	users, err := dao.GetUserDao(h.txn).GetUsers(apiRequest, c)
	if err != nil {
		fmt.Println(err.Error())
		return users, err
	}
	return users, nil
}

func (h *UserService) GetUser(apiRequest models.UserInfo, c echo.Context) ([]models.UserInfo, error) {
	users, err := dao.GetUserDao(h.txn).GetUser(apiRequest, c)
	if err != nil {
		fmt.Println(err.Error())
		return users, err
	}
	return users, nil
}

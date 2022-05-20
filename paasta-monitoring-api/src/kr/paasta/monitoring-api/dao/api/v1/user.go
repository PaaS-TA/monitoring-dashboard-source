package v1

import (
	models "GoEchoProject/models/api/v1"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

type UserDao struct {
	txn *gorm.DB
}

func GetUserDao(txn *gorm.DB) *UserDao {
	return &UserDao{
		txn: txn,
	}
}

func (u *UserDao) GetUsers(request models.UserInfo, c echo.Context) (tz []models.UserInfo, err error) {
	// 전달 받은 계정 정보로 데이터베이스에 계정이 존재하는지 확인한다. (test code)
	var t []models.UserInfo
	users := u.txn.Debug().Table("user_infos").
		Select(" * ").
		Find(&t)
	//fmt.Println(users)

	if users.Error != nil {
		fmt.Println(users.Error)
		return t, users.Error
	}
	//c.Logger().Info(users)

	return t, err
}

func (u *UserDao) GetUser(request models.UserInfo, c echo.Context) (tz []models.UserInfo, err error) {
	// 전달 받은 계정 정보로 데이터베이스에 계정이 존재하는지 확인한다. (test code)
	var t []models.UserInfo
	users := u.txn.Debug().Table("user_infos").
		Select(" * ").
		Where("username = ? ", request.Username).
		Find(&t)
	//fmt.Println(users)

	if users.Error != nil {
		fmt.Println(users.Error)
		return t, users.Error
	}
	//c.Logger().Info(users)

	return t, err
}

package v1

import (
	"fmt"
	"gorm.io/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	models "paasta-monitoring-api/models/api/v1"
)

type UserDao struct {
	db *gorm.DB
}

func GetUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (u *UserDao) GetMember(params models.MemberInfos) ([]models.MemberInfos, error) {
	var members []models.MemberInfos

	result := u.db.Where(params).Find(&members)
	if result.Error != nil {
		log.Errorf(result.Error.Error())
		return nil, result.Error
	}
	return members, nil
}


func (u *UserDao) GetUsers(request models.UserInfo, c echo.Context) (tz []models.UserInfo, err error) {
	// 전달 받은 계정 정보로 데이터베이스에 계정이 존재하는지 확인한다. (test code)
	var t []models.UserInfo
	users := u.db.Debug().Table("user_infos").
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
	users := u.db.Debug().Table("user_infos").
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


func (u *UserDao) GetMemberInfo(param models.MemberInfos, c echo.Context) (tz []models.MemberInfos, err error) {
	// 전달 받은 계정 정보로 데이터베이스에 계정이 존재하는지 확인한다. (test code)
	var t []models.MemberInfos
	users := u.db.Debug().Model(&models.MemberInfos{}).Where(param).Find(&t)
	//fmt.Println(users)

	if users.Error != nil {
		fmt.Println(users.Error)
		return t, users.Error
	}
	//c.Logger().Info(users)

	return t, err
}


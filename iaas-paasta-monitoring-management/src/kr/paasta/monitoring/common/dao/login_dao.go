package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	cm "kr/paasta/monitoring/common/model"
	"kr/paasta/monitoring/utils"
	//"strconv"
)

type LoginDao struct {
	txn *gorm.DB
}

func GetLoginDao(txn *gorm.DB) *LoginDao {
	return &LoginDao{
		txn: txn,
	}
}

//Dao
func (h *LoginDao) GetLoginMemberInfo(request cm.UserInfo , txn *gorm.DB) (cm.UserInfo, int, error) {

	fmt.Println("Get Call LoginDao GetLoginMemberInfo =====")

	pw := utils.GetSha256(request.Password)

	t := cm.UserInfo{}
	//var rowCount int
	status := txn.Debug().Table("member_infos").
		Select(" * ").
		Where("user_id = ? and user_pw = ? ", request.Username, pw).
		Find(&t)
	if status.Error != nil {
		return t, 0, status.Error
	}

	return t, 1, nil
}

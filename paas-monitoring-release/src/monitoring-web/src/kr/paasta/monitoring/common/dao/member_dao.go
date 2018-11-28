package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	cm "kr/paasta/monitoring/common/model"
	//"kr/paasta/monitoring/utils"
	"time"
	"kr/paasta/monitoring/utils"
)

type MemberDao struct {
	txn *gorm.DB
}

func GetMemberDao(txn *gorm.DB) *MemberDao {
	return &MemberDao{
		txn: txn,
	}
}



func (a *MemberDao) MemberJoinSave(request cm.UserInfo , txn *gorm.DB) error {

	pw := utils.GetSha256(request.UserPw)

	actionData := cm.MemberInfo{
		UserId: request.UserId,
		UserPw: pw,
		UserEmail: request.UserEmail,
		UserNm  : request.UserNm,
		IaasUserId: request.IaasUserId,
		IaasUserPw: request.IaasUserPw,
		PaasUserId: request.PaasUserId,
		PaasUserPw: request.PaasUserPw,
		PaasUserUseYn: request.PaasUserUseYn,
		IaasUserUseYn: request.IaasUserUseYn,
	}

	status := a.txn.Debug().Create(&actionData)

	if status.Error != nil{
		return  status.Error
	}
	return  status.Error
}

func (h *MemberDao) MemberAuthCheck(request cm.UserInfo , txn *gorm.DB) (cm.UserInfo, int, error) {

	t := cm.UserInfo{}

	status := txn.Debug().Table("member_infos").
		Select(" * ").
		Where("user_id = ? ", request.UserId).
		Find(&t)
	if status.Error != nil {
		return t, 0, status.Error
	}

	return t, 1, nil
}


func (h *MemberDao) MemberInfoView(request cm.UserInfo , txn *gorm.DB) (cm.UserInfo, int, error) {

	t := cm.UserInfo{}

	status := txn.Debug().Table("member_infos").
		Select(" * ").
		Where("user_id = ? ", request.UserId).
		Find(&t)
	if status.Error != nil {
		return t, 0, status.Error
	}

	return t, 1, nil
}

func (h *MemberDao) MemberInfoCheck(request cm.UserInfo , txn *gorm.DB) (cm.UserInfo, int, error) {

	t := cm.UserInfo{}

	status := txn.Debug().Table("member_infos").
		Select(" * ")

		if request.UserId !="" {
			status = status.Where("user_id = ? ", request.UserId)
		} else if request.UserEmail !="" {
			status = status.Where("user_email = ? ", request.UserEmail)
		}

	status = status.Find(&t)

    if status.RecordNotFound() {
		return t, 0, nil
	}else if status.Error != nil {
		return t, 0, status.Error
	}

	return t, 1, nil
}


func (h *MemberDao) MemberInfoUpdate(request cm.UserInfo , txn *gorm.DB) ( int, error) {

	if request.IaasUserUseYn !="Y" {
		request.IaasUserId = ""
		request.IaasUserPw = ""
	}

	if request.PaasUserUseYn !="Y" {
		request.PaasUserId = ""
		request.PaasUserPw = ""
	}

	status := txn.Debug().Table("member_infos").
		Where("user_id = ?  ", request.UserId).
			Updates(map[string]interface{}{
					//"user_pw":     pw,
					"user_email":  request.UserEmail,
					"user_nm":  request.UserNm,
					"iaas_user_id": request.IaasUserId,
					"iaas_user_pw": request.IaasUserPw,
					"paas_user_id": request.PaasUserId,
					"paas_user_pw": request.PaasUserPw,
					"iaas_user_use_yn": request.IaasUserUseYn,
					"paas_user_use_yn": request.PaasUserUseYn,
				    "updated_at": time.Now() })
	if status.Error != nil {
		return  0, status.Error
	}

    if request.UserPw != "" {
		pw := utils.GetSha256(request.UserPw)
		status := txn.Debug().Table("member_infos").
			Where("user_id = ?  ", request.UserId).
			Updates(map[string]interface{}{
			"user_pw":     pw,
			"updated_at": time.Now() })
		if status.Error != nil {
			return  0, status.Error
		}
	}


	return  1, nil
}

func (h *MemberDao) MemberInfoDelete(request cm.UserInfo , txn *gorm.DB) ( int, error) {

	fmt.Println("Get Call LoginDao Delete LoginMemberInfo =====")
	//var rowCount int
	status := txn.Debug().Table("member_infos").
		Where("user_id = ? ", request.UserId).Delete(&request)
	if status.Error != nil {
		return 0, status.Error
	}

	return  1, nil
}

func (h *MemberDao) MemberJoinCheckDuplicationIaasId(request cm.UserInfo , txn *gorm.DB) (cm.UserInfo, error) {
	var rows cm.UserInfo
	status := txn.Debug().Table("member_infos").Where("iaas_user_id = ?", request.IaasUserId).Limit(1).Find(&rows)
	if status.RecordNotFound() {
		return rows, nil
	}
	return rows, status.Error
}

func (h *MemberDao) MemberJoinCheckDuplicationPaasId(request cm.UserInfo , txn *gorm.DB) (cm.UserInfo, error) {
	var rows cm.UserInfo
	status := txn.Debug().Table("member_infos").Where("paas_user_id = ?", request.PaasUserId).Limit(1).Find(&rows)
	if status.RecordNotFound() {
		return rows, nil
	}
	return rows, status.Error
}
package dao

import (
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring/domain"
	"kr/paasta/monitoring/util"
	"fmt"
)

type AlarmPolicyDao struct {
	txn   *gorm.DB
}

func GetAlarmPolicyDao(txn *gorm.DB) *AlarmPolicyDao {
	return &AlarmPolicyDao{
		txn:   txn,
	}
}

//Dao
func (b AlarmPolicyDao) GetAlarmPolicyList() ([]domain.AlarmPolicyResponse, domain.ErrMessage) {

	policies := []domain.AlarmPolicyResponse{}
	status := b.txn.Debug().Table("alarm_policies").Find(&policies)
	err := util.GetError().DbCheckError(status.Error)

	if err != nil{
		fmt.Println("Error::", err )
	}
	return policies, err
}

//Dao
func (b AlarmPolicyDao) GetAlarmTargetsList() ([]domain.AlarmTargetsResponse, domain.ErrMessage) {

	targets := []domain.AlarmTargetsResponse{}
	status := b.txn.Debug().Table("alarm_targets").Find(&targets)
	err := util.GetError().DbCheckError(status.Error)

	if err != nil{
		fmt.Println("Error::", err )
	}
	return targets, err
}

/*func (b AlarmPolicyDao) CreateAlarmPolicy(request domain.AlarmPolicyRequest, txn *gorm.DB) (errMsg domain.ErrMessage)  {

	alarmPolicy := domain.AlarmPolicy{OriginType: request.OriginType, WarningThreshold: request.WarningThreshold ,
		CriticalThreshold: request.CriticalThreshold, RepeatTime: request.RepeatTime, Comment: request.Comment}

	status := txn.Debug().Create(&alarmPolicy).Table("alarm_policies")
	err := util.GetError().DbCheckError(status.Error)
	return err
}*/

func (b AlarmPolicyDao) UpdateAlarmPolicy(request domain.AlarmPolicyRequest) (errMsg domain.ErrMessage)  {

	alarmPolicy := domain.AlarmPolicyRequest{OriginType: request.OriginType, AlarmType: request.AlarmType, WarningThreshold: request.WarningThreshold, CriticalThreshold: request.CriticalThreshold,
		RepeatTime: request.RepeatTime, Comment: request.Comment}
	status := b.txn.Debug().Table("alarm_policies").Model(&alarmPolicy).Where("origin_type = ? and alarm_type = ?",
		alarmPolicy.OriginType, alarmPolicy.AlarmType ).Updates(map[string]interface{}{"warning_threshold": request.WarningThreshold, "critical_threshold":request.CriticalThreshold,
		"repeat_time": alarmPolicy.RepeatTime,  "modi_date": util.GetDBCurrentTime(), "modi_user": "test"})

	fmt.Println("STATUS::",status.Error)

	err := util.GetError().DbCheckError(status.Error)

	return err
}

func (b AlarmPolicyDao) UpdateAlarmTargets(request domain.AlarmPolicyRequest) (errMsg domain.ErrMessage)  {

	alarmTargets :=  domain.AlarmTargetsRequest	{OriginType: request.OriginType, MailAddress :request.MailAddress }
	status := b.txn.Debug().Table("alarm_targets").Model(&alarmTargets).Where("origin_type = ? ",
		alarmTargets.OriginType ).Updates(map[string]interface{}{"mail_address": alarmTargets.MailAddress, "modi_date": util.GetDBCurrentTime(), "modi_user": "test"})

	fmt.Println("STATUS::",status.Error)

	err := util.GetError().DbCheckError(status.Error)

	return err
}
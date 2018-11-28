package dao

import (
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring/paas/model"
	"kr/paasta/monitoring/paas/util"
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
func (b AlarmPolicyDao) GetAlarmPolicyList() ([]model.AlarmPolicyResponse, model.ErrMessage) {

	policies := []model.AlarmPolicyResponse{}
	status := b.txn.Debug().Table("alarm_policies").Find(&policies)
	err := util.GetError().DbCheckError(status.Error)

	if err != nil {
		fmt.Println("Error::", err )
	}
	return policies, err
}

//Dao
func (b AlarmPolicyDao) GetAlarmTargetsList() ([]model.AlarmTargetsResponse, model.ErrMessage) {

	targets := []model.AlarmTargetsResponse{}
	status := b.txn.Debug().Table("alarm_targets").Find(&targets)
	err := util.GetError().DbCheckError(status.Error)

	if err != nil {
		fmt.Println("Error::", err )
	}
	return targets, err
}

/*func (b AlarmPolicyDao) CreateAlarmPolicy(request model.AlarmPolicyRequest, txn *gorm.DB) (errMsg model.ErrMessage)  {

	alarmPolicy := model.AlarmPolicy{OriginType: request.OriginType, WarningThreshold: request.WarningThreshold ,
		CriticalThreshold: request.CriticalThreshold, RepeatTime: request.RepeatTime, Comment: request.Comment}

	status := txn.Debug().Create(&alarmPolicy).Table("alarm_policies")
	err := util.GetError().DbCheckError(status.Error)
	return err
}*/

func (b AlarmPolicyDao) UpdateAlarmPolicy(request model.AlarmPolicyRequest) (errMsg model.ErrMessage)  {

	alarmPolicy := model.AlarmPolicyRequest{
		OriginType			: request.OriginType,
		AlarmType			: request.AlarmType,
		WarningThreshold	: request.WarningThreshold,
		CriticalThreshold	: request.CriticalThreshold,
		RepeatTime			: request.RepeatTime,
		MeasureTime			: request.MeasureTime,
		Comment				: request.Comment}

	status := b.txn.Debug().Table("alarm_policies").
		Model(&alarmPolicy).
		Where("origin_type = ? and alarm_type = ?", alarmPolicy.OriginType, alarmPolicy.AlarmType ).
		Updates(map[string]interface{}{
			"warning_threshold": request.WarningThreshold,
			"critical_threshold":request.CriticalThreshold,
			"repeat_time": alarmPolicy.RepeatTime,
			"measure_time": alarmPolicy.MeasureTime,
			"modi_date": util.GetDBCurrentTime(),
			"modi_user": "test"})

	err := util.GetError().DbCheckError(status.Error)
	return err
}

func (b AlarmPolicyDao) UpdateAlarmTargets(request model.AlarmPolicyRequest) (errMsg model.ErrMessage)  {

	alarmTargets :=  model.AlarmTargetsRequest{
		OriginType	: request.OriginType,
		MailAddress : request.MailAddress,
		MailSendYn	: request.MailSendYn}

	status := b.txn.Debug().Table("alarm_targets").
		Model(&alarmTargets).
		Where("origin_type = ? ", alarmTargets.OriginType).
		Updates(map[string]interface{}{
			"mail_address": alarmTargets.MailAddress,
			"mail_send_yn": alarmTargets.MailSendYn,
			"modi_date": util.GetDBCurrentTime(),
			"modi_user": "test"})

	err := util.GetError().DbCheckError(status.Error)
	return err
}

func (b AlarmPolicyDao) UpdateAlarmSns(request model.AlarmPolicyRequest) (errMsg model.ErrMessage)  {

	alarmSns :=  model.AlarmSns{
		OriginType	: request.OriginType,
		SnsSendYn 	: request.SnsSendYn}

	status := b.txn.Debug().Table("alarm_sns").
		Model(&alarmSns).
		//Where("origin_type = ? ", alarmSns.OriginType).
		Updates(map[string]interface{}{
			"sns_send_yn": alarmSns.SnsSendYn,
			"modi_date": util.GetDBCurrentTime(),
			"modi_user": "system"})

	err := util.GetError().DbCheckError(status.Error)
	return err
}


func (b AlarmPolicyDao) GetAlarmSnsChannelList(request model.AlarmPolicyRequest) ([]model.AlarmSnsChannelResponse, model.ErrMessage) {

	var alarmSnsChannelList []model.AlarmSnsChannelResponse

	where := "1=1"
	if request.OriginType != "" {
		where += fmt.Sprintf(" and origin_type = '%s'", request.OriginType)
	}
	if request.SnsType != "" {
		where += fmt.Sprintf(" and sns_type = '%s'", request.SnsType)
	}
	status := b.txn.Debug().Table("alarm_sns").Where(where).Find(&alarmSnsChannelList)
	err := util.GetError().DbCheckError(status.Error)
	return alarmSnsChannelList, err
}

func (b AlarmPolicyDao) CreateAlarmSnsChannel(request model.AlarmPolicyRequest) model.ErrMessage {

	alarmSns := model.AlarmSns{
		OriginType:request.OriginType,
		SnsType:"telegram",
		SnsId:request.SnsId,
		Token:request.Token,
		Expl:request.Expl,
		SnsSendYn:request.SnsSendYn,
	}
	status := b.txn.Debug().Create(&alarmSns)
	err := util.GetError().DbCheckError(status.Error)
	return err
}

func (b AlarmPolicyDao) DeleteAlarmSnsChannel(request model.AlarmPolicyRequest) model.ErrMessage {

	alarmSns := model.AlarmSns{
		ChannelId:request.Id,
	}
	status := b.txn.Debug().Delete(&alarmSns)
	err := util.GetError().DbCheckError(status.Error)
	return err
}
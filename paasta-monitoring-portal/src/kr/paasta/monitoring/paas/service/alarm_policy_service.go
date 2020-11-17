package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring/paas/dao"
	"kr/paasta/monitoring/paas/model"
)

type AlarmPolicyService struct {
	txn *gorm.DB
}

func GetAlarmPolicyService(txn *gorm.DB) *AlarmPolicyService {
	return &AlarmPolicyService{
		txn: txn,
	}
}

//Alarm Policy Select
func (h *AlarmPolicyService) GetAlarmPolicyList() (result []model.AlarmPolicyResponse, err model.ErrMessage) {

	result, err = dao.GetAlarmPolicyDao(h.txn).GetAlarmPolicyList()

	if err != nil {
		return result, err
	}

	resultAddr, err := dao.GetAlarmPolicyDao(h.txn).GetAlarmTargetsList()

	if err != nil {
		return result, err
	}

	for i := 0; i < len(result); i++ {
		for _, val1 := range resultAddr {
			if result[i].OriginType == val1.OriginType {
				result[i].MailAddress = val1.MailAddress
				result[i].MailSendYn = val1.MailSendYn
			}
		}
	}

	return result, err
}

//Alarm Policy Update
func (h *AlarmPolicyService) UpdateAlarmPolicyList(apiRequest []model.AlarmPolicyRequest) (err model.ErrMessage) {
	i := 0
	for _, data := range apiRequest {

		if i < 3 {
			fmt.Println("Data::", data)
			err = dao.GetAlarmPolicyDao(h.txn).UpdateAlarmPolicy(data)
			if err != nil {
				fmt.Println("DB Error =========+>", err)
				return err
			}
		} else {

			err = dao.GetAlarmPolicyDao(h.txn).UpdateAlarmTargets(data)
			if err != nil {
				fmt.Println("DB Error =========+>", err)
				return err
			}

			err = dao.GetAlarmPolicyDao(h.txn).UpdateAlarmSns(data)
			if err != nil {
				fmt.Println("DB Error =========+>", err)
				return err
			}
		}
		i++
	}
	return nil
}

//Alarm Policy Select
func (h *AlarmPolicyService) GetAlarmSnsChannelList(request model.AlarmPolicyRequest) ([]model.AlarmSnsChannelResponse, model.ErrMessage) {

	result, err := dao.GetAlarmPolicyDao(h.txn).GetAlarmSnsChannelList(request)
	return result, err
}

func (h *AlarmService) CreateAlarmSnsChannel(request model.AlarmPolicyRequest, txn *gorm.DB) model.ErrMessage {

	err := dao.GetAlarmPolicyDao(h.txn).CreateAlarmSnsChannel(request)
	return err
}

func (h *AlarmService) DeleteAlarmSnsChannel(request model.AlarmPolicyRequest, txn *gorm.DB) model.ErrMessage {

	err := dao.GetAlarmPolicyDao(h.txn).DeleteAlarmSnsChannel(request)
	return err
}

package service

import (
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring/domain"
	"kr/paasta/monitoring/dao"
	"fmt"
)

type AlarmPolicyService struct {
	txn   *gorm.DB
}

func GetAlarmPolicyService(txn *gorm.DB) *AlarmPolicyService {
	return &AlarmPolicyService{
		txn:   txn,
	}
}

//Alarm Policy Select
func (h *AlarmPolicyService) GetAlarmPolicyList() (result []domain.AlarmPolicyResponse, err domain.ErrMessage) {

	result,  err = dao.GetAlarmPolicyDao(h.txn).GetAlarmPolicyList()

	if err != nil{
		return result , err
	}


	resultAddr,  err := dao.GetAlarmPolicyDao(h.txn).GetAlarmTargetsList()

	if err != nil{
		return result , err
	}

	for i :=0 ; i < len(result); i++ {

		for _, val1 := range resultAddr {

			if result[i].OriginType == val1.OriginType {

				result[i].MailAddress = val1.MailAddress

				//fmt.Println(result[i].OriginType+"/"+result[i].MailAddress)

			}
		}
	}


	return result, err
}

//Alarm Policy Update
func (h *AlarmPolicyService) UpdateAlarmPolicyList(apiRequest []domain.AlarmPolicyRequest) (err domain.ErrMessage) {
    i := 0
	for _, data := range apiRequest {

		if i<3 {
			fmt.Println("Data::", data)
			err = dao.GetAlarmPolicyDao(h.txn).UpdateAlarmPolicy(data)
			if err != nil {
				fmt.Println("DB Error =========+>", err)
				return err
			}
		}else {

			err = dao.GetAlarmPolicyDao(h.txn).UpdateAlarmTargets(data)
			if err != nil {
				fmt.Println("DB Error =========+>", err)
				return err
			}
		}
		i++
	}
	return nil
}
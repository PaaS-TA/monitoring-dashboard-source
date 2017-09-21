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

	return result, nil
}

//Alarm Policy Update
func (h *AlarmPolicyService) UpdateAlarmPolicyList(apiRequest []domain.AlarmPolicyRequest) (err domain.ErrMessage) {

	for _, data := range apiRequest {

		fmt.Println("Data::",data)
		err = dao.GetAlarmPolicyDao(h.txn).UpdateAlarmPolicy(data)
		if err != nil{
			fmt.Println("DB Error =========+>", err)
			return err
		}
	}
	return nil
}
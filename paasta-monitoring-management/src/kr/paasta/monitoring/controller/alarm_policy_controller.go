package controller

import (
	"net/http"
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring/service"
	"kr/paasta/monitoring/util"
	"encoding/json"
	"kr/paasta/monitoring/domain"
)

type alarmPolicyService struct {
	txn   *gorm.DB
}

func GetAlarmPolicyController(txn *gorm.DB) *alarmPolicyService {
	return &alarmPolicyService{
		txn:   txn,
	}
}

//Alarm 정책 조회
func (h *alarmPolicyService) GetAlarmPolicyList(w http.ResponseWriter, r *http.Request){

	alarmPolicyList, err := service.GetAlarmPolicyService(h.txn).GetAlarmPolicyList()

	if err != nil {
		util.RenderJsonResponse(err, w)
	}

	util.RenderJsonResponse(alarmPolicyList, w)
}

//Alarm정책 Update
func (h *alarmPolicyService) UpdateAlarmPolicyList(w http.ResponseWriter, r *http.Request){

	var apiRequest []domain.AlarmPolicyRequest


	json.NewDecoder(r.Body).Decode(&apiRequest)

	for _, data := range apiRequest{
		err := data.AlarmPolicyValidate(data)
		if err != nil{
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	}

	err := service.GetAlarmPolicyService(h.txn).UpdateAlarmPolicyList(apiRequest)

	if err != nil {
		util.RenderJsonResponse(err, w)
	}

	util.RenderJsonResponse(nil, w)
}



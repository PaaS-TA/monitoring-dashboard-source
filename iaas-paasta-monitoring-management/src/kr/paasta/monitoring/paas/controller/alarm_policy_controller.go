package controller

import (
	"net/http"
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring/paas/service"
	"kr/paasta/monitoring/paas/util"
	"encoding/json"
	"kr/paasta/monitoring/paas/model"
	"fmt"
	"strconv"
)

type AlarmPolicyService struct {
	txn   *gorm.DB
}

func GetAlarmPolicyController(txn *gorm.DB) *AlarmPolicyService {
	return &AlarmPolicyService{
		txn:   txn,
	}
}

// Alarm 정책 조회
func (h *AlarmPolicyService) GetAlarmPolicyList(w http.ResponseWriter, r *http.Request){

	alarmPolicyList, err := service.GetAlarmPolicyService(h.txn).GetAlarmPolicyList()

	if err != nil {
		util.RenderJsonResponse(err, w)
	}

	util.RenderJsonResponse(alarmPolicyList, w)
}

// Alarm정책 Update
func (h *AlarmPolicyService) UpdateAlarmPolicyList(w http.ResponseWriter, r *http.Request){

	var apiRequest []model.AlarmPolicyRequest

	json.NewDecoder(r.Body).Decode(&apiRequest)
    i:=0
	for _, data := range apiRequest {
		if i < 3 {
			err := data.AlarmPolicyValidate(data)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
		} else {
			err := data.AlarmEmailValidate(data)
			if err != nil{
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
		}
       	i++
	}

	err := service.GetAlarmPolicyService(h.txn).UpdateAlarmPolicyList(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(nil)

	return
}

func (h *AlarmPolicyService) GetAlarmSnsChannelList(w http.ResponseWriter, r *http.Request){

	var apiRequest model.AlarmPolicyRequest
	apiRequest.SnsType = r.URL.Query().Get("snsType")
	apiRequest.OriginType = r.URL.Query().Get("originType")

	alarmSnsChannelList, err := service.GetAlarmPolicyService(h.txn).GetAlarmSnsChannelList(apiRequest)

	if err != nil {
		fmt.Println(err)
		util.RenderJsonResponse(err, w)
	}

	util.RenderJsonResponse(alarmSnsChannelList, w)
}

func (h *AlarmPolicyService) CreateAlarmSnsChannel(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.AlarmPolicyRequest
	err := json.NewDecoder(r.Body).Decode(&apiRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	ErrMessage := service.GetAlarmService(h.txn).CreateAlarmSnsChannel(apiRequest, h.txn)
	if ErrMessage != nil {
		util.RenderJsonResponse(ErrMessage, w)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		w.Write([]byte("{\"status\":\"Created\"}"))
	}
}

func (h *AlarmPolicyService) DeleteAlarmSnsChannel(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.AlarmPolicyRequest
	id, _ := strconv.Atoi(r.FormValue(":id"))
	apiRequest.Id = uint(id)

	ErrMessage := service.GetAlarmService(h.txn).DeleteAlarmSnsChannel(apiRequest, h.txn)
	if ErrMessage != nil {
		util.RenderJsonResponse(ErrMessage, w)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(204)
		w.Write([]byte("{\"status\":\"No Content\"}"))
	}
}

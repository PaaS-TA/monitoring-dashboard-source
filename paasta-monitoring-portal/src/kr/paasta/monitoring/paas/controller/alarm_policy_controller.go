package controller

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"monitoring-portal/paas/model"
	"monitoring-portal/paas/service"
	"monitoring-portal/paas/util"
	"net/http"
	"strconv"
)

type AlarmPolicyService struct {
	txn *gorm.DB
}

func GetAlarmPolicyController(txn *gorm.DB) *AlarmPolicyService {
	return &AlarmPolicyService{
		txn: txn,
	}
}

// Alarm 정책 조회
func (h *AlarmPolicyService) GetAlarmPolicyList(w http.ResponseWriter, r *http.Request) {

	alarmPolicyList, err := service.GetAlarmPolicyService(h.txn).GetAlarmPolicyList()

	if err != nil {
		util.RenderJsonResponse(err, w)
	}

	util.RenderJsonResponse(alarmPolicyList, w)
}

// Alarm정책 Update
func (h *AlarmPolicyService) UpdateAlarmPolicyList(w http.ResponseWriter, r *http.Request) {
	var apiRequest []model.AlarmPolicyRequest
	data, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	err := json.Unmarshal(data, &apiRequest)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	i := 0
	for _, data := range apiRequest {
		if i < 3 {
			err := data.AlarmPolicyValidate(data)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
		} else {
			err := data.AlarmEmailValidate(data)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
		}
		i++
	}

	error := service.GetAlarmPolicyService(h.txn).UpdateAlarmPolicyList(apiRequest)
	if error != nil {
		util.RenderJsonResponse(err, w)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(nil)

	return
}

func (h *AlarmPolicyService) GetAlarmSnsChannelList(w http.ResponseWriter, r *http.Request) {

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

	err := service.GetAlarmService(h.txn).DeleteAlarmSnsChannel(apiRequest, h.txn)


	// TODO : 파라미터를 배열로 받아서 처리하는 방안 고민 필요..
	/*
		var err model.ErrMessage
		idArr, _ := r.URL.Query()["id"]
		for _, value := range idArr {
			id, _ := strconv.Atoi(value)
			apiRequest.Id = uint(id)
			err = service.GetAlarmService(h.txn).DeleteAlarmSnsChannel(apiRequest, h.txn)
			if (err != nil) {
				break
			}
		}
	*/

	util.RenderJsonResponse(err, w)
}


/**
2021.05.18 - PaaS 알람 SNS 정보 수정 기능 추가
*/
func (h *AlarmPolicyService) UpdateAlarmSnsChannel(w http.ResponseWriter, r *http.Request) {
	var apiRequest model.AlarmPolicyRequest
	err := json.NewDecoder(r.Body).Decode(&apiRequest)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	ErrMessage := service.GetAlarmService(h.txn).UpdateAlarmSnsChannel(apiRequest, h.txn)
	if ErrMessage != nil {
		util.RenderJsonResponse(ErrMessage, w)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		w.Write([]byte("{\"status\":\"Created\"}"))
	}
}

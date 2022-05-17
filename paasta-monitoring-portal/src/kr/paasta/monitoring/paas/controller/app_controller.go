package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"monitoring-portal/paas/model"
	"monitoring-portal/paas/service"
	"monitoring-portal/paas/util"
	"net/http"
	"strconv"
)

type AppController struct {
	txn *gorm.DB
}

func GetAppController(txn *gorm.DB) *AppController {
	return &AppController{
		txn: txn,
	}
}

func (s *AppController) UpdatePaasAppAutoScalingPolicy(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.AppAutoscalingPolicy
	err := json.NewDecoder(r.Body).Decode(&apiRequest)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	validation := apiRequest.DefaultAutoScalingPolicyValidate(apiRequest)
	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validation.Error()))
		return
	}
	result := service.GetAppService(s.txn).UpdatePaasAppAutoScalingPolicy(apiRequest)

	util.RenderJsonResponse(result, w)
}

func (s *AppController) GetPaasAppAutoScalingPolicy(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.AppAlarmReq

	apiRequest.AppGuid = r.URL.Query().Get("appGuid")

	if apiRequest.AppGuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.New("Required input value does not exist. [AppGuid]").Error()))
		return
	}

	result, err := service.GetAppService(s.txn).GetPaasAppAutoScalingPolicy(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(result, w)
	}

}

func (s *AppController) UpdatePaasAppPolicyInfo(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.AppAlarmPolicy
	err := json.NewDecoder(r.Body).Decode(&apiRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	validation := apiRequest.DefaultAlarmPolicyValidate(apiRequest)
	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validation.Error()))
		return
	}

	result := service.GetAppService(s.txn).UpdatePaasAppPolicyInfo(apiRequest)

	util.RenderJsonResponse(result, w)
}

func (s *AppController) GetPaasAppPolicyInfo(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.AppAlarmReq

	apiRequest.AppGuid = r.URL.Query().Get("appGuid")

	if apiRequest.AppGuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.New("Required input value does not exist. [AppGuid]").Error()))
		return
	}

	result, err := service.GetAppService(s.txn).GetPaasAppPolicyInfo(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(result, w)
	}

}

func (s *AppController) GetPaasAppAlarmList(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.AppAlarmReq

	apiRequest.ResourceType = r.URL.Query().Get("resourceType")
	pageItems, _ := strconv.Atoi(r.URL.Query().Get("pageItems"))
	pageIndex, _ := strconv.Atoi(r.URL.Query().Get("pageIndex"))
	apiRequest.PageItems = pageItems
	apiRequest.PageIndex = pageIndex
	apiRequest.AppGuid = r.URL.Query().Get("appGuid")
	apiRequest.AlarmLevel = r.URL.Query().Get("alarmLevel")
	apiRequest.SearchDateFrom = r.URL.Query().Get("searchDateFrom")
	apiRequest.SearchDateTo = r.URL.Query().Get("searchDateTo")

	validation := apiRequest.DefaultAlarmListValidate(apiRequest)
	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validation.Error()))
		return
	}

	result, err := service.GetAppService(s.txn).GetPaasAppAlarmList(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(result, w)
	}

}

func (s *AppController) DeletePaasAppPolicy(w http.ResponseWriter, r *http.Request) {

	guid := r.FormValue(":guid")
	fmt.Println(">>>>>>>>>>>>>>>> DeletePaasAppPolicy guid : ", guid)

	result := service.GetAppService(s.txn).DeletePaasAppPolicy(guid)

	util.RenderJsonResponse(result, w)

}

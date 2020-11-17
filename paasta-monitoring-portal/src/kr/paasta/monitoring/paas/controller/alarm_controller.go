package controller

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring/paas/model"
	"kr/paasta/monitoring/paas/service"
	"kr/paasta/monitoring/paas/util"
	"net/http"
	"strconv"
)

//Gorm Object Struct
type AlarmService struct {
	txn *gorm.DB
}

func GetAlarmController(txn *gorm.DB) *AlarmService {
	return &AlarmService{
		txn: txn,
	}
}

//Controller
func (h *AlarmService) Main(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Main API Called")

	url := "/public/index.html"
	http.Redirect(w, r, url, 302)
}

func (h *AlarmService) GetAlarmList(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.AlarmRequest

	apiRequest.PagingReq.PageIndex, _ = strconv.Atoi(r.FormValue("pageIndex")) // Page 번호
	apiRequest.PagingReq.PageItem, _ = strconv.Atoi(r.FormValue("pageItems"))  // Page당 보여주는 갯수

	apiRequest.OriginType = r.URL.Query().Get("originType")
	apiRequest.ResolveStatus = r.URL.Query().Get("resolveStatus")
	apiRequest.SearchDateFrom = r.URL.Query().Get("searchDateFrom")
	apiRequest.SearchDateTo = r.URL.Query().Get("searchDateTo")
	apiRequest.AlarmType = r.URL.Query().Get("alarmType")
	apiRequest.Level = r.URL.Query().Get("level")

	//service호출 (Gorm Obj 매개 변수)
	alarms, err := service.GetAlarmService(h.txn).GetAlarmList(apiRequest, h.txn)
	if err != nil {
		util.RenderJsonResponse(err, w)
		return
	}
	util.RenderJsonResponse(alarms, w)
}

func (h *AlarmService) GetAlarmListCount(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.AlarmRequest
	apiRequest.ResolveStatus = r.URL.Query().Get("resolveStatus")
	apiRequest.SearchDateFrom = r.URL.Query().Get("searchDateFrom")
	apiRequest.SearchDateTo = r.URL.Query().Get("searchDateTo")

	alarms, err := service.GetAlarmService(h.txn).GetAlarmListCount(apiRequest, h.txn)
	if err != nil {
		util.RenderJsonResponse(err, w)
		return
	}

	util.RenderJsonResponse(alarms, w)
}

func (h *AlarmService) GetAlarmResolveStatus(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.AlarmRequest
	apiRequest.ResolveStatus = r.FormValue(":resolveStatus")

	alarms, err := service.GetAlarmService(h.txn).GetAlarmResolveStatus(apiRequest, h.txn)
	if err != nil {
		util.RenderJsonResponse(err, w)
		return
	}
	util.RenderJsonResponse(alarms, w)
}

func (h *AlarmService) GetAlarmDetail(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.AlarmRequest
	id, _ := strconv.Atoi(r.FormValue(":id"))
	apiRequest.Id = uint(id)

	alarms, err := service.GetAlarmService(h.txn).GetAlarmDetail(apiRequest, h.txn)
	if err != nil {
		util.RenderJsonResponse(err, w)
		return
	}
	util.RenderJsonResponse(alarms, w)
}

func (h *AlarmService) UpdateAlarm(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.AlarmRequest
	err := json.NewDecoder(r.Body).Decode(&apiRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	id, _ := strconv.Atoi(r.FormValue(":id"))
	apiRequest.Id = uint(id)

	validation := apiRequest.AlarmValidate(apiRequest)
	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	ErrMessage := service.GetAlarmService(h.txn).UpdateAlarm(apiRequest, h.txn)
	if ErrMessage != nil {
		util.RenderJsonResponse(ErrMessage, w)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		w.Write([]byte("{\"status\":\"Created\"}"))
	}
}

func (h *AlarmService) CreateAlarmAction(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("enter CreateAlarmAction!!!")
	var apiRequest model.AlarmActionRequest
	err := json.NewDecoder(r.Body).Decode(&apiRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	ErrMessage := service.GetAlarmService(h.txn).CreateAlarmAction(apiRequest, h.txn)
	if ErrMessage != nil {
		util.RenderJsonResponse(ErrMessage, w)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		w.Write([]byte("{\"status\":\"Created\"}"))
	}
}

func (h *AlarmService) UpdateAlarmAction(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.AlarmActionRequest
	err := json.NewDecoder(r.Body).Decode(&apiRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	id, _ := strconv.Atoi(r.FormValue(":actionId"))
	apiRequest.Id = uint(id)

	validation := apiRequest.AlarmActionValidate(apiRequest)
	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	ErrMessage := service.GetAlarmService(h.txn).UpdateAlarmAction(apiRequest, h.txn)
	if ErrMessage != nil {
		util.RenderJsonResponse(ErrMessage, w)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		w.Write([]byte("{\"status\":\"Created\"}"))
	}
}

func (h *AlarmService) DeleteAlarmAction(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.AlarmActionRequest
	id, _ := strconv.Atoi(r.FormValue(":actionId"))
	apiRequest.Id = uint(id)

	ErrMessage := service.GetAlarmService(h.txn).DeleteAlarmAction(apiRequest, h.txn)
	if ErrMessage != nil {
		util.RenderJsonResponse(ErrMessage, w)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(204)
		w.Write([]byte("{\"status\":\"No Content\"}"))
	}
}

func (h *AlarmService) GetAlarmStat(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.AlarmStatRequest
	apiRequest.Period = r.URL.Query().Get("period")
	interval, _ := strconv.Atoi(r.URL.Query().Get("interval"))
	apiRequest.Interval = interval

	alarms, err := service.GetAlarmService(h.txn).GetAlarmStat(apiRequest, h.txn)
	if err != nil {
		util.RenderJsonResponse(err, w)
		return
	}
	util.RenderJsonResponse(alarms, w)
}

func (h *AlarmService) GetAlarmStatGraphTotal(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.AlarmStatRequest
	apiRequest.Period = r.URL.Query().Get("period")
	apiRequest.Interval, _ = strconv.Atoi(r.URL.Query().Get("interval"))
	apiRequest.Args = []model.AlarmStat{
		{model.ALARM_LEVEL_CRITICAL, model.ALARM_LEVEL_CRITICAL, "", ""},
		{model.ALARM_LEVEL_WARNING, model.ALARM_LEVEL_WARNING, "", ""}}
	result, err := service.GetAlarmService(h.txn).GetAlarmStatGraph(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
		return
	}
	util.RenderJsonResponse(result, w)
}

func (h *AlarmService) GetAlarmStatGraphService(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.AlarmStatRequest
	apiRequest.Period = r.URL.Query().Get("period")
	apiRequest.Interval, _ = strconv.Atoi(r.URL.Query().Get("interval"))
	apiRequest.Args = []model.AlarmStat{
		{model.ORIGIN_TYPE_BOSH + "-" + model.ALARM_LEVEL_CRITICAL, model.ALARM_LEVEL_CRITICAL, model.ORIGIN_TYPE_BOSH, ""},
		{model.ORIGIN_TYPE_BOSH + "-" + model.ALARM_LEVEL_WARNING, model.ALARM_LEVEL_WARNING, model.ORIGIN_TYPE_BOSH, ""},
		{model.ORIGIN_TYPE_PAASTA + "-" + model.ALARM_LEVEL_CRITICAL, model.ALARM_LEVEL_CRITICAL, model.ORIGIN_TYPE_PAASTA, ""},
		{model.ORIGIN_TYPE_PAASTA + "-" + model.ALARM_LEVEL_WARNING, model.ALARM_LEVEL_WARNING, model.ORIGIN_TYPE_PAASTA, ""},
		{model.ORIGIN_TYPE_CONTAINER + "-" + model.ALARM_LEVEL_CRITICAL, model.ALARM_LEVEL_CRITICAL, model.ORIGIN_TYPE_CONTAINER, ""},
		{model.ORIGIN_TYPE_CONTAINER + "-" + model.ALARM_LEVEL_WARNING, model.ALARM_LEVEL_WARNING, model.ORIGIN_TYPE_CONTAINER, ""}}
	result, err := service.GetAlarmService(h.txn).GetAlarmStatGraph(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
		return
	}
	util.RenderJsonResponse(result, w)
}

func (h *AlarmService) GetAlarmStatGraphMatrix(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.AlarmStatRequest
	apiRequest.Period = r.URL.Query().Get("period")
	apiRequest.Interval, _ = strconv.Atoi(r.URL.Query().Get("interval"))
	apiRequest.Args = []model.AlarmStat{
		{model.ALARM_TYPE_CPU + "-" + model.ALARM_LEVEL_CRITICAL, model.ALARM_LEVEL_CRITICAL, "", model.ALARM_TYPE_CPU},
		{model.ALARM_TYPE_CPU + "-" + model.ALARM_LEVEL_WARNING, model.ALARM_LEVEL_WARNING, "", model.ALARM_TYPE_CPU},
		{model.ALARM_TYPE_MEMORY + "-" + model.ALARM_LEVEL_CRITICAL, model.ALARM_LEVEL_CRITICAL, "", model.ALARM_TYPE_MEMORY},
		{model.ALARM_TYPE_MEMORY + "-" + model.ALARM_LEVEL_WARNING, model.ALARM_LEVEL_WARNING, "", model.ALARM_TYPE_MEMORY},
		{model.ALARM_TYPE_DISK + "-" + model.ALARM_LEVEL_CRITICAL, model.ALARM_LEVEL_CRITICAL, "", model.ALARM_TYPE_DISK},
		{model.ALARM_TYPE_DISK + "-" + model.ALARM_LEVEL_WARNING, model.ALARM_LEVEL_WARNING, "", model.ALARM_TYPE_DISK}}
	result, err := service.GetAlarmService(h.txn).GetAlarmStatGraph(apiRequest)
	if err != nil {
		util.RenderJsonResponse(err, w)
		return
	}
	util.RenderJsonResponse(result, w)
}

func (h *AlarmService) GetPaasAlarmRealTimeCount(w http.ResponseWriter, r *http.Request) {

	result, err := service.GetAlarmService(h.txn).GetPaasAlarmRealTimeCount()
	if err != nil {
		util.RenderJsonResponse(err, w)
		return
	}
	util.RenderJsonResponse(result, w)
}

func (h *AlarmService) GetPaasAlarmRealTimeList(w http.ResponseWriter, r *http.Request) {

	result, err := service.GetAlarmService(h.txn).GetPaasAlarmRealTimeList()
	if err != nil {
		util.RenderJsonResponse(err, w)
		return
	}
	util.RenderJsonResponse(result, w)
}

package controller

import (
	"github.com/monasca/golang-monascaclient/monascaclient"
	client "github.com/influxdata/influxdb/client/v2"
	mod "github.com/monasca/golang-monascaclient/monascaclient/models"
	"kr/paasta/monitoring/utils"
	"kr/paasta/monitoring/iaas/model"
	"kr/paasta/monitoring/iaas/service"
	"net/http"
	"strconv"
	"github.com/jinzhu/gorm"
	"encoding/json"
)

//Compute Node Controller
type AlarmStatusController struct{
	monClient     monascaclient.Client
	influxClient  client.Client
	txn           *gorm.DB
}

func NewAlarmStatusController(monClient monascaclient.Client, influxClient client.Client, txn *gorm.DB) *AlarmStatusController {
	return &AlarmStatusController{
		monClient: monClient,
		influxClient: influxClient,
		txn:   txn,
	}
}


func (s *AlarmStatusController)GetAlarmStatusCount(w http.ResponseWriter, r *http.Request) {

	var query mod.AlarmQuery
	state      := r.FormValue("state")
	if r.FormValue("state") != ""{
		query.State = &state
	}
	monClient,  err := utils.GetMonascaClient(r, s.monClient)
	result, err := services.GetAlarmStatusService(monClient, s.influxClient, s.txn).GetAlarmStatusCount(query)

	statusErr := utils.GetError().GetCheckErrorMessage(err)

	if statusErr != nil {
		utils.ErrRenderJsonResponse(statusErr, w)
	}else{
		utils.RenderJsonResponse(result, w)
	}
}

func (s *AlarmStatusController)GetAlarmStatusList(w http.ResponseWriter, r *http.Request) {

	var query mod.AlarmQuery

	severity  := r.FormValue("severity")
	state      := r.FormValue("state")
	offset, _  := strconv.Atoi(r.FormValue("offset"))
	limit, _   := strconv.Atoi(r.FormValue("limit"))
	orderBy := "state_updated_timestamp desc, state desc"

	if r.FormValue("severity") != ""{
		query.Severity = &severity
	}
	if r.FormValue("state") != ""{
		query.State = &state
	}
	if r.FormValue("offset") != ""{
		query.Offset = &offset
	}
	if r.FormValue("limit") != ""{
		query.Limit = &limit
	}

	query.SortBy = &orderBy
	monClient,  err := utils.GetMonascaClient(r, s.monClient)
	result, err := services.GetAlarmStatusService(monClient, s.influxClient, s.txn).GetAlarmStatusList(query)

	statusErr := utils.GetError().GetCheckErrorMessage(err)

	if statusErr != nil {
		utils.ErrRenderJsonResponse(statusErr, w)
	}else{
		utils.RenderJsonResponse(result, w)
	}

}

func (s *AlarmStatusController)GetAlarmStatus(w http.ResponseWriter, r *http.Request) {

	alarmId       := r.FormValue(":alarmId")
	monClient,  err := utils.GetMonascaClient(r, s.monClient)
	result, err := services.GetAlarmStatusService(monClient, s.influxClient, s.txn).GetAlarmStatus(alarmId)

	statusErr := utils.GetError().GetCheckErrorMessage(err)

	if statusErr != nil {
		utils.ErrRenderJsonResponse(statusErr, w)
	}else{
		utils.RenderJsonResponse(result, w)
	}

}

func (s *AlarmStatusController)GetAlarmHistoryList(w http.ResponseWriter, r *http.Request) {

	var alarmReq model.AlarmReq
	alarmId    := r.FormValue(":alarmId")
	timeRange  := r.FormValue("timeRange")

	alarmReq.AlarmId = alarmId

	//TimeRange는 1w, 1d, 10m 등으로 사용
	if r.FormValue("timeRange") != ""{
		alarmReq.TimeRange = timeRange
	}else{
		alarmReq.TimeRange = "1d"
	}
	monClient,  _  := utils.GetMonascaClient(r, s.monClient)
	result, err := services.GetAlarmStatusService(monClient, s.influxClient, s.txn).GetAlarmHistoryList(alarmReq)
	//statusErr := utils.GetError().GetCheckErrorMessage(err)

	if err != nil {
		utils.ErrRenderJsonResponse(err, w)
	}else{
		utils.RenderJsonResponse(result, w)
	}
}

func (s *AlarmStatusController) GetAlarmHistoryActionList(w http.ResponseWriter, r *http.Request) {

	alarmId    := r.FormValue(":alarmId")
	monClient,  _  := utils.GetMonascaClient(r, s.monClient)
	result, err := services.GetAlarmStatusService(monClient, s.influxClient, s.txn).GetAlarmHistoryActionList(alarmId)
	statusErr := utils.GetError().GetCheckErrorMessage(err)

	if statusErr != nil {
		utils.ErrRenderJsonResponse(statusErr, w)
	}else{
		utils.RenderJsonResponse(result, w)
	}

}

func (s *AlarmStatusController) CreateAlarmHistoryAction(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.AlarmActionRequest
	err := json.NewDecoder(r.Body).Decode(&apiRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	monClient,  _  := utils.GetMonascaClient(r, s.monClient)
	err  = services.GetAlarmStatusService( monClient, s.influxClient, s.txn).CreateAlarmHistoryAction(apiRequest)

	if err != nil {
		utils.ErrRenderJsonResponse(err, w)
	}else{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		w.Write([]byte("{\"status\":\"Created\"}"))
	}

}



func (s *AlarmStatusController) UpdateAlarmHistoryAction(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.AlarmActionRequest
	actionId , _   := strconv.Atoi(r.FormValue(":id"))
	apiRequest.Id = uint(actionId)

	err := json.NewDecoder(r.Body).Decode(&apiRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	monClient,  _  := utils.GetMonascaClient(r, s.monClient)
	err  = services.GetAlarmStatusService( monClient, s.influxClient, s.txn).UpdateAlarmAction(apiRequest)

	if err != nil {
		utils.ErrRenderJsonResponse(err, w)
	}else{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("{\"status\":\"Created\"}"))
	}

}


func (s *AlarmStatusController) DeleteAlarmHistoryAction(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.AlarmActionRequest
	actionId , _   := strconv.Atoi(r.FormValue(":id"))
	apiRequest.Id = uint(actionId)
	monClient,  _  := utils.GetMonascaClient(r, s.monClient)
	err  := services.GetAlarmStatusService( monClient, s.influxClient, s.txn).DeleteAlarmAction(apiRequest)

	if err != nil {
		utils.ErrRenderJsonResponse(err, w)
	}else{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("{\"status\":\"Deleted\"}"))
	}

}

func (s *AlarmStatusController)GetIaasAlarmRealTimeList(w http.ResponseWriter, r *http.Request) {

	var query mod.AlarmQuery
	targetState := model.ALARM_STATE_ALARM
	query.State = &targetState
	sortBy := "updated_timestamp desc"
	query.SortBy = &sortBy

	result, err := services.GetAlarmStatusService(s.monClient, s.influxClient, s.txn).GetAlarmStatusList(query)
	statusErr := utils.GetError().GetCheckErrorMessage(err)
	if statusErr != nil {
		utils.ErrRenderJsonResponse(statusErr, w)
	} else {
		utils.RenderJsonResponse(result, w)
	}
}

func (s *AlarmStatusController)GetIaasAlarmRealTimeCount(w http.ResponseWriter, r *http.Request) {
	monClient,  _  := utils.GetMonascaClient(r, s.monClient)
	result, err := services.GetAlarmStatusService(monClient, s.influxClient, s.txn).GetIaasAlarmRealTimeCount()

	statusErr := utils.GetError().GetCheckErrorMessage(err)

	if statusErr != nil {
		utils.ErrRenderJsonResponse(statusErr, w)
	}else{
		utils.RenderJsonResponse(result, w)
	}
}
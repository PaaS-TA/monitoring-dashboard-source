package controller

import (
	client "github.com/influxdata/influxdb/client/v2"
	"kr/paasta/monitoring/iaas/service"
	"kr/paasta/monitoring/utils"
	mod "github.com/monasca/golang-monascaclient/monascaclient/models"
	"github.com/monasca/golang-monascaclient/monascaclient"
	"net/http"
	"encoding/json"
	"errors"
	"strconv"
	"fmt"
)

//Compute Node Controller
type AlarmDefinitionController struct{
	monClient     monascaclient.Client
	influxClient  client.Client
}

func NewAlarmDefinitionController(monClient monascaclient.Client, influxClient client.Client) *AlarmDefinitionController {
	return &AlarmDefinitionController{
		monClient: monClient,
		influxClient: influxClient,
	}
}


func (s *AlarmDefinitionController)GetAlarmDefinitionList(w http.ResponseWriter, r *http.Request){

	var query mod.AlarmDefinitionQuery

	severity  := r.FormValue("severity")
	name  := r.FormValue("name")
	offset, _  := strconv.Atoi(r.FormValue("offset"))
	limit, _   := strconv.Atoi(r.FormValue("limit"))

	//AlarmDefinition Order by
	orderBy := "name"
	if r.FormValue("name") != ""{
		query.Name = &name
	}
	if r.FormValue("severity") != ""{
		query.Severity = &severity
	}
	if r.FormValue("offset") != ""{
		query.Offset = &offset
	}
	if r.FormValue("limit") != ""{
		query.Limit = &limit
	}
	query.SortBy = &orderBy

	monClient,  err := utils.GetMonascaClient(r, s.monClient)
	result, err := services.GetAlarmDefinitionService(monClient, s.influxClient).GetAlarmDefinitionList(query)

	if err != nil {
		utils.ErrRenderJsonResponse(err, w)
	}else{
		utils.RenderJsonResponse(result, w)
	}
}


func (s *AlarmDefinitionController)GetAlarmDefinition(w http.ResponseWriter, r *http.Request){

	definitionId  := r.FormValue(":id")
	monClient,  err := utils.GetMonascaClient(r, s.monClient)
	result, err := services.GetAlarmDefinitionService(monClient, s.influxClient).GetAlarmDefinition(definitionId)

	if err != nil {
		utils.ErrRenderJsonResponse(err, w)
	}else{
		utils.RenderJsonResponse(result, w)
	}
}

func (s *AlarmDefinitionController)CreateAlarmDefinition(w http.ResponseWriter, r *http.Request){

	fmt.Println("!!!!!!!!!!!!=========================>>")
	var query mod.AlarmDefinitionRequestBody
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	validation := definitionValidate(query)

	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validation.Error()))
		return
	}

	monClient,  _ := utils.GetMonascaClient(r, s.monClient)
	fmt.Println("Data Name::", *query.Name)
	err = services.GetAlarmDefinitionService(monClient, s.influxClient).CreateAlarmDefinition(query)
	fmt.Println("Data Err::",err)
	notiErr := utils.GetError().GetCheckErrorMessage(err)

	if notiErr != nil {
		utils.ErrRenderJsonResponse(notiErr, w)
	}else{
		utils.RenderJsonResponse(nil, w)
	}

}


func (s *AlarmDefinitionController)UpdateAlarmDefinition(w http.ResponseWriter, r *http.Request){

	var alarmDefinitionRequestBody mod.AlarmDefinitionRequestBody
	definitionId  := r.FormValue(":id")

	err := json.NewDecoder(r.Body).Decode(&alarmDefinitionRequestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	monClient,  _ := utils.GetMonascaClient(r, s.monClient)
	err = services.GetAlarmDefinitionService(monClient, s.influxClient).UpdateAlarmDefinition(definitionId, alarmDefinitionRequestBody)

	notiErr := utils.GetError().GetCheckErrorMessage(err)

	if notiErr != nil {
		utils.ErrRenderJsonResponse(notiErr, w)
	}else{
		utils.RenderJsonResponse(nil, w)
	}

}

func (s *AlarmDefinitionController)DeleteAlarmDefinition(w http.ResponseWriter, r *http.Request){

	notificationId  := r.FormValue(":id")
	monClient,  _ := utils.GetMonascaClient(r, s.monClient)
	err := services.GetAlarmDefinitionService(monClient, s.influxClient).DeleteAlarmDefinition(notificationId)

	if err != nil {
		utils.ErrRenderJsonResponse(err, w)
	}else{
		utils.RenderJsonResponse(nil, w)
	}
}


func definitionValidate(apiRequest mod.AlarmDefinitionRequestBody) error {

	if apiRequest.Name == nil{
		return errors.New("Required input value does not exist. [name]");
	}

	if apiRequest.Expression == nil{
		return errors.New("Required input value does not exist. [expression]");
	}

	if len(*apiRequest.AlarmActions) <= 0{
		return errors.New("Required input value does not exist. [alarm_actions]");
	}

	return nil
}

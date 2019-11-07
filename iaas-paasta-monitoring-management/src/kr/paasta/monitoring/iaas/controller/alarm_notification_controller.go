package controller

import (
	"encoding/json"
	"errors"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/monasca/golang-monascaclient/monascaclient"
	mod "github.com/monasca/golang-monascaclient/monascaclient/models"
	"kr/paasta/monitoring/iaas/service"
	"kr/paasta/monitoring/utils"
	"net/http"
	"strconv"
)

//Compute Node Controller
type NotificationController struct {
	monClient    monascaclient.Client
	influxClient client.Client
}

func NewNotificationController(monClient monascaclient.Client, influxClient client.Client) *NotificationController {
	return &NotificationController{
		monClient:    monClient,
		influxClient: influxClient,
	}
}

func (s *NotificationController) GetAlarmNotificationList(w http.ResponseWriter, r *http.Request) {
	var query mod.NotificationQuery

	offset := r.FormValue("offset")
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	orderBy := "name"

	if r.FormValue("offset") != "" {
		query.Offset = &offset
	}
	if r.FormValue("limit") != "" {
		query.Limit = &limit
	}
	query.SortBy = &orderBy

	monClient, err := utils.GetMonascaClient(r, s.monClient)
	result, err := services.GetNotificationService(monClient, s.influxClient).GetAlarmNotificationList(query)
	notiErr := utils.GetError().GetCheckErrorMessage(err)
	if notiErr != nil {
		utils.ErrRenderJsonResponse(notiErr, w)
	} else {
		utils.RenderJsonResponse(result, w)
	}
}

func (s *NotificationController) UpdateAlarmNotification(w http.ResponseWriter, r *http.Request) {

	var query mod.NotificationRequestBody
	notificationId := r.FormValue(":id")

	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	period := 0
	notiType := "EMAIL"
	query.Period = &period
	query.Type = &notiType

	validation := notificateValidate(query)

	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validation.Error()))
		return
	}
	monClient, err := utils.GetMonascaClient(r, s.monClient)
	err = services.GetNotificationService(monClient, s.influxClient).UpdateAlarmNotification(notificationId, query)
	notiErr := utils.GetError().GetCheckErrorMessage(err)

	if notiErr != nil {
		utils.ErrRenderJsonResponse(notiErr, w)
	} else {
		utils.RenderJsonResponse(nil, w)
	}

}

func (s *NotificationController) CreateAlarmNotification(w http.ResponseWriter, r *http.Request) {

	var query mod.NotificationRequestBody
	err := json.NewDecoder(r.Body).Decode(&query)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	period := 0
	notiType := "EMAIL"
	query.Period = &period
	query.Type = &notiType

	validation := notificateValidate(query)

	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validation.Error()))
		return
	}

	monClient, err := utils.GetMonascaClient(r, s.monClient)
	createdId, err := services.GetNotificationService(monClient, s.influxClient).CreateAlarmNotification(query)
	notiErr := utils.GetError().GetCheckErrorMessage(err)

	if notiErr != nil {
		utils.ErrRenderJsonResponse(notiErr, w)
	} else {
		utils.RenderJsonResponse(createdId, w)
	}

}

func (s *NotificationController) DeleteAlarmNotification(w http.ResponseWriter, r *http.Request) {

	notificationId := r.FormValue(":id")
	monClient, _ := utils.GetMonascaClient(r, s.monClient)
	err := services.GetNotificationService(monClient, s.influxClient).DeleteAlarmNotification(notificationId)
	notiErr := utils.GetError().GetCheckErrorMessage(err)
	if notiErr != nil {
		utils.ErrRenderJsonResponse(notiErr, w)
	} else {
		utils.RenderJsonResponse(nil, w)
	}
}

func notificateValidate(apiRequest mod.NotificationRequestBody) error {

	if apiRequest.Name == nil {
		return errors.New("Required input value does not exist. [name]")
	}

	if apiRequest.Address == nil {
		return errors.New("Required input value does not exist. [address]")
	}
	return nil
}

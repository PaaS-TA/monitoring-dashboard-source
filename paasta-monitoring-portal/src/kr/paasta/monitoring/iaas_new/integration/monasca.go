package integration

import (
	"github.com/monasca/golang-monascaclient/monascaclient"
	mod "github.com/monasca/golang-monascaclient/monascaclient/models"
	"kr/paasta/monitoring/iaas_new/model"
	"strings"
	"time"
)

type Monasca struct {
	monClient monascaclient.Client
}

func GetMonasca(monClient monascaclient.Client) *Monasca {
	return &Monasca{
		monClient: monClient,
	}
}

//Alarm Notification List조회
func (m *Monasca) GetAlarmNotificationList(query mod.NotificationQuery) (result []model.AlarmNotification, err error) {

	notificationList, notiErr := m.monClient.GetNotificationMethods(&query)

	if notiErr != nil {
		return nil, notiErr
	}

	for _, data := range notificationList.Elements {
		var notification model.AlarmNotification
		notification.Id = data.ID
		notification.Name = data.Name
		notification.Email = data.Address

		result = append(result, notification)
	}
	return result, err
}

func (m *Monasca) GetAlarmNotification(id string) (*mod.NotificationElement, error) {

	var notiQuery mod.NotificationQuery
	//query := new(mod.NotificationQuery)

	return m.monClient.GetNotificationMethod(id, &notiQuery)

}

func (m *Monasca) CreateAlarmNotification(query mod.NotificationRequestBody) (string, error) {
	res, err := m.monClient.CreateNotificationMethod(&query)
	return res.ResponseElement.ID, err
}

func (m *Monasca) UpdateAlarmNotification(notificationId string, notificationRequestBody mod.NotificationRequestBody) error {

	_, err := m.monClient.UpdateNotificationMethod(notificationId, &notificationRequestBody)
	return err
}

func (m *Monasca) DeleteAlarmNotification(notificationId string) error {

	err := m.monClient.DeleteNotificationMethod(notificationId)
	return err
}

func (m *Monasca) GetAlarmDefinitionList(query mod.AlarmDefinitionQuery) (result []model.AlarmDefinition, err error) {

	definitionList, defErr := m.monClient.GetAlarmDefinitions(&query)

	if defErr != nil {
		return nil, defErr
	}
	for _, data := range definitionList.Elements {

		var definition model.AlarmDefinition
		definition.Id = data.ID
		definition.Name = data.Name
		definition.AlarmAction = data.AlarmActions
		definition.OkAction = data.OkActions
		definition.Expression = data.Expression
		definition.MatchBy = data.MatchBy
		definition.Severity = data.Severity
		definition.Description = data.Description
		definition.UndeterminedActions = data.UndeterminedActions

		result = append(result, definition)
	}

	return result, err
}

func (m *Monasca) GetAlarmDefinition(id string) (result model.AlarmDefinition, err error) {

	definitionDetail, defErr := m.monClient.GetAlarmDefinition(id)

	if defErr != nil {
		return result, defErr
	}

	var definition model.AlarmDefinition
	definition.Id = definitionDetail.ID
	definition.Name = definitionDetail.Name
	definition.AlarmAction = definitionDetail.AlarmActions
	definition.OkAction = definitionDetail.OkActions
	definition.Expression = definitionDetail.Expression
	definition.MatchBy = definitionDetail.MatchBy
	definition.Severity = definitionDetail.Severity
	definition.Description = definitionDetail.Description
	definition.UndeterminedActions = definitionDetail.UndeterminedActions

	//definitionDetail.AlarmActions

	return definition, nil
}

func (m *Monasca) CreateAlarmDefinitionList(alarmDefinitionRequestBody mod.AlarmDefinitionRequestBody) error {

	_, err := m.monClient.CreateAlarmDefinition(&alarmDefinitionRequestBody)
	return err
}

func (m *Monasca) UpdateAlarmDefinitionList(alarmDefinitionId string, alarmDefinitionRequestBody mod.AlarmDefinitionRequestBody) error {

	//_, err := m.monClient.UpdateAlarmDefinition(alarmDefinitionId, &alarmDefinitionRequestBody)
	_, err := m.monClient.PatchAlarmDefinition(alarmDefinitionId, &alarmDefinitionRequestBody)
	return err
}

func (m *Monasca) DeleteAlarmDefinitionList(alarmDefinitionId string) error {

	err := m.monClient.DeleteAlarmDefinition(alarmDefinitionId)
	return err
}

func (m *Monasca) GetAlarmCount(query mod.AlarmQuery) ([][]int, error) {

	result, alarmErr := m.monClient.GetAlarmCount(&query)

	if alarmErr != nil {
		return result, alarmErr
	}
	return result, nil
}

func (m *Monasca) GetAlarmList(query mod.AlarmQuery) (result []model.AlarmStatus, err error) {

	alarmStatusList, alarmErr := m.monClient.GetAlarms(&query)

	if alarmErr != nil {
		return result, alarmErr
	}

	for _, data := range alarmStatusList.Elements {

		var status model.AlarmStatus
		status.Id = data.ID
		status.State = data.State
		status.UpdateTime = data.StateUpdatedTimestamp.Add(time.Duration(model.GmtTimeGap) * time.Hour).Format("2006-01-02 15:04:05")

		var metricNameList []string
		for _, dimensionData := range data.Metrics {
			metricNameList = append(metricNameList, dimensionData.Name)

			if dimensionData.Dimensions["component"] != "" {
				status.Type = dimensionData.Dimensions["component"]
				status.Zone = dimensionData.Dimensions["zone"]
			} else {
				status.Type = "Node"
			}
			status.HostName = dimensionData.Dimensions["hostname"]
		}
		status.MetricName = strings.Join(metricNameList[:], ",")
		result = append(result, status)
	}

	if alarmErr != nil {
		err = alarmErr
	}
	return result, err
}

func (m *Monasca) GetAlarm(id string) (result model.AlarmStatus, err error) {

	alarmStatus, alarmErr := m.monClient.GetAlarm(id)

	if alarmErr != nil {
		return result, alarmErr
	}

	result.Id = alarmStatus.ID
	result.State = alarmStatus.State
	//fmt.Println( alarmStatus.StateUpdatedTimestamp.Add(time.Duration(models.GmtTimeGap) * time.Hour).Format("2006-01-02 15:04:05"))
	result.UpdateTime = alarmStatus.StateUpdatedTimestamp.Add(time.Duration(model.GmtTimeGap) * time.Hour).Format("2006-01-02 15:04:05")

	var metricNameList []string
	for _, dimensionData := range alarmStatus.Metrics {
		metricNameList = append(metricNameList, dimensionData.Name)

		if dimensionData.Dimensions["component"] != "" {
			result.Type = dimensionData.Dimensions["component"]
			result.Zone = dimensionData.Dimensions["zone"]
		} else {
			result.Type = "Node"
		}
		result.HostName = dimensionData.Dimensions["hostname"]
	}
	result.MetricName = strings.Join(metricNameList[:], ",")

	return result, err
}

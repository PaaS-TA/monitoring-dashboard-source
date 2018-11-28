package services

import (
	client "github.com/influxdata/influxdb/client/v2"
	"github.com/monasca/golang-monascaclient/monascaclient"
	mod "github.com/monasca/golang-monascaclient/monascaclient/models"
	"kr/paasta/monitoring/iaas/integration"
	"kr/paasta/monitoring/iaas/model"
)

type AlarmDefinitionService struct {
	monClient     monascaclient.Client
	influxClient 	client.Client
}

func GetAlarmDefinitionService( monClient    monascaclient.Client, influxClient client.Client) *AlarmDefinitionService {
	return &AlarmDefinitionService{
		monClient: monClient,
		influxClient: 	influxClient,
	}
}

func (a *AlarmDefinitionService)GetAlarmDefinitionList(query mod.AlarmDefinitionQuery) (map[string]interface{}, error){

	var allQuery mod.AlarmDefinitionQuery

	if query.Severity != nil{
		allQuery.Severity = query.Severity
	}
	if query.Name != nil{
		allQuery.Name = query.Name
	}

	allData, _ := integration.GetMonasca(a.monClient).GetAlarmDefinitionList(allQuery)

	definitionList, _ := integration.GetMonasca(a.monClient).GetAlarmDefinitionList(query)

	resultData := map[string]interface{}{
		model.RESULT_CNT:  len(allData),
		model.RESULT_DATA: definitionList,
	}
	return resultData, nil

}

func (a *AlarmDefinitionService)GetAlarmDefinition(id string) (result model.AlarmDefinitionDetail, err error){

	definition, err := integration.GetMonasca(a.monClient).GetAlarmDefinition(id)
	//var result models.AlarmDefinitionDetail
	var notiList []model.AlarmNotification

	result.Id = definition.Id
	result.Name = definition.Name
	result.Description = definition.Description
	result.Expression = definition.Expression
	result.MatchBy = definition.MatchBy
	//result.OkAction = definition.OkAction
	result.Severity = definition.Severity
	result.UndeterminedActions = definition.UndeterminedActions

	for _, data := range definition.AlarmAction{
		//var notificationQuery mod.NotificationQuery

		notifications , _ := integration.GetMonasca(a.monClient).GetAlarmNotification(data)
		var noti model.AlarmNotification
		noti.Name = notifications.Name
		noti.Email = notifications.Address
		noti.Id = notifications.ID
		noti.Period = notifications.Period

		notiList = append(notiList, noti)
	}
	result.AlarmNotification = notiList
	return result, err


}

func (a *AlarmDefinitionService)CreateAlarmDefinition(query mod.AlarmDefinitionRequestBody) error{

	return integration.GetMonasca(a.monClient).CreateAlarmDefinitionList(query)

}

func (a *AlarmDefinitionService)UpdateAlarmDefinition(alarmDefinitionId string, alarmDefinitionRequestBody mod.AlarmDefinitionRequestBody) error{

	return integration.GetMonasca(a.monClient).UpdateAlarmDefinitionList(alarmDefinitionId, alarmDefinitionRequestBody)

}


func (a *AlarmDefinitionService)DeleteAlarmDefinition(alarmDefinitionId string) error{

	return integration.GetMonasca(a.monClient).DeleteAlarmDefinitionList(alarmDefinitionId)

}
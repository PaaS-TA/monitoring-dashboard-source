package services

import (
	client "github.com/influxdata/influxdb/client/v2"
	"github.com/monasca/golang-monascaclient/monascaclient"
	mod "github.com/monasca/golang-monascaclient/monascaclient/models"
	"kr/paasta/monitoring/iaas/integration"
	"kr/paasta/monitoring/iaas/model"
	"fmt"
)

type NotificationService struct {
	monClient     monascaclient.Client
	influxClient 	client.Client
}

func GetNotificationService( monClient    monascaclient.Client, influxClient client.Client) *NotificationService {
	return &NotificationService{
		monClient: monClient,
		influxClient: 	influxClient,
	}
}

func (n *NotificationService)GetAlarmNotificationList(query mod.NotificationQuery)(map[string]interface{}, error){

	var allQuery mod.NotificationQuery
	fmt.Println("MonCLinet====>>>>", n.monClient)
	allData, err := integration.GetMonasca(n.monClient).GetAlarmNotificationList(allQuery)
	result, err := integration.GetMonasca(n.monClient).GetAlarmNotificationList(query)

	if err != nil{
		fmt.Println("err==::", err)
		return nil, err
	}
	resultData := map[string]interface{}{
		model.RESULT_CNT:  len(allData),
		model.RESULT_DATA: result,
	}

	return resultData, err

}


func (a *NotificationService)UpdateAlarmNotification(notificationId string, notificationRequestBody mod.NotificationRequestBody)(error){

	return integration.GetMonasca(a.monClient).UpdateAlarmNotification(notificationId, notificationRequestBody)
}


func (a *NotificationService)CreateAlarmNotification(query mod.NotificationRequestBody)(string, error){
	return integration.GetMonasca(a.monClient).CreateAlarmNotification(query)
}

func (a *NotificationService)DeleteAlarmNotification(id string)(error){

	return integration.GetMonasca(a.monClient).DeleteAlarmNotification(id)
}



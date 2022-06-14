package iaas

import (
	"github.com/gophercloud/gophercloud"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	"paasta-monitoring-api/middlewares/zabbix-client/lib/go-zabbix"
	service "paasta-monitoring-api/services/api/v1/iaas"
)

type ZabbixController struct {
	ZabbixSession *zabbix.Session
	OpenstackProvider *gophercloud.ProviderClient
}


func GetZabbixController(zabbixSession *zabbix.Session, openstackProvider *gophercloud.ProviderClient) *ZabbixController {
	return &ZabbixController{
		ZabbixSession: zabbixSession,
		OpenstackProvider : openstackProvider,
	}
}


func (controller *ZabbixController) GetCpuUsage(ctx echo.Context) error {
	instanceId := ctx.QueryParam("instance_id")
	hypervisorName := ctx.QueryParam("host")

	result, err := service.GetZabbixService(controller.ZabbixSession, controller.OpenstackProvider).GetCpuUsage(instanceId, hypervisorName)

	resultMap := make(map[string]interface{})
	resultMap["label"] = "CPU"
	resultMap["data"] = result

	resultList := make([]interface{}, 1)
	resultList[0] = resultMap

	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get CPU usage.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", resultList)
	}
	return nil
}
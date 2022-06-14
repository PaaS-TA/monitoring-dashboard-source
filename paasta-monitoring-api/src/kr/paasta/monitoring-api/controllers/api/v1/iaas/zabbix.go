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



/**
	메모리 사용률 차트 데이터를 불러옴
*/
func (controller *ZabbixController) GetMemoryUsage(ctx echo.Context) error {
	instanceId := ctx.QueryParam("instance_id")
	hypervisorName := ctx.QueryParam("host")

	result, err := service.GetZabbixService(controller.ZabbixSession, controller.OpenstackProvider).GetMemoryUsage(instanceId, hypervisorName)

	resultMap := make(map[string]interface{})
	resultMap["label"] = "Memory"
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


/**
	디스크 사용률 차트 데이터를 불러옴
*/
func (controller *ZabbixController) GetDiskUsage(ctx echo.Context) error {
	instanceId := ctx.QueryParam("instance_id")
	hypervisorName := ctx.QueryParam("host")

	result, err := service.GetZabbixService(controller.ZabbixSession, controller.OpenstackProvider).GetDiskUsage(instanceId, hypervisorName)

	resultMap := make(map[string]interface{})
	resultMap["label"] = "Disk"
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

/**
	CPU Load Average 차트 데이터를 불러옴 (1분 단위, 5분 단위, 15분 단위)
*/
func (controller *ZabbixController) GetCpuLoadAverage(ctx echo.Context) error {
	instanceId := ctx.QueryParam("instance_id")
	hypervisorName := ctx.QueryParam("host")

	resultInterval1, err := service.GetZabbixService(controller.ZabbixSession, controller.OpenstackProvider).GetCpuLoadAverage(instanceId, hypervisorName, 1)
	resultInterval5, err := service.GetZabbixService(controller.ZabbixSession, controller.OpenstackProvider).GetCpuLoadAverage(instanceId, hypervisorName, 5)
	resultInterval15, err := service.GetZabbixService(controller.ZabbixSession, controller.OpenstackProvider).GetCpuLoadAverage(instanceId, hypervisorName, 15)

	resultMapInterval1 := make(map[string]interface{})
	resultMapInterval5 := make(map[string]interface{})
	resultMapInterval15 := make(map[string]interface{})

	resultMapInterval1["label"] = "1M"
	resultMapInterval1["data"] = resultInterval1

	resultMapInterval5["label"] = "5M"
	resultMapInterval5["data"] = resultInterval5

	resultMapInterval15["label"] = "15M"
	resultMapInterval15["data"] = resultInterval15

	resultList := make([]interface{}, 3)
	resultList[0] = resultMapInterval1
	resultList[1] = resultMapInterval5
	resultList[2] = resultMapInterval15

	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get CPU usage.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", resultList)
	}
	return nil
}


func (controller *ZabbixController) GetDiskIORate(ctx echo.Context) error {
	instanceId := ctx.QueryParam("instance_id")
	hypervisorName := ctx.QueryParam("host")

	resultReadRate, err := service.GetZabbixService(controller.ZabbixSession, controller.OpenstackProvider).GetDiskReadRate(instanceId, hypervisorName)
	resultWriteRate, err := service.GetZabbixService(controller.ZabbixSession, controller.OpenstackProvider).GetDiskWriteRate(instanceId, hypervisorName)

	resultMapReadRate := make(map[string]interface{})
	resultMapWriteRate := make(map[string]interface{})

	resultMapReadRate["label"] = "Disk read"
	resultMapReadRate["data"] = resultReadRate

	resultMapWriteRate["label"] = "Disk write"
	resultMapWriteRate["data"] = resultWriteRate

	resultList := make([]interface{}, 2)
	resultList[0] = resultMapReadRate
	resultList[1] = resultMapWriteRate

	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get CPU usage.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", resultList)
	}
	return nil
}


func (controller *ZabbixController) GetNetworkIOBytes(ctx echo.Context) error {
	instanceId := ctx.QueryParam("instance_id")
	hypervisorName := ctx.QueryParam("host")

	resultReceivedBytes, err := service.GetZabbixService(controller.ZabbixSession, controller.OpenstackProvider).GetNetworkBitReceived(instanceId, hypervisorName)
	resultSentBytes, err := service.GetZabbixService(controller.ZabbixSession, controller.OpenstackProvider).GetNetworkBitSent(instanceId, hypervisorName)

	resultMapReceivedBytes := make(map[string]interface{})
	resultMapSentBytes := make(map[string]interface{})

	resultMapReceivedBytes["label"] = "In"
	resultMapReceivedBytes["data"] = resultReceivedBytes

	resultMapSentBytes["label"] = "Out"
	resultMapSentBytes["data"] = resultSentBytes

	resultList := make([]interface{}, 2)
	resultList[0] = resultMapReceivedBytes
	resultList[1] = resultMapSentBytes

	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get CPU usage.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", resultList)
	}
	return nil
}
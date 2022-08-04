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
	ZabbixSession     *zabbix.Session
	OpenstackProvider *gophercloud.ProviderClient
}

func GetZabbixController(zabbixSession *zabbix.Session, openstackProvider *gophercloud.ProviderClient) *ZabbixController {
	return &ZabbixController{
		ZabbixSession:     zabbixSession,
		OpenstackProvider: openstackProvider,
	}
}

// GetCpuUsage
//  @tags         IaaS
//  @Summary      자빅스를 통해 모니터링되는 CPU 사용량 정보 가져오기
//  @Description  자빅스를 통해 모니터링되는 CPU 사용량 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        instance_id  query     string  true   "인스턴스 아이디 정보를 주입한다."  example(44b801fd-018d-45ff-87bf-cefa3c4b404f)
//  @Param        host         query     string  false  "호스트 정보를 주입한다."       example(null)
//  @Success      200          {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/iaas/instance/cpu/usage [get]
func (controller *ZabbixController) GetCpuUsage(ctx echo.Context) error {
	result, err := service.GetZabbixService(controller.ZabbixSession, controller.OpenstackProvider).GetCpuUsage(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get CPU usage.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", result)
	}
	return nil
}

// GetMemoryUsage
//  @tags         IaaS
//  @Summary      자빅스를 통해 모니터링되는 Memory 사용량 정보 가져오기
//  @Description  자빅스를 통해 모니터링되는 Memory 사용량 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        instance_id  query     string  true   "인스턴스 아이디 정보를 주입한다."  example(44b801fd-018d-45ff-87bf-cefa3c4b404f)
//  @Param        host         query     string  false  "호스트 정보를 주입한다."       example(null)
//  @Success      200          {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/iaas/instance/memory/usage [get]
func (controller *ZabbixController) GetMemoryUsage(ctx echo.Context) error {
	result, err := service.GetZabbixService(controller.ZabbixSession, controller.OpenstackProvider).GetMemoryUsage(ctx)
	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get CPU usage.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", result)
	}
	return nil
}

// GetDiskUsage
//  @tags         IaaS
//  @Summary      자빅스를 통해 모니터링되는 Disk 사용량 정보 가져오기
//  @Description  자빅스를 통해 모니터링되는 Disk 사용량 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        instance_id  query     string  true   "인스턴스 아이디 정보를 주입한다."  example(44b801fd-018d-45ff-87bf-cefa3c4b404f)
//  @Param        host         query     string  false  "호스트 정보를 주입한다."       example(null)
//  @Success      200          {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/iaas/instance/disk/usage [get]
func (controller *ZabbixController) GetDiskUsage(ctx echo.Context) error {
	result, err := service.GetZabbixService(controller.ZabbixSession, controller.OpenstackProvider).GetDiskUsage(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get CPU usage.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", result)
	}
	return nil
}

// GetCpuLoadAverage
//  @tags         IaaS
//  @Summary      자빅스를 통해 모니터링되는 CPU Load Average 정보 가져오기
//  @Description  자빅스를 통해 모니터링되는 CPU Load Average 사용량 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        instance_id  query     string  true   "인스턴스 아이디 정보를 주입한다."  example(44b801fd-018d-45ff-87bf-cefa3c4b404f)
//  @Param        host         query     string  false  "호스트 정보를 주입한다."       example(null)
//  @Success      200          {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/iaas/instance/cpu/load/average [get]
func (controller *ZabbixController) GetCpuLoadAverage(ctx echo.Context) error {
	resultInterval1m, err := service.GetZabbixService(controller.ZabbixSession, controller.OpenstackProvider).GetCpuLoadAverage(ctx, 1)
	resultInterval5m, err := service.GetZabbixService(controller.ZabbixSession, controller.OpenstackProvider).GetCpuLoadAverage(ctx, 5)
	resultInterval15m, err := service.GetZabbixService(controller.ZabbixSession, controller.OpenstackProvider).GetCpuLoadAverage(ctx, 15)

	resultMapInterval1m := make(map[string]interface{})
	resultMapInterval5m := make(map[string]interface{})
	resultMapInterval15m := make(map[string]interface{})

	resultMapInterval1m["label"] = "1M"
	resultMapInterval1m["data"] = resultInterval1m

	resultMapInterval5m["label"] = "5M"
	resultMapInterval5m["data"] = resultInterval5m

	resultMapInterval15m["label"] = "15M"
	resultMapInterval15m["data"] = resultInterval15m

	resultList := make([]interface{}, 3)
	resultList[0] = resultMapInterval1m
	resultList[1] = resultMapInterval5m
	resultList[2] = resultMapInterval15m

	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get CPU usage.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", resultList)
	}
	return nil
}

// GetDiskIORate
//  @tags         IaaS
//  @Summary      자빅스를 통해 모니터링되는 Disk IO Rate 정보 가져오기
//  @Description  자빅스를 통해 모니터링되는 Disk IO Rate 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        instance_id  query     string  true   "인스턴스 아이디 정보를 주입한다."  example(44b801fd-018d-45ff-87bf-cefa3c4b404f)
//  @Param        host         query     string  false  "호스트 정보를 주입한다."       example(null)
//  @Success      200          {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/iaas/instance/disk/io/rate [get]
func (controller *ZabbixController) GetDiskIORate(ctx echo.Context) error {
	result, err := service.GetZabbixService(controller.ZabbixSession, controller.OpenstackProvider).GetDiskIORate(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get CPU usage.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", result)
	}
	return nil
}

// GetNetworkIOBytes
//  @tags         IaaS
//  @Summary      자빅스를 통해 모니터링되는 Network IO Bytes 정보 가져오기
//  @Description  자빅스를 통해 모니터링되는 Network IO Bytes 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        instance_id  query     string  true   "인스턴스 아이디 정보를 주입한다."  example(44b801fd-018d-45ff-87bf-cefa3c4b404f)
//  @Param        host         query     string  false  "호스트 정보를 주입한다."       example(null)
//  @Success      200          {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/iaas/instance/network/io/bytes [get]
func (controller *ZabbixController) GetNetworkIOBytes(ctx echo.Context) error {
	resultReceivedBytes, err := service.GetZabbixService(controller.ZabbixSession, controller.OpenstackProvider).GetNetworkBitReceived(ctx)
	resultSentBytes, err := service.GetZabbixService(controller.ZabbixSession, controller.OpenstackProvider).GetNetworkBitSent(ctx)

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
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get CPU usage.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", resultList)
	}
	return nil
}

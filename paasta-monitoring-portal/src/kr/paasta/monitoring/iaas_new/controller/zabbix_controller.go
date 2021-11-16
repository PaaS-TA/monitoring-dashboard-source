package controller

import (
	"github.com/cavaliercoder/go-zabbix"
	"kr/paasta/monitoring/iaas_new/model"
	"kr/paasta/monitoring/iaas_new/service"
	"kr/paasta/monitoring/utils"
	"net/http"
)

type ZabbixController struct {
	ZabbixSession *zabbix.Session
	OpenstackProvider model.OpenstackProvider
}

func NewZabbixController(zabbixSession *zabbix.Session, openstackProvider model.OpenstackProvider) *ZabbixController {
	return &ZabbixController{
		ZabbixSession: zabbixSession,
		OpenstackProvider: openstackProvider,
	}
}

func (zabbix *ZabbixController) GetCpuUsage(w http.ResponseWriter, req *http.Request) {
	instanceId := req.URL.Query()["instance_id"][0]
	hypervisorName := req.URL.Query()["host"][0]

	result, err := service.GetZabbixService(zabbix.ZabbixSession, zabbix.OpenstackProvider).GetCpuUsage(instanceId, hypervisorName, req)

	resultMap := make(map[string]interface{})
	resultMap["label"] = "CPU"
	resultMap["data"] = result

	resultList := make([]interface{}, 1)
	resultList[0] = resultMap

	if err != nil {
		utils.Logger.Error(err)
		model.MonitLogger.Error("GetServerList error :", err)
		errMessage := model.ErrMessage{}
		errMessage["Message"] = err.Error()
		utils.ErrRenderJsonResponse(errMessage, w)
	} else {
		utils.RenderJsonResponse(resultList, w)
	}
}


/**
	메모리 사용률 차트 데이터를 불러옴
 */
func (zabbix *ZabbixController) GetMemoryUsage(w http.ResponseWriter, req *http.Request) {
	instanceId := req.URL.Query()["instance_id"][0]
	hypervisorName := req.URL.Query()["host"][0]

	result, err := service.GetZabbixService(zabbix.ZabbixSession, zabbix.OpenstackProvider).GetMemoryUsage(instanceId, hypervisorName, req)

	resultMap := make(map[string]interface{})
	resultMap["label"] = "Memory"
	resultMap["data"] = result

	resultList := make([]interface{}, 1)
	resultList[0] = resultMap

	if err != nil {
		utils.Logger.Error(err)
		model.MonitLogger.Error("GetServerList error :", err)
		errMessage := model.ErrMessage{}
		errMessage["Message"] = err.Error()
		utils.ErrRenderJsonResponse(errMessage, w)
	} else {
		utils.RenderJsonResponse(resultList, w)
	}
}


/**
	디스크 사용률 차트 데이터를 불러옴
 */
func (zabbix *ZabbixController) GetDiskUsage(w http.ResponseWriter, req *http.Request) {
	instanceId := req.URL.Query()["instance_id"][0]
	hypervisorName := req.URL.Query()["host"][0]

	result, err := service.GetZabbixService(zabbix.ZabbixSession, zabbix.OpenstackProvider).GetDiskUsage(instanceId, hypervisorName, req)

	resultMap := make(map[string]interface{})
	resultMap["label"] = "Disk"
	resultMap["data"] = result

	resultList := make([]interface{}, 1)
	resultList[0] = resultMap

	if err != nil {
		utils.Logger.Error(err)
		model.MonitLogger.Error("GetServerList error :", err)
		errMessage := model.ErrMessage{}
		errMessage["Message"] = err.Error()
		utils.ErrRenderJsonResponse(errMessage, w)
	} else {
		utils.RenderJsonResponse(resultList, w)
	}
}

/**
	CPU Load Average 차트 데이터를 불러옴 (1분 단위, 5분 단위, 15분 단위)
 */
func (zabbix *ZabbixController) GetCpuLoadAverage(w http.ResponseWriter, req *http.Request) {
	instanceId := req.URL.Query()["instance_id"][0]
	hypervisorName := req.URL.Query()["host"][0]

	resultInterval1, err := service.GetZabbixService(zabbix.ZabbixSession, zabbix.OpenstackProvider).GetCpuLoadAverage(instanceId, hypervisorName, req, 1)
	resultInterval5, err := service.GetZabbixService(zabbix.ZabbixSession, zabbix.OpenstackProvider).GetCpuLoadAverage(instanceId, hypervisorName, req, 5)
	resultInterval15, err := service.GetZabbixService(zabbix.ZabbixSession, zabbix.OpenstackProvider).GetCpuLoadAverage(instanceId, hypervisorName, req, 15)

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
		utils.Logger.Error(err)
		model.MonitLogger.Error("GetServerList error :", err)
		errMessage := model.ErrMessage{}
		errMessage["Message"] = err.Error()
		utils.ErrRenderJsonResponse(errMessage, w)
	} else {
		utils.RenderJsonResponse(resultList, w)
	}
}


func (zabbix *ZabbixController) GetDiskIORate(w http.ResponseWriter, req *http.Request) {
	instanceId := req.URL.Query()["instance_id"][0]
	hypervisorName := req.URL.Query()["host"][0]

	resultReadRate, err := service.GetZabbixService(zabbix.ZabbixSession, zabbix.OpenstackProvider).GetDiskReadRate(instanceId, hypervisorName, req)
	resultWriteRate, err := service.GetZabbixService(zabbix.ZabbixSession, zabbix.OpenstackProvider).GetDiskWriteRate(instanceId, hypervisorName, req)

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
		utils.Logger.Error(err)
		model.MonitLogger.Error("GetServerList error :", err)
		errMessage := model.ErrMessage{}
		errMessage["Message"] = err.Error()
		utils.ErrRenderJsonResponse(errMessage, w)
	} else {
		utils.RenderJsonResponse(resultList, w)
	}
}


func (zabbix *ZabbixController) GetNetworkIOBytes(w http.ResponseWriter, req *http.Request) {
	instanceId := req.URL.Query()["instance_id"][0]
	hypervisorName := req.URL.Query()["host"][0]

	resultReceivedBytes, err := service.GetZabbixService(zabbix.ZabbixSession, zabbix.OpenstackProvider).GetNetworkBitReceived(instanceId, hypervisorName, req)
	resultSentBytes, err := service.GetZabbixService(zabbix.ZabbixSession, zabbix.OpenstackProvider).GetNetworkBitSent(instanceId, hypervisorName, req)

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
		utils.Logger.Error(err)
		model.MonitLogger.Error("GetServerList error :", err)
		errMessage := model.ErrMessage{}
		errMessage["Message"] = err.Error()
		utils.ErrRenderJsonResponse(errMessage, w)
	} else {
		utils.RenderJsonResponse(resultList, w)
	}
}
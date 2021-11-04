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
	//data, _ := ioutil.ReadAll(r.Body)
	instanceId := req.FormValue(":instance_id")

	result, err := service.GetZabbixService(zabbix.ZabbixSession, zabbix.OpenstackProvider).GetCpuUsage(instanceId, req)

	resultMap := make(map[string]interface{})
	resultMap["label"] = "CPU"
	resultMap["data"] = result

	resultList := make([]interface{}, 1)
	resultList[0] = resultMap

	if err != nil {
		model.MonitLogger.Error("GetServerList error :", err)
		errMessage := model.ErrMessage{}
		errMessage["Message"] = err.Error()
		utils.ErrRenderJsonResponse(errMessage, w)
	} else {
		utils.RenderJsonResponse(resultList, w)
	}
}


func (zabbix *ZabbixController) GetMemoryUsage(w http.ResponseWriter, req *http.Request) {
	//data, _ := ioutil.ReadAll(r.Body)
	instanceId := req.FormValue(":instance_id")

	result, err := service.GetZabbixService(zabbix.ZabbixSession, zabbix.OpenstackProvider).GetMemoryUsage(instanceId, req)

	resultMap := make(map[string]interface{})
	resultMap["label"] = "Memory"
	resultMap["data"] = result

	resultList := make([]interface{}, 1)
	resultList[0] = resultMap

	if err != nil {
		model.MonitLogger.Error("GetServerList error :", err)
		errMessage := model.ErrMessage{}
		errMessage["Message"] = err.Error()
		utils.ErrRenderJsonResponse(errMessage, w)
	} else {
		utils.RenderJsonResponse(resultList, w)
	}
}


func (zabbix *ZabbixController) GetDiskUsage(w http.ResponseWriter, req *http.Request) {
	//data, _ := ioutil.ReadAll(r.Body)
	instanceId := req.FormValue(":instance_id")

	result, err := service.GetZabbixService(zabbix.ZabbixSession, zabbix.OpenstackProvider).GetDiskUsage(instanceId, req)

	resultMap := make(map[string]interface{})
	resultMap["label"] = "Disk"
	resultMap["data"] = result

	resultList := make([]interface{}, 1)
	resultList[0] = resultMap

	if err != nil {
		model.MonitLogger.Error("GetServerList error :", err)
		errMessage := model.ErrMessage{}
		errMessage["Message"] = err.Error()
		utils.ErrRenderJsonResponse(errMessage, w)
	} else {
		utils.RenderJsonResponse(resultList, w)
	}
}

func (zabbix *ZabbixController) GetCpuLoadAverage(w http.ResponseWriter, req *http.Request) {
	//data, _ := ioutil.ReadAll(r.Body)
	instanceId := req.FormValue(":instance_id")

	resultInterval1, err := service.GetZabbixService(zabbix.ZabbixSession, zabbix.OpenstackProvider).GetCpuLoadAverage(instanceId, req, 1)
	resultInterval5, err := service.GetZabbixService(zabbix.ZabbixSession, zabbix.OpenstackProvider).GetCpuLoadAverage(instanceId, req, 5)
	resultInterval15, err := service.GetZabbixService(zabbix.ZabbixSession, zabbix.OpenstackProvider).GetCpuLoadAverage(instanceId, req, 15)

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
		model.MonitLogger.Error("GetServerList error :", err)
		errMessage := model.ErrMessage{}
		errMessage["Message"] = err.Error()
		utils.ErrRenderJsonResponse(errMessage, w)
	} else {
		utils.RenderJsonResponse(resultList, w)
	}
}


func (zabbix *ZabbixController) GetDiskIORate(w http.ResponseWriter, req *http.Request) {
	//data, _ := ioutil.ReadAll(r.Body)
	instanceId := req.FormValue(":instance_id")

	resultReadRate, err := service.GetZabbixService(zabbix.ZabbixSession, zabbix.OpenstackProvider).GetDiskReadRate(instanceId, req)
	resultWriteRate, err := service.GetZabbixService(zabbix.ZabbixSession, zabbix.OpenstackProvider).GetDiskWriteRate(instanceId, req)

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
		model.MonitLogger.Error("GetServerList error :", err)
		errMessage := model.ErrMessage{}
		errMessage["Message"] = err.Error()
		utils.ErrRenderJsonResponse(errMessage, w)
	} else {
		utils.RenderJsonResponse(resultList, w)
	}
}

func (zabbix *ZabbixController) GetNetworkIOBytes(w http.ResponseWriter, req *http.Request) {
	//data, _ := ioutil.ReadAll(r.Body)
	instanceId := req.FormValue(":instance_id")

	resultReceivedBytes, err := service.GetZabbixService(zabbix.ZabbixSession, zabbix.OpenstackProvider).GetNetworkBitReceived(instanceId, req)
	resultSentBytes, err := service.GetZabbixService(zabbix.ZabbixSession, zabbix.OpenstackProvider).GetNetworkBitSent(instanceId, req)

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
		model.MonitLogger.Error("GetServerList error :", err)
		errMessage := model.ErrMessage{}
		errMessage["Message"] = err.Error()
		utils.ErrRenderJsonResponse(errMessage, w)
	} else {
		utils.RenderJsonResponse(resultList, w)
	}
}
package service

import (
	"fmt"
	"github.com/cavaliercoder/go-zabbix"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"kr/paasta/monitoring/iaas_new/model"
	"kr/paasta/monitoring/utils"
	"net/http"
	"strconv"
	"zabbix-client/common"
	"zabbix-client/history"
	"zabbix-client/host"
	"zabbix-client/item"
)

type ZabbixService struct {
	ZabbixSession *zabbix.Session
	OpenstackProvider model.OpenstackProvider
}

func GetZabbixService(zabbixSession *zabbix.Session, openstackProvider model.OpenstackProvider) *ZabbixService {
	return &ZabbixService{
		ZabbixSession: zabbixSession,
		OpenstackProvider: openstackProvider,
	}
}


func (zabbixService *ZabbixService) getHostIpAddress(req *http.Request, instanceId string) string {
	var ipAddress string

	provider, _, _ := utils.GetOpenstackProvider(req)
	computeClient, _ := utils.GetComputeClient(provider, zabbixService.OpenstackProvider.Region)

	result, _ := servers.Get(computeClient, instanceId).Extract()
	addressList := result.Addresses

	for _, address := range addressList {
		dataArr := address.([]interface{})
		dataMap := dataArr[0].(map[string]interface{})
		ipAddress = dataMap["addr"].(string)
	}
	return ipAddress
}

func (zabbixService *ZabbixService) getZabbixHostDetail(ipAddr string) (zabbix.Host, error) {
	// IP주소로 호스트 정보 조회
	hostParams := make(map[string]interface{}, 0)
	filterMap := make(map[string]interface{}, 0)
	filterMap["ip"] = ipAddr
	hostParams["filter"] = filterMap
	result, err := host.GetHostList(zabbixService.ZabbixSession, hostParams)
	if err != nil {
		utils.Logger.Error(err)
	}
	if len(result) == 0 {
		_result := zabbix.Host{}
		return _result, fmt.Errorf("%s is not exist host.", ipAddr)
	}
	return result[0], nil
}



func (zabbixService *ZabbixService) GetCpuUsage(instanceId string, req *http.Request) ([]zabbix.History, error) {

	hostIp := zabbixService.getHostIpAddress(req, instanceId)
	zabbixHost, err := zabbixService.getZabbixHostDetail(hostIp)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	hostIds := make([]string, 1)
	hostIds[0] = zabbixHost.HostID

	itemParams := make(map[string]interface{}, 0)
	keywordArr := make([]string, 1)
	keywordArr[0] = common.SYSTEM_CPU_UTIL
	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIds
	itemResult := item.GetItemList(zabbixService.ZabbixSession, itemParams)

	itemId := strconv.Itoa(itemResult[0].ItemID)
	itemType := itemResult[0].LastValueType

	params := make(map[string]interface{}, 0)
	params["itemId"] = itemId
	params["itemType"] = itemType
	params["offset"] = 10
	result := history.GetHistory(zabbixService.ZabbixSession, params)

	return result, err
}


func (zabbixService *ZabbixService) GetCpuLoadAverage(instanceId string, req *http.Request, interval int) ([]zabbix.History, error) {

	hostIp := zabbixService.getHostIpAddress(req, instanceId)
	zabbixHost, err := zabbixService.getZabbixHostDetail(hostIp)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	hostIds := make([]string, 1)
	hostIds[0] = zabbixHost.HostID


	itemParams := make(map[string]interface{}, 0)
	keywordArr := make([]string, 1)
	switch interval {
	case 1:
		keywordArr[0] = common.CPU_LOAD_AVERAGE_PER_1M
	case 5:
		keywordArr[0] = common.CPU_LOAD_AVERAGE_PER_5M
	case 15:
		keywordArr[0] = common.CPU_LOAD_AVERAGE_PER_15M
	default:
		keywordArr[0] = common.CPU_LOAD_AVERAGE_PER_1M
	}

	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIds
	itemResult := item.GetItemList(zabbixService.ZabbixSession, itemParams)

	itemId := strconv.Itoa(itemResult[0].ItemID)
	itemType := itemResult[0].LastValueType

	params := make(map[string]interface{}, 0)
	params["itemId"] = itemId
	params["itemType"] = itemType
	params["offset"] = 10
	result := history.GetHistory(zabbixService.ZabbixSession, params)

	return result, err
}


func (zabbixService *ZabbixService) GetMemoryUsage(instanceId string, req *http.Request) ([]zabbix.History, error) {
	hostIp := zabbixService.getHostIpAddress(req, instanceId)
	zabbixHost, err := zabbixService.getZabbixHostDetail(hostIp)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	hostIds := make([]string, 1)
	hostIds[0] = zabbixHost.HostID

	itemParams := make(map[string]interface{}, 0)
	keywordArr := make([]string, 1)
	keywordArr[0] = common.VM_MEMORY_UTILIZATION
	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIds
	itemResult := item.GetItemList(zabbixService.ZabbixSession, itemParams)

	itemId := strconv.Itoa(itemResult[0].ItemID)
	itemType := itemResult[0].LastValueType

	params := make(map[string]interface{}, 0)
	params["itemId"] = itemId
	params["itemType"] = itemType
	params["offset"] = 10
	result := history.GetHistory(zabbixService.ZabbixSession, params)

	return result, err
}


func (zabbixService *ZabbixService) GetDiskUsage(instanceId string, req *http.Request) ([]zabbix.History, error) {
	hostIp := zabbixService.getHostIpAddress(req, instanceId)
	zabbixHost, err := zabbixService.getZabbixHostDetail(hostIp)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	hostIds := make([]string, 1)
	hostIds[0] = zabbixHost.HostID


	itemParams := make(map[string]interface{}, 0)
	keywordArr := make([]string, 1)
	keywordArr[0] = common.SPACE_UTILIZATION
	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIds
	itemResult := item.GetItemList(zabbixService.ZabbixSession, itemParams)

	itemId := strconv.Itoa(itemResult[0].ItemID)
	itemType := itemResult[0].LastValueType

	params := make(map[string]interface{}, 0)
	params["itemId"] = itemId
	params["itemType"] = itemType
	params["offset"] = 10
	result := history.GetHistory(zabbixService.ZabbixSession, params)

	return result, err
}

func (zabbixService *ZabbixService) GetDiskReadRate(instanceId string, req *http.Request) ([]zabbix.History, error) {
	hostIp := zabbixService.getHostIpAddress(req, instanceId)
	zabbixHost, err := zabbixService.getZabbixHostDetail(hostIp)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	hostIds := make([]string, 1)
	hostIds[0] = zabbixHost.HostID


	itemParams := make(map[string]interface{}, 0)
	keywordArr := make([]string, 1)
	keywordArr[0] = common.DISK_READ_RATE
	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIds
	itemResult := item.GetItemList(zabbixService.ZabbixSession, itemParams)

	itemId := strconv.Itoa(itemResult[0].ItemID)
	itemType := itemResult[0].LastValueType

	params := make(map[string]interface{}, 0)
	params["itemId"] = itemId
	params["itemType"] = itemType
	params["offset"] = 10
	result := history.GetHistory(zabbixService.ZabbixSession, params)

	return result, err
}

func (zabbixService *ZabbixService) GetDiskWriteRate(instanceId string, req *http.Request) ([]zabbix.History, error) {
	hostIp := zabbixService.getHostIpAddress(req, instanceId)
	zabbixHost, err := zabbixService.getZabbixHostDetail(hostIp)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	hostIds := make([]string, 1)
	hostIds[0] = zabbixHost.HostID


	itemParams := make(map[string]interface{}, 0)
	keywordArr := make([]string, 1)
	keywordArr[0] = common.DISK_WRITE_RATE
	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIds
	itemResult := item.GetItemList(zabbixService.ZabbixSession, itemParams)

	itemId := strconv.Itoa(itemResult[0].ItemID)
	itemType := itemResult[0].LastValueType

	params := make(map[string]interface{}, 0)
	params["itemId"] = itemId
	params["itemType"] = itemType
	params["offset"] = 10
	result := history.GetHistory(zabbixService.ZabbixSession, params)

	return result, err
}


func (zabbixService *ZabbixService) GetNetworkBitReceived(instanceId string, req *http.Request) ([]zabbix.History, error) {
	hostIp := zabbixService.getHostIpAddress(req, instanceId)
	zabbixHost, err := zabbixService.getZabbixHostDetail(hostIp)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	hostIds := make([]string, 1)
	hostIds[0] = zabbixHost.HostID


	itemParams := make(map[string]interface{}, 0)
	keywordArr := make([]string, 1)
	keywordArr[0] = common.NETWORK_INPUT_PACKET
	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIds
	itemResult := item.GetItemList(zabbixService.ZabbixSession, itemParams)

	itemId := strconv.Itoa(itemResult[0].ItemID)
	itemType := itemResult[0].LastValueType

	params := make(map[string]interface{}, 0)
	params["itemId"] = itemId
	params["itemType"] = itemType
	params["offset"] = 10
	result := history.GetHistory(zabbixService.ZabbixSession, params)

	return result, err
}


func (zabbixService *ZabbixService) GetNetworkBitSent(instanceId string, req *http.Request) ([]zabbix.History, error) {
	hostIp := zabbixService.getHostIpAddress(req, instanceId)
	zabbixHost, err := zabbixService.getZabbixHostDetail(hostIp)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	hostIds := make([]string, 1)
	hostIds[0] = zabbixHost.HostID


	itemParams := make(map[string]interface{}, 0)
	keywordArr := make([]string, 1)
	keywordArr[0] = common.NETWORK_OUTPUT_PACKET
	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIds
	itemResult := item.GetItemList(zabbixService.ZabbixSession, itemParams)

	itemId := strconv.Itoa(itemResult[0].ItemID)
	itemType := itemResult[0].LastValueType

	params := make(map[string]interface{}, 0)
	params["itemId"] = itemId
	params["itemType"] = itemType
	params["offset"] = 10
	result := history.GetHistory(zabbixService.ZabbixSession, params)

	return result, err
}
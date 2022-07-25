package iaas

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"paasta-monitoring-api/middlewares/zabbix-client/common"
	"paasta-monitoring-api/middlewares/zabbix-client/history"
	"paasta-monitoring-api/middlewares/zabbix-client/host"
	"paasta-monitoring-api/middlewares/zabbix-client/item"
	"paasta-monitoring-api/middlewares/zabbix-client/lib/go-zabbix"

	"strconv"
)

type ZabbixService struct {
	ZabbixSession  *zabbix.Session
	OpenstackProvider *gophercloud.ProviderClient
}

func GetZabbixService(zabbixSession *zabbix.Session, openstackProvider *gophercloud.ProviderClient) *ZabbixService {
	return &ZabbixService{
		ZabbixSession: zabbixSession,
		OpenstackProvider : openstackProvider,
	}
}


func (service *ZabbixService) GetCpuUsage(ctx echo.Context) (map[string]interface{}, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	instanceId := ctx.QueryParam("instance_id")
	hypervisorName := ctx.QueryParam("host")
	paramMap := make(map[string]interface{})
	if instanceId != "" {
		hostIp := GetOpenstackService(service.OpenstackProvider).GetHostIpAddress(instanceId)
		paramMap["ip"] = hostIp
	}
	if hypervisorName != "" {
		paramMap["host"] = hypervisorName
	}

	zabbixHost, err := service.getZabbixHostDetail(paramMap)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	hostIds := make([]string, 1)
	hostIds[0] = zabbixHost.HostID

	itemParams := make(map[string]interface{}, 0)
	keywordArr := make([]string, 1)
	keywordArr[0] = common.SYSTEM_CPU_UTIL
	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIds
	itemResult, err := item.GetItemList(service.ZabbixSession, itemParams)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	itemId := strconv.Itoa(itemResult[0].ItemID)
	itemType := itemResult[0].LastValueType

	params := make(map[string]interface{}, 0)
	params["itemId"] = itemId
	params["itemType"] = itemType
	params["offset"] = 10
	result, err := history.GetHistory(service.ZabbixSession, params)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	resultMap := make(map[string]interface{})
	resultMap["label"] = "CPU"
	resultMap["data"] = result

	return resultMap, err
}


func (service *ZabbixService) GetCpuLoadAverage(ctx echo.Context, interval int) ([]zabbix.History, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	instanceId := ctx.QueryParam("instance_id")
	hypervisorName := ctx.QueryParam("host")

	paramMap := make(map[string]interface{})
	if instanceId != "" {
		hostIp := GetOpenstackService(service.OpenstackProvider).GetHostIpAddress(instanceId)
		paramMap["ip"] = hostIp
	}
	if hypervisorName != "" {
		paramMap["host"] = hypervisorName
	}

	zabbixHost, err := service.getZabbixHostDetail(paramMap)
	if err != nil {
		logger.Error(err)
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
	itemResult, err := item.GetItemList(service.ZabbixSession, itemParams)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	itemId := strconv.Itoa(itemResult[0].ItemID)
	itemType := itemResult[0].LastValueType

	params := make(map[string]interface{}, 0)
	params["itemId"] = itemId
	params["itemType"] = itemType
	params["offset"] = 10
	result, err := history.GetHistory(service.ZabbixSession, params)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return result, err
}


func (service *ZabbixService) GetMemoryUsage(ctx echo.Context) ([]interface{}, error) {
	instanceId := ctx.QueryParam("instance_id")
	hypervisorName := ctx.QueryParam("host")

	paramMap := make(map[string]interface{})
	if instanceId != "" {
		hostIp := GetOpenstackService(service.OpenstackProvider).GetHostIpAddress(instanceId)
		paramMap["ip"] = hostIp
	}
	if hypervisorName != "" {
		paramMap["host"] = hypervisorName
	}

	zabbixHost, err := service.getZabbixHostDetail(paramMap)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}

	hostIds := make([]string, 1)
	hostIds[0] = zabbixHost.HostID

	itemParams := make(map[string]interface{}, 0)
	keywordArr := make([]string, 1)
	keywordArr[0] = common.VM_MEMORY_UTILIZATION
	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIds
	itemResult, err := item.GetItemList(service.ZabbixSession, itemParams)
	if err != nil {
		return nil, err
	}

	itemId := strconv.Itoa(itemResult[0].ItemID)
	itemType := itemResult[0].LastValueType

	params := make(map[string]interface{}, 0)
	params["itemId"] = itemId
	params["itemType"] = itemType
	params["offset"] = 10
	result, err := history.GetHistory(service.ZabbixSession, params)

	resultMap := make(map[string]interface{})
	resultMap["label"] = "Memory"
	resultMap["data"] = result

	resultList := make([]interface{}, 1)
	resultList[0] = resultMap

	return resultList, err
}


func (service *ZabbixService) GetDiskUsage(ctx echo.Context) ([]interface{}, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	instanceId := ctx.QueryParam("instance_id")
	hypervisorName := ctx.QueryParam("host")

	paramMap := make(map[string]interface{})
	if instanceId != "" {
		hostIp := GetOpenstackService(service.OpenstackProvider).GetHostIpAddress(instanceId)
		paramMap["ip"] = hostIp
	}
	if hypervisorName != "" {
		paramMap["host"] = hypervisorName
	}

	zabbixHost, err := service.getZabbixHostDetail(paramMap)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	hostIds := make([]string, 1)
	hostIds[0] = zabbixHost.HostID

	itemParams := make(map[string]interface{}, 0)
	keywordArr := make([]string, 1)
	keywordArr[0] = common.SPACE_UTILIZATION
	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIds
	itemResult, err := item.GetItemList(service.ZabbixSession, itemParams)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	itemId := strconv.Itoa(itemResult[0].ItemID)
	itemType := itemResult[0].LastValueType

	params := make(map[string]interface{}, 0)
	params["itemId"] = itemId
	params["itemType"] = itemType
	params["offset"] = 10
	result, err := history.GetHistory(service.ZabbixSession, params)

	resultMap := make(map[string]interface{})
	resultMap["label"] = "Disk"
	resultMap["data"] = result

	resultList := make([]interface{}, 1)
	resultList[0] = resultMap

	return resultList, err
}


/**
	GetDiskIORate()
		Return data integrating disk read rate data & disk write rate date
 */
func (service *ZabbixService) GetDiskIORate(ctx echo.Context) ([]interface{}, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	resultReadRate, err := service.GetDiskReadRate(ctx)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	resultWriteRate, err := service.GetDiskWriteRate(ctx)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	resultList := make([]interface{}, 2)
	resultList[0] = resultReadRate
	resultList[1] = resultWriteRate

	return resultList, err
}


/**
	GetDiskReadRate
		Return data disk read rate metrics
 */
func (service *ZabbixService) GetDiskReadRate(ctx echo.Context) (map[string]interface{}, error) {
	instanceId := ctx.QueryParam("instance_id")
	hypervisorName := ctx.QueryParam("host")

	paramMap := make(map[string]interface{})
	if instanceId != "" {
		hostIp := GetOpenstackService(service.OpenstackProvider).GetHostIpAddress(instanceId)
		paramMap["ip"] = hostIp
	}
	if hypervisorName != "" {
		paramMap["host"] = hypervisorName
	}

	zabbixHost, err := service.getZabbixHostDetail(paramMap)
	if err != nil {
		return nil, err
	}

	hostIds := make([]string, 1)
	hostIds[0] = zabbixHost.HostID

	itemParams := make(map[string]interface{}, 0)
	keywordArr := make([]string, 1)
	keywordArr[0] = common.DISK_READ_RATE
	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIds
	itemResult, err := item.GetItemList(service.ZabbixSession, itemParams)
	if err != nil {
		return nil, err
	}

	itemId := strconv.Itoa(itemResult[0].ItemID)
	itemType := itemResult[0].LastValueType

	params := make(map[string]interface{}, 0)
	params["itemId"] = itemId
	params["itemType"] = itemType
	params["offset"] = 10
	result, err := history.GetHistory(service.ZabbixSession, params)

	resultMapReadRate := make(map[string]interface{})

	resultMapReadRate["label"] = "Disk read"
	resultMapReadRate["data"] = result

	return resultMapReadRate, err
}


/**
	GetDiskWriteRate
		Return data disk write rate metrics
*/
func (service *ZabbixService) GetDiskWriteRate(ctx echo.Context) (map[string]interface{}, error) {
	instanceId := ctx.QueryParam("instance_id")
	hypervisorName := ctx.QueryParam("host")

	paramMap := make(map[string]interface{})
	if instanceId != "" {
		hostIp := GetOpenstackService(service.OpenstackProvider).GetHostIpAddress(instanceId)
		paramMap["ip"] = hostIp
	}
	if hypervisorName != "" {
		paramMap["host"] = hypervisorName
	}

	zabbixHost, err := service.getZabbixHostDetail(paramMap)
	if err != nil {
		return nil, err
	}

	hostIds := make([]string, 1)
	hostIds[0] = zabbixHost.HostID

	itemParams := make(map[string]interface{}, 0)
	keywordArr := make([]string, 1)
	keywordArr[0] = common.DISK_WRITE_RATE
	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIds
	itemResult, err := item.GetItemList(service.ZabbixSession, itemParams)
	if err != nil {
		return nil, err
	}

	itemId := strconv.Itoa(itemResult[0].ItemID)
	itemType := itemResult[0].LastValueType

	params := make(map[string]interface{}, 0)
	params["itemId"] = itemId
	params["itemType"] = itemType
	params["offset"] = 10
	result, err := history.GetHistory(service.ZabbixSession, params)

	resultMapWriteRate := make(map[string]interface{})
	resultMapWriteRate["label"] = "Disk write"
	resultMapWriteRate["data"] = result

	return resultMapWriteRate, err
}


func (service *ZabbixService) GetNetworkBitReceived(ctx echo.Context) ([]zabbix.History, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	instanceId := ctx.QueryParam("instance_id")
	hypervisorName := ctx.QueryParam("host")
	paramMap := make(map[string]interface{})
	if instanceId != "" {
		hostIp := GetOpenstackService(service.OpenstackProvider).GetHostIpAddress(instanceId)
		paramMap["ip"] = hostIp
	}
	if hypervisorName != "" {
		paramMap["host"] = hypervisorName
	}

	zabbixHost, err := service.getZabbixHostDetail(paramMap)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}

	hostIds := make([]string, 1)
	hostIds[0] = zabbixHost.HostID

	itemParams := make(map[string]interface{}, 0)
	keywordArr := make([]string, 1)
	keywordArr[0] = common.NETWORK_INPUT_PACKET
	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIds
	itemResult, err := item.GetItemList(service.ZabbixSession, itemParams)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	itemId := strconv.Itoa(itemResult[0].ItemID)
	itemType := itemResult[0].LastValueType

	params := make(map[string]interface{}, 0)
	params["itemId"] = itemId
	params["itemType"] = itemType
	params["offset"] = 10
	result, err := history.GetHistory(service.ZabbixSession, params)

	return result, err
}


func (service *ZabbixService) GetNetworkBitSent(ctx echo.Context) ([]zabbix.History, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	instanceId := ctx.QueryParam("instance_id")
	hypervisorName := ctx.QueryParam("host")
	paramMap := make(map[string]interface{})
	if instanceId != "" {
		hostIp := GetOpenstackService(service.OpenstackProvider).GetHostIpAddress(instanceId)
		paramMap["ip"] = hostIp
	}
	if hypervisorName != "" {
		paramMap["host"] = hypervisorName
	}

	zabbixHost, err := service.getZabbixHostDetail(paramMap)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	hostIds := make([]string, 1)
	hostIds[0] = zabbixHost.HostID

	itemParams := make(map[string]interface{}, 0)
	keywordArr := make([]string, 1)
	keywordArr[0] = common.NETWORK_OUTPUT_PACKET
	itemParams["itemKey"] = keywordArr
	itemParams["hostIds"] = hostIds
	itemResult, err := item.GetItemList(service.ZabbixSession, itemParams)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	itemId := strconv.Itoa(itemResult[0].ItemID)
	itemType := itemResult[0].LastValueType

	params := make(map[string]interface{}, 0)
	params["itemId"] = itemId
	params["itemType"] = itemType
	params["offset"] = 10
	result, err := history.GetHistory(service.ZabbixSession, params)

	return result, err
}


/**
	Zabbix 호스트 정보를 조회
		- IP주소나 호스트 이름으로 조회 가능
*/
func (service *ZabbixService) getZabbixHostDetail(paramMap map[string]interface{}) (zabbix.Host, error) {
	// IP주소로 호스트 정보 조회
	hostParams := make(map[string]interface{}, 0)
	filterMap := make(map[string]interface{}, 0)

	ipAddr, ok := paramMap["ip"].(string)
	if ok {
		filterMap["ip"] = ipAddr
	}
	hostName, ok := paramMap["host"].(string)
	if ok {
		filterMap["host"] = hostName
	}
	hostParams["filter"] = filterMap

	fmt.Println(service.ZabbixSession.URL)

	result, err := host.GetHostList(service.ZabbixSession, hostParams)
	if err != nil {
		fmt.Println(err)
	}
	if len(result) == 0 {
		_result := zabbix.Host{}
		var target string
		if ipAddr == "" {
			target = hostName
		} else {
			target = ipAddr
		}
		return _result, fmt.Errorf("%s is not exist host.", target)
	}
	return result[0], nil
}
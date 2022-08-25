package service

import (
	"encoding/json"
	"github.com/influxdata/influxdb1-client/v2"
	"github.com/gophercloud/gophercloud"
	"monitoring-portal/iaas_new/dao"
	"monitoring-portal/iaas_new/integration"
	"monitoring-portal/iaas_new/model"
	"monitoring-portal/utils"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type ComputeNodeService struct {
	openstackProvider model.OpenstackProvider
	provider          *gophercloud.ProviderClient
	influxClient      client.Client
}

func GetComputeNodeService(openstackProvider model.OpenstackProvider, provider *gophercloud.ProviderClient, influxClient client.Client) *ComputeNodeService {
	return &ComputeNodeService{
		openstackProvider: openstackProvider,
		provider:          provider,
		influxClient:      influxClient,
	}
}

func (n ComputeNodeService) GetComputeNodeSummary(apiRequest model.NodeReq) ([]model.NodeResources, model.ErrMessage) {

	var result []model.NodeResources

	//Compute Node목록 및 Summary정보를 조회한다.
	computeInfoList, err := integration.GetNova(n.openstackProvider, n.provider).GetComputeNodeResources()

	errMsg := utils.GetError().GetCheckErrorMessage(err)
	if err != nil {
		return computeInfoList, errMsg
	}
	//Compute Node의 Status를 조회한다.
	computeResult := computeInfoList
	var errs []model.ErrMessage
	if len(computeInfoList) > 0 {

		for idx, computeNode := range computeInfoList {

			var req model.NodeReq
			req.HostName = computeNode.Hostname
			cpuData, memData, agentForwarderData, agentCollectorData, runningVmsCnt, err := getNodeSummary_Sub(req, n.influxClient)

			if err != nil {
				errs = append(errs, err)
			}

			cpuUsage := utils.GetDataFloatFromInterfaceSingle(cpuData)
			memUsage := utils.GetDataFloatFromInterfaceSingle(memData)

			agentForwarderStatus := utils.GetDataFloatFromInterfaceSingle(agentForwarderData)
			agentCollectorStatus := utils.GetDataFloatFromInterfaceSingle(agentCollectorData)

			computeResult[idx].CpuUsage = cpuUsage
			computeResult[idx].MemUsage = 100 - memUsage
			computeResult[idx].RunningVms = runningVmsCnt
			if agentForwarderStatus == 1 && agentCollectorStatus == 1 {
				computeResult[idx].AgentStatus = "OK"
			} else {
				if agentForwarderStatus != 1 {
					if agentForwarderStatus == 0 {
						computeResult[idx].AgentStatus = "Forwarder Down"
					} else if agentForwarderStatus == -1 {
						computeResult[idx].AgentStatus = "Forwarder UnKnown"
					}
				} else if agentCollectorStatus != 1 {
					if agentCollectorStatus == 0 {
						computeResult[idx].AgentStatus = "Collector Down"
					} else if agentForwarderStatus == -1 {
						computeResult[idx].AgentStatus = "Collector UnKnown"
					}
				}
			}
		}
	}
	//==========================================================================
	// Error가 여러건일 경우 대해 고려해야함.
	if len(errs) > 0 {
		var returnErrMessage string
		for _, err := range errs {
			returnErrMessage = returnErrMessage + " " + err["Message"].(string)
		}
		errMessage := model.ErrMessage{
			"Message": returnErrMessage,
		}
		return nil, errMessage
	}
	//==========================================================================

	//조회 조건
	if apiRequest.HostName != "" {
		for _, compute := range computeInfoList {

			if strings.Contains(compute.Hostname, apiRequest.HostName) {
				result = append(result, compute)
			}
		}
		return result, nil
	} else {
		return computeInfoList, nil
	}

}

//CPU 사용률
func (n ComputeNodeService) GetComputeNodeCpuUsageList(request model.DetailReq) (result []map[string]interface{}, _ model.ErrMessage) {

	cpuUsageResp, err := dao.GetNodeDao(n.influxClient).GetNodeCpuUsageList(request)
	if err != nil {
		model.MonitLogger.Error(err)
		return result, err
	} else {
		cpuUsage, _ := utils.GetResponseConverter().InfluxConverterList(cpuUsageResp, model.METRIC_NAME_CPU_USAGE)

		datamap := cpuUsage[model.RESULT_DATA_NAME].([]map[string]interface{})
		for _, data := range datamap {

			swapFree := data["usage"].(json.Number)
			convertData, _ := strconv.ParseFloat(swapFree.String(), 64)
			//swap 사용률을 구한다. ( 100 -  freeUsage)
			data["usage"] = utils.RoundFloatDigit2(convertData)
		}

		result = append(result, cpuUsage)
		return result, nil
	}
}

//CPU Load Avg_1m
func (n ComputeNodeService) GetComputeNodeCpuLoad1mList(request model.DetailReq) (result []map[string]interface{}, _ model.ErrMessage) {

	cpuLoad1mResp, err := dao.GetNodeDao(n.influxClient).GetNodeCpuLoadList(request, "1m")
	cpuLoad5mResp, err := dao.GetNodeDao(n.influxClient).GetNodeCpuLoadList(request, "5m")
	cpuLoad15mResp, err := dao.GetNodeDao(n.influxClient).GetNodeCpuLoadList(request, "15m")

	if err != nil {
		model.MonitLogger.Error(err)
		return result, err
	} else {
		cpu1mLoad, _ := utils.GetResponseConverter().InfluxConverterList(cpuLoad1mResp, model.METRIC_NAME_CPU_LOAD_1M)
		cpu5mLoad, _ := utils.GetResponseConverter().InfluxConverterList(cpuLoad5mResp, model.METRIC_NAME_CPU_LOAD_5M)
		cpu15mLoad, _ := utils.GetResponseConverter().InfluxConverterList(cpuLoad15mResp, model.METRIC_NAME_CPU_LOAD_15M)

		datamap1m := cpu1mLoad[model.RESULT_DATA_NAME].([]map[string]interface{})
		for _, data := range datamap1m {

			swapFree := data["usage"].(json.Number)
			convertData, _ := strconv.ParseFloat(swapFree.String(), 64)

			data["usage"] = utils.RoundFloatDigit2(convertData)
		}

		datamap5m := cpu5mLoad[model.RESULT_DATA_NAME].([]map[string]interface{})
		for _, data := range datamap5m {

			swapFree := data["usage"].(json.Number)
			convertData, _ := strconv.ParseFloat(swapFree.String(), 64)

			data["usage"] = utils.RoundFloatDigit2(convertData)
		}

		datamap15m := cpu15mLoad[model.RESULT_DATA_NAME].([]map[string]interface{})
		for _, data := range datamap15m {

			swapFree := data["usage"].(json.Number)
			convertData, _ := strconv.ParseFloat(swapFree.String(), 64)

			data["usage"] = utils.RoundFloatDigit2(convertData)
		}
		result = append(result, cpu1mLoad)
		result = append(result, cpu5mLoad)
		result = append(result, cpu15mLoad)

		return result, nil
	}
}

//Memory Swap Usage
func (n ComputeNodeService) GetComputeNodeMemoryUsageList(request model.DetailReq) (result []map[string]interface{}, _ model.ErrMessage) {

	memoryResp, err := dao.GetNodeDao(n.influxClient).GetNodeMemoryUsageList(request)
	if err != nil {
		model.MonitLogger.Error(err)
		return result, err
	} else {
		memoryUsage, _ := utils.GetResponseConverter().InfluxConverterList(memoryResp, model.METRIC_NAME_MEMORY_USAGE)

		datamap := memoryUsage[model.RESULT_DATA_NAME].([]map[string]interface{})

		for _, data := range datamap {
			data["usage"] = 100 - utils.TypeChecker_float64(data["usage"]).(float64)
		}

		result = append(result, memoryUsage)
		return result, nil
	}
}

//Memory Swap Usage
func (n ComputeNodeService) GetComputeNodeSwapUsageList(request model.DetailReq) (result []map[string]interface{}, _ model.ErrMessage) {

	cpuLoadResp, err := dao.GetNodeDao(n.influxClient).GetNodeSwapMemoryFreeUsageList(request)
	if err != nil {
		model.MonitLogger.Error(err)
		return result, err
	} else {
		swapFreeUsage, _ := utils.GetResponseConverter().InfluxConverterList(cpuLoadResp, model.METRIC_NAME_MEMORY_SWAP)

		datamap := swapFreeUsage[model.RESULT_DATA_NAME].([]map[string]interface{})
		for _, data := range datamap {

			swapFree := data["usage"].(json.Number)
			convertData, _ := strconv.ParseFloat(swapFree.String(), 64)
			//swap 사용률을 구한다. ( 100 -  freeUsage)
			data["usage"] = utils.RoundFloatDigit2(100 - convertData)
		}

		result = append(result, swapFreeUsage)
		return result, nil
	}
}

func getNodeSummary_Sub(request model.NodeReq, f client.Client) (map[string]interface{}, map[string]interface{}, map[string]interface{},
	map[string]interface{}, int, model.ErrMessage) {
	var cpuResp, memResp, agentForwarderResp, agentCollectorResp, instanceListResp client.Response

	var errs []model.ErrMessage
	var err model.ErrMessage
	var wg sync.WaitGroup
	wg.Add(5)

	for i := 0; i < 5; i++ {
		go func(wg *sync.WaitGroup, index int) {

			switch index {
			case 0:
				cpuResp, err = dao.GetNodeDao(f).GetNodeCpuUsage(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 1:
				memResp, err = dao.GetNodeDao(f).GetNodeMemoryUsage(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 2:
				agentForwarderResp, err = dao.GetNodeDao(f).GetAgentProcessStatus(request, "forwarder")
				if err != nil {
					errs = append(errs, err)
				}
			case 3:
				agentCollectorResp, err = dao.GetNodeDao(f).GetAgentProcessStatus(request, "collector")
				if err != nil {
					errs = append(errs, err)
				}
			case 4:
				instanceListResp, err = dao.GetNodeDao(f).GetAliveInstanceListByNodename(request, false)
				if err != nil {
					errs = append(errs, err)
				}
			default:
				break
			}
			wg.Done()
		}(&wg, i)
	}
	wg.Wait()

	//==========================================================================
	// Error가 여러건일 경우 대해 고려해야함.
	if len(errs) > 0 {
		var returnErrMessage string
		for _, err := range errs {
			returnErrMessage = returnErrMessage + " " + err["Message"].(string)
		}
		errMessage := model.ErrMessage{
			"Message": returnErrMessage,
		}
		return nil, nil, nil, nil, 0, errMessage
	}
	//==========================================================================
	cpuUsage, _ := utils.GetResponseConverter().InfluxConverter(cpuResp)
	memUsage, _ := utils.GetResponseConverter().InfluxConverter(memResp)
	agentForwarder, _ := utils.GetResponseConverter().InfluxConverter(agentForwarderResp)
	agentCollector, _ := utils.GetResponseConverter().InfluxConverter(agentCollectorResp)
	instanceList, _ := utils.GetResponseConverter().InfluxConverterToMap(instanceListResp)
	var instanceGuidList []string
	//valueList, _ := utils.GetResponseConverter().InfluxConverterToMap(instanceList)
	for _, value := range instanceList {

		instanceGuid := reflect.ValueOf(value["resource_id"]).String()

		if utils.StringArrayDistinct(instanceGuid, instanceGuidList) == false {
			instanceGuidList = append(instanceGuidList, instanceGuid)
		}
	}

	return cpuUsage, memUsage, agentForwarder, agentCollector, len(instanceGuidList), nil
}

//Memory Swap Usage
func (n ComputeNodeService) GetNodeDiskUsageList(request model.DetailReq) (result []map[string]interface{}, _ model.ErrMessage) {

	mountPointResp, err := dao.GetNodeDao(n.influxClient).GetMountPointList(request)
	if err != nil {
		model.MonitLogger.Error(err)
		return result, err
	} else {
		mountPointList, _ := utils.GetResponseConverter().GetMountPointList(mountPointResp)

		var mountPointSelectList []string
		for idx := range mountPointList {
			//Boot Mount Point는 제외
			if strings.Contains(mountPointList[idx], "/boot") == false {
				mountPointSelectList = append(mountPointSelectList, mountPointList[idx])
			}

		}

		for _, value := range mountPointSelectList {
			request.MountPoint = value

			diskResp, _ := dao.GetNodeDao(n.influxClient).GetNodeDiskUsage(request)
			diskUsage, _ := utils.GetResponseConverter().InfluxConverterList(diskResp, value)

			datamap := diskUsage[model.RESULT_DATA_NAME].([]map[string]interface{})
			for _, data := range datamap {

				swapFree := data["usage"].(json.Number)
				convertData, _ := strconv.ParseFloat(swapFree.String(), 64)
				//swap 사용률을 구한다. ( 100 -  freeUsage)
				data["usage"] = utils.RoundFloatDigit2(convertData)
			}

			result = append(result, diskUsage)
		}

		//result = mountPointList

		return result, nil
	}
}

//Disk IO Read Byte
func (n ComputeNodeService) GetNodeDiskIoReadList(request model.DetailReq) (result []map[string]interface{}, _ model.ErrMessage) {

	mountPointResp, err := dao.GetNodeDao(n.influxClient).GetMountPointList(request)
	if err != nil {
		model.MonitLogger.Error(err)
		return result, err
	} else {
		mountPointList, _ := utils.GetResponseConverter().GetMountPointList(mountPointResp)

		var mountPointSelectList []string
		for idx := range mountPointList {
			//Boot Mount Point는 제외
			if strings.Contains(mountPointList[idx], "/boot") == false {
				mountPointSelectList = append(mountPointSelectList, mountPointList[idx])
			}
		}

		for _, value := range mountPointSelectList {
			request.MountPoint = value

			diskResp, _ := dao.GetNodeDao(n.influxClient).GetNodeDiskIoReadKbyte(request)
			diskUsage, _ := utils.GetResponseConverter().InfluxConverterList(diskResp, value)

			datamap := diskUsage[model.RESULT_DATA_NAME].([]map[string]interface{})
			if len(datamap) == 0 {
				continue
			}
			for _, data := range datamap {

				swapFree := data["usage"].(json.Number)
				convertData, _ := strconv.ParseFloat(swapFree.String(), 64)
				//swap 사용률을 구한다. ( 100 -  freeUsage)
				data["usage"] = utils.RoundFloatDigit2(convertData)
			}

			result = append(result, diskUsage)
		}

		return result, nil
	}
}

//Disk IO Write Byte
func (n ComputeNodeService) GetNodeDiskIoWriteList(request model.DetailReq) (result []map[string]interface{}, _ model.ErrMessage) {

	mountPointResp, err := dao.GetNodeDao(n.influxClient).GetMountPointList(request)
	if err != nil {
		model.MonitLogger.Error(err)
		return result, err
	} else {
		mountPointList, _ := utils.GetResponseConverter().GetMountPointList(mountPointResp)

		var mountPointSelectList []string
		for idx := range mountPointList {
			//Boot Mount Point는 제외
			if strings.Contains(mountPointList[idx], "/boot") == false {
				mountPointSelectList = append(mountPointSelectList, mountPointList[idx])
			}

		}

		for _, value := range mountPointSelectList {
			request.MountPoint = value

			diskResp, _ := dao.GetNodeDao(n.influxClient).GetNodeDiskIoWriteKbyte(request)
			diskUsage, _ := utils.GetResponseConverter().InfluxConverterList(diskResp, value)

			datamap := diskUsage[model.RESULT_DATA_NAME].([]map[string]interface{})
			if len(datamap) == 0 {
				continue
			}
			for _, data := range datamap {

				swapFree := data["usage"].(json.Number)
				convertData, _ := strconv.ParseFloat(swapFree.String(), 64)
				//swap 사용률을 구한다. ( 100 -  freeUsage)
				data["usage"] = utils.RoundFloatDigit2(convertData)
			}

			result = append(result, diskUsage)
		}

		return result, nil
	}
}

//Disk IO Write Byte
func (n ComputeNodeService) GetNodeNetworkInOutKByteList(request model.DetailReq) (result []map[string]interface{}, _ model.ErrMessage) {

	networkInEthResp, err := dao.GetNodeDao(n.influxClient).GetNodeNetworkKbyte(request, "in", "en")
	networkInVxResp, err := dao.GetNodeDao(n.influxClient).GetNodeNetworkKbyte(request, "in", "vxlan")

	networkEthOutResp, err := dao.GetNodeDao(n.influxClient).GetNodeNetworkKbyte(request, "out", "en")
	networkVxOutResp, err := dao.GetNodeDao(n.influxClient).GetNodeNetworkKbyte(request, "out", "vxlan")

	if err != nil {
		model.MonitLogger.Error(err)
		return result, err
	} else {
		networkEthInUsage, _ := utils.GetResponseConverter().InfluxConverterList(networkInEthResp, model.METRIC_NAME_NETWORK_ETH_IN)
		networkVxInUsage, _ := utils.GetResponseConverter().InfluxConverterList(networkInVxResp, model.METRIC_NAME_NETWORK_VX_IN)

		networkEthOutUsage, _ := utils.GetResponseConverter().InfluxConverterList(networkEthOutResp, model.METRIC_NAME_NETWORK_ETH_OUT)
		networkVxOutUsage, _ := utils.GetResponseConverter().InfluxConverterList(networkVxOutResp, model.METRIC_NAME_NETWORK_VX_OUT)

		inEthDatamap := networkEthInUsage[model.RESULT_DATA_NAME].([]map[string]interface{})
		for _, data := range inEthDatamap {

			inByte := data["usage"].(json.Number)
			convertData, _ := strconv.ParseFloat(inByte.String(), 64)
			data["usage"] = utils.RoundFloatDigit2(convertData)
		}

		inVxDatamap := networkVxInUsage[model.RESULT_DATA_NAME].([]map[string]interface{})
		for _, data := range inVxDatamap {

			inByte := data["usage"].(json.Number)
			convertData, _ := strconv.ParseFloat(inByte.String(), 64)
			data["usage"] = utils.RoundFloatDigit2(convertData)
		}

		outEthDatamap := networkEthOutUsage[model.RESULT_DATA_NAME].([]map[string]interface{})
		for _, data := range outEthDatamap {

			outByte := data["usage"].(json.Number)
			convertData, _ := strconv.ParseFloat(outByte.String(), 64)

			data["usage"] = utils.RoundFloatDigit2(convertData)
		}

		outVxDatamap := networkVxOutUsage[model.RESULT_DATA_NAME].([]map[string]interface{})
		for _, data := range outVxDatamap {

			outByte := data["usage"].(json.Number)
			convertData, _ := strconv.ParseFloat(outByte.String(), 64)

			data["usage"] = utils.RoundFloatDigit2(convertData)
		}

		result = append(result, networkEthInUsage)
		result = append(result, networkVxInUsage)
		result = append(result, networkEthOutUsage)
		result = append(result, networkVxOutUsage)

		return result, nil
	}
}

//Network In/Out Error
func (n ComputeNodeService) GetNodeNetworkInOutErrorList(request model.DetailReq) (result []map[string]interface{}, _ model.ErrMessage) {

	networkInEthResp, err := dao.GetNodeDao(n.influxClient).GetNodeNetworkError(request, "in", "en")
	networkInVxResp, err := dao.GetNodeDao(n.influxClient).GetNodeNetworkError(request, "in", "vxlan")
	networkOuEthtResp, err := dao.GetNodeDao(n.influxClient).GetNodeNetworkError(request, "out", "en")
	networkOutVxResp, err := dao.GetNodeDao(n.influxClient).GetNodeNetworkError(request, "out", "vxlan")
	if err != nil {
		model.MonitLogger.Error(err)
		return result, err
	} else {
		networkInEthError, _ := utils.GetResponseConverter().InfluxConverterList(networkInEthResp, model.METRIC_NAME_NETWORK_ETH_IN_ERROR)
		networkInVxError, _ := utils.GetResponseConverter().InfluxConverterList(networkInVxResp, model.METRIC_NAME_NETWORK_VX_IN_ERROR)

		networkOutEthError, _ := utils.GetResponseConverter().InfluxConverterList(networkOuEthtResp, model.METRIC_NAME_NETWORK_ETH_OUT_ERROR)
		networkOutVxError, _ := utils.GetResponseConverter().InfluxConverterList(networkOutVxResp, model.METRIC_NAME_NETWORK_VX_OUT_ERROR)

		inEthDatamap := networkInEthError[model.RESULT_DATA_NAME].([]map[string]interface{})
		for _, data := range inEthDatamap {

			inByte := data["usage"].(json.Number)
			convertData, _ := strconv.ParseFloat(inByte.String(), 64)
			data["usage"] = utils.RoundFloatDigit2(convertData)
		}

		inVxDatamap := networkInVxError[model.RESULT_DATA_NAME].([]map[string]interface{})
		for _, data := range inVxDatamap {

			inByte := data["usage"].(json.Number)
			convertData, _ := strconv.ParseFloat(inByte.String(), 64)
			data["usage"] = utils.RoundFloatDigit2(convertData)
		}

		outEthDatamap := networkOutEthError[model.RESULT_DATA_NAME].([]map[string]interface{})
		for _, data := range outEthDatamap {

			outByte := data["usage"].(json.Number)
			convertData, _ := strconv.ParseFloat(outByte.String(), 64)

			data["usage"] = utils.RoundFloatDigit2(convertData)
		}

		outVxDatamap := networkOutVxError[model.RESULT_DATA_NAME].([]map[string]interface{})
		for _, data := range outVxDatamap {

			outByte := data["usage"].(json.Number)
			convertData, _ := strconv.ParseFloat(outByte.String(), 64)

			data["usage"] = utils.RoundFloatDigit2(convertData)
		}

		result = append(result, networkInEthError)
		result = append(result, networkInVxError)
		result = append(result, networkOutEthError)
		result = append(result, networkOutVxError)

		return result, nil
	}
}

//Network Dropped packets
func (n ComputeNodeService) GetNodeNetworkDropPacketList(request model.DetailReq) (result []map[string]interface{}, _ model.ErrMessage) {

	networkInEthResp, err := dao.GetNodeDao(n.influxClient).GetNodeNetworkDropPacket(request, "in", "en")
	networkInVxResp, err := dao.GetNodeDao(n.influxClient).GetNodeNetworkDropPacket(request, "in", "vxlan")

	networkOutEthResp, err := dao.GetNodeDao(n.influxClient).GetNodeNetworkDropPacket(request, "out", "en")
	networkOutVxResp, err := dao.GetNodeDao(n.influxClient).GetNodeNetworkDropPacket(request, "out", "vxlan")

	if err != nil {
		model.MonitLogger.Error(err)
		return result, err
	} else {
		networkInEthError, _ := utils.GetResponseConverter().InfluxConverterList(networkInEthResp, model.METRIC_NAME_NETWORK_ETH_IN_DROPPED_PACKET)
		networkInVxError, _ := utils.GetResponseConverter().InfluxConverterList(networkInVxResp, model.METRIC_NAME_NETWORK_VX_IN_DROPPED_PACKET)

		networkOutEthError, _ := utils.GetResponseConverter().InfluxConverterList(networkOutEthResp, model.METRIC_NAME_NETWORK_ETH_OUT_DROPPED_PACKET)
		networkOutVxError, _ := utils.GetResponseConverter().InfluxConverterList(networkOutVxResp, model.METRIC_NAME_NETWORK_VX_OUT_DROPPED_PACKET)

		inEthDatamap := networkInEthError[model.RESULT_DATA_NAME].([]map[string]interface{})
		for _, data := range inEthDatamap {

			inByte := data["usage"].(json.Number)
			convertData, _ := strconv.ParseFloat(inByte.String(), 64)
			data["usage"] = utils.RoundFloatDigit2(convertData)
		}

		inVxDatamap := networkInVxError[model.RESULT_DATA_NAME].([]map[string]interface{})
		for _, data := range inVxDatamap {

			inByte := data["usage"].(json.Number)
			convertData, _ := strconv.ParseFloat(inByte.String(), 64)
			data["usage"] = utils.RoundFloatDigit2(convertData)
		}

		outEthDatamap := networkOutEthError[model.RESULT_DATA_NAME].([]map[string]interface{})
		for _, data := range outEthDatamap {

			outByte := data["usage"].(json.Number)
			convertData, _ := strconv.ParseFloat(outByte.String(), 64)

			data["usage"] = utils.RoundFloatDigit2(convertData)
		}

		outVxDatamap := networkOutVxError[model.RESULT_DATA_NAME].([]map[string]interface{})
		for _, data := range outVxDatamap {

			outByte := data["usage"].(json.Number)
			convertData, _ := strconv.ParseFloat(outByte.String(), 64)

			data["usage"] = utils.RoundFloatDigit2(convertData)
		}

		result = append(result, networkInEthError)
		result = append(result, networkInVxError)
		result = append(result, networkOutEthError)
		result = append(result, networkOutVxError)

		return result, nil
	}
}

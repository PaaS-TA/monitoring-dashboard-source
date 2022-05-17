package service

import (
	"github.com/gophercloud/gophercloud"
	client "github.com/influxdata/influxdb1-client/v2"
	"monitoring-portal/iaas_new/dao"
	"monitoring-portal/iaas_new/integration"
	"monitoring-portal/iaas_new/model"
	"monitoring-portal/utils"
	"reflect"
	"sort"
	"strings"
	"sync"
)

type ManageNodeService struct {
	openstackProvider model.OpenstackProvider
	provider          *gophercloud.ProviderClient
	influxClient      client.Client
}

func GetManageNodeService(openstackProvider model.OpenstackProvider, provider *gophercloud.ProviderClient, influxClient client.Client) *ManageNodeService {
	return &ManageNodeService{
		openstackProvider: openstackProvider,
		provider:          provider,
		influxClient:      influxClient,
	}
}

func (n ManageNodeService) GetNodeList() (result []string, err error) {

	manageNodeNameList := getManageNodeName(n.openstackProvider, n.provider, n.influxClient)
	computeInfoList, err := integration.GetNova(n.openstackProvider, n.provider).GetComputeNodeResources()

	if err != nil {
		model.MonitLogger.Error(err)
		return result, err
	}

	for _, data := range computeInfoList {
		result = append(result, data.Hostname)
	}

	result = append(result, manageNodeNameList...)

	sort.Strings(result)

	return result, err

}

//CPU Top Process List
func (n ManageNodeService) GetTopProcessListByCpu(request model.DetailReq) (result []model.TopProcess, _ model.ErrMessage) {

	var topProcessList []model.TopProcess
	topProcess, err := dao.GetNodeDao(n.influxClient).GetNodeTopProcessByCpu(request)
	if err != nil {
		model.MonitLogger.Error(err)
		return result, err
	} else {

		topProcessResult, _ := utils.GetResponseConverter().InfluxConverter4TopProcess(topProcess)
		usage := []map[float64]string{}
		var usageArr []float64

		for _, value := range topProcessResult {
			aa := value.([]map[string]interface{})

			for _, vv := range aa {
				processData := map[float64]string{utils.TypeChecker_float64(vv["usage"]).(float64): utils.TypeChecker_string(vv["process_name"]).(string)}

				usage = append(usage, processData)
				usageArr = append(usageArr, utils.TypeChecker_float64(vv["usage"]).(float64))
			}
		}

		sort.Sort(sort.Reverse(sort.Float64Slice(usageArr)))

		for idx := range usageArr {

			for _, v := range usage {
				if v[usageArr[idx]] != "" {

					var topProcess model.TopProcess

					topProcess.Index = idx + 1
					topProcess.ProcessName = v[usageArr[idx]]
					topProcess.Usage = usageArr[idx]

					topProcessList = append(topProcessList, topProcess)

				}
			}
			if idx > 8 {
				break
			}

		}
		result = topProcessList
		return result, nil
	}
}

//CPU Top Process List
func (n ManageNodeService) GetTopProcessListByMemory(request model.DetailReq) (result []model.TopProcess, _ model.ErrMessage) {

	var topProcessList []model.TopProcess
	topProcess, err := dao.GetNodeDao(n.influxClient).GetNodeTopProcessByMemory(request)
	if err != nil {
		model.MonitLogger.Error(err)
		return result, err
	} else {

		topProcessResult, _ := utils.GetResponseConverter().InfluxConverter4TopProcess(topProcess)
		usage := []map[float64]string{}
		var usageArr []float64

		for _, value := range topProcessResult {
			aa := value.([]map[string]interface{})

			for _, vv := range aa {
				processData := map[float64]string{utils.TypeChecker_float64(vv["usage"]).(float64): utils.TypeChecker_string(vv["process_name"]).(string)}

				usage = append(usage, processData)
				//fmt.Println("processData:", processData)
				//fmt.Println( reflect.TypeOf(utils.TypeChecker_float64(vv["usage"]).(float64)))
				usageArr = append(usageArr, utils.TypeChecker_float64(vv["usage"]).(float64))
			}
		}

		sort.Sort(sort.Reverse(sort.Float64Slice(usageArr)))

		for idx := range usageArr {

			for _, v := range usage {
				if v[usageArr[idx]] != "" {

					var topProcess model.TopProcess

					topProcess.Index = idx + 1
					topProcess.ProcessName = v[usageArr[idx]]
					topProcess.Usage = usageArr[idx]

					topProcessList = append(topProcessList, topProcess)

				}
			}
			if idx > 8 {
				break
			}

		}
		result = topProcessList
		return result, nil
	}
}

func (n ManageNodeService) GetRabbitMqSummary() (result model.RabbitMQGlobalResource, err model.ErrMessage) {

	result, _ = integration.GetRabbitMq(n.openstackProvider).GetRabbitMQOverview()

	return result, nil
}

func (n ManageNodeService) GetManageNodeSummary(apiRequest model.NodeReq) ([]model.ManageNodeResources, model.ErrMessage) {

	manageNodeNameList := getManageNodeName(n.openstackProvider, n.provider, n.influxClient)

	var manageNodeList []model.ManageNodeResources
	var errs []model.ErrMessage

	for _, hostname := range manageNodeNameList {

		if apiRequest.HostName != "" {

			if strings.Contains(hostname, apiRequest.HostName) == false {
				continue
			}
		}

		var manageNodeResource model.ManageNodeResources
		var req model.NodeReq
		req.HostName = hostname

		cpuData, memTotData, memFreeData, diskTotalData, diskUsedData, agentForwarderData, agentCollectorData, err := getManageNodeSummary_Sub(req, n.influxClient)

		if err != nil {
			errs = append(errs, err)
		}

		cpuUsage := utils.GetDataFloatFromInterfaceSingle(cpuData)
		memTot := utils.GetDataFloatFromInterfaceSingle(memTotData)
		memFree := utils.GetDataFloatFromInterfaceSingle(memFreeData)
		memUsage := utils.RoundFloatDigit2(100 - ((memFree / memTot) * 100))
		diskTotal := utils.GetDataFloatFromInterfaceSingle(diskTotalData)
		diskUsed := utils.GetDataFloatFromInterfaceSingle(diskUsedData)
		diskUsage := utils.RoundFloatDigit2((diskUsed / diskTotal) * 100)
		agentForwarderStatus := utils.GetDataFloatFromInterfaceSingle(agentForwarderData)
		agentCollectorStatus := utils.GetDataFloatFromInterfaceSingle(agentCollectorData)

		manageNodeResource.Hostname = hostname
		manageNodeResource.CpuUsage = cpuUsage
		manageNodeResource.MemUsage = memUsage
		manageNodeResource.DiskUsage = diskUsage
		manageNodeResource.MemoryMbMax = memTot
		manageNodeResource.MemoryMbUsed = memTot - memFree
		manageNodeResource.DiskGbMax = diskTotal / 1024
		manageNodeResource.DiskGbUsed = diskUsed / 1024

		if agentForwarderStatus == 1 && agentCollectorStatus == 1 {
			manageNodeResource.AgentStatus = "OK"
		} else {
			manageNodeResource.AgentStatus = "UnKnown"
		}

		manageNodeList = append(manageNodeList, manageNodeResource)

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

	return manageNodeList, nil

}
func getManageNodeName(opts model.OpenstackProvider, provider *gophercloud.ProviderClient, client client.Client) []string {
	var computeNodeList, nodeList []string

	computeNodeInfo, _ := integration.GetNova(opts, provider).GetComputeNodeResources()

	for _, value := range computeNodeInfo {
		computeNodeList = append(computeNodeList, value.Hostname)
	}

	nodeListResp, _ := dao.GetMainDao(client).GetOpenstackNodeList()
	valueList, _ := utils.GetResponseConverter().InfluxConverterToMap(nodeListResp)

	for _, value := range valueList {

		hostname := reflect.ValueOf(value["hostname"]).String()
		if utils.StringArrayDistinct(hostname, nodeList) == false && utils.StringArrayDistinct(hostname, computeNodeList) == false {

			nodeList = append(nodeList, hostname)
		}
	}

	sort.Strings(nodeList)

	return nodeList
}

func getManageNodeSummary_Sub(request model.NodeReq, f client.Client) (map[string]interface{}, map[string]interface{}, map[string]interface{},
	map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, model.ErrMessage) {
	var cpuResp, memTotalResp, memFreeResp, diskTotalResp, diskUsedResp, agentForwarderResp, agentCollectorResp client.Response

	var errs []model.ErrMessage
	var err model.ErrMessage
	var wg sync.WaitGroup
	wg.Add(7)

	for i := 0; i < 7; i++ {
		go func(wg *sync.WaitGroup, index int) {

			switch index {
			case 0:
				cpuResp, err = dao.GetMainDao(f).GetNodeCpuUsage(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 1:
				memTotalResp, err = dao.GetMainDao(f).GetNodeTotalMemoryUsage(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 2:
				memFreeResp, err = dao.GetMainDao(f).GetNodeFreeMemoryUsage(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 3:
				agentForwarderResp, err = dao.GetMainDao(f).GetAgentProcessStatus(request, "forwarder")
				if err != nil {
					errs = append(errs, err)
				}
			case 4:
				agentCollectorResp, err = dao.GetMainDao(f).GetAgentProcessStatus(request, "collector")
				if err != nil {
					errs = append(errs, err)
				}
			case 5:
				diskTotalResp, err = dao.GetMainDao(f).GetNodeTotalDisk(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 6:
				diskUsedResp, err = dao.GetMainDao(f).GetNodeUsedDisk(request)
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
		return nil, nil, nil, nil, nil, nil, nil, errMessage
	}
	//==========================================================================
	cpuUsage, _ := utils.GetResponseConverter().InfluxConverter(cpuResp)
	memTotal, _ := utils.GetResponseConverter().InfluxConverter(memTotalResp)
	memFree, _ := utils.GetResponseConverter().InfluxConverter(memFreeResp)
	diskTotal, _ := utils.GetResponseConverter().InfluxConverter(diskTotalResp)
	diskUsed, _ := utils.GetResponseConverter().InfluxConverter(diskUsedResp)
	agentForwarder, _ := utils.GetResponseConverter().InfluxConverter(agentForwarderResp)
	agentCollector, _ := utils.GetResponseConverter().InfluxConverter(agentCollectorResp)

	//return cpuUsage,memUsage, agentResp, nil
	return cpuUsage, memTotal, memFree, diskTotal, diskUsed, agentForwarder, agentCollector, nil
}

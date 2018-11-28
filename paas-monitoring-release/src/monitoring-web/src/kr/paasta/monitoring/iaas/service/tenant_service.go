package services

import (
	client "github.com/influxdata/influxdb/client/v2"
	"github.com/rackspace/gophercloud"
	"kr/paasta/monitoring/utils"
	"kr/paasta/monitoring/iaas/dao"
	"kr/paasta/monitoring/iaas/integration"
	"kr/paasta/monitoring/iaas/model"
	"sync"
	"strings"
	"encoding/json"
	"strconv"
	"fmt"
)

type TenantService struct {
	openstackProvider model.OpenstackProvider
	provider          *gophercloud.ProviderClient
	influxClient      client.Client
}

func GetTenantService(openstackProvider model.OpenstackProvider, provider *gophercloud.ProviderClient,influxClient client.Client) *TenantService {
	return &TenantService{
		openstackProvider: openstackProvider,
		provider: provider,
		influxClient: 	influxClient,
	}
}



//CPU 사용률
func (n TenantService) GetInstanceCpuUsageList(request model.DetailReq)(result []map[string]interface{}, _ model.ErrMessage){

	cpuUsageResp, err := dao.GetInstanceDao(n.influxClient).GetInstanceCpuUsageList(request)
	if err != nil {
		model.MonitLogger.Error(err)
		return result, err
	}else {
		cpuUsage, _ := utils.GetResponseConverter().InfluxConverterList(cpuUsageResp, model.METRIC_NAME_CPU_USAGE)

		datamap := cpuUsage[model.RESULT_DATA_NAME].([]map[string]interface{})
		for _, data := range datamap{

			swapFree := data["usage"].(json.Number)
			convertData, _ := strconv.ParseFloat(swapFree.String(),64)
			//swap 사용률을 구한다. ( 100 -  freeUsage)
			data["usage"] = utils.RoundFloatDigit2(convertData)
		}

		result = append(result,cpuUsage )
		return result, nil
	}
}


//Instance Memory Usage
func (s TenantService) GetInstanceMemoryUsageList(request model.DetailReq)(result []map[string]interface{}, _ model.ErrMessage){

	memoryResp, err := dao.GetInstanceDao(s.influxClient).GetInstanceMemoryUsageList(request)
	if err != nil {
		model.MonitLogger.Error(err)
		return result, err
	}else {
		memoryUsage, _ := utils.GetResponseConverter().InfluxConverter4Usage(memoryResp, model.METRIC_NAME_MEMORY_USAGE)

		datamap := memoryUsage[model.RESULT_DATA_NAME].([]map[string]interface{})
		for _, data := range datamap{

			swapFree := data["usage"].(float64)
			data["usage"] = utils.RoundFloatDigit2(swapFree)
		}

		result = append(result, memoryUsage )
		return result, nil
	}
}


//Disk IO Read Kbyte
func (s TenantService) GetInstanceDiskIoKbyteList(request model.DetailReq, gubun string)(result []map[string]interface{}, _ model.ErrMessage){

	memoryResp, err := dao.GetInstanceDao(s.influxClient).GetInstanceDiskIoKbyte(request, gubun)
	if err != nil {
		model.MonitLogger.Error(err)
		return result, err
	}else {
		var resultName string

		if gubun == "read"{
			resultName = model.METRIC_NAME_DISK_READ_KBYTE
		}else{
			resultName = model.METRIC_NAME_DISK_WRITE_KBYTE
		}
		byte, _ := utils.GetResponseConverter().InfluxConverterList(memoryResp, resultName)

		datamap := byte[model.RESULT_DATA_NAME].([]map[string]interface{})
		fmt.Println(datamap)
		for _, data := range datamap{
			usage := utils.TypeChecker_float64(data["usage"]).(float64)
			data["usage"] = utils.RoundFloatDigit2(usage)
		}

		result = append(result, byte )
		return result, nil
	}
}

//Network IO Kbyte
func (s TenantService) GetInstanceNetworkIoKbyteList(request model.DetailReq)(result []map[string]interface{}, _ model.ErrMessage){

	networkInResp, err := dao.GetInstanceDao(s.influxClient).GetInstanceNetworkKbyte(request, "in")
	networkOutResp, err := dao.GetInstanceDao(s.influxClient).GetInstanceNetworkKbyte(request, "out")

	if err != nil {
		model.MonitLogger.Error(err)
		return result, err
	}else {


		inData, _ := utils.GetResponseConverter().InfluxConverterList(networkInResp, model.METRIC_NAME_NETWORK_IN)
		outData, _ := utils.GetResponseConverter().InfluxConverterList(networkOutResp, model.METRIC_NAME_NETWORK_OUT)

		inDatamap := inData[model.RESULT_DATA_NAME].([]map[string]interface{})
		for _, data := range inDatamap{
			usage := utils.TypeChecker_float64(data["usage"]).(float64)
			data["usage"] = utils.RoundFloatDigit2(usage)
		}

		outDatamap := outData[model.RESULT_DATA_NAME].([]map[string]interface{})
		for _, data := range outDatamap{
			usage := utils.TypeChecker_float64(data["usage"]).(float64)
			data["usage"] = utils.RoundFloatDigit2(usage)
		}

		result = append(result, inData )
		result = append(result, outData )
		return result, nil
	}
}

//Network Packets
func (s TenantService) GetInstanceNetworkPacketsList(request model.DetailReq)(result []map[string]interface{}, _ model.ErrMessage){

	networkInResp,  err := dao.GetInstanceDao(s.influxClient).GetInstanceNetworkPackets(request, "in")
	networkOutResp, err := dao.GetInstanceDao(s.influxClient).GetInstanceNetworkPackets(request, "out")

	if err != nil {
		model.MonitLogger.Error(err)
		return result, err
	}else {


		inData, _ := utils.GetResponseConverter().InfluxConverterList(networkInResp, model.METRIC_NAME_NETWORK_IN)
		outData, _ := utils.GetResponseConverter().InfluxConverterList(networkOutResp, model.METRIC_NAME_NETWORK_OUT)

		inDatamap := inData[model.RESULT_DATA_NAME].([]map[string]interface{})

		for _, data := range inDatamap{
			usage := utils.TypeChecker_float64(data["usage"]).(float64)
			data["usage"] = utils.RoundFloatDigit2(usage)
		}

		outDatamap := outData[model.RESULT_DATA_NAME].([]map[string]interface{})

		for _, data := range outDatamap{
			usage := utils.TypeChecker_float64(data["usage"]).(float64)
			data["usage"] = utils.RoundFloatDigit2(usage)
		}

		result = append(result, inData )
		result = append(result, outData )
		return result, nil
	}
}

func (n TenantService) GetTenantSummary(apiRequest model.TenantReq)(result []model.TenantSummaryInfo, err error){

	userId, err := integration.GetKeystone(n.openstackProvider, n.provider).GetUserIdByName(n.openstackProvider.Username)
	if err != nil {
		fmt.Println("Get UserId Error :", err)
		return result, err
	}

	//Get Tenant List by User Own
	tenantLists, err := integration.GetKeystone(n.openstackProvider,  n.provider).GetUserTenantList(userId)

	var searchTenantList []model.TenantInfo

	for _, tenant := range tenantLists{
		if apiRequest.TenantName != "" && strings.Contains(tenant.Name, apiRequest.TenantName){
			searchTenantList = append(searchTenantList, tenant)
		}else if apiRequest.TenantName == ""{
			searchTenantList = append(searchTenantList, tenant)
		}


	}

	if err != nil {
		fmt.Println("Get nodes resources error :", err)
	}

	var tenantSummaryInfos []model.TenantSummaryInfo

	for _, tenant := range searchTenantList{

		var tenantSummaryInfo model.TenantSummaryInfo

		tenantInstances,  tenantResourceLimit, tenantNetworkLimit, tenantFloatingIps, tenantSecurityGroups, tenantStorageResources,_ := getTenantSummary_Sub(n.openstackProvider, n.provider,  tenant.Id, tenant.Name)

		var total_vcpus, total_memory float64
		for _, instance :=range tenantInstances{

			total_vcpus = total_vcpus + instance.Vcpus
			total_memory = total_memory + instance.MemoryMb
			//total_disk = total_disk + instance.Disk_gb
		}

		tenantSummaryInfo.Name = tenant.Name
		tenantSummaryInfo.Id = tenant.Id
		tenantSummaryInfo.Enabled = tenant.Enabled
		tenantSummaryInfo.InstancesUsed = len(tenantInstances)
		tenantSummaryInfo.MemoryMbUsed = total_memory
		tenantSummaryInfo.VcpusUsed = total_vcpus

		tenantSummaryInfo.MemoryMbLimit = tenantResourceLimit.MemoryMbLimit
		tenantSummaryInfo.InstancesLimit = tenantResourceLimit.InstancesLimit
		tenantSummaryInfo.VcpusLimit  = tenantResourceLimit.CoresLimit

		tenantSummaryInfo.SecurityGroupsLimit = tenantNetworkLimit.SecurityGroupLimit
		tenantSummaryInfo.FloatingIpsLimit    = tenantNetworkLimit.FloatingIpsLimit

		tenantSummaryInfo.FloatingIpsUsed     = len(tenantFloatingIps)
		tenantSummaryInfo.SecurityGroupsUsed  = tenantSecurityGroups

		tenantSummaryInfo.VolumeStorageLimit   = tenantStorageResources.VolumesLimit
		tenantSummaryInfo.VolumeStorageUsed    = tenantStorageResources.Volumes
		tenantSummaryInfo.VolumeStorageLimitGb = tenantStorageResources.VolumeLimitGb
		tenantSummaryInfo.VolumeStorageUsedGb    = tenantStorageResources.VolumeGb

		tenantSummaryInfos = append(tenantSummaryInfos, tenantSummaryInfo)
	}

	return tenantSummaryInfos, nil

}

func  getTenantSummary_Sub(opts model.OpenstackProvider, provider *gophercloud.ProviderClient, tenantId, tenantName string)(tenantInstances []model.InstanceInfo,
			tenantResourcesLimit model.TenantResourcesLimit, tenantNetworkLimit model.TenantNetworkLimit,
			tenantFloatingIps []model.FloatingIPInfo, tenantSecurityGroups int, tenantStorageResources model.TenantStorageResources,
			_ model.ErrMessage) {

	var errs []model.ErrMessage
	var err error
	var wg sync.WaitGroup
	wg.Add(6)
	for i := 0; i < 6; i++ {
		go func(wg *sync.WaitGroup, index int) {
			switch index {
			case 0 :
				tenantInstances, err  = integration.GetNova(opts, provider).GetProjectInstances(tenantId)
				if err != nil {
					//errs = append(errs, err)
				}
			case 1 :
				tenantResourcesLimit, err = integration.GetNova(opts, provider).GetProjectResourcesLimit(tenantId)
				if err != nil {
					//errs = append(errs, err)
				}
			case 2 :
				tenantNetworkLimit, err = integration.GetNeutron(opts, provider).GetTenantNetworkLimit(tenantId)
				if err != nil {
					//errs = append(errs, err)
				}
			case 3 :
				tenantFloatingIps, err = integration.GetNeutron(opts, provider).GetTenantFloatingIps(tenantId)
				if err != nil {
					//errs = append(errs, err)
				}
			case 4 :
				tenantSecurityGroups, err = integration.GetNeutron(opts, provider).GetTenantSecurityGroups(tenantId)
				if err != nil {
					//errs = append(errs, err)
				}
			case 5 :
				tenantStorageResources, err = integration.GetCinder(opts, provider).GetTenantStorageResources(tenantId, tenantName)
				if err != nil {
					//errs = append(errs, err)
				}
			}

			wg.Done()
		}(&wg, i)
	}
	wg.Wait()

	//==========================================================================
	// Error가 여러건일 경우 대해 고려해야함.
	if len(errs) > 0 {
		var returnErrMessage string
		for _, err := range errs{
			returnErrMessage = returnErrMessage + " " + err["Message"].(string)
		}
		errMessage := model.ErrMessage{
			"Message": returnErrMessage ,
		}
		return tenantInstances, tenantResourcesLimit, tenantNetworkLimit, tenantFloatingIps, tenantSecurityGroups, tenantStorageResources, errMessage
	}
	//==========================================================================
	model.MonitLogger.Debug("tenantStorageResources::", tenantStorageResources)
	return tenantInstances, tenantResourcesLimit, tenantNetworkLimit, tenantFloatingIps, tenantSecurityGroups, tenantStorageResources, nil
}

func (n TenantService) GetTenantInstanceList(apiRequest model.TenantReq)(resultArr map[string]interface{}, err error){

	var result []model.InstanceInfo
	instanceMainList , _ := integration.GetNova(n.openstackProvider, n.provider).GetProjectInstancesList(apiRequest)
	instanceSubInfo , _ :=  integration.GetNova(n.openstackProvider, n.provider).GetProjectInstances(apiRequest.TenantId)

	totInstance := 0
	for _, instance := range instanceSubInfo{
		if apiRequest.HostName != "" && strings.Contains(instance.Name, apiRequest.HostName){
			totInstance += 1
		}else if apiRequest.HostName == ""{
			totInstance += 1
		}
	}

	var errs []model.ErrMessage

	for _, mainInstance := range instanceMainList{
		for _, subInstance := range instanceSubInfo{
			if mainInstance.InstanceId == subInstance.InstanceId{
				mainInstance.Flavor    = subInstance.Flavor
				mainInstance.TenantId  = subInstance.TenantId
				mainInstance.Vcpus     = subInstance.Vcpus
				mainInstance.MemoryMb  = subInstance.MemoryMb
				mainInstance.DiskGb    = subInstance.DiskGb
				mainInstance.StartedAt = subInstance.StartedAt
				mainInstance.Uptime    = subInstance.Uptime
			}
		}
		var req model.InstanceReq
		req.InstanceId = mainInstance.InstanceId
		cpuData, memTotData, memFreeData, err := getTenantInstanceStatus_Sub(req, n.influxClient)

		if err != nil {
			errs = append(errs, err)
		}

		cpuUsage  := utils.GetDataFloatFromInterfaceSingle(cpuData)
		memTot    := utils.GetDataFloatFromInterfaceSingle(memTotData)
		memFree   := utils.GetDataFloatFromInterfaceSingle(memFreeData)

		mainInstance.CpuUsage = cpuUsage
		if memTot > 0.0 && memFree > 0.0 {
			mainInstance.MemoryUsage = utils.RoundFloatDigit2(100 - ((memFree/memTot)*100))
		}

		result = append(result, mainInstance)
	}

	resultArr = map[string]interface{}{
		model.RESULT_CNT:        totInstance,
		model.RESULT_PROJECT_ID: apiRequest.TenantId,
		model.RESULT_DATA_NAME:  result,
	}

	return resultArr, nil

}

func  getTenantInstanceStatus_Sub( request model.InstanceReq, f client.Client )(map[string]interface{}, map[string]interface{}, map[string]interface{}, model.ErrMessage) {

	var errs []model.ErrMessage
	var err model.ErrMessage
	var wg sync.WaitGroup
	var cpuResp, memTotalResp, memFreeResp client.Response

	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func(wg *sync.WaitGroup, index int) {
			switch index {
			case 0 :
				cpuResp, err  = dao.GetInstanceDao(f).GetInstanceCpuUsage(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 1 :
				memTotalResp, err = dao.GetInstanceDao(f).GetInstanceTotalMemoryUsage(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 2 :
				memFreeResp, err = dao.GetInstanceDao(f).GetInstanceFreeMemoryUsage(request)
				if err != nil {
					errs = append(errs, err)
				}
			}

			wg.Done()
		}(&wg, i)
	}
	wg.Wait()

	//==========================================================================
	// Error가 여러건일 경우 대해 고려해야함.
	if len(errs) > 0 {
		var returnErrMessage string
		for _, err := range errs{
			returnErrMessage = returnErrMessage + " " + err["Message"].(string)
		}
		errMessage := model.ErrMessage{
			"Message": returnErrMessage ,
		}
		return nil, nil, nil  , errMessage
	}
	//==========================================================================

	cpuUsage, _   := utils.GetResponseConverter().InfluxConverter(cpuResp)
	memTotal, _   := utils.GetResponseConverter().InfluxConverter(memTotalResp)
	memFree,  _   := utils.GetResponseConverter().InfluxConverter(memFreeResp)


	return cpuUsage, memTotal , memFree, nil
}
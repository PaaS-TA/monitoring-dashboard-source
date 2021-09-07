package services

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/gophercloud/gophercloud"
	"kr/paasta/monitoring/iaas_new/dao"
	"kr/paasta/monitoring/iaas_new/integration"
	"kr/paasta/monitoring/iaas_new/model"
	"kr/paasta/monitoring/utils"
	"reflect"
)

type MainService struct {
	openstackProvider model.OpenstackProvider
	provider          *gophercloud.ProviderClient
	influxClient      client.Client
}

func GetMainService(openstackProvider model.OpenstackProvider, provider *gophercloud.ProviderClient, influxClient client.Client) *MainService {
	return &MainService{
		openstackProvider: openstackProvider,
		provider:          provider,
		influxClient:      influxClient,
	}
}

func (n MainService) GetOpenstackSummary(userName string) (model.HypervisorResources, error) {

	//Openstack Summary 정보 조회
	openstackSummaryInfo, err := integration.GetNova(n.openstackProvider, n.provider).GetOpenstackResources()

	if err != nil {
		model.MonitLogger.Error("Error Occur", err)
		return openstackSummaryInfo, err
	}
	model.MonitLogger.Debug("openstackSummaryInfo is ", openstackSummaryInfo)

	//Running Vms 정보 조회
	NodeList, err := integration.GetNova(n.openstackProvider, n.provider).GetComputeNodeResources()

	model.MonitLogger.Debug("node List is:", NodeList)

	if err != nil {
		return openstackSummaryInfo, err
	}

	//Available Total Instance Get
	userId, err := integration.GetKeystone(n.openstackProvider, n.provider).GetUserIdByName(userName)
	if err != nil {
		fmt.Println("Get UserId Error :", err)
		return openstackSummaryInfo, err
	}
	//Get Tenant List by User Own

	projectLists, err := integration.GetKeystone(n.openstackProvider, n.provider).GetUserTenantList(userId)

	var availableInstance int
	for _, project := range projectLists {
		limits, _ := integration.GetNova(n.openstackProvider, n.provider).GetProjectResourcesLimit(project.Id)
		availableInstance += limits.InstancesLimit
	}

	var instanceGuidList []string

	var noStatusList []string
	var runningList []string
	var idleList []string
	var pausedList []string
	var shutDownList []string
	var shutOffList []string
	var crashedList []string
	var powerOffList []string

	for _, computeNode := range NodeList {

		var request model.NodeReq
		var instanceListResp client.Response
		request.HostName = computeNode.Hostname

		instanceListResp, _ = dao.GetNodeDao(n.influxClient).GetAliveInstanceListByNodename(request, true)
		instanceList, _ := utils.GetResponseConverter().InfluxConverterToMap(instanceListResp)

		for _, value := range instanceList {
			instanceGuid := reflect.ValueOf(value["resource_id"]).String()
			if utils.StringArrayDistinct(instanceGuid, instanceGuidList) == false {

				instanceGuidList = append(instanceGuidList, instanceGuid)
				/*
					-1 : no status,
					0 : Running / OK,
					1 : Idle / blocked,
					2 : Paused,
					3 : Shutting down,
					4 : Shut off or Nova suspend
					5 : Crashed,
					6 : Power management suspend (S3 state)
				*/
				if utils.TypeChecker_int(value["value"]).(int64) == -1 {
					noStatusList = append(noStatusList, instanceGuid)
				} else if utils.TypeChecker_int(value["value"]).(int64) == 0 {
					runningList = append(runningList, instanceGuid)
				} else if utils.TypeChecker_int(value["value"]).(int64) == 1 {
					idleList = append(idleList, instanceGuid)
				} else if utils.TypeChecker_int(value["value"]).(int64) == 2 {
					pausedList = append(pausedList, instanceGuid)
				} else if utils.TypeChecker_int(value["value"]).(int64) == 3 {
					shutDownList = append(shutDownList, instanceGuid)
				} else if utils.TypeChecker_int(value["value"]).(int64) == 4 {
					shutOffList = append(shutOffList, instanceGuid)
				} else if utils.TypeChecker_int(value["value"]).(int64) == 5 {
					crashedList = append(crashedList, instanceGuid)
				} else if utils.TypeChecker_int(value["value"]).(int64) == 6 {
					powerOffList = append(powerOffList, instanceGuid)
				}
			}
		}
	}

	//VM의 상태 통계 정보 리턴
	vmStatusList := utils.GetVmStatusCount(noStatusList, runningList, idleList, pausedList, shutDownList, shutOffList, crashedList, powerOffList)

	openstackSummaryInfo.VmTotalLimit = availableInstance
	openstackSummaryInfo.VmRunning = len(runningList)
	openstackSummaryInfo.VmState = vmStatusList
	return openstackSummaryInfo, err
}

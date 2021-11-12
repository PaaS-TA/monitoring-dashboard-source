package service

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/hypervisors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/usage"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	_ "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	client "github.com/influxdata/influxdb1-client/v2"
	"kr/paasta/monitoring/iaas_new/model"
	"kr/paasta/monitoring/utils"
)

type OpenstackService struct {
	OpenstackProvider model.OpenstackProvider
	Provider          *gophercloud.ProviderClient
	InfluxClient      client.Client
}

func GetOpenstackService(openstackProvider model.OpenstackProvider, provider *gophercloud.ProviderClient, influxClient client.Client) *OpenstackService {
	return &OpenstackService{
		OpenstackProvider: openstackProvider,
		Provider:          provider,
		InfluxClient:      influxClient,
	}
}

func (service *OpenstackService) GetHypervisorStatistics(userName string) (map[string]interface{}, error) {

	client, err := utils.GetComputeClient(service.Provider, service.OpenstackProvider.Region)
	if err != nil {
		fmt.Println(err)
	}

	hypervisorStatistics, err := hypervisors.GetStatistics(client).Extract()
	if err != nil {
		fmt.Println(err)
	}

	result := make(map[string]interface{})

	result["runningVms"] = hypervisorStatistics.RunningVMs
	result["vcpu"] = hypervisorStatistics.VCPUs
	result["vcpuUsed"] = hypervisorStatistics.VCPUsUsed
	result["freeRam"] = hypervisorStatistics.FreeRamMB
	result["freeDisk"] = hypervisorStatistics.FreeDiskGB
	result["memory"] = hypervisorStatistics.MemoryMB
	result["memoryUsed"] = hypervisorStatistics.MemoryMBUsed
	result["disk"] = hypervisorStatistics.LocalGB
	result["diskUsed"] = hypervisorStatistics.LocalGBUsed

	result["hypervisorCount"] = hypervisorStatistics.Count

	//result, err := compute.GetHypervisorStatistics(service.osSession)

	return result, err
}


/**
	@Unused
	@parameter
	- params [map]
		- host [string] : 호스트 이름
		- name [string] : 서버 이름
		- status [string] : 상태
		- tenantId [string] : 프로젝트 ID
*/
func (service *OpenstackService) GetServerList(params map[string]interface{}) ([]interface{}, error) {
	client, err := utils.GetComputeClient(service.Provider, service.OpenstackProvider.Region)
	if err != nil {
		fmt.Println(err)
	}

	var listOpts servers.ListOpts
	if params != nil {
		host, ok := params["host"].(string)
		if ok {
			listOpts.Host = host
		}
		name, ok := params["name"].(string)
		if ok {
			listOpts.Name = name
		}
		status, ok := params["status"].(string)
		if ok {
			listOpts.Status = status
		}
		tenantId, ok := params["tenantId"].(string)
		if ok {
			listOpts.TenantID = tenantId
		}
		allTenants, ok := params["allTenants"].(bool)
		if ok {
			listOpts.AllTenants = allTenants
		}
	}

	result := servers.List(client, listOpts)
	resultPages, err := result.AllPages()
	if err != nil {
		fmt.Println(err)
	}
	resultBody := resultPages.GetBody()

	list := resultBody.(map[string][]interface{})["servers"]

	return list, err

}


func (service *OpenstackService) GetProjectList(params map[string]interface{}) ([]interface{}, error) {
	client := utils.GetKeystoneClient(service.Provider)
	networkClient := utils.GetNetworkClient(service.Provider, service.OpenstackProvider.Region)
	blockstorageClient := utils.GetBlockStorageClient(service.Provider, service.OpenstackProvider.Region)

	var listOpts projects.ListOpts
	result := projects.List(client, listOpts)
	resultPages, err := result.AllPages()

	if err != nil {
		fmt.Println(err)
	}
	resultBody := resultPages.GetBody()
	list := resultBody.(map[string][]interface{})["projects"]

	for _, item := range(list) {
		itemMap := item.(map[string]interface{})

		projectId := itemMap["id"].(string)

		// 프로젝트 내 Floating IP 개수 조회
		var fipListOpts floatingips.ListOpts
		fipListOpts.ProjectID = projectId

		allPages, err := floatingips.List(networkClient, fipListOpts).AllPages()
		if err != nil {
			panic(err)
		}
		allFloatingIPs, err := floatingips.ExtractFloatingIPs(allPages)
		if err != nil {
			panic(err)
		}
		itemMap["floatingIps"] = len(allFloatingIPs)

		// 프로젝트 내 보안그룹 개수 조회
		var secGroupListOpts groups.ListOpts
		secGroupListOpts.ProjectID = projectId
		secGroupPages, err := groups.List(networkClient, secGroupListOpts).AllPages()
		if err != nil {
			panic(err)
		}

		secGroups, err := groups.ExtractGroups(secGroupPages)
		if err != nil {
			panic(err)
		}
		itemMap["secGroups"] = len(secGroups)

		// 프로젝트 내 볼륨 개수 조회
		volumeListOpt := volumes.ListOpts{
			TenantID: projectId,
		}
		volumePages, _ := volumes.List(blockstorageClient, volumeListOpt).AllPages()
		volumeBody := volumePages.GetBody()
		volumeList := volumeBody.(map[string][]interface{})["volumes"]
		itemMap["volumes"] = len(volumeList)
		//utils.Logger.Debugf("len(result.Attachments) : %v\n", len(result.Attachments))

		//service.retrieveSingleProjectUsage(projectId)   // TODO 호출해도 조회안됨...


		/*
		var listOpts servers.ListOpts
		listOpts.TenantID = projectId

		result := servers.List(client, listOpts)
		resultPages, err := result.AllPages()
		resultBody := resultPages.GetBody()
		serverList := resultBody.(map[string][]interface{})["servers"]

		itemMap["instances"] = len(serverList)
		 */
	}
	return list, err
}


func (service *OpenstackService) RetrieveTenantUsage() []usage.TenantUsage {
	computeClient, _ := utils.GetComputeClient(service.Provider, service.OpenstackProvider.Region)

	allTenantsOpts := usage.AllTenantsOpts{
		Detailed: true,
	}

	usagePages, _ := usage.AllTenants(computeClient, allTenantsOpts).AllPages()
	usageList, _ := usage.ExtractAllTenants(usagePages)

	// For Test
	/*
	for _, item := range(usageList) {
		for _, server := range(item.ServerUsages) {
			id := server.InstanceID


			result, _ := servers.Get(computeClient, id).Extract()
			addressList := result.Addresses

			for _, address := range addressList {
				tee := address.([]interface{})
				tee2 := tee[0].(map[string]interface{})
				ipAddress := tee2["addr"]
				server.Name += " ("
				server.Name += ipAddress.(string)
				server.Name += ")"
			}

		}
	}
	*/

	//service.GetServerDetail("08769233-2599-4ba5-ae54-e8aa92bd11b9")

	return usageList
}


// Unused
func (service *OpenstackService) retrieveSingleProjectUsage(projectId string) {
	computeClient, _ := utils.GetComputeClient(service.Provider, service.OpenstackProvider.Region)
	//computeClient.Microversion = "2.40"
	opts := usage.SingleTenantOpts{

	}
	fmt.Println("projectId : " + projectId)
	usagePages, _ := usage.SingleTenant(computeClient, projectId, opts).AllPages()
	usageList, _ := usage.ExtractSingleTenant(usagePages)
	fmt.Println(usageList)

}


// Unused
func (service *OpenstackService) GetServerDetail(id string) *servers.Server {
	computeClient, _ := utils.GetComputeClient(service.Provider, service.OpenstackProvider.Region)

	result, _ := servers.Get(computeClient, id).Extract()
	return result
}


func (service *OpenstackService) GetHypervisorList() ([]interface{}, error) {
	computeClient, _ := utils.GetComputeClient(service.Provider, service.OpenstackProvider.Region)

	withServices := false
	opts := hypervisors.ListOpts{
		WithServers: &withServices,
	}

	allPages, err := hypervisors.List(computeClient, opts).AllPages()
	if err != nil {
		utils.Logger.Error(err)
	}

	resultBody := allPages.GetBody()

	hypervisorMap := resultBody.(map[string]interface{})["hypervisors"]

	hypervisorList := hypervisorMap.([]interface{})

	return hypervisorList, err
}

/**
	@Unused
 */
func (service *OpenstackService) getHypervisorDetail(hypervisorId string) (*hypervisors.Hypervisor, error) {
	computeClient, _ := utils.GetComputeClient(service.Provider, service.OpenstackProvider.Region)

	result, err := hypervisors.Get(computeClient, hypervisorId).Extract()

	return result, err
}
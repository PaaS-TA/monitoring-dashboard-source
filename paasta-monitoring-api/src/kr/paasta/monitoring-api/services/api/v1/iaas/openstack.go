package iaas

import (
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"log"
	utils "paasta-monitoring-api/helpers"
	"strconv"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/hypervisors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/usage"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
)


type OpenstackService struct {
	Provider  *gophercloud.ProviderClient
}

func GetOpenstackService(provider *gophercloud.ProviderClient) *OpenstackService {
	return &OpenstackService{
		Provider: provider,
	}
}

func (service *OpenstackService) GetHypervisorStatistics() (map[string]interface{}, error) {

	client, err := utils.GetComputeClient(service.Provider, "")
	if err != nil {
		log.Println(err.Error())
	}

	hypervisorStatistics, err := hypervisors.GetStatistics(client).Extract()
	if err != nil {
		log.Println(err.Error())
	}

	log.Println("Running VMs : " + strconv.Itoa(hypervisorStatistics.RunningVMs))

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


func (service *OpenstackService) GetHypervisorList() ([]interface{}, error) {
	computeClient, _ := utils.GetComputeClient(service.Provider, "")

	withServices := false
	opts := hypervisors.ListOpts{
		WithServers: &withServices,
	}

	allPages, err := hypervisors.List(computeClient, opts).AllPages()
	if err != nil {
		log.Println(err.Error())
	}

	resultBody := allPages.GetBody()

	hypervisorMap := resultBody.(map[string]interface{})["hypervisors"]

	hypervisorList := hypervisorMap.([]interface{})

	return hypervisorList, err
}

func (service *OpenstackService) GetProjectList(params map[string]interface{}) ([]interface{}, error) {
	client := utils.GetKeystoneClient(service.Provider)
	networkClient := utils.GetNetworkClient(service.Provider, "")
	blockstorageClient := utils.GetBlockStorageClient(service.Provider, "")

	var listOpts projects.ListOpts
	result := projects.List(client, listOpts)
	resultPages, err := result.AllPages()

	if err != nil {
		log.Println(err.Error())
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

func (service *OpenstackService) RetrieveTenantUsage() ([]usage.TenantUsage, error) {
	computeClient, _ := utils.GetComputeClient(service.Provider, "")

	allTenantsOpts := usage.AllTenantsOpts{
		Detailed: true,
	}

	usagePages, err := usage.AllTenants(computeClient, allTenantsOpts).AllPages()
	usageList, _ := usage.ExtractAllTenants(usagePages)

	// IP 정보는 별도로 조회해야 하지만.. 속도가 너무 느려짐..
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

					fmt.Println(server.Name)
				}

			}
		}
	*/

	//service.GetServerDetail("08769233-2599-4ba5-ae54-e8aa92bd11b9")

	return usageList, err
}


func (service *OpenstackService) GetHostIpAddress(instanceId string) string {
	var ipAddress string

	computeClient, _ := utils.GetComputeClient(service.Provider, "")

	result, _ := servers.Get(computeClient, instanceId).Extract()
	addressList := result.Addresses

	for _, address := range addressList {
		dataArr := address.([]interface{})
		dataMap := dataArr[0].(map[string]interface{})
		ipAddress = dataMap["addr"].(string)
	}
	return ipAddress
}
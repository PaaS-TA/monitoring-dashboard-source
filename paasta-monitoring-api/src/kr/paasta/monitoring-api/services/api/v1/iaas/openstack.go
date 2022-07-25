package iaas

import (
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
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

func (service *OpenstackService) GetHypervisorStatistics(ctx echo.Context) (map[string]interface{}, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	client, err := utils.GetComputeClient(service.Provider, "")
	if err != nil {
		logger.Error(err)
	}

	hypervisorStatistics, err := hypervisors.GetStatistics(client).Extract()
	if err != nil {
		logger.Error(err)
	}

	logger.Debug("Running VMs : " + strconv.Itoa(hypervisorStatistics.RunningVMs))

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


func (service *OpenstackService) GetHypervisorList(ctx echo.Context) ([]interface{}, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	computeClient, _ := utils.GetComputeClient(service.Provider, "")
	withServices := false
	opts := hypervisors.ListOpts{
		WithServers: &withServices,
	}

	allPages, err := hypervisors.List(computeClient, opts).AllPages()
	if err != nil {
		logger.Error(err)
	}

	resultBody := allPages.GetBody()
	hypervisorMap := resultBody.(map[string]interface{})["hypervisors"]
	hypervisorList := hypervisorMap.([]interface{})

	return hypervisorList, err
}


func (service *OpenstackService) GetProjectList(ctx echo.Context) ([]interface{}, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	client := utils.GetKeystoneClient(service.Provider)
	networkClient := utils.GetNetworkClient(service.Provider, "")
	blockstorageClient := utils.GetBlockStorageClient(service.Provider, "")

	var listOpts projects.ListOpts
	result := projects.List(client, listOpts)
	resultPages, err := result.AllPages()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	resultBody := resultPages.GetBody()
	list := resultBody.(map[string][]interface{})["projects"]

	for _, item := range(list) {
		itemMap := item.(map[string]interface{})

		projectId := itemMap["id"].(string)

		// 프로젝트 내 Floating IP 개수 조회
		var fipListOpts floatingips.ListOpts
		fipListOpts.ProjectID = projectId
		ipsPages, ipsErr := floatingips.List(networkClient, fipListOpts).AllPages()
		if ipsErr != nil {
			logger.Error(ipsErr)
			return nil, ipsErr
		}

		floatingIPs, fipsErr := floatingips.ExtractFloatingIPs(ipsPages)
		if fipsErr != nil {
			logger.Error(fipsErr)
			return nil, fipsErr
		}
		itemMap["floatingIps"] = len(floatingIPs)

		// 프로젝트 내 보안그룹 개수 조회
		var secGroupListOpts groups.ListOpts
		secGroupListOpts.ProjectID = projectId
		secGroupPages, gErr := groups.List(networkClient, secGroupListOpts).AllPages()
		if gErr != nil {
			logger.Error(gErr)
			return nil, gErr
		}

		secGroups, sgErr := groups.ExtractGroups(secGroupPages)
		if sgErr != nil {
			logger.Error(sgErr)
			return nil, sgErr
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

func (service *OpenstackService) RetrieveTenantUsage(ctx echo.Context) ([]usage.TenantUsage, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	serverParams := make(map[string]interface{}, 0)
	serverParams["allTenants"] = true
	tenantIdParam := ctx.Param("tenantId")
	if tenantIdParam != "" {
		serverParams["tenantId"] = tenantIdParam
	}

	computeClient, _ := utils.GetComputeClient(service.Provider, "")
	allTenantsOpts := usage.AllTenantsOpts{
		Detailed: true,
	}
	usagePages, err := usage.AllTenants(computeClient, allTenantsOpts).AllPages()
	if err != nil {
		logger.Error(err)
		return nil, err
	}
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
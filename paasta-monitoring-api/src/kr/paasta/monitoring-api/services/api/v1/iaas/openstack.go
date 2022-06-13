package iaas

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/hypervisors"
	utils "paasta-monitoring-api/helpers"
)


type OpenstackService struct {
	Provider          *gophercloud.ProviderClient
}

func GetOpenstackService(provider *gophercloud.ProviderClient) *OpenstackService {
	return &OpenstackService{
		Provider: provider,
	}
}

func (service *OpenstackService) GetHypervisorStatistics() (map[string]interface{}, error) {
	fmt.Println(service.Provider.TokenID)
	client, err := utils.GetComputeClient(service.Provider, "")
	if err != nil {
		fmt.Println(err)
	}

	hypervisorStatistics, err := hypervisors.GetStatistics(client).Extract()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(hypervisorStatistics.RunningVMs)

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
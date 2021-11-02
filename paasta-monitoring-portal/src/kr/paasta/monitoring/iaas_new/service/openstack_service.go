package service

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/hypervisors"
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

	hypervisorStatistics, err := hypervisors.GetStatistics(client).Extract();
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

	//result, err := compute.GetHypervisorStatistics(service.osSession)

	return result, err
}
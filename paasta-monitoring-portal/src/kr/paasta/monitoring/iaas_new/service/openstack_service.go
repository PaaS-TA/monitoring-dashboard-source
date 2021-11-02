package service

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/hypervisors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
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

	result["hypervisorCount"] = hypervisorStatistics.Count

	//result, err := compute.GetHypervisorStatistics(service.osSession)

	return result, err
}


/**
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

	var listOpts projects.ListOpts
	result := projects.List(client, listOpts)
	resultPages, err := result.AllPages()
	if err != nil {
		fmt.Println(err)
	}
	resultBody := resultPages.GetBody()
	list := resultBody.(map[string][]interface{})["projects"]

	return list, err
}
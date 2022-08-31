package compute

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/hypervisors"

	"monitoring-portal/openstack-client/session"
	"monitoring-portal/openstack-client/utils"
)

/********************
	Compute (Nova)
*********************/


func GetHypervisor(session session.OpenStackSession, hypervisorId string) {
	opts := gophercloud.EndpointOpts{ Region: "RegionOne" }
	client, err := openstack.NewComputeV2(session.Provider, opts)
	if err != nil {
		fmt.Println(err)
	}

	result := hypervisors.Get(client, hypervisorId)
	utils.PrintJson(result.Body)
}

/**
@parameter
	- params [map]
		- hypervisorHostname [string] : 호스트 이름
		- withServers [bool] : 서버 정보 포함 여부 (?)
*/
func GetHypervisorList(session session.OpenStackSession, params map[string]interface {}) {
	opts := gophercloud.EndpointOpts{ Region: "RegionOne" }
	client, err := openstack.NewComputeV2(session.Provider, opts)
	if err != nil {
		fmt.Println(err)
	}

	var listOpts hypervisors.ListOpts
	if params != nil {
		hypervisorHostname, ok := params["hypervisorHostname"].(string)
		if ok {
			listOpts.HypervisorHostnamePattern = &hypervisorHostname
			client.Microversion = "2.53"
		}

		withServers, ok := params["withServers"].(bool)
		if ok {
			listOpts.WithServers = &withServers
			client.Microversion = "2.53"
		}
	}

	resultPager := hypervisors.List(client, listOpts)
	list := utils.PagerToMap(resultPager)
	hypervisorList := list.(map[string]interface{})["hypervisors"]
	utils.PrintJson(hypervisorList)

	//fmt.Println(list)
}


func GetHypervisorStatistics(session session.OpenStackSession) (map[string]interface{}, error) {
	fmt.Println("IdentityEndpoint : " + session.Provider.IdentityEndpoint)

	opts := gophercloud.EndpointOpts{ Region: "RegionOne" }
	client, err := openstack.NewComputeV2(session.Provider, opts)
	if err != nil {
		fmt.Println(err)
	}

	hypervisorStatistics, err := hypervisors.GetStatistics(client).Extract();
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("========= Hypervisor Statistics =========")
	fmt.Printf("RunningVMs : %d\n", hypervisorStatistics.RunningVMs)
	fmt.Printf("VCPUs : %d\n", hypervisorStatistics.VCPUs)
	fmt.Printf("VCPUsUsed : %d\n", hypervisorStatistics.VCPUsUsed)
	fmt.Printf("FreeRam(MB)    : %d MiB \n", hypervisorStatistics.FreeRamMB)
	fmt.Printf("FreeDisk(GB)   : %d GiB \n", hypervisorStatistics.FreeDiskGB)
	fmt.Printf("Memory(MB)     : %d MiB \n", hypervisorStatistics.MemoryMB)
	fmt.Printf("MemoryUsed(MB) : %d MiB \n", hypervisorStatistics.MemoryMBUsed)

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

	//result["vcpuPercent"] = (hypervisorStatistics.VCPUsUsed / hypervisorStatistics.VCPUs) * 100
	//result["memoryPercent"] = (hypervisorStatistics.MemoryMBUsed / hypervisorStatistics.MemoryMB) * 100
	//result["diskPercent"] = (hypervisorStatistics.LocalGBUsed / hypervisorStatistics.LocalGB) * 100

	return result, err
}

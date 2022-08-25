package compute

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/pagination"

	"monitoring-portal/openstack-client/session"
	"monitoring-portal/openstack-client/utils"
)

/********************
	Compute (Nova)
*********************/

/**
	@parameter
		- params [map]
			- host [string] : 호스트 이름
			- name [string] : 서버 이름
			- status [string] : 상태
			- tenantId [string] : 프로젝트 ID
 */
func GetServerList(session session.OpenStackSession, params map[string]interface{}) []interface{} {
	// 서버 목록 조회
	opts := gophercloud.EndpointOpts{
		Region: "RegionOne",
	}
	computeClient, computeErr := openstack.NewComputeV2(session.Provider, opts)
	if computeErr != nil {
		fmt.Println(computeErr)
	}

	fmt.Println("Server List")
	servers.List(computeClient, nil).EachPage(func (page pagination.Page) (bool, error) {
		serverList, err := servers.ExtractServers(page)
		if err != nil {
			return false, err
		}

		utils.PrintJson(serverList)
		return true, nil;
	})

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

	result := servers.List(computeClient, listOpts)
	resultPages, err := result.AllPages()
	if err != nil {
		fmt.Println(err)
	}
	resultBody := resultPages.GetBody()

	list := resultBody.(map[string][]interface{})["servers"]

	return list
}

/**
	@parameter
		- serverId [string] : Server ID
 */
func GetServerDetail(session session.OpenStackSession, serverId string) {
	opts := gophercloud.EndpointOpts{
		Region: "RegionOne",
	}

	serverClient, err := openstack.NewComputeV2(session.Provider, opts)
	if err != nil {
		fmt.Println(err)
	}

	result := servers.Get(serverClient, serverId)

	utils.PrintJson(result.Body)
}

/**
	params
		- flavorId [string] : Flavor ID
 */
func GetFlavor(session session.OpenStackSession, flavorId string) {
	// 네트워크 목록 조회
	opts := gophercloud.EndpointOpts{Region: "RegionOne"}
	flavorClient, err := openstack.NewComputeV2(session.Provider, opts)
	if err != nil {
		fmt.Println(err)
	}

	result := flavors.Get(flavorClient, flavorId)


	utils.PrintJson(result.Body)
}
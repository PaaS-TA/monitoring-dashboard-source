package network

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"

	"monitoring-portal/openstack-client/session"
)


/***********************
	Network (Neutron)
************************/

/**
	params
		- projectId
 */
func GetNetworkList(session session.OpenStackSession, params map[string]string) []interface{} {
	// 네트워크 목록 조회
	opts := gophercloud.EndpointOpts{Region: "RegionOne"}
	networkClient, err := openstack.NewNetworkV2(session.Provider, opts)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println("Network List")
	//networks.List(networkClient, nil).EachPage(func (page pagination.Page) (bool, error) {
	//	networkList, err := networks.ExtractNetworks(page)
	//	if err != nil {
	//		return false, err
	//	}
	//
	//	utils.PrintJson(networkList)
	//	return true, nil;
	//})

	var listOpts networks.ListOpts
	if params != nil {
		projectId, ok := params["projectId"]
		if ok {
			listOpts.ProjectID = projectId
		}
		status, ok := params["status"]
		if ok {
			listOpts.Status = status
		}
	}
	result := networks.List(networkClient, listOpts)
	resultPages, err := result.AllPages()
	if err != nil {
		fmt.Println(err)
	}
	resultBody := resultPages.GetBody()
	//fmt.Println(reflect.TypeOf(resultBody))
	//fmt.Println(result3)

	networkList := resultBody.(map[string][]interface{})["networks"]
	//fmt.Println(floatingIpList)

	return networkList
}

/**
	params
		- projectId
*/
func GetFloatingIps(session session.OpenStackSession, params map[string]string) []interface{} {
	// Floating IP 조회
	opts := gophercloud.EndpointOpts{Region: "RegionOne"}
	client, err := openstack.NewNetworkV2(session.Provider, opts)
	if err != nil {
		fmt.Println(err)
	}

	var listOpts floatingips.ListOpts
	if params != nil {
		projectId, ok := params["projectId"]
		if ok {
			listOpts.ProjectID = projectId
		}
		status, ok := params["status"]
		if ok {
			listOpts.Status = status
		}
	}

	result := floatingips.List(client, listOpts)
	resultPages, err := result.AllPages()
	resultBody := resultPages.GetBody()
	//fmt.Println(reflect.TypeOf(resultBody))
	//fmt.Println(result3)

	floatingIpList := resultBody.(map[string][]interface{})["floatingips"]
	//fmt.Println(floatingIpList)

	return floatingIpList
}

/**
	params
		- projectId
*/
func GetSecurityGroups(session session.OpenStackSession, params map[string]string) []interface{} {
	// Security Group 조회
	opts := gophercloud.EndpointOpts{Region: "RegionOne"}
	client, err := openstack.NewNetworkV2(session.Provider, opts)
	if err != nil {
		fmt.Println(err)
	}

	var listOpts groups.ListOpts
	if params != nil {
		projectId, ok := params["projectId"]
		if ok {
			listOpts.ProjectID = projectId
		}
	}

	result := groups.List(client, listOpts)
	resultPages, err := result.AllPages()
	resultBody := resultPages.GetBody()

	securityGroupList := resultBody.(map[string][]interface{})["security_groups"]
	//fmt.Println(securityGroupList)

	return securityGroupList
}


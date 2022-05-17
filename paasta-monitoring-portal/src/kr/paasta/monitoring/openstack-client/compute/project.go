package compute

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"monitoring-portal/openstack-client/session"
	"monitoring-portal/openstack-client/utils"
)

func GetProjectList(session session.OpenStackSession) {
	opts := gophercloud.EndpointOpts{ Region: "RegionOne" }
	identityClient, err := openstack.NewIdentityV3(session.Provider, opts)
	if err != nil {
		fmt.Println(err)
	}

	var listOpts projects.ListOpts
	resultPager := projects.List(identityClient, listOpts)
	list := utils.PagerToMap(resultPager)
	projectList := list.(map[string][]interface{})["projects"]
	utils.PrintJson(projectList)
}
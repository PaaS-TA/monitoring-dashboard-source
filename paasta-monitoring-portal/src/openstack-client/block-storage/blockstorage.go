package block_storage

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"openstack-client/session"
)

/****************************
	Block Storage (Cinder)
*****************************/

/**
	params
		- tenantId
 */
func GetVolumeList(session session.OpenStackSession, params map[string]string) []interface{} {
	// 볼륨 목록 조회
	var opts = gophercloud.EndpointOpts{Region: "RegionOne"}
	client, err := openstack.NewBlockStorageV3(session.Provider, opts)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println("Volume List")
	//volumes.List(client, nil).EachPage(func (page pagination.Page) (bool, error) {
	//	volumeList, err := volumes.ExtractVolumes(page)
	//	if err != nil {
	//		return false, err
	//	}
	//
	//	utils.PrintJson(volumeList)
	//	return true, nil;
	//})

	var listOpts volumes.ListOpts
	if params != nil {
		tenantId, ok := params["tenantId"]
		if ok {
			listOpts.TenantID = tenantId
		}
		status, ok := params["status"]
		if ok {
			listOpts.Status = status
		}
	}
	resultPager := volumes.List(client, listOpts)
	if resultPager.Err != nil {
		fmt.Println(resultPager.Err.Error())
	}
	resultPages, err := resultPager.AllPages()
	if err != nil {
		fmt.Println(err.Error())
	}
	resultBody := resultPages.GetBody()
	//fmt.Println(reflect.TypeOf(resultBody))
	//fmt.Println(result3)

	volumeList := resultBody.(map[string][]interface{})["volumes"]

	return volumeList
}

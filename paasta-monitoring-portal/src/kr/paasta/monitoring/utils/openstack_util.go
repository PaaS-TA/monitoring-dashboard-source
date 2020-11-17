package utils

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"kr/paasta/monitoring/iaas/model"
	//"fmt"
)

func GetComputeClient(provider *gophercloud.ProviderClient, region string) (*gophercloud.ServiceClient, error) {

	//fmt.Println("provider======+>", provider)
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: region,
	})
	if err != nil {
		model.MonitLogger.Error("GetComputeClient::", err)
		return client, err
	}
	return client, nil
}

func GetKeystoneClient(provider *gophercloud.ProviderClient) *gophercloud.ServiceClient {

	client := openstack.NewIdentityV3(provider)

	return client
}

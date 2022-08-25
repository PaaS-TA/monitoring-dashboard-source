package utils

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"

	"monitoring-portal/iaas_new/model"
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

func GetNetworkClient(provider *gophercloud.ProviderClient, region string) *gophercloud.ServiceClient {

	//fmt.Println("provider======+>", provider)
	client, _ := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Region: region,
	})

	return client
}

func GetKeystoneClient(provider *gophercloud.ProviderClient) *gophercloud.ServiceClient {
	client, _ := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	return client
}

func GetBlockStorageClient(provider *gophercloud.ProviderClient, region string) *gophercloud.ServiceClient {
	client, _ := openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{
		Region: region,
	})
	return client
}

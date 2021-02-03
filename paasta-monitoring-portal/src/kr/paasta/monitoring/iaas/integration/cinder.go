package integration

import (
	"encoding/json"
	"fmt"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"io/ioutil"
	"kr/paasta/monitoring/iaas/model"
	"kr/paasta/monitoring/utils"
)

type Cinder struct {
	OpenstackProvider model.OpenstackProvider
	Provider          *gophercloud.ProviderClient
}

func GetCinder(openstack_provider model.OpenstackProvider, provider *gophercloud.ProviderClient) *Cinder {
	return &Cinder{
		OpenstackProvider: openstack_provider,
		Provider:          provider,
	}
}

/**
Description : Get project Storage Max & Used information
*/
func (n *Cinder) GetTenantStorageResources(project_id, project_name string) (result model.TenantStorageResources, err error) {
	var data interface{}

	n.OpenstackProvider.TenantName = project_name

	//client for Cinder API operation
	client, err := openstack.NewBlockStorageV2(n.Provider, gophercloud.EndpointOpts{
		Region: n.OpenstackProvider.Region,
	})

	response, err := client.Get(fmt.Sprintf("%s/%s/%s/limits", model.CinderUrl, model.CinderVersion, project_id), nil, nil)
	if err != nil {
		return result, err
	}
	rawdata, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}
	json.Unmarshal(rawdata, &data)
	msg := data.(map[string]interface{})

	sub_msg := msg["limits"].(map[string]interface{})

	if len(sub_msg) < 1 {
		return result, nil
	}

	resources := sub_msg["absolute"].(map[string]interface{})

	if len(resources) > 0 {
		result.VolumeLimitGb = utils.TypeChecker_int(resources["maxTotalVolumeGigabytes"]).(int)
		result.VolumeGb = utils.TypeChecker_int(resources["totalGigabytesUsed"]).(int)
		result.VolumesLimit = utils.TypeChecker_int(resources["maxTotalVolumes"]).(int)
		result.Volumes = utils.TypeChecker_int(resources["totalVolumesUsed"]).(int)
		result.SnapshotsLimit = utils.TypeChecker_int(resources["maxTotalSnapshots"]).(int)
		result.Snapshots = utils.TypeChecker_int(resources["totalSnapshotsUsed"]).(int)
		result.BackupsLimit = utils.TypeChecker_int(resources["maxTotalBackups"]).(int)
		result.Backups = utils.TypeChecker_int(resources["totalBackupsUsed"]).(int)
	}
	return result, nil
}

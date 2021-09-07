package integration

import (
	"encoding/json"
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"io/ioutil"
	"kr/paasta/monitoring/iaas_new/model"
	"kr/paasta/monitoring/utils"
	"net/http"
)

type Nova struct {
	OpenstackProvider model.OpenstackProvider
	Provider          *gophercloud.ProviderClient
}

func GetNova(openstackProvider model.OpenstackProvider, provider *gophercloud.ProviderClient) *Nova {
	return &Nova{
		OpenstackProvider: openstackProvider,
		Provider:          provider,
	}
}

/**
Description : Get openstack's total & used resources - to check why admin project id is needed.
*/
func (n *Nova) GetOpenstackResources() (result model.HypervisorResources, err error) {

	client, err := utils.GetComputeClient(n.Provider, n.OpenstackProvider.Region)
	response, err := client.Get(fmt.Sprintf("%s/%s/%s/os-hypervisors/statistics", model.NovaUrl, model.NovaVersion, model.DefaultTenantId), nil, nil)
	msg, err := utils.ResponseUnmarshal(response, err)

	resources := msg["hypervisor_statistics"].(map[string]interface{})
	if len(resources) > 0 {
		result.VcpuTotal, _ = utils.TypeChecker_float64(resources["vcpus"]).(float64)
		result.VcpuUsed, _ = utils.TypeChecker_float64(resources["vcpus_used"]).(float64)
		result.DiskGbTotal, _ = utils.TypeChecker_float64(resources["local_gb"]).(float64)
		result.DiskGbUsed, _ = utils.TypeChecker_float64(resources["local_gb_used"]).(float64)
		result.DiskGbFree, _ = utils.TypeChecker_float64(resources["free_disk_gb"]).(float64)
		result.DiskGbLeastAvailable, _ = utils.TypeChecker_float64(resources["disk_available_least"]).(float64)
		result.MemoryMbTotal, _ = utils.TypeChecker_float64(resources["memory_mb"]).(float64)
		result.MemoryMbFree, _ = utils.TypeChecker_float64(resources["free_ram_mb"]).(float64)
		result.MemoryMbUsed, _ = utils.TypeChecker_float64(resources["memory_mb_used"]).(float64)
		//Api 에서 제공되는 running_vms 는 Total Vms 이다.
		result.VmTotal, _ = utils.TypeChecker_int(resources["running_vms"]).(int)
	}
	return result, err
}

/**
Description : Compute Node Resource Summary
*/
func (n *Nova) GetComputeNodeResources() (result []model.NodeResources, err error) {

	/*provider, err := utils.GetAdminToken(n.OpenstackProvider)
	if err != nil {
		return result, err
	}
	*/
	//client for Compute API operation
	client, err := utils.GetComputeClient(n.Provider, n.OpenstackProvider.Region)
	if err != nil {
		return result, err
	}

	response, err := client.Get(fmt.Sprintf("%s/%s/%s/os-hypervisors/detail", model.NovaUrl, model.NovaVersion, model.DefaultTenantId), nil, nil)

	if err != nil {
		return result, err
	}

	msg, err := utils.ResponseUnmarshal(response, err)

	resources := msg["hypervisors"].([]interface{})
	if len(resources) > 0 {

		for _, compute_node := range resources {
			node := compute_node.(map[string]interface{})
			var node_info model.NodeResources
			node_info.Id = utils.TypeChecker_int(node["id"]).(int)
			node_info.Hostname = utils.TypeChecker_string(node["hypervisor_hostname"]).(string)
			node_info.HostIp = utils.TypeChecker_string(node["host_ip"]).(string)
			node_info.Type = utils.TypeChecker_string(node["hypervisor_type"]).(string)
			node_info.VcpusMax = utils.TypeChecker_int(node["vcpus"]).(int)
			node_info.VcpusUsed = utils.TypeChecker_int(node["vcpus_used"]).(int)
			node_info.MemoryMbMax = utils.TypeChecker_int(node["memory_mb"]).(int)
			node_info.MemoryMbUsed = utils.TypeChecker_int(node["memory_mb_used"]).(int)
			node_info.MemoryMbFree = utils.TypeChecker_int(node["free_ram_mb"]).(int)
			node_info.DiskGbMax = utils.TypeChecker_int(node["local_gb"]).(int)
			node_info.DiskGbUsed = utils.TypeChecker_int(node["local_gb_used"]).(int)
			node_info.DiskGbFree = utils.TypeChecker_int(node["free_disk_gb"]).(int)
			node_info.DiskAvailableLeast = utils.TypeChecker_int(node["disk_available_least"]).(int)
			node_info.TotalVms = utils.TypeChecker_int(node["running_vms"]).(int)
			node_info.State = utils.TypeChecker_string(node["state"]).(string)
			node_info.Status = utils.TypeChecker_string(node["status"]).(string)

			result = append(result, node_info)
		}
	}
	return result, nil
}

/**
Description : Get project resources limit metadata(include network)
*/
func (n *Nova) GetProjectResourcesLimit(project_id string) (result model.TenantResourcesLimit, err error) {
	var data interface{}
	/*provider, err := utils.GetAdminToken(n.OpenstackProvider)
	if err != nil {
		return result, err
	}*/

	//client for Compute API operation
	client, err := openstack.NewComputeV2(n.Provider, gophercloud.EndpointOpts{
		Region: n.OpenstackProvider.Region,
	})

	response, err := client.Get(fmt.Sprintf("%s/%s/%s/os-quota-sets/%s", model.NovaUrl, model.NovaVersion, model.DefaultTenantId, project_id), nil, nil)
	if err != nil {
		return result, err
	}
	rawdata, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	json.Unmarshal(rawdata, &data)
	msg := data.(map[string]interface{})
	resources := msg["quota_set"].(map[string]interface{})
	if len(resources) > 0 {
		result.InstancesLimit = utils.TypeChecker_int(resources["instances"]).(int)
		result.MemoryMbLimit = utils.TypeChecker_int(resources["ram"]).(int)
		result.CoresLimit = utils.TypeChecker_int(resources["cores"]).(int)
		result.ServerGroupsLimit = utils.TypeChecker_int(resources["server_groups"]).(int)
		result.KeyPairsLimit = utils.TypeChecker_int(resources["key_pairs"]).(int)
	}
	return result, nil

}

func (n *Nova) GetProjectInstancesList(apiRequest model.TenantReq) (result []model.InstanceInfo, err error) {

	//var instanceList []model.InstanceInfo
	var data interface{}

	/*provider, err := utils.GetAdminToken(n.OpenstackProvider)
	if err != nil {
		return result, err
	}*/

	//client for Compute API operation
	client, err := openstack.NewComputeV2(n.Provider, gophercloud.EndpointOpts{
		Region: n.OpenstackProvider.Region,
	})

	var response *http.Response
	var apiError error

	if apiRequest.Marker == "" {

		if apiRequest.HostName == "" {
			response, apiError = client.Get(fmt.Sprintf("%s/%s/%s/servers/detail?all_tenants=1&limit=%s&project_id=%s",
				model.NovaUrl, model.NovaVersion, model.DefaultTenantId, apiRequest.Limit, apiRequest.TenantId), nil, nil)
		} else {
			response, apiError = client.Get(fmt.Sprintf("%s/%s/%s/servers/detail?all_tenants=1&limit=%s&name=%s&project_id=%s",
				model.NovaUrl, model.NovaVersion, model.DefaultTenantId, apiRequest.Limit, apiRequest.HostName, apiRequest.TenantId), nil, nil)
		}

	} else {

		if apiRequest.HostName == "" {
			response, apiError = client.Get(fmt.Sprintf("%s/%s/%s/servers/detail?all_tenants=1&limit=%s&marker=%s&project_id=%s",
				model.NovaUrl, model.NovaVersion, model.DefaultTenantId, apiRequest.Limit, apiRequest.Marker, apiRequest.TenantId), nil, nil)
		} else {
			response, apiError = client.Get(fmt.Sprintf("%s/%s/%s/servers/detail?all_tenants=1&limit=%s&marker=%s&name=%s&project_id=%s",
				model.NovaUrl, model.NovaVersion, model.DefaultTenantId, apiRequest.Limit, apiRequest.Marker, apiRequest.HostName, apiRequest.TenantId), nil, nil)
		}

	}

	//fmt.Printf("%s/%s/%s/server/details/%s?limit=5", model.NovaUrl, model.NovaVersion, model.DefaultProjectId, project_id)
	if apiError != nil {
		fmt.Println("Err :", apiError)
		return result, err
	}
	rawdata, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println("Err :", err)
		return result, err
	}

	json.Unmarshal(rawdata, &data)
	msg := data.(map[string]interface{})

	servers := msg["servers"].([]interface{})

	for _, v := range servers {
		var instanceInfo model.InstanceInfo
		inst := v.(map[string]interface{})

		instanceInfo.Zone = utils.TypeChecker_string(inst["OS-EXT-AZ:availability_zone"]).(string)
		instanceInfo.Name = utils.TypeChecker_string(inst["name"]).(string)
		instanceInfo.InstanceId = utils.TypeChecker_string(inst["id"]).(string)
		instanceInfo.State = utils.TypeChecker_string(inst["status"]).(string)

		addInfos := inst["addresses"].(map[string]interface{})
		var ipAddress []string
		for _, addrDetail := range addInfos {

			addrList := addrDetail.([]interface{})
			for _, address := range addrList {
				instanceAddr := address.(map[string]interface{})

				ipAddress = append(ipAddress, utils.TypeChecker_string(instanceAddr["addr"]).(string))
			}

		}

		instanceInfo.Address = ipAddress

		result = append(result, instanceInfo)
	}

	return result, nil
}

/**
Description : Get Project's created instance list
*/
func (n *Nova) GetProjectInstances(project_id string) (result []model.InstanceInfo, err error) {
	//var return_value model.HypervisorResources
	var data interface{}

	/*provider, err := utils.GetAdminToken(n.OpenstackProvider)
	if err != nil {
		return result, err
	}*/

	//client for Compute API operation
	client, err := openstack.NewComputeV2(n.Provider, gophercloud.EndpointOpts{
		Region: n.OpenstackProvider.Region,
	})

	response, err := client.Get(fmt.Sprintf("%s/%s/%s/os-simple-tenant-usage/%s", model.NovaUrl, model.NovaVersion, model.DefaultTenantId, project_id), nil, nil)

	if err != nil {
		return result, err
	}
	rawdata, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}
	json.Unmarshal(rawdata, &data)
	msg := data.(map[string]interface{})
	map_tenant_usage := msg["tenant_usage"].(map[string]interface{})
	if len(map_tenant_usage) < 1 {
		return result, err
	}
	resources := map_tenant_usage["server_usages"].([]interface{})

	if len(resources) > 0 {
		for _, v := range resources {
			inst := v.(map[string]interface{})
			var instance model.InstanceInfo
			instance.TenantId = utils.TypeChecker_string(inst["tenant_id"]).(string)
			instance.InstanceId = utils.TypeChecker_string(inst["instance_id"]).(string)
			instance.Name = utils.TypeChecker_string(inst["name"]).(string)
			instance.Flavor = utils.TypeChecker_string(inst["flavor"]).(string)
			instance.Vcpus = utils.TypeChecker_float64(inst["vcpus"]).(float64)
			instance.DiskGb = utils.TypeChecker_float64(inst["local_gb"]).(float64)
			instance.MemoryMb = utils.TypeChecker_float64(inst["memory_mb"]).(float64)
			instance.State = utils.TypeChecker_string(inst["state"]).(string)
			instance.StartedAt = utils.TypeChecker_string(inst["started_at"]).(string)
			if len(instance.StartedAt) > 19 {
				instance.StartedAt = instance.StartedAt[0:19]
			}
			instance.EndedAt = utils.TypeChecker_string(inst["ended_at"]).(string)
			if len(instance.EndedAt) > 19 {
				instance.EndedAt = instance.EndedAt[0:19]
			}
			instance.Uptime = utils.TypeChecker_float64(inst["uptime"]).(float64)
			result = append(result, instance)
		}
	}
	return result, nil
}

/**
Description : Nova API Test
*/
func (n *Nova) GetInstanceDetail(instance_id, project_id string) (result model.InstanceDetail, err error) {
	var data interface{}

	/*provider, err := utils.GetAdminToken(n.OpenstackProvider)
	if err != nil {
		return result, err
	}*/

	//client for Compute API operation
	client, err := openstack.NewComputeV2(n.Provider, gophercloud.EndpointOpts{
		Region: n.OpenstackProvider.Region,
	})

	response, err := client.Get(fmt.Sprintf("%s/%s/%s/servers/%s", model.NovaUrl, model.NovaVersion, project_id, instance_id), nil, nil)
	if err != nil {
		return result, err
	}
	rawdata, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}
	json.Unmarshal(rawdata, &data)
	msg := data.(map[string]interface{})

	instance_info := msg["server"].(map[string]interface{})
	//########### Instance Bastic Info ##############
	result.Id = utils.TypeChecker_string(instance_info["id"]).(string)
	result.Name = utils.TypeChecker_string(instance_info["name"]).(string)
	result.ProcessName = utils.TypeChecker_string(instance_info["OS-EXT-SRV-ATTR:instance_name"]).(string)
	result.AvailabilityZone = utils.TypeChecker_string(instance_info["OS-EXT-AZ:availability_zone"]).(string)
	result.CreatedDate = utils.TypeChecker_string(instance_info["created"]).(string)

	metadata_info := instance_info["metadata"].(map[string]interface{})
	if len(metadata_info) > 0 {
		result.Deployment.Deployment = utils.TypeChecker_string(metadata_info["deployment"]).(string)
		result.Deployment.Director = utils.TypeChecker_string(metadata_info["director"]).(string)
		result.Deployment.Job = utils.TypeChecker_string(metadata_info["job"]).(string)
		result.Deployment.Name = utils.TypeChecker_string(metadata_info["name"]).(string)
	}

	address_info := instance_info["addresses"].(map[string]interface{})
	if len(address_info) > 0 {
		instance_network := make([]model.NetworkInfo, len(address_info["private_net"].([]interface{})))
		var index int
		for _, network := range address_info["private_net"].([]interface{}) {
			instance_network[index].Ip = utils.TypeChecker_string(network.(map[string]interface{})["addr"]).(string)
			instance_network[index].Type = utils.TypeChecker_string(network.(map[string]interface{})["OS-EXT-IPS:type"]).(string)
			instance_network[index].Mac_addr = utils.TypeChecker_string(network.(map[string]interface{})["mac_addr"]).(string)

			index = index + 1
		}
		result.Network = instance_network
	}

	//########### Instance Image Info ############
	image_info := instance_info["image"].(map[string]interface{})

	//########### Instacne Flavor Info ###########
	flavor_info := instance_info["flavor"].(map[string]interface{})

	//########### Instance Security Group ########
	security_groups := instance_info["security_groups"].([]interface{})

	if len(security_groups) > 0 {
		var s_groups string
		for _, security_group := range security_groups {
			s_g := security_group.(map[string]interface{})
			s_groups = s_groups + s_g["name"].(string) + " "
		}
		result.SecurityGroups = s_groups
	}

	response, err = client.Get(fmt.Sprintf("%s/%s/%s/flavors/%s", model.NovaUrl, model.NovaVersion, project_id, flavor_info["id"]), nil, nil)
	if err != nil {
		return result, err
	}
	rawdata, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}
	json.Unmarshal(rawdata, &data)
	msg = data.(map[string]interface{})

	flavor_detail := msg["flavor"].(map[string]interface{})
	if len(flavor_detail) > 0 {
		result.Flavor.Name = utils.TypeChecker_string(flavor_detail["name"]).(string)
		result.Flavor.Vcpu = utils.TypeChecker_int(flavor_detail["vcpus"]).(int)
		result.Flavor.Memory = utils.TypeChecker_int(flavor_detail["ram"]).(int)
		result.Flavor.Disk = utils.TypeChecker_int(flavor_detail["disk"]).(int)
	}

	response, err = client.Get(fmt.Sprintf("%s/%s/%s/images/%s", model.NovaUrl, model.NovaVersion, project_id, image_info["id"].(string)), nil, nil)
	if err != nil {
		return result, err
	}
	rawdata, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}
	json.Unmarshal(rawdata, &data)
	msg = data.(map[string]interface{})
	i_info := msg["image"].(map[string]interface{})

	if len(i_info) > 0 {
		result.Image.Id = utils.TypeChecker_string(i_info["id"]).(string)
		result.Image.Name = utils.TypeChecker_string(i_info["name"]).(string)

		resources := i_info["metadata"].(map[string]interface{})
		if len(resources) > 0 {
			result.Image.Version = utils.TypeChecker_string(i_info["version"]).(string)
			result.Image.OsType = utils.TypeChecker_string(i_info["os_type"]).(string)
			result.Image.OsKind = utils.TypeChecker_string(i_info["os_distro"]).(string)
			result.Image.HypervisorType = utils.TypeChecker_string(i_info["hypervisor_type"]).(string)
		}
	}
	return result, err
}

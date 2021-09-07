package integration

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"kr/paasta/monitoring/iaas_new/model"
	"kr/paasta/monitoring/utils"
)

type Keystone struct {
	OpenstackProvider model.OpenstackProvider
	provider          *gophercloud.ProviderClient
}

func GetKeystone(openstack_provider model.OpenstackProvider, provider *gophercloud.ProviderClient) *Keystone {
	return &Keystone{
		OpenstackProvider: openstack_provider,
		provider:          provider,
	}
}

/**
Description : Get openstack's registered project lists
*/
func (k *Keystone) GetTenantList() (tenant_lists []model.TenantInfo, err error) {

	//client for Keystone API operation
	client, _ := openstack.NewIdentityV3(k.provider, gophercloud.EndpointOpts{})
	response, err := client.Get(fmt.Sprintf("%s/%s/projects", model.KeystoneUrl, model.KeystoneVersion), nil, nil)

	msg, err := utils.ResponseUnmarshal(response, err)

	for _, v := range msg["projects"].([]interface{}) {
		var tenant_info model.TenantInfo
		v_detail := v.(map[string]interface{})
		tenant_info.Id = utils.TypeChecker_string(v_detail["id"]).(string)
		tenant_info.Name = utils.TypeChecker_string(v_detail["name"]).(string)
		tenant_info.IsDomain = v_detail["is_domain"].(bool)
		tenant_info.Description = utils.TypeChecker_string(v_detail["description"]).(string)
		tenant_info.Enabled = v_detail["enabled"].(bool)
		tenant_info.ParentId = utils.TypeChecker_string(v_detail["parent_id"]).(string)
		tenant_info.Links = v_detail["links"].(map[string]interface{})

		tenant_lists = append(tenant_lists, tenant_info)
	}
	return tenant_lists, nil
}

func (k *Keystone) GetUserIdByName(userName string) (userId string, err error) {

	/*provider, err := utils.GetAdminToken(k.OpenstackProvider)
	if err != nil {
		return userId, err
	}*/
	//client for Keystone API operation
	client, _ := openstack.NewIdentityV3(k.provider, gophercloud.EndpointOpts{})

	response, err := client.Get(fmt.Sprintf("%s/%s/users?name=%s", model.KeystoneUrl, model.KeystoneVersion, userName), nil, nil)
	msg, err := utils.ResponseUnmarshal(response, err)

	for _, v := range msg["users"].([]interface{}) {
		v_detail := v.(map[string]interface{})
		userId = utils.TypeChecker_string(v_detail["id"]).(string)
	}
	return userId, err
}

func (k *Keystone) GetUserTenantList(userId string) (project_lists []model.TenantInfo, err error) {

	client, _ := openstack.NewIdentityV3(k.provider, gophercloud.EndpointOpts{})

	response, err := client.Get(fmt.Sprintf("%s/%s/users/%s/projects", model.KeystoneUrl, model.KeystoneVersion, userId), nil, nil)
	msg, err := utils.ResponseUnmarshal(response, err)

	for _, v := range msg["projects"].([]interface{}) {
		var tenant_info model.TenantInfo
		v_detail := v.(map[string]interface{})
		tenant_info.Id = utils.TypeChecker_string(v_detail["id"]).(string)
		tenant_info.Name = utils.TypeChecker_string(v_detail["name"]).(string)
		tenant_info.IsDomain = v_detail["is_domain"].(bool)
		tenant_info.Description = utils.TypeChecker_string(v_detail["description"]).(string)
		tenant_info.Enabled = v_detail["enabled"].(bool)
		tenant_info.ParentId = utils.TypeChecker_string(v_detail["parent_id"]).(string)
		tenant_info.Links = v_detail["links"].(map[string]interface{})

		project_lists = append(project_lists, tenant_info)
	}
	return project_lists, nil
}

func (k *Keystone) RevokeToken() tokens.RevokeResult {

	client, _ := openstack.NewIdentityV3(k.provider, gophercloud.EndpointOpts{})
	result := tokens.Revoke(client, client.TokenID)

	return result
}

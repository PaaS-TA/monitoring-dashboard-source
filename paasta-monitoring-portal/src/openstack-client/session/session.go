package session

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

type OpenStackSession struct {
	identityEndpoint string
	username         string
	password         string
	domainId         string
	AuthOpts         gophercloud.AuthOptions
	Provider         *gophercloud.ProviderClient
}

func (session OpenStackSession) CreateSession(params map[string]string) {
	session.AuthOpts = gophercloud.AuthOptions{
		IdentityEndpoint: params["identityEndpoint"],
		Username: params["username"],
		Password: params["password"],
		DomainID: params["domainid"],
	}

	provider, err := openstack.AuthenticatedClient(session.AuthOpts)
	if err != nil {
		fmt.Println(err)
	}
	session.Provider = provider
}
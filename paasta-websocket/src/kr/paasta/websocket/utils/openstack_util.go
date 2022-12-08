package utils

import (
	"bufio"
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func GetComputeClient(provider *gophercloud.ProviderClient, region string) (*gophercloud.ServiceClient, error) {

	//fmt.Println("provider======+>", provider)
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: region,
	})
	if err != nil {
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


func GetOpenstackProvider(r *http.Request) (provider *gophercloud.ProviderClient, username string, err error) {

	//config, err := ReadConfig(`../../config.ini`) // test
	config, err := ReadConfig(`config.ini`) // real


	opts := gophercloud.AuthOptions{
		//IdentityEndpoint: config["keystone.url"],
		IdentityEndpoint: config["identity.endpoint"],
		Username:         config["default.username"],
		Password: config["default.password"],
		//TenantName:  config["default.tenant_name"],
		//DomainName:  config["default.domain"],
		TenantID:    config["default.tenant_id"],
		//TokenID:     val["iaasToken"],
		//AllowReauth: false,
	}

	//Provider is the top-level client that all of your OpenStack services
	providerClient, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		fmt.Println(err.Error())
		return nil, "", err
	}
	log.Println("providerClient.TokenID : " + providerClient.TokenID)
	//새로 로그인 되었으므로 변경된 토큰으로 변경하여 저장
	//rclient.HSet(reqToken, "iaasToken", providerClient.TokenID)

	return providerClient, opts.Username, err
}

type Config map[string]string

func ReadConfig(filename string) (Config, error) {
	// init with some bogus data
	config := Config{
		"server.ip":   "127.0.0.1",
		"server.port": "8888",
	}

	if len(filename) == 0 {
		return config, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		// check if the line has = sign
		// and process the line. Ignore the rest.
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				// assign the config map
				config[key] = value
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}

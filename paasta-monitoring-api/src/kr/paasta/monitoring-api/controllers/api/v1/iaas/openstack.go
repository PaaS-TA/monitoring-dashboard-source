package iaas

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v7"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"paasta-monitoring-api/apiHelpers"
	service "paasta-monitoring-api/services/api/v1/iaas"
)

type (
	OpenstackController struct {
		OpenstackProvider *gophercloud.ProviderClient
	}
)

const (
	CSRF_TOKEN_NAME   = "X-XSRF-TOKEN"
)

func GetOpenstackProvider(r *http.Request) (provider *gophercloud.ProviderClient, err error) {


	// .env 파일 로드
	envLoadErr := godotenv.Load(".env")
	if envLoadErr != nil {
		log.Fatal("Error loading .env file")
	}

	reqToken := r.Header.Get(CSRF_TOKEN_NAME)
	rclient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis_url"),
		Password: os.Getenv("redis_password"),
	})

	iaasTokenData := rclient.HGetAll(reqToken).Val()

	//fmt.Println("GetOpenstackProvider redis.iaas_userid:::",val[""],"/",val["iaasToken"])

	opts := gophercloud.AuthOptions{
		IdentityEndpoint : os.Getenv("openstack.identity_endpoint"),
		DomainName : 	   os.Getenv("openstack.domain"),
		TenantID :   	   os.Getenv("openstack.tenant_id"),
		TenantName : 	   os.Getenv("openstack.tenant_name"),
		Username : 	 	   iaasTokenData["iaasUserId"],
		TokenID :	 	   iaasTokenData["iaasToken"],
		AllowReauth : 	   false,
	}

	//Provider is the top-level client that all of your OpenStack services
	providerClient, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	//새로 로그인 되었으므로 변경된 토큰으로 변경하여 저장
	rclient.HSet(reqToken, "iaasToken", providerClient.TokenID)

	return providerClient, err
}


func GetOpenstackController(openstackProvider *gophercloud.ProviderClient) *OpenstackController {
	return &OpenstackController{
		OpenstackProvider: openstackProvider,
	}
}

func (controller *OpenstackController) GetHypervisorStatistics(ctx echo.Context) error {
	fmt.Println(controller.OpenstackProvider.TokenID)
	results, err := service.GetOpenstackService(controller.OpenstackProvider).GetHypervisorStatistics()
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "", results)
	}

	return nil
}

package connections

import (
	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/jinzhu/gorm"
	"os"
	"paasta-monitoring-api/helpers"
	"strings"
)

type Connections struct {
	DbInfo    *gorm.DB
	RedisInfo *redis.Client
	OpenstackProvider *gophercloud.ProviderClient
}

func setEnv(env map[string]string) {

	// Redis 설정
	env["redis_url"] = os.Getenv("redis_url")

	// PaaS DataBase 설정
	env["paas_db_type"] = os.Getenv("paas_db_type")
	env["paas_db_password"] = os.Getenv("paas_db_password")
	env["paas_db_username"] = os.Getenv("paas_db_username")
	env["paas_db_protocol"] = os.Getenv("paas_db_protocol")
	env["paas_db_host"] = os.Getenv("paas_db_host")
	env["paas_db_port"] = os.Getenv("paas_db_port")
	env["paas_db_name"] = os.Getenv("paas_db_name")
	env["paas_db_charset"] = os.Getenv("paas_db_charset")
	env["paas_db_parseTime"] = os.Getenv("paas_db_parseTime")

}

func redisConnection(env map[string]string) *redis.Client {

	// Redis 설정
	dsn := env["redis_url"]
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		panic(err)
	}
	return redisClient
}

func paasConnection(env map[string]string) *gorm.DB {

	// DB 설정
	paasDBClient, paasDBErr := gorm.Open(
		helpers.GetDBConnectionString(
			env["paas_db_type"],
			env["paas_db_username"],
			env["paas_db_password"],
			env["paas_db_protocol"],
			env["paas_db_host"],
			env["paas_db_port"],
			env["paas_db_name"],
			env["paas_db_charset"],
			env["paas_db_parseTime"]),
	)

	if paasDBErr != nil {
		panic(paasDBErr)
	}

	return paasDBClient
}

func saasConnection(env map[string]string) error {
	return nil
}

func caasConnection(env map[string]string) error {
	return nil
}

func iaasConnection(env map[string]string) error {
	return nil
}

func openstackConnection(connections Connections) *gophercloud.ProviderClient {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint : os.Getenv("openstack.identity_endpoint"),
		Username   : os.Getenv("openstack.username"),
		Password   : os.Getenv("openstack.password"),
		TenantID   :   	   os.Getenv("openstack.tenant_id"),
		TenantName : 	   os.Getenv("openstack.tenant_name"),
		DomainName : 	   os.Getenv("openstack.domain"),
		//TokenID :	 	   iaasTokenData["iaasToken"],
		//AllowReauth : 	   false,
	}

	providerClient, _ := openstack.AuthenticatedClient(opts)

	//openstackToken := providerClient.TokenID

	//새로 로그인 되었으므로 변경된 토큰으로 변경하여 저장
	//connections.RedisInfo.HSet(reqToken, "iaasToken", providerClient.TokenID)
	return providerClient
}

func SetupConnection() Connections {
	conn := Connections{}

	// 환경 변수 설정
	env := make(map[string]string)
	setEnv(env)

	// 내부 서비스 환경 설정
	conn.RedisInfo = redisConnection(env)

	services := os.Getenv("services")
	servicesArr := strings.Split(services, ",")

	// 외부 서비스 별 환경 설정
	for _, value := range servicesArr {
		switch value {
		case "PaaS":
			conn.DbInfo = paasConnection(env)
		case "SaaS":
			saasConnection(env)
		case "CaaS":
			caasConnection(env)
		case "IaaS":
			iaasConnection(env)
			conn.OpenstackProvider = openstackConnection(conn)
		}
	}
	return conn
}

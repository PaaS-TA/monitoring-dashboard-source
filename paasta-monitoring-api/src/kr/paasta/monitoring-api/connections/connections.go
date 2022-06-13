package connections

import (
	"fmt"
	"github.com/cloudfoundry-community/gogobosh"
	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/jinzhu/gorm"
	"net/http"
	"os"
	"paasta-monitoring-api/helpers"
	models "paasta-monitoring-api/models/api/v1"
	"strconv"
	"strings"
)

type Connections struct {
	DbInfo         *gorm.DB
	RedisInfo      *redis.Client
	InfluxDBClient client.Client
	BoshInfoList   []models.Bosh
	Env            map[string]interface{}
	OpenstackProvider *gophercloud.ProviderClient
}

/*
	Read for environment variables including variables of system and program
 */
func getEnv(envData []string, getKeyVal func(item string) (key, value string)) map[string]interface{} {
	envMap := make(map[string]interface{})
	for _, item := range envData {
		key, value := getKeyVal(item)
		envMap[key] = value
	}
	return envMap
}

func SetEnv() map[string]interface{} {
	envMap := getEnv(os.Environ(), func(item string) (key, value string) {
		keyValueSplit := strings.Split(item, "=")
		key = keyValueSplit[0]
		value = keyValueSplit[1]
		return
	})

	for key,value := range envMap {
		if strings.Contains(key, "openstack") {
			fmt.Println(key + " : " + value.(string))
		}
	}

	/*
	env := make(map[string]interface{})

	// Redis 설정
	env["redis_url"] = os.Getenv("redis_url")
	env["redis_db"], _ = strconv.Atoi(os.Getenv("redis_db"))

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
	*/
	return envMap
}

func GetBoshInfoList(env map[string]interface{}) []models.Bosh {
	// Bosh 설정
	BoshCount, _ := strconv.Atoi(os.Getenv("bosh_count"))
	var BoshList []models.Bosh
	for i := 0; i < BoshCount; i++ {
		var bosh models.Bosh
		bosh.UUID = os.Getenv("bosh_" + strconv.Itoa(i) + "_uuid")
		bosh.Name = os.Getenv("bosh_" + strconv.Itoa(i) + "_name")
		bosh.Ip = os.Getenv("bosh_" + strconv.Itoa(i) + "_ip")
		bosh.Deployname = os.Getenv("bosh_" + strconv.Itoa(i) + "_deployname")
		bosh.Address = os.Getenv("bosh_" + strconv.Itoa(i) + "_client_api_address")
		bosh.Username = os.Getenv("bosh_" + strconv.Itoa(i) + "_client_api_username")
		bosh.Password = os.Getenv("bosh_" + strconv.Itoa(i) + "_client_api_password")
		bosh.Client = BoshConnection(bosh)
		BoshList = append(BoshList, bosh)
	}
	return BoshList
}

func BoshConnection(bosh models.Bosh) *gogobosh.Client {

	// BOSH Client Config
	boshConfig := &gogobosh.Config{
		BOSHAddress:       bosh.Address,
		Username:          bosh.Username,
		Password:          bosh.Password,
		HttpClient:        http.DefaultClient,
		SkipSslValidation: true,
	}
	boshClient, err := gogobosh.NewClient(boshConfig)
	if err != nil {
		fmt.Errorf("Failed to create connection to the bosh server. err=", err)
	}

	return boshClient
}

func RedisConnection(env map[string]interface{}) *redis.Client {

	// Redis 설정
	dsn := env["redis_url"].(string)
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	redisDbName, _ := strconv.Atoi(env["redis_db"].(string))
	redisClient := redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
		Password: env["redis_password"].(string),
		DB: redisDbName,
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		panic(err)
	}
	return redisClient
}

func PaaSConnection(env map[string]interface{}) *gorm.DB {
	// DB 설정
	paasDBClient, paasDBErr := gorm.Open(
		helpers.GetDBConnectionString(
			env["paas_db_type"].(string),
			env["paas_db_username"].(string),
			env["paas_db_password"].(string),
			env["paas_db_protocol"].(string),
			env["paas_db_host"].(string),
			env["paas_db_port"].(string),
			env["paas_db_name"].(string),
			env["paas_db_charset"].(string),
			env["paas_db_parseTime"].(string)),
	)

	if paasDBErr != nil {
		panic(paasDBErr)
	}

	return paasDBClient
}

func SaasConnection(env map[string]interface{}) error {
	return nil
}

func CaaSConnection(env map[string]interface{}) error {
	return nil
}

func iaasConnection(env map[string]interface{}) error {
	return nil
}

func InfluxDBConnection(env map[string]interface{}) client.Client {

	influxDBClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:               env["paas_metric_db_url"].(string),
		Username:           env["paas_metric_db_username"].(string),
		Password:           env["paas_metric_db_password"].(string),
		InsecureSkipVerify: true,
	})
	if err != nil {
		panic(err)
	}
	return influxDBClient
}


func(connectionStruct *Connections ) openstackConnection(env map[string]interface{}) *gophercloud.ProviderClient {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint : env["openstack_identity_endpoint"].(string),
		Username         : env["openstack_username"].(string),
		Password         : env["openstack_password"].(string),
		TenantID         : env["openstack_tenant_id"].(string),
		AllowReauth      : true,

	}
	providerClient, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Openstack TokenID : " + providerClient.TokenID)
	//openstackToken := providerClient.TokenID

	// TODO Openstack 토큰 적재 방식 수립 필요
	//새로 로그인 되었으므로 변경된 토큰으로 변경하여 저장
	//connections.RedisInfo.HSet(reqToken, "iaasToken", providerClient.TokenID)
	return providerClient
}

func SetupConnection() Connections {

	Conn := Connections{}

	// 환경 변수 설정
	Conn.Env = SetEnv()

	// 내부 서비스 환경 설정
	Conn.RedisInfo = RedisConnection(Conn.Env)

	// 보쉬 정보 등록
	Conn.BoshInfoList = GetBoshInfoList(Conn.Env)

	services := os.Getenv("services")
	servicesArr := strings.Split(services, ",")

	// 외부 서비스 별 환경 설정
	for _, value := range servicesArr {
		switch value {
		case "PaaS":
			Conn.DbInfo = PaaSConnection(Conn.Env)
			Conn.InfluxDBClient = InfluxDBConnection(Conn.Env)
		case "SaaS":
			SaasConnection(Conn.Env)
		case "CaaS":
			CaaSConnection(Conn.Env)
		case "IaaS":
			iaasConnection(Conn.Env)
			Conn.OpenstackProvider = openstackConnection(Conn.Env)
		}
	}
	return Conn
}

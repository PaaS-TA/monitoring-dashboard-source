package connections

import (
	"crypto/tls"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/cloudfoundry-community/gogobosh"
	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"paasta-monitoring-api/helpers"
	"paasta-monitoring-api/middlewares/zabbix-client/lib/go-zabbix"
	models "paasta-monitoring-api/models/api/v1"
	"strconv"
	"strings"
)

type Connections struct {
	DbInfo            *gorm.DB
	RedisInfo         *redis.Client
	InfluxDbClient    models.InfluxDbClient
	BoshInfoList      []models.Bosh
	Env               map[string]interface{}
	OpenstackProvider *gophercloud.ProviderClient
	ZabbixSession     *zabbix.Session
	Logger            *logrus.Logger
	CfClient          *cfclient.Client
	CP                models.CP
	SaaS              models.SaaS
}

func SetupConnection(logger *logrus.Logger) Connections {
	connection := Connections{
		Env:    setEnv(),
		Logger: logger,
	}

	// 내부 서비스 환경 설정
	connection.redisConnection(connection.Env)

	// 보쉬 정보 등록
	connection.BoshInfoList = GetBoshInfoList(connection.Env)

	services := os.Getenv("services")
	servicesArr := strings.Split(services, ",")

	// 외부 서비스 별 환경 설정
	for _, value := range servicesArr {
		switch value {
		case "PaaS":
			connection.initPaasConfig(connection.Env)
			connection.initInfluxDbConfig(connection.Env)
		case "CaaS":
			connection.initCaaSConfig()
		case "IaaS":
			connection.initOpenstackProvider()
			connection.initZabbixSession()
		case "SaaS":
			connection.initSaaSConfig()
		}
	}
	return connection
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

func setEnv() map[string]interface{} {
	envMap := getEnv(os.Environ(), func(item string) (key, value string) {
		keyValueSplit := strings.Split(item, "=")
		key = keyValueSplit[0]
		value = keyValueSplit[1]
		return
	})
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

func (conn *Connections) redisConnection(env map[string]interface{}) {

	// Redis 설정
	dsn := env["redis_url"].(string)
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	redisDbName, _ := strconv.Atoi(env["redis_db"].(string))
	redisClient := redis.NewClient(&redis.Options{
		Addr:     dsn, //redis port
		Password: env["redis_password"].(string),
		DB:       redisDbName,
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		conn.Logger.Panic(err)
	}
	conn.RedisInfo = redisClient
}

func (conn *Connections) initPaasConfig(env map[string]interface{}) {
	// DB 설정
	dsn := helpers.GetDBConnectionString(
		env["paas_db_username"].(string),
		env["paas_db_password"].(string),
		env["paas_db_protocol"].(string),
		env["paas_db_host"].(string),
		env["paas_db_port"].(string),
		env["paas_db_name"].(string),
		env["paas_db_charset"].(string),
		env["paas_db_parseTime"].(string))

	paasDBClient, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		conn.Logger.Panic(err)
	}

	if strings.Compare(os.Getenv("mode"), "development") == 0 {
		paasDBClient.Debug()
	}
	conn.DbInfo = paasDBClient

	// Cloud Foundry Client Initialize
	config := &cfclient.Config{
		ApiAddress:        env["paas_cf_client_api_address"].(string),
		Username:          env["paas_cf_client_username"].(string),
		Password:          env["paas_cf_client_password"].(string),
		SkipSslValidation: true,
	}
	cloudFoundryClient, _ := cfclient.NewClient(config)
	conn.CfClient = cloudFoundryClient
}

func (conn *Connections) initInfluxDbConfig(env map[string]interface{}) {
	httpClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:               env["paas_metric_db_url"].(string),
		Username:           env["paas_metric_db_username"].(string),
		Password:           env["paas_metric_db_password"].(string),
		InsecureSkipVerify: true,
	})

	if err != nil {
		conn.Logger.Panic(err)
	}

	DbName := models.InfluxDbName{
		BoshDatabase:      env["paas_metric_db_name_bosh"].(string),
		PaastaDatabase:    env["paas_metric_db_name_paasta"].(string),
		ContainerDatabase: env["paas_metric_db_name_container"].(string),
		LoggingDatabase:   env["paas_metric_db_name_logging"].(string),
	}

	influxDbClient := models.InfluxDbClient{
		HttpClient: httpClient,
		DbName:     DbName,
	}
	conn.InfluxDbClient = influxDbClient
}

func (conn *Connections) initOpenstackProvider() {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: conn.Env["openstack_identity_endpoint"].(string),
		DomainName:       conn.Env["openstack_domain"].(string),
		Username:         conn.Env["openstack_username"].(string),
		Password:         conn.Env["openstack_password"].(string),
		TenantID:         conn.Env["openstack_tenant_id"].(string),
		AllowReauth:      true,
	}
	providerClient, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		log.Println(err.Error())
	}
	conn.Logger.Info("Openstack TokenID : " + providerClient.TokenID)
	//openstackToken := providerClient.TokenID

	// TODO Openstack 토큰 적재 방식 수립 필요
	//새로 로그인 되었으므로 변경된 토큰으로 변경하여 저장
	//connections.RedisInfo.HSet(reqToken, "iaasToken", providerClient.TokenID)
	conn.OpenstackProvider = providerClient
}

func (conn *Connections) initZabbixSession() {
	zabbixHost := conn.Env["zabbix_host"].(string)
	zabbixAdminId := conn.Env["zabbix_admin_id"].(string)
	zabbixAdminPw := conn.Env["zabbix_admin_pw"].(string)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	cache := zabbix.NewSessionFileCache().SetFilePath("./zabbix_session")
	zabbixSession, err := zabbix.CreateClient(zabbixHost).
		WithCache(cache).
		WithHTTPClient(client).
		WithCredentials(zabbixAdminId, zabbixAdminPw).Connect()
	if err != nil {
		fmt.Println(err)
	}
	conn.Logger.Info("Zabbix Token : " + zabbixSession.Token)
	conn.ZabbixSession = zabbixSession
}

func (conn *Connections) initCaaSConfig() {
	cp := models.CP{
		PromethusUrl:      conn.Env["prometheus_host"].(string),
		PromethusRangeUrl: conn.Env["prometheus_host"].(string) + "/api/v1/query_range?query=",
		K8sUrl:            conn.Env["kubernetes_host"].(string),
		K8sAdminToken:     conn.Env["kubernetes_admin_token"].(string),
	}

	conn.CP = cp
}

func (conn *Connections) initSaaSConfig() {
	saas := models.SaaS{
		PinpointWebUrl: conn.Env["pinpoint_web_url"].(string),
		Logger:         conn.Logger,
	}
	conn.SaaS = saas
}

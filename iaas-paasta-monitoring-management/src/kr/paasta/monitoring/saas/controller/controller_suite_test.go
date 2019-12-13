package controller

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/cihub/seelog"
	//"github.com/cloudfoundry-community/go-cfclient"
	"github.com/cloudfoundry-community/gogobosh"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gophercloud/gophercloud"
	"github.com/influxdata/influxdb1-client/v2"
	"github.com/jinzhu/gorm"
	"github.com/monasca/golang-monascaclient/monascaclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/olivere/elastic.v3"
	"io"
	"io/ioutil"
	comModels "kr/paasta/monitoring/common/model"
	iaasModels "kr/paasta/monitoring/iaas/model"
	bm "kr/paasta/monitoring/paas/model"
	paasModels "kr/paasta/monitoring/paas/model"
	"kr/paasta/monitoring/utils"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
)

type DBConfig struct {
	DbType       string
	UserName     string
	UserPassword string
	Host         string
	Port         string
	DbName       string
}

type Config map[string]string

var TestToken string

type Response struct {
	Token   string
	Content string
	Code    int
}

type PingResponse struct {
	Token string
	Code  int
}

type LoginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var Logger seelog.LoggerInterface

func TestController(t *testing.T) {
	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^")
	fmt.Printf("f : %v\n", t)
	tt = t
	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^")
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")
}

var _ = Describe("Controller BeforeSuite", func() {

	BeforeSuite(func() {

		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>...  1")

		//config, err := readConfig(`../../config.ini`)
		config, err := readConfig(`/home/ubuntu/workspace/user/arom/PaaS-TA-Monitoring/iaas-paasta-monitoring-management/src/kr/paasta/monitoring/config.ini`)
		if err != nil {
			fmt.Printf("read config file error: %s", err)
			fmt.Errorf("read config file error: %s", err)
			os.Exit(0)
		}

		xmlFile, err := ReadXmlConfig(`/home/ubuntu/workspace/user/arom/PaaS-TA-Monitoring/iaas-paasta-monitoring-management/src/kr/paasta/monitoring/log_config.xml`)
		if err != nil {
			log.Println(err)
			os.Exit(-1)
		}
		logger, err := seelog.LoggerFromConfigAsBytes([]byte(xmlFile))

		if err != nil {
			fmt.Println(err)
			return
		}
		iaasModels.MonitLogger = logger
		UseLogger(logger)

		// GMT Time Setting
		timeGap, _ := strconv.Atoi(config["gmt.time.gap"])
		iaasModels.GmtTimeGap = timeGap
		fmt.Println("String - timeGap : ", timeGap)

		// IaaS Database Connection
		iaasConfigDbCon := new(DBConfig)
		iaasConfigDbCon.DbType = config["iaas.monitoring.db.type"]
		iaasConfigDbCon.DbName = config["iaas.monitoring.db.dbname"]
		iaasConfigDbCon.UserName = config["iaas.monitoring.db.username"]
		iaasConfigDbCon.UserPassword = config["iaas.monitoring.db.password"]
		iaasConfigDbCon.Host = config["iaas.monitoring.db.host"]
		iaasConfigDbCon.Port = config["iaas.monitoring.db.port"]
		iaasConnectionString := utils.GetConnectionString(iaasConfigDbCon.Host, iaasConfigDbCon.Port, iaasConfigDbCon.UserName,
			iaasConfigDbCon.UserPassword, iaasConfigDbCon.DbName)

		fmt.Println("String:", iaasConnectionString)
		iaasDbAccessObj, dbErr := gorm.Open(iaasConfigDbCon.DbType, iaasConnectionString+"?charset=utf8&parseTime=true")
		if dbErr != nil {
			fmt.Println("err::", dbErr)
		}

		//iaasDbAccessObj.Debug().AutoMigrate(&iaasModels.AlarmActionHistory{})

		// paaS Database Connection
		paasConfigDbCon := new(DBConfig)
		paasConfigDbCon.DbType = config["paas.monitoring.db.type"]
		paasConfigDbCon.DbName = config["paas.monitoring.db.dbname"]
		paasConfigDbCon.UserName = config["paas.monitoring.db.username"]
		paasConfigDbCon.UserPassword = config["paas.monitoring.db.password"]
		paasConfigDbCon.Host = config["paas.monitoring.db.host"]
		paasConfigDbCon.Port = config["paas.monitoring.db.port"]
		paasConnectionString := utils.GetConnectionString(paasConfigDbCon.Host, paasConfigDbCon.Port, paasConfigDbCon.UserName,
			paasConfigDbCon.UserPassword, paasConfigDbCon.DbName)

		fmt.Println("String:", paasConnectionString)
		paasDbAccessObj, dbErr := gorm.Open(paasConfigDbCon.DbType, paasConnectionString+"?charset=utf8&parseTime=true")
		if dbErr != nil {
			fmt.Println("err::", dbErr)
		}

		paasDbAccessObj.Debug().AutoMigrate(&iaasModels.AlarmActionHistory{})

		// IaaS InfluxDB Info
		iaasUrl, _ := config["iaas.metric.db.url"]
		iaasUserName, _ := config["iaas.metric.db.username"]
		iaasPassword, _ := config["iaas.metric.db.password"]

		iaasInfluxServerClient, _ := client.NewHTTPClient(client.HTTPConfig{
			Addr:     iaasUrl,
			Username: iaasUserName,
			Password: iaasPassword,
		})

		// PaaS InfluxDB Info
		paasUrl, _ := config["paas.metric.db.url"]
		paasUserName, _ := config["paas.metric.db.username"]
		paasPassword, _ := config["paas.metric.db.password"]

		paasInfluxServerClient, _ := client.NewHTTPClient(client.HTTPConfig{
			Addr:     paasUrl,
			Username: paasUserName,
			Password: paasPassword,
		})

		// IaaS ElasticSearch
		iaasElasticUrl, _ := config["iaas.elastic.url"]
		iaasElasticClient, err := elastic.NewClient(
			elastic.SetURL(fmt.Sprintf("http://%s", iaasElasticUrl)),
			elastic.SetSniff(false),
		)

		// PaaS ElasticSearch
		paasElasticUrl, _ := config["paas.elastic.url"]
		paasElasticClient, err := elastic.NewClient(
			elastic.SetURL(fmt.Sprintf("http://%s", paasElasticUrl)),
			elastic.SetSniff(false),
		)

		var openstackProvider iaasModels.OpenstackProvider
		openstackProvider.Region, _ = config["default.region"]
		openstackProvider.Username, _ = config["default.username"]
		openstackProvider.Password, _ = config["default.password"]
		openstackProvider.Domain, _ = config["default.domain"]
		openstackProvider.TenantName, _ = config["default.tenant_name"]
		openstackProvider.AdminTenantId, _ = config["default.project_id"]
		openstackProvider.KeystoneUrl, _ = config["keystone.url"]
		openstackProvider.IdentityEndpoint, _ = config["identity.endpoint"]
		openstackProvider.RabbitmqUser, _ = config["rabbitmq.user"]
		openstackProvider.RabbitmqPass, _ = config["rabbitmq.pass"]
		openstackProvider.RabbitmqTargetNode, _ = config["rabbitmq.target.node"]

		iaasModels.MetricDBName, _ = config["iaas.metric.db.name"]
		iaasModels.NovaUrl, _ = config["nova.target.url"]
		iaasModels.NovaVersion, _ = config["nova.target.version"]
		iaasModels.NeutronUrl, _ = config["neutron.target.url"]
		iaasModels.NeutronVersion, _ = config["neutron.target.version"]
		iaasModels.KeystoneUrl, _ = config["keystone.target.url"]
		iaasModels.KeystoneVersion, _ = config["keystone.target.version"]
		iaasModels.CinderUrl, _ = config["cinder.target.url"]
		iaasModels.CinderVersion, _ = config["cinder.target.version"]
		iaasModels.GlanceUrl, _ = config["glance.target.url"]
		iaasModels.GlanceVersion, _ = config["glance.target.version"]
		iaasModels.DefaultTenantId, _ = config["default.project_id"]
		iaasModels.RabbitMqIp, _ = config["rabbitmq.ip"]
		iaasModels.RabbitMqPort, _ = config["rabbitmq.port"]
		iaasModels.GMTTimeGap, _ = strconv.ParseInt(config["gmt.time.gap"], 10, 64)

		monClient := monascaclient.New()
		monClient.SetBaseURL(config["monasca.url"])
		timeOut, _ := strconv.Atoi(config["monasca.connect.timeout"])
		monClient.SetTimeout(timeOut)

		tls, _ := strconv.ParseBool(config["monasca.secure.tls"])
		monClient.SetInsecure(tls)

		auth := gophercloud.AuthOptions{
			DomainName:       config["default.domain"],
			IdentityEndpoint: config["keystone.url"],
			Username:         config["default.username"],
			Password:         config["default.password"],
			TenantID:         config["default.project_id"],
		}
		iaasModels.TestUserName = auth.Username
		iaasModels.TestPassword = auth.Password
		iaasModels.TestTenantID = auth.TenantID
		iaasModels.TestDomainName = auth.DomainName
		iaasModels.TestIdentityEndpoint = auth.IdentityEndpoint

		// PaaS Database
		bosh_database, _ := config["paas.metric.db.name.bosh"]
		paasta_database, _ := config["paas.metric.db.name.paasta"]
		container_database, _ := config["paas.metric.db.name.container"]

		var databases paasModels.Databases
		databases.BoshDatabase = bosh_database
		databases.PaastaDatabase = paasta_database
		databases.ContainerDatabase = container_database

		// Redis Client
		rdClient := redis.NewClient(&redis.Options{
			Addr:     config["redis.addr"],
			Password: config["redis.password"],
		})

		sysType := config["system.monitoring.type"]

		// BOSH Client Config
		boshConfig := &gogobosh.Config{
			BOSHAddress:       config["bosh.client.api.address"],
			Username:          config["bosh.client.api.username"],
			Password:          config["bosh.client.api.password"],
			HttpClient:        http.DefaultClient,
			SkipSslValidation: true,
		}
		boshClient, err := gogobosh.NewClient(boshConfig)
		if err != nil {
			log.Fatalln("Failed to create connection to the bosh server. err=", err)
		}

		bm.GmtTimeGap = timeGap

		cfConfig := bm.CFConfig{
			Host:           config["paas.monitoring.cf.host"],
			CaasBrokerHost: config["caas.monitoring.broker.host"],
		}

		// Handler
		var handler http.Handler
		handler = NewHandler(openstackProvider, iaasInfluxServerClient, paasInfluxServerClient,
			iaasDbAccessObj, paasDbAccessObj, iaasElasticClient, paasElasticClient, *monClient, auth, databases,
			rdClient, sysType, boshClient, cfConfig)

		server = httptest.NewServer(handler)
		testUrl = server.URL

		// Login Test
		res, err := DoGetPing(testUrl + "/v2/ping")

		var userInfo comModels.UserInfo
		userInfo.Username = "admin"
		userInfo.Password = "1234"

		TestToken = res.Token

		fmt.Println("************************************************************")
		fmt.Println("testUrl : " + testUrl)
		fmt.Println("TestToken : " + TestToken)
		fmt.Println("************************************************************")

		data, _ := json.Marshal(userInfo)
		DoPost(testUrl+"/v2/login", TestToken, strings.NewReader(string(data)))
	})

})

func ReadXmlConfig(filename string) (string, error) {
	xmlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "", err
	}
	return string(xmlFile), nil
}

func UseLogger(newLogger seelog.LoggerInterface) {
	Logger = newLogger
}

func DoGetPing(url string) (*PingResponse, error) {

	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(iaasModels.TEST_TOKEN_NAME, iaasModels.TEST_TOKEN_VALUE)
	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	token := response.Header.Get(iaasModels.TEST_TOKEN_NAME)

	return &PingResponse{Token: string(token), Code: response.StatusCode}, nil
}

func DoGet(url string) (*Response, error) {

	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(iaasModels.CSRF_TOKEN_NAME, TestToken)
	req.Header.Add(iaasModels.TEST_TOKEN_NAME, TestToken)
	req.Header.Add("username", "admin")
	req.Header.Add("password", "1234")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
		fmt.Printf("err : %v\n", err)
		fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
		return nil, err
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &Response{Content: string(contents), Code: response.StatusCode}, nil
}

func DoPost(url, token string, body io.Reader) (*Response, error) {

	client := &http.Client{}

	req, _ := http.NewRequest("POST", url, body)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add(iaasModels.CSRF_TOKEN_NAME, token)
	req.Header.Add(iaasModels.TEST_TOKEN_NAME, token)

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return &Response{Content: string(contents), Code: response.StatusCode}, nil
}

func DoLogout(url, testToken string) (*Response, error) {

	client := &http.Client{}

	token, _ := utils.GenerateRandomString(32)

	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add(iaasModels.TEST_TOKEN_NAME, token)
	req.Header.Add(iaasModels.CSRF_TOKEN_NAME, token)

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return &Response{Content: string(contents), Code: response.StatusCode}, nil
}

func DoUpdate(url, token string, body io.Reader) (*Response, error) {

	client := &http.Client{}

	req, _ := http.NewRequest("PUT", url, body)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add(iaasModels.TEST_TOKEN_NAME, token)
	req.Header.Add(iaasModels.CSRF_TOKEN_NAME, token)

	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return &Response{Content: string(contents), Code: response.StatusCode}, nil
}

func DoPatch(url, token string, body io.Reader) (*Response, error) {

	client := &http.Client{}

	req, _ := http.NewRequest("PATCH", url, body)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add(iaasModels.TEST_TOKEN_NAME, token)
	req.Header.Add(iaasModels.CSRF_TOKEN_NAME, token)

	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return &Response{Content: string(contents), Code: response.StatusCode}, nil
}

func DoDelete(url, token string, body io.Reader) (*Response, error) {

	client := &http.Client{}

	req, _ := http.NewRequest("DELETE", url, body)
	req.Header.Add("Accept", "application/json")
	req.Header.Add(iaasModels.TEST_TOKEN_NAME, token)
	req.Header.Add(iaasModels.CSRF_TOKEN_NAME, token)

	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return &Response{Content: string(contents), Code: response.StatusCode}, nil
}

func DoDetail(url, token string, body io.Reader) (*Response, error) {

	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, body)

	req.Header.Add("Accept", "application/json")
	req.Header.Add(iaasModels.TEST_TOKEN_NAME, token)
	req.Header.Add(iaasModels.CSRF_TOKEN_NAME, token)

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &Response{Content: string(contents), Code: response.StatusCode}, nil
}

var (
	server  *httptest.Server
	testUrl string
	t       *testing.T
	tt      *testing.T
)

func readConfig(filename string) (Config, error) {
	fmt.Println("filename : " + filename)

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

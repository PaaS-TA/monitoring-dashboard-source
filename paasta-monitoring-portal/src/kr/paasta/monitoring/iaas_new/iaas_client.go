package iaas_new

import (
	"crypto/tls"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gophercloud/gophercloud"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/jinzhu/gorm"
	commonModel "kr/paasta/monitoring/common/model"
	"kr/paasta/monitoring/iaas_new/model"
	iaasModel "kr/paasta/monitoring/iaas_new/model"
	"kr/paasta/monitoring/utils"
	"net"
	"net/http"
	"strconv"
	"time"
)

func GetIaasClients(config map[string]string) (iaasDbAccessObj *gorm.DB, iaaSInfluxServerClient client.Client, iaasElasticClient *elasticsearch.Client, openstackProvider iaasModel.OpenstackProvider, auth gophercloud.AuthOptions, err error) {

	// Mysql
	iaasConfigDbCon := new(commonModel.DBConfig)
	iaasConfigDbCon.DbType = config["iaas.monitoring.db.type"]
	iaasConfigDbCon.DbName = config["iaas.monitoring.db.dbname"]
	iaasConfigDbCon.UserName = config["iaas.monitoring.db.username"]
	iaasConfigDbCon.UserPassword = config["iaas.monitoring.db.password"]
	iaasConfigDbCon.Host = config["iaas.monitoring.db.host"]
	iaasConfigDbCon.Port = config["iaas.monitoring.db.port"]

	iaasConnectionString := utils.GetConnectionString(iaasConfigDbCon.Host, iaasConfigDbCon.Port, iaasConfigDbCon.UserName, iaasConfigDbCon.UserPassword, iaasConfigDbCon.DbName)
	fmt.Println("String:", iaasConnectionString)
	iaasDbAccessObj, err = gorm.Open(iaasConfigDbCon.DbType, iaasConnectionString+"?charset=utf8&parseTime=true")

	// 2021.09.06 - 이거 왜 있는지?
	//Alarm 처리 내역 정보 Table 생성
	//iaasDbAccessObj.Debug().AutoMigrate(&model.AlarmActionHistory{})

	// InfluxDB
	iaasUrl, _ := config["iaas.metric.db.url"]
	iaasUserName, _ := config["iaas.metric.db.username"]
	iaasPassword, _ := config["iaas.metric.db.password"]

	iaaSInfluxServerClient, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     iaasUrl,
		Username: iaasUserName,
		Password: iaasPassword,
		InsecureSkipVerify: true,
	})

	elasticsearchUsername, _ := config["paas.elasticsearch.username"]
	elasticsearchPassword, _ := config["paas.elasticsearch.password"]
	elasticsearchUrl, _ := config["paas.elasticsearch.url"]
	elasticsearchHttpsEnabled, _ := strconv.ParseBool(config["paas.elasticsearch.https_enabled"])

	cfg := elasticsearch.Config{
		Username: elasticsearchUsername,
		Password: elasticsearchPassword,
		Addresses: []string{
			elasticsearchUrl,
		},
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MaxVersion:         tls.VersionTLS11,
				InsecureSkipVerify: elasticsearchHttpsEnabled,
			},
		},
	}
	iaasElasticClient, err = elasticsearch.NewClient(cfg)
	fmt.Println("iaasElasticClient::", iaasElasticClient)
	fmt.Println("err::", err)

	// ElasticSearch
	/*iaasElasticUrl, _ := config["iaas.elastic.url"]
	iaasElasticClient, err = elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("http://%s", iaasElasticUrl)),
		elastic.SetSniff(false),
	)*/

	// Openstack Info
	openstackProvider.Region, _ = config["default.region"]
	openstackProvider.Username, _ = config["default.username"]
	openstackProvider.Password, _ = config["default.password"]
	openstackProvider.Domain, _ = config["default.domain"]
	openstackProvider.TenantName, _ = config["default.tenant_name"]
	openstackProvider.AdminTenantId, _ = config["default.tenant_id"]
	openstackProvider.KeystoneUrl, _ = config["keystone.url"]
	openstackProvider.IdentityEndpoint, _ = config["identity.endpoint"]
	openstackProvider.RabbitmqUser, _ = config["rabbitmq.user"]
	openstackProvider.RabbitmqPass, _ = config["rabbitmq.pass"]
	openstackProvider.RabbitmqTargetNode, _ = config["rabbitmq.target.node"]

	model.MetricDBName, _ = config["iaas.metric.db.name"]
	model.NovaUrl, _ = config["nova.target.url"]
	model.NovaVersion, _ = config["nova.target.version"]
	model.NeutronUrl, _ = config["neutron.target.url"]
	model.NeutronVersion, _ = config["neutron.target.version"]
	model.KeystoneUrl, _ = config["keystone.target.url"]
	model.KeystoneVersion, _ = config["keystone.target.version"]
	model.CinderUrl, _ = config["cinder.target.url"]
	model.CinderVersion, _ = config["cinder.target.version"]
	model.GlanceUrl, _ = config["glance.target.url"]
	model.GlanceVersion, _ = config["glance.target.version"]
	model.DefaultTenantId, _ = config["default.tenant_id"]
	model.RabbitMqIp, _ = config["rabbitmq.ip"]
	model.RabbitMqPort, _ = config["rabbitmq.port"]
	model.GMTTimeGap, _ = strconv.ParseInt(config["gmt.time.gap"], 10, 64)


	auth = gophercloud.AuthOptions{
		DomainName:       config["default.domain"],
		IdentityEndpoint: config["keystone.url"],
		Username:         config["default.username"],
		Password:         config["default.password"],
		TenantID:         config["default.tenant_id"],
	}

	return
}

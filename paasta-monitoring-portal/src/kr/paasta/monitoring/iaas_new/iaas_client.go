package iaas_new

import (
    "fmt"
    "github.com/gophercloud/gophercloud"
    influx "github.com/influxdata/influxdb1-client/v2"
    "github.com/jinzhu/gorm"
    commonModel "kr/paasta/monitoring/common/model"
    iaasModel "kr/paasta/monitoring/iaas_new/model"
    "kr/paasta/monitoring/utils"
    "strconv"
)

type IaasClient struct {
    ConnectionPool *gorm.DB
    InfluxClient   influx.Client
    Provider       iaasModel.OpenstackProvider
    AuthOpts       gophercloud.AuthOptions
}

func GetIaasClients(config map[string]string) (client IaasClient, err error) {

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
    iaasDbAccessObj, _ := gorm.Open(iaasConfigDbCon.DbType, iaasConnectionString+"?charset=utf8&parseTime=true")

    // 2021.09.06 - 이거 왜 있는지?
    //Alarm 처리 내역 정보 Table 생성
    //iaasDbAccessObj.Debug().AutoMigrate(&model.AlarmActionHistory{})

    // InfluxDB
    iaasUrl, _ := config["iaas.metric.db.url"]
    iaasUserName, _ := config["iaas.metric.db.username"]
    iaasPassword, _ := config["iaas.metric.db.password"]

    iaaSInfluxServerClient, _ := influx.NewHTTPClient(influx.HTTPConfig{
        Addr:               iaasUrl,
        Username:           iaasUserName,
        Password:           iaasPassword,
        InsecureSkipVerify: true,
    })

    // Openstack 정보
    openstackProvider := iaasModel.OpenstackProvider{}
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

    iaasModel.MetricDBName, _ = config["iaas.metric.db.name"]
    iaasModel.NovaUrl, _ = config["nova.target.url"]
    iaasModel.NovaVersion, _ = config["nova.target.version"]
    iaasModel.NeutronUrl, _ = config["neutron.target.url"]
    iaasModel.NeutronVersion, _ = config["neutron.target.version"]
    iaasModel.KeystoneUrl, _ = config["keystone.target.url"]
    iaasModel.KeystoneVersion, _ = config["keystone.target.version"]
    iaasModel.CinderUrl, _ = config["cinder.target.url"]
    iaasModel.CinderVersion, _ = config["cinder.target.version"]
    iaasModel.GlanceUrl, _ = config["glance.target.url"]
    iaasModel.GlanceVersion, _ = config["glance.target.version"]
    iaasModel.DefaultTenantId, _ = config["default.tenant_id"]
    iaasModel.RabbitMqIp, _ = config["rabbitmq.ip"]
    iaasModel.RabbitMqPort, _ = config["rabbitmq.port"]
    iaasModel.GMTTimeGap, _ = strconv.ParseInt(config["gmt.time.gap"], 10, 64)

    auth := gophercloud.AuthOptions{
        DomainName: config["default.domain"],
        //IdentityEndpoint: config["keystone.url"],
        IdentityEndpoint: config["identity.endpoint"],
        Username:         config["default.username"],
        Password:         config["default.password"],
        TenantID:         config["default.tenant_id"],
    }

    client.ConnectionPool = iaasDbAccessObj
    client.InfluxClient = iaaSInfluxServerClient
    client.Provider = openstackProvider
    client.AuthOpts = auth

    return
}

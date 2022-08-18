package main

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	Connections "paasta-monitoring-api/connections"
	_ "paasta-monitoring-api/docs"
	Routers "paasta-monitoring-api/routers"
)

//Execution starts from main function
// @title            Monitoring Dashboard API
// @version          5.8.0
// @host             localhost:8395
// @tag.name         Common
// @tag.description  Common Module API (Alarm & Log)
// @tag.name         AP
// @tag.description  Application Platform API (Based on Cloud Foundry)
// @tag.name         CP
// @tag.description  Container Platform API (Based on Kubernetes/Prometheus)
// @tag.name         SaaS
// @tag.description  Application Performance Monitoring API (Based on Pinpoint)
// @tag.name         IaaS
// @tag.description  Infrastructure Monitoring API (Based on Openstack/Zabbix) - Only can use it, when you use IaaS option
func main() {
	logger := logrus.New()
	logger.SetFormatter(&nested.Formatter{
		CallerFirst: true,
	})
	logger.SetReportCaller(true)

	// .env 파일 로드
	err := godotenv.Load(".env")
	if err != nil {
		logger.Panic("Error loading .env file")
	}

	// connection 설정 (DB & API etc..)
	c := Connections.SetupConnection(logger)

	// Router 설정
	r := Routers.SetupRouter(c)

	webPort := os.Getenv("web_port")

	// For run on requested port
	if len(os.Args) > 1 {
		reqPort := os.Args[1]
		if reqPort != "" {
			webPort = reqPort
		}
	}

	if webPort == "" {
		webPort = "8080" //localhost
	}

	logger.Info(r.Start(":" + webPort))
}

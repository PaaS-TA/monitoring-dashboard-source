package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	Connections "paasta-monitoring-api/connections"
	_ "paasta-monitoring-api/docs"
	Routers "paasta-monitoring-api/routers"
)

//Execution starts from main function
// @title            Monitoring Dashboard API
// @version          5.8.0
// @host             localhost:8395
// @tag.name         AP
// @tag.description  Application Platform API (Cloud Foundry)
// @tag.name         CP
// @tag.description  Container Platform API (Kubernetes)
// @tag.name         SaaS
// @tag.description  APM (Pinpoint)
// @tag.name         IaaS
// @tag.description  Openstack(Zabbix)
func main() {
	// .env 파일 로드
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// connection 설정 (DB & API etc..)
	c := Connections.SetupConnection()

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

	r.Logger.Fatal(r.Start(":" + webPort))
}
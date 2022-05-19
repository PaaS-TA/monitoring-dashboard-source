package main

import (
	Connections "GoEchoProject/connections"
	_ "GoEchoProject/docs"
	Routers "GoEchoProject/routers"
	"github.com/joho/godotenv"
	"log"
	"os"
)

//Execution starts from main function
// @title    Monitoring Dashboard API
// @version  5.8.0
// @host     localhost:8395
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

package main

import (
	Connections "GoEchoProject/connections"
	Routers "GoEchoProject/routers"
	"os"
)

//Execution starts from main function
func main() {

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

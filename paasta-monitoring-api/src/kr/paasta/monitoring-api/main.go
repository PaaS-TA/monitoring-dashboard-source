package main

import (
	Connections "GoEchoProject/connections"
	Routers "GoEchoProject/routers"
	"fmt"
	"os"
)

//Execution starts from main function
func main() {

	// connection 설정 (DB & API etc..)
	c, err := Connections.SetupConnection()
	if err != nil {
		fmt.Println(err)
	}

	// Router 설정
	r := Routers.SetupRouter(c)

	port := os.Getenv("port")

	// For run on requested port
	if len(os.Args) > 1 {
		reqPort := os.Args[1]
		if reqPort != "" {
			port = reqPort
		}
	}

	if port == "" {
		port = "8080" //localhost
	}
	type Job interface {
		Run()
	}

	r.Logger.Fatal(r.Start(":" + port))
}

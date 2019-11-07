package main

import (
	"fmt"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Main Test", func() {

	Describe("Main Contents", func() {

		Context("Main Func", func() {
			var testConfig Config
			It("Main ReadConfig", func() {
				fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> 2")
				testConfig, _ = ReadConfig(`config.ini`)
				fmt.Println(testConfig)
			})

			It("Main ReadXmlConfig", func() {
				fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> 3")
				res, _ := ReadXmlConfig(`log_config.xml`)
				fmt.Println(res)
			})

			It("Main getIaasClients", func() {
				fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> 4")
				iaasDbAccessObj, iaaSInfluxServerClient, iaasElasticClient, openstackProvider, monClient, auth, _ := getIaasClients(testConfig)
				fmt.Println(iaasDbAccessObj, iaaSInfluxServerClient, iaasElasticClient, openstackProvider, monClient, auth)
			})

			It("Main getPaasClients", func() {
				fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> 5")
				paaSInfluxServerClient, paasElasticClient, databases, cfProvider, boshClient, _ := getPaasClients(testConfig)
				fmt.Println(paaSInfluxServerClient, paasElasticClient, databases, cfProvider, boshClient)
			})
		})
	})
})

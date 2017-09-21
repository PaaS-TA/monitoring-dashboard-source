package services_test

import (
	. "github.com/onsi/ginkgo"
	_ "github.com/go-sql-driver/mysql"
	"kr/paasta/monitoring/monit-batch/services"
)

//#Top-level 단위테스트 묶음.
var _ = Describe("Container Service Test", func() {

	It("Container Alarm Collector", func() {
		config, _ := ReadConfig(`../config_test.ini`)
		influxCon, configDbCon, portalDbCon, boshCon, mailConfig, thresholdConfig, config := GetObject(config)

		backendService := services.NewBackendServices(-9, influxCon, configDbCon, portalDbCon, boshCon, mailConfig, thresholdConfig, config)

		services.ContainerAlarmCollect(backendService)
		/*Eventually(func()  {
			services.BoshAlarmCollect(backendService)
		})*/
	})
})

package services_test


import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	_ "github.com/go-sql-driver/mysql"
	"kr/paasta/monitoring/monit-batch/services"
	"kr/paasta/monitoring/monit-batch/models"
	"kr/paasta/monitoring/monit-batch/dao"
)
//#Top-level 단위테스트 묶음.
var _ = Describe("Bosh Data Sync Test", func() {

	It("Bosh Data Sync", func() {
		config, _ := ReadConfig(`../config_test.ini`)
		influxCon, configDbCon, portalDbCon, boshCon, mailConfig, config := GetObject(config)

		backendService := services.NewBackendServices(-9, influxCon, configDbCon, portalDbCon, boshCon,  mailConfig, config)

		//Vm, Zone정보 초기화
		backendService.MonitoringDbClient.Where("").Delete(models.Vm{})
		backendService.MonitoringDbClient.Where("").Delete(models.Zone{})

		services.CreteUpdateBoshVms(backendService, *backendService.BoshConfig, backendService.MonitoringDbClient)
		vmList, _ := dao.GetBoshVmsDao(backendService.BoshClient).GetJobInfoList(backendService.MonitoringDbClient)
		if len(vmList) > 1 {

			dao.GetBoshVmsDao(backendService.BoshClient).DeleteVmInfo(vmList[0], backendService.MonitoringDbClient )
			vmList[1].Ip = "xxxx"
			dao.GetBoshVmsDao(backendService.BoshClient).UpdateVmData(vmList[1], backendService.MonitoringDbClient )

			testZone := models.Zone{Id: 9999, Name: "z3"}
			testVm := models.Vm{ZoneId: 9999, Name: "testVM", Ip: "127.0.0.1", VmType: "vms"}
			backendService.MonitoringDbClient.FirstOrCreate(&testZone)
			backendService.MonitoringDbClient.FirstOrCreate(&testVm)

			services.CreteUpdateBoshVms(backendService, *backendService.BoshConfig, backendService.MonitoringDbClient)
		}

		Eventually(func() bool {
			zoneList, _ := dao.GetBoshVmsDao(backendService.BoshClient).GetZoneInfosList(backendService.MonitoringDbClient)
			vmList, _ := dao.GetBoshVmsDao(backendService.BoshClient).GetJobInfoList(backendService.MonitoringDbClient)
			if len(zoneList) > 0 && len(vmList) > 0{
				return true
			}
			return false
		}).Should(BeTrue())

	})

})


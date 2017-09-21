package services_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	_ "github.com/go-sql-driver/mysql"
	"kr/paasta/monitoring/monit-batch/services"
	"kr/paasta/monitoring/monit-batch/models"
	"kr/paasta/monitoring/monit-batch/models/base"
	"kr/paasta/monitoring/monit-batch/dao"
)

//#Top-level 단위테스트 묶음.
var _ = Describe("PaasTa Service Test", func() {

	It("PaasTa Alarm Collector ", func() {
		config, _ := ReadConfig(`../config_test.ini`)
		influxCon, configDbCon, portalDbCon, boshCon, mailConfig, thresholdConfig, config := GetObject(config)

		backendService := services.NewBackendServices(-9, influxCon, configDbCon, portalDbCon, boshCon,  mailConfig, thresholdConfig, config)
		//Fail Data 생성
		testZone := models.Zone{Id: 1111, Name: "z9"}
		testVm := models.Vm{Id: 9999, ZoneId: 1111, Name: "testVM", Ip: "127.0.0.1", VmType: "vms"}
		backendService.MonitoringDbClient.FirstOrCreate(&testZone)
		backendService.MonitoringDbClient.FirstOrCreate(&testVm)

		services.PaasTaAlarmCollect(backendService)

		services.PaasTaAlarmCollect(backendService)

		Eventually (func()  bool{
			//Alarm이 정상 발생 되었는지 확인 조회
			//Fail Data생성
			failData := base.Alarm{OriginType: base.ORIGIN_TYPE_PAASTA, OriginId: 9999, AlarmType: base.ALARM_LEVEL_FAIL, Level: base.ALARM_LEVEL_FAIL}
			failExist, _ := dao.GetCommonDao(backendService.Influxclient).GetIsNotExistAlarm(failData, backendService.MonitoringDbClient)

			//PaasTa VM 중 CloudControlle Memory가 Warning 처리 되도록 임계치 수정
			alarmData := base.Alarm{OriginType: base.ORIGIN_TYPE_PAASTA}
			alarmExist, _ := dao.GetCommonDao(backendService.Influxclient).GetIsNotExistAlarmCheck(alarmData, backendService.MonitoringDbClient)

			if failExist == true && alarmExist == true{
				return true
			}
			return false
		}).ShouldNot(BeTrue())


		//발생한 Alarm Data 삭제
		backendService.MonitoringDbClient.Where("").Delete(base.Alarm{})

	})

})


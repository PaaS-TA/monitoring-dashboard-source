package services

import (
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring/monit-batch/models"
)

func CreateTable(dbClient *gorm.DB){

	dbClient.Debug().AutoMigrate(&models.Zone{}, &models.Vm{})
	dbClient.Debug().AutoMigrate(&models.Alarm{}, &models.AlarmAction{}, &models.AlarmPolicy{}, &models.AlarmTarget{})
}

/*func CreateTablePortal(dbClient *gorm.DB){
	dbClient.Debug().AutoMigrate(&models.AutoScaleConfig{})
}*/


func CreateAlarmPolicyInitialData(dbClient *gorm.DB){

	paasTaCpuData := models.AlarmPolicy{Id:1, OriginType: "pas", AlarmType: "cpu", WarningThreshold: 85, CriticalThreshold: 90, RepeatTime: 10 , Comment: "Initial Data"}
	paasTaMemData := models.AlarmPolicy{Id:2, OriginType: "pas", AlarmType: "memory", WarningThreshold: 85, CriticalThreshold: 90, RepeatTime: 10 , Comment: "Initial Data"}
	paasTaDiskData := models.AlarmPolicy{Id:3, OriginType: "pas", AlarmType: "disk", WarningThreshold: 85, CriticalThreshold: 90, RepeatTime: 1 , Comment: "Initial Data"}

	boshCpuData := models.AlarmPolicy{Id:4, OriginType: "bos", AlarmType: "cpu", WarningThreshold: 85, CriticalThreshold: 90, RepeatTime: 10 , Comment: "Initial Data"}
	boshMemData := models.AlarmPolicy{Id:5, OriginType: "bos", AlarmType: "memory", WarningThreshold:85, CriticalThreshold: 90, RepeatTime: 10 , Comment: "Initial Data"}
	boshDiskData := models.AlarmPolicy{Id:6, OriginType: "bos", AlarmType: "disk", WarningThreshold: 85, CriticalThreshold: 90, RepeatTime: 10 , Comment: "Initial Data"}

	appCpuData := models.AlarmPolicy{Id:7, OriginType: "con", AlarmType: "cpu", WarningThreshold: 85, CriticalThreshold: 90, RepeatTime: 10 , Comment: "Initial Data"}
	appMemData := models.AlarmPolicy{Id:8, OriginType: "con", AlarmType: "memory", WarningThreshold: 85, CriticalThreshold: 90, RepeatTime: 10 , Comment: "Initial Data"}
	appDiskData := models.AlarmPolicy{Id:9, OriginType: "con", AlarmType: "disk", WarningThreshold: 85, CriticalThreshold: 90, RepeatTime: 10 , Comment: "Initial Data"}

	alarmTagetBosh := models.AlarmTarget{Id:1, OriginType: "bos", MailAddress: "xxxx@nate.com" }
	alarmTagetPassTa := models.AlarmTarget{Id:2, OriginType: "pas", MailAddress: "xxxx@nate.com" }
	alarmTagetContainer := models.AlarmTarget{Id:3, OriginType: "con", MailAddress: "xxxx@nate.com" }


	dbClient.FirstOrCreate(&paasTaCpuData)
	dbClient.FirstOrCreate(&paasTaMemData)
	dbClient.FirstOrCreate(&paasTaDiskData)

	dbClient.FirstOrCreate(&boshCpuData)
	dbClient.FirstOrCreate(&boshMemData)
	dbClient.FirstOrCreate(&boshDiskData)

	dbClient.FirstOrCreate(&appCpuData)
	dbClient.FirstOrCreate(&appMemData)
	dbClient.FirstOrCreate(&appDiskData)

	dbClient.FirstOrCreate(&alarmTagetBosh)
	dbClient.FirstOrCreate(&alarmTagetPassTa)
	dbClient.FirstOrCreate(&alarmTagetContainer)

}

func CreatePortalInitialData(dbClient *gorm.DB){

	autoScaleData1 := models.AutoScaleConfig{No:1, Guid:"b7a14c50-4108-4df1-bb1f-c6c5f652d9e8", Org: "org", Space: "space", App: "spring-music", InstanceMaxCnt: 20, InstanceMinCnt: 2,
	  CpuThresholdMaxPer: 80, CpuThresholdMinPer: 20, MemoryMaxSize: 80, MemoryMinSize:20, CheckTimeSec: 30, AutoDecreaseYn: "Y", AutoIncreaseYn: "Y"}

	autoScaleData2 := models.AutoScaleConfig{No:2, Guid:"00b3b012-c6af-49d4-8849-d1b90d53c93f", Org: "org", Space: "space", App: "spring-music-2", InstanceMaxCnt: 20, InstanceMinCnt: 2,
		CpuThresholdMaxPer: 80, CpuThresholdMinPer: 20, MemoryMaxSize: 80, MemoryMinSize:20, CheckTimeSec: 30, AutoDecreaseYn: "Y", AutoIncreaseYn: "Y"}

	dbClient.FirstOrCreate(&autoScaleData1)
	dbClient.FirstOrCreate(&autoScaleData2)
}

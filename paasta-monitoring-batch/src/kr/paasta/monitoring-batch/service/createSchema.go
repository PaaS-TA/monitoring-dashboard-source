package service

import (
	"github.com/jinzhu/gorm"
	"monitoring-batch/model"
)

func CreateTable(dbClient *gorm.DB){

	dbClient.Debug().AutoMigrate(&model.Zone{}, &model.Vm{})
	dbClient.Debug().AutoMigrate(&model.Alarm{}, &model.AlarmAction{}, &model.AlarmPolicy{}, &model.AlarmTarget{},
								&model.MemberInfo{}, &model.AlarmSns{}, &model.AlarmSnsTarget{},
								&model.AppAutoScalingPolicy{}, &model.AppAlarmPolicy{}, &model.AppAlarmHistory{})
}

/*func CreateTablePortal(dbClient *gorm.DB){
	dbClient.Debug().AutoMigrate(&models.AutoScaleConfig{})
}*/


func CreateAlarmPolicyInitialData(dbClient *gorm.DB) {

	paasTaCpuData 	:= model.AlarmPolicy{Id:1, OriginType: "pas", AlarmType: "cpu", WarningThreshold: 85, CriticalThreshold: 90, RepeatTime: 10 , MeasureTime: 600 , Comment: "Initial Data"}
	paasTaMemData 	:= model.AlarmPolicy{Id:2, OriginType: "pas", AlarmType: "memory", WarningThreshold: 85, CriticalThreshold: 90, RepeatTime: 10 , MeasureTime: 600 , Comment: "Initial Data"}
	paasTaDiskData 	:= model.AlarmPolicy{Id:3, OriginType: "pas", AlarmType: "disk", WarningThreshold: 85, CriticalThreshold: 90, RepeatTime: 10 , MeasureTime: 600 , Comment: "Initial Data"}

	boshCpuData 	:= model.AlarmPolicy{Id:4, OriginType: "bos", AlarmType: "cpu", WarningThreshold: 85, CriticalThreshold: 90, RepeatTime: 10 , MeasureTime: 600 , Comment: "Initial Data"}
	boshMemData 	:= model.AlarmPolicy{Id:5, OriginType: "bos", AlarmType: "memory", WarningThreshold:85, CriticalThreshold: 90, RepeatTime: 10 , MeasureTime: 600 , Comment: "Initial Data"}
	boshDiskData 	:= model.AlarmPolicy{Id:6, OriginType: "bos", AlarmType: "disk", WarningThreshold: 85, CriticalThreshold: 90, RepeatTime: 10 , MeasureTime: 600 , Comment: "Initial Data"}

	appCpuData 		:= model.AlarmPolicy{Id:7, OriginType: "con", AlarmType: "cpu", WarningThreshold: 85, CriticalThreshold: 90, RepeatTime: 10 , MeasureTime: 600 , Comment: "Initial Data"}
	appMemData 		:= model.AlarmPolicy{Id:8, OriginType: "con", AlarmType: "memory", WarningThreshold: 85, CriticalThreshold: 90, RepeatTime: 10 , MeasureTime: 600 , Comment: "Initial Data"}
	appDiskData 	:= model.AlarmPolicy{Id:9, OriginType: "con", AlarmType: "disk", WarningThreshold: 85, CriticalThreshold: 90, RepeatTime: 10 , MeasureTime: 600 , Comment: "Initial Data"}

	alarmTagetBosh      := model.AlarmTarget{Id:1, OriginType: "bos", MailAddress: "adminUser@gmail.com", MailSendYn: "Y" }
	alarmTagetPassTa    := model.AlarmTarget{Id:2, OriginType: "pas", MailAddress: "adminUser@gmail.com", MailSendYn: "Y" }
	alarmTagetContainer := model.AlarmTarget{Id:3, OriginType: "con", MailAddress: "adminUser@gmail.com", MailSendYn: "Y" }


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

	autoScaleData1 := model.AutoScaleConfig{No:1, Guid:"b7a14c50-4108-4df1-bb1f-c6c5f652d9e8", Org: "org", Space: "space", App: "spring-music", InstanceMaxCnt: 20, InstanceMinCnt: 2,
		CpuThresholdMaxPer: 80, CpuThresholdMinPer: 20, MemoryMaxSize: 80, MemoryMinSize:20, CheckTimeSec: 30, AutoDecreaseYn: "Y", AutoIncreaseYn: "Y"}

	autoScaleData2 := model.AutoScaleConfig{No:2, Guid:"00b3b012-c6af-49d4-8849-d1b90d53c93f", Org: "org", Space: "space", App: "spring-music-2", InstanceMaxCnt: 20, InstanceMinCnt: 2,
		CpuThresholdMaxPer: 80, CpuThresholdMinPer: 20, MemoryMaxSize: 80, MemoryMinSize:20, CheckTimeSec: 30, AutoDecreaseYn: "Y", AutoIncreaseYn: "Y"}

	dbClient.FirstOrCreate(&autoScaleData1)
	dbClient.FirstOrCreate(&autoScaleData2)
}

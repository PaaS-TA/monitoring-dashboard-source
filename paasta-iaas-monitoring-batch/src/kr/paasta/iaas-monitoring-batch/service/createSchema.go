package service

import (
	"github.com/jinzhu/gorm"
	"iaas-monitoring-batch/model"
)

func CreateTable(dbClient *gorm.DB) {

	dbClient.Debug().AutoMigrate(&model.Zone{}, &model.Vm{})
	dbClient.Debug().AutoMigrate(&model.Alarm{}, &model.AlarmAction{}, &model.AlarmPolicy{}, &model.AlarmTarget{},
								&model.MemberInfo{}, &model.AlarmSns{}, &model.AlarmSnsTarget{},
								&model.AppAutoScalingPolicy{}, &model.AppAlarmPolicy{}, &model.AppAlarmHistory{})
}


func CreateAlarmPolicyInitialData(dbClient *gorm.DB) {
	cpuAlarmPolicy    := model.AlarmPolicy{Id: 1, OriginType: "ias", AlarmType: "cpu", WarningThreshold: 85, CriticalThreshold: 90, RepeatTime: 10 , MeasureTime: 600 , Comment: "Initial Data"}
	memoryAlarmPolicy := model.AlarmPolicy{Id: 2, OriginType: "ias", AlarmType: "memory", WarningThreshold: 85, CriticalThreshold: 90, RepeatTime: 10 , MeasureTime: 600 , Comment: "Initial Data"}
	diskAlarmPolicy   := model.AlarmPolicy{Id: 3, OriginType: "ias", AlarmType: "disk", WarningThreshold: 85, CriticalThreshold: 90, RepeatTime: 10 , MeasureTime: 600 , Comment: "Initial Data"}

	alarmTagetBosh    := model.AlarmTarget{Id:4, OriginType: "ias", MailAddress: "adminUser@gmail.com", MailSendYn: "Y" }

	dbClient.FirstOrCreate(&cpuAlarmPolicy)
	dbClient.FirstOrCreate(&memoryAlarmPolicy)
	dbClient.FirstOrCreate(&diskAlarmPolicy)

	dbClient.FirstOrCreate(&alarmTagetBosh)
}

func CreatePortalInitialData(dbClient *gorm.DB) {
	autoScaleData1 := model.AutoScaleConfig{No:1, Guid:"b7a14c50-4108-4df1-bb1f-c6c5f652d9e8", Org: "org", Space: "space", App: "spring-music", InstanceMaxCnt: 20, InstanceMinCnt: 2,
		CpuThresholdMaxPer: 80, CpuThresholdMinPer: 20, MemoryMaxSize: 80, MemoryMinSize:20, CheckTimeSec: 30, AutoDecreaseYn: "Y", AutoIncreaseYn: "Y"}

	autoScaleData2 := model.AutoScaleConfig{No:2, Guid:"00b3b012-c6af-49d4-8849-d1b90d53c93f", Org: "org", Space: "space", App: "spring-music-2", InstanceMaxCnt: 20, InstanceMinCnt: 2,
		CpuThresholdMaxPer: 80, CpuThresholdMinPer: 20, MemoryMaxSize: 80, MemoryMinSize:20, CheckTimeSec: 30, AutoDecreaseYn: "Y", AutoIncreaseYn: "Y"}

	dbClient.FirstOrCreate(&autoScaleData1)
	dbClient.FirstOrCreate(&autoScaleData2)
}
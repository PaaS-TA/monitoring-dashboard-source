package service

import (
	"fmt"
	"sync"
	"time"
	client "github.com/influxdata/influxdb/client/v2"
	"monitoring-batch/dao"
	"monitoring-batch/util"
	"monitoring-batch/alarm"
	mod "monitoring-batch/model"
	cb "monitoring-batch/model/base"
)

func PaasTaAlarmCollect(f *BackendServices) {

	var alarmConfig mod.AlarmConfig
	alarmConfig.SmtpHost = f.MailConfig.SmtpHost
	alarmConfig.Port = f.MailConfig.Port
	alarmConfig.MailSender = f.MailConfig.MailSender
	alarmConfig.SenderPwd = f.MailConfig.SenderPwd
	alarmConfig.MailResource = f.MailConfig.ResouceUrl
	alarmConfig.MailReceiver = f.MailConfig.MailReceiver
	alarmConfig.AlarmSend = f.MailConfig.AlarmSend
	alarmConfig.MailTlsSend = f.MailConfig.MailTlsSend

	//임계치 종류에는 6가지 이다.
	// CPU    - Critical/warning
	// Memory - Critical/warning
	// Disk   - Critical/warning
	var thresholdList []cb.AlarmThreshold
	alarmPolicy, _ := dao.GetPaasTaAlarmDao(f.Influxclient).GetPaastaAlarmPolicy(f.MonitoringDbClient)

	var measureTimeList []mod.AlarmItemMeasureTime

	for _, data := range alarmPolicy {

		var thresholdCritical cb.AlarmThreshold
		var warningThreshold cb.AlarmThreshold

		thresholdCritical.OriginType = cb.ORIGIN_TYPE_PAASTA
		thresholdCritical.AlarmType = data.AlarmType
		thresholdCritical.Level = cb.ALARM_LEVEL_CRITICAL
		thresholdCritical.Threshold = float64(data.CriticalThreshold)
		thresholdCritical.RepeatTime = data.RepeatTime

		warningThreshold.OriginType = cb.ORIGIN_TYPE_PAASTA
		warningThreshold.AlarmType = data.AlarmType
		warningThreshold.Level = cb.ALARM_LEVEL_WARNING
		warningThreshold.Threshold = float64(data.WarningThreshold)
		warningThreshold.RepeatTime = data.RepeatTime

		thresholdList = append(thresholdList, thresholdCritical)
		thresholdList = append(thresholdList, warningThreshold)

		measureTimeList = append(measureTimeList, mod.AlarmItemMeasureTime{data.AlarmType, data.MeasureTime})
	}

	paasTaList, dbErr := dao.GetPaasTaAlarmDao(f.Influxclient).GetPaaSTaList(f.MonitoringDbClient)

	//fmt.Println("=paasTaList=======>>>>", paasTaList)
	if dbErr != nil {
		fmt.Errorf("Error:", dbErr)
	}

	usageResponse := make([]mod.PaasTaResponse, len(paasTaList))

	var cpuUsageList []map[string]interface{}
	var memoryUsageList []map[string]interface{}
	var diskUsageList []map[string]interface{}
	var diskRootUsageList []map[string]interface{}

	//response := make([]mod.BoshResponse, len(paasTaList))

	var wg sync.WaitGroup
	var errs []cb.ErrMessage
	wg.Add(len(paasTaList))
	i := 0
	for _, paasTaInfo := range paasTaList {
		go func(wg *sync.WaitGroup, info mod.Vm) {
			var request mod.VmReq

			request.ServiceName = info.Name
			request.Ip = info.Ip
			request.MetricDatabase = f.InfluxConfig.PaastaDatabase
			request.DefaultTimeRange = f.InfluxConfig.DefaultTimeRange
			request.MeasureTimeList = measureTimeList

			cpuUsage, memResp, diskResp, diskRootResp, err := GetPaasTaSummary_Sub(f.Influxclient, request)

			cpuUsageList = append(cpuUsageList, cpuUsage)
			memoryUsageList = append(memoryUsageList, memResp)
			diskUsageList = append(diskUsageList, diskResp)
			diskRootUsageList = append(diskRootUsageList, diskRootResp)

			if err != nil {
				errs = append(errs, err)
			}
			i++
			wg.Done()
		}(&wg, paasTaInfo)
	}
	wg.Wait()

	//==========================================================================
	// Error가 여러건일 경우 대해 고려해야함.
	if len(errs) > 0 {
		var returnErrMessage string
		for _, err := range errs {
			returnErrMessage = returnErrMessage + " " + err["Message"].(string)
		}
		fmt.Errorf("Error Occur:", returnErrMessage)
	}
	//==========================================================================

	for idx, value := range paasTaList {

		usageResponse[idx].ServiceName = value.Name
		usageResponse[idx].Ip = value.Ip

		for idx_i, cpuValue := range cpuUsageList {
			if value.Name == cpuValue["serviceName"] {
				usageResponse[idx].CpuUsage = util.GetDataFloatFromInterface(cpuUsageList, idx_i)
			}
		}

		for idx_i, value1 := range memoryUsageList {
			if value.Name == value1["serviceName"] {
				usageResponse[idx].MemoryUsage = util.GetDataFromInterface(memoryUsageList, idx_i)
			}
		}

		for idx_i, value1 := range diskUsageList {
			if value.Name == value1["serviceName"] {
				usageResponse[idx].DiskUsage = util.GetDataFloatFromInterface(diskUsageList, idx_i)
			}
		}

		for idx_i, value1 := range diskRootUsageList {
			if value.Name == value1["serviceName"] {
				usageResponse[idx].DiskRootUsage = util.GetDataFloatFromInterface(diskRootUsageList, idx_i)
			}
		}
	}

	//fmt.Println("usageResponse==>", usageResponse)
	//Fail VM 조회
	//makePaasTaFailAlarmData(paasTaList, usageResponse, thresholdList, f, alarmConfig, cb.ALARM_TYPE_CPU)

	makePaasTaAlarmData(paasTaList, usageResponse, thresholdList, f, alarmConfig, cb.ALARM_TYPE_CPU)
	makePaasTaAlarmData(paasTaList, usageResponse, thresholdList, f, alarmConfig, cb.ALARM_TYPE_MEMORY)
	makePaasTaAlarmData(paasTaList, usageResponse, thresholdList, f, alarmConfig, cb.ALARM_TYPE_DISK)
	makePaasTaAlarmData(paasTaList, usageResponse, thresholdList, f, alarmConfig, cb.ALARM_TYPE_ROOTDISK)

}

func makePaasTaFailAlarmData(paasTaList []mod.Vm, systemUsageList []mod.PaasTaResponse, paasTaThresholdList []cb.AlarmThreshold,
	f *BackendServices, alarmConfig mod.AlarmConfig, alarmType string) bool {

	for _, systemUsage := range systemUsageList {

		if systemUsage.CpuUsage <= 0 && systemUsage.MemoryUsage <= 0 && systemUsage.DiskUsage <= 0 {

			for _, data := range paasTaList {

				if systemUsage.ServiceName == data.Name {
					var alarmSource cb.AlarmThreshold
					alarmSource.Level = cb.ALARM_LEVEL_FAIL

					alarmSource.OriginType = cb.ORIGIN_TYPE_PAASTA
					alarmSource.OriginId = data.Id
					alarmSource.AlarmType = cb.ALARM_LEVEL_FAIL
					alarmSource.ServiceName = systemUsage.ServiceName
					alarmSource.Level = cb.ALARM_LEVEL_FAIL
					alarmSource.Ip = systemUsage.Ip

					alarmData := alarm.GetAlarmService(alarmConfig).DBAlarmMessageBuild(alarmSource, 0.0)

					//Alarm전송 및 Alarm Data 생성 여부 조회
					notExist, existData := dao.GetCommonDao(f.Influxclient).GetIsNotExistAlarm(alarmData, f.MonitoringDbClient)

					//기존 AlarmData가 존재하지 않는다면 Alarm 신규 생성
					if notExist == true {
						dao.GetCommonDao(f.Influxclient).CreateAlarmData(alarmData, f.MonitoringDbClient)
						alarm.GetAlarmService(alarmConfig).AlarmSend(alarmSource, alarmData, f.MonitoringDbClient, f.Influxclient, alarmConfig, 0)
					}

					//ResolveStatus가 1이거나 신규 Data면 Mail전송
					if existData.ResolveStatus == "1" || notExist == false {

						if notExist == false {
							now := time.Now()
							alarmSendAvailableTime := existData.AlarmSendDate.Add(time.Duration(3) * time.Minute)
							//DB에서 받아온 시간은 GMT TIme 으로 받아 오기 때문에 9 시간을 뺀다. (config.ini  : gmt.time.hour.gap)
							availTime := alarmSendAvailableTime.Add(time.Duration(f.GmtTimeGapHour) * time.Hour).Unix()

							if now.Unix() >= availTime {
								alarm.GetAlarmService(alarmConfig).AlarmSend(alarmSource, alarmData, f.MonitoringDbClient, f.Influxclient, alarmConfig, 0)
								dao.GetContainerAlarmDao(f.Influxclient).UpdateZoneAlarmSendDate(alarmData, f.MonitoringDbClient)
							}
						}
					}
				}
			}
		}
	}

	return true
}

func makePaasTaAlarmData(paasTaList []mod.Vm, systemUsageList []mod.PaasTaResponse, paasTaThresholdList []cb.AlarmThreshold,
	f *BackendServices, alarmConfig mod.AlarmConfig, alarmType string) {

	var wg sync.WaitGroup
	wg.Add(len(systemUsageList))

	for _, systemUsage := range systemUsageList {
		go func(wg *sync.WaitGroup, systemUsage mod.PaasTaResponse) {

			var thresholdUsgae float64

			if alarmType == cb.ALARM_TYPE_CPU {
				thresholdUsgae = systemUsage.CpuUsage
			} else if alarmType == cb.ALARM_TYPE_MEMORY {
				thresholdUsgae = systemUsage.MemoryUsage
			} else if alarmType == cb.ALARM_TYPE_DISK {
				thresholdUsgae = systemUsage.DiskUsage
			} else if alarmType == cb.ALARM_TYPE_ROOTDISK {
				thresholdUsgae = systemUsage.DiskRootUsage
				alarmType = cb.ALARM_TYPE_DISK
			}

			for _, paasTaThreshold := range paasTaThresholdList {

				if paasTaThreshold.AlarmType == alarmType && paasTaThreshold.Level == cb.ALARM_LEVEL_CRITICAL &&
					paasTaThreshold.Threshold <= thresholdUsgae {

					for _, data := range paasTaList {
						if data.Name == systemUsage.ServiceName {
							paasTaThreshold.OriginId = data.Id
							paasTaThreshold.ServiceName = data.Name
							paasTaThreshold.Ip = data.Ip
							break
						}
					}

					alarmData := alarm.GetAlarmService(alarmConfig).DBAlarmMessageBuild(paasTaThreshold, thresholdUsgae)
					notExist, existData := dao.GetCommonDao(f.Influxclient).GetIsNotExistAlarm(alarmData, f.MonitoringDbClient)

					//fmt.Println("notExist----------------------------------------->",notExist)
					//fmt.Println("existData----------------------------------------->",existData)
					//기존 AlarmData가 존재하지 않는다면 Alarm 신규 생성
					if notExist == true {
						dao.GetCommonDao(f.Influxclient).CreateAlarmData(alarmData, f.MonitoringDbClient)
						alarm.GetAlarmService(alarmConfig).AlarmSend(paasTaThreshold, alarmData, f.MonitoringDbClient, f.Influxclient, alarmConfig, thresholdUsgae)
					}

					//ResolveStatus가 1이거나 신규 Data면 Mail전송
					//ResolveStatus가 2이면 처리중 (Alarm전송 안함)
					if existData.ResolveStatus == "1" && notExist == false {

						if notExist == false {
							now := time.Now()
							alarmSendAvailableTime := existData.AlarmSendDate.Add(time.Duration(paasTaThreshold.RepeatTime) * time.Minute)
							//DB에서 받아온 시간은 GMT TIme 으로 받아 오기 때문에 9 시간을 뺀다. (config.ini  : gmt.time.hour.gap)
							availTime := alarmSendAvailableTime.Add(time.Duration(f.GmtTimeGapHour) * time.Hour).Unix()
							fmt.Println("===========>>>>>>>>>> existData.AlarmSendDate.Unix() : ", existData.AlarmSendDate.Unix())
							fmt.Println("===========>>>>>>>>>> now : ", now.Unix())
							fmt.Println("===========>>>>>>>>>> alarmSendAvailableTime : ", alarmSendAvailableTime)
							fmt.Println("===========>>>>>>>>>> availTime : ", availTime)
							fmt.Println("===========>>>>>>>>>> f.GmtTimeGapHour : ", f.GmtTimeGapHour)
							if now.Unix() >= availTime {
								alarm.GetAlarmService(alarmConfig).AlarmSend(paasTaThreshold, alarmData, f.MonitoringDbClient, f.Influxclient, alarmConfig, thresholdUsgae)
								dao.GetCommonDao(f.Influxclient).UpdateSendDate(alarmData, f.MonitoringDbClient)
							}
						}
					}
					break

					/*if existData.ResolveStatus == "1" || isNotExist == true{
						alarm.GetAlarmService(alarmConfig).AlarmSend(paasTaThreshold, alarmData, f.MonitoringDbClient, f.Influxclient, alarmConfig, thresholdUsgae)
					}*/

				} else if paasTaThreshold.AlarmType == alarmType && paasTaThreshold.Level == cb.ALARM_LEVEL_WARNING &&
					paasTaThreshold.Threshold <= thresholdUsgae {

					//서비스명 설정
					for _, data := range paasTaList {
						if data.Name == systemUsage.ServiceName {
							paasTaThreshold.OriginId = data.Id
							paasTaThreshold.ServiceName = data.Name
							paasTaThreshold.Ip = data.Ip
							break
						}
					}

					alarmData := alarm.GetAlarmService(alarmConfig).DBAlarmMessageBuild(paasTaThreshold, thresholdUsgae)
					notExist, existData := dao.GetCommonDao(f.Influxclient).GetIsNotExistAlarm(alarmData, f.MonitoringDbClient)

					//기존 AlarmData가 존재하지 않는다면 Alarm 신규 생성
					if notExist == true {
						dao.GetCommonDao(f.Influxclient).CreateAlarmData(alarmData, f.MonitoringDbClient)
						alarm.GetAlarmService(alarmConfig).AlarmSend(paasTaThreshold, alarmData, f.MonitoringDbClient, f.Influxclient, alarmConfig, thresholdUsgae)
					}

					//ResolveStatus가 1이거나 신규 Data면 Mail전송
					//ResolveStatus가 2이면 처리중 (Alarm전송 안함)
					if existData.ResolveStatus == "1" && notExist == false {

						if notExist == false {
							now := time.Now()
							alarmSendAvailableTime := existData.AlarmSendDate.Add(time.Duration(paasTaThreshold.RepeatTime) * time.Minute)
							//DB에서 받아온 시간은 GMT TIme 으로 받아 오기 때문에 9 시간을 뺀다. (config.ini  : gmt.time.hour.gap)
							availTime := alarmSendAvailableTime.Add(time.Duration(f.GmtTimeGapHour) * time.Hour).Unix()

							fmt.Println("===========>>>>>>>>>> existData.AlarmSendDate.Unix() : ", existData.AlarmSendDate.Unix())
							fmt.Println("===========>>>>>>>>>> now : ", now.Unix())
							fmt.Println("===========>>>>>>>>>> alarmSendAvailableTime : ", alarmSendAvailableTime)
							fmt.Println("===========>>>>>>>>>> availTime : ", availTime)
							fmt.Println("===========>>>>>>>>>> f.GmtTimeGapHour : ", f.GmtTimeGapHour)

							if now.Unix() >= availTime {
								alarm.GetAlarmService(alarmConfig).AlarmSend(paasTaThreshold, alarmData, f.MonitoringDbClient, f.Influxclient, alarmConfig, thresholdUsgae)
								dao.GetCommonDao(f.Influxclient).UpdateSendDate(alarmData, f.MonitoringDbClient)
							}
						}
					}
					break
				}

			}

			wg.Done()
		}(&wg, systemUsage)
	}
	wg.Wait()
}

//Bosh 상태 목록 조회 - DAO 호출.
func GetPaasTaSummary_Sub(influxClient client.Client, request mod.VmReq) (map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, cb.ErrMessage) {
	var cpuResp, memResp, diskResp, diskRootResp client.Response
	var errs []cb.ErrMessage
	var err cb.ErrMessage
	var wg sync.WaitGroup
	wg.Add(4)
	for i := 0; i < 4; i++ {
		go func(wg *sync.WaitGroup, index int) {
			switch index {
			case 0:
				cpuResp, err = dao.GetPaasTaAlarmDao(influxClient).GetPaasTaCpuUsage(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 1:
				memResp, err = dao.GetPaasTaAlarmDao(influxClient).GetPaasTaMemoryUsage(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 2:
				diskResp, err = dao.GetPaasTaAlarmDao(influxClient).GetPaasTaDiskUsage(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 3:
				diskRootResp, err = dao.GetPaasTaAlarmDao(influxClient).GetPaasTaRootDiskUsage(request)
				if err != nil {
					errs = append(errs, err)
				}
			default:
				break
			}
			wg.Done()
		}(&wg, i)
	}
	wg.Wait()

	//fileSystemResp, err := dao.GetOvvDao(b.influxClient).GetFileSystemUsageList(msg)
	//==========================================================================
	// Error가 여러건일 경우 대해 고려해야함.
	if len(errs) > 0 {
		var returnErrMessage string
		for _, err := range errs {
			returnErrMessage = returnErrMessage + " " + err["Message"].(string)
		}
		errMessage := cb.ErrMessage{
			"Message": returnErrMessage,
		}
		return nil, nil, nil, nil, errMessage
	}
	//==========================================================================
	//fmt.Println("fileSystemResp==========>", fileSystemResp)

	cpuUsage, _ := util.GetResponseConverter().InfluxConverter(cpuResp, request.ServiceName)
	memoryUsage, _ := util.GetResponseConverter().InfluxConverter4Usage(memResp, request.ServiceName)
	diskUsage, _ := util.GetResponseConverter().InfluxConverter(diskResp, request.ServiceName)
	diskRootUsage, _ := util.GetResponseConverter().InfluxConverter(diskRootResp, request.ServiceName)

	return cpuUsage, memoryUsage, diskUsage, diskRootUsage, nil

}

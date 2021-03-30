package service

import (
	"fmt"
	client "github.com/influxdata/influxdb/client/v2"
	"kr/paasta/monitoring-batch/alarm"
	"kr/paasta/monitoring-batch/dao"
	mod "kr/paasta/monitoring-batch/model"
	cb "kr/paasta/monitoring-batch/model/base"
	"kr/paasta/monitoring-batch/util"
	"strings"
	"sync"
	"time"
)

func BoshAlarmCollect(f *BackendServices) {

	var boshInfo mod.Bosh
	var usageResponse mod.BoshResponse

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
	boshIp := strings.Split(f.BoshConfig.BoshUrl, ":")
	boshInfo.ServiceName = f.BoshConfig.ServiceName
	boshInfo.Ip = boshIp[0]

	alarmPolicy, _ := dao.GetBoshAlarmDao(f.Influxclient).GetBoshAlarmPolicy(f.MonitoringDbClient)

	var measureTimeList []mod.AlarmItemMeasureTime

	for _, data := range alarmPolicy {

		var thresholdCritical cb.AlarmThreshold
		var warningThreshold cb.AlarmThreshold

		thresholdCritical.OriginType = cb.ORIGIN_TYPE_BOSH
		thresholdCritical.OriginId = cb.DEFAULT_ORIGIN_ID
		thresholdCritical.AlarmType = data.AlarmType
		thresholdCritical.ServiceName = f.BoshConfig.ServiceName
		thresholdCritical.Level = cb.ALARM_LEVEL_CRITICAL
		thresholdCritical.Threshold = float64(data.CriticalThreshold)
		thresholdCritical.RepeatTime = data.RepeatTime

		warningThreshold.OriginType = cb.ORIGIN_TYPE_BOSH
		warningThreshold.OriginId = cb.DEFAULT_ORIGIN_ID
		warningThreshold.AlarmType = data.AlarmType
		warningThreshold.ServiceName = f.BoshConfig.ServiceName
		warningThreshold.Level = cb.ALARM_LEVEL_WARNING
		warningThreshold.Threshold = float64(data.WarningThreshold)
		warningThreshold.RepeatTime = data.RepeatTime

		thresholdList = append(thresholdList, thresholdCritical)
		thresholdList = append(thresholdList, warningThreshold)

		measureTimeList = append(measureTimeList, mod.AlarmItemMeasureTime{data.AlarmType, data.MeasureTime})

	}

	var cpuUsageList []map[string]interface{}
	var memoryUsageList []map[string]interface{}
	var diskUsageList []map[string]interface{}
	var diskRootUsageList []map[string]interface{}

	//response := make([]mod.BoshResponse, len(boshInfo))

	var wg sync.WaitGroup
	var errs []cb.ErrMessage
	wg.Add(1)
	i := 0

	//Bosh의 CPU/Memory/Disk 사용률을 조회한다.
	go func(wg *sync.WaitGroup, info mod.Bosh) {
		var request mod.BoshReq

		request.ServiceName = info.ServiceName
		request.Ip = info.Ip
		request.MetricDatabase = f.InfluxConfig.InfraDatabase
		request.DefaultTimeRange = f.InfluxConfig.DefaultTimeRange
		request.MeasureTimeList = measureTimeList

		cpuUsage, memUsage, diskUsage, diskRootUsage, err := GetBoshSummary_Sub(f.Influxclient, request)
		fmt.Println("########## get bosh memory usage =>", memUsage)
		cpuUsageList = append(cpuUsageList, cpuUsage)
		memoryUsageList = append(memoryUsageList, memUsage)
		diskUsageList = append(diskUsageList, diskUsage)
		diskRootUsageList = append(diskRootUsageList, diskRootUsage)

		if err != nil {
			errs = append(errs, err)
		}
		i++
		wg.Done()
	}(&wg, boshInfo)

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

	for idx_i, cpuValue := range cpuUsageList {
		if boshInfo.ServiceName == cpuValue["serviceName"] {
			usageResponse.CpuUsage = util.GetDataFloatFromInterface(cpuUsageList, idx_i)

		}
	}

	for idx_i, memoryValue := range memoryUsageList {
		if boshInfo.ServiceName == memoryValue["serviceName"] {
			usageResponse.MemoryUsage = util.GetDataFromInterface(memoryUsageList, idx_i)
		}
	}

	for idx_i, diskValue := range diskUsageList {
		if boshInfo.ServiceName == diskValue["serviceName"] {
			usageResponse.DiskUsage = util.GetDataFloatFromInterface(diskUsageList, idx_i)
		}
	}

	for idx_i, diskValue := range diskRootUsageList {
		if boshInfo.ServiceName == diskValue["serviceName"] {
			usageResponse.DiskRootUsage = util.GetDataFloatFromInterface(diskRootUsageList, idx_i)
		}
	}

	usageResponse.ServiceName = f.BoshConfig.ServiceName

	if makeBoshFailAlarmData(usageResponse, thresholdList, f, alarmConfig, boshInfo.Ip) == false {
		// Cpu/Memory/Disk임계치 초과 체크한다.
		var wg sync.WaitGroup
		wg.Add(4)
		for i := 0; i < 4; i++ {
			go func(wg *sync.WaitGroup, index int) {
				switch index {
				case 0:
					makeBoshAlarmData(usageResponse, thresholdList, f, alarmConfig, boshInfo.Ip, cb.ALARM_TYPE_CPU)
				case 1:
					makeBoshAlarmData(usageResponse, thresholdList, f, alarmConfig, boshInfo.Ip, cb.ALARM_TYPE_MEMORY)
				case 2:
					makeBoshAlarmData(usageResponse, thresholdList, f, alarmConfig, boshInfo.Ip, cb.ALARM_TYPE_DISK)
				case 3:
					makeBoshAlarmData(usageResponse, thresholdList, f, alarmConfig, boshInfo.Ip, cb.ALARM_TYPE_ROOTDISK)
				default:
					break
				}
				wg.Done()
			}(&wg, i)
		}
		wg.Wait()
	}

}

func makeBoshFailAlarmData(systemUsage mod.BoshResponse, boshThresholds []cb.AlarmThreshold, f *BackendServices, alarmConfig mod.AlarmConfig, boshIp string) bool {

	var alarmSource cb.AlarmThreshold
	alarmSource.OriginType = boshThresholds[0].OriginType
	alarmSource.OriginId = boshThresholds[0].OriginId
	alarmSource.AlarmType = cb.ALARM_LEVEL_FAIL
	alarmSource.ServiceName = boshThresholds[0].ServiceName
	alarmSource.Level = cb.ALARM_LEVEL_FAIL
	alarmSource.Ip = boshIp

	if systemUsage.CpuUsage == 0 && systemUsage.MemoryUsage == 0 && systemUsage.DiskUsage == 0 && systemUsage.DiskRootUsage == 0 {

		boshThresholds[0].Level = cb.ALARM_LEVEL_FAIL
		alarmData := alarm.GetAlarmService(alarmConfig).DBAlarmMessageBuild(alarmSource, 0.0)

		//Alarm전송 및 Alarm Data 생성 여부 조회
		notExist, existData := dao.GetCommonDao(f.Influxclient).GetIsNotExistAlarm(alarmData, f.MonitoringDbClient)

		//기존 AlarmData가 존재하지 않는다면 Alarm 신규 생성
		if notExist == true {
			dao.GetCommonDao(f.Influxclient).CreateAlarmData(alarmData, f.MonitoringDbClient)
			alarm.GetAlarmService(alarmConfig).AlarmSend(boshThresholds[0], alarmData, f.MonitoringDbClient, f.Influxclient, alarmConfig, 0)
		}

		//ResolveStatus가 1이거나 신규 Data면 Mail전송
		if existData.ResolveStatus == "1" || notExist == false {
			if notExist == false {
				now := time.Now()
				alarmSendAvailableTime := existData.AlarmSendDate.Add(time.Duration(boshThresholds[0].RepeatTime) * time.Minute)
				//DB에서 받아온 시간은 GMT TIme 으로 받아 오기 때문에 9 시간을 뺀다. (config.ini  : gmt.time.hour.gap)
				availTime := alarmSendAvailableTime.Add(time.Duration(f.GmtTimeGapHour) * time.Hour).Unix()

				if now.Unix() >= availTime {
					alarm.GetAlarmService(alarmConfig).AlarmSend(boshThresholds[0], alarmData, f.MonitoringDbClient, f.Influxclient, alarmConfig, 0)
					dao.GetCommonDao(f.Influxclient).UpdateSendDate(alarmData, f.MonitoringDbClient)
				}
			}
		}

		return true
	} else {
		return false
	}
}

func makeBoshAlarmData(systemUsage mod.BoshResponse, boshThresholds []cb.AlarmThreshold, f *BackendServices, alarmConfig mod.AlarmConfig, boshIp, alarmType string) {

	var thresholdUsgae float64
	if alarmType == cb.ALARM_TYPE_CPU {
		thresholdUsgae = systemUsage.CpuUsage
	} else if alarmType == cb.ALARM_TYPE_MEMORY {
		thresholdUsgae = systemUsage.MemoryUsage
	} else if alarmType == cb.ALARM_TYPE_DISK {
		thresholdUsgae = systemUsage.DiskUsage
	} else if alarmType == cb.ALARM_TYPE_ROOTDISK {
		thresholdUsgae = systemUsage.DiskUsage
		alarmType = cb.ALARM_TYPE_DISK
	}

	for _, boshThreshold := range boshThresholds {

		if systemUsage.ServiceName == boshThreshold.ServiceName {

			/*fmt.Println("::::",boshThreshold.Level, thresholdUsgae)
			fmt.Println("::::",alarmType, boshThreshold.Threshold)*/

			if boshThreshold.AlarmType == alarmType && boshThreshold.Level == cb.ALARM_LEVEL_CRITICAL &&
				boshThreshold.Threshold <= thresholdUsgae {

				boshThreshold.Ip = boshIp
				alarmData := alarm.GetAlarmService(alarmConfig).DBAlarmMessageBuild(boshThreshold, thresholdUsgae)
				notExist, existData := dao.GetCommonDao(f.Influxclient).GetIsNotExistAlarm(alarmData, f.MonitoringDbClient)

				//기존 AlarmData가 존재하지 않는다면 Alarm 신규 생성
				if notExist == true {
					dao.GetCommonDao(f.Influxclient).CreateAlarmData(alarmData, f.MonitoringDbClient)
					alarm.GetAlarmService(alarmConfig).AlarmSend(boshThreshold, alarmData, f.MonitoringDbClient, f.Influxclient, alarmConfig, thresholdUsgae)
				}

				if existData.ResolveStatus == "1" || notExist == false {
					now := time.Now()
					alarmSendAvailableTime := existData.AlarmSendDate.Add(time.Duration(boshThreshold.RepeatTime) * time.Minute)
					//DB에서 받아온 시간은 GMT TIme 으로 받아 오기 때문에 9 시간을 뺀다. (config.ini  : gmt.time.hour.gap)
					availTime := alarmSendAvailableTime.Add(time.Duration(f.GmtTimeGapHour) * time.Hour).Unix()

					if now.Unix() >= availTime {
						alarm.GetAlarmService(alarmConfig).AlarmSend(boshThreshold, alarmData, f.MonitoringDbClient, f.Influxclient, alarmConfig, thresholdUsgae)
						dao.GetCommonDao(f.Influxclient).UpdateSendDate(alarmData, f.MonitoringDbClient)
					}
				}

				//ResolveStatus가 1이거나 신규 Data면 Mail전송
				/*if existData.ResolveStatus == "1" || isNotExist == true{
					alarm.GetAlarmService(alarmConfig).AlarmSend(boshThreshold, alarmData, f.MonitoringDbClient, f.Influxclient, alarmConfig, thresholdUsgae)
				}*/
				break

			} else if boshThreshold.AlarmType == alarmType && boshThreshold.Level == cb.ALARM_LEVEL_WARNING &&
				boshThreshold.Threshold <= thresholdUsgae {

				if alarmType == cb.ALARM_TYPE_DISK {
					fmt.Println("Disk Usage:::::", boshThreshold.Threshold)
				}
				boshThreshold.Ip = boshIp

				alarmData := alarm.GetAlarmService(alarmConfig).DBAlarmMessageBuild(boshThreshold, thresholdUsgae)
				notExist, existData := dao.GetCommonDao(f.Influxclient).GetIsNotExistAlarm(alarmData, f.MonitoringDbClient)

				//기존 AlarmData가 존재하지 않는다면 Alarm 신규 생성
				if notExist == true {
					dao.GetCommonDao(f.Influxclient).CreateAlarmData(alarmData, f.MonitoringDbClient)
					alarm.GetAlarmService(alarmConfig).AlarmSend(boshThreshold, alarmData, f.MonitoringDbClient, f.Influxclient, alarmConfig, thresholdUsgae)
				}
				//ResolveStatus가 1이거나 신규 Data면 Mail전송
				//ResolveStatus가 2이면 처리중 (Alarm전송 안함)
				if existData.ResolveStatus == "1" && notExist == false {
					now := time.Now()
					alarmSendAvailableTime := existData.AlarmSendDate.Add(time.Duration(boshThreshold.RepeatTime) * time.Minute)
					//DB에서 받아온 시간은 GMT TIme 으로 받아 오기 때문에 9 시간을 뺀다. (config.ini  : gmt.time.hour.gap)
					availTime := alarmSendAvailableTime.Add(time.Duration(f.GmtTimeGapHour) * time.Hour).Unix()

					if now.Unix() >= availTime {
						alarm.GetAlarmService(alarmConfig).AlarmSend(boshThreshold, alarmData, f.MonitoringDbClient, f.Influxclient, alarmConfig, thresholdUsgae)
						dao.GetCommonDao(f.Influxclient).UpdateSendDate(alarmData, f.MonitoringDbClient)
					}
				}

				break
			}
		}
	}
}

//Bosh 상태 목록 조회 - DAO 호출.
func GetBoshSummary_Sub(influxClient client.Client, request mod.BoshReq) (map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, cb.ErrMessage) {

	var cpuResp, memUsageResp, diskResp, diskRootResp client.Response
	var errs []cb.ErrMessage
	var err cb.ErrMessage
	var wg sync.WaitGroup
	wg.Add(4)
	for i := 0; i < 4; i++ {
		go func(wg *sync.WaitGroup, index int) {
			switch index {
			case 0:
				cpuResp, err = dao.GetBoshAlarmDao(influxClient).GetBoshCpuUsage(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 1:
				memUsageResp, err = dao.GetBoshAlarmDao(influxClient).GetBoshMemoryUsage(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 2:
				diskResp, err = dao.GetBoshAlarmDao(influxClient).GetBoshDiskUsage(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 3:
				diskRootResp, err = dao.GetBoshAlarmDao(influxClient).GetBoshRootDiskUsage(request)
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

	cpuUsage, _ := util.GetResponseConverter().InfluxConverter(cpuResp, request.ServiceName)
	memUsage, _ := util.GetResponseConverter().InfluxConverter4Usage(memUsageResp, request.ServiceName)
	diskUsage, _ := util.GetResponseConverter().InfluxConverter(diskResp, request.ServiceName)
	diskRootUsage, _ := util.GetResponseConverter().InfluxConverter(diskRootResp, request.ServiceName)

	return cpuUsage, memUsage, diskUsage, diskRootUsage, nil

}

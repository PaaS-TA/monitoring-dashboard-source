package service

import (
	"fmt"
	"sync"
	"time"
	"reflect"
	"strconv"
	client "github.com/influxdata/influxdb/client/v2"
	"monitoring-batch/dao"
	"monitoring-batch/util"
	"monitoring-batch/alarm"
	mod "monitoring-batch/model"
	cb "monitoring-batch/model/base"
)

func ContainerAlarmCollect(f *BackendServices) {

	var alarmConfig mod.AlarmConfig
	alarmConfig.SmtpHost = f.MailConfig.SmtpHost
	alarmConfig.Port = f.MailConfig.Port
	alarmConfig.MailSender = f.MailConfig.MailSender
	alarmConfig.SenderPwd = f.MailConfig.SenderPwd
	alarmConfig.MailResource = f.MailConfig.ResouceUrl
	alarmConfig.MailReceiver = f.MailConfig.MailReceiver
	alarmConfig.AlarmSend = f.MailConfig.AlarmSend
	alarmConfig.MailTlsSend = f.MailConfig.MailTlsSend

	cellInfoList, dbErr := dao.GetContainerAlarmDao(f.Influxclient).GetCellList(f.MonitoringDbClient)

	if dbErr != nil {
		fmt.Errorf("Error:", dbErr)
	}

	var thresholdList []cb.AlarmThreshold
	alarmPolicy, _ := dao.GetContainerAlarmDao(f.Influxclient).GetContainerAlarmPolicy(f.MonitoringDbClient)

	var measureTimeList []mod.AlarmItemMeasureTime

	//임계치 종류에는 6가지 이다.
	// CPU    - Critical/warning
	// Memory - Critical/warning
	// Disk   - Critical/warning

	for _, data := range alarmPolicy {

		var thresholdCritical cb.AlarmThreshold
		var warningThreshold cb.AlarmThreshold

		thresholdCritical.OriginType = cb.ORIGIN_TYPE_CONTAINER
		thresholdCritical.AlarmType = data.AlarmType
		thresholdCritical.Level = cb.ALARM_LEVEL_CRITICAL
		thresholdCritical.Threshold = float64(data.CriticalThreshold)
		thresholdCritical.RepeatTime = data.RepeatTime

		warningThreshold.OriginType = cb.ORIGIN_TYPE_CONTAINER
		warningThreshold.AlarmType = data.AlarmType
		warningThreshold.Level = cb.ALARM_LEVEL_WARNING
		warningThreshold.Threshold = float64(data.WarningThreshold)
		warningThreshold.RepeatTime = data.RepeatTime

		thresholdList = append(thresholdList, thresholdCritical)
		thresholdList = append(thresholdList, warningThreshold)

		measureTimeList = append(measureTimeList, mod.AlarmItemMeasureTime{data.AlarmType, data.MeasureTime})
	}

	//MetricData에서 Cell에서 실행되는 App Container목록을 추출한다.
	cellMap := getZoneCellList(cellInfoList, f)
	//Container에 Cpu/Memory/Disk 사용률을 추출한다.
	metricCellList, _ := mapToTreeStruct(cellMap, f, measureTimeList)

	// Cpu/Memory/Disk임계치 초과 체크한다.
	makeContainerAlarmData(metricCellList, cellInfoList, thresholdList, f, alarmConfig, cb.ALARM_TYPE_CPU)
	makeContainerAlarmData(metricCellList, cellInfoList, thresholdList, f, alarmConfig, cb.ALARM_TYPE_MEMORY)
	makeContainerAlarmData(metricCellList, cellInfoList, thresholdList, f, alarmConfig, cb.ALARM_TYPE_DISK)
}

func makeContainerAlarmData(systemUsageList []mod.CellTileView, dbCellList []mod.ZoneCellInfo, containerThresholds []cb.AlarmThreshold, f *BackendServices, alarmConfig mod.AlarmConfig, alarmType string) {

	var wg sync.WaitGroup
	wg.Add(len(systemUsageList))

	for _, cellStatusData := range systemUsageList {
		go func(wg *sync.WaitGroup, cellStatusData mod.CellTileView) {

			for _, containerStatusData := range cellStatusData.ContainerTileView {

				var thresholdUsgae float64
				if alarmType == cb.ALARM_TYPE_CPU {
					thresholdUsgae = containerStatusData.CpuUsage
				} else if alarmType == cb.ALARM_TYPE_MEMORY {
					thresholdUsgae = containerStatusData.MemoryUsage
					fmt.Println(thresholdUsgae)
				} else if alarmType == cb.ALARM_TYPE_DISK {
					thresholdUsgae = containerStatusData.DiskUsage
				}

				for _, containerThreshold := range containerThresholds {

					if containerThreshold.AlarmType == alarmType && containerThreshold.Level == cb.ALARM_LEVEL_CRITICAL &&
						containerThreshold.Threshold <= thresholdUsgae {

						alarmData := alarm.GetAlarmService(alarmConfig).DBAppAlarmMessageBuild(containerThreshold, containerStatusData, cellStatusData.CellName, thresholdUsgae)
						notExist, existData := dao.GetContainerAlarmDao(f.Influxclient).GetContainerIsExistAlarm(alarmData, f.MonitoringDbClient)

						//기존 AlarmData가 존재하지 않는다면 Alarm 신규 생성
						if notExist == true {
							for _, data := range dbCellList {
								if data.CellName == cellStatusData.CellName {
									alarmData.OriginId = data.Id
									alarmData.Ip = data.Ip
									dao.GetContainerAlarmDao(f.Influxclient).CreateContainerAlarmData(alarmData, f.MonitoringDbClient)
									alarm.GetAlarmService(alarmConfig).AlarmSend(containerThreshold, alarmData, f.MonitoringDbClient, f.Influxclient, alarmConfig, thresholdUsgae)
									break
								}
							}

						}

						//ResolveStatus가 1이거나 신규 Data면 Mail전송
						//ResolveStatus가 2이면 처리중 (Alarm전송 안함)
						if existData.ResolveStatus == "1" && notExist == false {

							if notExist == false {
								now := time.Now()
								alarmSendAvailableTime := existData.AlarmSendDate.Add(time.Duration(containerThreshold.RepeatTime) * time.Minute)
								//DB에서 받아온 시간은 GMT TIme 으로 받아 오기 때문에 9 시간을 뺀다. (config.ini  : gmt.time.hour.gap)
								availTime := alarmSendAvailableTime.Add(time.Duration(f.GmtTimeGapHour) * time.Hour).Unix()

								if now.Unix() >= availTime {
									for _, data := range dbCellList {
										if data.CellName == cellStatusData.CellName {
											alarmData.OriginId = data.Id
											alarmData.Ip = data.Ip
											alarm.GetAlarmService(alarmConfig).AlarmSend(containerThreshold, alarmData, f.MonitoringDbClient, f.Influxclient, alarmConfig, thresholdUsgae)
											dao.GetContainerAlarmDao(f.Influxclient).UpdateZoneAlarmSendDate(alarmData, f.MonitoringDbClient)
											break
										}
									}

								}

							}
						}
						break

					} else if containerThreshold.AlarmType == alarmType && containerThreshold.Level == cb.ALARM_LEVEL_WARNING &&
						containerThreshold.Threshold <= thresholdUsgae {

						alarmData := alarm.GetAlarmService(alarmConfig).DBAppAlarmMessageBuild(containerThreshold, containerStatusData, cellStatusData.CellName, thresholdUsgae)

						notExist, existData := dao.GetContainerAlarmDao(f.Influxclient).GetContainerIsExistAlarm(alarmData, f.MonitoringDbClient)

						//기존 AlarmData가 존재하지 않는다면 Alarm 신규 생성
						if notExist == true {
							for _, data := range dbCellList {
								if data.CellName == cellStatusData.CellName {
									alarmData.OriginId = data.Id
									alarmData.Ip = data.Ip
									dao.GetContainerAlarmDao(f.Influxclient).CreateContainerAlarmData(alarmData, f.MonitoringDbClient)
									alarm.GetAlarmService(alarmConfig).AlarmSend(containerThreshold, alarmData, f.MonitoringDbClient, f.Influxclient, alarmConfig, thresholdUsgae)
									break
								}
							}

						}
						fmt.Println("============existData.ResolveStatus:::", existData.ResolveStatus)
						fmt.Println("============notExist:", notExist)
						//ResolveStatus가 1이거나 신규 Data면 Mail전송
						//ResolveStatus가 2이면 처리중 (Alarm전송 안함)
						if existData.ResolveStatus == "1" && notExist == false {
							fmt.Println("=======ininni=======")
							now := time.Now()
							alarmSendAvailableTime := existData.AlarmSendDate.Add(time.Duration(containerThreshold.RepeatTime) * time.Minute)
							//DB에서 받아온 시간은 GMT TIme 으로 받아 오기 때문에 9 시간을 뺀다. (config.ini  : gmt.time.hour.gap)
							availTime := alarmSendAvailableTime.Add(time.Duration(f.GmtTimeGapHour) * time.Hour).Unix()

							if now.Unix() >= availTime {
								for _, data := range dbCellList {
									if data.CellName == cellStatusData.CellName {
										alarmData.OriginId = data.Id
										alarmData.Ip = data.Ip
										alarm.GetAlarmService(alarmConfig).AlarmSend(containerThreshold, alarmData, f.MonitoringDbClient, f.Influxclient, alarmConfig, thresholdUsgae)
										dao.GetContainerAlarmDao(f.Influxclient).UpdateZoneAlarmSendDate(alarmData, f.MonitoringDbClient)
										break
									}
								}

							}
						}
						break
					}

				}
			}

			wg.Done()
		}(&wg, cellStatusData)
	}
	wg.Wait()
}

func getZoneCellList(cellInfos []mod.ZoneCellInfo, f *BackendServices) map[string]map[string]map[string]string {

	cellMap := make(map[string]map[string]map[string]string)

	//Zone에 존재하는 Cell들에 실행되고 있는 Container 목록을 받아온다.
	for _, cellInfo := range cellInfos {
		var request mod.ZonesReq
		request.CellIp = cellInfo.Ip
		request.MetricDatabase = f.InfluxConfig.ContainerDatabase
		containerResp, _ := dao.GetContainerAlarmDao(f.Influxclient).GetCellContainersList(request)
		valueList, _ := util.GetResponseConverter().InfluxConverterToMap(containerResp)

		appMap := make(map[string]map[string]string)
		for _, value := range valueList {

			containerMap := make(map[string]string)
			appName := reflect.ValueOf(value["application_name"]).String()
			containerName := reflect.ValueOf(value["container_interface"]).String()
			applicationIndex := reflect.ValueOf(value["application_index"]).String()

			containerMap[containerName] = applicationIndex

			//동일한 App의 Container는 AppMap에 Append 처리 한다.
			if exists, ok := appMap[appName]; ok {
				for k, v := range containerMap {
					exists[k] = v
					appMap[appName] = exists
				}
			} else {
				appMap[appName] = containerMap
			}

		}
		cellMap[cellInfo.CellName] = appMap
	}
	return cellMap
}

func mapToTreeStruct(mapData map[string]map[string]map[string]string, f *BackendServices, measureTimeList []mod.AlarmItemMeasureTime) ([]mod.CellTileView, cb.ErrMessage) {

	returnValue := make([]mod.CellTileView, len(mapData))
	cellInfo := make([]mod.CellTileView, len(mapData))

	c := 0

	for cellName, apps := range mapData {

		var containers []mod.ContainerTileView

		var wg sync.WaitGroup
		var errs []cb.ErrMessage
		wg.Add(len(apps))

		for appName, containerInfos := range apps {
			go func(wg *sync.WaitGroup, info map[string]string, containerAppName string) {

				var containerInfo mod.ContainerTileView

				fmt.Println("container Info >>>>", info)
				for name, index := range info {

					idx, _ := strconv.Atoi(index)
					containerInfo.AppIndex = idx
					containerInfo.ContainerName = name
					containerInfo.AppName = containerAppName

					var request mod.ZonesReq
					request.ContainerName = name
					request.MetricDatabase = f.InfluxConfig.ContainerDatabase
					request.MeasureTimeList = measureTimeList
					cpuData, memData, diskData, err := GetZoneSummary_Sub(request, f)

					if err != nil {
						errs = append(errs, err)
						fmt.Println("####Container Error####", err)
					}

					cpuUsage := util.GetDataFloatFromInterfaceSingle(cpuData)
					memUsage := util.GetDataFloatFromInterfaceSingle(memData)
					diskUsage := util.GetDataFloatFromInterfaceSingle(diskData)

					containerInfo.CpuUsage = cpuUsage
					containerInfo.MemoryUsage = memUsage
					containerInfo.DiskUsage = diskUsage

					fmt.Println("Container usage ===>", containerInfo)

					containers = append(containers, containerInfo)
					/*}else{
						fmt.Println("####Error#### containerInfo nil===>", containerInfo)
						fmt.Println("####Error#### cpuUsage===>", cpuUsage)
						fmt.Println("####Error#### memUsage===>", memUsage)
						fmt.Println("####Error#### diskUsage===>", diskUsage)
					}*/

				}
				wg.Done()
			}(&wg, containerInfos, appName)
		}
		wg.Wait()

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
			return nil, errMessage
		}
		//==========================================================================

		cellInfo[c].CellName = cellName
		cellInfo[c].ContainerTileView = containers
		c++

	}

	sortIdx := 0
	for cellName, _ := range mapData {
		for _, info := range cellInfo {
			if cellName == info.CellName {
				returnValue[sortIdx].CellName = cellName
				returnValue[sortIdx].ContainerTileView = info.ContainerTileView
			}
		}
		sortIdx++
	}

	return returnValue, nil
}

//Server 상태 목록 조회 - DAO 호출.
func GetZoneSummary_Sub(request mod.ZonesReq, f *BackendServices) (map[string]interface{}, map[string]interface{}, map[string]interface{}, cb.ErrMessage) {
	var cpuResp, memResp, diskResp client.Response
	var errs []cb.ErrMessage
	var err cb.ErrMessage
	var wg sync.WaitGroup
	wg.Add(3)

	for i := 0; i < 3; i++ {
		go func(wg *sync.WaitGroup, index int) {

			switch index {
			case 0:
				cpuResp, err = dao.GetContainerAlarmDao(f.Influxclient).GetContainerCpuUsage(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 1:
				memResp, err = dao.GetContainerAlarmDao(f.Influxclient).GetContainerMemoryUsage(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 2:
				diskResp, err = dao.GetContainerAlarmDao(f.Influxclient).GetContainerOvvDiskUsage(request)
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
		return nil, nil, nil, errMessage
	}
	//==========================================================================

	cpuUsage, _ := util.GetResponseConverter().InfluxConverter(cpuResp, request.ContainerName)
	memUsage, _ := util.GetResponseConverter().InfluxConverter(memResp, request.ContainerName)
	diskUsage, _ := util.GetResponseConverter().InfluxConverter(diskResp, request.ContainerName)

	return cpuUsage, memUsage, diskUsage, nil

}

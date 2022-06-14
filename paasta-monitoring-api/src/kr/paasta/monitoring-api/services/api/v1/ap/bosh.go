package ap

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/jinzhu/gorm"
	dao "paasta-monitoring-api/dao/api/v1/ap"
	"paasta-monitoring-api/helpers"
	models "paasta-monitoring-api/models/api/v1"
	"sort"
	"strconv"
	"sync"
	"time"
)

type ApBoshService struct {
	DbInfo         *gorm.DB
	InfluxDBClient client.Client
	BoshInfoList   []models.Bosh
	Env            map[string]interface{}
}

func GetApBoshService(DbInfo *gorm.DB, InfluxDBClient client.Client, BoshInfoList []models.Bosh, Env map[string]interface{}) *ApBoshService {
	return &ApBoshService{
		DbInfo:         DbInfo,
		InfluxDBClient: InfluxDBClient,
		BoshInfoList:   BoshInfoList,
		Env:            Env,
	}
}

func (b *ApBoshService) GetBoshInfoList() ([]models.Bosh, error) {
	// after Use Database
	results := b.BoshInfoList
	return results, nil
}

func (b *ApBoshService) GetBoshOverview() (models.BoshOverview, error) {
	var result models.BoshOverview
	boshSummary, err := b.GetBoshSummary()
	if err != nil {
		fmt.Println(err.Error())
		return result, err
	}

	// bosh overview
	totalCnt, failedCnt, criticalCnt, warningCnt := len(boshSummary), 0, 0, 0
	for _, value := range boshSummary {
		if value.BoshSummaryMetric.BoshState == models.BOSH_STATE_FAIL {
			failedCnt++
		} else if value.BoshSummaryMetric.BoshState == models.ALARM_LEVEL_CRITICAL {
			criticalCnt++
		} else if value.BoshSummaryMetric.BoshState == models.ALARM_LEVEL_WARNING {
			warningCnt++
		}
	}
	result.Total = strconv.Itoa(totalCnt)
	result.Running = strconv.Itoa(totalCnt - failedCnt - criticalCnt - warningCnt)
	result.Failed = strconv.Itoa(failedCnt)
	result.Critical = strconv.Itoa(criticalCnt)
	result.Warning = strconv.Itoa(warningCnt)

	return result, nil
}

func (b *ApBoshService) GetBoshSummary() ([]models.BoshSummary, error) {
	var results []models.BoshSummary
	//var errs []models.ErrMessage
	var resultsResponse map[string]interface{}

	//임계치 설정정보를 조회한다.
	serverThresholds, err := dao.GetApAlarmDao(b.DbInfo).GetAlarmPolicy()
	if err != nil {
		fmt.Println(err.Error())
		return results, err
	}
	//fmt.Println(serverThresholds)

	for _, boshInfo := range b.BoshInfoList {
		var boshSummary models.BoshSummary
		boshSummary.Name = boshInfo.Name
		boshSummary.Ip = boshInfo.Ip
		boshSummary.UUID = boshInfo.UUID
		cpuCoreData, cpuData, memTotData, memFreeData, diskTotalData, diskUsedData, diskDataTotalData, diskDataUsedData, errs := b.GetBoshSummaryMetric(boshSummary)
		fmt.Println(resultsResponse)
		fmt.Println(errs)

		cpuUsage := helpers.GetDataFloatFromInterfaceSingle(cpuData)
		memTot := helpers.GetDataFloatFromInterfaceSingle(memTotData)
		memFree := helpers.GetDataFloatFromInterfaceSingle(memFreeData)
		memUsage := helpers.RoundFloatDigit2(100 - ((memFree / memTot) * 100))
		diskTotal := helpers.GetDataFloatFromInterfaceSingle(diskTotalData)
		diskUsed := helpers.GetDataFloatFromInterfaceSingle(diskUsedData)
		diskUsage := 100 - ((diskTotal - diskUsed) / diskTotal * 100)

		diskDataTotal := helpers.GetDataFloatFromInterfaceSingle(diskDataTotalData)
		diskDataUsed := helpers.GetDataFloatFromInterfaceSingle(diskDataUsedData)
		diskDataUsage := 100 - ((diskDataTotal - diskDataUsed) / diskDataTotal * 100)

		var boshSummaryMetric models.BoshSummaryMetric
		boshSummaryMetric.Core = strconv.Itoa(len(cpuCoreData))
		boshSummaryMetric.CpuUsage = helpers.RoundFloat(cpuUsage, 2)
		boshSummaryMetric.TotalMemory = memTot / models.MB
		boshSummaryMetric.MemoryUsage = memUsage
		boshSummaryMetric.TotalDisk = diskTotal / models.MB
		boshSummaryMetric.DataDisk = diskDataTotal / models.MB

		if boshSummaryMetric.Core == "0" || boshSummaryMetric.TotalMemory == 0 {
			boshSummaryMetric.State, boshSummaryMetric.BoshState, boshSummaryMetric.CpuErrStat, boshSummaryMetric.MemErrStat = models.BOSH_STATE_FAIL, models.BOSH_STATE_FAIL, models.BOSH_STATE_FAIL, models.BOSH_STATE_FAIL
		}

		if boshSummaryMetric.TotalDisk == 0 || boshSummaryMetric.DataDisk == 0 {
			boshSummaryMetric.DiskStatus, boshSummaryMetric.BoshState, boshSummaryMetric.DiskRootErrStat, boshSummaryMetric.DiskDataErrStat = models.BOSH_STATE_FAIL, models.BOSH_STATE_FAIL, models.BOSH_STATE_FAIL, models.BOSH_STATE_FAIL
		}

		// bosh state setting
		if boshSummaryMetric.State != models.BOSH_STATE_FAIL {
			var alarmStatus []string

			cpuStatus := helpers.GetAlarmStatusByServiceName(models.ORIGIN_TYPE_BOSH, models.ALARM_TYPE_CPU, boshSummaryMetric.CpuUsage, serverThresholds)
			memStatus := helpers.GetAlarmStatusByServiceName(models.ORIGIN_TYPE_BOSH, models.ALARM_TYPE_MEMORY, boshSummaryMetric.MemoryUsage, serverThresholds)

			if cpuStatus != "" {
				alarmStatus = append(alarmStatus, cpuStatus)
				boshSummaryMetric.CpuErrStat = cpuStatus
			} else {
				boshSummaryMetric.CpuErrStat = models.BOSH_STATE_RUNNING
			}
			if memStatus != "" {
				alarmStatus = append(alarmStatus, memStatus)
				boshSummaryMetric.MemErrStat = memStatus
			} else {
				boshSummaryMetric.MemErrStat = models.BOSH_STATE_RUNNING
			}

			state := helpers.GetMaxAlarmLevel(alarmStatus)
			if state == "" {
				boshSummaryMetric.State = models.BOSH_STATE_RUNNING
			} else {
				boshSummaryMetric.State = state
			}
		}

		// bosh diskStatus setting
		if boshSummaryMetric.DiskStatus != models.BOSH_STATE_FAIL {
			var diskStatusList []string
			diskStatus := helpers.GetAlarmStatusByServiceName(models.ORIGIN_TYPE_BOSH, models.ALARM_TYPE_DISK, diskUsage, serverThresholds)
			if diskStatus != "" {
				diskStatusList = append(diskStatusList, diskStatus)
				boshSummaryMetric.DiskRootErrStat = diskStatus
			} else {
				boshSummaryMetric.DiskRootErrStat = models.BOSH_STATE_NORMAL
			}

			diskDataStatus := helpers.GetAlarmStatusByServiceName(models.ORIGIN_TYPE_BOSH, models.ALARM_TYPE_DISK, diskDataUsage, serverThresholds)
			if diskDataStatus != "" {
				diskStatusList = append(diskStatusList, diskDataStatus)
				boshSummaryMetric.DiskDataErrStat = diskDataStatus
			} else {
				boshSummaryMetric.DiskDataErrStat = models.BOSH_STATE_NORMAL
			}

			diskState := helpers.GetMaxAlarmLevel(diskStatusList)
			if diskState == "" {
				boshSummaryMetric.DiskStatus = models.BOSH_STATE_NORMAL
			} else {
				boshSummaryMetric.DiskStatus = diskState
			}
		}

		if boshSummaryMetric.State == models.BOSH_STATE_RUNNING && boshSummaryMetric.DiskStatus == models.BOSH_STATE_NORMAL {
			boshSummaryMetric.BoshState = models.BOSH_STATE_RUNNING
		} else if boshSummaryMetric.BoshState != models.BOSH_STATE_FAIL {
			var boshStatusList []string
			boshStatusList = append(boshStatusList, boshSummaryMetric.State)
			if boshSummaryMetric.DiskStatus == models.BOSH_STATE_NORMAL {
				boshStatusList = append(boshStatusList, models.BOSH_STATE_RUNNING)
			} else {
				boshStatusList = append(boshStatusList, boshSummaryMetric.DiskStatus)
			}
			boshSummaryMetric.BoshState = helpers.GetMaxAlarmLevel(boshStatusList)
			boshSummaryMetric.State = boshSummaryMetric.BoshState
		}

		boshSummary.BoshSummaryMetric = boshSummaryMetric
		results = append(results, boshSummary)
	}

	return results, nil
}

func (b *ApBoshService) GetBoshSummaryMetric(boshSummary models.BoshSummary) ([]map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, []models.ErrMessage) {
	var cpuResp, cpuCoreResp, memTotalResp, memFreeResp, diskTotalResp, diskUsedResp, diskDataTotalResp, diskDataUsedResp *client.Response
	var errs []models.ErrMessage
	var err models.ErrMessage
	var wg sync.WaitGroup

	wg.Add(8)
	for i := 0; i < 8; i++ {
		go func(wg *sync.WaitGroup, index int) {
			switch index {
			case 0:
				boshSummary.MetricName = models.MTR_CPU_CORE
				boshSummary.Time = "1m"
				boshSummary.SqlQuery = "select value from bosh_metrics where id = '%s' and time > now() - %s and metricname =~ /%s/ group by metricname order by time desc limit 1;"
				cpuCoreResp, err = dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshSummary(boshSummary)
				if err != nil {
					errs = append(errs, err)
				}
			case 1:
				boshSummary.MetricName = models.MTR_CPU_CORE
				boshSummary.Time = "1m"
				boshSummary.SqlQuery = "select mean(value) as value from bosh_metrics where id = '%s' and time > now() - %s and metricname =~ /%s/ ;"
				cpuResp, err = dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshSummary(boshSummary)
				if err != nil {
					errs = append(errs, err)
				}
			case 2:
				boshSummary.MetricName = models.MTR_MEM_TOTAL
				boshSummary.Time = "1m"
				boshSummary.SqlQuery = "select mean(value) as value from bosh_metrics where id = '%s' and time > now() - %s and metricname = '%s' ;"
				memTotalResp, err = dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshSummary(boshSummary)
				if err != nil {
					errs = append(errs, err)
				}
			case 3:
				boshSummary.MetricName = models.MTR_MEM_FREE
				boshSummary.Time = "1m"
				boshSummary.SqlQuery = "select mean(value) as value from bosh_metrics where id = '%s' and time > now() - %s and metricname = '%s' ;"
				memFreeResp, err = dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshSummary(boshSummary)
				if err != nil {
					errs = append(errs, err)
				}
			case 4:
				boshSummary.MetricName = models.MTR_DISK_TOTAL
				boshSummary.Time = "1m"
				boshSummary.SqlQuery = "select mean(value) as value from bosh_metrics where id = '%s' and time > now() - %s and metricname = '%s' ;"
				diskTotalResp, err = dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshSummary(boshSummary)
				if err != nil {
					errs = append(errs, err)
				}
			case 5:
				boshSummary.MetricName = models.MTR_DISK_USED
				boshSummary.Time = "1m"
				boshSummary.SqlQuery = "select mean(value) as value from bosh_metrics where id = '%s' and time > now() - %s and metricname = '%s' ;"
				diskUsedResp, err = dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshSummary(boshSummary)
				if err != nil {
					errs = append(errs, err)
				}
			case 6:
				boshSummary.MetricName = models.MTR_DISK_DATA_TOTAL
				boshSummary.Time = "1m"
				boshSummary.SqlQuery = "select mean(value) as value from bosh_metrics where id = '%s' and time > now() - %s and metricname = '%s' ;"
				diskDataTotalResp, err = dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshSummary(boshSummary)
				if err != nil {
					errs = append(errs, err)
				}
			case 7:
				boshSummary.MetricName = models.MTR_DISK_DATA_USED
				boshSummary.Time = "1m"
				boshSummary.SqlQuery = "select mean(value) as value from bosh_metrics where id = '%s' and time > now() - %s and metricname = '%s' ;"
				diskDataUsedResp, err = dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshSummary(boshSummary)
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
		/*var returnErrMessage string
		for _, err := range errs {
			returnErrMessage = returnErrMessage + " " + err["Message"].(string)
		}
		errMessage := models.ErrMessage{
			"Message": returnErrMessage,
		}*/
		return nil, nil, nil, nil, nil, nil, nil, nil, errs
	}
	//==========================================================================

	cpuCore, _ := helpers.InfluxConverterToMap(cpuCoreResp)
	memTotal, _ := helpers.InfluxConverter(memTotalResp)
	memFree, _ := helpers.InfluxConverter(memFreeResp)
	diskTotal, _ := helpers.InfluxConverter(diskTotalResp)
	cpuUsage, _ := helpers.InfluxConverter(cpuResp)
	diskUsage, _ := helpers.InfluxConverter(diskUsedResp)
	diskDataTotal, _ := helpers.InfluxConverter(diskDataTotalResp)
	diskDataUsage, _ := helpers.InfluxConverter(diskDataUsedResp)

	return cpuCore, cpuUsage, memTotal, memFree, diskTotal, diskUsage, diskDataTotal, diskDataUsage, nil
}

func (b *ApBoshService) GetBoshProcessByMemory() ([]models.BoshProcess, error) {
	var results []models.BoshProcess

	for _, BoshInfo := range b.BoshInfoList {
		resp, err := dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshProcessByMemory(BoshInfo.UUID)

		if err != nil {
			fmt.Println(err.Error())
			return results, err
		} else {
			valueList, _ := helpers.InfluxConverterToMap(resp)

			var resList []map[string]interface{}

			for z := 0; z < len(valueList); z++ {
				if len(resList) > 0 {
					chk := false
					for y := 0; y < len(resList); y++ {
						if resList[y][models.IFX_MTR_PROC_NAME] == valueList[z][models.IFX_MTR_PROC_NAME] && resList[y][models.IFX_MTR_PROC_PID] == valueList[z][models.IFX_MTR_PROC_PID] {
							chk = true
						}
					}
					if !chk {
						resList = append(resList, valueList[z])
					}
				} else {
					resList = append(resList, valueList[z])
				}
			}

			// mem sort
			sort.Slice(resList, func(i, j int) bool {
				return helpers.TypeChecker_float64(resList[j][models.IFX_MTR_MEM_USAGE]).(float64) < helpers.TypeChecker_float64(resList[i][models.IFX_MTR_MEM_USAGE]).(float64)
			})

			var idx int

			for _, vl := range resList {
				var BoshProcess models.BoshProcess

				BoshProcess.Index = strconv.Itoa(idx + 1)
				BoshProcess.Process = helpers.TypeChecker_string(vl[models.IFX_MTR_PROC_NAME]).(string)
				BoshProcess.Memory = helpers.TypeChecker_float64(vl[models.IFX_MTR_MEM_USAGE]).(float64) / models.MB
				BoshProcess.Pid = strconv.FormatFloat(helpers.TypeChecker_float64(vl[models.IFX_MTR_PROC_PID]).(float64), 'f', 0, 64)
				BoshProcess.Time = time.Unix(vl[models.IFX_MTR_TIME].(int64), 0).Format(time.RFC3339)[0:19]
				BoshProcess.UUID = BoshInfo.UUID
				results = append(results, BoshProcess)
				idx++
				if idx == 5 {
					break
				} //fixed 5row
			}
		}
	}

	return results, nil
}

func (b *ApBoshService) GetBoshChart(boshChart models.BoshChart) ([]models.BoshChart, error) {
	var results []models.BoshChart

	for _, boshInfo := range b.BoshInfoList {
		boshChart.MetricName = models.MTR_CPU_CORE
		cpuUsageResp, err := dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshCpuUsageList(boshChart)
		boshChart.MetricName = models.MTR_CPU_LOAD_1M
		cpuLoad1MResp, err := dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshCpuLoadList(boshChart)
		boshChart.MetricName = models.MTR_CPU_LOAD_5M
		cpuLoad5MResp, err := dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshCpuLoadList(boshChart)
		boshChart.MetricName = models.MTR_CPU_LOAD_15M
		cpuLoad15MResp, err := dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshCpuLoadList(boshChart)

		boshChart.MetricName = models.MTR_MEM_USAGE
		memoryUsageResp, err := dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshMemoryUsageList(boshChart)

		boshChart.MetricName = models.MTR_DISK_USAGE
		diskUsageRootResp, err := dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshDiskUsageList(boshChart)
		boshChart.MetricName = models.MTR_DISK_DATA_USAGE
		diskUsageVcapDataResp, err := dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshDiskUsageList(boshChart)

		boshChart.MetricName = "diskIOStats.\\/\\..*.readBytes"
		diskIoRootReadByteList, err := dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshDiskIoList(boshChart)
		boshChart.MetricName = "diskIOStats.\\/\\..*.writeBytes"
		diskIoRootWriteByteList, err := dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshDiskIoList(boshChart)
		boshChart.MetricName = "diskIOStats.\\/var\\/vcap\\/data\\..*.readBytes"
		diskIoVcapReadByteList, err := dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshDiskIoList(boshChart)
		boshChart.MetricName = "diskIOStats.\\/var\\/vcap\\/data\\..*.writeBytes"
		diskIoVcapWriteByteList, err := dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshDiskIoList(boshChart)

		boshChart.MetricName = "networkIOStats.eth0.bytesSent"
		networkByteSentList, err := dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshNetworkByteList(boshChart)
		boshChart.MetricName = "networkIOStats.eth0.bytesRecv"
		networkByteRecvList, err := dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshNetworkByteList(boshChart)
		boshChart.MetricName = "networkIOStats.eth0.packetSent"
		networkPacketSentList, err := dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshNetworkPacketList(boshChart)
		boshChart.MetricName = "networkIOStats.eth0.packetRecv"
		networkPacketRecvList, err := dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshNetworkPacketList(boshChart)
		boshChart.MetricName = "networkIOStats.eth0.dropIn"
		networkDropInResp, err := dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshNetworkDropList(boshChart)
		boshChart.MetricName = "networkIOStats.eth0.dropOut"
		networkDropOutResp, err := dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshNetworkDropList(boshChart)
		boshChart.MetricName = "networkIOStats.eth0.errIn"
		networkErrorInResp, err := dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshNetworkErrorList(boshChart)
		boshChart.MetricName = "networkIOStats.eth0.errOut"
		networkErrorOutResp, err := dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshNetworkErrorList(boshChart)
		if err != nil {
			fmt.Println(err.Error())
			return results, err
		}
		fmt.Println(boshInfo)

		cpuUsage, _ := helpers.InfluxConverterList(cpuUsageResp, models.RESP_DATA_CPU_NAME)
		cpuLoad1M, _ := helpers.InfluxConverterList(cpuLoad1MResp, models.RESP_DATA_LOAD_1M_NAME)
		cpuLoad5M, _ := helpers.InfluxConverterList(cpuLoad5MResp, models.RESP_DATA_LOAD_5M_NAME)
		cpuLoad15M, _ := helpers.InfluxConverterList(cpuLoad15MResp, models.RESP_DATA_LOAD_5M_NAME)
		memoryUsage, _ := helpers.InfluxConverter4Usage(memoryUsageResp, models.MTR_MEM_USAGE)
		diskRootUsage, _ := helpers.InfluxConverterList(diskUsageRootResp, models.MTR_MEM_USAGE)
		diskVcapDataUsage, _ := helpers.InfluxConverterList(diskUsageVcapDataResp, models.MTR_MEM_USAGE)

		diskIoRootReadByte, _ := helpers.InfluxConverterList(diskIoRootReadByteList, "/-read")
		diskIoRootWriteByte, _ := helpers.InfluxConverterList(diskIoRootWriteByteList, "/-write")
		diskIoVcapReadByte, _ := helpers.InfluxConverterList(diskIoVcapReadByteList, "data-read")
		diskIoVcapWriteByte, _ := helpers.InfluxConverterList(diskIoVcapWriteByteList, "data-write")

		networkByteSent, _ := helpers.InfluxConverterList(networkByteSentList, "sent")
		networkByteRecv, _ := helpers.InfluxConverterList(networkByteRecvList, "recv")
		networkPacketSent, _ := helpers.InfluxConverterList(networkPacketSentList, "in")
		networkPacketRecv, _ := helpers.InfluxConverterList(networkPacketRecvList, "out")
		networkDropIn, _ := helpers.InfluxConverterList(networkDropInResp, "in")
		networkDropOut, _ := helpers.InfluxConverterList(networkDropOutResp, "out")
		networkErrorIn, _ := helpers.InfluxConverterList(networkErrorInResp, "in")
		networkErrorOut, _ := helpers.InfluxConverterList(networkErrorOutResp, "out")

		MetricData := map[string]interface{}{
			"cpuUsage":            cpuUsage,
			"cpuLoad1M":           cpuLoad1M,
			"cpuLoad5M":           cpuLoad5M,
			"cpuLoad15M":          cpuLoad15M,
			"memoryUsage":         memoryUsage,
			"diskRootUsage":       diskRootUsage,
			"diskVcapDataUsage":   diskVcapDataUsage,
			"diskIoRootReadByte":  diskIoRootReadByte,
			"diskIoRootWriteByte": diskIoRootWriteByte,
			"diskIoVcapReadByte":  diskIoVcapReadByte,
			"diskIoVcapWriteByte": diskIoVcapWriteByte,
			"networkByteSent":     networkByteSent,
			"networkByteRecv":     networkByteRecv,
			"networkPacketSent":   networkPacketSent,
			"networkPacketRecv":   networkPacketRecv,
			"networkDropIn":       networkDropIn,
			"networkDropOut":      networkDropOut,
			"networkErrorIn":      networkErrorIn,
			"networkErrorOut":     networkErrorOut,
		}

		var resultBoshChart models.BoshChart
		resultBoshChart.UUID = boshChart.UUID
		resultBoshChart.DefaultTimeRange = boshChart.DefaultTimeRange
		resultBoshChart.TimeRangeFrom = boshChart.TimeRangeFrom
		resultBoshChart.TimeRangeTo = boshChart.TimeRangeTo
		resultBoshChart.GroupBy = boshChart.GroupBy
		resultBoshChart.MetricData = MetricData
		results = append(results, resultBoshChart)
	}

	return results, nil
}

func (b *ApBoshService) GetBoshLog(boshLog models.BoshLog) ([]models.BoshLog, error) {
	var results []models.BoshLog

	for _, boshInfo := range b.BoshInfoList {
		if boshInfo.UUID == boshLog.UUID {
			/**
			Period 파라미터가 존재하면 Period 값으로 DB 조회
			없으면 StartTime, EndTime 파라미터 값으로 DB조회
			*/
			if boshLog.Period == "" {
				/**
				날짜 시간 값을 DB에서 조회할 수 있는 포맷으로 변경
				*/
				if boshLog.StartTime == "" && boshLog.EndTime == "" {
					boshLog.StartTime = fmt.Sprintf("%sT%s", boshLog.TargetDate, "00:00:00")
					boshLog.EndTime = fmt.Sprintf("%sT%s", boshLog.TargetDate, "23:59:59")
				} else if boshLog.StartTime != "" && boshLog.EndTime == "" {
					boshLog.StartTime = fmt.Sprintf("%sT%s", boshLog.TargetDate, boshLog.StartTime)
					boshLog.EndTime = fmt.Sprintf("%sT%s", boshLog.TargetDate, "23:59:59")
				} else if boshLog.StartTime == "" && boshLog.EndTime != "" {
					boshLog.StartTime = fmt.Sprintf("%sT%s", boshLog.TargetDate, "00:00:00")
					boshLog.EndTime = fmt.Sprintf("%sT%s", boshLog.TargetDate, boshLog.EndTime)
				} else {
					boshLog.StartTime = fmt.Sprintf("%sT%s", boshLog.TargetDate, boshLog.StartTime)
					boshLog.EndTime = fmt.Sprintf("%sT%s", boshLog.TargetDate, boshLog.EndTime)
				}
				convert_start_time, _ := time.Parse(time.RFC3339, fmt.Sprintf("%s+09:00", boshLog.StartTime))
				convert_end_time, _ := time.Parse(time.RFC3339, fmt.Sprintf("%s+09:00", boshLog.EndTime))
				startTime := convert_start_time.Unix() - int64(models.GmtTimeGap)*60*60
				endTime := convert_end_time.Unix() - int64(models.GmtTimeGap)*60*60

				// Make RFC3339 date-time strings
				boshLog.StartTime = time.Unix(startTime, 0).Format(time.RFC3339)[0:19] + ".000000000Z"
				boshLog.EndTime = time.Unix(endTime, 0).Format(time.RFC3339)[0:19] + ".000000000Z"
			}
			logResp, err := dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshLog(boshLog)
			if err != nil {
				fmt.Println(err.Error())
				return results, err
			}
			messages, _ := helpers.InfluxConverterList(logResp, "")

			var resultBoshLog models.BoshLog
			resultBoshLog.UUID = boshLog.UUID
			resultBoshLog.LogType = boshLog.LogType
			resultBoshLog.Keyword = boshLog.Keyword
			resultBoshLog.TargetDate = boshLog.TargetDate
			resultBoshLog.Period = boshLog.Period
			resultBoshLog.StartTime = boshLog.StartTime
			resultBoshLog.EndTime = boshLog.EndTime
			resultBoshLog.Messages = messages["metric"]
			results = append(results, resultBoshLog)
		}
	}

	return results, nil
}

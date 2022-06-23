package ap

import (
	"encoding/json"
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	dao "paasta-monitoring-api/dao/api/v1/ap"
	Common "paasta-monitoring-api/dao/api/v1/common"
	"paasta-monitoring-api/helpers"
	models "paasta-monitoring-api/models/api/v1"
	"sort"
	"strconv"
	"sync"
	"time"
)

type ApPaastaService struct {
	DbInfo         *gorm.DB
	InfluxDbClient models.InfluxDbClient
}

func GetApPaastaService(DbInfo *gorm.DB, InfluxDbClient models.InfluxDbClient) *ApPaastaService {
	return &ApPaastaService{
		DbInfo:         DbInfo,
		InfluxDbClient: InfluxDbClient,
	}
}

func (p *ApPaastaService) GetPaastaInfoList() ([]models.Paasta, error) {
	// after Use Database
	results, err := dao.GetPaastaDao(p.DbInfo, p.InfluxDbClient).GetPaastaInfoList()
	if err != nil {
		fmt.Println(err.Error())
		return results, err
	}

	return results, nil
}

func (p *ApPaastaService) GetPaastaOverview(c echo.Context) (models.BoshOverview, error) {
	var result models.BoshOverview
	boshSummary, err := p.GetPaastaSummary(c)
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

func (p *ApPaastaService) GetPaastaSummary(c echo.Context) ([]models.BoshSummary, error) {
	var results []models.BoshSummary
	//var errs []models.ErrMessage
	//var resultsResponse map[string]interface{}

	//임계치 설정정보를 조회한다.
	var params models.AlarmPolicies
	serverThresholds, err := Common.GetAlarmPolicyDao(p.DbInfo).GetAlarmPolicy(params)
	if err != nil {
		fmt.Println(err.Error())
		return results, err
	}
	fmt.Println(serverThresholds)

	/*for _, boshInfo := range p.BoshInfoList {
		var boshSummary models.BoshSummary
		boshSummary.Name = boshInfo.Name
		boshSummary.Ip = boshInfo.Ip
		boshSummary.UUID = boshInfo.UUID
		cpuCoreData, cpuData, memTotData, memFreeData, diskTotalData, diskUsedData, diskDataTotalData, diskDataUsedData, errs := p.GetPaastaSummaryMetric(boshSummary)
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
	}*/

	return results, nil
}

func (p *ApPaastaService) GetPaastaSummaryMetric(boshSummary models.BoshSummary) ([]map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, []models.ErrMessage) {
	var cpuResp, cpuCoreResp, memTotalResp, memFreeResp, diskTotalResp, diskUsedResp, diskDataTotalResp, diskDataUsedResp *client.Response
	var errs []models.ErrMessage
	//var err models.ErrMessage
	var wg sync.WaitGroup

	wg.Add(8)
	/*for i := 0; i < 8; i++ {
		go func(wg *sync.WaitGroup, index int) {
			switch index {
			case 0:
				boshSummary.MetricName = models.MTR_CPU_CORE
				boshSummary.Time = "1m"
				boshSummary.SqlQuery = "select value from bosh_metrics where id = '%s' and time > now() - %s and metricname =~ /%s/ group by metricname order by time desc limit 1;"
				cpuCoreResp, err = dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshSummary(boshSummary)
				if err != nil {
					errs = append(errs, err)
				}
			case 1:
				boshSummary.MetricName = models.MTR_CPU_CORE
				boshSummary.Time = "1m"
				boshSummary.SqlQuery = "select mean(value) as value from bosh_metrics where id = '%s' and time > now() - %s and metricname =~ /%s/ ;"
				cpuResp, err = dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshSummary(boshSummary)
				if err != nil {
					errs = append(errs, err)
				}
			case 2:
				boshSummary.MetricName = models.MTR_MEM_TOTAL
				boshSummary.Time = "1m"
				boshSummary.SqlQuery = "select mean(value) as value from bosh_metrics where id = '%s' and time > now() - %s and metricname = '%s' ;"
				memTotalResp, err = dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshSummary(boshSummary)
				if err != nil {
					errs = append(errs, err)
				}
			case 3:
				boshSummary.MetricName = models.MTR_MEM_FREE
				boshSummary.Time = "1m"
				boshSummary.SqlQuery = "select mean(value) as value from bosh_metrics where id = '%s' and time > now() - %s and metricname = '%s' ;"
				memFreeResp, err = dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshSummary(boshSummary)
				if err != nil {
					errs = append(errs, err)
				}
			case 4:
				boshSummary.MetricName = models.MTR_DISK_TOTAL
				boshSummary.Time = "1m"
				boshSummary.SqlQuery = "select mean(value) as value from bosh_metrics where id = '%s' and time > now() - %s and metricname = '%s' ;"
				diskTotalResp, err = dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshSummary(boshSummary)
				if err != nil {
					errs = append(errs, err)
				}
			case 5:
				boshSummary.MetricName = models.MTR_DISK_USED
				boshSummary.Time = "1m"
				boshSummary.SqlQuery = "select mean(value) as value from bosh_metrics where id = '%s' and time > now() - %s and metricname = '%s' ;"
				diskUsedResp, err = dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshSummary(boshSummary)
				if err != nil {
					errs = append(errs, err)
				}
			case 6:
				boshSummary.MetricName = models.MTR_DISK_DATA_TOTAL
				boshSummary.Time = "1m"
				boshSummary.SqlQuery = "select mean(value) as value from bosh_metrics where id = '%s' and time > now() - %s and metricname = '%s' ;"
				diskDataTotalResp, err = dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshSummary(boshSummary)
				if err != nil {
					errs = append(errs, err)
				}
			case 7:
				boshSummary.MetricName = models.MTR_DISK_DATA_USED
				boshSummary.Time = "1m"
				boshSummary.SqlQuery = "select mean(value) as value from bosh_metrics where id = '%s' and time > now() - %s and metricname = '%s' ;"
				diskDataUsedResp, err = dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshSummary(boshSummary)
				if err != nil {
					errs = append(errs, err)
				}
			default:
				break
			}
			wg.Done()
		}(&wg, i)
	}*/
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

func (p *ApPaastaService) GetPaastaProcessByMemory(paastaProcess models.PaastaProcess) ([]models.PaastaProcess, error) {
	var results []models.PaastaProcess

	response, err := dao.GetPaastaDao(p.DbInfo, p.InfluxDbClient).GetPaastaProcessByMemory(paastaProcess)
	if err != nil {
		fmt.Println(err.Error())
		return results, err
	}
	fmt.Println(response)

	responseConvert, _ := helpers.InfluxConverter(response)
	for _, data := range responseConvert {

		switch data.(type) {
		case []map[string]interface{}:
			dataMap := data.([]map[string]interface{})

			cfProcessMetricsSlice := make([]map[string]interface{}, len(dataMap))
			for idxSlice, sliceValue := range dataMap {
				cfProcessMetricsSlice[idxSlice] = sliceValue
			}

			sort.Slice(cfProcessMetricsSlice, func(i, j int) bool {
				a, _ := strconv.Atoi(cfProcessMetricsSlice[i]["mem_usage"].(string))
				b, _ := strconv.Atoi(cfProcessMetricsSlice[j]["mem_usage"].(string))
				return a > b
			})

			results = make([]models.PaastaProcess, models.TOP_PROCESS_CNT)

			for i := 0; i < models.TOP_PROCESS_CNT; i++ {
				results[i].Index = int64(i) + 1
				results[i].Time = time.Unix(cfProcessMetricsSlice[i]["time"].(int64), 0).Format(time.RFC3339)
				memUsage, _ := strconv.Atoi(cfProcessMetricsSlice[i]["mem_usage"].(string))
				results[i].Memory = int64(helpers.Round(float64(memUsage) / models.MB))
				procPid, _ := cfProcessMetricsSlice[i]["proc_pid"].(json.Number).Int64()
				results[i].Pid = strconv.FormatInt(procPid, 10)
				results[i].Process = cfProcessMetricsSlice[i]["proc_name"].(string)
				results[i].UUID = paastaProcess.UUID
			}
		}
	}

	return results, nil
}

func (p *ApPaastaService) GetPaastaChart(boshChart models.BoshChart) ([]models.BoshChart, error) {
	var results []models.BoshChart

	/*for _, boshInfo := range p.BoshInfoList {
		boshChart.MetricName = models.MTR_CPU_CORE
		cpuUsageResp, err := dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshCpuUsageList(boshChart)
		boshChart.MetricName = models.MTR_CPU_LOAD_1M
		cpuLoad1MResp, err := dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshCpuLoadList(boshChart)
		boshChart.MetricName = models.MTR_CPU_LOAD_5M
		cpuLoad5MResp, err := dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshCpuLoadList(boshChart)
		boshChart.MetricName = models.MTR_CPU_LOAD_15M
		cpuLoad15MResp, err := dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshCpuLoadList(boshChart)

		boshChart.MetricName = models.MTR_MEM_USAGE
		memoryUsageResp, err := dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshMemoryUsageList(boshChart)

		boshChart.MetricName = models.MTR_DISK_USAGE
		diskUsageRootResp, err := dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshDiskUsageList(boshChart)
		boshChart.MetricName = models.MTR_DISK_DATA_USAGE
		diskUsageVcapDataResp, err := dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshDiskUsageList(boshChart)

		boshChart.MetricName = "diskIOStats.\\/\\..*.readBytes"
		diskIoRootReadByteList, err := dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshDiskIoList(boshChart)
		boshChart.MetricName = "diskIOStats.\\/\\..*.writeBytes"
		diskIoRootWriteByteList, err := dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshDiskIoList(boshChart)
		boshChart.MetricName = "diskIOStats.\\/var\\/vcap\\/data\\..*.readBytes"
		diskIoVcapReadByteList, err := dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshDiskIoList(boshChart)
		boshChart.MetricName = "diskIOStats.\\/var\\/vcap\\/data\\..*.writeBytes"
		diskIoVcapWriteByteList, err := dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshDiskIoList(boshChart)

		boshChart.MetricName = "networkIOStats.eth0.bytesSent"
		networkByteSentList, err := dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshNetworkByteList(boshChart)
		boshChart.MetricName = "networkIOStats.eth0.bytesRecv"
		networkByteRecvList, err := dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshNetworkByteList(boshChart)
		boshChart.MetricName = "networkIOStats.eth0.packetSent"
		networkPacketSentList, err := dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshNetworkPacketList(boshChart)
		boshChart.MetricName = "networkIOStats.eth0.packetRecv"
		networkPacketRecvList, err := dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshNetworkPacketList(boshChart)
		boshChart.MetricName = "networkIOStats.eth0.dropIn"
		networkDropInResp, err := dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshNetworkDropList(boshChart)
		boshChart.MetricName = "networkIOStats.eth0.dropOut"
		networkDropOutResp, err := dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshNetworkDropList(boshChart)
		boshChart.MetricName = "networkIOStats.eth0.errIn"
		networkErrorInResp, err := dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshNetworkErrorList(boshChart)
		boshChart.MetricName = "networkIOStats.eth0.errOut"
		networkErrorOutResp, err := dao.GetBoshDao(p.DbInfo, p.InfluxDbClient).GetBoshNetworkErrorList(boshChart)
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
	}*/

	return results, nil
}

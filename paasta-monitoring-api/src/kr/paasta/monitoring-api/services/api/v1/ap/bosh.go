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
	results := b.BoshInfoList
	return results, nil
}

func (b *ApBoshService) GetBoshOverview() ([]models.BoshOverview, error) {
	var results []models.BoshOverview

	return results, nil
}

func (b *ApBoshService) GetBoshSummary() ([]models.BoshSummary, error) {
	var results []models.BoshSummary
	return results, nil
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

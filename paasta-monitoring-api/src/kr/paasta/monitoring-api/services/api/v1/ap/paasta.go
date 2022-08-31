package ap

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	dao "paasta-monitoring-api/dao/api/v1/ap"
	Common "paasta-monitoring-api/dao/api/v1/common"
	"paasta-monitoring-api/helpers"
	models "paasta-monitoring-api/models/api/v1"
	"sort"
	"strconv"
	"strings"
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

func (p *ApPaastaService) GetPaastaInfoList(ctx echo.Context) ([]models.Paasta, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)
	results, err := dao.GetPaastaDao(p.DbInfo, p.InfluxDbClient).GetPaastaInfoList()
	if err != nil {
		logger.Error(err)
		return results, err
	}

	return results, nil
}

func (p *ApPaastaService) GetPaastaOverview(ctx echo.Context) (models.PaastaOverview, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)
	var paastaOverview models.PaastaOverview
	var paastaRequest models.PaastaRequest
	paastaSummary, err := p.GetPaastaSummary(paastaRequest)
	if err != nil {
		logger.Error(err)
		return paastaOverview, err
	}

	// paasta overview
	var overviewTotal, overviewFailed, overviewCritical, overviewWarning = len(paastaSummary), 0, 0, 0
	for _, value := range paastaSummary {
		if value.PaastaSummaryMetric.CpuState == models.STATE_FAILED || value.PaastaSummaryMetric.MemoryState == models.STATE_FAILED || value.PaastaSummaryMetric.DiskState == models.STATE_FAILED {
			overviewFailed++
		} else if value.PaastaSummaryMetric.CpuState == models.STATE_CRITICAL || value.PaastaSummaryMetric.MemoryState == models.STATE_CRITICAL || value.PaastaSummaryMetric.DiskState == models.STATE_CRITICAL {
			overviewCritical++
		} else if value.PaastaSummaryMetric.CpuState == models.STATE_WARNING || value.PaastaSummaryMetric.MemoryState == models.STATE_WARNING || value.PaastaSummaryMetric.DiskState == models.STATE_WARNING {
			overviewWarning++
		}
	}
	paastaOverview.Total = strconv.Itoa(overviewTotal)
	paastaOverview.Failed = strconv.Itoa(overviewFailed)
	paastaOverview.Critical = strconv.Itoa(overviewCritical)
	paastaOverview.Warning = strconv.Itoa(overviewWarning)
	paastaOverview.Running = strconv.Itoa(overviewTotal - overviewFailed - overviewCritical - overviewWarning)

	return paastaOverview, nil
}


func (p *ApPaastaService) GetPaastaSummary(paastaRequest models.PaastaRequest) ([]models.PaastaSummary, error) {
	var results []models.PaastaSummary

	var overviewTotal = 0
	var thresholdWarningCpu, thresholdCriticalCpu float64 = 0, 0
	var thresholdWarningMemory, thresholdCriticalMemory float64 = 0, 0
	var thresholdWarningDisk, thresholdCriticalDisk float64 = 0, 0

	//임계치 설정정보를 조회한다.
	var params models.AlarmPolicies
	policies, err := Common.GetAlarmPolicyDao(p.DbInfo).GetAlarmPolicy(params)
	for _, v := range policies {
		if models.ORIGIN_TYPE_PAAS == v.OriginType {
			switch v.AlarmType {
			case models.ALARM_TYPE_CPU:
				thresholdWarningCpu = float64(v.WarningThreshold)
				thresholdCriticalCpu = float64(v.CriticalThreshold)
			case models.ALARM_TYPE_MEMORY:
				thresholdWarningMemory = float64(v.WarningThreshold)
				thresholdCriticalMemory = float64(v.CriticalThreshold)
			case models.ALARM_TYPE_DISK:
				thresholdWarningDisk = float64(v.WarningThreshold)
				thresholdCriticalDisk = float64(v.CriticalThreshold)
			}
		}
	}
	if err != nil {
		fmt.Println(err.Error())
		return results, err
	}
	fmt.Println(policies)

	vms, err := dao.GetPaastaDao(p.DbInfo, p.InfluxDbClient).GetPaastaInfoList()
	if err != nil {
		fmt.Println(err.Error())
		return results, err
	}
	overviewTotal = len(vms)
	paasVms := make([]models.PaastaSummary, overviewTotal)
	conditionalPaasVms := make([]models.PaastaSummary, 0, overviewTotal)

	for vmIdx, vmValue := range vms {

		paasVms[vmIdx].Name = vmValue.Name
		paasVms[vmIdx].Address = vmValue.Ip
		paasVms[vmIdx].PaastaSummaryMetric.State = models.STATE_RUNNING

		resp, _ := dao.GetPaastaDao(p.DbInfo, p.InfluxDbClient).GetPaastaCfMetrics(vmValue.Ip)

		result, _ := helpers.InfluxConverter(resp)

		for _, data := range result {

			switch data.(type) {
			case []map[string]interface{}:
				dataMap := data.([]map[string]interface{})

				core := 0
				var sumCpuUsed float64 = 0
				var memoryTotal int64 = 0
				var memoryFree int64 = 0
				var wg sync.WaitGroup
				wg.Add(len(dataMap))
				for _, v := range dataMap {
					go func(wg *sync.WaitGroup, metric map[string]interface{}) {
						defer wg.Done()

						switch metric["metricname"] {
						case models.METRIC_NAME_TOTAL_MEMORY:
							memoryTotal, _ = metric["value"].(json.Number).Int64()
							paasVms[vmIdx].PaastaSummaryMetric.TotalMemory = memoryTotal / models.MB
							paasVms[vmIdx].UUID = metric["id"].(string)
						case models.METRIC_NAME_FREE_MEMORY:
							memoryFree, _ = metric["value"].(json.Number).Int64()
						case models.METRIC_NAME_TOTAL_DISK_ROOT:
							value, _ := metric["value"].(json.Number).Int64()
							paasVms[vmIdx].PaastaSummaryMetric.TotalDisk = value / models.MB
						case models.METRIC_NAME_TOTAL_DISK_VCAP:
							value, _ := metric["value"].(json.Number).Int64()
							paasVms[vmIdx].PaastaSummaryMetric.DataDisk = value / models.MB
						case models.METRIC_NAME_DISK_ROOT_USAGE:
							paasVms[vmIdx].PaastaSummaryMetric.TotalDiskUsage, _ = metric["value"].(json.Number).Float64()
							if paasVms[vmIdx].PaastaSummaryMetric.TotalDiskUsage > thresholdCriticalDisk {
								paasVms[vmIdx].PaastaSummaryMetric.TotalDiskState, paasVms[vmIdx].PaastaSummaryMetric.DiskState, paasVms[vmIdx].PaastaSummaryMetric.State = models.STATE_CRITICAL, models.STATE_CRITICAL, models.STATE_CRITICAL
							} else if paasVms[vmIdx].PaastaSummaryMetric.TotalDiskUsage > thresholdWarningDisk {
								paasVms[vmIdx].PaastaSummaryMetric.TotalDiskState, paasVms[vmIdx].PaastaSummaryMetric.DiskState, paasVms[vmIdx].PaastaSummaryMetric.State = models.STATE_WARNING, models.STATE_WARNING, models.STATE_WARNING
							} else if paasVms[vmIdx].PaastaSummaryMetric.TotalDiskUsage == 0 {
								paasVms[vmIdx].PaastaSummaryMetric.TotalDiskState, paasVms[vmIdx].PaastaSummaryMetric.DiskState, paasVms[vmIdx].PaastaSummaryMetric.State = models.STATE_FAILED, models.STATE_FAILED, models.STATE_FAILED
							} else {
								paasVms[vmIdx].PaastaSummaryMetric.TotalDiskState, paasVms[vmIdx].PaastaSummaryMetric.DiskState = models.DISK_STATE_NORMAL, models.DISK_STATE_NORMAL
							}
						case models.METRIC_NAME_DISK_VCAP_USAGE:
							paasVms[vmIdx].PaastaSummaryMetric.DataDiskUsage, _ = metric["value"].(json.Number).Float64()
							if paasVms[vmIdx].PaastaSummaryMetric.DataDiskUsage > thresholdCriticalDisk {
								paasVms[vmIdx].PaastaSummaryMetric.DataDiskState, paasVms[vmIdx].PaastaSummaryMetric.DiskState, paasVms[vmIdx].PaastaSummaryMetric.State = models.STATE_CRITICAL, models.STATE_CRITICAL, models.STATE_CRITICAL
							} else if paasVms[vmIdx].PaastaSummaryMetric.DataDiskUsage > thresholdWarningDisk {
								paasVms[vmIdx].PaastaSummaryMetric.DataDiskState, paasVms[vmIdx].PaastaSummaryMetric.DiskState, paasVms[vmIdx].PaastaSummaryMetric.State = models.STATE_WARNING, models.STATE_WARNING, models.STATE_WARNING
							} else if paasVms[vmIdx].PaastaSummaryMetric.DataDiskUsage == 0 {
								paasVms[vmIdx].PaastaSummaryMetric.DataDiskState, paasVms[vmIdx].PaastaSummaryMetric.DiskState, paasVms[vmIdx].PaastaSummaryMetric.State = models.STATE_FAILED, models.STATE_FAILED, models.STATE_FAILED
							} else {
								paasVms[vmIdx].PaastaSummaryMetric.DataDiskState, paasVms[vmIdx].PaastaSummaryMetric.DiskState = models.DISK_STATE_NORMAL, models.DISK_STATE_NORMAL
							}
						default:
							if strings.Contains(metric["metricname"].(string), models.METRIC_NAME_CPU_CORE_PREFIX) {
								core++
								value, _ := metric["value"].(json.Number).Float64()
								sumCpuUsed += value
							}
						}
					}(&wg, v)
				}
				wg.Wait()

				paasVms[vmIdx].PaastaSummaryMetric.MemoryUsage = 100.0 - (float64(memoryFree) / float64(memoryTotal) * 100.0)
				if paasVms[vmIdx].PaastaSummaryMetric.MemoryUsage > thresholdCriticalMemory {
					paasVms[vmIdx].PaastaSummaryMetric.MemoryState, paasVms[vmIdx].PaastaSummaryMetric.State = models.STATE_CRITICAL, models.STATE_CRITICAL
				} else if paasVms[vmIdx].PaastaSummaryMetric.MemoryUsage > thresholdWarningMemory {
					paasVms[vmIdx].PaastaSummaryMetric.MemoryState, paasVms[vmIdx].PaastaSummaryMetric.State = models.STATE_WARNING, models.STATE_WARNING
				} else if paasVms[vmIdx].PaastaSummaryMetric.MemoryUsage == 0 {
					paasVms[vmIdx].PaastaSummaryMetric.MemoryState, paasVms[vmIdx].PaastaSummaryMetric.State = models.STATE_FAILED, models.STATE_FAILED
				} else {
					paasVms[vmIdx].PaastaSummaryMetric.MemoryState = models.STATE_RUNNING
				}

				paasVms[vmIdx].PaastaSummaryMetric.Core = strconv.Itoa(core)
				if core > 0 {
					paasVms[vmIdx].PaastaSummaryMetric.CpuUsage = sumCpuUsed / float64(core)
					if paasVms[vmIdx].PaastaSummaryMetric.CpuUsage > thresholdCriticalCpu {
						paasVms[vmIdx].PaastaSummaryMetric.CpuState, paasVms[vmIdx].PaastaSummaryMetric.State = models.STATE_CRITICAL, models.STATE_CRITICAL
					} else if paasVms[vmIdx].PaastaSummaryMetric.CpuUsage > thresholdWarningCpu {
						paasVms[vmIdx].PaastaSummaryMetric.CpuState, paasVms[vmIdx].PaastaSummaryMetric.State = models.STATE_WARNING, models.STATE_WARNING
					} else {
						paasVms[vmIdx].PaastaSummaryMetric.CpuState = models.STATE_RUNNING
					}
				}

				results = append(results, paasVms[vmIdx])
			}
		}

		/* conditional search */
		if (paastaRequest.Ip != "" && strings.Contains(paasVms[vmIdx].Address, paastaRequest.Ip)) ||
			(paastaRequest.Name != "" && strings.Contains(paasVms[vmIdx].Name, paastaRequest.Name) ||
				paastaRequest.Status != "" && paasVms[vmIdx].PaastaSummaryMetric.State == paastaRequest.Status) {
			conditionalPaasVms = append(conditionalPaasVms, paasVms[vmIdx])
		}
	}

	return results, nil
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

func (p *ApPaastaService) GetPaastaChart(paastaChart models.PaastaChart) (models.PaastaChart, error) {
	var results models.PaastaChart

	// CPU

	paastaChart.MetricName = models.METRIC_NAME_CPU_CORE_PREFIX
	paastaChart.IsLikeQuery = true
	cpuUsageResponse, err := dao.GetPaastaDao(p.DbInfo, p.InfluxDbClient).GetPaastaCommonUsageByTime(paastaChart)
	paastaChart.MetricName = models.METRIC_NAME_CPU_LOAD_AVG_01_MIN
	cpuLoad1MResponse, err := dao.GetPaastaDao(p.DbInfo, p.InfluxDbClient).GetPaastaCommonUsageByTime(paastaChart)
	paastaChart.MetricName = models.METRIC_NAME_CPU_LOAD_AVG_05_MIN
	cpuLoad5MResponse, err := dao.GetPaastaDao(p.DbInfo, p.InfluxDbClient).GetPaastaCommonUsageByTime(paastaChart)
	paastaChart.MetricName = models.METRIC_NAME_CPU_LOAD_AVG_15_MIN
	cpuLoad15MResponse, err := dao.GetPaastaDao(p.DbInfo, p.InfluxDbClient).GetPaastaCommonUsageByTime(paastaChart)

	// Memory
	//paastaChart.MetricName = models.MTR_MEM_USAGE
	memoryUsageResponse, err := dao.GetPaastaDao(p.DbInfo, p.InfluxDbClient).GetPaastaMemoryUsage(paastaChart)

	// Disk

	paastaChart.IsLikeQuery = false
	paastaChart.MetricName = models.METRIC_NAME_DISK_ROOT_USAGE
	diskUsageRootResponse, err := dao.GetPaastaDao(p.DbInfo, p.InfluxDbClient).GetPaastaCommonUsageByTime(paastaChart)
	paastaChart.MetricName = models.METRIC_NAME_DISK_VCAP_USAGE
	diskUsageVcapDataResponse, err := dao.GetPaastaDao(p.DbInfo, p.InfluxDbClient).GetPaastaCommonUsageByTime(paastaChart)
	paastaChart.IsLikeQuery = true
	paastaChart.IsRespondKb = true
	paastaChart.IsNonNegativeDerivative = true
	paastaChart.MetricName = models.METRIC_NAME_DISK_IO_ROOT_READ_BYTES
	diskIoRootReadByteResponse, err := dao.GetPaastaDao(p.DbInfo, p.InfluxDbClient).GetPaastaCommonUsageByTime(paastaChart)
	paastaChart.MetricName = models.METRIC_NAME_DISK_IO_ROOT_WRITE_BYTES
	diskIoRootWriteByteResponse, err := dao.GetPaastaDao(p.DbInfo, p.InfluxDbClient).GetPaastaCommonUsageByTime(paastaChart)
	paastaChart.MetricName = models.METRIC_NAME_DISK_IO_VCAP_READ_BYTES
	diskIoVcapReadByteResponse, err := dao.GetPaastaDao(p.DbInfo, p.InfluxDbClient).GetPaastaCommonUsageByTime(paastaChart)
	paastaChart.MetricName = models.METRIC_NAME_DISK_IO_VCAP_WRITE_BYTES
	diskIoVcapWriteByteResponse, err := dao.GetPaastaDao(p.DbInfo, p.InfluxDbClient).GetPaastaCommonUsageByTime(paastaChart)

	// Network

	paastaChart.IsLikeQuery = false
	paastaChart.MetricName = models.METRIC_NETWORK_IO_BYTES_SENT
	networkByteSentResponse, err := dao.GetPaastaDao(p.DbInfo, p.InfluxDbClient).GetPaastaCommonUsageByTime(paastaChart)
	paastaChart.MetricName = models.METRIC_NETWORK_IO_BYTES_RECV
	networkByteRecvResponse, err := dao.GetPaastaDao(p.DbInfo, p.InfluxDbClient).GetPaastaCommonUsageByTime(paastaChart)
	paastaChart.MetricName = models.METRIC_NETWORK_IO_BYTES_SENT
	networkPacketSentResponse, err := dao.GetPaastaDao(p.DbInfo, p.InfluxDbClient).GetPaastaCommonUsageByTime(paastaChart)
	paastaChart.MetricName = models.METRIC_NETWORK_IO_BYTES_RECV
	networkPacketRecvResponse, err := dao.GetPaastaDao(p.DbInfo, p.InfluxDbClient).GetPaastaCommonUsageByTime(paastaChart)
	paastaChart.MetricName = models.METRIC_NETWORK_IO_DROP_IN
	networkDropInResponse, err := dao.GetPaastaDao(p.DbInfo, p.InfluxDbClient).GetPaastaCommonUsageByTime(paastaChart)
	paastaChart.MetricName = models.METRIC_NETWORK_IO_DROP_OUT
	networkDropOutResponse, err := dao.GetPaastaDao(p.DbInfo, p.InfluxDbClient).GetPaastaCommonUsageByTime(paastaChart)
	paastaChart.MetricName = models.METRIC_NETWORK_IO_ERR_IN
	networkErrorInResponse, err := dao.GetPaastaDao(p.DbInfo, p.InfluxDbClient).GetPaastaCommonUsageByTime(paastaChart)
	paastaChart.MetricName = models.METRIC_NETWORK_IO_ERR_OUT
	networkErrorOutResponse, err := dao.GetPaastaDao(p.DbInfo, p.InfluxDbClient).GetPaastaCommonUsageByTime(paastaChart)

	// CPU

	cpuUsage, _ := helpers.InfluxConverterList(cpuUsageResponse, models.RESP_DATA_CPU_NAME)
	cpuLoad1M, _ := helpers.InfluxConverterList(cpuLoad1MResponse, models.RESP_DATA_LOAD_1M_NAME)
	cpuLoad5M, _ := helpers.InfluxConverterList(cpuLoad5MResponse, models.RESP_DATA_LOAD_5M_NAME)
	cpuLoad15M, _ := helpers.InfluxConverterList(cpuLoad15MResponse, models.RESP_DATA_LOAD_5M_NAME)

	// Memory
	memoryUsage, _ := helpers.InfluxConverter4Usage(memoryUsageResponse, models.MTR_MEM_USAGE)

	// Disk

	diskRootUsage, _ := helpers.InfluxConverterList(diskUsageRootResponse, models.MTR_MEM_USAGE)
	diskVcapDataUsage, _ := helpers.InfluxConverterList(diskUsageVcapDataResponse, models.MTR_MEM_USAGE)
	diskIoRootReadByteUsage, _ := helpers.InfluxConverterList(diskIoRootReadByteResponse, "/-read")
	diskIoRootWriteByteUsage, _ := helpers.InfluxConverterList(diskIoRootWriteByteResponse, "/-write")
	diskIoVcapReadByteUsage, _ := helpers.InfluxConverterList(diskIoVcapReadByteResponse, "data-read")
	diskIoVcapWriteByteUsage, _ := helpers.InfluxConverterList(diskIoVcapWriteByteResponse, "data-write")

	// Network

	networkByteSent, _ := helpers.InfluxConverterList(networkByteSentResponse, "sent")
	networkByteRecv, _ := helpers.InfluxConverterList(networkByteRecvResponse, "recv")
	networkPacketSent, _ := helpers.InfluxConverterList(networkPacketSentResponse, "in")
	networkPacketRecv, _ := helpers.InfluxConverterList(networkPacketRecvResponse, "out")
	networkDropIn, _ := helpers.InfluxConverterList(networkDropInResponse, "in")
	networkDropOut, _ := helpers.InfluxConverterList(networkDropOutResponse, "out")
	networkErrorIn, _ := helpers.InfluxConverterList(networkErrorInResponse, "in")
	networkErrorOut, _ := helpers.InfluxConverterList(networkErrorOutResponse, "out")

	if err != nil {
		fmt.Println(err.Error())
		return results, err
	}

	MetricData := map[string]interface{}{
		"cpuUsage":                 cpuUsage,
		"cpuLoad1M":                cpuLoad1M,
		"cpuLoad5M":                cpuLoad5M,
		"cpuLoad15M":               cpuLoad15M,
		"memoryUsage":              memoryUsage,
		"diskRootUsage":            diskRootUsage,
		"diskVcapDataUsage":        diskVcapDataUsage,
		"diskIoRootReadByteUsage":  diskIoRootReadByteUsage,
		"diskIoRootWriteByteUsage": diskIoRootWriteByteUsage,
		"diskIoVcapReadByteUsage":  diskIoVcapReadByteUsage,
		"diskIoVcapWriteByteUsage": diskIoVcapWriteByteUsage,
		"networkByteSent":          networkByteSent,
		"networkByteRecv":          networkByteRecv,
		"networkPacketSent":        networkPacketSent,
		"networkPacketRecv":        networkPacketRecv,
		"networkDropIn":            networkDropIn,
		"networkDropOut":           networkDropOut,
		"networkErrorIn":           networkErrorIn,
		"networkErrorOut":          networkErrorOut,
	}
	paastaChart.MetricData = MetricData

	return paastaChart, nil
}

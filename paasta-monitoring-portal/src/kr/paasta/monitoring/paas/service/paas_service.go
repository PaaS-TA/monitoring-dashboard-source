package service

import (
	"encoding/json"
	"github.com/cihub/seelog"
	"github.com/cloudfoundry-community/gogobosh"
	"github.com/influxdata/influxdb1-client/v2"
	"github.com/jinzhu/gorm"
	"monitoring-portal/paas/dao"
	"monitoring-portal/paas/model"
	"monitoring-portal/paas/util"
	"monitoring-portal/utils"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type PaasService struct {
	txn          *gorm.DB
	influxClient client.Client
	databases    model.Databases
	boshClient   *gogobosh.Client
}

var logger seelog.LoggerInterface

func GetPaasService(txn *gorm.DB, influxClient client.Client, databases model.Databases, boshClent *gogobosh.Client) *PaasService {
	return &PaasService{
		txn:          txn,
		influxClient: influxClient,
		databases:    databases,
		boshClient:   boshClent,
	}
}

func (p *PaasService) GetPaasOverview(request model.PaasRequest) (paasOverview model.PaasOverview, err model.ErrMessage) {

	paasSummary, err := GetPaasService(p.txn, p.influxClient, p.databases, p.boshClient).GetPaasVms(request)
	if err != nil {
		logger.Error(err)
	}

	return paasSummary.PaasOverview, err
}

func (p *PaasService) GetPaasOverviewStatus(request model.PaasRequest) (paasOverviewStatus model.PaasOverviewStatus, err model.ErrMessage) {

	paasSummary, err := GetPaasService(p.txn, p.influxClient, p.databases, p.boshClient).GetPaasVms(request)
	if err != nil {
		logger.Error(err)
	}

	for _, v := range paasSummary.Data {
		if v.State == request.Status {
			paasOverviewStatus.Data = append(paasOverviewStatus.Data, v)
		}
	}

	return paasOverviewStatus, err
}

func (p *PaasService) GetPaasSummary(request model.PaasRequest) (model.PaasSummary, model.ErrMessage) {

	paasSummary, err := GetPaasService(p.txn, p.influxClient, p.databases, p.boshClient).GetPaasVms(request)
	if err != nil {
		return paasSummary, err
	}

	return paasSummary, err
}

func (p *PaasService) GetPaasVms(request model.PaasRequest) (paasSummary model.PaasSummary, errMsg model.ErrMessage) {

	ds := p.databases.PaastaDatabase

	var overviewTotal, overviewFailed, overviewCritical, overviewWarning = 0, 0, 0, 0
	var thresholdWarningCpu, thresholdCriticalCpu float64 = 0, 0
	var thresholdWarningMemory, thresholdCriticalMemory float64 = 0, 0
	var thresholdWarningDisk, thresholdCriticalDisk float64 = 0, 0

	// get state criteria from alarm_policies.
	policies, err := dao.GetAlarmPolicyDao(p.txn).GetAlarmPolicyList()
	for _, v := range policies {
		if model.ORIGIN_TYPE_PAAS == v.OriginType {
			switch v.AlarmType {
			case model.ALARM_TYPE_CPU:
				thresholdWarningCpu = float64(v.WarningThreshold)
				thresholdCriticalCpu = float64(v.CriticalThreshold)
			case model.ALARM_TYPE_MEMORY:
				thresholdWarningMemory = float64(v.WarningThreshold)
				thresholdCriticalMemory = float64(v.CriticalThreshold)
			case model.ALARM_TYPE_DISK:
				thresholdWarningDisk = float64(v.WarningThreshold)
				thresholdCriticalDisk = float64(v.CriticalThreshold)
			}
		}
	}

	vms, err := dao.GetPaasDao(p.txn, p.influxClient, ds).GetPaasVms()
	if err != nil {
		return paasSummary, err
	}
	overviewTotal = len(vms)

	paasVms := make([]model.PaasVm, overviewTotal)
	conditionalPaasVms := make([]model.PaasVm, 0, overviewTotal)

	for vmIdx, vmValue := range vms {

		paasVms[vmIdx].Name = vmValue.Name
		paasVms[vmIdx].Address = vmValue.Ip
		paasVms[vmIdx].State = model.STATE_RUNNING

		resp, _ := dao.GetPaasDao(p.txn, p.influxClient, ds).GetPaasCfMetrics(vmValue.Ip)

		result, _ := util.GetResponseConverter().InfluxConverter(resp, "")

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
						case model.METRIC_NAME_TOTAL_MEMORY:
							memoryTotal, _ = metric["value"].(json.Number).Int64()
							paasVms[vmIdx].TotalMemory = memoryTotal / model.MB
							paasVms[vmIdx].Id = metric["id"].(string)
						/*case model.METRIC_NAME_MEMORY_USAGE:
						paasVms[vmIdx].MemoryUsage, _ = metric["value"].(json.Number).Float64()
						if paasVms[vmIdx].MemoryUsage > thresholdCriticalMemory {
							paasVms[vmIdx].MemoryState, paasVms[vmIdx].State = model.STATE_CRITICAL, model.STATE_CRITICAL
						} else if paasVms[vmIdx].MemoryUsage > thresholdWarningMemory {
							paasVms[vmIdx].MemoryState, paasVms[vmIdx].State = model.STATE_WARNING, model.STATE_WARNING
						} else if paasVms[vmIdx].MemoryUsage == 0 {
							paasVms[vmIdx].MemoryState, paasVms[vmIdx].State = model.STATE_FAILED, model.STATE_FAILED
						} else {
							paasVms[vmIdx].MemoryState = model.STATE_RUNNING
						}*/
						case model.METRIC_NAME_FREE_MEMORY:
							memoryFree, _ = metric["value"].(json.Number).Int64()
						case model.METRIC_NAME_TOTAL_DISK_ROOT:
							value, _ := metric["value"].(json.Number).Int64()
							paasVms[vmIdx].TotalDisk = value / model.MB
						case model.METRIC_NAME_TOTAL_DISK_VCAP:
							value, _ := metric["value"].(json.Number).Int64()
							paasVms[vmIdx].DataDisk = value / model.MB
						case model.METRIC_NAME_DISK_ROOT_USAGE:
							paasVms[vmIdx].TotalDiskUsage, _ = metric["value"].(json.Number).Float64()
							if paasVms[vmIdx].TotalDiskUsage > thresholdCriticalDisk {
								paasVms[vmIdx].TotalDiskState, paasVms[vmIdx].DiskState, paasVms[vmIdx].State = model.STATE_CRITICAL, model.STATE_CRITICAL, model.STATE_CRITICAL
							} else if paasVms[vmIdx].TotalDiskUsage > thresholdWarningDisk {
								paasVms[vmIdx].TotalDiskState, paasVms[vmIdx].DiskState, paasVms[vmIdx].State = model.STATE_WARNING, model.STATE_WARNING, model.STATE_WARNING
							} else if paasVms[vmIdx].TotalDiskUsage == 0 {
								paasVms[vmIdx].TotalDiskState, paasVms[vmIdx].DiskState, paasVms[vmIdx].State = model.STATE_FAILED, model.STATE_FAILED, model.STATE_FAILED
							} else {
								paasVms[vmIdx].TotalDiskState, paasVms[vmIdx].DiskState = model.DISK_STATE_NORMAL, model.DISK_STATE_NORMAL
							}
						case model.METRIC_NAME_DISK_VCAP_USAGE:
							paasVms[vmIdx].DataDiskUsage, _ = metric["value"].(json.Number).Float64()
							if paasVms[vmIdx].DataDiskUsage > thresholdCriticalDisk {
								paasVms[vmIdx].DataDiskState, paasVms[vmIdx].DiskState, paasVms[vmIdx].State = model.STATE_CRITICAL, model.STATE_CRITICAL, model.STATE_CRITICAL
							} else if paasVms[vmIdx].DataDiskUsage > thresholdWarningDisk {
								paasVms[vmIdx].DataDiskState, paasVms[vmIdx].DiskState, paasVms[vmIdx].State = model.STATE_WARNING, model.STATE_WARNING, model.STATE_WARNING
							} else if paasVms[vmIdx].DataDiskUsage == 0 {
								paasVms[vmIdx].DataDiskState, paasVms[vmIdx].DiskState, paasVms[vmIdx].State = model.STATE_FAILED, model.STATE_FAILED, model.STATE_FAILED
							} else {
								paasVms[vmIdx].DataDiskState, paasVms[vmIdx].DiskState = model.DISK_STATE_NORMAL, model.DISK_STATE_NORMAL
							}
						default:
							if strings.Contains(metric["metricname"].(string), model.METRIC_NAME_CPU_CORE_PREFIX) {
								core++
								value, _ := metric["value"].(json.Number).Float64()
								sumCpuUsed += value
							}
						}
					}(&wg, v)
				}
				wg.Wait()

				paasVms[vmIdx].MemoryUsage = 100.0 - (float64(memoryFree) / float64(memoryTotal) * 100.0)
				if paasVms[vmIdx].MemoryUsage > thresholdCriticalMemory {
					paasVms[vmIdx].MemoryState, paasVms[vmIdx].State = model.STATE_CRITICAL, model.STATE_CRITICAL
				} else if paasVms[vmIdx].MemoryUsage > thresholdWarningMemory {
					paasVms[vmIdx].MemoryState, paasVms[vmIdx].State = model.STATE_WARNING, model.STATE_WARNING
				} else if paasVms[vmIdx].MemoryUsage == 0 {
					paasVms[vmIdx].MemoryState, paasVms[vmIdx].State = model.STATE_FAILED, model.STATE_FAILED
				} else {
					paasVms[vmIdx].MemoryState = model.STATE_RUNNING
				}

				paasVms[vmIdx].Core = strconv.Itoa(core)
				if core > 0 {
					paasVms[vmIdx].CpuUsage = sumCpuUsed / float64(core)
					if paasVms[vmIdx].CpuUsage > thresholdCriticalCpu {
						paasVms[vmIdx].CpuState, paasVms[vmIdx].State = model.STATE_CRITICAL, model.STATE_CRITICAL
					} else if paasVms[vmIdx].CpuUsage > thresholdWarningCpu {
						paasVms[vmIdx].CpuState, paasVms[vmIdx].State = model.STATE_WARNING, model.STATE_WARNING
					} else {
						paasVms[vmIdx].CpuState = model.STATE_RUNNING
					}
				}
			}
		}

		/* conditional search */
		if (request.Ip != "" && strings.Contains(paasVms[vmIdx].Address, request.Ip)) ||
			(request.Name != "" && strings.Contains(paasVms[vmIdx].Name, request.Name) ||
				request.Status != "" && paasVms[vmIdx].State == request.Status) {
			conditionalPaasVms = append(conditionalPaasVms, paasVms[vmIdx])
		}
	}

	if request.Ip != "" || request.Name != "" || request.Status != "" {
		paasSummary.Data = conditionalPaasVms
	} else {
		paasSummary.Data = paasVms
	}
	paasSummary.TotalCount = len(paasSummary.Data)
	paasSummary.PageItems = request.PagingReq.PageItem

	// data pagination
	if request.PagingReq.PageItem != 0 && request.PagingReq.PageIndex != 0 {
		offset := (request.PagingReq.PageIndex - 1) * request.PagingReq.PageItem
		limit := offset + request.PagingReq.PageItem
		if offset >= len(paasSummary.Data) {
			// invalid request
			paasSummary.Data = nil
		} else {
			if limit > len(paasSummary.Data) {
				limit = len(paasSummary.Data)
			}
			paasSummary.Data = paasSummary.Data[offset:limit]
		}
	}

	// get overview stats.
	for _, vm := range paasVms {
		if vm.CpuState == model.STATE_FAILED || vm.MemoryState == model.STATE_FAILED || vm.DiskState == model.STATE_FAILED {
			overviewFailed++
		} else if vm.CpuState == model.STATE_CRITICAL || vm.MemoryState == model.STATE_CRITICAL || vm.DiskState == model.STATE_CRITICAL {
			overviewCritical++
		} else if vm.CpuState == model.STATE_WARNING || vm.MemoryState == model.STATE_WARNING || vm.DiskState == model.STATE_WARNING {
			overviewWarning++
		}
	}
	paasSummary.PaasOverview.Total = strconv.Itoa(overviewTotal)
	paasSummary.PaasOverview.Failed = strconv.Itoa(overviewFailed)
	paasSummary.PaasOverview.Critical = strconv.Itoa(overviewCritical)
	paasSummary.PaasOverview.Warning = strconv.Itoa(overviewWarning)
	paasSummary.PaasOverview.Running = strconv.Itoa(overviewTotal - overviewFailed - overviewCritical - overviewWarning)

	return paasSummary, errMsg
}

func (p *PaasService) GetPaasTopProcessMemory(request model.PaasRequest) ([]model.PaasProcessUsage, model.ErrMessage) {

	ds := p.databases.PaastaDatabase

	resp, err := dao.GetPaasDao(p.txn, p.influxClient, ds).GetPaasTopProcessList(request)

	var topProcessUsage []model.PaasProcessUsage

	if err != nil {
		return topProcessUsage, err

	}

	result, _ := util.GetResponseConverter().InfluxConverter(resp, "")

	for _, data := range result {

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

			topProcessUsage = make([]model.PaasProcessUsage, model.TOP_PROCESS_CNT)

			for i := 0; i < model.TOP_PROCESS_CNT; i++ {
				topProcessUsage[i].Index = int64(i) + 1
				topProcessUsage[i].Time = time.Unix(cfProcessMetricsSlice[i]["time"].(int64), 0).Format(time.RFC3339)
				memUsage, _ := strconv.Atoi(cfProcessMetricsSlice[i]["mem_usage"].(string))
				topProcessUsage[i].Memory = int64(utils.Round(float64(memUsage) / model.MB))
				procPid, _ := cfProcessMetricsSlice[i]["proc_pid"].(json.Number).Int64()
				topProcessUsage[i].Pid = strconv.FormatInt(procPid, 10)
				topProcessUsage[i].Process = cfProcessMetricsSlice[i]["proc_name"].(string)
			}
		}
	}

	return topProcessUsage, err
}

func (p *PaasService) GetPaasMetricStats(request model.PaasRequest) (result []map[string]interface{}, err model.ErrMessage) {

	ds := p.databases.PaastaDatabase
	for _, v := range request.Args.([]model.MetricArg) {
		request.MetricName = v.Name
		daoResult, err := dao.GetPaasDao(p.txn, p.influxClient, ds).GetMetricUsageByTime(request)
		if err != nil {
			logger.Error(err)
			return result, err
		}
		convertedList, _ := utils.GetResponseConverter().InfluxConverterList(daoResult, v.Alias)
		result = append(result, convertedList)
	}
	return result, err
}

func (p *PaasService) GetPaasMemoryUsage(request model.PaasRequest) (result []map[string]interface{}, err model.ErrMessage) {

	ds := p.databases.PaastaDatabase

	daoResult, err := dao.GetPaasDao(p.txn, p.influxClient, ds).GetPaasMemoryUsage(request)
	if err != nil {
		logger.Error(err)
		return result, err
	}
	convertedList, _ := utils.GetResponseConverter().InfluxConverter4Usage(daoResult, request.Args.(model.MemoryMetricArg).Alias)
	result = append(result, convertedList)

	return result, err
}

func (p *PaasService) GetTopologicalView(request model.PaasRequest) (monitVms model.MonitVms, errMsg model.ErrMessage) {

	//Bosh Rest API: List all deployments
	resDeployments, err := p.boshClient.GetDeployments()
	if err != nil {
		errMsg = model.ErrMessage{
			"Message": "Failed to get deployment list from Bosh. err =" + err.Error(),
		}
		return monitVms, errMsg
	}

	var errs []error
	var boshDeployments []model.BoshDeployment
	var wg sync.WaitGroup
	wg.Add(len(resDeployments))
	for _, resDep := range resDeployments {
		go func(wg *sync.WaitGroup, resDep gogobosh.Deployment) {
			defer wg.Done()
			if model.BOSH_DEPLOYMENT_NAME_CF == resDep.Name {
				//Bosh Rest API: List all VMs(detail)
				vms, err := p.boshClient.GetDeploymentVMs(resDep.Name)
				if err != nil {
					errs = append(errs, err)
					logger.Error("Failed to get deployment's VM list from Bosh. err =" + err.Error())
				} else {
					deployment := model.BoshDeployment{}
					deployment.Name = resDep.Name
					deployment.VMS = vms
					boshDeployments = append(boshDeployments, deployment)
				}
			}
		}(&wg, resDep)
	}
	wg.Wait()
	if len(errs) > 0 {
		return monitVms, model.ErrMessage{"Message": err.Error()}
	}

	//Get PaaS-TA Summary
	paasSummary, errMsg := GetPaasService(p.txn, p.influxClient, p.databases, p.boshClient).GetPaasVms(request)
	if errMsg != nil {
		return monitVms, errMsg
	}

	//Make Topology Tree
	monitDeployments := make([]model.MonitDeployment, len(boshDeployments))
	for redDepIdx, resDep := range boshDeployments {

		monitDeployments[redDepIdx].Name = resDep.Name
		monitDeployments[redDepIdx].Status = model.STATE_RUNNING

		wg.Add(len(resDep.VMS))
		for _, resDepVm := range resDep.VMS {
			go func(wg *sync.WaitGroup, resDepVm gogobosh.VM) {
				defer wg.Done()
				for _, paasVm := range paasSummary.Data {
					if resDepVm.ID == paasVm.Id {

						usage := model.Usage{utils.FloattostrDigit2(paasVm.CpuUsage)}
						usages := []model.Usage{usage}
						vmCpuUsage := model.BoshVmUsage{model.USAGE_NAME_CPU, paasVm.CpuState, usages}

						usage = model.Usage{utils.FloattostrDigit2(paasVm.MemoryUsage)}
						usages = []model.Usage{usage}
						vmMemoryUsage := model.BoshVmUsage{model.USAGE_NAME_MEMORY, paasVm.MemoryState, usages}

						usage = model.Usage{utils.FloattostrDigit2(paasVm.TotalDiskUsage)}
						usages = []model.Usage{usage}
						vmDiskRootUsage := model.BoshVmUsage{model.USAGE_NAME_DISK_ROOT, paasVm.TotalDiskState, usages}

						usage = model.Usage{utils.FloattostrDigit2(paasVm.DataDiskUsage)}
						usages = []model.Usage{usage}
						vmDiskVcapUsage := model.BoshVmUsage{model.USAGE_NAME_DISK_VCAP, paasVm.DataDiskState, usages}

						boshVm := model.BoshVm{}
						boshVm.Name = paasVm.Name
						boshVm.Status = paasVm.State
						boshVm.BoshVmUsages = []model.BoshVmUsage{vmCpuUsage, vmMemoryUsage, vmDiskRootUsage, vmDiskVcapUsage}

						monitDeployments[redDepIdx].Status = model.DetermineVmState(monitDeployments[redDepIdx].Status, boshVm.Status)
						monitDeployments[redDepIdx].VMS = append(monitDeployments[redDepIdx].VMS, boshVm)
					}
				}
			}(&wg, resDepVm)
		}
		wg.Wait()
	}
	wg.Wait()

	monitVms.Name = model.BOSH_NAME
	monitVms.Deployments = monitDeployments

	return monitVms, errMsg
}

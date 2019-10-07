package service

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring/paas/dao"
	"kr/paasta/monitoring/paas/model"
	"kr/paasta/monitoring/paas/util"
	"kr/paasta/monitoring/utils"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type BoshStatusService struct {
	txn          *gorm.DB
	influxClient client.Client
	databases    model.Databases
}

func GetBoshStatusService(txn *gorm.DB, influxClient client.Client, databases model.Databases) *BoshStatusService {
	return &BoshStatusService{
		txn:          txn,
		influxClient: influxClient,
		databases:    databases,
	}
}

func (n BoshStatusService) GetBoshStatusOverview(request model.BoshSummaryReq) (model.BoshOverviewCntRes, model.ErrMessage) {
	boshSummary, err := GetBoshStatusService(n.txn, n.influxClient, n.databases).getBoshStatus(request)
	if err != nil {
		//log.Println(err)
		logger.Error(err)
	}

	return boshSummary.Overview, err
}

func (n BoshStatusService) GetBoshStatusSummary(request model.BoshSummaryReq) (res model.BoshStatusOverviewRes, err model.ErrMessage) {
	boshSummary, err := GetBoshStatusService(n.txn, n.influxClient, n.databases).getBoshStatus(request)
	if err != nil {
		//log.Println(err)
		logger.Error(err)
	}

	return boshSummary, err
}

func (n BoshStatusService) getBoshStatus(request model.BoshSummaryReq) (model.BoshStatusOverviewRes, model.ErrMessage) {

	fmt.Println("GetBoshStatusOverview Request API parameter =========+>", request)

	var boshOverviewRes model.BoshStatusOverviewRes

	config, readErr := util.ReadConfig(`config.ini`) // real
	//config, readErr := util.ReadConfig(`../../config.ini`) //test
	if readErr != nil {
		errMessage := model.ErrMessage{
			"Message": readErr.Error(),
		}
		return boshOverviewRes, errMessage
	}

	boshCnt, _ := strconv.Atoi(config["bosh.count"])

	var boshList []model.BoshSummaryReq

	for i := 0; i < boshCnt; i++ {
		var boshInfo model.BoshSummaryReq
		boshInfo.Name = config["bosh."+strconv.Itoa(i)+".name"]
		boshInfo.Ip = config["bosh."+strconv.Itoa(i)+".ip"]
		deployName := config["bosh."+strconv.Itoa(i)+".deployname"]

		boshIdResp, _ := dao.GetBoshStatusDao(n.txn, n.influxClient, n.databases).GetBoshId(deployName)
		boshId, _ := util.GetResponseConverter().InfluxConverterToMap(boshIdResp)
		if len(boshId) > 0 {
			boshInfo.Id = boshId[0]["id"].(string)
		}

		boshList = append(boshList, boshInfo)
	}

	//임계치 설정정보를 조회한다.
	serverThresholds, err := dao.GetAlarmPolicyDao(n.txn).GetAlarmPolicyList()
	if err != nil {
		return boshOverviewRes, err
	}

	var errs []model.ErrMessage

	for _, boshInfo := range boshList {

		var request model.BoshSummaryReq
		request.Name = boshInfo.Name
		request.Ip = boshInfo.Ip
		request.Id = boshInfo.Id

		cpuCoreData, cpuData, memTotData, memFreeData, diskTotalData, diskUsedData, diskDataTotalData, diskDataUsedData, err := n.GetBoshSummaryMetricData(request)

		if err != nil {
			errs = append(errs, err)
			continue
		}

		var boshData model.BoshSummaryRes

		boshData.Name = request.Name
		boshData.Address = request.Ip
		boshData.Id = request.Id

		cpuUsage := utils.GetDataFloatFromInterfaceSingle(cpuData)
		memTot := utils.GetDataFloatFromInterfaceSingle(memTotData)
		memFree := utils.GetDataFloatFromInterfaceSingle(memFreeData)
		memUsage := utils.RoundFloatDigit2(100 - ((memFree / memTot) * 100))
		diskTotal := utils.GetDataFloatFromInterfaceSingle(diskTotalData)
		diskUsed := utils.GetDataFloatFromInterfaceSingle(diskUsedData)
		var diskUsage = 100 - ((diskTotal - diskUsed) / diskTotal * 100)

		diskDataTotal := utils.GetDataFloatFromInterfaceSingle(diskDataTotalData)
		diskDataUsed := utils.GetDataFloatFromInterfaceSingle(diskDataUsedData)
		var diskDataUsage = 100 - ((diskDataTotal - diskDataUsed) / diskDataTotal * 100)

		boshData.Core = strconv.Itoa(len(cpuCoreData))
		boshData.CpuUsage = utils.RoundFloat(cpuUsage, 2)
		boshData.TotalMemory = memTot / model.MB
		boshData.MemoryUsage = memUsage
		boshData.TotalDisk = diskTotal / model.MB
		boshData.DataDisk = diskDataTotal / model.MB

		if boshData.Core == "0" || boshData.TotalMemory == 0 {
			boshData.State, boshData.BoshState, boshData.CpuErrStat, boshData.MemErrStat = model.BOSH_STATE_FAIL, model.BOSH_STATE_FAIL, model.BOSH_STATE_FAIL, model.BOSH_STATE_FAIL
		}

		if boshData.TotalDisk == 0 || boshData.DataDisk == 0 {
			boshData.DiskStatus, boshData.BoshState, boshData.DiskRootErrStat, boshData.DiskDataErrStat = model.BOSH_STATE_FAIL, model.BOSH_STATE_FAIL, model.BOSH_STATE_FAIL, model.BOSH_STATE_FAIL
		}

		// bosh state setting
		if boshData.State != model.BOSH_STATE_FAIL {
			var alarmStatus []string

			cpuStatus := util.GetAlarmStatusByServiceName(model.ORIGIN_TYPE_BOSH, model.ALARM_TYPE_CPU, boshData.CpuUsage, serverThresholds)
			memStatus := util.GetAlarmStatusByServiceName(model.ORIGIN_TYPE_BOSH, model.ALARM_TYPE_MEMORY, boshData.MemoryUsage, serverThresholds)

			if cpuStatus != "" {
				alarmStatus = append(alarmStatus, cpuStatus)
				boshData.CpuErrStat = cpuStatus
			} else {
				boshData.CpuErrStat = model.BOSH_STATE_RUNNING
			}
			if memStatus != "" {
				alarmStatus = append(alarmStatus, memStatus)
				boshData.MemErrStat = memStatus
			} else {
				boshData.MemErrStat = model.BOSH_STATE_RUNNING
			}

			state := util.GetMaxAlarmLevel(alarmStatus)
			if state == "" {
				boshData.State = model.BOSH_STATE_RUNNING
			} else {
				boshData.State = state
			}
		}

		// bosh diskStatus setting
		if boshData.DiskStatus != model.BOSH_STATE_FAIL {
			var diskStatusList []string
			diskStatus := util.GetAlarmStatusByServiceName(model.ORIGIN_TYPE_BOSH, model.ALARM_TYPE_DISK, diskUsage, serverThresholds)
			if diskStatus != "" {
				diskStatusList = append(diskStatusList, diskStatus)
				boshData.DiskRootErrStat = diskStatus
			} else {
				boshData.DiskRootErrStat = model.BOSH_STATE_NORMAL
			}

			diskDataStatus := util.GetAlarmStatusByServiceName(model.ORIGIN_TYPE_BOSH, model.ALARM_TYPE_DISK, diskDataUsage, serverThresholds)
			if diskDataStatus != "" {
				diskStatusList = append(diskStatusList, diskDataStatus)
				boshData.DiskDataErrStat = diskDataStatus
			} else {
				boshData.DiskDataErrStat = model.BOSH_STATE_NORMAL
			}

			diskState := util.GetMaxAlarmLevel(diskStatusList)
			if diskState == "" {
				boshData.DiskStatus = model.BOSH_STATE_NORMAL
			} else {
				boshData.DiskStatus = diskState
			}
		}

		if boshData.State == model.BOSH_STATE_RUNNING && boshData.DiskStatus == model.BOSH_STATE_NORMAL {
			boshData.BoshState = model.BOSH_STATE_RUNNING
		} else if boshData.BoshState != model.BOSH_STATE_FAIL {
			var boshStatusList []string
			boshStatusList = append(boshStatusList, boshData.State)
			if boshData.DiskStatus == model.BOSH_STATE_NORMAL {
				boshStatusList = append(boshStatusList, model.BOSH_STATE_RUNNING)
			} else {
				boshStatusList = append(boshStatusList, boshData.DiskStatus)
			}
			boshData.BoshState = util.GetMaxAlarmLevel(boshStatusList)
			boshData.State = boshData.BoshState
		}

		boshOverviewRes.Data = append(boshOverviewRes.Data, boshData)
	}

	if len(errs) > 0 {
		var returnErrMessage string
		for _, err := range errs {
			returnErrMessage = returnErrMessage + " " + err["Message"].(string)
		}
		errMessage := model.ErrMessage{
			"Message": returnErrMessage,
		}
		return boshOverviewRes, errMessage
	}

	// bosh overview
	totalCnt, failedCnt, criticalCnt, warningCnt := len(boshOverviewRes.Data), 0, 0, 0
	for _, value := range boshOverviewRes.Data {
		if value.BoshState == model.BOSH_STATE_FAIL {
			failedCnt++
		} else if value.BoshState == model.ALARM_LEVEL_CRITICAL {
			criticalCnt++
		} else if value.BoshState == model.ALARM_LEVEL_WARNING {
			warningCnt++
		}
	}

	boshOverviewRes.Overview.Total = strconv.Itoa(totalCnt)
	boshOverviewRes.Overview.Running = strconv.Itoa(totalCnt - failedCnt - criticalCnt - warningCnt)
	boshOverviewRes.Overview.Failed = strconv.Itoa(failedCnt)
	boshOverviewRes.Overview.Critical = strconv.Itoa(criticalCnt)
	boshOverviewRes.Overview.Warning = strconv.Itoa(warningCnt)

	boshOverviewRes.PageItem = request.PageItem
	boshOverviewRes.PageIndex = request.PageIndex
	boshOverviewRes.TotalCount = len(boshOverviewRes.Data)

	return boshOverviewRes, nil
}

func (n BoshStatusService) GetBoshSummaryMetricData(request model.BoshSummaryReq) ([]map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, model.ErrMessage) {
	var cpuResp, cpuCoreResp, memTotalResp, memFreeResp, diskTotalResp, diskUsedResp, diskDataTotalResp, diskDataUsedResp client.Response
	var errs []model.ErrMessage
	var err model.ErrMessage
	var wg sync.WaitGroup

	wg.Add(8)
	for i := 0; i < 8; i++ {
		go func(wg *sync.WaitGroup, index int) {
			switch index {
			case 0:
				request.MetricName = model.MTR_CPU_CORE
				request.Time = "1m"
				request.SqlQuery = "select value from bosh_metrics where id = '%s' and time > now() - %s and metricname =~ /%s/ group by metricname order by time desc limit 1;"
				cpuCoreResp, err = dao.GetBoshStatusDao(n.txn, n.influxClient, n.databases).GetDynamicBoshSummaryData(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 1:
				request.MetricName = model.MTR_CPU_CORE
				request.Time = "1m"
				request.SqlQuery = "select mean(value) as value from bosh_metrics where id = '%s' and time > now() - %s and metricname =~ /%s/ ;"
				cpuResp, err = dao.GetBoshStatusDao(n.txn, n.influxClient, n.databases).GetDynamicBoshSummaryData(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 2:
				request.MetricName = model.MTR_MEM_TOTAL
				request.Time = "1m"
				request.SqlQuery = "select mean(value) as value from bosh_metrics where id = '%s' and time > now() - %s and metricname = '%s' ;"
				memTotalResp, err = dao.GetBoshStatusDao(n.txn, n.influxClient, n.databases).GetDynamicBoshSummaryData(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 3:
				request.MetricName = model.MTR_MEM_FREE
				request.Time = "1m"
				request.SqlQuery = "select mean(value) as value from bosh_metrics where id = '%s' and time > now() - %s and metricname = '%s' ;"
				memFreeResp, err = dao.GetBoshStatusDao(n.txn, n.influxClient, n.databases).GetDynamicBoshSummaryData(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 4:
				request.MetricName = model.MTR_DISK_TOTAL
				request.Time = "1m"
				request.SqlQuery = "select mean(value) as value from bosh_metrics where id = '%s' and time > now() - %s and metricname = '%s' ;"
				diskTotalResp, err = dao.GetBoshStatusDao(n.txn, n.influxClient, n.databases).GetDynamicBoshSummaryData(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 5:
				request.MetricName = model.MTR_DISK_USED
				request.Time = "1m"
				request.SqlQuery = "select mean(value) as value from bosh_metrics where id = '%s' and time > now() - %s and metricname = '%s' ;"
				diskUsedResp, err = dao.GetBoshStatusDao(n.txn, n.influxClient, n.databases).GetDynamicBoshSummaryData(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 6:
				request.MetricName = model.MTR_DISK_DATA_TOTAL
				request.Time = "1m"
				request.SqlQuery = "select mean(value) as value from bosh_metrics where id = '%s' and time > now() - %s and metricname = '%s' ;"
				diskDataTotalResp, err = dao.GetBoshStatusDao(n.txn, n.influxClient, n.databases).GetDynamicBoshSummaryData(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 7:
				request.MetricName = model.MTR_DISK_DATA_USED
				request.Time = "1m"
				request.SqlQuery = "select mean(value) as value from bosh_metrics where id = '%s' and time > now() - %s and metricname = '%s' ;"
				diskDataUsedResp, err = dao.GetBoshStatusDao(n.txn, n.influxClient, n.databases).GetDynamicBoshSummaryData(request)
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
		errMessage := model.ErrMessage{
			"Message": returnErrMessage,
		}
		return nil, nil, nil, nil, nil, nil, nil, nil, errMessage
	}
	//==========================================================================

	cpuCore, _ := util.GetResponseConverter().InfluxConverterToMap(cpuCoreResp)
	memTotal, _ := utils.GetResponseConverter().InfluxConverter(memTotalResp)
	memFree, _ := utils.GetResponseConverter().InfluxConverter(memFreeResp)
	diskTotal, _ := utils.GetResponseConverter().InfluxConverter(diskTotalResp)
	cpuUsage, _ := utils.GetResponseConverter().InfluxConverter(cpuResp)
	diskUsage, _ := utils.GetResponseConverter().InfluxConverter(diskUsedResp)
	diskDataTotal, _ := utils.GetResponseConverter().InfluxConverter(diskDataTotalResp)
	diskDataUsage, _ := utils.GetResponseConverter().InfluxConverter(diskDataUsedResp)

	return cpuCore, cpuUsage, memTotal, memFree, diskTotal, diskUsage, diskDataTotal, diskDataUsage, nil
}

func (n BoshStatusService) GetTopProcessListByMemory(request model.BoshSummaryReq) (model.BoshTopprocessUsageRes, model.ErrMessage) {

	fmt.Println("GetTopProcessListByMemory Request API parameter =========+>", request)

	var topProcessRes model.BoshTopprocessUsageRes

	boshTopprocessResp, err := dao.GetBoshStatusDao(n.txn, n.influxClient, n.databases).GetBoshTopprocessList(request)

	if err != nil {
		fmt.Println(err)
		return topProcessRes, err
	} else {

		valueList, _ := utils.GetResponseConverter().InfluxConverterToMap(boshTopprocessResp)

		var resList []map[string]interface{}

		for z := 0; z < len(valueList); z++ {
			if len(resList) > 0 {
				chk := false
				for y := 0; y < len(resList); y++ {
					if resList[y][model.IFX_MTR_PROC_NAME] == valueList[z][model.IFX_MTR_PROC_NAME] && resList[y][model.IFX_MTR_PROC_PID] == valueList[z][model.IFX_MTR_PROC_PID] {
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
			return utils.TypeChecker_float64(resList[j][model.IFX_MTR_MEM_USAGE]).(float64) < utils.TypeChecker_float64(resList[i][model.IFX_MTR_MEM_USAGE]).(float64)
		})

		var idx int

		for _, vl := range resList {
			var topProcess model.BoshTopProcessUsage

			topProcess.Index = strconv.Itoa(idx + 1)
			topProcess.Process = utils.TypeChecker_string(vl[model.IFX_MTR_PROC_NAME]).(string)
			topProcess.Memory = utils.TypeChecker_float64(vl[model.IFX_MTR_MEM_USAGE]).(float64) / model.MB
			topProcess.Pid = strconv.FormatFloat(utils.TypeChecker_float64(vl[model.IFX_MTR_PROC_PID]).(float64), 'f', 0, 64)
			topProcess.Time = time.Unix(vl[model.IFX_MTR_TIME].(int64), 0).Format(time.RFC3339)[0:19]

			topProcessRes.Data = append(topProcessRes.Data, topProcess)
			idx++
			if idx == 5 {
				break
			} //fixed 5row
		}

		return topProcessRes, nil
	}
}

func (n BoshStatusService) GetBoshCpuUsageList(request model.BoshDetailReq) (result []map[string]interface{}, _ model.ErrMessage) {

	request.MetricName = model.MTR_CPU_CORE
	cpuUsageResp, err := dao.GetBoshStatusDao(n.txn, n.influxClient, n.databases).GetBoshCpuUsageList(request)

	if err != nil {
		fmt.Println(err)
		return result, err
	} else {
		cpuUsage, _ := utils.GetResponseConverter().InfluxConverterList(cpuUsageResp, model.RESP_DATA_CPU_NAME)
		result = append(result, cpuUsage)
		return result, nil
	}
}

func (n BoshStatusService) GetBoshCpuLoadList(request model.BoshDetailReq) (result []map[string]interface{}, _ model.ErrMessage) {

	request.MetricName = model.MTR_CPU_LOAD_1M
	cpuLoad1mResp, err := dao.GetBoshStatusDao(n.txn, n.influxClient, n.databases).GetBoshDetailList(request)
	request.MetricName = model.MTR_CPU_LOAD_5M
	cpuLoad5mResp, err := dao.GetBoshStatusDao(n.txn, n.influxClient, n.databases).GetBoshDetailList(request)
	request.MetricName = model.MTR_CPU_LOAD_15M
	cpuLoad15mResp, err := dao.GetBoshStatusDao(n.txn, n.influxClient, n.databases).GetBoshDetailList(request)

	if err != nil {
		fmt.Println(err)
		return result, err
	} else {
		cpu1mLoad, _ := utils.GetResponseConverter().InfluxConverterList(cpuLoad1mResp, model.RESP_DATA_LOAD_1M_NAME)
		cpu5mLoad, _ := utils.GetResponseConverter().InfluxConverterList(cpuLoad5mResp, model.RESP_DATA_LOAD_5M_NAME)
		cpu15mLoad, _ := utils.GetResponseConverter().InfluxConverterList(cpuLoad15mResp, model.RESP_DATA_LOAD_15M_NAME)

		result = append(result, cpu1mLoad)
		result = append(result, cpu5mLoad)
		result = append(result, cpu15mLoad)

		return result, nil
	}
}

func (n BoshStatusService) GetBoshMemoryUsageList(request model.BoshDetailReq) (result []map[string]interface{}, _ model.ErrMessage) {
	request.MetricName = model.MTR_MEM_USAGE
	memoryResp, err := dao.GetBoshStatusDao(n.txn, n.influxClient, n.databases).GetBoshMemUsageList(request)

	if err != nil {
		fmt.Println(err)
		return result, err
	} else {
		memoryUsage, _ := utils.GetResponseConverter().InfluxConverter4Usage(memoryResp, model.RESP_DATA_MEM_NAME)
		result = append(result, memoryUsage)
		return result, nil
	}
}

func (n BoshStatusService) GetBoshDiskUsageList(request model.BoshDetailReq) (result []map[string]interface{}, _ model.ErrMessage) {

	config, _ := util.ReadConfig(`config.ini`)
	var mountPoint = strings.Split(config["disk.mount.point"], ",")
	var mountPointMetricList []map[string]string

	for _, value := range mountPoint {
		var tmp = fmt.Sprintf(model.MTR_DISK_USAGE_STR, value)
		info := map[string]string{
			"metricname": tmp,
			"jsonname":   config["disk."+value+".resp.json.name"],
		}
		mountPointMetricList = append(mountPointMetricList, info)
	}

	for _, value := range mountPointMetricList {

		request.MetricName = value["metricname"]

		diskResp, err := dao.GetBoshStatusDao(n.txn, n.influxClient, n.databases).GetBoshDetailList(request)
		if err != nil {
			fmt.Println(err)
		} else {
			diskUsage, _ := utils.GetResponseConverter().InfluxConverterList(diskResp, value["jsonname"])
			result = append(result, diskUsage)
		}
	}
	return result, nil
}

func (n BoshStatusService) GetBoshDiskIoList(request model.BoshDetailReq) (result []map[string]interface{}, _ model.ErrMessage) {

	request.IsConvertKb = true
	config, _ := util.ReadConfig(`config.ini`)
	var mountPoint = strings.Split(config["disk.io.mount.point"], ",")

	var mountPointMetricList []map[string]string

	for _, value := range mountPoint {
		var tmp = fmt.Sprintf(model.MTR_DISK_IO_READ_STR, strings.Replace(value, "/", "\\/", -1)+"\\..*")
		info := map[string]string{
			"metricname": tmp,
			"jsonname":   config["disk.io."+value+".read.json.name"],
		}
		mountPointMetricList = append(mountPointMetricList, info)

		tmp = fmt.Sprintf(model.MTR_DISK_IO_WRITE_STR, strings.Replace(value, "/", "\\/", -1)+"\\..*")
		info = map[string]string{
			"metricname": tmp,
			"jsonname":   config["disk.io."+value+".write.json.name"],
		}
		mountPointMetricList = append(mountPointMetricList, info)
	}

	for _, value := range mountPointMetricList {

		request.MetricName = value["metricname"]

		diskResp, err := dao.GetBoshStatusDao(n.txn, n.influxClient, n.databases).GetBoshDetailList(request)

		if err != nil {
			fmt.Println(err)
		} else {
			diskUsage, _ := utils.GetResponseConverter().InfluxConverterList(diskResp, value["jsonname"])
			result = append(result, diskUsage)
		}
	}
	return result, nil
}

func (n BoshStatusService) GetBoshNetworkByteList(request model.BoshDetailReq) (result []map[string]interface{}, _ model.ErrMessage) {

	request.IsConvertKb = true

	config, _ := util.ReadConfig(`config.ini`)
	var mountPoint = strings.Split(config["network.monitor.item"], ",")

	var mountPointMetricList []map[string]string

	for _, value := range mountPoint {
		info := map[string]string{
			"metricname": fmt.Sprintf(model.MTR_NETWORK_BYTE_SENT, value),
			"jsonname":   model.RESP_DATA_NETWORK_IO_SENT_NAME,
		}
		mountPointMetricList = append(mountPointMetricList, info)

		info = map[string]string{
			"metricname": fmt.Sprintf(model.MTR_NETWORK_BYTE_RECV, value),
			"jsonname":   model.RESP_DATA_NETWORK_IO_RECV_NAME,
		}
		mountPointMetricList = append(mountPointMetricList, info)
	}

	for _, value := range mountPointMetricList {

		request.MetricName = value["metricname"]

		diskResp, err := dao.GetBoshStatusDao(n.txn, n.influxClient, n.databases).GetBoshDetailList(request)

		if err != nil {
			fmt.Println(err)
		} else {
			diskUsage, _ := utils.GetResponseConverter().InfluxConverterList(diskResp, value["jsonname"])
			result = append(result, diskUsage)
		}
	}
	return result, nil

}

func (n BoshStatusService) GetBoshNetworkPacketList(request model.BoshDetailReq) (result []map[string]interface{}, _ model.ErrMessage) {

	request.IsConvertKb = true

	config, _ := util.ReadConfig(`config.ini`)
	var mountPoint = strings.Split(config["network.monitor.item"], ",")

	var mountPointMetricList []map[string]string

	for _, value := range mountPoint {
		info := map[string]string{
			"metricname": fmt.Sprintf(model.MTR_NETWORK_PACKET_SENT, value),
			"jsonname":   model.RESP_DATA_NETWORK_IO_SENT_NAME,
		}
		mountPointMetricList = append(mountPointMetricList, info)

		info = map[string]string{
			"metricname": fmt.Sprintf(model.MTR_NETWORK_PACKET_RECV, value),
			"jsonname":   model.RESP_DATA_NETWORK_IO_RECV_NAME,
		}
		mountPointMetricList = append(mountPointMetricList, info)
	}

	for _, value := range mountPointMetricList {

		request.MetricName = value["metricname"]

		diskResp, err := dao.GetBoshStatusDao(n.txn, n.influxClient, n.databases).GetBoshDetailList(request)

		if err != nil {
			fmt.Println(err)
		} else {
			diskUsage, _ := utils.GetResponseConverter().InfluxConverterList(diskResp, value["jsonname"])
			result = append(result, diskUsage)
		}
	}
	return result, nil
}

func (n BoshStatusService) GetBoshNetworkDropList(request model.BoshDetailReq) (result []map[string]interface{}, _ model.ErrMessage) {

	config, _ := util.ReadConfig(`config.ini`)
	var mountPoint = strings.Split(config["network.monitor.item"], ",")

	var mountPointMetricList []map[string]string

	for _, value := range mountPoint {
		info := map[string]string{
			"metricname": fmt.Sprintf(model.MTR_NETWORK_DROP_IN, value),
			"jsonname":   model.RESP_DATA_NETWORK_IO_IN_NAME,
		}
		mountPointMetricList = append(mountPointMetricList, info)

		info = map[string]string{
			"metricname": fmt.Sprintf(model.MTR_NETWORK_DROP_OUT, value),
			"jsonname":   model.RESP_DATA_NETWORK_IO_OUT_NAME,
		}
		mountPointMetricList = append(mountPointMetricList, info)
	}

	for _, value := range mountPointMetricList {

		request.MetricName = value["metricname"]

		diskResp, err := dao.GetBoshStatusDao(n.txn, n.influxClient, n.databases).GetBoshDetailList(request)

		if err != nil {
			fmt.Println(err)
		} else {
			diskUsage, _ := utils.GetResponseConverter().InfluxConverterList(diskResp, value["jsonname"])
			result = append(result, diskUsage)
		}
	}
	return result, nil
}

func (n BoshStatusService) GetBoshNetworkErrorList(request model.BoshDetailReq) (result []map[string]interface{}, _ model.ErrMessage) {

	config, _ := util.ReadConfig(`config.ini`)
	var mountPoint = strings.Split(config["network.monitor.item"], ",")

	var mountPointMetricList []map[string]string

	for _, value := range mountPoint {
		info := map[string]string{
			"metricname": fmt.Sprintf(model.MTR_NETWORK_ERROR_IN, value),
			"jsonname":   model.RESP_DATA_NETWORK_IO_IN_NAME,
		}
		mountPointMetricList = append(mountPointMetricList, info)

		info = map[string]string{
			"metricname": fmt.Sprintf(model.MTR_NETWORK_ERROR_OUT, value),
			"jsonname":   model.RESP_DATA_NETWORK_IO_OUT_NAME,
		}
		mountPointMetricList = append(mountPointMetricList, info)
	}

	for _, value := range mountPointMetricList {

		request.MetricName = value["metricname"]

		diskResp, err := dao.GetBoshStatusDao(n.txn, n.influxClient, n.databases).GetBoshDetailList(request)

		if err != nil {
			fmt.Println(err)
		} else {
			diskUsage, _ := utils.GetResponseConverter().InfluxConverterList(diskResp, value["jsonname"])
			result = append(result, diskUsage)
		}
	}
	return result, nil
}

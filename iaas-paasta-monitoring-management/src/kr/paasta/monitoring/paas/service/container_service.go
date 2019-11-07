package service

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring/paas/dao"
	"kr/paasta/monitoring/paas/model"
	"kr/paasta/monitoring/paas/util"
	"kr/paasta/monitoring/utils"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type ContainerService struct {
	txn          *gorm.DB
	influxClient client.Client
	databases    model.Databases
}

func GetContainerService(txn *gorm.DB, influxClient client.Client, databases model.Databases) *ContainerService {
	return &ContainerService{
		txn:          txn,
		influxClient: influxClient,
		databases:    databases,
	}
}

//Cell에 배포된 App Container 배포 현황을 조회한다.
func (h ContainerService) GetContainerDeploy() ([]model.CellTileView, model.ErrMessage) {

	var cellInfos []model.CellTileView

	cellList, err := dao.GetContainerDao(h.txn, h.influxClient, h.databases).GetCellList()

	if err != nil {
		return cellInfos, err
	}

	cellMap := getZoneCellList(cellList, h)

	cellMapStruct, _ := mapToTreeStruct(cellMap, cellList, h)

	var cellInfoResult []model.CellTileView

	//Cell Name Sorting 위해 For Loop
	for _, cellInfo := range cellList {

		for _, cellMapInfo := range cellMapStruct {
			if cellInfo.CellName == cellMapInfo.CellName {
				cellMapInfo.ZoneName = cellInfo.ZoneName
				cellInfoResult = append(cellInfoResult, cellMapInfo)
			}
		}

	}

	return cellInfoResult, nil
}

//DB의  Cell정보와 MetricDB의 Container정보를 조합하여 구조화된 Map 정보구성
// cell -app1 - container1
//            - container2
//      - app2 - container3
func getZoneCellList(cellInfos []model.ZoneCellInfo, b ContainerService) map[string]map[string]map[string]string {

	cellMap := make(map[string]map[string]map[string]string)

	//Zone에 존재하는 Cell들에 실행되고 있는 Container 목록을 받아온다.
	for _, cellInfo := range cellInfos {

		containerResp, _ := dao.GetContainerDao(b.txn, b.influxClient, b.databases).GetCellContainersList(cellInfo.Ip)
		valueList, _ := util.GetResponseConverter().InfluxConverterToMap(containerResp)

		appMap := make(map[string]map[string]string)
		for _, value := range valueList {

			containerMap := make(map[string]string)
			appName := reflect.ValueOf(value["application_name"]).String()
			containerId := reflect.ValueOf(value["container_interface"]).String()
			if strings.Contains(containerId, model.CON_MTR_ID_PREFIX) {
				containerId = strings.Replace(containerId, model.CON_MTR_ID_PREFIX, "", 1)
			}

			applicationIndex := reflect.ValueOf(value["application_index"]).String()

			containerMap[containerId] = applicationIndex

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

func mapToTreeStruct(mapData map[string]map[string]map[string]string, dbCellInfo []model.ZoneCellInfo, b ContainerService) ([]model.CellTileView, model.ErrMessage) {

	returnValue := make([]model.CellTileView, len(mapData))
	cellInfo := make([]model.CellTileView, len(mapData))

	c := 0

	for cellName, apps := range mapData {

		var containerList []model.ContainerTileView
		for appName, containerInfos := range apps {

			var container model.ContainerTileView

			for key, data := range containerInfos {

				container.AppName = appName
				container.AppIndex = data
				container.ContainerId = key
				containerList = append(containerList, container)
			}
		}

		cellInfo[c].CellName = cellName
		cellInfo[c].ContainerTileView = containerList
		c++

	}

	sortIdx := 0
	for cellName, _ := range mapData {
		for _, info := range cellInfo {
			if cellName == info.CellName {

				for _, data := range dbCellInfo {
					if data.CellName == cellName {
						returnValue[sortIdx].Ip = data.Ip
						break
					}
				}
				returnValue[sortIdx].CellName = cellName
				returnValue[sortIdx].ContainerTileView = info.ContainerTileView
			}
		}
		sortIdx++
	}

	return returnValue, nil
}

func (h ContainerService) GetCellOverview(request model.ContainerReq) (model.OverviewCntRes, model.ErrMessage) {

	var result model.OverviewCntRes

	zoneList, err := GetContainerService(h.txn, h.influxClient, h.databases).GetContainerDeploy()

	if err != nil {
		fmt.Println(err)
		return result, err
	}

	//임계치 설정정보를 조회한다.
	serverThresholds, err := dao.GetAlarmPolicyDao(h.txn).GetAlarmPolicyList()
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	cellUsageList, _ := GetContainerService(h.txn, h.influxClient, h.databases).getCellUsageList(zoneList, serverThresholds)

	totalCnt, failedCnt, criticalCnt, warningCnt := len(cellUsageList), 0, 0, 0

	for _, value := range cellUsageList {
		if value.CellState == model.STATE_FAILED {
			failedCnt++
		} else if value.CellState == model.ALARM_LEVEL_CRITICAL {
			criticalCnt++
		} else if value.CellState == model.ALARM_LEVEL_WARNING {
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

func (h ContainerService) GetCellOverviewStatusList(request model.ContainerReq) ([]model.CellOverviewRes, model.ErrMessage) {

	var result []model.CellOverviewRes

	zoneList, err := GetContainerService(h.txn, h.influxClient, h.databases).GetContainerDeploy()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	//임계치 설정정보를 조회한다.
	serverThresholds, err := dao.GetAlarmPolicyDao(h.txn).GetAlarmPolicyList()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	cellUsageList, _ := GetContainerService(h.txn, h.influxClient, h.databases).getCellUsageList(zoneList, serverThresholds)

	for _, value := range cellUsageList {
		if value.CellState == request.Status {
			result = append(result, value)
		}
	}

	// sort
	if len(result) > 0 {
		sort.Slice(result, func(i, j int) bool {
			return result[i].CellName+result[i].Ip < result[j].CellName+result[j].Ip
		})
	}

	return result, nil
}

func (h ContainerService) GetContainerOverview(request model.ContainerReq) (model.OverviewCntRes, model.ErrMessage) {

	var result model.OverviewCntRes
	zonelist, err := dao.GetContainerDao(h.txn, h.influxClient, h.databases).GetZoneList()
	cellist, err := GetContainerService(h.txn, h.influxClient, h.databases).GetContainerDeploy()

	if err != nil {
		fmt.Println(err)
		return result, err
	}

	//임계치 설정정보를 조회한다.
	serverThresholds, err := dao.GetAlarmPolicyDao(h.txn).GetAlarmPolicyList()
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	var containerUsageList []model.ContainerOverviewRes

	var wg sync.WaitGroup
	wg.Add(len(zonelist))

	// zone 정보 별로 List 추출
	for _, zone := range zonelist {
		go func(wg *sync.WaitGroup, zone model.ZoneCellInfo) {
			_, containers, _ := getZoneItemCount(cellist, zone.ZoneName)

			for _, container := range containers {
				var apiRequest model.ContainerReq
				apiRequest.CellIp = container.Ip
				apiRequest.AppName = container.AppName
				apiRequest.AppIndex = container.AppIndex

				// select container monitoring
				containerUsageRes, err := h.getContainerUsageState(apiRequest, serverThresholds)
				if err != nil {
					fmt.Println(err)
				}
				containerUsageList = append(containerUsageList, containerUsageRes)
			}
			wg.Done()
		}(&wg, zone)
	}
	wg.Wait()

	totalCnt, failedCnt, criticalCnt, warningCnt := len(containerUsageList), 0, 0, 0

	for _, value := range containerUsageList {
		if value.Status == model.STATE_FAILED {
			failedCnt++
		} else if value.Status == model.ALARM_LEVEL_CRITICAL {
			criticalCnt++
		} else if value.Status == model.ALARM_LEVEL_WARNING {
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

func (h ContainerService) GetContainerOverviewStatusList(request model.ContainerReq) ([]model.ContainerOverviewRes, model.ErrMessage) {

	var result []model.ContainerOverviewRes
	zonelist, err := dao.GetContainerDao(h.txn, h.influxClient, h.databases).GetZoneList()
	cellist, err := GetContainerService(h.txn, h.influxClient, h.databases).GetContainerDeploy()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	//임계치 설정정보를 조회한다.
	serverThresholds, err := dao.GetAlarmPolicyDao(h.txn).GetAlarmPolicyList()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var containerUsageList []model.ContainerOverviewRes

	var wg sync.WaitGroup
	wg.Add(len(zonelist))
	// zone 정보 별로 List 추출
	for _, zone := range zonelist {
		go func(wg *sync.WaitGroup, zone model.ZoneCellInfo) {
			_, containers, _ := getZoneItemCount(cellist, zone.ZoneName)

			for _, container := range containers {
				var apiRequest model.ContainerReq
				apiRequest.CellIp = container.Ip
				apiRequest.AppName = container.AppName
				apiRequest.AppIndex = container.AppIndex

				// select container monitoring
				containerUsageRes, err := h.getContainerUsageState(apiRequest, serverThresholds)

				if err != nil {
					fmt.Println(err)
				}
				containerUsageRes.ZoneName = zone.ZoneName
				containerUsageRes.CellName = container.CellName
				containerUsageRes.ContainerName = container.ContainerId
				containerUsageList = append(containerUsageList, containerUsageRes)
			}
			wg.Done()
		}(&wg, zone)
	}
	wg.Wait()

	for _, value := range containerUsageList {
		if value.Status == request.Status {
			result = append(result, value)
		}
	}

	// sort
	if len(result) > 0 {
		sort.Slice(result, func(i, j int) bool {
			return result[i].ZoneName+result[i].CellName+result[i].AppName+result[i].AppIndex < result[j].ZoneName+result[j].CellName+result[j].AppName+result[j].AppIndex
		})
	}

	return result, nil
}

func (h ContainerService) GetContainerSummary(request model.ContainerReq, searchZone string) (model.ContainerSummaryPagingRes, model.ErrMessage) {
	var pagingResult model.ContainerSummaryPagingRes

	var result []model.ContainerSummary
	var tmpZoneList []model.ZoneCellInfo
	zonelist, err := dao.GetContainerDao(h.txn, h.influxClient, h.databases).GetZoneList()

	if searchZone != "" {
		for _, tmpZone := range zonelist {
			if strings.Contains(tmpZone.ZoneName, searchZone) {
				tmpZoneList = append(tmpZoneList, tmpZone)
			}
		}
	} else {
		tmpZoneList = zonelist
	}

	cellist, err := GetContainerService(h.txn, h.influxClient, h.databases).GetContainerDeploy()

	if err != nil {
		fmt.Println(err)
		return pagingResult, err
	}

	//임계치 설정정보를 조회한다.
	serverThresholds, err := dao.GetAlarmPolicyDao(h.txn).GetAlarmPolicyList()
	if err != nil {
		fmt.Println(err)
		return pagingResult, err
	}

	var totalContainerUsageList []model.ContainerOverviewRes

	var wg sync.WaitGroup
	wg.Add(len(tmpZoneList))

	// zone 정보 별로 List 추출
	for _, zone := range tmpZoneList {
		go func(wg *sync.WaitGroup, zone model.ZoneCellInfo) {
			var res model.ContainerSummary
			var containerUsageList []model.ContainerOverviewRes
			cells, containers, appcnt := getZoneItemCount(cellist, zone.ZoneName)

			res.ZoneName = zone.ZoneName
			res.CellCnt = strconv.Itoa(len(cells))
			res.ContainerCnt = strconv.Itoa(len(containers))
			res.AppCnt = strconv.Itoa(appcnt)

			for _, container := range containers {

				var apiRequest model.ContainerReq
				apiRequest.CellIp = container.Ip
				apiRequest.AppName = container.AppName
				apiRequest.AppIndex = container.AppIndex

				// select container monitoring
				containerUsageRes, err := h.getContainerUsageState(apiRequest, serverThresholds)
				if err != nil {
					fmt.Println(err)
				}
				containerUsageList = append(containerUsageList, containerUsageRes)
				totalContainerUsageList = append(totalContainerUsageList, containerUsageRes)
			}

			totalCnt, failedCnt, criticalCnt, warningCnt := len(containerUsageList), 0, 0, 0

			for _, value := range containerUsageList {
				if value.Status == model.STATE_FAILED {
					failedCnt++
				} else if value.Status == model.ALARM_LEVEL_CRITICAL {
					criticalCnt++
				} else if value.Status == model.ALARM_LEVEL_WARNING {
					warningCnt++
				}
			}

			res.RunningCnt = strconv.Itoa(totalCnt - failedCnt - criticalCnt - warningCnt)
			res.FailCnt = strconv.Itoa(failedCnt)
			res.CriticalCnt = strconv.Itoa(criticalCnt)
			res.WarningCnt = strconv.Itoa(warningCnt)

			result = append(result, res)
			wg.Done()
		}(&wg, zone)
	}
	wg.Wait()

	// sort
	if len(result) > 0 {
		sort.Slice(result, func(i, j int) bool {
			return result[i].ZoneName < result[j].ZoneName
		})
	}

	// data pagination
	if request.PageItems != 0 && request.PageIndex != 0 {
		offset := (request.PageIndex - 1) * request.PageItems
		limit := offset + request.PageItems
		if offset >= len(result) {
			// invalid request
			result = nil
		} else {
			if limit > len(result) {
				limit = len(result)
			}
			result = result[offset:limit]
		}
	}
	pagingResult.TotalCount = len(result)
	pagingResult.PageItems = request.PageItems

	pagingResult.CotainerSummaryList = result

	totalCnt, failedCnt, criticalCnt, warningCnt := len(totalContainerUsageList), 0, 0, 0

	for _, value := range totalContainerUsageList {
		if value.Status == model.STATE_FAILED {
			failedCnt++
		} else if value.Status == model.ALARM_LEVEL_CRITICAL {
			criticalCnt++
		} else if value.Status == model.ALARM_LEVEL_WARNING {
			warningCnt++
		}
	}

	pagingResult.Overview.Total = strconv.Itoa(totalCnt)
	pagingResult.Overview.Running = strconv.Itoa(totalCnt - failedCnt - criticalCnt - warningCnt)
	pagingResult.Overview.Failed = strconv.Itoa(failedCnt)
	pagingResult.Overview.Critical = strconv.Itoa(criticalCnt)
	pagingResult.Overview.Warning = strconv.Itoa(warningCnt)

	return pagingResult, nil
}

func (h ContainerService) GetContainerRelationship(request model.ContainerReq) ([]model.ContainerRelationshipRes, model.ErrMessage) {

	var result []model.ContainerRelationshipRes
	zonelist, err := dao.GetContainerDao(h.txn, h.influxClient, h.databases).GetZoneList()
	cellist, err := GetContainerService(h.txn, h.influxClient, h.databases).GetContainerDeploy()

	if err != nil {
		fmt.Println(err)
		return result, err
	}

	//임계치 설정정보를 조회한다.
	serverThresholds, err := dao.GetAlarmPolicyDao(h.txn).GetAlarmPolicyList()
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	// zone 정보 별로 List 추출
	for _, zone := range zonelist {
		if zone.ZoneName == request.Name {

			cells, _, _ := getZoneItemCount(cellist, zone.ZoneName)

			var wg sync.WaitGroup
			wg.Add(len(cells))

			for _, cell := range cells {
				go func(wg *sync.WaitGroup, cell model.ZoneCellInfo) {
					var res model.ContainerRelationshipRes

					res.ZoneName = zone.ZoneName
					res.CellName = cell.CellName
					res.Ip = cell.Ip

					_, containers, _ := getZoneItemCount(cellist, zone.ZoneName)
					if len(containers) > 0 {
						var appInfoList []model.AppStatusInfo

						for _, container := range containers {
							if cell.CellName == container.CellName {
								var apiRequest model.ContainerReq
								apiRequest.CellIp = container.Ip
								apiRequest.AppName = container.AppName
								apiRequest.AppIndex = container.AppIndex
								apiRequest.ContainerName = container.ContainerId
								// select container monitoring
								containerUsageRes, err := h.getContainerUsageState(apiRequest, serverThresholds)
								if err != nil {
									fmt.Println(err)
								}

								var appstate model.AppStatusInfo
								appstate.AppName = containerUsageRes.AppName
								appstate.AppIndex = containerUsageRes.AppIndex
								appstate.Status = containerUsageRes.Status
								appstate.ContainerId = containerUsageRes.ContainerName
								appInfoList = append(appInfoList, appstate)
							}
						}

						totalCnt, failedCnt, criticalCnt, warningCnt := len(appInfoList), 0, 0, 0

						for _, value := range appInfoList {
							if value.Status == model.STATE_FAILED {
								failedCnt++
							} else if value.Status == model.ALARM_LEVEL_CRITICAL {
								criticalCnt++
							} else if value.Status == model.ALARM_LEVEL_WARNING {
								warningCnt++
							}
						}

						res.RunningCnt = strconv.Itoa(totalCnt - failedCnt - criticalCnt - warningCnt)
						res.FailCnt = strconv.Itoa(failedCnt)
						res.CriticalCnt = strconv.Itoa(criticalCnt)
						res.WarningCnt = strconv.Itoa(warningCnt)

						res.ContainerCnt = strconv.Itoa(len(appInfoList))
						res.AppCnt = strconv.Itoa(getContainerAppCount(appInfoList))
						res.AppInfoList = appInfoList
					}

					result = append(result, res)

					wg.Done()
				}(&wg, cell)
			}
			wg.Wait()
		}
	}

	// sort
	if len(result) > 0 {
		sort.Slice(result, func(i, j int) bool {
			return result[i].CellName < result[j].CellName
		})

		for _, value := range result {
			sort.Slice(value.AppInfoList, func(i, j int) bool {
				return value.AppInfoList[i].AppName+value.AppInfoList[i].AppIndex < value.AppInfoList[j].AppName+value.AppInfoList[j].AppIndex
			})
		}
	}

	return result, nil
}

func (h ContainerService) GetPaasMainContainerView(request model.ContainerReq) ([]model.ContainerRelationshipRes, model.ErrMessage) {

	var result []model.ContainerRelationshipRes

	cellist, err := GetContainerService(h.txn, h.influxClient, h.databases).GetContainerDeploy()

	if err != nil {
		fmt.Println(err)
		return result, err
	}

	//임계치 설정정보를 조회한다.
	serverThresholds, err := dao.GetAlarmPolicyDao(h.txn).GetAlarmPolicyList()
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	// zone 정보 별로 List 추출
	var wg sync.WaitGroup
	wg.Add(len(cellist))

	for _, cell := range cellist {
		go func(wg *sync.WaitGroup, cell model.CellTileView) {
			var res model.ContainerRelationshipRes

			res.ZoneName = cell.ZoneName
			res.CellName = cell.CellName
			res.Ip = cell.Ip

			if len(cell.ContainerTileView) > 0 {
				var appInfoList []model.AppStatusInfo

				for _, container := range cell.ContainerTileView {

					var apiRequest model.ContainerReq
					apiRequest.CellIp = cell.Ip
					apiRequest.AppName = container.AppName
					apiRequest.AppIndex = container.AppIndex
					apiRequest.ContainerName = container.ContainerId
					// select container monitoring
					containerUsageRes, err := h.getContainerUsageState(apiRequest, serverThresholds)
					if err != nil {
						fmt.Println(err)
					}

					var appstate model.AppStatusInfo
					appstate.AppName = containerUsageRes.AppName
					appstate.AppIndex = containerUsageRes.AppIndex
					appstate.Status = containerUsageRes.Status
					appstate.ContainerId = containerUsageRes.ContainerName
					appInfoList = append(appInfoList, appstate)

				}

				totalCnt, failedCnt, criticalCnt, warningCnt := len(appInfoList), 0, 0, 0

				for _, value := range appInfoList {
					if value.Status == model.STATE_FAILED {
						failedCnt++
					} else if value.Status == model.ALARM_LEVEL_CRITICAL {
						criticalCnt++
					} else if value.Status == model.ALARM_LEVEL_WARNING {
						warningCnt++
					}
				}

				res.RunningCnt = strconv.Itoa(totalCnt - failedCnt - criticalCnt - warningCnt)
				res.FailCnt = strconv.Itoa(failedCnt)
				res.CriticalCnt = strconv.Itoa(criticalCnt)
				res.WarningCnt = strconv.Itoa(warningCnt)

				res.ContainerCnt = strconv.Itoa(len(appInfoList))
				res.AppCnt = strconv.Itoa(getContainerAppCount(appInfoList))
				res.AppInfoList = appInfoList
			}

			result = append(result, res)

			wg.Done()
		}(&wg, cell)
	}
	wg.Wait()

	// sort
	if len(result) > 0 {
		sort.Slice(result, func(i, j int) bool {
			return result[i].CellName < result[j].CellName
		})

		for _, value := range result {
			sort.Slice(value.AppInfoList, func(i, j int) bool {
				return value.AppInfoList[i].AppName+value.AppInfoList[i].AppIndex < value.AppInfoList[j].AppName+value.AppInfoList[j].AppIndex
			})
		}
	}

	return result, nil
}

func (h ContainerService) getCellUsageList(cellList []model.CellTileView, serverThresholds []model.AlarmPolicyResponse) ([]model.CellOverviewRes, model.ErrMessage) {
	var resultList []model.CellOverviewRes

	var wg sync.WaitGroup
	wg.Add(len(cellList))

	for _, cellValue := range cellList {
		go func(wg *sync.WaitGroup, cellValue model.CellTileView) {
			var result model.CellOverviewRes
			result.CellName = cellValue.CellName
			result.Ip = cellValue.Ip

			var request model.ContainerReq
			request.CellIp = cellValue.Ip

			cpuCoreData, cpuData, memTotData, memFreeData, diskTotalData, diskUsageData, err := h.GetCellSummaryMetricData(request)
			if err != nil {
				fmt.Println(err)
			}

			// get cell id for paas detail view
			cellIdResp, _ := dao.GetContainerDao(h.txn, h.influxClient, h.databases).GetCellIdForDetail(request)
			cellId, _ := util.GetResponseConverter().InfluxConverterToMap(cellIdResp)
			if len(cellId) > 0 {
				result.CellId = cellId[0]["id"].(string)
			}

			cpuUsage := utils.GetDataFloatFromInterfaceSingle(cpuData)
			memTot := utils.GetDataFloatFromInterfaceSingle(memTotData)
			memFree := utils.GetDataFloatFromInterfaceSingle(memFreeData)
			memUsage := utils.RoundFloatDigit2(100 - ((memFree / memTot) * 100))
			diskTotal := utils.GetDataFloatFromInterfaceSingle(diskTotalData)
			diskUsage := utils.GetDataFloatFromInterfaceSingle(diskUsageData)

			result.Core = strconv.Itoa(len(cpuCoreData))
			result.CpuUsage = utils.RoundFloat(cpuUsage, 2)
			result.TotalDisk = diskTotal / model.MB
			result.TotalMemory = memTot / model.MB
			if memUsage < 0 {
				result.MemoryUsage = 0
			} else {
				result.MemoryUsage = memUsage
			}

			if result.Core == "0" || result.TotalMemory == 0 {
				result.State, result.CellState, result.CpuState, result.MemoryState = model.STATE_FAILED, model.STATE_FAILED, model.STATE_FAILED, model.STATE_FAILED
			}

			if result.TotalDisk == 0 {
				result.DiskStatus, result.CellState, result.TotalDiskState = model.STATE_FAILED, model.STATE_FAILED, model.STATE_FAILED
			}

			if result.State != model.STATE_FAILED {
				var alarmStatus []string

				cpuStatus := util.GetAlarmStatusByServiceName(model.ORIGIN_TYPE_CONTAINER, model.ALARM_TYPE_CPU, result.CpuUsage, serverThresholds)
				memStatus := util.GetAlarmStatusByServiceName(model.ORIGIN_TYPE_CONTAINER, model.ALARM_TYPE_MEMORY, result.MemoryUsage, serverThresholds)

				if cpuStatus != "" {
					alarmStatus = append(alarmStatus, cpuStatus)
					result.CpuState = cpuStatus
				} else {
					result.CpuState = model.STATE_RUNNING
				}
				if memStatus != "" {
					alarmStatus = append(alarmStatus, memStatus)
					result.MemoryState = memStatus
				} else {
					result.MemoryState = model.STATE_RUNNING
				}

				state := util.GetMaxAlarmLevel(alarmStatus)
				if state == "" {
					result.State = model.STATE_RUNNING
				} else {
					result.State = state
				}
			}

			if result.DiskStatus != model.STATE_FAILED {
				var diskStatusList []string
				diskStatus := util.GetAlarmStatusByServiceName(model.ORIGIN_TYPE_CONTAINER, model.ALARM_TYPE_DISK, diskUsage, serverThresholds)
				if diskStatus != "" {
					diskStatusList = append(diskStatusList, diskStatus)
					result.TotalDiskState = diskStatus
				} else {
					result.TotalDiskState = model.DISK_STATE_NORMAL
				}

				diskState := util.GetMaxAlarmLevel(diskStatusList)
				if diskState == "" {
					result.DiskStatus = model.DISK_STATE_NORMAL
				} else {
					result.DiskStatus = diskState
				}
			}

			if result.State == model.STATE_RUNNING && result.DiskStatus == model.DISK_STATE_NORMAL {
				result.CellState = model.STATE_RUNNING
			} else if result.CellState != model.STATE_FAILED {
				var statusList []string
				statusList = append(statusList, result.State)
				if result.DiskStatus == model.DISK_STATE_NORMAL {
					statusList = append(statusList, model.STATE_RUNNING)
				} else {
					statusList = append(statusList, result.DiskStatus)
				}
				result.CellState = util.GetMaxAlarmLevel(statusList)
				result.State = result.CellState
			}

			resultList = append(resultList, result)
			wg.Done()
		}(&wg, cellValue)
	}

	wg.Wait()

	return resultList, nil
}

func (h ContainerService) GetCellSummaryMetricData(request model.ContainerReq) ([]map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, model.ErrMessage) {
	var cpuResp, cpuCoreResp, memTotalResp, memFreeResp, diskTotalResp, diskUsageResp client.Response
	var errs []model.ErrMessage
	var err model.ErrMessage
	var wg sync.WaitGroup

	wg.Add(6)
	for i := 0; i < 6; i++ {
		go func(wg *sync.WaitGroup, index int) {
			switch index {
			case 0:
				request.MetricName = model.MTR_CPU_CORE
				request.Time = "1m"
				request.SqlQuery = "select value from cf_metrics where ip = '%s' and time > now() - %s and metricname =~ /%s/ group by metricname order by time desc limit 1;"
				cpuCoreResp, err = dao.GetContainerDao(h.txn, h.influxClient, h.databases).GetCellSummaryData(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 1:
				request.MetricName = model.MTR_CPU_CORE
				request.Time = "1m"
				request.SqlQuery = "select mean(value) as value from cf_metrics where ip = '%s' and time > now() - %s and metricname =~ /%s/ ;"
				cpuResp, err = dao.GetContainerDao(h.txn, h.influxClient, h.databases).GetCellSummaryData(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 2:
				request.MetricName = model.MTR_MEM_TOTAL
				request.Time = "1m"
				request.SqlQuery = "select mean(value) as value from cf_metrics where ip = '%s' and time > now() - %s and metricname = '%s' ;"
				memTotalResp, err = dao.GetContainerDao(h.txn, h.influxClient, h.databases).GetCellSummaryData(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 3:
				request.MetricName = model.MTR_MEM_FREE
				request.Time = "1m"
				request.SqlQuery = "select mean(value) as value from cf_metrics where ip = '%s' and time > now() - %s and metricname = '%s' ;"
				memFreeResp, err = dao.GetContainerDao(h.txn, h.influxClient, h.databases).GetCellSummaryData(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 4:
				request.MetricName = model.MTR_DISK_TOTAL
				request.Time = "1m"
				request.SqlQuery = "select mean(value) as value from cf_metrics where ip = '%s' and time > now() - %s and metricname = '%s' ;"
				diskTotalResp, err = dao.GetContainerDao(h.txn, h.influxClient, h.databases).GetCellSummaryData(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 5:
				request.MetricName = model.MTR_DISK_USAGE
				request.Time = "1m"
				request.SqlQuery = "select mean(value) as value from cf_metrics where ip = '%s' and time > now() - %s and metricname = '%s' ;"
				diskUsageResp, err = dao.GetContainerDao(h.txn, h.influxClient, h.databases).GetCellSummaryData(request)
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

	if len(errs) > 0 {
		var returnErrMessage string
		for _, err := range errs {
			returnErrMessage = returnErrMessage + " " + err["Message"].(string)
		}
		errMessage := model.ErrMessage{
			"Message": returnErrMessage,
		}
		return nil, nil, nil, nil, nil, nil, errMessage
	}

	cpuCore, _ := util.GetResponseConverter().InfluxConverterToMap(cpuCoreResp)
	memTotal, _ := utils.GetResponseConverter().InfluxConverter(memTotalResp)
	memFree, _ := utils.GetResponseConverter().InfluxConverter(memFreeResp)
	diskTotal, _ := utils.GetResponseConverter().InfluxConverter(diskTotalResp)
	cpuUsage, _ := utils.GetResponseConverter().InfluxConverter(cpuResp)
	diskUsage, _ := utils.GetResponseConverter().InfluxConverter(diskUsageResp)

	return cpuCore, cpuUsage, memTotal, memFree, diskTotal, diskUsage, nil
}

func getZoneItemCount(cellInfoList []model.CellTileView, zoneName string) ([]model.ZoneCellInfo, []model.ContainerTileView, int) {

	var cells []model.ZoneCellInfo
	var containers []model.ContainerTileView
	var apps []string

	for _, cell := range cellInfoList {
		if zoneName == cell.ZoneName {

			var tmpCell model.ZoneCellInfo
			tmpCell.CellName = cell.CellName
			tmpCell.Ip = cell.Ip
			cells = append(cells, tmpCell)

			if cell.ContainerTileView != nil {
				for _, container := range cell.ContainerTileView {
					container.Ip = cell.Ip
					container.CellName = cell.CellName
					containers = append(containers, container)
				}
			}
		}
	}

	if len(containers) > 0 {
		for _, app := range containers {
			var existChk bool = false
			if len(apps) > 0 {
				for _, name := range apps {
					if app.AppName == name {
						existChk = true
					}
				}
			}
			if !existChk {
				apps = append(apps, app.AppName)
			}
		}
	}

	return cells, containers, len(apps)
}

func getContainerAppCount(appList []model.AppStatusInfo) int {

	var apps []string

	for _, app := range appList {
		var existChk bool = false
		if len(apps) > 0 {
			for _, name := range apps {
				if app.AppName == name {
					existChk = true
				}
			}
		}
		if !existChk {
			apps = append(apps, app.AppName)
		}
	}

	return len(apps)
}

func (h ContainerService) getContainerUsageState(request model.ContainerReq, serverThresholds []model.AlarmPolicyResponse) (model.ContainerOverviewRes, model.ErrMessage) {

	var result model.ContainerOverviewRes

	cpuData, memData, diskData, err := h.GetContainerummaryMetricData(request)

	if err != nil {
		return result, err
	}

	result.AppName = request.AppName
	result.AppIndex = request.AppIndex
	result.Ip = request.CellIp
	result.ContainerName = request.ContainerName

	result.CpuUsage = utils.GetDataFloatFromInterfaceSingle(cpuData)
	result.MemoryUsage = utils.GetDataFloatFromInterfaceSingle(memData)
	result.DiskUsage = utils.GetDataFloatFromInterfaceSingle(diskData)

	var alarmStatus []string
	cpuStatus := util.GetAlarmStatusByServiceName(model.ORIGIN_TYPE_CONTAINER, model.ALARM_TYPE_CPU, result.CpuUsage, serverThresholds)
	if cpuStatus != "" {
		result.CpuState = cpuStatus
	} else {
		result.CpuState = model.STATE_RUNNING
	}
	alarmStatus = append(alarmStatus, cpuStatus)

	memStatus := util.GetAlarmStatusByServiceName(model.ORIGIN_TYPE_CONTAINER, model.ALARM_TYPE_MEMORY, result.MemoryUsage, serverThresholds)
	if memStatus != "" {
		result.MemoryState = memStatus
	} else {
		result.MemoryState = model.STATE_RUNNING
	}
	alarmStatus = append(alarmStatus, memStatus)

	diskStatus := util.GetAlarmStatusByServiceName(model.ORIGIN_TYPE_CONTAINER, model.ALARM_TYPE_DISK, result.DiskUsage, serverThresholds)
	if diskStatus != "" {
		result.DiskState = diskStatus
	} else {
		result.DiskState = model.STATE_RUNNING
	}
	alarmStatus = append(alarmStatus, diskStatus)

	result.Status = util.GetMaxAlarmLevel(alarmStatus)

	if result.Status == "" {
		result.Status = model.STATE_RUNNING
	}

	return result, nil
}

func (h ContainerService) GetContainerummaryMetricData(request model.ContainerReq) (map[string]interface{}, map[string]interface{}, map[string]interface{}, model.ErrMessage) {
	var cpuResp, memResp, diskResp client.Response
	var errs []model.ErrMessage
	var err model.ErrMessage
	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func(wg *sync.WaitGroup, index int) {

			switch index {
			case 0:
				request.MetricName = "cpu_usage_total"
				request.Service = model.ALARM_TYPE_CPU
				cpuResp, err = dao.GetContainerDao(h.txn, h.influxClient, h.databases).GetContainerUsage(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 1:
				request.MetricName = "memory_usage"
				request.Service = model.ALARM_TYPE_MEMORY
				memResp, err = dao.GetContainerDao(h.txn, h.influxClient, h.databases).GetContainerUsage(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 2:
				request.MetricName = "disk_usage"
				request.Service = model.ALARM_TYPE_DISK
				diskResp, err = dao.GetContainerDao(h.txn, h.influxClient, h.databases).GetContainerUsage(request)
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
	if len(errs) > 0 {
		var returnErrMessage string
		for _, err := range errs {
			returnErrMessage = returnErrMessage + " " + err["Message"].(string)
		}
		errMessage := model.ErrMessage{
			"Message": returnErrMessage,
		}
		return nil, nil, nil, errMessage
	}
	//==========================================================================

	cpuUsage, _ := utils.GetResponseConverter().InfluxConverter(cpuResp)
	memUsage, _ := utils.GetResponseConverter().InfluxConverter(memResp)
	diskUsage, _ := utils.GetResponseConverter().InfluxConverter(diskResp)

	return cpuUsage, memUsage, diskUsage, nil
}

func (h ContainerService) GetPaasContainerUsages(request model.ContainerReq) (result []map[string]interface{}, err model.ErrMessage) {

	for _, item := range request.Item {
		request.MetricName = item.Name

		resp, err := dao.GetContainerDao(h.txn, h.influxClient, h.databases).GetPaasContainerDetailUsages(request)

		if err != nil {
			fmt.Println(err)
			return result, err
		} else {
			usage, _ := utils.GetResponseConverter().InfluxConverterList(resp, item.ResName)
			result = append(result, usage)
		}
	}

	return result, err
}

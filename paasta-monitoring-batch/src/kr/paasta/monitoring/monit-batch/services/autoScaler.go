package services

import (
	client "github.com/influxdata/influxdb/client/v2"
	"kr/paasta/monitoring/monit-batch/dao"
	cb "kr/paasta/monitoring/monit-batch/models/base"
	mod "kr/paasta/monitoring/monit-batch/models"
	"kr/paasta/monitoring/monit-batch/util"
	"fmt"
	"reflect"
	"sync"
	"strconv"
	"encoding/json"
)

func PortalAutoScale(f *BackendServices){

	cellList, _ := dao.GetContainerAlarmDao(f.Influxclient).GetCellList(f.MonitoringDbClient)
	//AutoScale 설정된 App정보 조회
	autoScaleAppListConfig, _ := dao.GetAutoScaleAppDao(f.Influxclient).GetAutoScaleAppList(f.PortalDbClient)

	containerList := getAppCellList(cellList, f)
	//Cpu,Memory AutoScale 임계치 체크
	AppUsageList, _ := checkSystemUsage(containerList, autoScaleAppListConfig, f)

	autoScaleAppList := getAutoScaleList( autoScaleAppListConfig, AppUsageList)
	fmt.Println("autoScaleAppList===>", autoScaleAppList)

	var cpu float64
	var mem float64
	cpu = 5.5
	mem = 5.5
	var autoScaleTestApp mod.AutoScaleAction
	autoScaleTestApp.AppGuid = "6ab3fa59-b784-400c-a456-9e1dfd2dd7d7"
	autoScaleTestApp.AppName = "autoscailing-test1"
	autoScaleTestApp.Cause = "cpu"
	autoScaleTestApp.Instance = "3"
	autoScaleTestApp.CpuUsage = util.Floattostr(cpu)
	autoScaleTestApp.MemoryUsage = util.Floattostr(mem)
	autoScaleTestApp.Action = "I"

	//Portal Check
	err := util.PortalExistCHeck()
	if err == nil{
		body, _ := json.Marshal(autoScaleTestApp)
		_, _, err :=util.HttpRequest("/external/app/updateApp",  "POST", nil,  body, *mod.PortalClient)

		if err != nil{
			fmt.Println("AutoScale Request Err:", err)
		}

		/*
			for _, data := range autoScaleAppList  {
				body, _ := json.Marshal(data)
				response, code, err :=util.HttpRequest("/external/app/updateApp",  "POST", nil,  body, *mod.PortalClient)
				fmt.Println("res:::", response)
				fmt.Println("code:::", code)
				fmt.Println("err:::", err)
			}
		*/
	}

}

func getAutoScaleList(configAppList []mod.AutoScaleConfig, AppUsageList []mod.CellTileView) ([]mod.AutoScaleAction){

	var autoScaleAppList []mod.AutoScaleAction

	for _, configData := range configAppList{
		for _, cellData := range AppUsageList{

			for _, appUsageData := range cellData.ContainerTileView{

				//AutoScale Config등록된 App만 Auto Scale 대상
				if appUsageData.AppGuid == configData.Guid{

					//AutoScale Out Check (autoScaleIncrease Y 인경우)
					if configData.AutoIncreaseYn == "Y"{
						if appUsageData.CpuUsage > float64(configData.CpuThresholdMaxPer * 100){

							if isExistApp(autoScaleAppList, configData.App) == false{
								var autoScale mod.AutoScaleAction
								autoScale.AppName = configData.App
								autoScale.AppGuid = configData.Guid
								autoScale.CpuUsage = util.Floattostr(appUsageData.CpuUsage)
								autoScale.MemoryUsage = util.Floattostr(appUsageData.MemoryUsage)
								autoScale.Action = "O"
								autoScale.Instance = getAppInstances(configData.Guid, AppUsageList)
								autoScale.Cause = "CPU"
								autoScaleAppList = append(autoScaleAppList, autoScale)
							}
						}
						if appUsageData.MemoryUsage > float64(configData.MemoryMaxSize * 100){
							if isExistApp(autoScaleAppList, configData.App) == false{
								var autoScale mod.AutoScaleAction
								autoScale.AppName = configData.App
								autoScale.AppGuid = configData.Guid
								autoScale.CpuUsage = util.Floattostr(appUsageData.CpuUsage)
								autoScale.MemoryUsage = util.Floattostr(appUsageData.MemoryUsage)
								autoScale.Action = "O"
								autoScale.Instance = getAppInstances(configData.Guid, AppUsageList)
								autoScale.Cause = "MEM"
								autoScaleAppList = append(autoScaleAppList, autoScale)
							}
						}
					}

					//AutoScale In Check (autoScaleDecerease Y 인경우)
					if configData.AutoDecreaseYn == "Y"{

						if appUsageData.CpuUsage < float64(configData.CpuThresholdMinPer){
							if isExistApp(autoScaleAppList, configData.App) == false{
								var autoScale mod.AutoScaleAction
								autoScale.AppName = configData.App
								autoScale.AppGuid = configData.Guid
								autoScale.CpuUsage = util.Floattostr(appUsageData.CpuUsage)
								autoScale.MemoryUsage = util.Floattostr(appUsageData.MemoryUsage)
								autoScale.Action = "I"
								autoScale.Instance = getAppInstances(configData.Guid, AppUsageList)
								autoScale.Cause = "CPU"
								autoScaleAppList = append(autoScaleAppList, autoScale)
							}
						}

						if appUsageData.MemoryUsage < float64(configData.MemoryMinSize){
							if isExistApp(autoScaleAppList, configData.App) == false{
								var autoScale mod.AutoScaleAction
								autoScale.AppName = configData.App
								autoScale.AppGuid = configData.Guid
								autoScale.CpuUsage = util.Floattostr(appUsageData.CpuUsage)
								autoScale.MemoryUsage = util.Floattostr(appUsageData.MemoryUsage)
								autoScale.Action = "I"
								autoScale.Instance = getAppInstances(configData.Guid, AppUsageList)
								autoScale.Cause = "MEM"
								autoScaleAppList = append(autoScaleAppList, autoScale)
							}
						}
					}

				}
			}

		}
	}
	return autoScaleAppList
}

func getAppInstances (appGUid string, AppUsageList []mod.CellTileView) string{

	instancesCnt := 0
	for _, cellInfos := range AppUsageList{

		for _, appInfos := range cellInfos.ContainerTileView{

			if appGUid == appInfos.AppGuid{
				instancesCnt += 1
			}
		}
	}
	return strconv.Itoa(instancesCnt)
}

func isExistApp(appList []mod.AutoScaleAction, appName string) bool{

	for _, data := range appList{
		if data.AppName == appName{
			return true
		}
	}
	return false
}


func getAppCellList(cellInfos []mod.ZoneCellInfo, f *BackendServices) map[string]map[string]map[string]string{

	cellMap := make(map[string]map[string]map[string]string)

	//Zone에 존재하는 Cell들에 실행되고 있는 Container 목록을 받아온다.
	for _, cellInfo := range cellInfos{
		var request mod.ZonesReq
		request.CellIp = cellInfo.Ip
		request.MetricDatabase = f.InfluxConfig.ContainerDatabase
		containerResp, _ := dao.GetAutoScaleAppDao(f.Influxclient).GetAppContainersList(request)
		valueList, _ := util.GetResponseConverter().InfluxConverterToMap(containerResp)

		appMap := make(map[string]map[string]string)
		for _ , value := range valueList{

			containerMap     := make(map[string]string)
			/*appName 	 := reflect.ValueOf(value["application_name"]).String()*/
			appGuid 	 := reflect.ValueOf(value["application_id"]).String()
			containerName 	 := reflect.ValueOf(value["container_interface"]).String()
			applicationIndex := reflect.ValueOf(value["application_index"]).String()

			containerMap[containerName] = applicationIndex

			//동일한 App의 Container는 AppMap에 Append 처리 한다.
			if exists, ok := appMap[appGuid]; ok{
				for k, v := range containerMap {
					exists[k] = v
					appMap[appGuid] = exists
				}
			}else{
				appMap[appGuid] = containerMap
			}

		}
		cellMap[cellInfo.CellName] = appMap
	}
	return cellMap
}

func checkSystemUsage(mapData map[string]map[string]map[string]string, autoScaleList []mod.AutoScaleConfig, f *BackendServices) ([]mod.CellTileView, cb.ErrMessage){

	returnValue := make([]mod.CellTileView, len(mapData))
	cellInfo := make([]mod.CellTileView, len(mapData))

	c := 0

	for cellName, apps := range mapData{

		var containers []mod.ContainerTileView

		var wg sync.WaitGroup
		var errs []cb.ErrMessage
		wg.Add(len(apps))

		for appGuid, containerInfos := range apps {
			go func(wg *sync.WaitGroup, info map[string]string, containerAppGuid string){

				for _, data := range autoScaleList{
					//AutoScale 설정된 App만 체크
					if containerAppGuid == data.Guid{
						var containerInfo mod.ContainerTileView

						for name, index := range info {

							idx, _ := strconv.Atoi(index)
							containerInfo.AppIndex = idx
							containerInfo.ContainerName = name
							containerInfo.AppGuid = containerAppGuid

							var request mod.ZonesReq
							request.ContainerName = name
							request.MetricDatabase = f.InfluxConfig.ContainerDatabase
							request.CheckTIme = strconv.Itoa(data.CheckTimeSec)

							cpuData, memData,  err := GetContainerSummary_Sub(request, f)

							if err != nil {
								errs = append(errs, err)
							}

							cpuUsage  := util.GetDataFloatFromInterfaceSingle(cpuData)
							memUsage  := util.GetDataFloatFromInterfaceSingle(memData)

							containerInfo.CpuUsage = cpuUsage
							containerInfo.MemoryUsage = memUsage

							containers = append(containers, containerInfo)
						}
					}
				}



				wg.Done()
			}(&wg, containerInfos, appGuid)
		}
		wg.Wait()

		//==========================================================================
		// Error가 여러건일 경우 대해 고려해야함.
		if len(errs) > 0 {
			var returnErrMessage string
			for _, err := range errs{
				returnErrMessage = returnErrMessage + " " + err["Message"].(string)
			}
			errMessage := cb.ErrMessage{
				"Message": returnErrMessage ,
			}
			return nil, errMessage
		}
		//==========================================================================

		cellInfo[c].CellName = cellName
		cellInfo[c].ContainerTileView = containers
		c++

	}

	sortIdx := 0
	for cellName, _ := range mapData{
		for  _, info := range cellInfo{
			if cellName == info.CellName {
				returnValue[sortIdx].CellName = cellName
				returnValue[sortIdx].ContainerTileView =  info.ContainerTileView
			}
		}
		sortIdx++
	}

	return returnValue, nil
}

//Server 상태 목록 조회 - DAO 호출.
func  GetContainerSummary_Sub(request mod.ZonesReq, f *BackendServices) (map[string]interface{}, map[string]interface{}, cb.ErrMessage) {
	var cpuResp, memResp client.Response
	var errs []cb.ErrMessage
	var err cb.ErrMessage
	var wg sync.WaitGroup
	wg.Add(2)

	for i := 0; i < 2; i++ {
		go func(wg *sync.WaitGroup, index int) {

			switch index {
			case 0 :
				cpuResp, err = dao.GetAutoScaleAppDao(f.Influxclient).GetContainerCpuUsage(request)
				if err != nil {
					errs = append(errs, err)
				}
			case 1 :
				memResp, err = dao.GetAutoScaleAppDao(f.Influxclient).GetContainerMemoryUsage(request)
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
		for _, err := range errs{
			returnErrMessage = returnErrMessage + " " + err["Message"].(string)
		}
		errMessage := cb.ErrMessage{
			"Message": returnErrMessage ,
		}
		return nil, nil, errMessage
	}
	//==========================================================================

	cpuUsage,   _   := util.GetResponseConverter().InfluxConverter(cpuResp, request.ContainerName)
	memUsage,   _   := util.GetResponseConverter().InfluxConverter(memResp, request.ContainerName)

	return cpuUsage,memUsage, nil

}
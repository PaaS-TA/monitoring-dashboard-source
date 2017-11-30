package services_test

import (
	. "github.com/onsi/ginkgo"
	_ "github.com/go-sql-driver/mysql"
	"kr/paasta/monitoring/monit-batch/services"
	"kr/paasta/monitoring/monit-batch/models"
	"kr/paasta/monitoring/monit-batch/dao"
	"kr/paasta/monitoring/monit-batch/util"
	"reflect"
	"fmt"
)

var _ = Describe("Bosh Service Test", func() {

	It("Portal AutoScale", func() {
		config, _ := ReadConfig(`../config_test.ini`)
		influxCon, configDbCon, portalDbCon, boshCon, mailConfig, config := GetObject(config)

		backendService := services.NewBackendServices(-9,  influxCon, configDbCon, portalDbCon, boshCon,  mailConfig, config)

		var cellInfoList []models.ZoneCellInfo
		var cellInfo models.ZoneCellInfo
		cellInfo.Id = 1
		//cellInfo.Ip = "10.244.16.106"
		cellInfo.Ip = "10.10.4.61"
		cellInfo.CellName = "cell_z1/0"
		cellInfo.ZoneName = "z1"
		cellInfoList = append(cellInfoList, cellInfo)
		result := GetAppCellList(cellInfoList, backendService)
		appGuidList := GetAppGuid(result)

		fmt.Println("result==appGuid==>", appGuidList)
		if len(appGuidList) >= 4 {

			//CPU ScaleOut
			autoScaleData1 := models.AutoScaleConfig{No:9000, Guid:appGuidList[0], Org: "org", Space: "space", App: "test1", InstanceMaxCnt: 3, InstanceMinCnt: 1,
				CpuThresholdMaxPer: 0, CpuThresholdMinPer: 0, MemoryMaxSize: 90, MemoryMinSize:10, CheckTimeSec: 30, AutoDecreaseYn: "Y", AutoIncreaseYn: "Y"}
			//CPU ScaleIn
			autoScaleData2 := models.AutoScaleConfig{No:9001, Guid:appGuidList[1], Org: "org", Space: "space", App: "test2", InstanceMaxCnt: 3, InstanceMinCnt: 1,
				CpuThresholdMaxPer: 90, CpuThresholdMinPer: 80, MemoryMaxSize: 99, MemoryMinSize:98, CheckTimeSec: 30, AutoDecreaseYn: "Y", AutoIncreaseYn: "Y"}
			//Memory ScaleOut
			autoScaleData3 := models.AutoScaleConfig{No:9002, Guid:appGuidList[2], Org: "org", Space: "space", App: "test3", InstanceMaxCnt: 3, InstanceMinCnt: 1,
				CpuThresholdMaxPer: 90, CpuThresholdMinPer: 1, MemoryMaxSize: 10, MemoryMinSize:8, CheckTimeSec: 30, AutoDecreaseYn: "Y", AutoIncreaseYn: "Y"}
			//Memory ScaleIn
			autoScaleData4 := models.AutoScaleConfig{No:9003, Guid:appGuidList[3], Org: "org", Space: "space", App: "test4", InstanceMaxCnt: 3, InstanceMinCnt: 1,
				CpuThresholdMaxPer: 90, CpuThresholdMinPer: 1, MemoryMaxSize: 95, MemoryMinSize:80, CheckTimeSec: 30, AutoDecreaseYn: "Y", AutoIncreaseYn: "Y"}

			backendService.PortalDbClient.FirstOrCreate(&autoScaleData1)
			backendService.PortalDbClient.FirstOrCreate(&autoScaleData2)
			backendService.PortalDbClient.FirstOrCreate(&autoScaleData3)
			backendService.PortalDbClient.FirstOrCreate(&autoScaleData4)

			services.PortalAutoScale(backendService)
			backendService.PortalDbClient.Delete(&autoScaleData1)
			backendService.PortalDbClient.Delete(&autoScaleData2)
			backendService.PortalDbClient.Delete(&autoScaleData3)
			backendService.PortalDbClient.Delete(&autoScaleData4)

		}
		/*Eventually(func()  {
			services.BoshAlarmCollect(backendService)
		})*/
	})
})

func GetAppGuid(cellInfos map[string]map[string]map[string]string) []string{

	var appGuidList []string
	for _, apps := range cellInfos{
		for appGuid, _ := range apps {
			if appGuid != ""{
				appGuidList = append(appGuidList, appGuid)
			}
		}
	}

	return appGuidList
}

func GetAppCellList(cellInfos []models.ZoneCellInfo, f *services.BackendServices) map[string]map[string]map[string]string{

	cellMap := make(map[string]map[string]map[string]string)

	//Zone에 존재하는 Cell들에 실행되고 있는 Container 목록을 받아온다.
	for _, cellInfo := range cellInfos{
		var request models.ZonesReq
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




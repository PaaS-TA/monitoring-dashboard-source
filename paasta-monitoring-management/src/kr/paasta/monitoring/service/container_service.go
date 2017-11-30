package service

import (
	"github.com/jinzhu/gorm"
	client "github.com/influxdata/influxdb/client/v2"
	"kr/paasta/monitoring/dao"
	"kr/paasta/monitoring/util"
	"kr/paasta/monitoring/domain"
	"reflect"
)

type ContainerService struct {
	txn   *gorm.DB
	influxClient 	client.Client
}

func GetContainerService(txn *gorm.DB, influxClient client.Client) *ContainerService {
	return &ContainerService{
		txn:   txn,
		influxClient: 	influxClient,
	}
}

//Cell에 배포된 App Container 배포 현황을 조회한다.
func (h ContainerService) GetContainerDeploy()([]domain.CellTileView, domain.ErrMessage){

	var cellInfos []domain.CellTileView

	cellList, err := dao.GetContainerDao(h.txn, h.influxClient).GetCellList()

	if err != nil{
		return cellInfos, err
	}

	cellMap   := getZoneCellList(cellList, h)

	cellMapStruct , _ := mapToTreeStruct(cellMap, cellList, h)

	var cellInfoResult []domain.CellTileView

	//Cell Name Sorting 위해 For Loop
	for _, cellInfo := range cellList{

		for _, cellMapInfo := range cellMapStruct{
			if cellInfo.CellName == cellMapInfo.CellName{
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
func getZoneCellList(cellInfos []domain.ZoneCellInfo, b ContainerService) map[string]map[string]map[string]string{

	cellMap := make(map[string]map[string]map[string]string)

	//Zone에 존재하는 Cell들에 실행되고 있는 Container 목록을 받아온다.
	for _, cellInfo := range cellInfos{

		containerResp, _ := dao.GetContainerDao(b.txn, b.influxClient).GetCellContainersList(cellInfo.Ip)
		valueList, _ := util.GetResponseConverter().InfluxConverterToMap(containerResp)

		appMap := make(map[string]map[string]string)
		for _ , value := range valueList{

			containerMap     := make(map[string]string)
			appName 	 := reflect.ValueOf(value["application_name"]).String()
			containerName 	 := reflect.ValueOf(value["container_interface"]).String()
			applicationIndex := reflect.ValueOf(value["application_index"]).String()

			containerMap[containerName] = applicationIndex

			//동일한 App의 Container는 AppMap에 Append 처리 한다.
			if exists, ok := appMap[appName]; ok{
				for k, v := range containerMap {
					exists[k] = v
					appMap[appName] = exists
				}
			}else{
				appMap[appName] = containerMap
			}

		}
		cellMap[cellInfo.CellName] = appMap
	}

	return cellMap
}

func mapToTreeStruct(mapData map[string]map[string]map[string]string, dbCellInfo []domain.ZoneCellInfo, b ContainerService) ([]domain.CellTileView, domain.ErrMessage){

	returnValue := make([]domain.CellTileView, len(mapData))
	cellInfo    := make([]domain.CellTileView, len(mapData))

	c := 0

	for cellName, apps := range mapData{

		var containerList []domain.ContainerTileView
		for appName, containerInfos := range apps {

			var container domain.ContainerTileView
			for _, data := range containerInfos{
				container.AppName = appName
				container.AppIndex = data

				containerList = append(containerList, container)
			}
		}

		cellInfo[c].CellName = cellName
		cellInfo[c].ContainerTileView = containerList
		c++

	}

	sortIdx := 0
	for cellName, _ := range mapData{
		for  _, info := range cellInfo{
			if cellName == info.CellName {

				for _, data := range dbCellInfo{
					if data.CellName == cellName{
						returnValue[sortIdx].Ip = data.Ip
						break
					}
				}
				returnValue[sortIdx].CellName = cellName
				returnValue[sortIdx].ContainerTileView =  info.ContainerTileView
			}
		}
		sortIdx++
	}

	return returnValue, nil
}
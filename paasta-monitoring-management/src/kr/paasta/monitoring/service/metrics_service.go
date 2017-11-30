package service

import (
	"encoding/json"
	"kr/paasta/monitoring/domain"
	"kr/paasta/monitoring/dao"
	"kr/paasta/monitoring/util"
	"github.com/influxdata/influxdb/client/v2"
	"fmt"
	"strconv"
)

const(
	DATA_NAME     string = "data"
	PERSON        string = "person"
	SERVICE_NAME  string  = "serviceName"
)

type MetricsService struct {
	influxClient 	client.Client
	databases 	domain.Databases
}

func GetMetricsService(influxClient client.Client, databases domain.Databases) *MetricsService {
	return &MetricsService{
		influxClient:   influxClient,
		databases:      databases,
	}
}


func (b MetricsService) GetDiskIOList(request domain.MetricsRequest) (domain.DiskIOUsage, domain.ErrMessage) {
	var ds string
	readName := "diskIOStats.vda1.readBytes"
	writeName := "diskIOStats.vda1.writeBytes"
	/*
	config, _ := util.ReadConfig(`../config.ini`)
	if request.Origin == "bos" {
		ds , _ = config["metric.infra.db_name"]
	} else if request.Origin == "ctl" {
		ds , _ = config["metric.controller.db_name"]
	} else if request.Origin == "ctn" {
		ds , _ = config["metric.controller.db_name"]
	}*/
	if request.Origin == "bos" {
		ds = b.databases.BoshDatabase
	} else if request.Origin == "ctl" {
		ds = b.databases.PaastaDatabase
	} else if request.Origin == "ctn" {
		ds = b.databases.ContainerDatabase
	}

	var diskIOUsage domain.DiskIOUsage
	var diskIOUsageList []domain.DiskIOUsageList
	DiskIORdResp, err := dao.GetMetricsDao(b.influxClient, ds).GetDiskIOList(request, readName)
	DiskIOWrtResp, err := dao.GetMetricsDao(b.influxClient, ds).GetDiskIOList(request, writeName)
	var diskIORdUsage []domain.MultiUsage
	var diskIOWrtUsage []domain.MultiUsage

	if err != nil {
		return diskIOUsage, err
	}else {

		DiskIORdUsage, _ := DiskIORdResp, request.ServiceName
		DiskIOWrtUsage, _ := DiskIOWrtResp, request.ServiceName

		readUsage, _ := util.GetResponseConverter().InfluxConverter(DiskIORdUsage, request.ServiceName)
		writeUsage, _ := util.GetResponseConverter().InfluxConverter(DiskIOWrtUsage, request.ServiceName)

		for _, data := range readUsage{

			switch data.(type){
			case []map[string]interface {}:
				datamap := data.([]map[string]interface{})
				diskIORdUsage = make([]domain.MultiUsage, len(datamap))

				for idx_i, value := range datamap{
					diskIORdUsage[idx_i].Time = value["time"].(int64)
					diskIORdUsage[idx_i].Usage = value["totalUsage"].(json.Number)
				}
			}
		}

		for _, data := range writeUsage{

			switch data.(type){
			case []map[string]interface {}:
				datamap := data.([]map[string]interface{})
				diskIOWrtUsage = make([]domain.MultiUsage, len(datamap))

				for idx_i, value := range datamap{
					diskIOWrtUsage[idx_i].Time = value["time"].(int64)
					diskIOWrtUsage[idx_i].Usage = value["totalUsage"].(json.Number)
				}
			}
		}

		diskIOUsageList = make([]domain.DiskIOUsageList, 2)

		diskIOUsageList[0].Name = "I/O Read"
		diskIOUsageList[0].Data   = diskIORdUsage

		diskIOUsageList[1].Name = "I/O Write"
		diskIOUsageList[1].Data   = diskIOWrtUsage

		diskIOUsage.DiskIOUsage = diskIOUsageList

		return diskIOUsage, nil

	}
}

func (b MetricsService) GetNetworkIOList(request domain.MetricsRequest) (domain.NetworkIOUsage, domain.ErrMessage) {
	var ds string
	rxName := "networkIOStats.eth0.bytesRecv"
	txName := "networkIOStats.eth0.bytesSent"
	/*
	config, _ := util.ReadConfig(`../config.ini`)
	if request.Origin == "bos" {
		ds , _ = config["metric.infra.db_name"]
	} else if request.Origin == "ctl" {
		ds , _ = config["metric.controller.db_name"]
	} else if request.Origin == "ctn" {
		ds , _ = config["metric.controller.db_name"]
	}*/
	if request.Origin == "bos" {
		ds = b.databases.BoshDatabase
	} else if request.Origin == "ctl" {
		ds = b.databases.PaastaDatabase
	} else if request.Origin == "ctn" {
		ds = b.databases.ContainerDatabase
	}

	var networkIOUsage domain.NetworkIOUsage
	var networkIOUsageList []domain.NetworkIOUsageList
	NetworkRxIOResp, err := dao.GetMetricsDao(b.influxClient, ds).GetNetworkIOUsageList(request, rxName)
	NetworkTxIOResp, err := dao.GetMetricsDao(b.influxClient, ds).GetNetworkIOUsageList(request, txName)

	var networkRxIOUsage []domain.NetworkUsage
	var networkTxIOUsage []domain.NetworkUsage

	if err != nil {
		return networkIOUsage, err
	}else {

		NetworkRxIOUsage, _ := NetworkRxIOResp, request.ServiceName
		NetworkTxIOUsage, _ := NetworkTxIOResp, request.ServiceName

		rxUsage, _ := util.GetResponseConverter().InfluxConverter(NetworkRxIOUsage, request.ServiceName)
		txUsage, _ := util.GetResponseConverter().InfluxConverter(NetworkTxIOUsage, request.ServiceName)

		var maxValue float64
		var unit     string
		for _, data := range rxUsage{
			switch data.(type){
			case []map[string]interface {}:
				datamap := data.([]map[string]interface{})
				for _, value := range datamap{
					usage, _ := strconv.ParseFloat(value["totalUsage"].(json.Number).String(),64)
					if usage > maxValue{
						maxValue = usage
					}
				}
			}
		}

		for _, data := range txUsage{
			switch data.(type){
			case []map[string]interface {}:
				datamap := data.([]map[string]interface{})
				for _, value := range datamap{

					usage, _ := strconv.ParseFloat(value["totalUsage"].(json.Number).String(),64)
					if usage > maxValue{
						maxValue = usage
					}
				}
			}
		}

		if maxValue >= 1000000{
			unit = "(M)"
		}else if maxValue >= 1000{
			unit = "(K)"
		}else{
			unit = ""
		}

		for _, data := range rxUsage{

			switch data.(type){
			case []map[string]interface {}:
				datamap := data.([]map[string]interface{})
				networkRxIOUsage = make([]domain.NetworkUsage, len(datamap))

				for idx_i, value := range datamap{

					usage, _ := strconv.ParseFloat(value["totalUsage"].(json.Number).String(),64)

					networkRxIOUsage[idx_i].Time = value["time"].(int64)
					if unit == "(M)"{
						networkRxIOUsage[idx_i].Usage = usage/1000000
					}else if unit == "(K)"{
						networkRxIOUsage[idx_i].Usage = usage/1000
					}else{
						networkRxIOUsage[idx_i].Usage = usage
					}
					networkRxIOUsage[idx_i].Usage = usage

					if usage > maxValue{
						maxValue = usage
					}
				}
			}
		}

		for _, data := range txUsage{

			switch data.(type){
			case []map[string]interface {}:
				datamap := data.([]map[string]interface{})
				networkTxIOUsage = make([]domain.NetworkUsage, len(datamap))

				for idx_i, value := range datamap{

					usage, _ := strconv.ParseFloat(value["totalUsage"].(json.Number).String(),64)
					networkTxIOUsage[idx_i].Time = value["time"].(int64)
					if unit == "(M)"{
						networkTxIOUsage[idx_i].Usage = usage/1000000
					}else if unit == "(K)"{
						networkTxIOUsage[idx_i].Usage = usage/1000
					}else{
						networkTxIOUsage[idx_i].Usage = usage
					}

					if usage > maxValue{
						maxValue = usage
					}
				}
			}
		}

		networkIOUsageList = make([]domain.NetworkIOUsageList, 2)

		networkIOUsageList[0].Name = "RxPackets"
		networkIOUsageList[0].Unit = unit
		networkIOUsageList[0].Data   = networkRxIOUsage

		networkIOUsageList[1].Name = "TxPackets"
		networkIOUsageList[1].Unit = unit
		networkIOUsageList[1].Data   = networkTxIOUsage


		networkIOUsage.NetworkIOUsage = networkIOUsageList

		return networkIOUsage, nil

	}
}

func (b MetricsService) GetTopProcessList(request domain.MetricsRequest) (domain.TopProcessList, domain.ErrMessage) {
	var ds string
	/*
	config, _ := util.ReadConfig(`../config.ini`)
	if request.Origin == "bos" {
		ds , _ = config["metric.infra.db_name"]
	} else if request.Origin == "ctl" {
		ds , _ = config["metric.controller.db_name"]
	} else if request.Origin == "ctn" {
		ds , _ = config["metric.controller.db_name"]
	} else if request.Origin == "app" {
		ds , _ = config["metric.container.db_name"]
	}*/
	if request.Origin == "bos" {
		ds = b.databases.BoshDatabase
	} else if request.Origin == "ctl" {
		ds = b.databases.PaastaDatabase
	} else if request.Origin == "ctn" {
		ds = b.databases.PaastaDatabase
	} else if request.Origin == "app" {
		ds = b.databases.ContainerDatabase
	}


	resp , err := dao.GetMetricsDao(b.influxClient, ds).GetTopProcessList(request)

	var topProcessList domain.TopProcessList
	var topProcessUsage []domain.TopProcessUsage

	if err != nil {
		return topProcessList, err
	}else{

		result, _ := util.GetResponseConverter().InfluxConverter(resp, request.ServiceName)
		var serviceName string
		for _, data := range result{

			switch data.(type){
			case []map[string]interface {}:
				datamap := data.([]map[string]interface{})
				topProcessUsage = make([]domain.TopProcessUsage, 10)

				for idx_i := 0; idx_i < len(datamap); idx_i++ {
					value := datamap[idx_i]
					topProcessUsage[idx_i].Index = strconv.Itoa(idx_i + 1)
					topProcessUsage[idx_i].Memory = value["memory"].(json.Number)
					topProcessUsage[idx_i].Pid = value["pid"].(json.Number)
					topProcessUsage[idx_i].Process = value["process"].(string)
					topProcessUsage[idx_i].Time = value["time"].(int64)
				}
			case string :
				serviceName = data.(string)
			}

		}

		topProcessList.ServiceName   = serviceName
		topProcessList.Data   = topProcessUsage

		return topProcessList, nil
	}
}


func (b MetricsService) GetAppCpuUsage(request domain.MetricsRequest) (map[string]interface{}, domain.ErrMessage) {

	resp , err := dao.GetMetricsDao(b.influxClient, b.databases.ContainerDatabase).GetAppCpuUsage(request)

	if err != nil {
		return nil, err
	}else{
		result, err := util.GetResponseConverter().InfluxConverter(resp, request.ServiceName)
		if err != nil {
			return nil, err
		}
		fmt.Println("#### Cpu variation :", result)
		return result, nil
	}
}

func (b MetricsService) GetAppMemoryUsage(request domain.MetricsRequest) (map[string]interface{}, domain.ErrMessage) {

	resp , err := dao.GetMetricsDao(b.influxClient, b.databases.ContainerDatabase).GetAppMemoryUsage(request)
	if err != nil {
		return nil, err
	}else{
		result, err := util.GetResponseConverter().InfluxConverter(resp, request.ServiceName)
		if err != nil {
			return nil, err
		}
		fmt.Println("#### Memory variation :", result)
		return result, nil
	}
}

func (b MetricsService) GetDiskUsage(request domain.MetricsRequest) (map[string]interface{}, domain.ErrMessage) {

	resp , err := dao.GetMetricsDao(b.influxClient, b.databases.ContainerDatabase).GetAppDiskUsage(request)
	if err != nil {
		return nil, err
	}else{
		result, err := util.GetResponseConverter().InfluxConverter(resp, request.ServiceName)
		if err != nil {
			return nil, err
		}
		fmt.Println("#### Cpu variation :", result)
		return result, nil
	}
}

func (b MetricsService) GetApplicationResources(request domain.MetricsRequest) (domain.ApplicationResources, domain.ErrMessage) {
	resp , err := dao.GetMetricsDao(b.influxClient, b.databases.ContainerDatabase).GetApplicationResources(request)

	var appResources domain.ApplicationResources
	if err != nil {
		return appResources, err
	}else{
		result, _ := util.GetResponseConverter().InfluxConverter(resp, "resources")
		for _, resources := range result{
			switch resources.(type){
			case []map[string]interface {}:
				datamap := resources.([]map[string]interface{})
				for _, data := range datamap{
					if data["name"] == "cpu_usage_total" {
						appResources.CpuUsage = data["value"].(json.Number)
					}
					if data["name"] == "memory_usage" {
						appResources.MemUsage = data["value"].(json.Number)
					}
					if data["name"] == "disk_usage" {
						appResources.DiskUsage = data["value"].(json.Number)
					}
				}
			}
		}
		return appResources, nil
	}
}

func (b MetricsService) GetApplicationResourcesAll(request domain.MetricsRequest) (domain.ApplicationResources, domain.ErrMessage) {
	resp , err := dao.GetMetricsDao(b.influxClient, b.databases.ContainerDatabase).GetApplicationResourcesAll(request)

	var appResources domain.ApplicationResources
	var appInfo []domain.ApplicationInfo

	if err != nil {
		return appResources, err
	}else{
		result, _ := util.GetResponseConverter().InfluxConverter(resp, "resources")

		for _, data := range result{

			switch data.(type){
			case []map[string]interface {}:
				datamap := data.([]map[string]interface{})
				appInfo = make([]domain.ApplicationInfo, 10)

				for idx_i := 0; idx_i < len(datamap); idx_i++ {
					value := datamap[idx_i]
					appInfo[idx_i].Time = value["time"].(int64)
					appInfo[idx_i].Id = value["application_id"].(string)
					appInfo[idx_i].Index = value["application_index"].(string)
					appInfo[idx_i].Name = value["name"].(string)
					appInfo[idx_i].Value = value["value"].(json.Number)
				}
			}

		}
		appResources.Data   = appInfo
		return appResources, nil
	}
}


func (b MetricsService) GetAppNetworkKByte(request domain.MetricsRequest) (domain.AppNetworkIOKbyte, domain.ErrMessage) {

	result := domain.AppNetworkIOKbyte{Name: "", Data: make([]map[string]interface{}, 2)}
	//Application Network Rx (Receive)
	var result_rx, result_tx map[string]interface{}
	resp_rx , err := dao.GetMetricsDao(b.influxClient, b.databases.ContainerDatabase).GetAppNetworkKByte(request, "rx_bytes")
	if err != nil {
		return result, err
	}else{
		result_rx, err = util.GetResponseConverter().InfluxConverter(resp_rx, "Rx_Network")
		if err != nil {
			return result, err
		}

	}

	//Applicatino Network Tx (Transfer)
	resp_tx , err := dao.GetMetricsDao(b.influxClient, b.databases.ContainerDatabase).GetAppNetworkKByte(request, "tx_bytes")
	if err != nil {
		return result, err
	}else{
		result_tx, err = util.GetResponseConverter().InfluxConverter(resp_tx, "Tx_Network")
		if err != nil {
			return result, err
		}

	}
	result.Name = request.ServiceName
	result.Data[0] = result_rx
	result.Data[1] = result_tx
	fmt.Println("#### Network variation :", result)
	return result, nil
}

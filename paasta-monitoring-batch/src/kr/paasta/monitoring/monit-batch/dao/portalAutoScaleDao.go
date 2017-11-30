package dao

import (
	client "github.com/influxdata/influxdb/client/v2"
	"github.com/jinzhu/gorm"
	cb "kr/paasta/monitoring/monit-batch/models/base"
	"kr/paasta/monitoring/monit-batch/util"
	mod "kr/paasta/monitoring/monit-batch/models"
	"fmt"
)


type autoScaleStruct struct {
	influxClient 	client.Client
}


func GetAutoScaleAppDao(influxClient client.Client) *autoScaleStruct{

	return &autoScaleStruct{
		influxClient: 	influxClient,
	}
}



//Auto Scale할 대상 App 목록 조회
func (f autoScaleStruct) GetAutoScaleAppList(txn *gorm.DB) ([]mod.AutoScaleConfig, cb.ErrMessage){

	appList := []mod.AutoScaleConfig{}
	status := txn.Debug().Table("auto_scaling_config").Find(&appList)
	err := util.GetError().DbCheckError(status.Error)

	if err != nil{
		fmt.Println("Error::", err )
	}

	return appList, err
}


func (b autoScaleStruct) GetAppContainersList(request mod.ZonesReq) (_ client.Response, errMsg cb.ErrMessage)  {


	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = cb.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	getContainerListSql := "select application_name, application_id, application_index, container_interface, value from container_metrics where cell_ip = '%s' and \"name\" = 'load_average'  and container_id <> '/' and time > now() - 90s";

	var q client.Query

	q = client.Query{
		Command:  fmt.Sprintf( getContainerListSql,
			request.CellIp),
		Database: request.MetricDatabase,
	}
	fmt.Println("GetCellContainerList Sql======>", q)

	resp, err := b.influxClient.Query(q)

	if err != nil{
		errLogMsg = err.Error()
	}

	return util.GetError().CheckError(*resp, err)
}



func (b autoScaleStruct) GetContainerCpuUsage(request mod.ZonesReq) (_ client.Response, errMsg cb.ErrMessage)  {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = cb.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	totalUsage := "select mean(usage) as usage from (select non_negative_derivative(mean(value),30s)/30000000000 * 100000000000 as usage from container_metrics where \"name\" = 'cpu_usage_total' and container_interface = '%s' "

	var q client.Query

	totalUsage += " and time > now() - %ss  group by time(30s));"
	q = client.Query{
		Command:  fmt.Sprintf( totalUsage,
			request.ContainerName, request.CheckTIme),
		Database: request.MetricDatabase,
	}

	fmt.Println("cpu sql:", q)
	resp, err := b.influxClient.Query(q)

	if err != nil{
		errLogMsg = err.Error()
	}

	return util.GetError().CheckError(*resp, err)
}

func (b autoScaleStruct) GetContainerMemoryUsage(request mod.ZonesReq) (_ client.Response, errMsg cb.ErrMessage)  {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = cb.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	MemorySql := "select mean(usage) as usage from (select value/app_mem * 100 as usage from container_metrics where \"name\" = 'memory_usage' and container_interface = '%s' ";
	var q client.Query
	MemorySql += " and time > now() - %ss) ;"

	q = client.Query{
		Command:  fmt.Sprintf( MemorySql,
			request.ContainerName, request.CheckTIme),
		Database: request.MetricDatabase,

	}
	fmt.Println("memory sql:", q)
	resp, err := b.influxClient.Query(q)
	if err != nil{
		errLogMsg = err.Error()
	}
	return util.GetError().CheckError(*resp, err)
}

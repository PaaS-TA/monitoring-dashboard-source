package dao

import (
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring/paas/model"
	"kr/paasta/monitoring/paas/util"
	client "github.com/influxdata/influxdb/client/v2"
	"fmt"
	"strings"
)

type ContainerDao struct {
	txn   *gorm.DB
	influxClient 	client.Client
	databases       model.Databases
}

func GetContainerDao(txn *gorm.DB, influxClient client.Client, databases model.Databases) *ContainerDao {
	return &ContainerDao{
		txn:   txn,
		influxClient: 	influxClient,
		databases:    databases,
	}
}

//Cell 목록 조회
func (b ContainerDao) GetCellList() ([]model.ZoneCellInfo, model.ErrMessage){

	cells := []model.ZoneCellInfo{}

	status := b.txn.Debug().Table("zones").Order("cell_name asc").
		Select("zones.name as zone_name, vms.id,  vms.name as cell_name, vms.ip").
		Joins("inner join vms on zones.id = vms.zone_id and vms.vm_type = 'Cel' ").Find(&cells)

	err := util.GetError().DbCheckError(status.Error)
	if err != nil{
		return nil, err
	}

	return cells, err
}

//zone 목록 조회
func (b ContainerDao) GetZoneList() ([]model.ZoneCellInfo, model.ErrMessage){

	zones := []model.ZoneCellInfo{}

	status := b.txn.Debug().Table("zones").Order("zone_name asc").
		Select("name as zone_name, id").Find(&zones)

	err := util.GetError().DbCheckError(status.Error)
	if err != nil{
		return nil, err
	}

	return zones, err
}


func (b ContainerDao) GetCellContainersList(cellIp string) (_ client.Response, errMsg model.ErrMessage)  {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	getContainerListSql := "select application_name, application_index, container_interface, value from container_metrics where cell_ip = '%s' and \"name\" = 'load_average'  and container_id <> '/' and time > now() - %s";

	var q client.Query

	q = client.Query{
		Command:  fmt.Sprintf( getContainerListSql,
			cellIp, "60s"),
		Database: b.databases.ContainerDatabase,
	}

	//fmt.Println("GetCellContainerList Sql======>", q)
	resp, err := b.influxClient.Query(q)
	if err != nil{
		errLogMsg = err.Error()
	}

	return util.GetError().CheckError(*resp, err)
}

func (b ContainerDao) GetCellSummaryData(request model.ContainerReq) (_ client.Response, errMsg model.ErrMessage)  {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = model.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	sql := request.SqlQuery

	var q client.Query

	q = client.Query{
		Command:  fmt.Sprintf( sql, request.CellIp, request.Time, request.MetricName),
		Database: b.databases.PaastaDatabase,
	}

	fmt.Println("GetCellSummaryData Sql======>", q)
	resp, err := b.influxClient.Query(q)
	if err != nil{
		errLogMsg = err.Error()
	}

	return util.GetError().CheckError(*resp, err)
}

func (b ContainerDao) GetContainerUsage(request model.ContainerReq) (_ client.Response, errMsg model.ErrMessage)  {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = model.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	sql := "select mean(totalUsage) as value from ( "
	sql += "select mean(value) as totalUsage from container_metrics "
	sql += "where \"name\" = '%s' and application_index = '%s' and application_name = '%s' and cell_ip = '%s'"

	switch request.Service {
	case model.ALARM_TYPE_CPU:
		sql = strings.Replace(sql, "mean(value)", "non_negative_derivative(mean(value),30s)/30000000000*100000000000", 1)
	case model.ALARM_TYPE_MEMORY:
		sql = strings.Replace(sql, "mean(value)", "value/app_mem * 100", 1)
	case model.ALARM_TYPE_DISK:
		sql = strings.Replace(sql, "mean(value)", "value/app_disk * 100", 1)
	}

	if request.Service == model.ALARM_TYPE_CPU {
		sql += " and time > now() - 90s  group by time(1m) "
	}else{
		sql += " and time > now() - 90s "
	}
	sql += ");"

	var q client.Query
	q = client.Query{
		Command:  fmt.Sprintf(sql, request.MetricName, request.AppIndex, request.AppName, request.CellIp),
		Database: b.databases.ContainerDatabase,
	}

	fmt.Println("GetContainerUsage Sql======>", q)
	resp, err := b.influxClient.Query(q)
	if err != nil{
		errLogMsg = err.Error()
	}

	return util.GetError().CheckError(*resp, err)
}

func (b ContainerDao) GetPaasContainerDetailUsages(request model.ContainerReq) (_ client.Response, errMsg model.ErrMessage)  {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = model.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	sql := "select mean(value) as usage from container_metrics "
	sql += "where \"name\" = '%s' and container_interface = '%s'"

	switch request.MetricName {
	case model.CON_MTR_CPU_USAGE:
		//sql = strings.Replace(sql, "mean(value)", "non_negative_derivative(mean(value))*100", 1)
		sql = "select non_negative_derivative(usage,30s)/30000000000*100000000000 as usage from (" + sql
	case model.CON_MTR_LOAD_AVG:
		sql = strings.Replace(sql, "mean(value)", "value", 1)
		sql = "select mean(usage) as usage from ( " + sql
	case model.CON_MTR_MEM_USAGE:
		sql = strings.Replace(sql, "mean(value)", "value/app_mem * 100", 1)
		sql = "select mean(usage) as usage from ( " + sql
	case model.CON_MTR_DISK_USAGE:
		sql = strings.Replace(sql, "mean(value)", "value/app_disk * 100", 1)
		sql = "select mean(usage) as usage from ( " + sql
	default:
		sql = strings.Replace(sql, "mean(value)", "value/1024", 1)
		sql = strings.Replace(sql, "container_interface", "container_id", 1)
		sql = "select mean(usage) as usage from ( " + sql
	}

	// container_interface -> container_id 사용으로 전환 - id /garden/으로 시작
	//if !strings.Contains(request.ContainerName, model.CON_MTR_ID_PREFIX){
	//	request.ContainerName = model.CON_MTR_ID_PREFIX + request.ContainerName
	//}

	switch request.MetricName {
	case model.CON_MTR_CPU_USAGE:
		if request.DefaultTimeRange != "" {
			sql += " and time > now() - %s  group by time(%s) );"
			sql =  fmt.Sprintf( sql, request.MetricName, request.ContainerName, request.DefaultTimeRange, request.GroupBy)
		} else {
			sql += " and time < now() - %s and time > now() - %s  group by time(%s) );"
			sql =   fmt.Sprintf( sql, request.MetricName, request.ContainerName, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy)
		}
	default:
		if request.DefaultTimeRange != "" {
			sql += " and time > now() - %s ) where time > now() - %s group by time(%s) ;"
			sql =  fmt.Sprintf( sql, request.MetricName, request.ContainerName, request.DefaultTimeRange, request.DefaultTimeRange, request.GroupBy)
		} else {
			sql += " and time < now() - %s and time > now() - %s ) where time < now() - %s and time > now() - %s group by time(%s) ;"
			sql =   fmt.Sprintf( sql, request.MetricName, request.ContainerName, request.TimeRangeFrom, request.TimeRangeTo, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy)
		}
	}

	var q = client.Query {
		Command:  sql,
		Database: b.databases.ContainerDatabase,
	}

	fmt.Println("GetPaasContainerDetailUsages Sql======>", q)
	resp, err := b.influxClient.Query(q)
	if err != nil{
		errLogMsg = err.Error()
	}

	return util.GetError().CheckError(*resp, err)
}

func (b ContainerDao) GetCellIdForDetail(request model.ContainerReq) (_ client.Response, errMsg model.ErrMessage)  {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = model.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	sql := "select id, value from cf_metrics where ip = '%s' and time > now() - 1m order by time desc limit 1 ;"

	var q client.Query

	q = client.Query{
		Command:  fmt.Sprintf( sql, request.CellIp),
		Database: b.databases.PaastaDatabase,
	}

	//fmt.Println("GetCellIdForDetail Sql======>", q)
	resp, err := b.influxClient.Query(q)
	if err != nil{
		errLogMsg = err.Error()
	}

	return util.GetError().CheckError(*resp, err)
}
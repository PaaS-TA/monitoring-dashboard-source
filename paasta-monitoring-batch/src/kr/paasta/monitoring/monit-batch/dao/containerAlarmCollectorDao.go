package dao

import (
	client "github.com/influxdata/influxdb/client/v2"
	"github.com/jinzhu/gorm"
	cb "kr/paasta/monitoring/monit-batch/models/base"
	mod "kr/paasta/monitoring/monit-batch/models"
	"kr/paasta/monitoring/monit-batch/util"
	"fmt"
	"time"
)

type containerAlarmStruct struct {
	influxClient 	client.Client
}

func GetContainerAlarmDao(influxClient client.Client) *containerAlarmStruct{

	return &containerAlarmStruct{
		influxClient: 	influxClient,
	}
}

func (b containerAlarmStruct) GetContainerAlarmPolicy(txn *gorm.DB) ([]mod.AlarmPolicy, cb.ErrMessage) {

	var alarmPolicy []mod.AlarmPolicy

	status := txn.Debug().Model(&alarmPolicy).Where("origin_type = ? ", cb.ORIGIN_TYPE_CONTAINER).Find(&alarmPolicy)

	err := util.GetError().DbCheckError(status.Error)
	if err != nil{
		fmt.Println("Error::", err )
		return   nil, err
	}

	return alarmPolicy, nil
}


func (b containerAlarmStruct) UpdateZoneAlarmSendDate(alarm cb.Alarm, txn *gorm.DB) (cb.ErrMessage) {

	status := txn.Debug().Model(&alarm).Where("origin_type = ? and origin_id = ? and alarm_type = ? and level = ? and app_name = ? and app_index = ? and resolve_status = '1'", alarm.OriginType, alarm.OriginId, alarm.AlarmType, alarm.Level, alarm.AppName, alarm.AppIndex ).
		Updates(map[string]interface{}{ "alarm_send_date": time.Now(), "modi_date": time.Now(), "modi_user": cb.BAT_USER})
	err := util.GetError().DbCheckError(status.Error)
	if err != nil{
		fmt.Println("Error::", err )
		return   err
	}
	return  err
}

func (f containerAlarmStruct) GetContainerIsExistAlarm( alarm cb.Alarm,  txn *gorm.DB) (bool, cb.Alarm){

	var alarmData cb.Alarm
	isNew := txn.Debug().Model(&alarmData).Where("origin_type = ? and alarm_type = ? and app_name = ? and app_index = ? and (resolve_status = '1' || resolve_status = '2') and level = ? ", alarm.OriginType, alarm.AlarmType,  alarm.AppName, alarm.AppIndex, alarm.Level ).
		Find(&alarmData).RecordNotFound()
	return isNew, alarmData
}

func (b containerAlarmStruct) GetAlarmData(alarm cb.Alarm, txn *gorm.DB) (bool, mod.Alarm) {

	var alarmData mod.Alarm
	isNew := txn.Debug().Model(&alarm).Where("origin_type = ? and alarm_type = ? and app_name = ? and app_index = ? and resolve_status = '1' and level = ? ", alarm.OriginType, alarm.AlarmType,  alarm.AppName, alarm.AppIndex, alarm.Level).
		Find(&alarmData).RecordNotFound()
	return isNew, alarmData

}

func (b containerAlarmStruct) GetAlarmTarget(orginType string, txn *gorm.DB) (mod.AlarmTarget) {

	var alarmTarget mod.AlarmTarget
	txn.Debug().Table("alarm_targets").Where("origin_type = ? ", orginType).
		Find(&alarmTarget).RecordNotFound()
	return alarmTarget

}

func (f containerAlarmStruct) CreateContainerAlarmData(alarm cb.Alarm, txn *gorm.DB) cb.ErrMessage{

	eventData := cb.Alarm{OriginId: alarm.OriginId, OriginType: alarm.OriginType, AlarmType: alarm.AlarmType, Level: alarm.Level, AppName: alarm.AppName, Ip: alarm.Ip, AlarmTitle: alarm.AlarmTitle ,AppIndex: alarm.AppIndex, ContainerName: alarm.ContainerName ,
		AppYn: alarm.AppYn, AlarmMessage: alarm.AlarmMessage , ResolveStatus: alarm.ResolveStatus, AlarmCnt: 1, RegDate: time.Now(), RegUser: "Bat",
		ModiUser: cb.BAT_USER, ModiDate: time.Now(), AlarmSendDate: time.Now()}
	status := txn.Debug().Create(&eventData)
	err := util.GetError().DbCheckError(status.Error)
	if err != nil{
		fmt.Println("Error::", err )
	}
	return  err
}

//Zone 상태 목록 조회
func (f containerAlarmStruct) GetCellList(txn *gorm.DB) ([]mod.ZoneCellInfo, cb.ErrMessage){

	cells := []mod.ZoneCellInfo{}

	status := txn.Debug().Table("zones").
		Select("zones.name as zone_name, vms.id,  vms.name as cell_name, vms.ip").
		Joins("inner join vms on zones.id = vms.zone_id and vms.vm_type = 'Cel' ").Find(&cells)

	err := util.GetError().DbCheckError(status.Error)
	if err != nil{
		return nil, err
	}

	return cells, err
}

func (b containerAlarmStruct) GetCellContainersList(request mod.ZonesReq) (_ client.Response, errMsg cb.ErrMessage)  {


	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = cb.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	getContainerListSql := "select application_name, application_index, container_interface, value from container_metrics where cell_ip = '%s' and \"name\" = 'load_average'  and container_id <> '/' and time > now() - 90s";

	var q client.Query

	q = client.Query{
		Command:  fmt.Sprintf( getContainerListSql,
			request.CellIp),
		Database: request.MetricDatabase,
	}
	//fmt.Println("GetCellContainerList Sql======>", q)

	resp, err := b.influxClient.Query(q)

	if err != nil{
		errLogMsg = err.Error()
	}

	return util.GetError().CheckError(*resp, err)
}


func (b containerAlarmStruct) GetContainerCpuUsage(request mod.ZonesReq) (_ client.Response, errMsg cb.ErrMessage)  {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = cb.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	totalUsage := "select non_negative_derivative(mean(value),30s)/30000000000 * 100000000000 as usage from container_metrics where \"name\" = 'cpu_usage_total' and container_interface = '%s' "

	var q client.Query

	totalUsage += " and time > now() - 90s  group by time(30s);"
	q = client.Query{
		Command:  fmt.Sprintf( totalUsage,
			request.ContainerName),
		Database: request.MetricDatabase,
	}

	fmt.Println("cpu sql:", q)
	resp, err := b.influxClient.Query(q)

	if err != nil{
		errLogMsg = err.Error()
	}

	return util.GetError().CheckError(*resp, err)
}

func (b containerAlarmStruct) GetContainerMemoryUsage(request mod.ZonesReq) (_ client.Response, errMsg cb.ErrMessage)  {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = cb.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	MemorySql := "select value/app_mem * 100 as usage from container_metrics where \"name\" = 'memory_usage' and container_interface = '%s' ";
	var q client.Query
	MemorySql += " and time > now() - 90s order by time desc limit 1;"

	q = client.Query{
		Command:  fmt.Sprintf( MemorySql,
			request.ContainerName),
		Database: request.MetricDatabase,

	}
	fmt.Println("memory sql:", q)
	resp, err := b.influxClient.Query(q)
	if err != nil{
		errLogMsg = err.Error()
	}
	return util.GetError().CheckError(*resp, err)
}

//Server Disk에 Mounted된 FileSystem 목록을 조회한다.
func (b containerAlarmStruct) GetContainerOvvDiskUsage(request mod.ZonesReq) (_ client.Response, errMsg cb.ErrMessage)  {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = cb.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	//httpClient, influxConfig, err :=util.GetHttpClient().GetHttpClient(ZoneDtvmetricDataSource)
	DiskSql := "select value/app_disk * 100 as usage from container_metrics where \"name\" = 'disk_usage' and container_interface = '%s' "

	var q client.Query

	DiskSql += " and time > now() - 90s order by time desc limit 1;"
	q = client.Query{
		Command:  fmt.Sprintf( DiskSql,
			request.ContainerName),
		Database: request.MetricDatabase,
	}
	fmt.Println("GetContainerDisk Sql======>", q)

	resp, err := b.influxClient.Query(q)
	if err != nil{
		errLogMsg = err.Error()
	}
	return util.GetError().CheckError(*resp, err)
}
package dao

import (
	client "github.com/influxdata/influxdb/client/v2"
	"github.com/jinzhu/gorm"
	cb "kr/paasta/monitoring/monit-batch/models/base"
	"kr/paasta/monitoring/monit-batch/util"
	mod "kr/paasta/monitoring/monit-batch/models"
	"fmt"
)

type paasTaAlarmStruct struct {
	influxClient 	client.Client
}


func GetPaasTaAlarmDao(influxClient client.Client) *paasTaAlarmStruct{

	return &paasTaAlarmStruct{
		influxClient: 	influxClient,
	}
}

//Server 상태 목록 조회
func (f paasTaAlarmStruct) GetPaaSTaList(txn *gorm.DB) ([]mod.Vm, cb.ErrMessage){

	vms := []mod.Vm{}

	status := txn.Debug().Find(&vms)
	err := util.GetError().DbCheckError(status.Error)

	if err != nil{
		fmt.Println("Error::", err )
	}

	return vms, err
}

func (b paasTaAlarmStruct) GetPaastaAlarmPolicy(txn *gorm.DB) ([]mod.AlarmPolicy, cb.ErrMessage) {

	var alarmPolicy []mod.AlarmPolicy

	status := txn.Debug().Model(&alarmPolicy).Where("origin_type = ? ", cb.ORIGIN_TYPE_PAASTA).Find(&alarmPolicy)

	err := util.GetError().DbCheckError(status.Error)
	if err != nil{
		fmt.Println("Error::", err )
		return   nil, err
	}

	return alarmPolicy, nil
}


func (b paasTaAlarmStruct) GetPaasTaCpuUsage(request mod.VmReq) (_ client.Response, errMsg cb.ErrMessage)  {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = cb.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	cpuUsageSql := "select mean(value) as usage  from cf_metrics where ip = '%s' and metricname =~ /cpuStats.*/ and time > now() - 2m group by time(1m) order by time desc limit 1"
	var q client.Query

	q = client.Query{
		Command:  fmt.Sprintf( cpuUsageSql , request.Ip),
		Database: request.MetricDatabase,
	}

	fmt.Println("CPU Sql==>%s", q)

	resp, err := b.influxClient.Query(q)

	if err != nil{
		errLogMsg = err.Error()
	}
	return util.GetError().CheckError(*resp, err)
}

func (b paasTaAlarmStruct) GetPaasTaMemoryUsage(request mod.VmReq) (_ client.Response, errMsg cb.ErrMessage)  {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = cb.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	memoryUsageSql := "select value as usage from cf_metrics where ip = '%s' and metricname = 'memoryStats.UsedPercent' and time > now() - 2m  order by time desc limit 1 "
	var q client.Query

	q = client.Query{
		Command:  fmt.Sprintf( memoryUsageSql , request.Ip),
		Database: request.MetricDatabase,
	}

	resp, err := b.influxClient.Query(q)

	fmt.Println("Memory Sql==>%s", q)
	if err != nil{
		errLogMsg = err.Error()
	}
	return util.GetError().CheckError(*resp, err)
}



func (b paasTaAlarmStruct) GetPaasTaDiskUsage(request mod.VmReq) (_ client.Response, errMsg cb.ErrMessage)  {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = cb.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	memoryUsageSql := "select value as usage from cf_metrics where ip = '%s' and metricname = 'diskStats.Usage' and time > now() - 2m  order by time desc limit 1"
	var q client.Query

	q = client.Query{
		Command:  fmt.Sprintf( memoryUsageSql , request.Ip),
		Database: request.MetricDatabase,
	}
	fmt.Println("Disk Sql==>", q)
	resp, err := b.influxClient.Query(q)

	if err != nil{
		errLogMsg = err.Error()
	}
	return util.GetError().CheckError(*resp, err)
}


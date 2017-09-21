package dao

import (
	client "github.com/influxdata/influxdb/client/v2"
	cb "kr/paasta/monitoring/monit-batch/models/base"
	mod "kr/paasta/monitoring/monit-batch/models"
	"kr/paasta/monitoring/monit-batch/util"
	"fmt"
	"github.com/jinzhu/gorm"
)

type boshAlarmStruct struct {
	influxClient 	client.Client
}


func GetBoshAlarmDao(influxClient client.Client) *boshAlarmStruct{

	return &boshAlarmStruct{
		influxClient: 	influxClient,
	}
}


func (b boshAlarmStruct) GetBoshAlarmPolicy(txn *gorm.DB) ([]mod.AlarmPolicy, cb.ErrMessage) {

	var alarmPolicy []mod.AlarmPolicy

	status := txn.Debug().Model(&alarmPolicy).Where("origin_type = ? ", cb.ORIGIN_TYPE_BOSH).Find(&alarmPolicy)

	err := util.GetError().DbCheckError(status.Error)
	if err != nil{
		fmt.Println("Error::", err )
		return   nil, err
	}

	return alarmPolicy, nil
}

func (b boshAlarmStruct) GetBoshCpuUsage(request mod.BoshReq) (_ client.Response, errMsg cb.ErrMessage)  {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = cb.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	cpuUsageSql := "select mean(value) as usage  from bosh_metrics where origin = '%s' and metricname =~ /cpuStats.core*/ and time > now() - 2m group by time(1m) order by time desc limit 1"
	var q client.Query

	q = client.Query{
		Command:  fmt.Sprintf( cpuUsageSql , request.ServiceName),
		Database: request.MetricDatabase,
	}

	fmt.Println("CPU Sql======>", q)


	resp, err := b.influxClient.Query(q)

	if err != nil{
		errLogMsg = err.Error()
	}
	return util.GetError().CheckError(*resp, err)
}

func (b boshAlarmStruct) GetBoshMemoryUsage(request mod.BoshReq) (_ client.Response, errMsg cb.ErrMessage)  {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = cb.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	memoryUsageSql := "select value as usage from bosh_metrics where origin = '%s' and metricname = 'memoryStats.UsedPercent' and time > now() - 2m order by time desc limit 1"
	var q client.Query

	q = client.Query{
		Command:  fmt.Sprintf( memoryUsageSql , request.ServiceName),
		Database: request.MetricDatabase,
	}

	resp, err := b.influxClient.Query(q)

	fmt.Println("Memory Sql==>", q)

	if err != nil{
		errLogMsg = err.Error()
	}
	return util.GetError().CheckError(*resp, err)
}

func (b boshAlarmStruct) GetBoshDiskUsage(request mod.BoshReq) (_ client.Response, errMsg cb.ErrMessage)  {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = cb.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	memoryUsageSql := "select value as usage from bosh_metrics where origin = '%s' and metricname = 'diskStats.Usage' and time > now() - 2m order by time desc limit 1"
	var q client.Query

	q = client.Query{
		Command:  fmt.Sprintf( memoryUsageSql , request.ServiceName),
		Database: request.MetricDatabase,
	}

	resp, err := b.influxClient.Query(q)
	fmt.Println("Disk Sql==>%s", q)
	if err != nil{
		errLogMsg = err.Error()
	}
	return util.GetError().CheckError(*resp, err)
}

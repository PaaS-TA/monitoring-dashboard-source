package dao

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/jinzhu/gorm"
	mod "kr/paasta/monitoring-batch/model"
	cb "kr/paasta/monitoring-batch/model/base"
	"kr/paasta/monitoring-batch/util"
	"strconv"
)

type paasTaAlarmStruct struct {
	influxClient client.Client
}

func GetPaasTaAlarmDao(influxClient client.Client) *paasTaAlarmStruct {

	return &paasTaAlarmStruct{
		influxClient: influxClient,
	}
}

//Server 상태 목록 조회
func (f paasTaAlarmStruct) GetPaaSTaList(txn *gorm.DB) ([]mod.Vm, cb.ErrMessage) {

	vms := []mod.Vm{}

	status := txn.Debug().Find(&vms)
	err := util.GetError().DbCheckError(status.Error)

	if err != nil {
		fmt.Println("Error::", err)
	}

	return vms, err
}

func (b paasTaAlarmStruct) GetPaastaAlarmPolicy(txn *gorm.DB) ([]mod.AlarmPolicy, cb.ErrMessage) {

	var alarmPolicy []mod.AlarmPolicy

	status := txn.Debug().Model(&alarmPolicy).Where("origin_type = ? ", cb.ORIGIN_TYPE_PAASTA).Find(&alarmPolicy)

	err := util.GetError().DbCheckError(status.Error)
	if err != nil {
		fmt.Println("Error::", err)
		return nil, err
	}

	return alarmPolicy, nil
}

func (b paasTaAlarmStruct) GetPaasTaCpuUsage(request mod.VmReq) (_ client.Response, errMsg cb.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = cb.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	// alarm measure time default : 120s
	measureTime := "120s"

	for _, value := range request.MeasureTimeList {
		if value.Item == cb.ALARM_TYPE_CPU {
			measureTime = strconv.Itoa(value.MeasureTime) + "s"
		}
	}

	cpuUsageSql := "select mean(value) as usage  from cf_metrics where ip = '%s' and metricname =~ /cpuStats.*/ and time > now() - %s"
	var q client.Query

	q = client.Query{
		Command:  fmt.Sprintf(cpuUsageSql, request.Ip, measureTime),
		Database: request.MetricDatabase,
	}

	fmt.Println("CPU Sql==>%s", q)

	resp, err := b.influxClient.Query(q)

	if err != nil {
		errLogMsg = err.Error()
	}
	return util.GetError().CheckError(*resp, err)
}

func (b paasTaAlarmStruct) GetPaasTaMemoryUsage(request mod.VmReq) (_ client.Response, errMsg cb.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = cb.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	// alarm measure time default : 120s
	measureTime := "120s"

	for _, value := range request.MeasureTimeList {
		if value.Item == cb.ALARM_TYPE_MEMORY {
			measureTime = strconv.Itoa(value.MeasureTime) + "s"
		}
	}

	memoryTotalSql := "select mean(value) as usage from cf_metrics where ip = '%s' and metricname = 'memoryStats.TotalMemory' and time > now() - %s ;"
	memoryFreeSql := "select mean(value) as memUsage from cf_metrics where ip = '%s' and metricname = 'memoryStats.FreeMemory' and time > now() - %s ;"

	var q client.Query

	q = client.Query{
		Command:  fmt.Sprintf(memoryTotalSql+memoryFreeSql, request.Ip, measureTime, request.Ip, measureTime),
		Database: request.MetricDatabase,
	}

	resp, err := b.influxClient.Query(q)

	fmt.Println("Memory Sql==>%s", q)
	if err != nil {
		errLogMsg = err.Error()
	}
	return util.GetError().CheckError(*resp, err)
}

func (b paasTaAlarmStruct) GetPaasTaDiskUsage(request mod.VmReq) (_ client.Response, errMsg cb.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = cb.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	// alarm measure time default : 120s
	measureTime := "120s"

	for _, value := range request.MeasureTimeList {
		if value.Item == cb.ALARM_TYPE_DISK {
			measureTime = strconv.Itoa(value.MeasureTime) + "s"
		}
	}

	diskUsageSql := "select mean(value) as usage from cf_metrics where ip = '%s' and metricname = 'diskStats./.Usage' and time > now() - %s"
	var q client.Query

	q = client.Query{
		Command:  fmt.Sprintf(diskUsageSql, request.Ip, measureTime),
		Database: request.MetricDatabase,
	}
	fmt.Println("Disk Sql==>", q)
	resp, err := b.influxClient.Query(q)

	if err != nil {
		errLogMsg = err.Error()
	}
	return util.GetError().CheckError(*resp, err)
}

func (b paasTaAlarmStruct) GetPaasTaRootDiskUsage(request mod.VmReq) (_ client.Response, errMsg cb.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = cb.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	// alarm measure time default : 120s
	measureTime := "120s"

	for _, value := range request.MeasureTimeList {
		if value.Item == cb.ALARM_TYPE_DISK {
			measureTime = strconv.Itoa(value.MeasureTime) + "s"
		}
	}

	diskUsageSql := "select mean(value) as usage from cf_metrics where ip = '%s' and metricname = 'diskStats./var/vcap/data.Usage' and time > now() - %s"
	var q client.Query

	q = client.Query{
		Command:  fmt.Sprintf(diskUsageSql, request.Ip, measureTime),
		Database: request.MetricDatabase,
	}
	fmt.Println("Disk Sql==>", q)
	resp, err := b.influxClient.Query(q)

	if err != nil {
		errLogMsg = err.Error()
	}
	return util.GetError().CheckError(*resp, err)
}

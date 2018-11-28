package dao

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring/paas/model"
	"kr/paasta/monitoring/paas/util"
	"strings"
)

type PaasDao struct {
	txn                      *gorm.DB
	influxClient             client.Client
	databases                model.Databases
	InfraDtvmetricDataSource string
}

func GetPaasDao(txn *gorm.DB, influxClient client.Client, ds string) *PaasDao {
	return &PaasDao{
		txn:                      txn,
		influxClient:             influxClient,
		InfraDtvmetricDataSource: ds,
	}
}

var logger seelog.LoggerInterface

func (p *PaasDao) GetPaasVms() ([]model.ResultVm, model.ErrMessage) {

	var vms []model.ResultVm

	status := p.txn.Debug().Table("vms").Order("name").Find(&vms)
	err := util.GetError().DbCheckError(status.Error)

	return vms, err
}

func (p *PaasDao) GetPaasCfMetrics(ip string) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("GetPaasCfMetrics error :", errLogMsg)
			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	sql := "select time, id, ip, metricname, origin, value from cf_metrics " +
		"where time > now() - 1m and ip = '%s' group by metricname order by time desc limit 1"

	var q client.Query
	q = client.Query{
		Command:  fmt.Sprintf(sql, ip),
		Database: p.InfraDtvmetricDataSource,
	}

	resp, err := p.influxClient.Query(q)

	if err != nil {
		errLogMsg = err.Error()
	}

	return util.GetError().CheckError(*resp, err)
}

func (p *PaasDao) GetPaasTopProcessList(request model.PaasRequest) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("GetPaasTopProcessList error :", errLogMsg)
			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	sql := "select ip, mem_usage, proc_name, proc_pid from cf_process_metrics " +
		"where time > now() - 1m and id = '%s' group by proc_name order by time desc limit 1"

	var q client.Query
	q = client.Query{
		Command:  fmt.Sprintf(sql, request.Id),
		Database: p.InfraDtvmetricDataSource,
	}

	resp, err := p.influxClient.Query(q)

	if err != nil {
		errLogMsg = err.Error()
	}

	return util.GetError().CheckError(*resp, err)
}

func (p *PaasDao) GetMetricUsageByTime(request model.PaasRequest) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string

	defer func() {
		if r := recover(); r != nil {
			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	field := "mean(value)"
	if request.IsNonNegativeDerivative {
		field = fmt.Sprintf("non_negative_derivative(%s)", field)
	}
	if request.IsRespondKb {
		field = field + "/1024"
	}

	sql := "select %s as usage from cf_metrics where id = '%s' and metricname = '%s'"
	if request.IsLikeQuery {
		sql = strings.Replace(sql, "metricname = '%s'", "metricname =~ /%s/", 1)
	}

	var q client.Query
	if request.DefaultTimeRange != "" {
		sql += " and time > now() - %s group by time(%s)"
		q = client.Query{
			Command:  fmt.Sprintf(sql, field, request.Id, request.MetricName, request.DefaultTimeRange, request.GroupBy),
			Database: p.InfraDtvmetricDataSource,
		}
	} else {
		sql += " and time < now() - %s and time > now() - %s group by time(%s)"
		q = client.Query{
			Command:  fmt.Sprintf(sql, field, request.Id, request.MetricName, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: p.InfraDtvmetricDataSource,
		}
	}

	resp, err := p.influxClient.Query(q)

	return util.GetError().CheckError(*resp, err)
}

func (p *PaasDao) GetPaasMemoryUsage(request model.PaasRequest)(_ client.Response, errMsg model.ErrMessage){
	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()
	memoryTotalSql := "select mean(value) as usage from cf_metrics where id = '%s' and metricname = '%s' "
	memoryFreeSql := "select mean(value) as usage from cf_metrics where id = '%s' and metricname = '%s' "
	var q client.Query
	if request.DefaultTimeRange != "" {
		memoryTotalSql += " and time > now() - %s group by time(%s);"
		memoryFreeSql += " and time > now() - %s group by time(%s);"
		q = client.Query{
			Command: fmt.Sprintf(memoryTotalSql + memoryFreeSql,
				request.Id, request.Args.(model.MemoryMetricArg).NameMemoryTotal, request.DefaultTimeRange, request.GroupBy,
				request.Id, request.Args.(model.MemoryMetricArg).NameMemoryFree, request.DefaultTimeRange, request.GroupBy),
			Database: p.InfraDtvmetricDataSource,
		}
	} else {
		memoryTotalSql += " and time < now() - %s and time > now() - %s group by time(%s);"
		memoryFreeSql  += " and time < now() - %s and time > now() - %s group by time(%s);"
		q = client.Query{
			Command:  fmt.Sprintf(memoryTotalSql + memoryFreeSql,
				request.Id, request.Args.(model.MemoryMetricArg).NameMemoryTotal, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy,
				request.Id, request.Args.(model.MemoryMetricArg).NameMemoryFree, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: p.InfraDtvmetricDataSource,
		}
	}
	resp, err := p.influxClient.Query(q)
	return util.GetError().CheckError(*resp, err)
}

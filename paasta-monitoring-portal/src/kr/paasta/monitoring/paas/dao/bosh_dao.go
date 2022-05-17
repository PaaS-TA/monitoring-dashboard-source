package dao

import (
	"fmt"
	"github.com/influxdata/influxdb1-client/v2"
	"github.com/jinzhu/gorm"
	"monitoring-portal/paas/model"
	"monitoring-portal/paas/util"
	"strings"
)

type BoshStatusDao struct {
	txn          *gorm.DB
	influxClient client.Client
	databases    model.Databases
}

func GetBoshStatusDao(txn *gorm.DB, influxClient client.Client, databases model.Databases) *BoshStatusDao {
	return &BoshStatusDao{
		txn:          txn,
		influxClient: influxClient,
		databases:    databases,
	}
}

func (b BoshStatusDao) GetBoshTopprocessList(request model.BoshSummaryReq) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	getBoshTopprocessListSql := "select proc_name as process_name, time, proc_index, proc_pid, mem_usage from bosh_process_metrics where id = '%s' and time > now() - %s order by time desc"

	var q client.Query

	q = client.Query{
		Command:  fmt.Sprintf(getBoshTopprocessListSql, request.Id, "1m"),
		Database: b.databases.BoshDatabase,
	}

	fmt.Println("GetBoshTopprocessList Sql======>", q)
	resp, err := b.influxClient.Query(q)
	if err != nil {
		errLogMsg = err.Error()
	}

	return util.GetError().CheckError(*resp, err)
}

func (b BoshStatusDao) GetDynamicBoshSummaryData(request model.BoshSummaryReq) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	sql := request.SqlQuery

	var q client.Query

	q = client.Query{
		Command:  fmt.Sprintf(sql, request.Id, request.Time, request.MetricName),
		Database: b.databases.BoshDatabase,
	}

	fmt.Println("GetDynamicBoshSummaryData Sql======>", q)
	resp, err := b.influxClient.Query(q)
	if err != nil {
		errLogMsg = err.Error()
	}

	return util.GetError().CheckError(*resp, err)
}

func (b BoshStatusDao) GetBoshCpuUsageList(request model.BoshDetailReq) (_ client.Response, errMsg model.ErrMessage) {

	fmt.Println("GetBoshCpuUsageList Request API parameter =========+>", request)

	var errLogMsg string

	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	cpuUsageSql := "select mean(value) as usage from bosh_metrics where id = '%s' and metricname =~ /%s/ "

	var q client.Query
	if request.DefaultTimeRange != "" {

		cpuUsageSql += " and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command: fmt.Sprintf(cpuUsageSql,
				request.Id, request.MetricName, request.DefaultTimeRange, request.GroupBy),
			Database: b.databases.BoshDatabase,
		}
	} else {

		cpuUsageSql += " and time < now() - %s and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command: fmt.Sprintf(cpuUsageSql,
				request.Id, request.MetricName, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: b.databases.BoshDatabase,
		}
	}

	fmt.Println("GetBoshCpuUsageList Sql======>", q)
	resp, err := b.influxClient.Query(q)

	return util.GetError().CheckError(*resp, err)
}

func (b BoshStatusDao) GetBoshMemUsageList(request model.BoshDetailReq) (_ client.Response, errMsg model.ErrMessage) {

	fmt.Println("GetBoshMemUsageList Request API parameter =========+>", request)

	var errLogMsg string

	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	totalMemorySql := "select mean(value) as usage from bosh_metrics where id = '%s' and metricname = 'memoryStats.TotalMemory'"
	freeMemorySql := "select mean(value) as memUsage from bosh_metrics  where id = '%s' and metricname = 'memoryStats.FreeMemory'"

	var q client.Query
	if request.DefaultTimeRange != "" {

		totalMemorySql += " and time > now() - %s  group by time(%s);"
		freeMemorySql += " and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command: fmt.Sprintf(totalMemorySql+freeMemorySql,
				request.Id, request.DefaultTimeRange, request.GroupBy,
				request.Id, request.DefaultTimeRange, request.GroupBy),
			Database: b.databases.BoshDatabase,
		}
	} else {

		totalMemorySql += " and time < now() - %s and time > now() - %s  group by time(%s);"
		freeMemorySql += " and time < now() - %s and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command: fmt.Sprintf(totalMemorySql+freeMemorySql,
				request.Id, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy,
				request.Id, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: b.databases.BoshDatabase,
		}
	}

	fmt.Println("GetBoshMemUsageList Sql======>", q)
	resp, err := b.influxClient.Query(q)

	return util.GetError().CheckError(*resp, err)
}

func (b BoshStatusDao) GetBoshDetailList(request model.BoshDetailReq) (_ client.Response, errMsg model.ErrMessage) {

	fmt.Println("GetBoshDetailList Request API parameter =========+>", request)

	var errLogMsg string

	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	detailSql := "select mean(value) as usage from bosh_metrics where id = '%s' and metricname = '%s' "

	if strings.Contains(request.MetricName, "bytesRecv") || strings.Contains(request.MetricName, "bytesSent") {
		detailSql = strings.Replace(detailSql, "mean(value)", "non_negative_derivative(mean(value))/1024", 1)
	} else if strings.Contains(request.MetricName, "packetRecv") || strings.Contains(request.MetricName, "packetSent") {
		detailSql = strings.Replace(detailSql, "mean(value)", "non_negative_derivative(mean(value))", 1)
	} else if strings.Contains(request.MetricName, "err") || strings.Contains(request.MetricName, "drop") {
		detailSql = strings.Replace(detailSql, "mean(value)", "non_negative_derivative(sum(value))", 1)
	} else if strings.Contains(request.MetricName, "diskIOStats") {
		detailSql = strings.Replace(detailSql, "mean(value)", "non_negative_derivative(mean(value),1m)/1024", 1)
		detailSql = strings.Replace(detailSql, "metricname = '%s'", "metricname =~ /%s/", 1)
	} else {
		if request.IsConvertKb {
			detailSql = strings.Replace(detailSql, "mean(value)", "mean(value)/1024", 1)
		}
	}

	if request.DefaultTimeRange != "" {
		detailSql += " and time > now() - %s  group by time(%s);"
		detailSql = fmt.Sprintf(detailSql, request.Id, request.MetricName, request.DefaultTimeRange, request.GroupBy)
	} else {
		detailSql += " and time < now() - %s and time > now() - %s  group by time(%s);"
		detailSql = fmt.Sprintf(detailSql, request.Id, request.MetricName, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy)
	}

	var q = client.Query{
		Command:  detailSql,
		Database: b.databases.BoshDatabase,
	}

	fmt.Println("GetBoshDetailList Sql======>", q)
	resp, err := b.influxClient.Query(q)

	return util.GetError().CheckError(*resp, err)
}

func (b BoshStatusDao) GetBoshId(deployName string) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	//sql := "select id, value from bosh_metrics where deployment = '%s' and time > now() - 1m order by time desc limit 1 ;"
	sql := "select id, value from bosh_metrics where deployment = '" + deployName + "' and time > now() - 1m order by time desc limit 1 ;"

	var q client.Query

	q = client.Query{
		//Command:  fmt.Sprintf(sql, deployName),
		Command:  sql,
		Database: b.databases.BoshDatabase,
	}

	fmt.Println("sql : " + sql)
	fmt.Printf("b.databases.BoshDatabase : %v\n", b.databases.BoshDatabase)

	resp, err := b.influxClient.Query(q)
	if err != nil {
		errLogMsg = err.Error()
	}

	return util.GetError().CheckError(*resp, err)
}

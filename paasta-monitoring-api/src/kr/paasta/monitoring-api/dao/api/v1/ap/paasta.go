package ap

import (
	"fmt"
	"github.com/influxdata/influxdb1-client/v2"
	"gorm.io/gorm"
	models "paasta-monitoring-api/models/api/v1"
	"strings"
)

type PaastaDao struct {
	DbInfo         *gorm.DB
	InfluxDbClient models.InfluxDbClient
}

func GetPaastaDao(DbInfo *gorm.DB, InfluxDbClient models.InfluxDbClient) *PaastaDao {
	return &PaastaDao{
		DbInfo:         DbInfo,
		InfluxDbClient: InfluxDbClient,
	}
}

func (p *PaastaDao) GetPaastaInfoList() ([]models.Paasta, error) {
	var response []models.Paasta
	results := p.DbInfo.Debug().Table("vms").
		Select("*").
		Order("name").
		Find(&response)

	if results.Error != nil {
		fmt.Println(results.Error)
		return response, results.Error
	}

	return response, results.Error
}

func (p *PaastaDao) GetPaastaCfMetrics(ip string) (*client.Response, models.ErrMessage) {
	var errLogMsg string
	var errMsg models.ErrMessage
	defer func() {
		if r := recover(); r != nil {
			errMsg = models.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	sql := "select time, id, ip, metricname, origin, value from cf_metrics " +
		"where time > now() - 2m and ip = '%s' group by metricname order by time desc limit 1"

	var q client.Query
	q = client.Query{
		Command:  fmt.Sprintf(sql, ip),
		Database: p.InfluxDbClient.DbName.PaastaDatabase,
	}

	resp, err := p.InfluxDbClient.HttpClient.Query(q)
	if err != nil {
		errLogMsg = err.Error()
		return resp, errMsg
	}
	fmt.Println("GetPaastaCfMetrics resp======>", resp)

	return resp, nil
}

func (p *PaastaDao) GetPaastaSummary(boshSummary models.BoshSummary) (*client.Response, models.ErrMessage) {
	var errLogMsg string
	var errMsg models.ErrMessage
	defer func() {
		if r := recover(); r != nil {
			errMsg = models.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	getBoshSummarySql := boshSummary.SqlQuery
	q := client.Query{
		Command:  fmt.Sprintf(getBoshSummarySql, boshSummary.UUID, boshSummary.Time, boshSummary.MetricName),
		Database: p.InfluxDbClient.DbName.BoshDatabase,
	}

	resp, err := p.InfluxDbClient.HttpClient.Query(q)
	if err != nil {
		errLogMsg = err.Error()
		return resp, errMsg
	}
	fmt.Println("GetBoshProcessByMemory resp======>", resp)

	return resp, nil
}

func (p *PaastaDao) GetPaastaProcessByMemory(paastaProcess models.PaastaProcess) (*client.Response, error) {
	getPaastaTopprocessListSql := "select ip, mem_usage, proc_name, proc_pid from cf_process_metrics " +
		"where id = '%s' and time > now() - %s group by proc_name order by time desc limit 1"

	q := client.Query{
		Command:  fmt.Sprintf(getPaastaTopprocessListSql, paastaProcess.UUID, "1m"),
		Database: p.InfluxDbClient.DbName.PaastaDatabase,
	}
	fmt.Println("GetPaastaProcessByMemory Sql======>", q)

	resp, err := p.InfluxDbClient.HttpClient.Query(q)
	if err != nil {
		return resp, err
	}
	fmt.Println("GetPaastaProcessByMemory resp======>", resp)

	return resp, err
}

func (p *PaastaDao) GetPaastaCommonUsageByTime(paastaChart models.PaastaChart) (*client.Response, error) {
	field := "mean(value)"
	if paastaChart.IsNonNegativeDerivative {
		field = fmt.Sprintf("non_negative_derivative(%s)", field)
	}
	if paastaChart.IsRespondKb {
		field = field + "/1024"
	}

	sql := "select %s as usage from cf_metrics where id = '%s' and metricname = '%s'"
	if paastaChart.IsLikeQuery {
		sql = strings.Replace(sql, "metricname = '%s'", "metricname =~ /%s/", 1)
	}

	var q client.Query
	if paastaChart.DefaultTimeRange != "" {
		sql += " and time > now() - %s group by time(%s)"
		q = client.Query{
			Command:  fmt.Sprintf(sql, field, paastaChart.UUID, paastaChart.MetricName, paastaChart.DefaultTimeRange, paastaChart.GroupBy),
			Database: p.InfluxDbClient.DbName.PaastaDatabase,
		}
	} else {
		sql += " and time < now() - %s and time > now() - %s group by time(%s)"
		q = client.Query{
			Command:  fmt.Sprintf(sql, field, paastaChart.UUID, paastaChart.MetricName, paastaChart.TimeRangeFrom, paastaChart.TimeRangeTo, paastaChart.GroupBy),
			Database: p.InfluxDbClient.DbName.PaastaDatabase,
		}
	}

	response, err := p.InfluxDbClient.HttpClient.Query(q)
	if err != nil {
		return response, err
	}
	return response, err
}

func (p *PaastaDao) GetPaastaMemoryUsage(paastaChart models.PaastaChart) (*client.Response, error) {
	memoryTotalSql := "select mean(value) as usage from cf_metrics where id = '%s' and metricname = '%s' "
	memoryFreeSql := "select mean(value) as usage from cf_metrics where id = '%s' and metricname = '%s' "
	var q client.Query
	if paastaChart.DefaultTimeRange != "" {
		memoryTotalSql += " and time > now() - %s group by time(%s);"
		memoryFreeSql += " and time > now() - %s group by time(%s);"
		q = client.Query{
			Command: fmt.Sprintf(memoryTotalSql+memoryFreeSql,
				paastaChart.UUID, models.METRIC_NAME_TOTAL_MEMORY, paastaChart.DefaultTimeRange, paastaChart.GroupBy,
				paastaChart.UUID, models.METRIC_NAME_FREE_MEMORY, paastaChart.DefaultTimeRange, paastaChart.GroupBy),
			Database: p.InfluxDbClient.DbName.PaastaDatabase,
		}
	} else {
		memoryTotalSql += " and time < now() - %s and time > now() - %s group by time(%s);"
		memoryFreeSql += " and time < now() - %s and time > now() - %s group by time(%s);"
		q = client.Query{
			Command: fmt.Sprintf(memoryTotalSql+memoryFreeSql,
				paastaChart.UUID, models.METRIC_NAME_TOTAL_MEMORY, paastaChart.TimeRangeFrom, paastaChart.TimeRangeTo, paastaChart.GroupBy,
				paastaChart.UUID, models.METRIC_NAME_FREE_MEMORY, paastaChart.TimeRangeFrom, paastaChart.TimeRangeTo, paastaChart.GroupBy),
			Database: p.InfluxDbClient.DbName.PaastaDatabase,
		}
	}
	response, err := p.InfluxDbClient.HttpClient.Query(q)
	if err != nil {
		return response, err
	}
	return response, err
}

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

func (p *PaastaDao) GetPaastaCpuUsageList(boshChart models.BoshChart) (*client.Response, error) {
	cpuUsageSql := "select mean(value) as usage from bosh_metrics where id = '%s' and metricname =~ /%s/ "

	if boshChart.DefaultTimeRange != "" {
		cpuUsageSql += " and time > now() - %s  group by time(%s);"
		cpuUsageSql = fmt.Sprintf(cpuUsageSql, boshChart.UUID, boshChart.MetricName, boshChart.DefaultTimeRange, boshChart.GroupBy)

	} else {
		cpuUsageSql += " and time < now() - %s and time > now() - %s  group by time(%s);"
		cpuUsageSql = fmt.Sprintf(cpuUsageSql, boshChart.UUID, boshChart.MetricName, boshChart.TimeRangeFrom, boshChart.TimeRangeTo, boshChart.GroupBy)
	}
	fmt.Println("GetBoshCpuUsageList Sql======>", cpuUsageSql)

	q := client.Query{
		Command:  cpuUsageSql,
		Database: p.InfluxDbClient.DbName.BoshDatabase,
	}
	resp, err := p.InfluxDbClient.HttpClient.Query(q)
	if err != nil {
		return resp, err
	}
	fmt.Println("GetBoshCpuUsageList resp======>", resp)

	return resp, err
}

func (p *PaastaDao) GetPaastaCpuLoadList(boshChart models.BoshChart) (*client.Response, error) {
	cpuLoadSql := "select mean(value) as usage from bosh_metrics where id = '%s' and metricname = '%s' "

	if boshChart.DefaultTimeRange != "" {
		cpuLoadSql += " and time > now() - %s  group by time(%s);"
		cpuLoadSql = fmt.Sprintf(cpuLoadSql, boshChart.UUID, boshChart.MetricName, boshChart.DefaultTimeRange, boshChart.GroupBy)
	} else {
		cpuLoadSql += " and time < now() - %s and time > now() - %s  group by time(%s);"
		cpuLoadSql = fmt.Sprintf(cpuLoadSql, boshChart.UUID, boshChart.MetricName, boshChart.TimeRangeFrom, boshChart.TimeRangeTo, boshChart.GroupBy)
	}
	fmt.Println("GetBoshCpuLoadList Sql======>", cpuLoadSql)

	q := client.Query{
		Command:  cpuLoadSql,
		Database: p.InfluxDbClient.DbName.BoshDatabase,
	}
	resp, err := p.InfluxDbClient.HttpClient.Query(q)
	if err != nil {
		return resp, err
	}
	fmt.Println("GetBoshCpuLoadList resp======>", resp)

	return resp, err
}

func (p *PaastaDao) GetPaastaMemoryUsageList(boshChart models.BoshChart) (*client.Response, error) {
	totalMemorySql := "select mean(value) as usage from bosh_metrics where id = '%s' and metricname = 'memoryStats.TotalMemory'"
	freeMemorySql := "select mean(value) as memUsage from bosh_metrics  where id = '%s' and metricname = 'memoryStats.FreeMemory'"
	var Sql string

	var q client.Query
	if boshChart.DefaultTimeRange != "" {
		totalMemorySql += " and time > now() - %s  group by time(%s);"
		freeMemorySql += " and time > now() - %s  group by time(%s);"
		Sql = fmt.Sprintf(totalMemorySql+freeMemorySql,
			boshChart.UUID, boshChart.DefaultTimeRange, boshChart.GroupBy,
			boshChart.UUID, boshChart.DefaultTimeRange, boshChart.GroupBy)
	} else {
		totalMemorySql += " and time < now() - %s and time > now() - %s  group by time(%s);"
		freeMemorySql += " and time < now() - %s and time > now() - %s  group by time(%s);"
		Sql = fmt.Sprintf(totalMemorySql+freeMemorySql,
			boshChart.UUID, boshChart.TimeRangeFrom, boshChart.TimeRangeTo, boshChart.GroupBy,
			boshChart.UUID, boshChart.TimeRangeFrom, boshChart.TimeRangeTo, boshChart.GroupBy)
	}

	q = client.Query{
		Command:  Sql,
		Database: p.InfluxDbClient.DbName.BoshDatabase,
	}
	fmt.Println("GetBoshMemUsageList Sql======>", q)
	resp, err := p.InfluxDbClient.HttpClient.Query(q)

	return resp, err
}

func (p *PaastaDao) GetPaastaDiskUsageList(boshChart models.BoshChart) (*client.Response, error) {
	boshDiskUsageSql := "select mean(value) as usage from bosh_metrics where id = '%s' and metricname = '%s' "

	if boshChart.DefaultTimeRange != "" {
		boshDiskUsageSql += " and time > now() - %s  group by time(%s);"
		boshDiskUsageSql = fmt.Sprintf(boshDiskUsageSql, boshChart.UUID, boshChart.MetricName, boshChart.DefaultTimeRange, boshChart.GroupBy)
	} else {
		boshDiskUsageSql += " and time < now() - %s and time > now() - %s  group by time(%s);"
		boshDiskUsageSql = fmt.Sprintf(boshDiskUsageSql, boshChart.UUID, boshChart.MetricName, boshChart.TimeRangeFrom, boshChart.TimeRangeTo, boshChart.GroupBy)
	}
	fmt.Println("GetBoshCpuLoadList Sql======>", boshDiskUsageSql)

	q := client.Query{
		Command:  boshDiskUsageSql,
		Database: p.InfluxDbClient.DbName.BoshDatabase,
	}
	resp, err := p.InfluxDbClient.HttpClient.Query(q)
	if err != nil {
		return resp, err
	}
	fmt.Println("GetBoshCpuLoadList resp======>", resp)

	return resp, err
}

func (p *PaastaDao) GetPaastaDiskIoList(boshChart models.BoshChart) (*client.Response, error) {
	boshDiskIoSql := "select mean(value) as usage from bosh_metrics where id = '%s' and metricname = '%s' "
	boshDiskIoSql = strings.Replace(boshDiskIoSql, "mean(value)", "non_negative_derivative(mean(value),1m)/1024", 1)
	boshDiskIoSql = strings.Replace(boshDiskIoSql, "metricname = '%s'", "metricname =~ /%s/", 1)

	if boshChart.DefaultTimeRange != "" {
		boshDiskIoSql += " and time > now() - %s  group by time(%s);"
		boshDiskIoSql = fmt.Sprintf(boshDiskIoSql, boshChart.UUID, boshChart.MetricName, boshChart.DefaultTimeRange, boshChart.GroupBy)
	} else {
		boshDiskIoSql += " and time < now() - %s and time > now() - %s  group by time(%s);"
		boshDiskIoSql = fmt.Sprintf(boshDiskIoSql, boshChart.UUID, boshChart.MetricName, boshChart.TimeRangeFrom, boshChart.TimeRangeTo, boshChart.GroupBy)
	}
	fmt.Println("GetBoshNetworkPacketList Sql======>", boshDiskIoSql)

	q := client.Query{
		Command:  boshDiskIoSql,
		Database: p.InfluxDbClient.DbName.BoshDatabase,
	}
	resp, err := p.InfluxDbClient.HttpClient.Query(q)
	if err != nil {
		return resp, err
	}
	fmt.Println("GetBoshNetworkPacketList resp======>", resp)

	return resp, err
}

func (p *PaastaDao) GetPaastaNetworkByteList(boshChart models.BoshChart) (*client.Response, error) {
	boshNetworkByteSql := "select mean(value) as usage from bosh_metrics where id = '%s' and metricname = '%s' "
	boshNetworkByteSql = strings.Replace(boshNetworkByteSql, "mean(value)", "non_negative_derivative(mean(value))/1024", 1)

	if boshChart.DefaultTimeRange != "" {
		boshNetworkByteSql += " and time > now() - %s  group by time(%s);"
		boshNetworkByteSql = fmt.Sprintf(boshNetworkByteSql, boshChart.UUID, boshChart.MetricName, boshChart.DefaultTimeRange, boshChart.GroupBy)
	} else {
		boshNetworkByteSql += " and time < now() - %s and time > now() - %s  group by time(%s);"
		boshNetworkByteSql = fmt.Sprintf(boshNetworkByteSql, boshChart.UUID, boshChart.MetricName, boshChart.TimeRangeFrom, boshChart.TimeRangeTo, boshChart.GroupBy)
	}
	fmt.Println("GetBoshNetworkPacketList Sql======>", boshNetworkByteSql)

	q := client.Query{
		Command:  boshNetworkByteSql,
		Database: p.InfluxDbClient.DbName.BoshDatabase,
	}
	resp, err := p.InfluxDbClient.HttpClient.Query(q)
	if err != nil {
		return resp, err
	}
	fmt.Println("GetBoshNetworkPacketList resp======>", resp)

	return resp, err
}

func (p *PaastaDao) GetPaastaNetworkPacketList(boshChart models.BoshChart) (*client.Response, error) {
	boshNetworkPacketSql := "select mean(value) as usage from bosh_metrics where id = '%s' and metricname = '%s' "
	boshNetworkPacketSql = strings.Replace(boshNetworkPacketSql, "mean(value)", "non_negative_derivative(mean(value))", 1)

	if boshChart.DefaultTimeRange != "" {
		boshNetworkPacketSql += " and time > now() - %s  group by time(%s);"
		boshNetworkPacketSql = fmt.Sprintf(boshNetworkPacketSql, boshChart.UUID, boshChart.MetricName, boshChart.DefaultTimeRange, boshChart.GroupBy)
	} else {
		boshNetworkPacketSql += " and time < now() - %s and time > now() - %s  group by time(%s);"
		boshNetworkPacketSql = fmt.Sprintf(boshNetworkPacketSql, boshChart.UUID, boshChart.MetricName, boshChart.TimeRangeFrom, boshChart.TimeRangeTo, boshChart.GroupBy)
	}
	fmt.Println("GetBoshNetworkPacketList Sql======>", boshNetworkPacketSql)

	q := client.Query{
		Command:  boshNetworkPacketSql,
		Database: p.InfluxDbClient.DbName.BoshDatabase,
	}
	resp, err := p.InfluxDbClient.HttpClient.Query(q)
	if err != nil {
		return resp, err
	}
	fmt.Println("GetBoshNetworkPacketList resp======>", resp)

	return resp, err
}

func (p *PaastaDao) GetPaastaNetworkDropList(boshChart models.BoshChart) (*client.Response, error) {
	boshNetworkDropSql := "select mean(value) as usage from bosh_metrics where id = '%s' and metricname = '%s' "
	boshNetworkDropSql = strings.Replace(boshNetworkDropSql, "mean(value)", "non_negative_derivative(sum(value))", 1)

	if boshChart.DefaultTimeRange != "" {
		boshNetworkDropSql += " and time > now() - %s  group by time(%s);"
		boshNetworkDropSql = fmt.Sprintf(boshNetworkDropSql, boshChart.UUID, boshChart.MetricName, boshChart.DefaultTimeRange, boshChart.GroupBy)
	} else {
		boshNetworkDropSql += " and time < now() - %s and time > now() - %s  group by time(%s);"
		boshNetworkDropSql = fmt.Sprintf(boshNetworkDropSql, boshChart.UUID, boshChart.MetricName, boshChart.TimeRangeFrom, boshChart.TimeRangeTo, boshChart.GroupBy)
	}
	fmt.Println("GetBoshNetworkDropList Sql======>", boshNetworkDropSql)

	q := client.Query{
		Command:  boshNetworkDropSql,
		Database: p.InfluxDbClient.DbName.BoshDatabase,
	}
	resp, err := p.InfluxDbClient.HttpClient.Query(q)
	if err != nil {
		return resp, err
	}
	fmt.Println("GetBoshNetworkDropList resp======>", resp)

	return resp, err
}

func (p *PaastaDao) GetPaastaNetworkErrorList(boshChart models.BoshChart) (*client.Response, error) {
	boshNetworkErrorSql := "select mean(value) as usage from bosh_metrics where id = '%s' and metricname = '%s' "
	boshNetworkErrorSql = strings.Replace(boshNetworkErrorSql, "mean(value)", "non_negative_derivative(sum(value))", 1)

	if boshChart.DefaultTimeRange != "" {
		boshNetworkErrorSql += " and time > now() - %s  group by time(%s);"
		boshNetworkErrorSql = fmt.Sprintf(boshNetworkErrorSql, boshChart.UUID, boshChart.MetricName, boshChart.DefaultTimeRange, boshChart.GroupBy)
	} else {
		boshNetworkErrorSql += " and time < now() - %s and time > now() - %s  group by time(%s);"
		boshNetworkErrorSql = fmt.Sprintf(boshNetworkErrorSql, boshChart.UUID, boshChart.MetricName, boshChart.TimeRangeFrom, boshChart.TimeRangeTo, boshChart.GroupBy)
	}
	fmt.Println("GetBoshNetworkErrorList Sql======>", boshNetworkErrorSql)

	q := client.Query{
		Command:  boshNetworkErrorSql,
		Database: p.InfluxDbClient.DbName.BoshDatabase,
	}
	resp, err := p.InfluxDbClient.HttpClient.Query(q)
	if err != nil {
		return resp, err
	}
	fmt.Println("GetBoshNetworkErrorList resp======>", resp)

	return resp, err
}

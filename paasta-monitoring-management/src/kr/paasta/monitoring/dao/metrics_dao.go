package dao

import (
	"kr/paasta/monitoring/domain"
	"kr/paasta/monitoring/util"
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
)

type serverDtvMetricStruct struct {
	domain.MetricsRequest
	influxClient 	client.Client
	InfraDtvmetricDataSource string
}

func GetMetricsDao(influxClient client.Client, ds string) *serverDtvMetricStruct {

	return &serverDtvMetricStruct{
		influxClient: 	influxClient,
		InfraDtvmetricDataSource: ds,
	}
}


func (b serverDtvMetricStruct) GetDiskIOList(request domain.MetricsRequest, metricname string) (_ client.Response, errMsg domain.ErrMessage) {
	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("GetDiskIOList error :", errLogMsg)
			errMsg = domain.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	var totalDiskIORdSql string
	if request.Origin == "bos" {
		totalDiskIORdSql = "SELECT derivative(mean(value), 30s) as totalUsage FROM bosh_metrics WHERE origin = '%s' and metricname = '%s'";
	} else if request.Origin == "ctl" {
		totalDiskIORdSql = "SELECT derivative(mean(value), 30s) as totalUsage FROM cf_metrics WHERE job = '%s' and metricname = '%s'";
	} else if request.Origin == "ctn" {
		totalDiskIORdSql = "SELECT derivative(mean(value), 30s) as totalUsage FROM cf_metrics WHERE job = '%s' and metricname = '%s'";
	} else if request.Origin == "app" {
		totalDiskIORdSql = "SELECT derivative(mean(value), 30s) as totalUsage FROM container_metrics WHERE container_id = '%s' and \"name\" = '%s'";
	}

	var q client.Query
	if request.DefaultTimeRange != "" {
		totalDiskIORdSql += " and time > now() - %s  group by time(%s);"
		q = client.Query{
			Command: fmt.Sprintf( totalDiskIORdSql, request.ServiceName, metricname, request.DefaultTimeRange, request.GroupBy),
			Database: b.InfraDtvmetricDataSource,
		}
	}else{
		totalDiskIORdSql += " and time < now() - %s and time > now() - %s  group by time(%s);"
		q = client.Query{
			Command: fmt.Sprintf( totalDiskIORdSql, request.ServiceName, metricname, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: b.InfraDtvmetricDataSource,
		}
	}
	//fmt.Printf("DiskIO %s  ", metricname)
	//fmt.Println("Sql==> ", q)
	resp, err := b.influxClient.Query(q)
	if err != nil{
		errLogMsg = err.Error()
	}
	return util.GetError().CheckError(*resp, err)
}


func (b serverDtvMetricStruct) GetNetworkIOUsageList(request domain.MetricsRequest, metricname string) (_ client.Response, errMsg domain.ErrMessage) {
	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("GetNetworkIOUsageList error :", errLogMsg)
			errMsg = domain.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	var totalNetworkIORdSql string
	if request.Origin == "bos" {
		totalNetworkIORdSql = "SELECT derivative(mean(value), 30s) as totalUsage FROM bosh_metrics WHERE origin = '%s' and metricname = '%s'";
	} else if request.Origin == "ctl" {
		totalNetworkIORdSql = "SELECT derivative(mean(value), 30s) as totalUsage FROM cf_metrics WHERE job = '%s' and metricname = '%s'";
	} else if request.Origin == "ctn" {
		totalNetworkIORdSql = "SELECT derivative(mean(value), 30s) as totalUsage FROM cf_metrics WHERE job = '%s' and metricname = '%s'";
	} else if request.Origin == "app" {
		totalNetworkIORdSql = "SELECT derivative(mean(value), 30s) as totalUsage FROM container_metrics WHERE container_id = '%s' and \"name\" = '%s'";
	}

	var q client.Query
	if request.DefaultTimeRange != "" {
		totalNetworkIORdSql += " and time > now() - %s  group by time(%s);"
		q = client.Query{
			Command: fmt.Sprintf( totalNetworkIORdSql, request.ServiceName, metricname, request.DefaultTimeRange, request.GroupBy),
			Database: b.InfraDtvmetricDataSource,
		}
	}else{
		totalNetworkIORdSql += " and time < now() - %s and time > now() - %s  group by time(%s);"
		q = client.Query{
			Command: fmt.Sprintf( totalNetworkIORdSql, request.ServiceName, metricname, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: b.InfraDtvmetricDataSource,
		}
	}
	//fmt.Printf("NetworkIO %s  ", metricname)
	//fmt.Println("Sql==> ", q)
	resp, err := b.influxClient.Query(q)
	if err != nil{
		errLogMsg = err.Error()
	}
	return util.GetError().CheckError(*resp, err)
}

func (b serverDtvMetricStruct) GetTopProcessList(request domain.MetricsRequest) (_ client.Response, errMsg domain.ErrMessage) {
	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("GetTopProcessList error :", errLogMsg)
			errMsg = domain.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	var totalTopProcessRdSql string
	if request.Origin == "bos" {
		totalTopProcessRdSql = "select proc_pid as pid, proc_name as process, (top(mem_usage, 20)/1024/1000) as memory, proc_index  as index from bosh_metrics where origin = '%s' and ip = '%s'"
	} else if request.Origin == "ctl" {
		totalTopProcessRdSql = "select proc_pid as pid, proc_name as process, (top(mem_usage, 20)/1024/1000) as memory, proc_index  as index from cf_metrics where job = '%s' and ip = '%s'"
	} else if request.Origin == "ctn" {
		totalTopProcessRdSql = "select proc_pid as pid, proc_name as process, (top(mem_usage, 20)/1024/1000) as memory, proc_index  as index from container_metrics where container_id = '%s' and ip = '%s'"
	}

	var q client.Query
	if request.DefaultTimeRange != "" {
		totalTopProcessRdSql += " and time > now() - %s  group by time(%s);"
		q = client.Query{
			Command: fmt.Sprintf( totalTopProcessRdSql, request.ServiceName, request.Addr, request.DefaultTimeRange, request.GroupBy),
			Database: b.InfraDtvmetricDataSource,
		}
	}else{
		totalTopProcessRdSql += " and time < now() - %s and time > now() - %s  group by time(%s);"
		q = client.Query{
			Command: fmt.Sprintf( totalTopProcessRdSql, request.ServiceName, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: b.InfraDtvmetricDataSource,
		}
	}
	//fmt.Println("TopProcess Sql==> ", q)
	resp, err := b.influxClient.Query(q)
	if err != nil{
		errLogMsg = err.Error()
	}
	return util.GetError().CheckError(*resp, err)
}



//Container(Application) Resource info - Cpu, Memory, Disk
func (b serverDtvMetricStruct) GetApplicationResources(request domain.MetricsRequest) (_ client.Response, errMsg domain.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("GetApplicationResources error :", errLogMsg)
			errMsg = domain.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	appResourceSql := "SELECT \"name\", value  FROM \"container_metrics\" WHERE application_id = '%s' AND application_index = '%s'  AND (\"name\" = 'cpu_usage_total' OR \"name\" = 'memory_usage' OR \"name\" = 'disk_usage') and time > now() - %s;"
	q := client.Query{
		Command: fmt.Sprintf( appResourceSql, request.ServiceName, request.Index, request.DefaultTimeRange),
		Database: b.InfraDtvmetricDataSource,
	}
	resp, err := b.influxClient.Query(q)
	if err != nil{
		errLogMsg = err.Error()
	}
	return util.GetError().CheckError(*resp, err)
}

//Cpu variation of Container(Application)
func (b serverDtvMetricStruct) GetAppCpuUsage(request domain.MetricsRequest) (_ client.Response, errMsg domain.ErrMessage) {
	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("GetCpuVariation error :", errLogMsg)
			errMsg = domain.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	cpuVariationSql := "select non_negative_derivative(mean(value),30s)/30000000000 * 100000000000 as value from container_metrics where \"name\" = 'cpu_usage_total' and application_id = '%s' and application_index = '%s' "
	var q client.Query
	if request.DefaultTimeRange != "" {
		cpuVariationSql += " and time > now() - %s group by time(%s);"
		q = client.Query{
			Command:  fmt.Sprintf( cpuVariationSql, request.ServiceName, request.Index, request.DefaultTimeRange, request.GroupBy),
			Database: b.InfraDtvmetricDataSource,
		}
	}else{
		cpuVariationSql += " and time < now() - %s and time > now() - %s group by time(%s);"
		q = client.Query{
			Command:  fmt.Sprintf( cpuVariationSql, request.ServiceName,  request.Index, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: b.InfraDtvmetricDataSource,

		}
	}
	fmt.Println("GetAppCpuVariation", q)
	resp, err := b.influxClient.Query(q)
	if err != nil{
		errLogMsg = err.Error()
	}
	return util.GetError().CheckError(*resp, err)
}

//Memory variation of Container(Application)
func (b serverDtvMetricStruct) GetAppMemoryUsage(request domain.MetricsRequest) (_ client.Response, errMsg domain.ErrMessage) {
	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("GetMemoryVariation error :", errLogMsg)
			errMsg = domain.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	//Unit : byte
	memVariationSql := "select mean(value)/mean(app_mem) * 100 as value from container_metrics where \"name\" = 'memory_usage' and application_id = '%s' and application_index = '%s' ";
	var q client.Query
	if request.DefaultTimeRange != "" {
		memVariationSql += " and time > now() - %s group by time(%s);"
		q = client.Query{
			Command:  fmt.Sprintf( memVariationSql, request.ServiceName, request.Index, request.DefaultTimeRange, request.GroupBy),
			Database: b.InfraDtvmetricDataSource,
		}
	}else{
		memVariationSql += " and time < now() - %s and time > now() - %s group by time(%s);"
		q = client.Query{
			Command:  fmt.Sprintf(memVariationSql, request.ServiceName, request.Index, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: b.InfraDtvmetricDataSource,
		}
	}
	fmt.Println("GetAppMemoryVariation:",q)
	resp, err := b.influxClient.Query(q)
	if err != nil{
		errLogMsg = err.Error()
	}

	return util.GetError().CheckError(*resp, err)
}

func (b serverDtvMetricStruct) GetAppDiskUsage(request domain.MetricsRequest) (_ client.Response, errMsg domain.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("GetMemoryVariation error :", errLogMsg)
			errMsg = domain.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	//Unit : byte
	memVariationSql := "select mean(value)/mean(app_disk) * 100 as value from container_metrics where \"name\" = 'disk_usage' and application_id = '%s' and application_index = '%s' ";
	var q client.Query
	if request.DefaultTimeRange != "" {
		memVariationSql += " and time > now() - %s group by time(%s);"
		q = client.Query{
			Command:  fmt.Sprintf( memVariationSql, request.ServiceName, request.Index, request.DefaultTimeRange, request.GroupBy),
			Database: b.InfraDtvmetricDataSource,
		}
	}else{
		memVariationSql += " and time < now() - %s and time > now() - %s group by time(%s);"
		q = client.Query{
			Command:  fmt.Sprintf(memVariationSql, request.ServiceName, request.Index, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: b.InfraDtvmetricDataSource,
		}
	}

	resp, err := b.influxClient.Query(q)
	if err != nil{
		errLogMsg = err.Error()
	}

	return util.GetError().CheckError(*resp, err)
}


//Network Rx of Container(Application)
func (b serverDtvMetricStruct) GetAppNetworkKByte(request domain.MetricsRequest, name string) (_ client.Response, errMsg domain.ErrMessage) {
	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("GetAppNetworkRxVariation error :", errLogMsg)
			errMsg = domain.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()
	var q client.Query

	//Get Container_Interface
	containerInterface := "select * from container_metrics where \"name\" = 'cpu_usage_total' and application_id = '%s' and application_index = '%s' and time > now() - 1m"
	q = client.Query{
		Command:  fmt.Sprintf( containerInterface, request.ServiceName, request.Index),
		Database: b.InfraDtvmetricDataSource,
	}

	resp, err := b.influxClient.Query(q)
	result, _ := util.GetResponseConverter().InfluxConverter(*resp, request.ServiceName)

	if len(result["data"].([]map[string]interface{})) > 0 {
		container_interface := result["data"].([]map[string]interface{})[0]["container_interface"]
		networkSql := "select non_negative_derivative(sum(value),30s)/1024 as value from container_metrics where \"name\" = '%s' and container_id = '%s'";
		if request.DefaultTimeRange != "" {
			networkSql += " and time > now() - %s  group by time(%s);"
			q = client.Query{
				Command:  fmt.Sprintf( networkSql,  name, container_interface,  request.DefaultTimeRange, request.GroupBy),
				Database: b.InfraDtvmetricDataSource,
			}
		}else{
			networkSql += " and time < now() - %s and time > now() - %s  group by time(%s);"
			q = client.Query{
				Command:  fmt.Sprintf( networkSql, name, container_interface, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
				Database: b.InfraDtvmetricDataSource,
			}
		}
		resp, err = b.influxClient.Query(q)
		fmt.Println("GetAppNetworkKByte::",q)
		if err != nil{
			errLogMsg = err.Error()
		}
	}else {
		errLogMsg = "There is no result for your request. Please try again or check the request data."
	}
	return util.GetError().CheckError(*resp, err)
}
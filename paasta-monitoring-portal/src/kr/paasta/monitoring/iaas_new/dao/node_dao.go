package dao

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"monitoring-portal/iaas_new/model"
	"monitoring-portal/utils"
)

type NodeDao struct {
	influxClient client.Client
}

func GetNodeDao(influxClient client.Client) *NodeDao {
	return &NodeDao{
		influxClient: influxClient,
	}
}

// Node의 현재 CPU 사용률을 조회
func (d NodeDao) GetNodeCpuUsage(request model.NodeReq) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string

	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	cpuUsageSql := "select value from \"cpu.percent\"  where time > now() - 2m and hostname = '%s' order by time desc limit 1"

	var q client.Query

	q = client.Query{
		Command: fmt.Sprintf(cpuUsageSql,
			request.HostName),
		Database: model.MetricDBName,
	}

	model.MonitLogger.Debug("GetNodeCpuUsage Sql======>", q)
	resp, err := d.influxClient.Query(q)
	if err != nil {
		errLogMsg = err.Error()
	}

	return utils.GetError().CheckError(*resp, err)
}

// Node의 현재 CPU 사용률을 목록 조회
func (d NodeDao) GetNodeCpuUsageList(request model.DetailReq) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string

	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	cpuUsageSql := "select mean(value) as usage from \"cpu.percent\"  where hostname = '%s' "

	model.MonitLogger.Debugf("defaultTimeRange: %s, timeRangeFrom: %s, timeRangeTo:%s", request.DefaultTimeRange, request.TimeRangeFrom, request.TimeRangeTo)

	var q client.Query
	if request.DefaultTimeRange != "" {

		cpuUsageSql += " and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command: fmt.Sprintf(cpuUsageSql,
				request.HostName, request.DefaultTimeRange, request.GroupBy),
			Database: model.MetricDBName,
		}
	} else {

		cpuUsageSql += " and time < now() - %s and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command: fmt.Sprintf(cpuUsageSql,
				request.HostName, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: model.MetricDBName,
		}
	}

	model.MonitLogger.Debug("GetNodeCpuUsageList Sql==>", q)
	resp, err := d.influxClient.Query(q)

	return utils.GetError().CheckError(*resp, err)
}

// Node의 현재 Memory 사용률을 조회
func (d NodeDao) GetNodeMemoryUsageList(request model.DetailReq) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string

	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	memoryTotalSql := "select mean(value) as usage from \"mem.usable_perc\"  where hostname = '%s' "
	//memoryFreeSql := "select mean(value) as usage from \"mem.free_mb\"  where hostname = '%s' ";

	model.MonitLogger.Debugf("defaultTimeRange: %s, timeRangeFrom: %s, timeRangeTo:%s", request.DefaultTimeRange, request.TimeRangeFrom, request.TimeRangeTo)

	var q client.Query
	if request.DefaultTimeRange != "" {

		memoryTotalSql += " and time > now() - %s  group by time(%s);"
		//memoryFreeSql += " and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command: fmt.Sprintf(memoryTotalSql,
				request.HostName, request.DefaultTimeRange, request.GroupBy),
			Database: model.MetricDBName,
		}
	} else {

		memoryTotalSql += " and time < now() - %s and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command: fmt.Sprintf(memoryTotalSql,
				request.HostName, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: model.MetricDBName,
		}
	}

	model.MonitLogger.Debug("GetNodeMemoryUsageList Sql==>", q)
	resp, err := d.influxClient.Query(q)

	return utils.GetError().CheckError(*resp, err)
}

//Node의 현재 CPU사용률을 조회한다.
func (d NodeDao) GetNodeSwapMemoryFreeUsageList(request model.DetailReq) (_ client.Response, errMsg model.ErrMessage) {
	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	swapFreeUsageSql := "select mean(value) as usage from \"mem.swap_free_perc\"  where hostname = '%s' "

	model.MonitLogger.Debugf("defaultTimeRange: %s, timeRangeFrom: %s, timeRangeTo:%s", request.DefaultTimeRange, request.TimeRangeFrom, request.TimeRangeTo)

	var q client.Query
	if request.DefaultTimeRange != "" {

		swapFreeUsageSql += " and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command: fmt.Sprintf(swapFreeUsageSql,
				request.HostName, request.DefaultTimeRange, request.GroupBy),
			Database: model.MetricDBName,
		}
	} else {

		swapFreeUsageSql += " and time < now() - %s and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command: fmt.Sprintf(swapFreeUsageSql,
				request.HostName, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: model.MetricDBName,
		}
	}

	model.MonitLogger.Debug("GetNodeSwapMemoryFreeUsageList Sql==>", q)
	resp, err := d.influxClient.Query(q)

	return utils.GetError().CheckError(*resp, err)
}

//Node의 현재 CPU사용률을 조회한다.
func (d NodeDao) GetNodeCpuLoadList(request model.DetailReq, minute string) (_ client.Response, errMsg model.ErrMessage) {
	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	var cpuUsageSql string
	if minute == "1m" {
		cpuUsageSql = "select mean(value) as usage from \"load.avg_1_min\"  where hostname = '%s' "
	} else if minute == "5m" {
		cpuUsageSql = "select mean(value) as usage from \"load.avg_5_min\"  where hostname = '%s' "
	} else if minute == "15m" {
		cpuUsageSql = "select mean(value) as usage from \"load.avg_15_min\"  where hostname = '%s' "
	}

	model.MonitLogger.Debugf("defaultTimeRange: %s, timeRangeFrom: %s, timeRangeTo:%s", request.DefaultTimeRange, request.TimeRangeFrom, request.TimeRangeTo)

	var q client.Query
	if request.DefaultTimeRange != "" {

		cpuUsageSql += " and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command: fmt.Sprintf(cpuUsageSql,
				request.HostName, request.DefaultTimeRange, request.GroupBy),
			Database: model.MetricDBName,
		}
	} else {

		cpuUsageSql += " and time < now() - %s and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command: fmt.Sprintf(cpuUsageSql,
				request.HostName, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: model.MetricDBName,
		}
	}

	model.MonitLogger.Debug("GetNodeCpuLoadList Sql==>", q)
	resp, err := d.influxClient.Query(q)

	return utils.GetError().CheckError(*resp, err)
}

//Node의 현재 Total Memory을 조회한다.
func (d NodeDao) GetNodeMemoryUsage(request model.NodeReq) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	totalMemSql := "select value from \"mem.usable_perc\" where time > now() - 2m and hostname = '%s' order by time desc limit 1;"

	var q client.Query

	q = client.Query{
		Command: fmt.Sprintf(totalMemSql,
			request.HostName),
		Database: model.MetricDBName,
	}

	model.MonitLogger.Debug("GetNodeMemoryUsage Sql======>", q)
	resp, err := d.influxClient.Query(q)
	if err != nil {
		errLogMsg = err.Error()
	}

	return utils.GetError().CheckError(*resp, err)
}

//Node의 현재 Total Memory을 조회한다.
func (d NodeDao) GetNodeTotalMemoryUsage(request model.NodeReq) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	totalMemSql := "select value from \"mem.total_mb\" where time > now() - 2m and hostname = '%s' order by time desc limit 1;"

	var q client.Query

	q = client.Query{
		Command: fmt.Sprintf(totalMemSql,
			request.HostName),
		Database: model.MetricDBName,
	}

	model.MonitLogger.Debug("GetNodeCpuUsage Sql======>", q)
	resp, err := d.influxClient.Query(q)
	if err != nil {
		errLogMsg = err.Error()
	}

	return utils.GetError().CheckError(*resp, err)
}

//Node의 현재 Total Memory을 조회한다.
func (d NodeDao) GetNodeFreeMemoryUsage(request model.NodeReq) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	freeMemSql := "select value from \"mem.free_mb\"  where time > now() - 2m and hostname = '%s' order by time desc limit 1;"

	var q client.Query

	q = client.Query{
		Command: fmt.Sprintf(freeMemSql,
			request.HostName),
		Database: model.MetricDBName,
	}

	model.MonitLogger.Debug("GetNodeCpuUsage Sql======>", q)
	resp, err := d.influxClient.Query(q)
	if err != nil {
		errLogMsg = err.Error()
	}

	return utils.GetError().CheckError(*resp, err)
}

//Node의 현재 Total Memory을 조회한다.
func (d NodeDao) GetNodeTotalDisk(request model.NodeReq) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	totalMemSql := "select value from \"disk.total_space_mb\" where time > now() - 2m and hostname = '%s' order by time desc limit 1;"

	var q client.Query

	q = client.Query{
		Command: fmt.Sprintf(totalMemSql,
			request.HostName),
		Database: model.MetricDBName,
	}

	model.MonitLogger.Debug("GetNodeTotalDisk Sql======>", q)
	resp, err := d.influxClient.Query(q)
	if err != nil {
		errLogMsg = err.Error()
	}

	return utils.GetError().CheckError(*resp, err)
}

//Node의 현재 Total Memory을 조회한다.
func (d NodeDao) GetNodeUsedDisk(request model.NodeReq) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	totalMemSql := "select value from \"disk.total_used_space_mb\" where time > now() - 2m and hostname = '%s' order by time desc limit 1;"

	var q client.Query

	q = client.Query{
		Command: fmt.Sprintf(totalMemSql,
			request.HostName),
		Database: model.MetricDBName,
	}

	model.MonitLogger.Debug("GetNodeUsedDisk Sql======>", q)
	resp, err := d.influxClient.Query(q)
	if err != nil {
		errLogMsg = err.Error()
	}

	return utils.GetError().CheckError(*resp, err)
}

//Monasca Agent Forwarder 현재 상태 조회
func (d NodeDao) GetAgentProcessStatus(request model.NodeReq, processName string) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	agentStatusSql := "select value from \"supervisord.process.status\" where hostname = '%s' and supervisord_process = '%s' and time > now() - 2m order by time desc limit 1;"

	var q client.Query

	q = client.Query{
		Command: fmt.Sprintf(agentStatusSql,
			request.HostName, processName),
		Database: model.MetricDBName,
	}

	model.MonitLogger.Debug("AgentProcess Status Sql======>", q)
	resp, err := d.influxClient.Query(q)
	if err != nil {
		errLogMsg = err.Error()
	}

	return utils.GetError().CheckError(*resp, err)
}

//VM Instance가 Running인 VM만 조회
func (d NodeDao) GetAliveInstanceListByNodename(request model.NodeReq, allStatus bool) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	/*     VM Instance Status
	       -1 : no status,
		0 : Running / OK,
		1 : Idle / blocked,
		2 : Paused,
		3 : Shutting down,
		4 : Shut off or Nova suspend
		5 : Crashed,
		6 : Power management suspend (S3 state)
	*/
	var instanceListStatusSql string
	if allStatus == true {
		instanceListStatusSql = "select resource_id, value from \"vm.host_alive_status\" where time > now() - 2m and hostname = '%s' ;"
	} else {
		instanceListStatusSql = "select resource_id, value from \"vm.host_alive_status\" where time > now() - 2m and hostname = '%s' and value = 0 ;"
	}

	var q client.Query

	q = client.Query{
		Command: fmt.Sprintf(instanceListStatusSql,
			request.HostName),
		Database: model.MetricDBName,
	}

	model.MonitLogger.Debug("GetInstanceList Sql======>", q)
	resp, err := d.influxClient.Query(q)
	if err != nil {
		errLogMsg = err.Error()
	}

	return utils.GetError().CheckError(*resp, err)
}

func (b NodeDao) GetMountPointList(request model.DetailReq) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	mountPointListSql := "select  mount_point, value from \"disk.space_used_perc\"  where time > now() - 150s and hostname = '%s'"
	var q client.Query

	q = client.Query{
		Command: fmt.Sprintf(mountPointListSql,
			request.HostName),
		Database: model.MetricDBName,
	}

	model.MonitLogger.Debug("GetServiceFileSystems Sql======>", q)
	resp, err := b.influxClient.Query(q)

	if err != nil {
		errLogMsg = err.Error()
	}

	return utils.GetError().CheckError(*resp, err)
}

//Node의 현재 Total Memory을 조회한다.
func (d NodeDao) GetNodeDiskUsage(request model.DetailReq) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	diskUsageSql := "select mean(value) as usage from \"disk.space_used_perc\"  where hostname = '%s' and mount_point = '%s' "
	model.MonitLogger.Debugf("defaultTimeRange: %s, timeRangeFrom: %s, timeRangeTo:%s", request.DefaultTimeRange, request.TimeRangeFrom, request.TimeRangeTo)

	var q client.Query
	if request.DefaultTimeRange != "" {

		diskUsageSql += " and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command: fmt.Sprintf(diskUsageSql,
				request.HostName, request.MountPoint, request.DefaultTimeRange, request.GroupBy),
			Database: model.MetricDBName,
		}
	} else {

		diskUsageSql += " and time < now() - %s and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command: fmt.Sprintf(diskUsageSql,
				request.HostName, request.MountPoint, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: model.MetricDBName,
		}
	}
	model.MonitLogger.Debug("GetNodeDiskUsage Sql==>", q)
	resp, err := d.influxClient.Query(q)

	return utils.GetError().CheckError(*resp, err)

}

//Node의 disk io read Kbyte를 조회한다.
func (d NodeDao) GetNodeDiskIoReadKbyte(request model.DetailReq) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	diskUsageSql := "select mean(value) as usage from \"io.read_kbytes_sec\"  where hostname = '%s' and mount_point = '%s' "
	model.MonitLogger.Debugf("defaultTimeRange: %s, timeRangeFrom: %s, timeRangeTo:%s", request.DefaultTimeRange, request.TimeRangeFrom, request.TimeRangeTo)

	var q client.Query
	if request.DefaultTimeRange != "" {

		diskUsageSql += " and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command: fmt.Sprintf(diskUsageSql,
				request.HostName, request.MountPoint, request.DefaultTimeRange, request.GroupBy),
			Database: model.MetricDBName,
		}
	} else {

		diskUsageSql += " and time < now() - %s and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command: fmt.Sprintf(diskUsageSql,
				request.HostName, request.MountPoint, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: model.MetricDBName,
		}
	}
	model.MonitLogger.Debug("GetNodeDiskIoReadKbyte Sql==>", q)
	resp, err := d.influxClient.Query(q)

	return utils.GetError().CheckError(*resp, err)

}

//Node의 disk io read Kbyte를 조회한다.
func (d NodeDao) GetNodeDiskIoWriteKbyte(request model.DetailReq) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	diskUsageSql := "select mean(value) as usage from \"io.write_kbytes_sec\"  where hostname = '%s' and mount_point = '%s' "
	model.MonitLogger.Debugf("defaultTimeRange: %s, timeRangeFrom: %s, timeRangeTo:%s", request.DefaultTimeRange, request.TimeRangeFrom, request.TimeRangeTo)

	var q client.Query
	if request.DefaultTimeRange != "" {

		diskUsageSql += " and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command: fmt.Sprintf(diskUsageSql,
				request.HostName, request.MountPoint, request.DefaultTimeRange, request.GroupBy),
			Database: model.MetricDBName,
		}
	} else {

		diskUsageSql += " and time < now() - %s and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command: fmt.Sprintf(diskUsageSql,
				request.HostName, request.MountPoint, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: model.MetricDBName,
		}
	}
	model.MonitLogger.Debug("GetNodeDiskIoReadKbyte Sql==>", q)
	resp, err := d.influxClient.Query(q)

	return utils.GetError().CheckError(*resp, err)

}

//Node의 disk io read Kbyte를 조회한다.
func (d NodeDao) GetNodeNetworkKbyte(request model.DetailReq, inOut, device string) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	var diskUsageSql string
	if inOut == "in" {
		diskUsageSql = "select sum(value)/1024 as usage from \"net.in_bytes_sec\"  where hostname = '%s' and device =~ /%s/"
	} else if inOut == "out" {
		diskUsageSql = "select sum(value)/1024 as usage from \"net.out_bytes_sec\"  where hostname = '%s' and device =~ /%s/"
	}

	model.MonitLogger.Debugf("defaultTimeRange: %s, timeRangeFrom: %s, timeRangeTo:%s", request.DefaultTimeRange, request.TimeRangeFrom, request.TimeRangeTo)

	var q client.Query
	if request.DefaultTimeRange != "" {

		diskUsageSql += " and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command: fmt.Sprintf(diskUsageSql,
				request.HostName, device, request.DefaultTimeRange, request.GroupBy),
			Database: model.MetricDBName,
		}
	} else {

		diskUsageSql += " and time < now() - %s and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command: fmt.Sprintf(diskUsageSql,
				request.HostName, device, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: model.MetricDBName,
		}
	}
	model.MonitLogger.Debug("GetNodeNetworkInOutKbyte Sql==>", q)
	resp, err := d.influxClient.Query(q)

	return utils.GetError().CheckError(*resp, err)

}

//Node의 Network Error를 조회한다.
func (d NodeDao) GetNodeNetworkError(request model.DetailReq, inOut, device string) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	var diskUsageSql string
	if inOut == "in" {
		diskUsageSql = "select sum(value) as usage from \"net.in_errors_sec\"  where hostname = '%s' and device =~ /%s/"
	} else if inOut == "out" {
		diskUsageSql = "select sum(value) as usage from \"net.out_errors_sec\"  where hostname = '%s' and device =~ /%s/"
	}

	model.MonitLogger.Debugf("defaultTimeRange: %s, timeRangeFrom: %s, timeRangeTo:%s", request.DefaultTimeRange, request.TimeRangeFrom, request.TimeRangeTo)

	var q client.Query
	if request.DefaultTimeRange != "" {

		diskUsageSql += " and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command: fmt.Sprintf(diskUsageSql,
				request.HostName, device,
				request.DefaultTimeRange, request.GroupBy),
			Database: model.MetricDBName,
		}
	} else {

		diskUsageSql += " and time < now() - %s and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command: fmt.Sprintf(diskUsageSql,
				request.HostName, device,
				request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: model.MetricDBName,
		}
	}
	model.MonitLogger.Debug("GetNodeNetworkError Sql==>", q)
	resp, err := d.influxClient.Query(q)

	return utils.GetError().CheckError(*resp, err)

}

//Node의 Network Dropped packets를 조회한다.
func (d NodeDao) GetNodeNetworkDropPacket(request model.DetailReq, inOut, device string) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	var diskUsageSql string
	if inOut == "in" {
		diskUsageSql = "select sum(value) as usage from \"net.in_packets_dropped_sec\"  where hostname = '%s'and device =~ /%s/ "
	} else if inOut == "out" {
		diskUsageSql = "select sum(value) as usage from \"net.out_packets_dropped_sec\"  where hostname = '%s' and device =~ /%s/"
	}

	model.MonitLogger.Debugf("defaultTimeRange: %s, timeRangeFrom: %s, timeRangeTo:%s", request.DefaultTimeRange, request.TimeRangeFrom, request.TimeRangeTo)

	var q client.Query
	if request.DefaultTimeRange != "" {
		diskUsageSql += " and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command: fmt.Sprintf(diskUsageSql, request.HostName, device,
				request.DefaultTimeRange, request.GroupBy),
			Database: model.MetricDBName,
		}
	} else {
		diskUsageSql += " and time < now() - %s and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command: fmt.Sprintf(diskUsageSql, request.HostName, device,
				request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: model.MetricDBName,
		}
	}
	model.MonitLogger.Debug("GetNodeNetworkDropPacket Sql==>", q)
	resp, err := d.influxClient.Query(q)

	return utils.GetError().CheckError(*resp, err)
}

//Node의 disk io read Kbyte를 조회한다.
func (d NodeDao) GetNodeTopProcessByCpu(request model.DetailReq) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	cpuTopProcessSql := "select mean(value) as usage from \"process.cpu_perc\"  where time > now() - 2m and hostname = '%s' group by process_name "

	var q client.Query

	q = client.Query{
		Command: fmt.Sprintf(cpuTopProcessSql,
			request.HostName),
		Database: model.MetricDBName,
	}

	model.MonitLogger.Debug("GetNodeTopProcessByCpu Sql==>", q)

	resp, err := d.influxClient.Query(q)
	if err != nil {
		errLogMsg = err.Error()
	}

	return utils.GetError().CheckError(*resp, err)

}

//Node의 disk io read Kbyte를 조회한다.
func (d NodeDao) GetNodeTopProcessByMemory(request model.DetailReq) (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	cpuTopProcessSql := "select mean(value) as usage from \"process.mem.rss_mbytes\"  where time > now() - 2m and hostname = '%s' group by process_name "

	var q client.Query

	q = client.Query{
		Command: fmt.Sprintf(cpuTopProcessSql,
			request.HostName),
		Database: model.MetricDBName,
	}

	model.MonitLogger.Debug("GetNodeTopProcessByMemory Sql==>", q)

	resp, err := d.influxClient.Query(q)
	if err != nil {
		errLogMsg = err.Error()
	}

	return utils.GetError().CheckError(*resp, err)
}

package dao

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"monitoring-portal/iaas_new/model"
	"monitoring-portal/utils"
)

type MainDao struct {
	influxClient client.Client
}

func GetMainDao(influxClient client.Client) *MainDao {
	return &MainDao{
		influxClient: influxClient,
	}
}

//Node의 현재 CPU사용률을 조회한다.
func (d MainDao) GetNodeCpuUsage(request model.NodeReq) (_ client.Response, errMsg model.ErrMessage) {

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

//Node의 현재 Total Memory을 조회한다.
func (d MainDao) GetNodeTotalMemoryUsage(request model.NodeReq) (_ client.Response, errMsg model.ErrMessage) {

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
func (d MainDao) GetNodeFreeMemoryUsage(request model.NodeReq) (_ client.Response, errMsg model.ErrMessage) {

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
func (d MainDao) GetNodeTotalDisk(request model.NodeReq) (_ client.Response, errMsg model.ErrMessage) {

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
func (d MainDao) GetNodeUsedDisk(request model.NodeReq) (_ client.Response, errMsg model.ErrMessage) {

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
func (d MainDao) GetAgentProcessStatus(request model.NodeReq, processName string) (_ client.Response, errMsg model.ErrMessage) {

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
func (d MainDao) GetAliveInstanceListByNodename(request model.NodeReq) (_ client.Response, errMsg model.ErrMessage) {

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

	instanceListStatusSql := "select resource_id, value from \"vm.host_alive_status\" where time > now() - 2m and hostname = '%s' and value = 0 ;"

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

//Monasca Agent Forwarder 현재 상태 조회
func (d MainDao) GetOpenstackNodeList() (_ client.Response, errMsg model.ErrMessage) {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	nodeListSql := "select hostname, value from \"cpu.percent\" where time > now() - 2m;"

	var q client.Query

	q = client.Query{
		Command:  fmt.Sprintf(nodeListSql),
		Database: model.MetricDBName,
	}

	model.MonitLogger.Debug("NodeList Sql =====>", q)
	resp, err := d.influxClient.Query(q)
	if err != nil {
		errLogMsg = err.Error()
	}

	return utils.GetError().CheckError(*resp, err)
}

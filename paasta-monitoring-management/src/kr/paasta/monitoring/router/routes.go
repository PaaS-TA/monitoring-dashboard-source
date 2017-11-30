package router

import "github.com/tedsuo/rata"

const (
	Main = "Main"

	// Alarm Policy
	AlarmPolicyList = "AlarmPolicyList"
	AlarmPolicyUpdate = "AlarmPolicyUpdate"

	// Alarm
	AlarmList = "AlarmList"
	AlarmResolveStatus = "AlarmResolveStatus"
	AlarmDetail = "AlarmDetail"
	UpdateAlarm = "UpdateAlarm"
	CreateAlarmAction = "CreateAlarmAction"
	UpdateAlarmAction = "UpdateAlarmAction"
	DeleteAlarmAction = "DeleteAlarmAction"

	GetAlarmStat = "GetAlarmStat"
	GetContainerDeploy = "GetContainerDeploy"

	GetDiskIOList = "GetDiskIOList"
	GetNetworkIOList = "GetNetworkIOList"
	GetTopProcessList = "GetTopProcessList"

	// Web Resource
	Static           = "Static"

	// Application Resources(CPU, Memory, Disk usage)
	GetAppResource = "GetAppResource"
	GetAppResourceAll = "GetAppResourceAll"
	GetAppCpuUsage = "GetAppCpuVariation"
	GetAppMemoryUsage = "GetAppMemoryVariation"
	GetDiskUsage = "GetDiskUsage"
	GetAppNetworkIoByte = "GetAppNetworkVariation"
)

var Routes = rata.Routes{
	{Path: "/", Method: "GET", Name: Main},

	// AlarmPolicy
	{Path: "/alarmsPolicy", Method: "GET", Name: AlarmPolicyList},
	{Path: "/alarmsPolicy", Method: "PUT", Name: AlarmPolicyUpdate},

	// Alarm
	{Path: "/alarms", 		       Method: "GET", Name: AlarmList},
	{Path: "/alarms/status/:resolveStatus", Method: "GET", Name: AlarmResolveStatus},
	{Path: "/alarms/:id", 		       Method: "GET", Name: AlarmDetail},
	{Path: "/alarms/:id", 		       Method: "PUT", Name: UpdateAlarm},
	{Path: "/alarmsAction",   	       Method: "POST", Name: CreateAlarmAction},
	{Path: "/alarmsAction/:actionId",      Method: "PUT", Name: UpdateAlarmAction},
	{Path: "/alarmsAction/:actionId",      Method: "DELETE", Name: DeleteAlarmAction},

	{Path: "/alarmsStat",                  Method: "GET", Name: GetAlarmStat},
	{Path: "/containerDeploy",             Method: "GET", Name: GetContainerDeploy},

	{Path: "/diskIO/:origin", Method: "GET", Name: GetDiskIOList},
	{Path: "/networkIO/:origin", Method: "GET", Name: GetNetworkIOList},
	{Path: "/topProcess/:origin", Method: "GET", Name: GetTopProcessList},

	//Portal API for application monitoring
	{Path: "/app/resources",                Method: "GET", Name: GetAppResource},
	{Path: "/app/resources/all",            Method: "GET", Name: GetAppResourceAll},
	{Path: "/app/:guid/:idx/cpuUsage",      Method: "GET", Name: GetAppCpuUsage},
	{Path: "/app/:guid/:idx/memoryUsage",   Method: "GET", Name: GetAppMemoryUsage},
	{Path: "/app/:guid/:idx/diskUsage",     Method: "GET", Name: GetDiskUsage},
	{Path: "/app/:guid/:idx/networkIoKByte", Method: "GET", Name: GetAppNetworkIoByte},

	// Web Resource
	{Path: "/public/", Method: "GET", Name: Static},
}

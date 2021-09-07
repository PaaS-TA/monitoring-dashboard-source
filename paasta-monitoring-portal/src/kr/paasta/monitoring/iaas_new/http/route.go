package http

import "github.com/tedsuo/rata"


const (
	MEMBER_JOIN_CHECK_DUPLICATION_IAAS_ID = "MEMBER_JOIN_CHECK_DUPLICATION_IAAS_ID"
	MEMBER_JOIN_CHECK_IAAS                = "MEMBER_JOIN_CHECK_IAAS"

	//DASHBOARD = "DASHBOARD"

	IAAS_MAIN_SUMMARY                   = "IAAS_MAIN_SUMMARY"
	IAAS_NODE_MANAGE_SUMMARY            = "IAAS_NODE_MANAGE_SUMMARY"
	IAAS_NODE_COMPUTE_SUMMARY           = "IAAS_NODE_COMPUTE_SUMMARY"
	IAAS_NODE_RABBITMQ_SUMMARY_OVERVIEW = "IAAS_NODE_RABBITMQ_SUMMARY_OVERVIEW"
	//IAAS_NODE_RABBITMQ_SUMMARY            = "IAAS_NODE_RABBITMQ_SUMMARY"

	IAAS_NODES = "IAAS_NODES"

	IAAS_NODE_TOPPROCESS_CPU    = "IAAS_NODE_TOPPROCESS_CPU"
	IAAS_NODE_TOPPROCESS_MEMORY = "IAAS_NODE_TOPPROCESS_MEMORY"

	IAAS_NODE_CPU_USAGE_LIST           = "IAAS_NODE_CPU_USAGE_LIST"
	IAAS_NODE_CPU_LOAD_LIST            = "IAAS_NODE_CPU_LOAD_LIST"
	IAAS_NODE_MEMORY_SWAP_LIST         = "IAAS_NODE_MEMORY_SWAP_LIST"
	IAAS_NODE_MEMORY_USAGE_LIST        = "IAAS_NODE_MEMORY_USAGE_LIST"
	IAAS_NODE_DISK_USAGE_LIST          = "IAAS_NODE_DISK_USAGE_LIST"
	IAAS_NODE_DISK_READ_LIST           = "IAAS_NODE_DISK_READ_LIST"
	IAAS_NODE_DISK_WRITE_LIST          = "IAAS_NODE_DISK_WRITE_LIST"
	IAAS_NODE_NETWORK_KBYTE_LIST       = "IAAS_NODE_NETWORK_KBYTE_LIST"
	IAAS_NODE_NETWORK_ERROR_LIST       = "IAAS_NODE_NETWORK_ERROR_LIST"
	IAAS_NODE_NETWORK_DROP_PACKET_LIST = "IAAS_NODE_NETWORK_DROP_PACKET_LIST"

	IAAS_TENANT_SUMMARY             = "IAAS_TENANT_SUMMARY"
	IAAS_TENANT_INSTANCE_LIST       = "IAAS_TENANT_INSTANCE_LIST"
	IAAS_TENANT_CPU_USAGE_LIST      = "IAAS_TENANT_CPU_USAGE_LIST"
	IAAS_TENANT_MEMORY_USAGE_LIST   = "IAAS_TENANT_MEMORY_USAGE_LIST"
	IAAS_TENANT_DISK_READ_LIST      = "IAAS_TENANT_DISK_READ_LIST"
	IAAS_TENANT_DISK_WRITE_LIST     = "IAAS_TENANT_DISK_WRITE_LIST"
	IAAS_TENANT_NETWORK_IO_LIST     = "IAAS_TENANT_NETWORK_IO_LIST"
	IAAS_TENANT_NETWORK_PACKET_LIST = "IAAS_TENANT_NETWORK_PACKET_LIST"
	IAAS_LOG_RECENT                 = "IAAS_LOG_RECENT"
	IAAS_LOG_SPECIFIC               = "IAAS_LOG_SPECIFIC"
	IAAS_ALARM_NOTIFICATION_LIST    = "IAAS_ALARM_NOTIFICATION_LIST"
	IAAS_ALARM_NOTIFICATION_CREATE  = "IAAS_ALARM_NOTIFICATION_CREATE"
	IAAS_ALARM_NOTIFICATION_UPDATE  = "IAAS_ALARM_NOTIFICATION_UPDATE"
	IAAS_ALARM_NOTIFICATION_DELETE  = "IAAS_ALARM_NOTIFICATION_DELETE"
	IAAS_ALARM_POLICY_LIST          = "IAAS_ALARM_POLICY_LIST"
	IAAS_ALARM_POLICY               = "IAAS_ALARM_POLICY"
	IAAS_ALARM_POLICY_CREATE        = "IAAS_ALARM_POLICY_CREATE"
	IAAS_ALARM_POLICY_UPDATE        = "IAAS_ALARM_POLICY_UPDATE"
	IAAS_ALARM_POLICY_DELETE        = "IAAS_ALARM_POLICY_DELETE"
	IAAS_ALARM_STATUS_COUNT         = "IAAS_ALARM_STATUS_COUNT"
	IAAS_ALARM_STATUS_LIST          = "IAAS_ALARM_STATUS_LIST"
	IAAS_ALARM_STATUS               = "IAAS_ALARM_STATUS"
	IAAS_ALARM_HISTORY_LIST         = "IAAS_ALARM_HISTORY_LIST"
	IAAS_ALARM_ACTION_LIST          = "IAAS_ALARM_ACTION_LIST"
	IAAS_ALARM_ACTION_CREATE        = "IAAS_ALARM_ACTION_CREATE"
	IAAS_ALARM_ACTION_UPDATE        = "IAAS_ALARM_ACTION_UPDATE"
	IAAS_ALARM_ACTION_DELETE        = "IAAS_ALARM_ACTION_DELETE"

	IAAS_ALARM_REALTIME_COUNT = "IAAS_ALARM_REALTIME_COUNT"
	IAAS_ALARM_REALTIME_LIST  = "IAAS_ALARM_REALTIME_LIST"
)

var IaasRoutes = rata.Routes{
	//{Path: "/v2/dashboard",								Method: "GET", Name: DASHBOARD                              },
	//{Path: "/v2/member/join/check/duplication/iaas/:id", Method: "GET", Name: MEMBER_JOIN_CHECK_DUPLICATION_IAAS_ID},
	//{Path: "/v2/member/join/check/iaas", Method: "POST", Name: MEMBER_JOIN_CHECK_IAAS},

	{Path: "/v2/iaas/main/summary", Method: "GET", Name: IAAS_MAIN_SUMMARY},
	{Path: "/v2/iaas/node/manage/summary", Method: "GET", Name: IAAS_NODE_MANAGE_SUMMARY},
	{Path: "/v2/iaas/node/compute/summary", Method: "GET", Name: IAAS_NODE_COMPUTE_SUMMARY},
	{Path: "/v2/iaas/node/rabbitmq/summary", Method: "GET", Name: IAAS_NODE_RABBITMQ_SUMMARY_OVERVIEW},
	//{Path: "/v2/iaas/node/:hostname/rabbitmq/summary",  	Method: "GET", Name: IAAS_NODE_RABBITMQ_SUMMARY             },

	{Path: "/v2/iaas/nodes", Method: "GET", Name: IAAS_NODES},

	{Path: "/v2/iaas/node/topprocess/:hostname/cpu", Method: "GET", Name: IAAS_NODE_TOPPROCESS_CPU},
	{Path: "/v2/iaas/node/topprocess/:hostname/memory", Method: "GET", Name: IAAS_NODE_TOPPROCESS_MEMORY},

	{Path: "/v2/iaas/node/cpu/:hostname/usages", Method: "GET", Name: IAAS_NODE_CPU_USAGE_LIST},
	{Path: "/v2/iaas/node/cpu/:hostname/loads", Method: "GET", Name: IAAS_NODE_CPU_LOAD_LIST},
	{Path: "/v2/iaas/node/memory/:hostname/swaps", Method: "GET", Name: IAAS_NODE_MEMORY_SWAP_LIST},
	{Path: "/v2/iaas/node/memory/:hostname/usages", Method: "GET", Name: IAAS_NODE_MEMORY_USAGE_LIST},
	{Path: "/v2/iaas/node/disk/:hostname/usages", Method: "GET", Name: IAAS_NODE_DISK_USAGE_LIST},
	{Path: "/v2/iaas/node/disk/:hostname/reads", Method: "GET", Name: IAAS_NODE_DISK_READ_LIST},
	{Path: "/v2/iaas/node/disk/:hostname/writes", Method: "GET", Name: IAAS_NODE_DISK_WRITE_LIST},
	{Path: "/v2/iaas/node/network/:hostname/kbytes", Method: "GET", Name: IAAS_NODE_NETWORK_KBYTE_LIST},
	{Path: "/v2/iaas/node/network/:hostname/errors", Method: "GET", Name: IAAS_NODE_NETWORK_ERROR_LIST},
	{Path: "/v2/iaas/node/network/:hostname/droppackets", Method: "GET", Name: IAAS_NODE_NETWORK_DROP_PACKET_LIST},

	{Path: "/v2/iaas/tenant/summary", Method: "GET", Name: IAAS_TENANT_SUMMARY},
	{Path: "/v2/iaas/tenant/:instanceId/instances", Method: "GET", Name: IAAS_TENANT_INSTANCE_LIST},
	{Path: "/v2/iaas/tenant/cpu/:instanceId/usages", Method: "GET", Name: IAAS_TENANT_CPU_USAGE_LIST},
	{Path: "/v2/iaas/tenant/memory/:instanceId/usages", Method: "GET", Name: IAAS_TENANT_MEMORY_USAGE_LIST},
	{Path: "/v2/iaas/tenant/disk/:instanceId/reads", Method: "GET", Name: IAAS_TENANT_DISK_READ_LIST},
	{Path: "/v2/iaas/tenant/disk/:instanceId/writes", Method: "GET", Name: IAAS_TENANT_DISK_WRITE_LIST},
	{Path: "/v2/iaas/tenant/network/:instanceId/ios", Method: "GET", Name: IAAS_TENANT_NETWORK_IO_LIST},
	{Path: "/v2/iaas/tenant/network/:instanceId/packets", Method: "GET", Name: IAAS_TENANT_NETWORK_PACKET_LIST},
	{Path: "/v2/iaas/log/recent", Method: "GET", Name: IAAS_LOG_RECENT},
	{Path: "/v2/iaas/log/specific", Method: "GET", Name: IAAS_LOG_SPECIFIC},

	//{Path: "/v2/iaas/alarm/notifications", Method: "GET", Name: IAAS_ALARM_NOTIFICATION_LIST},
	//{Path: "/v2/iaas/alarm/notification", Method: "POST", Name: IAAS_ALARM_NOTIFICATION_CREATE},
	//{Path: "/v2/iaas/alarm/notification/:id", Method: "PUT", Name: IAAS_ALARM_NOTIFICATION_UPDATE},
	//{Path: "/v2/iaas/alarm/notification/:id", Method: "DELETE", Name: IAAS_ALARM_NOTIFICATION_DELETE},
	//{Path: "/v2/iaas/alarm/policies", Method: "GET", Name: IAAS_ALARM_POLICY_LIST},
	//{Path: "/v2/iaas/alarm/policy/:id", Method: "GET", Name: IAAS_ALARM_POLICY},
	//{Path: "/v2/iaas/alarm/policy", Method: "POST", Name: IAAS_ALARM_POLICY_CREATE},
	//{Path: "/v2/iaas/alarm/policy/:id", Method: "PATCH", Name: IAAS_ALARM_POLICY_UPDATE},
	//{Path: "/v2/iaas/alarm/policy/:id", Method: "DELETE", Name: IAAS_ALARM_POLICY_DELETE},
	//{Path: "/v2/iaas/alarm/status/count", Method: "GET", Name: IAAS_ALARM_STATUS_COUNT},
	//{Path: "/v2/iaas/alarm/statuses", Method: "GET", Name: IAAS_ALARM_STATUS_LIST},
	//{Path: "/v2/iaas/alarm/status/:alarmId", Method: "GET", Name: IAAS_ALARM_STATUS},
	//{Path: "/v2/iaas/alarm/histories/:alarmId", Method: "GET", Name: IAAS_ALARM_HISTORY_LIST},
	//{Path: "/v2/iaas/alarm/actions/:alarmId", Method: "GET", Name: IAAS_ALARM_ACTION_LIST},
	//{Path: "/v2/iaas/alarm/action", Method: "POST", Name: IAAS_ALARM_ACTION_CREATE},
	//{Path: "/v2/iaas/alarm/action/:id", Method: "PUT", Name: IAAS_ALARM_ACTION_UPDATE},
	//{Path: "/v2/iaas/alarm/action/:id", Method: "DELETE", Name: IAAS_ALARM_ACTION_DELETE},
	//{Path: "/v2/iaas/alarm/realtime/count", Method: "GET", Name: IAAS_ALARM_REALTIME_COUNT},
	//{Path: "/v2/iaas/alarm/realtime/list", Method: "GET", Name: IAAS_ALARM_REALTIME_LIST},
}

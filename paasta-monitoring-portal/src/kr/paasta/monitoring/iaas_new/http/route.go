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

	IAAS_ALARM_POLICY               = "IAAS_ALARM_POLICY"
	IAAS_ALARM_POLICY_CREATE        = "IAAS_ALARM_POLICY_CREATE"

	IAAS_ALARM_POLICY_DELETE        = "IAAS_ALARM_POLICY_DELETE"

	IAAS_ALARM_STATUS               = "IAAS_ALARM_STATUS"
	IAAS_ALARM_HISTORY_LIST         = "IAAS_ALARM_HISTORY_LIST"
	IAAS_ALARM_ACTION_LIST          = "IAAS_ALARM_ACTION_LIST"


	IAAS_ALARM_REALTIME_COUNT = "IAAS_ALARM_REALTIME_COUNT"
	IAAS_ALARM_REALTIME_LIST  = "IAAS_ALARM_REALTIME_LIST"

	IAAS_ALARM_POLICY_LIST   = "IAAS_ALARM_POLICY_LIST"
	IAAS_ALARM_POLICY_UPDATE = "IAAS_ALARM_POLICY_UPDATE"

	IAAS_ALARM_SNS_CHANNEL_LIST   = "IAAS_ALARM_SNS_CHANNEL_LIST"
	IAAS_ALARM_SNS_CHANNEL_CREATE = "IAAS_ALARM_SNS_CHANNEL_CREATE"
	IAAS_ALARM_SNS_CHANNEL_DELETE = "IAAS_ALARM_SNS_CHANNEL_DELETE"
	IAAS_ALARM_SNS_CHANNEL_UPDATE = "IAAS_ALARM_SNS_CHANNEL_UPDATE"

	IAAS_ALARM_STATUS_LIST    = "IAAS_ALARM_STATUS_LIST"
	IAAS_ALARM_STATUS_COUNT   = "IAAS_ALARM_STATUS_COUNT"
	IAAS_ALARM_STATUS_RESOLVE = "IAAS_ALARM_STATUS_RESOLVE"
	IAAS_ALARM_STATUS_DETAIL  = "IAAS_ALARM_DETAIL"
	IAAS_ALARM_STATUS_UPDATE  = "IAAS_ALARM_UPDATE"

	IAAS_ALARM_ACTION_CREATE = "IAAS_ALARM_ACTION_CREATE"
	IAAS_ALARM_ACTION_UPDATE = "IAAS_ALARM_ACTION_UPDATE"
	IAAS_ALARM_ACTION_DELETE = "IAAS_ALARM_ACTION_DELETE"

	IAAS_ALARM_STATISTICS               = "IAAS_ALARM_STATISTICS"
	IAAS_ALARM_STATISTICS_GRAPH_TOTAL   = "IAAS_ALARM_STATISTICS_GRAPH_TOTAL"
	IAAS_ALARM_STATISTICS_GRAPH_SERVICE = "IAAS_ALARM_STATISTICS_GRAPH_SERVICE"
	IAAS_ALARM_STATISTICS_GRAPH_MATRIX  = "IAAS_ALARM_STATISTICS_GRAPH_MATRIX"
	IAAS_ALARM_CONTAINER_DEPLOY         = "IAAS_ALARM_CONTAINER_DEPLOY"
	IAAS_ALARM_DISK_IO_LIST             = "IAAS_ALARM_DISK_IO_LIST"
	IAAS_ALARM_NETWORK_IO_LIST          = "IAAS_ALARM_NETWORK_IO_LIST"
	IAAS_ALARM_TOPPROCESS_LIST          = "IAAS_ALARM_TOPPROCESS_LIST"
	IAAS_ALARM_APP_RESOURCES            = "IAAS_ALARM_APP_RESOURCES"
	IAAS_ALARM_APP_RESOURCES_ALL        = "IAAS_ALARM_APP_RESOURCES_ALL"
	IAAS_ALARM_APP_USAGES               = "IAAS_ALARM_APP_USAGES"
	IAAS_ALARM_APP_MEMORY_USAGES        = "IAAS_ALARM_APP_MEMORY_USAGES"
	IAAS_ALARM_APP_DISK_USAGES          = "IAAS_ALARM_APP_DISK_USAGES"
	IAAS_ALARM_APP_NETWORK_USAGES       = "IAAS_ALARM_APP_NETWORK_USAGES"

	// TODO 2021.11.01 - IAAS 모니터링
	IAAS_GET_HYPER_STATISTICS = "IAAS_GET_HYPER_STATISTICS"
	IAAS_GET_SERVER_LIST = "IAAS_GET_SERVER_LIST"
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


	{Path: "/v2/iaas/alarm/realtime/count", Method: "GET", Name: IAAS_ALARM_REALTIME_COUNT},
	{Path: "/v2/iaas/alarm/realtime/list", Method: "GET", Name: IAAS_ALARM_REALTIME_LIST},

	{Path: "/v2/iaas/alarm/policies", Method: "GET", Name: IAAS_ALARM_POLICY_LIST},
	{Path: "/v2/iaas/alarm/policy", Method: "PUT", Name: IAAS_ALARM_POLICY_UPDATE},

	{Path: "/v2/iaas/alarm/sns/channel", Method: "POST", Name: IAAS_ALARM_SNS_CHANNEL_CREATE},
	{Path: "/v2/iaas/alarm/sns/channel/list", Method: "GET", Name: IAAS_ALARM_SNS_CHANNEL_LIST},
	{Path: "/v2/iaas/alarm/sns/channel/:id", Method: "DELETE", Name: IAAS_ALARM_SNS_CHANNEL_DELETE},
	{Path: "/v2/iaas/alarm/sns/channel", Method: "PUT", Name: IAAS_ALARM_SNS_CHANNEL_UPDATE},  // 2021.05.18 - PaaS 채널 SNS 수정 기능 추가

	{Path: "/v2/iaas/alarm/statuses", Method: "GET", Name: IAAS_ALARM_STATUS_LIST},
	{Path: "/v2/iaas/alarm/status/count", Method: "GET", Name: IAAS_ALARM_STATUS_COUNT},
	{Path: "/v2/iaas/alarm/status/:id", Method: "GET", Name: IAAS_ALARM_STATUS_DETAIL},
	{Path: "/v2/iaas/alarm/status/:id", Method: "PUT", Name: IAAS_ALARM_STATUS_UPDATE},
	{Path: "/v2/iaas/alarm/status/:resolveStatus", Method: "GET", Name: IAAS_ALARM_STATUS_RESOLVE},

	{Path: "/v2/iaas/alarm/action", Method: "POST", Name: IAAS_ALARM_ACTION_CREATE},
	{Path: "/v2/iaas/alarm/action/:actionId", Method: "PATCH", Name: IAAS_ALARM_ACTION_UPDATE},
	{Path: "/v2/iaas/alarm/action/:actionId", Method: "DELETE", Name: IAAS_ALARM_ACTION_DELETE},

	{Path: "/v2/iaas/alarm/statistics", Method: "GET", Name: IAAS_ALARM_STATISTICS},
	{Path: "/v2/iaas/alarm/statistics/graph/total", Method: "GET", Name: IAAS_ALARM_STATISTICS_GRAPH_TOTAL},
	{Path: "/v2/iaas/alarm/statistics/graph/service", Method: "GET", Name: IAAS_ALARM_STATISTICS_GRAPH_SERVICE},
	{Path: "/v2/iaas/alarm/statistics/graph/matrix", Method: "GET", Name: IAAS_ALARM_STATISTICS_GRAPH_MATRIX},

	//{Path: "/v2/iaas/alarm/container/deploy", Method: "GET", Name: IAAS_ALARM_CONTAINER_DEPLOY},
	//{Path: "/v2/iaas/alarm/disk/io/:origin", Method: "GET", Name: IAAS_ALARM_DISK_IO_LIST},
	//{Path: "/v2/iaas/alarm/network/io/:origin", Method: "GET", Name: IAAS_ALARM_NETWORK_IO_LIST},
	//{Path: "/v2/iaas/alarm/topprocess/:origin", Method: "GET", Name: IAAS_ALARM_TOPPROCESS_LIST},
	//{Path: "/v2/iaas/alarm/app/resources", Method: "GET", Name: IAAS_ALARM_APP_RESOURCES},
	//{Path: "/v2/iaas/alarm/app/resources/all", Method: "GET", Name: IAAS_ALARM_APP_RESOURCES_ALL},
	//{Path: "/v2/iaas/alarm/app/cpu/:guid/:idx/usages", Method: "GET", Name: IAAS_ALARM_APP_USAGES},
	//{Path: "/v2/iaas/alarm/app/memory/:guid/:idx/usages", Method: "GET", Name: IAAS_ALARM_APP_MEMORY_USAGES},
	//{Path: "/v2/iaas/alarm/app/disk/:guid/:idx/usages", Method: "GET", Name: IAAS_ALARM_APP_DISK_USAGES},
	//{Path: "/v2/iaas/alarm/app/network/:guid/:idx/usages", Method: "GET", Name: IAAS_ALARM_APP_NETWORK_USAGES},

	//TODO
	{Path: "/v2/iaas/hyper/statistics", Method: "GET", Name: IAAS_GET_HYPER_STATISTICS},
	{Path: "/v2/iaas/server/list", Method: "GET", Name: IAAS_GET_SERVER_LIST},
}

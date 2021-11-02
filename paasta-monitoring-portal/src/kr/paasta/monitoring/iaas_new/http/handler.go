package http

import (
	"github.com/jinzhu/gorm"
	"net/http"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/influxdata/influxdb1-client/v2"
	"github.com/tedsuo/rata"

	"kr/paasta/monitoring/iaas_new/controller"
	"kr/paasta/monitoring/iaas_new/model"
)

func InitHandler(dbConn *gorm.DB, provider model.OpenstackProvider, influx client.Client, elasticsearch *elasticsearch.Client) rata.Handlers {
	mainController := controller.NewMainController(provider, influx)
	computeController := controller.NewComputeController(provider, influx)
	manageNodeController := controller.NewManageNodeController(provider, influx)
	tenantController := controller.NewOpenstackTenantController(provider, influx)
	//notificationController := controller.NewNotificationController(monsClient, iaasInfluxClient)
	//definitionController := controller.NewAlarmDefinitionController(monsClient, iaasInfluxClient)
	//stautsController := controller.NewAlarmStatusController(monsClient, iaasInfluxClient, iaasTxn)
	logController := controller.NewLogController(provider, influx, elasticsearch)

	alarmController := controller.GetAlarmController(dbConn)
	alarmPolicyController := controller.GetAlarmPolicyController(dbConn)

	openstackController := controller.NewOpenstackController(provider)

	handlers := rata.Handlers {
		//routes.MEMBER_JOIN_CHECK_DUPLICATION_IAAS_ID: route(memberController.MemberJoinCheckDuplicationIaasId),
		//routes.MEMBER_JOIN_CHECK_IAAS:                route(memberController.MemberCheckIaaS),

		//Integrated with routes
		IAAS_MAIN_SUMMARY:                  route(mainController.OpenstackSummary),
		IAAS_NODE_COMPUTE_SUMMARY:          route(computeController.NodeSummary),
		IAAS_NODES:                         route(manageNodeController.GetNodeList),

		IAAS_NODE_CPU_USAGE_LIST:           route(computeController.GetCpuUsageList),
		IAAS_NODE_CPU_LOAD_LIST:            route(computeController.GetCpuLoadList),
		IAAS_NODE_MEMORY_SWAP_LIST:         route(computeController.GetMemorySwapList),
		IAAS_NODE_MEMORY_USAGE_LIST:        route(computeController.GetMemoryUsageList),
		IAAS_NODE_DISK_USAGE_LIST:          route(computeController.GetDiskUsageList),
		IAAS_NODE_DISK_READ_LIST:           route(computeController.GetDiskIoReadList),
		IAAS_NODE_DISK_WRITE_LIST:          route(computeController.GetDiskIoWriteList),
		IAAS_NODE_NETWORK_KBYTE_LIST:       route(computeController.GetNetworkInOutKByteList),
		IAAS_NODE_NETWORK_ERROR_LIST:       route(computeController.GetNetworkInOutErrorList),
		IAAS_NODE_NETWORK_DROP_PACKET_LIST: route(computeController.GetNetworkDroppedPacketList),

		IAAS_NODE_MANAGE_SUMMARY:            route(manageNodeController.ManageNodeSummary),
		IAAS_NODE_RABBITMQ_SUMMARY_OVERVIEW: route(manageNodeController.ManageRabbitMqSummary),
		IAAS_NODE_TOPPROCESS_CPU:            route(manageNodeController.GetTopProcessByCpu),
		IAAS_NODE_TOPPROCESS_MEMORY:         route(manageNodeController.GetTopProcessByMemory),

		IAAS_TENANT_SUMMARY:             route(tenantController.TenantSummary),
		IAAS_TENANT_INSTANCE_LIST:       route(tenantController.GetTenantInstanceList),
		IAAS_TENANT_CPU_USAGE_LIST:      route(tenantController.GetInstanceCpuUsageList),
		IAAS_TENANT_MEMORY_USAGE_LIST:   route(tenantController.GetInstanceMemoryUsageList),
		IAAS_TENANT_DISK_READ_LIST:      route(tenantController.GetInstanceDiskReadList),
		IAAS_TENANT_DISK_WRITE_LIST:     route(tenantController.GetInstanceDiskWriteList),
		IAAS_TENANT_NETWORK_IO_LIST:     route(tenantController.GetInstanceNetworkIoList),
		IAAS_TENANT_NETWORK_PACKET_LIST: route(tenantController.GetInstanceNetworkPacketsList),

		IAAS_LOG_RECENT:   route(logController.GetDefaultRecentLog),
		IAAS_LOG_SPECIFIC: route(logController.GetSpecificTimeRangeLog),

		// TODO
		IAAS_GET_HYPER_STATISTICS: route(openstackController.GetHypervisorStatistics),
		IAAS_GET_SERVER_LIST: route(openstackController.GetServerList),


		//iaas.IAAS_ALARM_NOTIFICATION_LIST:   route(notificationController.GetAlarmNotificationList),
		//iaas.IAAS_ALARM_NOTIFICATION_CREATE: route(notificationController.CreateAlarmNotification),
		//iaas.IAAS_ALARM_NOTIFICATION_UPDATE: route(notificationController.UpdateAlarmNotification),
		//iaas.IAAS_ALARM_NOTIFICATION_DELETE: route(notificationController.DeleteAlarmNotification),

		//IAAS_ALARM_POLICY_LIST:   route(definitionController.GetAlarmDefinitionList),
		//IAAS_ALARM_POLICY:        route(definitionController.GetAlarmDefinition),
		//IAAS_ALARM_POLICY_CREATE: route(definitionController.CreateAlarmDefinition),
		//IAAS_ALARM_POLICY_UPDATE: route(definitionController.UpdateAlarmDefinition),
		//IAAS_ALARM_POLICY_DELETE: route(definitionController.DeleteAlarmDefinition),

		//IAAS_ALARM_STATUS_LIST:  route(stautsController.GetAlarmStatusList),
		//IAAS_ALARM_STATUS:       route(stautsController.GetAlarmStatus),
		//IAAS_ALARM_HISTORY_LIST: route(stautsController.GetAlarmHistoryList),
		//iaas.IAAS_ALARM_STATUS_COUNT: route(stautsController.GetAlarmStatusCount),

		//iaas.IAAS_ALARM_ACTION_LIST:   route(stautsController.GetAlarmHistoryActionList),
		//iaas.IAAS_ALARM_ACTION_CREATE: route(stautsController.CreateAlarmHistoryAction),
		//iaas.IAAS_ALARM_ACTION_UPDATE: route(stautsController.UpdateAlarmHistoryAction),
		//iaas.IAAS_ALARM_ACTION_DELETE: route(stautsController.DeleteAlarmHistoryAction),

		//iaas.IAAS_ALARM_REALTIME_COUNT: route(stautsController.GetIaasAlarmRealTimeCount),
		//iaas.IAAS_ALARM_REALTIME_LIST:  route(stautsController.GetIaasAlarmRealTimeList),

		IAAS_ALARM_REALTIME_COUNT: route(alarmController.GetPaasAlarmRealTimeCount),
		IAAS_ALARM_REALTIME_LIST:  route(alarmController.GetPaasAlarmRealTimeList),

		IAAS_ALARM_POLICY_LIST:   route(alarmPolicyController.GetAlarmPolicyList),
		IAAS_ALARM_POLICY_UPDATE: route(alarmPolicyController.UpdateAlarmPolicyList),

		IAAS_ALARM_SNS_CHANNEL_LIST:   route(alarmPolicyController.GetAlarmSnsChannelList),
		IAAS_ALARM_SNS_CHANNEL_CREATE: route(alarmPolicyController.CreateAlarmSnsChannel),
		IAAS_ALARM_SNS_CHANNEL_DELETE: route(alarmPolicyController.DeleteAlarmSnsChannel),
		IAAS_ALARM_SNS_CHANNEL_UPDATE: route(alarmPolicyController.UpdateAlarmSnsChannel),  // 2021.05.18 - PaaS 채널 SNS 정보 수정 기능 추가

		IAAS_ALARM_STATUS_LIST:    route(alarmController.GetAlarmList),
		IAAS_ALARM_STATUS_COUNT:   route(alarmController.GetAlarmListCount),
		IAAS_ALARM_STATUS_RESOLVE: route(alarmController.GetAlarmResolveStatus),
		IAAS_ALARM_STATUS_DETAIL:  route(alarmController.GetAlarmDetail),
		IAAS_ALARM_STATUS_UPDATE:  route(alarmController.UpdateAlarm),
		IAAS_ALARM_ACTION_CREATE:  route(alarmController.CreateAlarmAction),
		IAAS_ALARM_ACTION_UPDATE:  route(alarmController.UpdateAlarmAction),
		IAAS_ALARM_ACTION_DELETE:  route(alarmController.DeleteAlarmAction),

		IAAS_ALARM_STATISTICS:               route(alarmController.GetAlarmStat),
		IAAS_ALARM_STATISTICS_GRAPH_TOTAL:   route(alarmController.GetAlarmStatGraphTotal),
		IAAS_ALARM_STATISTICS_GRAPH_SERVICE: route(alarmController.GetAlarmStatGraphService),
		IAAS_ALARM_STATISTICS_GRAPH_MATRIX:  route(alarmController.GetAlarmStatGraphMatrix),
	}
	return handlers
}

func route(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(f)
}
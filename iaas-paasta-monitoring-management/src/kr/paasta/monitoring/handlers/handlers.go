package handlers

import (
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/cloudfoundry-community/gogobosh"
	"github.com/go-redis/redis"
	monascagopher "github.com/gophercloud/gophercloud"
	"github.com/influxdata/influxdb1-client/v2"
	"github.com/jinzhu/gorm"
	"github.com/monasca/golang-monascaclient/monascaclient"
	"github.com/rackspace/gophercloud"
	/*tokens3 "github.com/rackspace/gophercloud/openstack/identity/v3/tokens"*/
	"github.com/tedsuo/rata"
	"gopkg.in/olivere/elastic.v3"
	"io"
	caasContoller "kr/paasta/monitoring/caas/controller"
	"kr/paasta/monitoring/common/controller"
	iaasContoller "kr/paasta/monitoring/iaas/controller"
	"kr/paasta/monitoring/iaas/model"
	paasContoller "kr/paasta/monitoring/paas/controller"
	pm "kr/paasta/monitoring/paas/model"
	"kr/paasta/monitoring/routes"
	saasContoller "kr/paasta/monitoring/saas/controller"
	"kr/paasta/monitoring/utils"
	"net/http"
	/*"strings"*/
	"time"
)

func NewHandler(openstack_provider model.OpenstackProvider, iaasInfluxClient client.Client, paasInfluxClient client.Client,
	iaasTxn *gorm.DB, paasTxn *gorm.DB, iaasElasticClient *elastic.Client, paasElasticClient *elastic.Client, monsClient monascaclient.Client,
	auth monascagopher.AuthOptions, databases pm.Databases, cfProvider cfclient.Config, rdClient *redis.Client, sysType string, boshClient *gogobosh.Client) http.Handler {

	//Controller선언
	var loginController *controller.LoginController
	var memberController *controller.MemberController

	// SaaS Metrics
	var applicationController *saasContoller.SaasController

	if sysType == utils.SYS_TYPE_IAAS {
		loginController = controller.NewIaasLoginController(openstack_provider, monsClient, auth, paasTxn, rdClient, sysType)
		memberController = controller.NewIaasMemberController(openstack_provider, paasTxn, rdClient, sysType)
	} else if sysType == utils.SYS_TYPE_PAAS {
		loginController = controller.NewPaasLoginController(cfProvider, paasTxn, rdClient, sysType)
		memberController = controller.NewPaasMemberController(cfProvider, paasTxn, rdClient, sysType)
	} else if sysType == utils.SYS_TYPE_SAAS {
		applicationController = saasContoller.GetSaasController(paasTxn)
	} else {
		loginController = controller.NewLoginController(openstack_provider, monsClient, auth, cfProvider, paasTxn, rdClient, sysType)
		memberController = controller.NewMemberController(openstack_provider, cfProvider, paasTxn, rdClient, sysType)
	}

	var mainController *iaasContoller.OpenstackServices
	var computeController *iaasContoller.OpenstackComputeNode
	var manageNodeController *iaasContoller.OpenstackManageNode
	var tenantController *iaasContoller.OpenstackTenant
	var notificationController *iaasContoller.NotificationController
	var definitionController *iaasContoller.AlarmDefinitionController
	var stautsController *iaasContoller.AlarmStatusController
	var logController *iaasContoller.OpenstackLog
	var caasMetricsController *caasContoller.MetricController

	var iaasActions rata.Handlers

	if sysType == utils.SYS_TYPE_IAAS || sysType == utils.SYS_TYPE_ALL {
		mainController = iaasContoller.NewMainController(openstack_provider, iaasInfluxClient)
		computeController = iaasContoller.NewComputeController(openstack_provider, iaasInfluxClient)
		manageNodeController = iaasContoller.NewManageNodeController(openstack_provider, iaasInfluxClient)
		tenantController = iaasContoller.NewOpenstackTenantController(openstack_provider, iaasInfluxClient)
		notificationController = iaasContoller.NewNotificationController(monsClient, iaasInfluxClient)
		definitionController = iaasContoller.NewAlarmDefinitionController(monsClient, iaasInfluxClient)
		stautsController = iaasContoller.NewAlarmStatusController(monsClient, iaasInfluxClient, iaasTxn)
		logController = iaasContoller.NewLogController(openstack_provider, iaasInfluxClient, iaasElasticClient)

		iaasActions = rata.Handlers{
			routes.MEMBER_JOIN_CHECK_DUPLICATION_IAAS_ID: route(memberController.MemberJoinCheckDuplicationIaasId),
			routes.MEMBER_JOIN_CHECK_IAAS:                route(memberController.MemberCheckIaaS),

			//Integrated with routes
			routes.IAAS_MAIN_SUMMARY:         route(mainController.OpenstackSummary),
			routes.IAAS_NODE_COMPUTE_SUMMARY: route(computeController.NodeSummary),
			routes.IAAS_NODES:                route(manageNodeController.GetNodeList),

			routes.IAAS_NODE_CPU_USAGE_LIST:           route(computeController.GetCpuUsageList),
			routes.IAAS_NODE_CPU_LOAD_LIST:            route(computeController.GetCpuLoadList),
			routes.IAAS_NODE_MEMORY_SWAP_LIST:         route(computeController.GetMemorySwapList),
			routes.IAAS_NODE_MEMORY_USAGE_LIST:        route(computeController.GetMemoryUsageList),
			routes.IAAS_NODE_DISK_USAGE_LIST:          route(computeController.GetDiskUsageList),
			routes.IAAS_NODE_DISK_READ_LIST:           route(computeController.GetDiskIoReadList),
			routes.IAAS_NODE_DISK_WRITE_LIST:          route(computeController.GetDiskIoWriteList),
			routes.IAAS_NODE_NETWORK_KBYTE_LIST:       route(computeController.GetNetworkInOutKByteList),
			routes.IAAS_NODE_NETWORK_ERROR_LIST:       route(computeController.GetNetworkInOutErrorList),
			routes.IAAS_NODE_NETWORK_DROP_PACKET_LIST: route(computeController.GetNetworkDroppedPacketList),

			routes.IAAS_NODE_MANAGE_SUMMARY:            route(manageNodeController.ManageNodeSummary),
			routes.IAAS_NODE_RABBITMQ_SUMMARY_OVERVIEW: route(manageNodeController.ManageRabbitMqSummary),
			routes.IAAS_NODE_TOPPROCESS_CPU:            route(manageNodeController.GetTopProcessByCpu),
			routes.IAAS_NODE_TOPPROCESS_MEMORY:         route(manageNodeController.GetTopProcessByMemory),

			routes.IAAS_TENANT_SUMMARY:             route(tenantController.TenantSummary),
			routes.IAAS_TENANT_INSTANCE_LIST:       route(tenantController.GetTenantInstanceList),
			routes.IAAS_TENANT_CPU_USAGE_LIST:      route(tenantController.GetInstanceCpuUsageList),
			routes.IAAS_TENANT_MEMORY_USAGE_LIST:   route(tenantController.GetInstanceMemoryUsageList),
			routes.IAAS_TENANT_DISK_READ_LIST:      route(tenantController.GetInstanceDiskReadList),
			routes.IAAS_TENANT_DISK_WRITE_LIST:     route(tenantController.GetInstanceDiskWriteList),
			routes.IAAS_TENANT_NETWORK_IO_LIST:     route(tenantController.GetInstanceNetworkIoList),
			routes.IAAS_TENANT_NETWORK_PACKET_LIST: route(tenantController.GetInstanceNetworkPacketsList),

			routes.IAAS_LOG_RECENT:   route(logController.GetDefaultRecentLog),
			routes.IAAS_LOG_SPECIFIC: route(logController.GetSpecificTimeRangeLog),

			routes.IAAS_ALARM_NOTIFICATION_LIST:   route(notificationController.GetAlarmNotificationList),
			routes.IAAS_ALARM_NOTIFICATION_CREATE: route(notificationController.CreateAlarmNotification),
			routes.IAAS_ALARM_NOTIFICATION_UPDATE: route(notificationController.UpdateAlarmNotification),
			routes.IAAS_ALARM_NOTIFICATION_DELETE: route(notificationController.DeleteAlarmNotification),

			routes.IAAS_ALARM_POLICY_LIST:   route(definitionController.GetAlarmDefinitionList),
			routes.IAAS_ALARM_POLICY:        route(definitionController.GetAlarmDefinition),
			routes.IAAS_ALARM_POLICY_CREATE: route(definitionController.CreateAlarmDefinition),
			routes.IAAS_ALARM_POLICY_UPDATE: route(definitionController.UpdateAlarmDefinition),
			routes.IAAS_ALARM_POLICY_DELETE: route(definitionController.DeleteAlarmDefinition),

			routes.IAAS_ALARM_STATUS_LIST:  route(stautsController.GetAlarmStatusList),
			routes.IAAS_ALARM_STATUS:       route(stautsController.GetAlarmStatus),
			routes.IAAS_ALARM_HISTORY_LIST: route(stautsController.GetAlarmHistoryList),
			routes.IAAS_ALARM_STATUS_COUNT: route(stautsController.GetAlarmStatusCount),

			routes.IAAS_ALARM_ACTION_LIST:   route(stautsController.GetAlarmHistoryActionList),
			routes.IAAS_ALARM_ACTION_CREATE: route(stautsController.CreateAlarmHistoryAction),
			routes.IAAS_ALARM_ACTION_UPDATE: route(stautsController.UpdateAlarmHistoryAction),
			routes.IAAS_ALARM_ACTION_DELETE: route(stautsController.DeleteAlarmHistoryAction),

			routes.IAAS_ALARM_REALTIME_COUNT: route(stautsController.GetIaasAlarmRealTimeCount),
			routes.IAAS_ALARM_REALTIME_LIST:  route(stautsController.GetIaasAlarmRealTimeList),
		}
	}

	var alarmController *paasContoller.AlarmService
	var alarmPolicyController *paasContoller.AlarmPolicyService
	var containerController *paasContoller.ContainerService
	var metricsController *paasContoller.InfluxServerClient
	var boshStatusController *paasContoller.BoshStatusService
	var paasController *paasContoller.PaasController
	var paasLogController *paasContoller.PaasLogController
	var appController *paasContoller.AppController

	var paasActions rata.Handlers

	if sysType == utils.SYS_TYPE_PAAS || sysType == utils.SYS_TYPE_ALL {
		alarmController = paasContoller.GetAlarmController(paasTxn)
		alarmPolicyController = paasContoller.GetAlarmPolicyController(paasTxn)
		containerController = paasContoller.GetContainerController(paasTxn, paasInfluxClient, databases)
		metricsController = paasContoller.GetMetricsController(paasInfluxClient, databases)
		boshStatusController = paasContoller.GetBoshStatusController(paasTxn, paasInfluxClient, databases)
		paasController = paasContoller.GetPaasController(paasTxn, paasInfluxClient, databases, boshClient)
		paasLogController = paasContoller.NewLogController(paasInfluxClient, paasElasticClient)
		appController = paasContoller.GetAppController(paasTxn)

		paasActions = rata.Handlers{
			routes.MEMBER_JOIN_CHECK_DUPLICATION_PAAS_ID: route(memberController.MemberJoinCheckDuplicationPaasId),
			routes.MEMBER_JOIN_CHECK_PAAS:                route(memberController.MemberCheckPaaS),

			////PAAS///////////////////////////////////////////////////////////////////////
			routes.PAAS_ALARM_REALTIME_COUNT: route(alarmController.GetPaasAlarmRealTimeCount),
			routes.PAAS_ALARM_REALTIME_LIST:  route(alarmController.GetPaasAlarmRealTimeList),

			routes.PAAS_ALARM_POLICY_LIST:   route(alarmPolicyController.GetAlarmPolicyList),
			routes.PAAS_ALARM_POLICY_UPDATE: route(alarmPolicyController.UpdateAlarmPolicyList),

			routes.PAAS_ALARM_SNS_CHANNEL_LIST:   route(alarmPolicyController.GetAlarmSnsChannelList),
			routes.PAAS_ALARM_SNS_CHANNEL_CREATE: route(alarmPolicyController.CreateAlarmSnsChannel),
			routes.PAAS_ALARM_SNS_CHANNEL_DELETE: route(alarmPolicyController.DeleteAlarmSnsChannel),

			routes.PAAS_ALARM_STATUS_LIST:    route(alarmController.GetAlarmList),
			routes.PAAS_ALARM_STATUS_COUNT:   route(alarmController.GetAlarmListCount),
			routes.PAAS_ALARM_STATUS_RESOLVE: route(alarmController.GetAlarmResolveStatus),
			routes.PAAS_ALARM_STATUS_DETAIL:  route(alarmController.GetAlarmDetail),
			routes.PAAS_ALARM_STATUS_UPDATE:  route(alarmController.UpdateAlarm),
			routes.PAAS_ALARM_ACTION_CREATE:  route(alarmController.CreateAlarmAction),
			routes.PAAS_ALARM_ACTION_UPDATE:  route(alarmController.UpdateAlarmAction),
			routes.PAAS_ALARM_ACTION_DELETE:  route(alarmController.DeleteAlarmAction),

			routes.PAAS_ALARM_STATISTICS:               route(alarmController.GetAlarmStat),
			routes.PAAS_ALARM_STATISTICS_GRAPH_TOTAL:   route(alarmController.GetAlarmStatGraphTotal),
			routes.PAAS_ALARM_STATISTICS_GRAPH_SERVICE: route(alarmController.GetAlarmStatGraphService),
			routes.PAAS_ALARM_STATISTICS_GRAPH_MATRIX:  route(alarmController.GetAlarmStatGraphMatrix),
			routes.PAAS_ALARM_CONTAINER_DEPLOY:         route(containerController.GetContainerDeploy),

			// bosh
			routes.PAAS_BOSH_STATUS_OVERVIEW:     route(boshStatusController.GetBoshStatusOverview),
			routes.PAAS_BOSH_STATUS_SUMMARY:      route(boshStatusController.GetBoshStatusSummary),
			routes.PAAS_BOSH_STATUS_TOPPROCESS:   route(boshStatusController.GetBoshStatusTopprocess),
			routes.PAAS_BOSH_CPU_USAGE_LIST:      route(boshStatusController.GetBoshCpuUsageList),
			routes.PAAS_BOSH_CPU_LOAD_LIST:       route(boshStatusController.GetBoshCpuLoadList),
			routes.PAAS_BOSH_MEMORY_USAGE_LIST:   route(boshStatusController.GetBoshMemoryUsageList),
			routes.PAAS_BOSH_DISK_USAGE_LIST:     route(boshStatusController.GetBoshDiskUsageList),
			routes.PAAS_BOSH_DISK_IO_LIST:        route(boshStatusController.GetBoshDiskIoList),
			routes.PAAS_BOSH_NETWORK_BYTE_LIST:   route(boshStatusController.GetBoshNetworkByteList),
			routes.PAAS_BOSH_NETWORK_PACKET_LIST: route(boshStatusController.GetBoshNetworkPacketList),
			routes.PAAS_BOSH_NETWORK_DROP_LIST:   route(boshStatusController.GetBoshNetworkDropList),
			routes.PAAS_BOSH_NETWORK_ERROR_LIST:  route(boshStatusController.GetBoshNetworkErrorList),

			//Application Resources 조회 (2017-08-14 추가)
			//Application cpu, memory, disk usage 정보 조회
			routes.PAAS_ALARM_APP_RESOURCES:     route(metricsController.GetApplicationResources),
			routes.PAAS_ALARM_APP_RESOURCES_ALL: route(metricsController.GetApplicationResourcesAll),
			//Application cpu variation 정보 조회
			routes.PAAS_ALARM_APP_USAGES: route(metricsController.GetAppCpuUsage),
			//Application memory variation 정보 조회
			routes.PAAS_ALARM_APP_MEMORY_USAGES: route(metricsController.GetAppMemoryUsage),
			//Application disk variation 정보 조회
			routes.PAAS_ALARM_APP_DISK_USAGES: route(metricsController.GetDiskUsage),

			//Application network variation 정보 조회
			routes.PAAS_ALARM_APP_NETWORK_USAGES: route(metricsController.GetAppNetworkIoKByte),
			// influxDB에서 조회
			routes.PAAS_ALARM_DISK_IO_LIST:    route(metricsController.GetDiskIOList),
			routes.PAAS_ALARM_NETWORK_IO_LIST: route(metricsController.GetNetworkIOList),
			routes.PAAS_ALARM_TOPPROCESS_LIST: route(metricsController.GetTopProcessList),

			// PaaS Overview
			routes.PAAS_PAASTA_OVERVIEW:          route(paasController.GetPaasOverview),
			routes.PAAS_PAASTA_SUMMARY:           route(paasController.GetPaasSummary),
			routes.PAAS_PAASTA_TOPPROCESS_MEMORY: route(paasController.GetPaasTopProcessMemory),
			routes.PAAS_PAASTA_OVERVIEW_STATUS:   route(paasController.GetPaasOverviewStatus),

			// PaaS Detail
			routes.PAAS_PAASTA_CPU_USAGE:      route(paasController.GetPaasCpuUsage),
			routes.PAAS_PAASTA_CPU_LOAD:       route(paasController.GetPaasCpuLoad),
			routes.PAAS_PAASTA_MEMORY_USAGE:   route(paasController.GetPaasMemoryUsage),
			routes.PAAS_PAASTA_DISK_USAGE:     route(paasController.GetPaasDiskUsage),
			routes.PAAS_PAASTA_DISK_IO:        route(paasController.GetPaasDiskIO),
			routes.PAAS_PAASTA_NETWORK_BYTE:   route(paasController.GetPaasNetworkByte),
			routes.PAAS_PAASTA_NETWORK_PACKET: route(paasController.GetPaasNetworkPacket),
			routes.PAAS_PAASTA_NETWORK_DROP:   route(paasController.GetPaasNetworkDrop),
			routes.PAAS_PAASTA_NETWORK_ERROR:  route(paasController.GetPaasNetworkError),

			// PaaS Dashboard
			routes.PAAS_TOPOLOGICAL_VIEW: route(paasController.GetTopologicalView),

			// Container Overview
			routes.PAAS_CELL_OVERVIEW:          route(containerController.GetCellOverview),
			routes.PAAS_CONTAINER_OVERVIEW:     route(containerController.GetContainerOverview),
			routes.PAAS_CONTAINER_SUMMARY:      route(containerController.GetContainerSummary),
			routes.PAAS_CONTAINER_RELATIONSHIP: route(containerController.GetContainerRelationship),

			routes.PAAS_CELL_OVERVIEW_STATE_LIST:      route(containerController.GetCellOverviewStatusList),
			routes.PAAS_CONTAINER_OVERVIEW_STATE_LIST: route(containerController.GetContainerOverviewStatusList),

			routes.PAAS_CONTAINER_OVERVIEW_MAIN: route(containerController.GetPaasMainContainerView),

			routes.PAAS_CONTAINER_CPU_USAGE_LIST:     route(containerController.GetPaasContainerCpuUsages),
			routes.PAAS_CONTAINER_CPU_LOADS_LIST:     route(containerController.GetPaasContainerCpuLoads),
			routes.PAAS_CONTAINER_MEMORY_USAGE_LIST:  route(containerController.GetPaasContainerMemoryUsages),
			routes.PAAS_CONTAINER_DISK_USAGE_LIST:    route(containerController.GetPaasContainerDiskUsages),
			routes.PAAS_CONTAINER_NETWORK_BYTE_LIST:  route(containerController.GetPaasContainerNetworkBytes),
			routes.PAAS_CONTAINER_NETWORK_DROP_LIST:  route(containerController.GetPaasContainerNetworkDrops),
			routes.PAAS_CONTAINER_NETWORK_ERROR_LIST: route(containerController.GetPaasContainerNetworkErrors),

			routes.PAAS_LOG_RECENT:   route(paasLogController.GetDefaultRecentLog),
			routes.PAAS_LOG_SPECIFIC: route(paasLogController.GetSpecificTimeRangeLog),

			// potal - paas api

			routes.PAAS_APP_CPU_USAGES:     route(metricsController.GetAppCpuUsage),
			routes.PAAS_APP_MEMORY_USAGES:  route(metricsController.GetAppMemoryUsage),
			routes.PAAS_APP_NETWORK_USAGES: route(metricsController.GetAppNetworkIoKByte),

			routes.PAAS_APP_AUTOSCALING_POLICY_UPDATE: route(appController.UpdatePaasAppAutoScalingPolicy),
			routes.PAAS_APP_AUTOSCALING_POLICY_INFO:   route(appController.GetPaasAppAutoScalingPolicy),
			routes.PAAS_APP_POLICY_UPDATE:             route(appController.UpdatePaasAppPolicyInfo),
			routes.PAAS_APP_POLICY_INFO:               route(appController.GetPaasAppPolicyInfo),
			routes.PAAS_APP_ALARM_LIST:                route(appController.GetPaasAppAlarmList),
			routes.PAAS_APP_POLICY_DELETE:             route(appController.DeletePaasAppPolicy),
		}
	}

	commonActions := rata.Handlers{

		routes.PING:   route(loginController.Ping),
		routes.LOGIN:  route(loginController.Login),
		routes.LOGOUT: route(loginController.Logout),

		routes.MEMBER_JOIN_INFO:        route(memberController.MemberJoinInfo),
		routes.MEMBER_JOIN_SAVE:        route(memberController.MemberJoinSave),
		routes.MEMBER_JOIN_CHECK_ID:    route(memberController.MemberCheckId),
		routes.MEMBER_JOIN_CHECK_EMAIL: route(memberController.MemberCheckEmail),

		routes.MEMBER_AUTH_CHECK:  route(memberController.MemberAuthCheck),
		routes.MEMBER_INFO_VIEW:   route(memberController.MemberInfoView),
		routes.MEMBER_INFO_UPDATE: route(memberController.MemberInfoUpdate),
		routes.MEMBER_INFO_DELETE: route(memberController.MemberInfoDelete),

		// Html
		routes.Main: route(loginController.Main),
		//routes.Main: route(mainController.Main),
		routes.Static: route(StaticHandler),
	}

	// =================================================
	// SaaS Metrics
	// =================================================
	var saasActions rata.Handlers
	if sysType == utils.SYS_TYPE_SAAS || sysType == utils.SYS_TYPE_ALL {
		saasActions = rata.Handlers{
			routes.SAAS_API_APPLICATION_LIST:   route(applicationController.GetApplicationList),
			routes.SAAS_API_APPLICATION_STATUS: route(applicationController.GetAgentStatus),
			routes.SAAS_API_APPLICATION_GAUGE:  route(applicationController.GetAgentGaugeTot),
		}
	}

	// =================================================
	// CaaS Metrics
	// =================================================
	caasActions := rata.Handlers{
		routes.CAAS_K8S_CLUSTER_AVG:         route(caasMetricsController.GetClusterAvg),
		routes.CAAS_WORK_NODE_LIST:          route(caasMetricsController.GetWorkNodeList),
		routes.CAAS_WORK_NODE_INFO:          route(caasMetricsController.GetWorkNodeInfo),
		routes.CAAS_CONTIANER_LIST:          route(caasMetricsController.GetContainerList),
		routes.CAAS_CONTIANER_INFO:          route(caasMetricsController.GetContainerInfo),
		routes.CAAS_CONTIANER_LOG:           route(caasMetricsController.GetContainerLog),
		routes.CAAS_CLUSTER_OVERVIEW:        route(caasMetricsController.GetClusterOverView),
		routes.CAAS_WORKLOADS_STATUS:        route(caasMetricsController.GetWorkloadsStatus),
		routes.CAAS_MASTER_NODE_USAGE:       route(caasMetricsController.GetMasterNodeUsage),
		routes.CAAS_WORKLOADS_CONTI_SUMMARY: route(caasMetricsController.GetWorkloadsContiSummary),
		routes.CAAS_WORKLOADS_USAGE:         route(caasMetricsController.GetWorkloadsUsage),
		routes.CAAS_POD_STAT:                route(caasMetricsController.GetPodStatList),
		routes.CAAS_POD_LIST:                route(caasMetricsController.GetPodMetricList),
		routes.CAAS_POD_INFO:                route(caasMetricsController.GetPodInfo),
	}

	var actions rata.Handlers
	var actionlist []rata.Handlers

	var route rata.Routes
	var routeList []rata.Routes

	if sysType == utils.SYS_TYPE_IAAS {
		actionlist = append(actionlist, commonActions)
		actionlist = append(actionlist, iaasActions)
		actions = getActions(actionlist)

		routeList = append(routeList, routes.Routes)
		routeList = append(routeList, routes.IaasRoutes)
		route = getRoutes(routeList)
	} else if sysType == utils.SYS_TYPE_PAAS {
		actionlist = append(actionlist, commonActions)
		actionlist = append(actionlist, paasActions)
		actions = getActions(actionlist)

		routeList = append(routeList, routes.Routes)
		routeList = append(routeList, routes.PaasRoutes)
		route = getRoutes(routeList)
	} else if sysType == utils.SYS_TYPE_SAAS {
		actionlist = append(actionlist, commonActions)
		actionlist = append(actionlist, saasActions)
		actions = getActions(actionlist)

		routeList = append(routeList, routes.Routes)
		routeList = append(routeList, routes.SaasRoutes)
		route = getRoutes(routeList)
	} else if sysType == utils.SYS_TYPE_CAAS {
		actionlist = append(actionlist, commonActions)
		actionlist = append(actionlist, caasActions)
		actions = getActions(actionlist)

		routeList = append(routeList, routes.Routes)
		routeList = append(routeList, routes.CaasRoutes)
		route = getRoutes(routeList)
	} else {
		actionlist = append(actionlist, commonActions)
		actionlist = append(actionlist, iaasActions)
		actionlist = append(actionlist, paasActions)
		actionlist = append(actionlist, saasActions)
		actionlist = append(actionlist, caasActions)

		actions = getActions(actionlist)

		routeList = append(routeList, routes.Routes)
		routeList = append(routeList, routes.PaasRoutes)
		routeList = append(routeList, routes.IaasRoutes)
		routeList = append(routeList, routes.SaasRoutes)
		routeList = append(routeList, routes.CaasRoutes)
		route = getRoutes(routeList)
	}

	handler, err := rata.NewRouter(route, actions)
	if err != nil {
		panic("unable to create router: " + err.Error())
	}
	fmt.Println("Monit Application Started")
	return HttpWrap(handler, rdClient, openstack_provider, cfProvider)
}

//CAAS 모니터링 테스트를 위한 임시방편
//func CaasHandler() http.Handler {
//	var metricsController *caasContoller.MetricController
//
//	var caasActions rata.Handlers
//
//	caasActions = rata.Handlers{
//		routes.CAAS_K8S_CLUSTER_AVG: route(metricsController.GetClusterAvg),
//		routes.CAAS_WORK_NODE_LIST:  route(metricsController.GetWorkNodeList),
//		routes.CAAS_WORK_NODE_INFO:  route(metricsController.GetWorkNodeInfo),
//		routes.CAAS_CONTIANER_LIST:  route(metricsController.GetContainerList),
//		routes.CAAS_CONTIANER_INFO:  route(metricsController.GetContainerInfo),
//		routes.CAAS_CONTIANER_LOG:   route(metricsController.GetContainerLog),
//	}
//
//	var actions rata.Handlers
//	var actionlist []rata.Handlers
//
//	var route rata.Routes
//	var routeList []rata.Routes
//
//	actionlist = append(actionlist, caasActions)
//	actions = getActions(actionlist)
//
//	//routeList = append(routeList, routes.Routes)
//	routeList = append(routeList, routes.CaasRoutes)
//	route = getRoutes(routeList)
//
//	handler, err := rata.NewRouter(route, actions)
//	if err != nil {
//		panic("unable to create router: " + err.Error())
//	}
//	fmt.Println("Caas Monit Application Started")
//	return handler
//}

func getActions(list []rata.Handlers) rata.Handlers {
	actions := make(map[string]http.Handler)

	for _, value := range list {
		for key, val := range value {
			actions[key] = val
		}
	}
	return actions
}

func getRoutes(list []rata.Routes) rata.Routes {
	var rList []rata.Route

	for _, value := range list {
		for _, val := range value {
			rList = append(rList, val)
		}
	}
	return rList
}

func HttpWrap(handler http.Handler, rdClient *redis.Client, openstack_provider model.OpenstackProvider, cfProvider cfclient.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, X-XSRF-TOKEN, Accept-Encoding, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Expose-Headers", "X-XSRF-TOKEN")
		}

		// Stop here if its Preflighted OPTIONS request
		if r.Method == "OPTIONS" {
			return
		}

		// token Pass
		//if r.RequestURI != "/v2/login" && r.RequestURI != "/v2/logout" && !strings.Contains(r.RequestURI, "/v2/member/join") && r.RequestURI != "/v2/ping" && r.RequestURI != "/" && !strings.Contains(r.RequestURI, "/public/") && !strings.Contains(r.RequestURI, "/v2/paas/app/") {
		//	fmt.Println("Request URI :: ", r.RequestURI)
		//
		//	reqToken := r.Header.Get(model.CSRF_TOKEN_NAME)
		//	if reqToken == "0" || reqToken == "null" {
		//		fmt.Println("HttpWrap Hander reqToken is null ")
		//		errMessage := model.ErrMessage{"Message": "UnAuthrized"}
		//		utils.RenderJsonUnAuthResponse(errMessage, http.StatusUnauthorized, w)
		//	} else {
		//		//fmt.Println("HttpWrap Hander reqToken =",len(reqToken),":",reqToken)
		//		//모든 경로의 redis 의 토큰 정보를 확인한다
		//		val := rdClient.HGetAll(reqToken).Val()
		//		if val == nil || len(val) == 0 { // redis 에서 token 정보가 expire 된경우 로그인 화면으로 돌아간다
		//			//fmt.Println("HttpWrap Hander redis.iaas_userid is null ")
		//			errMessage := model.ErrMessage{"Message": "UnAuthrized"}
		//			utils.RenderJsonUnAuthResponse(errMessage, http.StatusUnauthorized, w)
		//		} else {
		//
		//			if strings.Contains(r.RequestURI, "/v2/member") && val["userId"] != "" {
		//
		//				handler.ServeHTTP(w, r)
		//
		//			} else if strings.Contains(r.RequestURI, "/v2/iaas") && val["iaasToken"] != "" && val["iaasUserId"] != "" { // IaaS 토큰 정보가 있는경우
		//
		//				provider1, _, err := utils.GetOpenstackProvider(r)
		//				if err != nil || provider1 == nil {
		//					errMessage := model.ErrMessage{"Message": "UnAuthrized"}
		//					utils.RenderJsonUnAuthResponse(errMessage, http.StatusUnauthorized, w)
		//				} else {
		//					v3Client := NewIdentityV3(provider1)
		//
		//					//IaaS, token 검증
		//					bool, err := tokens3.Validate(v3Client, val["iaasToken"])
		//					if err != nil || bool == false {
		//						//errMessage := model.ErrMessage{"Message": "UnAuthrized"}
		//						//utils.RenderJsonUnAuthResponse(errMessage, http.StatusUnauthorized, w)
		//						fmt.Println("iaas token validate error::", err)
		//						handler.ServeHTTP(w, r)
		//					} else {
		//						//두개 token 이 없는 경우도 고려 해야함
		//						rdClient.Expire(reqToken, 30*60*time.Second)
		//						handler.ServeHTTP(w, r)
		//					}
		//				}
		//
		//			} else if strings.Contains(r.RequestURI, "/v2/paas") && val["paasToken"] != "" && val["paasUserId"] != "" { // PaaS 토큰 정보가 있는경우
		//
		//				// Pass token 검증 로직 추가
		//				//get paas token
		//				cfProvider.Token = val["paasToken"]
		//				client_test, err := cfclient.NewClient(&cfProvider)
		//				errMessage := model.ErrMessage{"Message": "UnAuthrized"}
		//				if err != nil {
		//					utils.RenderJsonUnAuthResponse(errMessage, http.StatusUnauthorized, w)
		//				} else {
		//					_, err01 := client_test.GetToken() // cf token 을 refresh 함
		//					if err01 != nil {
		//						utils.RenderJsonUnAuthResponse(errMessage, http.StatusUnauthorized, w)
		//						return
		//					}
		//					/*
		//						fmt.Println("paas hander token ::: ",token)
		//
		//						token01, err02 := client_test.ListApps()
		//						if err02 != nil {
		//							fmt.Println("paas ListApps error::",token01,":::",err02.Error())
		//						}else{
		//							fmt.Println("paas ListApps info  ::",token01)
		//						}
		//					*/
		//					rdClient.Expire(reqToken, 30*60*time.Second)
		//					handler.ServeHTTP(w, r)
		//				}
		//
		//			} else {
		//				fmt.Println("URL Not All")
		//				//rdClient.Expire(reqToken, 30*60*time.Second)
		//				//handler.ServeHTTP(w, r)
		//			}
		//		}
		//	}
		//} else {
		//	fmt.Println("url pass ::", r.RequestURI)
		//	handler.ServeHTTP(w, r)
		//}
		handler.ServeHTTP(w, r)
	}

}

func route(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(f)
}

const STATIC_URL string = "/public/"
const STATIC_ROOT string = "public/"

func StaticHandler(w http.ResponseWriter, req *http.Request) {
	static_file := req.URL.Path[len(STATIC_URL):]
	if len(static_file) != 0 {
		f, err := http.Dir(STATIC_ROOT).Open(static_file)
		if err == nil {
			content := io.ReadSeeker(f)
			http.ServeContent(w, req, static_file, time.Now(), content)
			return
		}
	}
	http.NotFound(w, req)
}
func NewIdentityV3(client *gophercloud.ProviderClient) *gophercloud.ServiceClient {
	v3Endpoint := client.IdentityBase + "v3/"

	return &gophercloud.ServiceClient{
		ProviderClient: client,
		Endpoint:       v3Endpoint,
	}
}

/* global malarkey:false, moment:false */
(function() {
  'use strict';

  //var apiHost = "http://localhost:8080";
  var apiHost = window.location.protocol+"//"+window.location.host;
  var apiUris = {

    /* -------------------------------------- Common --------------------------------------*/

    join:                           apiHost + "/v2/member/join",
    joinCheckId:                    apiHost + "/v2/member/join/check/id/:id",
    joinCheckEmail:                 apiHost + "/v2/member/join/check/email/:email",
    joinCheckIaas:                  apiHost + "/v2/member/join/check/iaas",
    joinCheckPaas:                  apiHost + "/v2/member/join/check/paas",

    joinCheckDuplicationIaas:       apiHost + "/v2/member/join/check/duplication/iaas/:id",
    joinCheckDuplicationPaas:       apiHost + "/v2/member/join/check/duplication/paas/:id",

    authCheck:                      apiHost + "/v2/member/auth/check/:id",
    memberInfoView:                 apiHost + "/v2/member/info/view",
    memberInfoSave:                 apiHost + "/v2/member/info",

    ping:                           apiHost + "/v2/ping",

    login:                          apiHost + "/v2/login",
    logout:                         apiHost + "/v2/logout",




    /* ------------------------------------ Dashboard ------------------------------------*/

    iaasAlarmRealtimeCount:         apiHost + "/v2/iaas/alarm/realtime/count",
    paasAlarmRealtimeCount:         apiHost + "/v2/paas/alarm/realtime/count",

    /* -------------------------------------- IaaS --------------------------------------*/

    // Main
    iaasMainSummary:                apiHost + "/v2/iaas/main/summary",
    iaasAlarmRealtimeList:          apiHost + "/v2/iaas/alarm/realtime/list",


    // Node(Manage, Compute)
    iaasComputeNodeSummary:         apiHost + "/v2/iaas/node/compute/summary",
    iaasManageNodeSummary:          apiHost + "/v2/iaas/node/manage/summary",

    iaasNodeTopProcessCpu:          apiHost + "/v2/iaas/node/topprocess/:hostname/cpu",
    iaasNodeTopProcessMemory:       apiHost + "/v2/iaas/node/topprocess/:hostname/memory",

    iaasNodeRabbitMqSummary:        apiHost + "/v2/iaas/node/rabbitmq/summary",
    iaasNodeRabbitMqList:           apiHost + "/v2/iaas/node/:hostname/rabbitmq/summary",

    iaasNodeList:                   apiHost + "/v2/iaas/nodes",
    iaasNodeCpuUsageList:           apiHost + "/v2/iaas/node/cpu/:hostname/usages",
    iaasNodeCpuLoadList:            apiHost + "/v2/iaas/node/cpu/:hostname/loads",
    iaasNodeMemorySwapList:         apiHost + "/v2/iaas/node/memory/:hostname/swaps",
    iaasNodeMemoryUsageList:        apiHost + "/v2/iaas/node/memory/:hostname/usages",
    iaasNodeDiskUsageList:          apiHost + "/v2/iaas/node/disk/:hostname/usages",
    iaasNodeDiskReadList:           apiHost + "/v2/iaas/node/disk/:hostname/reads",
    iaasNodeDiskWriteList:          apiHost + "/v2/iaas/node/disk/:hostname/writes",
    iaasNodeNetworkKByteList:       apiHost + "/v2/iaas/node/network/:hostname/kbytes",
    iaasNodeNetworkErrorList:       apiHost + "/v2/iaas/node/network/:hostname/errors",
    iaasNodeNetworkDropPacketList:  apiHost + "/v2/iaas/node/network/:hostname/droppackets",


    // Tenant
    iaasTenantSummary:              apiHost + "/v2/iaas/tenant/summary",

    iaasTenantInstanceList:         apiHost + "/v2/iaas/tenant/:instanceId/instances",
    iaasTenantCpuUsageList:         apiHost + "/v2/iaas/tenant/cpu/:instanceId/usages",
    iaasTenantMemoryUsageList:      apiHost + "/v2/iaas/tenant/memory/:instanceId/usages",
    iaasTenantDiskReadList:         apiHost + "/v2/iaas/tenant/disk/:instanceId/reads",
    iaasTenantDiskWriteList:        apiHost + "/v2/iaas/tenant/disk/:instanceId/writes",
    iaasTenantNetworkKByteList:     apiHost + "/v2/iaas/tenant/network/:instanceId/ios",
    iaasTenantNetworkPacketList:    apiHost + "/v2/iaas/tenant/network/:instanceId/packets",


    // Logs
    iaasDefaultRecentLogs:          apiHost + "/v2/iaas/log/recent",
    iaasSpecificTimeRangeLogs:      apiHost + "/v2/iaas/log/specific",


    // Alarm
    iaasAlarmNotificationList:      apiHost + "/v2/iaas/alarm/notifications",
    iaasAlarmNotification:          apiHost + "/v2/iaas/alarm/notification",
    iaasAlarmNotificationId:        apiHost + "/v2/iaas/alarm/notification/:id",

    iaasAlarmPolicyList:            apiHost + "/v2/iaas/alarm/policies",
    iaasAlarmPolicy:                apiHost + "/v2/iaas/alarm/policy",
    iaasAlarmPolicyId:              apiHost + "/v2/iaas/alarm/policy/:id",

    iaasAlarmStatusList:            apiHost + "/v2/iaas/alarm/statuses",
    // iaasAlarmStatus:                apiHost + "/v2/iaas/alarm/status",
    iaasAlarmStatusCount:           apiHost + "/v2/iaas/alarm/status/count",
    iaasAlarmStatusId:              apiHost + "/v2/iaas/alarm/status/:id",
    iaasAlarmStatusHistoryList:     apiHost + "/v2/iaas/alarm/histories/:alarmId",
    iaasAlarmActionList:            apiHost + "/v2/iaas/alarm/actions/:alarmId",
    iaasAlarmAction:                apiHost + "/v2/iaas/alarm/action",
    iaasAlarmActionId:              apiHost + "/v2/iaas/alarm/action/:alarmId",

    /* -------------------------------------- PaaS --------------------------------------*/

    // Main
    paasTopology:                   apiHost + "/v2/paas/main/topological",
    paasAlarmRealtimeList:          apiHost + "/v2/paas/alarm/realtime/list",
    paasZoneContainerRelationship:  apiHost + "/v2/paas/container/relationship",


    // Bosh
    paasBoshOverview:               apiHost + "/v2/paas/bosh/overview",
    paasBoshSummary:                apiHost + "/v2/paas/bosh/summary",
    paasBoshTopProcessMemory:       apiHost + "/v2/paas/bosh/topprocess/:id/memory",

    paasBoshCpuUsageList:           apiHost + "/v2/paas/bosh/cpu/:id/usages",
    paasBoshCpuLoadAverageList:     apiHost + "/v2/paas/bosh/cpu/:id/loads",
    paasBoshMemoryUsageLis:         apiHost + "/v2/paas/bosh/memory/:id/usages",
    paasBoshDiskUsageList:          apiHost + "/v2/paas/bosh/disk/:id/usages",
    paasBoshDiskIOList:             apiHost + "/v2/paas/bosh/disk/:id/ios",
    paasBoshNetworkIoByteList:      apiHost + "/v2/paas/bosh/network/:id/bytes",
    paasBoshNetworkIoPackteList:    apiHost + "/v2/paas/bosh/network/:id/packets",
    paasBoshNetworkIoDropList:      apiHost + "/v2/paas/bosh/network/:id/drops",
    paasBoshNetworkIoErrorList:     apiHost + "/v2/paas/bosh/network/:id/errors",


    // PaaS-TA
    paasPaastaOverview:             apiHost + "/v2/paas/paasta/overview",
    paasPaastaOverviewList:         apiHost + "/v2/paas/paasta/overview/:status",
    paasPaastaSummary:              apiHost + "/v2/paas/paasta/summary",
    paasPaastaTopProcessMemory:     apiHost + "/v2/paas/paasta/topprocess/:id/memory",

    paasPaastaCpuUsageList:         apiHost + "/v2/paas/paasta/cpu/:id/usages",
    paasPaastaCpuLoadAverageList:   apiHost + "/v2/paas/paasta/cpu/:id/loads",
    paasPaastaMemoryUsageLis:       apiHost + "/v2/paas/paasta/memory/:id/usages",
    paasPaastaDiskUsageList:        apiHost + "/v2/paas/paasta/disk/:id/usages",
    paasPaastaDiskIOList:           apiHost + "/v2/paas/paasta/disk/:id/ios",
    paasPaastaNetworkIoByteList:    apiHost + "/v2/paas/paasta/network/:id/bytes",
    paasPaastaNetworkIoPackteList:  apiHost + "/v2/paas/paasta/network/:id/packets",
    paasPaastaNetworkIoDropList:    apiHost + "/v2/paas/paasta/network/:id/drops",
    paasPaastaNetworkIoErrorList:   apiHost + "/v2/paas/paasta/network/:id/errors",


    // Container
    paasCellOverview:               apiHost + "/v2/paas/cell/overview",
    paasContainerOverview:          apiHost + "/v2/paas/container/overview",
    paasContainerSummary:           apiHost + "/v2/paas/container/summary",
    paasContainerRelationship:      apiHost + "/v2/paas/container/relationship/:name",

    paasCellOverviewList:           apiHost + "/v2/paas/cell/overview/:status",
    paasContainerOverviewList:      apiHost + "/v2/paas/container/overview/:status",

    paasContainerCpuUsageList:      apiHost + "/v2/paas/container/cpu/:id/usages",
    paasContainerCpuLoadList:       apiHost + "/v2/paas/container/cpu/:id/loads",
    paasContainerMemoryUsageList:   apiHost + "/v2/paas/container/memory/:id/usages",
    paasContainerDiskUsageList:     apiHost + "/v2/paas/container/disk/:id/usages",
    paasContainerNetworkIoByteList: apiHost + "/v2/paas/container/network/:id/bytes",
    paasContainerNetworkIoDropList: apiHost + "/v2/paas/container/network/:id/drops",
    paasContainerNetworkIoErrorList:apiHost + "/v2/paas/container/network/:id/errors",


    // Logs
    paasDefaultRecentLogs:          apiHost + "/v2/paas/log/recent",
    paasSpecificTimeRangeLogs:      apiHost + "/v2/paas/log/specific",


    // Alarm
    paasAlarmPolicyList:            apiHost + "/v2/paas/alarm/policies",
    paasAlarmPolicy:                apiHost + "/v2/paas/alarm/policy",

    paasAlarmSnsChannelRegist:      apiHost + "/v2/paas/alarm/sns/channel",
    paasAlarmSnsChannelList:        apiHost + "/v2/paas/alarm/sns/channel/list",
    paasAlarmSnsChannelDelete:      apiHost + "/v2/paas/alarm/sns/channel/:id",

    paasAlarmStatusList:            apiHost + "/v2/paas/alarm/statuses",
    paasAlarmStatusCount:           apiHost + "/v2/paas/alarm/status/count",
    paasAlarmStatus:                apiHost + "/v2/paas/alarm/status",
    paasAlarmStatusId:              apiHost + "/v2/paas/alarm/status/:id",
    paasAlarmStatusResolve:         apiHost + "/v2/paas/alarm/status/resolve/:resolveStatus",

    paasAlarmAction:                apiHost + "/v2/paas/alarm/action",
    paasAlarmActionId:              apiHost + "/v2/paas/alarm/action/:actionId",

    paasAlarmStatisticList:         apiHost + "/v2/paas/alarm/statistics",
    paasAlarmStatisticTotal:        apiHost + "/v2/paas/alarm/statistics/graph/total",
    paasAlarmStatisticService:      apiHost + "/v2/paas/alarm/statistics/graph/service",
    paasAlarmStatisticMatrix:       apiHost + "/v2/paas/alarm/statistics/graph/matrix"

  };


  // IaaS Chart
  var nodeChartConfig = [
    {id: 1, name: 'CPU Usage',               func: 'nodeCpuUsageList',             type: 'lineChart', percent: true,  axisLabel: '%'},
    {id: 2, name: 'CPU Load Average',        func: 'nodeCpuLoad1mList',            type: 'lineChart', percent: false, axisLabel: 'Count per 1 minute'},
    {id: 3, name: 'Swap Usage',              func: 'nodeMemorySwapList',           type: 'lineChart', percent: true,  axisLabel: '%'},
    {id: 4, name: 'Memory Usage',            func: 'nodeMemoryUsageList',          type: 'lineChart', percent: true,  axisLabel: '%'},
    {id: 5, name: 'Disk Usage',              func: 'nodeDiskUsageList',            type: 'lineChart', percent: true,  axisLabel: '%'},
    {id: 6, name: 'Disk IO Read',            func: 'nodeDiskIOReadList',           type: 'lineChart', percent: false, axisLabel: 'KB'},
    {id: 7, name: 'Disk IO Write',           func: 'nodeDiskIOWriteList',          type: 'lineChart', percent: false, axisLabel: 'KB'},
    {id: 8, name: 'Network IO KByte',        func: 'nodeNetworkIOKByteList',       type: 'lineChart', percent: false, axisLabel: 'KB'},
    {id: 9, name: 'Network Error',           func: 'nodeNetworkErrorList',         type: 'lineChart', percent: false, axisLabel: 'Count'},
    {id: 10, name: 'Network Dropped Packet', func: 'nodeNetworkDroppedPacketList', type: 'lineChart', percent: false, axisLabel: 'Count'}
  ];

  var tenantChartConfig = [
    {id: 1, name: 'CPU Usage',          func: 'instanceCpuUsageList',       type: 'lineChart', percent: true,  axisLabel: '%'},
    {id: 2, name: 'Memory Usage',       func: 'instanceMemoryUsageList',    type: 'lineChart', percent: true,  axisLabel: '%'},
    {id: 3, name: 'Disk IO Read',       func: 'instanceDiskIOReadList',     type: 'lineChart', percent: false, axisLabel: 'KB'},
    {id: 4, name: 'Disk IO Write',      func: 'instanceDiskIOWriteList',    type: 'lineChart', percent: false, axisLabel: 'KB'},
    {id: 5, name: 'Network IO KByte',   func: 'instanceNetworkIOKByteList', type: 'lineChart', percent: false, axisLabel: 'KB'},
    {id: 6, name: 'Network IO Packet',  func: 'instanceNetworkPacketList',  type: 'lineChart', percent: false, axisLabel: 'Count'}
  ];


  // PaaS Chart
  var boshChartConfig = [
    {id: 1, name: 'CPU Usage',          func: 'boshCpuUsageList',         type: 'lineChart', percent: true,  axisLabel: '%'},
    {id: 2, name: 'CPU Load Average',   func: 'boshCpuLoadAverageList',   type: 'lineChart', percent: false, axisLabel: 'Count per 1 minute'},
    {id: 3, name: 'Memory Usage',       func: 'boshMemoryUsageList',      type: 'lineChart', percent: true,  axisLabel: '%'},
    {id: 4, name: 'Disk Usage',         func: 'boshDiskUsageList',        type: 'lineChart', percent: true,  axisLabel: '%'},
    {id: 5, name: 'Disk IO',            func: 'boshDiskIOList',           type: 'lineChart', percent: false, axisLabel: 'KByte'},
    {id: 6, name: 'Network IO Byte',    func: 'boshNetworkIoByteList',    type: 'lineChart', percent: false, axisLabel: 'KByte'},
    {id: 7, name: 'Network IO Packtes', func: 'boshNetworkIoPackteList',  type: 'lineChart', percent: false, axisLabel: 'Packets/Sec'},
    {id: 8, name: 'Network IO Drop',    func: 'boshNetworkIoDropList',    type: 'lineChart', percent: false, axisLabel: 'Count'},
    {id: 9, name: 'Network IO Error',   func: 'boshNetworkIoErrorList',   type: 'lineChart', percent: false, axisLabel: 'Count'}
  ];

  var paastaChartConfig = [
    {id: 1, name: 'CPU Usage',          func: 'paastaCpuUsageList',         type: 'lineChart', percent: true,  axisLabel: '%'},
    {id: 2, name: 'CPU Load Average',   func: 'paastaCpuLoadAverageList',   type: 'lineChart', percent: false, axisLabel: 'Count per 1 minute'},
    {id: 3, name: 'Memory Usage',       func: 'paastaMemoryUsageList',      type: 'lineChart', percent: true,  axisLabel: '%'},
    {id: 4, name: 'Disk Usage',         func: 'paastaDiskUsageList',        type: 'lineChart', percent: true,  axisLabel: '%'},
    {id: 5, name: 'Disk IO',            func: 'paastaDiskIOList',           type: 'lineChart', percent: false, axisLabel: 'KByte'},
    {id: 6, name: 'Network IO Byte',    func: 'paastaNetworkIoByteList',    type: 'lineChart', percent: false, axisLabel: 'KByte'},
    {id: 7, name: 'Network IO Packtes', func: 'paastaNetworkIoPackteList',  type: 'lineChart', percent: false, axisLabel: 'Packets/Sec'},
    {id: 8, name: 'Network IO Drop',    func: 'paastaNetworkIoDropList',    type: 'lineChart', percent: false, axisLabel: 'Count'},
    {id: 9, name: 'Network IO Error',   func: 'paastaNetworkIoErrorList',   type: 'lineChart', percent: false, axisLabel: 'Count'}
  ];

  var containerChartConfig = [
    {id: 1, name: 'CPU Usage',          func: 'containerCpuUsageList',       type: 'lineChart', percent: true,  axisLabel: '%'},
    {id: 2, name: 'CPU Load Average',   func: 'containerCpuLoadList',        type: 'lineChart', percent: false, axisLabel: 'Count per 1 minute'},
    {id: 3, name: 'Memory Usage',       func: 'containerMemoryUsageList',    type: 'lineChart', percent: true,  axisLabel: '%'},
    {id: 4, name: 'Disk Usage',         func: 'containerDiskUsageList',      type: 'lineChart', percent: true,  axisLabel: '%'},
    {id: 5, name: 'Network IO Byte',    func: 'containerNetworkIoByteList',  type: 'lineChart', percent: false, axisLabel: 'KByte'},
    {id: 6, name: 'Network IO Drop',    func: 'containerNetworkIoDropList',  type: 'lineChart', percent: false, axisLabel: 'Count'},
    {id: 7, name: 'Network IO Error',   func: 'containerNetworkIoErrorList', type: 'lineChart', percent: false, axisLabel: 'Count'}
  ];

  var paasAlarmStatisticChartConfig = [
    {id: 1, name: 'Total',    func: 'alarmStatisticTotal',    type: 'lineChart', percent: false, axisLabel: 'Count'},
    {id: 2, name: 'Service',  func: 'alarmStatisticService',  type: 'lineChart', percent: false, axisLabel: 'Count'},
    {id: 3, name: 'Matrix',   func: 'alarmStatisticMatrix',   type: 'lineChart', percent: false, axisLabel: 'Count'}
  ];

  var constants = {
    version: '0.0.1',
    expire: '0.0.1'
  };

  angular
    .module('monitoring')
    .constant('malarkey', malarkey)
    .constant('moment', moment)
    .constant('apiUris', apiUris)
    .constant('nodeChartConfig', nodeChartConfig)
    .constant('tenantChartConfig', tenantChartConfig)
    .constant('boshChartConfig', boshChartConfig)
    .constant('paastaChartConfig', paastaChartConfig)
    .constant('containerChartConfig', containerChartConfig)
    .constant('paasAlarmStatisticChartConfig', paasAlarmStatisticChartConfig)
    .constant('constants', constants);

})();

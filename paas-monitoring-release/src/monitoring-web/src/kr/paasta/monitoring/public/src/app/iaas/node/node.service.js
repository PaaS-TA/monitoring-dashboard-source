(function() {
  'use strict';

  angular
    .module('monitoring')
    .factory('manageNodeService', ManageNodeService);

  /** @ngInject */
  function ManageNodeService($http, apiUris) {
    var service = {};

    service.manageNodeSummary = function(hostname){
      var config = {
        params: {'hostname': hostname},
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasManageNodeSummary, config);
    };

    service.manageTopProcessByCpu = function(condition){
      return $http.get(apiUris.iaasNodeTopProcessCpu.replace(":hostname", condition.hostname));
    };

    service.manageTopProcessByMemory = function(condition){
      return $http.get(apiUris.iaasNodeTopProcessMemory.replace(":hostname", condition.hostname));
    };

    service.manageRabbitMqSummary = function(){
      return $http.get(apiUris.iaasNodeRabbitMqSummary);
    };

    return service;
  }

  angular
    .module('monitoring')
    .factory('computeNodeService', ComputeNodeService);

  /** @ngInject */
  function ComputeNodeService($http, apiUris, common) {
    var service = {};

    service.computeNodeSummary = function(hostname){
      var config = {
        params: {'hostname': hostname},
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasComputeNodeSummary, config);
    };

    service.computeTopProcessByCpu = function(condition) {
      return $http.get(apiUris.iaasNodeTopProcessCpu.replace(":hostname", condition.hostname));
    };

    service.computeTopProcessByMemory = function(condition) {
      return $http.get(apiUris.iaasNodeTopProcessMemory.replace(":hostname", condition.hostname));
    };

    service.nodeCpuUsageList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasNodeCpuUsageList.replace(":hostname", condition.hostname), config);
    };

    service.nodeCpuLoad1mList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasNodeCpuLoadList.replace(":hostname", condition.hostname), config);
    };

    service.nodeMemorySwapList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasNodeMemorySwapList.replace(":hostname", condition.hostname), config);
    };

    service.nodeMemoryUsageList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasNodeMemoryUsageList.replace(":hostname", condition.hostname), config);
    };

    service.nodeDiskUsageList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasNodeDiskUsageList.replace(":hostname", condition.hostname), config);
    };

    service.nodeDiskIOReadList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasNodeDiskReadList.replace(":hostname", condition.hostname), config);
    };

    service.nodeDiskIOWriteList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasNodeDiskWriteList.replace(":hostname", condition.hostname), config);
    };

    service.nodeNetworkIOKByteList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasNodeNetworkKByteList.replace(":hostname", condition.hostname), config);
    };

    service.nodeNetworkErrorList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasNodeNetworkErrorList.replace(":hostname", condition.hostname), config);
    };

    service.nodeNetworkDroppedPacketList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasNodeNetworkDropPacketList.replace(":hostname", condition.hostname), config);
    };

    service.nodeRabbitMQList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasNodeRabbitMqList.replace(":hostname", condition.hostname), config);
    };

    return service;
  }
})();

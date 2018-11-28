(function() {
  'use strict';

  angular
    .module('monitoring')
    .factory('tenantService', TenantService);

  /** @ngInject */
  function TenantService($http, apiUris, common) {
    var service = {};

    service.tenantSummary = function(tenantName){
      var config = {
        params: {'tenantName': tenantName},
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasTenantSummary, config);
    };

    service.tenantInstanceList = function(id, params){
      var config = {
        params: params,
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasTenantInstanceList.replace(":instanceId", id), config);
    };

    service.instanceCpuUsageList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasTenantCpuUsageList.replace(":instanceId", condition.instanceId), config);
    };

    service.instanceMemoryUsageList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasTenantMemoryUsageList.replace(":instanceId", condition.instanceId), config);
    };

    service.instanceDiskIOReadList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasTenantDiskReadList.replace(":instanceId", condition.instanceId), config);
    };

    service.instanceDiskIOWriteList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasTenantDiskWriteList.replace(":instanceId", condition.instanceId), config);
    };

    service.instanceNetworkIOKByteList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasTenantNetworkKByteList.replace(":instanceId", condition.instanceId), config);
    };

    service.instanceNetworkPacketList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasTenantNetworkPacketList.replace(":instanceId", condition.instanceId), config);
    };

    return service;
  }
})();

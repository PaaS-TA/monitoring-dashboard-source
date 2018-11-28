(function() {
  'use strict';

  angular
    .module('monitoring')
    .factory('paasBoshService', PaasBoshService);

  /** @ngInject */
  function PaasBoshService($http, apiUris, common) {
    var service = {};

    service.boshOverview = function() {
      var config = {
        headers: {'Accept': 'application/json'}
      };
      return $http.get(apiUris.paasBoshOverview, config);
    };

    service.boshSummary = function(params) {
      var config = {
        params: params,
        headers: {'Accept': 'application/json'}
      };
      return $http.get(apiUris.paasBoshSummary, config);
    };

    service.boshTopProcessMemory = function(condition) {
      var config = {
        params: condition,
        headers: {'Accept': 'application/json'}
      };
      return $http.get(apiUris.paasBoshTopProcessMemory.replace(":id", condition.id), config);
    };


    service.boshCpuUsageList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasBoshCpuUsageList.replace(":id", condition.id), config);
    };

    service.boshCpuLoadAverageList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasBoshCpuLoadAverageList.replace(":id", condition.id), config);
    };

    service.boshMemoryUsageList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasBoshMemoryUsageLis.replace(":id", condition.id), config);
    };

    service.boshDiskUsageList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasBoshDiskUsageList.replace(":id", condition.id), config);
    };

    service.boshDiskIOList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasBoshDiskIOList.replace(":id", condition.id), config);
    };

    service.boshNetworkIoByteList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasBoshNetworkIoByteList.replace(":id", condition.id), config);
    };

    service.boshNetworkIoPackteList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasBoshNetworkIoPackteList.replace(":id", condition.id), config);
    };

    service.boshNetworkIoDropList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasBoshNetworkIoDropList.replace(":id", condition.id), config);
    };

    service.boshNetworkIoErrorList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasBoshNetworkIoErrorList.replace(":id", condition.id), config);
    };

    return service;
  }

})();

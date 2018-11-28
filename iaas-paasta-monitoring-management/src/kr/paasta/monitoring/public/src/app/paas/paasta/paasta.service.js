(function() {
  'use strict';

  angular
    .module('monitoring')
    .factory('paasPaastaService', PaasPaastaService);

  /** @ngInject */
  function PaasPaastaService($http, apiUris, common) {
    var service = {};

    service.paastaOverview = function() {
      var config = {
        headers: {'Accept': 'application/json'}
      };
      return $http.get(apiUris.paasPaastaOverview, config);
    };

    service.paastaOverviewList = function(status) {
      var config = {
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasPaastaOverviewList.replace(":status", status), config);
    };

    service.paastaSummary = function(params) {
      var config = {
        params: params,
        headers: {'Accept': 'application/json'}
      };
      return $http.get(apiUris.paasPaastaSummary, config);
    };

    service.paastaTopProcessMemory = function(condition) {
      var config = {
        params: condition,
        headers: {'Accept': 'application/json'}
      };
      return $http.get(apiUris.paasPaastaTopProcessMemory.replace(":id", condition.id), config);
    };


    service.paastaCpuUsageList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasPaastaCpuUsageList.replace(":id", condition.id), config);
    };

    service.paastaCpuLoadAverageList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasPaastaCpuLoadAverageList.replace(":id", condition.id), config);
    };

    service.paastaMemoryUsageList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasPaastaMemoryUsageLis.replace(":id", condition.id), config);
    };

    service.paastaDiskUsageList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasPaastaDiskUsageList.replace(":id", condition.id), config);
    };

    service.paastaDiskIOList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasPaastaDiskIOList.replace(":id", condition.id), config);
    };

    service.paastaNetworkIoByteList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasPaastaNetworkIoByteList.replace(":id", condition.id), config);
    };

    service.paastaNetworkIoPackteList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasPaastaNetworkIoPackteList.replace(":id", condition.id), config);
    };

    service.paastaNetworkIoDropList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasPaastaNetworkIoDropList.replace(":id", condition.id), config);
    };

    service.paastaNetworkIoErrorList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasPaastaNetworkIoErrorList.replace(":id", condition.id), config);
    };

    return service;
  }

})();

(function() {
  'use strict';

  angular
    .module('monitoring')
    .factory('paasContainerService', PaasContainerService);

  /** @ngInject */
  function PaasContainerService($http, apiUris) {
    var service = {};

    service.cellOverview = function() {
      var config = {
        headers: {'Accept': 'application/json'}
      };
      return $http.get(apiUris.paasCellOverview, config);
    };

    service.containerOverview = function() {
      var config = {
        headers: {'Accept': 'application/json'}
      };
      return $http.get(apiUris.paasContainerOverview, config);
    };

    service.containerSummary = function(params) {
      var config = {
        params: params,
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasContainerSummary, config);
    };

    service.cellOverviewList = function(status) {
      var config = {
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasCellOverviewList.replace(":status", status), config);
    };

    service.containerOverviewList = function(status) {
      var config = {
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasContainerOverviewList.replace(":status", status), config);
    };

    service.containerRelationship = function(name) {
      var config = {
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasContainerRelationship.replace(":name", name), config);
    };

    service.zoneContainerRelationship = function() {
      var config = {
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasZoneContainerRelationship, config);
    };

    return service;
  }

  angular
    .module('monitoring')
    .factory('paasContainerDetailService', PaasContainerDetailService);

  /** @ngInject */
  function PaasContainerDetailService($http, apiUris, common) {
    var service = {};

    service.containerCpuUsageList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasContainerCpuUsageList.replace(":id", condition.id), config);
    };

    service.containerCpuLoadList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasContainerCpuLoadList.replace(":id", condition.id), config);
    };

    service.containerMemoryUsageList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasContainerMemoryUsageList.replace(":id", condition.id), config);
    };

    service.containerDiskUsageList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasContainerDiskUsageList.replace(":id", condition.id), config);
    };

    service.containerNetworkIoByteList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasContainerNetworkIoByteList.replace(":id", condition.id), config);
    };

    service.containerNetworkIoDropList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasContainerNetworkIoDropList.replace(":id", condition.id), config);
    };

    service.containerNetworkIoErrorList = function(condition) {
      var config = {
        params: common.setDtvParam(condition),
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasContainerNetworkIoErrorList.replace(":id", condition.id), config);
    };

    return service;
  }

})();

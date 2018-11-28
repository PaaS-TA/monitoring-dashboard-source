(function() {
  'use strict';

  angular
    .module('monitoring')
    .factory('paasAlarmPolicyService', paasAlarmPolicyService);

  /** @ngInject */
  function paasAlarmPolicyService($http, apiUris) {
    var service = {};

    /********** 알람설정 **********/
    service.alarmPolicyList = function() {
      var config = {
        headers: {'Accept': 'application/json'}
      };
      return $http.get(apiUris.paasAlarmPolicyList, config);
    };

    service.alarmSnsChannelList = function() {
      var config = {
        headers: {'Accept': 'application/json'}
      };
      return $http.get(apiUris.paasAlarmSnsChannelList, config);
    };

    service.channelRegist = function (data) {
      var config = {
        headers : {'Accept': 'application/json'}
      };
      return $http.post(apiUris.paasAlarmSnsChannelRegist, data, config);
    };

    service.channelDelete = function (id) {
      var config = {
        headers : {'Accept': 'application/json'}
      };
      return $http.delete(apiUris.paasAlarmSnsChannelDelete.replace(":id", id), config);
    };

    service.updateAlarmSetup = function(data) {
      var config = {
        headers : {'Content-Type': 'application/json', 'Accept': 'application/json'}
      };
      return $http.put(apiUris.paasAlarmPolicy, data, config);
    };

    return service;
  }


  angular
    .module('monitoring')
    .factory('paasAlarmStatusService', paasAlarmStatusService);

  /** @ngInject */
  function paasAlarmStatusService($http, apiUris) {
    var service = {};

    service.alarmStatusList = function(params) {
      var config = {
        headers: {'Accept': 'application/json'}
      };
      return $http.get(apiUris.paasAlarmStatusList + params, config);
    };

    service.alarmStatusCount = function(params) {
      var config = {
        params: params,
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasAlarmStatusCount, config);
    };

    service.alarmStatusDetail = function(id) {
      var config = {
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasAlarmStatusId.replace(":id", id), config);
    };

    service.alarmStatusUpdate = function(id, data) {
      var config = {
        headers : {'Content-Type': 'application/json', 'Accept': 'application/json'}
      };
      return $http.put(apiUris.paasAlarmStatusId.replace(":id", id), data, config);
    };

    service.alarmActionCreate = function(data){
      var config = {
        headers : {'Accept' : 'application/json'}
      };
      return $http.post(apiUris.paasAlarmAction, data, config);
    };

    service.alarmActionUpdate = function(actionId, data){
      var config = {
        headers : {'Content-Type': 'application/json', 'Accept': 'application/json'}
      };
      return $http.patch(apiUris.paasAlarmActionId.replace(":actionId", actionId), data, config);
    };

    service.alarmActionDelete = function(actionId){
      var config = {
        headers : {'Accept' : 'application/json'}
      };
      return $http.delete(apiUris.paasAlarmActionId.replace(":actionId", actionId), config);
    };

    return service;
  }


  angular
    .module('monitoring')
    .factory('paasAlarmStatisticsService', paasAlarmStatisticsService);

  /** @ngInject */
  function paasAlarmStatisticsService($http, apiUris) {
    var service = {};

    service.alarmStatistics = function(params) {
      var config = {
        headers: {'Accept': 'application/json'}
      };
      return $http.get(apiUris.paasAlarmStatisticList + params, config);
    };

    service.alarmStatisticTotal = function(params) {
      var config = {
        headers: {'Accept': 'application/json'}
      };
      return $http.get(apiUris.paasAlarmStatisticTotal + params, config);
    };

    service.alarmStatisticService = function(params) {
      var config = {
        headers: {'Accept': 'application/json'}
      };
      return $http.get(apiUris.paasAlarmStatisticService + params, config);
    };

    service.alarmStatisticMatrix = function(params) {
      var config = {
        headers: {'Accept': 'application/json'}
      };
      return $http.get(apiUris.paasAlarmStatisticMatrix + params, config);
    };





    return service;
  }

})();


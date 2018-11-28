(function() {
  'use strict';

  angular
    .module('monitoring')
    .factory('iaasAlarmNotificationService', iaasAlarmNotificationService);

  /** @ngInject */
  function iaasAlarmNotificationService($http, apiUris) {
    var service = {};

    service.alarmNotificationList = function(params){
      var config = {
        params: params,
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasAlarmNotificationList, config);
    };

    service.insertAlarmNotification = function(data){
      var config = {
        headers : {'Content-Type': 'application/json', 'Accept': 'application/json'}
      };
      return $http.post(apiUris.iaasAlarmNotification, data, config);
    };

    service.updateAlarmNotification = function(data){
      var config = {
        headers : {'Content-Type': 'application/json', 'Accept': 'application/json'}
      };
      return $http.put(apiUris.iaasAlarmNotificationId.replace(":id", data.id), data, config);
    };

    service.deleteAlarmNotification = function(id){
      var config = {
        headers : {'Accept' : 'application/json'}
      };
      return $http.delete(apiUris.iaasAlarmNotificationId.replace(":id", id), config);
    };

    return service;
  }


  angular
    .module('monitoring')
    .factory('iaasAlarmPolicyService', iaasAlarmPolicyService);

  /** @ngInject */
  function iaasAlarmPolicyService($http, apiUris) {
    var service = {};

    service.alarmPolicyList = function(params){
      var config = {
        params: params,
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasAlarmPolicyList, config);
    };

    service.alarmPolicy = function(id){
      var config = {
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasAlarmPolicyId.replace(":id", id), config);
    };

    service.insertAlarmPolicy = function(data){
      var config = {
        headers : {'Accept' : 'application/json'}
      };
      return $http.post(apiUris.iaasAlarmPolicy, data, config);
    };

    service.updateAlarmPolicy = function(data){
      var config = {
        headers : {'Content-Type': 'application/json', 'Accept': 'application/json'}
      };
      return $http.patch(apiUris.iaasAlarmPolicyId.replace(":id", data.id), data, config);
    };

    service.deleteAlarmPolicy = function(id){
      var config = {
        headers : {'Accept' : 'application/json'}
      };
      return $http.delete(apiUris.iaasAlarmPolicyId.replace(":id", id), config);
    };

    service.nodeList = function(){
      var config = {
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasNodeList, config);
    };

    return service;
  }


  angular
    .module('monitoring')
    .factory('iaasAlarmStatusService', iaasAlarmStatusService);

  /** @ngInject */
  function iaasAlarmStatusService($http, apiUris) {
    var service = {};

    service.alarmStatusList = function(params){
      var config = {
        params: params,
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasAlarmStatusList, config);
    };

    service.alarmStatusCount = function(params){
      var config = {
        params: params,
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasAlarmStatusCount, config);
    };

    service.alarmStatus = function(id){
      var config = {
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasAlarmStatusId.replace(":id", id), config);
    };

    service.alarmStatusHistoryList = function(alarmId, params){
      var config = {
        params: params,
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasAlarmStatusHistoryList.replace(":alarmId", alarmId), config);
    };

    service.alarmActionList = function(alarmId){
      var config = {
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasAlarmActionList.replace(":alarmId", alarmId), config);
    };

    service.insertAlarmAction = function(data){
      var config = {
        headers : {'Accept' : 'application/json'}
      };
      return $http.post(apiUris.iaasAlarmAction, data, config);
    };

    service.updateAlarmAction = function(id, data){
      var config = {
        headers : {'Content-Type' : 'application/json', 'Accept' : 'application/json'}
      };
      return $http.put(apiUris.iaasAlarmActionId.replace(":alarmId", id), data, config);
    };

    service.deleteAlarmAction = function(id){
      var config = {
        headers : {'Accept' : 'application/json'}
      };
      return $http.delete(apiUris.iaasAlarmActionId.replace(":alarmId", id), config);
    };

    return service;
  }

})();

(function() {
  'use strict';

  angular
    .module('monitoring')
    .factory('iaaSMainService', IaaSMainService);

  /** @ngInject */
  function IaaSMainService($http, apiUris) {
    var service = {};

    service.openStackSummary = function() {
      return $http.get(apiUris.iaasMainSummary);
    };

    service.alarmRealtimeList = function(param) {
      var config = {
        params:param,
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasAlarmRealtimeList, config);
    };

    return service;
  }
})();

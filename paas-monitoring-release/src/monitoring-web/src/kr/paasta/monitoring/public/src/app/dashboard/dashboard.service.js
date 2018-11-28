(function() {
  'use strict';

  angular
    .module('monitoring')
    .factory('dashboardService', DashboardService);

  /** @ngInject */
  function DashboardService($http, apiUris) {
    var service = {};

    service.iaasAlarmRealtimeCount = function() {
      var config = {
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasAlarmRealtimeCount, config);
    };

    service.paasAlarmRealtimeCount = function() {
      var config = {
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasAlarmRealtimeCount, config);
    };


    return service;
  }
})();

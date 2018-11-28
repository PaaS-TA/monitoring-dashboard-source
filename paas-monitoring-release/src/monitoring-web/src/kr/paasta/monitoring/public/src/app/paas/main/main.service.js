(function() {
  'use strict';

  angular
    .module('monitoring')
    .factory('paaSMainService', PaaSMainService);

  /** @ngInject */
  function PaaSMainService($http, apiUris) {
    var service = {};

    service.alarmRealtimeList = function(param) {
      var config = {
        params : param,
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasAlarmRealtimeList, config);
    };

    return service;
  }
})();

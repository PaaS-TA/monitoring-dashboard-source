(function() {
  'use strict';

  angular
    .module('monitoring')
    .factory('loginService', LoginService);

  /** @ngInject */
  function LoginService($http, apiUris) {
    var service = {};

    service.ping = function() {
      var config = {
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.ping, config);
    };

    service.authenticate = function(data) {
      var config = {
        headers : {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      };
      return $http.post(apiUris.login, data, config);
    };

    service.logout = function() {
      return $http.post(apiUris.logout);
    };

    return service;
  }
})();

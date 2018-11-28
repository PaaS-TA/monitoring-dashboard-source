(function() {
  'use strict';

  angular
    .module('monitoring')
    .service('authInterceptor', AuthInterceptor)
    .service('XSRFInterceptor', XSRFInterceptor);

  function AuthInterceptor($location, $window, $q, cache) {
    return {
      request: function(config) {
        return config;
      },
      requestError: function(rejection) {
        return $q.reject(rejection);
      },
      response: function(response) {
        return response;
      },
      responseError: function(rejection) {
        if(rejection.status == 401) {
          cache.clear();
          if($location.path() == '/login') {
            return false;
          } else {
            $window.location.href = '#/login';
            $window.location.reload();
            // $location.path('/login');
            // return false;
          }
        }
        return $q.reject(rejection);
      }
    };
  }

  function XSRFInterceptor($location, $window, $q, cache) {
    var XSRFToken;
    return {
      request: function(config) {
        if(!XSRFToken) {
          XSRFToken = cache.getToken();
        }
        config.headers['X-XSRF-TOKEN'] = XSRFToken;
        return config;
      },
      response: function(response) {
        var newToken = response.headers('X-XSRF-TOKEN');
        if(newToken) {
          XSRFToken = newToken;
          cache.putToken(newToken);
        }
        return response;
      },
      responseError: function(rejection) {
        if(rejection.status == 403) {
          cache.clear();
          if($location.path() == '/login') {
            return false;
          } else {
            $window.location.href = '#/login';
            $window.location.reload();
            // $location.path('/login');
            // return false;
          }
        }
        return $q.reject(rejection);
      }
    };
  }
})();

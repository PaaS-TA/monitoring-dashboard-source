'use strict';

angular.module('app')
    .factory('authInterceptor', function($location, $q, $window, cache) {
        var authInterceptor = {
            request: function(config) {
                /*if(!cache.isAuthenticated()) {
                    $q.reject();
                    if($location.path() != '/login') {
                        // $location.path('/login');
                        $window.location.href = '#/login';
                        $window.location.reload();
                    }
                }*/
                return config;
            },
            requestError: function(rejection) {
                return $q.reject(rejection);
            },
            response: function(response) {
                return response;
            },
            responseError: function(rejection) {
                return $q.reject(rejection);
            }
        };
        return authInterceptor;
});
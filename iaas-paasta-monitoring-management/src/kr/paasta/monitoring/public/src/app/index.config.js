(function() {
  'use strict';

  angular
    .module('monitoring')
    .config(config)
    .config(localStorageService);

  /** @ngInject */
  function config($httpProvider, $logProvider) {
    $httpProvider.interceptors.push('authInterceptor');

    $httpProvider.defaults.withCredentials = true;
    $httpProvider.interceptors.push('XSRFInterceptor');

    // Enable log
    $logProvider.debugEnabled(true);
  }

  /** @ngInject */
  function localStorageService($provide, constants) {
    $provide.decorator('localStorageService', function($delegate) {
      //store original get & set methods
      var originalGet = $delegate.get,
        originalSet = $delegate.set;

      /**
       * extending the localStorageService get method
       *
       * @param key
       * @returns {*}
       */
      $delegate.get = function(key) {
        if(originalGet(key)) {
          var data = originalGet(key);

          if(data.expire) {
            var now = Date.now();

            // delete the key if it timed out
            if(data.expire < now) {
              $delegate.remove(key);
              return null;
            } else {
              var expiryDate = Date.now() + (1000 * 60 * 60 * constants.expire);
              originalSet(key, {
                data: data.data,
                expire: expiryDate
              });
            }

            return data.data;
          } else {
            return data;
          }
        } else {
          return null;
        }
      };

      /**
       * set
       * @param key               key
       * @param val               value to be stored
       * @param {int} expires     hours until the localStorage expires (hour)
       */
      $delegate.set = function(key, val, expires) {
        var expiryDate = null;

        if(angular.isNumber(expires)) {
          expiryDate = Date.now() + (1000 * 60 * 60 * expires);
          originalSet(key, {
            data: val,
            expire: expiryDate
          });
        } else {
          originalSet(key, val);
        }
      };

      return $delegate;
    });
  }

})();

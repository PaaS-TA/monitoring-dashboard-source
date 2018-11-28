(function() {
  'use strict';

  angular
    .module('monitoring')
    .run(runBlock);

  /** @ngInject */
  function runBlock(cookies, $rootScope, $window) {
    if(angular.isUndefined(cookies.getDefaultTimeRange())) cookies.putDefaultTimeRange('15m');
    if(angular.isUndefined(cookies.getGroupBy())) cookies.putGroupBy('1m');
    if(angular.isUndefined(cookies.getRefreshTime())) cookies.putRefreshTime('off');

    var unregister = $rootScope.$on('$stateChangeSuccess', function() {
      $window.scrollTo(0, 0);
    });

    if(unregister) { return; }
  }

})();

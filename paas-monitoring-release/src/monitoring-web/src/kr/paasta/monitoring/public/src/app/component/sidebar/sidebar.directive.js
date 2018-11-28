(function() {
  'use strict';

  angular
    .module('monitoring')
    .directive('acmeIaasSidebar', acmeIaasSidebar)
    .directive('acmePaasSidebar', acmePaasSidebar);

  /** @ngInject */
  function acmeIaasSidebar() {
    var directive = {
      restrict: 'E',
      templateUrl: 'app/component/sidebar/iaasSidebar.html',
      scope: {
        creationDate: '=',
        eventHandler: '&ngClick'
      },
      controller: IaaSSidebarController,
      controllerAs: 'iaasSide',
      bindToController: true
    };

    return directive;

    /** @ngInject */
    function IaaSSidebarController($scope, $rootScope, $location) {
      $scope.alarmSelected = false;

      if($location.path().indexOf("/manage_node") > -1) {
        $scope.selected = 'imn';
      } else if($location.path().indexOf("/compute_node") > -1) {
        $scope.selected = 'icn';
      } else if($location.path().indexOf("/tenant") > -1) {
        $scope.selected = 'itt';
      } else if($location.path().indexOf("/alarm/notification") > -1) {
        $scope.selected = 'ial';
        $scope.subSelected = 'ialn';
        $scope.alarmSelected = true;
      } else if($location.path().indexOf("/alarm/policy") > -1) {
        $scope.selected = 'ial';
        $scope.subSelected = 'iapc';
        $scope.alarmSelected = true;
      } else if($location.path().indexOf("/alarm/status") > -1) {
        $scope.selected = 'ial';
        $scope.subSelected = 'ials';
        $scope.alarmSelected = true;
      } else {
        $scope.selected = '';
        $scope.subSelected = '';
      }


      // Go to Page
      $scope.goPage = function(url) {
        $location.path(url);
      };


      // Reload
      $scope.reload = function() {
        $rootScope.$broadcast('broadcast:reload');
      };

    }
  }


  /** @ngInject */
  function acmePaasSidebar() {
    var directive = {
      restrict: 'E',
      templateUrl: 'app/component/sidebar/paasSidebar.html',
      scope: {
        creationDate: '=',
        eventHandler: '&ngClick'
      },
      controller: PaaSSidebarController,
      controllerAs: 'paasSide',
      bindToController: true
    };

    return directive;

    /** @ngInject */
    function PaaSSidebarController($scope, $rootScope, $location) {
      $scope.alarmSelected = false;

      if($location.path().indexOf("/bosh") > -1) {
        $scope.selected = 'pbs';
      } else if($location.path().indexOf("/paasta") > -1) {
        $scope.selected = 'ppt';
      } else if($location.path().indexOf("/container") > -1) {
        $scope.selected = 'pct';
      } else if($location.path().indexOf("/alarm/policy") > -1) {
        $scope.selected = 'pal';
        $scope.subSelected = 'papc';
        $scope.alarmSelected = true;
      } else if($location.path().indexOf("/alarm/status") > -1) {
        $scope.selected = 'pal';
        $scope.subSelected = 'pasu';
        $scope.alarmSelected = true;
      } else if($location.path().indexOf("/alarm/statistics") > -1) {
        $scope.selected = 'pal';
        $scope.subSelected = 'past';
        $scope.alarmSelected = true;
      } else {
        $scope.selected = '';
        $scope.subSelected = '';
      }


      // Go to Page
      $scope.goPage = function(url) {
        $location.path(url);
      };


      // Reload
      $scope.reload = function() {
        $rootScope.$broadcast('broadcast:reload');
      };

    }
  }

})();

(function() {
  'use strict';

  angular
    .module('monitoring')
    .controller('PaaSMainController', PaaSMainController);


  /** @ngInject */
  function PaaSMainController($scope, $timeout, $interval, $location, $exceptionHandler,
                              common, nvd3Generator, cookies,
                              paaSMainService, paasBoshService, paasPaastaService, paasContainerService) {
    var vm = this;
    vm.scope = $scope;
    vm.common = common;
    vm.Math = Math;

    vm.scope.loading = true;

    // Chart Init
    vm.scope.bosh = {
      chart: {
        options: nvd3Generator.donutChart.options(),
        api: {}
      }
    };

    vm.scope.bosh.chart.options.title = {enable: true, text: 'Bosh'};
    vm.scope.bosh.chart.options.chart.height = 150;
    vm.scope.bosh.chart.options.chart.showLabels = false;
    vm.scope.bosh.chart.options.chart.color = ['#00aacc','#e66b6b','#f0a141','#ad6de8'];
    vm.scope.bosh.chart.options.chart.pie = { dispatch: { elementClick: function () {
      $scope.$apply(function() {
        $location.path('/paas/bosh');
      });
    }}};

    vm.scope.paasta = {
      chart: {
        options: nvd3Generator.donutChart.options(),
        api: {}
      }
    };

    vm.scope.paasta.chart.options.title = {enable: true, text: 'PaaS-TA'};
    vm.scope.paasta.chart.options.chart.height = 150;
    vm.scope.paasta.chart.options.chart.showLabels = false;
    vm.scope.paasta.chart.options.chart.color = ['#00aacc','#e66b6b','#f0a141','#ad6de8'];
    vm.scope.paasta.chart.options.chart.pie = { dispatch: { elementClick: function () {
      $scope.$apply(function() {
        $location.path('/paas/paasta');
      });
    }}};

    vm.scope.container = {
      chart: {
        options: nvd3Generator.donutChart.options(),
        api: {}
      }
    };

    vm.scope.container.chart.options.title = {enable: true, text: 'Container'};
    vm.scope.container.chart.options.chart.height = 150;
    vm.scope.container.chart.options.chart.showLabels = false;
    vm.scope.container.chart.options.chart.color = ['#00aacc','#e66b6b','#f0a141','#ad6de8'];
    vm.scope.container.chart.options.chart.pie = { dispatch: { elementClick: function () {
      $scope.$apply(function() {
        $location.path('/paas/container');
      });
    }}};


    // widget events
    vm.scope.events = {
      resize: function(e, scope){
        $timeout(function(){
          if (scope.api && scope.api.update) scope.api.update();
        },200)
      }
    };


    // Summary List
    var summaryCnt = 0;
    (vm.getSummaryList = function() {
      paasBoshService.boshOverview().then(
        function(result) {
          summaryCnt++;
          vm.scope.bosh.chart.data = [];
          if(result.data.total > 0) {
            vm.scope.bosh.chart.data.push({key: "Running", value: result.data.running});
            vm.scope.bosh.chart.data.push({key: "Fail", value: result.data.failed});
            vm.scope.bosh.chart.data.push({key: "Warning", value: result.data.warning});
            vm.scope.bosh.chart.data.push({key: "Critical", value: result.data.critical});
            vm.scope.bosh.chart.options.chart.title = result.data.total;
          }
          vm.scope.loadingSummaryBosh = false;
        },
        function(reason) {
          vm.scope.loadingSummaryBosh = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );

      paasPaastaService.paastaOverview().then(
        function(result) {
          summaryCnt++;
          vm.scope.paasta.chart.data = [];
          if(result.data.total > 0) {
            vm.scope.paasta.chart.data.push({key: "Running", value: result.data.running});
            vm.scope.paasta.chart.data.push({key: "Fail", value: result.data.failed});
            vm.scope.paasta.chart.data.push({key: "Warning", value: result.data.warning});
            vm.scope.paasta.chart.data.push({key: "Critical", value: result.data.critical});
            vm.scope.paasta.chart.options.chart.title = result.data.total;
          }
          vm.scope.loadingSummaryPaasta = false;
        },
        function(reason) {
          vm.scope.loadingSummaryPaasta = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );

      paasContainerService.containerOverview().then(
        function(result) {
          summaryCnt++;
          vm.scope.container.chart.data = [];
          if(result.data.total > 0) {
            vm.scope.container.chart.data.push({key: "Running", value: result.data.running});
            vm.scope.container.chart.data.push({key: "Fail", value: result.data.failed});
            vm.scope.container.chart.data.push({key: "Warning", value: result.data.warning});
            vm.scope.container.chart.data.push({key: "Critical", value: result.data.critical});
            vm.scope.container.chart.options.chart.title = result.data.total;
          }
          vm.scope.loadingSummaryContainer = false;
        },
        function(reason) {
          vm.scope.loadingSummaryContainer = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );

    })();

    // All Chart Drawing Loading
    var stop = $interval(function() {
      if(summaryCnt == 3 || vm.scope.loading == false) {
        $interval.cancel(stop);
        vm.scope.loading = false;
      }
    }, 500, 60);


    // PaaS-TA Service List
    vm.scope.itemsPerPage = 100;
    vm.scope.currentPage = 1;
    (vm.getPaastaSummary = function() {
      if(vm.scope.loading == false) {
        vm.scope.loading = true;
      }
      var params = {
        'pageItems': vm.scope.itemsPerPage,
        'pageIndex': vm.scope.currentPage
      };
      paasPaastaService.paastaSummary(params).then(
        function(result) {
          vm.paastaSummary = result.data.data;
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loadingSummary = false;
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    })();


    // Cell & Container Relationship
    vm.scope.zoneTree = {};
    vm.scope.zoneTreeData = [];
    vm.scope.doing_async = true;

    (vm.getZoneContainerRelationInfo = function() {
      var params = {};
      paasContainerService.zoneContainerRelationship(params).then(
        function(result) {
          vm.relationshipList = result.data;
          vm.scope.loadingRelationship = false;
        },
        function(reason) {
          vm.scope.loadingRelationship = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    })();

    // Tile Color
    vm.getTileColor = function(status) {
      var tileClass = '';
      switch(status) {
        case 'fail' : tileClass = 'tile-fail'; break;
        case 'warning' : tileClass = 'tile-warning'; break;
        case 'critical' : tileClass = 'tile-critical'; break;
        default: tileClass = 'tile'; break;
      }
      return tileClass;
    };

    // Go Container Detail
    vm.goContainerDetail = function(container) {
      if(container.status != 'fail') {
        $location.path('/paas/container/'+container.containerId).search({'name' : container.appName, 'appIndex' : container.appIndex});
      }
    };


    // Alarm View List
    (vm.getAlarmRealtimeList = function() {
      var params = '?pageItems=5&pageIndex=1';
      paaSMainService.alarmRealtimeList(params).then(
        function(result) {
          vm.alarmViewList = result.data.data;
          vm.scope.loadingAlramRealtime = false;
        },
        function(reason) {
          vm.scope.loadingAlramRealtime = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    })();


    // Reload
    vm.scope.$on('broadcast:reload', function() {
      vm.scope.selSearchCondition = 'name';
      vm.scope.searchKeyword = '';

      vm.scope.loadingSummaryContainer = true;
      vm.scope.loadingRelationship = true;
      vm.scope.loadingTree = true;
      vm.scope.loadingAlramRealtime = true;
      vm.scope.loading = true;

      vm.getSummaryList();
      vm.getZoneContainerRelationInfo();
      vm.getPaastaSummary();
      vm.getAlarmRealtimeList();
    });

  }

})();

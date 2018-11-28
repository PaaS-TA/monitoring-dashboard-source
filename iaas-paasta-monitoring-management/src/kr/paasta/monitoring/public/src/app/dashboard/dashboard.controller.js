(function() {
  'use strict';

  angular
    .module('monitoring')
    .controller('DashboardController', DashboardController);

  /** @ngInject */
  function DashboardController($scope, $timeout, $interval, $location, $exceptionHandler,
                                common, nvd3Generator, cache,
                                dashboardService,
                                iaaSMainService, paasBoshService, paasPaastaService, paasContainerService) {
    var vm = this;
    vm.scope = $scope;
    vm.Math = Math;
    vm.common = common;

    vm.scope.loading = true;
    var summaryCnt = 0;

    /** ----------------------- IaaS Summary -----------------------**/
    (vm.getOpenStackSummary = function() {
      iaaSMainService.openStackSummary().then(
        function (result) {
          vm.summary = result.data;
          vm.usageVcpu = Math.round((vm.summary.vcpuUsed / vm.summary.vcpuTotal) * 100);
          vm.usageMemory = Math.round((vm.summary.memoryMbUsed / vm.summary.memoryMbTotal) * 100);
          vm.usageDisk = Math.round((vm.summary.diskGbUsed / vm.summary.diskGbTotal) * 100);
          vm.usageVms = Math.round((vm.summary.vmRunning / vm.summary.vmTotal) * 100);
          vm.init(vm.usageVcpu, vm.usageMemory, vm.usageDisk, vm.usageVms);
          vm.scope.loadingSummary = false;
          vm.scope.loading = false;
          summaryCnt++;
        },
        function(reason) {
          vm.scope.loadingSummary = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    })();

    vm.init = function(cpu, mem, disk, inst) {
      vm.scope.cpuPercent = cpu;
      vm.scope.memPercent = mem;
      vm.scope.diskPercent = disk;
      vm.scope.instance = inst;
    };


    /** ----------------------- PaaS Summary -----------------------**/
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
      resize: function(e, scope) {
        $timeout(function(){
          if(scope.api && scope.api.update) scope.api.update();
        },200)
      }
    };

    vm.scope.config = { visible: false };

    // make chart visible after grid have been created
    $timeout(function() {
      $scope.config.visible = true;
    }, 200);

    // subscribe widget on window resize event
    angular.element(window).on('resize', function(){
      $scope.$broadcast('resize');
    });


    // Summary List
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


    // PaaS Alarm View
    (vm.getIaasAlarmRealtimeCount = function() {
      dashboardService.iaasAlarmRealtimeCount().then(
        function(result) {
          vm.iaasAlarmView = result.data;
          vm.scope.loadingIaasAlram = false;
        },
        function(reason) {
          vm.scope.loadingIaasAlram = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    })();


    // PaaS Alarm View
    (vm.getPaasAlarmRealtimeCount = function() {
      dashboardService.paasAlarmRealtimeCount().then(
        function(result) {
          vm.paasAlarmView = result.data;
          vm.scope.loadingPaasAlram = false;
        },
        function(reason) {
          vm.scope.loadingPaasAlram = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    })();


    // Alarm Status List Move
    vm.scope.goAlarmPage = function (type, status) {
      var url = '/' + type + '/alarm/status'

      var paramName = '';
      if(type == 'iaas') { paramName = 'severity'; }
      else { paramName = '_l'; }

      if(status == 'total') { status = ''; }

      $location.path(url).search(paramName, status);
    }


    // Reload
    vm.scope.$on('broadcast:reload', function() {
      summaryCnt = 0;

      vm.scope.loading = true;
      vm.scope.loadingSummary = true;
      vm.scope.loadingSummaryBosh = true;
      vm.scope.loadingSummaryPaasta = true;
      vm.scope.loadingSummaryContainer = true;
      vm.scope.loadingIaasAlram = true;
      vm.scope.loadingPaasAlram = true;

      vm.getOpenStackSummary();
      vm.getSummaryList();
      vm.getIaasAlarmRealtimeCount();
      vm.getPaasAlarmRealtimeCount();
    });

  }

})();

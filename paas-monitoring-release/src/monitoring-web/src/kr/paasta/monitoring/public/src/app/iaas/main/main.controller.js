(function() {
  'use strict';

  angular
    .module('monitoring')
    .controller('IaaSMainController', IaaSMainController);

  /** @ngInject */
  function IaaSMainController($scope, $timeout, $exceptionHandler, iaaSMainService, computeNodeService, manageNodeService, tenantService) {
    var vm = this;
    vm.scope = $scope;
    vm.Math = Math;

    vm.scope.loading = true;

    (vm.getOpenStackSummary = function() {
      iaaSMainService.openStackSummary().then(
        function (result) {
          vm.summary = result.data;
          vm.usageVcpu = Math.round((vm.summary.vcpuUsed / vm.summary.vcpuTotal) * 100);
          vm.usageMemory = Math.round((vm.summary.memoryMbUsed / vm.summary.memoryMbTotal) * 100);
          vm.usageDisk = Math.round((vm.summary.diskGbUsed / vm.summary.diskGbTotal) * 100);
          vm.usageVms = Math.round((vm.summary.vmRunning / vm.summary.vmTotal) * 100);
          vm.init(vm.usageVcpu, vm.usageMemory, vm.usageDisk, vm.usageVms);
          vm.scope.loading = false;
        },
        function(reason) {
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    })();

    vm.init = function(cpu, mem, disk, inst) {
      vm.scope.cpuPercent = cpu;
      vm.scope.memPercent = mem;
      vm.scope.diskPercent = disk;
      vm.scope.instance = inst;
      /*if(!$scope.$$phase) { // Error: $digest already in progress
        vm.scope.$apply(function() {
          vm.scope.cpuPercent = cpu;
          vm.scope.memPercent = mem;
          vm.scope.diskPercent = disk;
          vm.scope.instance = inst;
        });
      }*/
    };

    (vm.getComputeNodeSummary = function() {
      computeNodeService.computeNodeSummary().then(
        function (result) {
          vm.computeNodeSummary = result.data;
        },
        function(reason) {
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      )
    })();

    (vm.getManageNodeSummary = function() {
      manageNodeService.manageNodeSummary().then(
        function (result) {
          vm.manageNodeSummary = result.data;
        },
        function(reason) {
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      )
    })();

    // Tenant Summary
    (vm.getTenantSummary = function() {
      tenantService.tenantSummary().then(
        function (result) {
          vm.tenantSummary = result.data;
        },
        function(reason) {
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      )
    })();


    // Alarm View List
    (vm.getAlarmRealtimeList = function() {
      var params = '?pageItems=5&pageIndex=1';
      iaaSMainService.alarmRealtimeList(params).then(
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

    /********** reload **********/
    vm.scope.$on('broadcast:reload', function() {
      vm.scope.loading = true;
      vm.getOpenStackSummary();
      vm.getComputeNodeSummary();
      vm.getManageNodeSummary();
      vm.getTenantSummary();
    });
  }
})();

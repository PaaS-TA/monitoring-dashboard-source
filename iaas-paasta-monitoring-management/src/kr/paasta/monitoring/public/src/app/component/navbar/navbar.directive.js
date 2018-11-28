(function() {
  'use strict';

  angular
    .module('monitoring')
    .directive('acmeNavbar', acmeNavbar);

  /** @ngInject */
  function acmeNavbar() {
    var directive = {
      restrict: 'E',
      templateUrl: 'app/component/navbar/navbar.html',
      scope: {
        creationDate: '=',
        eventHandler: '&ngClick'
      },
      controller: NavbarController,
      controllerAs: 'vm',
      bindToController: true/*,
      link: function(scope, elem, attrs) {
        // click 1
        elem.bind('click', function(e) {
          scope.$apply(function () {
            console.log(angular.element(e.target).parents('.dropdown-menu').is('ul'));
            if(angular.element(e.target).parents('.dropdown-menu').is('ul')) {
              // angular.element(e.target).parent().children("active").remove();
              angular.element(e.target).parent().addClass('active');
            }
          });
        });

        // click 2
        scope.eventHandler = function(e) {
          angular.element(e.target).parent().addClass('active');
        };
      }*/
    };

    return directive;

    /** @ngInject */
    function NavbarController($scope, $rootScope, $location, $interval, $timeout, $exceptionHandler,
                              moment, common, cookies, cache, constants,
                              loginService, iaasAlarmStatusService, paasAlarmStatusService) {
      var vm = this;
      vm.scope = $scope;

      if(cache.getUser() != null && cache.getUser() != "") {

        // 사용자아이디
        vm.scope.username = cache.getUser().name;
        vm.scope.email = cache.getUser().email;

        // 시스템 설정 ( IaaS : IaaS 만 사용 , PaaS : PaaS 만 사용, ALL : IaaS, PaaS 모두 사용)
        vm.scope.sysType = cache.getUser().sysType;

        // 로고 이동 URL 변경
        if(vm.scope.sysType == "IaaS") {
          vm.scope.brandUrl = "#/iaas/main";
        } else if(vm.scope.sysType == "PaaS") {
          vm.scope.brandUrl = "#/paas/main";
        } else {
          if(cache.getUser().i2 == "S" && cache.getUser().p2 == "F") {
            vm.scope.brandUrl = "#/iaas/main";
          } else if(cache.getUser().i2 == "F" && cache.getUser().p2 == "S") {
            vm.scope.brandUrl = "#/paas/main";
          } else if(cache.getUser().i2 == "F" && cache.getUser().p2 == "F") {
            vm.scope.brandUrl = "#/member/info";
          } else {
            vm.scope.brandUrl = "#/";
          }
        }

        vm.scope.iaasMenuAuth = cache.getUser().i2;
        vm.scope.paasMenuAuth = cache.getUser().p2;
      }


      // Menu Seleted
      if($location.path().indexOf("/iaas") > -1) {
        vm.scope.selected = 'iaas';
      } else if($location.path().indexOf("/paas") > -1) {
        vm.scope.selected = 'paas';
      } else {
        vm.scope.selected = '';
      }


      // Alarm Total Count
      (vm.scope.getAlarms = function() {
        var params = {
          state: 'ALARM',
          resolveStatus: 1
        };
        vm.scope.alarms = 0;

        if (cache.getUser().sysType == "ALL" || cache.getUser().sysType == "IaaS") {
          iaasAlarmStatusService.alarmStatusCount(params).then(
            function(result) {
              vm.scope.iaasAlarms = result.data.totalCnt;
              vm.scope.alarms = vm.scope.alarms + vm.scope.iaasAlarms;
            },
            function(reason) {
              $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
            }
          );
        }

        if (cache.getUser().sysType == "ALL" || cache.getUser().sysType == "PaaS") {
          paasAlarmStatusService.alarmStatusCount(params).then(
            function(result) {
              vm.scope.paasAlarms = result.data.totalCnt;
              vm.scope.alarms = vm.scope.alarms + vm.scope.paasAlarms;
            },
            function(reason) {
              $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
            }
          );
        }

      })();


      // Spinning
      (vm.scope.spinning = function() {
        vm.scope.spin = true;
        var stop = $interval(function() {
          if(angular.element('body').find('.fa-spinner').is(':visible') == false) {
            $interval.cancel(stop);
            vm.scope.spin = false;
          }
        }, 500);
      })();


      // Loout
      vm.scope.logout = function() {
        loginService.logout().then(
          function() {
            cache.clear();
            $location.path('/login');
          },
          function(reason) {
            $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
            cache.clear();
            $location.path('/login');
          }
        );
      };


      // Reload
      vm.scope.reload = function() {
        vm.scope.getAlarms();
        vm.scope.spinning();
        $rootScope.$broadcast('broadcast:reload');
      };

      /********** TimeRange & GroupBy **********/
      vm.scope.selTimeRange = cookies.getDefaultTimeRange();
      vm.scope.selGroupBy = cookies.getGroupBy();
      vm.scope.selRefreshTime = cookies.getRefreshTime();
      // Time Range 수동설정 달력
      vm.scope.timeRangeFrom = moment();
      vm.scope.timeRangeTo = moment();
      vm.scope.timeRangeTo.startOf('day').fromNow();
      vm.scope.optionsFrom = {format: 'YYYY.MM.DD HH:mm'};
      vm.scope.optionsTo = {format: 'YYYY.MM.DD HH:mm'};
      vm.scope.updateTimeRange = function (dateFrom, dateTo) {
        vm.scope.optionsFrom.maxDate = dateTo;
        vm.scope.optionsTo.minDate = dateFrom;
        vm.scope.optionsFromDate = vm.scope.optionsFrom.maxDate._d;
        vm.scope.optionsToDate = vm.scope.optionsTo.minDate._d;
        if(vm.scope.selTimeRange == 'custom') {
          vm.scope.selGroupBy = common.selectGroupingByCustomTimeRange(dateTo, dateFrom);
        }
      };
      $timeout(function() {
        vm.scope.updateTimeRange(vm.scope.timeRangeTo, vm.scope.timeRangeFrom);
      });

      (vm.scope.getTimeRangeString = function() {
        $timeout(function() {
          if(vm.scope.selTimeRange == 'custom') {
            var toMonth = (new Date(vm.scope.optionsToDate).getMonth()+1).toString().length === 1 ? '0'+(new Date(vm.scope.optionsToDate).getMonth()+1).toString() : (new Date(vm.scope.optionsToDate).getMonth()+1).toString();
            var toDate = new Date(vm.scope.optionsToDate).getDate().toString().length === 1 ? '0'+new Date(vm.scope.optionsToDate).getDate().toString() : new Date(vm.scope.optionsToDate).getDate().toString();
            var toHours = new Date(vm.scope.optionsToDate).getHours().toString().length === 1 ? '0'+new Date(vm.scope.optionsToDate).getHours().toString() : new Date(vm.scope.optionsToDate).getHours().toString();
            var toMinutes = new Date(vm.scope.optionsToDate).getMinutes().toString().length === 1 ? '0'+new Date(vm.scope.optionsToDate).getMinutes().toString() : new Date(vm.scope.optionsToDate).getMinutes().toString();
            var toSeconds = new Date(vm.scope.optionsToDate).getSeconds().toString().length === 1 ? '0'+new Date(vm.scope.optionsToDate).getSeconds().toString() : new Date(vm.scope.optionsToDate).getSeconds().toString();

            var fromMonth = (new Date(vm.scope.optionsFromDate).getMonth()+1).toString().length === 1 ? '0'+(new Date(vm.scope.optionsFromDate).getMonth()+1).toString() : (new Date(vm.scope.optionsFromDate).getMonth()+1).toString();
            var fromDate = new Date(vm.scope.optionsFromDate).getDate().toString().length === 1 ? '0'+new Date(vm.scope.optionsFromDate).getDate().toString() : new Date(vm.scope.optionsFromDate).getDate().toString();
            var fromHours = new Date(vm.scope.optionsFromDate).getHours().toString().length === 1 ? '0'+new Date(vm.scope.optionsFromDate).getHours().toString() : new Date(vm.scope.optionsFromDate).getHours().toString();
            var fromMinutes = new Date(vm.scope.optionsFromDate).getMinutes().toString().length === 1 ? '0'+new Date(vm.scope.optionsFromDate).getMinutes().toString() : new Date(vm.scope.optionsFromDate).getMinutes().toString();
            var fromSeconds = new Date(vm.scope.optionsFromDate).getSeconds().toString().length === 1 ? '0'+new Date(vm.scope.optionsFromDate).getSeconds().toString() : new Date(vm.scope.optionsFromDate).getSeconds().toString();

            var to = new Date(vm.scope.optionsToDate).getFullYear()+
              '.' +toMonth+
              '.' +toDate+
              ' ' +toHours+
              ':' +toMinutes+
              ':' +toSeconds;
            var from = new Date(vm.scope.optionsFromDate).getFullYear()+
              '.'+fromMonth+
              '.'+fromDate+
              ' '+fromHours+
              ':'+fromMinutes+
              ':'+fromSeconds;
            vm.scope.timeRangeString = (to)+' to '+(from);
          } else {
            vm.scope.timeRangeString = angular.element("input[name='radioTimeRange']:checked").parent().text();
          }
        });
      })();


      // 조회주기 및 GroupBy 설정
      vm.scope.saveTimeRange = function () {
        if(vm.scope.selTimeRange == 'custom') {
          cookies.putDefaultTimeRange(vm.scope.selTimeRange);
          cookies.putTimeRangeFrom(Number(vm.scope.timeRangeFrom));
          cookies.putTimeRangeTo(Number(vm.scope.timeRangeTo));
        } else {
          cookies.putDefaultTimeRange(vm.scope.selTimeRange);
          cookies.putGroupBy(vm.scope.selGroupBy);
        }
        cookies.putRefreshTime(vm.scope.selRefreshTime);
        var datas = {
          selTimeRange: vm.scope.selTimeRange,
          selGroupBy: vm.scope.selGroupBy,
          selRefreshTime: vm.scope.selRefreshTime,
          timeRangeFrom: vm.scope.timeRangeFrom,
          timeRangeTo: vm.scope.timeRangeTo
        };
        angular.element('body').find('.modal-backdrop').hide();
        $rootScope.$broadcast('broadcast:saveTimeRange', datas);
        vm.scope.getTimeRangeString();
      };


      // time range 선택 시 그에 해당하는 group by 선택
      vm.scope.selectGroupBy = function() {
        vm.scope.selGroupBy = common.getGroupingByTimeRange(vm.scope.selTimeRange, vm.scope.timeRangeFrom, vm.scope.timeRangeTo);
      };


      /********** Auto Refresh **********/
      var refreshInterval;
      (vm.scope.runRefreshInterval = function() {
        if(cookies.getRefreshTime() !== 'off' && angular.isUndefined(cookies.getRefreshTime())) {
          var refreshTime = common.getMillisecondsRefreshTime(cookies.getRefreshTime());
          refreshInterval = $interval(function() {
            $rootScope.$broadcast('broadcast:reload');
          }, refreshTime);
        }
      })();
      vm.scope.$on('broadcast:reload', function(){
        $interval.cancel(refreshInterval);
        vm.scope.runRefreshInterval();
      });
      vm.scope.$on('$stateChangeStart', function(){
        $interval.cancel(refreshInterval);
      });


      /********** modal **********/
      vm.scope.timeRangeTop = function($event){
        var offsetTop = angular.element($event.target).prop('offsetTop');
        angular.element('.time-range').css('top',(offsetTop+180)-angular.element(window).scrollTop());
      };

    }
  }

})();

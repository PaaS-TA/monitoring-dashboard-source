(function() {
  'use strict';

  angular
    .module('monitoring')
    .controller('TenantController', TenantController);

  /** @ngInject */
  function TenantController($scope, $timeout, $location, $sce, $exceptionHandler, tenantService) {
    var vm = this;
    vm.scope = $scope;
    vm.Math = Math;

    var hash = $location.hash();
    $location.hash('');
    if(hash) {
      vm.searchCondition = hash;
    }

    (vm.getTenantSummary = function() {
      vm.scope.loading = true;
      vm.selectedTenant = false;
      vm.tenantInstanceList = null;
      tenantService.tenantSummary(vm.searchCondition).then(
        function(result) {
          vm.tenantSummary = result.data;
          if(hash) {
            angular.forEach(vm.tenantSummary, function(tenantSummary) {
              if(tenantSummary.name == hash) {
                vm.selectTenant(tenantSummary);
                hash = undefined;
              }
            });
          } else {
            if(vm.tenantSummary) {
              vm.selectTenant(vm.tenantSummary[0]);
            } else {
              vm.selectedTenantName = null;
              vm.searchInstanceName = '';
              vm.tenantInstanceList = null;
              vm.tenantInstanceTotalCount = 0;
            }
          }
          // vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    })();

    var oldObj = {};
    vm.selectTenant = function(obj) {
      oldObj['select'] = '';
      obj['select'] = 'active';
      vm.selectedTenant = obj;
      vm.selectedTenantName = $sce.trustAsHtml('[ <span style="color:#77ae33;">' + obj.name + '</span> ]');
      vm.searchInstanceName = '';
      vm.tenantInstanceList = null;
      vm.tenantInstanceTotalCount = 0;
      marker = '';
      vm.getInstanceList();
      oldObj = obj;
    };

    var limit = 10;
    var marker = '';
    vm.searchInstance = function() {
      if(vm.selectedTenant == false) {
        var message = 'Tenant가 선택되지 않았습니다.';
        $exceptionHandler(message, {code: null, message: message});
        return;
      }
      marker = '';
      vm.tenantInstanceList = null;
      vm.tenantInstanceTotalCount = 0;
      vm.getInstanceList();
    };
    vm.getInstanceList = function() {
      if(vm.selectedTenant == false) {
        var message = 'Tenant가 선택되지 않았습니다.';
        $exceptionHandler(message, {code: null, message: message});
        return;
      }
      vm.scope.loading = true;

      var params = {
        'hostname': vm.searchInstanceName,
        'limit': limit,
        'marker': marker
      };
      tenantService.tenantInstanceList(vm.selectedTenant.id, params).then(
        function(result) {
          if(result.data.metric && vm.selectedTenant.id == result.data.tenantId) {
            vm.tenantInstanceList = vm.tenantInstanceList == null ? [] : vm.tenantInstanceList;
            vm.tenantInstanceList = vm.tenantInstanceList.concat(result.data.metric);
            vm.tenantInstanceTotalCount = result.data.totalCnt;
            vm.moreButton = '<strong>더 보 기</strong> (총 ' + vm.tenantInstanceTotalCount + ' 건)';
            if(vm.tenantInstanceList) {
              marker = vm.tenantInstanceList[(vm.tenantInstanceList.length-1)].instance_id;
            }
            if(vm.tenantInstanceList.length >= vm.tenantInstanceTotalCount) {
              vm.moreButton = '(총 ' + vm.tenantInstanceTotalCount + '건)';
            }
          }
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };

    /********** reload **********/
    vm.scope.$on('broadcast:reload', function() {
      vm.getTenantSummary();
    });
  }


  angular
    .module('monitoring')
    .controller('TenantDetailController', TenantDetailController);

  /** @ngInject */
  function TenantDetailController($scope, $log, $stateParams, $interval, $timeout, $location,
                                   tenantService, tenantChartConfig, common, cookies, nvd3Generator, $exceptionHandler) {
    var vm = this;
    vm.scope = $scope;
    vm.Math = Math;

    vm.scope.loading = true;

    vm.scope.name = $location.search().name;

    /********** Chart **********/
    vm.scope.gridsterOpts = {
      margins: [20, 20],
      columns: 3,
      rowHeight: 250,
      mobileModeEnabled: false,
      swapping: true,
      draggable: {
        handle: 'h4',
        stop: function(event, $element, widget) {
          $timeout(function(){
            $log.log(event+','+$element+','+widget);
          },400)
        }
      },
      resizable: {
        enabled: true,
        handles: ['n', 'e', 's', 'w', 'ne', 'se', 'sw', 'nw'],
        minWidth: 200,
        layoutChanged: function() {
        },

        // optional callback fired when resize is started
        start: function(event, $element, widget) {
          $timeout(function(){
            $log.log(event+','+$element+','+widget);
          },400)
        },

        // optional callback fired when item is resized,
        resize: function(event, $element, widget) {
          if (widget.chart.api) widget.chart.api.update();
        },

        // optional callback fired when item is finished resizing
        stop: function(event, $element, widget) {
          $timeout(function(){
            if (widget.chart.api) widget.chart.api.update();
          },400)
        }
      }
    };

    vm.scope.dashboard = {
      widgets: []
    };

    vm.scope.events = {
      resize: function(e, scope){
        $timeout(function(){
          if (scope.api && scope.api.update) scope.api.update();
        },200)
      }
    };

    vm.scope.config = { visible: false };

    $timeout(function(){
      vm.scope.config.visible = true;
    }, 200);

    angular.element(window).on('resize', function(){
      vm.scope.$broadcast('resize');
    });

    var charts = tenantChartConfig;

    // 조회조건 설정
    var instanceId = $stateParams.instanceId;
    var condition = {
      instanceId: instanceId,
      groupBy: cookies.getGroupBy()
    };
    if(cookies.getDefaultTimeRange() == 'custom') {
      condition['timeRangeFrom'] = common.timeDifference(new Date().getTime(), cookies.getTimeRangeFrom());
      condition['timeRangeTo'] = common.timeDifference(new Date().getTime(), cookies.getTimeRangeTo());
    } else {
      condition['defaultTimeRange'] = cookies.getDefaultTimeRange();
    }

    var count = 0;
    for(var i in charts) {
      var chartOpt = charts[i];
      if(chartOpt.func) {
        (function(opt, cnt) {
          tenantService[opt.func](condition).then(
            function (result) {
              // if(opt.func == 'nodeNetworkDroppedPacketList') console.info(JSON.stringify(result.data));
              vm.scope.setWidget(cnt, opt, result.data);
            },
            function (reason) {
              vm.scope.setWidget(cnt, opt);
              $log.error(reason);
              $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
            }
          );
        })(chartOpt, count);
        count++;
      }
    }
    var stop = $interval(function() {
      if(count == (vm.scope.dashboard.widgets).length) {
        $interval.cancel(stop);
        vm.scope.dashboard.widgets.sort(common.CompareForSort);
        vm.scope.loading = false;
      }
    }, 500);

    vm.scope.setWidget = function(index, dtvOpt, jsonArr) {
      var col = dtvOpt.col==undefined?Math.floor(index%3):dtvOpt.col;
      var row = dtvOpt.row==undefined?Math.floor(index/3):dtvOpt.row;
      var sizeX = dtvOpt.sizeX==undefined?1:dtvOpt.sizeX;
      var sizeY = dtvOpt.sizeY==undefined?1:dtvOpt.sizeY;

      var widget = {
        col: col, row: row, sizeX: sizeX, sizeY: sizeY, name: dtvOpt.name, id: dtvOpt.id, type: dtvOpt.type,
        chart: {
          options: nvd3Generator[dtvOpt.type].options(),
          api: {}
        },
        func: dtvOpt.func
      };
      if(dtvOpt.percent && jsonArr) {
        widget.chart.options.chart.forceY = [0, 100];
        widget.chart.options.chart.yAxis.tickFormat = function (d) { return d3.format('.0%')(d/100); };
        widget.percent = dtvOpt.percent;
      }
      if(dtvOpt.axisLabel) {
        widget.chart.options.chart.yAxis.axisLabel = dtvOpt.axisLabel;
        widget.chart.options.chart.yAxis.axisLabelDistance = -5;
        widget.chart.options.chart.margin.left = 55;
        widget.axisLabel = dtvOpt.axisLabel;
      }
      if(jsonArr) {
        var value, arr = [];
        var checkData = 0;
        for(var i in jsonArr) {
          value = jsonArr[i].metric==null?[{time:0,usage:0}]:jsonArr[i].metric;
          arr.push({values: value, key: jsonArr[i].name});
          for (var j in value) {
            if(value[j] != null) {
              if (checkData < value[j].usage) {
                checkData = value[j].usage;
              }
            }
          }
        }
        widget.chart.data = arr;

        if(dtvOpt.type != 'list') {
          if(checkData < 5 && widget.chart.options.chart.forceY == undefined) {
            widget.chart.options.chart.forceY = [0, 5];
          }
          if(checkData > 10000) {
            widget.chart.options.chart.yAxis.axisLabelDistance = 20;
            widget.chart.options.chart.margin.left = 80;
          }
        }
      } else {
        widget.chart.options.chart.forceY = false;
      }
      vm.scope.dashboard.widgets.push(widget);
    };

    // 조회주기 및 GroupBy 설정
    var savedCustom = false;
    vm.scope.saveTimeRange = function () {
      var instanceId = $stateParams.instanceId;
      var condition = {
        instanceId: instanceId,
        groupBy: cookies.getGroupBy()
      };
      if(vm.scope.selTimeRange == 'custom') {
        condition['timeRangeFrom'] = common.timeDifference(new Date().getTime(), vm.scope.timeRangeFrom);
        condition['timeRangeTo'] = common.timeDifference(new Date().getTime(), vm.scope.timeRangeTo);
        savedCustom = true;
      } else {
        condition['defaultTimeRange'] = vm.scope.selTimeRange;

        condition.defaultTimeRange = vm.scope.selTimeRange;
        condition.groupBy = vm.scope.selGroupBy;
        savedCustom = false;
      }
      angular.forEach(vm.scope.dashboard.widgets, function(widget, index) {
        (function(opt, idx) {
          tenantService[opt.func](condition).then(
            function (result) {
              if(result) {
                var jsonArr = result.data;
                var value, arr = [];
                for(var i in jsonArr) {
                  value = jsonArr[i].metric==null?[{time:0,usage:0}]:jsonArr[i].metric;
                  arr.push({values: value, key: jsonArr[i].name});
                }
                vm.scope.dashboard.widgets[idx].chart.data = arr;
                vm.scope.dashboard.widgets[idx].loading = false;
              }
            },
            function (reason, status) {
              $timeout(function() { $exceptionHandler(reason.Message, {code: status, message: reason.Message}); }, 500);
            }
          );
        })(widget, index);
      });
    };

    // time range 선택 시 그에 해당하는 group by 선택
    vm.scope.selectGroupBy = function() {
      vm.scope.selGroupBy = common.getGroupingByTimeRange(vm.scope.selTimeRange, vm.scope.timeRangeFrom, vm.scope.timeRangeTo);
    };

    // 팝업에서 save 하지 않은 경우 원래 값을 유지
    angular.element('#timeRange').on('hidden.bs.modal', function () {
      if(savedCustom != true) {
        vm.scope.selTimeRange = cookies.getDefaultTimeRange();
        vm.scope.selGroupBy = cookies.getGroupBy();
      }
    });

    // Overview Page
    vm.scope.goMovePage = function() {
      $location.path("/iaas/tenant");
    };

    // Reload
    vm.scope.$on('broadcast:reload', function() {
      if(vm.scope.selectedTab == 'logs') {
        vm.scope.logSearch();
      } else {
        angular.forEach(vm.scope.dashboard.widgets, function(widget) {
          widget.loading = true;
        });
        vm.scope.saveTimeRange();
      }
    });
    vm.scope.$on('broadcast:saveTimeRange', function (event, data) {
      vm.scope.selTimeRange = data.selTimeRange;
      vm.scope.selGroupBy = data.selGroupBy;
      vm.scope.selRefreshTime = data.selRefreshTime;
      vm.scope.timeRangeFrom = data.timeRangeFrom;
      vm.scope.timeRangeTo = data.timeRangeTo;
      vm.scope.saveTimeRange();
    });
  }
})();

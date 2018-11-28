(function() {
  'use strict';

  angular
    .module('monitoring')
    .controller('ManageNodeController', ManageNodeController);

  /** @ngInject */
  function ManageNodeController($scope, $timeout, $window, $sce, $exceptionHandler, manageNodeService) {
    var vm = this;
    vm.scope = $scope;
    vm.Math = Math;

    vm.scope.loading = true;

    (vm.getManageNodeSummary = function() {
      if(vm.loadingSummary != undefined) {
        vm.loadingSummary = true;
      }
      manageNodeService.manageNodeSummary(vm.searchCondition).then(
        function(result) {
          vm.manageNodeSummary = result.data;
          vm.scope.loading = false;
          vm.loadingSummary = false;
          if(vm.manageNodeSummary) {
            vm.scope.getTopProcessByHostname(vm.manageNodeSummary[0]);
          } else {
            vm.selectedManageNode = null;
            vm.topProcessCpuList = null;
            vm.topProcessMemoryList = null;
          }
        },
        function(reason) {
          vm.loadingSummary = true;
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    })();

    var oldObj = {};
    vm.scope.getTopProcessByHostname = function(obj) {
      vm.scope.loading = true;
      oldObj['select'] = '';
      obj['select'] = 'active';
      vm.selectedManageNode = $sce.trustAsHtml('[ <span style="color:#77ae33;">' + obj.hostname + '</span> ]');
      vm.scope.getTopProcessCpuList(obj.hostname);
      vm.scope.getTopProcessMemoryList(obj.hostname);
      oldObj = obj;
    };
    /***** cpu used by process *****/
    vm.scope.getTopProcessCpuList = function(hostname) {
      var condition = {'hostname': hostname};
      manageNodeService.manageTopProcessByCpu(condition).then(
        function(result) {
          vm.topProcessCpuList = result.data;
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };
    /***** memory used by process *****/
    vm.scope.getTopProcessMemoryList = function(hostname) {
      var condition = {'hostname': hostname};
      manageNodeService.manageTopProcessByMemory(condition).then(
        function(result) {
          vm.topProcessMemoryList = result.data;
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };

    (vm.getMRabbitMqSummary = function() {
      manageNodeService.manageRabbitMqSummary().then(
        function(result) {
          vm.messageQueue = result.data;
          var rsc = result.data.NodeResources;
          vm.filePerc = Math.round((rsc.fileDescriptorUsed / rsc.fileDescriptorTotal) * 100);
          vm.socketPerc = Math.round((rsc.socketsUsed / rsc.socketsLimit) * 100);
          vm.erlangPerc = Math.round((rsc.erlangProcessUsed / rsc.erlangProcessLimit) * 100);
          vm.memoryPerc = Math.round((rsc.memoryMbUsed / rsc.memoryMbLimit) * 100);
          vm.diskPerc = Math.round((rsc.diskMbLimitFree / rsc.diskMbFree) * 100);
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    })();

    vm.messageQueueHost = 'http://115.68.151.175:15672';
    vm.goMessageQueue = function(uri) {
      var messageQueueUrl = vm.messageQueueHost + '/#/' + uri;
      $window.open(messageQueueUrl, 'messageQueue');
    };

    /********** reload **********/
    vm.scope.$on('broadcast:reload', function() {
      vm.scope.loading = true;
      vm.getManageNodeSummary();
      vm.getMRabbitMqSummary();
    });
  }


  angular
    .module('monitoring')
    .controller('ComputeNodeController', ComputeNodeController);

  /** @ngInject */
  function ComputeNodeController($scope, $timeout, $sce, $exceptionHandler, iaaSMainService, computeNodeService) {
    var vm = this;
    vm.scope = $scope;
    vm.Math = Math;

    vm.scope.loading = true;

    (vm.getOpenStackSummary = function() {
      iaaSMainService.openStackSummary().then(
        function(result) {
          vm.summary = result.data;
          vm.usageVcpu = Math.round((vm.summary.vcpuUsed / vm.summary.vcpuTotal) * 100);
          vm.usageMemory = Math.round((vm.summary.memoryMbUsed / vm.summary.memoryMbTotal) * 100);
          vm.usageDisk = Math.round((vm.summary.diskGbUsed / vm.summary.diskGbTotal) * 100);
          vm.usageVms = Math.round((vm.summary.vmRunning / vm.summary.vmTotal) * 100);
          vm.init(vm.usageVcpu, vm.usageMemory, vm.usageDisk, vm.usageVms);
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
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

    (vm.getComputeNodeSummary = function() {
      computeNodeService.computeNodeSummary(vm.searchCondition).then(
        function(result) {
          vm.computeNodeSummary = result.data;
          vm.scope.loading = false;
          vm.loadingSummary = false;
          if(vm.computeNodeSummary) {
            vm.scope.getTopProcessByHostname(vm.computeNodeSummary[0]);
          } else {
            vm.selectedComputeNode = null;
            vm.topProcessCpuList = null;
            vm.topProcessMemoryList = null;
          }
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    })();

    var oldObj = {};
    vm.scope.getTopProcessByHostname = function(obj) {
      vm.scope.loading = true;
      oldObj['select'] = '';
      obj['select'] = 'active';
      vm.selectedComputeNode = $sce.trustAsHtml('[ <span style="color:#77ae33;">' + obj.hostname + '</span> ]');
      vm.scope.getTopProcessCpuList(obj.hostname);
      vm.scope.getTopProcessMemoryList(obj.hostname);
      oldObj = obj;
    };

    /***** cpu used by process *****/
    vm.scope.getTopProcessCpuList = function(hostname) {
      var condition = {'hostname': hostname};
      computeNodeService.computeTopProcessByCpu(condition).then(
        function(result) {
          vm.topProcessCpuList = result.data;
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };

    /***** memory used by process *****/
    vm.scope.getTopProcessMemoryList = function(hostname) {
      var condition = {'hostname': hostname};
      computeNodeService.computeTopProcessByMemory(condition).then(
        function(result) {
          vm.topProcessMemoryList = result.data;
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
      vm.scope.loading = true;
      vm.getOpenStackSummary();
      vm.getComputeNodeSummary();
    });
  }


  angular
    .module('monitoring')
    .controller('NodeDetailController', NodeDetailController);

  /** @ngInject */
  function NodeDetailController($scope, $log, $stateParams, $interval, $timeout, $location, $window, $exceptionHandler,
                                    common, cookies, nvd3Generator, nodeChartConfig, computeNodeService, iaasLogService) {
    var vm = this;
    vm.scope = $scope;
    vm.Math = Math;

    vm.scope.loading = true;

    vm.scope.selTab = function(tab) {
      vm.scope.selectedTab = tab;
      if(tab == 'chart') {
        /*if(widgets.length > 0) {
          $scope.dashboard.widgets = widgets;
          $timeout(function() {
            $scope.$broadcast('resize');
          }, 200);
        }*/
      } else if(tab == 'logs') {
        if(vm.scope.recentLogs == null) {
          vm.scope.defaultRecentLog();
        }
      }
    };

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

    var charts = nodeChartConfig;

    // 조회조건 설정
    var hostname = $stateParams.hostname;
    vm.scope.name = hostname;

    var condition = {
      hostname: hostname,
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
          computeNodeService[opt.func](condition).then(
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

      if(cookies.getDefaultTimeRange() == '7d' || cookies.getDefaultTimeRange() == '30d') {
        widget.chart.options.chart.xAxis.tickFormat = function (d) { return d3.time.format('%y-%m-%d %H:%M:%S')(new Date(d*1000)); };
      } else if (cookies.getDefaultTimeRange() == 'custom' && (cookies.getGroupBy() == '1h36m' || cookies.getGroupBy() == '11h12m' || cookies.getGroupBy() == '48h')) {
        widget.chart.options.chart.xAxis.tickFormat = function (d) { return d3.time.format('%y-%m-%d %H:%M:%S')(new Date(d*1000)); };
      } else {
        widget.chart.options.chart.xAxis.tickFormat = function (d) { return d3.time.format('%H:%M:%S')(new Date(d*1000)); };
      }

      vm.scope.dashboard.widgets.push(widget);
    };

    // 조회주기 및 GroupBy 설정
    var savedCustom = false;
    vm.scope.saveTimeRange = function () {
      var condition = {
        hostname: hostname,
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
      cookies.putGroupBy(vm.scope.selGroupBy);

      angular.forEach(vm.scope.dashboard.widgets, function(widget, index) {
        (function(opt, idx) {
          computeNodeService[opt.func](condition).then(
            function (result) {
              if(result) {
                if(cookies.getDefaultTimeRange() == '7d' || cookies.getDefaultTimeRange() == '30d') {
                  vm.scope.dashboard.widgets[idx].chart.options.chart.xAxis.tickFormat = function (d) { return d3.time.format('%y-%m-%d %H:%M:%S')(new Date(d*1000)); };
                } else if (cookies.getDefaultTimeRange() == 'custom' && (cookies.getGroupBy() == '1h36m' || cookies.getGroupBy() == '11h12m' || cookies.getGroupBy() == '48h')) {
                  vm.scope.dashboard.widgets[idx].chart.options.chart.xAxis.tickFormat = function (d) { return d3.time.format('%y-%m-%d %H:%M:%S')(new Date(d*1000)); };
                } else {
                  vm.scope.dashboard.widgets[idx].chart.options.chart.xAxis.tickFormat = function (d) { return d3.time.format('%H:%M:%S')(new Date(d*1000)); };
                }

                var jsonArr = result.data;
                var value, arr = [];
                for(var i in jsonArr) {
                  value = jsonArr[i].metric == null ? [{time: 0, usage: 0}] : jsonArr[i].metric;
                  arr.push({values: value, key: jsonArr[i].name});
                }

                vm.scope.dashboard.widgets[idx].chart.data = arr;
                vm.scope.dashboard.widgets[idx].chart.api.refresh();
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


    /********** Logs **********/
    vm.scope.targetDateOptions = {format: 'YYYY.MM.DD'};
    vm.scope.timeOptions = {format: 'LT'};

    // pagination
    vm.scope.totalItems = 0;
    vm.scope.currentPage = 1;
    vm.scope.maxSize = 10;
    vm.scope.pageItems = 50;

    vm.scope.hostname = hostname;

    // 최근 로그 조회
    vm.scope.defaultRecentLog = function(message, optParam) {
      vm.scope.loading = true;
      var param = {
        'hostname': vm.scope.hostname,
        'pageItems': vm.scope.pageItems,
        'pageIndex': vm.scope.currentPage,
        'logType': 'log',
        'period': '5m'
      };
      if(message) param['keyword'] = message;
      if(optParam) {
        for(var key in optParam){
          if(optParam.hasOwnProperty(key)){
            param[key]=optParam[key];
          }
        }
      }
      iaasLogService.dtvDefaultRecentLog(param).then(
        function(result) {
          vm.scope.recentLogs = result.data.messages;

          vm.scope.optLogstashIndex = result.data.logstashIndex;
          vm.scope.startTime = result.data.startTime;
          vm.scope.endTime = result.data.endTime;
          vm.scope.targetDate = common.getDate(result.data.targetDate);
          vm.scope.totalItems = result.data.totalCount;
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };

    // 특정 날짜와 시간에 발생한 로그 조회
    vm.scope.specificTimeRangeLog = function(message, optParam) {
      var param = {
        'hostname': vm.scope.hostname,
        'pageItems': vm.scope.pageItems,
        'pageIndex': vm.scope.currentPage,
        'logType': 'log',
        'period': '5m'
      };
      if(message) param['keyword'] = message;
      if(optParam) {
        for(var key in optParam){
          if(optParam.hasOwnProperty(key)){
            param[key]=optParam[key];
          }
        }
      }

      iaasLogService.dtvSpecificTimeRangeLog(param).then(
        function(result) {
          vm.scope.recentLogs = result.data.messages;

          vm.scope.optLogstashIndex = result.data.logstashIndex;
          vm.scope.startTime = result.data.startTime;
          vm.scope.endTime = result.data.endTime;
          vm.scope.totalItems = result.data.totalCount;
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );

    };

    var getType = 'default';
    var message = '';
    vm.scope.logSearch = function() {
      vm.scope.loading = true;
      vm.scope.currentPage = 1;
      message = vm.scope.message;
      if(vm.scope.targetDate) {
        getType = 'specific';
        vm.scope.startTime == undefined ? vm.scope.optStartTime = '' : vm.scope.optStartTime = common.getTime(vm.scope.startTime);
        vm.scope.endTime == undefined ? vm.scope.optEndTime = '' : vm.scope.optEndTime = common.getTime(vm.scope.endTime);
        var optParam = {
          'targetDate': common.getDate(vm.scope.targetDate),
          'startTime': vm.scope.optStartTime,
          'endTime': vm.scope.optEndTime
        };
        vm.scope.specificTimeRangeLog(message, optParam);
      } else {
        getType = 'default';
        vm.scope.defaultRecentLog(message);
      }
    };
    vm.scope.pageChanged = function() {
      vm.scope.loading = true;

      var optParam = {
        'logstashIndex': vm.scope.optLogstashIndex,
        'startTime': vm.scope.optStartTime,
        'endTime': vm.scope.optEndTime
      };
      if(getType == 'specific') {
        optParam['targetDate'] = common.getDate(vm.scope.targetDate);
        vm.scope.specificTimeRangeLog(message, optParam);
      } else {
        vm.scope.defaultRecentLog(message, optParam);
      }

      $window.scrollTo(0, 0);
    };


    // Search Time Reset
    vm.scope.searchReset = function() {
      vm.scope.targetDate = '';
      vm.scope.startTime = '';
      vm.scope.endTime = '';
    }


    // Overview Page
    vm.scope.goMovePage = function() {
      var url;
      if($location.path().indexOf("/iaas/manage") > -1) {
        url = "/iaas/manage_node";
      } else {
        url = "/iaas/compute_node";
      }
      $location.path(url);
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

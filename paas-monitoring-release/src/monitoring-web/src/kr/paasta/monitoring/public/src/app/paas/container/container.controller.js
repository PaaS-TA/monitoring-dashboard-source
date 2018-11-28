(function() {
  'use strict';

  angular
    .module('monitoring')
    .controller('PaasContainerController', PaasContainerController);

  /** @ngInject */
  function PaasContainerController($scope, $timeout, $sce, $location, $interval, $exceptionHandler, common, paasContainerService) {
    var vm = this;
    vm.scope = $scope;
    vm.Math = Math;
    vm.scope.loading = true;
    vm.common = common;

    // angular-gridster option
    vm.scope.cellOverviewOptions = {
      margins: [0, 15],
      columns: 5,
      rowHeight: 110,
      widget_base_dimensions: [140, 140],
      mobileModeEnabled: false,
      draggable: false,
      resizable: {enabled: false}
    };

    vm.scope.cellOverviewItems = [
      {id: 'running',  col: 0, row: 0, sizeY: 1, sizeX: 1, name: "Running",  count: 'cellRunning',  color: "#00aacc"},
      {id: 'failed',   col: 1, row: 0, sizeY: 1, sizeX: 1, name: "Failed",   count: 'cellFailed',   color: "#e66b6b"},
      {id: 'critical', col: 2, row: 0, sizeY: 1, sizeX: 1, name: "Critical", count: 'cellCritical', color: "#ad6de8"},
      {id: 'warning',  col: 3, row: 0, sizeY: 1, sizeX: 1, name: "Warning",  count: 'cellWarning',  color: "#f0a141"},
      {id: 'total',    col: 4, row: 0, sizeY: 1, sizeX: 1, name: "Total",    count: 'cellTotal',    color: "#909090"}
    ];

    vm.scope.containerOverviewOptions = {
      margins: [0, 15],
      columns: 4,
      rowHeight: 110,
      widget_base_dimensions: [140, 140],
      mobileModeEnabled: false,
      draggable: false,
      resizable: {enabled: false}
    };

    vm.scope.containerOverviewItems = [
      {id: 'running',  col: 0, row: 0, sizeY: 1, sizeX: 1, name: "Running",  count: 'containerRunning',  color: "#00aacc"},
      {id: 'critical', col: 1, row: 0, sizeY: 1, sizeX: 1, name: "Critical", count: 'containerCritical', color: "#ad6de8"},
      {id: 'warning',  col: 2, row: 0, sizeY: 1, sizeX: 1, name: "Warning",  count: 'containerWarning',  color: "#f0a141"},
      {id: 'total',    col: 3, row: 0, sizeY: 1, sizeX: 1, name: "Total",    count: 'containerTotal',    color: "#909090"}
    ];


    // Cell Overview
    (vm.getCellOverview = function() {
      paasContainerService.cellOverview().then(
        function(result) {
          vm.scope.cellRunning = result.data.running;
          vm.scope.cellFailed = result.data.failed;
          vm.scope.cellCritical = result.data.critical;
          vm.scope.cellWarning = result.data.warning;
          vm.scope.cellTotal = result.data.total;
          vm.scope.loadingCellOverview = false;
        },
        function(reason) {
          vm.scope.loadingCellOverview = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    })();


    // Container Overview
    (vm.getContainerOverview = function() {
      paasContainerService.containerOverview().then(
        function(result) {
          vm.scope.containerRunning = result.data.running;
          vm.scope.containerCritical = result.data.critical;
          vm.scope.containerWarning = result.data.warning;
          vm.scope.containerTotal = result.data.total;
          vm.scope.loadingContainerOverview = false;
        },
        function(reason) {
          vm.scope.loadingContainerOverview = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    })();


    // Container Summary
    vm.scope.currentPage = 1;             // 현재페이지
    vm.scope.itemsPerPage = 10;           // 페이지당 목록 건수
    vm.scope.selSearchCondition = 'name';

    (vm.getContainerSummary = function(summaryObj, flag) {
      if(summaryObj == undefined || summaryObj == null) {
        vm.scope.currentPage = 1;
        summaryObj = {"currentPage":vm.scope.currentPage, "itemsPerPage":vm.scope.itemsPerPage};
      }

      var params = {
        'pageItems': vm.scope.itemsPerPage,
        'pageIndex': summaryObj.currentPage
      };

      if(vm.scope.searchKeyword) {
        if(flag == 'name') {
          params['zoneName'] = vm.scope.searchKeyword;
        }
      }

      paasContainerService.containerSummary(params).then(
        function(result) {
          vm.containerSummary = result.data.data;

          vm.scope.totalItems = result.data.totalCount;
          vm.scope.pageItems = result.data.pageItems;
          vm.scope.totalPages = Math.ceil(result.data.totalCount / result.data.pageItems);

          if(vm.containerSummary) {
            vm.scope.containerRunning = result.data.overview.running;
            vm.scope.containerCritical = result.data.overview.critical;
            vm.scope.containerWarning = result.data.overview.warning;
            vm.scope.containerTotal = result.data.overview.total;

            vm.getCellContainerRelationInfo(vm.containerSummary[0]);
          } else {
            vm.selectedPaasta = null;
            vm.topProcessMemoryList = null;
            vm.scope.loading = false;
          }

          vm.scope.loadingContainerSummary = false;
        },
        function(reason) {
          vm.scope.loadingContainerSummary = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    })();

    // Page Change
    vm.scope.pageChanged = function() {
      var summaryObj = {"currentPage":vm.scope.currentPage};
      vm.getCellContainerRelationInfo(summaryObj);
    };


    // Cell & Container Relationship
    var oldObj = {};
    vm.scope.doing_async = true;

    vm.getCellContainerRelationInfo = function(obj) {
      oldObj['select'] = '';
      obj['select'] = 'active';
      vm.selectedContainer = $sce.trustAsHtml('[ <span style="color:#77ae33;">' + obj.zoneName + '</span> ]');

      paasContainerService.containerRelationship(obj.zoneName).then(
        function(result) {
          vm.cellSummary = result.data;
          vm.relationshipList = result.data;
          vm.scope.loadingRelationship = false;
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loadingRelationship = false;
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
      oldObj = obj;
    };


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


    // Cell Overview List Pop
    vm.getCellOverviewList = function(item, $event) {
      if(item.id == 'running' || item.id == 'total') {
        $event.stopPropagation();
        return;
      }

      vm.cellOverviewPopTitle = item.name;
      // vm.cellOverviewListPop = false;
      vm.cellOverviewNullPop = false;
      vm.scope.loadingStatusListPop = true;

      paasContainerService.cellOverviewList(item.id).then(
        function(result) {
          vm.cellOverviewListPop = result.data;
          if(result.data == null) {
            vm.cellOverviewNullPop = 'No Data Available.';
          }
          vm.scope.loadingOverviewListPop = false;
        },
        function(reason) {
          vm.scope.loadingOverviewListPop = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };

    // Go Cell Detail View
    vm.goCellDetailView = function(item) {
      angular.element('#cellOverviewList').modal('hide');
      var close = $interval(function() {
        if(item != null && item.status != 'fail') {
          if(angular.element('#cellOverviewList').hasClass('in') == false) {
            $interval.cancel(close);
            $location.path('/paas/paasta/' + item.cellId);
          }
        }
      }, 300);
    };


    // Container Overview List Pop
    vm.getContainerOverviewList = function(status, title, $event) {
      if(status == 'running' || status == 'total') {
        $event.stopPropagation();
        return;
      }

      vm.containerOverviewPopTitle = title;
      vm.containerOverviewListPop = false;
      vm.containerOverviewNullPop = false;
      vm.scope.loadingStatusListPop = true;

      paasContainerService.containerOverviewList(status).then(
        function(result) {
          vm.containerOverviewListPop = result.data;
          if(result.data == null) {
            vm.containerOverviewNullPop = 'No Data Available.';
          }
          vm.scope.loadingOverviewListPop = false;
        },
        function(reason) {
          vm.scope.loadingOverviewListPop = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };

    // Go Container Detail View
    vm.goContainerDetailView = function(item) {
      angular.element('#containerOverviewList').modal('hide');
      var close = $interval(function() {
        if(item != null && item.status != 'fail') {
          if(angular.element('#containerOverviewList').hasClass('in') == false) {
            $interval.cancel(close);
            $location.path('/paas/contatiner/' + item.containerId).search('addr', item.ip);
          }
        }
      }, 300);
    };



    // Table Value Color Change
    var cpu = {};
    var memory = {};
    var disk = {};
    vm.scope.cpuUsageStyle = function(value) {
      var style = {};
      if(cpu.critical <= value) {
        style = {color: 'red'};
      } else if(cpu.warning <= value && value < cpu.critical) {
        style = {color: 'orange'};
      }
      return style;
    };
    vm.scope.memoryUsageStyle = function(value) {
      var style = {};
      if(memory.critical <= value) {
        style = {color: 'red'};
      } else if(memory.warning <= value && value < memory.critical) {
        style = {color: 'orange'};
      }
      return style;
    };
    vm.scope.diskUsageStyle = function(value) {
      var style = {};
      if(disk.critical <= value) {
        style = {color: 'red'};
      } else if(disk.warning <= value && value < disk.critical) {
        style = {color: 'orange'};
      }
      return style;
    };
    vm.scope.textStatusStyle = function(value) {
      var style = {};
      if('fail' == value) {
        style = {color: '#de4d58'};
      } else if('warning' == value) {
        style = {color: '#fbae42'};
      } else if('critical' == value) {
        style = {color: '#a76dd8'};
      }
      return style;
    };


    vm.scope.errStateStyle = function(value) {
      var style = {};
      if('fail' == value) {
        style = {color: '#de4d58'};
      } else if('warning' == value) {
        style = {color: '#fbae42'};
      } else if('critical' == value) {
        style = {color: '#a76dd8'};
      } else {
        style = {color: '#8a8c99'};
      }
      return style;
    };


    // Reload
    vm.scope.$on('broadcast:reload', function() {
      vm.scope.selSearchCondition = 'name';
      vm.scope.searchKeyword = '';

      vm.scope.loadingCellOverview = true;
      vm.scope.loadingContainerOverview = true;
      vm.scope.loadingContainerSummary = true;
      vm.scope.loadingCellSummary = true;
      vm.scope.loadingRelationship = true;
      vm.scope.loading = true;

      vm.getCellOverview();
      vm.getContainerOverview();
      vm.getContainerSummary(null, '');
    });

  }


  angular
    .module('monitoring')
    .controller('PaasContainerDetailController', PaasContainerDetailController);

  /** @ngInject */
  function PaasContainerDetailController($scope, $log, $stateParams, $interval, $timeout, $location, $exceptionHandler
                                      , paasContainerDetailService, containerChartConfig
                                      , common, cookies, nvd3Generator) {
    var vm = this;
    vm.scope = $scope;
    vm.Math = Math;

    vm.scope.loading = true;

    vm.scope.name = $location.search().name;
    vm.scope.appIndex = $location.search().appIndex;

    // Top Tab Change
    vm.scope.selTab = function(tab) {
      vm.scope.selectedTab = tab;
      if (tab == 'chart') {
        /*if(widgets.length > 0) {
          $scope.dashboard.widgets = widgets;
          $timeout(function() {
            $scope.$broadcast('resize');
          }, 200);
        }*/
      } else if (tab == 'logs') {
        if (vm.scope.recentLogs == null) {
          vm.scope.defaultRecentLog();
        }
      }
    };

    // Gridster (Layout)
    vm.scope.gridsterOpts = {
      margins: [20, 20],
      columns: 3,
      rowHeight: 250,
      mobileModeEnabled: false,
      swapping: true,
      draggable: {
        handle: 'h4',
        stop: function (event, $element, widget) {
          $timeout(function () {
            $log.log(event + ',' + $element + ',' + widget);
          }, 400)
        }
      },
      resizable: {
        enabled: true,
        handles: ['n', 'e', 's', 'w', 'ne', 'se', 'sw', 'nw'],
        minWidth: 200,
        layoutChanged: function () {
        },

        // optional callback fired when resize is started
        start: function (event, $element, widget) {
          $timeout(function () {
            $log.log(event + ',' + $element + ',' + widget);
          }, 400)
        },

        // optional callback fired when item is resized,
        resize: function (event, $element, widget) {
          if (widget.chart.api) widget.chart.api.update();
        },

        // optional callback fired when item is finished resizing
        stop: function (event, $element, widget) {
          $timeout(function () {
            if (widget.chart.api) widget.chart.api.update();
          }, 400)
        }
      }
    };

    vm.scope.dashboard = {
      widgets: []
    };

    vm.scope.events = {
      resize: function (e, scope) {
        $timeout(function () {
          if (scope.api && scope.api.update) scope.api.update();
        }, 200)
      }
    };

    vm.scope.config = {visible: false};

    $timeout(function () {
      vm.scope.config.visible = true;
    }, 200);

    angular.element(window).on('resize', function () {
      vm.scope.$broadcast('resize');
    });


    // Chart
    var charts = containerChartConfig;

    // Search Condition Setting
    var id = $stateParams.id;
    var condition = {
      id: id,
      groupBy: cookies.getGroupBy()
    };
    if (cookies.getDefaultTimeRange() == 'custom') {
      condition['timeRangeFrom'] = common.timeDifference(new Date().getTime(), cookies.getTimeRangeFrom());
      condition['timeRangeTo'] = common.timeDifference(new Date().getTime(), cookies.getTimeRangeTo());
    } else {
      condition['defaultTimeRange'] = cookies.getDefaultTimeRange();
    }

    // Service Calling & Chart Setting
    var count = 0;
    for (var i in charts) {
      var chartOpt = charts[i];
      if (chartOpt.func) {
        (function (opt, cnt) {
          paasContainerDetailService[opt.func](condition).then(
            function (result) {
              vm.scope.setWidget(cnt, opt, result.data);
            },
            function (reason) {
              vm.scope.setWidget(cnt, opt);
              $log.error(reason);
              $timeout(function () {
                $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message});
              }, 500);
            }
          );
        })(chartOpt, count);
        count++;
      }
    }
    var stop = $interval(function () {
      if (count == (vm.scope.dashboard.widgets).length) {
        $interval.cancel(stop);
        vm.scope.dashboard.widgets.sort(common.CompareForSort);
        vm.scope.loading = false;
      }
    }, 500);

    // Line Chart Create
    vm.scope.setWidget = function (index, dtvOpt, jsonArr) {
      var col = dtvOpt.col == undefined ? Math.floor(index % 3) : dtvOpt.col;
      var row = dtvOpt.row == undefined ? Math.floor(index / 3) : dtvOpt.row;
      var sizeX = dtvOpt.sizeX == undefined ? 1 : dtvOpt.sizeX;
      var sizeY = dtvOpt.sizeY == undefined ? 1 : dtvOpt.sizeY;

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
        widget.chart.options.chart.yAxis.tickFormat = function (d) {
          return d3.format('.0%')(d / 100);
        };
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
        for (var i in jsonArr) {
          value = jsonArr[i].metric == null ? [{time: 0, usage: 0}] : jsonArr[i].metric;
          arr.push({values: value, key: jsonArr[i].name});
          for (var j in value) {
            if (value[j] != null) {
              if (checkData < value[j].usage) {
                checkData = value[j].usage;
              }
            }
          }
        }
        widget.chart.data = arr;

        if (dtvOpt.type != 'list') {
          if (checkData < 5 && widget.chart.options.chart.forceY == undefined) {
            widget.chart.options.chart.forceY = [0, 5];
          }
          if (checkData > 10000) {
            widget.chart.options.chart.yAxis.axisLabelDistance = 20;
            widget.chart.options.chart.margin.left = 80;
          }
        }
      } else {
        widget.chart.options.chart.forceY = false;
      }

      if (cookies.getDefaultTimeRange() == '7d' || cookies.getDefaultTimeRange() == '30d') {
        widget.chart.options.chart.xAxis.tickFormat = function (d) {
          return d3.time.format('%y-%m-%d %H:%M:%S')(new Date(d * 1000));
        };
      } else if (cookies.getDefaultTimeRange() == 'custom' && (cookies.getGroupBy() == '11h12m' || cookies.getGroupBy() == '48h')) {
        widget.chart.options.chart.xAxis.tickFormat = function (d) {
          return d3.time.format('%y-%m-%d %H:%M:%S')(new Date(d * 1000));
        };
      } else {
        widget.chart.options.chart.xAxis.tickFormat = function (d) {
          return d3.time.format('%H:%M:%S')(new Date(d * 1000));
        };
      }

      vm.scope.dashboard.widgets.push(widget);
    };

    // 조회주기 및 GroupBy 설정
    var savedCustom = false;
    vm.scope.saveTimeRange = function () {
      var condition = {
        id: id,
        groupBy: cookies.getGroupBy()
      };

      if (vm.scope.selTimeRange == 'custom') {
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

      angular.forEach(vm.scope.dashboard.widgets, function (widget, index) {
        (function (opt, idx) {
          paasContainerDetailService[opt.func](condition).then(
            function (result) {
              if (result) {
                if (cookies.getDefaultTimeRange() == '7d' || cookies.getDefaultTimeRange() == '30d') {
                  vm.scope.dashboard.widgets[idx].chart.options.chart.xAxis.tickFormat = function (d) {
                    return d3.time.format('%y-%m-%d %H:%M:%S')(new Date(d * 1000));
                  };
                } else if (cookies.getDefaultTimeRange() == 'custom' && (cookies.getGroupBy() == '11h12m' || cookies.getGroupBy() == '48h')) {
                  vm.scope.dashboard.widgets[idx].chart.options.chart.xAxis.tickFormat = function (d) {
                    return d3.time.format('%y-%m-%d %H:%M:%S')(new Date(d * 1000));
                  };
                } else {
                  vm.scope.dashboard.widgets[idx].chart.options.chart.xAxis.tickFormat = function (d) {
                    return d3.time.format('%H:%M:%S')(new Date(d * 1000));
                  };
                }

                var jsonArr = result.data;
                var value, arr = [];
                for (var i in jsonArr) {
                  value = jsonArr[i].metric == null ? [{time: 0, usage: 0}] : jsonArr[i].metric;
                  arr.push({values: value, key: jsonArr[i].name});
                }

                vm.scope.dashboard.widgets[idx].chart.data = arr;
                vm.scope.dashboard.widgets[idx].chart.api.refresh();
                vm.scope.dashboard.widgets[idx].loading = false;
              }
            },
            function (reason, status) {
              $timeout(function () {
                $exceptionHandler(reason.Message, {code: status, message: reason.Message});
              }, 500);
            }
          );
        })(widget, index);
      });
    };

    // time range 선택 시 그에 해당하는 group by 선택
    vm.scope.selectGroupBy = function () {
      vm.scope.selGroupBy = common.getGroupingByTimeRange(vm.scope.selTimeRange, vm.scope.timeRangeFrom, vm.scope.timeRangeTo);
    };

    // 팝업에서 save 하지 않은 경우 원래 값을 유지
    angular.element('#timeRange').on('hidden.bs.modal', function () {
      if (savedCustom != true) {
        vm.scope.selTimeRange = cookies.getDefaultTimeRange();
        vm.scope.selGroupBy = cookies.getGroupBy();
      }
    });


    // Overview Page
    vm.scope.goMovePage = function() {
      $location.path("/paas/container");
    };


    // Reload
    vm.scope.$on('broadcast:reload', function () {
      if (vm.scope.selectedTab == 'logs') {
        vm.scope.logSearch();
      } else {
        angular.forEach(vm.scope.dashboard.widgets, function (widget) {
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


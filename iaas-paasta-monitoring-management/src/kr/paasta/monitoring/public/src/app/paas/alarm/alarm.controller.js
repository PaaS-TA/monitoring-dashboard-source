(function() {
  'use strict';

  angular
    .module('monitoring')
    .controller('PaasAlarmPolicyController', PaasAlarmPolicyController);

  /** @ngInject */
  function PaasAlarmPolicyController($scope, $timeout, $document, $window, $exceptionHandler, paasAlarmPolicyService) {
    var vm = this;
    vm.scope = $scope;

    vm.scope.loading = true;

    (vm.getAlarmPolicyList = function() {
      paasAlarmPolicyService.alarmPolicyList().then(
        function(result) {

          var setup = {
            pas: {
              cpu: {threshold: {critical: 0, warning: 0}, repeatTime: 0},
              memory: {threshold: {critical: 0, warning: 0}, repeatTime: 0},
              disk: {threshold: {critical: 0, warning: 0}, repeatTime: 0},
              measureTime : 0,
              email: {emailName: '', emailAddr: '', emailType: '', emailSendYn: true}
            },
            bos: {
              cpu: {threshold: {critical: 0, warning: 0}, repeatTime: 0},
              memory: {threshold: {critical: 0, warning: 0}, repeatTime: 0},
              disk: {threshold: {critical: 0, warning: 0}, repeatTime: 0},
              email: {emailName: '', emailAddr: '', emailType: '', emailSendYn: true}
            },
            con: {
              cpu: {threshold: {critical: 0, warning: 0}, repeatTime: 0},
              memory: {threshold: {critical: 0, warning: 0}, repeatTime: 0},
              disk: {threshold: {critical: 0, warning: 0}, repeatTime: 0},
              email: {emailName: '', emailAddr: '', emailType: '', emailSendYn: true}
            }
          };

          for(var i in result.data) {
            setup[result.data[i].originType][result.data[i].alarmType] = {
              threshold: {
                critical: result.data[i].criticalThreshold,
                warning: result.data[i].warningThreshold
              },
              repeatTime: result.data[i].repeatTime
            };

            // 측정시간
            setup[result.data[i].originType].measureTime = result.data[i].measureTime;

            if(result.data[i].mailAddress != null && result.data[i].mailAddress != "") {
              if(result.data[i].mailAddress.indexOf("@") > -1) {
                var array = result.data[i].mailAddress.split('@');
                var vBool = true;
                var vType = array[1];
                for(var j = 0; j < vm.mailList.length; j++) {
                  if(vm.mailList[j].code == vType) { vBool = false; }
                }
                if(vBool) { vType = ''; }
                setup[result.data[i].originType].email = {
                  emailName: array[0],
                  emailAddr: array[1],
                  emailType: vType
                }
              }
            }

            // 이메일 발송 여부
            if(result.data[i].mailSendYn == "Y") {
              setup[result.data[i].originType].email.emailSendYn = true;
            } else {
              setup[result.data[i].originType].email.emailSendYn = false;
            }
          }

          vm.setup = setup;
          vm.selType = vm.selType == null ? 'bos' : vm.selType;
          vm.setTypeObject();
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    })();


    // Alarm Channel List
    vm.snsSendYn = false;
    (vm.getAlarmSnsChannelList = function() {
      paasAlarmPolicyService.alarmSnsChannelList().then(
        function(result) {
          if(result.data != '' && result.data != null) {
            vm.alarmSnsChannelList = result.data;
            if(result.data[0].snsSendYn == "Y") {
              vm.snsSendYn = true;
            }
          } else {
            angular.element('#snsChannelList').empty();
            vm.alarmSnsChannelList = null;
          }
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    })();


    vm.channel = {
      type : '',
      id : '',
      token : '',
      description : ''
    }

    // Sns Channel Regist
    vm.channelRegist = function () {
      vm.scope.loading = true;

      if(!(vm.channel.type != '' && vm.channel.id && vm.channel.token)) {
        vm.scope.loading = false;
        var vMsg = '';
        if (vm.channel.type == '') {
          vMsg = 'Please Select Channel Type';
        } else if (vm.channel.id == '') {
          vMsg = 'Please Enter Channel ID';
        } else if (vm.channel.token == '') {
          vMsg = 'Please Enter Token Value ';
        }
        $exceptionHandler(vMsg, {code: status, message: vMsg});
        return false;
      }

      var snsSendYn = "Y";
      if(!vm.snsSendYn) { snsSendYn = "N"; }

      var body = {
        originType : vm.channel.type,
        snsId : vm.channel.id,
        token : vm.channel.token,
        expl : vm.channel.description,
        sendSnsYn : snsSendYn
      }

      paasAlarmPolicyService.channelRegist(body).then(
        function(result) {
          if(result.data.status == 'Created') {
            vm.channel.type = '';
            vm.channel.id = '';
            vm.channel.token = '';
            vm.channel.description = '';
          }
          vm.getAlarmSnsChannelList();
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data, {code: 500, message: reason.data}); });
        }
      );
    };

    // Sns Channel Delete
    vm.channelDelete = function(id) {
      paasAlarmPolicyService.channelDelete(id).then(
        function() {
          vm.getAlarmSnsChannelList();
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data, {code: 500, message: reason.data}); });
        }
      );
    };

    // Policy Save
    vm.saveSetup = function() {
      vm.scope.loading = true;

      if(!(vm.origin.email.emailName != '' && vm.origin.email.emailAddr != '') && vm.origin.email.emailSendYn == 'Y') {
        vm.scope.loading = false;
        var vMsg = '';
        if (vm.origin.email.emailName == '') {
          vMsg = 'Please enter Email Name';
        } else {
          vMsg = 'Please enter Email Addr or selete Email Addr ';
        }
        $exceptionHandler(vMsg, {code: status, message: vMsg});
      }

      var mailSendYn = "Y";
      if(!vm.origin.email.emailSendYn) { mailSendYn = "N"; }

      var snsSendYn = "Y";
      if(!vm.snsSendYn) { snsSendYn = "N"; }

      var body = [
        {
          originType: vm.selType,
          alarmType: "cpu",
          warningThreshold: vm.origin.cpu.threshold.warning,
          criticalThreshold: vm.origin.cpu.threshold.critical,
          repeatTime: vm.origin.cpu.repeatTime,
          measureTime: vm.origin.measureTime
        },
        {
          originType: vm.selType,
          alarmType: "memory",
          warningThreshold: vm.origin.memory.threshold.warning,
          criticalThreshold: vm.origin.memory.threshold.critical,
          repeatTime: vm.origin.memory.repeatTime,
          measureTime: vm.origin.measureTime
        },
        {
          originType: vm.selType,
          alarmType: "disk",
          warningThreshold: vm.origin.disk.threshold.warning,
          criticalThreshold: vm.origin.disk.threshold.critical,
          repeatTime: vm.origin.disk.repeatTime,
          measureTime: vm.origin.measureTime
        },
        {
          originType: vm.selType,
          mailAddress: vm.origin.email.emailName + '@' + vm.origin.email.emailAddr,
          mailSendYn: mailSendYn,
          snsSendYn: snsSendYn
        }
      ];

      paasAlarmPolicyService.updateAlarmSetup(body).then(
        function() {
          vm.getAlarmPolicyList();
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data, {code: 500, message: reason.data}); });
        }
      );
    };

    vm.mailList = [
      {name:"gmail",    code:"gmail.com",   $hashKey:"gmail.com"},
      {name:"yahoo",    code:"yahoo.com",   $hashKey:"yahoo.com"},
      {name:"hotmail",  code:"hotmail.com", $hashKey:"hotmail.com"},
      {name:"naver",    code:"naver.com",   $hashKey:"naver.com"},
      {name:"nate",     code:"nate.com",    $hashKey:"nate.com"},
      {name:"daum",     code:"hanmail.net", $hashKey:"hanmail.net"}
    ];

    vm.setEmailAddr = function() {
      if(vm.origin.email.emailType == '') {
        vm.origin.email.emailAddr = '';
      } else {
        vm.origin.email.emailAddr = vm.origin.email.emailType;
      }
    };

    vm.setTypeObject = function() {
      if(vm.selType == 'pas') {
        vm.origin = vm.setup.pas;
      } else if(vm.selType == 'bos') {
        vm.origin = vm.setup.bos;
      } else if(vm.selType == 'con') {
        vm.origin = vm.setup.con;
      }
    };

    // SNS Regist Method
    vm.snsRegistMethod = function() {
      var innerContents = $document[0].getElementById('snsRegistMethod').innerHTML;
      var popupWinindow = $window.open('', '_blank', 'width=724,height=700,scrollbars=no,menubar=no,toolbar=no,location=no,status=no,titlebar=no');
      popupWinindow.document.open();
      popupWinindow.document.write(
        '<html >' +
        '<head>' +
        '<link rel="stylesheet" type="text/css" href="bower_components/components-font-awesome/css/font-awesome.css">'+
        '<link rel="stylesheet" type="text/css" href="styles/vendor-min.css">' +
        '<link rel="stylesheet" type="text/css" href="styles/app-min.css">' +
        '<style>\n' +
        '.contents {font-size:14px;}\n' +
        '.arrow {text-align:center;margin-bottom:15px;}\n' +
        '.thumbnail {padding:15px;margin-bottom:20px;border:1px solid #ddd;border-radius:5px;}\n' +
        '.fa.fa-caret-right {margin-right:5px;}\n' +
        '</style>' +
        '</head>' +
        '<body>' + innerContents + '</body>' +
        '</html>'
      );
      popupWinindow.document.close();
    };

    /********** reload **********/
    vm.scope.$on('broadcast:reload', function() {
      vm.scope.loading = true;
      vm.getAlarmPolicyList();
    });
  }


  angular
    .module('monitoring')
    .controller('PaasAlarmStatusController', PaasAlarmStatusController);

  /** @ngInject */
  function PaasAlarmStatusController($scope, $timeout, $location, $exceptionHandler, common, paasAlarmStatusService) {
    var vm = this;
    vm.scope = $scope;
    vm.scope.loading = true;
    vm.common = common;

    vm.scope.optionsFrom = {format: 'YYYY.MM.DD'};
    vm.scope.optionsTo = {format: 'YYYY.MM.DD'};

    // pagination
    vm.scope.currentPage = 1;             // 현재페이지
    vm.scope.itemsPerPage = 10;           // 페이지당 목록 건수
    vm.scope.maxSize = 10;

    (vm.getAlarmStatusList = function(cPage) {
      vm.stateParams = $location.search();

      if(vm.scope.loading == false) {
        vm.scope.loading = true;
      }

      if(cPage == undefined || cPage == null) {
        vm.scope.currentPage = 1;
      }

      if(vm.stateParams._o) {
        vm.scope.selOriginType = vm.stateParams._o;
      }
      if(vm.stateParams._a) {
        vm.scope.selAlarmType = vm.stateParams._a;
      }
      if(vm.stateParams._l) {
        vm.scope.selLevel = vm.stateParams._l;
      }
      if(vm.stateParams._r) {
        vm.scope.selResolveStatus = vm.stateParams._r;
      }
      if(vm.stateParams._f && vm.stateParams._t) {
        vm.scope.dateFrom = moment(parseInt(vm.stateParams._f));
        vm.scope.dateTo = moment(parseInt(vm.stateParams._t));
      }
      if(vm.stateParams._p) {
        vm.scope.currentPage = vm.stateParams._p;
      }

      var param = '?pageItems=' + vm.scope.itemsPerPage;
          param += '&pageIndex=' + vm.scope.currentPage;

      if(vm.scope.selOriginType) {
        param += '&originType=' + vm.scope.selOriginType;
      }
      if(vm.scope.selAlarmType) {
        param += '&alarmType='+vm.scope.selAlarmType;
      }
      if(vm.scope.selResolveStatus) {
        param += '&resolveStatus='+vm.scope.selResolveStatus;
      }
      if(vm.scope.dateFrom) {
        param += '&searchDateFrom='+vm.scope.dateFrom.startOf('day');
      }
      if(vm.scope.dateTo) {
        param += '&searchDateTo='+vm.scope.dateTo.endOf('day');
      }
      if(vm.scope.selLevel) {
        param += '&level='+vm.scope.selLevel;
      }

      paasAlarmStatusService.alarmStatusList(param).then(
        function(result) {
          vm.issues = result.data;

          vm.scope.totalItems = result.data.totalCount;
          vm.scope.pageItems = result.data.pageItem;
          vm.scope.totalPages = Math.ceil(result.data.totalCount / result.data.pageItem);

          $location.search({});
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    })();

    vm.scope.pageChanged = function() {
      vm.getAlarmStatusList(vm.scope.currentPage);
    };


    // Go Container Detail
    vm.goAlarmStatusDetail = function(itemId) {
      var params = {
        'selOriginType' : vm.scope.selOriginType,
        'selAlarmType' : vm.scope.selAlarmType,
        'selLevel' : vm.scope.selLevel,
        'selResolveStatus' : vm.scope.selResolveStatus,
        'dateFrom' : vm.scope.dateFrom,
        'dateTo' : vm.scope.dateTo,
        'page' : vm.scope.currentPage
      }
      $location.path('/paas/alarm/status/'+itemId).search(params);
    };


    // Reload
    vm.scope.$on('broadcast:reload', function() {
      vm.scope.loading = true;
      vm.getAlarmStatusList(vm.scope.currentPage);
    });
  }


  angular
    .module('monitoring')
    .controller('PaasAlarmStatusDetailController', PaasAlarmStatusDetailController);

  /** @ngInject */
  function PaasAlarmStatusDetailController($scope, $timeout, $stateParams, $location, $exceptionHandler, common, paasAlarmStatusService) {
    var vm = this;
    vm.scope = $scope;
    vm.scope.loading = true;

    vm.alarmId = $stateParams;
    vm.stateParams = $location.search();
    vm.common = common;

    // Alarm Status Detail
    (vm.getAlarmStatusDetail = function() {
      paasAlarmStatusService.alarmStatusDetail(vm.alarmId.id).then(
        function(result) {
          vm.alarm = result.data;
          vm.alarmLevelClass = vm.alarm.level;
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    })();

    // Modal Setting
    vm.confirmReceipt = function() {
      vm.modalReset();
      vm.modalTitle = '접수';
      vm.modalMessage = '접수 하시겠습니까?';
      vm.confirmCallback = vm.alarmDetailReceipt;
    };
    vm.confirmActionModify = function(action) {
      vm.modalReset();
      vm.modalTitle = '수정';
      vm.modalMessage = '조치내용을 수정 하시겠습니까?';
      vm.actionId = action.id;
      vm.confirmCallback = vm.alarmActionUpdate;
      vm.callbackParam = action;
    };
    vm.confirmActionDelete = function(actionId) {
      vm.modalReset();
      vm.modalTitle = '삭제';
      vm.modalMessage = '조치내용을 삭제 하시겠습니까?';
      vm.actionId = actionId;
      vm.confirmCallback = vm.alarmActionDelete;
    };
    vm.modalReset = function() {
      vm.modalTitle = '';
      vm.modalMessage = '';
      vm.actionId = '';
      vm.confirmCallback = undefined;
      vm.callbackParam = undefined;
    };


    // Alarm Detail Cancel
    vm.cancel = function(path) {
      var params = {
        _o : vm.stateParams.selOriginType,
        _a : vm.stateParams.selAlarmType,
        _l : vm.stateParams.selLevel,
        _r : vm.stateParams.selResolveStatus,
        _f : vm.stateParams.dateFrom,
        _t : vm.stateParams.dateTo,
        _p : vm.stateParams.page
      }
      $location.path(path).search(params);
    };

    // Alarm Detail Receipt
    vm.alarmDetailReceipt = function() {
      vm.scope.loading = true;

      var body = {resolveStatus: "2"};

      paasAlarmStatusService.alarmStatusUpdate(vm.alarm.id, body).then(
        function() {
          vm.getAlarmStatusDetail();
          // $rootScope.$broadcast('broadcast:getAlarms');
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };

    // Alarm Action Complete
    vm.alarmActionComplete = function() {
      vm.scope.loading = true;

      if(vm.actionDesc == '' || vm.actionDesc == null || vm.actionDesc == undefined) {
        $exceptionHandler('', {code: null, message: '조치내용이 입력되지 않았습니다.'});
        return;
      }

      var body = {resolveStatus: "3"};

      paasAlarmStatusService.alarmStatusUpdate(vm.alarm.id, body).then(
        function() {
          vm.alarmActionCreate();
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };

    // Alarm Action Create
    vm.alarmActionCreate = function() {
      vm.scope.loading = true;

      if(vm.actionDesc == '' || vm.actionDesc == null || vm.actionDesc == undefined) {
        $exceptionHandler('', {code: null, message: '조치내용이 입력되지 않았습니다.'});
        return;
      }

      var body = {alarmId: vm.alarm.id, alarmActionDesc: vm.actionDesc};

      paasAlarmStatusService.alarmActionCreate(body).then(
        function() {
          vm.getAlarmStatusDetail();
          vm.actionDesc = '';
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };

    //
    vm.setUpdating = function(index) {
      vm.updating = index;
    };

    // Alarm Action Update
    vm.alarmActionUpdate = function() {
      vm.scope.loading = true;
      vm.updating = -1;

      var body = {alarmActionDesc: vm.callbackParam.alarmActionDesc};

      paasAlarmStatusService.alarmActionUpdate(vm.callbackParam.id, body).then(
        function() {
          vm.getAlarmStatusDetail();
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };

    // Alarm Action Delete
    vm.alarmActionDelete = function() {
      vm.scope.loading = true;

      paasAlarmStatusService.alarmActionDelete(vm.actionId).then(
        function() {
          vm.getAlarmStatusDetail();
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };

    // Reload
    vm.scope.$on('broadcast:reload', function() {
      vm.scope.loading = true;
      vm.getAlarmStatusDetail();
    });

  }


  angular
    .module('monitoring')
    .controller('PaasAlarmStatisticsController', PaasAlarmStatisticsController);

  /** @ngInject */
  function PaasAlarmStatisticsController($scope, $timeout, $location, $exceptionHandler, $document, $window, Excel,
                                         common, nvd3Generator, paasAlarmStatisticsService) {
    var vm = this;
    vm.scope = $scope;
    vm.scope.loading = true;

    var today = '';
    vm.interval = 1;
    vm.period = 'd';

    vm.scope.totalWidgets = {};
    vm.scope.servieWidgets = {};
    vm.scope.matrixWidgets = {};

    // Alarm Statistics
    (vm.getAlarmStatistics = function(condition) {
      var from, to;
      vm.period = condition;

      if(condition == 'd') {
        today = common.convertTimestampToDate(moment().add((Number(vm.interval)-1)*-1, 'days'));
        vm.strPeriod = today;
      } else if(condition == 'w') {
        from = common.convertTimestampToDate(moment().add((Number(vm.interval)-1)*-1, 'week').add(-1, 'week'));
        to = common.convertTimestampToDate(moment().add((Number(vm.interval)-1)*-1, 'week'));
        vm.strPeriod = from.substring(0, today.lastIndexOf(' ')) + ' ~ ' + to.substring(0, today.lastIndexOf(' '));
      } else if (condition == 'm') {
        from = common.convertTimestampToDate(moment().add((Number(vm.interval)-1)*-1, 'month').add(-1, 'month'));
        to = common.convertTimestampToDate(moment().add((Number(vm.interval)-1)*-1, 'month'));
        vm.strPeriod = from.substring(0, today.lastIndexOf(' ')) + ' ~ ' + to.substring(0, today.lastIndexOf(' '));
      } else if (condition == 'y') {
        from = common.convertTimestampToDate(moment().add((Number(vm.interval)-1)*-1, 'year').add(-1, 'year'));
        to = common.convertTimestampToDate(moment().add((Number(vm.interval)-1)*-1, 'year'));
        vm.strPeriod = from.substring(0, today.lastIndexOf(' ')) + ' ~ ' + to.substring(0, today.lastIndexOf(' '));
      }

      var params = '?period=' + vm.period;
      params += '&interval=' + vm.interval;

      // Statistics Info
      paasAlarmStatisticsService.alarmStatistics(params).then(
        function(result) {
          vm.stats = result.data;
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );


      // Total Chart
      paasAlarmStatisticsService.alarmStatisticTotal(params).then(
        function(result) {
          vm.scope.setWidget(1, "Total", "alarmStatisticTotal", result.data);
        },
        function(reason) {
          vm.scope.setWidget(1, "Total", "alarmStatisticTotal", null);
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );

      // Service Chart
      paasAlarmStatisticsService.alarmStatisticService(params).then(
        function(result) {
          vm.scope.setWidget(1, "Service", "alarmStatisticService", result.data);
        },
        function(reason) {
          vm.scope.setWidget(2, "Service", "alarmStatisticService", null);
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );

      // Matrix Chart
      paasAlarmStatisticsService.alarmStatisticMatrix(params).then(
        function(result) {
          vm.scope.setWidget(1, "Matrix", "alarmStatisticMatrix", result.data);
        },
        function(reason) {
          vm.scope.setWidget(3, "Matrix", "alarmStatisticMatrix", null);
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );

    })(vm.period);


    vm.scope.config = {visible: true};

    $timeout(function () {
      vm.scope.config.visible = true;
    }, 200);

    // Line Chart Create
    vm.scope.setWidget = function (id, title, func, jsonArr) {
      var widget = {
        col: 0, row: 0, sizeX: 1, sizeY: 1, id: id, type: "lineChart",
        chart: {
          options: nvd3Generator["lineChart"].options()
        },
        func: func
      };

      if(jsonArr) {
        var value, arr = [];
        var checkData = 0;
        for (var i in jsonArr) {
          value = jsonArr[i].stat == null ? [{time: 0, count: 0}] : jsonArr[i].stat;
          arr.push({values: value, key: jsonArr[i].name});
          for(var j in value) {
            if(value[j] != null) {
              if(checkData < value[j].count) {
                checkData = value[j].count;
              }
            }
          }
        }
        widget.chart.data = arr;
      } else {
        widget.chart.options.chart.forceY = false;
      }

      widget.chart.options.chart.xAxis.tickFormat = function(d) {
        if(vm.period == "d") {
          return d3.time.format('%H:%M')(new Date(d * 1000));
        } else if(vm.period == "w") {
          return d3.time.format('%y-%m-%d')(new Date(d * 1000));
        } else if(vm.period == "m") {
          return d3.time.format('%y-%m-%d')(new Date(d * 1000));
        } else if(vm.period == "y") {
          return d3.time.format('%y-%m')(new Date(d * 1000));
        }
      };

      widget.chart.options.chart.y = function (d){ return d.count; };
      widget.chart.options.chart.yAxis.tickFormat = function(d) { return d3.format("")(d); };
      widget.chart.options.chart.yAxis.axisLabel = "Count";
      widget.chart.options.chart.yAxis.axisLabelDistance = -5;
      widget.chart.options.chart.margin.left = 55;
      widget.axisLabel = "Count";

      if(title == "Total") {
        vm.scope.totalWidgets = widget;
      } else if(title == "Service") {
        vm.scope.servieWidgets = widget;
      } else if(title == "Matrix") {
        vm.scope.matrixWidgets = widget;
      }
    };

    // Interval Alarm Statistics
    vm.setInterval = function(arithmetic) {
      if(arithmetic == '-') {
        vm.interval--;
      } else {
        vm.interval++;
      }
      vm.getAlarmStatistics(vm.period);
    };


    // Alarm Status Page Move(Origin)
    vm.goListOriginType = function(originType, level) {
      goList(originType, '', level);
    };


    // Alarm Status Page Move(Type)
    vm.goListAlarmType = function(alarmType, level) {
      goList('', alarmType, level);
    };


    function goList(originType, alarmType, level) {
      var param = {
        _o: originType,
        _a: alarmType,
        _l: level
      };

      var from, to;
      var tmpYear, tmpMonth, tmpDay;
      if(vm.period == 'd') {
        var tmpDate = vm.strPeriod;
        from = new Date(tmpDate.slice(0,4), tmpDate.slice(5,7)-1, tmpDate.slice(8,10)).getTime();
        to = new Date(tmpDate.slice(0,4), tmpDate.slice(5,7)-1, tmpDate.slice(8,10)).getTime();
      } else if(vm.period == 'w') {
        var tmpFrom = vm.strPeriod.slice(0, vm.strPeriod.indexOf('~')-1);
        var tmpTo = vm.strPeriod.slice(vm.strPeriod.indexOf('~')+2);
        from = new Date(tmpFrom.slice(0,4), tmpFrom.slice(5,7)-1, tmpFrom.slice(8,10)).getTime();
        to = new Date(tmpTo.slice(0,4), tmpTo.slice(5,7)-1, tmpTo.slice(8,10)).getTime();
      } else if (vm.period == 'm') {
        tmpYear = vm.strPeriod.slice(0, 4);
        tmpMonth = parseInt(vm.strPeriod.slice(5), 10)-1;
        tmpDay = new Date().getDate();
        from = new Date(tmpYear, tmpMonth, tmpDay).getTime();
        to = new Date(tmpYear, tmpMonth+1, tmpDay).getTime();
      } else if(vm.period == 'y') {
        tmpYear = parseInt(vm.strPeriod.slice(0, 4))+1;
        tmpMonth = parseInt(vm.strPeriod.slice(5), 10)-1;
        tmpDay = new Date().getDate();
        from = new Date(tmpYear-1, tmpMonth, tmpDay).getTime();
        to = new Date(tmpYear, tmpMonth, tmpDay).getTime();
      }

      param['_f'] = from;
      param['_t'] = to;
      $location.path('/paas/alarm/status').search(param);
    }


    // Statistics Print
    vm.statisticsPrint = function(divName) {
      var dateContents = '<p class="date_box">' + vm.strPeriod + '</p>'
      var innerContents = $document[0].getElementById(divName).innerHTML;
      var popupWinindow = $window.open('', '_blank', 'width=1200,height=700,scrollbars=no,menubar=no,toolbar=no,location=no,status=no');
      popupWinindow.document.open();
      popupWinindow.document.write(
        '<!doctype html>' +
        '<head>' +
        '<meta charset="utf-8">' +
        '<title>IaaS / PaaS Monitoring</title>' +
        '<style>' +
        'p.date_box {width:280px;height:60px;font-size:20px;padding:16px 10px;text-align:center;border:1px solid #cccccc;margin-bottom:30px;}\n' +
        '.panel{margin-bottom:30px;}' +
        '</style>' +
        '<link rel="stylesheet" type="text/css" href="styles/vendor-min.css">' +
        '<link rel="stylesheet" type="text/css" href="styles/app-min.css">' +
        '</head>' +
        '<body onload="javascript:window.print();">'+ dateContents + innerContents +
        '</html>'
      );
      popupWinindow.document.close();
    };


    // Statistics Excel Download
    vm.excelDownload = function(divName) {
      var exportHref = Excel.tableToExcel(divName, 'Statistics');
      var link = $document[0].createElement('a');
      link.href = exportHref;
      link.download = 'Alarm_Statistics.xls';
      link.click();
    };

  }

})();

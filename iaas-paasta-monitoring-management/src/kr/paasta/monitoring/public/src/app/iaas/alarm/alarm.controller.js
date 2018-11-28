(function() {
  'use strict';

  angular
    .module('monitoring')
    .controller('IaasAlarmNotificationController', IaasAlarmNotificationController);

  /** @ngInject */
  function IaasAlarmNotificationController($scope, $timeout, $window, $sce, $exceptionHandler, iaasAlarmNotificationService) {
    var vm = this;
    vm.scope = $scope;
    vm.Math = Math;

    vm.scope.loading = true;

    vm.checkedCnt = 0;

    var limit = 10;
    var offset = 0;
    vm.alarmNotificationList = [];
    (vm.getAlarmNotificationList = function() {
      var params = {
        'offset': offset,
        'limit': limit
      };
      iaasAlarmNotificationService.alarmNotificationList(params).then(
        function(result) {
          if(result.data.data) vm.alarmNotificationList = vm.alarmNotificationList.concat(result.data.data);
          vm.totalCount = result.data.totalCnt;
          vm.moreButton = '<strong>더 보 기</strong> (총 ' + vm.totalCount + ' 건)';
          if(vm.alarmNotificationList) {
            offset = vm.alarmNotificationList.length;
          }
          if(vm.alarmNotificationList.length >= vm.totalCount) {
            vm.moreButton = '(총 ' + vm.totalCount + '건)';
          }
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    })();

    vm.checkAllNotification = function() {
      angular.forEach(vm.alarmNotificationList, function(notification) {
        if(vm.selectAll) {
          notification.select = true;
          vm.checkedCnt++;
        } else {
          notification.select = false;
          vm.checkedCnt--;
        }
      });
    };

    vm.checkNotification = function(obj) {
      if(obj.select) {
        vm.checkedCnt++;
      } else {
        vm.checkedCnt--;
      }
    };

    vm.getAlarmNotification = function(obj) {
      vm.detail = angular.copy(obj);
    };

    vm.saveAlarmNotification = function() {
      vm.scope.loading = true;

      var data = {
        name: vm.detail.name,
        address: vm.detail.email
      };
      var func = null;
      if(vm.detail.id) {
        data.id = vm.detail.id;
        data.period = vm.detail.period;
        func = iaasAlarmNotificationService.updateAlarmNotification(data);
      } else {
        func = iaasAlarmNotificationService.insertAlarmNotification(data);
      }
      func.then(
        function() {
          vm.getAlarmNotificationList();
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };

    vm.deleteAlarmNotification = function() {
      var deleteCnt = 0;
      var deletedCnt = 0;
      angular.forEach(vm.alarmNotificationList, function(notification, index) {
        if(notification.select) {
          deleteCnt++;
          iaasAlarmNotificationService.deleteAlarmNotification(notification.id).then(
            function() {
              vm.alarmNotificationList.splice(index, 1);
              vm.checkedCnt--;
              deletedCnt++;
              if(deleteCnt == deletedCnt) {
                vm.getAlarmNotificationList();
              }
            },
            function(reason) {
              vm.scope.loading = false;
              $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
            }
          );
        }
      });
    };

    /********** reload **********/
    vm.scope.$on('broadcast:reload', function() {
      vm.scope.loading = true;
      vm.checkedCnt = 0;
      limit = 10;
      offset = 0;
      vm.alarmNotificationList = [];
      vm.getAlarmNotificationList();
    });
  }


  angular
    .module('monitoring')
    .controller('IaasAlarmPolicyController', IaasAlarmPolicyController);

  /** @ngInject */
  function IaasAlarmPolicyController($scope, $timeout, $window, $exceptionHandler, iaasAlarmPolicyService) {
    var vm = this;
    vm.scope = $scope;
    vm.Math = Math;

    vm.scope.loading = true;

    vm.checkedCnt = 0;

    var limit = 10;
    var offset = 0;
    vm.alarmPolicyList = [];
    vm.searchAlarmPolicy = function() {
      offset = 0;
      vm.alarmPolicyList = [];
      vm.totalCount = 0;
      vm.getAlarmPolicyList();
    };

    (vm.getAlarmPolicyList = function() {
      var params = {
        'offset': offset,
        'limit': limit
      };
      if(vm.searchCondition) {
        params['name'] = vm.searchCondition;
      }
      if(vm.selectedSeverity) {
        params['severity'] = vm.selectedSeverity;
      }
      iaasAlarmPolicyService.alarmPolicyList(params).then(
        function(result) {
          if(result.data.data) {
            vm.alarmPolicyList = vm.alarmPolicyList.concat(result.data.data);

            if(vm.alarmPolicyList.length > 0) {
              for(var i = 0; i < vm.alarmPolicyList.length; i++) {
                if(vm.alarmPolicyList[i].severity == 'HIGH') {
                  vm.alarmPolicyList[i].severity = 'WARNING'
                }
              }
            }
          }
          vm.totalCount = result.data.totalCnt;
          vm.moreButton = '<strong>더 보 기</strong> (총 ' + vm.totalCount + ' 건)';
          if(vm.alarmPolicyList) {
            offset = vm.alarmPolicyList.length;
          }
          if(vm.alarmPolicyList.length >= vm.totalCount) {
            vm.moreButton = '(총 ' + vm.totalCount + '건)';
          }
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    })();

    vm.checkAllPolicy = function() {
      angular.forEach(vm.alarmPolicyList, function(policy) {
        if(vm.selectAll) {
          policy.select = true;
          vm.checkedCnt++;
        } else {
          policy.select = false;
          vm.checkedCnt--;
        }
      });
    };

    vm.checkPolicy = function(obj) {
      if(obj.select) {
        vm.checkedCnt++;
      } else {
        vm.checkedCnt--;
      }
    };

    vm.scope.alarmSeverityClass = function(severity) {
      var severityClass = '';
      if(severity == 'CRITICAL') {
        severityClass = 'critical';
      } else if (severity == 'WARNING') {
        severityClass = 'warning';
      }
      return severityClass;
    };

    vm.deleteAlarmPolicy = function() {
      var deleteCnt = 0;
      var deletedCnt = 0;
      angular.forEach(vm.alarmPolicyList, function(policy) {
        if(policy.select) {
          deleteCnt++;
          iaasAlarmPolicyService.deleteAlarmPolicy(policy.id).then(
            function() {
              vm.checkedCnt--;
              deletedCnt++;
              if(deleteCnt == deletedCnt) {
                $window.location.reload();
                vm.scope.loading = false;
              }
            },
            function(reason) {
              vm.scope.loading = false;
              $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
            }
          );
        }
      });
    };

    /********** reload **********/
    vm.scope.$on('broadcast:reload', function() {
      vm.scope.loading = true;
      vm.searchAlarmPolicy();
    });

  }


  angular
    .module('monitoring')
    .controller('AlarmPolicyDetailController', AlarmPolicyDetailController);

  /** @ngInject */
  function AlarmPolicyDetailController($scope, $stateParams, $timeout, $location, $exceptionHandler
                                       , iaasAlarmPolicyService, tenantService, iaasAlarmNotificationService) {
    var vm = this;
    vm.scope = $scope;
    vm.Math = Math;

    vm.scope.loading = true;

    vm.pageTitle = $stateParams.id == 'new' ? null : 'Alarm Policy Update';
    vm.expressionYn = $stateParams.id == 'new' ? false : true;
    vm.alarmPolicy = {};
    vm.alarmPolicy['severity'] = 'HIGH';
    vm.alarmPolicy['matchBy'] = 'hostname';

    /***** get alarm policy *****/
    vm.getAlarmPolicy = function() {
      iaasAlarmPolicyService.alarmPolicy($stateParams.id).then(
        function(result) {
          vm.alarmPolicy = result.data;

          // expression
          var obj_expression = vm.alarmPolicy.expression;
          var gate = obj_expression.indexOf(' and ') >= 0 ? ' and ' : ' or ';
          var gate2 = gate == ' and ' ? ' or ' : ' and ';
          var tmp_expression = obj_expression.split(gate);
          for(var i in tmp_expression) {
            if(tmp_expression[i].indexOf(gate2) >= 0) {
              var or_expression = tmp_expression[i].split(gate2);
              tmp_expression.splice(i, 1);
              for(var j in or_expression) {
                tmp_expression.push(or_expression[j]);
              }
            }
          }
          var arr_expression = [];
          angular.forEach(tmp_expression, function(expression) {
            expression = expression.replace(/(\s*)/g, '');
            var json = {};
            json['func'] = expression.substr(0, expression.indexOf('('));
            var metric = expression.slice(expression.indexOf('(')+1, expression.indexOf(')'));
            json['metric'] = metric;
            if(metric.indexOf('{') >= 0) {
              json['metric'] = metric.substr(0, metric.indexOf('{'));
              json['dimension'] = metric.slice(metric.indexOf('{')+1, metric.indexOf('}'));
            }
            var tmp = expression.substr(expression.indexOf(')')+1);
            var len = tmp.indexOf('=') < 0 ? 1 : 2;
            json['operation'] = tmp.substr(0, len);
            json['value'] = parseInt(tmp.substr(len));
            arr_expression.push(json);
          });
          vm.alarmPolicy.arrExpression = arr_expression;

          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };

    /***** dimension setup modal *****/
    vm.measurementType = '';
    vm.scope.showDimensionModal = false;
    vm.selDimension = 'hostname';
    var modalDimensionIndex = 0;
    vm.setDimensionModal = function(obj, index) {
      vm.scope.loadingModal = true;

      var modalDimension = obj.dimension;
      modalDimensionIndex = index;
      vm.scope.dimensionTitle = modalDimension == undefined ? 'Set Dimension' : modalDimension;
      if(modalDimension) {
        vm.scope.dimensionTitle = modalDimension;
        vm.modalDimension = modalDimension.substr(0, modalDimension.indexOf('='));
        vm.modalDimensionValue = modalDimension.substr(modalDimension.indexOf('=')+1);
      } else {
        vm.modalDimension = 'hostname';
        vm.scope.dimensionTitle = 'Set Dimension';
      }
      vm.scope.showDimensionModal = !vm.scope.showDimensionModal;

      switch (obj.metric){
        case 'cpu.percent':
        case 'mem.usable_perc':
        case 'disk.space_used_perc':
          vm.measurementType = 'node';
          break;
        case 'vm.cpu.utilization_norm_perc':
        case 'vm.mem.free_perc':
          vm.measurementType = 'vm';
          break;
        default:
          vm.measurementType = 'node';
          break;
      }

      if(vm.measurementType == 'node') {
        iaasAlarmPolicyService.nodeList().then(
          function(result) {
            vm.nodeList = result.data;
            vm.selDimensionValue1 = vm.nodeList[0];
            vm.scope.loadingModal = false;
          },
          function(reason) {
            vm.scope.loadingModal = false;
            $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
          }
        );
      } else {
        getTenantList();
      }
    };
    vm.setDimension = function() {
      var dimensionStr = vm.measurementType == 'node' ? 'hostname=' + vm.selDimensionValue1 : 'hostname=' + vm.selDimensionValue2.name;
      vm.alarmPolicy.arrExpression[modalDimensionIndex].dimension = dimensionStr;
    };

    /***** alarm receiver setup modal *****/
    vm.scope.showAlarmActionModal = false;
    vm.scope.alarmActionTitle = 'Alarm Receiver';
    vm.setAlarmActionModal = function() {
      vm.scope.loadingModal = true;
      vm.scope.showAlarmActionModal = !vm.scope.showAlarmActionModal;
      rLimit = 10;
      rOffset = 0;
      vm.rAlarmNotificationList = [];
      vm.rGetAlarmNotificationList();
    };
    // Show selected rows
    vm.selectedAlarmAction = [];
    vm.selectAlarmAction = function(obj) {
      if(vm.selectedAlarmAction.indexOf(obj) >= 0) {
        obj['select'] = '';
        vm.selectedAlarmAction.splice(vm.selectedAlarmAction.indexOf(obj), 1);
      } else {
        obj['select'] = 'active';
        vm.selectedAlarmAction.push(obj);
      }
    };
    vm.setAlarmAction = function() {
      if(vm.selectedAlarmAction) {
        if(vm.alarmPolicy.alarmAction) {
          vm.alarmPolicy.alarmAction = vm.alarmPolicy.alarmAction.concat(vm.selectedAlarmAction);
        } else {
          vm.alarmPolicy['alarmAction'] = [];
          vm.alarmPolicy.alarmAction = vm.selectedAlarmAction;
        }
        vm.selectedAlarmAction = [];
      }
    };
    vm.deleteAlarmAction = function(index) {
      vm.alarmPolicy.alarmAction.splice(index, 1);
    };

    /***** init *****/
    vm.addAlarmPolicy = function() {
      vm.alarmPolicy.arrExpression.push(
        {
          func: 'max',
          metric: 'cpu.percent',
          operation: '>'
        }
      );
      vm.scope.loading = false;
    };
    if($stateParams.id == 'new') {
      vm.alarmPolicy['arrExpression'] = [];
      vm.addAlarmPolicy();
    } else {
      vm.getAlarmPolicy();
    }

    /***** save alarm policy *****/
    vm.saveAlarmPolicy = function() {
      vm.scope.loadingModal = true;

      // make expression
      var expression = '';
      var arrExpression = vm.alarmPolicy.arrExpression;
      for(var i in arrExpression) {
        expression += arrExpression[i].func + '(' + arrExpression[i].metric;
        if(arrExpression[i].dimension) {
          expression += '{' + arrExpression[i].dimension + '}';
        }
        expression += ') ' + arrExpression[i].operation + ' ' + arrExpression[i].value;
        if(arrExpression[i].value == null) {
          alert('Expression의 Value 값을 입력해 주세요.');
          angular.element('#expressionValue').focus();
          vm.scope.loadingModal = false;
          return false;
        }
        if((i+1) < arrExpression.length) {
          expression += ' ' + arrExpression[i].gate + ' ';
        }
      }

      var alarm_actions = [];
      angular.forEach(vm.alarmPolicy.alarmAction, function(alarmAction) {
        alarm_actions.push(alarmAction.id);
      });

      if(alarm_actions.length == 0) {
        alert('Alarm Receiver를 등록해 주세요.');
        angular.element('#btnPlus').focus();
        vm.scope.loadingModal = false;
        return false;
      }

      var body = {
        name: vm.alarmPolicy.name,
        severity: vm.alarmPolicy.severity,
        expression: expression,
        alarm_actions: alarm_actions,
        description: vm.alarmPolicy.description
      };

      var func = '';
      if($stateParams.id == 'new') {
        body['match_by'] = [vm.alarmPolicy.matchBy];
        func = 'insertAlarmPolicy';
      } else {
        body['id'] = $stateParams.id;
        body['match_by'] = vm.alarmPolicy.matchBy;
        func = 'updateAlarmPolicy';
      }
      iaasAlarmPolicyService[func](body).then(
        function() {
          vm.scope.loadingModal = false;
          $location.path('/iaas/alarm/policy');
        },
        function(reason) {
          vm.scope.loadingModal = false;
          $timeout(function() { $exceptionHandler(reason.data, {code: reason.status, message: reason.data.message}); }, 500);
        }
      );
      vm.scope.loadingModal = false;
    };

    /***** get project & instance list => for dimension setup modal *****/
    vm.getTenantList = function() {
      tenantService.tenantSummary().then(
        function(result) {
          vm.tenantList = result.data;
          vm.selDimensionValue1 = vm.tenantList[0];
          getInstanceList(vm.selDimensionValue1.id);
        },
        function(reason) {
          vm.scope.loadingModal = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };
    vm.getInstanceList = function(id) {
      vm.scope.loadingModal = true;
      var params = {
        'hostname': '',
        'limit': 100,
        'marker': ''
      };
      tenantService.tenantInstanceList(id, params).then(
        function(result) {
          vm.instanceList = result.data.metric;
          vm.selDimensionValue2 = vm.instanceList[0];
          vm.scope.loadingModal = false;
        },
        function(reason) {
          vm.scope.loadingModal = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };
    function getTenantList() {
      vm.getTenantList();
    }
    function getInstanceList(id) {
      vm.getInstanceList(id);
    }

    /***** get project & instance list => for alarm receiver setup modal *****/
    var rLimit = 10;
    var rOffset = 0;
    vm.rAlarmNotificationList = [];
    vm.rGetAlarmNotificationList = function() {
      var params = {
        'offset': rOffset,
        'limit': rLimit
      };
      iaasAlarmNotificationService.alarmNotificationList(params).then(
        function(result) {
          if(result.data.data) vm.rAlarmNotificationList = vm.rAlarmNotificationList.concat(result.data.data);
          vm.rTotalCount = result.data.totalCnt;
          vm.rMoreButton = '<strong>더 보 기</strong> (총 ' + vm.rTotalCount + ' 건)';
          if(vm.rAlarmNotificationList) {
            rOffset = vm.rAlarmNotificationList.length;
            // Disable already existing action from list
            if(vm.alarmPolicy.alarmAction) {
              angular.forEach(vm.alarmPolicy.alarmAction, function(alarmAction) {
                angular.forEach(vm.rAlarmNotificationList, function(notification) {
                  if(alarmAction.id == notification.id) {
                    notification['disabled'] = true;
                  }
                });
              });
            }
          }
          if(vm.rAlarmNotificationList.length >= vm.rTotalCount) {
            vm.rMoreButton = '(총 ' + vm.rTotalCount + '건)';
          }
          vm.scope.loadingModal = false;
        },
        function(reason) {
          vm.scope.loadingModal = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };

  }


  angular
    .module('monitoring')
    .controller('IaasAlarmStatusController', IaasAlarmStatusController);

  /** @ngInject */
  function IaasAlarmStatusController($scope, $timeout, $location, $exceptionHandler, iaasAlarmStatusService) {
    var vm = this;
    vm.scope = $scope;
    vm.Math = Math;

    vm.scope.loading = true;

    if($location.search().severity != 'undefined' && $location.search().severity != '') {
      vm.selectedSeverity = $location.search().severity;
    }

    var limit = 10;
    var offset = 0;
    vm.alarmStatusList = [];
    vm.searchAlarmStatus = function() {
      offset = 0;
      vm.alarmStatusList = [];
      vm.totalCount = 0;
      vm.getAlarmStatusList();
    };
    vm.selectedState = 'ALARM';
    (vm.getAlarmStatusList = function() {
      var params = {
        'offset': offset,
        'limit': limit,
        'state': vm.selectedState
      };
      if(vm.selectedSeverity) {
        params['severity'] = vm.selectedSeverity;

      }
      iaasAlarmStatusService.alarmStatusList(params).then(
        function(result) {
          if(result.data.data) {
            vm.alarmStatusList = vm.alarmStatusList.concat(result.data.data);

            if(vm.alarmStatusList.length > 0) {
              for(var i = 0; i < vm.alarmStatusList.length; i++) {
                if(vm.alarmStatusList[i].severity == 'HIGH') {
                  vm.alarmStatusList[i].severity = 'WARNING'
                }
              }
            }
          }

          vm.totalCount = result.data.totalCnt;
          vm.moreButton = '<strong>더 보 기</strong> (총 ' + vm.totalCount + ' 건)';
          if(vm.alarmStatusList) {
            offset = vm.alarmStatusList.length;
          }
          if(vm.alarmStatusList.length >= vm.totalCount) {
            vm.moreButton = '(총 ' + vm.totalCount + '건)';
          }

          $location.search('severity', null);
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    })();


    vm.scope.alarmSeverityClass = function(severity) {
      var severityClass = '';
      if(severity == 'CRITICAL') {
        severityClass = 'critical';
      } else if (severity == 'WARNING') {
        severityClass = 'warning';
      }
      return severityClass;
    };


    /********** reload **********/
    vm.scope.$on('broadcast:reload', function() {
      vm.scope.loading = true;
      vm.searchAlarmStatus();
    });
  }


  angular
    .module('monitoring')
    .controller('AlarmStatusDetailController', AlarmStatusDetailController);

  /** @ngInject */
  function AlarmStatusDetailController($scope, $stateParams, $window, $timeout, $exceptionHandler, iaasAlarmStatusService) {
    var vm = this;
    vm.scope = $scope;
    vm.Math = Math;

    vm.scope.loading = true;

    var alarmId = $stateParams.id;
    (vm.getAlarmStatus = function() {
      iaasAlarmStatusService.alarmStatus(alarmId).then(
        function(result) {
          vm.detail = result.data;
          if(vm.detail.severity != '' && vm.detail.severity == 'HIGH') {
            vm.detail.severity = 'WARNING';
          }
          vm.getAlarmStatusHistory('1d');
          vm.getAlarmActionList();
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    })();

    // set panel height of alarm history
    angular.element('.alarmHistory').height(angular.element('.alarmDetail').height());
    angular.element($window).on('resize', function () {
      angular.element('.alarmHistory').height(angular.element('.alarmDetail').height());
    });

    vm.getAlarmStatusHistory = function(timeRange) {
      vm.timeRange = timeRange;
      var params = {
        timeRange: timeRange
      };
      iaasAlarmStatusService.alarmStatusHistoryList(alarmId, params).then(
        function(result) {
          vm.alarmStatusHistoryList = result.data;
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };

    // 조치 이력 조회
    vm.getAlarmActionList = function() {
      iaasAlarmStatusService.alarmActionList(alarmId).then(
        function(result) {
          vm.alarmActionList = result.data;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };
    // 조치 이력 등록
    vm.insertAlarmAction = function() {
      vm.scope.loading = true;
      if(vm.alarmActionDesc == '' || vm.alarmActionDesc == null || vm.alarmActionDesc == undefined) {
        vm.scope.loading = false;
        $exceptionHandler('', {code: null, message: '조치내용이 입력되지 않았습니다.'});
        return;
      }
      var body = {
        alarmId: alarmId,
        alarmActionDesc: vm.alarmActionDesc
      };
      iaasAlarmStatusService.insertAlarmAction(body).then(
        function() {
          vm.getAlarmActionList();
          vm.alarmActionDesc = '';
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };

    // 조치내용 수정 필드 활성화
    vm.setModifying = function(index) {
      vm.modifying = index;
    };
    // 조치내용 수정
    vm.updateAction = function(obj) {
      vm.scope.loading = true;
      vm.modifying = -1;
      var body = {
        alarmActionDesc: obj.alarmActionDesc
      };
      iaasAlarmStatusService.updateAlarmAction(obj.id, body).then(
        function() {
          vm.getAlarmActionList();
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };

    // 조치내용 삭제
    vm.deleteAction = function(id) {
      vm.scope.loading = true;
      iaasAlarmStatusService.deleteAlarmAction(id).then(
        function() {
          vm.getAlarmActionList();
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };
  }

})();

'use strict';

angular.module('monitor.controllers')

    /******************** 알람 현황 ********************/
    .controller('almLstCtrl', function($scope, $stateParams, $timeout, $exceptionHandler, common, almSrvc){
        $scope.common = common;

        $scope.optionsFrom = {format: 'YYYY.MM.DD'};
        $scope.optionsTo = {format: 'YYYY.MM.DD'};

        // pagination
        $scope.totalItems = 0;
        $scope.currentPage = 1; // 현재 페이지
        $scope.maxSize = 10; // 페이지 당 목록 건수
        $scope.pageItems = 10;

        ($scope.getAlarms = function(page) {
            $scope.loading = true;

            if (page == undefined || page == null) {
                $scope.currentPage = 1;
                page = {"currentPage":$scope.currentPage, "maxSize":$scope.maxSize};
            }

            if($stateParams._o) {
                $scope.selOriginType = $stateParams._o;
            }
            if($stateParams._a) {
                $scope.selAlarmType = $stateParams._a;
            }
            if($stateParams._f && $stateParams._t) {
                $scope.dateFrom = moment(parseInt($stateParams._f));
                $scope.dateTo = moment(parseInt($stateParams._t));
            }
            if($stateParams._l) {
                $scope.selLevel = $stateParams._l;
            }

            var param = '?pageItems='+$scope.maxSize;
            param += '&pageIndex='+page.currentPage;
            if($scope.selOriginType) {
                param += '&originType='+$scope.selOriginType;
            }
            if($scope.selAlarmType) {
                param += '&alarmType='+$scope.selAlarmType;
            }
            if($scope.selResolveStatus) {
                param += '&resolveStatus='+$scope.selResolveStatus;
            }
            if($scope.dateFrom) {
                param += '&searchDateFrom='+$scope.dateFrom.startOf('day');
            }
            if($scope.dateTo) {
                param += '&searchDateTo='+$scope.dateTo.endOf('day');
            }
            if($scope.selLevel) {
                param += '&level='+$scope.selLevel;
            }

            var alarmPromise = almSrvc.alarmList(param);
            alarmPromise.success(function (result, status, headers) {
                $scope.issues = result.data;
                $scope.totalItems = result.totalCount; // pagination
                $scope.totalPages = Math.ceil(result.totalCount / result.pageItem); // pagination
                $scope.loading = false;
            });
            alarmPromise.error(function (reason, status, headers) {
                $scope.loading = false;
                $timeout(function() { $exceptionHandler(reason.Message, {code: status, message: reason.Message}); }, 500);
            });
        })();

        $scope.pageChanged = function() {
            var obj = {"currentPage":$scope.currentPage, "maxSize":$scope.maxSize};
            $scope.getAlarms(obj);
        };
    })

    /******************** 알람 현황 상세(조치이력 관리) ********************/
    .controller('almLstDtlCtrl', function($scope, $rootScope, $stateParams, $timeout, $exceptionHandler, $location, common, almSrvc){
        $scope.common = common;

        ($scope.getAlarmDetail = function() {
            $scope.loading = true;
            var alarmPromise = almSrvc.getAlarmDetail($stateParams.id);
            alarmPromise.success(function (result, status, headers) {
                $scope.alarm = result;
                $scope.alarmLevelClass = $scope.alarm.level;
                $scope.loading = false;
            });
            alarmPromise.error(function (reason, status, headers) {
                $scope.loading = false;
                $timeout(function() { $exceptionHandler(reason.Message, {code: status, message: reason.Message}); }, 500);
            });
            angular.element('body').find('.modal-backdrop').hide();
        })();

        $scope.confirmAccept = function() {
            $scope.cancel();
            $scope.modalTitle = '접수';
            $scope.modalMessage = '접수 하시겠습니까?';
            $scope.confirmCallback = $scope.accept;

        };
        $scope.confirmActionModify = function(action) {
            $scope.cancel();
            $scope.modalTitle = '수정';
            $scope.modalMessage = '조치내용을 수정 하시겠습니까?';
            $scope.actionId = action.id;
            $scope.confirmCallback = $scope.updateAction;
            $scope.callbackParam = action;
        };
        $scope.confirmActionDelete = function(actionId) {
            $scope.cancel();
            $scope.modalTitle = '삭제';
            $scope.modalMessage = '조치내용을 삭제 하시겠습니까?';
            $scope.actionId = actionId;
            $scope.confirmCallback = $scope.deleteAction;
        };
        $scope.cancel = function() {
            $scope.modalTitle = '';
            $scope.modalMessage = '';
            $scope.actionId = '';
            $scope.confirmCallback = undefined;
            $scope.callbackParam = undefined;
        };
        $scope.go = function(path) {
            $location.path(path);
        };

        // 접수
        $scope.accept = function() {
            $scope.loading = true;
            var body = {resolveStatus: "2"};
            var acceptPromise = almSrvc.updateAlarmDetail($scope.alarm.id, body);
            acceptPromise.success(function (result, status, headers) {
                $scope.getAlarmDetail();
                $rootScope.$broadcast('broadcast:getAlarms');
            });
            acceptPromise.error(function (reason, status, headers) {
                $scope.loading = false;
                $timeout(function() { $exceptionHandler(reason.Message, {code: status, message: reason.Message}); }, 500);
            });
        };

        // 조치 완료
        $scope.resolve = function() {
            $scope.loading = true;
            var body = {resolveStatus: "3"};
            var acceptPromise = almSrvc.updateAlarmDetail($scope.alarm.id, body);
            acceptPromise.success(function (result, status, headers) {
                $scope.createAction();
            });
            acceptPromise.error(function (reason, status, headers) {
                $scope.loading = false;
                $timeout(function() { $exceptionHandler(reason.Message, {code: status, message: reason.Message}); }, 500);
            });
        };

        // 조치내용 입력
        $scope.createAction = function() {
            $scope.loading = true;
            if($scope.actionDesc == '' || $scope.actionDesc == null || $scope.actionDesc == undefined) {
                $exceptionHandler('', {code: null, message: '조치내용이 입력되지 않았습니다.'});
                return;
            }
            var body = {alarmId: $scope.alarm.id, alarmActionDesc: $scope.actionDesc};
            var actionPromise = almSrvc.createAlarmAction(body);
            actionPromise.success(function (result, status, headers) {
                $scope.getAlarmDetail();
                $scope.actionDesc = '';
            });
            actionPromise.error(function (reason, status, headers) {
                $scope.loading = false;
                $timeout(function() { $exceptionHandler(reason.Message, {code: status, message: reason.Message}); }, 500);
            });
        };

        // 조치내용 수정 필드 활성화
        $scope.setModifying = function(index) {
            var regUser = $scope.alarm.data[index].regUser;
            $scope.modifying = index;
        };
        // 조치내용 수정
        $scope.updateAction = function() {
            $scope.loading = true;
            $scope.modifying = -1;
            var body = {alarmActionDesc: $scope.callbackParam.alarmActionDesc};
            var actionPromise = almSrvc.updateAlarmAction($scope.callbackParam.id, body);
            actionPromise.success(function (result, status, headers) {
                $scope.getAlarmDetail();
            });
            actionPromise.error(function (reason, status, headers) {
                $scope.loading = false;
                $timeout(function() { $exceptionHandler(reason.Message, {code: status, message: reason.Message}); }, 500);
            });
        };

        // 조치내용 삭제
        $scope.deleteAction = function() {
            $scope.loading = true;
            var actionPromise = almSrvc.deleteAlarmAction($scope.actionId);
            actionPromise.success(function (result, status, headers) {
                $scope.getAlarmDetail();
            });
            actionPromise.error(function (reason, status, headers) {
                $scope.loading = false;
                $timeout(function() { $exceptionHandler(reason.Message, {code: status, message: reason.Message}); }, 500);
            });
        };
    })

    /******************** 알람 통계 ********************/
    .controller('almSttCtrl', function($scope, $timeout, $location, $exceptionHandler, common, almSrvc){

        // 달력
        $scope.optionsFrom = {format: 'YYYY.MM.DD'};
        $scope.optionsTo = {format: 'YYYY.MM.DD'};
        $scope.updateTimeRange = function (dateFrom, dateTo) {
            if(dateFrom!= null && dateTo == null) {
                $scope.dateTo = (dateFrom+1000);
            }
            if(dateFrom) {
                $scope.optionsTo.minDate = dateFrom;
                $scope.optionsToDate = $scope.optionsTo.minDate._d;
            }
            if(dateTo) {
                $scope.optionsFrom.maxDate = dateTo;
                $scope.optionsFromDate = $scope.optionsFrom.maxDate._d;
            }
        };
        $timeout(function() {
            $scope.updateTimeRange($scope.dateFrom, $scope.dateTo);
        });
        
        $scope.interval = 1;
        $scope.period = 'd';
        var today = '';
        ($scope.stat = function(condition) {
            $scope.loading = true;
            $scope.period = condition;
            if(condition == 'd') {
                today = common.convertTimestampToDate(moment().add((Number($scope.interval)-1)*-1, 'days'));
                $scope.strPeriod = today;
            } else if(condition == 'w') {
                var from = common.convertTimestampToDate(moment().add((Number($scope.interval)-1)*-1, 'week').add(-1, 'week'));
                var to = common.convertTimestampToDate(moment().add((Number($scope.interval)-1)*-1, 'week'));
                $scope.strPeriod = from + ' ~ ' + to;
            } else if (condition == 'm') {
                today = common.convertTimestampToDate(moment().add((Number($scope.interval)-1)*-1, 'month'));
                $scope.strPeriod = today.substring(0, today.lastIndexOf('.'));
            }
            var param = '?period=' + $scope.period;
            if(condition == 'custom') {
                param += '&searchDateFrom=' + $scope.dateFrom.startOf('day') + '&searchDateTo=' + $scope.dateTo.endOf('day');
            } else {
                param += '&interval=' + $scope.interval;
            }
            var statPromise = almSrvc.alarmStat(param);
            statPromise.success(function (result, status, headers) {
                $scope.stats = result;
                $scope.loading = false;
            });
            statPromise.error(function (reason, status, headers) {
                $scope.loading = false;
                $timeout(function() { $exceptionHandler(reason.Message, {code: status, message: reason.Message}); }, 500);
            });
        })($scope.period);

        $scope.setInterval = function(arithmetic) {
            if(arithmetic == '-') {
                $scope.interval--;
            } else {
                $scope.interval++;
            }
            $scope.stat($scope.period);
        };

        $scope.goListOriginType = function(originType, level) {
            goList(originType, '', level);
        };
        $scope.goListAlarmType = function(alarmType, level) {
            goList('', alarmType, level);
        };
        function goList(originType, alarmType, level) {
            var param = {
                _o: originType,
                _a: alarmType,
                _l: level
            };
            var from, to;
            if($scope.period == 'custom') {
                from = $scope.dateFrom;
                to = $scope.dateTo;
            } else if($scope.period == 'd') {
                var tmpDate = $scope.strPeriod;
                from = new Date(tmpDate.slice(0,4), tmpDate.slice(5,7)-1, tmpDate.slice(8,10)).getTime();
                to = new Date(tmpDate.slice(0,4), tmpDate.slice(5,7)-1, tmpDate.slice(8,10)).getTime();
            } else if($scope.period == 'w') {
                var tmpFrom = $scope.strPeriod.slice(0, $scope.strPeriod.indexOf('~')-1);
                var tmpTo = $scope.strPeriod.slice($scope.strPeriod.indexOf('~')+2);
                from = new Date(tmpFrom.slice(0,4), tmpFrom.slice(5,7)-1, tmpFrom.slice(8,10)).getTime();
                to = new Date(tmpTo.slice(0,4), tmpTo.slice(5,7)-1, tmpTo.slice(8,10)).getTime();
            } else if ($scope.period == 'm') {
                var tmpYear = $scope.strPeriod.slice(0, 4);
                var tmpMonth = parseInt($scope.strPeriod.slice(5), 10)-1;
                var tmpDay = new Date().getDate();
                from = new Date(tmpYear, tmpMonth-1, tmpDay).getTime();
                to = new Date(tmpYear, tmpMonth, tmpDay).getTime();
            }
            param['_f'] = from;
            param['_t'] = to;
            $location.path('/almLst').search(param);
        }
    })

    /******************** 알람 설정 ********************/
    .controller('almSetCtrl', function($scope, $timeout, $exceptionHandler, common, almSrvc){
        ($scope.getAlarmSetups = function() {
            $scope.loading = true;
            var containerPromise = almSrvc.getAlarmSetup();
            containerPromise.success(function (result) {
                var setup = {
                    pas: {
                        cpu: {threshold: {critical: 0, warning: 0}, repeatTime: 0},
                        memory: {threshold: {critical: 0, warning: 0}, repeatTime: 0},
                        disk: {threshold: {critical: 0, warning: 0}, repeatTime: 0}
                    },
                    bos: {
                        cpu: {threshold: {critical: 0, warning: 0}, repeatTime: 0},
                        memory: {threshold: {critical: 0, warning: 0}, repeatTime: 0},
                        disk: {threshold: {critical: 0, warning: 0}, repeatTime: 0}
                    },
                    con: {
                        cpu: {threshold: {critical: 0, warning: 0}, repeatTime: 0},
                        memory: {threshold: {critical: 0, warning: 0}, repeatTime: 0},
                        disk: {threshold: {critical: 0, warning: 0}, repeatTime: 0}
                    }
                };
                for(var i in result) {
                    setup[result[i].originType][result[i].alarmType] = {
                        threshold: {
                            critical: result[i].criticalThreshold,
                            warning: result[i].warningThreshold
                        },
                        repeatTime: result[i].repeatTime
                    };
                }
                $scope.setup = setup;
                $scope.selType = $scope.selType == null ? 'pas' : $scope.selType;
                $scope.setTypeObject();
                $scope.loading = false;
            });
            containerPromise.error(function (reason, status, headers) {
                $scope.loading = false;
                $timeout(function() { $exceptionHandler(reason.Message, {code: status, message: reason.Message}); }, 500);
            });
        })();

        $scope.setTypeObject = function() {
            if($scope.selType == 'pas') {
                $scope.origin = $scope.setup.pas;
            } else if($scope.selType == 'bos') {
                $scope.origin = $scope.setup.bos;
            } else if($scope.selType == 'con') {
                $scope.origin = $scope.setup.con;
            }
        };

        $scope.saveSetup = function() {
            $scope.loading = true;
            var body = [
                {
                    originType: $scope.selType,
                    alarmType: "cpu",
                    warningThreshold: $scope.origin.cpu.threshold.warning,
                    criticalThreshold: $scope.origin.cpu.threshold.critical,
                    repeatTime: $scope.origin.cpu.repeatTime
                },
                {
                    originType: $scope.selType,
                    alarmType: "memory",
                    warningThreshold: $scope.origin.memory.threshold.warning,
                    criticalThreshold: $scope.origin.memory.threshold.critical,
                    repeatTime: $scope.origin.memory.repeatTime
                },
                {
                    originType: $scope.selType,
                    alarmType: "disk",
                    warningThreshold: $scope.origin.disk.threshold.warning,
                    criticalThreshold: $scope.origin.disk.threshold.critical,
                    repeatTime: $scope.origin.disk.repeatTime
                }
            ];
            var containerPromise = almSrvc.updateAlarmSetup(body);
            containerPromise.success(function (result) {
                $scope.getAlarmSetups();
            });
            containerPromise.error(function (reason, status, headers) {
                $scope.loading = false;
                $timeout(function() { $exceptionHandler(reason.Message, {code: status, message: reason.Message}); }, 500);
            });
        };
    })

    /******************** 컨테이너 배치현황 ********************/
    .controller('ctnPstCtrl', function($scope, $timeout, $exceptionHandler, almSrvc){
        $scope.loading = true;

        var param = '';
        var containerPromise = almSrvc.containerDeploy(param);
        containerPromise.success(function (result, status, headers) {
            $scope.cells = result;
            $scope.loading = false;
        });
        containerPromise.error(function (reason, status, headers) {
            $scope.loading = false;
            $timeout(function() { $exceptionHandler(reason.Message, {code: status, message: reason.Message}); }, 500);
        });
    })
;

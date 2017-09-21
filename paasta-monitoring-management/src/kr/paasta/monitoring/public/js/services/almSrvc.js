'use strict';

angular.module('monitor.services')
    .factory('almSrvc', ['$http', 'common', 'cache', 'CONSTANTS', function ($http, common, cache, CONSTANTS) {
        var almSrvc = {};
        var serverUrl = CONSTANTS.apiServer + CONSTANTS.context;

        /********** 알람현황 **********/
        almSrvc.alarmList = function(param) {
            return common.retrieveResource(common.resourcePromise(serverUrl+'alarms'+param, 'GET'));
        };

        almSrvc.alarmResolveStatus = function(resolveStatus) {
            return common.retrieveResource(common.resourcePromise(serverUrl+'alarms/status/'+resolveStatus, 'GET'));
        };

        almSrvc.getAlarmDetail = function(id) {
            return common.retrieveResource(common.resourcePromise(serverUrl+'alarms/'+id, 'GET'));
        };

        almSrvc.updateAlarmDetail = function(id, body) {
            return common.retrieveResource(common.resourcePromise(serverUrl+'alarms/'+id, 'PUT', body));
        };

        almSrvc.createAlarmAction = function(body) {
            return common.retrieveResource(common.resourcePromise(serverUrl+'alarmsAction', 'POST', body));
        };

        almSrvc.updateAlarmAction = function(actionId, body) {
            return common.retrieveResource(common.resourcePromise(serverUrl+'alarmsAction/'+actionId, 'PUT', body));
        };

        almSrvc.deleteAlarmAction = function(actionId) {
            return common.retrieveResource(common.resourcePromise(serverUrl+'alarmsAction/'+actionId, 'DELETE'));
        };

        /********** 알람통계 **********/
        almSrvc.alarmStat = function(param) {
            return common.retrieveResource(common.resourcePromise(serverUrl+'alarmsStat'+param, 'GET'));
        };

        /********** 알람설정 **********/
        almSrvc.getAlarmSetup = function() {
            return common.retrieveResource(common.resourcePromise(serverUrl+'alarmsPolicy', 'GET'));
        };
        almSrvc.updateAlarmSetup = function(body) {
            return common.retrieveResource(common.resourcePromise(serverUrl+'alarmsPolicy', 'PUT', body));
        };

        /********** 컨테이너 배치현황 **********/
        almSrvc.containerDeploy = function(param) {
            return common.retrieveResource(common.resourcePromise(serverUrl+'containerDeploy'+param, 'GET'));
        };

        return almSrvc;
    }])
;

'use strict';

angular.module('monitor.services', ['LocalStorageModule'])
    .factory('common', function ($http, cache) {
        var common = {};

        common.resourcePromise = function (endpoint, method, body) {
            var config = {
                method: method,
                url: endpoint,
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json'}
            };
            if (body !== null) {
                config.data = JSON.stringify(body);
            }
            var promise = $http(config);
            return promise;
        };

        common.retrieveResource = function (promise, cacheFn) {
            promise.success = function (fn) {
                promise.then(function (response) {
                    if(response == undefined) response = {data:null,status:null,headers:null};
                    if (cacheFn) {
                        cacheFn(response.data);
                    }
                    fn(response.data, response.status, response.headers);
                });
                return promise;
            };
            promise.error = function (fn) {
                promise.then(null, function (response) {
                    fn(response.data, response.status, response.headers);
                });
                return promise;
            };
            return promise;
        };

        common.getDateTime = function(str) {
            // yyyy-mm-dd hh:MM:ss => 1970-01-01 00:00:00
            var date = new Date(str);
            var year = date.getFullYear();
            var month = ("0"+(date.getMonth()+1)).substr(-2);
            var day = ("0"+date.getDate()).substr(-2);
            var hour = ("0"+date.getHours()).substr(-2);
            var minutes = ("0"+date.getMinutes()).substr(-2);
            var seconds = ("0"+date.getSeconds()).substr(-2);
            return year+'-'+month+'-'+day+' '+hour+':'+minutes+':'+seconds;
        };common.getDateTTime = function(str) {
            // yyyy-mm-dd'T'hh:MM:ss => 1970-01-01T00:00:00
            var date = new Date(str);
            var year = date.getFullYear();
            var month = ("0"+(date.getMonth()+1)).substr(-2);
            var day = ("0"+date.getDate()).substr(-2);
            var hour = ("0"+date.getHours()).substr(-2);
            var minutes = ("0"+date.getMinutes()).substr(-2);
            var seconds = ("0"+date.getSeconds()).substr(-2);
            return year+'-'+month+'-'+day+'T'+hour+':'+minutes+':'+seconds;
        };
        common.getDate = function(str) {
            // yyyy-mm-dd => 1970-01-01
            var date = new Date(str);
            var year = date.getFullYear();
            var month = ("0"+(date.getMonth()+1)).substr(-2);
            var day = ("0"+date.getDate()).substr(-2);
            return year+'-'+month+'-'+day;
        };
        common.getTime = function(str) {
            // hh:MM:ss => 00:00:00
            var date = new Date(str);
            var hour = ("0"+date.getHours()).substr(-2);
            var minutes = ("0"+date.getMinutes()).substr(-2);
            var seconds = ("0"+date.getSeconds()).substr(-2);
            return hour+':'+minutes+':'+seconds;
        };

        // InfluxDB 조회 쿼리 용 시간
        common.timeDifference = function(timestamp1, timestamp2) {
            var result = '';
            var difference = timestamp1 - timestamp2;

            var daysDifference = Math.floor(difference/1000/60/60/24);
            difference -= daysDifference*1000*60*60*24;
            if(daysDifference > 0) result += daysDifference + 'd';

            var hoursDifference = Math.floor(difference/1000/60/60);
            difference -= hoursDifference*1000*60*60;
            if(hoursDifference > 0) result += hoursDifference + 'h';

            var minutesDifference = Math.floor(difference/1000/60);
            difference -= minutesDifference*1000*60;
            if(minutesDifference > 0) result += minutesDifference + 'm';

            var secondsDifference = Math.floor(difference/1000);
            if(secondsDifference > 0) result += secondsDifference + 's';

            return result;
        };

        common.getGroupingByTimeRange = function(timeRange, from, to) {
            var grouping = '';
            switch(timeRange) {
                case '15m':
                    grouping = '1m';
                    break;
                case '30m':
                    grouping = '2m';
                    break;
                case '1h':
                    grouping = '4m';
                    break;
                case '3h':
                    grouping = '12m';
                    break;
                case '6h':
                    grouping = '24m';
                    break;
                case '12h':
                    grouping = '48m';
                    break;
                case '1d':
                    grouping = '1h36m';
                    break;
                case '7d':
                    grouping = '11h12m';
                    break;
                case '1M':
                    grouping = '48h';
                    break;
                case 'custom':
                    grouping = common.selectGroupingByCustomTimeRange(from, to);
                    break;
                default:
                    grouping = '1m';
            }
            return grouping;
        };

        common.selectGroupingByCustomTimeRange = function(from, to) {
            var subtraction = Math.round((from - to) / 120);
            var grouping = '';
            if(subtraction <= 18000) {
                grouping = '1m';
            } else if (18000 < subtraction && subtraction <= 36000) {
                grouping = '2m';
            } else if (36000 < subtraction && subtraction <= 108000) {
                grouping = '4m';
            } else if (108000 < subtraction && subtraction <= 216000) {
                grouping = '12m';
            } else if (216000 < subtraction && subtraction <= 432000) {
                grouping = '24m';
            } else if (432000 < subtraction && subtraction <= 864000) {
                grouping = '48m';
            } else if (864000 < subtraction && subtraction <= 6048000) {
                grouping = '1h36m';
            } else if (6048000 < subtraction && subtraction <= 26784000) {
                grouping = '11h12m';
            } else if (26784000 < subtraction) {
                grouping = '48h';
            }
            return grouping;
        };

        common.getOriginTypeString = function(type) {
            var strType = '';
            switch(type) {
                case 'bos' :
                    strType = 'Bosh';
                    break;
                case 'pas' :
                    strType = 'PaaS-TA';
                    break;
                case 'con' :
                    strType = 'Container';
                    break;
            }
            return strType;
        };

        common.getResourceTypeString = function(type) {
            var strType = '';
            switch(type) {
                case 'cpu' :
                    strType = 'CPU';
                    break;
                case 'memory' :
                    strType = 'Memory';
                    break;
                case 'disk' :
                    strType = 'Disk';
                    break;
            }
            return strType;
        };
        
        common.getMillisecondsRefreshTime = function(refreshTime) {
            var interval = 0;
            var num = refreshTime.substring(0, refreshTime.length-1);
            var unit = refreshTime.substring(refreshTime.length-1);
            switch(unit) {
                case 'm' :
                    interval = num * 1000 * 60;
                    break;
                case 'h' :
                    interval = num * 1000 * 60 * 60;
                    break;
                case 'd' :
                    interval = num * 1000 * 60 * 60 * 24;
                    break;
            }
            return interval;
        };

        common.convertTimestampToDateTime = function(sourceTimestamp) {
            var source = moment(Number(sourceTimestamp)).unix() * 1000;

            var toMonth = (new Date(source).getMonth()+1).toString().length === 1 ? '0'+(new Date(source).getMonth()+1).toString() : (new Date(source).getMonth()+1).toString();
            var toDate = new Date(source).getDate().toString().length === 1 ? '0'+new Date(source).getDate().toString() : new Date(source).getDate().toString();
            var toHours = new Date(source).getHours().toString().length === 1 ? '0'+new Date(source).getHours().toString() : new Date(source).getHours().toString();
            var toMinutes = new Date(source).getMinutes().toString().length === 1 ? '0'+new Date(source).getMinutes().toString() : new Date(source).getMinutes().toString();
            var toSeconds = new Date(source).getSeconds().toString().length === 1 ? '0'+new Date(source).getSeconds().toString() : new Date(source).getSeconds().toString();

            var rtDate = new Date(source).getFullYear()+ '.' +toMonth+ '.' +toDate+ ' ' +toHours+ ':' +toMinutes;
            return rtDate;
        };
        common.convertTimestampToDate = function(sourceTimestamp) {
            var source = moment(Number(sourceTimestamp)).unix() * 1000;

            var toMonth = (new Date(source).getMonth()+1).toString().length === 1 ? '0'+(new Date(source).getMonth()+1).toString() : (new Date(source).getMonth()+1).toString();
            var toDate = new Date(source).getDate().toString().length === 1 ? '0'+new Date(source).getDate().toString() : new Date(source).getDate().toString();

            var rtDate = new Date(source).getFullYear()+ '.' +toMonth+ '.' +toDate;
            return rtDate;
        };


        common.alarmLevelStyle = function(status) {
            var style = {};
            if(status == 'critical') {
                style = 'critical';
            } else if(status == 'warning') {
                style = 'warning';
            }
            return style;
        };

        common.resolveStatusStyle = function(status) {
            var style = {};
            if(status == 1) {
                style = 'waiting';
            } else if(status == 2) {
                style = 'processing';
            } else {
                style = '';
            }
            return style;
        };

        return common;
    })
    .factory('cache', function (localStorageService) {
        var cache = {};

        cache.setUser = function(user, expires) {
            localStorageService.set('monitor.user', user, expires);
        };
        cache.getUser = function() {
            return localStorageService.get("monitor.user");
        };

        cache.clear = function() {
            localStorageService.clearAll();
        };

        cache.isAuthenticated = function() {
            return (cache.getUser() != null);
        };

        return cache;
    })
    .factory('cookies', function ($cookies) {
        var cookies = {};

        cookies.putDefaultTimeRange = function(defaultTimeRange) {
            $cookies.put("defaultTimeRange", defaultTimeRange);
        };
        cookies.getDefaultTimeRange = function() {
            return $cookies.get("defaultTimeRange");
        };
        cookies.removeDefaultTimeRange = function() {
            $cookies.remove("defaultTimeRange");
        };

        cookies.putTimeRangeFrom = function(timeRangeFrom) {
            $cookies.put("timeRangeFrom", timeRangeFrom);
        };
        cookies.getTimeRangeFrom = function() {
            return $cookies.get("timeRangeFrom");
        };
        cookies.removeTimeRangeFrom = function() {
            $cookies.remove("timeRangeFrom");
        };

        cookies.putTimeRangeTo = function(timeRangeTo) {
            $cookies.put("timeRangeTo", timeRangeTo);
        };
        cookies.getTimeRangeTo = function() {
            return $cookies.get("timeRangeTo");
        };
        cookies.removeTimeRangeTo = function() {
            $cookies.remove("timeRangeTo");
        };

        cookies.putGroupBy = function(groupBy) {
            $cookies.put("groupBy", groupBy);
        };
        cookies.getGroupBy = function() {
            return $cookies.get("groupBy");
        };
        cookies.removeGroupBy = function() {
            $cookies.remove("groupBy");
        };

        cookies.putRefreshTime = function(refreshTime) {
            $cookies.put("refreshTime", refreshTime);
        };
        cookies.getRefreshTime = function() {
            return $cookies.get("refreshTime");
        };
        cookies.removeRefreshTime = function() {
            $cookies.remove("refreshTime");
        };

        return cookies;
    })
;
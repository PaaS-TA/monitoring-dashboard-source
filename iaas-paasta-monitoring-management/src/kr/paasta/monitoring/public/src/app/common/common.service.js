(function() {
  'use strict';

  angular
    .module('monitoring')
    .service('common', CommonService)
    .service('cookies', CookieService)
    .service('cache', CacheService)
    .service('anchorSmoothScroll', ScrollService)
    .service('iaasLogService', IaasLogService)
    .service('paasLogService', PaasLogService);

  /** @ngInject */
  function CommonService($http) {
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
        config.data = angular.fromJson(body);
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
    };
    common.getDateTTime = function(str) {
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

    /**
     * @return {number}
     */
    common.CompareForSort = function(first, second){
      var n = 0;
      if (first.id == second.id)
        n =  0;
      else if (first.id < second.id)
        n =  -1;
      else
        n =  1;
      return n;
    };

    common.setDtvParam = function(condition) {
      var param = {};
      if(condition.timeRangeFrom && condition.timeRangeTo) {
        param['timeRangeFrom'] = condition.timeRangeFrom;
        param['timeRangeTo'] = condition.timeRangeTo;
      } else {
        param['defaultTimeRange'] = condition.defaultTimeRange==undefined?'15m':condition.defaultTimeRange;
      }
      param['groupBy'] = condition.groupBy==undefined?'1m':condition.groupBy;

      return param;
    };

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
        case '30d':
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

    common.alarmLevelStyle = function(status) {
      var style = {};
      if(status == 'critical') {
        style = 'critical';
      } else if(status == 'warning') {
        style = 'warning';
      } else if(status == 'fail') {
        style = 'fail';
      } else if(status == 'running' || status == 'healthy') {
        style = 'running';
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

    common.convertTimestampToDate = function(sourceTimestamp) {
      var source = moment(Number(sourceTimestamp)).unix() * 1000;

      var toMonth = (new Date(source).getMonth()+1).toString().length === 1 ? '0'+(new Date(source).getMonth()+1).toString() : (new Date(source).getMonth()+1).toString();
      var toDate = new Date(source).getDate().toString().length === 1 ? '0'+new Date(source).getDate().toString() : new Date(source).getDate().toString();
      var toHours = new Date(source).getHours().toString().length === 1 ? '0'+new Date(source).getHours().toString() : new Date(source).getHours().toString();
      var toMinutes = new Date(source).getMinutes().toString().length === 1 ? '0'+new Date(source).getMinutes().toString() : new Date(source).getMinutes().toString();
      //var toSeconds = new Date(source).getSeconds().toString().length === 1 ? '0'+new Date(source).getSeconds().toString() : new Date(source).getSeconds().toString();

      var rtDate = new Date(source).getFullYear()+
        '.' +toMonth+
        '.' +toDate+
        ' ' +toHours+
        ':' +toMinutes;
      return rtDate;
    };

    return common;
  }

  function CookieService($cookies) {

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

      // server detail view layout
      cookies.putSvrDtvLayout = function(num) {
        $cookies.put("svrDtvLayout", num);
      };
      cookies.getSvrDtvLayout = function() {
        return $cookies.get("svrDtvLayout");
      };
      cookies.removeSvrDtvLayout = function() {
        $cookies.remove("svrDtvLayout");
      };

      // server detail view options
      var svrDtvOpt = [];
      cookies.putSvrDtvOpt = function(dtvOpt) {
        svrDtvOpt.push(dtvOpt);
        $cookies.putObject("svr.dtvOpt", svrDtvOpt);
      };
      cookies.getSvrDtvOpt = function() {
        return $cookies.getObject("svr.dtvOpt");
      };
      cookies.removeSvrDtvOpt = function() {
        svrDtvOpt = [];
        $cookies.remove("svr.dtvOpt");
      };

      // infrastructure detail view layout
      cookies.putIfrDtvLayout = function(num) {
        $cookies.put("ifrDtvLayout", num);
      };
      cookies.getIfrDtvLayout = function() {
        return $cookies.get("ifrDtvLayout");
      };
      cookies.removeIfrDtvLayout = function() {
        $cookies.remove("ifrDtvLayout");
      };

      // infrastructure detail view options
      var ifrDtvOpt = [];
      cookies.putIfrDtvOpt = function(dtvOpt) {
        ifrDtvOpt.push(dtvOpt);
        $cookies.putObject("ifr.dtvOpt", ifrDtvOpt);
      };
      cookies.getIfrDtvOpt = function() {
        return $cookies.getObject("ifr.dtvOpt");
      };
      cookies.removeIfrDtvOpt = function() {
        ifrDtvOpt = [];
        $cookies.remove("ifr.dtvOpt");
      };

      // Controller detail view layout
      cookies.putCtlDtvLayout = function(num) {
        $cookies.put("ctlDtvLayout", num);
      };
      cookies.getCtlDtvLayout = function() {
        return $cookies.get("ctlDtvLayout");
      };
      cookies.removeCtlDtvLayout = function() {
        $cookies.remove("ctlDtvLayout");
      };

      // Controller detail view options
      var ctlDtvOpt = [];
      cookies.putCtlDtvOpt = function(dtvOpt) {
        ctlDtvOpt.push(dtvOpt);
        $cookies.putObject("ctl.dtvOpt", ctlDtvOpt);
      };
      cookies.getCtlDtvOpt = function() {
        return $cookies.getObject("ctl.dtvOpt");
      };
      cookies.removeCtlDtvOpt = function() {
        ctlDtvOpt = [];
        $cookies.remove("ctl.dtvOpt");
      };

      // Zone detail view layout
      cookies.putZonDtvLayout = function(num) {
        $cookies.put("zonDtvLayout", num);
      };
      cookies.getZonDtvLayout = function() {
        return $cookies.get("zonDtvLayout");
      };
      cookies.removeZonDtvLayout = function() {
        $cookies.remove("zonDtvLayout");
      };

      // Zone detail view options
      var zonDtvOpt = [];
      cookies.putZonDtvOpt = function(dtvOpt) {
        zonDtvOpt.push(dtvOpt);
        $cookies.putObject("zon.dtvOpt", zonDtvOpt);
      };
      cookies.getZonDtvOpt = function() {
        return $cookies.getObject("zon.dtvOpt");
      };
      cookies.removeZonDtvOpt = function() {
        zonDtvOpt = [];
        $cookies.remove("zon.dtvOpt");
      };

      // Service detail view layout
      cookies.putSvcDtvLayout = function(num) {
        $cookies.put("svcDtvLayout", num);
      };
      cookies.getSvcDtvLayout = function() {
        return $cookies.get("svcDtvLayout");
      };
      cookies.removeSvcDtvLayout = function() {
        $cookies.remove("svcDtvLayout");
      };

      // Service detail view options
      var svcDtvOpt = [];
      cookies.putSvcDtvOpt = function(dtvOpt) {
        svcDtvOpt.push(dtvOpt);
        $cookies.putObject("svc.dtvOpt", svcDtvOpt);
      };
      cookies.getSvcDtvOpt = function() {
        return $cookies.getObject("svc.dtvOpt");
      };
      cookies.removeSvcDtvOpt = function() {
        svcDtvOpt = [];
        $cookies.remove("svc.dtvOpt");
      };

      // Service detail view layout(Postgres)
      cookies.putSvcDtvPgrLayout = function(num) {
        $cookies.put("svcDtvPgrLayout", num);
      };
      cookies.getSvcDtvPgrLayout = function() {
        return $cookies.get("svcDtvPgrLayout");
      };
      cookies.removeSvcDtvPgrLayout = function() {
        $cookies.remove("svcDtvPgrLayout");
      };

      // Service detail view options(Postgres)
      var svcDtvPgrOpt = [];
      cookies.putSvcDtvPgrOpt = function(svcOpt) {
        svcDtvPgrOpt.push(svcOpt);
        $cookies.putObject("svc.dtvPgrOpt", svcDtvPgrOpt);
      };
      cookies.getSvcDtvPgrOpt = function() {
        return $cookies.getObject("svc.dtvPgrOpt");
      };
      cookies.removeSvcDtvPgrOpt = function() {
        svcDtvPgrOpt = [];
        $cookies.remove("svc.dtvPgrOpt");
      };

      // Service detail view layout(mysql)
      cookies.putSvcDtvMysqlLayout = function(num) {
        $cookies.put("svcDtvMysqlLayout", num);
      };
      cookies.getSvcDtvMysqlLayout = function() {
        return $cookies.get("svcDtvMysqlLayout");
      };
      cookies.removeSvcDtvMysqlLayout = function() {
        $cookies.remove("svcDtvMysqlLayout");
      };

      // Service detail view options(mysql)
      var svcDtvMysqlOpt = [];
      cookies.putSvcDtvMysqlOpt = function(svcOpt) {
        svcDtvMysqlOpt.push(svcOpt);
        $cookies.putObject("svc.dtvMysqlOpt", svcDtvMysqlOpt);
      };
      cookies.getSvcDtvMysqlOpt = function() {
        return $cookies.getObject("svc.dtvMysqlOpt");
      };
      cookies.removeSvcDtvMysqlOpt = function() {
        svcDtvMysqlOpt = [];
        $cookies.remove("svc.dtvMysqlOpt");
      };

      // Service detail view layout(redis)
      cookies.putSvcDtvRedisLayout = function(num) {
        $cookies.put("svcDtvRedisLayout", num);
      };
      cookies.getSvcDtvRedisLayout = function() {
        return $cookies.get("svcDtvRedisLayout");
      };
      cookies.removeSvcDtvRedisLayout = function() {
        $cookies.remove("svcDtvRedisLayout");
      };

      // Service detail view options(Redis)
      var svcDtvRedisOpt = [];
      cookies.putSvcDtvRedisOpt = function(svcOpt) {
        svcDtvRedisOpt.push(svcOpt);
        $cookies.putObject("svc.dtvRedisOpt", svcDtvRedisOpt);
      };
      cookies.getSvcDtvRedisOpt = function() {
        return $cookies.getObject("svc.dtvRedisOpt");
      };
      cookies.removeSvcDtvRedisOpt = function() {
        svcDtvRedisOpt = [];
        $cookies.remove("svc.dtvRedisOpt");
      };

      // IaaS Infrastructure Threshold
      cookies.putIfrThreshold = function(obj) {
        $cookies.putObject("ifr.threshold", obj);
      };
      cookies.getIfrThreshold = function() {
        return $cookies.getObject("ifr.threshold");
      };
      cookies.removeIfrThreshold  = function() {
        $cookies.remove("ifr.threshold");
      };

      // PaaS Controller Threshold
      cookies.putCtlThreshold = function(obj) {
        $cookies.putObject("ctl.threshold", obj);
      };
      cookies.getCtlThreshold = function() {
        return $cookies.getObject("ctl.threshold");
      };
      cookies.removeCtlThreshold  = function() {
        $cookies.remove("ctl.threshold");
      };

      // Application Zone Threshold
      cookies.putZonThreshold = function(obj) {
        $cookies.putObject("zon.threshold", obj);
      };
      cookies.getZonThreshold = function() {
        return $cookies.getObject("zon.threshold");
      };
      cookies.removeZonThreshold  = function() {
        $cookies.remove("zon.threshold");
      };

      // External Service Threshold
      cookies.putSvcThreshold = function(obj) {
        $cookies.putObject("svc.threshold", obj);
      };
      cookies.getSvcThreshold = function() {
        return $cookies.getObject("svc.threshold");
      };
      cookies.removeSvcThreshold  = function() {
        $cookies.remove("svc.threshold");
      };

      // Hardware Threshold
      cookies.putSvrThreshold = function(obj) {
        $cookies.putObject("svr.threshold", obj);
      };
      cookies.getSvrThreshold = function() {
        return $cookies.getObject("svr.threshold");
      };
      cookies.removeSvrThreshold  = function() {
        $cookies.remove("svr.threshold");
      };

      return cookies;
  }

  function CacheService(localStorageService) {

    var cache = {};

    cache.setUser = function(user, expires) {
      localStorageService.set('user', user, expires);
    };
    cache.getUser = function() {
      return localStorageService.get("user");
    };
    cache.clear = function() {
      localStorageService.clearAll();
    };
    cache.isAuthenticated = function() {
      return (cache.getUser() != null);
    };

    cache.putToken = function(token) {
      localStorageService.set('token', token);
    };
    cache.getToken = function() {
      return localStorageService.get('token');
    };
    cache.removeToken = function() {
      localStorageService.remove('token');
    };

    cache.setUserAuth = function(auth, expires) {
      localStorageService.set('auth', auth, expires);
    };

    cache.getUserAuth = function() {
      return localStorageService.get("auth");
    };

    cache.setSysType = function(sysType, expires) {
      localStorageService.set('sysType', sysType, expires);
    };

    cache.getSysType = function() {
      return localStorageService.get("sysType");
    };

    return cache;

  }

  function ScrollService($document) {
    this.scrollTo = function(eID) {

      // This scrolling function
      // is from http://www.itnewb.com/tutorial/Creating-the-Smooth-Scroll-Effect-with-JavaScript

      var startY = currentYPosition();
      var stopY = elmYPosition(eID);
      var distance = stopY > startY ? stopY - startY : startY - stopY;
      if (distance < 100) {
        scrollTo(0, stopY); return;
      }
      var speed = Math.round(distance / 100);
      if (speed >= 20) speed = 20;
      var step = Math.round(distance / 25);
      var leapY = stopY > startY ? startY + step : startY - step;
      var timer = 0;
      if (stopY > startY) {
        for ( var i=startY; i<stopY; i+=step ) {
          angular.element('body').animate({scrollTop: stopY}, "slow");
          leapY += step; if (leapY > stopY) leapY = stopY; timer++;
        } return;
      }
      for ( var ii=startY; ii>stopY; ii-=step ) {
        angular.element('body').animate({scrollTop: stopY}, "slow");
        leapY -= step; if (leapY < stopY) leapY = stopY; timer++;
      }

      function currentYPosition() {
        // Firefox, Chrome, Opera, Safari
        if (self.pageYOffset) return self.pageYOffset;
        // Internet Explorer 6, 7 and 8
        if ($document[0].body.scrollTop) return $document[0].body.scrollTop;
        return 0;
      }

      function elmYPosition(eID) {
        var elm = $document[0].getElementById(eID);
        var y = elm.offsetTop;
        var node = elm;
        while (node.offsetParent && node.offsetParent != $document[0].body) {
          node = node.offsetParent;
          y += node.offsetTop;
        } return y;
      }

    };
  }

  function IaasLogService($http, apiUris) {
    var log = {};

    log.dtvDefaultRecentLog = function(param) {
      var config = {
        params: param,
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasDefaultRecentLogs, config);
    };
    log.dtvSpecificTimeRangeLog = function(param) {
      var config = {
        params: param,
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.iaasSpecificTimeRangeLogs, config);
    };

    return log;
  }

  function PaasLogService($http, apiUris) {
    var log = {};

    log.dtvDefaultRecentLog = function(param) {
      var config = {
        params: param,
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasDefaultRecentLogs, config);
    };
    log.dtvSpecificTimeRangeLog = function(param) {
      var config = {
        params: param,
        headers : {'Accept' : 'application/json'}
      };
      return $http.get(apiUris.paasSpecificTimeRangeLogs, config);
    };

    return log;
  }


})();

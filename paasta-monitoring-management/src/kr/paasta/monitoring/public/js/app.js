angular.module('app', [
    'monitor.controllers', 'monitor.services', 'exception.service',
    'ui.router', 'ui.bootstrap', 'ngCookies',
    'ae-datetimepicker', 'simplePagination'
])

/**
 * Config -------------------------------------------------------------------------
 */
    .config(['$httpProvider', '$stateProvider', '$urlRouterProvider', function($httpProvider, $stateProvider, $urlRouterProvider) {

        $httpProvider.interceptors.push('authInterceptor');

        $httpProvider.defaults.useXDomain = true;
        delete $httpProvider.defaults.headers.common['X-Requested-With'];

        $urlRouterProvider.otherwise('/');

        var navigation = {
            templateUrl: 'partials/navigation.html',
            controller: 'navigationCtrl'
        };

        $stateProvider.state('main', {
            url: '/',
            views: {
                'navigation': navigation,
                'body': {
                    controller: 'mainCtrl'
                }
            }
        });
        $stateProvider.state('almLst', {
            url: '/almLst?_f&_t&_o&_a&_l',
            views: {
                'navigation': navigation,
                'body': {
                    templateUrl: 'partials/almLst.html',
                    controller: 'almLstCtrl'
                }
            }
        });
        $stateProvider.state('almLstDtl', {
            url: '/almLst/:id',
            views: {
                'navigation': navigation,
                'body': {
                    templateUrl: 'partials/almLstDtl.html',
                    controller: 'almLstDtlCtrl'
                }
            }
        });
        $stateProvider.state('almStt', {
            url: '/almStt',
            views: {
                'navigation': navigation,
                'body': {
                    templateUrl: 'partials/almStt.html',
                    controller: 'almSttCtrl'
                }
            }
        });
        $stateProvider.state('almSet', {
            url: '/almSet',
            views: {
                'navigation': navigation,
                'body': {
                    templateUrl: 'partials/almSet.html',
                    controller: 'almSetCtrl'
                }
            }
        });
        $stateProvider.state('ctnPst', {
            url: '/ctnPst',
            views: {
                'navigation': navigation,
                'body': {
                    templateUrl: 'partials/ctnPst.html',
                    controller: 'ctnPstCtrl'
                }
            }
        });
    }])

    /**
     * localStorageService -------------------------------------------------------------------------
     */
    .config(function ($provide, CONSTANTS) {
        $provide.decorator('localStorageService', function($delegate) {
            //store original get & set methods
            var originalGet = $delegate.get,
                originalSet = $delegate.set;

            /**
             * extending the localStorageService get method
             *
             * @param key
             * @returns {*}
             */
            $delegate.get = function(key) {
                if(originalGet(key)) {
                    var data = originalGet(key);

                    if(data.expire) {
                        var now = Date.now();

                        // delete the key if it timed out
                        if(data.expire < now) {
                            $delegate.remove(key);
                            return null;
                        } else {
                            var expiryDate = Date.now() + (1000 * 60 * 60 * CONSTANTS.expire);
                            originalSet(key, {
                                data: data.data,
                                expire: expiryDate
                            });
                        }

                        return data.data;
                    } else {
                        return data;
                    }
                } else {
                    return null;
                }
            };

            /**
             * set
             * @param key               key
             * @param val               value to be stored
             * @param {int} expires     hours until the localStorage expires (hour)
             */
            $delegate.set = function(key, val, expires) {
                var expiryDate = null;

                if(angular.isNumber(expires)) {
                    expiryDate = Date.now() + (1000 * 60 * 60 * expires);
                    originalSet(key, {
                        data: val,
                        expire: expiryDate
                    });
                } else {
                    originalSet(key, val);
                }
            };

            return $delegate;
        });
    })

    /**
     * Constants -------------------------------------------------------------------------
     */
    .constant('CONSTANTS', {
        version: '0.0.1',
        apiServer: '',
//        apiServer: 'http://localhost:8080',
        context: '/',
        expire: 1,
        ifr: 'IaaS Interface',
        ctl: 'PaaS Controller',
        zon: 'Application Zone',
        svc: 'External Service',
        svr: 'Hardware',
        ChannelTypeCode: 'chnlType',
        portalApiServer: 'http://localhost:8080'
//        portalApiServer: 'http://portal-api.paasxpert.com:8080'
    })
    .constant('DTVOPT', [
        {id: 1, name: 'CPU', func: 'dtvCpuUsage', type: 'lineChart', percent: true, axisLabel: '%'},
        {id: 2, name: 'Memory', func: 'dtvMemoryUsage', type: 'lineChart', percent: true, axisLabel: '%'},
        {id: 3, name: 'Disk (Mounted FileSystem Occupied)', func: 'dtvFileSystemUsage', type: 'lineChart', percent: true, axisLabel: '%'},
        {id: 4, name: 'Network', func: 'dtvNetworkUsage', type: 'lineChart', percent: false, axisLabel: 'Kb/Sec'},
        {id: 5, name: 'Network Packets', func: 'dtvNetworkPacketsUsage', type: 'lineChart', percent: false, axisLabel: 'Packets/Sec'},
        {id: 6, name: 'CPU Load (1min)', func: 'dtvCpuLoadAvgUsage', type: 'lineChart', percent: false},
        {id: 7, name: 'Disk IO', func: 'dtvDiskIOUsage', type: 'lineChart', percent: false, axisLabel: 'Operations/Sec'},
        {id: 8, name: 'Disk Utilization', func: 'dtvDiskUtilizationUsage', type: 'lineChart', percent: true, axisLabel: '%'},
        {id: 9, name: 'Disk Device (Physical Disk Occupied)', func: 'dtvDiskDeviceUsage', type: 'lineChart', percent: true, axisLabel: '%'},
        {id: 10, name: 'Swap Memory', func: 'dtvSwapMemoryUsage', type: 'lineChart', percent: true, axisLabel: '%'},
        {id: 11, name: 'Network Drop Packets', func: 'dtvNetworkDropPacketsUsage', type: 'multiBarChart', percent: false, axisLabel: 'Count'},
        {id: 12, name: 'Network Packets Error', func: 'dtvNetworkPacketsErrUsage', type: 'lineChart', percent: false, axisLabel: 'Count'},
        /*{id: 13, name: 'Server Boot Time', func: 'dtvBootTimeUsage', type: 'lineChart', percent: false, axisLabel: 'Restart Time'},*/
        {id: 13, name: 'Top Process', func: 'dtvTopProcessUsage', type: 'list'}
    ])
    .constant('ZONDTVOPT', [
        {id: 1, name: 'CPU', func: 'dtvCpuUsage', type: 'lineChart', percent: true, axisLabel: '%'},
        {id: 2, name: 'Memory', func: 'dtvMemoryUsage', type: 'lineChart', percent: true, axisLabel: '%'},
        {id: 3, name: 'Disk', func: 'dtvDiskUsage', type: 'lineChart', percent: true, axisLabel: '%'},
        {id: 4, name: 'CPU Load (1min)', func: 'dtvCpuLoadAvgUsage', type: 'lineChart', percent: false},
        {id: 5, name: 'Total Memory', func: 'dtvTotalMemoryUsage', type: 'lineChart', percent: false},
        {id: 6, name: 'Total Disk', func: 'dtvTotalDiskUsage', type: 'lineChart', percent: false},
        {id: 7, name: 'Network', func: 'dtvNetworkUsage', type: 'lineChart', percent: false, axisLabel: 'Kb/Sec'},
        {id: 8, name: 'Network Drop Packets', func: 'dtvNetworkDropPacketsUsage', type: 'multiBarChart', percent: false, axisLabel: 'Count'},
        {id: 9, name: 'Network Packets Error', func: 'dtvNetworkPacketsErrUsage', type: 'lineChart', percent: false, axisLabel: 'Count'}
    ])
    .constant('POSTGRES_OPT', [
        {id: 1, name: 'Statistics', func: 'dtvStatistics', type: 'multiBarChart', percent: false, axisLabel: 'count'},
        {id: 2, name: 'Runtime', func: 'dtvRuntime', type: 'lineChart', percent: false, axisLabel: 'count'},
        {id: 3, name: 'Buffer', func: 'dtvBuffer', type: 'lineChart', percent: false, axisLabel: 'count'},
        {id: 4, name: 'CheckPoint(ms)', func: 'dtvCheckPointMs', type: 'lineChart', percent: false, axisLabel: 'MilliSeconds'},
        {id: 5, name: 'CheckPoint(count)', func: 'dtvCheckPointCount', type: 'lineChart', percent: false, axisLabel: 'count'},
        {id: 6, name: 'Blocks', func: 'dtvBlocks', type: 'multiBarChart', percent: false, axisLabel: 'Per/Second'},
        {id: 7, name: 'BlockHits', func: 'dtvBlockHits', type: 'lineChart', percent: false, axisLabel: 'Per/Second'},
        {id: 8, name: 'Conflicts', func: 'dtvConflicts', type: 'multiBarChart', percent: false, axisLabel: 'Count'},
        {id: 9, name: 'Deadlocks', func: 'dtvDeadlocks', type: 'multiBarChart', percent: false, axisLabel: 'Count'},
        {id: 10, name: 'Temporary', func: 'dtvTemporary', type: 'lineChart', percent: false, axisLabel: 'Per/Second'},
        {id: 11, name: 'Rows', func: 'dtvRows', type: 'lineChart', percent: false, axisLabel: 'Per/Second'},
        {id: 12, name: 'Transaction', func: 'dtvTransaction', type: 'lineChart', percent: false, axisLabel: 'Per/Second'}
    ])
    .constant('MYSQL_OPT', [
        {id: 1, name: 'Connections', func: 'dtvMysqlConnections', type: 'multiBarChart', percent: false, axisLabel: 'connection/second'},
        {id: 2, name: 'Innodb Read', func: 'dtvMysqlInnodbread', type: 'lineChart', percent: false, axisLabel: 'read/second'},
        {id: 3, name: 'Innodb Write', func: 'dtvMysqlInnodbwrite', type: 'lineChart', percent: false, axisLabel: 'writes/second'},
        {id: 4, name: 'Log File Sync', func: 'dtvMysqlLogfilesync', type: 'lineChart', percent: false, axisLabel: 'write/second'},
        {id: 5, name: 'Slow Queries', func: 'dtvMysqlSlowqueries', type: 'lineChart', percent: false, axisLabel: 'query/second'},
        {id: 6, name: 'Row', func: 'dtvMysqlRow', type: 'multiBarChart', percent: false, axisLabel: 'per/second'},
        {id: 7, name: 'Buffer Pools', func: 'dtvMysqlBufferpools', type: 'lineChart', percent: false, axisLabel: 'count'},
        {id: 8, name: 'Threads', func: 'dtvMysqlThreads', type: 'multiBarChart', percent: false, axisLabel: 'count'},
        {id: 9, name: 'Table Locks', func: 'dtvMysqlTablelocks', type: 'multiBarChart', percent: false, axisLabel: 'count'},
        {id: 10, name: 'Temp Tables', func: 'dtvMysqlTemptables', type: 'lineChart', percent: false, axisLabel: 'per/second'},
        {id: 11, name: 'Tempfiles', func: 'dtvMysqlTempfiles', type: 'lineChart', percent: false, axisLabel: 'per/second'},
        {id: 12, name: 'Network Usage', func: 'dtvMysqlNetworkusage', type: 'lineChart', percent: false, axisLabel: 'kbs'},
        {id: 13, name: 'Aborts', func: 'dtvMysqlAborts', type: 'lineChart', percent: false, axisLabel: 'count'},
        {id: 14, name: 'Transactions', func: 'dtvMysqlTransactions', type: 'lineChart', percent: false, axisLabel: 'count'}
    ])

    /**
     * Run -------------------------------------------------------------------------
     */
    .run(function(cookies) {

        /********** default start **********/
        if(cookies.getDefaultTimeRange() === undefined) cookies.putDefaultTimeRange('15m');
        if(cookies.getGroupBy() === undefined) cookies.putGroupBy('1m');
        if(cookies.getRefreshTime() === undefined) cookies.putRefreshTime('off');
    })
;
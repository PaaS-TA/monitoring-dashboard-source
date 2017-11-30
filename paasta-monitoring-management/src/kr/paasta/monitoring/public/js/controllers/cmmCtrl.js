'use strict';

angular.module('monitor.controllers', [])

    /******************** Main ********************/
    .controller('mainCtrl', function($scope, $location){
        $location.path('/almLst');
    })

    /******************** Navigation ********************/
    .controller('navigationCtrl', function($scope, $rootScope, $location, $stateParams, common, almSrvc){
        $scope.common = common;

        switch($location.path()){
            case '/almLst' :
            case '/almLst/'+$stateParams.id :
                $scope.selected = 'als';
                break;
            case '/almStt' :
                $scope.selected = 'ast';
                break;
            case '/almSet' :
                $scope.selected = 'ase';
                break;
            case '/ctnPst' :
                $scope.selected = 'con';
                break;
            default:
                $scope.selected = '';
        }
    })
;

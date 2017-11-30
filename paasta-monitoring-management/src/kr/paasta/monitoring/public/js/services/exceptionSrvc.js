'use strict';

angular.module('exception.service', [])
    .factory('$exceptionHandler', function($injector, $log) {
    return function exceptionHandler(exception, cause) {

        var $rootScope = $injector.get('$rootScope');
        if($rootScope){
            $rootScope.loading = false;
            angular.element('.loader-container').hide();
            if(cause != undefined) {
                $rootScope.errorCode = cause.code;
                if(cause.message != undefined) {
                    $rootScope.errorMessage = cause.message;
                }
            }
            if(exception != undefined) {
                if($rootScope.errorMessage == undefined) {
                    $rootScope.errorMessage = exception;
                }
                $log.error('Error!!! ' + $rootScope.errorMessage + '. Cause: ' + JSON.stringify(cause));
                angular.element('#errorMessage').html($rootScope.errorMessage);
                angular.element('#errorModal').modal('show');
            }

            angular.element('#errorModal').on('hidden.bs.modal', function () {
                $rootScope.errorMessage = undefined;
                $rootScope.errorCode = undefined;
            });
        }
    };
});
(function() {
  'use strict';

  angular
    .module('monitoring')
    .controller('LoginController', LoginController);

  /** @ngInject */
  function LoginController($scope, $timeout, $location, $http, $exceptionHandler, apiUris, cache, constants, loginService) {
    var vm = this;
    vm.scope = $scope;
    vm.Math = Math;

    angular.element('.wrapper').css('min-height',0);
    angular.element('.wrapper').css('margin-top',0);
    angular.element('.loginWrapper').css('height','100%');

    vm.credentials = {
      username: '',
      password: '',
      rememberMe: false
    };

    vm.login = function() {
      vm.scope.loading = true;

      var credentials = {
        username: vm.email,
        password: vm.password
      };

      loginService.ping().then(
        function() {
          loginService.authenticate(credentials).then(
            function(response) {
              if (!response) {
                vm.scope.loading = false;
                vm.scope.login.password = '';
                var vMsg = 'Username or password is incorrect. Please check back again.';
                $timeout(function() { $exceptionHandler(vMsg, {code: 401, message: vMsg}); }, 500);
                $location.path('/login');
              } else {

                var data = response.data;
                cache.setUser({
                  name: data.username,
                  email: data.userEmail,
                  sysType: data.sysType,
                  i1: data.authI1,        // (시스템 설정이 IaaS 를 사용하는 경우 : I , 미사용 : N)
                  i2: data.authI2,        // (회원가입시 사용자 IaaS 계정이 정상인경우 : S , 계정인증 실패 : F)
                  p1: data.authP1,        // (시스템 설정이 IaaS 를 사용하는 경우 : P , 미사용 : N)
                  p2: data.authP2         // (회원가입시 사용자 PaaS 계정이 정상인경우 : S , 계정인증 실패 : F)
                }, constants.expire);

                var url = "";
                if(cache.getUser().sysType == "IaaS") {
                  if(cache.getUser().i2 == "F") {
                    url = "/member/info";
                  } else {
                    url = "/iaas/main";
                  }
                } else if(cache.getUser().sysType == "PaaS") {
                  if(cache.getUser().p2 == "F") {
                    url = "/member/info";
                  } else {
                    url = "/paas/main";
                  }
                } else {
                  if(((cache.getUser().p1 == "P" && cache.getUser().p2 == "S") && (cache.getUser().i1 == "I" && cache.getUser().i2 == "F"))
                    || (cache.getUser().i1 == "F" && cache.getUser().i2 == "F")) {
                    url = "/paas/main";
                  } else if(((cache.getUser().i1 == "I" && cache.getUser().i2 == "S") && (cache.getUser().p1 == "P" && cache.getUser().p2 == "F"))
                    || (cache.getUser().p1 == "F" && cache.getUser().p2 == "F")) {
                    url = "/iaas/main";
                  } else if(cache.getUser().i2 == "S" && cache.getUser().p2 == "S") {
                    url = "/";
                  } else {
                    url = "/member/info";
                  }
                }

                $location.path(url);

                angular.element('.wrapper').css('min-height', 100);
                angular.element('.wrapper').css('margin-top', 70);
              }
            },
            function(reason) {
              vm.scope.loading = false;
              $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
            }
          );
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };


    // 아이디 찾기
    vm.findId = function() {

    };


    // 비밀번호 찾기
    vm.findPassword = function() {

    };

  }
})();

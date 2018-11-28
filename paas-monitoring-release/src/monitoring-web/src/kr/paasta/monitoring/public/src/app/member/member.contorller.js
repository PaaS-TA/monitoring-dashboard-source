(function() {
  'use strict';

  angular
    .module('monitoring')
    .controller('MemberController', MemberController);

  /** @ngInject */
  function MemberController($scope, $timeout, $location, $window, $document, $exceptionHandler, cache, constants, memberService) {
    var vm = this;
    vm.scope = $scope;
    vm.memberInfo = {};
    vm.scope.sysConfig = "";


    // Case Join Css Change
    if($location.path().indexOf("/join") > -1) {
      angular.element('.wrapper').css('margin-top', 0);
      vm.scope.type = false;
    } else {
      vm.scope.type = true;
    }


    // Member Join Init
    (vm.getInit = function() {
      memberService.init().then(
        function(result) {
          vm.scope.sysConfig = result.data;

          vm.iaasCertificationMsg = "인증이 필요합니다.";
          vm.paasCertificationMsg = "인증이 필요합니다.";

          if(vm.scope.sysConfig == "IaaS") {
            vm.rightClass = "col-md-0";
            vm.leftClass = "col-md-12";
            vm.rightShow = false;
            vm.leftShow = true;
            vm.memberInfo.iaasUserUseYn = "Y";
            vm.cancelUrl = "/iaas/main";
          } else if(vm.scope.sysConfig == "PaaS") {
            vm.rightClass = "col-md-12";
            vm.leftClass = "col-md-0";
            vm.rightShow = true;
            vm.leftShow = false;
            vm.memberInfo.paasUserUseYn = "Y";
            vm. cancelUrl = "/paas/main";
          } else {
            vm.rightClass = "col-md-6";
            vm.leftClass = "col-md-6";
            vm.rightShow = true;
            vm.leftShow = true;
            vm.pwShow = true;
            vm.cancelUrl = "/";
          }

          if(cache.getUser() != null) {
            vm.pwShow = false;
            vm.getMemberInfoView();

            if(cache.getUser().i2 == "S" && cache.getUser().p2 == "F") {
              vm.cancelUrl = "/iaas/main";
            } else if(cache.getUser().i2 == "F" && cache.getUser().p2 == "S") {
              vm.cancelUrl = "/paas/main";
            } else if(cache.getUser().i2 == "F" && cache.getUser().p2 == "F") {
              vm.cancelUrl = "/member/info";
            } else {
              vm.cancelUrl = "/";
            }
          }
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    })();


    // Member Info Selete
    vm.getMemberInfoView = function() {
      var userId = cache.getUser().name;
      memberService.memberInfoView(userId).then(
        function(result) {
          vm.memberInfo = result.data;
          vm.memberInfo.userPw = "";
          vm.memberInfo.userPwConfirm = "";

          if((cache.getUser().sysType == "ALL" && result.data.iaasUserUseYn == "Y") || (cache.getUser().sysType == "IaaS" && result.data.iaasUserUseYn == "Y")) {
            vm.certificationConfirm("IaaS", true);
          }

          if((cache.getUser().sysType == "ALL" && result.data.paasUserUseYn == "Y") || (cache.getUser().sysType == "PaaS" && result.data.paasUserUseYn == "Y")) {
            vm.certificationConfirm("PaaS", true);
          }

          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    }


    // Duplicate Check
    vm.duplicateConfirm = function(val) {

      if(cache.getUser() == null) {
        if(val == "userId" && vm.memberInfo.userId != "" && vm.memberInfo.userId != undefined) {
          memberService.duplicateConfirmId(vm.memberInfo.userId).then(
            function(result) {
              if (result.data != null && result.data != "") {
                vm.alertModal('회원가입', '중복된 아이디 입니다.', undefined, false, true);
                vm.memberInfo.userId = "";
              }
            },
            function (reason) {
              $timeout(function () {
                $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message});
              }, 500);
            }
          );
        }
      }

      if(val == "iaasUser" && vm.memberInfo.iaasUserId != "" && vm.memberInfo.iaasUserId != undefined) {
        if(vm.memberInfo.iaasUserUseYn == 'Y') {
          memberService.iaasDuplicateCheckId(vm.memberInfo.iaasUserId).then(
            function (result) {
              if(result.data == vm.memberInfo.iaasUserId) {
                vm.memberInfo.iaasUserChck = false;
                vm.iaasCertificationMsg = '이미 사용중이 아이디 입니다.';
                vm.memberInfo.iaasUserId = '';
              } else {
                vm.iaasCertificationMsg = '인증이 필요합니다.';
              }
            },
            function (reason) {
              $timeout(function () {
                $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message});
              }, 500);
            }
          );
        } else {
          vm.memberInfo.iaasUserId = '';
          alert('IaaS 사용여부를 선택해주세요.');
        }
      } else if(val == "paasUser" && vm.memberInfo.paasUserId != "" && vm.memberInfo.paasUserId != undefined) {
        if(vm.memberInfo.paasUserUseYn == 'Y') {
          memberService.paasDuplicateCheckId(vm.memberInfo.paasUserId).then(
            function (result) {
              if (result.data == vm.memberInfo.paasUserId) {
                vm.memberInfo.paasUserChck = false;
                vm.paasCertificationMsg = '이미 사용중이 아이디 입니다.';
                vm.memberInfo.paasUserId = '';
              } else {
                vm.paasCertificationMsg = '인증이 필요합니다.';
              }
            },
            function (reason) {
              $timeout(function () {
                $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message});
              }, 500);
            }
          );
        } else {
          vm.memberInfo.paasUserId = '';
          alert('PaaS 사용여부를 선택해주세요.');
        }
      }
    };


    // Member Join Action
    vm.join = function() {
      vm.scope.loading = true;
      memberService.join(vm.memberInfo).then(
        function() {
          vm.alertModal('회원가입', '회원 가입이 완료 되었습니다. 로그인 화면으로 이동 합니다.', vm.movePage, true, false);
          vm.scope.loading = false;
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };


    // Member Info Update
    vm.save = function() {
      vm.scope.loading = true;
      memberService.save(vm.memberInfo).then(
        function(response) {
          var data = response.data;
          cache.setUser({
            name: data.username,
            email: data.userEmail,
            sysType: data.sysType,
            i1: data.authI1,
            i2: data.authI2,
            p1: data.authP1,
            p2: data.authP2
          }, constants.expire);

          vm.scope.loading = false;
          vm.alertModal('회원정보 수정', '회원 정보 수정이 완료 되었습니다.', vm.movePage, true, false);
        },
        function(reason) {
          vm.scope.loading = false;
          $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
        }
      );
    };


    // Remove IaaS/PaaS User ID/PW
    vm.userUseConfirm = function(type) {
      if(type == 'iaas') {
        if(vm.memberInfo.iaasUserUseYn == 'N') {
          vm.memberInfo.iaasUserId = "";
          vm.memberInfo.iaasUserPw = "";
          vm.memberInfo.iaasUserChck = false;
        }
      } else {
        if(vm.memberInfo.paasUserUseYn == 'N') {
          vm.memberInfo.paasUserId = "";
          vm.memberInfo.paasUserPw = "";
          vm.memberInfo.paasUserChck = false;
        }
      }
    };


    // Alert Modal
    vm.alertModal = function(title, message, callback, aShow, cShow) {
      vm.modalReset();
      vm.modalTitle = title;
      vm.modalMessage = message;
      vm.confirmCallback = callback;
      vm.alertShow = aShow;
      vm.closeShow = cShow;
      angular.element('#alertModal').modal('show');
    }


    // Move Page
    vm.movePage = function(url) {
      if(url != "" && url != null) {
        $location.path(url);
      } else {
        if(!vm.scope.type) {
          angular.element('.wrapper').css('margin-top', 70);
          angular.element($document[0].querySelector(".modal-backdrop")).removeClass("modal-backdrop");
          $location.path('/login');
        } else {
          $window.location.reload();
        }
      }
    };


    // Certification Confirm
    vm.certificationConfirm = function(certifyType, bool) {
      if(certifyType == "IaaS") {
        memberService.iaasCertificationConfirm(vm.memberInfo).then(
          function(response) {
            vm.scope.loading = false;
            if(response.data != "" && response.data != "unauthorized") {
              vm.memberInfo.iaasUserChck = true;
            } else {
              vm.memberInfo.iaasUserChck = false;
              if(!bool) {
                if(response.data == "unauthorized") {
                  vm.iaasCertificationMsg = "에 일치하는 사용자 정보가 존재 하지 않습니다. ID/PASSWORD를 확인하세요!"
                } else {
                  vm.iaasCertificationMsg = "인증에 실패 하였습니다.";
                }
              }
            }
          },
          function(reason) {
            vm.scope.loading = false;
            $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
          }
        );
      } else if(certifyType == "PaaS") {
        memberService.paasCertificationConfirm(vm.memberInfo).then(
          function(response) {
            vm.scope.loading = false;
            if(response.data != "" && response.data != "bad_credentials" && response.data != "account_locked" && response.data != "not_admin_account") {
              vm.memberInfo.paasUserChck = true;
            } else {
              vm.memberInfo.paasUserChck = false;
              if(!bool) {
                if(response.data == "bad_credentials") {
                  vm.paasCertificationMsg = "에 일치하는 사용자 정보가 존재 하지 않습니다. ID/PASSWORD를 확인하세요!"
                } else if(response.data == "account_locked") {
                  vm.paasCertificationMsg = "계정이 잠겨있습니다. 잠시 후 다시 시도하시기 바랍니다."
                } else if(response.data == "not_admin_account") {
                  vm.paasCertificationMsg = "관리자 권한이 없는 사용자 입니다."
                } else {
                  vm.paasCertificationMsg = "인증에 실패 하였습니다."
                }
              }
            }
          },
          function(reason) {
            vm.scope.loading = false;
            $timeout(function() { $exceptionHandler(reason.data.message, {code: reason.data.HttpStatus, message: reason.data.message}); }, 500);
          }
        );
      }
    };

    // Certify Check Change
    vm.certifyChange = function(type) {
      if(type == "IaaS") {
        vm.memberInfo.iaasUserChck = false;
      } else if(type == "PaaS") {
        vm.memberInfo.paasUserChck = false;
      }
    };


    // Modal Setting
    vm.confirmJoin = function() {
      vm.modalReset();
      vm.modalTitle = '회원 등록';
      if(vm.validation()) {
        vm.confirmShow = true;
        vm.modalMessage = '회원 등록을 하시겠습니까?';
        vm.confirmCallback = vm.join;
      } else {
        vm.confirmShow = false;
      }
    };
    vm.confirmSave = function() {
      vm.modalReset();
      vm.modalTitle = '회원정보 수정';
      if(vm.validation()) {
        vm.confirmShow = true;
        vm.modalMessage = '회원정보를 수정 하시겠습니까?';
        vm.confirmCallback = vm.save;
      } else {
        vm.confirmShow = false;
      }
    };
    vm.modalReset = function() {
      vm.modalTitle = '';
      vm.modalMessage = '';
      vm.confirmCallback = undefined;
    };


    // Validation Check
    vm.validation = function() {
      if(vm.scope.sysConfig == "IaaS") {
        if(vm.memberInfo.iaasUserUseYn == "N" || vm.memberInfo.iaasUserUseYn == undefined) {
          vm.modalMessage = 'IaaS 인증을 받으셔야 합니다.';
          return false;
        }
      } else if(vm.scope.sysConfig == "PaaS") {
        if(vm.memberInfo.paasUserUseYn == "N" || vm.memberInfo.paasUserUseYn == undefined) {
          vm.modalMessage = 'PaaS 인증을 받으셔야 합니다.';
          return false;
        }
      } else {
        if((vm.memberInfo.iaasUserUseYn == "N" || vm.memberInfo.iaasUserUseYn == undefined)
          && (vm.memberInfo.paasUserUseYn == "N" || vm.memberInfo.paasUserUseYn == undefined)) {
          vm.modalMessage = 'IaaS 또는 PaaS 인증을 받으셔야 합니다.';
          return false;
        } else {
          if(vm.memberInfo.iaasUserUseYn == "Y" && !vm.memberInfo.iaasUserChck) {
            vm.modalMessage = 'IaaS 아이디와 비밀번호 입력 후 인증을 해야합니다.';
            return false;
          } else if(vm.memberInfo.paasUserUseYn == "Y" && !vm.memberInfo.paasUserChck) {
            vm.modalMessage = 'PaaS 아이디와 비밀번호 입력 후 인증을 해야합니다.';
            return false;
          }
        }
      }

      return true;
    }

  }

})();

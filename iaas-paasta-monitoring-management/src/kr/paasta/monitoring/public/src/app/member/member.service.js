(function() {
  'use strict';

  angular
    .module('monitoring')
    .factory('memberService', MemberService);

  /** @ngInject */
  function MemberService($http, apiUris) {
    var service = {};

    service.init = function() {
      var config = {
        headers : {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      };
      return $http.get(apiUris.join, config);
    };

    service.memberInfoView = function(data) {
      var config = {
        headers : {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      };
      return $http.post(apiUris.memberInfoView, {"userid":data}, config);
    };

    service.duplicateConfirmId = function(data) {
      return $http.get(apiUris.joinCheckId.replace(":id", data));
    };

    service.iaasDuplicateCheckId = function(id) {
      var config = {
        headers : {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      };
      return $http.get(apiUris.joinCheckDuplicationIaas.replace(":id", id), config);
    };

    service.paasDuplicateCheckId = function(id) {
      var config = {
        headers : {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      };
      return $http.get(apiUris.joinCheckDuplicationPaas.replace(":id", id), config);
    };

    service.join = function(data) {
      var config = {
        headers : {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      };
      return $http.post(apiUris.join, data, config);
    };

    service.save = function(data) {
      var config = {
        headers : {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      };
      return $http.patch(apiUris.memberInfoSave, data, config);
    };

    service.iaasCertificationConfirm = function(data){
      var config = {
        headers : {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      };
      return $http.post(apiUris.joinCheckIaas, data, config);
    };

    service.paasCertificationConfirm = function(data){
      var config = {
        headers : {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      };
      return $http.post(apiUris.joinCheckPaas, data, config);
    };

    return service;
  }

})();

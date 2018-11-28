(function() {
  'use strict';

  angular
    .module('monitoring')
    .directive('modal', modal);

  /** @ngInject */
  function modal($timeout){
    return {
      templateUrl: 'app/component/modal/modal.html',
      restrict: 'E',
      transclude: true,
      replace: true,
      scope: {
        modalTitle: '=',
        modalVisible: '='
      },
      link: function(scope, element, attrs) {
        scope.modalTitle = attrs.modalTitle;

        $timeout(function() {
          scope.$watch('modalVisible', function(value){
            if(value == true)
              element.modal('show');
            else
              element.modal('hide');
          });
          scope.$watch('modalTitle', function(value){
            scope.modalTitle = value;
          });
        });

        element.on('shown.bs.modal', function(){
          scope.$apply(function(){
            scope.$parent[attrs.modalVisible] = true;
          });
        });

        element.on('hidden.bs.modal', function(){
          scope.$apply(function(){
            scope.$parent[attrs.modalVisible] = false;
          });
        });
      }
    };
  }

})();

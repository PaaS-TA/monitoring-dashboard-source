(function() {
  'use strict';

  angular
    .module('monitoring')
    .directive('loading', loading);

  /** @ngInject */
  function loading($window, $interval) {
    var directive = {
      restrict: 'A',
      scope: {
        isShown: '=',
        isModal: '='
      },
      link: linkFunc
    };
    return directive;

    function linkFunc(scope, elem) {
      var stop = null;
      scope.$watch('isShown', function(newValue, oldValue) {
        if(scope.isModal != true && newValue) {
          elem.css('top', (($window.innerHeight * 0.4) + $window.pageYOffset) + 'px');
          elem.parent().css('margin-top', '-20px');

          stop = $interval(function() {
            if($window.innerHeight < elem.parent().parent()[0].scrollHeight) {
              elem.parent().height(elem.parent().parent()[0].scrollHeight);
            } else {
              elem.parent().height($window.innerHeight-51);
            }
          }, 10, 10);
        }

        if(oldValue != newValue && oldValue == true) {
          $interval.cancel(stop);
        }
      });
    }
  }
})();

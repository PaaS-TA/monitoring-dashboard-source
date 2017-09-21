angular.module('simplePagination', [])
    .constant('simplePaginationConfig', {
        totalItems: 1, // 전체 목록 건수
        currentPage: 1, // 현재 페이지
        itemsPerPage: 5 // 페이지 당 목록 건수
    })
    .controller('SimplePaginationController', ['$scope', 'simplePaginationConfig',
        function ($scope, simplePaginationConfig) {
            
            this.setNavigator = function(obj) {
                var nodes = $(".portal-pagination-block").children();
                nodes.each(function(){
                    setEnabled($(this));
                });
                if(obj.currentPage == 1){
                    setDisabled($(".start"));
                    setDisabled($(".left"));
                }
                if(obj.currentPage == obj.totalPages){
                    setDisabled($(".right"));
                    setDisabled($(".end"));
                }
            };
            function setEnabled(obj) {
                $(obj).removeClass("not-allowed-cursor").addClass("cursor");
            };
            function setDisabled(obj) {
                $(obj).removeClass("cursor").addClass("not-allowed-cursor");
            };
            
        }])
    .directive("simplePagination", function(simplePaginationConfig) {
        return {
            restrict: 'E',
            controller: 'SimplePaginationController',
            scope: {
                totalItems: '=',
                currentPage: '=',
                itemsPerPage: '=',
                totalPages: '=',
                getList: '&'
            },
            templateUrl: 'partials/pagination.html',
            replace: true,
            link: function(scope, element, attrs, ctrl) {
                scope.$watch('totalItems', function(newValue, oldValue) {
                    if (newValue == oldValue)
                        return;
                    //var totalPages = Math.ceil(scope.totalItems / scope.itemsPerPage);
                    //scope.totalPages = totalPages;
                    ctrl.setNavigator(scope);
                }, true);

                scope.firstPage = function () {
                    scope.currentPage = 1;
                    ctrl.setNavigator(scope);
                    scope.getList(scope);
                };
                scope.prevPage = function () {
                    if(scope.currentPage == 1) return;
                    scope.currentPage--;
                    ctrl.setNavigator(scope);
                    scope.getList(scope);
                };
                scope.nextPage = function () {
                    if(scope.currentPage == scope.totalPages) return;
                    scope.currentPage++;
                    ctrl.setNavigator(scope);
                    scope.getList(scope);
                };
                scope.lastPage = function () {
                    scope.currentPage = scope.totalPages;
                    ctrl.setNavigator(scope);
                    scope.getList(scope);
                };
            }   
        }
    });
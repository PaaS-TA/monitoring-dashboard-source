(function() {
  'use strict';

  angular
    .module('monitoring')
    .config(routerConfig);

  /** @ngInject */
  function routerConfig($stateProvider, $urlRouterProvider) {
    var navbar = {
      template: '<div><acme-navbar></acme-navbar></div>'
    };

    var iaasSidebar = {
      template: '<div><acme-iaas-sidebar></acme-iaas-sidebar></div>'
    };

    var paasSidebar = {
      template: '<div><acme-paas-sidebar></acme-paas-sidebar></div>'
    };

    $stateProvider
      .state('Main Dashboard', {
        url: '/',
        views: {
          'navbar': navbar,
          'body'  : {
            templateUrl: 'app/dashboard/dashboard.html',
            controller: 'DashboardController',
            controllerAs: 'dashboard'
          }
        }
      });

    $stateProvider
      .state('login', {
        url: '/login',
        views: {
          'login': {
            templateUrl: 'app/login/login.html',
            controller: 'LoginController',
            controllerAs: 'login'
          }
        }
      });

    $stateProvider
      .state('member_join', {
        url: '/member/join',
        views: {
          'body': {
            templateUrl: 'app/member/member.html',
            controller: 'MemberController',
            controllerAs: 'member'
          }
        }
      });

    $stateProvider
      .state('member_info', {
        url: '/member/info',
        views: {
          'navbar': navbar,
          'body': {
            templateUrl: 'app/member/member.html',
            controller: 'MemberController',
            controllerAs: 'member'
          }
        }
      });

    /* ----------------------------------- IaaS ----------------------------------- */

    $stateProvider
      .state('iaas_main', {
        url: '/iaas/main',
        views: {
          'navbar': navbar,
          'sidebar': iaasSidebar,
          'body': {
            templateUrl: 'app/iaas/main/main.html',
            controller: 'IaaSMainController',
            controllerAs: 'iaasMain'
          }
        }
      });

    $stateProvider
      .state('compute_node', {
        url: '/iaas/compute_node',
        views: {
          'navbar': navbar,
          'sidebar': iaasSidebar,
          'body': {
            templateUrl: 'app/iaas/node/compute_node.html',
            controller: 'ComputeNodeController',
            controllerAs: 'cnd'
          }
        }
      });

    $stateProvider
      .state('compute_node_detail', {
        url: '/iaas/compute_node/:hostname',
        views: {
          'navbar': navbar,
          'sidebar': iaasSidebar,
          'body': {
            templateUrl: 'app/iaas/node/node_detail.html',
            controller: 'NodeDetailController'
          }
        }
      });

    $stateProvider
      .state('manage_node', {
        url: '/iaas/manage_node',
        views: {
          'navbar': navbar,
          'sidebar': iaasSidebar,
          'body': {
            templateUrl: 'app/iaas/node/manage_node.html',
            controller: 'ManageNodeController',
            controllerAs: 'mnd'
          }
        }
      });

    $stateProvider
      .state('manage_node_detail', {
        url: '/iaas/manage_node/:hostname',
        views: {
          'navbar': navbar,
          'sidebar': iaasSidebar,
          'body': {
            templateUrl: 'app/iaas/node/node_detail.html',
            controller: 'NodeDetailController'
          }
        }
      });

    $stateProvider
      .state('tenant', {
        url: '/iaas/tenant',
        views: {
          'navbar': navbar,
          'sidebar': iaasSidebar,
          'body': {
            templateUrl: 'app/iaas/tenant/tenant.html',
            controller: 'TenantController',
            controllerAs: 'tnt'
          }
        }
      });

    $stateProvider
      .state('tenant_detail', {
        url: '/iaas/tenant/:instanceId',
        views: {
          'navbar': navbar,
          'sidebar': iaasSidebar,
          'body': {
            templateUrl: 'app/iaas/tenant/tenant_detail.html',
            controller: 'TenantDetailController'
          }
        }
      });

    $stateProvider
      .state('alarm_notification', {
        url: '/iaas/alarm/notification',
        views: {
          'navbar': navbar,
          'sidebar': iaasSidebar,
          'body': {
            templateUrl: 'app/iaas/alarm/alarm_notification.html',
            controller: 'IaasAlarmNotificationController',
            controllerAs: 'aln'
          }
        }
      });

    $stateProvider
      .state('alarm_policy', {
        url: '/iaas/alarm/policy',
        views: {
          'navbar': navbar,
          'sidebar': iaasSidebar,
          'body': {
            templateUrl: 'app/iaas/alarm/alarm_policy.html',
            controller: 'IaasAlarmPolicyController',
            controllerAs: 'alp'
          }
        }
      });

    $stateProvider
      .state('alarm_policy_detail', {
        url: '/iaas/alarm/policy/:id',
        views: {
          'navbar': navbar,
          'sidebar': iaasSidebar,
          'body': {
            templateUrl: 'app/iaas/alarm/alarm_policy_detail.html',
            controller: 'AlarmPolicyDetailController',
            controllerAs: 'alp'
          }
        }
      });

    $stateProvider
      .state('alarm_status', {
        url: '/iaas/alarm/status',
        views: {
          'navbar': navbar,
          'sidebar': iaasSidebar,
          'body': {
            templateUrl: 'app/iaas/alarm/alarm_status.html',
            controller: 'IaasAlarmStatusController',
            controllerAs: 'ast'
          }
        }
      });

    $stateProvider
      .state('alarm_status_detail', {
        url: '/iaas/alarm/status/:id',
        views: {
          'navbar': navbar,
          'sidebar': iaasSidebar,
          'body': {
            templateUrl: 'app/iaas/alarm/alarm_status_detail.html',
            controller: 'AlarmStatusDetailController',
            controllerAs: 'ast'
          }
        }
      });

    /* ----------------------------------- PaaS ----------------------------------- */

    $stateProvider
      .state('paas_main', {
        url: '/paas/main',
        views: {
          'navbar': navbar,
          'sidebar': paasSidebar,
          'body': {
            templateUrl: 'app/paas/main/main.html',
            controller: 'PaaSMainController',
            controllerAs: 'paasMain'
          }
        }
      });

    $stateProvider
      .state('paas_bosh', {
        url: '/paas/bosh',
        views: {
          'navbar': navbar,
          'sidebar': paasSidebar,
          'body': {
            templateUrl: 'app/paas/bosh/bosh.html',
            controller: 'PaasBoshController',
            controllerAs: 'pbsh'
          }
        }
      });

    $stateProvider
      .state('paas_bosh_detail', {
        url: '/paas/bosh/:id',
        views: {
          'navbar': navbar,
          'sidebar': paasSidebar,
          'body': {
            templateUrl: 'app/paas/bosh/bosh_detail.html',
            controller: 'PaasBoshDetailController',
            controllerAs: 'pbsh'
          }
        }
      });

    $stateProvider
      .state('paas_paasta', {
        url: '/paas/paasta',
        views: {
          'navbar': navbar,
          'sidebar': paasSidebar,
          'body': {
            templateUrl: 'app/paas/paasta/paasta.html',
            controller: 'PaasPaastaController',
            controllerAs: 'ppst'
          }
        }
      });

    $stateProvider
      .state('paas_paasta_detail', {
        url: '/paas/paasta/:id',
        views: {
          'navbar': navbar,
          'sidebar': paasSidebar,
          'body': {
            templateUrl: 'app/paas/paasta/paasta_detail.html',
            controller: 'PaasPaastaDetailController',
            controllerAs: 'ppst'
          }
        }
      });

    $stateProvider
      .state('paas_container', {
        url: '/paas/container',
        views: {
          'navbar': navbar,
          'sidebar': paasSidebar,
          'body': {
            templateUrl: 'app/paas/container/container.html',
            controller: 'PaasContainerController',
            controllerAs: 'pctn'
          }
        }
      });

    $stateProvider
      .state('paas_container_detail', {
        url: '/paas/container/:id',
        views: {
          'navbar': navbar,
          'sidebar': paasSidebar,
          'body': {
            templateUrl: 'app/paas/container/container_detail.html',
            controller: 'PaasContainerDetailController',
            controllerAs: 'pctn'
          }
        }
      });

    $stateProvider
      .state('paas_alarm_policy', {
        url: '/paas/alarm/policy',
        views: {
          'navbar': navbar,
          'sidebar': paasSidebar,
          'body': {
            templateUrl: 'app/paas/alarm/alarm_policy.html',
            controller: 'PaasAlarmPolicyController',
            controllerAs: 'papc'
          }
        }
      });

    $stateProvider
      .state('paas_alarm_status', {
        url: '/paas/alarm/status',
        views: {
          'navbar': navbar,
          'sidebar': paasSidebar,
          'body': {
            templateUrl: 'app/paas/alarm/alarm_status.html',
            controller: 'PaasAlarmStatusController',
            controllerAs: 'pasu'
          }
        }
      });

    $stateProvider
      .state('paas_alarm_status_detail', {
        url: '/paas/alarm/status/:id',
        views: {
          'navbar': navbar,
          'sidebar': paasSidebar,
          'body': {
            templateUrl: 'app/paas/alarm/alarm_status_detail.html',
            controller: 'PaasAlarmStatusDetailController',
            controllerAs: 'pasu'
          }
        }
      });

    $stateProvider
      .state('paas_alarm_statistics', {
        url: '/paas/alarm/statistics',
        views: {
          'navbar': navbar,
          'sidebar': paasSidebar,
          'body': {
            templateUrl: 'app/paas/alarm/alarm_statistics.html',
            controller: 'PaasAlarmStatisticsController',
            controllerAs: 'past'
          }
        }
      });

    $urlRouterProvider.otherwise('/');
  }

})();

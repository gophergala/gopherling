require './modules/app'

# define our configuration class
class Config
  constructor: ($routeProvider, $locationProvider) ->
    $locationProvider.html5Mode false

    $routeProvider
    .when '/',
      templateUrl: 'views/home.html'

    .when '/new',
      templateUrl: 'views/new_test.html'
      controller: 'NewTestController'

    .when '/tests',
      templateUrl: 'views/all_tests.html'
      controller: 'AllTestsController'

    .when '/tests/:id',
      templateUrl: 'views/start_test.html'
      controller: 'StartTestController'

    .when '/tests/:id/edit',
      templateUrl: 'views/new_test.html'
      controller: 'NewTestController'

    .otherwise
      redirectTo: '/'

# register our configuration class
angular.module('app').config ['$routeProvider', '$locationProvider', Config]

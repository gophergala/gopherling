require './modules/app'

# define our configuration class
class Config
  constructor: ($routeProvider, $locationProvider) ->
    $locationProvider.html5Mode false

    $routeProvider
    .when '/',
      templateUrl: 'views/home.html'

    .when '/new',
      templateUrl: 'views/new.html'

    .when '/tests',
      templateUrl: 'views/tests.html'

    .when '/settings',
      templateUrl: 'views/settings.html'

    .otherwise
      redirectTo: '/'

# register our configuration class
angular.module('app').config ['$routeProvider', '$locationProvider', Config]

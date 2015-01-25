require 'angular'
require 'angular-route'
require 'angular-websocket'

require './app.controllers'

app = angular.module 'app', ['ngRoute', 'ngWebSocket', 'app.controllers']

module.exports = app

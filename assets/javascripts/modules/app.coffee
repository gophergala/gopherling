require 'angular'
require 'angular-route'

require './app.controllers'

app = angular.module 'app', ['ngRoute', 'app.controllers']

module.exports = app

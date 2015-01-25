_ = require 'underscore'

module.exports = class
  @$inject: ['$scope', '$http', '$routeParams', '$location', '$websocket']
  constructor: (@scope, @http, @params, @location, @socket) ->
    @scope.tests = []

    @http.get '/api/tests'
    .success (res) =>
      @scope.tests = res

    angular.extend @scope,
      delete: @delete

  delete: (id) =>
    @http.delete '/api/tests/'+id
    .success (res) =>
      @scope.tests = _(@scope.tests).reject (test) ->
        id is test._id

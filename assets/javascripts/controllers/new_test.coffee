_ = require 'underscore'

module.exports = class
  @$inject: ['$scope', '$http', '$location']
  constructor: (@scope, @http, @location) ->
    @scope.test =
      name: ''
      description: ''
      baseUrl: ''
      requests: 0
      concurrency: 0
      tasks: []

    @addTask()

    angular.extend @scope,
      save: @save
      addTask: @addTask
      removeTask: @removeTask
      addHeader: @addHeader
      removeHeader: @removeHeader

  addTask: () =>
    @scope.test.tasks.push
      method: 'GET'
      path: ''
      headers: []
      rawBody: ''

  removeTask: (task) =>
    @scope.test.tasks = _(@scope.test.tasks).reject (t) ->
      t is task

  addHeader: (task) =>
    task.headers.push
      field: ''
      value: ''

  removeHeader: (task, header) =>
    task.headers = _(task.headers).reject (h) ->
      h is header

  save: (run = false) =>
    console.log @scope.test
    @http.post '/api/tests', @scope.test
    .success (res) =>
      if run is true
        @location.path '/tests/'+res._id
      else
        @location.path '/tests'
      console.debug 'Your test has been saved'
    .error (err) =>
      console.debug 'An error occured'

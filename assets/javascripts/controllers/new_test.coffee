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

  addTask: () =>
    @scope.test.tasks.push
      method: 'GET'
      path: ''

  save: (run = false) =>
    @http.post '/api/tests', @scope.test
    .success (res) =>
      if run is true
        @location.path '/tests/'+res._id
        console.debug 'Test ('+res._id+') will be started'
      else
        @location.path '/tests'
      console.debug 'Your test has been saved'
    .error (err) =>
      console.debug 'An error occured'

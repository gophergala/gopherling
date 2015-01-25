module.exports = class
  @$inject: ['$scope', '$http', '$routeParams', '$websocket']
  constructor: (@scope, @http, @params, @socket) ->
    @scope.test =
      name: ''
      tasks: []

    @scope.total =
      requests: 0
      success: 0
      failures: 0
      min: 0
      mean: 0
      max: 0
      rps: 0

    @http.get '/api/tests/'+@params.id
    .success (res) =>
      @scope.test = res

      for task in @scope.test.tasks
        task.requests = 0
        task.success = 0
        task.failures = 0
        task.min = 0
        task.mean = 0
        task.max = 0
        task.rps = 0

      @stream = @socket 'ws://127.0.0.1:9410/api/tests/'+@params.id+'/start'

      @stream.onMessage (message) =>
        data = angular.fromJson(message.data)

        # We increment the number of requests
        @scope.test.tasks[data.task].requests++
        @scope.total.requests++

        # Is the request a success or a failure?
        if data.statusCode isnt 0
          @scope.test.tasks[data.task].success++
          @scope.total.success++
        else
          @scope.test.tasks[data.task].failures++
          @scope.total.failures++

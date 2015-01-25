module.exports = class
  @$inject: ['$scope', '$http', '$routeParams', '$location', '$websocket']
  constructor: (@scope, @http, @params, @location, @socket) ->
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
      rps: '..'

    @start = null

    @host = @location.host()
    @host = @host + ':' + @location.port() if @location.port() isnt 80 and @location.port() isnt 443

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
        task.rps = '..'

      @stream = @socket 'ws://'+@host+'/api/tests/'+@params.id+'/start'

      @stream.onOpen () =>
        @start = Date.now()

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

        if data.duration < @scope.test.tasks[data.task].min or @scope.test.tasks[data.task].min is 0
          @scope.test.tasks[data.task].min = data.duration

        if data.duration > @scope.test.tasks[data.task].max or @scope.test.tasks[data.task].max is 0
          @scope.test.tasks[data.task].max = data.duration

        @scope.test.tasks[data.task].mean = (@scope.test.tasks[data.task].min + @scope.test.tasks[data.task].max) / 2

      @stream.onClose () =>
        if @start?
          time = (Date.now() - @start) / 1000
          @scope.total.rps = 0
          for task in @scope.test.tasks
            task.rps = (task.requests / time).toFixed(2)
            @scope.total.rps += task.requests / time
          @scope.total.rps = @scope.total.rps.toFixed(2)
          @scope.$apply()

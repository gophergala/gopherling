module.exports = class
  @$inject: ['$scope', '$http', '$routeParams', '$websocket']
  constructor: (@scope, @http, @params, @socket) ->
    @stream = @socket 'ws://127.0.0.1:9410/api/tests/'+@params.id+'/start'

    @stream.onMessage (message) =>
      console.log message.data

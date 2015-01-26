app = angular.module 'app.controllers', []

app.controller 'NewTestController', require '../controllers/new_test'
app.controller 'EditTestController', require '../controllers/edit_test'
app.controller 'AllTestsController', require '../controllers/all_tests'
app.controller 'StartTestController', require '../controllers/start_test'

module.exports = app

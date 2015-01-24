gulp       = require 'gulp'
connect    = require 'gulp-connect'
stylus     = require 'gulp-stylus'
prefix     = require 'gulp-autoprefixer'
browserify = require 'browserify'
source     = require 'vinyl-source-stream'

jeet       = require 'jeet'
citrine    = require 'citrine'

handleError = (err) ->
  console.log err.toString()
  this.emit('end')

gulp.task 'js', () ->
  browserify {entries: ['./assets/javascripts/index.coffee'], extensions: ['.coffee'], debug : false}
  .transform 'coffeeify'
  .bundle()
  .on 'error', handleError
  .pipe source 'app.js'
  .pipe gulp.dest './static/javascripts'
  .pipe connect.reload()

gulp.task 'css', () ->
  gulp.src './assets/stylesheets/*.styl'
  .pipe stylus
    use: [jeet(), citrine()]
  .pipe prefix "last 4 versions", "> 1%", "ie 9", "ie 8", { cascade: true }
  .pipe gulp.dest './static/stylesheets'
  .pipe connect.reload()

gulp.task 'fonts', () ->
  gulp.src './assets/fonts/**/*'
  .pipe gulp.dest './static/fonts'
  .pipe connect.reload()

gulp.task 'html', () ->
  gulp.src './assets/**/*.html'
  .pipe gulp.dest './static'
  .pipe connect.reload()

gulp.task 'webserver', () ->
  connect.server
    root: 'static'
    livereload: true
    port: 4000

gulp.task 'watch', () ->
    gulp.watch 'assets/javascripts/**/*.coffee', ['js']
    gulp.watch 'assets/stylesheets/**/*.styl', ['css']
    gulp.watch 'assets/fonts/**/*', ['fonts']
    gulp.watch 'assets/**/*', ['html']

gulp.task 'build', ['js', 'css', 'fonts', 'html']
gulp.task 'server', ['build', 'webserver', 'watch']

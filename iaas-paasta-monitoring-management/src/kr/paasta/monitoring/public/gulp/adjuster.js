'use strict';

var path = require('path');
var gulp = require('gulp');
var conf = require('./conf');
var del = require('del');
var urlAdjuster = require('gulp-css-url-adjuster');
var cssmin = require('gulp-cssmin');

gulp.task('url:adjuster', function() {
  return gulp.src(path.join(conf.paths.dist, '/styles/app-*.css'))
    .pipe(urlAdjuster({
      prepend: '/public/dist'
    }))
    .pipe(gulp.dest(path.join(conf.paths.tmp, '/styles/')));
});

gulp.task('del:css', function () {
  return del([path.join(conf.paths.dist, '/styles/app-*.css')]);
});

gulp.task('css:compress', ['del:css', 'url:adjuster'], function() {
  gulp.src(path.join(conf.paths.tmp, '/styles/app-*.css'))
    .pipe(cssmin())
    .pipe(gulp.dest(path.join(conf.paths.dist, '/styles')));
});

gulp.task('css:rebuild', ['css:compress'], function () {
  return del([path.join(conf.paths.tmp, '/styles')]);
});

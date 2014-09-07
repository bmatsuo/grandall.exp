package main

import "net/http"

var indexRaw = []byte(`
<!DOCTYPE html>
<html lang="en" ng-app="grandall">
<head>
	<meta charset="utf-8">
	<title>Grandall</title>

	<!-- Latest Bootstrap CSS -->
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.2.0/css/bootstrap.min.css">
	<!-- Latest Angular -->
	<script src="https://ajax.googleapis.com/ajax/libs/angularjs/1.3.0-rc.0/angular.min.js"></script>
</head>
<body ng-controller="AliasListCtrl">
	<div class="container">
		<h1 class="row">Aliases</h1>
		<ul class="list-unstyled row">
			<li ng-repeat="alias in aliases">
			<a href={{alias.url}}><strong>{{alias.name}}</string></a>
			<span>{{alias.description}}</span>
			</li>
		</ul>
	</div>
	<script src="/static/app.js"></script>
</body>
</html>
`)

var appJSRaw = []byte(`
var grandallApp = angular.module('grandall', []);

grandallApp.controller('AliasListCtrl', ['$scope', '$http', function($scope, $http) {
	$scope.aliases = [];
	$http.get('/.api/v1/aliases').success(function(data) {
		$scope.aliases = data;
	});
}]);
`)

func UI(s []*Site) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/static/app.js", staticHandler("text/javascript", appJSRaw))
	mux.Handle("/", staticHandler("text/html", indexRaw))
	return mux
}

func staticHandler(mtype string, p []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body.Close()
		if mtype != "" {
			w.Header().Set("Content-Type", mtype)
		}
		w.Write(p)
	})
}

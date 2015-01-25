Gopherling
=========

Gopherling is a stress/load testing webapp similar to ab for the gophergala 2015.

Prerequisites
--------------

- Go
- mongoDB
- node env (if you want to edit the js or css)

Installation
--------------

```sh
$ cd $GOTPATH
$ mkdir -p src/github.com/gophergala
$ cd src/github.com/gophergala
$ git clone https://github.com/gophergala/gopherling.git
$ cd gopherling
$ go get ./...
```

Starting the server
--------------

```sh
$ cd $GOPATH/src/github.com/gophergala/gopherling
$ go install && gopherling --dbHost=127.0.0.1 --dbPort=27017
```

Using the app
--------------

All you have to do is open this url on your browser `http://127.0.0.1:9410`

License
----

Gopherling is under the GNU GPL v3.0 license.

Made by Frostwind <hi@frostwind.me>

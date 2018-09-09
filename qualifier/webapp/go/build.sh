#!/bin/bash

go get github.com/go-martini/martini
go get github.com/go-sql-driver/mysql
go get github.com/martini-contrib/render
go get github.com/martini-contrib/sessions
go get github.com/gomodule/redigo/redis
go build -o golang-webapp .

go run ../redis/init.go
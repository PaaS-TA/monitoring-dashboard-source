#!/usr/bin/env bash

export GOPATH=$PWD
export PATH=$GOPATH/bin:$PATH

#Mysql Driver and Orm Library Download
go get github.com/go-sql-driver/mysql
go get github.com/jinzhu/gorm

go get github.com/thoas/go-funk
go get gopkg.in/gomail.v2
go get github.com/go-telegram-bot-api/telegram-bot-api
go get github.com/jinzhu/gorm
go get github.com/mileusna/crontab
go get github.com/thoas/go-funk
go get github.com/tidwall/gjson

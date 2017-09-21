#!/usr/bin/env bash

. .envrc
#go get github.com/cloudfoundry-incubator/runtime-schema
go get github.com/tedsuo/ifrit

#Mysql Driver and Orm Library Download
go get github.com/go-sql-driver/mysql
go get github.com/jinzhu/gorm

#InfluxDB Library Download
go get github.com/influxdata/influxdb/client/v2

#Bosh Client
go get github.com/cloudfoundry-community/gogobosh
go get golang.org/x/oauth2
go get golang.org/x/net/context

#Test Code
go get github.com/onsi/ginkgo
go get github.com/onsi/gomega

#godep
go get github.com/tools/godep
go get golang.org/x/sys/unix
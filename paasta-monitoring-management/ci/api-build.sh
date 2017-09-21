#!/usr/bin/env bash

cd git
echo $(pwd)

cd paasta-monitoring-management

export GOPATH=$PWD
export PATH=$GOPATH/bin:$PATH


#go get github.com/cloudfoundry-incubator/runtime-schema
go get github.com/tedsuo/ifrit

#rata
go get github.com/tedsuo/rata

#Mysql Driver and Orm Library Download
go get github.com/go-sql-driver/mysql
go get github.com/jinzhu/gorm

#InfluxDB Library Download
go get github.com/influxdata/influxdb/client/v2


#TestCode
go get github.com/gorilla/handlers
go get github.com/gorilla/mux
go get github.com/stretchr/testify/assert
go get github.com/onsi/ginkgo
go get github.com/onsi/gomega

#godep
go get github.com/tools/godep
go get golang.org/x/sys/unix
go get github.com/davecgh/go-spew/spew
go get github.com/pmezard/go-difflib/difflib


go build src/kr/paasta/monitoring/main.go

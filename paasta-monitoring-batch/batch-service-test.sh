#!/usr/bin/env bash

export GOPATH=$PWD
export PATH=$GOPATH/bin:$PATH:.

go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega


cd src/kr/paasta/monitoring/monit-batch/services/

go test -v services_suite_test.go backend_services_test.go boshAlarmCollector_test.go

#go test -v services_suite_test.go  autoScaler_test.go
#go test autoScaler_test.go
#go test backend_services_test.go
#go test boshAlarmCollector_test.go
#go test BoshVmsInfoCollector_test.go
#go test containerAlarmCollector_test.go
#go test createSchema_test.go
#go test paastaAlarmCollector_test.go

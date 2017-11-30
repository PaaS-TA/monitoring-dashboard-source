#!/usr/bin/env bash

cd git
echo $(pwd)

cd paasta-monitoring-management

export GOPATH=$PWD
export PATH=$GOPATH/bin:$PATH


go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega


cd src/kr/paasta/monitoring/controller

go test -coverprofile=coverage.out
go tool cover -func=coverage.out
#go tool cover -html=coverage.out
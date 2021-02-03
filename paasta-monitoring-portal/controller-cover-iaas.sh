#!/usr/bin/env bash

export GOPATH=$PWD
export PATH=$GOPATH/bin:$PATH

go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega

echo "=================================== test start ========================================"
cd src/kr/paasta/monitoring/iaas/controller

go test -coverprofile=coverage.out
go tool cover -html=coverage.out
#go tool cover -func=coverage.out

echo "=================================== test end ==============================================="

#!/usr/bin/env bash

export GOPATH=$PWD
export PATH=$GOPATH/bin:$PATH:.

go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega

cd src/kr/paasta/monitoring/

go test -v ./...


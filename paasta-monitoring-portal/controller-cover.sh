#!/usr/bin/env bash

export GOPATH=$PWD
export PATH=$GOPATH/bin:$PATH

go get src/github.com/onsi/ginkgo/ginkgo
go get src/github.com/onsi/gomega

cd src/kr/paasta/monitoring/test

go test -coverprofile=coverage.out
go tool cover -html=coverage.out
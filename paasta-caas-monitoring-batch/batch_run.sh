#!/usr/bin/env bash

export GOPATH=$PWD
export PATH=$GOPATH/bin:$PATH

cd src/kr/paasta/batch
go run main.go
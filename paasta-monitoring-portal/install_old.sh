#!/usr/bin/env bash

. .envrc
go get github.com/tedsuo/ifrit

#rata
go get github.com/tedsuo/rata

#InfluxDB Library Download
go get github.com/influxdata/influxdb1-client/v2

#go client for openstack
go get github.com/rackspace/gophercloud

#go client for redis
go get github.com/go-redis/redis

#Mysql Driver and Orm Library Download
go get github.com/go-sql-driver/mysql
go get github.com/jinzhu/gorm
go get github.com/cihub/seelog
go get github.com/monasca/golang-monascaclient/monascaclient
go get github.com/gophercloud/gophercloud/
go get github.com/alexedwards/scs

# elastic search
go get gopkg.in/olivere/elastic.v3

#Test Code
go get github.com/onsi/ginkgo
go get github.com/onsi/gomega
go get github.com/stretchr/testify

# BOSH client
go get github.com/cloudfoundry-community/gogobosh

# Functional lib
go get github.com/thoas/go-funk

# json parser
go get github.com/tidwall/gjson

#monasca client Bug Fix Src
cp ./lib-bugFix-src/alarm_definitions.go ./src/github.com/monasca/golang-monascaclient/monascaclient
cp ./lib-bugFix-src/notifications.go ./src/github.com/monasca/golang-monascaclient/monascaclient
cp ./lib-bugFix-src/alarms.go ./src/github.com/monasca/golang-monascaclient/monascaclient

cp ./lib-bugFix-src/monascaclient/client.go ./src/github.com/monasca/golang-monascaclient/monascaclient
cp ./lib-bugFix-src/gophercloud/requests.go ./src/github.com/rackspace/gophercloud/openstack/identity/v3/tokens
cp ./lib-bugFix-src/gophercloud/results.go ./src/github.com/rackspace/gophercloud/openstack/identity/v3/tokens
cp ./lib-bugFix-src/gophercloud/client.go ./src/github.com/rackspace/gophercloud/openstack
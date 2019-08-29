#!/usr/bin/env bash
set GOPATH=D:\paas_iaas\workspace\PaaS-TA-Monitoring-4.0\paasta-monitoring-batch
set PATH=%PATH%;%GOPATH%bin

// Dependencies Module Download
go get github.com/tedsuo/ifrit
go get github.com/tedsuo/rata
go get github.com/influxdata/influxdb1-client/v2
go get github.com/rackspace/gophercloud
go get github.com/go-redis/redis
go get github.com/go-sql-driver/mysql
go get github.com/jinzhu/gorm
go get github.com/cihub/seelog
go get github.com/monasca/golang-monascaclient/monascaclient
go get github.com/gophercloud/gophercloud/
go get github.com/alexedwards/scs
go get gopkg.in/olivere/elastic.v3
go get github.com/onsi/ginkgo
go get github.com/onsi/gomega
go get github.com/stretchr/testify
go get github.com/cloudfoundry-community/gogobosh
go get github.com/thoas/go-funk
go get github.com/tidwall/gjson
# Functional lib
go get github.com/thoas/go-funk

# json parser
go get github.com/tidwall/gjson

xcopy .\lib-bugFix-src\alarm_definitions.go .\src\github.com\monasca\golang-monascaclient\monascaclient
xcopy .\lib-bugFix-src\notifications.go .\src\github.com\monasca\golang-monascaclient\monascaclient
xcopy .\lib-bugFix-src\alarms.go .\src\github.com\monasca\golang-monascaclient\monascaclient

xcopy .\lib-bugFix-src\monascaclient\client.go .\src\github.com\monasca\golang-monascaclient\monascaclient
xcopy .\lib-bugFix-src\gophercloud\requests.go .\src\github.com\rackspace\gophercloud\openstack\identity\v3\tokens
xcopy .\lib-bugFix-src\gophercloud\results.go .\src\github.com\rackspace\gophercloud\openstack\identity\v3\tokens
xcopy .\lib-bugFix-src\gophercloud\client.go .\src\github.com\rackspace\gophercloud\openstack
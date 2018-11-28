#!/usr/bin/env bash
set GOPATH=D:\paas_iaas\workspace\PaaS-TA-Monitoring-4.0\paasta-monitoring-batch
set PATH=%PATH%;%GOPATH%bin

// Dependencies Module Download
go get github.com/tedsuo/ifrit
go get github.com/go-sql-driver/mysql
go get github.com/jinzhu/gorm
go get github.com/influxdata/influxdb/client/v2
go get github.com/cloudfoundry-community/gogobosh
go get golang.org/x/oauth2
go get golang.org/x/net/context
go get github.com/onsi/ginkgo
go get github.com/onsi/gomega
go get github.com/tools/godep
go get golang.org/x/sys/unix
go get github.com/go-telegram-bot-api/telegram-bot-api
go get github.com/go-redis/redis

xcopy .\lib-bugFix-src\alarm_definitions.go .\src\github.com\monasca\golang-monascaclient\monascaclient
xcopy .\lib-bugFix-src\notifications.go .\src\github.com\monasca\golang-monascaclient\monascaclient
xcopy .\lib-bugFix-src\alarms.go .\src\github.com\monasca\golang-monascaclient\monascaclient

xcopy .\lib-bugFix-src\monascaclient\client.go .\src\github.com\monasca\golang-monascaclient\monascaclient
xcopy .\lib-bugFix-src\gophercloud\requests.go .\src\github.com\rackspace\gophercloud\openstack\identity\v3\tokens
xcopy .\lib-bugFix-src\gophercloud\results.go .\src\github.com\rackspace\gophercloud\openstack\identity\v3\tokens
xcopy .\lib-bugFix-src\gophercloud\client.go .\src\github.com\rackspace\gophercloud\openstack
xcopy .\lib-bugFix-src\go-cfclient\client.go .\src\github.com\cloudfoundry-community\go-cfclient
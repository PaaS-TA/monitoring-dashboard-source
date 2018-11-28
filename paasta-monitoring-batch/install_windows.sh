#!/usr/bin/env bash
set GOPATH=D:\paas_iaas\workspace\PaaS-TA-Monitoring-4.0\iaas-paasta-monitoring-management
set PATH=%PATH%;%GOPATH%bin

//DependenciesModuleDownload
go get github.com/tedsuo/ifrit
go get github.com/tedsuo/rata
go get github.com/influxdata/influxdb/client/v2
go get github.com/rackspace/gophercloud
go get github.com/cloudfoundry-community/go-cfclient
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
package main

import (
	"caas-monitoring-batch/caas"
	config2 "caas-monitoring-batch/config"
	"caas-monitoring-batch/dao"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"sync"
)

func main() {
	config := config2.GetConfiguration()

	connectionStr := config2.GetDBConnectionStr()
	dbConn, err := gorm.Open("mysql", connectionStr)
	if err != nil {
		log.Fatalln("err::", err)
	}

	dao.CreateTable(dbConn)

	// CaaS 스케쥴 실행
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	//caasAlarms := funk.Filter(alarmInfos, func(x model.BatchAlarmInfo) bool {
	//	return x.ServiceType == "CaaS"
	//}).([]model.BatchAlarmInfo)

	caas.Startschedule(dbConn, config.CaasApiUrl)

	waitGroup.Wait()
}

package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/thoas/go-funk"
	"kr/paasta/batch/caas"
	"kr/paasta/batch/dao"
	"kr/paasta/batch/model"
	"kr/paasta/batch/util"
	"log"
	"os"
	"sync"
)

func main() {
	// 기본적인 프로퍼티 설정 정보 읽어오기
	config, err := util.ReadConfig(`config.ini`)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	dbAccessObj := dao.GetdbAccessObj()
	//dao.CreateTable(dbAccessObj)
	alarmInfos := dao.GetBatchAlarmInfo(dbAccessObj)

	// CaaS 스케쥴 실행
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	caasAlarms := funk.Filter(alarmInfos, func(x model.BatchAlarmInfo) bool {
		return x.ServiceType == "CaaS"
	}).([]model.BatchAlarmInfo)

	caas.Startschedule(dbAccessObj, caasAlarms, config["caas.monitoring.api.url"])

	waitGroup.Wait()
}

package main

import (
	"github.com/thoas/go-funk"
	"kr/paasta/batch/dao"
	"kr/paasta/batch/model"
	"kr/paasta/batch/saas"
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
	dao.CreateTable(dbAccessObj)
	alarmInfos := dao.GetBatchAlarmInfo(dbAccessObj)

	// SaaS 스케쥴 실행
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	saasAlarms := funk.Filter(alarmInfos, func(x model.BatchAlarmInfo) bool {
		return x.ServiceType == "SaaS"
	}).([]model.BatchAlarmInfo)

	saas.Startschedule(dbAccessObj, saasAlarms, config["saas.pinpoint.url"])

	waitGroup.Wait()
}

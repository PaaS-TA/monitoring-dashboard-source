package main

import (
	"kr/paasta/monitoring-batch/dao"
	"kr/paasta/monitoring-batch/saas"
	"kr/paasta/monitoring-batch/util"
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

	// SaaS 스케쥴 실행
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	saas.Startschedule(dbAccessObj, config["saas.pinpoint.url"])

	waitGroup.Wait()
}

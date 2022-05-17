package main

import (
	"log"
	"os"
	"saas-monitoring-batch/dao"
	"saas-monitoring-batch/saas"
	"saas-monitoring-batch/util"
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

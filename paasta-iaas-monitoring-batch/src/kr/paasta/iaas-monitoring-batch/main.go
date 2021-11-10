package main

import (
	"fmt"
	"kr/paasta/iaas-monitoring-batch/config"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/grouper"
	"github.com/tedsuo/ifrit/sigmon"
)

func main() {

	//debugserver.AddFlags(flag.CommandLine)
	//cflager.AddFlags(flag.CommandLine)
	//===============================================================
	//catch or finally

	defer func() { //catch or finally
		if err := recover(); err != nil { //catch
			fmt.Fprintf(os.Stderr, "Main: Exception: %v\n", err)
			os.Exit(1)
		}
	}()
	//===============================================================
	//logger, _ := cflager.New("paasta-iaas-monitoring-batch")

	var startTime time.Time
	//============================================

	//============================================
	// Channel for Singal Checkig
	sigs := make(chan os.Signal, 1)
	//Waiting to be notified
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println("returned signal:", sig)
		elapsed := time.Since(startTime)
		fmt.Println("# ElapsedTime in seconds:", elapsed)

		//When unexpected signal happens, defer function doesn't work.
		//So, go func has a role to be notified signal and do defer function execute
		os.Exit(0)
	}()
	//============================================



	pid := os.Getpid()
	f, err := os.Create(".pid")
	defer f.Close()
	if err != nil {
		log.Fatalln("Main: failt to create pid file.", err.Error())
		panic(err)
	}
	f.WriteString(strconv.Itoa(pid))
	//=======================================================================

	//logger.Info("##### process id :", lager.Data{"process_id ":pid})

	//fmt.Println("gmtTimeGapHour::::", gmtTimeGapHour)

	// 설정데이터 초기화
	configData := config.InitializeConfig()

	members := grouper.Members {
		{"monitoring_batch_processor", NewIfritRunner(configData) },
	}

	log.Println("#monitoring batch processor started")
	group := grouper.NewOrdered(os.Interrupt, members)
	monitor := ifrit.Invoke(sigmon.New(group))
	monit_err := <-monitor.Wait()

	if monit_err != nil {
		log.Fatalln("#Main: exited-with-failure", monit_err)
		panic(monit_err)
	}

	elapsed := time.Since(startTime)
	log.Println("#ElapsedTime in seconds:", map[string]interface{}{"elapsed_time": elapsed, })
	log.Println("#monitoring batch processor exited")
}
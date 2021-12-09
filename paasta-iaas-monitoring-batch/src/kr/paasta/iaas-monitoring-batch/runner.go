package main

import (
	"fmt"
	"kr/paasta/iaas-monitoring-batch/config"
	"os"
	"sync"
	"time"

	"github.com/tedsuo/ifrit"

	"kr/paasta/iaas-monitoring-batch/service"
)

type Runner struct {
	configData *config.Config
}

/**
	Ifrit : A process model for Go
 */
func NewIfritRunner(configParam *config.Config) ifrit.Runner {
	return &Runner {
		configData: configParam,
	}
}


/**
	ifrit의 Runner 인터페이스를 구현함 (Runner 인터페이스의 Run 메서드를 구현 - ifrit/runner.go 참고)
 */
func (runner *Runner) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	//catch or finally
	defer func() {
		//catch or finally
		if err := recover(); err != nil {
			fmt.Println(err)
			//catch
			os.Exit(0)
		}
	}()
	close(ready)

	alarmService := service.AlarmServiceBuilder(runner.configData)

	ticker := time.NewTicker(time.Duration(runner.configData.ExecuteInterval) * time.Second)

	alarmService.StopProcess()

	for {
		select {
		case <- ticker.C:

			// WaitGroup :
			var wait sync.WaitGroup
			wait.Add(1)

			go func(wait *sync.WaitGroup) {
				defer wait.Done()
				alarmService.RunScheduler()
			}(&wait)
			wait.Wait()
		}
	}
}


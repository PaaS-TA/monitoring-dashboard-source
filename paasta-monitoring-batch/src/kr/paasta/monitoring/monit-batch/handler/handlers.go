package handler

import (
	"os"
	"sync"
	"time"
	"github.com/tedsuo/ifrit"
	"kr/paasta/monitoring/monit-batch/services"

)

type backend_server struct {
	batchInterval   int
	gmtTimeGapHour  int64
	influxCon 	*services.InfluxConfig
	configDbCon     *services.DBConfig
	portalDbCon     *services.DBConfig
	boshCon         *services.BoshConfig
	config          map[string]string
	mailConfig      *services.MailConfig
	//thresholdConfig   *services.ThresholdConfig

}

func New(batchInterval int, gmtTimeGapHour int64, influxCon *services.InfluxConfig, configDbCon *services.DBConfig, portalDbCon *services.DBConfig,
         boshCon *services.BoshConfig, config map[string]string,  mailConfig *services.MailConfig) ifrit.Runner {
	return &backend_server{

		batchInterval: batchInterval,
		gmtTimeGapHour: gmtTimeGapHour,
		influxCon:   influxCon,
		configDbCon: configDbCon,
		portalDbCon: portalDbCon,
		boshCon:     boshCon,
		config:      config,
		mailConfig: mailConfig,
		//thresholdConfig: thresholdConfig,
	}
}


func (n *backend_server) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	//===============================================================
	//catch or finally
	defer func() { //catch or finally
		if err := recover(); err != nil { //catch
			os.Exit(0)
		}
	}()
	//===============================================================
	close(ready)

	backend_service := services.NewBackendServices( n.gmtTimeGapHour, n.influxCon, n.configDbCon, n.portalDbCon, n.boshCon,  n.mailConfig,  n.config)

	//최초 실행시 Bosh VM정보 동기화
	backend_service.CreateUpdateBoshData(*n.boshCon)

	ticker := time.NewTicker(time.Duration(n.batchInterval) * time.Second)
	var index int
	backend_service.StopProcess()
	for {
		select {
		case <- ticker.C:

			go func(){
				backend_service.CreateUpdateBoshData(*n.boshCon)
			}()

			var wg sync.WaitGroup
			wg.Add(2)

			//임계치 초과시 Alarm 전송
			go func(wg *sync.WaitGroup){
				defer wg.Done()
				backend_service.CreateAlarmData()
			}(&wg)
			//AutoScale
			go func(wg *sync.WaitGroup){
				defer wg.Done()
				backend_service.AutoScale()
			}(&wg)

			wg.Wait()
		}
		index = index + 1
	}
}


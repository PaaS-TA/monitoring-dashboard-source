package handler

import (
	"os"
	"github.com/tedsuo/ifrit"
	"kr/paasta/monitoring-batch/service"
	"time"
	"sync"
)

type backend_server struct {
	batchInterval   int
	gmtTimeGapHour  int64
	influxCon 	*service.InfluxConfig
	configDbCon     *service.DBConfig
	portalDbCon     *service.DBConfig
	boshCon         *service.BoshConfig
	redisCon		*service.RedisConfig
	config          map[string]string
	mailConfig      *service.MailConfig
	//thresholdConfig   *services.ThresholdConfig
}

func New(batchInterval int, gmtTimeGapHour int64, influxCon *service.InfluxConfig, configDbCon *service.DBConfig, portalDbCon *service.DBConfig,
         boshCon *service.BoshConfig, redisCon *service.RedisConfig, config map[string]string,  mailConfig *service.MailConfig) ifrit.Runner {
	return &backend_server{

		batchInterval: batchInterval,
		gmtTimeGapHour: gmtTimeGapHour,
		influxCon:   influxCon,
		configDbCon: configDbCon,
		portalDbCon: portalDbCon,
		boshCon:     boshCon,
		redisCon:	redisCon,
		config:      config,
		mailConfig: mailConfig,
		//thresholdConfig: thresholdConfig,
	}
}


func (n *backend_server) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	//===============================================================
	//catch or finally
	defer func() {
		//catch or finally
		if err := recover(); err != nil {
			//catch
			os.Exit(0)
		}
	}()
	//===============================================================
	close(ready)

	backend_service := service.NewBackendServices(n.gmtTimeGapHour, n.influxCon, n.configDbCon, n.portalDbCon, n.boshCon, n.redisCon, n.mailConfig, n.config)

	//최초 실행시 Bosh VM정보 동기화
	backend_service.CreateUpdateBoshData(*n.boshCon)

	ticker := time.NewTicker(time.Duration(n.batchInterval) * time.Second)

	backend_service.StopProcess()

	for {
		select {
		case <- ticker.C:

			var wg sync.WaitGroup
			wg.Add(3)

			go func(wg *sync.WaitGroup){
				defer wg.Done()
				backend_service.CreateUpdateBoshData(*n.boshCon)
			}(&wg)

			//임계치 초과시 Alarm 전송
			go func(wg *sync.WaitGroup){
				defer wg.Done()
				backend_service.CreateAlarmData()
			}(&wg)

			go func(wg *sync.WaitGroup){
				defer wg.Done()
				backend_service.UserPortalService()
			}(&wg)

			wg.Wait()
		}
	}
}


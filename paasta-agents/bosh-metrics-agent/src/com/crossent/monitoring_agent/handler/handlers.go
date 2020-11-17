package handler

import (
	"os"

	"code.cloudfoundry.org/lager"
	"com/crossent/monitoring_agent/services"

	"github.com/tedsuo/ifrit"
	"time"
)

const (
	statsInterval = 30 * time.Second
)

type metricsSenderServer struct {
	logger    lager.Logger
	influxCon *services.InfluxConfig
	origin    *services.OriginConfig
	cellIp    string
}

func New(logger lager.Logger, influxCon *services.InfluxConfig, origin *services.OriginConfig) ifrit.Runner {
	return &metricsSenderServer{
		logger:    logger,
		influxCon: influxCon,
		origin:    origin,
	}
}

func (n *metricsSenderServer) Run(<-chan os.Signal, chan<- struct{}) error {
	//===============================================================
	// Call Service
	metricsSender := services.NewMetricSender(n.logger, n.influxCon, n.origin, statsInterval)
	err := metricsSender.SendMetricsToInfluxDb(nil)
	//===============================================================
	return err
}

package v1

import "go.uber.org/zap"

type SaaS struct {
	PinpointWebUrl string
	Logger *zap.Logger
}
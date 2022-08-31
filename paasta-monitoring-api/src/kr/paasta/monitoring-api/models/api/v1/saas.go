package v1

import "github.com/sirupsen/logrus"

type SaaS struct {
	PinpointWebUrl string
	Logger *logrus.Logger
}
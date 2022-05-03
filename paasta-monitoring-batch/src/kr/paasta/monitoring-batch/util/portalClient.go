package util

import (
	"time"
	"net"
	"net/http"
)

func NewPortalClient() *http.Client{
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{

			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   100000000 * time.Second,
				KeepAlive: 300000000 * time.Second,
			}).Dial,
		},
	}
}
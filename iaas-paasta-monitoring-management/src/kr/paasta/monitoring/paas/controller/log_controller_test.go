package controller

import (
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"net/http"
)

var _ = Describe("PaaS Log Controller", func() {

	Describe("Log Recent & Specific", func() {

		Context("Logs View", func() {

			It("Log Recent List", func() {
				res, err := DoGet(testUrl + "/v2/paas/log/recent?id=0811ed1c-4f17-4c2b-8260-34d59dc51eda&logType=cf&pageIndex=1&pageItems=100&period=5m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

			It("Log Specific List", func() {
				res, err := DoGet(testUrl + "/v2/paas/log/specific?id=0811ed1c-4f17-4c2b-8260-34d59dc51eda&logType=cf&pageIndex=1&pageItems=1000&period=5m&startTime=18:12:41&endTime=19:12:41&targetDate=2019-11-18")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

		})
	})

})

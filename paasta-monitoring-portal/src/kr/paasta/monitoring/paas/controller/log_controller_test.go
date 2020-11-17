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
				res, err := DoGet(testUrl + "/v2/paas/log/recent?id=22a12601-c141-4f05-8d15-41a8ab3786bf&logType=cf&pageIndex=1&pageItems=100&period=5m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

			It("Log Specific List", func() {
				res, err := DoGet(testUrl + "/v2/paas/log/specific?id=22a12601-c141-4f05-8d15-41a8ab3786bf&logType=cf&pageIndex=1&pageItems=1000&period=5m&startTime=07:11:29&endTime=18:11:29&targetDate=2020-11-12")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

		})
	})

})

package controller

import (
	"github.com/stretchr/testify/assert"
	. "github.com/onsi/ginkgo"
	"net/http"
)

var _ = Describe("IaaS Log Controller", func() {

	Describe("Log Recent & Specific", func() {

		Context("Logs View", func() {

			It("Log Recent List", func() {
				res, err := DoGet(testUrl + "/v2/iaas/log/recent?hostname=compute1&logType=log&pageIndex=1&pageItems=50&period=5m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NotEmpty(t, res.Content)
			})

			It("Log Specific List", func() {
				res, err := DoGet(testUrl + "/v2/iaas/log/specific?endTime=13:26:26&hostname=compute1&logType=log&pageIndex=1&pageItems=50&period=5m&startTime=01:21:00&targetDate=2018-05-15")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NotEmpty(t, res.Content)
			})

		})

	})

})
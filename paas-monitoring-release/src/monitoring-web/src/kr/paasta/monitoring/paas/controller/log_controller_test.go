package controller

import (
	"github.com/stretchr/testify/assert"
	. "github.com/onsi/ginkgo"
	"net/http"
)

var _ = Describe("PaaS Log Controller", func() {

	Describe("Log Recent & Specific", func() {

		Context("Logs View", func() {

			It("Log Recent List", func() {
				res, err := DoGet(testUrl + "/v2/paas/log/recent?id=2a1ba4bc-ecf5-41ea-afca-2b24d33e853f&logType=cf&pageIndex=1&pageItems=50&period=5m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NotEmpty(t, res.Content)
			})

			It("Log Specific List", func() {
				res, err := DoGet(testUrl + "/v2/paas/log/specific?endTime=01:17:47&id=4a89a8ed-c9c2-4263-9651-04a012c2f336&logType=cf&pageIndex=1&pageItems=50&period=5m&startTime=01:12:00&targetDate=2018-05-16")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NotEmpty(t, res.Content)
			})

		})
	})

})

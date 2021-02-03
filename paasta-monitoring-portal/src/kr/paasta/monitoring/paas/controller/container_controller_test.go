package controller

import (
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"net/http"
)

var containerId = "whj6uom6feb2-0"
var zoneName = "z1"
var testStatus = "critical"

var _ = Describe("ContainerController", func() {

	Describe("Container", func() {
		Context("Container Overview", func() {

			It("PAAS_CELL_OVERVIEW", func() {
				res, err := DoGet(testUrl + "/v2/paas/cell/overview")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)

			})

			It("PAAS_CONTAINER_OVERVIEW", func() {
				res, err := DoGet(testUrl + "/v2/paas/container/overview")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("PAAS_CONTAINER_SUMMARY", func() {
				res, err := DoGet(testUrl + "/v2/paas/container/summary")
				assert.Nil(tt, err)

				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("PAAS_CONTAINER_OVERVIEW_MAIN", func() {
				res, err := DoGet(testUrl + "/v2/paas/container/relationship")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("PAAS_CONTAINER_RELATIONSHIP", func() {
				res, err := DoGet(testUrl + "/v2/paas/container/relationship/" + zoneName)
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("PAAS_CELL_OVERVIEW_STATE_LIST", func() {
				res, err := DoGet(testUrl + "/v2/paas/cell/overview/" + testStatus)
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("PAAS_CONTAINER_OVERVIEW_STATE_LIST", func() {
				res, err := DoGet(testUrl + "/v2/paas/container/overview/" + testStatus)
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})
		})

		Context("Container Detail", func() {

			It("Get Instance Cpu Usage", func() {
				res, err := DoGet(testUrl + "/v2/paas/container/cpu/" + containerId + "/usages?defaultTimeRange=15m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

			It("Get Instance Cpu Load", func() {
				res, err := DoGet(testUrl + "/v2/paas/container/cpu/" + containerId + "/loads?defaultTimeRange=15m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

			It("Get Instance Memory Usage", func() {
				res, err := DoGet(testUrl + "/v2/paas/container/memory/" + containerId + "/usages?defaultTimeRange=15m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

			It("Get Instance disk Usage", func() {
				res, err := DoGet(testUrl + "/v2/paas/container/disk/" + containerId + "/usages?defaultTimeRange=15m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

			It("Get Instance network Bytes", func() {
				res, err := DoGet(testUrl + "/v2/paas/container/network/" + containerId + "/bytes?defaultTimeRange=15m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

			It("Get Instance network drops", func() {
				res, err := DoGet(testUrl + "/v2/paas/container/network/" + containerId + "/drops?defaultTimeRange=15m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

			It("Get Instance network errors", func() {
				res, err := DoGet(testUrl + "/v2/paas/container/network/" + containerId + "/errors?defaultTimeRange=15m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

		})
	})

})

package controller

import (
	"github.com/stretchr/testify/assert"
	. "github.com/onsi/ginkgo"
	"net/http"
)

var containerId = "w9jobutarvce-0"
var zoneName = "z2"
var testStatus = "critical"

var _ = Describe("ContainerController", func() {

	Describe("Container", func() {
		Context("Container Overview", func() {

			It("PAAS_CELL_OVERVIEW", func() {
				res, err := DoGet(testUrl + "/v2/paas/cell/overview")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_CONTAINER_OVERVIEW", func() {
				res, err := DoGet(testUrl + "/v2/paas/container/overview")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_CONTAINER_SUMMARY", func() {
				res, err := DoGet(testUrl + "/v2/paas/container/summary")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_CONTAINER_OVERVIEW_MAIN", func() {
				res, err := DoGet(testUrl + "/v2/paas/container/relationship")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_CONTAINER_RELATIONSHIP", func() {
				res, err := DoGet(testUrl + "/v2/paas/container/relationship/" + zoneName )
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_CELL_OVERVIEW_STATE_LIST", func() {
				res, err := DoGet(testUrl + "/v2/paas/cell/overview/" + testStatus )
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_CONTAINER_OVERVIEW_STATE_LIST", func() {
				res, err := DoGet(testUrl + "/v2/paas/container/overview/" + testStatus )
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})
		})

		Context("Container Detail", func() {

			It("Get Instance Cpu Usage", func() {
				res, err := DoGet(testUrl + "/v2/paas/container/cpu/" + containerId + "/usages?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NotEmpty(t, res.Content)
			})

			It("Get Instance Cpu Load", func() {
				res, err := DoGet(testUrl + "/v2/paas/container/cpu/" + containerId + "/loads?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NotEmpty(t, res.Content)
			})

			It("Get Instance Memory Usage", func() {
				res, err := DoGet(testUrl + "/v2/paas/container/memory/" + containerId + "/usages?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NotEmpty(t, res.Content)
			})

			It("Get Instance disk Usage", func() {
				res, err := DoGet(testUrl + "/v2/paas/container/disk/" + containerId + "/usages?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NotEmpty(t, res.Content)
			})

			It("Get Instance network Bytes", func() {
				res, err := DoGet(testUrl + "/v2/paas/container/network/" + containerId + "/bytes?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NotEmpty(t, res.Content)
			})

			It("Get Instance network drops", func() {
				res, err := DoGet(testUrl + "/v2/paas/container/network/" + containerId + "/drops?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NotEmpty(t, res.Content)
			})

			It("Get Instance network errors", func() {
				res, err := DoGet(testUrl + "/v2/paas/container/network/" + containerId + "/errors?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NotEmpty(t, res.Content)
			})

		})
	})

})

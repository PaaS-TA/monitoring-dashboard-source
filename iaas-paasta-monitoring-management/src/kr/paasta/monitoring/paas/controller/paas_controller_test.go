package controller

import (
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"net/http"
)

var vms = "dc73c6bf-9fa7-4689-41af-6da0281f5cf3/memory"
var vmsId = "ab63c109-37c9-402b-8367-abcb3a00f62e"

var _ = Describe("PaastaController", func() {

	Describe("Paasta", func() {

		var testStatus = "critical"

		Context("Paasta Overview", func() {
			It("PAAS_PAASTA_OVERVIEW", func() {
				res, err := DoGet(testUrl + "/v2/paas/paasta/overview")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("PAAS_PAASTA_SUMMARY", func() {
				res, err := DoGet(testUrl + "/v2/paas/paasta/summary")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("PAAS_PAASTA_TOPPROCESS_MEMORY", func() {
				res, err := DoGet(testUrl + "/v2/paas/paasta/topprocess/" + vms + "/memory")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("PAAS_PAASTA_OVERVIEW_STATUS", func() {
				res, err := DoGet(testUrl + "/v2/paas/paasta/overview/" + testStatus)
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})
		})

		Context("Paasta Detail", func() {

			It("Get Instance Cpu Usage", func() {
				res, err := DoGet(testUrl + "/v2/paas/paasta/cpu/" + vmsId + "/usages?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

			It("Get Instance Cpu Load", func() {
				res, err := DoGet(testUrl + "/v2/paas/paasta/cpu/" + vmsId + "/loads?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

			It("Get Instance Memory Usage", func() {
				res, err := DoGet(testUrl + "/v2/paas/paasta/memory/" + vmsId + "/usages?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

			It("Get Instance disk Usage", func() {
				res, err := DoGet(testUrl + "/v2/paas/paasta/disk/" + vmsId + "/usages?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

			It("Get Instance disk Ios", func() {
				res, err := DoGet(testUrl + "/v2/paas/paasta/disk/" + vmsId + "/ios?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

			It("Get Instance network Bytes", func() {
				res, err := DoGet(testUrl + "/v2/paas/paasta/network/" + vmsId + "/bytes?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

			It("Get Instance network Packets", func() {
				res, err := DoGet(testUrl + "/v2/paas/paasta/network/" + vmsId + "/packets?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

			It("Get Instance network drops", func() {
				res, err := DoGet(testUrl + "/v2/paas/paasta/network/" + vmsId + "/drops?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

			It("Get Instance network errors", func() {
				res, err := DoGet(testUrl + "/v2/paas/paasta/network/" + vmsId + "/errors?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

		})
	})

})

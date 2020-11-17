package controller

import (
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"os"
)

var id = "77e3669f-7850-437c-a098-3fea73b8c384"

var _ = Describe("BoshController", func() {

	Describe("Bosh", func() {
		Context("Bosh Overview", func() {

			It("PAAS_BOSH_STATUS_OVERVIEW", func() {
				res, err := DoGet(testUrl + "/v2/paas/bosh/overview")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)

				f, qerr := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
				if qerr != nil {
					log.Fatalf("error opening file: %v", qerr)
				}
				defer f.Close()

				log.SetOutput(f)
				log.Println("err :", err)
				log.Println("res :", res)
				log.Println("tt :", tt)
			})

			It("PAAS_BOSH_STATUS_SUMMARY", func() {
				res, err := DoGet(testUrl + "/v2/paas/bosh/summary")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("PAAS_BOSH_STATUS_TOPPROCESS", func() {
				res, err := DoGet(testUrl + "/v2/paas/bosh/topprocess/" + id + "/memory")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)

			})

		})

		Context("Bosh Detail", func() {

			It("Get Instance Cpu Usage", func() {
				res, err := DoGet(testUrl + "/v2/paas/bosh/cpu/" + id + "/usages?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

			It("Get Instance Cpu Load", func() {
				res, err := DoGet(testUrl + "/v2/paas/bosh/cpu/" + id + "/loads?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)

				assert.NotEmpty(tt, res.Content)
			})

			It("Get Instance Memory Usage", func() {
				res, err := DoGet(testUrl + "/v2/paas/bosh/memory/" + id + "/usages?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

			It("Get Instance disk Usage", func() {
				res, err := DoGet(testUrl + "/v2/paas/bosh/disk/" + id + "/usages?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

			It("Get Instance disk Ios", func() {
				res, err := DoGet(testUrl + "/v2/paas/bosh/disk/" + id + "/ios?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

			It("Get Instance network Bytes", func() {
				res, err := DoGet(testUrl + "/v2/paas/bosh/network/" + id + "/bytes?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

			It("Get Instance network Packets", func() {
				res, err := DoGet(testUrl + "/v2/paas/bosh/network/" + id + "/packets?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

			It("Get Instance network drops", func() {
				res, err := DoGet(testUrl + "/v2/paas/bosh/network/" + id + "/drops?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

			It("Get Instance network errors", func() {
				res, err := DoGet(testUrl + "/v2/paas/bosh/network/" + id + "/errors?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
				assert.NotEmpty(tt, res.Content)
			})

		})
	})

})

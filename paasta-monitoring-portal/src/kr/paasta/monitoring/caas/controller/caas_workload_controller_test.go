package controller

import (
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"net/http"
)

var _ = Describe("CaaS_Workload", func() {

	Describe("Workload", func() {
		Context("Workload Info", func() {

			It("CAAS_WORKLOAD_STATUS", func() {
				res, err := DoGet(testUrl + "/v2/caas/monitoring/workloadsStatus")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("CAAS_WORKLOAD_SUMMARY", func() {
				res, err := DoGet(testUrl + "/v2/caas/monitoring/workerloadsConainerSummary")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("CAAS_WORKLOAD_DEPLOYMENT_CONTAINER", func() {
				res, err := DoGet(testUrl + "/v2/caas/monitoring/contiList?WorkloadsName=deployment")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("CAAS_WORKLOAD_STATEFULSET_CONTAINER", func() {
				res, err := DoGet(testUrl + "/v2/caas/monitoring/contiList?WorkloadsName=statefulset")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("CAAS_WORKLOAD_DAEMONSET_CONTAINER", func() {
				res, err := DoGet(testUrl + "/v2/caas/monitoring/contiList?WorkloadsName=daemonset")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("CAAS_WORKLOAD_DEPLOYMENT_GRAPH", func() {
				res, err := DoGet(testUrl + "/v2/caas/monitoring/workloadsGraph?WorkloadsName=deployment")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("CAAS_WORKLOAD_STATEFULSET_GRAPH", func() {
				res, err := DoGet(testUrl + "/v2/caas/monitoring/workloadsGraph?WorkloadsName=statefulset")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("CAAS_WORKLOAD_DAEMONSET_GRAPH", func() {
				res, err := DoGet(testUrl + "/v2/caas/monitoring/workloadsGraph?WorkloadsName=daemonset")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

		})

	})

})

package controller

import (
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"net/http"
)

var node_name = "ip-10-0-121-121.ap-northeast-2.compute.internal"
var instance = "10.0.121.121:9100"

var _ = Describe("CaaS_Cluster", func() {

	Describe("Cluster", func() {
		Context("Cluster Info", func() {

			It("CAAS_CLUSTER_AVG", func() {
				res, err := DoGet(testUrl + "/v2/caas/monitoring/clusterAvg")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("CAAS_CLUSTER_OVERVIEW", func() {
				res, err := DoGet(testUrl + "/v2/caas/monitoring/clusterOverview")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("CAAS_WORKLOAD_STATUS", func() {
				res, err := DoGet(testUrl + "/v2/caas/monitoring/workloadsStatus")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("CAAS_WORKLOAD_LIST", func() {
				res, err := DoGet(testUrl + "/v2/caas/monitoring/workerNodeList")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("CAAS_WORKLOAD_GRAPH", func() {
				res, err := DoGet(testUrl + "/v2/caas/monitoring/workerNodeGraph?NodeName=" + node_name + "&Instance=" + instance)
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("CAAS_WORKLOAD_GRAPH_LIST", func() {
				res, err := DoGet(testUrl + "/v2/caas/monitoring/workerNodeGraphList?NodeName=" + node_name + "&Instance=" + instance)
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

		})

	})

})

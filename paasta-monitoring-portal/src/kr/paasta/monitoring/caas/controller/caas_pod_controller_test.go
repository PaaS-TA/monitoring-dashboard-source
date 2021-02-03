package controller

import (
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"net/http"
)

var pod_name = "alertmanager-prometheus-prometheus-oper-alertmanager-0"
var name_space = "prometheus"
var container_name = "alertmanager"

var _ = Describe("CaaS_Pod", func() {

	Describe("Pod", func() {
		Context("Pod Info", func() {

			It("CAAS_POD_STATUS", func() {
				res, err := DoGet(testUrl + "/v2/caas/monitoring/podStat")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("CAAS_POD_LIST", func() {
				res, err := DoGet(testUrl + "/v2/caas/monitoring/podList")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("CAAS_POD_CONTAINER", func() {
				res, err := DoGet(testUrl + "/v2/caas/monitoring/contiList?PodName=" + pod_name)
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("CAAS_POD_GRAPH", func() {
				res, err := DoGet(testUrl + "/v2/caas/monitoring/podGraph?PodName=" + pod_name)
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("CAAS_POD_CONTAINER_LOG", func() {
				res, err := DoGet(testUrl + "/v2/caas/monitoring/contiInfoLog?ContainerName=" + container_name + "&NameSpace=" + name_space + "&PodName=" + pod_name)
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

		})

	})

})

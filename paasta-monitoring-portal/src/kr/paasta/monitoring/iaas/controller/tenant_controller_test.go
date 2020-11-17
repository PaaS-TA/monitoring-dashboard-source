package controller

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"kr/paasta/monitoring/iaas/model"
	"net/http"
)

var _ = Describe("IaaS Tenant Controller", func() {

	Describe("IaaS Tenant ", func() {
		var instanceId string
		Context("Tenant Summary", func() {

			var projectId string
			It("Tenant Summry", func() {
				res, err := DoGet(testUrl + "/v2/iaas/tenant/summary")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NotEmpty(t, res.Content)

				projectInfoList := &[]model.TenantSummaryInfo{}
				json.Unmarshal([]byte(res.Content), projectInfoList)
				assert.True(t, len((*projectInfoList)) > 0)
				projectId = (*projectInfoList)[0].Id
			})

			It("Instance From Tenant", func() {
				res, err := DoGet(testUrl + "/v2/iaas/tenant/" + projectId + "/instances?limit=1")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NotEmpty(t, res.Content)
				//Response Data형변환
				var msgMapTemplate interface{}
				json.Unmarshal([]byte(res.Content), &msgMapTemplate)
				msgMap := msgMapTemplate.(map[string]interface{})
				tmp := msgMap["metric"].([]interface{})
				data := tmp[0].(map[string]interface{})

				assert.True(t, data["instance_id"] != "")
				instanceId = data["instance_id"].(string)
			})

		})

		Context("Instance Metric Detail", func() {

			It("Get Instance Cpu Usage", func() {
				res, err := DoGet(testUrl + "/v2/iaas/tenant/cpu/" + instanceId + "/usages?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NotEmpty(t, res.Content)
			})

			It("Get Instance Memory Usage", func() {
				res, err := DoGet(testUrl + "/v2/iaas/tenant/memory/" + instanceId + "/usages?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NotEmpty(t, res.Content)
			})

			It("Get Instance diskRead", func() {
				res, err := DoGet(testUrl + "/v2/iaas/tenant/disk/" + instanceId + "/reads?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NotEmpty(t, res.Content)
			})

			It("Get Instance diskWrite", func() {
				res, err := DoGet(testUrl + "/v2/iaas/tenant/disk/" + instanceId + "/writes?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NotEmpty(t, res.Content)
			})

			It("Get Instance networkIo", func() {
				res, err := DoGet(testUrl + "/v2/iaas/tenant/network/" + instanceId + "/ios?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NotEmpty(t, res.Content)
			})

			It("Get Instance networkPackets", func() {
				res, err := DoGet(testUrl + "/v2/iaas/tenant/network/" + instanceId + "/packets?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NotEmpty(t, res.Content)
			})

		})

	})

})

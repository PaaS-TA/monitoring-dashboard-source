package controller_test

import (
	"net/http"
	"github.com/stretchr/testify/assert"
	. "github.com/onsi/ginkgo"
	_ "github.com/go-sql-driver/mysql"
	cb "kr/paasta/monitoring/domain"
	"encoding/json"
)

var _ = Describe("MetricsController", func() {

	Describe("Metrics", func() {
		Context("Disk I/O", func() {
			It("get bosh disk I/O usage list", func() {
				param := "?defaultTimeRange=30m&groupBy=1m&serviceName=micro-bosh"
				res, err := DoGet(testUrl + "/diskIO/bos" + param)
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)

				metrics := &cb.DiskIOUsageList{}
				json.Unmarshal([]byte(res.Content), metrics)
				assert.True(t, len(metrics.Data) == 2)
			})
			It("get PaaS-TA controller disk I/O usage list", func() {
				param := "?defaultTimeRange=30m&groupBy=1m&serviceName=api_worker"
				res, err := DoGet(testUrl + "/diskIO/ctl" + param)
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)

				metrics := &cb.DiskIOUsageList{}
				json.Unmarshal([]byte(res.Content), metrics)
				assert.True(t, len(metrics.Data) == 2)
			})
			It("get PaaS-TA container disk I/O usage list", func() {
				param := "?defaultTimeRange=30m&groupBy=1m&serviceName=access"
				res, err := DoGet(testUrl + "/diskIO/ctn" + param)
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)

				metrics := &cb.DiskIOUsageList{}
				json.Unmarshal([]byte(res.Content), metrics)
				assert.True(t, len(metrics.Data) == 2)
			})
		})

		Context("Network I/O", func() {
			It("get bosh network I/O usage list", func() {
				param := "?defaultTimeRange=30m&groupBy=1m&serviceName=micro-bosh"
				res, err := DoGet(testUrl + "/networkIO/bos" + param)
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)

				metrics := &cb.NetworkIOUsageList{}
				json.Unmarshal([]byte(res.Content), metrics)
				assert.True(t, len(metrics.Data) == 2)
			})
			It("get PaaS-TA controller network I/O usage list", func() {
				param := "?defaultTimeRange=30m&groupBy=1m&serviceName=api_worker"
				res, err := DoGet(testUrl + "/networkIO/ctl" + param)
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)

				metrics := &cb.NetworkIOUsageList{}
				json.Unmarshal([]byte(res.Content), metrics)
				assert.True(t, len(metrics.Data) == 2)
			})
			It("get PaaS-TA container network I/O usage list", func() {
				param := "?defaultTimeRange=30m&groupBy=1m&serviceName=access"
				res, err := DoGet(testUrl + "/networkIO/ctn" + param)
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)

				metrics := &cb.NetworkIOUsageList{}
				json.Unmarshal([]byte(res.Content), metrics)
				assert.True(t, len(metrics.Data) == 2)
			})
			It("get application container network I/O usage list", func() {
				param := "?defaultTimeRange=30m&groupBy=1m"
				res, err := DoGet(testUrl + "/networkIO/ctn" + param)
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)

				metrics := &cb.NetworkIOUsageList{}
				json.Unmarshal([]byte(res.Content), metrics)
				assert.True(t, len(metrics.Data) == 2)
			})
		})
		Context("Top Process", func() {
			It("get bosh Top Process list", func() {
				param := "?defaultTimeRange=30m&groupBy=1m&serviceName=micro-bosh&addr=10.10.10.4"
				res, err := DoGet(testUrl + "/topProcess/bos" + param)
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)

				metrics := &cb.TopProcessList{}
				json.Unmarshal([]byte(res.Content), metrics)
				assert.True(t, len(metrics.Data) == 10)
			})
			It("get PaaS-TA controller Top Process list", func() {
				param := "?defaultTimeRange=30m&groupBy=1m&serviceName=api_worker&addr=10.244.0.29"
				res, err := DoGet(testUrl + "/topProcess/ctl" + param)
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)

				metrics := &cb.TopProcessList{}
				json.Unmarshal([]byte(res.Content), metrics)
				assert.True(t, len(metrics.Data) == 10)
			})
			It("get PaaS-TA container Top Process list", func() {
				param := "?defaultTimeRange=30m&groupBy=1m&serviceName=eth0"
				res, err := DoGet(testUrl + "/topProcess/ctn" + param)
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)

				metrics := &cb.TopProcessList{}
				json.Unmarshal([]byte(res.Content), metrics)
				assert.True(t, len(metrics.Data) == 10)
			})
		})

	})

})

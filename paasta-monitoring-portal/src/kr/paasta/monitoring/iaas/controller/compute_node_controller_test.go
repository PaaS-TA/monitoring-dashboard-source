package controller

import (
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"net/http"
)

var _ = Describe("IaaSComputeNodeController", func() {

	Describe("IaaS Compute Node Overview", func() {

		Context("GET IaaS Compute Node Overview", func() {

			It("IaaS Main Summary", func() {
				res, err := DoGet(testUrl + "/v2/iaas/main/summary")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("IaaS Compute Node List", func() {
				res, err := DoGet(testUrl + "/v2/iaas/nodes")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("IaaS Compute Node Cpu Top Process", func() {
				res, err := DoGet(testUrl + "/v2/iaas/node/topprocess/compute2/cpu")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("IaaS Compute Node Memory Top Process", func() {
				res, err := DoGet(testUrl + "/v2/iaas/node/topprocess/compute2/memory")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

		})

	})

	Describe("IaaS Compute Node Detail", func() {

		Context("GET IaaS Compute Node Detail", func() {

			It("IaaS Node Cpu Usage", func() {
				res, err := DoGet(testUrl + "/v2/iaas/node/cpu/compute2/usages?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("IaaS Node Cpu Load", func() {
				res, err := DoGet(testUrl + "/v2/iaas/node/cpu/compute2/loads?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("IaaS Node Memory Usage", func() {
				res, err := DoGet(testUrl + "/v2/iaas/node/memory/compute2/swaps?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("IaaS Node Swap", func() {
				res, err := DoGet(testUrl + "/v2/iaas/node/memory/compute2/usages?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("IaaS Node Disk Usage", func() {
				res, err := DoGet(testUrl + "/v2/iaas/node/disk/compute2/usages?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("IaaS Node Disk Read", func() {
				res, err := DoGet(testUrl + "/v2/iaas/node/disk/compute2/reads?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("Node Disk Write", func() {
				res, err := DoGet(testUrl + "/v2/iaas/node/disk/compute2/writes?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("Node Network IO", func() {
				res, err := DoGet(testUrl + "/v2/iaas/node/network/compute2/kbytes?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("Node Network Err", func() {
				res, err := DoGet(testUrl + "/v2/iaas/node/network/compute2/errors?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("Node Network DropPacket", func() {
				res, err := DoGet(testUrl + "/v2/iaas/node/network/compute2/droppackets?defaultTimeRange=10m&groupBy=1m")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

		})

	})

})

package test

import (
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"net/http"
)

var tenant_id = "2275753136b2479fa3c6a576279f25cd"

var _ = Describe("Openstack API", func() {

	Describe("Openstack", func() {

		Context("Openstack hypervisor GET", func() {
			It("IaaS alarm policy GET", func() {
				res, err := DoGet(testUrl + "/v2/iaas/hypervisor/list")

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("IaaS alarm policy GET", func() {
				res, err := DoGet(testUrl + "/v2/iaas/hyper/statistics")

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("IaaS alarm policy GET", func() {
				res, err := DoGet(testUrl + "/v2/iaas/server/list?tenant_id=" + tenant_id)

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("IaaS alarm policy GET", func() {
				res, err := DoGet(testUrl + "/v2/iaas/project/list")

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("IaaS alarm policy GET", func() {
				res, err := DoGet(testUrl + "/v2/iaas/instance/usage/list")

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})
		})

	})
})

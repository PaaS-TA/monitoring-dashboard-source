package test

import (
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"net/http"
)

var instance_id = "1783bc89-c3c3-4c16-a48a-587c1e8acb"
var host = ""

var _ = Describe("Zabbix API", func() {

	Describe("Zabbix", func() {

		Context("Instance CPU usage", func() {
			It("IaaS alarm policy GET", func() {
				res, err := DoGet(testUrl + "/v2/iaas/instance/cpu/usage?instance_id=" + instance_id +"&host=" + host)

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("Instance Memory usage", func() {
				res, err := DoGet(testUrl + "/v2/iaas/instance/memory/usage?instance_id=" + instance_id +"&host=" + host)

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("Instance Disk usage", func() {
				res, err := DoGet(testUrl + "/v2/iaas/instance/disk/usage?instance_id=" + instance_id +"&host=" + host)

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("Instance CPU load average usage", func() {
				res, err := DoGet(testUrl + "/v2/iaas/instance/cpu/load/average?instance_id=" + instance_id +"&host=" + host)

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("Instance Disk IO rate usage", func() {
				res, err := DoGet(testUrl + "/v2/iaas/instance/disk/io/rate?instance_id=" + instance_id +"&host=" + host)

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("Instance Network IO usage", func() {
				res, err := DoGet(testUrl + "/v2/iaas/instance/network/io/bytes?instance_id=" + instance_id +"&host=" + host)

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})
		})

	})
})

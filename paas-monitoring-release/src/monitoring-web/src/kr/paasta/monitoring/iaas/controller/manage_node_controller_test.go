package controller

import (
	"github.com/stretchr/testify/assert"
	. "github.com/onsi/ginkgo"
	"net/http"
)

var _ = Describe("IaaSManageNodeController", func() {

	Describe("IaaS Manage Node Overview", func() {

		Context("GET IaaS Manage Node Overview", func() {

			It("IaaS Manage Node Summary", func() {
				res, err := DoGet(testUrl + "/v2/iaas/node/manage/summary")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("IaaS Cpu Top Process", func() {
				res, err := DoGet(testUrl + "/v2/iaas/node/topprocess/controller/cpu")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("IaaS Memory Top Process", func() {
				res, err := DoGet(testUrl + "/v2/iaas/node/topprocess/controller/memory")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("IaaS RabbitMq Status", func() {
				res, err := DoGet(testUrl + "/v2/iaas/node/rabbitmq/summary")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

		})

	})

})
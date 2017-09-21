package controller_test

import (
	"net/http"
	"github.com/stretchr/testify/assert"
	. "github.com/onsi/ginkgo"
	_ "github.com/go-sql-driver/mysql"
)

var _ = Describe("ContainerController", func() {

	Describe("Container Placement State", func() {
		It("get container placement state", func() {
			res, err := DoGet(testUrl + "/containerDeploy")
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, res.Code)
		})

		It("not allowed method", func() {
			res, err := Do("DELETE", testUrl + "/containerDeploy", nil)
			assert.Nil(t, err)
			assert.Equal(t, http.StatusMethodNotAllowed, res.Code)
		})
	})

})

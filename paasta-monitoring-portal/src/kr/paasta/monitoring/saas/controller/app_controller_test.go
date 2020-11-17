package controller

import (
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"net/http"
)

var _ = Describe("ApplicationController", func() {

	Describe("Container", func() {
		Context("Application List", func() {

			It("PAAS_CELL_OVERVIEW", func() {
				res, err := DoGet(testUrl + "/v2/saas/app/application/list")
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

		})

	})

})

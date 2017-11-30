package controller_test

import (
	"net/http"
	"github.com/stretchr/testify/assert"
	. "github.com/onsi/ginkgo"
	_ "github.com/go-sql-driver/mysql"
	"bytes"
)

var _ = Describe("AlarmPolicyController", func() {

	Describe("AlarmPolicy", func() {
		Context("Get", func() {
			It("all list", func() {
				res, err := DoGet(testUrl + "/alarmsPolicy")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})
		})

		Context("Update", func() {
			It("policy", func() {
				client := &http.Client{}

				json := `[{"originType":"pas","alarmType":"memory","warningThreshold":80,"criticalThreshold":90,"repeatTime":5}]`
				body := bytes.NewReader([]byte(json))
				request, err := http.NewRequest("PUT", testUrl + "/alarmsPolicy", body)
				res, err := client.Do(request)
				assert.Nil(t, err)
				assert.Equal(t, http.StatusCreated, res.StatusCode)
			})

			It("Error - Bad Request", func() {
				client := &http.Client{}

				json := `[{"warningThreshold":75,"criticalThreshold":80,"repeatTime":1}]`
				body := bytes.NewReader([]byte(json))
				request, err := http.NewRequest("PUT", testUrl + "/alarmsPolicy", body)
				res, err := client.Do(request)
				assert.Nil(t, err)
				assert.Equal(t, http.StatusBadRequest, res.StatusCode)
			})
		})
	})

})
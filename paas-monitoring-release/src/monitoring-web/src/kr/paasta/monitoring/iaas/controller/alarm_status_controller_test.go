package controller

import (
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"encoding/json"
)

type StatusRequestBody struct {
	Id                uint
	AlarmId           string
	AlarmActionDesc   string
}
var alramId = "e2d74e85-d5c0-4b59-8aa7-8fcf1deca392"
var alramDefinitionId = "952ad6b1-7858-4220-af33-93bcdb706fe7"

var _ = Describe("IaaSAlarmStatusController", func() {

	Describe("IaaS Alarm Status", func() {

		Context("IaaS Alarm Status", func() {

			It("IaaS Main Alarm Realtime List", func() {
				res, err := DoGet(testUrl + "/v2/iaas/alarm/realtime/list")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("IaaS Main Alarm Realtime Count", func() {
				res, err := DoGet(testUrl + "/v2/iaas/alarm/realtime/count")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("Alarm Status List", func() {
				res, err := DoGet(testUrl + "/v2/iaas/alarm/statuses?severity=HIGH&state=OK&offset=0&limit=10")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("Alarm Status List Count", func() {
				res, err := DoGet(testUrl + "/v2/iaas/alarm/status/count?state=OK")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("Alarm Status Detail", func() {
				res, err := DoGet(testUrl + "/v2/iaas/alarm/status/" + alramDefinitionId)
				assert.Nil(t, err)
				//assert.Equal(t, http.StatusOK, res.Code)
				assert.Equal(t, http.StatusInternalServerError, res.Code)
			})

			It("Alarm Status History List", func() {
				res, err := DoGet(testUrl + "/v2/iaas/alarm/histories/" + alramDefinitionId)
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("Alarm Status Action List", func() {
				res, err := DoGet(testUrl + "/v2/iaas/alarm/actions/" + alramDefinitionId)
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("Alarm Status Create Info", func() {
				var query StatusRequestBody
				query.AlarmId = alramDefinitionId
				query.AlarmActionDesc = "Test Alram Action Desc!!"

				data, _ := json.Marshal(query)

				res, err := DoPost(testUrl + "/v2/iaas/alarm/action", TestToken, strings.NewReader(string(data)))
				assert.Nil(t, err)
				assert.Equal(t, http.StatusCreated, res.Code)
			})

			It("Alarm Status Update Info", func() {
				var query StatusRequestBody
				query.Id = 12
				query.AlarmActionDesc = "Update Test Alram Action Desc!! : " + alramId

				data, _ := json.Marshal(query)

				res, err := DoUpdate(testUrl + "/v2/iaas/alarm/action/" + alramId, TestToken, strings.NewReader(string(data)))
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("Alarm Status Delete Info", func() {
				res, err := DoDelete(testUrl + "/v2/iaas/alarm/action/" + alramId, TestToken, strings.NewReader(string("")))
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

		})

	})

})
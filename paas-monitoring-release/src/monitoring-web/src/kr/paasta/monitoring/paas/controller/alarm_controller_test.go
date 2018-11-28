package controller

import (
	"net/http"
	"github.com/stretchr/testify/assert"
	//"github.com/monasca/golang-monascaclient/monascaclient/models"
	. "github.com/onsi/ginkgo"
	"encoding/json"
	"strings"
	"time"
	"fmt"
)

type AlarmPolicyRequestBody struct {
	Id					uint     `json:"id"`
	OriginType         	string   `json:"originType"`
	AlarmType          	string   `json:"alarmType"`
	WarningThreshold   	int      `json:"warningThreshold"`
	CriticalThreshold  	int      `json:"criticalThreshold"`
	RepeatTime         	int      `json:"repeatTime"`
	Comment            	string   `json:"comment"`
	MeasureTime        	int   	 `json:"measureTime"`
	MailAddress        	string   `json:"mailAddress"`
	SnsType				string   `json:"snsType"`
	SnsId				string   `json:"snsId"`
	Token				string   `json:"token"`
	Expl				string   `json:"expl"`
	MailSendYn			string   `json:"mailSendYn"`
	SnsSendYn        	string   `json:"snsSendYn"`
	ModiDate           	time.Time `json:"modiDate"`
	ModiUser           	string   `json:"modiUser"`
}

type AlarmRequestBody struct {
	Id                uint
	OriginType        string
	OriginId          uint
	AlarmType         string
	Level             string
	AlarmTitle        string
	ResolveStatus     string
	SearchDateFrom    string
	SearchDateTo      string
}

type AlarmActionRequestBody struct {
	Id                uint
	AlarmId           uint
	AlarmActionDesc   string
	RegDate           time.Time
	RegUser           string
	ModiDate          time.Time
	ModiUser          string
}

var channel_id = "123"
var alarm_id = "44"
var resolveStatus = "2"
var action_id = "6"
var appId = "af9c7835-dd86-42f9-b105-dd4a3bae3f3c"
var appIndex = "1"

var _ = Describe("Alarm API", func() {

	Describe("Paas Alarm", func() {

		Context("Alarm Status", func() {
			It("PAAS_ALARM_REALTIME_COUNT", func() {
				res, err := DoGet(testUrl + "/v2/paas/alarm/realtime/count")

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_ALARM_REALTIME_LIST", func() {
				res, err := DoGet(testUrl + "/v2/paas/alarm/realtime/list")

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_ALARM_STATUS_LIST", func() {
				res, err := DoGet(testUrl + "/v2/paas/alarm/statuses?pageIndex=1&pageItems=10")

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_ALARM_STATUS_COUNT", func() {
				res, err := DoGet(testUrl + "/v2/paas/alarm/status/count?resolveStatus=1&state=ALARM")

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_ALARM_STATUS_DETAIL", func() {
				res, err := DoGet(testUrl + "/v2/paas/alarm/status/" + alarm_id)

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_ALARM_STATUS_UPDATE (PUT)", func() {

				query := AlarmRequestBody{
					ResolveStatus: resolveStatus,
				}

				data, _ := json.Marshal(query)

				fmt.Println(">>>>>>>>>>> data :", data)
				fmt.Println(">>>>>>>>>>> TestToken :", TestToken)

				res, _ := DoUpdate(testUrl + "/v2/paas/alarm/status/" + alarm_id, TestToken, strings.NewReader(string(data)))

				//assert.Nil(t, err)
				assert.Equal(t, http.StatusCreated, res.Code)
			})

			It("PAAS_ALARM_STATUS_RESOLVE", func() {
				res, err := DoGet(testUrl + "/v2/paas/alarm/status/" + resolveStatus)

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_ALARM_ACTION_CREATE (POST)", func() {

				query := AlarmActionRequestBody{
					AlarmId : 44,
					AlarmActionDesc :"test description",
				}

				data, _ := json.Marshal(query)

				res, _ := DoPost(testUrl + "/v2/paas/alarm/action", TestToken, strings.NewReader(string(data)))
				fmt.Println(res)
				//assert.Nil(t, err)
				//assert.Equal(t, http.StatusCreated, res.Code)
			})

			It("PAAS_ALARM_ACTION_UPDATE (PATCH)", func() {

				query := AlarmActionRequestBody{
					AlarmId : 44,
					AlarmActionDesc :"update test description",
				}

				data, _ := json.Marshal(query)

				res, _ := DoPatch(testUrl + "/v2/paas/alarm/action/" + action_id, TestToken, strings.NewReader(string(data)))
				fmt.Println(res)
				//assert.Nil(t, err)
				//assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_ALARM_ACTION_DELETE(DELETE)", func() {
				res, _ := DoDelete(testUrl + "/v2/paas/alarm/action/" + action_id, TestToken, strings.NewReader(string("")))
				fmt.Println(res)
				//assert.Nil(t, err)
				//assert.Equal(t, http.StatusOK, res.Code)
			})
		})


		Context("Alarm Statistics", func() {

			It("PAAS_ALARM_STATISTICS", func() {
				res, err := DoGet(testUrl + "/v2/paas/alarm/statistics?period=d&interval=1")

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_ALARM_STATISTICS_GRAPH_TOTAL", func() {
				res, err := DoGet(testUrl + "/v2/paas/alarm/statistics/graph/total?period=d&interval=1")

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_ALARM_STATISTICS_GRAPH_SERVICE", func() {
				res, err := DoGet(testUrl + "/v2/paas/alarm/statistics/graph/service?period=d&interval=1")

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_ALARM_STATISTICS_GRAPH_MATRIX", func() {
				res, err := DoGet(testUrl + "/v2/paas/alarm/statistics/graph/matrix?period=d&interval=1")

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})
		})


		Context("Alarm API Etc.", func() {

			It("PAAS_ALARM_CONTAINER_DEPLOY", func() {
				res, err := DoGet(testUrl + "/v2/paas/alarm/container/deploy")

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_ALARM_APP_RESOURCES", func() {
				res, err := DoGet(testUrl + "/v2/paas/alarm/app/resources?app_id=" + appId + "&app_index=" + appIndex)

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_ALARM_APP_RESOURCES_ALL", func() {
				res, err := DoGet(testUrl + "/v2/paas/alarm/app/resources/all?limit=10")

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_ALARM_APP_USAGES", func() {
				res, err := DoGet(testUrl + "/v2/paas/alarm/app/cpu/" + appId + "/" + appIndex + "/usages?defaultTimeRange=15m&groupBy=1m")

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_ALARM_APP_MEMORY_USAGES", func() {
				res, err := DoGet(testUrl + "/v2/paas/alarm/app/memory/" + appId + "/" + appIndex + "/usages?defaultTimeRange=15m&groupBy=1m")

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_ALARM_APP_DISK_USAGES", func() {
				res, err := DoGet(testUrl + "/v2/paas/alarm/app/disk/" + appId + "/" + appIndex + "/usages?defaultTimeRange=15m&groupBy=1m")

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_ALARM_APP_NETWORK_USAGES", func() {
				res, err := DoGet(testUrl + "/v2/paas/alarm/app/network/" + appId + "/" + appIndex + "/usages?defaultTimeRange=15m&groupBy=1m")

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

		})
	})

})


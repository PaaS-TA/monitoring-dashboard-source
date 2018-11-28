package controller

import (
	"net/http"
	"github.com/stretchr/testify/assert"
	//"github.com/monasca/golang-monascaclient/monascaclient/models"
	. "github.com/onsi/ginkgo"
	"encoding/json"
	"strings"
	"fmt"
)

var _ = Describe("Alarm API", func() {

	Describe("Paas Alarm", func() {
		Context("Alarm Policy", func() {
			It("PAAS_ALARM_POLICY_LIST", func() {
				res, err := DoGet(testUrl + "/v2/paas/alarm/policies")

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_ALARM_POLICY_UPDATE (PUT)", func() {
				var reqBody []AlarmPolicyRequestBody
				query := AlarmPolicyRequestBody{
					OriginType: "bos",
					AlarmType: "cpu",
					WarningThreshold: 85,
					CriticalThreshold: 90,
					RepeatTime: 10,
					MeasureTime: 100,
				}
				query1 := AlarmPolicyRequestBody{
					OriginType: "bos",
					AlarmType: "memory",
					WarningThreshold: 85,
					CriticalThreshold: 90,
					RepeatTime: 10,
					MeasureTime: 10,
				}
				query2 := AlarmPolicyRequestBody{
					OriginType: "bos",
					AlarmType: "disk",
					WarningThreshold: 85,
					CriticalThreshold: 90,
					RepeatTime: 10,
					MeasureTime: 100,
				}
				query3 := AlarmPolicyRequestBody{
					OriginType: "bos",
					MailAddress: "adminUser@gmail.com",
					MailSendYn: "Y",
					SnsSendYn: "N",
				}

				reqBody = append(reqBody, query)
				reqBody = append(reqBody, query1)
				reqBody = append(reqBody, query2)
				reqBody = append(reqBody, query3)

				data, _ := json.Marshal(reqBody)

				res, _ := DoUpdate(testUrl + "/v2/paas/alarm/policy", TestToken, strings.NewReader(string(data)))

				//assert.Nil(t, err)
				assert.Equal(t, http.StatusCreated, res.Code)
			})

			It("PAAS_ALARM_SNS_CHANNEL_LIST", func() {
				res, err := DoGet(testUrl + "/v2/paas/alarm/sns/channel/list")

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_ALARM_SNS_CHANNEL_CREATE(POST)", func() {

				query := AlarmPolicyRequestBody{
					OriginType:"pas",
					SnsId:"pas_test_bot",
					Token:"595845637:AAGUgw96sfoyTO3RUZoX-i06OHh9ZX0yt",
					Expl:"test pas bot",
					SnsSendYn:"N",
				}

				data, _ := json.Marshal(query)

				res, _ := DoPost(testUrl + "/v2/paas/alarm/sns/channel", TestToken, strings.NewReader(string(data)))
				fmt.Println(res)
				//assert.Nil(t, err)
				//assert.Equal(t, http.StatusCreated, res.Code)
			})

			It("PAAS_ALARM_SNS_CHANNEL_DELETE(DELETE)", func() {
				res, _ := DoDelete(testUrl + "/v2/paas/alarm/sns/channel/" + channel_id, TestToken, strings.NewReader(string("")))
				fmt.Println(res)
				//assert.Nil(t, err)
				//assert.Equal(t, http.StatusOK, res.Code)
			})

		})

	})

})


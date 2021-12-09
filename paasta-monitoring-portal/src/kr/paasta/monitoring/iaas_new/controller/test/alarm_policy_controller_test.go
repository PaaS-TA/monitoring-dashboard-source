package test

import (
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"kr/paasta/monitoring/paas/model"
	"net/http"
	"strings"
)

var alarm_sns_id = "1"

var _ = Describe("Alarm API", func() {

	Describe("IaaS Alarm", func() {

		Context("Iaas Alarm Policy GET", func() {
			It("IaaS alarm policy GET", func() {
				res, err := DoGet(testUrl + "/v2/iaas/alarm/policies")

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("IaaS Alarm Policy UPDATE", func() {

				var reqBody []model.AlarmPolicyRequest
				query := model.AlarmPolicyRequest{
					OriginType:        "ias",
					AlarmType:         "cpu",
					WarningThreshold:  85,
					CriticalThreshold: 90,
					RepeatTime:        10,
					MeasureTime:       100,
				}
				query1 := model.AlarmPolicyRequest{
					OriginType:        "ias",
					AlarmType:         "memory",
					WarningThreshold:  85,
					CriticalThreshold: 90,
					RepeatTime:        10,
					MeasureTime:       10,
				}
				query2 := model.AlarmPolicyRequest{
					OriginType:        "ias",
					AlarmType:         "disk",
					WarningThreshold:  85,
					CriticalThreshold: 90,
					RepeatTime:        10,
					MeasureTime:       100,
				}
				query3 := model.AlarmPolicyRequest{
					OriginType:  "bos",
					MailAddress: "adminUser@gmail.com",
					MailSendYn:  "Y",
					SnsSendYn:   "N",
				}

				reqBody = append(reqBody, query)
				reqBody = append(reqBody, query1)
				reqBody = append(reqBody, query2)
				reqBody = append(reqBody, query3)

				data, _ := json.Marshal(reqBody)

				res, _ := DoUpdate(testUrl + "/v2/iaas/alarm/policy", TestToken, strings.NewReader(string(data)))

				//assert.Nil(tt, err)
				assert.Equal(tt, http.StatusCreated, res.Code)
			})
		})


		Context("Alarm Sns Channel", func() {

			It("Alarm Sns Channel READ", func() {
				res, err := DoGet(testUrl + "/v2/iaas/alarm/sns/channel/list")

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("Alarm Sns Channel CREATE", func() {
				query := model.AlarmPolicyRequest{
					OriginType: "ias",
					SnsId:      "pas_test_bot",
					Token:      "595845637:AAGUgw96sfoyTO3RUZoX-i06OHh9ZX0yt",
					Expl:       "test pas bot",
					SnsSendYn:  "N",
				}

				data, _ := json.Marshal(query)
				res, _ := DoPost(testUrl + "/v2/iaas/alarm/sns/channel", TestToken, strings.NewReader(string(data)))
				fmt.Println(res)

				//assert.Nil(tt, err)
				//assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("Alarm Sns Channel DELETE", func() {
				query := model.AlarmPolicyRequest{
					MailAddress: "test@test.com",
					MailSendYn: "Y",
				}

				data, _ := json.Marshal(query)
				res, _ := DoDelete(testUrl + "/v2/iaas/alarm/sns/channel" + alarm_sns_id, TestToken, strings.NewReader(string(data)))
				fmt.Println(res)
				//assert.Nil(tt, err)
				//assert.Equal(tt, http.StatusOK, res.Code)
			})
		})
	})
})

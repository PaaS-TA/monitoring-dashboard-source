package controller

import (
	"encoding/json"
	"fmt"
	"monitoring-portal/saas/model"
	"strings"

	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"net/http"
)

type AlarmPolicyRequest struct {
	Id                uint   `json:"id"`
	OriginType        string `json:"originType"`
	AlarmType         string `json:"alarmType"`
	WarningThreshold  int    `json:"warningThreshold"`
	CriticalThreshold int    `json:"criticalThreshold"`
	RepeatTime        int    `json:"repeatTime"`
	Comment           string `json:"comment"`
	MeasureTime       int    `json:"measureTime"`
	MailAddress       string `json:"mailAddress"`
	SnsType           string `json:"snsType"`
	SnsId             string `json:"snsId"`
	Token             string `json:"token"`
	Expl              string `json:"expl"`
	MailSendYn        string `json:"mailSendYn"`
	SnsSendYn         string `json:"snsSendYn"`
}

var channel_id = "1"

var _ = Describe("Alarm API", func() {

	Describe("Saas Alarm", func() {
		Context("Alarm Policy", func() {
			It("SAAS_ALARM_POLICY_LIST", func() {
				res, err := DoGet(testUrl + "/v2/saas/app/application/alarmInfo")

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("SAAS_ALARM_POLICY_UPDATE (PUT)", func() {
				var reqBody []AlarmPolicyRequest
				query := AlarmPolicyRequest{
					OriginType:        "SaaS",
					AlarmType:         "SYSTEM_CPU",
					WarningThreshold:  85,
					CriticalThreshold: 91,
					RepeatTime:        10,
					MeasureTime:       100,
				}
				query1 := AlarmPolicyRequest{
					OriginType:        "SaaS",
					AlarmType:         "JVM_CPU",
					WarningThreshold:  85,
					CriticalThreshold: 91,
					RepeatTime:        10,
					MeasureTime:       10,
				}
				query2 := AlarmPolicyRequest{
					OriginType:        "SaaS",
					AlarmType:         "HEAP_MEMORY",
					WarningThreshold:  85,
					CriticalThreshold: 91,
					RepeatTime:        10,
					MeasureTime:       100,
				}
				query3 := AlarmPolicyRequest{
					OriginType:  "SaaS",
					MailAddress: "chotom73@daum.net",
					MailSendYn:  "Y",
					SnsSendYn:   "Y",
				}

				reqBody = append(reqBody, query)
				reqBody = append(reqBody, query1)
				reqBody = append(reqBody, query2)
				reqBody = append(reqBody, query3)

				data, _ := json.Marshal(reqBody)

				res, err := DoUpdate(testUrl+"/v2/saas/app/application/alarmUpdate", TestToken, strings.NewReader(string(data)))

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusCreated, res.Code)
			})

			It("SAAS_ALARM_SNS_CHANNEL_LIST", func() {
				res, err := DoGet(testUrl + "/v2/saas/app/application/snsChannel/list  ")

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("SAAS_ALARM_SNS_CHANNEL_CREATE(POST)", func() {

				query := model.BatchAlarmSnsRequest{
					OriginType: "SaaS",
					SnsId:      "SaaSMonitoringBot",
					Token:      "857394197:AAEnkrLk7S_-1dvcLCo1tSzYLomTkoPsUFA",
					Expl:       "test saas bot",
					SnsSendYn:  "Y",
				}

				data, _ := json.Marshal(query)

				res, err := DoPost(testUrl+"/v2/saas/app/application/alarmSnsSave", TestToken, strings.NewReader(string(data)))
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusCreated, res.Code)
			})
			//
			It("SAAS_ALARM_SNS_CHANNEL_DELETE(DELETE)", func() {
				res, err := DoDelete(testUrl+"/v2/saas/app/application/snsChannel/"+channel_id, TestToken, strings.NewReader(string("")))
				fmt.Println(res)
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})
		})

	})

})

package controller

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"time"

	//"github.com/monasca/golang-monascaclient/monascaclient/models"
	"encoding/json"
	. "github.com/onsi/ginkgo"
	"strings"
)

type AlarmPolicyRequest struct {
	Id                uint     `json:"id"`
	OriginType        string   `json:"originType"`
	AlarmType         string   `json:"alarmType"`
	WarningThreshold  int      `json:"warningThreshold"`
	CriticalThreshold int      `json:"criticalThreshold"`
	RepeatTime        int      `json:"repeatTime"`
	Comment           string   `json:"comment"`
	MeasureTime       int      `json:"measureTime"`
	MailAddress       string   `json:"mailAddress"`
	SnsType           string   `json:"snsType"`
	SnsId             string   `json:"snsId"`
	Token             string   `json:"token"`
	Expl              string   `json:"expl"`
	MailSendYn        string   `json:"mailSendYn"`
	SnsSendYn         string   `json:"snsSendYn"`
	ModiDate          time.Time `json:"modiDate"`
	ModiUser          string   `json:"modiUser"`
}

var _ = Describe("Potal Application API", func() {

	Describe("Potal Application Policy", func() {
		Context("Potal Application", func() {

			var appId = "52d78f94-5275-4a7e-8a67-bdf5e58a3246"

			It("PAAS_APP_AUTOSCALING_POLICY_INFO", func() {
				res, err := DoGet(testUrl + "/v2/paas/app/autoscaling/policy?appGuid=" + appId)

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("PAAS_APP_AUTOSCALING_POLICY_UPDATE (POST)", func() {

				query := AlarmPolicyRequest{

				}

				data, _ := json.Marshal(query)

				res, err := DoPost(testUrl+"/v2/paas/app/autoscaling/policy", TestToken, strings.NewReader(string(data)))

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("PAAS_APP_POLICY_INFO", func() {
				res, err := DoGet(testUrl + "/v2/paas/app/alarm/policy?appGuid=" + appId)

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("PAAS_APP_POLICY_UPDATE (POST)", func() {

				query := AppAlarmPolicyBody{
					AppGuid:                 appId,
					CpuWarningThreshold:     80,
					CpuCriticalThreshold:    90,
					MemoryWarningThreshold:  80,
					MemoryCriticalThreshold: 90,
					MeasureTimeSec:          500,
					Email:                   "test@test.com",
					EmailSendYn:             "N",
					AlarmUseYn:              "Y",
				}

				data, _ := json.Marshal(query)

				res, err := DoPost(testUrl+"/v2/paas/app/alarm/policy", TestToken, strings.NewReader(string(data)))

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("PAAS_APP_ALARM_LIST", func() {
				res, err := DoGet(testUrl + "/v2/paas/app/alarm/list?appGuid=" + appId + "&pageItems=10&pageIndex=1")

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("PAAS_APP_POLICY_DELETE (DELETE)", func() {
				res, err := DoDelete(testUrl+"/v2/paas/app/policy/"+appId, TestToken, strings.NewReader(string("")))

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

		})
	})

})

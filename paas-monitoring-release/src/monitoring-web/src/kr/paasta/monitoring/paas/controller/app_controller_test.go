package controller

import (
	"net/http"
	"github.com/stretchr/testify/assert"
	//"github.com/monasca/golang-monascaclient/monascaclient/models"
	. "github.com/onsi/ginkgo"
	"encoding/json"
	"strings"
)

type AppAutoscalingPolicyBody struct {
	AppGuid					string      `json:"appGuid"`
	InstanceMinCnt         	int   		`json:"instanceMinCnt"`
	InstanceMaxCnt          int  		`json:"instanceMaxCnt"`
	CpuMinThreshold   		int      	`json:"cpuMinThreshold"`
	CpuMaxThreshold  		int      	`json:"cpuMaxThreshold"`
	MemoryMinThreshold      int      	`json:"memoryMinThreshold"`
	MemoryMaxThreshold      int   		`json:"memoryMaxThreshold"`
	InstanceScalingUnit     int   	 	`json:"instanceVariationUnit"`
	MeasureTimeSec        	int   		`json:"measureTimeSec"`
	AutoScalingOutYn		string   	`json:"autoScalingOutYn"`
	AutoScalingInYn			string   	`json:"autoScalingInYn"`
}

type AppAlarmPolicyBody struct {
	AppGuid						string      `json:"appGuid"`
	CpuWarningThreshold        	int   		`json:"cpuWarningThreshold"`
	CpuCriticalThreshold        int  		`json:"cpuCriticalThreshold"`
	MemoryWarningThreshold   	int      	`json:"memoryWarningThreshold"`
	MemoryCriticalThreshold  	int      	`json:"memoryCriticalThreshold"`
	MeasureTimeSec      		int      	`json:"measureTimeSec"`
	Email      					string   	`json:"email"`
	EmailSendYn   				string   	`json:"emailSendYn"`
	AlarmUseYn        			string   	`json:"alarmUseYn"`
}

var _ = Describe("Potal Application API", func() {

	Describe("Potal Application Policy", func() {
		Context("Potal Application", func() {

			It("PAAS_APP_AUTOSCALING_POLICY_INFO", func() {
				res, err := DoGet(testUrl + "/v2/paas/app/autoscaling/policy?appGuid=" + appId)

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_APP_AUTOSCALING_POLICY_UPDATE (POST)", func() {

				query := AppAutoscalingPolicyBody{
					AppGuid: appId,
					InstanceMinCnt: 1,
					InstanceMaxCnt: 20,
					CpuMinThreshold: 50,
					CpuMaxThreshold: 80,
					MemoryMinThreshold: 50,
					MemoryMaxThreshold: 90,
					InstanceScalingUnit: 30,
					MeasureTimeSec: 500,
					AutoScalingOutYn: "N",
					AutoScalingInYn: "N",
				}

				data, _ := json.Marshal(query)

				res, err := DoPost(testUrl + "/v2/paas/app/autoscaling/policy", TestToken, strings.NewReader(string(data)))

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_APP_POLICY_INFO", func() {
				res, err := DoGet(testUrl + "/v2/paas/app/alarm/policy?appGuid=" + appId)

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_APP_POLICY_UPDATE (POST)", func() {

				query := AppAlarmPolicyBody{
					AppGuid: appId,
					CpuWarningThreshold: 80,
					CpuCriticalThreshold: 90,
					MemoryWarningThreshold: 80,
					MemoryCriticalThreshold: 90,
					MeasureTimeSec: 500,
					Email: "test@test.com",
					EmailSendYn: "N",
					AlarmUseYn: "Y",
				}

				data, _ := json.Marshal(query)

				res, err := DoPost(testUrl + "/v2/paas/app/alarm/policy", TestToken, strings.NewReader(string(data)))

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_APP_ALARM_LIST", func() {
				res, err := DoGet(testUrl + "/v2/paas/app/alarm/list?appGuid=" + appId + "&pageItems=10&pageIndex=1")

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("PAAS_APP_POLICY_DELETE (DELETE)", func() {
				res, err := DoDelete(testUrl + "/v2/paas/app/policy/" + appId, TestToken, strings.NewReader(string("")))

				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

		})
	})

})


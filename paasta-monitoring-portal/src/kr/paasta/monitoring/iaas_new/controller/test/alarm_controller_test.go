package test

import (
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
)

type AlarmPolicyRequestBody struct {
	Id                uint      `json:"id"`
	OriginType        string    `json:"originType"`
	AlarmType         string    `json:"alarmType"`
	WarningThreshold  int       `json:"warningThreshold"`
	CriticalThreshold int       `json:"criticalThreshold"`
	RepeatTime        int       `json:"repeatTime"`
	Comment           string    `json:"comment"`
	MeasureTime       int       `json:"measureTime"`
	MailAddress       string    `json:"mailAddress"`
	SnsType           string    `json:"snsType"`
	SnsId             string    `json:"snsId"`
	Token             string    `json:"token"`
	Expl              string    `json:"expl"`
	MailSendYn        string    `json:"mailSendYn"`
	SnsSendYn         string    `json:"snsSendYn"`
	ModiDate          time.Time `json:"modiDate"`
	ModiUser          string    `json:"modiUser"`
}

type AlarmRequestBody struct {
	Id             uint
	OriginType     string
	OriginId       uint
	AlarmType      string
	Level          string
	AlarmTitle     string
	ResolveStatus  string
	SearchDateFrom string
	SearchDateTo   string
}


var channel_id = "123"
var alarm_id = "1"
var resolveStatus = "2"
var action_id = "6"


var _ = Describe("Alarm API", func() {

	Describe("IaaS Alarm", func() {

		Context("Alarm Status", func() {
			It("IaaS alarm status list", func() {
				res, err := DoGet(testUrl + "/v2/iaas/alarm/statuses")

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("IaaS alarm status list count", func() {
				res, err := DoGet(testUrl + "/v2/iaas/alarm/status/count")

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("IaaS alarm status get", func() {
				res, err := DoGet(testUrl + "/v2/iaas/alarm/status/" + alarm_id)

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("IaaS alarm status get resolve status", func() {
				res, err := DoGet(testUrl + "/v2/iaas/alarm/status/" + resolveStatus)

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			/*
			It("IaaS alarm status update", func() {
				query := AlarmRequestBody{
					ResolveStatus: resolveStatus,
				}

				data, _ := json.Marshal(query)

				fmt.Println(">>>>>>>>>>> data :", data)
				fmt.Println(">>>>>>>>>>> TestToken :", TestToken)

				res, err := DoUpdate(testUrl+ "/v2/iaas/alarm/status/" +alarm_id, TestToken, strings.NewReader(string(data)))

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})
			*/
		})


		Context("Alarm Statistics", func() {

			It("PAAS_ALARM_STATISTICS", func() {
				res, err := DoGet(testUrl + "/v2/iaas/alarm/statistics?period=d&interval=1")

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("PAAS_ALARM_STATISTICS_GRAPH_TOTAL", func() {
				res, err := DoGet(testUrl + "/v2/iaas/alarm/statistics/graph/total?period=d&interval=1")

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("PAAS_ALARM_STATISTICS_GRAPH_SERVICE", func() {
				res, err := DoGet(testUrl + "/v2/iaas/alarm/statistics/graph/service?period=d&interval=1")

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("PAAS_ALARM_STATISTICS_GRAPH_MATRIX", func() {
				res, err := DoGet(testUrl + "/v2/iaas/alarm/statistics/graph/matrix?period=d&interval=1")

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})
		})

	})

})

package controller

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"kr/paasta/monitoring/caas/model"
	"net/http"
	"strconv"
	"strings"

	//"github.com/monasca/golang-monascaclient/monascaclient/models"
	. "github.com/onsi/ginkgo"
	"time"
)

//var channel_id = "123"
var alarm_id = "1"
var resolveStatus = "2"
var action_id = "2"

//var appId = "af9c7835-dd86-42f9-b105-dd4a3bae3f3c"
//var appIndex = "1"

var _ = Describe("Alarm API", func() {

	Describe("CaaS Alarm", func() {

		Context("Alarm Status", func() {
			It("CAAS_ALARM_REALTIME_COUNT", func() {
				res, err := DoGet(testUrl + "/v2/caas/monitoring/alarmCount")

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("CAAS_ALARM_TERM_COUNT", func() {
				toDay := time.Now().Format("2006-01-02")
				fromDay := time.Now().AddDate(0, -1, 0).Format("2006-01-02")
				res, err := DoGet(testUrl + "/v2/caas/monitoring/alarmCount?searchDateFrom=" + fromDay + "&searchDateTo=" + toDay)

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("CAAS_ALARM_STATUS_LIST", func() {
				res, err := DoGet(testUrl + "/v2/caas/monitoring/alarmLog")

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("CAAS_ALARM_STATUS_DETAIL", func() {
				res, err := DoGet(testUrl + "/v2/caas/monitoring/alarmAction/" + alarm_id)

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("CAAS_ALARM_STATUS_UPDATE (PUT)", func() {

				query := model.AlarmrRsolveRequest{
					ResolveStatus: resolveStatus,
				}

				data, _ := json.Marshal(query)

				fmt.Println(">>>>>>>>>>> data :", data)
				fmt.Println(">>>>>>>>>>> TestToken :", TestToken)

				res, err := DoUpdate(testUrl+"/v2/caas/monitoring/alarmStatus/"+alarm_id, TestToken, strings.NewReader(string(data)))

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusCreated, res.Code)
			})

			It("CAAS_ALARM_ACTION_CREATE (POST)", func() {
				id, _ := strconv.Atoi(alarm_id)
				query := model.AlarmrRsolveRequest{
					Id:              uint64(id),
					AlarmActionDesc: "test description",
				}

				data, _ := json.Marshal(query)

				res, err := DoPost(testUrl+"/v2/caas/monitoring/alarmAction", TestToken, strings.NewReader(string(data)))
				fmt.Println(res)
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusCreated, res.Code)
			})

			It("CAAS_ALARM_ACTION_UPDATE (PATCH)", func() {
				id, _ := strconv.Atoi(action_id)
				query := model.AlarmrRsolveRequest{
					Id:              uint64(id),
					AlarmActionDesc: "update test description",
				}

				data, _ := json.Marshal(query)

				res, err := DoPatch(testUrl+"/v2/caas/monitoring/alarmAction/"+action_id, TestToken, strings.NewReader(string(data)))
				fmt.Println(res)
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("CAAS_ALARM_ACTION_DELETE(DELETE)", func() {
				res, err := DoDelete(testUrl+"/v2/caas/monitoring/alarmAction/"+action_id, TestToken, strings.NewReader(string("")))
				fmt.Println(res)
				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusOK, res.Code)
			})

			It("CAAS_ALARM_STATUS_UPDATE_COMPLETE (PUT)", func() {

				query := model.AlarmrRsolveRequest{
					ResolveStatus: "3",
				}

				data, _ := json.Marshal(query)

				fmt.Println(">>>>>>>>>>> data :", data)
				fmt.Println(">>>>>>>>>>> TestToken :", TestToken)

				res, err := DoUpdate(testUrl+"/v2/caas/monitoring/alarmStatus/1", TestToken, strings.NewReader(string(data)))

				assert.Nil(tt, err)
				assert.Equal(tt, http.StatusCreated, res.Code)
			})
		})
	})

})

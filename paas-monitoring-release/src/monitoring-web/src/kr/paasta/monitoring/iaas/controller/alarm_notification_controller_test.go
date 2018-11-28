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

type NotificationRequestBody struct {
	Id    	string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Period  int    `json:"period,omitempty"`
	Type    string `json:"type,omitempty"`
	Address string `json:"address,omitempty"`
}

var notification_id = "c540e624-5bf3-4271-b79c-26e1605257ac"

var _ = Describe("IaaSAlarmNotificationController", func() {

	Describe("IaaS Alarm Notification", func() {

		Context("Get IaaS Alarm Notification List", func() {

			It("Notification List", func() {
				res, err := DoGet(testUrl + "/v2/iaas/alarm/notifications?offset=0&limit=10")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("Notification Create Info", func() {

				query := NotificationRequestBody{
					Name: "testname",
					Address: "testmail@gmail.com",
				}

				data, _ := json.Marshal(query)

				res, _ := DoPost(testUrl + "/v2/iaas/alarm/notification", TestToken, strings.NewReader(string(data)))
				fmt.Println(res)
				//assert.Nil(t, err)
				//assert.Equal(t, http.StatusOK, res.Code)


				notification_id = strings.Replace(res.Content, "\"", "", -1)
			})

			It("Notification Update Info", func() {
				var query NotificationRequestBody
				query.Name = "testname2"
				query.Address = "testmail2@gmail.com"
				query.Id = notification_id
				query.Period = 0
				data, _ := json.Marshal(query)

				res, _ := DoUpdate(testUrl + "/v2/iaas/alarm/notification/" + notification_id, TestToken, strings.NewReader(string(data)))
				fmt.Println(res)
				//assert.Nil(t, err)
				//assert.Equal(t, http.StatusOK, res.Code)
			})

			It("Notification Delete Info", func() {
				res, _ := DoDelete(testUrl + "/v2/iaas/alarm/notification/" + notification_id, TestToken, strings.NewReader(string("")))
				fmt.Println(res)
				//assert.Nil(t, err)
				//assert.Equal(t, http.StatusOK, res.Code)
			})

		})

	})

})


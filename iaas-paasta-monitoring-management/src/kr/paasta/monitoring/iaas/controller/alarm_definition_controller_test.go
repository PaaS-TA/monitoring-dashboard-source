package controller

import (
	"net/http"
	"github.com/stretchr/testify/assert"
	//"github.com/monasca/golang-monascaclient/monascaclient/models"
	. "github.com/onsi/ginkgo"
	"encoding/json"
	"strings"
)

type DefinitionRequestBody struct {
	Id	                string   `json:"id,omitempty"`
	Name                string   `json:"name,omitempty"`
	Description         string   `json:"description,omitempty"`
	Severity            string   `json:"severity,omitempty"`
	MatchBy             []string `json:"match_by,omitempty"`
	Expression          string   `json:"expression,omitempty"`
	ActionsEnabled      bool     `json:"actions_enabled,omitempty"`
	AlarmActions        []string `json:"alarm_actions,omitempty"`
	OkActions           []string `json:"ok_actions,omitempty"`
	UndeterminedActions []string `json:"undetermined_actions,omitempty"`
}

var definition_id = "952ad6b1-7858-4220-af33-93bcdb706fe7"
var action_id = "c540e624-5bf3-4271-b79c-26e1605257ac"

var _ = Describe("IaaS Alarm Policy Controller", func() {

	Describe("IaaS Alarm Policy List & Detail & Create & Patch & Delete", func() {

		Context("IaaS Alarm Policy View", func() {

			It("Policy List", func() {
				res, err := DoGet(testUrl + "/v2/iaas/alarm/policies?offset=0&limit=10")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("Policy Create", func() {
				var query DefinitionRequestBody

				query.Name = "mon_test4_mem"
				query.Expression = "max(mem.usable_perc) < 60"
				query.Severity = "CRITICAL"
				query.Id = definition_id
				query.AlarmActions = []string{action_id}
				query.MatchBy = []string{"hostname"}
				query.Description = "Test Description"

				data, _ := json.Marshal(query)

				res, err := DoPost(testUrl + "/v2/iaas/alarm/policy", TestToken, strings.NewReader(string(data)))
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				//fmt.Println(">>>>>>>>>>>>> res :: ",res)
			})

			It("Policy Detail", func() {
				//definition_id = "f6d76867-6004-4031-92d1-c6ba3897e9af"
				res, err := DoGet(testUrl + "/v2/iaas/alarm/policy/"+definition_id)
				//res, err := DoDetail(testUrl + "/v2/iaas/alarm/policy/"+definition_id, TestToken, strings.NewReader(string("")))
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("Policy Update", func() {
				definition_id = "d136c1e5-be64-471b-b7fc-95e18f4528d1"

				var query DefinitionRequestBody
				query.Expression = "max(mem.usable_perc) < 60"
				query.Description = ""
				query.Name = "node_MEM"
				query.Severity = "CRITICAL"
				query.Id = definition_id
				query.AlarmActions = []string{action_id}
				query.MatchBy = []string{"hostname"}

				data, _ := json.Marshal(query)

				res, err := DoPatch(testUrl + "/v2/iaas/alarm/policy/" + definition_id, TestToken, strings.NewReader(string(data)))
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				//fmt.Println(res)
			})

			It("Policy Delete", func() {
				definition_id = "d136c1e5-be64-471b-b7fc-95e18f4528d1"

				res, err := DoDelete(testUrl + "/v2/iaas/alarm/policy/" + definition_id, TestToken, strings.NewReader(string("")))
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
				//fmt.Println(res)
			})

		})

	})

})

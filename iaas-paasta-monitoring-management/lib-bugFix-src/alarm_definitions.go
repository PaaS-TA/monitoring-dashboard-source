// Copyright 2017 Hewlett Packard Enterprise Development LP
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package monascaclient

import (
	"github.com/monasca/golang-monascaclient/monascaclient/models"
)

const (
	alarmDefinitionsBasePath = "v2.0/alarm-definitions"
)

func GetAlarmDefinitions(alarmDefinitionQuery *models.AlarmDefinitionQuery) (*models.AlarmDefinitionsResponse, error) {
	return monClient.GetAlarmDefinitions(alarmDefinitionQuery)
}

func GetAlarmDefinition(alarmDefinitionID string) (*models.AlarmDefinitionElement, error) {
	return monClient.GetAlarmDefinition(alarmDefinitionID)
}

func CreateAlarmDefinition(alarmDefinitionRequestBody *models.AlarmDefinitionRequestBody) (*models.AlarmDefinitionElement, error) {
	return monClient.CreateAlarmDefinition(alarmDefinitionRequestBody)
}

func UpdateAlarmDefinition(alarmDefinitionID string, alarmDefinitionRequestBody *models.AlarmDefinitionRequestBody) (*models.AlarmDefinitionElement, error) {
	return monClient.UpdateAlarmDefinition(alarmDefinitionID, alarmDefinitionRequestBody)
}

func PatchAlarmDefinition(alarmDefinitionID string, alarmDefinitionRequestBody *models.AlarmDefinitionRequestBody) (*models.AlarmDefinitionElement, error) {
	return monClient.PatchAlarmDefinition(alarmDefinitionID, alarmDefinitionRequestBody)
}

func DeleteAlarmDefinition(alarmDefinitionID string) error {
	return monClient.DeleteAlarmDefinition(alarmDefinitionID)
}

func (c *Client) GetAlarmDefinitions(alarmDefinitionQuery *models.AlarmDefinitionQuery) (*models.AlarmDefinitionsResponse, error) {
	alarmDefinitionsResponse := new(models.AlarmDefinitionsResponse)
	err := c.callMonascaGet(alarmDefinitionsBasePath, "", alarmDefinitionQuery, alarmDefinitionsResponse)
	if err != nil {
		return nil, err
	}

	return alarmDefinitionsResponse, nil
}

func (c *Client) GetAlarmDefinition(alarmDefinitionID string) (*models.AlarmDefinitionElement, error) {
	alarmDefinitionElement := new(models.AlarmDefinitionElement)
	/******************** bug Fix start*****************/
	query := new(models.AlarmDefinitionQuery)
	/******************** bug Fix End*****************/

	err := c.callMonascaGet(alarmDefinitionsBasePath, alarmDefinitionID, query, alarmDefinitionElement)
	if err != nil {
		return nil, err
	}

	return alarmDefinitionElement, nil
}

func (c *Client) CreateAlarmDefinition(alarmDefinitionRequestBody *models.AlarmDefinitionRequestBody) (*models.AlarmDefinitionElement, error) {
	return c.sendAlarmDefinition("", "POST", alarmDefinitionRequestBody)
}

func (c *Client) UpdateAlarmDefinition(alarmDefinitionID string, alarmDefinitionRequestBody *models.AlarmDefinitionRequestBody) (*models.AlarmDefinitionElement, error) {
	return c.sendAlarmDefinition(alarmDefinitionID, "PUT", alarmDefinitionRequestBody)
}

func (c *Client) PatchAlarmDefinition(alarmDefinitionID string, alarmDefinitionRequestBody *models.AlarmDefinitionRequestBody) (*models.AlarmDefinitionElement, error) {
	return c.sendAlarmDefinition(alarmDefinitionID, "PATCH", alarmDefinitionRequestBody)
}

func (c *Client) sendAlarmDefinition(alarmDefinitionID string, method string, alarmDefinitionRequestBody *models.AlarmDefinitionRequestBody) (*models.AlarmDefinitionElement, error) {
	alarmDefinitionsElement := new(models.AlarmDefinitionElement)
	err := c.callMonascaWithBody(alarmDefinitionsBasePath, alarmDefinitionID, method, alarmDefinitionRequestBody, alarmDefinitionsElement)
	if err != nil {
		return nil, err
	}

	return alarmDefinitionsElement, nil
}

func (c *Client) DeleteAlarmDefinition(alarmDefinitionID string) error {
	return c.callMonascaDelete(alarmDefinitionsBasePath, alarmDefinitionID)
}

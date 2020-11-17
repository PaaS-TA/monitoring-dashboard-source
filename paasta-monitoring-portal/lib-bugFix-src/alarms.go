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
	//"fmt"
)

const (
	alarmsBasePath = "v2.0/alarms"
)

func GetAlarms(alarmQuery *models.AlarmQuery) (*models.AlarmsResponse, error) {
	return monClient.GetAlarms(alarmQuery)
}

func GetAlarm(alarmID string) (*models.Alarm, error) {
	return monClient.GetAlarm(alarmID)
}

func UpdateAlarm(alarmID string, alarmRequestBody *models.AlarmRequestBody) (*models.Alarm, error) {
	return monClient.UpdateAlarm(alarmID, alarmRequestBody)
}

func PatchAlarm(alarmID string, alarmRequestBody *models.AlarmRequestBody) (*models.Alarm, error) {
	return monClient.PatchAlarm(alarmID, alarmRequestBody)
}

func DeleteAlarm(alarmID string) error {
	return monClient.DeleteAlarm(alarmID)
}

func (c *Client) GetAlarmCount(alarmQuery *models.AlarmQuery) ([][]int, error) {

	type alarmCount struct {
		Count [][]int `json:"counts"`
	}
	count := new(alarmCount)

	err := c.callMonascaGet(alarmsBasePath+"/count", "", alarmQuery, count)
	if err != nil {
		return nil, err
	}
	//fmt.Println("==>",count)
	return count.Count, nil
}

func (c *Client) GetAlarms(alarmQuery *models.AlarmQuery) (*models.AlarmsResponse, error) {
	alarmsResponse := new(models.AlarmsResponse)
	err := c.callMonascaGet(alarmsBasePath, "", alarmQuery, alarmsResponse)
	if err != nil {

		return nil, err
	}

	return alarmsResponse, nil
}

func (c *Client) GetAlarm(alarmID string) (*models.Alarm, error) {

	alarm := new(models.Alarm)
	query := new(models.AlarmQuery)

	err := c.callMonascaGet(alarmsBasePath, alarmID, query, alarm)

	if err != nil {
		return nil, err
	}

	return alarm, nil
}

func (c *Client) UpdateAlarm(alarmID string, alarmRequestBody *models.AlarmRequestBody) (*models.Alarm, error) {
	return c.sendAlarm(alarmID, "PUT", alarmRequestBody)
}

func (c *Client) PatchAlarm(alarmID string, alarmRequestBody *models.AlarmRequestBody) (*models.Alarm, error) {
	return c.sendAlarm(alarmID, "PATCH", alarmRequestBody)
}

func (c *Client) sendAlarm(alarmID string, method string, alarmRequestBody *models.AlarmRequestBody) (*models.Alarm, error) {
	alarmsElement := new(models.Alarm)
	err := c.callMonascaWithBody(alarmsBasePath, alarmID, method, alarmRequestBody, alarmsElement)
	if err != nil {
		return nil, err
	}

	return alarmsElement, nil
}

func (c *Client) DeleteAlarm(alarmID string) error {
	return c.callMonascaDelete(alarmsBasePath, alarmID)
}

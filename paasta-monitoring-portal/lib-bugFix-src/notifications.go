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
	notificationsBasePath = "v2.0/notification-methods"
)

func GetNotificationMethods(notificationQuery *models.NotificationQuery) (*models.NotificationResponse, error) {
	return monClient.GetNotificationMethods(notificationQuery)
}

func GetNotificationMethod(notificationMethodID string, notificationQuery *models.NotificationQuery) (*models.NotificationElement, error) {
	return monClient.GetNotificationMethod(notificationMethodID, notificationQuery)
}

func CreateNotificationMethod(notificationRequestBody *models.NotificationRequestBody) (*models.NotificationElement, error) {
	return monClient.CreateNotificationMethod(notificationRequestBody)
}

func UpdateNotificationMethod(notificationID string, notificationRequestBody *models.NotificationRequestBody) (*models.NotificationElement, error) {
	return monClient.UpdateNotificationMethod(notificationID, notificationRequestBody)
}

func PatchNotificationMethod(notificationID string, notificationRequestBody *models.NotificationRequestBody) (*models.NotificationElement, error) {
	return monClient.PatchNotificationMethod(notificationID, notificationRequestBody)
}

func DeleteNotificationMethod(notificationID string) error {
	return monClient.DeleteNotificationMethod(notificationID)
}

func (c *Client) GetNotificationMethods(notificationQuery *models.NotificationQuery) (*models.NotificationResponse, error) {
	notificationsResponse := new(models.NotificationResponse)
	err := c.callMonascaGet(notificationsBasePath, "", notificationQuery, notificationsResponse)
	if err != nil {
		return nil, err
	}

	return notificationsResponse, nil
}

func (c *Client) GetNotificationMethod(notificationMethodID string, notificationQuery *models.NotificationQuery) (*models.NotificationElement, error) {
	notificationElement := new(models.NotificationElement)
	/******************** bug Fix start*****************/
	err := c.callMonascaGet(notificationsBasePath, notificationMethodID, notificationQuery, notificationElement)
	//err := c.callMonascaGet(notificationsBasePath, notificationMethodID, nil, notificationElement)
	/******************** bug Fix start*****************/
	if err != nil {
		return nil, err
	}

	return notificationElement, nil
}

func (c *Client) CreateNotificationMethod(notificationRequestBody *models.NotificationRequestBody) (*models.NotificationElement, error) {
	return c.sendNotification("", "POST", notificationRequestBody)
}

func (c *Client) UpdateNotificationMethod(notificationID string, notificationRequestBody *models.NotificationRequestBody) (*models.NotificationElement, error) {
	return c.sendNotification(notificationID, "PUT", notificationRequestBody)
}

func (c *Client) PatchNotificationMethod(notificationID string, notificationRequestBody *models.NotificationRequestBody) (*models.NotificationElement, error) {
	return c.sendNotification(notificationID, "PATCH", notificationRequestBody)
}

func (c *Client) sendNotification(notificationID string, method string, notificationRequestBody *models.NotificationRequestBody) (*models.NotificationElement, error) {
	notificationsElement := new(models.NotificationElement)
	err := c.callMonascaWithBody(notificationsBasePath, notificationID, method, notificationRequestBody, notificationsElement)
	if err != nil {
		return nil, err
	}

	return notificationsElement, nil
}

func (c *Client) DeleteNotificationMethod(notificationID string) error {
	return c.callMonascaDelete(notificationsBasePath, notificationID)
}

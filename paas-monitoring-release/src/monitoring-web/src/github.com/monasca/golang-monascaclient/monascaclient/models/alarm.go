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

package models

import "time"

type AlarmsResponse struct {
	Links    []Link  `json:"links"`
	Elements []Alarm `json:"elements"`
}

type Alarm struct {
	Metrics               []Metric  `json:"metrics"`
	LifecycleState        string    `json:"lifecycle_state"`
	State                 string    `json:"state"`
	Link                  string    `json:"link"`
	UpdatedTimestamp      time.Time `json:"updated_timestamp"`
	CreatedTimestamp      time.Time `json:"created_timestamp"`
	StateUpdatedTimestamp time.Time `json:"state_updated_timestamp"`
	ResponseElement
}

type AlarmQuery struct {
	AlarmDefinitionID     *string            `queryParameter:"alarm_definition_id"`
	MetricName            *string            `queryParameter:"metric_name"`
	MetricDimensions      *map[string]string `queryParameter:"metric_dimensions"`
	State                 *string            `queryParameter:"state"`
	Severity              *string            `queryParameter:"severity"`
	LifecycleState        *string            `queryParameter:"lifecycle_state"`
	Link                  *string            `queryParameter:"link"`
	StateUpdatedStartTime *time.Time         `queryParameter:"state_updated_start_time"`
	SortBy                *string            `queryParameter:"sort_by"`
	Offset                *int               `queryParameter:"offset"`
	Limit                 *int               `queryParameter:"limit"`
}

type AlarmRequestBody struct {
	State          *string `json:"state,omitempty"`
	LifecycleState *string `json:"lifecycle_state,omitempty"`
	Link           *string `json:"link,omitempty"`
}

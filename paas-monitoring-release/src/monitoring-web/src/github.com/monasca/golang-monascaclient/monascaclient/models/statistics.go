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

type StatisticsResponse struct {
	Links    []Link             `json:"links"`
	Elements []StatisticElement `json:"elements"`
}

type StatisticElement struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Dimensions map[string]string `json:"dimensions"`
	Columns    []string          `json:"columns"`
	Statistics [][]interface{}   `json:"statistics"`
}

type StatisticQuery struct {
	TenantID   *string            `queryParameter:"tenant_id"`
	Name       *string            `queryParameter:"name"`
	Dimensions *map[string]string `queryParameter:"dimensions"`
	Statistics *string            `queryParameter:"statistics"`
	StartTime  *time.Time         `queryParameter:"start_time"`
	EndTime    *time.Time         `queryParameter:"end_time"`
	Period     *int               `queryParameter:"period"`
	Offset     *int               `queryParameter:"offset"`
	Limit      *int               `queryParameter:"limit"`
	Merge      *bool              `queryParameter:"merge_metrics"`
	GroupBy    *string            `queryParameter:"group_by"`
}

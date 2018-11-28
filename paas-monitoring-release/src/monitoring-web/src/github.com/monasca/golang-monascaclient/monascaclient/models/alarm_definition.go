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

type AlarmDefinitionsResponse struct {
	Links    []Link                   `json:"links"`
	Elements []AlarmDefinitionElement `json:"elements"`
}

type AlarmDefinitionElement struct {
	AlarmDefinition
	// ActionsEnabled 	bool     `json:"actions_enabled"`
	Deterministic bool `json:"deterministic"`
	ResponseElement
}

type AlarmDefinition struct {
	Description         string   `json:"description,omitempty"`
	Severity            string   `json:"severity,omitempty"`
	AlarmActions        []string `json:"alarm_actions,omitempty"`
	OkActions           []string `json:"ok_actions,omitempty"`
	MatchBy             []string `json:"match_by,omitempty"`
	UndeterminedActions []string `json:"undetermined_actions,omitempty"`
	Expression          string   `json:"expression,omitempty"`
	Name                string   `json:"name,omitempty"`
}

type AlarmDefinitionQuery struct {
	Name       *string            `queryParameter:"name"`
	Dimensions *map[string]string `queryParameter:"dimensions"`
	Severity   *string            `queryParameter:"severity"`
	SortBy     *string            `queryParameter:"sort_by"`
	Offset     *int               `queryParameter:"offset"`
	Limit      *int               `queryParameter:"limit"`
}

type AlarmDefinitionRequestBody struct {
	Description         *string   `json:"description,omitempty"`
	Severity            *string   `json:"severity,omitempty"`
	AlarmActions        *[]string `json:"alarm_actions,omitempty"`
	OkActions           *[]string `json:"ok_actions,omitempty"`
	MatchBy             *[]string `json:"match_by,omitempty"`
	UndeterminedActions *[]string `json:"undetermined_actions,omitempty"`
	Expression          *string   `json:"expression,omitempty"`
	Name                *string   `json:"name,omitempty"`
	ActionsEnabled      *bool     `json:"actions_enabled,omitempty"`
}

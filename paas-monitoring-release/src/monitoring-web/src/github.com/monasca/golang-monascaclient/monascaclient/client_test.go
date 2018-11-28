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
	"fmt"
	"net/url"
	"testing"
)

func TestUrlCreation(t *testing.T) {
	const baseURL = "http://fred.com:7072"
	SetBaseURL(baseURL)
	const path = "v2.0/metrics/statistics"
	const metricName = "cpu.idle_perc"
	const period = 3600
	const startTime = "2017-02-27T06:00:00Z"
	const endTime = "2017-02-27T08:00:00Z"
	urlValues := url.Values{}
	urlValues.Add("name", metricName)
	urlValues.Add("statistics", "avg")
	urlValues.Add("start_time", startTime)
	urlValues.Add("end_time", endTime)
	urlValues.Add("period", fmt.Sprintf("%d", period))
	urlValues.Add("dimensions", "aggregation_period:hourly,host:all")

	monascaURL, err := monClient.createMonascaAPIURL(path, urlValues)
	if err != nil {
		t.Errorf("Error %s when creating Monasca API", err.Error())
	}
	expected := baseURL + "/" + path + "?" + "dimensions=aggregation_period%3Ahourly%2Chost%3Aall&end_time=2017-02-27T08%3A00%3A00Z&name=cpu.idle_perc&period=3600&start_time=2017-02-27T06%3A00%3A00Z&statistics=avg"
	if expected != monascaURL {
		t.Errorf("Expected '%v' but was '%v'", expected, monascaURL)
	}
}

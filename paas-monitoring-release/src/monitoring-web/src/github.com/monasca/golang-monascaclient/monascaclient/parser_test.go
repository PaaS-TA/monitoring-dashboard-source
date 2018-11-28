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
	"net/url"
	"reflect"
	"testing"
	"time"
)

type TestStruct struct {
	TestString *string            `queryParameter:"test_string"`
	TestTime   *time.Time         `queryParameter:"test_time"`
	TestMap    *map[string]string `queryParameter:"test_map"`
	TestInt    *int               `queryParameter:"test_int"`
}

func TestStructConversionEmptyStruct(t *testing.T) {
	testStructInput := TestStruct{}
	urlValuesExpected := url.Values{}
	urlValuesReturned := convertStructToQueryParameters(&testStructInput)
	if !reflect.DeepEqual(urlValuesExpected, urlValuesReturned) {
		t.Errorf("URL Values %s returned from method do not match expect values %s ", urlValuesReturned,
			urlValuesExpected)
	}
}

func TestStructConversion(t *testing.T) {
	inputString := "inputstring"
	inputTime := time.Now()
	inputMap := map[string]string{"key1": "value1", "key2": "value2", "key3": "value3"}
	inputInt := 3
	testStructInput := TestStruct{&inputString, &inputTime, &inputMap,
		&inputInt}
	urlValuesExpected := url.Values{}
	urlValuesExpected.Add("test_string", "inputstring")
	urlValuesExpected.Add("test_time", inputTime.UTC().Format(timeFormat))
	urlValuesExpected.Add("test_map", "key1:value1,key2:value2,key3:value3")
	urlValuesExpected.Add("test_int", "3")
	urlValuesReturned := convertStructToQueryParameters(&testStructInput)
	if !reflect.DeepEqual(urlValuesExpected, urlValuesReturned) {
		t.Errorf("URL Values %s returned from method do not match expect values %s ", urlValuesReturned,
			urlValuesExpected)
	}
}

func TestStructConversionPartiallyDefinedStruct(t *testing.T) {
	inputMap := map[string]string{"key1": "value1", "key2": "value2", "key3": "value3"}
	inputInt := 3
	testStructInput := TestStruct{TestMap: &inputMap, TestInt: &inputInt}
	urlValuesExpected := url.Values{}
	urlValuesExpected.Add("test_map", "key1:value1,key2:value2,key3:value3")
	urlValuesExpected.Add("test_int", "3")
	urlValuesReturned := convertStructToQueryParameters(&testStructInput)
	if !reflect.DeepEqual(urlValuesExpected, urlValuesReturned) {
		t.Errorf("URL Values %s returned from method do not match expect values %s ", urlValuesReturned,
			urlValuesExpected)
	}
}

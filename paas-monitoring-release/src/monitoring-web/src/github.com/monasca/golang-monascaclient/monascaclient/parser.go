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
	"reflect"
	"sort"
	"strings"
	"time"
)

func convertStructToQueryParameters(inputStruct interface{}) url.Values {
	urlValues := url.Values{}
	values := reflect.ValueOf(inputStruct)
	if values.IsNil() {
		return urlValues
	}
	values = values.Elem()
	typ := values.Type()
	// Loop through the struct
	for i := 0; i < typ.NumField(); i++ {
		currentValue := values.Field(i)
		currentType := typ.Field(i)
		// Get Query Parameter Name
		queryParameterKey := currentType.Tag.Get("queryParameter")
		if currentValue.Kind() == reflect.Ptr {
			if currentValue.IsNil() {
				continue
			}
			currentValue = currentValue.Elem()
		}
		addQueryParameter(currentValue, queryParameterKey, &urlValues)
	}
	return urlValues
}

func addQueryParameter(value reflect.Value, key string, values *url.Values) {
	if value.Kind() == reflect.Bool || value.Kind() == reflect.String || value.Kind() == reflect.Int {
		(*values).Add(key, fmt.Sprint(value.Interface()))
	} else if value.Type() == reflect.TypeOf(time.Time{}) {
		timeValue := value.Interface().(time.Time)
		(*values).Add(key, timeValue.UTC().Format(timeFormat))
	} else if value.Kind() == reflect.Map {
		mapValue := value.Interface().(map[string]string)
		if len(mapValue) > 0 {
			dimensionsSlice := make([]string, 0, len(mapValue))
			for key := range mapValue {
				dimensionsSlice = append(dimensionsSlice, key+":"+mapValue[key])
			}
			// Make sure dimensions are always in correct order to ensure tests pass
			sort.Strings(dimensionsSlice)
			(*values).Add(key, strings.Join(dimensionsSlice, ","))
		}
	}
}

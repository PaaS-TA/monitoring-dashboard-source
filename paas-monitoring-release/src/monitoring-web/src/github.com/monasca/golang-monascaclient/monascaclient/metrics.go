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
	"encoding/json"
	"github.com/monasca/golang-monascaclient/monascaclient/models"
	"net/url"
)

const (
	timeFormat      = "2006-01-02T15:04:05Z"
	metricsBasePath = "v2.0/metrics"
)

func GetMetrics(metricQuery *models.MetricQuery) ([]models.Metric, error) {
	return monClient.GetMetrics(metricQuery)
}

func GetMetricNames(metricQuery *models.MetricNameQuery) ([]string, error) {
	return monClient.GetMetricNames(metricQuery)
}

func GetDimensionValues(dimensionQuery *models.DimensionValueQuery) ([]string, error) {
	return monClient.GetDimensionValues(dimensionQuery)
}

func GetDimensionNames(dimensionQuery *models.DimensionNameQuery) ([]string, error) {
	return monClient.GetDimensionNames(dimensionQuery)
}

func GetStatistics(statisticsQuery *models.StatisticQuery) (*models.StatisticsResponse, error) {
	return monClient.GetStatistics(statisticsQuery)
}

func GetMeasurements(measurementQuery *models.MeasurementQuery) (*models.MeasurementsResponse, error) {
	return monClient.GetMeasurements(measurementQuery)
}

func CreateMetric(tenantID *string, metricRequestBody *models.MetricRequestBody) error {
	return monClient.CreateMetric(tenantID, metricRequestBody)
}

func (c *Client) CreateMetric(tenantID *string, metricRequestBody *models.MetricRequestBody) error {
	urlValues := url.Values{}
	if tenantID != nil {
		urlValues.Add("tenant_id", *tenantID)
	}
	monascaURL, URLerr := c.createMonascaAPIURL(metricsBasePath, urlValues)
	if URLerr != nil {
		return URLerr
	}
	byteInput, marshalErr := json.Marshal(*metricRequestBody)
	if marshalErr != nil {
		return marshalErr
	}
	return c.callMonascaNoContent(monascaURL, "POST", &byteInput)
}

func (c *Client) GetMetrics(metricQuery *models.MetricQuery) ([]models.Metric, error) {
	metricsResponse := new(models.MetricsResponse)
	err := c.callMonascaGet(metricsBasePath, "", metricQuery, metricsResponse)
	if err != nil {
		return []models.Metric{}, err
	}

	return metricsResponse.Elements, nil
}

func (c *Client) GetDimensionValues(dimensionQuery *models.DimensionValueQuery) ([]string, error) {
	return c.getDimensionQuery("/dimensions/names/values", dimensionQuery)
}

func (c *Client) GetDimensionNames(dimensionQuery *models.DimensionNameQuery) ([]string, error) {
	return c.getDimensionQuery("/dimensions/names/names", dimensionQuery)
}

func (c *Client) getDimensionQuery(path string, dimensionQuery interface{}) ([]string, error) {
	response := new(models.DimensionValueResponse)
	err := c.callMonascaGet(metricsBasePath+path, "", dimensionQuery, response)
	if err != nil {
		return []string{}, err
	}

	results := []string{}
	for _, value := range response.Elements {
		results = append(results, value.Value)
	}

	return results, nil
}

func (c *Client) GetMetricNames(metricQuery *models.MetricNameQuery) ([]string, error) {
	response := new(models.MetricNameResponse)
	err := c.callMonascaGet(metricsBasePath+"/names", "", metricQuery, response)
	if err != nil {
		return []string{}, err
	}

	results := []string{}
	for _, value := range response.Elements {
		results = append(results, value["name"])
	}
	return results, nil
}

func (c *Client) GetStatistics(statisticsQuery *models.StatisticQuery) (*models.StatisticsResponse, error) {
	statisticsResponse := new(models.StatisticsResponse)
	err := c.callMonascaGet(metricsBasePath+"/statistics", "", statisticsQuery, statisticsResponse)
	if err != nil {
		return nil, err
	}

	return statisticsResponse, nil
}

func (c *Client) GetMeasurements(measurementsQuery *models.MeasurementQuery) (*models.MeasurementsResponse, error) {
	measurementsResponse := new(models.MeasurementsResponse)
	err := c.callMonascaGet(metricsBasePath+"/measurements", "", measurementsQuery, measurementsResponse)
	if err != nil {
		return nil, err
	}

	return measurementsResponse, nil
}

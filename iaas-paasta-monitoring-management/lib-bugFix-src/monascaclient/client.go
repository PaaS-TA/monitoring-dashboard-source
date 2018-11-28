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
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	defaultURL      = "http://localhost:8070"
	defaultTimeout  = 60
	defaultInsecure = false
)

var (
	monClient = &Client{
		baseURL:        defaultURL,
		requestTimeout: defaultTimeout,
		allowInsecure:  defaultInsecure,
		headers:        http.Header{},
	}
)

func SetBaseURL(url string) {
	monClient.SetBaseURL(url)
}

func SetDefaultBaseURL(url string) {
	defaultURL = url
}

func SetInsecure(insecure bool) {
	monClient.SetInsecure(insecure)
}

func SetDefaultInsecure(insecure bool) {
	defaultInsecure = insecure
}

func SetTimeout(timeout int) {
	monClient.SetTimeout(timeout)
}

func SetDefaultTimeout(timeout int) {
	defaultTimeout = timeout
}

func SetHeaders(headers http.Header) {
	monClient.SetHeaders(headers)
}

func SetKeystoneConfig(config *gophercloud.AuthOptions) {
	monClient.SetKeystoneConfig(config)
}

type Client struct {
	baseURL        string
	requestTimeout int
	allowInsecure  bool
	headers        http.Header
	keystoneConfig *gophercloud.AuthOptions
}

func New() *Client {
	return &Client{
		baseURL:        defaultURL,
		requestTimeout: defaultTimeout,
		allowInsecure:  defaultInsecure,
		headers:        http.Header{},
	}
}

func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}

func (c *Client) SetInsecure(insecure bool) {
	c.allowInsecure = insecure
}

func (c *Client) SetTimeout(timeout int) {
	c.requestTimeout = timeout
}

func (c *Client) SetHeaders(headers http.Header) {
	c.headers = headers
}

func (c *Client) SetKeystoneConfig(config *gophercloud.AuthOptions) error {
	if config == nil {
		tmpConfig, err := openstack.AuthOptionsFromEnv()
		if err != nil {
			return fmt.Errorf("Failed to get keystone config from env: %v", err)
		}
		config = &tmpConfig
	}
	c.keystoneConfig = config
	return nil
}

func (c *Client) callMonasca(monascaURL string, method string, requestBody *[]byte) (*http.Response, error) {
	var req *http.Request
	var reqErr error

	if requestBody == nil {
		req, reqErr = http.NewRequest(method, monascaURL, nil)
	} else {
		req, reqErr = http.NewRequest(method, monascaURL, bytes.NewBuffer(*requestBody))
	}

	if reqErr != nil {
		return nil, reqErr
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if(c.keystoneConfig != nil){
		c.setKeystoneToken()
		c.headers.Set("X-Auth-Token",c.keystoneConfig.TokenID)
	}

	c.applyHeaders(req)

	timeout := time.Duration(c.requestTimeout) * time.Second
	var client *http.Client
	if !c.allowInsecure {
		client = &http.Client{Timeout: timeout}
	} else {
		transCfg := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
		}

		client = &http.Client{Timeout: timeout, Transport: transCfg}
	}
	resp, respErr := client.Do(req)

	// If response is 401, check for expired token and retry
	if respErr == nil && resp != nil && resp.StatusCode == 401 && c.keystoneConfig != nil {
		c.setKeystoneToken()
		c.headers.Set("X-Auth-Token",c.keystoneConfig.TokenID)
		c.applyHeaders(req)
		resp, respErr = client.Do(req)
	}

	return resp, respErr
}

func (c *Client) applyHeaders(req *http.Request) {
	for header, values := range c.headers {
		for index := range values {
			value := values[index]
			if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
				value = value[1 : len(value)-1]
			}
			req.Header.Set(header, value)
		}
	}
}

func (c *Client) callMonascaNoContent(monascaURL string, method string, requestBody *[]byte) error {
	resp, err := c.callMonasca(monascaURL, method, requestBody)
	if err != nil || resp == nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		return fmt.Errorf("Error: %d %s", resp.StatusCode, body)
	}
	return nil
}

func makePath(basePath string, id string) string {
	path := basePath
	if id == "" {
		return path
	}
	return path + "/" + id
}

func (c *Client) callMonascaGet(basePath string, id string, queryStruct interface{}, returned interface{}) error {

	urlValues := convertStructToQueryParameters(queryStruct)

	monascaURL, URLerr := c.createMonascaAPIURL(makePath(basePath, id), urlValues)
	if URLerr != nil {
		return URLerr
	}

	body, monascaErr := c.callMonascaReturnBody(monascaURL, "GET", nil)

	if (len(body)==0){
        //조회결과가 없는경우 err 가 아니고 nil 을 리턴하도록 변경
		//fmt.Println("len(body)==0")
		return nil

	}else if monascaErr != nil {

		return monascaErr

	}

	err := json.Unmarshal(body, returned)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) callMonascaWithBody(basePath string, id string, method string, toSend interface{}, returned interface{}) error {
	monascaURL, URLerr := c.createMonascaAPIURL(makePath(basePath, id), nil)
	if URLerr != nil {
		return URLerr
	}

	byteInput, marshalErr := json.Marshal(toSend)
	if marshalErr != nil {
		return marshalErr
	}
	body, monascaErr := c.callMonascaReturnBody(monascaURL, method, &byteInput)
	if monascaErr != nil {
		return monascaErr
	}

	err := json.Unmarshal(body, returned)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) callMonascaDelete(path string, id string) error {
	monascaURL, URLerr := c.createMonascaAPIURL(path+"/"+id, nil)
	if URLerr != nil {
		return URLerr
	}

	return c.callMonascaNoContent(monascaURL, "DELETE", nil)
}

func (c *Client) callMonascaReturnBody(monascaURL string, method string, requestBody *[]byte) ([]byte, error) {
	resp, err := c.callMonasca(monascaURL, method, requestBody)
	if err != nil || resp == nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, fmt.Errorf("Error: %d %s", resp.StatusCode, body)
	}

	return body, nil
}

func (c *Client) createMonascaAPIURL(path string, urlValues url.Values) (string, error) {

	monascaURL, parseErr := url.Parse(c.baseURL)
	if parseErr != nil {
		return "", parseErr
	}
	monascaURL.Path = path

	if urlValues != nil {
		monascaURL.RawQuery = urlValues.Encode()
	}

	return monascaURL.String(), nil
}

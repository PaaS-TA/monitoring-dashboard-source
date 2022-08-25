package util

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	md "iaas-monitoring-batch/model"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	//"fmt"
)

type UaaToken struct {
	Token            string `json:"access_token"`
	Scope            string `json:"scope"`
	Expire           int64  `json:"expires_in"`
	Refresh          string `json:"refresh_token"`
	ExpireTime       time.Time
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

/**
	AccessToken 발급을 위한 API 호출 함수
 */
func GetUaaToken(n md.CFConfig) (token md.UaaToken) {

	if len(n.UserId) > 0 && len(n.UserPw) > 0 {
		apiUrl := n.Host
		resource := "/oauth/token"
		data := url.Values{}
		data.Set("grant_type", "password")
		data.Set("response_type", "token")
		//data.Set("token_format", "opaque")
		data.Set("username", n.UserId)
		data.Set("password", n.UserPw)
		data.Set("client_id", n.ClientId)
		//data.Set("client_secret", n.ClientPw)

		u, _ := url.ParseRequestURI(apiUrl)
		u.Path = resource
		urlStr := u.String() // "https://api.com/user/"

		fmt.Println(">>>>> [cf_util.go / GetUaaToken()] API URL : " + urlStr)

		tp := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

		client := &http.Client{Transport: tp}
		r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
		r.Header.Add("Accept", "application/json")
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
		r.SetBasicAuth(url.QueryEscape("cf"), url.QueryEscape(""))
		resp, err := client.Do(r)

		fmt.Println("====== cf_util.go / GetUaaToken() - Request to ", urlStr, resp.Status)

		if err != nil {
			fmt.Println(err)
			time.Sleep(3 * time.Second)
			GetUaaToken(n)
		}

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		//bodyString := string(bodyBytes)
		//fmt.Println("====== cf_util.go / GetUaaToken() - response Body : ", bodyString)

		var uaaToken UaaToken
		json.Unmarshal(bodyBytes, &uaaToken)

		fmt.Println("====== cf_util.go / GetUaaToken() - uaaToken.Token : ", uaaToken.Token)

		if resp.StatusCode != http.StatusOK {
			fmt.Println(uaaToken)
			//return uaaToken.ErrorDescription
		}

		t := time.Now().Local().Add(time.Second * time.Duration(uaaToken.Expire))
		token.Token = uaaToken.Token
		token.ExpireTime = t
		token.Refresh = uaaToken.Refresh
		//token = token{
		//	Token:uaaToken.Token,
		//	ExpireTime:uaaToken.ExpireTime,
		//	Refresh:uaaToken.Refresh,
		//}

		//uaaToken.ExpireTime = t
		defer resp.Body.Close()

	} else {
		//result = "parameter is none"
	}

	return token
}

func GetUaaReFreshToken(n md.CFConfig, refreshToken string) (string, string) {

	result := ""
	//HttpClient := http.DefaultClient

	apiUrl := n.Host
	resource := "/uaa/oauth/token"
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("token_format", "opaque")
	data.Set("client_id", n.ClientId)
	data.Set("client_secret", n.ClientPw)
	data.Set("refresh_token", refreshToken)

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String() // "https://api.com/user/"

	fmt.Println("====== cf_util.go / GetUaaReFreshToken() - API URL : " + urlStr)

	tp := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	client := &http.Client{Transport: tp}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, _ := client.Do(r)
	fmt.Println(resp.Status)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)

	var uaaToken UaaToken
	json.Unmarshal(bodyBytes, &uaaToken)

	if resp.StatusCode != http.StatusOK {
		fmt.Println(uaaToken)
		//nil = uaaToken.ErrorDescription
		return "", uaaToken.ErrorDescription
	}

	t := time.Now().Local().Add(time.Second * time.Duration(uaaToken.Expire))
	uaaToken.ExpireTime = t
	defer resp.Body.Close()

	return result, ""
}

/**
	CF API를 통해 App 정보를 조회하는 함수
 */
func GetAppByGuid(n md.CFConfig, m md.UaaToken, guid string) (md.Resource, error) {
	var processResource md.ProcessResource
	var resource md.Resource
	apiUrl := n.ApiHost
	path := "/v3/apps/" + guid + "/processes"
	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = path
	urlStr := u.String() // "https://api.com/user/"

	tp := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	client := &http.Client{Transport: tp}
	r, _ := http.NewRequest("GET", urlStr, nil) // URL-encoded payload
	//r.Header.Add("Accept", "application/json")
	//r.Header.Add("Content-Type", "application/json")

	//fmt.Println(">>>>> [cf_util.go / GetAppByGuid()] CF API accessToken : " + m.Token)
	fmt.Println("====== cf_util.go / GetAppByGuid() - CF API URL : " + urlStr)
	r.Header.Add("Authorization", "bearer "+m.Token)
	//r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)

	if err != nil {
		return md.Resource{}, errors.Wrap(err, "Error requesting apps")
	}

	fmt.Println("====== cf_util.go / GetAppByGuid() - Request to ", urlStr, resp.Status)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return md.Resource{}, errors.Wrap(err, "Error reading app response body")
	}
	//bodyString := string(bodyBytes)
	//fmt.Println(">>>>> [cf_util.go / GetAppByGuid()]] resp.Body : ", bodyString)

	json.Unmarshal(bodyBytes, &processResource)
	//fmt.Println("====== cf_util.go / GetAppByGuid() - processResource : ", processResource)

	if resp.StatusCode != http.StatusOK {
		fmt.Println(processResource)

	}
	fmt.Println(processResource.Resources)
	for _, item := range processResource.Resources {
		//n.mergeAppResource(item)
		if item.Guid == guid {
			fmt.Printf("====== cf_util.go / GetAppByGuid() - Find Suitable App..!! GUID : %s, Instance count : %d\n", item.Guid, item.Instances)
			resource = item
			break
		}
		//return resource, nil
	}
	//json.Unmarshal(bodyBytes,  &processResource.Resources)
	//
	//if err != nil {
	//	return md.Process{}, errors.Wrap(err, "Error unmarshalling app")
	//}
	//
	//if resp.StatusCode != http.StatusOK {
	//	fmt.Println(processResource)
	//	//return uaaToken.ErrorDescription
	//}

	return resource, nil
}

//func mergeAppResource(app md.AppResource) md.App {
//	app.Entity.Guid = app.Meta.Guid
//	app.Entity.CreatedAt = app.Meta.CreatedAt
//	app.Entity.UpdatedAt = app.Meta.UpdatedAt
//	app.Entity.SpaceData.Entity.Guid = app.Entity.SpaceData.Meta.Guid
//	app.Entity.SpaceData.Entity.OrgData.Entity.Guid = app.Entity.SpaceData.Entity.OrgData.Meta.Guid
//	return app.Entity
//}

func UpdateApp(n md.CFConfig, m md.UaaToken, guid string, aur md.ScaleProcess) (md.Resource, error) {
	var Resource md.Resource

	apiUrl := n.ApiHost
	path := "/v3/processes/" + guid + "/actions/scale"
	//resource := fmt.Sprintf("/v3/apps/%s", guid)
	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = path
	urlStr := u.String() // "https://api.com/user/"

	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(aur)
	if err != nil {
		return md.Resource{}, err
	}

	tp := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	//fmt.Println(">>>>> [cf_util.go / UpdateApp()] CF API accessToken : " + m.Token)
	fmt.Println(">>>>> [cf_util.go / UpdateApp()] CF API URL : " + urlStr)

	client := &http.Client{Transport: tp}
	r, _ := http.NewRequest("POST", urlStr, buf) // URL-encoded payload
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "bearer "+m.Token)
	//r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)

	//req := c.NewRequestWithBody("PUT", fmt.Sprintf("/v2/apps/%s", guid), buf)
	//resp, err := c.DoRequest(req)
	if err != nil {
		return md.Resource{}, err
	}
	if resp.StatusCode != http.StatusAccepted {
		return md.Resource{}, fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}

	fmt.Println(">>>>> [cf_util.go / UpdateApp()] Request to ", urlStr, resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return md.Resource{}, err
	}

	bodyString := string(body)
	fmt.Println(bodyString)

	err = json.Unmarshal(body, &Resource)
	if err != nil {
		return md.Resource{}, err
	}
	return Resource, nil
}
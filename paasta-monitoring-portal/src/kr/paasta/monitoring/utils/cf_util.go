package utils

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"io/ioutil"
	cm "kr/paasta/monitoring/common/model"
	"kr/paasta/monitoring/iaas_new/model"
	pm "kr/paasta/monitoring/paas/model"
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
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func GetUaaToken(apiRequest cm.UserInfo, reqCsrfToken string, cfConfig pm.CFConfig, RdClient *redis.Client) (string, error) {

	//n.CfProvider.Username = strings.TrimSpace(apiRequest.PaasUserId)
	//n.CfProvider.Password = strings.TrimSpace(apiRequest.PaasUserPw)
	result := ""
	//HttpClient := http.DefaultClient
	if cfConfig.Type == "PAAS" {
		cfConfig.UserId = strings.TrimSpace(apiRequest.PaasUserId)
		cfConfig.UserPw = strings.TrimSpace(apiRequest.PaasUserPw)
	} else {
		cfConfig.UserId = strings.TrimSpace(apiRequest.CaasUserId)
		cfConfig.UserPw = strings.TrimSpace(apiRequest.CaasUserPw)
	}

	if len(cfConfig.UserId) > 0 && len(cfConfig.UserPw) > 0 {
		apiUrl := cfConfig.Host //+ ":" + cfConfig.Port
		//resource := "/uaa/oauth/token"
		resource := "/oauth/token"
		data := url.Values{}
		data.Set("grant_type", "password")
		data.Set("username", cfConfig.UserId)
		data.Set("password", cfConfig.UserPw)
		//data.Set("token_format", "opaque")
		data.Set("client_id", "cf")
		data.Set("response_type", "token")
		//data.Set("client_secret", cfConfig.ClientPw)
		u, _ := url.ParseRequestURI(apiUrl)
		u.Path = resource
		urlStr := u.String() // "https://api.com/user/"

		tp := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

		client := &http.Client{Transport: tp}
		r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
		r.Header.Add("Accept", "application/json")
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		//
		r.SetBasicAuth(url.QueryEscape("cf"), url.QueryEscape(""))

		resp, _ := client.Do(r)
		fmt.Println(resp.Status)

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)

		var uaaToken UaaToken
		json.Unmarshal(bodyBytes, &uaaToken)

		if resp.StatusCode != http.StatusOK {
			fmt.Println(uaaToken)
			return "", errors.New(uaaToken.ErrorDescription)
		}

		t := time.Now().Local().Add(time.Second * time.Duration(uaaToken.Expire))

		defer resp.Body.Close()

		if len(uaaToken.Scope) > 0 {
			if cfConfig.Type == "PAAS" {
				//apps, err1 := client.ListApps()
				//
				//fmt.Printf("cf-apps=[%v], cf-appsErr=[%v]\n", apps, err1)
				//
				//if err1 != nil {
				//	fmt.Println("paas ListApps err ::: ", err1.Error())
				//} else {
				//	//fmt.Println("paas ListApps client1 ::: ", client1)
				//}
				cfChk := false
				if strings.Contains(uaaToken.Scope, "cloud_controller.admin") {
					cfChk = true
				}

				if cfChk {
					//redis 에 사용자 정보를 수정 저장한다.
					//redis 에 사용자 정보를 저장한다.

					model.MonitLogger.Debug(RdClient)
					rdresult := RdClient.Get(reqCsrfToken)

					if rdresult.Val() == "" {
						model.MonitLogger.Debug("redis check")
						RdClient.HSet(reqCsrfToken, "paasUserId", cfConfig.UserId)
						RdClient.HSet(reqCsrfToken, "paasToken", uaaToken.Token)
						RdClient.HSet(reqCsrfToken, "paasAdminYn", 'Y')
						RdClient.HSet(reqCsrfToken, "paasExpire", t)
						RdClient.HSet(reqCsrfToken, "paasRefreshToken", uaaToken.Refresh)
					} else {
						fmt.Println(rdresult)
					}

					result = uaaToken.Token
				} else {
					return "", errors.New("not_admin_account")

				}
			} else {
				rdresult := RdClient.Get(reqCsrfToken)
				if rdresult.Val() == "" {
					RdClient.HSet(reqCsrfToken, "caasUserId", cfConfig.UserId)
					RdClient.HSet(reqCsrfToken, "caasToken", uaaToken.Token)
					RdClient.HSet(reqCsrfToken, "caasExpire", t)
					RdClient.HSet(reqCsrfToken, "caasRefreshToken", uaaToken.Refresh)
				}

				result = uaaToken.Token
			}

		} else {
			return "", errors.New("scope len is none")
		}

	} else {
		return "", errors.New("parameter is none")
		//result = "parameter is none"
	}

	return result, nil
}

func GetUaaReFreshToken(reqCsrfToken string, cfConfig pm.CFConfig, RdClient *redis.Client) (string, string) {

	result := ""
	//HttpClient := http.DefaultClient

	fmt.Println(RdClient.HGet(reqCsrfToken, "paasRefreshToken").Val())
	apiUrl := cfConfig.Host
	resource := "/oauth/token"
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("response_type", "token")
	data.Set("client_id", "cf")
	data.Set("refresh_token", RdClient.HGet(reqCsrfToken, "paasRefreshToken").Val())

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String() // "https://api.com/user/"

	tp := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	client := &http.Client{Transport: tp}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	r.SetBasicAuth(url.QueryEscape("cf"), url.QueryEscape(""))
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

	defer resp.Body.Close()

	if len(uaaToken.Scope) > 0 {
		if cfConfig.Type == "PAAS" {
			//apps, err1 := client.ListApps()
			//
			//fmt.Printf("cf-apps=[%v], cf-appsErr=[%v]\n", apps, err1)
			//
			//if err1 != nil {
			//	fmt.Println("paas ListApps err ::: ", err1.Error())
			//} else {
			//	//fmt.Println("paas ListApps client1 ::: ", client1)
			//}
			cfChk := false
			if strings.Contains(uaaToken.Scope, "cloud_controller.admin") {
				cfChk = true
			}

			if cfChk {
				//redis 에 사용자 정보를 수정 저장한다.
				//redis 에 사용자 정보를 저장한다.

				fmt.Println(RdClient)
				rdresult := RdClient.Get(reqCsrfToken)

				if rdresult.Val() == "" {
					RdClient.HSet(reqCsrfToken, "paasToken", uaaToken.Token)
					RdClient.HSet(reqCsrfToken, "paasExpire", t)
				}

				result = uaaToken.Token
			} else {
				//return "not_admin_account"
				//nil = "not_admin_account"
				return "", "not_admin_account"
			}
		} else {
			rdresult := RdClient.Get(reqCsrfToken)
			if rdresult.Val() == "" {
				RdClient.HSet(reqCsrfToken, "caasToken", uaaToken.Token)
				RdClient.HSet(reqCsrfToken, "caasExpire", t)
			}

			result = uaaToken.Token
		}

	} else {
		result = ""
	}

	return result, ""
}

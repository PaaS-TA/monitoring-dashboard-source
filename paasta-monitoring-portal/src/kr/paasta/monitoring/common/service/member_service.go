package services

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gophercloud/gophercloud"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"monitoring-portal/common/dao"
	cm "monitoring-portal/common/model"
	"monitoring-portal/iaas_new/model"
	pm "monitoring-portal/paas/model"
	"monitoring-portal/utils"
	"net/http"
	"net/url"
	"strings"
)

type MemberService struct {
	openstackProvider model.OpenstackProvider
	//CfProvider        cfclient.Config
	txn      *gorm.DB
	RdClient *redis.Client
}

type UaaToken struct {
	Token            string `json:"access_token"`
	Scope            string `json:"scope"`
	Expire           int64  `json:"expires_in"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
type CaasBroker struct {
	ResultCode string `json:"resultCode"`
	Items      []Item
}
type Item struct {
	UserId      string `json:"userId"`
	AccountName string `json:"caasAccountName"`
}

func GetMemberService(openstackProvider model.OpenstackProvider, txn *gorm.DB, rdClient *redis.Client) *MemberService {
	return &MemberService{
		openstackProvider: openstackProvider,
		//CfProvider : cfProvider,
		txn:      txn,
		RdClient: rdClient,
	}
}

func GetIaasMemberService(openstackProvider model.OpenstackProvider, txn *gorm.DB, rdClient *redis.Client) *MemberService {
	return &MemberService{
		openstackProvider: openstackProvider,
		txn:               txn,
		RdClient:          rdClient,
	}
}

func GetPaasMemberService(txn *gorm.DB, rdClient *redis.Client) *MemberService {
	return &MemberService{
		//CfProvider : cfProvider,
		txn:      txn,
		RdClient: rdClient,
	}
}

func (n MemberService) MemberJoinInfo() (val string) {

	return ""
}

func (n MemberService) MemberJoinSave(req cm.UserInfo) error {

	dbErr := dao.GetMemberDao(n.txn).MemberJoinSave(req, n.txn)

	if dbErr != nil {
		return dbErr
	}

	return nil
}

/*func (n MemberService) MemberAuthCheck(req cm.UserInfo) (result cm.UserInfo,  err error) {


	//member info 로그인 아이디 패스워드를 회원디비에서 검색하여 iaas, paas 로그인 정보를 가지고 온다.
	result, _ , dbErr := dao.GetMemberDao(n.txn).MemberInfoView(req, n.txn)

	if dbErr != nil {
	   return result, dbErr
	}

	//get iaas token
	if result.IaasUserId != "" && result.IaasUserPw != "" {

		n.openstackProvider.Username = result.IaasUserId
		n.openstackProvider.Password = result.IaasUserPw
		provider, err := utils.GetOpenstackToken(n.openstackProvider)
		if err != nil {
			fmt.Println("MemberAuthCheck iaas error::", err.Error())
			//return req, provider, err
		} else {
			result.IaasToken = provider.TokenID
			result.IaasToken = "Y"
			//fmt.Println("MemberAuthCheck iaas token ::: ", result.IaasToken)
		}
	}

	if result.PaasUserId != "" && result.PaasUserPw != "" {
		//get paas token
		n.CfProvider.Username = result.PaasUserId
		n.CfProvider.Password = result.PaasUserPw
		client, err := cfclient.NewClient(&n.CfProvider)
		if err != nil {
			fmt.Println("paas error::", err.Error())
			//return req, provider, err
		} else {
			token, _ := client.GetToken()
			//result.PaasToken = strings.Replace(token, "bearer ", "", -1)
			result.PaasToken = token
			//fmt.Println("MemberAuthCheck paas token ::: ", result.PaasToken)
			//fmt.Println("MemberAuthCheck paas Scope ::: ", client.Config.Scope)

			if len(client.Config.Scope) > 0{
				cfChk := false
				for _, value := range client.Config.Scope{
					if strings.Contains(value,"cloud_controller.admin") {
						cfChk = true
					}
				}

				if cfChk {
					result.PaasAdminYn = "Y"
					result.PaasToken = "Y"
					fmt.Println("MemberAuthCheck paas admin ok  !!! ")
				} else {
					result.PaasAdminYn = "N"
					result.PaasToken = "Y"
					fmt.Println("MemberAuthCheck paas admin fail  !!! ")
				}
			} else {
				result.PaasAdminYn = "N"
				result.PaasToken = "Y"
				fmt.Println("MemberAuthCheck paas admin fail  !!! ")
			}
		}
	}

	//fmt.Println("MemberAuthCheck preq.Token==>",req.Token)
	//redis 에 사용자 정보가 있으면 변경하여 저장한다.
	if req.Token != "" && result.IaasToken != "" {
		n.RdClient.HSet(req.Token,"iaasToken",result.IaasToken)
		n.RdClient.Expire(req.Token, 30 * 60  * time.Second )
	}

	if req.Token != "" && result.IaasToken != "" {
		n.RdClient.HSet(req.Token,"paasToken",result.PaasToken)
		n.RdClient.Expire(req.Token, 30 * 60  * time.Second )
	}



	return result,  err
}*/

func (n MemberService) MemberInfoCheck(req cm.UserInfo) (userInfo cm.UserInfo, provider *gophercloud.ProviderClient, err error) {

	result, _, dbErr := dao.GetMemberDao(n.txn).MemberInfoCheck(req, n.txn)

	if dbErr != nil {
		return result, provider, dbErr
	}

	return result, provider, dbErr
}

func (n MemberService) MemberInfoView(req cm.UserInfo) (userInfo cm.UserInfo, provider *gophercloud.ProviderClient, err error) {
	fmt.Println("MemberInfoView service req.userid", req.UserId)
	result, _, dbErr := dao.GetMemberDao(n.txn).MemberInfoView(req, n.txn)

	if dbErr != nil {
		return result, provider, dbErr
	}

	return result, provider, dbErr
}

func (n MemberService) MemberInfoUpdate(req cm.UserInfo) (userInfo cm.UserInfo, provider *gophercloud.ProviderClient, err error) {

	var nullInfo cm.UserInfo

	_, dbErr := dao.GetMemberDao(n.txn).MemberInfoUpdate(req, n.txn)

	if dbErr != nil {
		return nullInfo, provider, dbErr
	}

	result, _, dbErr := dao.GetMemberDao(n.txn).MemberInfoView(req, n.txn)

	if dbErr != nil {
		return nullInfo, provider, dbErr
	}

	return result, provider, dbErr
}

func (n MemberService) MemberInfoDelete(req cm.UserInfo) (cnt int, err error) {

	result, dbErr := dao.GetMemberDao(n.txn).MemberInfoDelete(req, n.txn)

	if dbErr != nil {
		return result, dbErr
	}

	return result, dbErr
}

func (s MemberService) GetIaasToken(apiRequest cm.UserInfo, reqCsrfToken string) string {
	//get iaas token
	s.openstackProvider.Username = apiRequest.IaasUserId
	s.openstackProvider.Password = apiRequest.IaasUserPw

	provider, err := utils.GetOpenstackToken(s.openstackProvider)
	result := ""

	if err != nil {
		utils.Logger.Error(err.Error())

		if strings.Contains(err.Error(), "Unauthorized") {
			return "unauthorized"
		} else {
			fmt.Println("unexpected_fail: return empty string")
			return ""
		}
		//utils.RenderJsonLogoutResponse("", w)
	} else {
		userAuthInfo := s.RdClient.Get(reqCsrfToken)

		if userAuthInfo.Val() == "" {
			s.RdClient.HSet(reqCsrfToken, "iaasUserId", apiRequest.IaasUserId)
			s.RdClient.HSet(reqCsrfToken, "iaasToken", provider.TokenID)
		}

		//utils.RenderJsonLogoutResponse(provider.TokenID, w)
		//result = provider.TokenID
	}
	return result
}

func shallowDefaultTransport() *http.Transport {
	defaultTransport := http.DefaultTransport.(*http.Transport)
	return &http.Transport{
		Proxy:                 defaultTransport.Proxy,
		TLSHandshakeTimeout:   defaultTransport.TLSHandshakeTimeout,
		ExpectContinueTimeout: defaultTransport.ExpectContinueTimeout,
	}
}

func (n MemberService) CaasServiceCheck(apiRequest cm.UserInfo, reqCsrfToken string, cfConfig pm.CFConfig) string {

	result := ""

	apiUrl := cfConfig.CaasBrokerHost
	resource := "/users"
	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String() // "https://api.com/user/"

	client := &http.Client{}
	r, _ := http.NewRequest("GET", urlStr, nil) // URL-encoded payload
	r.SetBasicAuth("admin", "PaaS-TA")
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")
	resp, _ := client.Do(r)
	fmt.Println(resp.Status)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)

	var caasBroker CaasBroker
	json.Unmarshal(bodyBytes, &caasBroker.Items)

	if resp.StatusCode != http.StatusOK {
		fmt.Println(caasBroker)
		return caasBroker.ResultCode
	}
	fmt.Println(caasBroker.Items)
	for _, item := range caasBroker.Items {
		//n.mergeAppResource(item)
		if item.UserId == apiRequest.CaasUserId {

			if strings.Contains(item.AccountName, "-admin") {
				result = "adm"
			}
		}
	}

	return result
}



func (n MemberService) MemberJoinCheckDuplicationIaasId(req cm.UserInfo) (userInfo cm.UserInfo, err error) {
	return dao.GetMemberDao(n.txn).MemberJoinCheckDuplicationIaasId(req, n.txn)
}

func (n MemberService) MemberJoinCheckDuplicationPaasId(req cm.UserInfo) (userInfo cm.UserInfo, err error) {
	return dao.GetMemberDao(n.txn).MemberJoinCheckDuplicationPaasId(req, n.txn)
}


func (n MemberService) MemberJoinCheckDuplicationCaasId(req cm.UserInfo) (userInfo cm.UserInfo, err error) {
	return dao.GetMemberDao(n.txn).MemberJoinCheckDuplicationCaasId(req, n.txn)
}

func (n MemberService) DeleteIaasToken(reqCsrfToken string) {

	beforeToken := ""

	val := n.RdClient.HGetAll(reqCsrfToken).Val()
	for key, value := range val {
		if key == "iaasToken" {
			beforeToken = value
		}
	}

	if beforeToken != "" {
		n.RdClient.HDel(reqCsrfToken, "iaasToken")
		n.RdClient.HDel(reqCsrfToken, "iaasUserId")
	}
}

func (n MemberService) DeletePaasToken(reqCsrfToken string) {

	beforeToken := ""

	val := n.RdClient.HGetAll(reqCsrfToken).Val()
	for key, value := range val {
		if key == "paasToken" {
			beforeToken = value
		}
	}

	if beforeToken != "" {
		n.RdClient.HDel(reqCsrfToken, "paasToken")
		n.RdClient.HDel(reqCsrfToken, "paasAdminYn")
		n.RdClient.HDel(reqCsrfToken, "paasUserId")
	}
}

/*func (n MemberS/*ervice) DeleteSaasToken(reqCsrfToken string) {

	beforeToken := ""

	val := n.RdClient.HGetAll(reqCsrfToken).Val()
	for key, value := range val{
		if key == "saasToken" {
			beforeToken = value
		}
	}

	if beforeToken != "" {
		n.RdClient.HDel(reqCsrfToken,"saasToken" )
		n.RdClient.HDel(reqCsrfToken,"saasUserId" )
	}
}*/
func (n MemberService) DeleteCaasToken(reqCsrfToken string) {

	beforeToken := ""

	val := n.RdClient.HGetAll(reqCsrfToken).Val()
	for key, value := range val {
		if key == "caasToken" {
			beforeToken = value
		}
	}

	if beforeToken != "" {
		n.RdClient.HDel(reqCsrfToken, "caasToken")
		n.RdClient.HDel(reqCsrfToken, "caasAdminYn")
		n.RdClient.HDel(reqCsrfToken, "caasUserId")
	}

}

///////////////////////test////////////////////////////////////////////
//func (n MemberService) GetAppByGuid(cfConfig pm.CFConfig,guid string,hkey string) (cm.Resource, error) {
//	var processResource cm.ProcessResource
//	var resource cm.Resource
//	apiUrl := "https://api.15.164.20.58.xip.io"
//	path := "/v3/apps/"+guid+"/processes"
//	u, _ := url.ParseRequestURI(apiUrl)
//	u.Path = path
//	urlStr := u.String() // "https://api.com/user/"
//
//
//
//	tp := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true},}
//
//	val := n.RdClient.HGetAll(hkey).Val()
//	client := &http.Client{Transport: tp}
//	r, _ := http.NewRequest("GET", urlStr,nil) // URL-encoded payload
//	//r.Header.Add("Accept", "application/json")
//	//r.Header.Add("Content-Type", "application/json")
//	r.Header.Add("Authorization", "bearer "+val["paasToken"])
//	//r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
//	fmt.Println(n.RdClient.HMGet(hkey,"paasToken").Val())
//
//	resp, err := client.Do(r)
//
//	if err != nil {
//		return cm.Resource{}, errors.Wrap(err, "Error requesting apps")
//	}
//
//
//	fmt.Println(resp.Status)
//
//	bodyBytes, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		fmt.Println(err)
//		return cm.Resource{}, errors.Wrap(err, "Error reading app response body")
//	}
//	bodyString := string(bodyBytes)
//	fmt.Println(bodyString)
//
//
//
//	json.Unmarshal(bodyBytes, &processResource)
//	fmt.Println(processResource)
//	fmt.Println("ccccccc")
//	fmt.Println(processResource.Resources)
//	fmt.Println("mmmmmmm")
//	if resp.StatusCode != http.StatusOK {
//		fmt.Println(processResource)
//
//	}
//	//fmt.Println(processResource.Resources)
//	for _, resource  = range processResource.Resources {
//		//n.mergeAppResource(item)
//		fmt.Println("pppppppppp")
//		fmt.Println(resource)
//		//return resource, nil
//	}
//	fmt.Println(resource)
//	//json.Unmarshal(bodyBytes,  &processResource.Resources)
//	//
//	//if err != nil {
//	//	return md.Process{}, errors.Wrap(err, "Error unmarshalling app")
//	//}
//	//
//	//if resp.StatusCode != http.StatusOK {
//	//	fmt.Println(processResource)
//	//	//return uaaToken.ErrorDescription
//	//}
//
//	return resource, nil
//}
//
////func mergeAppResource(app md.AppResource) md.App {
////	app.Entity.Guid = app.Meta.Guid
////	app.Entity.CreatedAt = app.Meta.CreatedAt
////	app.Entity.UpdatedAt = app.Meta.UpdatedAt
////	app.Entity.SpaceData.Entity.Guid = app.Entity.SpaceData.Meta.Guid
////	app.Entity.SpaceData.Entity.OrgData.Entity.Guid = app.Entity.SpaceData.Entity.OrgData.Meta.Guid
////	return app.Entity
////}
//
//func (n MemberService) UpdateApp(cfConfig pm.CFConfig,guid string, aur cm.ScaleProcess,hkey string) (cm.Resource, error) {
//	var Resource cm.Resource
//
//	apiUrl := "https://api.15.164.20.58.xip.io"
//	path := "/v3/processes/"+guid+"/actions/scale"
//	//resource := fmt.Sprintf("/v3/apps/%s", guid)
//	u, _ := url.ParseRequestURI(apiUrl)
//	u.Path = path
//	urlStr := u.String() // "https://api.com/user/"
//
//	buf := bytes.NewBuffer(nil)
//	err := json.NewEncoder(buf).Encode(aur)
//	if err != nil {
//		return cm.Resource{}, err
//	}
//
//	tp := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true},}
//	val := n.RdClient.HGetAll(hkey).Val()
//
//	client := &http.Client{Transport: tp}
//	r, _ := http.NewRequest("POST", urlStr,buf) // URL-encoded payload
//	r.Header.Add("Accept", "application/json")
//	r.Header.Add("Content-Type", "application/json")
//	r.Header.Add("Authorization", "bearer "+val["paasToken"])
//	//r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
//
//	resp, err := client.Do(r)
//
//	//req := c.NewRequestWithBody("PUT", fmt.Sprintf("/v2/apps/%s", guid), buf)
//	//resp, err := c.DoRequest(req)
//	if err != nil {
//		return cm.Resource{}, err
//	}
//	if resp.StatusCode != http.StatusCreated {
//		return cm.Resource{}, fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
//	}
//
//	body, err := ioutil.ReadAll(resp.Body)
//	defer resp.Body.Close()
//	if err != nil {
//		return cm.Resource{}, err
//	}
//	err = json.Unmarshal(body, &Resource)
//	if err != nil {
//		return cm.Resource{}, err
//	}
//	return Resource, nil
//}

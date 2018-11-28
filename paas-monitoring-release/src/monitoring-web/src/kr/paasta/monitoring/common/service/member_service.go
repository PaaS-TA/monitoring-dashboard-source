package services

import (
	"github.com/rackspace/gophercloud"
	"kr/paasta/monitoring/iaas/model"
	cm "kr/paasta/monitoring/common/model"
	"kr/paasta/monitoring/common/dao"
	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/jinzhu/gorm"
	"github.com/go-redis/redis"
	"kr/paasta/monitoring/utils"
	"fmt"
	"strings"
	"time"
)

type MemberService struct {
	openstackProvider model.OpenstackProvider
	CfProvider        cfclient.Config
	txn *gorm.DB
	RdClient *redis.Client
}

func GetMemberService(openstackProvider model.OpenstackProvider, cfProvider  cfclient.Config, txn *gorm.DB,  rdClient *redis.Client ) *MemberService {
	return &MemberService{
		openstackProvider: openstackProvider,
		CfProvider : cfProvider,
		txn: txn,
		RdClient: rdClient,
	}
}

func GetIaasMemberService(openstackProvider model.OpenstackProvider, txn *gorm.DB,  rdClient *redis.Client ) *MemberService {
	return &MemberService{
		openstackProvider: openstackProvider,
		txn: txn,
		RdClient: rdClient,
	}
}

func GetPaasMemberService(cfProvider  cfclient.Config, txn *gorm.DB,  rdClient *redis.Client ) *MemberService {
	return &MemberService{
		CfProvider : cfProvider,
		txn: txn,
		RdClient: rdClient,
	}
}

func (n MemberService) MemberJoinInfo() (val string) {



	return ""
}

func (n MemberService) MemberJoinSave(req cm.UserInfo) (error) {

	 dbErr := dao.GetMemberDao(n.txn).MemberJoinSave(req, n.txn)

	if dbErr != nil {
		return dbErr
	}

	return nil
}

func (n MemberService) MemberAuthCheck(req cm.UserInfo) (result cm.UserInfo,  err error) {


	//member info 로그인 아이디 패스워드를 회원디비에서 검색하여 iaas, paas 로그인 정보를 가지고 온다.
	result, _ , dbErr := dao.GetMemberDao(n.txn).MemberInfoView(req, n.txn)

	if dbErr != nil {
	   return result, dbErr
	}

	//get iaas token
	if result.IaasUserId != "" && result.IaasUserPw != "" {

		n.openstackProvider.Username = result.IaasUserId
		n.openstackProvider.Password = result.IaasUserPw
		provider, err := utils.GetAdminToken(n.openstackProvider)
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
}


func (n MemberService) MemberInfoCheck(req cm.UserInfo) (userInfo cm.UserInfo, provider *gophercloud.ProviderClient, err error) {

	result, _ , dbErr := dao.GetMemberDao(n.txn).MemberInfoCheck(req, n.txn)

	if dbErr != nil {
		return result, provider, dbErr
	}

	return result, provider, dbErr
}


func (n MemberService) MemberInfoView(req cm.UserInfo) (userInfo cm.UserInfo, provider *gophercloud.ProviderClient, err error) {
	fmt.Println("MemberInfoView service req.userid",req.UserId)
	result, _ , dbErr := dao.GetMemberDao(n.txn).MemberInfoView(req, n.txn)

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

	result, _ , dbErr := dao.GetMemberDao(n.txn).MemberInfoView(req, n.txn)

	if dbErr != nil {
		return nullInfo, provider, dbErr
	}

	return result, provider, dbErr
}

func (n MemberService) MemberInfoDelete(req cm.UserInfo) (cnt  int,  err error) {

	result , dbErr := dao.GetMemberDao(n.txn).MemberInfoDelete(req, n.txn)

	if dbErr != nil {
		return result,  dbErr
	}

	return result,  dbErr
}


func (s MemberService) GetIaasToken( apiRequest cm.UserInfo, reqCsrfToken string) string {
	//get iaas token
	s.openstackProvider.Username = apiRequest.IaasUserId
	s.openstackProvider.Password = apiRequest.IaasUserPw

	provider, err := utils.GetAdminToken(s.openstackProvider)
	result := ""

	if err != nil {
		fmt.Println("IaaS(OpenStack) login result:", err.Error())

		if strings.Contains(err.Error(), "Unauthorized") {
			return "unauthorized"
		} else {
			fmt.Println("unexpected_fail: return empty string")
			return ""
		}
		//utils.RenderJsonLogoutResponse("", w)
	} else {
		rdresult := s.RdClient.Get(reqCsrfToken)
		//fmt.Println("rdresult.Val() :::", rdresult.Val())
		//fmt.Println("RdClient.Get :::", reqCsrfToken, ", provider.TokenID :::", provider.TokenID)

		if rdresult.Val() == "" {
			s.RdClient.HSet(reqCsrfToken, "iaasUserId", apiRequest.IaasUserId)
			s.RdClient.HSet(reqCsrfToken, "iaasToken", provider.TokenID)
		}

		//utils.RenderJsonLogoutResponse(provider.TokenID, w)
		result = provider.TokenID
	}
	return result
}


func (n MemberService) GetPaasToken(apiRequest cm.UserInfo,  reqCsrfToken string) string {

	n.CfProvider.Username = strings.TrimSpace(apiRequest.PaasUserId)
	n.CfProvider.Password = strings.TrimSpace(apiRequest.PaasUserPw)
	result := ""
	//fmt.Println("paas Username :::>",len(n.CfProvider.Username),"</>",len(n.CfProvider.Password))
	//fmt.Println("paas CfProvider ::: ",	n.CfProvider)
	if len(n.CfProvider.Username) > 0 && len(n.CfProvider.Password) > 0 {

		client, err := cfclient.NewClient(&n.CfProvider)

		if err != nil {
			fmt.Println("PaaS(cf) login result:", err.Error())

			if strings.Contains(err.Error(), "Bad credentials") {
				fmt.Println("bad_credentials")
				return "bad_credentials"
			} else if strings.Contains(err.Error(), "account has been locked") {
				fmt.Println("account_locked")
				return "account_locked"
			} else {
				fmt.Println("unexpected_fail: return empty string")
				return ""
			}
		}

		if err != nil || len(n.CfProvider.Username) == 0 || len(n.CfProvider.Password) == 0 {
			fmt.Println("paas error::", err.Error())
			//utils.RenderJsonLogoutResponse("", w)
			//utils.ErrRenderJsonResponse(err, w)
			result = ""
		} else {
			token, tokenErr := client.GetToken()

			fmt.Printf("cf-token=[%v], cf-tokenErr=[%v]\n", token, tokenErr)

			apps, err1 := client.ListApps()

			fmt.Printf("cf-apps=[%v], cf-appsErr=[%v]\n", apps, err1)

			if err1 != nil {
				fmt.Println("paas ListApps err ::: ", err1.Error())
			} else {
				//fmt.Println("paas ListApps client1 ::: ", client1)
			}
			apiRequest.PaasToken = strings.Replace(token, "bearer ", "", -1)
			//result.PaasToken = token
			//fmt.Println("paas token ::: ", apiRequest.PaasToken)
			//fmt.Println("paas Scope ::: ",client.Config.Scope)

			fmt.Printf("client.Config.Scope=[%v]\n", client.Config.Scope)

			if len(client.Config.Scope) > 0 {
				cfChk := false
				for _, value := range client.Config.Scope {
					if strings.Contains(value,"cloud_controller.admin") {
						cfChk = true
					}
				}

				if cfChk {
					//redis 에 사용자 정보를 수정 저장한다.
					//redis 에 사용자 정보를 저장한다.
					rdresult := n.RdClient.Get(reqCsrfToken)
					if rdresult.Val() == "" {
						n.RdClient.HSet(reqCsrfToken, "paasUserId", n.CfProvider.Username)
						n.RdClient.HSet(reqCsrfToken, "paasToken", apiRequest.PaasToken)
						n.RdClient.HSet(reqCsrfToken, "paasAdminYn", 'Y')
					}

					result = apiRequest.PaasToken
				} else {
					return "not_admin_account"
				}
			} else {
				result = ""
			}
		}
	} else {
		result = ""
	}

	return result
}

func (n MemberService) MemberJoinCheckDuplicationIaasId(req cm.UserInfo) (userInfo cm.UserInfo, err error) {
	return dao.GetMemberDao(n.txn).MemberJoinCheckDuplicationIaasId(req, n.txn)
}

func (n MemberService) MemberJoinCheckDuplicationPaasId(req cm.UserInfo) (userInfo cm.UserInfo, err error) {
	return dao.GetMemberDao(n.txn).MemberJoinCheckDuplicationPaasId(req, n.txn)
}

func (n MemberService) DeleteIaasToken(reqCsrfToken string) {

	beforeToken := ""

	val := n.RdClient.HGetAll(reqCsrfToken).Val()
	for key, value := range val{
		if key == "iaasToken" {
			beforeToken = value
		}
	}

	if beforeToken != "" {
		n.RdClient.HDel(reqCsrfToken,"iaasToken" )
		n.RdClient.HDel(reqCsrfToken,"iaasUserId" )
	}
}

func (n MemberService) DeletePaasToken(reqCsrfToken string) {

	beforeToken := ""

	val := n.RdClient.HGetAll(reqCsrfToken).Val()
	for key, value := range val{
		if key == "paasToken" {
			beforeToken = value
		}
	}

	if beforeToken != "" {
		n.RdClient.HDel(reqCsrfToken,"paasToken" )
		n.RdClient.HDel(reqCsrfToken,"paasAdminYn" )
		n.RdClient.HDel(reqCsrfToken,"paasUserId" )
	}
}
package services

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/identity/v3/tokens"
	"kr/paasta/monitoring/iaas/integration"
	"kr/paasta/monitoring/iaas/model"
	cm "kr/paasta/monitoring/common/model"
	"kr/paasta/monitoring/common/dao"
	"kr/paasta/monitoring/utils"
	"github.com/cloudfoundry-community/go-cfclient"
	"fmt"
	"github.com/jinzhu/gorm"
	//"strings"
	"strings"
	"github.com/go-redis/redis"
	"time"
)

type LoginService struct {
	openstackProvider model.OpenstackProvider
	CfProvider        cfclient.Config
	txn *gorm.DB
	RdClient *redis.Client
	sysType string
}

func GetLoginService(openstackProvider model.OpenstackProvider, cfProvider  cfclient.Config, txn *gorm.DB,  rdClient *redis.Client, sysType string ) *LoginService {
	return &LoginService{
		openstackProvider: openstackProvider,
		CfProvider : cfProvider,
		txn: txn,
		RdClient: rdClient,
		sysType: sysType,
	}
}

func GetIaasLoginService(openstackProvider model.OpenstackProvider, txn *gorm.DB,  rdClient *redis.Client, sysType string ) *LoginService {
	return &LoginService{
		openstackProvider: openstackProvider,
		txn: txn,
		RdClient: rdClient,
		sysType: sysType,
	}
}

func GetPaasLoginService(cfProvider  cfclient.Config, txn *gorm.DB,  rdClient *redis.Client, sysType string ) *LoginService {
	return &LoginService{
		CfProvider : cfProvider,
		txn: txn,
		RdClient: rdClient,
		sysType: sysType,
	}
}

func (n *LoginService) Logout(provider *gophercloud.ProviderClient, reqCsrfToken string) tokens.RevokeResult {
	result := integration.GetKeystone(n.openstackProvider, provider).RevokeToken()

	return result
}

func (n *LoginService) Login(req cm.UserInfo) (userInfo cm.UserInfo, provider *gophercloud.ProviderClient, err error) {

	//member info 로그인 아이디 패스워드를 회원디비에서 검색하여 iaas, paas 로그인 정보를 가지고 온다.
	result, _ , dbErr := dao.GetLoginDao(n.txn).GetLoginMemberInfo(req, n.txn)

	if dbErr != nil {
		if strings.Contains(dbErr.Error(), "record not found"){
			return result, provider, fmt.Errorf("Member certification information is invalid.")
		}else{
			return result, provider, dbErr
		}
	}

	if n.sysType == utils.SYS_TYPE_IAAS {
		if result.IaasUserUseYn == "Y"{
			//get iaas token
			n.openstackProvider.Username = result.IaasUserId
			n.openstackProvider.Password = result.IaasUserPw
			provider, err = utils.GetAdminToken(n.openstackProvider)
			if err != nil {
				fmt.Println("iaas error::",err.Error())
				//return req, provider, err
			}else{
				result.IaasToken = provider.TokenID
				//fmt.Println("iaas token ::: ", result.IaasToken )
			}

			n.RdClient.HSet(req.Token,"iaasToken",result.IaasToken)
			n.RdClient.HSet(req.Token,"iaasUserId",result.IaasUserId)
		}

		result.PaasAdminYn = "N"
	}else if n.sysType == utils.SYS_TYPE_PAAS {
		if result.PaasUserUseYn == "Y"{
			//get paas token
			n.CfProvider.Username = result.PaasUserId
			n.CfProvider.Password = result.PaasUserPw
			client, err := cfclient.NewClient(&n.CfProvider)
			if err != nil {
				fmt.Println("paas error::",err.Error())
				//return req, provider, err
			}else {
				token, _ := client.GetToken()
				result.PaasToken = strings.Replace(token, "bearer ", "", -1)
				//result.PaasToken = token
				//fmt.Println("paas token ::: ",result.PaasToken)
				//fmt.Println("paas Scope ::: ",client.Config.Scope)

				if len(client.Config.Scope) > 0{
					cfChk := false
					for _, value := range client.Config.Scope{
						if(strings.Contains(value,"cloud_controller.admin")) {
							cfChk = true
						}
					}

					if cfChk {
						result.PaasAdminYn = "Y"
					}else{
						result.PaasAdminYn = "N"
					}
				}else{
					result.PaasAdminYn = "N"
				}
			}

			n.RdClient.HSet(req.Token,"paasToken",result.PaasToken)
			n.RdClient.HSet(req.Token,"paasUserId",result.PaasUserId)
		}
	}else{
		if result.IaasUserUseYn == "Y"{
			//get iaas token
			n.openstackProvider.Username = result.IaasUserId
			n.openstackProvider.Password = result.IaasUserPw
			provider, err = utils.GetAdminToken(n.openstackProvider)
			if err != nil {
				fmt.Println("iaas error::",err.Error())
				//return req, provider, err
			}else{
				result.IaasToken = provider.TokenID
				//fmt.Println("iaas token ::: ", result.IaasToken )
			}

			n.RdClient.HSet(req.Token,"iaasToken",result.IaasToken)
			n.RdClient.HSet(req.Token,"iaasUserId",result.IaasUserId)
		}

		if result.PaasUserUseYn == "Y"{
			//get paas token
			n.CfProvider.Username = result.PaasUserId
			n.CfProvider.Password = result.PaasUserPw
			client, err := cfclient.NewClient(&n.CfProvider)
			if err != nil {
				fmt.Println("paas error::",err.Error())
				//return req, provider, err
			}else {
				token, _ := client.GetToken()
				result.PaasToken = strings.Replace(token, "bearer ", "", -1)
				//result.PaasToken = token
				//fmt.Println("paas token ::: ",result.PaasToken)
				//fmt.Println("paas Scope ::: ",client.Config.Scope)

				if len(client.Config.Scope) > 0{
					cfChk := false
					for _, value := range client.Config.Scope{
						if(strings.Contains(value,"cloud_controller.admin")) {
							cfChk = true
						}
					}

					if cfChk {
						result.PaasAdminYn = "Y"
					}else{
						result.PaasAdminYn = "N"
					}
				}else{
					result.PaasAdminYn = "N"
				}
			}

			n.RdClient.HSet(req.Token,"paasToken",result.PaasToken)
			n.RdClient.HSet(req.Token,"paasUserId",result.PaasUserId)
		}
	}

	//redis 에 사용자 정보를 저장한다.
	n.RdClient.HSet(req.Token,"userId",result.UserId)
	n.RdClient.HSet(req.Token,"userPw",result.UserPw)
	n.RdClient.HSet(req.Token,"paasAdminYn",result.PaasAdminYn)
	n.RdClient.Expire(req.Token, 30 * 60  * time.Second )

	//if err != nil {
	//	panic(err)
	//}

	return result, provider, dbErr
}

func (s *LoginService) SetUserInfoCache(userInfo *cm.UserInfo,  reqCsrfToken string) {

	var strcd1 = ""
	var strcd2 = ""
	var strcd3 = ""
	var strcd4 = ""

	userInfo.Username = userInfo.UserId

	if s.sysType == utils.SYS_TYPE_IAAS {
		if (userInfo.IaasUserUseYn == "Y") {
			strcd1 = "I"
			var tokenRequest cm.UserInfo
			tokenRequest.IaasUserId = userInfo.IaasUserId
			tokenRequest.IaasUserPw = userInfo.IaasUserPw
			result := GetIaasMemberService(s.openstackProvider, s.txn, s.RdClient).GetIaasToken(tokenRequest, reqCsrfToken)

			if (result != "") {
				strcd2 = "S"
			} else {
				strcd2 = "F"
			}

		} else {
			strcd1 = "F"
			strcd2 = "F"
			GetIaasMemberService(s.openstackProvider, s.txn, s.RdClient).DeleteIaasToken(reqCsrfToken)
		}
	}else if s.sysType == utils.SYS_TYPE_PAAS {
		if (userInfo.PaasUserUseYn == "Y") {
			strcd3 = "P"
			var tokenRequest cm.UserInfo
			tokenRequest.PaasUserId = userInfo.PaasUserId
			tokenRequest.PaasUserPw = userInfo.PaasUserPw
			result := GetPaasMemberService(s.CfProvider, s.txn, s.RdClient).GetPaasToken(tokenRequest, reqCsrfToken)
			if (result != "") {
				strcd4 = "S"
			} else {
				strcd4 = "F"
			}

		} else {
			strcd3 = "F"
			strcd4 = "F"
			GetPaasMemberService(s.CfProvider, s.txn, s.RdClient).DeletePaasToken(reqCsrfToken)
		}
	}else{
		if (userInfo.IaasUserUseYn == "Y") {
			strcd1 = "I"
			var tokenRequest cm.UserInfo
			tokenRequest.IaasUserId = userInfo.IaasUserId
			tokenRequest.IaasUserPw = userInfo.IaasUserPw
			result := GetMemberService(s.openstackProvider, s.CfProvider, s.txn, s.RdClient).GetIaasToken(tokenRequest, reqCsrfToken)

			if (result != "") {
				strcd2 = "S"
			} else {
				strcd2 = "F"
			}

		} else {
			strcd1 = "F"
			strcd2 = "F"
			GetMemberService(s.openstackProvider, s.CfProvider, s.txn, s.RdClient).DeleteIaasToken(reqCsrfToken)
		}

		if (userInfo.PaasUserUseYn == "Y") {
			strcd3 = "P"
			var tokenRequest cm.UserInfo
			tokenRequest.PaasUserId = userInfo.PaasUserId
			tokenRequest.PaasUserPw = userInfo.PaasUserPw
			result := GetMemberService(s.openstackProvider, s.CfProvider, s.txn, s.RdClient).GetPaasToken(tokenRequest, reqCsrfToken)
			if (result != "") {
				strcd4 = "S"
			} else {
				strcd4 = "F"
			}

		} else {
			strcd3 = "F"
			strcd4 = "F"
			GetMemberService(s.openstackProvider, s.CfProvider, s.txn, s.RdClient).DeletePaasToken(reqCsrfToken)
		}
	}

	//fmt.Println("login User Auth ::: ", strcd1, strcd2, strcd3, strcd4)

	userInfo.AuthI1 = strcd1
	userInfo.AuthI2 = strcd2
	userInfo.AuthP1 = strcd3
	userInfo.AuthP2 = strcd4

	userInfo.IaasUserPw = ""
	userInfo.PaasUserPw = ""
	userInfo.UserPw = ""
}
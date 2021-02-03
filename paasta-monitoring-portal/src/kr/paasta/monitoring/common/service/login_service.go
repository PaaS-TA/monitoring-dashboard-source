package services

import (
	"fmt"
	//"github.com/cloudfoundry-community/go-cfclient"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/identity/v3/tokens"
	"kr/paasta/monitoring/common/dao"
	cm "kr/paasta/monitoring/common/model"
	"kr/paasta/monitoring/iaas/integration"
	"kr/paasta/monitoring/iaas/model"
	pm "kr/paasta/monitoring/paas/model"
	"kr/paasta/monitoring/utils"
	ua "kr/paasta/monitoring/utils"
	//"strings"
	"strings"
	"time"
)

type LoginService struct {
	openstackProvider model.OpenstackProvider
	//CfProvider        cfclient.Config
	txn      *gorm.DB
	RdClient *redis.Client
	sysType  string
}

func GetLoginService(openstackProvider model.OpenstackProvider, txn *gorm.DB, rdClient *redis.Client, sysType string) *LoginService {
	return &LoginService{
		openstackProvider: openstackProvider,
		//	CfProvider : cfProvider,
		txn:      txn,
		RdClient: rdClient,
		sysType:  sysType,
	}
}

//func GetIaasLoginService(openstackProvider model.OpenstackProvider, txn *gorm.DB,  rdClient *redis.Client, sysType string ) *LoginService {
//	return &LoginService{
//		openstackProvider: openstackProvider,
//		txn: txn,
//		RdClient: rdClient,
//		sysType: sysType,
//	}
//}
//
//func GetPaasLoginService(cfProvider  cfclient.Config, txn *gorm.DB,  rdClient *redis.Client, sysType string ) *LoginService {
//	return &LoginService{
//		CfProvider : cfProvider,
//		txn: txn,
//		RdClient: rdClient,
//		sysType: sysType,
//	}
//}

func (n *LoginService) Logout(provider *gophercloud.ProviderClient, reqCsrfToken string) tokens.RevokeResult {
	result := integration.GetKeystone(n.openstackProvider, provider).RevokeToken()

	return result
}

func (n *LoginService) Login(req cm.UserInfo, reqCsrfToken string, cfConfig pm.CFConfig) (userInfo cm.UserInfo, provider *gophercloud.ProviderClient, err error) {

	//member info 로그인 아이디 패스워드를 회원디비에서 검색하여 iaas, paas 로그인 정보를 가지고 온다.
	result, _, dbErr := dao.GetLoginDao(n.txn).GetLoginMemberInfo(req, n.txn)
	resultString := ""
	if dbErr != nil {
		if strings.Contains(dbErr.Error(), "record not found") {
			return result, provider, fmt.Errorf("Member certification information is invalid.")
		} else {
			return result, provider, dbErr
		}
	}

	if strings.Contains(n.sysType, utils.SYS_TYPE_IAAS) || strings.Contains(n.sysType, utils.SYS_TYPE_ALL) {
		if result.IaasUserUseYn == "Y" {
			//get iaas token
			n.openstackProvider.Username = result.IaasUserId
			n.openstackProvider.Password = result.IaasUserPw
			provider, err = utils.GetAdminToken(n.openstackProvider)
			if err != nil {
				fmt.Println("iaas error::", err.Error())
				//return req, provider, err
			} else {
				result.IaasToken = provider.TokenID
				//fmt.Println("iaas token ::: ", result.IaasToken )
			}

			n.RdClient.HSet(req.Token, "iaasToken", result.IaasToken)
			n.RdClient.HSet(req.Token, "iaasUserId", result.IaasUserId)
		}

		//result.PaasAdminYn = "N"
	}

	if strings.Contains(n.sysType, utils.SYS_TYPE_PAAS) || strings.Contains(n.sysType, utils.SYS_TYPE_ALL) {
		if result.PaasUserUseYn == "Y" {
			cfConfig.UserId = result.PaasUserId
			cfConfig.UserPw = result.PaasUserPw
			cfConfig.Type = "PAAS"
			resultString, err = ua.GetUaaToken(result, reqCsrfToken, cfConfig, n.RdClient)

			if err != nil {
				fmt.Println("uaa token::", err.Error())
				//return req, provider, err
			} else {
				req.PaasToken = resultString
			}

			//resultString = GetMemberService(n.openstackProvider, n.CfProvider, n.txn, n.RdClient).GetUaaToken(req, reqCsrfToken, cfConfig)
			fmt.Println("paas token ::", resultString)

		}
	}

	if strings.Contains(n.sysType, utils.SYS_TYPE_CAAS) || strings.Contains(n.sysType, utils.SYS_TYPE_ALL) {
		if result.CaasUserUseYn == "Y" {
			cfConfig.UserId = result.CaasUserId
			cfConfig.UserPw = result.CaasUserPw
			cfConfig.Type = "CAAS"
			//resultCaasString := ""
			//resultCaasString = GetMemberService(n.openstackProvider, n.txn, n.RdClient).CaasServiceCheck(result, reqCsrfToken, cfConfig)
			//if resultCaasString == "adm" {
			resultString, err = ua.GetUaaToken(result, reqCsrfToken, cfConfig, n.RdClient)

			if err != nil {
				fmt.Println("uaa token::", err.Error())
				//return req, provider, err
			} else {
				req.CaasToken = resultString
			}

			//} else {
			//	fmt.Println("caas token get fail ::" + resultCaasString)
			//}
			//resultString = ua.GetUaaToken(req, reqCsrfToken, cfConfig,n.RdClient)
			//resultString = GetMemberService(n.openstackProvider, n.CfProvider, n.txn, n.RdClient).GetUaaToken(req, reqCsrfToken, cfConfig)
			//fmt.Println("caas token ::",resultString)
			//get caas token
			//n.CfProvider.Username = result.CaasUserId
			//n.CfProvider.Password = result.CaasUserPw
			//
			//
			//n.RdClient.HSet(req.Token,"caasToken",result.CaasToken)
			//n.RdClient.HSet(req.Token,"caasUserId",result.CaasUserId)
		}
	}

	//redis 에 사용자 정보를 저장한다.
	n.RdClient.HSet(req.Token, "userId", result.UserId)
	n.RdClient.HSet(req.Token, "userPw", result.UserPw)
	n.RdClient.HSet(req.Token, "paasAdminYn", n.RdClient.HMGet(req.Token, "paasAdminYn"))
	n.RdClient.Expire(req.Token, 30*60*time.Second)

	return result, provider, dbErr
}

func (s *LoginService) SetUserInfoCache(userInfo *cm.UserInfo, reqCsrfToken string, cfConfig pm.CFConfig) {

	var strcd1 = ""
	var strcd2 = ""
	var strcd3 = ""
	var strcd4 = ""
	var strcd5 = ""
	var strcd6 = ""

	userInfo.Username = userInfo.UserId

	if strings.Contains(s.sysType, utils.SYS_TYPE_IAAS) || strings.Contains(s.sysType, utils.SYS_TYPE_ALL) {
		if userInfo.IaasUserUseYn == "Y" {
			strcd1 = "I"
			var tokenRequest cm.UserInfo
			tokenRequest.IaasUserId = userInfo.IaasUserId
			tokenRequest.IaasUserPw = userInfo.IaasUserPw
			result := GetIaasMemberService(s.openstackProvider, s.txn, s.RdClient).GetIaasToken(tokenRequest, reqCsrfToken)

			if result != "" {
				strcd2 = "S"
			} else {
				strcd2 = "F"
			}

		} else {
			strcd1 = "F"
			strcd2 = "F"
			GetIaasMemberService(s.openstackProvider, s.txn, s.RdClient).DeleteIaasToken(reqCsrfToken)
		}
	}
	if strings.Contains(s.sysType, utils.SYS_TYPE_PAAS) || strings.Contains(s.sysType, utils.SYS_TYPE_ALL) {
		if userInfo.PaasUserUseYn == "Y" {
			strcd3 = "P"
			//var tokenRequest cm.UserInfo
			//tokenRequest.PaasUserId = userInfo.PaasUserId
			//tokenRequest.PaasUserPw = userInfo.PaasUserPw
			//cfConfig.Type = "PAAS"
			//result,err := ua.GetUaaToken(tokenRequest, reqCsrfToken, cfConfig, s.RdClient)
			//
			//if err != nil {
			//	fmt.Println("uaa token::", err.Error())
			//	//return req, provider, err
			//}

			//result := GetPaasMemberService(s.CfProvider, s.txn, s.RdClient).GetUaaToken(tokenRequest, reqCsrfToken, cfConfig)
			if userInfo.PaasToken != "" {
				strcd4 = "S"
			} else {
				strcd4 = "F"
			}

		} else {
			strcd3 = "F"
			strcd4 = "F"
			GetPaasMemberService(s.txn, s.RdClient).DeletePaasToken(reqCsrfToken)
		}
	}
	if strings.Contains(s.sysType, utils.SYS_TYPE_CAAS) || strings.Contains(s.sysType, utils.SYS_TYPE_ALL) {
		if userInfo.CaasUserUseYn == "Y" {
			strcd5 = "C"
			//var tokenRequest cm.UserInfo
			//tokenRequest.PaasUserId = userInfo.PaasUserId
			//tokenRequest.PaasUserPw = userInfo.PaasUserPw
			//cfConfig.Type = "PAAS"
			//result,err := ua.GetUaaToken(tokenRequest, reqCsrfToken, cfConfig, s.RdClient)
			//
			//if err != nil {
			//	fmt.Println("uaa token::", err.Error())
			//	//return req, provider, err
			//}

			//result := GetPaasMemberService(s.CfProvider, s.txn, s.RdClient).GetUaaToken(tokenRequest, reqCsrfToken, cfConfig)
			if userInfo.PaasToken != "" {
				strcd6 = "S"
			} else {
				strcd6 = "F"
			}

		} else {
			strcd5 = "F"
			strcd6 = "F"
			GetPaasMemberService(s.txn, s.RdClient).DeletePaasToken(reqCsrfToken)
		}
	}

	//fmt.Println("login User Auth ::: ", strcd1, strcd2, strcd3, strcd4)

	userInfo.AuthI1 = strcd1
	userInfo.AuthI2 = strcd2
	userInfo.AuthP1 = strcd3
	userInfo.AuthP2 = strcd4
	userInfo.AuthC1 = strcd5
	userInfo.AuthC2 = strcd6

	userInfo.IaasUserPw = ""
	userInfo.PaasUserPw = ""
	userInfo.CaasUserPw = ""
	userInfo.UserPw = ""
}

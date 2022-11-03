package login

import (
	"fmt"
	"monitoring-portal/common/dao/login"
	"monitoring-portal/common/service/member"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	//"github.com/cloudfoundry-community/go-cfclient"
	"github.com/jinzhu/gorm"

	commonModel "monitoring-portal/common/model"
	"monitoring-portal/iaas_new/integration"
	"monitoring-portal/iaas_new/model"
	paasModel "monitoring-portal/paas/model"
	"monitoring-portal/utils"
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

func (n *LoginService) Login(req commonModel.UserInfo, reqCsrfToken string, cfConfig paasModel.CFConfig) (userInfo commonModel.UserInfo, provider *gophercloud.ProviderClient, err error) {

	result, _, dbErr := login.GetLoginDao(n.txn).GetLoginMemberInfo(req, n.txn) // 회원 정보 조회
	resultString := ""
	if dbErr != nil {
		if strings.Contains(dbErr.Error(), "record not found") {
			return result, provider, fmt.Errorf("Member certification information is invalid.")
		} else {
			return result, provider, dbErr
		}
	}

	// IaaS용 오픈스택 토큰 발급받기
	if strings.Contains(n.sysType, utils.SYS_TYPE_IAAS) || strings.Contains(n.sysType, utils.SYS_TYPE_ALL) {
		if result.IaasUserUseYn == "Y" {

			n.openstackProvider.Username = result.IaasUserId
			n.openstackProvider.Password = result.IaasUserPw

			provider, err = utils.GetOpenstackToken(n.openstackProvider)
			if err != nil {
				utils.Logger.Error(err.Error())
				//return req, provider, err
			} else {
				result.IaasToken = provider.TokenID
				fmt.Println("iaas token ::: ", result.IaasToken )
			}

			utils.Logger.Debugf("req.Token(CSRF_TOKEN) : %v\n", req.Token)
			utils.Logger.Debugf("result.IaasToken : %v\n", result.IaasToken)

			n.RdClient.HSet(req.Token, "iaasToken", result.IaasToken)
			n.RdClient.HSet(req.Token, "iaasUserId", result.IaasUserId)
		}

		//result.PaasAdminYn = "N"
	}

	// PaaS용 UAA 토큰 발급받기
	if strings.Contains(n.sysType, utils.SYS_TYPE_PAAS) || strings.Contains(n.sysType, utils.SYS_TYPE_ALL) {
		if result.PaasUserUseYn == "Y" {
			cfConfig.UserId = result.PaasUserId
			cfConfig.UserPw = result.PaasUserPw
			cfConfig.Type = "PAAS"
			resultString, err = utils.GetUaaToken(result, reqCsrfToken, cfConfig, n.RdClient)

			if err != nil {
				utils.Logger.Error(err.Error())
				//return req, provider, err
			} else {
				req.PaasToken = resultString
			}

			//resultString = GetMemberService(n.openstackProvider, n.CfProvider, n.txn, n.RdClient).GetUaaToken(req, reqCsrfToken, cfConfig)
			//utils.Logger.Debugf("paas token : %v\n", resultString)

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
			resultString, err = utils.GetUaaToken(result, reqCsrfToken, cfConfig, n.RdClient)

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

func (s *LoginService) SetUserInfoCache(userInfo *commonModel.UserInfo, reqCsrfToken string, cfConfig paasModel.CFConfig) {

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
			var tokenRequest commonModel.UserInfo
			tokenRequest.IaasUserId = userInfo.IaasUserId
			tokenRequest.IaasUserPw = userInfo.IaasUserPw

			result := member.GetIaasMemberService(s.openstackProvider, s.txn, s.RdClient).GetIaasToken(tokenRequest, reqCsrfToken)
			model.MonitLogger.Debug(result)

			if result != "" {
				strcd2 = "S"
			} else {
				strcd2 = "F"
			}

		} else {
			strcd1 = "F"
			strcd2 = "F"
			member.GetIaasMemberService(s.openstackProvider, s.txn, s.RdClient).DeleteIaasToken(reqCsrfToken)
		}
	}
	if strings.Contains(s.sysType, utils.SYS_TYPE_PAAS) || strings.Contains(s.sysType, utils.SYS_TYPE_ALL) {
		if userInfo.PaasUserUseYn == "Y" {
			strcd3 = "P"
			//var tokenRequest commonModel.UserInfo
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
			member.GetPaasMemberService(s.txn, s.RdClient).DeletePaasToken(reqCsrfToken)
		}
	}
	if strings.Contains(s.sysType, utils.SYS_TYPE_CAAS) || strings.Contains(s.sysType, utils.SYS_TYPE_ALL) {
		if userInfo.CaasUserUseYn == "Y" {
			strcd5 = "C"
			//var tokenRequest commonModel.UserInfo
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
			member.GetPaasMemberService(s.txn, s.RdClient).DeletePaasToken(reqCsrfToken)
		}
	}

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

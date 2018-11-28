package controller

import (
	"kr/paasta/monitoring/iaas/model"
	cm "kr/paasta/monitoring/common/model"
	"kr/paasta/monitoring/common/service"
	"kr/paasta/monitoring/utils"
	"net/http"
	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/jinzhu/gorm"
	"fmt"
	"github.com/go-redis/redis"
	"encoding/json"
	//"strings"
)

//Compute Node Controller
type MemberController struct {
	OpenstackProvider model.OpenstackProvider
	CfProvider        cfclient.Config
	txn *gorm.DB
	RdClient *redis.Client
	sysType string
}

func NewMemberController(openstackProvider model.OpenstackProvider, cfProvider  cfclient.Config, txn *gorm.DB,  rdClient *redis.Client , sysType string ) *MemberController {
	return &MemberController{
		OpenstackProvider: openstackProvider,
		CfProvider: cfProvider,
		txn: txn,
		RdClient: rdClient,
		sysType : sysType,
	}

}

func NewIaasMemberController(openstackProvider model.OpenstackProvider, txn *gorm.DB,  rdClient *redis.Client , sysType string) *MemberController {
	return &MemberController{
		OpenstackProvider: openstackProvider,
		txn: txn,
		RdClient: rdClient,
		sysType : sysType,
	}

}

func NewPaasMemberController(cfProvider  cfclient.Config, txn *gorm.DB,  rdClient *redis.Client , sysType string) *MemberController {
	return &MemberController{
		CfProvider: cfProvider,
		txn: txn,
		RdClient: rdClient,
		sysType : sysType,
	}

}

func (s *MemberController) MemberJoinInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("MemberController MemberJoinInfo enter!!")

	utils.RenderJsonLogoutResponse(s.sysType, w)
}

func (s *MemberController) MemberJoinSave(w http.ResponseWriter, r *http.Request) {

	fmt.Println("MemberController MemberJoinSave enter!!")

	var apiRequest cm.UserInfo
	err := json.NewDecoder(r.Body).Decode(&apiRequest)
	if err != nil {
		fmt.Println("MemberController MemberJoinSave error !!",err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	} else {

		var loginErr model.ErrMessage

		if s.sysType == utils.SYS_TYPE_IAAS{
			err := services.GetIaasMemberService(s.OpenstackProvider, s.txn, s.RdClient).MemberJoinSave(apiRequest)
			loginErr = utils.GetError().GetCheckErrorMessage(err)
		}else if s.sysType == utils.SYS_TYPE_PAAS{
			err := services.GetPaasMemberService(s.CfProvider, s.txn, s.RdClient).MemberJoinSave(apiRequest)
			loginErr = utils.GetError().GetCheckErrorMessage(err)
		}else{
			err := services.GetMemberService(s.OpenstackProvider, s.CfProvider, s.txn, s.RdClient).MemberJoinSave(apiRequest)
			loginErr = utils.GetError().GetCheckErrorMessage(err)
		}

		if loginErr != nil {
			utils.ErrRenderJsonResponse(loginErr, w)
			return
		}

		utils.RenderJsonLogoutResponse(nil, w)
	}

}


func (s *MemberController) MemberAuthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("MemberController MemberAuthCheck enter!!")

	reqCsrfToken := r.Header.Get(model.CSRF_TOKEN_NAME)
	fmt.Println("CSRF_TOKEN_NAME=>",reqCsrfToken)
/*
	var apiRequest cm.UserInfo
	id := r.FormValue(":id")
	apiRequest.UserId = id
	apiRequest.Token = reqCsrfToken
	userInfo,  err := services.GetMemberService(s.OpenstackProvider, s.CfProvider, s.txn, s.RdClient).MemberAuthCheck(apiRequest)
	loginErr := utils.GetError().GetCheckErrorMessage(err)

     var apiRequest01 cm.UserInfo
	apiRequest01.UserId = userInfo.UserId
	apiRequest01.IaasToken = userInfo.IaasToken
	apiRequest01.PaasToken = userInfo.PaasToken

	if loginErr != nil {
		utils.ErrRenderJsonResponse(loginErr, w)
		return
	}
*/
	var apiRequest cm.UserInfo
	id := r.FormValue(":id")
	apiRequest.UserId = id

	//캐쉬 정보중 사용자 정보 가져오기


	var userInfo cm.UserInfo
	var err error

	if s.sysType == utils.SYS_TYPE_IAAS{
		userInfo, _, err = services.GetIaasMemberService(s.OpenstackProvider, s.txn, s.RdClient).MemberInfoView(apiRequest)
	}else if s.sysType == utils.SYS_TYPE_PAAS{
		userInfo, _, err = services.GetPaasMemberService(s.CfProvider, s.txn, s.RdClient).MemberInfoView(apiRequest)
	}else{
		userInfo, _, err = services.GetMemberService(s.OpenstackProvider, s.CfProvider, s.txn, s.RdClient).MemberInfoView(apiRequest)
	}

	if err != nil {
		utils.ErrRenderJsonResponse(err, w)
	}else{
		//캐쉬 정보 생성
		services.GetLoginService(s.OpenstackProvider, s.CfProvider , s.txn, s.RdClient, s.sysType).SetUserInfoCache(&userInfo, reqCsrfToken)
		userInfo.SysType = s.sysType
		utils.RenderJsonResponse(userInfo, w)
	}


}



func (s *MemberController) MemberCheckId(w http.ResponseWriter, r *http.Request) {
	fmt.Println("MemberController MemberCheckId enter!!")
	var apiRequest cm.UserInfo
	id := r.FormValue(":id")
	apiRequest.UserId = id

	var userInfo cm.UserInfo
	var loginErr model.ErrMessage
	var err error

	if s.sysType == utils.SYS_TYPE_IAAS{
		userInfo, _, err = services.GetIaasMemberService(s.OpenstackProvider, s.txn, s.RdClient).MemberInfoCheck(apiRequest)
		loginErr = utils.GetError().GetCheckErrorMessage(err)
	}else if s.sysType == utils.SYS_TYPE_PAAS{
		userInfo, _, err = services.GetPaasMemberService(s.CfProvider, s.txn, s.RdClient).MemberInfoCheck(apiRequest)
		loginErr = utils.GetError().GetCheckErrorMessage(err)
	}else{
		userInfo, _, err = services.GetMemberService(s.OpenstackProvider, s.CfProvider, s.txn, s.RdClient).MemberInfoCheck(apiRequest)
		loginErr = utils.GetError().GetCheckErrorMessage(err)
	}

	if loginErr != nil {
		utils.ErrRenderJsonResponse(loginErr, w)
		return
	}

	utils.RenderJsonLogoutResponse(userInfo.UserId, w)
}

func (s *MemberController) MemberCheckEmail(w http.ResponseWriter, r *http.Request) {
	fmt.Println("MemberController MemberCheckEmail enter!!")
	var apiRequest cm.UserInfo
	email := r.FormValue(":email")
	apiRequest.UserEmail = email

	var userInfo cm.UserInfo
	var loginErr model.ErrMessage
	var err error

	if s.sysType == utils.SYS_TYPE_IAAS{
		userInfo, _, err = services.GetIaasMemberService(s.OpenstackProvider, s.txn, s.RdClient).MemberInfoCheck(apiRequest)
		loginErr = utils.GetError().GetCheckErrorMessage(err)
	}else if s.sysType == utils.SYS_TYPE_PAAS{
		userInfo, _, err = services.GetPaasMemberService(s.CfProvider, s.txn, s.RdClient).MemberInfoCheck(apiRequest)
		loginErr = utils.GetError().GetCheckErrorMessage(err)
	}else{
		userInfo, _, err = services.GetMemberService(s.OpenstackProvider, s.CfProvider, s.txn, s.RdClient).MemberInfoCheck(apiRequest)
		loginErr = utils.GetError().GetCheckErrorMessage(err)
	}

	if loginErr != nil {
		utils.ErrRenderJsonResponse(loginErr, w)
		return
	}

	utils.RenderJsonLogoutResponse(userInfo.UserEmail, w)
}

func (s *MemberController) MemberJoinCheckDuplicationIaasId(w http.ResponseWriter, r *http.Request) {
	var apiRequest cm.UserInfo
	apiRequest.IaasUserId = r.FormValue(":id")

	var userInfo cm.UserInfo
	var loginErr model.ErrMessage
	var err error

	if s.sysType == utils.SYS_TYPE_IAAS{
		userInfo, err = services.GetIaasMemberService(s.OpenstackProvider, s.txn, s.RdClient).MemberJoinCheckDuplicationIaasId(apiRequest)
		loginErr = utils.GetError().GetCheckErrorMessage(err)
	}else if s.sysType == utils.SYS_TYPE_PAAS{
		userInfo, err = services.GetPaasMemberService(s.CfProvider, s.txn, s.RdClient).MemberJoinCheckDuplicationIaasId(apiRequest)
		loginErr = utils.GetError().GetCheckErrorMessage(err)
	}else{
		userInfo, err = services.GetMemberService(s.OpenstackProvider, s.CfProvider, s.txn, s.RdClient).MemberJoinCheckDuplicationIaasId(apiRequest)
		loginErr = utils.GetError().GetCheckErrorMessage(err)
	}

	if loginErr != nil {
		utils.ErrRenderJsonResponse(loginErr, w)
		return
	}
	utils.RenderJsonLogoutResponse(userInfo.IaasUserId, w)
}

func (s *MemberController) MemberJoinCheckDuplicationPaasId(w http.ResponseWriter, r *http.Request) {
	var apiRequest cm.UserInfo
	apiRequest.PaasUserId = r.FormValue(":id")

	var userInfo cm.UserInfo
	var loginErr model.ErrMessage
	var err error

	if s.sysType == utils.SYS_TYPE_IAAS{
		userInfo, err = services.GetIaasMemberService(s.OpenstackProvider, s.txn, s.RdClient).MemberJoinCheckDuplicationPaasId(apiRequest)
		loginErr = utils.GetError().GetCheckErrorMessage(err)
	}else if s.sysType == utils.SYS_TYPE_PAAS{
		userInfo, err = services.GetPaasMemberService(s.CfProvider, s.txn, s.RdClient).MemberJoinCheckDuplicationPaasId(apiRequest)
		loginErr = utils.GetError().GetCheckErrorMessage(err)
	}else{
		userInfo, err = services.GetMemberService(s.OpenstackProvider, s.CfProvider, s.txn, s.RdClient).MemberJoinCheckDuplicationPaasId(apiRequest)
		loginErr = utils.GetError().GetCheckErrorMessage(err)
	}

	if loginErr != nil {
		utils.ErrRenderJsonResponse(loginErr, w)
		return
	}
	utils.RenderJsonLogoutResponse(userInfo.PaasUserId, w)
}

func (s *MemberController) MemberCheckIaaS(w http.ResponseWriter, r *http.Request) {
	fmt.Println("MemberController MemberCheckIaaS enter!!")
	reqCsrfToken := r.Header.Get(model.CSRF_TOKEN_NAME)
    result := ""
	var apiRequest cm.UserInfo
	_ = json.NewDecoder(r.Body).Decode(&apiRequest)

	if s.sysType == utils.SYS_TYPE_IAAS{
		result = services.GetIaasMemberService(s.OpenstackProvider,s.txn, s.RdClient).GetIaasToken( apiRequest, reqCsrfToken)
	}else if s.sysType == utils.SYS_TYPE_PAAS{
		result = ""
	}else{
		result = services.GetMemberService(s.OpenstackProvider, s.CfProvider, s.txn, s.RdClient).GetIaasToken( apiRequest, reqCsrfToken)
	}

	utils.RenderJsonLogoutResponse(result, w)

}


func (s *MemberController) MemberCheckPaaS(w http.ResponseWriter, r *http.Request) {
	fmt.Println("MemberController MemberCheckPaaS enter!!")
	reqCsrfToken := r.Header.Get(model.CSRF_TOKEN_NAME)
	result := ""
	var apiRequest cm.UserInfo
	_ = json.NewDecoder(r.Body).Decode(&apiRequest)

	if s.sysType == utils.SYS_TYPE_IAAS{
		result = ""
	}else if s.sysType == utils.SYS_TYPE_PAAS{
		result = services.GetPaasMemberService( s.CfProvider, s.txn, s.RdClient).GetPaasToken(apiRequest, reqCsrfToken)
	}else{
		result = services.GetMemberService(s.OpenstackProvider, s.CfProvider, s.txn, s.RdClient).GetPaasToken(apiRequest, reqCsrfToken)
	}



	utils.RenderJsonLogoutResponse(result, w)
}


func (s *MemberController) MemberInfoView(w http.ResponseWriter, r *http.Request) {
	fmt.Println("MemberController MemberInfoView enter!!")
	var apiRequest cm.UserInfo
	err := json.NewDecoder(r.Body).Decode(&apiRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	} else {

		var userInfo cm.UserInfo
		var loginErr model.ErrMessage
		var err error

		if s.sysType == utils.SYS_TYPE_IAAS{
			userInfo, _, err = services.GetIaasMemberService(s.OpenstackProvider, s.txn, s.RdClient).MemberInfoView(apiRequest)
			loginErr = utils.GetError().GetCheckErrorMessage(err)
		}else if s.sysType == utils.SYS_TYPE_PAAS{
			userInfo, _, err = services.GetPaasMemberService(s.CfProvider, s.txn, s.RdClient).MemberInfoView(apiRequest)
			loginErr = utils.GetError().GetCheckErrorMessage(err)
		}else{
			userInfo, _, err = services.GetMemberService(s.OpenstackProvider, s.CfProvider, s.txn, s.RdClient).MemberInfoView(apiRequest)
			loginErr = utils.GetError().GetCheckErrorMessage(err)
		}

		if loginErr != nil {
			utils.ErrRenderJsonResponse(loginErr, w)
			return
		}

		utils.RenderJsonLogoutResponse(userInfo, w)
	}
}

func (s *MemberController) MemberInfoUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("MemberController MemberInfoUpdate enter!!")
	reqCsrfToken := r.Header.Get(model.CSRF_TOKEN_NAME)
	var apiRequest cm.UserInfo
	err := json.NewDecoder(r.Body).Decode(&apiRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	} else {

		var userInfo cm.UserInfo
		var loginErr model.ErrMessage
		var err error

		if s.sysType == utils.SYS_TYPE_IAAS{
			userInfo, _, err = services.GetIaasMemberService(s.OpenstackProvider, s.txn, s.RdClient).MemberInfoUpdate(apiRequest)
			loginErr = utils.GetError().GetCheckErrorMessage(err)
		}else if s.sysType == utils.SYS_TYPE_PAAS{
			userInfo, _, err = services.GetPaasMemberService( s.CfProvider, s.txn, s.RdClient).MemberInfoUpdate(apiRequest)
			loginErr = utils.GetError().GetCheckErrorMessage(err)
		}else{
			userInfo, _, err = services.GetMemberService(s.OpenstackProvider, s.CfProvider, s.txn, s.RdClient).MemberInfoUpdate(apiRequest)
			loginErr = utils.GetError().GetCheckErrorMessage(err)
		}

		if loginErr != nil {
			utils.ErrRenderJsonResponse(loginErr, w)
			return
		}

		//회원정보 수정후 변경된 정보를 캐쉬에 넣기위해 캐쉬 정보 생성
		services.GetLoginService(s.OpenstackProvider, s.CfProvider , s.txn, s.RdClient, s.sysType).SetUserInfoCache(&userInfo, reqCsrfToken)

		userInfo.SysType = s.sysType

		utils.RenderJsonLogoutResponse(userInfo, w)
	}
}

func (s *MemberController) MemberInfoDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("MemberController MemberInfoDelete enter!!")
	var apiRequest cm.UserInfo
	id := r.FormValue(":id")
	apiRequest.UserId = id

	var cnt int
	var loginErr model.ErrMessage
	var err error

	if s.sysType == utils.SYS_TYPE_IAAS{
		cnt, err = services.GetIaasMemberService(s.OpenstackProvider, s.txn, s.RdClient).MemberInfoDelete(apiRequest)
		loginErr = utils.GetError().GetCheckErrorMessage(err)
	}else if s.sysType == utils.SYS_TYPE_PAAS{
		cnt, err = services.GetPaasMemberService( s.CfProvider, s.txn, s.RdClient).MemberInfoDelete(apiRequest)
		loginErr = utils.GetError().GetCheckErrorMessage(err)
	}else{
		cnt, err = services.GetMemberService(s.OpenstackProvider, s.CfProvider, s.txn, s.RdClient).MemberInfoDelete(apiRequest)
		loginErr = utils.GetError().GetCheckErrorMessage(err)
	}

	if loginErr != nil {
		utils.ErrRenderJsonResponse(loginErr, w)
		return
	}

	utils.RenderJsonLogoutResponse(cnt, w)
}


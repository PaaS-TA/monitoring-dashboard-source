package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	//"github.com/cloudfoundry-community/go-cfclient"
	monascagopher "github.com/gophercloud/gophercloud"
	"github.com/monasca/golang-monascaclient/monascaclient"
	//"github.com/rackspace/gophercloud"
	//"github.com/cloudfoundry-community/go-cfclient"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	cm "kr/paasta/monitoring/common/model"
	"kr/paasta/monitoring/common/service"
	"kr/paasta/monitoring/iaas/model"
	pm "kr/paasta/monitoring/paas/model"
	"kr/paasta/monitoring/utils"
	"net/http"
)

//Compute Node Controller
type LoginController struct {
	OpenstackProvider model.OpenstackProvider
	MonAuth           monascagopher.AuthOptions
	MonClient         monascaclient.Client
	//CfProvider        cfclient.Config
	txn      *gorm.DB
	RdClient *redis.Client
	sysType  string
	CfConfig pm.CFConfig
}

func NewLoginController(openstackProvider model.OpenstackProvider, monsClient monascaclient.Client, auth monascagopher.AuthOptions, txn *gorm.DB, rdClient *redis.Client, sysType string, cfConfig pm.CFConfig) *LoginController {
	return &LoginController{
		OpenstackProvider: openstackProvider,
		MonAuth:           auth,
		MonClient:         monsClient,
		//CfProvider: cfProvider,
		txn:      txn,
		RdClient: rdClient,
		sysType:  sysType,
		CfConfig: cfConfig,
	}

}

func NewIaasLoginController(openstackProvider model.OpenstackProvider, monsClient monascaclient.Client, auth monascagopher.AuthOptions, txn *gorm.DB, rdClient *redis.Client, sysType string) *LoginController {
	return &LoginController{
		OpenstackProvider: openstackProvider,
		MonAuth:           auth,
		MonClient:         monsClient,
		txn:               txn,
		RdClient:          rdClient,
		sysType:           sysType,
	}

}

func NewPaasLoginController(txn *gorm.DB, rdClient *redis.Client, sysType string) *LoginController {
	return &LoginController{
		//CfProvider: cfProvider,
		txn:      txn,
		RdClient: rdClient,
		sysType:  sysType,
	}

}

func (s *LoginController) Ping(w http.ResponseWriter, r *http.Request) {

	token, _ := utils.GenerateRandomString(32)
	//session := model.SessionManager.Load(r)

	testToken := r.Header.Get(model.TEST_TOKEN_NAME)
	if testToken != "" {
		w.Header().Add(model.TEST_TOKEN_NAME, token)
	} else {
		//fmt.Println("pint Token::::", token)
		//session.PutString(w, token, token)
		w.Header().Add(model.CSRF_TOKEN_NAME, token)
	}

	utils.RenderJsonResponse(nil, w)

}

func (s *LoginController) Login(w http.ResponseWriter, r *http.Request) {

	reqCsrfToken := r.Header.Get(model.CSRF_TOKEN_NAME)

	fmt.Println("Login Test CSRF_TOKEN: !!!", reqCsrfToken)

	var apiRequest cm.UserInfo
	apiRequest.Token = reqCsrfToken

	err := json.NewDecoder(r.Body).Decode(&apiRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	} else {

		err := loginValidate(apiRequest)
		if err != nil {
			loginErr := utils.GetError().GetCheckErrorMessage(err)
			utils.ErrRenderJsonResponse(loginErr, w)
			return
		}

		var userInfo cm.UserInfo
		var loginErr model.ErrMessage

		// check saas,caas
		userInfo, _, err = services.GetLoginService(s.OpenstackProvider, s.txn, s.RdClient, s.sysType).Login(apiRequest, reqCsrfToken, s.CfConfig)
		loginErr = utils.GetError().GetCheckErrorMessage(err)

		//if s.sysType == utils.SYS_TYPE_IAAS{
		//	userInfo, _, err = services.GetIaasLoginService(s.OpenstackProvider, s.txn, s.RdClient, s.sysType).Login(apiRequest)
		//	loginErr = utils.GetError().GetCheckErrorMessage(err)
		//}else if s.sysType == utils.SYS_TYPE_PAAS{
		//	userInfo, _, err = services.GetPaasLoginService( s.CfProvider , s.txn, s.RdClient, s.sysType).Login(apiRequest)
		//	loginErr = utils.GetError().GetCheckErrorMessage(err)
		//}else{
		//	userInfo, _, err = services.GetLoginService(s.OpenstackProvider, s.CfProvider , s.txn, s.RdClient, s.sysType).Login(apiRequest)
		//	loginErr = utils.GetError().GetCheckErrorMessage(err)
		//}

		if loginErr != nil {
			utils.ErrRenderJsonResponse(loginErr, w)
			return
		} else {

			services.GetLoginService(s.OpenstackProvider, s.txn, s.RdClient, s.sysType).SetUserInfoCache(&userInfo, reqCsrfToken, s.CfConfig)
			userInfo.SysType = s.sysType
			utils.RenderJsonResponse(userInfo, w)
			return
		}
	}
}

func (s *LoginController) Logout(w http.ResponseWriter, r *http.Request) {

	reqCsrfToken := r.Header.Get(model.CSRF_TOKEN_NAME)

	fmt.Println("logout reqCsrfToken=", reqCsrfToken)

	s.RdClient.Del(reqCsrfToken)

	//provider, _, _ := utils.GetOpenstackProvider(r)
	//services.GetLoginService(s.OpenstackProvider, s.CfProvider, s.txn, s.RdClient).Logout(provider,reqCsrfToken)
	//utils.RenderJsonLogoutResponse(nil, w)

}

func loginValidate(apiRequest cm.UserInfo) error {

	if apiRequest.Username == "" {
		return errors.New("Required input value does not exist. [username]")
	}

	if apiRequest.Password == "" {
		return errors.New("Required input value does not exist. [password]")
	}

	return nil
}

func (s *LoginController) Join(w http.ResponseWriter, r *http.Request) {

	services.GetLoginService(s.OpenstackProvider, s.txn, s.RdClient, s.sysType)
	//if s.sysType == utils.SYS_TYPE_IAAS{
	//	services.GetIaasLoginService(s.OpenstackProvider, s.txn, s.RdClient, s.sysType)
	//}else if s.sysType == utils.SYS_TYPE_PAAS{
	//	services.GetPaasLoginService(s.CfProvider, s.txn, s.RdClient, s.sysType)
	//}else{
	//	services.GetLoginService(s.OpenstackProvider, s.CfProvider, s.txn, s.RdClient, s.sysType)
	//}

	utils.RenderJsonLogoutResponse(nil, w)

}

func (s *LoginController) Main(w http.ResponseWriter, r *http.Request) {
	model.MonitLogger.Debug("Main API Called")

	url := "/public/dist/index.html"
	http.Redirect(w, r, url, 302)
}

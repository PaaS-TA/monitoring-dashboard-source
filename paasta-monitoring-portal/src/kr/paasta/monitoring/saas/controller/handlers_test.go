package controller

import (
	"fmt"
	"monitoring-portal/common/controller/login"
	"monitoring-portal/common/controller/member"
	"strings"

	"github.com/cloudfoundry-community/gogobosh"
	"github.com/go-redis/redis"
	monascagopher "github.com/gophercloud/gophercloud"
	"github.com/influxdata/influxdb1-client/v2"
	"github.com/jinzhu/gorm"
	"github.com/monasca/golang-monascaclient/monascaclient"
	"github.com/rackspace/gophercloud"
	tokens3 "github.com/rackspace/gophercloud/openstack/identity/v3/tokens"
	"github.com/tedsuo/rata"
	"gopkg.in/olivere/elastic.v3"
	"io"
	"monitoring-portal/iaas_new/model"
	pm "monitoring-portal/paas/model"
	"monitoring-portal/routes"
	"monitoring-portal/utils"
	"net/http"
	"time"
)

func NewHandler(openstack_provider model.OpenstackProvider, iaasInfluxClient client.Client, paasInfluxClient client.Client,
	iaasTxn *gorm.DB, paasTxn *gorm.DB, iaasElasticClient *elastic.Client, paasElasticClient *elastic.Client, monsClient monascaclient.Client,
	auth monascagopher.AuthOptions, databases pm.Databases, rdClient *redis.Client, sysType string, boshClient *gogobosh.Client, cfConfig pm.CFConfig) http.Handler {

	//Controller선언
	var loginController *login.LoginController
	var memberController *member.MemberController

	// SaaS Metrics
	var applicationController *SaasController

	loginController = login.NewLoginController(openstack_provider, monsClient, auth, paasTxn, rdClient, sysType, cfConfig)
	memberController = member.NewMemberController(openstack_provider, paasTxn, rdClient, sysType, cfConfig)

	var saasActions rata.Handlers
	// add SAAS
	if strings.Contains(sysType, utils.SYS_TYPE_SAAS) || sysType == utils.SYS_TYPE_ALL {
		applicationController = NewSaasController(paasTxn)

		saasActions = rata.Handlers{
			routes.SAAS_API_APPLICATION_LIST:   route(applicationController.GetApplicationList),
			routes.SAAS_API_APPLICATION_STATUS: route(applicationController.GetAgentStatus),
			routes.SAAS_API_APPLICATION_GAUGE:  route(applicationController.GetAgentGaugeTot),
			routes.SAAS_API_APPLICATION_REMOVE: route(applicationController.RemoveApplication),

			routes.SAAS_ALARM_INFO:     route(applicationController.GetAlarmInfo),
			routes.SAAS_ALARM_UPDATE:   route(applicationController.GetAlarmUpdate),
			routes.SAAS_ALARM_LOG:      route(applicationController.GetAlarmLog),
			routes.SAAS_ALARM_SNS_INFO: route(applicationController.GetSnsInfo),
			routes.SAAS_ALARM_COUNT:    route(applicationController.GetAlarmCount),
			routes.SAAS_ALARM_SNS_SAVE: route(applicationController.GetlarmSnsSave),

			routes.SAAS_ALARM_STATUS_UPDATE:      route(applicationController.UpdateAlarmState),
			routes.SAAS_ALARM_ACTION:             route(applicationController.CreateAlarmResolve),
			routes.SAAS_ALARM_ACTION_DELETE:      route(applicationController.DeleteAlarmResolve),
			routes.SAAS_ALARM_ACTION_UPDATE:      route(applicationController.UpdateAlarmResolve),
			routes.SAAS_ALARM_SNS_CHANNEL_LIST:   route(applicationController.GetAlarmSnsReceiver),
			routes.SAAS_ALARM_SNS_CHANNEL_DELETE: route(applicationController.DeleteAlarmSnsChannel),
			routes.SAAS_ALARM_ACTION_LIST:        route(applicationController.GetAlarmActionList),
		}
	}

	commonActions := rata.Handlers{

		routes.PING:   route(loginController.Ping),
		routes.LOGIN:  route(loginController.Login),
		routes.LOGOUT: route(loginController.Logout),

		routes.MEMBER_JOIN_INFO:        route(memberController.MemberJoinInfo),
		routes.MEMBER_JOIN_SAVE:        route(memberController.MemberJoinSave),
		routes.MEMBER_JOIN_CHECK_ID:    route(memberController.MemberCheckId),
		routes.MEMBER_JOIN_CHECK_EMAIL: route(memberController.MemberCheckEmail),

		routes.MEMBER_AUTH_CHECK:  route(memberController.MemberAuthCheck),
		routes.MEMBER_INFO_VIEW:   route(memberController.MemberInfoView),
		routes.MEMBER_INFO_UPDATE: route(memberController.MemberInfoUpdate),
		routes.MEMBER_INFO_DELETE: route(memberController.MemberInfoDelete),

		// Html
		routes.Main: route(loginController.Main),
		//routes.Main: route(mainController.Main),
		routes.Static: route(StaticHandler),
	}

	var actions rata.Handlers
	var actionlist []rata.Handlers

	var route rata.Routes
	var routeList []rata.Routes

	// add SAAS , CAAS routes
	actionlist = append(actionlist, commonActions)

	if strings.Contains(sysType, utils.SYS_TYPE_SAAS) || sysType == utils.SYS_TYPE_ALL {
		actionlist = append(actionlist, saasActions)
		routeList = append(routeList, routes.SaasRoutes)
	}

	actions = getActions(actionlist)

	routeList = append(routeList, routes.Routes)
	route = getRoutes(routeList)

	handler, err := rata.NewRouter(route, actions)
	if err != nil {
		panic("unable to create router: " + err.Error())
	}
	fmt.Println("Monit Application Started")
	return HttpWrap(handler, rdClient, openstack_provider, cfConfig)
}

func getActions(list []rata.Handlers) rata.Handlers {
	actions := make(map[string]http.Handler)

	for _, value := range list {
		for key, val := range value {
			actions[key] = val
		}
	}
	return actions
}

func getRoutes(list []rata.Routes) rata.Routes {
	var rList []rata.Route

	for _, value := range list {
		for _, val := range value {
			rList = append(rList, val)
		}
	}
	return rList
}

func HttpWrap(handler http.Handler, rdClient *redis.Client, openstack_provider model.OpenstackProvider, cfConfig pm.CFConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, X-XSRF-TOKEN, Accept-Encoding, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Expose-Headers", "X-XSRF-TOKEN")
		}

		// Stop here if its Preflighted OPTIONS request
		if r.Method == "OPTIONS" {
			return
		}

		// token Pass
		if r.RequestURI != "/v2/login" && r.RequestURI != "/v2/logout" && !strings.Contains(r.RequestURI, "/v2/member/join") && r.RequestURI != "/v2/ping" && r.RequestURI != "/" && !strings.Contains(r.RequestURI, "/public/") && !strings.Contains(r.RequestURI, "/v2/paas/app/") && !strings.Contains(r.RequestURI, "/v2/caas/monitoring/podList") {
			fmt.Println("Request URI :: ", r.RequestURI)

			reqToken := r.Header.Get(model.CSRF_TOKEN_NAME)
			if reqToken == "0" || reqToken == "null" {
				fmt.Println("HttpWrap Hander reqToken is null ")
				errMessage := model.ErrMessage{"Message": "UnAuthrized"}
				utils.RenderJsonUnAuthResponse(errMessage, http.StatusUnauthorized, w)
			} else {
				//fmt.Println("HttpWrap Hander reqToken =",len(reqToken),":",reqToken)
				//모든 경로의 redis 의 토큰 정보를 확인한다
				val := rdClient.HGetAll(reqToken).Val()
				if val == nil || len(val) == 0 { // redis 에서 token 정보가 expire 된경우 로그인 화면으로 돌아간다
					fmt.Println("HttpWrap Hander redis.iaas_userid is null ")
					errMessage := model.ErrMessage{"Message": "UnAuthrized"}
					utils.RenderJsonUnAuthResponse(errMessage, http.StatusUnauthorized, w)
				} else {

					if strings.Contains(r.RequestURI, "/v2/member") && val["userId"] != "" {

						handler.ServeHTTP(w, r)

					} else if strings.Contains(r.RequestURI, "/v2/iaas") && val["iaasToken"] != "" && val["iaasUserId"] != "" { // IaaS 토큰 정보가 있는경우

						provider1, _, err := utils.GetOpenstackProvider(r)
						if err != nil || provider1 == nil {
							errMessage := model.ErrMessage{"Message": "UnAuthrized"}
							utils.RenderJsonUnAuthResponse(errMessage, http.StatusUnauthorized, w)
						} else {
							v3Client := NewIdentityV3(provider1)

							//IaaS, token 검증
							bool, err := tokens3.Validate(v3Client, val["iaasToken"])
							if err != nil || bool == false {
								//errMessage := model.ErrMessage{"Message": "UnAuthrized"}
								//utils.RenderJsonUnAuthResponse(errMessage, http.StatusUnauthorized, w)
								fmt.Println("iaas token validate error::", err)
								handler.ServeHTTP(w, r)
							} else {
								//두개 token 이 없는 경우도 고려 해야함
								rdClient.Expire(reqToken, 30*60*time.Second)
								handler.ServeHTTP(w, r)
							}
						}

					} else if strings.Contains(r.RequestURI, "/v2/paas") && val["paasRefreshToken"] != "" { // PaaS 토큰 정보가 있는경우

						// Pass token 검증 로직 추가
						//get paas token
						//cfProvider.Token = val["paasToken"]
						t1, _ := time.Parse(time.RFC3339, val["paasExpire"])
						if t1.Before(time.Now()) {
							fmt.Println("paas time : " + t1.String())

							cfConfig.Type = "PAAS"
							result, err := utils.GetUaaReFreshToken(reqToken, cfConfig, rdClient)
							//client_test, err := cfclient.NewClient(&cfProvider)
							fmt.Println("paas token : " + result)
							errMessage := model.ErrMessage{"Message": "UnAuthrized"}

							if err != "" {
								utils.RenderJsonUnAuthResponse(errMessage, http.StatusUnauthorized, w)
							} else {
								rdClient.Expire(reqToken, 30*60*time.Second)
								handler.ServeHTTP(w, r)
							}
						} else {
							rdClient.Expire(reqToken, 30*60*time.Second)
							handler.ServeHTTP(w, r)
						}

					} else if strings.Contains(r.RequestURI, "/v2/caas") && val["caasRefreshToken"] != "" { // PaaS 토큰 정보가 있는경우

						rdClient.Expire(reqToken, 30*60*time.Second)
						handler.ServeHTTP(w, r)

					} else if strings.Contains(r.RequestURI, "/v2/saas") { // PaaS 토큰 정보가 있는경우

						rdClient.Expire(reqToken, 30*60*time.Second)
						handler.ServeHTTP(w, r)
					} else {
						fmt.Println("URL Not All")
						//rdClient.Expire(reqToken, 30*60*time.Second)
						//handler.ServeHTTP(w, r)
					}
				}
			}
		} else {
			fmt.Println("url pass ::", r.RequestURI)
			handler.ServeHTTP(w, r)
		}
		//handler.ServeHTTP(w, r)
	}

}

func route(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(f)
}

const STATIC_URL string = "/public/"
const STATIC_ROOT string = "public/"

func StaticHandler(w http.ResponseWriter, req *http.Request) {
	static_file := req.URL.Path[len(STATIC_URL):]
	if len(static_file) != 0 {
		f, err := http.Dir(STATIC_ROOT).Open(static_file)
		if err == nil {
			content := io.ReadSeeker(f)
			http.ServeContent(w, req, static_file, time.Now(), content)
			return
		}
	}
	http.NotFound(w, req)
}
func NewIdentityV3(client *gophercloud.ProviderClient) *gophercloud.ServiceClient {
	v3Endpoint := client.IdentityBase + "v3/"

	return &gophercloud.ServiceClient{
		ProviderClient: client,
		Endpoint:       v3Endpoint,
	}
}

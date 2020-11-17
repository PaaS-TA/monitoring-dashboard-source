package controller

import (
	"fmt"
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
	"kr/paasta/monitoring/common/controller"
	"kr/paasta/monitoring/iaas/model"
	pm "kr/paasta/monitoring/paas/model"
	"kr/paasta/monitoring/routes"
	"kr/paasta/monitoring/utils"
	"net/http"
	"time"
)

func NewHandler(openstack_provider model.OpenstackProvider, iaasInfluxClient client.Client, paasInfluxClient client.Client,
	iaasTxn *gorm.DB, paasTxn *gorm.DB, iaasElasticClient *elastic.Client, paasElasticClient *elastic.Client, monsClient monascaclient.Client,
	auth monascagopher.AuthOptions, databases pm.Databases, rdClient *redis.Client, sysType string, boshClient *gogobosh.Client, cfConfig pm.CFConfig) http.Handler {

	//Controller선언
	var loginController *controller.LoginController
	var memberController *controller.MemberController

	// CaaS Metrics
	var caasMetricsController *MetricController

	loginController = controller.NewLoginController(openstack_provider, monsClient, auth, paasTxn, rdClient, sysType, cfConfig)
	memberController = controller.NewMemberController(openstack_provider, paasTxn, rdClient, sysType, cfConfig)

	var caasActions rata.Handlers
	// add CAAS
	if strings.Contains(sysType, utils.SYS_TYPE_CAAS) || sysType == utils.SYS_TYPE_ALL {
		caasMetricsController = NewMetricControllerr(paasTxn)

		caasActions = rata.Handlers{
			routes.MEMBER_JOIN_CHECK_DUPLICATION_CAAS_ID: route(memberController.MemberJoinCheckDuplicationCaasId),
			routes.MEMBER_JOIN_CHECK_CAAS:                route(memberController.MemberCheckCaaS),
			routes.CAAS_K8S_CLUSTER_AVG:                  route(caasMetricsController.GetClusterAvg),
			routes.CAAS_WORK_NODE_LIST:                   route(caasMetricsController.GetWorkNodeList),
			routes.CAAS_WORK_NODE_INFO:                   route(caasMetricsController.GetWorkNodeInfo),
			routes.CAAS_CONTIANER_LIST:                   route(caasMetricsController.GetContainerList),
			routes.CAAS_CONTIANER_INFO:                   route(caasMetricsController.GetContainerInfo),
			routes.CAAS_CONTIANER_LOG:                    route(caasMetricsController.GetContainerLog),
			routes.CAAS_CLUSTER_OVERVIEW:                 route(caasMetricsController.GetClusterOverView),
			routes.CAAS_WORKLOADS_STATUS:                 route(caasMetricsController.GetWorkloadsStatus),
			routes.CAAS_MASTER_NODE_USAGE:                route(caasMetricsController.GetMasterNodeUsage),
			routes.CAAS_WORK_NODE_AVG:                    route(caasMetricsController.GetWorkNodeAvg),
			routes.CAAS_WORKLOADS_CONTI_SUMMARY:          route(caasMetricsController.GetWorkloadsContiSummary),
			routes.CAAS_WORKLOADS_USAGE:                  route(caasMetricsController.GetWorkloadsUsage),
			routes.CAAS_POD_STAT:                         route(caasMetricsController.GetPodStatList),
			routes.CAAS_POD_LIST:                         route(caasMetricsController.GetPodMetricList),
			routes.CAAS_POD_INFO:                         route(caasMetricsController.GetPodInfo),
			routes.CAAS_WORK_NODE_GRAPH:                  route(caasMetricsController.GetWorkNodeInfoGraph),
			routes.CAAS_WORKLOADS_GRAPH:                  route(caasMetricsController.GetWorkloadsInfoGraph),
			routes.CAAS_POD_GRAPH:                        route(caasMetricsController.GetPodInfoGraph),
			routes.CAAS_CONTIANER_GRAPH:                  route(caasMetricsController.GetContainerInfoGraph),

			routes.CAAS_ALARM_INFO:          route(caasMetricsController.GetAlarmInfo),
			routes.CAAS_ALARM_UPDATE:        route(caasMetricsController.GetAlarmUpdate),
			routes.CAAS_ALARM_LOG:           route(caasMetricsController.GetAlarmLog),
			routes.CAAS_WORK_NODE_GRAPHLIST: route(caasMetricsController.GetWorkNodeInfoGraphList),
			routes.CAAS_ALARM_SNS_INFO:      route(caasMetricsController.GetSnsInfo),
			routes.CAAS_ALARM_COUNT:         route(caasMetricsController.GetAlarmCount),
			routes.CAAS_ALARM_SNS_SAVE:      route(caasMetricsController.GetlarmSnsSave),

			routes.CAAS_ALARM_STATUS_UPDATE:      route(caasMetricsController.UpdateAlarmState),
			routes.CAAS_ALARM_ACTION:             route(caasMetricsController.CreateAlarmResolve),
			routes.CAAS_ALARM_ACTION_DELETE:      route(caasMetricsController.DeleteAlarmResolve),
			routes.CAAS_ALARM_ACTION_UPDATE:      route(caasMetricsController.UpdateAlarmResolve),
			routes.CAAS_ALARM_SNS_CHANNEL_LIST:   route(caasMetricsController.GetAlarmSnsReceiver),
			routes.CAAS_ALARM_SNS_CHANNEL_DELETE: route(caasMetricsController.DeleteAlarmSnsChannel),
			routes.CAAS_ALARM_ACTION_LIST:        route(caasMetricsController.GetAlarmActionList),
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

	if strings.Contains(sysType, utils.SYS_TYPE_CAAS) || sysType == utils.SYS_TYPE_ALL {
		actionlist = append(actionlist, caasActions)
		routeList = append(routeList, routes.CaasRoutes)
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

package utils

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"

	"monitoring-portal/iaas_new/model"
	pm "monitoring-portal/paas/model"
)

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
		if r.RequestURI != "/v2/login" && r.RequestURI != "/v2/logout" && !strings.Contains(r.RequestURI, "/v2/member/join") && r.RequestURI != "/v2/ping" && r.RequestURI != "/" && !strings.Contains(r.RequestURI, "/public/") && !strings.Contains(r.RequestURI, "/v2/paas/app/") && !strings.Contains(r.RequestURI, "/v2/caas/monitoring/podList") && !strings.Contains(r.RequestURI, "/v2/paas/diagram") {
			if !strings.Contains(r.RequestURI, "favicon.ico") {
				Logger.Info("Request URI :: ", r.RequestURI)
			}


			reqToken := r.Header.Get(model.CSRF_TOKEN_NAME)
			if reqToken == "0" || reqToken == "null" {
				Logger.Info("HttpWrap Hander reqToken is null ")
				errMessage := model.ErrMessage{"Message": "UnAuthrized"}
				RenderJsonUnAuthResponse(errMessage, http.StatusUnauthorized, w)
			} else {
				//fmt.Println("HttpWrap Hander reqToken =",len(reqToken),":",reqToken)
				//모든 경로의 redis 의 토큰 정보를 확인한다

				Logger.Debugf("reqToken : %v\n", reqToken)
				val := rdClient.HGetAll(reqToken).Val()
				Logger.Debugf("iaasToken : %v\n", val["iaasToken"])

				if val == nil || len(val) == 0 { // redis 에서 token 정보가 expire 된경우 로그인 화면으로 돌아간다
					Logger.Info("HttpWrap Hander redis.iaas_userid is null ")
					errMessage := model.ErrMessage{"Message": "UnAuthrized"}
					RenderJsonUnAuthResponse(errMessage, http.StatusUnauthorized, w)
				} else {

					if strings.Contains(r.RequestURI, "/v2/member") && val["userId"] != "" {
						handler.ServeHTTP(w, r)

					} else if strings.Contains(r.RequestURI, "/v2/iaas") && val["iaasToken"] != "" && val["iaasUserId"] != "" { // IaaS 토큰 정보가 있는경우

						provider1, _, err := GetOpenstackProvider(r)
						if err != nil || provider1 == nil {
							Logger.Debug(err)
							errMessage := model.ErrMessage{"Message": "UnAuthrized"}
							RenderJsonUnAuthResponse(errMessage, http.StatusUnauthorized, w)
						} else {
							v3Client := NewIdentityV3(provider1)

							//IaaS, token 검증
							bool, err := tokens.Validate(v3Client, val["iaasToken"])
							if err != nil || bool == false {
								//errMessage := model.ErrMessage{"Message": "UnAuthrized"}
								//utils.RenderJsonUnAuthResponse(errMessage, http.StatusUnauthorized, w)
								Logger.Info("iaas token validate error::", err)
								handler.ServeHTTP(w, r)
							} else {
								//두개 token 이 없는 경우도 고려 해야함
								rdClient.Expire(reqToken, 30*60*time.Second)
								handler.ServeHTTP(w, r)
							}
						}
					} else if strings.Contains(r.RequestURI, "/v2/paas") && val["paasRefreshToken"] != "" { // PaaS 토큰 정보가 있는경우

						// Pass token 검증 로직 추가
						// get paas token
						//cfProvider.Token = val["paasToken"]
						t1, _ := time.Parse(time.RFC3339, val["paasExpire"])
						if t1.Before(time.Now()) {
							Logger.Info("paas time : " + t1.String())

							cfConfig.Type = "PAAS"
							result, err := GetUaaReFreshToken(reqToken, cfConfig, rdClient)
							//client_test, err := cfclient.NewClient(&cfProvider)
							Logger.Info("paas token : " + result)
							errMessage := model.ErrMessage{"Message": "UnAuthrized"}

							if err != "" {
								RenderJsonUnAuthResponse(errMessage, http.StatusUnauthorized, w)
							} else {
								//_, err01 := client_test.GetToken() // cf token 을 refresh 함
								//if err01 != nil {
								//	utils.RenderJsonUnAuthResponse(errMessage, http.StatusUnauthorized, w)
								//	return
								//}
								/*
									fmt.Println("paas hander token ::: ",token)

									token01, err02 := client_test.ListApps()
									if err02 != nil {
										fmt.Println("paas ListApps error::",token01,":::",err02.Error())
									}else{
										fmt.Println("paas ListApps info  ::",token01)
									}
								*/
								rdClient.Expire(reqToken, 30*60*time.Second)
								handler.ServeHTTP(w, r)
							}
						} else {
							rdClient.Expire(reqToken, 30*60*time.Second)
							handler.ServeHTTP(w, r)
						}

					} else if strings.Contains(r.RequestURI, "/v2/caas") && val["caasRefreshToken"] != "" { // PaaS 토큰 정보가 있는경우

						// Pass token 검증 로직 추가
						//get paas token
						//cfProvider.Token = val["paasToken"]
						//t1, _ := time.Parse(time.RFC3339, val["caasExpire"])
						//if t1.Before(time.Now()) {
						//	fmt.Println("caas time : " + t1.String())
						//
						//	cfConfig.Type = "CAAS"
						//	result, err := utils.GetUaaReFreshToken(reqToken, cfConfig, rdClient)
						//	//client_test, err := cfclient.NewClient(&cfProvider)
						//	fmt.Println("caas token : " + result)
						//	errMessage := model.ErrMessage{"Message": "UnAuthrized"}
						//	if err != "" {
						//		utils.RenderJsonUnAuthResponse(errMessage, http.StatusUnauthorized, w)
						//	} else {
						//		//_, err01 := client_test.GetToken() // cf token 을 refresh 함
						//		//if err01 != nil {
						//		//	utils.RenderJsonUnAuthResponse(errMessage, http.StatusUnauthorized, w)
						//		//	return
						//		//}
						//		/*
						//			fmt.Println("paas hander token ::: ",token)
						//
						//			token01, err02 := client_test.ListApps()
						//			if err02 != nil {
						//				fmt.Println("paas ListApps error::",token01,":::",err02.Error())
						//			}else{
						//				fmt.Println("paas ListApps info  ::",token01)
						//			}
						//		*/
						//		rdClient.Expire(reqToken, 30*60*time.Second)
						//		handler.ServeHTTP(w, r)
						//	}
						//}else{
						//	rdClient.Expire(reqToken, 30*60*time.Second)
						//	handler.ServeHTTP(w, r)
						//}
						rdClient.Expire(reqToken, 30*60*time.Second)
						handler.ServeHTTP(w, r)

					} else if strings.Contains(r.RequestURI, "/v2/saas") { // PaaS 토큰 정보가 있는경우

						rdClient.Expire(reqToken, 30*60*time.Second)
						handler.ServeHTTP(w, r)
					} else {
						//fmt.Println("URL Not All")
						//rdClient.Expire(reqToken, 30*60*time.Second)
						//handler.ServeHTTP(w, r)
					}
				}
			}
		} else {
			Logger.Info("url pass ::", r.RequestURI)
			handler.ServeHTTP(w, r)
		}
		//handler.ServeHTTP(w, r)
	}

}
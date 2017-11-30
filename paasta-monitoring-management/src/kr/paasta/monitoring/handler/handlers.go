package handler

import (
	"io"
	"net/http"
	"time"
	"github.com/tedsuo/rata"
	"kr/paasta/monitoring/router"
	"kr/paasta/monitoring/controller"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	client "github.com/influxdata/influxdb/client/v2"
	"kr/paasta/monitoring/domain"
)

type Context struct {
	Title  string
	Static string
}

func NewHandler(txn *gorm.DB, influxClient client.Client, databases domain.Databases) http.Handler {

	alarmController := controller.GetAlarmController(txn)
	alarmPolicyController := controller.GetAlarmPolicyController(txn)
	containerController := controller.GetContainerController(txn, influxClient)
	metricsController := controller.GetMetricsController(influxClient, databases)

	actions := rata.Handlers{

		//Controller선언
		router.Main: route(alarmController.Main),

		router.AlarmPolicyList: route(alarmPolicyController.GetAlarmPolicyList),
		router.AlarmPolicyUpdate: route(alarmPolicyController.UpdateAlarmPolicyList),

		router.AlarmList: route(alarmController.GetAlarmList),
		router.AlarmResolveStatus: route(alarmController.GetAlarmResolveStatus),
		router.AlarmDetail: route(alarmController.GetAlarmDetail),
		router.UpdateAlarm: route(alarmController.UpdateAlarm),
		router.CreateAlarmAction: route(alarmController.CreateAlarmAction),
		router.UpdateAlarmAction: route(alarmController.UpdateAlarmAction),
		router.DeleteAlarmAction: route(alarmController.DeleteAlarmAction),

		router.GetAlarmStat: route(alarmController.GetAlarmStat),
		router.GetContainerDeploy: route(containerController.GetContainerDeploy),

		//Application Resources 조회 (2017-08-14 추가)
		//Application cpu, memory, disk usage 정보 조회
		router.GetAppResource: route(metricsController.GetApplicationResources),
		router.GetAppResourceAll: route(metricsController.GetApplicationResourcesAll),
		//Application cpu variation 정보 조회
		router.GetAppCpuUsage: route(metricsController.GetAppCpuUsage),
		//Application memory variation 정보 조회
		router.GetAppMemoryUsage: route(metricsController.GetAppMemoryUsage),
		//Application disk variation 정보 조회
		router.GetDiskUsage: route(metricsController.GetDiskUsage),

		//Application network variation 정보 조회
		router.GetAppNetworkIoByte: route(metricsController.GetAppNetworkIoKByte),
		// influxDB에서 조회
		router.GetDiskIOList: route(metricsController.GetDiskIOList),
		router.GetNetworkIOList: route(metricsController.GetNetworkIOList),
		router.GetTopProcessList: route(metricsController.GetTopProcessList),


		// Html
		router.Static:    route(StaticHandler),


	}

	handler, err := rata.NewRouter(router.Routes, actions)
	if err != nil {
		panic("unable to create router: " + err.Error())
	}

	return HttpWrap(handler)
}

func HttpWrap(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//log.Println("HttpWrap Called")
		handler.ServeHTTP(w, r)
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
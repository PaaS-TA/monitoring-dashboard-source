package controller

import (
	"encoding/json"
	"fmt"
	"github.com/thoas/go-funk"
	"github.com/tidwall/gjson"
	"kr/paasta/monitoring/saas/model"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	//"fmt"ApplicationStat
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"kr/paasta/monitoring/paas/util"
	"kr/paasta/monitoring/saas/service"
	"log"
	"net/http"
)

type SaasController struct {
	txn         *gorm.DB
	pinpointUrl string
	config      map[string]string
}

type ApplicationList struct {
	PinpointUrl string            `json:"pinpointUrl"`
	Data        []ApplicationStat `json:"data"`
}

func NewSaasController(txn *gorm.DB) *SaasController {
	config, _ := util.ReadConfig(`config.ini`)
	return &SaasController{
		txn:    txn,
		config: config,
	}
}

type ApplicationStat struct {
	AppName        string  `json:"appName"`
	AgentId        string  `json:"agentId"`
	ServiceType    string  `json:"serviceType"`
	JvmCpuRate     float64 `json:"jvmCpuRate"`
	SystemCpuRate  float64 `json:"systemCpuRate"`
	HaepMemory     float64 `json:"haepMemory"`
	NoneHeapMemory float64 `json:"noneHeapMemory"`
	ActiveThread   float64 `json:"activeThread"`
	ResponseTime   float64 `json:"responseTime"`
}

type ApplicationGaugeTot struct {
	AgentTotCnt        int64   `json:"agentTotCnt"`
	AgentUseCnt        int64   `json:"agentUseCnt"`
	SystemCpuRate      float64 `json:"systemCpuRate"`
	HeapMemoryRate     float64 `json:"heapMemoryRate"`
	HaepMaxMemory      float64 `json:"haepMaxMemory"`
	HaepMemory         float64 `json:"haepMemory"`
	NoneHeapMemoryRate float64 `json:"noneHeapMemoryRate"`
	NoneHeapMaxMemory  float64 `json:"noneHeapMaxMemory"`
	NoneHeapMemory     float64 `json:"noneHeapMemory"`
}

type AgentStatus struct {
	AgentCnt   int `json:"AgentCnt"`
	Running    int `json:"Running"`
	Disconnect int `json:"Disconnect"`
	Shutdown   int `json:"Shutdown"`
}

func GetSaasController(txn *gorm.DB) *SaasController {
	config, err := util.ReadConfig(`config.ini`)
	pinpointUrl := config["saas.pinpoint.url"]
	if err != nil {
		log.Println(err)
	}
	return &SaasController{
		txn:         txn,
		pinpointUrl: pinpointUrl,
	}
}

func (p *SaasController) GetApplicationList(w http.ResponseWriter, r *http.Request) {
	applications, data, pinpointUrl := appNameList(p.config)

	pAppName, _ := r.URL.Query()["appName"]
	pSortKey, _ := r.URL.Query()["sortKey"]
	pSortVal, _ := r.URL.Query()["sortVal"]

	resutChan := make(chan ApplicationStat, getMakeChannelCount(applications, pAppName, string(data)))

	var waitGroup sync.WaitGroup
	var agentCount int64 = 0

	for appName, _ := range applications {
		agentIds := gjson.Get(string(data), appName+".#.agentId")
		serviceType := gjson.Get(string(data), appName+".#.serviceType").Array()[0]

		if pAppName != nil && len(pAppName[0]) > 0 {
			if strings.Contains(appName, pAppName[0]) {
				for _, agentId := range agentIds.Array() {
					waitGroup.Add(1)
					go goroutinAppList(resutChan, &waitGroup, appName, agentId.String(), pinpointUrl, &agentCount, serviceType.String())
				}
			}
		} else {
			for _, agentId := range agentIds.Array() {
				waitGroup.Add(1)
				go goroutinAppList(resutChan, &waitGroup, appName, agentId.String(), pinpointUrl, &agentCount, serviceType.String())
			}
		}
	}

	waitGroup.Wait()

	var k int64 = 0
	applicationStats := make([]ApplicationStat, agentCount, agentCount)
	for k = 0; k < agentCount; k++ {
		applicationStats[k] = <-resutChan
	}

	close(resutChan)

	//pSortKey
	//pSortVal
	sort.Slice(applicationStats, func(i, j int) bool {
		if pSortKey != nil && len(pSortKey) > 0 && pSortVal != nil && len(pSortVal) > 0 {
			if pSortKey[0] == "appName" && pSortVal[0] == "A" {
				return applicationStats[i].AppName < applicationStats[j].AppName
			} else if pSortKey[0] == "appName" && pSortVal[0] == "D" {
				return applicationStats[i].AppName > applicationStats[j].AppName
			}

			if pSortKey[0] == "jvmCpuRate" && pSortVal[0] == "A" {
				return applicationStats[i].AppName < applicationStats[j].AppName
			} else if pSortKey[0] == "jvmCpuRate" && pSortVal[0] == "D" {
				return applicationStats[i].AppName > applicationStats[j].AppName
			}

			if pSortKey[0] == "systemCpuRate" && pSortVal[0] == "A" {
				return applicationStats[i].AppName < applicationStats[j].AppName
			} else if pSortKey[0] == "systemCpuRate" && pSortVal[0] == "D" {
				return applicationStats[i].AppName > applicationStats[j].AppName
			}

			if pSortKey[0] == "haepMemory" && pSortVal[0] == "A" {
				return applicationStats[i].AppName < applicationStats[j].AppName
			} else if pSortKey[0] == "haepMemory" && pSortVal[0] == "D" {
				return applicationStats[i].AppName > applicationStats[j].AppName
			}

			if pSortKey[0] == "noneHeapMemory" && pSortVal[0] == "A" {
				return applicationStats[i].AppName < applicationStats[j].AppName
			} else if pSortKey[0] == "noneHeapMemory" && pSortVal[0] == "D" {
				return applicationStats[i].AppName > applicationStats[j].AppName
			}

			if pSortKey[0] == "activeThread" && pSortVal[0] == "A" {
				return applicationStats[i].AppName < applicationStats[j].AppName
			} else if pSortKey[0] == "activeThread" && pSortVal[0] == "D" {
				return applicationStats[i].AppName > applicationStats[j].AppName
			}

			if pSortKey[0] == "responseTime" && pSortVal[0] == "A" {
				return applicationStats[i].AppName < applicationStats[j].AppName
			} else if pSortKey[0] == "responseTime" && pSortVal[0] == "D" {
				return applicationStats[i].AppName > applicationStats[j].AppName
			}
		}
		return applicationStats[i].AppName < applicationStats[j].AppName
	})

	util.RenderJsonResponse(ApplicationList{PinpointUrl: pinpointUrl, Data: applicationStats}, w)
}

func getMakeChannelCount(applications map[string]string, pAppName []string, data string) int {
	count := 0
	for appName, _ := range applications {
		agentIds := gjson.Get(data, appName+".#.agentId")
		if pAppName != nil && len(pAppName[0]) > 0 {
			if strings.Contains(appName, pAppName[0]) {
				count += len(agentIds.Array())
			}
		} else {
			count += len(agentIds.Array())
		}
	}
	return count
}

func (p *SaasController) GetAgentStatus(w http.ResponseWriter, r *http.Request) {
	applications, data, _ := appNameList(p.config)

	agentStatus := AgentStatus{
		AgentCnt:   0,
		Running:    0,
		Disconnect: 0,
		Shutdown:   0,
	}

	for appName, _ := range applications {
		jpath := appName + ".#.status.state.code"
		json := gjson.Get(string(data), jpath)

		funk.ForEach(json.Array(), func(x gjson.Result) {
			if x.Num == 100 {
				agentStatus.Running += 1
			} else if x.Num == 300 {
				agentStatus.Disconnect += 1
			} else if x.Num == 200 || x.Num == 201 {
				agentStatus.Shutdown += 1
			}
		})
	}

	agentStatus.AgentCnt = len(applications)

	util.RenderJsonResponse(agentStatus, w)
}

func (p *SaasController) GetAgentGaugeTot(w http.ResponseWriter, r *http.Request) {
	applications, data, pinpointUrl := appNameList(p.config)

	resutChan := make(chan ApplicationGaugeTot, getMakeChannelCount(applications, nil, string(data)))

	var waitGroup sync.WaitGroup
	var agentCount int64 = 0
	var agentUseCnt int64 = 0

	for appName, _ := range applications {
		agentIds := gjson.Get(string(data), appName+".#.agentId")
		jpath := appName + ".#.status.state.code"
		json := gjson.Get(string(data), jpath)

		funk.ForEach(json.Array(), func(x gjson.Result) {
			if x.Num == 100 {
				agentUseCnt += 1
			}
		})

		for _, agentId := range agentIds.Array() {
			waitGroup.Add(1)
			go goroutinAppGaugeTot(resutChan, &waitGroup, appName, agentId.String(), pinpointUrl, &agentCount)
		}
	}

	waitGroup.Wait()

	var k int64 = 0
	var sumSystemCpuRate float64 = 0
	var sumHaepMaxMemory float64 = 0
	var sumHaepMemory float64 = 0
	var sumNoneHeapMaxMemory float64 = 0
	var sumNoneHeapMemory float64 = 0

	var avgSystemCpuRate float64 = 0
	var avgHaepMemoryRate float64 = 0
	var avgNoneHaepMemoryRate float64 = 0

	for k = 0; k < agentCount; k++ {
		applicationGaugeTot := <-resutChan
		sumSystemCpuRate += applicationGaugeTot.SystemCpuRate
		sumHaepMaxMemory += applicationGaugeTot.HaepMaxMemory
		sumHaepMemory += applicationGaugeTot.HaepMemory
		sumNoneHeapMaxMemory += applicationGaugeTot.NoneHeapMaxMemory
		sumNoneHeapMemory += applicationGaugeTot.NoneHeapMemory
	}
	close(resutChan)

	if agentCount > 0 {
		avgSystemCpuRate = sumSystemCpuRate / float64(agentCount)
		avgHaepMemoryRate = sumHaepMemory / sumHaepMaxMemory * 100
		avgNoneHaepMemoryRate = sumNoneHeapMemory / sumNoneHeapMaxMemory * 100

	}

	avgSystemCpuRate, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", avgSystemCpuRate), 0)
	sumHaepMaxMemory, _ = strconv.ParseFloat(fmt.Sprintf("%.0f", sumHaepMaxMemory), 0)
	sumHaepMemory, _ = strconv.ParseFloat(fmt.Sprintf("%.0f", sumHaepMemory), 0)
	sumNoneHeapMaxMemory, _ = strconv.ParseFloat(fmt.Sprintf("%.0f", sumNoneHeapMaxMemory), 0)
	sumNoneHeapMemory, _ = strconv.ParseFloat(fmt.Sprintf("%.0f", sumNoneHeapMemory), 0)
	avgHaepMemoryRate, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", avgHaepMemoryRate), 0)
	avgNoneHaepMemoryRate, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", avgNoneHaepMemoryRate), 0)

	resultData := ApplicationGaugeTot{
		AgentTotCnt:        agentCount,
		AgentUseCnt:        agentUseCnt,
		SystemCpuRate:      avgSystemCpuRate,
		HeapMemoryRate:     avgHaepMemoryRate,
		NoneHeapMemoryRate: avgNoneHaepMemoryRate,
		HaepMaxMemory:      sumHaepMaxMemory,
		HaepMemory:         sumHaepMemory,
		NoneHeapMaxMemory:  sumNoneHeapMaxMemory,
		NoneHeapMemory:     sumNoneHeapMemory,
	}

	util.RenderJsonResponse(resultData, w)
}

func appNameList(config map[string]string) (map[string]string, []byte, string) {
	pinpointUrl := config["saas.pinpoint.url"]
	url := pinpointUrl + "/getAgentList.pinpoint"

	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	var applications map[string]string

	json.Unmarshal(data, &applications)

	return applications, data, pinpointUrl
}

func goroutinAppList(c chan ApplicationStat, waitGroup *sync.WaitGroup, appName string, agentId string, pinpointUrl string, agentCount *int64, serviceType string) {
	defer waitGroup.Done()

	from := strconv.FormatInt(time.Now().Add(-60*time.Second).UTC().Unix(), 10) + "000"
	to := strconv.FormatInt(time.Now().UTC().Unix(), 10) + "000"

	// =================================================================================================
	// CPU JVM / System
	// =================================================================================================
	data, err := getRestCall(pinpointUrl + "/getAgentStat/cpuLoad/chart.pinpoint?agentId=" + agentId + "&from=" + from + "&to=" + to + "")
	if err != nil {
		return
	}

	jpath := "charts.y.CPU_LOAD_JVM.#.2"
	jvmCpuRate := getAvgReustData(data, jpath, 1)

	jpath = "charts.y.CPU_LOAD_SYSTEM.#.2"
	systemCpuRate := getAvgReustData(data, jpath, 1)

	// =================================================================================================
	// VM GC
	// =================================================================================================
	data, err = getRestCall(pinpointUrl + "/getAgentStat/jvmGc/chart.pinpoint?agentId=" + agentId + "&from=" + from + "&to=" + to + "")
	if err != nil {
		return
	}

	jpath = "charts.y.JVM_MEMORY_HEAP_USED.#.2"
	haepMemory := getAvgReustData(data, jpath, 1024*1024)

	jpath = "charts.y.JVM_MEMORY_NON_HEAP_USED.#.2"
	noneHeapMemory := getAvgReustData(data, jpath, 1024*1024)

	// =================================================================================================
	// Active Thread
	// =================================================================================================
	data, err = getRestCall(pinpointUrl + "/getAgentStat/activeTrace/chart.pinpoint?agentId=" + agentId + "&from=" + from + "&to=" + to + "")
	if err != nil {
		return
	}

	jpath = "charts.y.ACTIVE_TRACE_FAST.#.3"
	activeThread := getAvgReustData(data, jpath, 1)

	// =================================================================================================
	// Response  Time
	// =================================================================================================
	data, err = getRestCall(pinpointUrl + "/getAgentStat/responseTime/chart.pinpoint?agentId=" + agentId + "&from=" + from + "&to=" + to + "")
	if err != nil {
		return
	}

	jpath = "charts.y.AVG.#.2"
	responseTime := getAvgReustData(data, jpath, 1)

	applicationStat := ApplicationStat{
		AppName:        appName,
		AgentId:        agentId,
		ServiceType:    serviceType,
		JvmCpuRate:     jvmCpuRate,
		SystemCpuRate:  systemCpuRate,
		HaepMemory:     haepMemory,
		NoneHeapMemory: noneHeapMemory,
		ActiveThread:   activeThread,
		ResponseTime:   responseTime,
	}
	c <- applicationStat
	atomic.AddInt64(agentCount, 1)

}

func goroutinAppGaugeTot(c chan ApplicationGaugeTot, waitGroup *sync.WaitGroup, appName string, agentId string, pinpointUrl string, agentCount *int64) {
	defer waitGroup.Done()

	from := strconv.FormatInt(time.Now().Add(-600*time.Second).UTC().Unix(), 10) + "000"
	to := strconv.FormatInt(time.Now().UTC().Unix(), 10) + "000"

	// =================================================================================================
	// System CPU
	// =================================================================================================
	systemCpuRate := getRestCallAvgReustData(pinpointUrl+"/getAgentStat/cpuLoad/chart.pinpoint?agentId="+agentId+"&from="+from+"&to="+to+"", "charts.y.CPU_LOAD_SYSTEM.#.2", 1)

	// =================================================================================================
	// VM GC
	// =================================================================================================
	data, err := getRestCall(pinpointUrl + "/getAgentStat/jvmGc/chart.pinpoint?agentId=" + agentId + "&from=" + from + "&to=" + to + "")
	if err != nil {
		return
	}

	jpath := "charts.y.JVM_MEMORY_HEAP_MAX.#.2"
	haepMaxMemory := getAvgReustData(data, jpath, 1024*1024)

	jpath = "charts.y.JVM_MEMORY_HEAP_USED.#.2"
	haepMemory := getAvgReustData(data, jpath, 1024*1024)

	jpath = "charts.y.JVM_MEMORY_NON_HEAP_MAX.#.2"
	noneHeapMaxMemory := getAvgReustData(data, jpath, 1024*1024)

	jpath = "charts.y.JVM_MEMORY_NON_HEAP_USED.#.2"
	noneHeapMemory := getAvgReustData(data, jpath, 1024*1024)

	applicationSGaugeTot := ApplicationGaugeTot{
		SystemCpuRate:     systemCpuRate,
		HaepMaxMemory:     haepMaxMemory,
		HaepMemory:        haepMemory,
		NoneHeapMaxMemory: noneHeapMaxMemory,
		NoneHeapMemory:    noneHeapMemory,
	}
	c <- applicationSGaugeTot
	atomic.AddInt64(agentCount, 1)
}

/**
Rest Api call
*/
func getRestCall(url string) (string, error) {
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	return string(data), err
}

func getAvgReustData(data string, jpath string, unit float64) float64 {
	json := gjson.Get(data, jpath)

	mapData := funk.Map(json.Array(), func(x gjson.Result) float64 {
		if jpath == "charts.y.ACTIVE_TRACE_NORMAL.#.4" {
			num, _ := strconv.ParseFloat(strings.Replace(x.String(), "s", "", 1), 0)
			return num
		}
		return x.Num
	})

	filterData := funk.Filter(mapData, func(x float64) bool {
		return x > 0.0
	}).([]float64)

	var resultData float64
	if len(filterData) > 0 {
		resultData = funk.SumFloat64(filterData) / float64(len(filterData)) / unit
		resultData, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", resultData), 0)
	}
	return resultData
}

func (p *SaasController) RemoveApplication(w http.ResponseWriter, r *http.Request) {
	pinpointWasUrl := p.config["saas.pinpointWas.url"]

	fmt.Println("pinpointWasUrl : " + pinpointWasUrl)

	applicationName := r.URL.Query().Get("applicationName")
	agentId := r.URL.Query().Get("agentId")

	result, err := getRestCall(pinpointWasUrl + "/admin/removeAgentId.pinpoint?applicationName=" + applicationName + "&agentId=" + agentId + "&password=admin")

	if err != nil {
		util.RenderJsonResponse(err, w)
		return
	}
	util.RenderJsonResponse(result, w)
}

func getRestCallAvgReustData(url string, jpath string, unit float64) float64 {
	body, err := getRestCall(url)

	if err != nil {
		return -1
	}

	return getAvgReustData(body, jpath, unit)
}

/**
Alarm Api call
*/
func (p *SaasController) GetAlarmInfo(w http.ResponseWriter, r *http.Request) {
	//service호출
	result, err := service.GetAlarmService(p.txn).GetAlarmInfo()
	if err != nil {
		util.RenderJsonResponse(err, w)
		return
	}
	util.RenderJsonResponse(result, w)
}

func (p *SaasController) GetAlarmUpdate(w http.ResponseWriter, r *http.Request) {
	var apiRequest []model.AlarmPolicyRequest
	data, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	err := json.Unmarshal(data, &apiRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	service.GetAlarmService(p.txn).GetAlarmUpdate(apiRequest)
}

func (p *SaasController) GetAlarmLog(w http.ResponseWriter, r *http.Request) {
	searchDateFrom := r.URL.Query().Get("searchDateFrom")
	searchDateTo := r.URL.Query().Get("searchDateTo")

	alarmType := r.URL.Query().Get("alarmType")
	alarmStatus := r.URL.Query().Get("alarmStatus")
	resolveStatus := r.URL.Query().Get("resolveStatus")

	result, err := service.GetAlarmService(p.txn).GetAlarmLog(searchDateFrom, searchDateTo, alarmType, alarmStatus, resolveStatus)
	if err != nil {
		util.RenderJsonResponse(err, w)
		return
	}
	util.RenderJsonResponse(result, w)
}

func (p *SaasController) GetSnsInfo(w http.ResponseWriter, r *http.Request) {
	result, err := service.GetAlarmService(p.txn).GetSnsInfo()
	if err != nil {
		util.RenderJsonResponse(err, w)
		return
	}
	util.RenderJsonResponse(result, w)
}

func (p *SaasController) GetAlarmCount(w http.ResponseWriter, r *http.Request) {
	result, err := service.GetAlarmService(p.txn).GetAlarmCount()

	if err != nil {
		util.RenderJsonResponse(err, w)
		return
	}
	util.RenderJsonResponse(result, w)
}

func (p *SaasController) GetlarmSnsSave(w http.ResponseWriter, r *http.Request) {
	var alarmSns model.BatchAlarmSnsRequest
	err := json.NewDecoder(r.Body).Decode(&alarmSns)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	error := service.GetAlarmService(p.txn).GetlarmSnsSave(alarmSns)

	util.RenderJsonResponse(error, w)
}

func (h *SaasController) UpdateAlarmState(w http.ResponseWriter, r *http.Request) {

	var alarmrRsolveRequest model.AlarmrRsolveRequest
	err := json.NewDecoder(r.Body).Decode(&alarmrRsolveRequest)
	defer r.Body.Close()

	id, _ := strconv.Atoi(r.FormValue(":id"))
	alarmrRsolveRequest.Id = uint64(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	error := service.GetAlarmService(h.txn).UpdateAlarmSate(alarmrRsolveRequest)
	util.RenderJsonResponse(error, w)
}

func (h *SaasController) CreateAlarmResolve(w http.ResponseWriter, r *http.Request) {
	var alarmrRsolveRequest model.AlarmrRsolveRequest
	err := json.NewDecoder(r.Body).Decode(&alarmrRsolveRequest)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	error := service.GetAlarmService(h.txn).CreateAlarmResolve(alarmrRsolveRequest)
	util.RenderJsonResponse(error, w)
}

func (h *SaasController) UpdateAlarmResolve(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue(":id"))

	var alarmrRsolveRequest model.AlarmrRsolveRequest
	err := json.NewDecoder(r.Body).Decode(&alarmrRsolveRequest)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	alarmrRsolveRequest.Id = uint64(id)
	error := service.GetAlarmService(h.txn).UpdateAlarmResolve(alarmrRsolveRequest)
	util.RenderJsonResponse(error, w)
	return
}

func (h *SaasController) DeleteAlarmResolve(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue(":id"))

	error := service.GetAlarmService(h.txn).DeleteAlarmResolve(uint64(id))
	util.RenderJsonResponse(error, w)
	return
}

func (h *SaasController) GetAlarmSnsReceiver(w http.ResponseWriter, r *http.Request) {
	alarmReceiver, _ := service.GetAlarmService(h.txn).GetAlarmSnsReceiver()
	util.RenderJsonResponse(alarmReceiver, w)
}

func (h *SaasController) DeleteAlarmSnsChannel(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue(":id"))

	err := service.GetAlarmService(h.txn).DeleteAlarmSnsChannel(id)
	util.RenderJsonResponse(err, w)
}

func (h *SaasController) GetAlarmActionList(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(r.FormValue(":id"))

	result, err := service.GetAlarmService(h.txn).GetAlarmActionList(id)

	if err != nil {
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(result, w)
	}
}

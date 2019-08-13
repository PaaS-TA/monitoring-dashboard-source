package controller

import (
	"encoding/json"
	"fmt"
	"github.com/thoas/go-funk"
	"github.com/tidwall/gjson"
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
	AgentTotCnt       int64   `json:"agentTotCnt"`
	AgentUseCnt       int64   `json:"agentUseCnt"`
	SystemCpuRate     float64 `json:"systemCpuRate"`
	HaepMaxMemory     float64 `json:"haepMaxMemory"`
	HaepMemory        float64 `json:"haepMemory"`
	NoneHeapMaxMemory float64 `json:"noneHeapMaxMemory"`
	NoneHeapMemory    float64 `json:"noneHeapMemory"`
}

type AgentStatus struct {
	AgentCnt   int `json:"agentCnt"`
	Running    int `json:"running"`
	Disconnect int `json:"disconnect"`
	Shutdown   int `json:"shutdown"`
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
	applications, data, pinpointUrl := appNameList()

	resutChan := make(chan ApplicationStat, len(applications)*3)

	var waitGroup sync.WaitGroup
	var agentCount int64 = 0

	pAppName, _ := r.URL.Query()["appName"]
	pSortKey, _ := r.URL.Query()["sortKey"]
	pSortVal, _ := r.URL.Query()["sortVal"]

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

	util.RenderJsonResponse(applicationStats, w)
}

func (p *SaasController) GetAgentStatus(w http.ResponseWriter, r *http.Request) {
	applications, data, _ := appNameList()

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
	applications, data, pinpointUrl := appNameList()

	resutChan := make(chan ApplicationGaugeTot, len(applications)*3)

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

	//applicationGaugeTot := make([]ApplicationGaugeTot, agentCount, agentCount)

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
	}

	avgSystemCpuRate, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", avgSystemCpuRate), 0)
	sumHaepMaxMemory, _ = strconv.ParseFloat(fmt.Sprintf("%.0f", sumHaepMaxMemory), 0)
	sumHaepMemory, _ = strconv.ParseFloat(fmt.Sprintf("%.0f", sumHaepMemory), 0)
	sumNoneHeapMaxMemory, _ = strconv.ParseFloat(fmt.Sprintf("%.0f", sumNoneHeapMaxMemory), 0)
	sumNoneHeapMemory, _ = strconv.ParseFloat(fmt.Sprintf("%.0f", sumNoneHeapMemory), 0)

	resultData := ApplicationGaugeTot{
		AgentTotCnt:       agentCount,
		AgentUseCnt:       agentUseCnt,
		SystemCpuRate:     avgSystemCpuRate,
		HaepMaxMemory:     sumHaepMaxMemory,
		HaepMemory:        sumHaepMemory,
		NoneHeapMaxMemory: sumNoneHeapMaxMemory,
		NoneHeapMemory:    sumNoneHeapMemory,
	}

	util.RenderJsonResponse(resultData, w)
}

func appNameList() (map[string]string, []byte, string) {
	config, _ := util.ReadConfig(`config.ini`)

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

	from := strconv.FormatInt(time.Now().Add(-600*time.Second).UTC().Unix(), 10) + "000"
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

	jpath = "charts.y.ACTIVE_TRACE_NORMAL.#.4"
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
	result, err := service.GetAlarmService().GetAlarmInfo()
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

func (p *SaasController) GetAlarmUpdate(w http.ResponseWriter, r *http.Request) {
	//service호출
	//result, err := service.GetAlarmService().GetAlarmUpdate(r)
	service.GetAlarmService().GetAlarmUpdate(r)
	//if err != nil {
	//	util.RenderJsonResponse(err, w)
	//}
	//util.RenderJsonResponse(result, w)
}

func (p *SaasController) GetAlarmLog(w http.ResponseWriter, r *http.Request) {
	//service호출
	result, err := service.GetAlarmService().GetAlarmLog()
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(result, w)
}

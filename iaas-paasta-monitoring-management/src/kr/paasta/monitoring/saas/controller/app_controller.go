package controller

import (
	"encoding/json"
	"fmt"
	"github.com/thoas/go-funk"
	"github.com/tidwall/gjson"
	_ "github.com/tidwall/sjson"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	//"fmt"ApplicationStat
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"kr/paasta/monitoring/paas/util"
	"log"
	"net/http"
)

type SaasController struct {
	txn *gorm.DB
	//config map[string]string

}

type ApplicationStat struct {
	AppName        string  `json:"appName"`
	AgentId        string  `json:"agentId"`
	JvmCpuRate     float64 `json:"jvmCpuRate"`
	SystemCpuRate  float64 `json:"systemCpuRate"`
	HaepMemory     float64 `json:"haepMemory"`
	NoneHeapMemory float64 `json:"noneHeapMemory"`
	ActiveThread   float64 `json:"activeThread"`
	ResponseTime   float64 `json:"responseTime"`
}

func GetSaasController(txn *gorm.DB) *SaasController {
	//config, err := util.ReadConfig(`config.ini`)

	//if err != nil {
	//	log.Println(err)
	//}
	return &SaasController{
		txn: txn,
		//config: config,
	}
}

func (p *SaasController) GetApplicationList(w http.ResponseWriter, r *http.Request) {
	config, err := util.ReadConfig(`config.ini`)
	pinpointUrl, _ := config["saas.pinpoint.url"]
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

	resutChan := make(chan ApplicationStat, len(applications)*3)

	var waitGroup sync.WaitGroup
	var agentCount int64 = 0

	for appName, _ := range applications {
		agentIds := gjson.Get(string(data), appName+".#.agentId")

		for _, agentId := range agentIds.Array() {
			waitGroup.Add(1)
			go goroutinAppList(resutChan, &waitGroup, appName, agentId.String(), pinpointUrl, &agentCount)
		}
	}

	waitGroup.Wait()

	var k int64 = 0
	fmt.Printf("agentCount : %d\n", agentCount)
	applicationStats := make([]ApplicationStat, agentCount, agentCount)
	for k = 0; k < agentCount; k++ {
		applicationStats[k] = <-resutChan
	}

	close(resutChan)

	fmt.Printf("applicationStats : %v\n", applicationStats)
	fmt.Println("Finish ============================")

	util.RenderJsonResponse(applicationStats, w)
}

func goroutinAppList(c chan ApplicationStat, waitGroup *sync.WaitGroup, appName string, agentId string, pinpointUrl string, agentCount *int64) {
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

/**
Rest Api call
*/
func getRestCall(url string) (string, error) {
	fmt.Println(url)
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

	if jpath == "charts.y.AVG.#.2" {
		fmt.Println("mapData : %v\n", mapData)
	}

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

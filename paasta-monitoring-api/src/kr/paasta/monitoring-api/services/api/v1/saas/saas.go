package saas

import (
	"github.com/tidwall/gjson"
	"fmt"
	"paasta-monitoring-api/helpers"
	models "paasta-monitoring-api/models/api/v1"
	"strconv"
	"time"
)

type SaasService struct {
	SaaS models.SaaS
}

func GetSaasService(saas models.SaaS) *SaasService {
	return &SaasService{
		SaaS: saas,
	}
}


func (service *SaasService) GetApplicationStatus() (map[string]int, error) {
	result := make(map[string]int)
	result["running"] = 0
	result["shutdown"] = 0
	result["disconnect"] = 0

	resultBytes, _ := helpers.RequestHttpGet(service.SaaS.PinpointWebUrl+"/getAgentList.pinpoint", "","")
	resultJson := gjson.Parse(string(resultBytes))

	resultJson.ForEach(func(key, value gjson.Result) bool {
		result["agentCount"]++
		for _, item := range value.Array() {
			statusCode := item.Get("status.state.code").Int()
			switch statusCode {
			case 100:
				result["running"]++
			case 200:
			case 201:
				result["shutdown"]++
			case 300:
				result["disconnect"]++
			}
		}
		return true
	})

	return result, nil
}


func (service *SaasService) GetApplicationUsage(period string)(map[string]interface{}, error) {
	result := make(map[string]interface{})
	var from string
	var to string
	var cpuUsage float64
	var heapUsage float64
	var heapMaxUsage float64
	var noneHeapUsage float64
	var noneHeapMaxUsage float64
	var cpuUsageDataCount int

	to = strconv.FormatInt(time.Now().UTC().Unix(), 10) + "000"
	if len(period) == 0 {
		from = strconv.FormatInt(time.Now().Add(-600*time.Second).UTC().Unix(), 10) + "000"
	} else {
		periodNum, _ := strconv.Atoi(period[0:1])
		periodUnit := period[1:2]
		switch periodUnit {
		case "m" :
			periodNum = periodNum
		case "h" :
			periodNum = 60*periodNum;
		case "d" :
			periodNum = 1400*periodNum;
		}
		from = strconv.FormatInt(time.Now().Add(time.Duration(-periodNum)*time.Minute).UTC().Unix(), 10) + "000"
	}


	resultBytes, _ := helpers.RequestHttpGet(service.SaaS.PinpointWebUrl+"/getAgentList.pinpoint", "","")
	resultJson := gjson.Parse(string(resultBytes))
	resultJson.ForEach(func(key, value gjson.Result) bool {
		for _, item := range value.Array() {
			//statusCode := item.Get("status.state.code").Int()
			agentId := item.Get("agentId").String()
			queryString := "agentId="+agentId +"&from=" + from + "&to=" + to

			cpuUsageBytes, _ := helpers.RequestHttpGet(service.SaaS.PinpointWebUrl+"/getAgentStat/cpuLoad/chart.pinpoint", queryString,"")
			cpuUsageArray := gjson.Get(string(cpuUsageBytes), "charts.y.CPU_LOAD_SYSTEM.#.2").Array()
			cpuUsageSum := summaryValueInGjsonArray(cpuUsageArray)
			cpuUsageDataCount += len(cpuUsageArray)
			cpuUsage += cpuUsageSum

			memoryUsageBytes, _ := helpers.RequestHttpGet(service.SaaS.PinpointWebUrl+"/getAgentStat/jvmGc/chart.pinpoint", queryString,"")
			heapUsageArray := gjson.Get(string(memoryUsageBytes), "charts.y.JVM_MEMORY_HEAP_USED.#.2").Array()
			heapMaxUsageArray := gjson.Get(string(memoryUsageBytes), "charts.y.JVM_MEMORY_HEAP_MAX.#.2").Array()
			noneHeapUsageArray := gjson.Get(string(memoryUsageBytes), "charts.y.JVM_MEMORY_NON_HEAP_USED.#.2").Array()
			noneHeapMaxUsageArray := gjson.Get(string(memoryUsageBytes), "charts.y.JVM_MEMORY_NON_HEAP_MAX.#.2").Array()
			heapUsageSum := summaryValueInGjsonArray(heapUsageArray)
			heapUsage += heapUsageSum
			heapMaxUsageSum := summaryValueInGjsonArray(heapMaxUsageArray)
			heapMaxUsage += heapMaxUsageSum
			noneHeapUsageSum := summaryValueInGjsonArray(noneHeapUsageArray)
			noneHeapUsage += noneHeapUsageSum
			noneHeapMaxUsageSum := summaryValueInGjsonArray(noneHeapMaxUsageArray)
			noneHeapMaxUsage += noneHeapMaxUsageSum
		}
		return true
	})
	cpuPercent, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", cpuUsage/float64(cpuUsageDataCount)), 0)

	result["systemCpuRate"] = cpuPercent
	result["heapMemory"] = heapUsage
	result["heapMaxMemory"] = heapMaxUsage
	result["noneHeapMemory"] = noneHeapUsage
	result["noneHeapMaxMemory"] = noneHeapMaxUsage

	return result, nil
}

func summaryValueInGjsonArray(array []gjson.Result) float64 {
	var total float64
	for _, usage := range array {
		usageVal := usage.Float()
		total += usageVal
	}
	return total
}
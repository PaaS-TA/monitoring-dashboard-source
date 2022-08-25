package saas

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
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


func (service *SaasService) GetApplicationStatus(ctx echo.Context) (map[string]int, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	result := make(map[string]int)
	result["running"] = 0
	result["shutdown"] = 0
	result["disconnect"] = 0

	resultBytes, err := helpers.RequestHttpGet(service.SaaS.PinpointWebUrl+"/getAgentList.pinpoint", "","")
	if err != nil {
		logger.Error(err)
	}
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


func (service *SaasService) GetApplicationUsage(ctx echo.Context)(map[string]interface{}, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	period := ctx.QueryParam("period")
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

	resultBytes, err := helpers.RequestHttpGet(service.SaaS.PinpointWebUrl+"/getAgentList.pinpoint", "","")
	if err != nil {
		logger.Error(err)
	}

	resultJson := gjson.Parse(string(resultBytes))
	resultJson.ForEach(func(key, value gjson.Result) bool {
		for _, item := range value.Array() {
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


func (service *SaasService) GetApplicationUsageList(ctx echo.Context)([]map[string]interface{}, error) {
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	var results []map[string]interface{}
	var from string
	var to string
	to = strconv.FormatInt(time.Now().UTC().Unix(), 10) + "000"
	period := ctx.QueryParam("period")
	if len(period) == 0 {
		from = strconv.FormatInt(time.Now().Add(-600*time.Second).UTC().Unix(), 10) + "000"
	} else {
		periodNum, _ := strconv.Atoi(period[0:1])
		periodUnit := period[1:2]
		switch periodUnit {
		case "m" :
			periodNum = periodNum
		case "h" :
			periodNum = 60*periodNum
		case "d" :
			periodNum = 1400*periodNum
		}
		from = strconv.FormatInt(time.Now().Add(time.Duration(-periodNum)*time.Minute).UTC().Unix(), 10) + "000"
	}

	resultBytes, err := helpers.RequestHttpGet(service.SaaS.PinpointWebUrl+"/getAgentList.pinpoint", "","")
	if err != nil {
		logger.Error(err)
	}

	resultJson := gjson.Parse(string(resultBytes))
	resultJson.ForEach(func(key, value gjson.Result) bool {
		appDataMap := make(map[string]interface{})
		for _, item := range value.Array() {
			//statusCode := item.Get("status.state.code").Int()
			agentId := item.Get("agentId").String()
			appName := item.Get("applicationName").String()
			ipAddr := item.Get("ip").String()
			queryString := "agentId="+agentId +"&from=" + from + "&to=" + to

			cpuUsageBytes, _ := helpers.RequestHttpGet(service.SaaS.PinpointWebUrl+"/getAgentStat/cpuLoad/chart.pinpoint", queryString,"")
			cpuUsageArray := gjson.Get(string(cpuUsageBytes), "charts.y.CPU_LOAD_SYSTEM.#.2").Array()
			cpuUsageSum := summaryValueInGjsonArray(cpuUsageArray)
			cpuPercent, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", cpuUsageSum/float64(len(cpuUsageArray))), 0)
			jvmCpuUsageArray := gjson.Get(string(cpuUsageBytes), "charts.y.CPU_LOAD_JVM.#.2").Array()
			jvmCpuUsageSum := summaryValueInGjsonArray(jvmCpuUsageArray)
			jvmCpuPercent, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", jvmCpuUsageSum/float64(len(jvmCpuUsageArray))), 0)

			memoryUsageBytes, _ := helpers.RequestHttpGet(service.SaaS.PinpointWebUrl+"/getAgentStat/jvmGc/chart.pinpoint", queryString,"")
			heapUsageArray := gjson.Get(string(memoryUsageBytes), "charts.y.JVM_MEMORY_HEAP_USED.#.2").Array()
			heapUsageSum := summaryValueInGjsonArray(heapUsageArray)
			heapUsage, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", heapUsageSum/float64(len(heapUsageArray))), 0)

			noneHeapUsageArray := gjson.Get(string(memoryUsageBytes), "charts.y.JVM_MEMORY_NON_HEAP_USED.#.2").Array()
			noneHeapUsageSum := summaryValueInGjsonArray(noneHeapUsageArray)
			noneHeapUsage, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", noneHeapUsageSum/float64(len(noneHeapUsageArray))), 0)

			activeTraceBytes, _ := helpers.RequestHttpGet(service.SaaS.PinpointWebUrl+"/getAgentStat/activeTrace/chart.pinpoint", queryString,"")
			activeTraceArray := gjson.Get(string(activeTraceBytes), "charts.y.ACTIVE_TRACE_FAST.#.3").Array()
			activeTraceSum := summaryValueInGjsonArray(activeTraceArray)
			activeTrace, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", activeTraceSum/float64(len(activeTraceArray))), 0)
			if activeTrace < 0 {
				activeTrace = 0
			}

			responseTimeBytes, _ := helpers.RequestHttpGet(service.SaaS.PinpointWebUrl+"/getAgentStat/responseTime/chart.pinpoint", queryString,"")
			responseTimeArray := gjson.Get(string(responseTimeBytes), "charts.y.AVG.#.2").Array()
			responseTimeSum := summaryValueInGjsonArray(responseTimeArray)
			responseTime, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", responseTimeSum/float64(len(responseTimeArray))), 0)
			if responseTime < 0 {
				responseTime = 0
			}

			appDataMap["agentId"] = agentId
			appDataMap["applicationName"] = appName
			appDataMap["ip"] = ipAddr
			appDataMap["cpuUsage"] = cpuPercent
			appDataMap["jvmCpuUsage"] = jvmCpuPercent
			appDataMap["heapUsage"] = heapUsage
			appDataMap["noneHeapUsage"] = noneHeapUsage
			appDataMap["activeTrace"] = activeTrace
			appDataMap["responseTime"] = responseTime
			results = append(results, appDataMap)
		}
		return true
	})
	return results, nil
}


func summaryValueInGjsonArray(array []gjson.Result) float64 {
	var total float64
	for _, usage := range array {
		usageVal := usage.Float()
		total += usageVal
	}
	return total
}
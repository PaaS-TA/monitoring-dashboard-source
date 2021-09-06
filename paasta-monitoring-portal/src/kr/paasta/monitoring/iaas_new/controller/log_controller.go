package controller

import (
	client "github.com/influxdata/influxdb1-client/v2"
	/*"gopkg.in/olivere/elastic.v3"*/
	"kr/paasta/monitoring/iaas/model"
	"kr/paasta/monitoring/iaas/service"
	"kr/paasta/monitoring/utils"
	"net/http"
	"strconv"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
)

type OpenstackLog struct {
	OpenstackProvider model.OpenstackProvider
	influxClient      client.Client
	ElasticClient     *elasticsearch.Client
}

func NewLogController(openstackProvider model.OpenstackProvider, influxClient client.Client, elasticClient *elasticsearch.Client) *OpenstackLog {
	s := &OpenstackLog{
		OpenstackProvider: openstackProvider,
		influxClient:      influxClient,
		ElasticClient:     elasticClient,
	}
	return s
}

func (s *OpenstackLog) GetDefaultRecentLog(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.LogMessage
	apiRequest.Hostname = r.URL.Query().Get("hostname")
	apiRequest.LogType = r.URL.Query().Get("logType")
	pageItems, _ := strconv.Atoi(r.URL.Query().Get("pageItems"))
	pageIndex, _ := strconv.Atoi(r.URL.Query().Get("pageIndex"))
	apiRequest.PageItems = pageItems
	apiRequest.PageIndex = pageIndex
	apiRequest.Keyword = r.URL.Query().Get("keyword")
	apiRequest.StartTime = r.URL.Query().Get("startTime")
	apiRequest.EndTime = r.URL.Query().Get("endTime")
	//apiRequest.Index = r.URL.Query().Get("logstashIndex")

	period := r.URL.Query().Get("period")

	if period != "" {
		time_unit := period[len(period)-1:]
		if time_unit == "h" {
			apiRequest.Period, _ = strconv.ParseInt(period[:len(period)-1], 10, 64)
			apiRequest.Period = apiRequest.Period * 60
		} else if time_unit == "m" {
			apiRequest.Period, _ = strconv.ParseInt(period[:len(period)-1], 10, 64)
		} else {
			errMessage := map[string]interface{}{"Persons": "Time unit is only allowed 'm' and 'h' -(ex) 5m, 5h"}
			utils.RenderJsonResponse(errMessage, w)
		}
	}
	validation := apiRequest.DefaultLogValidate(apiRequest)
	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validation.Error()))
		return
	}

	logInfo, err := services.GetLogService(s.ElasticClient).GetDefaultRecentLog(apiRequest, true)
	if err != nil {
		utils.ErrRenderJsonResponse(err, w)
	} else {
		utils.RenderJsonResponse(logInfo, w)
	}

}

func (s *OpenstackLog) GetSpecificTimeRangeLog(w http.ResponseWriter, r *http.Request) {

	var apiRequest model.LogMessage
	apiRequest.Hostname = r.URL.Query().Get("hostname")
	apiRequest.LogType = r.URL.Query().Get("logType")
	pageItems, _ := strconv.Atoi(r.URL.Query().Get("pageItems"))
	pageIndex, _ := strconv.Atoi(r.URL.Query().Get("pageIndex"))
	apiRequest.PageItems = pageItems
	apiRequest.PageIndex = pageIndex
	apiRequest.TargetDate = r.URL.Query().Get("targetDate")
	apiRequest.Keyword = r.URL.Query().Get("keyword")
	apiRequest.StartTime = r.URL.Query().Get("startTime")
	apiRequest.EndTime = r.URL.Query().Get("endTime")
	//apiRequest.Index = r.URL.Query().Get("logstashIndex")

	validation := apiRequest.SpecificTimeRangeLogValidate(apiRequest)
	if validation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(validation.Error()))
		return
	}

	cpuUsageList, err := services.GetLogService(s.ElasticClient).GetSpecificTimeRangeLog(apiRequest, true)
	if err != nil {
		utils.ErrRenderJsonResponse(err, w)
	} else {
		utils.RenderJsonResponse(cpuUsageList, w)
	}

}

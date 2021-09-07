package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

	/*"gopkg.in/olivere/elastic.v3"*/
	"kr/paasta/monitoring/iaas_new/model"
	"math"
	"strings"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	esapi "github.com/elastic/go-elasticsearch/v7/esapi"
)

type LogDao struct {
	elasticClient *elasticsearch.Client
}

func GetLogDao(elasticClient *elasticsearch.Client) *LogDao {
	return &LogDao{
		elasticClient: elasticClient,
	}
}

func setInitRequest(initRequest model.LogMessage) (request model.LogMessage){
	//initRequest.Period = 1
	//initRequest.PageIndex = 3
	//initRequest.PageItems = 10

	// test
	//if initRequest.Index == "" {
	//
	//	// Default Target vm's recent log - recent 30 minutes.
	//	now := time.Now().Local()
	//	current := now.Unix() - int64(9)*60*60 //9 hour difference Between Local PC and GMT(Greenwich Mean Time).
	//	//current := now.Unix()
	//	/*
	//		조회 주기정보를 전달받아 로그를 조회한다.(period - '분'단위)
	//	*/
	//	//Current Time Stamp
	//	before := now.Unix() - initRequest.Period*60 //화면에서 설정한 조회주기(분) (ex: 30 * 60 seconds)
	//	before = before - int64(9)*60*60 //9시간	=> if time zone is equal to Logsearch Instance, it must be removed.
	//	//9 hour difference Between Local PC and Virginia.
	//	initRequest.StartTime = time.Unix(before, 0).Format(time.RFC3339)[0:19]
	//	initRequest.EndTime = time.Unix(current, 0).Format(time.RFC3339)[0:19]
	//	//request.Index = fmt.Sprintf("logstash-%s", fmt.Sprintf("%d.%s.%s", time.Unix(current, 0).Year(), attachZero(int(time.Unix(current, 0).Month())), attachZero(time.Unix(current, 0).Day())))
	//	initRequest.Index = fmt.Sprintf("logstash-2021.03.02")
	//	initRequest.TargetDate = fmt.Sprintf("%s", fmt.Sprintf("%d.%s.%s", time.Unix(current, 0).Year(), attachZero(int(time.Unix(current, 0).Month())), attachZero(time.Unix(current, 0).Day())))
	//}

	// origin
	if initRequest.Index == "" {
		// Default Target vm's recent log - recent 30 minutes.
		now := time.Now().Local()
		//current := now.Unix() - int64(iaasmodel.GMTTimeGap)*60*60 //9 hour difference Between Local PC and GMT(Greenwich Mean Time).
		current := now.Unix()
		/*
			조회 주기정보를 전달받아 로그를 조회한다.(period - '분'단위)
		*/
		//Current Time Stamp
		before := now.Unix() - initRequest.Period*60 //화면에서 설정한 조회주기(분) (ex: 30 * 60 seconds)
		//before = before - int64(iaasmodel.GMTTimeGap)*60*60 //9시간	=> if time zone is equal to Logsearch Instance, it must be removed.
		//9 hour difference Between Local PC and Virginia.
		initRequest.StartTime = time.Unix(before, 0).Format(time.RFC3339)[0:19]
		initRequest.EndTime = time.Unix(current, 0).Format(time.RFC3339)[0:19]
		initRequest.Index = fmt.Sprintf("filebeat-%s", fmt.Sprintf("%d.%s.%s", time.Unix(current, 0).Year(), attachZero(int(time.Unix(current, 0).Month())), attachZero(time.Unix(current, 0).Day())))
		initRequest.TargetDate = fmt.Sprintf("%s", fmt.Sprintf("%d.%s.%s", time.Unix(current, 0).Year(), attachZero(int(time.Unix(current, 0).Month())), attachZero(time.Unix(current, 0).Day())))
	}

	fmt.Println(initRequest)
	return initRequest
}
func retrieveElasticClientInfo(elasticClient *elasticsearch.Client) (response *esapi.Response, errMsg model.ErrMessage){
	response, err := elasticClient.Info()
	if err != nil {
		fmt.Println("Error retrieveElasticClientInfo: %s", err)
		errMsg = model.ErrMessage{
			"Message": err.Error(),
		}
		return response, errMsg
	}
	defer response.Body.Close()
	// check elasticClient Info
	bytes, _:= ioutil.ReadAll(response.Body)
	str := string(bytes)
	fmt.Println(str)
	return response, nil
}
func setDefaultRecentLogQuery(request model.LogMessage, paging bool, search_count int)(query string){
	//example variable
	//request.Keyword = "guardian.list-containers.finished"
	//request.LogType = "bosh"
	//request.Id = "022eab7c-f0ad-4351-8616-43473aa6bd84"

	//!! not execute (term - syslog_sd_params.id), so use match_phrase @raw !!
	query = `{
		"query" : {
			"bool": {
				"must": [`

	if request.Keyword != ""{
		query +=
						`{
							"match_phrase" : {
								"message" : "`+request.Keyword+`"
							}
						},`
	}
	query +=
						`{
							"term": {
						  		"beat.hostname": "`+request.Hostname+`"
							}
						}
				],
				"filter": {
					"range": {
							"@timestamp": {
						  		"gte": "`+request.StartTime+`",
						  		"lte": "`+request.EndTime+`"
							}
					 	}
				},
				"boost" : 5.0
			}
		},
		"sort": {
			"@timestamp": {
			  "order": "desc"
			}
		}`
	if paging == true {
		query += `,
					"from" : "` + strconv.Itoa((request.PageIndex-1)*request.PageItems) + `",
					"size" : "` + strconv.Itoa(search_count) + `"`
	}else {
		query += `,
					"from" : "0",
					"size" : "` + strconv.Itoa(search_count) + `"`
	}

	query += `}`
	fmt.Println(query)
	return query
}
func retrieveElasticClientSearch(elasticClient *elasticsearch.Client, request model.LogMessage, query string) (body []byte, errMsg model.ErrMessage) {
	// Perform the search request.
	response, err := elasticClient.Search(
		elasticClient.Search.WithContext(context.Background()),
		elasticClient.Search.WithIndex(request.Index),
		elasticClient.Search.WithBody(strings.NewReader(query)),
		elasticClient.Search.WithTrackTotalHits(true),
		elasticClient.Search.WithPretty(),
	)
	if err != nil {
		fmt.Println("Error retrieveElasticClientSearch : %s", err)
		errMsg = model.ErrMessage{
			"Message": err.Error(),
		}
		if strings.Contains(err.Error(), "timeout awaiting response headers")  {
			fmt.Println("retry retrieveElasticClientSearch")
			response, err = elasticClient.Search(
				elasticClient.Search.WithContext(context.Background()),
				elasticClient.Search.WithIndex(request.Index),
				elasticClient.Search.WithBody(strings.NewReader(query)),
				elasticClient.Search.WithTrackTotalHits(true),
				elasticClient.Search.WithPretty(),
			)
		}else{
			return nil, errMsg
		}
	}
	defer response.Body.Close()

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error retrieveElasticClientSearch ioutil.ReadAll : %s", err)
		errMsg = model.ErrMessage{
			"Message": err.Error(),
		}
		return nil, errMsg
	}
	//str := string(body)
	//fmt.Println(str)

	return body, errMsg
}
func searchDataJsonUnmarshal(body []byte)(hits map[string]interface{}, totalHits float64, errMsg model.ErrMessage){
	// body ([]byte) -> JsonUnmarsal -> data(map[string]interface)
	var data map[string]interface{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error searchDataJsonUnmarshal json.Unmarshal : %s", err)
		errMsg = model.ErrMessage{
			"Message": err.Error(),
		}
		return nil, 0, errMsg
	}
	// response data
	//fmt.Printf("Results: %v\n", data)
	// 1 depth
	//fmt.Printf("Results: %v\n", data["_shards"])
	//fmt.Printf("Results: %v\n", data["hits"])
	//fmt.Printf("Results: %v\n", data["total"])

	hits = data["hits"].(map[string]interface{}) // 1 depth
	//fmt.Println(hits["total"])
	total := hits["total"].(map[string]interface{}) // 2 depth
	//fmt.Println(total["value"])
	totalHits = total["value"].(float64)

	return hits, totalHits, errMsg
}
func setLogInfoList(hits map[string]interface{})(logInfoList []model.LogInfo, errMsg model.ErrMessage){
	for _, hit := range hits["hits"].([]interface {}) {
		//fmt.Println("=======-------",hit)
		//_id := hit.(map[string]interface {})["_id"]
		_source := hit.(map[string]interface {})["_source"]
		messageObj := _source.(map[string]interface {})["@message"]
		timestamp := _source.(map[string]interface {})["@timestamp"]

		var logInfo model.LogInfo

		resTime := strings.Replace(timestamp.(string), "\\", "", -1)
		resTime = strings.Replace(resTime, "\"", "", -1)
		convert_time, err := time.Parse(time.RFC3339, resTime)
		if err != nil {
			fmt.Println("Recent TimeLogs - Time Conversion error :", err)
			errMsg = model.ErrMessage{
				"Message": err.Error(),
			}
			return nil, errMsg
		}
		logInfo.Time = time.Unix(convert_time.Unix()+int64(9)*60*60, 0).Format(time.RFC3339)[0:19]
		logInfo.Message = messageObj.(string)
		logInfoList = append(logInfoList, logInfo)
	}

	return logInfoList, nil
}


//IaaS 선택된 항목의 최근 로그를 조회한다.
func (log LogDao) GetDefaultRecentLog(request model.LogMessage, paging bool) (_ model.LogMessage, errMsg model.ErrMessage) {

	request = setInitRequest(request)
	/*response, err := retrieveElasticClientInfo(log.elasticClient)
	fmt.Println(response)
	if err != nil {
		return request, err
	}*/

	query := setDefaultRecentLogQuery(request, false, 0)
	body, err := retrieveElasticClientSearch(log.elasticClient, request, query)
	hits, totalHits, err := searchDataJsonUnmarshal(body)
	if err != nil {
		return request, err
	}

	if paging {
		totalPages := int(totalHits) / request.PageItems

		var search_count int
		if request.PageIndex > totalPages {
			search_count = int(math.Mod(float64(totalHits), float64(request.PageItems)))
		} else {
			search_count = request.PageItems
		}

		totalCount := int(totalHits)
		if totalCount > 10000 {
			totalCount = 9999
		}
		request.TotalCount = totalCount
		request.CurrentItems = search_count

		query := setDefaultRecentLogQuery(request, true, search_count)
		body, err = retrieveElasticClientSearch(log.elasticClient, request, query)
		hits, totalHits, err = searchDataJsonUnmarshal(body)
		if err != nil {
			return request, err
		}
	} else {
		search_count := int(totalHits)
		if search_count > 10000 {
			search_count = 9999
		}

		query := setDefaultRecentLogQuery(request, false, search_count)
		body, err = retrieveElasticClientSearch(log.elasticClient, request, query)
		hits, totalHits, err = searchDataJsonUnmarshal(body)
		if err != nil {
			return request, err
		}
	}

	request.Messages, err = setLogInfoList(hits)
	if err != nil {
		return request, err
	}

	return request, nil
}

//PaaS 선택된 항목의 시간 별 로그를 조회한다.
func (log LogDao) GetSpecificTimeRangeLog(request model.LogMessage, paging bool) (_ model.LogMessage, errMsg model.ErrMessage) {

	// setInitRequest
	if request.Index == "" {

		//To get Index name do not use TargetDate. Instead, use startTime.
		date_array := strings.Split(request.TargetDate, "-")
		if len(date_array) != 3 {
			errMessage := model.ErrMessage{
				"Message": errors.New("request target date error:" + request.TargetDate),
			}
			return request, errMessage
		}

		if request.StartTime == "" && request.EndTime == "" {
			request.StartTime = fmt.Sprintf("%sT%s", request.TargetDate, "00:00:00")
			request.EndTime = fmt.Sprintf("%sT%s", request.TargetDate, "23:59:59")
		} else if request.StartTime != "" && request.EndTime == "" {
			request.StartTime = fmt.Sprintf("%sT%s", request.TargetDate, request.StartTime)
			request.EndTime = fmt.Sprintf("%sT%s", request.TargetDate, "23:59:59")
		} else if request.StartTime == "" && request.EndTime != "" {
			request.StartTime = fmt.Sprintf("%sT%s", request.TargetDate, "00:00:00")
			request.EndTime = fmt.Sprintf("%sT%s", request.TargetDate, request.EndTime)
		} else {
			request.StartTime = fmt.Sprintf("%sT%s", request.TargetDate, request.StartTime)
			request.EndTime = fmt.Sprintf("%sT%s", request.TargetDate, request.EndTime)
		}

		//=================================================================================================
		// It will be deleted later.  Now it needs only for Time-zone difference between Local and Virginia
		//=================================================================================================
		convert_start_time, err := time.Parse(time.RFC3339, fmt.Sprintf("%s+09:00", request.StartTime))
		if err != nil {
			//fmt.Println("SpecificTimeLogs - Time Conversion error :", err)
			errMessage := model.ErrMessage{
				"Message": err.Error(),
			}
			return request, errMessage
		}
		convert_end_time, err := time.Parse(time.RFC3339, fmt.Sprintf("%s+09:00", request.EndTime))
		if err != nil {
			//fmt.Println("SpecificTimeLogs - Time Conversion error :", err)
			errMessage := model.ErrMessage{
				"Message": err.Error(),
			}
			return request, errMessage
		}

		//end := convert_end_time.Unix() 	- int64(model.GMTTimeGap)*60*60  //9 hour difference Between Local PC and GMT(Greenwich Mean Time).
		//start := convert_start_time.Unix() - int64(model.GMTTimeGap)*60*60 //9 hour difference Between Local PC and GMT(Greenwich Mean Time).
		end := convert_end_time.Unix()
		start := convert_start_time.Unix()

		request.StartTime = time.Unix(start, 0).Format(time.RFC3339)[0:19]
		request.EndTime = time.Unix(end, 0).Format(time.RFC3339)[0:19]

		request.Index = fmt.Sprintf("filebeat-%s", fmt.Sprintf("%d.%s.%s", time.Unix(start, 0).Year(), attachZero(int(time.Unix(start, 0).Month())), attachZero(time.Unix(start, 0).Day())))
	}

	// retrieveElasticClientInfo
	/*response, err := retrieveElasticClientInfo(log.elasticClient)
	fmt.Println(response)
	if err != nil {
		return request, err
	}*/

	query := setDefaultRecentLogQuery(request, false, 0)
	body, err := retrieveElasticClientSearch(log.elasticClient, request, query)
	hits, totalHits, err := searchDataJsonUnmarshal(body)
	if err != nil {
		return request, err
	}

	if paging {
		totalPages := int(totalHits) / request.PageItems

		var search_count int
		if request.PageIndex > totalPages {
			search_count = int(math.Mod(float64(totalHits), float64(request.PageItems)))
		} else {
			search_count = request.PageItems
		}

		totalCount := int(totalHits)
		if totalCount > 9001 {
			totalCount = 9000
		}
		request.TotalCount = totalCount
		request.CurrentItems = search_count

		query := setDefaultRecentLogQuery(request, true, search_count)
		body, err = retrieveElasticClientSearch(log.elasticClient, request, query)
		hits, totalHits, err = searchDataJsonUnmarshal(body)
		if err != nil {
			return request, err
		}
	} else {
		search_count := int(totalHits)
		if search_count > 9001 {
			search_count = 9000
		}

		query := setDefaultRecentLogQuery(request, false, search_count)
		body, err = retrieveElasticClientSearch(log.elasticClient, request, query)
		hits, totalHits, err = searchDataJsonUnmarshal(body)
		if err != nil {
			return request, err
		}
	}

	request.Messages, err = setLogInfoList(hits)
	if err != nil {
		return request, err
	}

	return request, nil
}

func attachZero(num int) string {
	if num < 10 {
		return fmt.Sprintf("0%d", num)
	} else {
		return fmt.Sprintf("%d", num)
	}
}

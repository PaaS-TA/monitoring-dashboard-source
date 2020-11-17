package dao

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/olivere/elastic.v3"
	iaasmodel "kr/paasta/monitoring/iaas/model"
	"kr/paasta/monitoring/paas/model"
	"math"
	"strings"
	"time"
)

type PaasLogDao struct {
	elasticClient *elastic.Client
}

func GetPaasLogDao(elasticClient *elastic.Client) *PaasLogDao {
	return &PaasLogDao{
		elasticClient: elasticClient,
	}
}

//Node의 현재 CPU사용률을 조회한다.
func (log PaasLogDao) GetDefaultRecentLog(request model.LogMessage, paging bool) (_ model.LogMessage, errMsg iaasmodel.ErrMessage) {

	if request.Index == "" {
		// Default Target vm's recent log - recent 30 minutes.
		now := time.Now().Local()
		//current := now.Unix() - int64(iaasmodel.GMTTimeGap)*60*60 //9 hour difference Between Local PC and GMT(Greenwich Mean Time).
		current := now.Unix()
		/*
			조회 주기정보를 전달받아 로그를 조회한다.(period - '분'단위)
		*/
		//Current Time Stamp
		before := now.Unix() - request.Period*60 //화면에서 설정한 조회주기(분) (ex: 30 * 60 seconds)
		//before = before - int64(iaasmodel.GMTTimeGap)*60*60 //9시간	=> if time zone is equal to Logsearch Instance, it must be removed.
		//9 hour difference Between Local PC and Virginia.
		request.StartTime = time.Unix(before, 0).Format(time.RFC3339)[0:19]
		request.EndTime = time.Unix(current, 0).Format(time.RFC3339)[0:19]
		request.Index = fmt.Sprintf("logs-platform-%s", fmt.Sprintf("%d.%s.%s", time.Unix(current, 0).Year(), attachZero(int(time.Unix(current, 0).Month())), attachZero(time.Unix(current, 0).Day())))
		request.TargetDate = fmt.Sprintf("%s", fmt.Sprintf("%d.%s.%s", time.Unix(current, 0).Year(), attachZero(int(time.Unix(current, 0).Month())), attachZero(time.Unix(current, 0).Day())))
	}

	exists, err := log.elasticClient.IndexExists(request.Index).Do()
	if err != nil || !exists {
		//fmt.Println("index doesn't exists :", err)
		errMsg = iaasmodel.ErrMessage{
			"Message":    fmt.Sprintf("target index - %s - doesn't exists.", request.Index),
			"HttpStatus": 200,
		}
		return request, errMsg
	}

	//fmt.Println("elastic request =====>>>>", request)

	if exists {
		reqQuery := elastic.NewBoolQuery()
		if request.Keyword != "" {
			reqQuery = reqQuery.Must(elastic.NewMatchPhraseQuery("@message", request.Keyword))
		}

		if request.LogType == model.BOSH_NAME {
			reqQuery = reqQuery.Must(elastic.NewTermQuery("syslog_sd_params.deployment", request.LogType))
		} else {
			reqQuery = reqQuery.Must(elastic.NewTermQuery("syslog_sd_params.id", request.Id))
		}

		reqQuery = reqQuery.Filter(elastic.NewRangeQuery("@timestamp").Gte(request.StartTime).Lte(request.EndTime))
		reqQuery = reqQuery.Boost(5)
		reqQuery = reqQuery.DisableCoord(true)
		reqQuery = reqQuery.QueryName("Default_Recent_logs")

		// Search with a term query for totalcount
		searchResult, err := log.elasticClient.Search().
			Index(request.Index).
			Query(reqQuery).           // specify the query
			Sort("@timestamp", false). // sort by "timestamp" field, ascending - true, descending - false
			Pretty(true).              // pretty print request and response JSON
			Do()                       // execute

		if err != nil {
			errMessage := iaasmodel.ErrMessage{
				"Message": err.Error(),
			}
			return request, errMessage
		}

		if paging {
			totalPages := int(searchResult.TotalHits()) / request.PageItems

			var search_count int
			if request.PageIndex > totalPages {
				search_count = int(math.Mod(float64(searchResult.TotalHits()), float64(request.PageItems)))
			} else {
				search_count = request.PageItems
			}

			totalCount := int(searchResult.TotalHits())
			if totalCount > 10000 {
				totalCount = 9999
			}
			request.TotalCount = totalCount
			request.CurrentItems = search_count

			// Search with a term query
			searchResult, err = log.elasticClient.Search().
				Index(request.Index).
				Query(reqQuery).           // specify the query
				Sort("@timestamp", false). // sort by "timestamp" field, ascending - true, descending - false
				From((request.PageIndex - 1) * request.PageItems).Size(search_count).
				Pretty(true). // pretty print request and response JSON
				Do()          // execute
		} else {
			search_count := int(searchResult.TotalHits())
			if search_count > 10000 {
				search_count = 9999
			}
			// Search with a term query
			searchResult, err = log.elasticClient.Search().
				Index(request.Index).
				Query(reqQuery).           // specify the query
				Sort("@timestamp", false). // sort by "timestamp" field, ascending - true, descending - false
				From(0).Size(search_count).
				Pretty(true). // pretty print request and response JSON
				Do()          // execute
		}

		var logInfoList []model.LogInfo

		for _, hit := range searchResult.Hits.Hits {
			//convert the result of searching to Map interface
			rawData := make(map[string]json.RawMessage)
			jsondata, err := hit.Source.MarshalJSON()
			if err != nil {
				//fmt.Println("Json Marshal error :", err)
				errMessage := iaasmodel.ErrMessage{
					"Message": err.Error(),
				}
				return request, errMessage
			} else {
				//fmt.Println("original source :", string(jsondata))
				err = json.Unmarshal(jsondata, &rawData)
				if err != nil {
					//fmt.Println("#### Json Unmarshal error :", err)
					errMessage := iaasmodel.ErrMessage{
						"Message": err.Error(),
					}
					return request, errMessage
				} else {

					var logInfo model.LogInfo

					for key, value := range rawData {
						if strings.Contains(key, "timestamp") {
							resTime := strings.Replace(string(value), "\\", "", -1)
							resTime = strings.Replace(resTime, "\"", "", -1)
							convert_time, err := time.Parse(time.RFC3339, resTime)
							if err != nil {
								fmt.Println("Recent TimeLogs - Time Conversion error :", err)
							}
							logInfo.Time = time.Unix(convert_time.Unix()+int64(iaasmodel.GMTTimeGap)*60*60, 0).Format(time.RFC3339)[0:19]
						}

						if strings.Contains(key, "message") {
							logInfo.Message = strings.Replace(string(value), "\\", "", -1)
						}
					}
					logInfoList = append(logInfoList, logInfo)
				}
			}
		}

		request.Messages = logInfoList
	}
	return request, nil
}

//Node의 현재 CPU사용률을 조회한다.
func (log PaasLogDao) GetSpecificTimeRangeLog(request model.LogMessage, paging bool) (_ model.LogMessage, errMsg iaasmodel.ErrMessage) {

	if request.Index == "" {

		//To get Index name do not use TargetDate. Instead, use startTime.
		date_array := strings.Split(request.TargetDate, "-")
		if len(date_array) != 3 {
			errMessage := iaasmodel.ErrMessage{
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
			errMessage := iaasmodel.ErrMessage{
				"Message": err.Error(),
			}
			return request, errMessage
		}
		convert_end_time, err := time.Parse(time.RFC3339, fmt.Sprintf("%s+09:00", request.EndTime))
		if err != nil {
			//fmt.Println("SpecificTimeLogs - Time Conversion error :", err)
			errMessage := iaasmodel.ErrMessage{
				"Message": err.Error(),
			}
			return request, errMessage
		}

		//end := convert_end_time.Unix() 	- int64(iaasmodel.GMTTimeGap)*60*60  //9 hour difference Between Local PC and GMT(Greenwich Mean Time).
		//start := convert_start_time.Unix() - int64(iaasmodel.GMTTimeGap)*60*60 //9 hour difference Between Local PC and GMT(Greenwich Mean Time).
		end := convert_end_time.Unix()
		start := convert_start_time.Unix()

		request.StartTime = time.Unix(start, 0).Format(time.RFC3339)[0:19]
		request.EndTime = time.Unix(end, 0).Format(time.RFC3339)[0:19]

		request.Index = fmt.Sprintf("logs-platform-%s", fmt.Sprintf("%d.%s.%s", time.Unix(start, 0).Year(), attachZero(int(time.Unix(start, 0).Month())), attachZero(time.Unix(start, 0).Day())))
	}

	fmt.Println("elastic request =====>>>>", request)

	exists, err := log.elasticClient.IndexExists(request.Index).Do()
	if err != nil || !exists {
		//fmt.Println("SpecificTimeLogs - index doesn't exists :", err)
		errMsg = iaasmodel.ErrMessage{
			"Message": fmt.Sprintf("target index - %s - doesn't exists.", request.Index),
		}
		return request, errMsg
	}

	if exists {

		reqQuery := elastic.NewBoolQuery()
		if request.Keyword != "" {
			reqQuery = reqQuery.Must(elastic.NewMatchPhraseQuery("@message", request.Keyword))
		}

		if request.LogType == model.BOSH_NAME {
			reqQuery = reqQuery.Must(elastic.NewTermQuery("syslog_sd_params.deployment", request.LogType))
		} else {
			reqQuery = reqQuery.Must(elastic.NewTermQuery("syslog_sd_params.id", request.Id))
		}

		reqQuery = reqQuery.Filter(elastic.NewRangeQuery("@timestamp").Gte(request.StartTime).Lte(request.EndTime))
		reqQuery = reqQuery.Boost(5)
		reqQuery = reqQuery.DisableCoord(true)
		reqQuery = reqQuery.QueryName("Target_Time_range_logs")

		// Search with a term query for totalcount
		searchResult, err := log.elasticClient.Search().
			Index(request.Index).
			Query(reqQuery).           // specify the query
			Sort("@timestamp", false). // sort by "timestamp" field, ascending - true, descending - false
			Pretty(true).              // pretty print request and response JSON
			Do()                       // execute

		if err != nil {
			//fmt.Println("searching error :", err)
			errMessage := iaasmodel.ErrMessage{
				"Message": err.Error(),
			}
			return request, errMessage
		}

		if paging {
			totalPages := int(searchResult.TotalHits()) / request.PageItems

			var search_count int
			if request.PageIndex > totalPages {
				search_count = int(math.Mod(float64(searchResult.TotalHits()), float64(request.PageItems)))
			} else {
				search_count = request.PageItems
			}

			totalCount := int(searchResult.TotalHits())
			if totalCount > 10000 {
				totalCount = 9999
			}
			request.TotalCount = totalCount
			request.CurrentItems = search_count

			// Search with a term query
			searchResult, err = log.elasticClient.Search().
				Index(request.Index).
				Query(reqQuery).           // specify the query
				Sort("@timestamp", false). // sort by "timestamp" field, ascending - true, descending - false
				From((request.PageIndex - 1) * request.PageItems).Size(search_count).
				Pretty(true). // pretty print request and response JSON
				Do()          // execute
		} else {
			search_count := int(searchResult.TotalHits())
			if search_count > 10000 {
				search_count = 9999
			}

			// Search with a term query
			searchResult, err = log.elasticClient.Search().
				Index(request.Index).
				Query(reqQuery).           // specify the query
				Sort("@timestamp", false). // sort by "timestamp" field, ascending - true, descending - false
				From(0).Size(search_count).
				Pretty(true). // pretty print request and response JSON
				Do()          // execute
		}

		var logInfoList []model.LogInfo

		for _, hit := range searchResult.Hits.Hits {
			//convert the result of searching to Map interface
			rawData := make(map[string]json.RawMessage)
			jsondata, err := hit.Source.MarshalJSON()
			if err != nil {
				//fmt.Println("SpecificTimeLogs - Json Marshal error :", err)
				errMessage := iaasmodel.ErrMessage{
					"Message": err.Error(),
				}
				return request, errMessage
			} else {
				//fmt.Println("original source :", string(jsondata))
				err = json.Unmarshal(jsondata, &rawData)
				if err != nil {
					//fmt.Println("#### SpecificTimeLogs - Json Unmarshal error :", err)
					errMessage := iaasmodel.ErrMessage{
						"Message": err.Error(),
					}
					return request, errMessage
				} else {
					var logInfo model.LogInfo

					for key, value := range rawData {
						if strings.Contains(key, "timestamp") {
							resTime := strings.Replace(string(value), "\\", "", -1)
							resTime = strings.Replace(resTime, "\"", "", -1)
							convert_time, err := time.Parse(time.RFC3339, resTime)
							if err != nil {
								fmt.Println("Specific TimeLogs - Time Conversion error :", err)
							}
							logInfo.Time = time.Unix(convert_time.Unix()+int64(iaasmodel.GMTTimeGap)*60*60, 0).Format(time.RFC3339)[0:19]
						}

						if strings.Contains(key, "message") {
							logInfo.Message = strings.Replace(string(value), "\\", "", -1)
						}
					}
					logInfoList = append(logInfoList, logInfo)
				}
			}
			//index = index + 1
		}
		request.Messages = logInfoList
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

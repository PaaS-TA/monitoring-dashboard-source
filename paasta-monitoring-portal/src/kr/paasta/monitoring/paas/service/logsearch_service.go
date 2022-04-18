package service

import (
	"fmt"
	"github.com/influxdata/influxdb1-client/v2"
	"kr/paasta/monitoring/paas/dao"
	"kr/paasta/monitoring/paas/model"
	"time"
)

type LogsearchService struct {
	influxClient client.Client
	databases    model.Databases
}

func GetLogsearchService(influxClient client.Client, databases model.Databases) *LogsearchService {
	return &LogsearchService{
		influxClient: influxClient,
		databases:    databases,
	}
}

func (service LogsearchService) GetLogData(param model.NewLogMessage) (response client.Response, errMsg model.ErrMessage) {
	/**
		Period 파라미터가 존재하면 Period 값으로 DB 조회
		없으면 StartTime, EndTime 파라미터 값으로 DB조회
	 */
	if param.Period == "" {
		/**
			날짜 시간 값을 DB에서 조회할 수 있는 포맷으로 변경
		 */
		if param.StartTime == "" && param.EndTime == "" {
			param.StartTime = fmt.Sprintf("%sT%s", param.TargetDate, "00:00:00")
			param.EndTime = fmt.Sprintf("%sT%s", param.TargetDate, "23:59:59")
		} else if param.StartTime != "" && param.EndTime == "" {
			param.StartTime = fmt.Sprintf("%sT%s", param.TargetDate, param.StartTime)
			param.EndTime = fmt.Sprintf("%sT%s", param.TargetDate, "23:59:59")
		} else if param.StartTime == "" && param.EndTime != "" {
			param.StartTime = fmt.Sprintf("%sT%s", param.TargetDate, "00:00:00")
			param.EndTime = fmt.Sprintf("%sT%s", param.TargetDate, param.EndTime)
		} else {
			param.StartTime = fmt.Sprintf("%sT%s", param.TargetDate, param.StartTime)
			param.EndTime = fmt.Sprintf("%sT%s", param.TargetDate, param.EndTime)
		}

		convert_start_time, _ := time.Parse(time.RFC3339, fmt.Sprintf("%s+09:00", param.StartTime))
		convert_end_time, _ := time.Parse(time.RFC3339, fmt.Sprintf("%s+09:00", param.EndTime))
		startTime := convert_start_time.Unix() - int64(model.GmtTimeGap)*60*60
		endTime := convert_end_time.Unix() - int64(model.GmtTimeGap)*60*60

		//param.StartTime = strconv.Itoa(int(startTime))
		//param.EndTime = strconv.Itoa(int(endTime))

		// Make RFC3339 date-time strings
		param.StartTime = time.Unix(startTime, 0).Format(time.RFC3339)[0:19] + ".000000000Z"
		param.EndTime = time.Unix(endTime, 0).Format(time.RFC3339)[0:19] + ".000000000Z"
	}

	/*
	result, err := dao.GetLogsearchDao(service.influxClient, service.databases.LoggingDatabase).GetLogData(param)
	if err == nil{


		for idx, result := range result.Results[0].Series {
			result.Series


		}
		result.Results

	}
	*/

	return dao.GetLogsearchDao(service.influxClient, service.databases.LoggingDatabase).GetLogData(param)
}
package dao

import (
	"fmt"
	"github.com/influxdata/influxdb1-client/v2"
	"kr/paasta/monitoring/paas/model"
	"kr/paasta/monitoring/paas/util"
)

type LogsearchDao struct {
	influxClient client.Client
	measurementName string
}

func GetLogsearchDao(influxClient client.Client, measurementName string) *LogsearchDao {
	return &LogsearchDao{
		influxClient: influxClient,
		measurementName: measurementName,
	}
}

func (dao LogsearchDao) GetLogData(param model.NewLogMessage) (response client.Response, errMsg model.ErrMessage) {
	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = model.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	sqlStr := "select * from \"logging_measurement\""
	if param.Period != "" {
		sqlStr += " where \"time\" <= now() + " + param.Period
	}
	if param.StartTime != "" && param.EndTime != "" {
		sqlStr += " where \"time\" >= '" + param.StartTime + "' and \"time\" <= '" + param.EndTime + "'"
	}
	if param.Id != "" {
		sqlStr += " and \"extradata\" =~ /" + param.Id + "*/"
	}
	sqlStr += " ORDER BY \"time\" DESC;"

	fmt.Println(sqlStr)
	influxQuery := client.Query{
		Command: fmt.Sprint(sqlStr),
		Database: dao.measurementName,
	}
	result, err := dao.influxClient.Query(influxQuery)
	if err != nil {
		errLogMsg = err.Error()
	}

	return util.GetError().CheckError(*result, err)
}
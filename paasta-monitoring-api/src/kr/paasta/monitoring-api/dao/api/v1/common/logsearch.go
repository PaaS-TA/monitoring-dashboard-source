package common

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/jinzhu/gorm"
	models "paasta-monitoring-api/models/api/v1"
)

type LogSearchDao struct {
	DbInfo         *gorm.DB
	InfluxDbClient models.InfluxDbClient
}

func GetLogSearchDao(DbInfo *gorm.DB, InfluxDbClient models.InfluxDbClient) *LogSearchDao {
	return &LogSearchDao{
		DbInfo:         DbInfo,
		InfluxDbClient: InfluxDbClient,
	}
}

func (l *LogSearchDao) GetLogs(params models.Logs) (*client.Response, error) {
	command := "select * from \"logging_measurement\""
	if params.Period != "" {
		command += " where \"time\" <= now() + " + params.Period
	}
	if params.StartTime != "" && params.EndTime != "" {
		command += " where \"time\" >= '" + params.StartTime + "' and \"time\" <= '" + params.EndTime + "'"
	}
	if params.UUID != "" {
		command += " and \"extradata\" =~ /" + params.UUID + "*/"
	}
	if params.Keyword != "" {
		command += " and \"message\" =~ /" + params.Keyword + "/"
	}
	command += " ORDER BY \"time\" DESC limit 100;"
	fmt.Println("GetLogs command======>", command)

	query := client.Query{
		Command:  command,
		Database: l.InfluxDbClient.DbName.LoggingDatabase,
	}
	response, err := l.InfluxDbClient.HttpClient.Query(query)
	if err != nil {
		return response, err
	}
	fmt.Println("GetLogs response======>", response)
	return response, err
}

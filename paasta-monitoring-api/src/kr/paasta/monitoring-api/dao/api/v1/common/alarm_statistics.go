package common

import (
	"github.com/jinzhu/gorm"
	"log"
	models "paasta-monitoring-api/models/api/v1"
)
type AlarmStatisticsDao struct {
	DbInfo *gorm.DB
}

func GetAlarmStatisticsDao(DbInfo *gorm.DB) *AlarmStatisticsDao {
	return &AlarmStatisticsDao {
		DbInfo: DbInfo,
	}
}


func (dao *AlarmStatisticsDao) GetAlarmStatisticsForGraphByTime(param models.AlarmStatisticsParam) ([]map[string]interface{}, error) {
	var countByTimeline []models.CountByTimeline
	extraParams := param.ExtraParams.([]models.AlarmStatisticsCriteriaRequest)
	var response []map[string]interface{}

	for _, v := range extraParams {
		whereRight := `level = '` + v.AlarmLevel + `'`
		if v.Service != "" {
			whereRight += ` AND origin_type = '` + v.Service + `'`
		}
		if v.Resource != "" {
			whereRight += ` AND alarm_type = '` + v.Resource + `'`
		}

		sqlLeft := `
(
WITH RECURSIVE AggregateTable
AS
  (
         SELECT DATE_FORMAT(NOW(), '` + param.DateFormat + `') AS TimelineA
         UNION ALL
         SELECT DATE_FORMAT(DATE_SUB(AggregateTable.TimelineA, INTERVAL 1 ` + param.TimeCriterion + `), '` + param.DateFormat + `') AS TimelineB
         FROM   AggregateTable )
  SELECT   DATE_FORMAT(AggregateTable.TimelineA, '` + param.DateFormat + `') AS Timeline
  FROM     AggregateTable
  WHERE    AggregateTable.TimelineA > DATE_SUB(NOW(), INTERVAL 1 ` + param.Period + `)
  ORDER BY Timeline ASC ) L`

		sqlRight := `
LEFT JOIN
(
         SELECT   DATE_FORMAT(reg_date, '` + param.DateFormat + `') AS Timeline,
                  COUNT(*)                                    AS COUNT
         FROM     alarms
         WHERE    DATE_FORMAT(reg_date, '%Y-%m-%d') > DATE_SUB(NOW(), INTERVAL 1 ` + param.Period + `)
         AND      DATE_FORMAT(reg_date, '%Y-%m-%d') <= NOW()
         AND      ` + whereRight + `
         GROUP BY Timeline
         ORDER BY Timeline ASC ) R ON L.Timeline = R.Timeline`

		results := dao.DbInfo.Debug().Table(sqlLeft).Joins(sqlRight).
			Select("UNIX_TIMESTAMP(L.Timeline) AS timeline, IFNULL(R.Count, 0) AS count").
			Order("timeline ASC").
			Find(&countByTimeline)

		if results.Error != nil {
			log.Println(results.Error)
			return response, results.Error
		}

		tmp := map[string]interface{}{"level": v.Alias, "statistics": countByTimeline}
		response = append(response, tmp)
	}

	return response, nil
}
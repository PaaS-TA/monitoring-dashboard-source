package ap

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	models "paasta-monitoring-api/models/api/v1"
	"time"
)

type ApAlarmDao struct {
	DbInfo *gorm.DB
}

func GetApAlarmDao(DbInfo *gorm.DB) *ApAlarmDao {
	return &ApAlarmDao{
		DbInfo: DbInfo,
	}
}

func (ap *ApAlarmDao) GetAlarmStatus() ([]models.Alarms, error) {
	var response []models.Alarms
	results := ap.DbInfo.Debug().Table("alarms").
		Select("*").
		Find(&response)

	if results.Error != nil {
		fmt.Println(results.Error)
		return response, results.Error
	}

	return response, nil
}

func (ap *ApAlarmDao) GetAlarmPolicy() ([]models.AlarmPolicies, error) {
	var response []models.AlarmPolicies
	results := ap.DbInfo.Debug().Table("alarm_policies").
		Select("*").
		Find(&response)

	if results.Error != nil {
		fmt.Println(results.Error)
		return response, results.Error
	}

	return response, nil
}

func (ap *ApAlarmDao) UpdateAlarmPolicy(request models.AlarmPolicyRequest) error {
	results := ap.DbInfo.Debug().Table("alarm_policies").
		Where("origin_type = ? and alarm_type = ?", request.OriginType, request.AlarmType).
		Updates(map[string]interface{}{
			"warning_threshold":  request.WarningThreshold,
			"critical_threshold": request.CriticalThreshold,
			"repeat_time":        request.RepeatTime,
			"measure_time":       request.MeasureTime,
			"modi_date":          time.Now().UTC().Add(time.Hour * 9),
			"modi_user":          "admin"})

	if results.Error != nil {
		fmt.Println(results.Error)
		return results.Error
	}

	return nil
}

func (ap *ApAlarmDao) UpdateAlarmTarget(request models.AlarmTargetRequest) error {
	results := ap.DbInfo.Debug().Table("alarm_targets").
		Where("origin_type = ?", request.OriginType).
		Updates(map[string]interface{}{
			"mail_address": request.MailAddress,
			"mail_send_yn": request.MailSendYN,
			"modi_date":    time.Now().UTC().Add(time.Hour * 9),
			"modi_user":    "admin"})

	if results.Error != nil {
		fmt.Println(results.Error)
		return results.Error
	}

	return nil
}

func (ap *ApAlarmDao) RegisterSnsAccount(request models.SnsAccountRequest) error {
	results := ap.DbInfo.Debug().Table("alarm_sns").Create(&request)

	if results.Error != nil {
		fmt.Println(results.Error)
		return results.Error
	}

	return nil
}

func (ap *ApAlarmDao) GetSnsAccount() ([]models.AlarmSns, error) {
	var response []models.AlarmSns
	results := ap.DbInfo.Debug().Table("alarm_sns").
		Select("*").
		Find(&response)

	if results.Error != nil {
		fmt.Println(results.Error)
		return response, results.Error
	}

	return response, nil
}

func (ap *ApAlarmDao) DeleteSnsAccount(request models.SnsAccountRequest) error {
	results := ap.DbInfo.Debug().Table("alarm_sns").
		Where("sns_id = ?", request.SnsId).
		Delete(&request)

	if results.Error != nil {
		fmt.Println(results.Error)
		return results.Error
	}

	return nil
}

func (ap *ApAlarmDao) UpdateSnsAccount(request models.SnsAccountRequest) error {
	results := ap.DbInfo.Debug().Table("alarm_sns").
		Where("sns_id = ?", request.SnsId).
		Updates(map[string]interface{}{
			"origin_type": request.OriginType,
			"sns_type":    request.SnsType,
			"sns_id":      request.SnsId,
			"token":       request.Token,
			"expl":        request.Expl,
			"sns_send_yn": request.SnsSendYN,
			"modi_date":   time.Now().UTC().Add(time.Hour * 9),
			"modi_user":   "admin"})

	if results.Error != nil {
		fmt.Println(results.Error)
		return results.Error
	}

	return nil
}

func (ap *ApAlarmDao) CreateAlarmAction(request models.AlarmActionRequest) error {
	results := ap.DbInfo.Debug().Table("alarm_actions").Create(&request)

	if results.Error != nil {
		fmt.Println(results.Error)
		return results.Error
	}

	return nil
}

func (ap *ApAlarmDao) GetAlarmAction() ([]models.AlarmActions, error) {
	var response []models.AlarmActions
	results := ap.DbInfo.Debug().Table("alarm_actions").
		Select("*").
		Find(&response)

	if results.Error != nil {
		fmt.Println(results.Error)
		return response, results.Error
	}

	return response, nil
}

func (ap *ApAlarmDao) UpdateAlarmAction(request models.AlarmActionRequest) error {
	results := ap.DbInfo.Debug().Table("alarm_actions").
		Where("id = ?", request.Id).
		Updates(map[string]interface{}{
			"alarm_id":          request.AlarmId,
			"alarm_action_desc": request.AlarmActionDesc,
			"modi_date":         time.Now().UTC().Add(time.Hour * 9),
			"modi_user":         "admin"})

	if results.Error != nil {
		fmt.Println(results.Error)
		return results.Error
	}

	return nil
}

func (ap *ApAlarmDao) DeleteAlarmAction(request models.AlarmActionRequest) error {
	results := ap.DbInfo.Debug().Table("alarm_actions").
		Where("id = ?", request.Id).
		Delete(&request)

	if results.Error != nil {
		fmt.Println(results.Error)
		return results.Error
	}

	return nil
}

func (ap *ApAlarmDao) GetAlarmStatisticsTotal(c echo.Context) ([]map[string]interface{}, error) {
	request := []models.AlarmStatisticsCriteriaRequest{
		{"Warning", "warning", "", ""},
		{"Critical", "critical", "", ""},
	}

	response, err := GetAlarmStatisticsForGraphByTime(ap, c, request)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (ap *ApAlarmDao) GetAlarmStatisticsService(c echo.Context) ([]map[string]interface{}, error) {
	request := []models.AlarmStatisticsCriteriaRequest{
		{"bos-Warning", "warning", "bos", ""},
		{"bos-Critical", "critical", "bos", ""},
		{"pas-Warning", "warning", "pas", ""},
		{"pas-Critical", "critical", "pas", ""},
		{"con-Warning", "warning", "con", ""},
		{"con-Critical", "critical", "con", ""},
	}

	response, err := GetAlarmStatisticsForGraphByTime(ap, c, request)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (ap *ApAlarmDao) GetAlarmStatisticsResource(c echo.Context) ([]map[string]interface{}, error) {
	request := []models.AlarmStatisticsCriteriaRequest{
		{"cpu-Warning", "warning", "", "cpu"},
		{"cpu-Critical", "critical", "", "cpu"},
		{"memory-Warning", "warning", "", "memory"},
		{"memory-Critical", "critical", "", "memory"},
		{"disk-Warning", "warning", "", "disk"},
		{"disk-Critical", "critical", "", "disk"},
	}

	response, err := GetAlarmStatisticsForGraphByTime(ap, c, request)
	if err != nil {
		return response, err
	}

	return response, nil
}

func GetAlarmStatisticsForGraphByTime(ap *ApAlarmDao, c echo.Context, request []models.AlarmStatisticsCriteriaRequest) ([]map[string]interface{}, error) {
	var period string
	var timeCriterion string
	var dateFormat string

	switch c.QueryParam("period") {
	case "d":
		period = "DAY"
		timeCriterion = "HOUR"
		dateFormat = "%Y-%m-%d %H:00"
	case "w":
		period = "WEEK"
		timeCriterion = "DAY"
		dateFormat = "%Y-%m-%d"
	case "m":
		period = "MONTH"
		timeCriterion = "DAY"
		dateFormat = "%Y-%m-%d"
	case "y":
		period = "YEAR"
		timeCriterion = "MONTH"
		dateFormat = "%Y-%m-01"
	}

	var countByTimeline []models.CountByTimeline
	var response []map[string]interface{}

	for _, v := range request {
		whereRight := `level = '` + v.AlarmLevel + `'`
		if v.Service != "" {
			whereRight += ` AND origin_type = '` + v.Service + `'`
		}
		if v.Resource != "" {
			whereRight += ` AND alarm_type = '` + v.Resource + `'`
		}

		SQLLeft := `
(
  WITH RECURSIVE AggregateTable AS (
    SELECT
      DATE_FORMAT(NOW(), '` + dateFormat + `') AS TimelineA
    UNION ALL
    SELECT
      DATE_FORMAT(DATE_SUB(AggregateTable.TimelineA, INTERVAL 1 ` + timeCriterion + `), '` + dateFormat + `') AS TimelineB
    FROM
      AggregateTable
  )
  SELECT
    DATE_FORMAT(AggregateTable.TimelineA, '` + dateFormat + `') AS Timeline
  FROM
    AggregateTable
  WHERE
    AggregateTable.TimelineA > DATE_SUB(NOW(), INTERVAL 1 ` + period + `)
  ORDER BY Timeline ASC
) L`

		SQLRight := `
LEFT JOIN
(
  SELECT
    DATE_FORMAT(reg_date, '` + dateFormat + `') AS Timeline,
    COUNT(*) AS Count
  FROM
    alarms
  WHERE
    DATE_FORMAT(reg_date, '%Y-%m-%d') > DATE_SUB(NOW(), INTERVAL 1 ` + period + `)
  AND
    DATE_FORMAT(reg_date, '%Y-%m-%d') <= NOW()
  AND
    ` + whereRight + `
  GROUP BY Timeline
  ORDER BY Timeline ASC
) R
ON L.Timeline = R.Timeline`

		results := ap.DbInfo.Debug().Table(SQLLeft).Joins(SQLRight).
			Select("UNIX_TIMESTAMP(L.Timeline) AS timeline, IFNULL(R.Count, 0) AS count").
			Order("timeline ASC").
			Find(&countByTimeline)

		if results.Error != nil {
			fmt.Println(results.Error)
			return response, results.Error
		}

		tmp := map[string]interface{}{"level": v.Alias, "statistics": countByTimeline}
		response = append(response, tmp)
	}

	return response, nil
}

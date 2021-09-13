package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring/paas/model"
	"kr/paasta/monitoring/paas/util"
	"strconv"
)

type AlarmDao struct {
	txn *gorm.DB
}

func GetAlarmDao(txn *gorm.DB) *AlarmDao {
	return &AlarmDao{
		txn: txn,
	}
}

//Dao
func (h *AlarmDao) GetAlarmList(request model.AlarmRequest, txn *gorm.DB) ([]model.AlarmResponse, int, model.ErrMessage) {

	t := []model.AlarmResponse{}

	if request.PagingReq.PageIndex != 0 && request.PagingReq.PageItem != 0 {

		//Page 를 계산한다.
		//Mysql 은 Limit을 제공함. LIMIT: Page당 조회 건수, OffSet: 조회시작 DataRow
		var rowCount int
		var startDataRow int
		endDataRow := request.PagingReq.PageItem * request.PagingReq.PageIndex
		if request.PagingReq.PageIndex == 1 {
			startDataRow = 0
		} else if request.PagingReq.PageIndex > 1 {
			startDataRow = endDataRow - request.PagingReq.PageItem
		}

		var queryWhere = " ( EXISTS ( SELECT id FROM vms WHERE id = A.origin_id  and A.origin_type= 'ias' ) or A.origin_type != 'ias' ) and"

		if request.OriginType != "" {
			queryWhere += " origin_type = '" + request.OriginType + "' and"
		}
		if request.AlarmType != "" {
			queryWhere += " alarm_type = '" + request.AlarmType + "' and"
		}
		if request.Level != "" {
			queryWhere += " level = '" + request.Level + "' and"
		}
		if request.ResolveStatus != "" {
			queryWhere += " resolve_status = '" + request.ResolveStatus + "' and"
		}
		if len(request.SearchDateFrom) > 0 && len(request.SearchDateTo) > 0 {
			queryWhere += " reg_date BETWEEN '" + request.SearchDateFrom + " 00:00:00' AND '" + request.SearchDateTo + " 23:59:59' and"
		}
		//if request.SearchDateFrom != "" && request.SearchDateTo != "" {
		//	//DB에 저장된 시간은 GMT Time기준
		//	//UI에서 요청한 Local Time이 GmtTime Gap보다 9시간 빠르다.
		//	//요청한 일시에서 9시간을 빼어 조회 요청한다.
		//
		//	dateFromUint, _ := strconv.ParseUint(request.SearchDateFrom, 10, 0)
		//	dateToUint, _ := strconv.ParseUint(request.SearchDateTo, 10, 0)
		//
		//	gmtTimeGap := uint64(model.GmtTimeGap)
		//
		//	dateFrom := strconv.FormatUint(dateFromUint-(60*60*gmtTimeGap), 10)
		//	dateTo := strconv.FormatUint(dateToUint-(60*60*gmtTimeGap), 10)
		//
		//	queryWhere += " unix_timestamp(reg_date)*1000 between '" + dateFrom + "' and '" + dateTo + "' and"
		//}
		//조건이 한가지라도 있다면
		if queryWhere != "" {
			queryWhere = queryWhere[0 : len(queryWhere)-3] //and 조건 제거
		}

		status := txn.Debug().Limit(request.PageItem).Table("alarms A").
			Select("id, origin_type, origin_id, alarm_type, level, " +
				"app_yn, app_name, app_index, container_name, alarm_title, " +
				"( " +
				"	CASE " +
				"		WHEN origin_type= 'bos' THEN 'micro-bosh' " +
				"		WHEN origin_type= 'pas' THEN (select name from vms where vms.id = origin_id) " +
				"		WHEN origin_type= 'ias' THEN 'IaaS' " +
				"		ELSE app_name " +
				"	END " +
				") origin_name, " +
				"( " +
				"	CASE " +
				"		WHEN resolve_status= '1' THEN 'Alarm 발생' " +
				"		WHEN resolve_status= '2' THEN 'Alarm 처리중' " +
				"		ELSE 'Alarm 처리완료'  " +
				"	END " +
				") resolve_status_name, " +
				"resolve_status, alarm_message, ip, alarm_cnt, " +
				"reg_date + INTERVAL " + strconv.Itoa(model.GmtTimeGap) +
				" HOUR as reg_date, alarm_send_date + INTERVAL " + strconv.Itoa(model.GmtTimeGap) +
				" HOUR as alarm_send_date, reg_user, 'admin' user_name ").
			Order("reg_date desc").
			Offset(startDataRow).
			Where(queryWhere).
			Find(&t)
		err := util.GetError().DbCheckError(status.Error)

		//status = txn.Debug().Model(&model.Alarm{}).Where(queryWhere).Count(&rowCount)
		status = txn.Debug().Table("alarms A").
			Where(queryWhere).
			Count(&rowCount)

		if err != nil {
			return nil, 0, err
		}
		return t, rowCount, err
	} else {
		return nil, 0, nil
	}
}

func (h *AlarmDao) GetAlarmListCount(request model.AlarmRequest, txn *gorm.DB) (model.AlarmStatusCountResponse, model.ErrMessage) {

	t := model.AlarmStatusCountResponse{}

	var queryWhere = " ( EXISTS ( SELECT id FROM vms WHERE id = A.origin_id  and A.origin_type= 'ias' ) or A.origin_type != 'ias' ) and resolve_status = '" + request.ResolveStatus + "'"

	if len(request.SearchDateFrom) != 0 && len(request.SearchDateTo) != 0 {
		queryWhere += " AND alarm_send_date BETWEEN '" + request.SearchDateFrom + " 00:00:00' AND '" + request.SearchDateTo + " 23:59:59' "
	}

	status := txn.Debug().Table("alarms A").
		Select("count(id) as totalCnt").
		Where(queryWhere).
		Count(&t.TotalCnt)

	err := util.GetError().DbCheckError(status.Error)
	if err != nil {
		return t, err
	}

	return t, err
}

func (h *AlarmDao) GetAlarmResolveStatus(request model.AlarmRequest, txn *gorm.DB) ([]model.AlarmResponse, model.ErrMessage) {

	t := []model.AlarmResponse{}

	var queryWhere = "resolve_status = '" + request.ResolveStatus + "'"

	status := txn.Debug().Table("alarms").
		Select("id, origin_type, origin_id, alarm_type, level, " +
			"app_yn, app_name, app_index, container_name, alarm_title, " +
			"( " +
			"	CASE " +
			"		WHEN origin_type= 'bos' THEN 'micro-bosh' " +
			"		WHEN origin_type= 'pas' THEN (select name from vms where vms.id = origin_id) " +
			"		WHEN origin_type= 'ias' THEN 'IaaS' " +
			"		ELSE app_name " +
			"	END " +
			") origin_name, " +
			"( " +
			"	CASE " +
			"		WHEN resolve_status= '1' THEN 'Alarm 발생' " +
			"		WHEN resolve_status= '2' THEN 'Alarm 처리중' " +
			"		ELSE 'Alarm 처리완료'  " +
			"	END " +
			") resolve_status_name, " +
			"resolve_status, alarm_message, ip, alarm_cnt, " +
			"reg_date + INTERVAL " + strconv.Itoa(model.GmtTimeGap) + " HOUR as reg_date, alarm_send_date, reg_user, 'admin' user_name ").Order("reg_date desc").
		Where(queryWhere).
		Find(&t)
	err := util.GetError().DbCheckError(status.Error)
	if err != nil {
		return nil, err
	}
	return t, err
}

func (h *AlarmDao) GetAlarmDetail(request model.AlarmRequest, txn *gorm.DB) (model.AlarmDetailResponse, model.ErrMessage) {

	t := model.AlarmDetailResponse{}
	status := txn.Debug().Table("alarms").
		Select("id, origin_type, origin_id, alarm_type, level, "+
			"app_yn, app_name, app_index, container_name, alarm_title, "+
			"( "+
			"	CASE "+
			"		WHEN origin_type= 'bos' THEN 'micro-bosh' "+
			"		WHEN origin_type= 'pas' THEN (select name from vms where vms.id = origin_id) "+
			"		WHEN origin_type= 'ias' THEN 'IaaS' " +
			"		ELSE app_name "+
			"	END "+
			") origin_name, "+
			"( "+
			"	CASE "+
			"		WHEN resolve_status= '1' THEN 'Alarm 발생' "+
			"		WHEN resolve_status= '2' THEN 'Alarm 처리중' "+
			"		ELSE 'Alarm 처리완료'  "+
			"	END "+
			") resolve_status_name, "+
			"resolve_status, alarm_message, ip, alarm_cnt, "+
			"reg_date + INTERVAL "+strconv.Itoa(model.GmtTimeGap)+" HOUR as reg_date, alarm_send_date  + INTERVAL "+strconv.Itoa(model.GmtTimeGap)+" HOUR as alarm_send_date, reg_user, 'admin' user_name ").
		Where("id = ?", request.Id).
		Find(&t)
	err := util.GetError().DbCheckError(status.Error)
	if err != nil {
		return t, err
	}
	return t, nil
}

func (h *AlarmDao) GetAlarmsAction(request model.AlarmRequest, txn *gorm.DB) ([]model.AlarmActionResponse, model.ErrMessage) {

	t := []model.AlarmActionResponse{}
	status := txn.Debug().Table("alarm_actions").
		Select("id, alarm_id,  alarm_action_desc, reg_date + INTERVAL "+strconv.Itoa(model.GmtTimeGap)+" HOUR  as reg_date , reg_user, modi_date, modi_user").
		Where("alarm_id = ?", request.Id).
		Find(&t)
	err := util.GetError().DbCheckError(status.Error)
	if err != nil {
		return t, err
	}
	return t, nil
}

func (h *AlarmDao) UpdateAlarm(request model.AlarmRequest, txn *gorm.DB) model.ErrMessage {

	var err model.ErrMessage
	status := txn.Debug().Table("alarms").Where("id = ? ", request.Id).
		Updates(map[string]interface{}{"resolve_status": request.ResolveStatus, "modi_date": util.GetDBCurrentTime(), "modi_user": "system"})
	err = util.GetError().DbCheckError(status.Error)

	if request.ResolveStatus == "3" {
		txn.Debug().Table("alarms").Where("id = ? ", request.Id).
			Updates(map[string]interface{}{"complete_date": util.GetDBCurrentTime(), "complete_user": "system"})
	}

	return err
}

func (h *AlarmDao) CreateAlarmAction(request model.AlarmActionRequest, txn *gorm.DB) model.ErrMessage {

	actionData := model.AlarmAction{AlarmId: request.AlarmId, AlarmActionDesc: request.AlarmActionDesc}

	status := txn.Debug().Create(&actionData)
	err := util.GetError().DbCheckError(status.Error)
	if err != nil {
		return err
	}
	return err
}

func (h *AlarmDao) UpdateAlarmAction(request model.AlarmActionRequest, txn *gorm.DB) model.ErrMessage {

	var err model.ErrMessage
	status := txn.Debug().Table("alarm_actions").Where("id = ? ", request.Id).
		Updates(map[string]interface{}{"alarm_action_desc": request.AlarmActionDesc, "modi_date": util.GetDBCurrentTime()})
	err = util.GetError().DbCheckError(status.Error)
	return err
}

func (h *AlarmDao) DeleteAlarmAction(request model.AlarmActionRequest, txn *gorm.DB) model.ErrMessage {

	var err model.ErrMessage
	status := txn.Debug().Table("alarm_actions").Where("id = ? ", request.Id).Delete(&request)
	err = util.GetError().DbCheckError(status.Error)
	return err
}

func (h *AlarmDao) GetAlarmStat(request model.AlarmStatRequest, txn *gorm.DB) (model.AlarmStatResponse, model.ErrMessage) {

	t := model.AlarmStatResponse{}

	var queryWhere = ""

	period := request.Period
	if period == "d" {
		period = "day"
	} else if period == "w" {
		period = "week"
	} else if period == "m" {
		period = "month"
	} else if period == "y" {
		period = "year"
	}
	queryWhere = " date_sub(reg_date, interval " + strconv.Itoa(model.GmtTimeGap) + " hour) > date_sub(now(), interval " + fmt.Sprint(request.Interval) + " " + period + ") " +
		"	and date_sub(reg_date, interval " + strconv.Itoa(model.GmtTimeGap) + " hour) < date_sub(now(), interval " + fmt.Sprint(request.Interval-1) + " " + period + ")"

	status := txn.Debug().Table("alarms").
		Select(
			//"count(level) as 'total_cnt', " +
			"	ifnull(sum(case when level != '" + model.ALARM_LEVEL_FAIL + "' then 1 else 0 end),0) as 'total_cnt', " +
				"	ifnull(sum(case when level = '" + model.ALARM_LEVEL_WARNING + "' then 1 else 0 end),0) as 'warning_cnt', " +
				"	ifnull(sum(case when level = '" + model.ALARM_LEVEL_CRITICAL + "' then 1 else 0 end),0) as 'critical_cnt', " +
				"	ifnull(sum(case when (level = '" + model.ALARM_LEVEL_WARNING + "' and origin_type = '" + model.ORIGIN_TYPE_PAASTA + "') then 1 else 0 end),0) as 'paasta_warning_cnt', " +
				"	ifnull(sum(case when (level = '" + model.ALARM_LEVEL_CRITICAL + "' and origin_type = '" + model.ORIGIN_TYPE_PAASTA + "') then 1 else 0 end),0) as 'paasta_critical_cnt', " +
				"	ifnull(sum(case when (level = '" + model.ALARM_LEVEL_WARNING + "' and origin_type = '" + model.ORIGIN_TYPE_BOSH + "') then 1 else 0 end),0) as 'bosh_warning_cnt', " +
				"	ifnull(sum(case when (level = '" + model.ALARM_LEVEL_CRITICAL + "' and origin_type = '" + model.ORIGIN_TYPE_BOSH + "') then 1 else 0 end),0) as 'bosh_critical_cnt', " +
				"	ifnull(sum(case when (level = '" + model.ALARM_LEVEL_WARNING + "' and origin_type = '" + model.ORIGIN_TYPE_CONTAINER + "') then 1 else 0 end),0) as 'container_warning_cnt', " +
				"	ifnull(sum(case when (level = '" + model.ALARM_LEVEL_CRITICAL + "' and origin_type = '" + model.ORIGIN_TYPE_CONTAINER + "') then 1 else 0 end),0) as 'container_critical_cnt', " +
				"	ifnull(sum(case when resolve_status = '3' then 1 else 0 end),0) as 'total_resolve_cnt', " +
				"	ifnull(sum(case when (level = '" + model.ALARM_LEVEL_WARNING + "' and resolve_status = '3') then 1 else 0 end),0) as 'warning_resolve_cnt', " +
				"	ifnull(sum(case when (level = '" + model.ALARM_LEVEL_CRITICAL + "' and resolve_status = '3') then 1 else 0 end),0) as 'critical_resolve_cnt', " +
				"	ifnull(sum(case when (level = '" + model.ALARM_LEVEL_WARNING + "' and alarm_type = '" + model.ALARM_TYPE_CPU + "') then 1 else 0 end),0) as 'cpu_warning_cnt', " +
				"	ifnull(sum(case when (level = '" + model.ALARM_LEVEL_CRITICAL + "' and alarm_type = '" + model.ALARM_TYPE_CPU + "') then 1 else 0 end),0) as 'cpu_critical_cnt', " +
				"	ifnull(sum(case when (level = '" + model.ALARM_LEVEL_WARNING + "' and alarm_type = '" + model.ALARM_TYPE_MEMORY + "') then 1 else 0 end),0) as 'memory_warning_cnt', " +
				"	ifnull(sum(case when (level = '" + model.ALARM_LEVEL_CRITICAL + "' and alarm_type = '" + model.ALARM_TYPE_MEMORY + "') then 1 else 0 end),0) as 'memory_critical_cnt', " +
				"	ifnull(sum(case when (level = '" + model.ALARM_LEVEL_WARNING + "' and alarm_type = '" + model.ALARM_TYPE_DISK + "') then 1 else 0 end),0) as 'disk_warning_cnt', " +
				"	ifnull(sum(case when (level = '" + model.ALARM_LEVEL_CRITICAL + "' and alarm_type = '" + model.ALARM_TYPE_DISK + "') then 1 else 0 end),0) as 'disk_critical_cnt' ").
		Where(queryWhere).
		Find(&t)
	err := util.GetError().DbCheckError(status.Error)
	if err != nil {
		return t, err
	}
	return t, nil
}

/* left outer join (criterion time + alarm count per time) */
func (h *AlarmDao) GetAlarmListByPeriod(request model.AlarmStatRequest) ([]model.CountPerTime, model.ErrMessage) {

	if request.Interval == 0 {
		request.Interval = 1
	}

	timeCriterion := "hour"
	period := "day"
	dateFormat := "%Y-%m-%d %H"
	toConvertUnixTimestamp := "L.time"

	switch request.Period {
	case "w":
		timeCriterion = "day"
		period = "week"
		dateFormat = "%Y-%m-%d"
	case "m":
		timeCriterion = "day"
		period = "month"
		dateFormat = "%Y-%m-%d"
	case "y":
		timeCriterion = "month"
		period = "year"
		dateFormat = "%Y-%m"
		toConvertUnixTimestamp = "CONCAT(L.time,'-01')"
	}

	intervalFrom := strconv.Itoa(request.Interval) + " " + period
	intervalTo := strconv.Itoa(request.Interval-1) + " " + period
	timeRangeL := "a.time > date_sub(now(), interval " + intervalFrom + ") and a.time <= date_sub(now(), interval " + intervalTo + ")"
	timeRangeR := "reg_date > date_sub(now(), interval " + intervalFrom + ") and reg_date <= date_sub(now(), interval " + intervalTo + ")"

	if period == "year" {
		timeRangeL = "a.time > date_format(date_sub(now(), interval " + intervalFrom + "), '%Y-%m') and a.time <= date_format(date_sub(now(), interval " + intervalTo + "), '%Y-%m')"
		timeRangeR = "date_format(reg_date, '%Y-%m') > date_format(date_sub(now(), interval " + intervalFrom + "), '%Y-%m') and date_format(reg_date, '%Y-%m') <= date_format(date_sub(now(), interval " + intervalTo + "), '%Y-%m')"
	}

	whereR := timeRangeR + " and level = '" + request.Level + "'"
	if request.Origin != "" {
		whereR += " and origin_type = '" + request.Origin + "'"
	}
	if request.Type != "" {
		whereR += " and alarm_type = '" + request.Type + "'"
	}

	sqlLeft := "(select a.time from (select DATE_FORMAT(date_sub(now(), interval (a.a + (10 * b.a) + (100 * c.a)) " + timeCriterion + "), '" + dateFormat + "') as time " +
		"from (select 0 as a union all select 1 union all select 2 union all select 3 union all select 4 union all select 5 union all select 6 union all select 7 union all select 8 union all select 9) as a " +
		"cross join (select 0 as a union all select 1 union all select 2 union all select 3 union all select 4 union all select 5 union all select 6 union all select 7 union all select 8 union all select 9) as b " +
		"cross join (select 0 as a union all select 1 union all select 2 union all select 3 union all select 4 union all select 5 union all select 6 union all select 7 union all select 8 union all select 9) as c" +
		") a where " + timeRangeL + " order by a.time asc) L"
	sqlRight := "left join (select date_format(reg_date, '" + dateFormat + "') as time, count(*) as count from alarms " +
		"where " + whereR + " group by time order by time asc) R on L.time = R.time"

	var result []model.CountPerTime
	status := h.txn.Debug().Table(sqlLeft).Joins(sqlRight).Select("ROUND(UNIX_TIMESTAMP(" + toConvertUnixTimestamp + " + interval 9 hour)) as time, ifnull(R.count, 0) as count").Order("time asc").Find(&result)
	err := util.GetError().DbCheckError(status.Error)

	return result, err
}

func (h *AlarmDao) GetPaasAlarmRealTimeList() ([]model.AlarmResponse, model.ErrMessage) {

	var list []model.AlarmResponse

	var queryWhere = " ( EXISTS ( SELECT id FROM vms WHERE id = A.origin_id  and A.origin_type= 'ias' ) or A.origin_type != 'ias' ) and resolve_status != '3'"
	status := h.txn.Debug().Table("alarms A").
		Select("id, origin_type, origin_id, alarm_type, level, ip, app_yn, app_name, app_index, container_name, alarm_title, alarm_message, resolve_status,alarm_cnt, reg_date + INTERVAL 9 HOUR as reg_date, reg_user, modi_date, modi_user, alarm_send_date + INTERVAL 9 HOUR as alarm_send_date").
		Where(queryWhere).
		Order("reg_date desc").
		Find(&list)
	err := util.GetError().DbCheckError(status.Error)

	return list, err
}

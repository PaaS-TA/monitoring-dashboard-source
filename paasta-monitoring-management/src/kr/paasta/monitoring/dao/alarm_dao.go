package dao

import (
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring/domain"
	"kr/paasta/monitoring/util"
	"fmt"
	"strconv"
)

type AlarmDao struct {
	txn   *gorm.DB
}

func GetAlarmDao(txn *gorm.DB) *AlarmDao {
	return &AlarmDao{
		txn:   txn,
	}
}

//Dao
func (h *AlarmDao) GetAlarmList(request domain.AlarmRequest, txn *gorm.DB) ([]domain.AlarmResponse, int, domain.ErrMessage) {

	fmt.Println("Get Call Dao=====")
	t := []domain.AlarmResponse{}

	if request.PageIndex != 0 && request.PageItem != 0{

		//Page 를 계산한다.
		//Mysql 은 Limit을 제공함. LIMIT: Page당 조회 건수, OffSet: 조회시작 DataRow
		var rowCount int
		var startDataRow int
		endDataRow    := request.PageItem * request.PageIndex
		if request.PageIndex == 1{
			startDataRow  = 0
		}else if request.PageIndex > 1{
			startDataRow  = endDataRow - request.PageItem
		}
		var queryWhere = " ( EXISTS ( SELECT id FROM vms WHERE id = A.origin_id  and A.origin_type= 'pas' ) or A.origin_type != 'pas' ) and"

		if request.OriginType != ""{
			queryWhere += " origin_type = '" + request.OriginType + "' and"
		}

		if request.AlarmType != ""{
			queryWhere += " alarm_type = '" + request.AlarmType + "' and"
		}

		if request.Level != ""{
			queryWhere += " level = '" + request.Level + "' and"
		}

		if request.ResolveStatus != ""{
			queryWhere += " resolve_status = '" + request.ResolveStatus + "' and"
		}

		if request.SearchDateFrom != "" && request.SearchDateTo != ""{

			//DB에 저장된 시간은 GMT Time기준
			//UI에서 요청한 Local Time이 GmtTime Gap보다 9시간 빠르다.
			//요청한 일시에서 9시간을 빼어 조회 요청한다.

			dateFromUint, _ := strconv.ParseUint(request.SearchDateFrom, 10, 0)
			dateToUint, _ := strconv.ParseUint(request.SearchDateTo, 10, 0)

			gmtTimeGap := uint64(domain.GmtTimeGap)

			dateFrom := strconv.FormatUint(dateFromUint - (60 * 60 * gmtTimeGap), 10)
			dateTo := strconv.FormatUint(dateToUint - (60 * 60 * gmtTimeGap), 10)

			queryWhere += " unix_timestamp(reg_date)*1000 between '" + dateFrom + "' and '" + dateTo + "' and"
		}
		//조건이 한가지라도 있다면
		if queryWhere != "" {
			queryWhere = queryWhere[0:len(queryWhere)-3] //and 조건 제거
		}

		status := txn.Debug().Limit(request.PageItem).Table("alarms A").
			Select("id, origin_type, origin_id, alarm_type, level, " +
			"app_yn, app_name, app_index, container_name, alarm_title, " +
			"( " +
			"	CASE " +
			"		WHEN origin_type= 'bos' THEN 'micro-bosh' " +
			"		WHEN origin_type= 'pas' THEN (select name from vms where vms.id = origin_id) " +
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
			"reg_date + INTERVAL " + strconv.Itoa(domain.GmtTimeGap) +" HOUR as reg_date, alarm_send_date + INTERVAL " + strconv.Itoa(domain.GmtTimeGap) +" HOUR as alarm_send_date, reg_user, 'admin' user_name ").Order("reg_date desc").Offset(startDataRow).
			Where(queryWhere).
		//Order("c")
			Find(&t)
		err := util.GetError().DbCheckError(status.Error)

		//status = txn.Debug().Model(&domain.Alarm{}).Where(queryWhere).Count(&rowCount)
		status = txn.Debug().Table("alarms A").Where(queryWhere).Count(&rowCount)

		if err != nil{
			return nil, 0, err
		}
		return t,  rowCount, err
	}else{
		return nil, 0, nil
	}
}

func (h *AlarmDao) GetAlarmResolveStatus(request domain.AlarmRequest, txn *gorm.DB) ([]domain.AlarmResponse, domain.ErrMessage) {

	t := []domain.AlarmResponse{}

	var queryWhere = "resolve_status = '" + request.ResolveStatus + "'"

	status := txn.Debug().Table("alarms").
		Select("id, origin_type, origin_id, alarm_type, level, " +
		"app_yn, app_name, app_index, container_name, alarm_title, " +
		"( " +
		"	CASE " +
		"		WHEN origin_type= 'bos' THEN 'micro-bosh' " +
		"		WHEN origin_type= 'pas' THEN (select name from vms where vms.id = origin_id) " +
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
		"reg_date + INTERVAL " + strconv.Itoa(domain.GmtTimeGap) +" HOUR as reg_date, alarm_send_date, reg_user, 'admin' user_name ").Order("reg_date desc").
		Where(queryWhere).
		Find(&t)
	err := util.GetError().DbCheckError(status.Error)
	if err != nil{
		return nil, err
	}
	return t,  err
}

func (h *AlarmDao) GetAlarmDetail(request domain.AlarmRequest, txn *gorm.DB) (domain.AlarmDetailResponse, domain.ErrMessage) {

	t := domain.AlarmDetailResponse{}
	status := txn.Debug().Table("alarms").
		Select("id, origin_type, origin_id, alarm_type, level, " +
		"app_yn, app_name, app_index, container_name, alarm_title, " +
		"( " +
		"	CASE " +
		"		WHEN origin_type= 'bos' THEN 'micro-bosh' " +
		"		WHEN origin_type= 'pas' THEN (select name from vms where vms.id = origin_id) " +
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
		"reg_date + INTERVAL " + strconv.Itoa(domain.GmtTimeGap) +" HOUR as reg_date, alarm_send_date  + INTERVAL " + strconv.Itoa(domain.GmtTimeGap) +" HOUR as alarm_send_date, reg_user, 'admin' user_name ").
		Where("id = ?", request.Id).
		Find(&t)
	err := util.GetError().DbCheckError(status.Error)
	if err != nil{
		return t,  err
	}
	return t, nil
}

func (h *AlarmDao) GetAlarmsAction(request domain.AlarmRequest, txn *gorm.DB) ([]domain.AlarmActionResponse, domain.ErrMessage) {

	t := []domain.AlarmActionResponse{}
	status := txn.Debug().Table("alarm_actions").
		Select("id, alarm_id,  alarm_action_desc, reg_date + INTERVAL " + strconv.Itoa(domain.GmtTimeGap) +" HOUR  as reg_date , reg_user, modi_date, modi_user").
		Where("alarm_id = ?", request.Id).
		Find(&t)
	err := util.GetError().DbCheckError(status.Error)
	if err != nil{
		return t, err
	}
	return t, nil
}

func (h *AlarmDao) UpdateAlarm(request domain.AlarmRequest, txn *gorm.DB) domain.ErrMessage {

	var err domain.ErrMessage
	status := txn.Debug().Table("alarms").Where("id = ? " , request.Id ).
		Updates(map[string]interface{}{ "resolve_status": request.ResolveStatus, "modi_date": util.GetDBCurrentTime(), "modi_user":"system"})
	err = util.GetError().DbCheckError(status.Error)
	return err
}

func (h *AlarmDao) CreateAlarmAction(request domain.AlarmActionRequest, txn *gorm.DB) domain.ErrMessage {

	actionData := domain.AlarmAction{AlarmId: request.AlarmId, AlarmActionDesc: request.AlarmActionDesc}

	status := txn.Debug().Create(&actionData)
	err := util.GetError().DbCheckError(status.Error)
	if err != nil{
		return  err
	}
	return  err
}

func (h *AlarmDao) UpdateAlarmAction(request domain.AlarmActionRequest, txn *gorm.DB) domain.ErrMessage {

	var err domain.ErrMessage
	status := txn.Debug().Table("alarm_actions").Where("id = ? ", request.Id).
		Updates(map[string]interface{}{ "alarm_action_desc": request.AlarmActionDesc, "modi_date": util.GetDBCurrentTime() })
	err = util.GetError().DbCheckError(status.Error)
	return err
}

func (h *AlarmDao) DeleteAlarmAction(request domain.AlarmActionRequest, txn *gorm.DB) domain.ErrMessage {

	var err domain.ErrMessage
	status := txn.Debug().Table("alarm_actions").Where("id = ? ", request.Id).Delete(&request)
	err = util.GetError().DbCheckError(status.Error)
	return err
}

func (h *AlarmDao) GetAlarmStat(request domain.AlarmStatRequest, txn *gorm.DB) (domain.AlarmStatResponse, domain.ErrMessage) {

	t := domain.AlarmStatResponse{}

	var queryWhere = ""

	period := request.Period
	if period == "custom" {
		queryWhere = " unix_timestamp(reg_date)*1000 between '" + request.SearchDateFrom + "' and '" + request.SearchDateTo + "' "
	} else {
		period := request.Period
		if period == "d" {
			period = "day"
		} else if period == "w" {
			period = "week"
		} else if period == "m" {
			period = "month"
		}
		queryWhere = " date_sub(reg_date, interval " + strconv.Itoa(domain.GmtTimeGap) +" hour) > date_sub(now(), interval " + fmt.Sprint(request.Interval) + " " + period +") " +
			"	and date_sub(reg_date, interval " + strconv.Itoa(domain.GmtTimeGap) +" hour) < date_sub(now(), interval " + fmt.Sprint(request.Interval-1) + " " + period +")"
	}

	status := txn.Debug().Table("alarms").
		Select("count(level) as 'total_cnt', " +
		"	ifnull(sum(case when level = '" + domain.ALARM_LEVEL_WARNING + "' then 1 else 0 end),0) as 'warning_cnt', " +
		"	ifnull(sum(case when level = '" + domain.ALARM_LEVEL_CRITICAL + "' then 1 else 0 end),0) as 'critical_cnt', " +
		"	ifnull(sum(case when (level = '" + domain.ALARM_LEVEL_WARNING + "' and origin_type = '" + domain.ORIGIN_TYPE_PAASTA + "') then 1 else 0 end),0) as 'paasta_warning_cnt', " +
		"	ifnull(sum(case when (level = '" + domain.ALARM_LEVEL_CRITICAL + "' and origin_type = '" + domain.ORIGIN_TYPE_PAASTA + "') then 1 else 0 end),0) as 'paasta_critical_cnt', " +
		"	ifnull(sum(case when (level = '" + domain.ALARM_LEVEL_WARNING + "' and origin_type = '" + domain.ORIGIN_TYPE_BOSH + "') then 1 else 0 end),0) as 'bosh_warning_cnt', " +
		"	ifnull(sum(case when (level = '" + domain.ALARM_LEVEL_CRITICAL + "' and origin_type = '" + domain.ORIGIN_TYPE_BOSH + "') then 1 else 0 end),0) as 'bosh_critical_cnt', " +
		"	ifnull(sum(case when (level = '" + domain.ALARM_LEVEL_WARNING + "' and origin_type = '" + domain.ORIGIN_TYPE_CONTAINER + "') then 1 else 0 end),0) as 'container_warning_cnt', " +
		"	ifnull(sum(case when (level = '" + domain.ALARM_LEVEL_CRITICAL + "' and origin_type = '" + domain.ORIGIN_TYPE_CONTAINER + "') then 1 else 0 end),0) as 'container_critical_cnt', " +
		"	ifnull(sum(case when resolve_status = '3' then 1 else 0 end),0) as 'total_resolve_cnt', " +
		"	ifnull(sum(case when (level = '" + domain.ALARM_LEVEL_WARNING + "' and resolve_status = '3') then 1 else 0 end),0) as 'warning_resolve_cnt', " +
		"	ifnull(sum(case when (level = '" + domain.ALARM_LEVEL_CRITICAL + "' and resolve_status = '3') then 1 else 0 end),0) as 'critical_resolve_cnt', " +
		"	ifnull(sum(case when (level = '" + domain.ALARM_LEVEL_WARNING + "' and alarm_type = '" + domain.ALARM_TYPE_CPU + "') then 1 else 0 end),0) as 'cpu_warning_cnt', " +
		"	ifnull(sum(case when (level = '" + domain.ALARM_LEVEL_CRITICAL + "' and alarm_type = '" + domain.ALARM_TYPE_CPU + "') then 1 else 0 end),0) as 'cpu_critical_cnt', " +
		"	ifnull(sum(case when (level = '" + domain.ALARM_LEVEL_WARNING + "' and alarm_type = '" + domain.ALARM_TYPE_MEMORY + "') then 1 else 0 end),0) as 'memory_warning_cnt', " +
		"	ifnull(sum(case when (level = '" + domain.ALARM_LEVEL_CRITICAL + "' and alarm_type = '" + domain.ALARM_TYPE_MEMORY + "') then 1 else 0 end),0) as 'memory_critical_cnt', " +
		"	ifnull(sum(case when (level = '" + domain.ALARM_LEVEL_WARNING + "' and alarm_type = '" + domain.ALARM_TYPE_DISK + "') then 1 else 0 end),0) as 'disk_warning_cnt', " +
		"	ifnull(sum(case when (level = '" + domain.ALARM_LEVEL_CRITICAL + "' and alarm_type = '" + domain.ALARM_TYPE_DISK + "') then 1 else 0 end),0) as 'disk_critical_cnt' ").
		Where(queryWhere).
		Find(&t)
	err := util.GetError().DbCheckError(status.Error)
	if err != nil{
		return t, err
	}
	return t,  nil
}

package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring/paas/model"
	"kr/paasta/monitoring/paas/util"
	"strconv"
)

type AppDao struct {
	txn *gorm.DB
}

func GetAppDao(txn *gorm.DB) *AppDao {
	return &AppDao{
		txn: txn,
	}
}

func (b *AppDao) UpdatePaasAppAutoScalingPolicy(request model.AppAutoscalingPolicy) string {

	req := model.AppAutoScalingPolicies{
		AppGuid:             request.AppGuid,
		InstanceMinCnt:      uint(request.InstanceMinCnt),
		InstanceMaxCnt:      uint(request.InstanceMaxCnt),
		CpuMinThreshold:     uint(request.CpuMinThreshold),
		CpuMaxThreshold:     uint(request.CpuMaxThreshold),
		MemoryMinThreshold:  uint(request.MemoryMinThreshold),
		MemoryMaxThreshold:  uint(request.MemoryMaxThreshold),
		InstanceScalingUnit: uint(request.InstanceScalingUnit),
		MeasureTimeSec:      uint(request.MeasureTimeSec),
		AutoScalingOutYn:    request.AutoScalingOutYn,
		AutoScalingInYn:     request.AutoScalingInYn,
		AutoScalingCpuYn:    request.AutoScalingCpuYn,
		AutoScalingMemoryYn: request.AutoScalingMemoryYn,
		RegDate:             util.GetDBCurrentTime(),
		RegUser:             "system",
		ModiDate:            util.GetDBCurrentTime(),
		ModiUser:            "system",
	}

	updateQuery := "on duplicate key update " +
		"instance_min_cnt ='" + strconv.Itoa(request.InstanceMinCnt) + "'," +
		"instance_max_cnt ='" + strconv.Itoa(request.InstanceMaxCnt) + "'," +
		"cpu_min_threshold ='" + strconv.Itoa(request.CpuMinThreshold) + "'," +
		"cpu_max_threshold ='" + strconv.Itoa(request.CpuMaxThreshold) + "'," +
		"memory_min_threshold ='" + strconv.Itoa(request.MemoryMinThreshold) + "'," +
		"memory_max_threshold ='" + strconv.Itoa(request.MemoryMaxThreshold) + "'," +
		"instance_scaling_unit ='" + strconv.Itoa(request.InstanceScalingUnit) + "'," +
		"measure_time_sec ='" + strconv.Itoa(request.MeasureTimeSec) + "'," +
		"auto_scaling_out_yn ='" + request.AutoScalingOutYn + "'," +
		"auto_scaling_in_yn ='" + request.AutoScalingInYn + "'," +
		"auto_scaling_cpu_yn ='" + request.AutoScalingCpuYn + "'," +
		"auto_scaling_memory_yn ='" + request.AutoScalingMemoryYn + "'," +
		"modi_date = now(), modi_user = 'system'"

	status := b.txn.Debug().Table("app_auto_scaling_policies").
		Set("gorm:insert_option", updateQuery).Create(&req)

	if status.Error != nil {
		return status.Error.Error()
	} else {
		return ""
	}
}

func (b *AppDao) GetPaasAppAutoScalingPolicy(request model.AppAlarmReq) (model.AppAutoscalingPolicy, model.ErrMessage) {
	t := model.AppAutoscalingPolicy{}

	status := b.txn.Debug().Table("app_auto_scaling_policies").
		Select("app_guid, instance_min_cnt, instance_max_cnt, cpu_min_threshold, cpu_max_threshold, "+
			"memory_min_threshold, memory_max_threshold, instance_scaling_unit, measure_time_sec, "+
			"auto_scaling_out_yn, auto_scaling_in_yn, auto_scaling_cpu_yn, auto_scaling_memory_yn ").
		Where("app_guid = ? ", request.AppGuid).
		Find(&t)
	err := util.GetError().DbCheckError(status.Error)

	return t, err

}

func (b *AppDao) UpdatePaasAppPolicyInfo(request model.AppAlarmPolicy) string {

	req := model.AppAlarmPolicies{
		AppGuid:                 request.AppGuid,
		CpuWarningThreshold:     uint(request.CpuWarningThreshold),
		CpuCriticalThreshold:    uint(request.CpuCriticalThreshold),
		MemoryWarningThreshold:  uint(request.MemoryWarningThreshold),
		MemoryCriticalThreshold: uint(request.MemoryCriticalThreshold),
		MeasureTimeSec:          uint(request.MeasureTimeSec),
		Email:                   request.Email,
		EmailSendYn:             request.EmailSendYn,
		AlarmUseYn:              request.AlarmUseYn,
		RegDate:                 util.GetDBCurrentTime(),
		RegUser:                 "system",
		ModiDate:                util.GetDBCurrentTime(),
		ModiUser:                "system",
	}

	updateQuery := "on duplicate key update " +
		"cpu_warning_threshold ='" + strconv.Itoa(request.CpuWarningThreshold) + "'," +
		"cpu_critical_threshold ='" + strconv.Itoa(request.CpuCriticalThreshold) + "'," +
		"memory_warning_threshold ='" + strconv.Itoa(request.MemoryWarningThreshold) + "'," +
		"memory_critical_threshold ='" + strconv.Itoa(request.MemoryCriticalThreshold) + "'," +
		"measure_time_sec ='" + strconv.Itoa(request.MeasureTimeSec) + "'," +
		"email ='" + request.Email + "'," +
		"email_send_yn ='" + request.EmailSendYn + "'," +
		"alarm_use_yn ='" + request.AlarmUseYn + "'," + "modi_date = now(), modi_user = 'system'"

	status := b.txn.Debug().Table("app_alarm_policies").
		Set("gorm:insert_option", updateQuery).Create(&req)

	if status.Error != nil {
		return status.Error.Error()
	} else {
		return ""
	}
}

func (b *AppDao) GetPaasAppPolicyInfo(request model.AppAlarmReq) (model.AppAlarmPolicy, model.ErrMessage) {
	t := model.AppAlarmPolicy{}

	status := b.txn.Debug().Table("app_alarm_policies").
		Select("app_guid, cpu_warning_threshold, cpu_critical_threshold, memory_warning_threshold, memory_critical_threshold, "+
			"measure_time_sec, email, email_send_yn, alarm_use_yn").
		Where("app_guid = ? ", request.AppGuid).
		Find(&t)
	err := util.GetError().DbCheckError(status.Error)

	return t, err

}

func (b *AppDao) GetPaasAppAlarmList(request model.AppAlarmReq) ([]model.AppAlarm, int, model.ErrMessage) {

	t := []model.AppAlarm{}

	if request.PageIndex != 0 && request.PageItems != 0 {

		//Page 를 계산한다.
		//Mysql 은 Limit을 제공함. LIMIT: Page당 조회 건수, OffSet: 조회시작 DataRow
		var rowCount int
		var startDataRow int
		endDataRow := request.PageItems * request.PageIndex
		if request.PageIndex == 1 {
			startDataRow = 0
		} else if request.PageIndex > 1 {
			startDataRow = endDataRow - request.PageItems
		}

		var queryWhere = " app_guid = '" + request.AppGuid + "' and"

		if request.ResourceType != "" {
			queryWhere += " resource_type = '" + request.ResourceType + "' and"
		}
		if request.AlarmLevel != "" {
			queryWhere += " alarm_level = '" + request.AlarmLevel + "' and"
		}

		if request.SearchDateFrom != "" && request.SearchDateTo != "" {
			//DB에 저장된 시간은 GMT Time기준
			//UI에서 요청한 Local Time이 GmtTime Gap보다 9시간 빠르다.
			//요청한 일시에서 9시간을 빼어 조회 요청한다.

			dateFromUint, _ := strconv.ParseUint(request.SearchDateFrom, 10, 0)
			dateToUint, _ := strconv.ParseUint(request.SearchDateTo, 10, 0)

			gmtTimeGap := uint64(model.GmtTimeGap)

			dateFrom := strconv.FormatUint(dateFromUint-(60*60*gmtTimeGap), 10)
			dateTo := strconv.FormatUint(dateToUint-(60*60*gmtTimeGap), 10)

			queryWhere += " unix_timestamp(reg_date)*1000 between '" + dateFrom + "' and '" + dateTo + "' and"
		}
		//조건이 한가지라도 있다면
		if queryWhere != "" {
			queryWhere = queryWhere[0 : len(queryWhere)-3] //and 조건 제거
		}

		status := b.txn.Debug().Limit(request.PageItems).Table("app_alarm_histories").
			Select("alarm_id, app_guid, app_idx, app_name, resource_type, " +
				"alarm_level, alarm_title, alarm_message, " +
				"reg_date + INTERVAL " + strconv.Itoa(model.GmtTimeGap) + " HOUR as reg_date").
			Order("reg_date desc").
			Offset(startDataRow).
			Where(queryWhere).
			Find(&t)
		err := util.GetError().DbCheckError(status.Error)

		status = b.txn.Debug().Table("app_alarm_histories").
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

func (b *AppDao) DeletePaasAppPolicy(guid string) string {

	req := model.AppAlarmReq{
		AppGuid: guid,
	}
	fmt.Println("===== DeletePaasAppAlarmHistory =====")
	status := b.txn.Debug().Table("app_alarm_histories").
		Where("app_guid = ? ", req.AppGuid).Delete(&req)

	fmt.Println("===== DeletePaasAppPolicy =====")
	status = b.txn.Debug().Table("app_alarm_policies").
		Where("app_guid = ? ", req.AppGuid).Delete(&req)

	fmt.Println("===== DeletePaasAppAutoScaling Policy =====")
	status = b.txn.Debug().Table("app_auto_scaling_policies").
		Where("app_guid = ? ", req.AppGuid).Delete(&req)

	if status.Error != nil {
		return status.Error.Error()
	} else {
		return ""
	}
}

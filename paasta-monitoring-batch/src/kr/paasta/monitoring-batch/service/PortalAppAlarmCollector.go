package service

import (
	"monitoring-batch/dao"
	"monitoring-batch/model"
	"monitoring-batch/model/base"
	"monitoring-batch/util"
	"fmt"
	"reflect"
	"time"
	"strconv"
	"sync"
	sender "monitoring-batch/alarm"
)

type PortalAppAlarmStruct struct {
	b *BackendServices
}

func PortalAppAlarm(backendServices *BackendServices) *PortalAppAlarmStruct {
	return &PortalAppAlarmStruct{
		b: backendServices,
	}
}

/*
	참고-01. ContainerInterface + ResourceType 이 동일하면 하나의 알람 건이다.
	참고-02. 종결되었다가 다시 발생하면 ContainerInterface 와 ResourceType 이 동일하더라도 신규 알람으로 생성한다.
	참고-03. 종결되지않고 알람등급이 변경된 건은 이전 알람과 동일 건으로 취급(과다한 알람전송 방지)
	참고-04. GUID:알람건수 = 1:N
	참고-05. 종결되지 않으면 10분마다 이메일 재발송한다.
*/
func (p *PortalAppAlarmStruct) PortalAppAlarmCollect(){

	//App 별 알람 정책 조회
	listAlarmPolicy := p.getAppAlarmPolicy()

	var listNowAlarm []model.AppAlarmHistory	//현재 발생 알람 모델 변수

	//App 알람 정책 별 사용량 조회
	var wg sync.WaitGroup
	wg.Add(len(listAlarmPolicy))
	for _, policy := range listAlarmPolicy {
		go func(wg *sync.WaitGroup, policy model.AppAlarmPolicy) {
			defer wg.Done()

			//InfluxDB 통해 앱GUID 별 컨테이너(복수) 정보(container_interface 등) 획득
			appInfo := p.GetAppInfo(policy.AppGuid)

			//리소스 사용량 SET & 알람 대상 추출
			p.setResourceUsage(&appInfo, policy, &listNowAlarm)

		}(&wg, policy)

	}
	wg.Wait()

	//1. 이전 미종결알람리스트와 현재알람리스트를 비교하여 현재알람리스트에 없는 이전알람건은 종결되었음을 의미하므로 종결여부를 'Y'로 업데이트.
	//2. 계속발생건 Update
	//3. 신규발생건 Insert
	p.updateDbAlarmHistory(listNowAlarm)

	//알람발송대상(신규건+재발송건) 조회
	listSendTargetAlarm := p.getSendTargetAppAlarm()

	//알람발송(redis, email)
	p.SendAlarm(listSendTargetAlarm)
}


func (p *PortalAppAlarmStruct) SendAlarm(listSendTargetAlarm []model.SendTargetAppAlarmHistory) {

	var wg sync.WaitGroup
	wg.Add(len(listSendTargetAlarm))
	for _, alarm := range listSendTargetAlarm {
		go func(wg *sync.WaitGroup, alarm model.SendTargetAppAlarmHistory) {
			defer wg.Done()

			isNew := false
			if alarm.AlarmSendDate.IsZero() {
				isNew = true
			}
			alarm.AlarmSendDate = time.Now()

			//redis(HMSET + LPUSH)
			p.setRedisAppAlarm(alarm, isNew)

			//이메일 발송
			if alarm.EmailSendYn == "Y" {
				p.sendEmailAppAlarm(alarm)
			}

			//알람발송시각 DB 업데이트
			dbErr := dao.GetPortalAppAlarmDao(p.b.MonitoringDbClient, p.b.Influxclient, p.b.InfluxConfig.ContainerDatabase).
				UpdateAlarmSendDate(alarm)
			if dbErr != nil {
				fmt.Errorf(">>>>> dbErr:%v", dbErr)
			}
		}(&wg, alarm)
	}
	wg.Wait()
}

func (p *PortalAppAlarmStruct) sendEmailAppAlarm(alarm model.SendTargetAppAlarmHistory) {

	var alarmConfig model.AlarmConfig
	alarmConfig.SmtpHost = p.b.MailConfig.SmtpHost
	alarmConfig.Port = p.b.MailConfig.Port
	alarmConfig.MailSender = p.b.MailConfig.MailSender
	alarmConfig.SenderPwd = p.b.MailConfig.SenderPwd
	alarmConfig.MailResource = p.b.MailConfig.ResouceUrl
	alarmConfig.MailReceiver = p.b.MailConfig.MailReceiver
	alarmConfig.AlarmSend = p.b.MailConfig.AlarmSend
	alarmConfig.MailTlsSend = p.b.MailConfig.MailTlsSend

	receiver := []model.AlarmChannelInfoResp{{Email:alarm.Email}}
	body := model.MailContents{
		Status: alarm.AlarmLevel,
		Message: alarm.AlarmMessage,
		//StatusDetail: "",
		ServerName: alarm.AppName+"["+strconv.Itoa(int(alarm.AppIdx))+"]",
		AlarmDate: util.TimeToGeneralFormat(alarm.RegDate),
		ElapseTime: strconv.FormatFloat(time.Now().Sub(alarm.RegDate).Minutes(), 'f', 0, 64),
	}

	sender.GetAlarmService(alarmConfig).SendMail(alarm.AlarmTitle, body, receiver)
}

func (p *PortalAppAlarmStruct) setRedisAppAlarm(alarm model.SendTargetAppAlarmHistory, isNew bool) {

	//HMSET(앱알람데이터등록)
	keyAlarm := "appAlarm:" + alarm.AppGuid + ":" + strconv.Itoa(int(alarm.AlarmId))
	statusCmd := p.b.RedisClient.HMSet(keyAlarm, map[string]interface{}{
		"alarmId": strconv.Itoa(int(alarm.AlarmId)),
		"appGuid": alarm.AppGuid,
		"appIdx": strconv.Itoa(int(alarm.AppIdx)),
		"appName": alarm.AppName,
		"resourceType": alarm.ResourceType,
		"alarmLevel": alarm.AlarmLevel,
		"alarmTitle": alarm.AlarmTitle,
		"alarmMessage": alarm.AlarmMessage,
		"alarmSendDate": util.TimeToGeneralFormat(alarm.AlarmSendDate),
		"regDate": util.TimeToGeneralFormat(alarm.RegDate),
	})
	fmt.Println(">>>>> redis HMSET StatusCmd:", statusCmd.String())

	if isNew {
		//LPUSH(미확인앱알람리스트)
		keyAlarmList := "appAlarmList:" + alarm.AppGuid
		intCmd := p.b.RedisClient.LPush(keyAlarmList, keyAlarm)
		fmt.Println(">>>>> redis LPUSH IntCmd:", intCmd.String())
	}
}

func (p *PortalAppAlarmStruct) updateDbAlarmHistory(listNow []model.AppAlarmHistory) {

	//종결되지 않은 알람 이력 조회
	listPre := p.getNotTerminatedAppAlarm()

	var listTerminated []model.AppAlarmHistory
	var listUpdated []model.AppAlarmHistory
	var wg sync.WaitGroup
	var mutexUpdated, mutexTerminated sync.Mutex

	//이전 미종결알람들을 종결건과 계속발생건으로 분류
	wg.Add(len(listPre))
	for _, pre := range listPre {
		go func(wg *sync.WaitGroup, pre model.AppAlarmHistory) {
			defer wg.Done()
			isTerminated := true
			for _, now := range listNow {
				if pre.ContainerInterface == now.ContainerInterface && pre.ResourceType == now.ResourceType {
					isTerminated = false
					mutexUpdated.Lock()
					listUpdated = append(listUpdated, now)
					mutexUpdated.Unlock()
					break
				}
			}
			if isTerminated {
				pre.TerminateYn = "Y"
				mutexTerminated.Lock()
				listTerminated = append(listTerminated, pre)
				mutexTerminated.Unlock()
			}
		}(&wg, pre)
	}
	wg.Wait()

	//종결건 DB 업데이트
	if len(listTerminated) > 0 {
		dbErr := dao.GetPortalAppAlarmDao(p.b.MonitoringDbClient, p.b.Influxclient, p.b.InfluxConfig.ContainerDatabase).
			UpdateTerminatedAlarm(listTerminated)
		if dbErr != nil {
			fmt.Errorf(">>>>> dbErr:%v", dbErr)
		}
	}

	//계속발생건 DB 업데이트
	wg.Add(len(listUpdated))
	for _, updated := range listUpdated {
		go func(wg *sync.WaitGroup, updated model.AppAlarmHistory) {
			defer wg.Done()
			dbErr := dao.GetPortalAppAlarmDao(p.b.MonitoringDbClient, p.b.Influxclient, p.b.InfluxConfig.ContainerDatabase).
				UpdateContinuousAppAlarm(updated)
			if dbErr != nil {
				fmt.Errorf(">>>>> dbErr:%v", dbErr)
			}
		}(&wg, updated)
	}
	wg.Wait()

	//신규발생건 DB Insert
	wg.Add(len(listNow))
	for _, now := range listNow {
		go func(wg *sync.WaitGroup, now model.AppAlarmHistory) {
			defer wg.Done()
			isNew := true
			for _, updated := range listUpdated {
				if now.ContainerInterface == updated.ContainerInterface && now.ResourceType == updated.ResourceType {
					isNew = false
					break
				}
			}
			if isNew {
				dbErr := dao.GetPortalAppAlarmDao(p.b.MonitoringDbClient, p.b.Influxclient, p.b.InfluxConfig.ContainerDatabase).
					InsertNewAppAlarm(now)
				if dbErr != nil {
					fmt.Errorf(">>>>> dbErr:%v", dbErr)
				}
			}
		}(&wg, now)
	}
	wg.Wait()
}

//자원 사용률 조회 후 임계치와 비교하여 초과 시 알람데이터 생성하여 리스트에 추가
func (p *PortalAppAlarmStruct) setResourceUsage(appInfo *model.ApplicationInfo, policy model.AppAlarmPolicy, listAlarm *[]model.AppAlarmHistory) {

	for _, container := range appInfo.ApplicationContainerInfo {
		//CPU 사용률 체크
		container.CpuUsage = p.GetContainerCpuUsage(container, policy.MeasureTimeSec)
		if container.CpuUsage > float64(policy.CpuCriticalThreshold) {
			*listAlarm = append(*listAlarm, generateAlarmData(container, base.ALARM_TYPE_CPU, base.ALARM_LEVEL_CRITICAL, policy.CpuCriticalThreshold))
		} else if container.CpuUsage > float64(policy.CpuWarningThreshold) {
			*listAlarm = append(*listAlarm, generateAlarmData(container, base.ALARM_TYPE_CPU, base.ALARM_LEVEL_WARNING, policy.CpuWarningThreshold))
		}

		//MEMORY 사용률 체크
		container.MemoryUsage = p.GetContainerMemoryUsage(container, policy.MeasureTimeSec)
		if container.MemoryUsage > float64(policy.MemoryCriticalThreshold) {
			*listAlarm = append(*listAlarm, generateAlarmData(container, base.ALARM_TYPE_MEMORY, base.ALARM_LEVEL_CRITICAL, policy.MemoryCriticalThreshold))
		} else if container.MemoryUsage > float64(policy.MemoryWarningThreshold) {
			*listAlarm = append(*listAlarm, generateAlarmData(container, base.ALARM_TYPE_MEMORY, base.ALARM_LEVEL_WARNING, policy.MemoryWarningThreshold))
		}
	}
}

func generateAlarmData(container model.ApplicationContainerInfo, resource string, level string, threshold uint) model.AppAlarmHistory{

	var appAlarm model.AppAlarmHistory
	appAlarm.AppGuid = container.ApplicationId
	appIdx, _ := strconv.Atoi(container.ApplicationIndex)
	appAlarm.AppIdx = uint(appIdx)
	appAlarm.ResourceType = resource
	appAlarm.AlarmLevel = level
	appAlarm.AppName = container.ApplicationName
	appAlarm.CellIp = container.CellIp
	appAlarm.ContainerId = container.ContainerId
	appAlarm.ContainerInterface = container.ContainerInterface
	appAlarm.TerminateYn = "N"
	title := "%s[%d] App의 %s 상태 - %s"
	appAlarm.AlarmTitle = fmt.Sprintf(title, container.ApplicationName, appIdx, resource, level)
	var usage float64
	if resource == base.ALARM_TYPE_CPU {
		usage = container.CpuUsage
	} else {
		usage = container.MemoryUsage
	}
	message := "%s[%d] App의 %s 사용률이 임계치 [%d]%%를 초과했습니다.\n현재 사용률은 [%s]%%입니다."
	appAlarm.AlarmMessage = fmt.Sprintf(message, container.ApplicationName, appIdx, resource, int(threshold), util.FloattostrDigit2(usage))

	return appAlarm
}

func (p *PortalAppAlarmStruct) GetContainerCpuUsage(container model.ApplicationContainerInfo, measureTimeSec uint) float64 {

	request := model.ZonesReq{
		MeasureTimeList: []model.AlarmItemMeasureTime{{Item:base.ALARM_TYPE_CPU, MeasureTime: int(measureTimeSec)}},
		ContainerName: container.ContainerInterface,
		MetricDatabase: p.b.InfluxConfig.ContainerDatabase,
	}
	resp, err := dao.GetContainerAlarmDao(p.b.Influxclient).GetContainerCpuUsage(request)
	if err != nil {
		fmt.Errorf(">>>>> err:%v", err)
	}
	mapConverted, _ := util.GetResponseConverter().InfluxConverter(resp, request.ContainerName)
	cpuUsage := util.GetDataFloatFromInterfaceSingle(mapConverted)
	return cpuUsage
}

func (p *PortalAppAlarmStruct) GetContainerMemoryUsage(container model.ApplicationContainerInfo, measureTimeSec uint) float64 {

	request := model.ZonesReq{
		MeasureTimeList: []model.AlarmItemMeasureTime{{Item:base.ALARM_TYPE_MEMORY, MeasureTime: int(measureTimeSec)}},
		ContainerName: container.ContainerInterface,
		MetricDatabase: p.b.InfluxConfig.ContainerDatabase,
	}
	resp, err := dao.GetContainerAlarmDao(p.b.Influxclient).GetContainerMemoryUsage(request)
	if err != nil {
		fmt.Errorf(">>>>> err:%v", err)
	}
	mapConverted, _ := util.GetResponseConverter().InfluxConverter(resp, request.ContainerName)
	memoryUsage := util.GetDataFloatFromInterfaceSingle(mapConverted)
	return memoryUsage
}

func (p *PortalAppAlarmStruct) getNotTerminatedAppAlarm() []model.AppAlarmHistory {

	listNotTerminatedAppAlarm, err := dao.GetPortalAppAlarmDao(p.b.MonitoringDbClient, p.b.Influxclient, p.b.InfluxConfig.ContainerDatabase).
		GetNotTerminatedAppAlarm()
	if err != nil {
		fmt.Errorf(">>>>> error:%v", err)
	}

	return listNotTerminatedAppAlarm
}

func (p *PortalAppAlarmStruct) getAppAlarmPolicy() []model.AppAlarmPolicy {

	listAppAlarmPolicy, err := dao.GetPortalAppAlarmDao(p.b.MonitoringDbClient, p.b.Influxclient, p.b.InfluxConfig.ContainerDatabase).
		GetAppAlarmPolicy()
	if err != nil {
		fmt.Errorf(">>>>> error:%v", err)
	}

	return listAppAlarmPolicy
}

func (p *PortalAppAlarmStruct) GetAppInfo(appGuid string) model.ApplicationInfo {

	resp, err := dao.GetPortalAppAlarmDao(p.b.MonitoringDbClient, p.b.Influxclient, p.b.InfluxConfig.ContainerDatabase).
		GetAppInfo(appGuid)
	if err != nil {
		fmt.Errorf(">>>>> error:%v", err)
	}

	appInfo := model.ApplicationInfo{ApplicationId:appGuid}

	fmt.Println(">>>>> GetAppInfo resp :", resp)

	for _, v := range resp.Results[0].Series{

		var appContainerInfo model.ApplicationContainerInfo
		appContainerInfo.ApplicationId = reflect.ValueOf(v.Values[0][1]).String()
		appContainerInfo.ApplicationName = reflect.ValueOf(v.Values[0][2]).String()
		appContainerInfo.ApplicationIndex = reflect.ValueOf(v.Values[0][3]).String()
		appContainerInfo.CellIp = reflect.ValueOf(v.Values[0][4]).String()
		appContainerInfo.ContainerId = reflect.ValueOf(v.Values[0][5]).String()
		appContainerInfo.ContainerInterface = reflect.ValueOf(v.Values[0][6]).String()

		appInfo.ApplicationContainerInfo = append(appInfo.ApplicationContainerInfo, appContainerInfo)
	}

	return appInfo
}

func (p *PortalAppAlarmStruct) getSendTargetAppAlarm() []model.SendTargetAppAlarmHistory {

	interval, _ := strconv.Atoi(p.b.config["user.portal.alarm.interval"])
	listSendTargetAlarm, err := dao.GetPortalAppAlarmDao(p.b.MonitoringDbClient, p.b.Influxclient, p.b.InfluxConfig.ContainerDatabase).
		GetSendTargetAppAlarm(interval)
	if err != nil {
		fmt.Errorf(">>>>> error:%v", err)
	}

	return listSendTargetAlarm
}


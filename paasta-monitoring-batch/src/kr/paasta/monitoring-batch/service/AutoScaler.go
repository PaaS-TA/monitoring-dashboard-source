package service

import (
	"fmt"
	"kr/paasta/monitoring-batch/model/base"
	"kr/paasta/monitoring-batch/model"
	"kr/paasta/monitoring-batch/util"
	"kr/paasta/monitoring-batch/dao"
	"sync"
	"strconv"
	"github.com/cloudfoundry-community/go-cfclient"
)

type AutoScalerStruct struct {
	b *BackendServices
	p *PortalAppAlarmStruct
}

func AutoScaler(backendServices *BackendServices) *AutoScalerStruct {
	return &AutoScalerStruct{
		b: backendServices,
	}
}

func (a *AutoScalerStruct) AutoScale(){

	a.p = PortalAppAlarm(a.b)

	//AutoScaling 정책 조회
	listAutoScalePolicy, err := dao.AutoScalerDao(a.b.MonitoringDbClient, a.b.Influxclient, a.b.InfluxConfig.ContainerDatabase).
		GetAutoScalePolicy()
	if err != nil {
		fmt.Errorf(">>>>> error:%v", err)
		return
	}
	fmt.Println(">>>>> AUTO_SCALING_POLICY:", listAutoScalePolicy)

	var listAutoScaleTarget []model.AutoScaleTarget

	//App 알람 정책 별 사용량 조회
	var wg sync.WaitGroup
	wg.Add(len(listAutoScalePolicy))
	for _, policy := range listAutoScalePolicy {
		go func(wg *sync.WaitGroup, policy model.AppAutoScalingPolicy) {
			defer wg.Done()

			//InfluxDB 통해 앱GUID 별 컨테이너(복수) 정보(container_interface 등) 획득
			appInfo := a.p.GetAppInfo(policy.AppGuid)

			//리소스 사용량 SET & 오토스케일링 대상 추출
			a.setResourceUsage(&appInfo, policy, &listAutoScaleTarget)
			fmt.Println(">>>>> APP_INFO:", appInfo)

		}(&wg, policy)
	}
	wg.Wait()

	//Request AutoScale API
	fmt.Println(">>>>> LIST_AUTO_SCALING_TARGET:", listAutoScaleTarget)
	for _, target := range listAutoScaleTarget {
		a.requestAutoScale(target)
	}
}

func (a *AutoScalerStruct) requestAutoScale(target model.AutoScaleTarget) {

	/*cfApp, cfErr := a.b.CfClient.GetAppByGuid(target.AppGuid)
	if cfErr != nil {
		fmt.Errorf(">>>>> cf API(GetAppByGuid) error:%v", cfErr)
		return
	}*/

	fmt.Printf(">>>>> Request cf AutoScaling: guid=[%v], instances=[%v]\n", target.AppGuid, target.InstanceCnt)

	var aur cfclient.AppUpdateResource
	aur.Instances, _ = strconv.Atoi(target.InstanceCnt)

	updateResp, updateErr := a.b.CfClient.UpdateApp(target.AppGuid, aur)
	if updateErr != nil {
		fmt.Errorf(">>>>> cf API(UpdateApp) error:%v", updateErr)
		return
	}
	fmt.Println(">>>>> cf API(UpdateApp) resp:", updateResp)

/*
	err := util.PortalExistCHeck()
	if err != nil {
		fmt.Errorf(">>>>> error:%v", err)
		return
	}
	body, _ := json.Marshal(target)
	fmt.Println(">>>>> AUTO_SCALE_REQUEST_BODY:", body)
	resp, status, errMessage := util.HttpRequest(base.SCALE_API_URI,  "POST", nil,  body, *model.PortalClient)
	fmt.Printf(">>>>> RESULT_AUTO_SCALE_API: http.status=[%v], err=[%v], resp=[%v]\n", status, errMessage, resp)
*/
}

//1개의 App에서 n개의 오토스케일 조건이 발생 될 경우 1개의 조건에 대해서만 오토스케일링API를 호출한다.
func (a *AutoScalerStruct) setResourceUsage(appInfo *model.ApplicationInfo, policy model.AppAutoScalingPolicy, listAutoScaleTarget *[]model.AutoScaleTarget) {

	//Scale-Out 조건: 복수개의 인스턴스 중 1개라도 임계치를 초과했을 경우
	//Scale-In  조건: 복수개의 인스턴스 전체가 임계치 미만일 경우

	cfApp, cfErr := a.b.CfClient.GetAppByGuid(appInfo.ApplicationId)
	if cfErr != nil {
		fmt.Errorf(">>>>> cf API(GetAppByGuid) error:%v", cfErr)
		return
	}

	isAppended := false
	belowCpuCnt, belowMemoryCnt := 0, 0

	for _, container := range appInfo.ApplicationContainerInfo {

		if isAppended {
			continue
		}

		container.CpuUsage = a.p.GetContainerCpuUsage(container, policy.MeasureTimeSec)
		fmt.Printf(">>>>> application_id=[%v], container_interface=[%v], cpu_usage=[%v], measure_time=[%v]\n", container.ApplicationId, container.ContainerInterface, container.CpuUsage, policy.MeasureTimeSec)
		container.MemoryUsage = a.p.GetContainerMemoryUsage(container, policy.MeasureTimeSec)
		fmt.Printf(">>>>> application_id=[%v], container_interface=[%v], memory_usage=[%v], measure_time=[%v]\n", container.ApplicationId, container.ContainerInterface, container.MemoryUsage, policy.MeasureTimeSec)

		//오토스케일링 후 목표 인스턴스 개수
		var instanceCntAfterAutoScale uint

		//Append to Scale-Out List
		if policy.AutoScalingOutYn == "Y" {

			if cfApp.Instances + int(policy.InstanceScalingUnit) > int(policy.InstanceMaxCnt) {
				instanceCntAfterAutoScale = policy.InstanceMaxCnt
			} else {
				instanceCntAfterAutoScale = uint(cfApp.Instances) + policy.InstanceScalingUnit
			}

			//임계치 비교하여 Scale-Out 대상에 추가
			if container.CpuUsage > float64(policy.CpuMaxThreshold) && cfApp.Instances < int(policy.InstanceMaxCnt) && policy.AutoScalingCpuYn == "Y" && !isAppended {
				*listAutoScaleTarget = append(*listAutoScaleTarget, generateAutoScaleTarget(container, instanceCntAfterAutoScale, base.SCALE_OUT, base.SCALE_RESOURCE_CPU))
				isAppended = true
			}
			if container.MemoryUsage > float64(policy.MemoryMaxThreshold) && cfApp.Instances < int(policy.InstanceMaxCnt) &&  policy.AutoScalingMemoryYn == "Y" && !isAppended {
				*listAutoScaleTarget = append(*listAutoScaleTarget, generateAutoScaleTarget(container, instanceCntAfterAutoScale, base.SCALE_OUT, base.SCALE_RESOURCE_MEM))
				isAppended = true
			}
		}

		//Append to Scale-In List
		if policy.AutoScalingInYn == "Y" {

			if cfApp.Instances - int(policy.InstanceScalingUnit) < int(policy.InstanceMinCnt) {
				instanceCntAfterAutoScale = policy.InstanceMinCnt
			} else {
				instanceCntAfterAutoScale = uint(cfApp.Instances) - policy.InstanceScalingUnit
			}

			//임계치 비교하여 Scale-In 대상에 추가
			if container.CpuUsage < float64(policy.CpuMinThreshold) && cfApp.Instances > int(policy.InstanceMinCnt) && policy.AutoScalingCpuYn == "Y" && !isAppended {
				//*listAutoScaleTarget = append(*listAutoScaleTarget, generateAutoScaleTarget(container, instanceCntAfterAutoScale, base.SCALE_IN, base.SCALE_RESOURCE_CPU))
				//isAppended = true
				belowCpuCnt++
			}
			if container.MemoryUsage < float64(policy.MemoryMinThreshold) && cfApp.Instances > int(policy.InstanceMinCnt) && policy.AutoScalingMemoryYn == "Y" && !isAppended {
				//*listAutoScaleTarget = append(*listAutoScaleTarget, generateAutoScaleTarget(container, instanceCntAfterAutoScale, base.SCALE_IN, base.SCALE_RESOURCE_MEM))
				//isAppended = true
				belowMemoryCnt++
			}
		}

		//Append to Scale-In List (복수개의 인스턴스 전체가 임계치 미만일 경우)
		if belowCpuCnt == cfApp.Instances && !isAppended {
			*listAutoScaleTarget = append(*listAutoScaleTarget, generateAutoScaleTarget(container, instanceCntAfterAutoScale, base.SCALE_IN, base.SCALE_RESOURCE_CPU))
			isAppended = true
		} else if belowMemoryCnt == cfApp.Instances && !isAppended {
			*listAutoScaleTarget = append(*listAutoScaleTarget, generateAutoScaleTarget(container, instanceCntAfterAutoScale, base.SCALE_IN, base.SCALE_RESOURCE_MEM))
			isAppended = true
		}
	}

}

func generateAutoScaleTarget(container model.ApplicationContainerInfo, instanceCntAfterAutoScale uint, action string, cause string) model.AutoScaleTarget{

	var autoScaleTarget model.AutoScaleTarget
	autoScaleTarget.AppName = container.ApplicationName
	autoScaleTarget.AppGuid = container.ApplicationId
	autoScaleTarget.CpuUsage = util.Floattostr(container.CpuUsage)
	autoScaleTarget.MemoryUsage = util.Floattostr(container.MemoryUsage)
	autoScaleTarget.InstanceCnt = strconv.Itoa(int(instanceCntAfterAutoScale))
	autoScaleTarget.Action = action
	autoScaleTarget.Cause = cause

	return autoScaleTarget
}


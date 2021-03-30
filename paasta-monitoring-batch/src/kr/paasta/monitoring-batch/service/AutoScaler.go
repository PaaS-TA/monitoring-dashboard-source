package service

import (
	"fmt"
	"time"
	//"github.com/cloudfoundry-community/go-cfclient"
	"kr/paasta/monitoring-batch/dao"
	"kr/paasta/monitoring-batch/model"
	"kr/paasta/monitoring-batch/model/base"
	"kr/paasta/monitoring-batch/util"
	"strconv"
	"sync"
	//"github.com/cloudfoundry-community/go-cfclient"
	md "kr/paasta/monitoring-batch/model"
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

func (a *AutoScalerStruct) AutoScale() {

	a.p = PortalAppAlarm(a.b)

	//AutoScaling 정책 조회 (연관 테이블 : app_auto_scaling_policies)
	listAutoScalePolicy, err := dao.AutoScalerDao(a.b.MonitoringDbClient, a.b.Influxclient, a.b.InfluxConfig.ContainerDatabase).
		GetAutoScalePolicy()
	if err != nil {
		fmt.Errorf(">>>>> error:%v", err)
		return
	}
	fmt.Println(">>>>> [AutoScaler.go] AUTO_SCALING_POLICY:", listAutoScalePolicy)

	//cfApp, _ := util.GetAppByGuid(a.b.CfConfig,a.b.CfClientToken,"bf60a3b5-c937-4d9f-ae97-3f7a7ef81d24")
	//
	//var aur md.ScaleProcess
	//aur.Instances = 1
	//
	//updateResp, _ := util.UpdateApp(a.b.CfConfig,a.b.CfClientToken,cfApp.Guid, aur)
	//
	//fmt.Println(">>>>>dfsfef:", updateResp)

	// CF API 호출을 위한 Access Token 획득하기
	var listAutoScaleTarget []model.AutoScaleTarget
	if a.b.CfClientToken.Token == "" {
		a.b.CfClientToken = util.GetUaaToken(a.b.CfConfig)
	} else {
		fmt.Println("time:", a.b.CfClientToken.ExpireTime)
		if a.b.CfClientToken.ExpireTime.Before(time.Now()) {
			a.b.CfClientToken = util.GetUaaToken(a.b.CfConfig)
			//fmt.Println(">>>>> cf token:", cfToken.Token)
		}
	}

	//t1, _ := time.Parse(time.RFC3339, a.b.CfClientToken.ExpireTime.String())
	//fmt.Println("time:", a.b.CfClientToken.ExpireTime)
	//fmt.Println("time now:",time.Now() )
	//if a.b.CfClientToken.ExpireTime.Before(time.Now()) {
	//	cfToken := util.GetUaaToken(a.b.CfConfig)
	//fmt.Println(">>>>> cf token:", cfToken.Token)
	//}
	//cfToken := util.GetUaaToken(a.b.CfConfig)
	//fmt.Println(">>>>> cf token:", cfToken.Token)

	//App 알람 정책 별 사용량 조회
	var wg sync.WaitGroup
	wg.Add(len(listAutoScalePolicy))
	for _, policy := range listAutoScalePolicy {
		go func(wg *sync.WaitGroup, policy model.AppAutoScalingPolicy) {
			defer wg.Done()

			//InfluxDB의 container_metrices Measurement에서 앱 정보를 조회함
			appInfo := a.p.GetAppInfo(policy.AppGuid)

			//리소스 사용량 SET & 오토스케일링 대상 추출
			a.setResourceUsage(&appInfo, policy, &listAutoScaleTarget)
			fmt.Println(">>>>> [AutoScaler.go] APP_INFO:", appInfo)

		}(&wg, policy)
	}
	wg.Wait()

	//Request AutoScale API
	fmt.Println(">>>>> [AutoScaler.go] LIST_AUTO_SCALING_TARGET:", listAutoScaleTarget)
	for _, target := range listAutoScaleTarget {
		a.requestAutoScale(target)
	}
}

func (a *AutoScalerStruct) requestAutoScale(target model.AutoScaleTarget) {

	//cfApp, cfErr := a.b.CfClient.GetAppByGuid(target.AppGuid)
	//if cfErr != nil {
	//	fmt.Errorf(">>>>> cf API(GetAppByGuid) error:%v", cfErr)
	//	return
	//}

	fmt.Printf(">>>>> [AutoScaler.go] Request cf AutoScaling: guid=[%v], instances=[%v]\n", target.AppGuid, target.InstanceCnt)

	var aur md.ScaleProcess
	aur.Instances, _ = strconv.Atoi(target.InstanceCnt)

	updateResp, updateErr := util.UpdateApp(a.b.CfConfig, a.b.CfClientToken, target.AppGuid, aur)
	if updateErr != nil {
		fmt.Errorf(">>>>> cf API(UpdateApp) error:%v", updateErr)
		return
	}
	fmt.Println(">>>>> [AutoScaler.go] cf API(UpdateApp) resp:", updateResp)

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

/**
	1개의 App에서 n개의 오토스케일 조건이 발생 될 경우 1개의 조건에 대해서만 오토스케일링 API를 호출한다.
		- Scale-Out 조건: 복수개의 인스턴스 중 1개라도 임계치를 초과했을 경우
		- Scale-In  조건: 복수개의 인스턴스 전체가 임계치 미만일 경우
 */
func (a *AutoScalerStruct) setResourceUsage(appInfo *model.ApplicationInfo, policy model.AppAutoScalingPolicy, listAutoScaleTarget *[]model.AutoScaleTarget) {

	fmt.Println(">>>>> [AutoScaler.go] CfClientToken : ", a.b.CfClientToken)

	cfApp, cfErr := util.GetAppByGuid(a.b.CfConfig, a.b.CfClientToken, appInfo.ApplicationId)
	if cfErr != nil {
		fmt.Errorf(">>>>> [AutoScaler.go]  cf API(GetAppByGuid) error:%v", cfErr)
		return
	}

	isAppended := false
	belowCpuCnt, belowMemoryCnt := 0, 0

	for _, container := range appInfo.ApplicationContainerInfo {

		if isAppended {
			continue
		}

		container.CpuUsage = a.p.GetContainerCpuUsage(container, policy.MeasureTimeSec)
		container.MemoryUsage = a.p.GetContainerMemoryUsage(container, policy.MeasureTimeSec)

		fmt.Printf(">>>>> [AutoScaler.go] application_id=[%v], container_interface=[%v], measure_time=[%v]\n",container.ApplicationId, container.ContainerInterface, policy.MeasureTimeSec)
		fmt.Printf(">>>>> [AutoScaler.go] cpu_usage=[%v]\n", container.CpuUsage)
		fmt.Printf(">>>>> [AutoScaler.go] memory_usage=[%v]\n", container.MemoryUsage)

		// 오토스케일링 후 목표 인스턴스 개수
		var instanceCntAfterAutoScale uint

		//Append to Scale-Out List
		if policy.AutoScalingOutYn == "Y" {

			if cfApp.Instances+int(policy.InstanceScalingUnit) > int(policy.InstanceMaxCnt) {
				instanceCntAfterAutoScale = policy.InstanceMaxCnt
			} else {
				instanceCntAfterAutoScale = uint(cfApp.Instances) + policy.InstanceScalingUnit
			}

			fmt.Println(">>>>> [AutoScaler.go] cfApp : ", cfApp)
			fmt.Println(">>>>> [AutoScaler.go] cfApp.Instances : ", cfApp.Instances)
			fmt.Println(">>>>> [AutoScaler.go] policy.InstanceScalingUnit : ", policy.InstanceScalingUnit)
			fmt.Println(">>>>> [AutoScaler.go] policy.InstanceMaxCnt : ", policy.InstanceMaxCnt)
			fmt.Println(">>>>> [AutoScaler.go] instanceCntAfterAutoScale : ", instanceCntAfterAutoScale)

			//임계치 비교하여 Scale-Out 대상에 추가
			if container.CpuUsage > float64(policy.CpuMaxThreshold) && cfApp.Instances < int(policy.InstanceMaxCnt) && policy.AutoScalingCpuYn == "Y" && !isAppended {
				*listAutoScaleTarget = append(*listAutoScaleTarget, generateAutoScaleTarget(container, instanceCntAfterAutoScale, base.SCALE_OUT, base.SCALE_RESOURCE_CPU))
				isAppended = true
			}
			if container.MemoryUsage > float64(policy.MemoryMaxThreshold) && cfApp.Instances < int(policy.InstanceMaxCnt) && policy.AutoScalingMemoryYn == "Y" && !isAppended {
				*listAutoScaleTarget = append(*listAutoScaleTarget, generateAutoScaleTarget(container, instanceCntAfterAutoScale, base.SCALE_OUT, base.SCALE_RESOURCE_MEM))
				isAppended = true
			}
		}

		//Append to Scale-In List
		if policy.AutoScalingInYn == "Y" {

			if cfApp.Instances-int(policy.InstanceScalingUnit) < int(policy.InstanceMinCnt) {
				instanceCntAfterAutoScale = policy.InstanceMinCnt
			} else {
				instanceCntAfterAutoScale = uint(cfApp.Instances) - policy.InstanceScalingUnit
			}

			// 임계치 비교하여 Scale-In 대상에 추가
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

func generateAutoScaleTarget(container model.ApplicationContainerInfo, instanceCntAfterAutoScale uint, action string, cause string) model.AutoScaleTarget {

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

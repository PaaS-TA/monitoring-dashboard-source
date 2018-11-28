package service

import (
	"kr/paasta/monitoring-batch/dao"
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
	"kr/paasta/monitoring-batch/util"
	"kr/paasta/monitoring-batch/model"
	cb "kr/paasta/monitoring-batch/model/base"
	"strconv"
	"time"
	"net/http"
	"github.com/cloudfoundry-community/gogobosh"
	"net"
)

func CreteUpdateBoshVms(f *BackendServices, boshConfig BoshConfig, dbClient *gorm.DB){

	//Bosh 접속가능 여부 체크
	conn, err := net.DialTimeout("tcp", boshConfig.BoshUrl ,  3*time.Second)
	defer conn.Close()
	if err != nil {
		fmt.Errorf("#######Bosh Connect Error: ", boshConfig.BoshUrl, err)
		return
	}

	vmsInfos, boshErrs  := dao.GetBoshVmsDao(f.BoshClient).GetDeploymets()
	// Error가 여러건일 경우 대해 고려해야함.
	if len(boshErrs) > 0 {
		var returnErrMessage string
		for _, err := range boshErrs{

			returnErrMessage = returnErrMessage + " " + err.Error()
			fmt.Errorf("GetDeploymets_Error Message:::",returnErrMessage)
			//Bosh 접속시 간헐적으로 The token has been revoked(401) 에러 발생한다.
			//에러 발생시 토큰 재발급 받는다.
			//그리고 다음 Batch 실행시 정상작동하도록 변경
			boshClientConfig := &gogobosh.Config{
				BOSHAddress: 		fmt.Sprintf("https://%s", boshConfig.BoshUrl),
				Username:    		boshConfig.BoshId,
				Password:    		boshConfig.BoshPass,
				HttpClient:        	http.DefaultClient,
				SkipSslValidation: 	true,
			}
			boshClient, _ := gogobosh.NewClient(boshClientConfig)
			f.BoshClient = boshClient
			return
		}

	}


	var zoneNames []string
	//Bosh에서 Zone정보를 추출한다.
	for _, value := range vmsInfos{

		//Deployment Name이 config.ini에 정의된
		// bosh.cf.deployment.name만 등록
		if value.Name == boshConfig.CfDeploymentName {

			for _, data := range value.VMS{

				fmt.Println( "ZonNme::", data.AZ)
				fmt.Println( "ID::", data.ID)
				//fmt.Println( "vmCid::", data.VMCID)
				fmt.Println( "jobName::", data.JobName)

				if data.AZ != "" {
					//fmt.Println("AWS CF deployments")
					if util.StringInSlice(data.AZ, zoneNames) == false{
						zoneNames = append( zoneNames, data.AZ)
					}

				}else{
					//fmt.Println("CF deployments")
					vmNames := strings.Split(data.JobName, "_")
					zoneName := vmNames[len(vmNames)-1]

					if util.StringInSlice(zoneName, zoneNames) == false{
						zoneNames = append( zoneNames, zoneName)
					}
				}
			}
		}
	}

	dbZoneInfos, dbErr := dao.GetBoshVmsDao(f.BoshClient).GetZoneInfosByZoneNames(zoneNames, dbClient)

	if dbErr != nil{
		fmt.Errorf("BoshVMS Collect Error Occur : ", dbErr)
	}

	//Zone정보가 없으면 새로 insert한다.
	for _, zoneName := range zoneNames{
		isExist := false
		for _, dbZoneInfo := range dbZoneInfos{

			if dbZoneInfo.Name == zoneName{
				isExist = true
			}
		}
		if isExist == false{
			fmt.Println("CreateZoneData zoneName =>", zoneName)

			createErr := dao.GetBoshVmsDao(f.BoshClient).CreateZoneData(zoneName, dbClient)
			if createErr != nil{
				fmt.Errorf("CreateZone Error Occured : ", createErr)
			}
		}
	}

	for _, value := range vmsInfos{

		//Deployment 명이 CF/DIego만 받아 VMS 에 저저장된다.
		if value.Name == boshConfig.CfDeploymentName || value.Name == boshConfig.DiegoDeploymentName{

			for _, data := range value.VMS{
				var vmsInfo model.Vm

				var zoneName string
				var vmName string
				var vmNames []string

				if data.AZ != "" {
					vmName = data.JobName
					zoneName = data.AZ
				}else{
					//Zone명을 추출한다. jobName의 끝자리에 zone명을 사용
					// ex) api_z1
					vmNames = strings.Split(data.JobName, "_")
					zoneName = vmNames[len(vmNames)-1]
				}

				jobExist, vmErr := dao.GetBoshVmsDao(f.BoshClient).IsExistJobName(data.JobName+"/"+strconv.Itoa(data.Index), dbClient)

				if vmErr != nil{
					fmt.Errorf("CreateVms Error Occured : ", vmErr)
				}

				//Job이 존재 하지 않으면 CF/Diego정보 INsert
				if jobExist == false{

					zoneInfo, _ := dao.GetBoshVmsDao(f.BoshClient).GetZoneInfosByZoneName(zoneName, dbClient)

					if len(vmNames) > 0 {
						fmt.Println("CF deployments get cell type")
						//VM명이 CELL 이면 VM_TYPE 은 CEL 이 된다.
						if vmNames[0] == boshConfig.CellNamePrefix{
							vmsInfo.VmType = cb.VM_TYPE_CEL
						}else{
							vmsInfo.VmType = cb.VM_TYPE
						}
					}else{
						fmt.Println("AWS CF deployments get cell type")
						//VM명이 CELL 이면 VM_TYPE 은 CEL 이 된다.
						if strings.Contains(vmName, boshConfig.CellNamePrefix){
							vmsInfo.VmType = cb.VM_TYPE_CEL
						}else{
							vmsInfo.VmType = cb.VM_TYPE
						}
					}

					if len(data.IPs) > 0{
						vmsInfo.ZoneId = zoneInfo.Id
						vmsInfo.Name = data.JobName +"/"+strconv.Itoa(data.Index)
						vmsInfo.Ip   = data.IPs[0]

						fmt.Println("CreateVmData vmsInfo =>", vmsInfo)
						dao.GetBoshVmsDao(f.BoshClient).CreateVmData(vmsInfo, dbClient)
					}



				}else{
					// 이미 존재하면 IP 가 변경 되었는지 체크하여 IP가 변경 되어 있으면 Update한다.
					jobInfo, _ := dao.GetBoshVmsDao(f.BoshClient).GetJobInfo(data.JobName+"/"+strconv.Itoa(data.Index), dbClient )

					if len(data.IPs) > 0{
						if jobInfo.Ip != data.IPs[0] {
							vmsInfo.Id   = jobInfo.Id
							vmsInfo.ZoneId = jobInfo.ZoneId
							vmsInfo.Ip   = data.IPs[0]
							vmsInfo.RegDate = jobInfo.RegDate
							vmsInfo.RegUser = jobInfo.RegUser
							vmsInfo.ModiDate = time.Now().Local()
							vmsInfo.ModiUser = "batch"

							fmt.Println("UpdateVmData vmsInfo =>", vmsInfo)
							dao.GetBoshVmsDao(f.BoshClient).UpdateVmData( vmsInfo , dbClient)
						}
					}

					// zone 정보 변경 시, Update
					zoneInfo, _ := dao.GetBoshVmsDao(f.BoshClient).GetZoneInfosByZoneId(int(jobInfo.ZoneId), dbClient)

					if data.AZ != zoneInfo.Name {
						//fmt.Println(">>>>>>>>>>>>>> Update zone info data =>", data.AZ, zoneInfo.Name)
						tmpzone, _ := dao.GetBoshVmsDao(f.BoshClient).GetZoneInfosByZoneName(data.AZ, dbClient)

						vmsInfo.Id   = jobInfo.Id
						vmsInfo.ZoneId = tmpzone.Id
						vmsInfo.Ip   = jobInfo.Ip
						vmsInfo.RegDate = jobInfo.RegDate
						vmsInfo.RegUser = jobInfo.RegUser
						vmsInfo.ModiDate = time.Now().Local()
						vmsInfo.ModiUser = "batch"

						//fmt.Println(">>>>>>>>>>>>>> Update zone info vms =>", vmsInfo)
						dao.GetBoshVmsDao(f.BoshClient).UpdateVmZoneData( vmsInfo , dbClient)
					}
				}

			}
		}
	}

	dbVmsInfoList, dbErr := dao.GetBoshVmsDao(f.BoshClient).GetJobInfoList(dbClient)
	// Error가 여러건일 경우 대해 고려해야함.
	if dbErr != nil {
		fmt.Errorf(" GetJobInfoList Error:::", dbErr)

	}

	var deleteVmList []model.Vm
	for _, dbVmData := range dbVmsInfoList {
		isExist := false
		for _, vmData := range vmsInfos {
			//Deployment 명이 CF/DIego만
			if vmData.Name == boshConfig.CfDeploymentName || vmData.Name == boshConfig.DiegoDeploymentName {

				for _, data := range vmData.VMS{
					if dbVmData.Name == data.JobName+"/"+strconv.Itoa(data.Index){
						isExist = true
						break
					}
				}

			}
		}
		if isExist == false{
			deleteVmList = append(deleteVmList, dbVmData)
		}
	}

	//Bosh VM 정보 조회시 오류 발생하면
	//VM 정보가 모두 삭제되는 경우 발생한다.
	//이를 방지하게 위해 아래 조건 추가( boshErrs == 0  )
	if len(boshErrs) == 0 {
		for _, data := range deleteVmList{
			//fmt.Println("DeleteVmInfo data =>", data)
			dao.GetBoshVmsDao(f.BoshClient).DeleteVmInfo(data, dbClient)
		}

		dbZoneInfoList, dbErr := dao.GetBoshVmsDao(f.BoshClient).GetZoneInfosList(dbClient)

		if dbErr != nil{
			fmt.Errorf(" GetZoneInfosList Error:::", dbErr)
		}

		//CF에서 삭제된 Zone Data DB에서 삭제
		var deleteZoneList []model.Zone
		for _, data := range dbZoneInfoList{
			isExist := false
			for _, zoneData := range zoneNames{
				if data.Name == zoneData{
					isExist = true
				}
			}
			if isExist == false {
				deleteZoneList = append(deleteZoneList, data)
			}
		}

		for _, data := range deleteZoneList{
			fmt.Println("DeleteZoneInfo data =>", data)
			dao.GetBoshVmsDao(f.BoshClient).DeleteZoneInfo(data, dbClient)
		}
		fmt.Println("end =========================================================================")
	}

}



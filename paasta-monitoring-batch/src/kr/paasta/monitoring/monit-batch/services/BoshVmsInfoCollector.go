package services

import (
	"kr/paasta/monitoring/monit-batch/dao"
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
	"kr/paasta/monitoring/monit-batch/util"
	"kr/paasta/monitoring/monit-batch/models"
	cb "kr/paasta/monitoring/monit-batch/models/base"
	"strconv"
	"time"
)

func CreteUpdateBoshVms(f *BackendServices, boshConfig BoshConfig, dbClient *gorm.DB){

	vmsInfos, errs  := dao.GetBoshVmsDao(f.BoshClient).GetDeploymets()
	// Error가 여러건일 경우 대해 고려해야함.
	if len(errs) > 0 {
		var returnErrMessage string
		for _, err := range errs{
			returnErrMessage = returnErrMessage + " " + err.Error()
			fmt.Errorf("Error Message:::",returnErrMessage)
		}
	}

	var zoneNames []string
	//Bosh에서 Zone정보를 추출한다.
	for _, value := range vmsInfos{
		//Deployment Name이 config.ini에 정의된
		// bosh.cf.deployment.name, bosh.diego.deployment.name 만 등록
		if value.Name == boshConfig.CfDeploymentName || value.Name == boshConfig.DiegoDeploymentName{

			for _, data := range value.VMS{

				vmNames := strings.Split(data.JobName, "_")
				zoneName := vmNames[len(vmNames)-1]

				if util.StringInSlice(zoneName, zoneNames) == false{
					zoneNames = append( zoneNames, zoneName)
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
				var vmsInfo models.Vm
				//Zone명을 추출한다. jobName의 끝자리에 zone명을 사용
				// ex) api_z1
				vmNames := strings.Split(data.JobName, "_")
				zoneName := vmNames[len(vmNames)-1]

				jobExist, vmErr := dao.GetBoshVmsDao(f.BoshClient).IsExistJobName(data.JobName+"/"+strconv.Itoa(data.Index), dbClient)

				if vmErr != nil{
					fmt.Errorf("CreateVms Error Occured : ", vmErr)
				}

				//Job이 존재 하지 않으면 CF/Diego정보 INsert
				if jobExist == false{

					zoneInfo, _ := dao.GetBoshVmsDao(f.BoshClient).GetZoneInfosByZoneName(zoneName, dbClient)

					if len(vmNames) > 0 {
						//VM명이 CELL 이면 VM_TYPE 은 CEL 이 된다.
						if vmNames[0] == boshConfig.CellNamePrefix{
							vmsInfo.VmType = cb.VM_TYPE_CEL
						}else{
							vmsInfo.VmType = cb.VM_TYPE
						}

						if len(data.IPs) > 0{
							vmsInfo.ZoneId = zoneInfo.Id
							vmsInfo.Name = data.JobName +"/"+strconv.Itoa(data.Index)
							vmsInfo.Ip   = data.IPs[0]
							dao.GetBoshVmsDao(f.BoshClient).CreateVmData(vmsInfo, dbClient)
						}

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

							dao.GetBoshVmsDao(f.BoshClient).UpdateVmData( vmsInfo , dbClient)
						}
					}
				}

			}
		}
	}

	dbVmsInfoList, dbErr := dao.GetBoshVmsDao(f.BoshClient).GetJobInfoList(dbClient)
	// Error가 여러건일 경우 대해 고려해야함.
	if len(dbErr) > 0 {
		var returnErrMessage string
		for _, err := range errs{
			returnErrMessage = returnErrMessage + " " + err.Error()
			fmt.Errorf("Error Message:::",returnErrMessage)
		}
	}

	var deleteVmList []models.Vm
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

	for _, data := range deleteVmList{
		dao.GetBoshVmsDao(f.BoshClient).DeleteVmInfo(data, dbClient)
	}

	dbZoneInfoList, dbErr := dao.GetBoshVmsDao(f.BoshClient).GetZoneInfosList(dbClient)


	// Error가 여러건일 경우 대해 고려해야함.
	if len(dbErr) > 0 {
		var returnErrMessage string
		for _, err := range errs{
			returnErrMessage = returnErrMessage + " " + err.Error()
			fmt.Errorf("Error Message:::",returnErrMessage)
		}
	}

	//CF에서 삭제된 Zone Data DB에서 삭제
	var deleteZoneList []models.Zone
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
		dao.GetBoshVmsDao(f.BoshClient).DeleteZoneInfo(data, dbClient)
	}


}
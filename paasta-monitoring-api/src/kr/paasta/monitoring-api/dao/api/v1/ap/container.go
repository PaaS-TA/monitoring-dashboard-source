package ap

import (
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient"
	"gorm.io/gorm"
	"paasta-monitoring-api/dao/api/v1/common"
	models "paasta-monitoring-api/models/api/v1"
)

type ApContainerDao struct {
	DbInfo         *gorm.DB
	InfluxDbClient models.InfluxDbClient
	CfClient       *cfclient.Client
}

func GetApContainerDao(DbInfo *gorm.DB, InfluxDbClient models.InfluxDbClient, CfClient *cfclient.Client) *ApContainerDao {
	return &ApContainerDao{
		DbInfo:         DbInfo,
		InfluxDbClient: InfluxDbClient,
		CfClient:       CfClient,
	}
}

func (ap *ApContainerDao) GetZoneInfo() ([]models.ZoneInfo, error) {
	var response []models.ZoneInfo
	results := ap.DbInfo.Debug().Table("zones").
		Select("zones.name AS zone_name, COUNT(*) AS cell_cnt").
		Joins("INNER JOIN vms ON zones.id = vms.zone_id").
		Where("vm_type = ?", "cel").
		Group("zones.name").
		Find(&response)

	if results.Error != nil {
		fmt.Println(results.Error)
		return response, results.Error
	}

	return response, nil
}

func (ap *ApContainerDao) GetCellInfo() ([]models.CellInfo, error) {
	var response []models.CellInfo
	results := ap.DbInfo.Debug().Table("zones").
		Select("zones.name AS zone_name, vms.name AS cell_name, vms.ip, vms.id").
		Joins("INNER JOIN vms ON zones.id = vms.zone_id").
		Where("vm_type = ?", "cel").
		Order("cell_name ASC").
		Find(&response)

	if results.Error != nil {
		fmt.Println(results.Error)
		return response, results.Error
	}

	// For counting AppCnt & ContainerCnt
	apps, _ := ap.CfClient.ListApps()
	for _, app := range apps {
		appStats, _ := ap.CfClient.GetAppStats(app.Guid)
		for i, cell := range response {
			if appStats["0"].Stats.Host == cell.Ip {
				response[i].AppCnt += 1
				response[i].ContainerCnt += uint(app.Instances)
			}
		}
	}

	return response, nil
}

func (ap *ApContainerDao) GetAppInfo() ([]models.AppInfo, error) {
	var response []models.AppInfo
	apps, _ := ap.CfClient.ListApps()
	cells, _ := ap.GetCellInfo()

	// Range App List
	for _, app := range apps {
		appStats, _ := ap.CfClient.GetAppStats(app.Guid)
		appEnvs, _ := ap.CfClient.GetAppEnv(app.Guid)
		appEnv := appEnvs.ApplicationEnv["VCAP_APPLICATION"].(map[string]interface{})

		// Range Cell List (For mapping CellName into AppInfo struct)
		for _, cell := range cells {
			if appStats["0"].Stats.Host == cell.Ip {
				tmp := models.AppInfo{
					CellName:  cell.CellName,
					AppName:   appEnv["application_name"].(string),
					Uri:       appEnv["application_uris"].([]interface{})[0].(string),
					Buildpack: app.Buildpack,
					Instances: app.Instances,
					Status:    app.State,
					Memory:    app.Memory,
					DiskQuota: app.DiskQuota,
					CfApi:     appEnv["cf_api"].(string),
					CreatedAt: app.CreatedAt,
					UpdatedAt: app.UpdatedAt,
				}
				response = append(response, tmp)
			}
		}
	}

	return response, nil
}

func (ap *ApContainerDao) GetContainerPageOverview() (models.Overview, error) {
	var response models.Overview
	zones, _ := ap.GetZoneInfo()
	cells, _ := ap.GetCellInfo()
	apps, _ := ap.GetAppInfo()

	for i, zone := range zones {
		for j, cell := range cells {
			for _, app := range apps {
				if zone.ZoneName == cell.ZoneName {
					if cell.CellName == app.CellName {
						zones[i].CellInfo = cells
						cells[j].AppInfo = apps
					}
				}
			}
		}
	}
	response.ZoneInfo = zones

	return response, nil
}

func (ap *ApContainerDao) GetContainerStatus() (models.Status, error) {
	var response models.Status
	var status models.StatusByResource
	var statuses []models.StatusByResource

	params := models.AlarmPolicies{
		OriginType: "con",
	}
	policies, _ := common.GetAlarmPolicyDao(ap.DbInfo).GetAlarmPolicy(params)
	apps, _ := ap.CfClient.ListApps()

	// For inserting status data by container resources
	for _, app := range apps {
		appStats, _ := ap.CfClient.GetAppStats(app.Guid)
		for _, appStat := range appStats {
			for _, policy := range policies {
				switch policy.AlarmType {
				case "cpu":
					if appStat.Stats.Usage.CPU >= float64(policy.CriticalThreshold) {
						status.CpuStatus = "Critical"
					} else if appStat.Stats.Usage.CPU >= float64(policy.WarningThreshold) {
						status.CpuStatus = "Warning"
					} else {
						status.CpuStatus = "Running"
					}
				case "memory":
					if float64(appStat.Stats.Usage.Mem/appStat.Stats.MemQuota) >= float64(policy.CriticalThreshold) {
						status.MemoryStatus = "Critical"
					} else if float64(appStat.Stats.Usage.Mem/appStat.Stats.MemQuota) >= float64(policy.WarningThreshold) {
						status.MemoryStatus = "Warning"
					} else {
						status.MemoryStatus = "Running"
					}
				case "disk":
					if float64(appStat.Stats.Usage.Disk/appStat.Stats.DiskQuota) >= float64(policy.CriticalThreshold) {
						status.DiskStatus = "Critical"
					} else if float64(appStat.Stats.Usage.Disk/appStat.Stats.DiskQuota) >= float64(policy.WarningThreshold) {
						status.DiskStatus = "Warning"
					} else {
						status.DiskStatus = "Running"
					}
				}
			}
			statuses = append(statuses, status)
		}
	}

	// For defining container status
	for i, status := range statuses {
		if status.CpuStatus == "Critical" || status.MemoryStatus == "Critical" || status.DiskStatus == "Critical" {
			statuses[i].TotalStatus = "Critical"
		} else if status.CpuStatus == "Warning" || status.MemoryStatus == "Warning" || status.DiskStatus == "Warning" {
			statuses[i].TotalStatus = "Warning"
		} else {
			statuses[i].TotalStatus = "Running"
		}
	}

	// For counting container status
	for _, status := range statuses {
		switch status.TotalStatus {
		case "Critical":
			response.Critical++
		case "Warning":
			response.Warning++
		case "Running":
			response.Running++
		}
	}

	return response, nil
}

func (ap *ApContainerDao) GetCellStatus() (models.Status, error) {
	var response models.Status

	return response, nil
}

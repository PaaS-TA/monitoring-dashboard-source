package ap

import (
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient"
	client "github.com/influxdata/influxdb1-client/v2"
	"gorm.io/gorm"
	"paasta-monitoring-api/dao/api/v1/common"
	"paasta-monitoring-api/helpers"
	models "paasta-monitoring-api/models/api/v1"
	"reflect"
	"sync"
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

func (ap *ApContainerDao) GetContainerInfo() ([]models.ContainerInfo, error) {
	var request models.InfluxDbQueryRequest
	var response []models.ContainerInfo
	cells, _ := ap.GetCellInfo()

	// For making appMap that contains Container infos
	appMap := make(map[string]map[string]string)
	for _, cell := range cells {
		request.CellIp = cell.Ip
		request.Sql = "SELECT application_name, application_index, container_interface, value FROM container_metrics " +
			"WHERE cell_ip = '%s' AND \"name\" = 'load_average' AND container_id <> '/' AND time > NOW() - 2m"
		results, _ := ap.GetQueryResultsFromInfluxDb(request)
		values, _ := helpers.InfluxConverterToMap(&results)

		for _, value := range values {
			containerMap := make(map[string]string)

			appName := reflect.ValueOf(value["application_name"]).String()
			containerId := reflect.ValueOf(value["container_interface"]).String()
			applicationIndex := reflect.ValueOf(value["application_index"]).String()
			containerMap[applicationIndex] = containerId

			if exists, ok := appMap[appName]; ok {
				for k, v := range containerMap {
					exists[k] = v
					appMap[appName] = exists
				}
			} else {
				appMap[appName] = containerMap
			}
		}
	}

	// For containing containerInfo into ContainerInfo struct by appName
	for appName, containerMap := range appMap {
		var containers []models.Container

		for AppIndex, containerId := range containerMap {
			tmp := models.Container{
				AppIndex: AppIndex, ContainerId: containerId,
			}
			containers = append(containers, tmp)
		}

		containerInfo := models.ContainerInfo{
			AppName: appName, Container: containers,
		}
		response = append(response, containerInfo)
	}

	return response, nil
}

func (ap *ApContainerDao) GetContainerPageOverview() (models.Overview, error) {
	var response models.Overview
	zones, _ := ap.GetZoneInfo()
	cells, _ := ap.GetCellInfo()
	apps, _ := ap.GetAppInfo()
	containers, _ := ap.GetContainerInfo()

	for i, zone := range zones {
		for j, cell := range cells {
			for k, app := range apps {
				for l, container := range containers {
					if zone.ZoneName == cell.ZoneName {
						if cell.CellName == app.CellName {
							if app.AppName == container.AppName {
								zones[i].CellInfo = cells
								cells[j].AppInfo = apps
								apps[k].ContainerInfo = &containers[l]
							}
						}
					}
				}
			}
		}
	}
	response.ZoneInfo = zones

	return response, nil
}

func (ap *ApContainerDao) GetContainerStatus() (models.Status, error) {
	var status models.StatusByResource
	var statuses []models.StatusByResource

	params := models.AlarmPolicies{
		OriginType: "con",
	}
	policies, _ := common.GetAlarmPolicyDao(ap.DbInfo).GetAlarmPolicy(params)
	apps, _ := ap.CfClient.ListApps()

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

	response := helpers.SetStatus(statuses)
	return response, nil
}

func (ap *ApContainerDao) GetCellStatus() (models.Status, error) {
	var request models.InfluxDbQueryRequest
	var cellsMetricData []models.CellMetricData
	var status models.StatusByResource
	var statuses []models.StatusByResource

	alarmPolicyParam := models.AlarmPolicies{
		OriginType: "pas",
	}
	policies, _ := common.GetAlarmPolicyDao(ap.DbInfo).GetAlarmPolicy(alarmPolicyParam)
	cells, _ := ap.GetCellInfo()

	for _, cell := range cells {
		request.CellIp = cell.Ip
		tmp := ap.GetCellMetricData(request)
		cellsMetricData = append(cellsMetricData, tmp)
	}

	convertedMetricData := helpers.ConvertDataFormatForCellMetricData(cellsMetricData)

	for _, cellMetricData := range convertedMetricData {
		for _, policy := range policies {
			switch policy.AlarmType {
			case "cpu":
				if cellMetricData.CpuUsage >= float64(policy.CriticalThreshold) {
					status.CpuStatus = "Critical"
				} else if cellMetricData.CpuUsage >= float64(policy.WarningThreshold) {
					status.CpuStatus = "Warning"
				} else if cellMetricData.CpuUsage == 0 {
					status.CpuStatus = "Failed"
				} else {
					status.CpuStatus = "Running"
				}
			case "memory":
				if cellMetricData.MemUsage >= float64(policy.CriticalThreshold) {
					status.MemoryStatus = "Critical"
				} else if cellMetricData.MemUsage >= float64(policy.WarningThreshold) {
					status.MemoryStatus = "Warning"
				} else if cellMetricData.MemUsage == 0 {
					status.MemoryStatus = "Failed"
				} else {
					status.MemoryStatus = "Running"
				}
			case "disk":
				if cellMetricData.DiskUsage >= float64(policy.CriticalThreshold) {
					status.DiskStatus = "Critical"
				} else if cellMetricData.DiskUsage >= float64(policy.WarningThreshold) {
					status.DiskStatus = "Warning"
				} else if cellMetricData.DiskUsage == 0 {
					status.DiskStatus = "Failed"
				} else {
					status.DiskStatus = "Running"
				}
			}
		}
		statuses = append(statuses, status)
	}

	response := helpers.SetStatus(statuses)
	return response, nil
}

func (ap *ApContainerDao) GetCellMetricData(request models.InfluxDbQueryRequest) models.CellMetricData {
	var response models.CellMetricData
	var cpuCore, cpuUsage, memTotal, memFree, diskTotal, diskUsage client.Response
	var wg sync.WaitGroup

	wg.Add(6)
	for i := 0; i < 6; i++ {
		go func(index int) {
			switch index {
			case 0:
				request.MetricName = "cpuStats.core"
				request.Sql = "SELECT value FROM cf_metrics WHERE ip = '%s' AND time > NOW() - 1m AND metricname =~ /%s/ GROUP BY metricname ORDER BY time DESC LIMIT 1"
				cpuCore, _ = ap.GetQueryResultsFromInfluxDb(request)
			case 1:
				request.MetricName = "cpuStats.core"
				request.Sql = "SELECT MEAN(value) AS value FROM cf_metrics WHERE ip = '%s' AND time > NOW() - 1m AND metricname =~ /%s/"
				cpuUsage, _ = ap.GetQueryResultsFromInfluxDb(request)
			case 2:
				request.MetricName = "memoryStats.TotalMemory"
				request.Sql = "SELECT MEAN(value) AS value FROM cf_metrics WHERE ip = '%s' AND time > NOW() - 1m AND metricname = '%s'"
				memTotal, _ = ap.GetQueryResultsFromInfluxDb(request)
			case 3:
				request.MetricName = "memoryStats.FreeMemory"
				request.Sql = "SELECT MEAN(value) AS value FROM cf_metrics WHERE ip = '%s' AND time > NOW() - 1m AND metricname = '%s'"
				memFree, _ = ap.GetQueryResultsFromInfluxDb(request)
			case 4:
				request.MetricName = "diskStats./.Total"
				request.Sql = "SELECT MEAN(value) AS value FROM cf_metrics WHERE ip = '%s' AND time > NOW() - 1m AND metricname = '%s'"
				diskTotal, _ = ap.GetQueryResultsFromInfluxDb(request)
			case 5:
				request.MetricName = "diskStats./.Usage"
				request.Sql = "SELECT MEAN(value) AS value FROM cf_metrics WHERE ip = '%s' AND time > NOW() - 1m AND metricname = '%s'"
				diskUsage, _ = ap.GetQueryResultsFromInfluxDb(request)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	response.CpuCore, _ = helpers.InfluxConverterToMap(&cpuCore)
	response.CpuUsage, _ = helpers.InfluxConverter(&cpuUsage)
	response.MemTotal, _ = helpers.InfluxConverter(&memTotal)
	response.MemFree, _ = helpers.InfluxConverter(&memFree)
	response.DiskTotal, _ = helpers.InfluxConverter(&diskTotal)
	response.DiskUsage, _ = helpers.InfluxConverter(&diskUsage)

	return response
}

func (ap *ApContainerDao) GetQueryResultsFromInfluxDb(request models.InfluxDbQueryRequest) (_ client.Response, errMsg models.ErrMessage) {
	var errLogMsg string
	var query client.Query

	defer func() {
		if r := recover(); r != nil {
			errMsg = models.ErrMessage{
				"Message": errLogMsg,
			}
		}
	}()

	if request.MetricName != "" {
		query = client.Query{
			Command:  fmt.Sprintf(request.Sql, request.CellIp, request.MetricName),
			Database: ap.InfluxDbClient.DbName.PaastaDatabase,
		}
	} else {
		query = client.Query{
			Command:  fmt.Sprintf(request.Sql, request.CellIp),
			Database: ap.InfluxDbClient.DbName.ContainerDatabase,
		}
	}

	response, err := ap.InfluxDbClient.HttpClient.Query(query)
	if err != nil {
		errLogMsg = err.Error()
	}

	return *response, errMsg
}

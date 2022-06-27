package ap

import (
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient"
	"gorm.io/gorm"
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

func (ap *ApContainerDao) GetCellInfo() ([]models.CellInfo, error) {
	var response []models.CellInfo
	results := ap.DbInfo.Debug().Table("vms").
		Select("*").
		Where("vm_type = ?", "cel").
		Find(&response)

	if results.Error != nil {
		fmt.Println(results.Error)
		return response, results.Error
	}

	return response, nil
}

func (ap *ApContainerDao) GetZoneInfo() ([]models.ZoneInfo, error) {
	var response []models.ZoneInfo
	results := ap.DbInfo.Debug().Table("zones").
		Select("*").
		Find(&response)

	if results.Error != nil {
		fmt.Println(results.Error)
		return response, results.Error
	}

	return response, nil
}

func (ap *ApContainerDao) GetAppInfo() ([]models.AppInfo, error) {
	var response []models.AppInfo
	apps, _ := ap.CfClient.ListApps()

	for _, app := range apps {
		appEnvs, _ := ap.CfClient.GetAppEnv(app.Guid)
		appEnv := appEnvs.ApplicationEnv["VCAP_APPLICATION"].(map[string]interface{})

		tmp := models.AppInfo{
			Name:      appEnv["application_name"].(string),
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

	return response, nil
}

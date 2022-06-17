package ap

import (
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient"
	influxDb "github.com/influxdata/influxdb1-client/v2"
	"github.com/jinzhu/gorm"
	models "paasta-monitoring-api/models/api/v1"
)

type ApContainerDao struct {
	DbInfo             *gorm.DB
	InfluxDbInfo       influxDb.Client
	Databases          models.Databases
	CloudFoundryClient *cfclient.Client
}

func GetApContainerDao(DbInfo *gorm.DB, InfluxDbInfo influxDb.Client, Databases models.Databases, CloudFoundryClient *cfclient.Client) *ApContainerDao {
	return &ApContainerDao{
		DbInfo:             DbInfo,
		InfluxDbInfo:       InfluxDbInfo,
		Databases:          Databases,
		CloudFoundryClient: CloudFoundryClient,
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

func (ap *ApContainerDao) GetAppInfo() ([]cfclient.App, error) {
	return ap.CloudFoundryClient.ListApps()
	// -----------------------------------------------------------------------------------------------------------
	/*var aaa interface{}
	cellInfos, _ := ap.GetCellInfo()
	query := "select application_name, application_index, container_interface, value from container_metrics " +
		"where cell_ip = '%s' and \"name\" = 'load_average'  and container_id <> '/' and time > now() - %s"
	var q influxDb.Query
	for _, cellInfo := range cellInfos {
		var resp *influxDb.Response
		q = influxDb.Query{
			Command:  fmt.Sprintf(query, cellInfo.Ip, "120s"),
			Database: ap.Databases.ContainerDatabase,
		}

		fmt.Println("GetCellContainerList Sql======>", q)
		resp, err := ap.InfluxDbInfo.Query(q)
		if err != nil {
			return resp, err
		}
		aaa = resp
	}
	temp := aaa.(*influxDb.Response)
	return temp, nil*/
	// -----------------------------------------------------------------------------------------------------------
}

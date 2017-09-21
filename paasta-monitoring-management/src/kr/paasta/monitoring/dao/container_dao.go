package dao

import (
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring/domain"
	"kr/paasta/monitoring/util"
	client "github.com/influxdata/influxdb/client/v2"
	"fmt"
)

var zonDataSource = "container_metric_db"

type ContainerDao struct {
	txn   *gorm.DB
	influxClient 	client.Client
}

func GetContainerDao(txn *gorm.DB, influxClient client.Client) *ContainerDao {
	return &ContainerDao{
		txn:   txn,
		influxClient: 	influxClient,
	}
}

//Cell 목록 조회
func (f ContainerDao) GetCellList() ([]domain.ZoneCellInfo, domain.ErrMessage){

	cells := []domain.ZoneCellInfo{}

	status := f.txn.Debug().Table("zones").Order("cell_name asc").
		Select("zones.name as zone_name, vms.id,  vms.name as cell_name, vms.ip").
		Joins("inner join vms on zones.id = vms.zone_id and vms.vm_type = 'Cel' ").Find(&cells)

	err := util.GetError().DbCheckError(status.Error)
	if err != nil{
		return nil, err
	}

	return cells, err
}


func (b ContainerDao) GetCellContainersList(cellIp string) (_ client.Response, errMsg domain.ErrMessage)  {

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = domain.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	getContainerListSql := "select application_name, application_index, container_interface, value from container_metrics where cell_ip = '%s' and \"name\" = 'load_average'  and container_id <> '/' and time > now() - %s";

	var q client.Query

	q = client.Query{
		Command:  fmt.Sprintf( getContainerListSql,
			cellIp, "60s"),
		Database: zonDataSource,
	}

	fmt.Println("GetCellContainerList Sql======>", q)
	resp, err := b.influxClient.Query(q)
	if err != nil{
		errLogMsg = err.Error()
	}

	return util.GetError().CheckError(*resp, err)
}
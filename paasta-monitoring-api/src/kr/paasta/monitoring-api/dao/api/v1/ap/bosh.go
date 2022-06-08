package ap

import (
	"fmt"
	"github.com/influxdata/influxdb1-client/v2"
	"github.com/jinzhu/gorm"
	models "paasta-monitoring-api/models/api/v1"
)

type BoshDao struct {
	DbInfo         *gorm.DB
	InfluxDBClient client.Client
	BoshInfoList   []models.Bosh
	Env            map[string]interface{}
}

func GetBoshDao(DbInfo *gorm.DB, InfluxDBClient client.Client, BoshInfoList []models.Bosh, Env map[string]interface{}) *BoshDao {
	return &BoshDao{
		DbInfo:         DbInfo,
		InfluxDBClient: InfluxDBClient,
		BoshInfoList:   BoshInfoList,
		Env:            Env,
	}
}

func (b *BoshDao) GetBoshProcessByMemory(uuid string) (*client.Response, error) {
	// 전달 받은 계정 정보로 데이터베이스에 계정이 존재하는지 확인한다. (test code)
	//var results []models.BoshProcess

	// InfluxDB 프로세스 목록 조회
	getBoshTopprocessListSql := "select proc_name as process_name, time, proc_index, proc_pid, mem_usage from bosh_process_metrics where id = '%s' and time > now() - %s order by time desc"
	q := client.Query{
		Command:  fmt.Sprintf(getBoshTopprocessListSql, uuid, "1m"),
		Database: b.Env["paas_metric_db_name_bosh"].(string),
	}
	fmt.Println("GetBoshProcessByMemory Sql======>", q)

	resp, err := b.InfluxDBClient.Query(q)
	if err != nil {
		return resp, err
	}
	fmt.Println("GetBoshProcessByMemory resp======>", resp)

	return resp, err
}

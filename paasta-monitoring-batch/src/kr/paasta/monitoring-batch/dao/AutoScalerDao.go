package dao

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/jinzhu/gorm"
	"kr/paasta/monitoring-batch/model"
	"kr/paasta/monitoring-batch/model/base"
	"kr/paasta/monitoring-batch/util"
)

type AutoScalerDaoStruct struct {
	txn          *gorm.DB
	influxClient client.Client
	influxDbName string
}

func AutoScalerDao(txn *gorm.DB, influxClient client.Client, influxDbName string) *AutoScalerDaoStruct {
	return &AutoScalerDaoStruct{
		influxClient: influxClient,
		txn:          txn,
		influxDbName: influxDbName,
	}
}

func (p AutoScalerDaoStruct) GetAutoScalePolicy() ([]model.AppAutoScalingPolicy, base.ErrMessage) {

	var listAutoScalePolicy []model.AppAutoScalingPolicy
	status := p.txn.Debug().Table("app_auto_scaling_policies").Find(&listAutoScalePolicy)
	//Where("auto_scaling_out_yn = 'Y' OR auto_scaling_in_yn = 'N'").Find(&listAutoScalePolicy)
	err := util.GetError().DbCheckError(status.Error)
	return listAutoScalePolicy, err
}

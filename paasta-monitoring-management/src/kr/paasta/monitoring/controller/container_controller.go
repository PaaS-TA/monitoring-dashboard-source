package controller

import (
	"net/http"
	"github.com/jinzhu/gorm"
	client "github.com/influxdata/influxdb/client/v2"
	"kr/paasta/monitoring/util"
	"kr/paasta/monitoring/service"
)

//Gorm Object Struct
type containerService struct {
	txn   *gorm.DB
	influxClient 	client.Client
}

func GetContainerController(txn *gorm.DB, influxClient client.Client) *containerService {
	return &containerService{
		txn:   txn,
		influxClient: 	influxClient,
	}
}



func (h *containerService) GetContainerDeploy(w http.ResponseWriter, r *http.Request){

	containerDeployList, err := service.GetContainerService(h.txn, h.influxClient).GetContainerDeploy()
	if err != nil {
		util.RenderJsonResponse(err, w)
	}
	util.RenderJsonResponse(containerDeployList, w)
}
package service

import (
	/*"gopkg.in/olivere/elastic.v3"*/
	iaasmodel "monitoring-portal/iaas_new/model"
	"monitoring-portal/paas/dao"
	"monitoring-portal/paas/model"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
)

type PaasLogService struct {
	elasticClient *elasticsearch.Client
}

func GetPaasLogService(elasticClient *elasticsearch.Client) *PaasLogService {
	return &PaasLogService{
		elasticClient: elasticClient,
	}
}

func (log PaasLogService) GetDefaultRecentLog(request model.LogMessage, paging bool) (model.LogMessage, iaasmodel.ErrMessage) {

	//최근 로그 조회
	return dao.GetPaasLogDao(log.elasticClient).GetDefaultRecentLog(request, paging)
}

func (log PaasLogService) GetSpecificTimeRangeLog(request model.LogMessage, paging bool) (model.LogMessage, iaasmodel.ErrMessage) {

	//특정 시간대 로그 조회
	return dao.GetPaasLogDao(log.elasticClient).GetSpecificTimeRangeLog(request, paging)
}

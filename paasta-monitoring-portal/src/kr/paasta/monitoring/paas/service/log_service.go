package service

import (
	"gopkg.in/olivere/elastic.v3"
	iaasmodel "kr/paasta/monitoring/iaas/model"
	"kr/paasta/monitoring/paas/dao"
	"kr/paasta/monitoring/paas/model"
)

type PaasLogService struct {
	elasticClient *elastic.Client
}

func GetPaasLogService(elasticClient *elastic.Client) *PaasLogService {
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

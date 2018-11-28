package services

import (
	"gopkg.in/olivere/elastic.v3"
	"kr/paasta/monitoring/iaas/dao"
	"kr/paasta/monitoring/iaas/model"
)

type LogServiceStruct struct {
	elasticClient *elastic.Client
}

func GetLogService(elasticClient *elastic.Client) *LogServiceStruct{
	return &LogServiceStruct{
		elasticClient: 	elasticClient,
	}
}

func (log LogServiceStruct) GetDefaultRecentLog(request model.LogMessage, paging bool) (model.LogMessage, model.ErrMessage) {

	//최근 로그 조회
	return dao.GetLogDao(log.elasticClient).GetDefaultRecentLog(request, paging)
}

func (log LogServiceStruct) GetSpecificTimeRangeLog(request model.LogMessage, paging bool) (model.LogMessage, model.ErrMessage) {

	//특정 시간대 로그 조회
	return dao.GetLogDao(log.elasticClient).GetSpecificTimeRangeLog(request, paging)
}

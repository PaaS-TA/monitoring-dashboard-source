package service

import (
	/*"gopkg.in/olivere/elastic.v3"*/
	"kr/paasta/monitoring/iaas_new/dao"
	"kr/paasta/monitoring/iaas_new/model"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
)

type LogServiceStruct struct {
	elasticClient *elasticsearch.Client
}

func GetLogService(elasticClient *elasticsearch.Client) *LogServiceStruct {
	return &LogServiceStruct{
		elasticClient: elasticClient,
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

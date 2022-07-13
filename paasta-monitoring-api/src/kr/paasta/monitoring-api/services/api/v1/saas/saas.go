package saas

import (
	"github.com/tidwall/gjson"
	"paasta-monitoring-api/helpers"
	models "paasta-monitoring-api/models/api/v1"
)

type SaasService struct {
	SaaS models.SaaS
}

func GetSaasService(saas models.SaaS) *SaasService {
	return &SaasService{
		SaaS: saas,
	}
}


func (service *SaasService) GetApplicationStatus() (map[string]int, error) {
	result := make(map[string]int)

	resultBytes, _ := helpers.RequestHttpGet(service.SaaS.PinpointWebUrl+"/getAgentList.pinpoint", "","")
	resultJson := gjson.Parse(string(resultBytes))

	resultJson.ForEach(func(key, value gjson.Result) bool {
		result["agentCount"]++
		for _, item := range value.Array() {
			statusCode := item.Get("status.state.code").Int()
			switch statusCode {
			case 100:
				result["Running"]++
			case 200:
			case 201:
				result["shutdown"]++
			case 300:
				result["disconnect"]++
			}
		}
		return true
	})

	return result, nil
}
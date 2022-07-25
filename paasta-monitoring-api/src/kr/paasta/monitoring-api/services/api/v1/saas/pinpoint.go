package saas

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"paasta-monitoring-api/helpers"
	models "paasta-monitoring-api/models/api/v1"
	"strconv"
	"time"
)

type PinpointService struct {
	SaaS models.SaaS
}

func GetPinpointService(saas models.SaaS) *PinpointService {
	return &PinpointService{
		SaaS: saas,
	}
}


func (service *PinpointService) GetAgentList(ctx echo.Context) (map[string]interface{}, error){
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	result := make(map[string]interface{})
	resultBytes, err := helpers.RequestHttpGet(service.SaaS.PinpointWebUrl+"/getAgentList.pinpoint", "","")
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	json.Unmarshal(resultBytes, &result)
	return result, nil
}


func (service *PinpointService) GetAgentStat(ctx echo.Context) (map[string]interface{}, error){
	logger := ctx.Request().Context().Value("LOG").(*logrus.Entry)

	result := make(map[string]interface{})

	chartType := ctx.Param("chartType")
	agentId := ctx.QueryParam("agentId")
	period := ctx.QueryParam("period")
	periodNum, err := strconv.Atoi(period[0:1])
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	periodUnit := period[1:2]
	switch periodUnit {
	case "m" :
		periodNum = periodNum
	case "h" :
		periodNum = 60*periodNum
	case "d" :
		periodNum = 1400*periodNum
	}

	from := strconv.FormatInt(time.Now().Add(time.Duration(-periodNum)*time.Minute).UTC().Unix(), 10) + "000"
	to := strconv.FormatInt(time.Now().UTC().Unix(), 10) + "000"

	queryString := "agentId="+agentId+"&from="+from +"&to="+to
	resultBytes, resultErr := helpers.RequestHttpGet(service.SaaS.PinpointWebUrl+"/getAgentStat/"+chartType+"/chart.pinpoint", queryString,"")
	if resultErr != nil {
		logger.Error(resultErr)
		return nil, resultErr
	}

	json.Unmarshal(resultBytes, &result)
	return result, nil
}

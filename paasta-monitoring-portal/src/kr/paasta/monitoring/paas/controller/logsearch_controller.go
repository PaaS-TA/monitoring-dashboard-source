package controller

import (
	"github.com/influxdata/influxdb1-client/v2"
	"kr/paasta/monitoring/paas/model"
	"kr/paasta/monitoring/paas/service"
	"kr/paasta/monitoring/paas/util"
	"kr/paasta/monitoring/utils"
	"net/http"
)

type LogsearchController struct {
	client    client.Client
	databases model.Databases
}

func GetLogsearchController(client client.Client, databases model.Databases) *LogsearchController {
	return &LogsearchController{
		client:    client,
		databases: databases,
	}
}

func (h LogsearchController) GetLogData(w http.ResponseWriter, r *http.Request) {

	var param model.NewLogMessage

	param.Id = r.URL.Query().Get("id")
	param.Keyword = r.URL.Query().Get("keyword")
	param.TargetDate = r.URL.Query().Get("targetDate")
	param.StartTime = r.URL.Query().Get("startTime")
	param.EndTime = r.URL.Query().Get("endTime")
	period := r.URL.Query().Get("period")
	param.Period = period
	/*
	if period != "" {
		time_unit := period[len(period)-1:]
		var periodNum int64
		if time_unit == "h" {
			periodNum, _ = strconv.ParseInt(period[:len(period)-1], 10, 64)
			periodNum = periodNum * 60
		} else if time_unit == "m" {
			periodNum, _ = strconv.ParseInt(period[:len(period)-1], 10, 64)
		} else {
			errMessage := map[string]interface{}{"Persons": "Time unit is only allowed 'm' and 'h' -(ex) 5m, 5h"}
			utils.RenderJsonResponse(errMessage, w)
		}

		now := time.Now().Local()
		current := now.Unix()
		before := now.Unix() - periodNum*60 //화면에서 설정한 조회주기(분) (ex: 30 * 60 seconds)
		param.StartTime = strconv.Itoa(int(before))
		param.EndTime = strconv.Itoa(int(current))

	}
	*/

	//service호출 (Gorm Obj 매개 변수)
	result, err := service.GetLogsearchService(h.client, h.databases).GetLogData(param)
	if err != nil {
		utils.Logger.Error(err)
		util.RenderJsonResponse(err, w)
	} else {
		util.RenderJsonResponse(result, w)
	}

}
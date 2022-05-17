package history

import (
	"github.com/cavaliercoder/go-zabbix"
	"log"
	"iaas-monitoring-batch/zabbix-client/utils"
)


var isDebug bool

func init() {
	isDebug = false
}

/**
getHistory
	parameters
		- itemId (string) : 조회할 item의 ID
		- offset (int)    : 불러올 데이터 개수
		- type (int)      : 데이터 타입
*/
func GetHistory(session *zabbix.Session, params map[string]interface{}) ([]zabbix.History, error) {
	var historyParams zabbix.HistoryGetParams

	itemArr := make([]string, 1)
	itemId, ok := params["itemId"].(string)
	if ok {
		itemArr[0] = itemId
		historyParams.ItemIDs = itemArr
	}

	offset, ok := params["offset"].(int)
	if ok {
		historyParams.ResultLimit = offset
	} else {
		historyParams.ResultLimit = 10
	}

	itemType, ok := params["itemType"].(int)
	if ok {
		historyParams.History = itemType
	}

	historyParams.SortField = zabbix.SelectFields{"clock"}
	historyParams.SortOrder = "ASC"

	result, err := session.GetHistories(historyParams)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	if isDebug {
		utils.PrintJson(result)
	}
	return result, err
}

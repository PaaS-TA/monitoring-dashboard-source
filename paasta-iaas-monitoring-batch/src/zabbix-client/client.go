package zabbix_client

import (
	"fmt"
	"github.com/cavaliercoder/go-zabbix"
	"strconv"
	"zabbix-client/history"
	"zabbix-client/host"
	"zabbix-client/item"
)

func GetHostInfo(session *zabbix.Session, ipAddr string) (zabbix.Host, error) {
	// IP주소로 호스트 정보 조회
	hostParams := make(map[string]interface{}, 0)
	filterMap := make(map[string]interface{}, 0)
	filterMap["ip"] = ipAddr
	hostParams["filter"] = filterMap
	result := host.GetHostList(session, hostParams)

	if len(result) == 0 {
		_result := zabbix.Host{}
		return _result, fmt.Errorf("%s is not exist host.", ipAddr)
	}
	return result[0], nil
}

func GetHistory(session *zabbix.Session, itemKey string, hostId string) []zabbix.History {
	itemParams := make(map[string]interface{}, 0)
	keywordArr := make([]string, 2)
	keywordArr[0] = itemKey
	itemParams["itemKey"] = keywordArr

	hostIds := make([]string, 1)
	hostIds[0] = hostId
	itemParams["hostIds"] = hostIds
	itemResult := item.GetItemList(session, itemParams)

	itemId := strconv.Itoa(itemResult[0].ItemID)
	itemType := itemResult[0].LastValueType

	params := make(map[string]interface{}, 0)
	params["itemId"] = itemId
	params["itemType"] = itemType
	params["offset"] = 10
	result := history.GetHistory(session, params)
	return result
}
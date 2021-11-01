package item

import (
	"fmt"
	"github.com/cavaliercoder/go-zabbix"
	"log"
)

var isDebug bool

func init() {
	isDebug = false
}

/**
Item (수집 항목) 정보 조회
Parameters
	- params [map]
		hostIds ([]string) : 호스트 ID
		itemKey (string)   : item의 key값으로 검색 (와일드카드 사용가능)

*/
func GetItemList(session *zabbix.Session, params map[string]interface{}) []zabbix.Item {

	var itemParams zabbix.ItemGetParams

	// 2021.10.25 - Host의 IP 정보도 가져올 수 있도록 추가함
	itemParams.SelectInterfaces = zabbix.SelectFields{"ip"}
	
	filterMap := make(map[string]interface{}, 0)
	searchMap := make(map[string]interface{}, 0)


	itemParams.Filter = filterMap

	hostIds, ok := params["hostIds"]
	if ok {
		itemParams.HostIDs = hostIds.([]string)
	}

	itemKey, ok := params["itemKey"]
	if ok {

		searchMap["key_"] = itemKey.([]string)
	}

	itemIds, ok := params["itemIds"]
	if ok {
		itemParams.ItemIDs = itemIds.([]string)
	}

	itemParams.Filter = filterMap
	itemParams.TextSearch = searchMap
	itemParams.EnableTextSearchWildcards = true

	result, err := session.GetItems(itemParams)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	//utils.PrintJson(result)

	if isDebug {
		for idx, item := range result {
			fmt.Printf("[%d] %d (%s) : %s - %s\n", idx, item.ItemID, item.ItemName, item.LastValue, item.ItemDescr)
		}
	}


	return result
}

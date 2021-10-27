package hostgroup

import (
	"github.com/cavaliercoder/go-zabbix"
	"log"
	"zabbix-client/utils"
)

var isDebug bool

func init() {
	isDebug = false
}

func GetHostgroup(session *zabbix.Session, params map[string]interface{}) []zabbix.Hostgroup {
	var hostgroupParams zabbix.HostgroupGetParams

	name, ok := params["name"]
	if ok {
		textsearchMap := make(map[string]interface{}, 0)
			keywordArr := make([]string, 1)
			keywordArr[0] = name.(string)
		textsearchMap["name"] = keywordArr
		hostgroupParams.TextSearch = textsearchMap
	}

	result, err := session.GetHostgroups(hostgroupParams)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	if isDebug {
		utils.PrintJson(result)
	}

	return result
}
package config

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	commonModel "kr/paasta/monitoring/common/model"
	"os"
	"strings"
)

func ImportConfig(fileName string) (map[string]string, error) {
	config := make(map[string]string, 0)
	/*
		config := Config{
			"server.ip":   "127.0.0.1",
			"server.port": "8888",
		}
	*/

	if len(fileName) == 0 {
		return config, nil
	}
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		// check if the line has = sign
		// and process the line. Ignore the rest.
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				// assign the config map
				config[key] = value
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}


func ConvertXmlToString(fileName string) (string, error) {
	xmlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "", err
	}
	return string(xmlFile), nil
}


func InitDBConnectionConfig(configMap map[string]string) commonModel.DBConfig {
	dbConnInfo := commonModel.DBConfig{}
	dbConnInfo.DbType = configMap["paas.monitoring.db.type"]
	dbConnInfo.DbName = configMap["paas.monitoring.db.dbname"]
	dbConnInfo.UserName = configMap["paas.monitoring.db.username"]
	dbConnInfo.UserPassword = configMap["paas.monitoring.db.password"]
	dbConnInfo.Host = configMap["paas.monitoring.db.host"]
	dbConnInfo.Port = configMap["paas.monitoring.db.port"]

	return dbConnInfo
}
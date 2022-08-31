package config

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Configuration struct {
	DbType string
	DbName string
	DbUser string
	DbPw string
	DbHost string
	DbPort string
	CaasApiUrl string
	SmtpHost string
	SmtpPort int
	Sender string
	SenderPw string
	MailResourceUrl string
}

var config Configuration

func init() {
	configData, err := readConfig(`config.ini`)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	smtpPortNum, _ := strconv.Atoi(configData["mail.smtp.port"])
	config = Configuration{
		DbType: configData["monitoring.db.type"],
		DbName : configData["monitoring.db.dbname"],
		DbUser : configData["monitoring.db.username"],
		DbPw : configData["monitoring.db.password"],
		DbHost : configData["monitoring.db.host"],
		DbPort : configData["monitoring.db.port"],
		CaasApiUrl : configData["caas.monitoring.api.url"],
		SmtpHost : configData["mail.smtp.host"],
		SmtpPort : smtpPortNum,
		Sender : configData["mail.sender"],
		SenderPw : configData["mail.sender.password"],
		MailResourceUrl: configData["mail.resource.url"],
	}
}

func GetDBConnectionStr() string {
	connectionString := fmt.Sprintf("%s:%s@%s([%s]:%s)/%s%s", config.DbUser, config.DbPw, "tcp", config.DbHost, config.DbPort, config.DbName, "?charset=utf8&parseTime=true")
	log.Println(connectionString)
	return connectionString
}

func GetConfiguration() Configuration {
	return config
}

func readConfig(filename string) (map[string]string, error) {
	// init with some bogus data
	config := make(map[string]string)

	if len(filename) == 0 {
		return config, nil
	}
	file, err := os.Open(filename)

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
package config

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)


type Config struct {
	DbType string
	DbName string
	DbUser string
	DbPasswd string
	DbHost string
	DbPort string

	SmtpHost     string
	Port         string
	MailSender   string
	SenderPwd    string
	ResouceUrl   string
	MailReceiver string
	AlarmSend    bool
	MailTlsSend  bool

	ZabbixHost string
	ZabbixAdminId string
	ZabbixAdminPw string

	GmtTimeGapHour int64
	ExecuteInterval int
}


func InitializeConfig() *Config {
	filePath, _ := filepath.Abs("src/kr/paasta/iaas-monitoring-batch/config.ini")
	configMap, err := readConfig(filePath)
	if err != nil {
		os.Exit(0)
	}

	//Monitoring configDB Configuration
	config := new(Config)
	config.DbType = configMap["monitoring.db.type"]
	config.DbName = configMap["monitoring.db.dbname"]
	config.DbUser = configMap["monitoring.db.username"]
	config.DbPasswd = configMap["monitoring.db.password"]
	config.DbHost   = configMap["monitoring.db.host"]
	config.DbPort   = configMap["monitoring.db.port"]

	isAlarmSend, _ := strconv.ParseBool(configMap["mail.alarm.send"])
	isMailTlsSend, _ := strconv.ParseBool(configMap["mail.tls.send"])
	config.SmtpHost = configMap["mail.smtp.host"]
	config.Port = configMap["mail.smtp.port"]
	config.MailSender = configMap["mail.sender"]
	config.SenderPwd = configMap["mail.sender.password"]
	config.ResouceUrl = configMap["mail.resource.url"]
	config.AlarmSend = isAlarmSend
	config.MailTlsSend = isMailTlsSend

	gmtTimeGapHour,  _ := strconv.ParseInt(configMap["gmt.time.hour.gap"], 10, 64)
	executeInterval, _ := strconv.Atoi(configMap["batch.interval.second"])
	config.GmtTimeGapHour = gmtTimeGapHour
	config.ExecuteInterval = executeInterval

	/*
	redisConfig := new(service.RedisConfig)
	redisConfig.RedisAddr = config["redis.addr"]
	redisConfig.RedisPassword = config["redis.password"]

	model.PortalUrl = config["portal.api.url"]
	model.PortalClient = util.NewPortalClient()
	*/

	config.ZabbixHost = configMap["zabbix.host"]
	config.ZabbixAdminId = configMap["zabbix.admin.id"]
	config.ZabbixAdminPw = configMap["zabbix.admin.pw"]

	return config
}


func readConfig(filename string) (map[string]string, error) {
	// init with some bogus data
	config := make(map[string]string, 0)
	config["server.port"] = "9999"

	if len(filename) == 0 {
		return config, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
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
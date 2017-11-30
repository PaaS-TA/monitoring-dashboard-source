package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"encoding/json"
	"kr/paasta/monitoring/handler"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	client "github.com/influxdata/influxdb/client/v2"
	//"com/crossent/monitoring/datasource"
	"kr/paasta/monitoring/util"
	"strconv"
	"kr/paasta/monitoring/domain"
)

type Config map[string]string

type DBConfig struct {
	DbType string
	UserName string
	UserPassword string
	Host string
	Port string
	DbName string
}

func main() {
	//var dbconfig *datasource.DBConfig
	var handlers http.Handler

	fmt.Println("##### Monitoring Management start!!!")
	//============================================
	// Sample VCAP_SERVICE INFO Parsing
/*	os_args := os.Getenv("VCAP_SERVICES")
	readOSEnvironment(os_args)*/
	//============================================

	//============================================
	// 기본적인 프로퍼티 설정 정보 읽어오기
	config, err := ReadConfig(`config.ini`)
	if err != nil {
		fmt.Println(err)
	}
	//============================================


	//maxConnection, err := strconv.Atoi(config["mysql.maxconn"])
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	apiPort, _ := strconv.Atoi(config["server.port"])

	configDbCon := new(DBConfig)
	configDbCon.DbType        = config["monitoring.db.type"]
	configDbCon.DbName        = config["monitoring.db.dbname"]
	configDbCon.UserName      = config["monitoring.db.username"]
	configDbCon.UserPassword  = config["monitoring.db.password"]
	configDbCon.Host          = config["monitoring.db.host"]
	configDbCon.Port          = config["monitoring.db.port"]

	connectionString := util.GetConnectionString(configDbCon.Host , configDbCon.Port, configDbCon.UserName, configDbCon.UserPassword, configDbCon.DbName )

	fmt.Println("String:",connectionString)

	dbAccessObj, dbErr := gorm.Open(configDbCon.DbType, connectionString + "?charset=utf8&parseTime=true")

	if dbErr != nil{
		fmt.Println("err::",err)
	}

	url     ,  _ := config["metric.db.url"]
	userName,  _ := config["metric.db.username"]
	password,  _ := config["metric.db.password"]

	InfluxServerClient, _ := client.NewHTTPClient(client.HTTPConfig{
		Addr: url,
		Username: userName,
		Password: password,
	})

	/**
	Newly Added - 2017.08.14

	Get InfluxDB Database Name
	 Bosh - metric.infra.db_name
	 PaaSTA - metric.controller.db_name
	 Container - metric.container.db_name
	*/
	bosh_database, _ := config["metric.infra.db_name"]
	paasta_database, _ := config["metric.controller.db_name"]
	container_database, _ := config["metric.container.db_name"]

	var databases domain.Databases
	databases.BoshDatabase = bosh_database
	databases.PaastaDatabase = paasta_database
	databases.ContainerDatabase = container_database

	timeGap, _ := strconv.Atoi(config["gmt.time.hour.gap"])

	domain.GmtTimeGap = timeGap
	// Web Resource 핸들링
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))

	// Route Path 정보와 처리 서비스 연결
	handlers = handler.NewHandler(dbAccessObj, InfluxServerClient, databases)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", apiPort), handlers); err != nil {
		log.Fatalln(err)
	}
}

// Config 파일 읽어 오기
func ReadConfig(filename string) (Config, error) {
	// init with some bogus data
	config := Config{
		"server.ip":     "127.0.0.1",
		"server.port":   "8888",
		"mysql.dburl":   "",
		"mysql.userid":  "",
		"mysql.userpwd": "",
		"mysql.maxconn": "",
	}

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

// cf app 배포시 시스템 환경 정보 읽어 오기
func readOSEnvironment(service string) (string, string) {
	args := os.Getenv("VCAP_SERVICES")
	var uri, name string
	if args != "" {
		var vcap_env map[string]interface{}
		if err := json.Unmarshal([]byte(args), &vcap_env); err != nil {
			log.Panic(err.Error())
		}

		//Service instance (not name) - for example : p-mysql - type check !!! []interface{}
		if vcap_env[service] != nil {
			sub_env := vcap_env[service].([]interface{})
			credentials := (sub_env[0].(map[string]interface{}))["credentials"].(map[string]interface{})
			uri = credentials["uri"].(string)
			if credentials["name"] != nil {
				name = credentials["name"].(string)
			}
		}
	}
	return uri, name
}

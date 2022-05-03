package connections

import (
	"GoEchoProject/helpers"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"os"
	"strings"
)

type Connections struct {
	DbInfo    *gorm.DB
	RedisInfo *redis.Client
}

func SetEnv(env map[string]string) {

	// Redis 설정
	env["redis_url"] = os.Getenv("redis_url")

	// PaaS DataBase 설정
	env["paas_db_type"] = os.Getenv("paas_db_type")
	env["paas_db_password"] = os.Getenv("paas_db_password")
	env["paas_db_username"] = os.Getenv("paas_db_username")
	env["paas_db_protocol"] = os.Getenv("paas_db_protocol")
	env["paas_db_host"] = os.Getenv("paas_db_host")
	env["paas_db_port"] = os.Getenv("paas_db_port")
	env["paas_db_name"] = os.Getenv("paas_db_name")
	env["paas_db_charset"] = os.Getenv("paas_db_charset")
	env["paas_db_parseTime"] = os.Getenv("paas_db_parseTime")
}

func RedisConnection(env map[string]string) *redis.Client {

	// Redis 설정
	dsn := env["redis_url"]
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		panic(err)
	}
	return redisClient
}

func PaaSConnection(env map[string]string) *gorm.DB {

	// DB 설정
	paasDBClient, paasDBErr := gorm.Open(
		helpers.GetDBConnectionString(
			env["paas_db_type"],
			env["paas_db_username"],
			env["paas_db_password"],
			env["paas_db_protocol"],
			env["paas_db_host"],
			env["paas_db_port"],
			env["paas_db_name"],
			env["paas_db_charset"],
			env["paas_db_parseTime"]),
	)

	if paasDBErr != nil {
		panic(paasDBErr)
	}

	return paasDBClient
}

func SaaSConnection(env map[string]string) error {
	return nil
}

func CaaSConnection(env map[string]string) error {
	return nil
}

func IaaSConnection(env map[string]string) error {
	return nil
}

func SetupConnection() Connections {

	conn := Connections{}

	// 환경 변수 설정
	env := make(map[string]string)
	SetEnv(env)

	// 내부 서비스 환경 설정
	conn.RedisInfo = RedisConnection(env)

	services := os.Getenv("services")
	servicesArr := strings.Split(services, ",")

	// 외부 서비스 별 환경 설정
	for _, value := range servicesArr {
		switch value {
		case "PaaS":
			conn.DbInfo = PaaSConnection(env)
		case "SaaS":
			SaaSConnection(env)
		case "CaaS":
			CaaSConnection(env)
		case "IaaS":
			IaaSConnection(env)
		}
	}
	return conn
}

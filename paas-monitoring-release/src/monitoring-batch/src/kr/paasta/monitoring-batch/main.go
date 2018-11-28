package main

import (
	"fmt"
	"os"
	"time"
	"bufio"
	"strings"
	"io"
	"os/signal"
	"syscall"
	"strconv"
	"github.com/tedsuo/ifrit/grouper"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/sigmon"
	"kr/paasta/monitoring-batch/service"
	"kr/paasta/monitoring-batch/util"
	"kr/paasta/monitoring-batch/model"
	"kr/paasta/monitoring-batch/handler"
)

type Config map[string]string

func main() {

	//debugserver.AddFlags(flag.CommandLine)
	//cflager.AddFlags(flag.CommandLine)
	//===============================================================
	//catch or finally

	defer func() { //catch or finally
		if err := recover(); err != nil { //catch
			fmt.Fprintf(os.Stderr, "Main: Exception: %v\n", err)
			os.Exit(1)
		}
	}()
	//===============================================================
	//logger, _ := cflager.New("paasta-monitoring-batch")

	var startTime time.Time
	//============================================
	// 기본적인 프로퍼티 설정 정보 읽어오기
	config, err := readConfig(`config.ini`)
	if err != nil {
		//logger.Fatal("read config file error :", err)
		os.Exit(0)
	}
	//============================================

	//============================================
	// Channel for Singal Checkig
	sigs := make(chan os.Signal, 1)
	//Waiting to be notified
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println("returned signal:", sig)
		elapsed := time.Since(startTime)
		fmt.Println("# ElapsedTime in seconds:", elapsed)

		//When unexpected signal happens, defer function doesn't work.
		//So, go func has a role to be notified signal and do defer function execute
		os.Exit(0)
	}()
	//============================================


	//==============================================
	//Influx Configuration
	influxCon := new(service.InfluxConfig)
	influxCon.InfluxUrl  = config["influx.url"]
	influxCon.InfluxUser = config["influx.user"]
	influxCon.InfluxPass = config["influx.pass"]
	influxCon.InfraDatabase     = config["influx.bosh.db_name"]
	influxCon.PaastaDatabase    = config["influx.paasta.db_name"]
	influxCon.ContainerDatabase = config["influx.container.db_name"]
	influxCon.DefaultTimeRange = config["influx.defaultTimeRange"]
	//==============================================

	//==============================================
	//Monitoring configDB Configuration
	configDbCon := new(service.DBConfig)
	configDbCon.DbType = config["monitoring.db.type"]
	configDbCon.DbName = config["monitoring.db.dbname"]
	configDbCon.UserName      = config["monitoring.db.username"]
	configDbCon.UserPassword  = config["monitoring.db.password"]
	configDbCon.Host          = config["monitoring.db.host"]
	configDbCon.Port          = config["monitoring.db.port"]
	//==============================================

	//==============================================
	//Monitoring configDB Configuration
	portalDbCon := new(service.DBConfig)
	portalDbCon.DbType = config["portal.db.type"]
	portalDbCon.DbName = config["portal.db.dbname"]
	portalDbCon.UserName      = config["portal.db.username"]
	portalDbCon.UserPassword  = config["portal.db.password"]
	portalDbCon.Host          = config["portal.db.host"]
	portalDbCon.Port          = config["portal.db.port"]
	//==============================================

	//==============================================
	//configDB Configuration
	boshCon := new(service.BoshConfig)
	boshCon.BoshUrl  = config["bosh.api.url"]
	boshCon.BoshIp  = config["bosh.ip"]
	boshCon.BoshId   = config["bosh.admin"]
	boshCon.BoshPass = config["bosh.password"]
	boshCon.CfDeploymentName    = config["bosh.cf.deployment.name"]
	boshCon.CellNamePrefix      = config["bosh.cell.name.prefix"]
	boshCon.ServiceName         = config["bosh.service.name"]
	//==============================================

	mailConfig := new(service.MailConfig)
	mailConfig.SmtpHost   = config["mail.smtp.host"]
	mailConfig.Port       = config["mail.smtp.port"]
	mailConfig.MailSender = config["mail.sender"]
	mailConfig.SenderPwd  = config["mail.sender.password"]
	mailConfig.ResouceUrl = config["mail.resource.url"]
	//mailConfig.MailReceiver = config["mail.receiver.email"]
	isAlarmSend, _ := strconv.ParseBool(config["mail.alarm.send"])
	mailConfig.AlarmSend    = isAlarmSend
	//tls send
	isMailTlsSend, _ := strconv.ParseBool(config["mail.tls.send"])
	mailConfig.MailTlsSend    = isMailTlsSend

	gmtTimeGapHour,  _ := strconv.ParseInt(config["gmt.time.hour.gap"], 10, 64)
	//======================= Save Process ID to .pid file ==================
	batInterval, _ := strconv.Atoi(config["batch.interval.second"])

	//redis
	redisConfig := new(service.RedisConfig)
	redisConfig.RedisAddr = config["redis.addr"]
	redisConfig.RedisPassword = config["redis.password"]

	//portal
	model.PortalUrl = config["portal.api.url"]
	model.PortalClient = util.NewPortalClient()

	pid := os.Getpid()
	f, err := os.Create(".pid")
	defer f.Close()
	if err != nil {
		fmt.Println("Main: failt to create pid file.", err.Error())
		panic(err)
	}
	f.WriteString(strconv.Itoa(pid))
	//=======================================================================

	//logger.Info("##### process id :", lager.Data{"process_id ":pid})

	//fmt.Println("gmtTimeGapHour::::", gmtTimeGapHour)

	members := grouper.Members{
		{"monitoring_batch_processor", handler.New( batInterval, gmtTimeGapHour, influxCon, configDbCon, portalDbCon, boshCon, redisConfig, config, mailConfig)},
	}

	fmt.Println("#monitoring batch processor started")
	group := grouper.NewOrdered(os.Interrupt, members)
	monitor := ifrit.Invoke(sigmon.New(group))
	monit_err := <-monitor.Wait()

	if monit_err != nil {
		fmt.Println("#Main: exited-with-failure", monit_err)
		panic(monit_err)
	}
	elapsed := time.Since(startTime)
	fmt.Println("#ElapsedTime in seconds:", map[string]interface{}{"elapsed_time": elapsed, })
	fmt.Println("#monitoring batch processor exited")
}



func readConfig(filename string) (Config, error) {
	// init with some bogus data
	config := Config{
		"server.port": "9999",
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

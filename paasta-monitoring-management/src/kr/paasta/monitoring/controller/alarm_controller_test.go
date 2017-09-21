package controller_test

import (
	"io/ioutil"
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
	. "github.com/onsi/ginkgo"
	"io"
	"net/url"
	"bytes"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	cb "kr/paasta/monitoring/domain"
	"fmt"
	"time"
	"os"
	"bufio"
	"strings"
	"strconv"
	"encoding/json"
	"net/http/httptest"
	"github.com/influxdata/influxdb/client/v2"
	"kr/paasta/monitoring/handler"
)

type Response struct {
	Content string
	Code    int
}
type AlarmDao struct {
	txn   *gorm.DB
}
type Config map[string]string

var (
	server  *httptest.Server
	testUrl string
	t *testing.T
	dbAccessObj *gorm.DB
	alarm *cb.Alarm
	alarmActionId uint
)
var _ = Describe("AlarmController", func() {
	var almFormValues url.Values
	var almActFormValues url.Values

	BeforeSuite(func() {
		config, err := readConfig(`../config.ini`)
		if err != nil {
			fmt.Errorf("read config file error: %s", err)
			os.Exit(0)
		}

		testUrl = config["server.url"]

		// alarms 데이터 생성
		originId := uint(9999)
		almFormValues = make(url.Values)
		almFormValues.Add("origin_type", "pas")
		almFormValues.Add("alarm_type", "cpu")
		almFormValues.Add("level", "critical")
		almFormValues.Add("resolve_status", "1")
		almFormValues.Add("app_yn", "N")
		almFormValues.Add("ip", "10.244.0.0")
		almFormValues.Add("alarm_title", "test")
		almFormValues.Add("alarm_message", "This is test data.")
		almFormValues.Add("resolve_status", "1")

		DbType := config["monitoring.db.type"]
		DbName := config["monitoring.db.dbname"]
		UserName := config["monitoring.db.username"]
		UserPassword := config["monitoring.db.password"]
		Host := config["monitoring.db.host"]
		Port := config["monitoring.db.port"]

		var dbErr error
		dbAccessObj, dbErr = gorm.Open(DbType, fmt.Sprintf("%s:%s@%s([%s]:%s)/%s%s",
			UserName, UserPassword, "tcp", Host, Port, DbName, "") + "?charset=utf8&parseTime=true")
		if dbErr != nil{
			fmt.Errorf("Open error: %s", dbErr)
		}

		eventData := cb.Alarm{OriginId: originId, OriginType: almFormValues.Get("origin_type"), AlarmType: almFormValues.Get("alarm_type"), Level: almFormValues.Get("level"),
			AppYn: almFormValues.Get("app_yn"), Ip: almFormValues.Get("ip"), AlarmTitle: almFormValues.Get("alarm_title"), AlarmMessage: almFormValues.Get("alarm_message"),
			ResolveStatus: almFormValues.Get("resolve_status"), AlarmCnt: 1, AlarmSendDate: GetInsertCurrentTime()}
		status := dbAccessObj.Debug().Create(&eventData)
		alarm = status.Value.(*cb.Alarm)
		fmt.Println("[Test] alarm id: ", alarm.Id)
		if status.Error != nil{
			fmt.Errorf("Create error: %s", status.Error)
		}


		url     ,  _ := config["metric.db.url"]
		userName,  _ := config["metric.db.username"]
		password,  _ := config["metric.db.password"]

		InfluxServerClient, _ := client.NewHTTPClient(client.HTTPConfig{
			Addr: url,
			Username: userName,
			Password: password,
		})

		bosh_database, _ := config["metric.infra.db_name"]
		paasta_database, _ := config["metric.controller.db_name"]
		container_database, _ := config["metric.container.db_name"]

		var databases cb.Databases
		databases.BoshDatabase = bosh_database
		databases.PaastaDatabase = paasta_database
		databases.ContainerDatabase = container_database

		var handlers http.Handler
		handlers = handler.NewHandler(dbAccessObj, InfluxServerClient, databases)
		server = httptest.NewServer(handlers);
		testUrl = server.URL
	})

	AfterSuite(func() {
		// alarms 데이터 제거
		status := dbAccessObj.Debug().Table("alarms").Where("id = ? ", alarm.Id).Delete(alarm)
		if status.Error != nil{
			fmt.Errorf("Delete error: %s", status.Error)
		}

		var alarmAction cb.AlarmAction
		alarmAction.AlarmId = alarm.Id
		actionStatus := dbAccessObj.Debug().Table("alarm_actions").Where("alarm_id = ? ", alarmAction.AlarmId).Delete(alarmAction)
		if actionStatus.Error != nil{
			fmt.Errorf("Action Delete error: %s", actionStatus.Error)
		}
	})

	Describe("Alarms", func() {
		Context("Get", func() {
			It("all list", func() {
				res, err := DoGet(testUrl + "/alarms")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("non existent paths", func() {
				res, err := DoGet(testUrl + "/non-exist-path")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusNotFound, res.Code, "404 Not Found")
			})

			It("list by resolve status", func() {
				res, err := DoGet(testUrl + "/alarms/status/1")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("list by origin type", func() {
				res, err := DoGet(testUrl + "/alarms?originType=pas")
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, res.Code)
			})

			It("not allowed method", func() {
				req, err :=  http.NewRequest("POST", testUrl + "/alarms", nil)
				var client http.Client
				resp, err := client.Do(req)

				assert.Nil(t, err)
				assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

				resp.Body.Close()
			})

			It("detail by alarm id", func() {
				res, err := DoGet(testUrl + "/alarms/" + strconv.Itoa(int(alarm.Id)))
				assert.Nil(t, err, "")
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NotEmpty(t, res.Content)
			})
		})

		Context("Update", func() {
			It("resolve status", func() {
				body := bytes.NewReader([]byte(`{"resolveStatus":"2"}`))
				res, err := Do("PUT", testUrl + "/alarms/" + strconv.Itoa(int(alarm.Id)), body)
				assert.Nil(t, err)
				assert.Equal(t, http.StatusCreated, res.Code, "status code")
				res, err = DoGet(testUrl + "/alarms/" + strconv.Itoa(int(alarm.Id)))

				alarm := &cb.AlarmResponse{}
				json.Unmarshal([]byte(res.Content), alarm)
				assert.Equal(t, "2", alarm.ResolveStatus)
			})

			It("If a required input value does not exist, update will fail.", func() {
				req, err :=  http.NewRequest("PUT", testUrl + "/alarms/" + strconv.Itoa(int(alarm.Id)), nil)
				var client http.Client
				resp, err := client.Do(req)

				assert.Nil(t, err)
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				resp.Body.Close()
			})
		})
	})

	Describe("Alarm Actions", func() {
		BeforeEach(func() {
			almActFormValues = make(url.Values)
			almActFormValues.Add("alarmId", strconv.Itoa(int(alarm.Id)))
		})

		It("regist alarm action", func() {
			body := bytes.NewReader([]byte(fmt.Sprintf("{\"alarmId\":%d}", alarm.Id)))
			res, err := DoPost(testUrl + "/alarmsAction", body)
			//res, err := DoPostForm(testUrl + "/alarmsAction", almActFormValues)
			assert.Nil(t, err)
			assert.Equal(t, http.StatusCreated, res.Code)
		})

		It("If a required input value does not exist, regist will fail.", func() {
			res, err := DoPost(testUrl + "/alarmsAction", nil)
			assert.Nil(t, err)
			assert.Equal(t, http.StatusBadRequest, res.Code)
		})

		It("get alarm actions by alarm id", func() {
			res, err := DoGet(testUrl + "/alarms/" + strconv.Itoa(int(alarm.Id)))
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, res.Code)
			assert.NotEmpty(t, res.Content)

			alarmDetail := &cb.AlarmDetailResponse{}
			json.Unmarshal([]byte(res.Content), alarmDetail)
			assert.True(t, len(alarmDetail.AlarmActionResponse) > 0)

			alarmActionId = alarmDetail.AlarmActionResponse[len(alarmDetail.AlarmActionResponse)-1].Id
		})

		It("update alarm action", func() {
			body := bytes.NewReader([]byte(`{"alarmActionDesc":"Description update test."}`))
			req, err :=  http.NewRequest("PUT", testUrl + "/alarmsAction/" + strconv.Itoa(int(alarmActionId)), body)
			var client http.Client
			resp, err := client.Do(req)

			assert.Nil(t, err)
			assert.Equal(t, http.StatusCreated, resp.StatusCode)

			resp.Body.Close()

			res, err := DoGet(testUrl + "/alarms/" + strconv.Itoa(int(alarm.Id)))
			alarmDetail := &cb.AlarmDetailResponse{}
			json.Unmarshal([]byte(res.Content), alarmDetail)
			desc := alarmDetail.AlarmActionResponse[len(alarmDetail.AlarmActionResponse)-1].AlarmActionDesc
			assert.Equal(t, "Description update test.", desc)
		})

		It("If a necessary value can not be found, alarm action updating has failed.", func() {
			req, err :=  http.NewRequest("PUT", testUrl + "/alarmsAction/0", nil)
			var client http.Client
			resp, err := client.Do(req)

			assert.Nil(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			resp.Body.Close()
		})

		It("alarm action delete", func() {
			res, err := Do("DELETE", testUrl + "/alarmsAction/" + strconv.Itoa(int(alarmActionId)), nil)
			assert.Nil(t, err)
			assert.Equal(t, http.StatusNoContent, res.Code)
		})
	})

	Describe("Alarm Stats", func() {
		It("get alarm stats", func() {
			param := "?period=m&interval=1"
			res, err := DoGet(testUrl + "/alarmsStat" + param)
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, res.Code)

			alarmStats := &cb.AlarmStatResponse{}
			json.Unmarshal([]byte(res.Content), alarmStats)
			assert.True(t, alarmStats.TotalCnt > 0)
		})

		It("get alarm stats with custom period", func() {
			tt := time.Now()
			dateFrom := strconv.Itoa(int(tt.AddDate(0, 0, -1).Unix())*1000)
			dateTo := strconv.Itoa(int(tt.Unix())*1000)

			param := "?period=custom&searchDateFrom="+dateFrom+"&searchDateTo="+dateTo
			res, err := DoGet(testUrl + "/alarmsStat" + param)
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, res.Code)

			alarmStats := &cb.AlarmStatResponse{}
			json.Unmarshal([]byte(res.Content), alarmStats)
			assert.True(t, alarmStats.TotalCnt > 0)
		})

		It("wrong interval", func() {
			param := "?period=m&interval=d"
			res, err := DoGet(testUrl + "/alarmsStat" + param)
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, res.Code)

			alarmStats := &cb.AlarmStatResponse{}
			json.Unmarshal([]byte(res.Content), alarmStats)
			assert.True(t, alarmStats.TotalCnt == 0)
		})

		It("wrong period", func() {
			tt := time.Now()
			dateFrom := strconv.Itoa(int(tt.AddDate(0, 0, -1).Unix())*1000)
			dateTo := strconv.Itoa(int(tt.Unix())*1000)

			param := "?period=custom&searchDateFrom="+dateTo+"&searchDateTo="+dateFrom
			res, err := DoGet(testUrl + "/alarmsStat" + param)
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, res.Code)

			alarmStats := &cb.AlarmStatResponse{}
			json.Unmarshal([]byte(res.Content), alarmStats)
			assert.True(t, alarmStats.TotalCnt == 0)
		})

		It("not allowed method", func() {
			res, err := DoPost(testUrl + "/alarmsStat", nil)
			assert.Nil(t, err)
			assert.Equal(t, http.StatusMethodNotAllowed, res.Code)
		})
	})

})

func DoGet(url string) (*Response, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return &Response{Content: string(contents), Code: response.StatusCode}, nil
}

func DoPost(url string, body io.Reader) (*Response, error) {
	response, err := http.Post(url, "application/json", body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return &Response{Content: string(contents), Code: response.StatusCode}, nil
}

func Do(method string, url string, body io.Reader) (*Response, error) {
	req, err :=  http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	var client http.Client
	response, err := client.Do(req)
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return &Response{Content: string(contents), Code: response.StatusCode}, nil
}

func GetInsertCurrentTime() time.Time{
	now := time.Now()
	t := now.Format("2006-01-02T15:04:05+00:00")
	currentTime, _ := time.Parse(time.RFC3339,t)
	return currentTime
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

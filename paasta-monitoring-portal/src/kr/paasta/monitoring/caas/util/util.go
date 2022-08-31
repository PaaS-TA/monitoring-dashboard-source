package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"io"
	"log"
	"monitoring-portal/caas/model"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	SERVICE_NAME string = "serviceName"
	DATA_NAME    string = "data"
)

type Config map[string]string

func GetConnectionString(host, port, user, pass, dbname string) string {

	return fmt.Sprintf("%s:%s@%s([%s]:%s)/%s%s",
		user, pass, "tcp", host, port, dbname, "")

}

func RenderJsonResponse(data interface{}, w http.ResponseWriter) {

	/* NaN 데이터가 있는경우 json.Marshal 에서 panic 발행후 프로세스가 정지 되므로 사전에 체크하여 우회한다. */
	str := fmt.Sprint(data)

	if strings.Contains(str, "NaN") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(""))
		return
	}

	js, err := json.Marshal(data)
	if err != nil {
		log.Fatalln("Error writing JSON:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(js)
	return
}

type errorMessage struct {
	model.ErrMessage
}

func GetError() *errorMessage {
	return &errorMessage{}
}

func (e errorMessage) DbCheckError(err error) model.ErrMessage {
	if err != nil {
		errMessage := model.ErrMessage{
			"Message": err.Error(),
		}
		return errMessage
	} else {
		return nil
	}
}

//오류 체크 모듈로 오류 발생시 오류 메시지 리턴
func (e errorMessage) CheckError(resp client.Response, err error) (client.Response, model.ErrMessage) {

	if err != nil {
		errMessage := model.ErrMessage{
			"Message": err.Error(),
		}
		return resp, errMessage

	} else if resp.Error() != nil {
		errMessage := model.ErrMessage{
			"Message": resp.Err,
		}
		return resp, errMessage
	} else {

		return resp, nil
	}
}

type clientResponse struct {
	client.Response
}

func GetResponseConverter() *clientResponse {
	return &clientResponse{}
}

func rfc3339ToUnixTimestamp(metricDataTime string) int64 {
	t, _ := time.Parse(time.RFC3339, metricDataTime)
	return t.Unix()
}

//조회한 결과List를 Map으로 변환한다.
func (r clientResponse) InfluxConverterToMap(resp client.Response) ([]map[string]interface{}, model.ErrMessage) {

	if len(resp.Results) != 1 {
		return nil, nil
	} else {
		//UI로 Return할 결과값
		var returnValues []map[string]interface{}
		var columns []string

		for _, v := range resp.Results[0].Series {
			for _, vc := range v.Columns {
				columns = append(columns, vc)
			}

			for i := 0; i < len(v.Values); i++ {
				row := make(map[string]interface{})

				//InfluxDB에서 Value 값이 nil인 경우 해당 row는 값을 보내주지 않는다.
				if v.Values[i][1] != nil {
					for kv, vv := range v.Values[i] {
						if vv != nil {
							//Time Column Case convert to UnixTimestamp
							if kv == 0 {
								t := rfc3339ToUnixTimestamp(reflect.ValueOf(vv).String())
								row[columns[kv]] = t
							} else {
								row[columns[kv]] = vv
							}

						} else {
							row[columns[kv]] = ""
						}
					}
					returnValues = append(returnValues, row)
				}
			}

		}

		return returnValues, nil
	}

}

//조회한 결과List를 Map으로 변환한다.
func (r clientResponse) InfluxConverter(resp client.Response, serviceName string) (map[string]interface{}, model.ErrMessage) {

	if len(resp.Results) != 1 {
		return nil, nil
	} else {
		//UI로 Return할 결과값
		var returnValues []map[string]interface{}
		var columns []string

		for _, v := range resp.Results[0].Series {
			for _, vc := range v.Columns {
				columns = append(columns, vc)
			}

			for i := 0; i < len(v.Values); i++ {
				row := make(map[string]interface{})

				//InfluxDB에서 Value 값이 nil인 경우 해당 row는 값을 보내주지 않는다.
				if v.Values[i][1] != nil {
					for kv, vv := range v.Values[i] {
						if vv != nil {
							//Time Column Case convert to UnixTimestamp
							if kv == 0 {
								t := rfc3339ToUnixTimestamp(reflect.ValueOf(vv).String())
								row[columns[kv]] = t
							} else {
								row[columns[kv]] = vv
							}

						} else {
							row[columns[kv]] = ""
						}
					}
					returnValues = append(returnValues, row)
				}
			}

		}

		result := map[string]interface{}{
			SERVICE_NAME: serviceName,
			DATA_NAME:    returnValues,
		}

		return result, nil
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

// change unit
func ConByteToGB(metric string) string {
	metricByte, _ := strconv.ParseFloat(metric, 64)

	size := metricByte / 1024 / 1024 / 1024

	Data := fmt.Sprintf("%.2f", size)

	return Data
}

func ConByteToMB(metric string) string {
	metricByte, _ := strconv.ParseFloat(metric, 64)

	size := metricByte / 1024 / 1024

	return strconv.FormatFloat(size, 'f', 1, 64)
}

func ConByteToTB(metric string) string {
	metricByte, _ := strconv.ParseFloat(metric, 64)

	size := metricByte / 1024 / 1024 / 1024 / 1024

	Data := fmt.Sprintf("%.2f", size)

	return Data
}

func GetUnixTimeFromTo(interval int64) (string, string) {
	currentTime := time.Now().Unix()
	previousTime := currentTime - interval

	pTime := strconv.FormatInt(previousTime, 10)
	cTime := strconv.FormatInt(currentTime, 10)

	return pTime, cTime
}

func GetPromqlFromToParameter(interval int64, timeStep string) string {
	fromTime, toTime := GetUnixTimeFromTo(interval)

	parameter := "&start=" + fromTime + "&end=" + toTime + "&step=" + timeStep
	return parameter
}

func GetDBCurrentTime() time.Time {
	now := time.Now()
	t := now.Format(model.DB_DATE_FORMAT)
	currentTime, _ := time.Parse(time.RFC3339, t)
	return currentTime
}

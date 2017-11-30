package util

import (
	"kr/paasta/monitoring/domain"
	client "github.com/influxdata/influxdb/client/v2"
	"net/http"
	"fmt"
	"encoding/json"
	"log"
	"time"
	"reflect"
	"os"
	"bufio"
	"strings"
	"io"
)

const(
	SERVICE_NAME  string  = "serviceName"
	DATA_NAME     string = "data"
)

type Config map[string]string

func GetConnectionString(host, port, user, pass , dbname string) string {

	return fmt.Sprintf("%s:%s@%s([%s]:%s)/%s%s",
		user, pass, "tcp", host, port, dbname, "")

}

func RenderJsonResponse(data interface{}, w http.ResponseWriter) {

	js, err := json.Marshal(data)
	if err != nil {
		log.Fatalln("Error writing JSON:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(js)
	return
}

type errorMessage struct{
	domain.ErrMessage
}

func GetError() *errorMessage{
	return &errorMessage{}
}

func (e errorMessage) DbCheckError(err error) (domain.ErrMessage) {
	if err != nil{
		errMessage := domain.ErrMessage{
			"Message": err.Error() ,
		}
		return  errMessage
	}else{
		return nil
	}
}

//오류 체크 모듈로 오류 발생시 오류 메시지 리턴
func (e errorMessage) CheckError(resp client.Response, err error) (client.Response, domain.ErrMessage) {


	if err != nil {
		errMessage := domain.ErrMessage{
			"Message": err.Error() ,
		}
		return resp , errMessage

	}else if resp.Error() != nil {
		errMessage := domain.ErrMessage{
			"Message": resp.Err ,
		}
		return resp, errMessage
	}else {

		return resp, nil
	}
}


func GetDBCurrentTime() time.Time{
	now := time.Now()
	t := now.Format(domain.DB_DATE_FORMAT)
	currentTime, _ := time.Parse(time.RFC3339,t)
	return currentTime
}

type clientResponse struct{
	client.Response
}

func GetResponseConverter() *clientResponse{
	return &clientResponse{}
}

func rfc3339ToUnixTimestamp(metricDataTime string) int64{
	t, _ := time.Parse(time.RFC3339, metricDataTime)
	return t.Unix()
}


//조회한 결과List를 Map으로 변환한다.
func (r clientResponse) InfluxConverterToMap(resp client.Response) ([]map[string]interface{}, domain.ErrMessage){

	if len(resp.Results) != 1 {
		return nil , nil
	}else{
		//UI로 Return할 결과값
		var returnValues []map[string]interface{}
		var columns []string

		for _, v := range resp.Results[0].Series{
			for _, vc :=range v.Columns{
				columns = append(columns, vc)
			}

			for i :=0 ; i < len(v.Values); i++ {
				row := make(map[string]interface{})

				//InfluxDB에서 Value 값이 nil인 경우 해당 row는 값을 보내주지 않는다.
				if v.Values[i][1] != nil{
					for kv, vv := range v.Values[i]{
						if vv != nil{
							//Time Column Case convert to UnixTimestamp
							if kv == 0 {
								t := rfc3339ToUnixTimestamp(reflect.ValueOf(vv).String())
								row[columns[kv]] = t
							}else{
								row[columns[kv]] = vv
							}

						}else{
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
func (r clientResponse) InfluxConverter(resp client.Response, serviceName string) (map[string]interface{}, domain.ErrMessage){

	if len(resp.Results) != 1 {
		return nil , nil
	}else{
		//UI로 Return할 결과값
		var returnValues []map[string]interface{}
		var columns []string

		for _, v := range resp.Results[0].Series{
			for _, vc :=range v.Columns{
				columns = append(columns, vc)
			}

			for i :=0 ; i < len(v.Values); i++ {
				row := make(map[string]interface{})

				//InfluxDB에서 Value 값이 nil인 경우 해당 row는 값을 보내주지 않는다.
				if v.Values[i][1] != nil{
					for kv, vv := range v.Values[i]{
						if vv != nil{
							//Time Column Case convert to UnixTimestamp
							if kv == 0 {
								t := rfc3339ToUnixTimestamp(reflect.ValueOf(vv).String())
								row[columns[kv]] = t
							}else{
								row[columns[kv]] = vv
							}

						}else{
							row[columns[kv]] = ""
						}
					}
					returnValues = append(returnValues, row)
				}
			}

		}



		result := map[string]interface{}{
			SERVICE_NAME : serviceName,
			DATA_NAME: returnValues,
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
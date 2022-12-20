package utils

import (
	"github.com/influxdata/influxdb1-client/v2"
	"paasta-websocket/model"
	"reflect"
	"time"
)

const (
	SERVICE_NAME string = "serviceName"
	DATA_NAME    string = "data"
)

//오류 체크 모듈로 오류 발생시 오류 메시지 리턴
type errorMessage struct {
	model.ErrMessage
}

func GetError() *errorMessage {
	return &errorMessage{}
}

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

func rfc3339ToUnixTimestamp(metricDataTime string) int64 {
	t, _ := time.Parse(time.RFC3339, metricDataTime)
	return t.Unix()
}
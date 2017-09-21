package util

import (
	cb "kr/paasta/monitoring/monit-batch/models/base"
	client "github.com/influxdata/influxdb/client/v2"
	"reflect"
	"strconv"
	"time"
	"encoding/json"
)

const(
	SERVICE_NAME  string  = "serviceName"
	DATA_NAME     string = "data"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

type errorMessage struct{
	cb.ErrMessage
}

func GetError() *errorMessage{
	return &errorMessage{}
}

func (e errorMessage) DbCheckError(err error) (cb.ErrMessage) {
	if err != nil{
		errMessage := cb.ErrMessage{
			"Message": err.Error() ,
		}
		return  errMessage
	}else{
		return nil
	}
}

//오류 체크 모듈로 오류 발생시 오류 메시지 리턴
func (e errorMessage) CheckError(resp client.Response, err error) (client.Response, cb.ErrMessage) {


	if err != nil {
		errMessage := cb.ErrMessage{
			"Message": err.Error() ,
		}
		return resp , errMessage

	}else if resp.Error() != nil {
		errMessage := cb.ErrMessage{
			"Message": resp.Err ,
		}
		return resp, errMessage
	}else {

		return resp, nil
	}
}


type clientResponse struct{
	client.Response
}

func rfc3339ToUnixTimestamp(metricDataTime string) int64{
	t, _ := time.Parse(time.RFC3339, metricDataTime)
	return t.Unix()
}

func GetResponseConverter() *clientResponse{
	return &clientResponse{}
}


//CPU Count를 반환한다.
func (r clientResponse) GetFileSystems(resp client.Response) ([]string, cb.ErrMessage){

	var columns []string
	var mountPointList []string

	if len(resp.Results) != 1 {
		return nil , nil
	}else{
		for _, v := range resp.Results[0].Series{
			for _, vc :=range v.Columns{
				columns = append(columns, vc)
			}

			//mountPointList = make([]string, len(v.Values))


			for i :=0 ; i < len(v.Values); i++ {
				row := make(map[string]interface{})

				if v.Values[i][1] != nil{

					for kv, vv := range v.Values[i]{

						if vv != nil{

							if kv == 1 {
								row[columns[kv]] = vv
							}
						}
					}
					//FileSystem이 중복으로 들어 올때 중복 제거
					if StringInSlice(reflect.ValueOf(row["mountpoint"]).String(), mountPointList) == false{
						mountPointList = append(mountPointList, reflect.ValueOf(row["mountpoint"]).String())
					}

				}
			}
		}

		return mountPointList, nil
	}

}


//CPU Count를 반환한다.
func (r clientResponse) GetCpuCores(resp client.Response, serviceName string) (map[string]interface{}, cb.ErrMessage){

	if len(resp.Results) != 1 {
		return nil , nil
	}else{
		var columns []string
		cpuCnt := 0
		for _, v := range resp.Results[0].Series{
			for _, vc :=range v.Columns{
				columns = append(columns, vc)
			}

			for i :=0 ; i < len(v.Values); i++ {
				row := make(map[string]interface{})

				if v.Values[i][1] != nil{

					for kv, vv := range v.Values[i]{

						if vv != nil{

							if kv == 1 {
								row[columns[kv]] = vv
							}

						}
					}
					//Cpu 정보는 cpu0, cpu1, cpu2 .... 순으로 넘어온다.
					//cpu정보를 Substring 후 number 값만 Integer로 변환하여 cpu 최대 숫자를 구하여
					//CPU 갯수를 구한다.
					cpuNm := reflect.ValueOf(row["cpu"]).String()

					cpuNmSize := len(cpuNm)
					cpuCoreNum := cpuNm[3:cpuNmSize]
					cnt , _ := strconv.Atoi(cpuCoreNum)

					if cpuCnt <= cnt{
						cpuCnt = cnt + 1
					}

				}
			}
		}

		result := map[string]interface{}{
			SERVICE_NAME: serviceName,
			"cpu": cpuCnt,
		}

		return result, nil
	}

}

//조회한 결과List를 Map으로 변환한다.
func (r clientResponse) InfluxConverterToMap(resp client.Response) ([]map[string]interface{}, cb.ErrMessage){

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
func (r clientResponse) InfluxConverter(resp client.Response, serviceName string) (map[string]interface{}, cb.ErrMessage){

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


//조회한 결과List를 Map으로 변환한다.
func (r clientResponse) InfluxConverterList(resp client.Response) ([]map[string]interface{}, cb.ErrMessage){

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


func GetDataFromInterface(data []map[string]interface {}, idx int)(float64){

	var returnValue float64;

	datamap := data[idx]["data"].([]map[string]interface{})
	for _, data := range datamap{
		returnValue = data["usage"].(float64)
	}
	return returnValue
}

func GetDataFloatFromInterface(data []map[string]interface {}, idx int)(float64){

	var jsonValue json.Number;

	datamap := data[idx]["data"].([]map[string]interface{})
	for _, data := range datamap{
		jsonValue = data["usage"].(json.Number)
	}
	returnValue, _ := strconv.ParseFloat(jsonValue.String(),64)

	return returnValue
}

func GetDataFloatFromInterfaceSingle(data map[string]interface {})(float64){

	var jsonValue json.Number;

	datamap := data["data"].([]map[string]interface{})
	for _, data := range datamap{
		jsonValue = data["usage"].(json.Number)
	}
	returnValue, _ := strconv.ParseFloat(jsonValue.String(),64)

	return returnValue
}

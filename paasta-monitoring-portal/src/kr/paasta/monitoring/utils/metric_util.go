package utils

import (
	"encoding/json"
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"kr/paasta/monitoring/iaas/model"
	"reflect"
	"strconv"
	"time"
)

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

func convertString2Time(metricDataTime string) time.Time {
	t, _ := time.Parse(time.RFC3339, metricDataTime)
	return t
}

func GetDataFloatFromInterfaceSingle(data map[string]interface{}) float64 {

	var jsonValue json.Number

	fmt.Printf("model.RESULT_DATA_NAME : %v\n", model.RESULT_DATA_NAME)
	fmt.Printf("data : %v\n", data[model.RESULT_DATA_NAME])

	// 임시 오류 처리
	if data[model.RESULT_DATA_NAME] == nil {
		return 0
	}

	datamap := data[model.RESULT_DATA_NAME].([]map[string]interface{})
	for _, data := range datamap {
		jsonValue = data["value"].(json.Number)
	}
	returnValue, _ := strconv.ParseFloat(jsonValue.String(), 64)

	return returnValue
}

//조회한 결과List를 Map으로 변환한다.
func (r clientResponse) InfluxConverterList(resp client.Response, name string) (map[string]interface{}, model.ErrMessage) {

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

								/*datamap := vv.(json.Number)
								returnValue, _ := strconv.ParseFloat(datamap.String(),64)
								row[columns[kv]] = 100 - returnValue*/
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
			model.RESULT_NAME:      name,
			model.RESULT_DATA_NAME: returnValues,
		}
		return result, nil
	}

}

//조회한 결과List를 Map으로 변환한다.
func (r clientResponse) InfluxConverter(resp client.Response) (map[string]interface{}, model.ErrMessage) {

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
			model.RESULT_DATA_NAME: returnValues,
		}

		return result, nil
	}

}

//조회한 결과를 Map으로 변환한다.
func (r clientResponse) InfluxConverter4Usage(resp client.Response, name string) (map[string]interface{}, model.ErrMessage) {

	fmt.Println(resp)
	if len(resp.Results) != 2 {
		return nil, nil
	} else {
		//UI로 Return할 결과값
		//var returnValues      map[string]interface{}
		//MetricDB에서 받은 결과 값
		var resultValues []map[string]interface{}
		var returnValuesTotal []map[string]interface{}

		var columns []string

		for _, v := range resp.Results[0].Series {
			for _, vc := range v.Columns {
				columns = append(columns, vc)
			}

			for i := 0; i < len(v.Values); i++ {
				row := make(map[string]interface{})
				for kv, vv := range v.Values[i] {
					if vv != nil {
						row[columns[kv]] = vv
					} else {
						row[columns[kv]] = ""
					}
				}
				returnValuesTotal = append(returnValuesTotal, row)
			}
		}

		//revel.TRACE.Printf("returnValues1 ===>%s" , returnValues1)
		for _, v := range resp.Results[1].Series {
			for _, vc := range v.Columns {
				columns = append(columns, vc)
			}

			//만약 return된 두개의 결과 Data건수가 다를경우 작은 Data를 기준으로 건수계산
			resultDataCnt := 0

			if len(v.Values) != len(returnValuesTotal) {
				if len(v.Values) > len(returnValuesTotal) {
					resultDataCnt = len(returnValuesTotal)
				} else if len(v.Values) < len(returnValuesTotal) {
					resultDataCnt = len(v.Values)
				}

			} else {
				resultDataCnt = len(returnValuesTotal)
			}

			for i := 0; i < resultDataCnt; i++ {
				row := make(map[string]interface{})

				for kv, vv := range v.Values[i] {

					if kv == 0 {
						//동일한 일시 DateTime인지 Check한다
						//사용률이 null 이거나 "" 이면 백분률 계산에서 제외한다.
						//시간을 10초단위로 동일 Data 체크
						for _, totalData := range returnValuesTotal {
							time1 := vv.(string)
							time2 := totalData["time"].(string)
							if vv != nil && v.Values[i][1] != nil && totalData["usage"] != "" && time1[0:18] == time2[0:18] {
								isNegative := false

								if kv == 0 {

									//return된 Type이 Interface{}이므로 String으로 변환 후 Integer로 변환한다.
									total := reflect.ValueOf(totalData["usage"]).String()
									idle := reflect.ValueOf(v.Values[i][1]).String()
									totalUsage, _ := strconv.ParseFloat(total, 64)
									idleUsage, _ := strconv.ParseFloat(idle, 64)

									t := rfc3339ToUnixTimestamp(reflect.ValueOf(vv).String())

									//사용률 계산한다.
									result := idleUsage / totalUsage * 100

									//DiskSize인 경우 간헐적으로 비정상적인 Data가 들어온다.
									//Ex) totalUsage > idleUsage
									//이런 비정상적인 Data는 Skip한다.
									if idleUsage > totalUsage {
										isNegative = true
									} else {
										row[columns[0]] = t
										row[columns[1]] = 100 - result
									}

								}
								if isNegative == false {
									resultValues = append(resultValues, row)
									isNegative = false
								} else {
									isNegative = false
								}

							}
						}

					}
				}
			}
		}

		result := map[string]interface{}{
			model.RESULT_NAME:      name,
			model.RESULT_DATA_NAME: resultValues,
		}
		//returnValues = append(returnValues, result)
		return result, nil
	}

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
								row[columns[kv]] = t //reflect.ValueOf(vv).String()
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

//Mount Point 목록을 반환한다.
func (r clientResponse) GetMountPointList(resp client.Response) ([]string, model.ErrMessage) {

	var columns []string
	var mountPointList []string

	if len(resp.Results) != 1 {
		return nil, nil
	} else {
		for _, v := range resp.Results[0].Series {
			for _, vc := range v.Columns {
				columns = append(columns, vc)
			}

			for i := 0; i < len(v.Values); i++ {
				row := make(map[string]interface{})

				if v.Values[i][1] != nil {

					for kv, vv := range v.Values[i] {

						if vv != nil {

							if kv == 1 {
								row[columns[kv]] = vv
							}
						}
					}
					//FileSystem이 중복으로 들어 올때 중복 제거
					if StringArrayDistinct(reflect.ValueOf(row["mount_point"]).String(), mountPointList) == false {
						mountPointList = append(mountPointList, reflect.ValueOf(row["mount_point"]).String())
					}

				}
			}
		}

		return mountPointList, nil
	}
}

//조회한 결과List를 Map으로 변환한다.
func (r clientResponse) InfluxConverter4TopProcess(resp client.Response) (map[string]interface{}, model.ErrMessage) {

	if len(resp.Results) != 1 {
		return nil, nil
	} else {
		//UI로 Return할 결과값
		var returnValues []map[string]interface{}
		var columns []string

		for _, v := range resp.Results[0].Series {

			for _, vc := range v.Columns {
				if vc == "time" {
					columns = append(columns, "process_name")
				} else {
					columns = append(columns, vc)
				}

			}

			for i := 0; i < len(v.Values); i++ {
				row := make(map[string]interface{})

				//InfluxDB에서 Value 값이 nil인 경우 해당 row는 값을 보내주지 않는다.
				if v.Values[i][1] != nil {
					for kv, vv := range v.Values[i] {

						if vv != nil {
							//Time Column Case convert to UnixTimestamp
							if kv == 0 {
								//t := rfc3339ToUnixTimestamp(reflect.ValueOf(vv).String())
								row[columns[kv]] = v.Tags["process_name"]
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
			model.RESULT_DATA_NAME: returnValues,
		}

		return result, nil
	}

}

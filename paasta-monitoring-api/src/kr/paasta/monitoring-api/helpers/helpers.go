package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"math"
	models "paasta-monitoring-api/models/api/v1"
	"reflect"
	"strconv"
	"strings"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

//Int64ToString function convert a float number to a string
func Int64ToString(inputNum int64) string {
	return strconv.FormatInt(inputNum, 10)
}

func GetDBConnectionString(user, password, protocol, host, port, dbname, charset, parseTime string) string {
	return fmt.Sprintf("%s:%s@%s([%s]:%s)/%s?charset=%s&parseTime=%s",
		user, password, protocol, host, port, dbname, charset, parseTime)
}

// BindJsonAndCheckValid :: 클라이언트로부터 요청된 JSON 데이터를
// 매개인자 request로 들어온 구조체에 바인딩하고 바인딩된 구조체 데이터의 유효성을 검사한 결과를 반환한다.
func BindJsonAndCheckValid(c echo.Context, request interface{}) error {
	bindErr := c.Bind(&request)
	if bindErr != nil {
		return bindErr
	}

	validErr := CheckValid(&request)
	if validErr != nil {
		return validErr
	}

	return nil
}

// CheckValid :: 매개인자 request로 들어온 구조체 데이터의 유효성을 검사한 결과를 반환한다.
func CheckValid(request interface{}) error {
	v := validator.New()
	return v.Var(request, "dive")
}

func rfc3339ToUnixTimestamp(metricDataTime string) int64 {
	t, _ := time.Parse(time.RFC3339, metricDataTime)
	return t.Unix()
}

func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func RoundFloat(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num*output)) / output
}

func RoundFloatDigit2(num float64) float64 {
	return RoundFloat(num, 2)
}

func GetDataFloatFromInterfaceSingle(data map[string]interface{}) float64 {

	var jsonValue json.Number

	fmt.Printf("model.RESULT_DATA_NAME : %v\n", models.RESULT_DATA_NAME)
	fmt.Printf("data : %v\n", data[models.RESULT_DATA_NAME])

	// 임시 오류 처리
	if data[models.RESULT_DATA_NAME] == nil {
		return 0
	}

	datamap := data[models.RESULT_DATA_NAME].([]map[string]interface{})
	for _, data := range datamap {
		jsonValue = data["value"].(json.Number)
	}
	returnValue, _ := strconv.ParseFloat(jsonValue.String(), 64)

	return returnValue
}

//조회한 결과List를 Map으로 변환한다.
func InfluxConverterList(resp *client.Response, name string) (map[string]interface{}, error) {

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
			/*models.RESULT_NAME:      name,*/
			models.RESULT_DATA_NAME: returnValues,
		}
		return result, nil
	}

}

//조회한 결과List를 Map으로 변환한다.
func InfluxConverterToMap(resp *client.Response) ([]map[string]interface{}, error) {

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

//조회한 결과List를 Map으로 변환한다.
func InfluxConverter(resp *client.Response) (map[string]interface{}, error) {

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
			models.RESULT_DATA_NAME: returnValues,
		}
		return result, nil
	}
}

//조회한 결과를 Map으로 변환한다.
func InfluxConverter4Usage(resp *client.Response, name string) (map[string]interface{}, error) {

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
			/*models.RESULT_NAME:      name,*/
			models.RESULT_DATA_NAME: resultValues,
		}
		//returnValues = append(returnValues, result)
		return result, nil
	}

}

func TypeChecker_float64(target interface{}) interface{} {

	switch target.(type) {
	case int:
		// v is an int here, so e.g. v + 1 is possible.
		return float64(target.(int))
	case float64:
		// v is a float64 here, so e.g. v + 1.0 is possible.
		return target.(float64)
	case string:
		// v is a string here, so e.g. v + " Yeah!" is possible.
		f, _ := strconv.ParseFloat(target.(string), 64)
		return f
	case nil:
		// v is a string here, so e.g. v + " Yeah!" is possible.
		return float64(0)
	case json.Number:
		jsonValue := target.(json.Number)
		f, _ := strconv.ParseFloat(jsonValue.String(), 64)
		return f

	default:
		// And here I'm feeling dumb. ;)
		return float64(0)
	}
}

func TypeChecker_string(target interface{}) interface{} {
	switch target.(type) {
	case int:
		// v is an int here, so e.g. v + 1 is possible.
		return fmt.Sprintf("%d", target)
	case float64:
		// v is a float64 here, so e.g. v + 1.0 is possible.
		return fmt.Sprintf("%f", target)
	case string:
		// v is a string here, so e.g. v + " Yeah!" is possible.
		return target.(string)
	case nil:
		// v is a string here, so e.g. v + " Yeah!" is possible.
		return ""
	default:
		// And here I'm feeling dumb. ;)
		return ""
	}
}

func GetMaxAlarmLevel(alarmLevels []string) string {
	var status string

	for _, alarmLevel := range alarmLevels {
		if alarmLevel == models.ALARM_LEVEL_CRITICAL {
			status = alarmLevel
		} else if alarmLevel == models.ALARM_LEVEL_WARNING {
			if status != models.ALARM_LEVEL_CRITICAL {
				status = alarmLevel
			}
		}
	}

	return status
}

func GetAlarmStatusByServiceName(originType, alarmType string, usage float64, thresholds []models.AlarmPolicies) string {

	for _, threshold := range thresholds {
		if threshold.AlarmType == alarmType && threshold.OriginType == originType {
			if usage >= float64(threshold.CriticalThreshold) {
				return models.ALARM_LEVEL_CRITICAL
			} else if usage >= float64(threshold.WarningThreshold) {
				return models.ALARM_LEVEL_WARNING
			}
		}
	}

	return ""
}


func FindStructFieldWithBlankValues(object interface{}) string {
	var result []string
	elem := reflect.ValueOf(object).Elem()
	fieldCount := elem.NumField()
	for i := 0; i < fieldCount; i++ {
		value := elem.Field(i).Interface()
		name := elem.Type().Field(i).Name
		if value == "" {
			result = append(result, name)
		}
	}
	return strings.Join(result, ",")
}
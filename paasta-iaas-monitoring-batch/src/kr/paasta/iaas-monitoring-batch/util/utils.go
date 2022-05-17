package util

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	cb "iaas-monitoring-batch/model/base"
	"time"
	"net/url"
	"net/http"
	"iaas-monitoring-batch/model"
	"strings"
	"net"
	"encoding/json"
)

func TypeChecker_float64(target interface{}) interface{}{

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
		f, _ := strconv.ParseFloat(jsonValue.String(),64)
		return f

	default:
		// And here I'm feeling dumb. ;)
		return float64(0)
	}
}

func GetInsertCurrentTime() time.Time{
	now := time.Now()
	t := now.Format(cb.INSERT_DATE_FORMAT)
	currentTime, _ := time.Parse(time.RFC3339,t)
	return currentTime
}

func TimeToGeneralFormat(time time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", time.Year(), time.Month(),
		time.Day(), time.Hour(), time.Minute(), time.Second())
}

func GetConnectionString(host, port, user, pass , dbname string) string {

	return fmt.Sprintf("%s:%s@%s([%s]:%s)/%s%s",
		user, pass, "tcp", host, port, dbname, "")

}

func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func RoundFloat(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num * output)) / output
}

func RoundFloatDigit2(num float64) float64 {
	return RoundFloat(num , 2)
}

func FloattostrDigit2(fv float64) string {
	return strconv.FormatFloat(RoundFloatDigit2(fv), 'f', 2, 64)
}


func Floattostr(fv float64) string {
	return strconv.FormatFloat(fv, 'f', 2, 64)
}

func Floattostrwithprec(fv float64, prec int) string {
	return strconv.FormatFloat(fv, 'f', prec, 64)
}

func isExistArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}


func PortalExistCHeck()  error{

	timeOut := time.Duration(3) * time.Second
	portalUrl := model.PortalUrl

	fmt.Println("portalUrl:::", portalUrl[7 : len(portalUrl)])
	_, err := net.DialTimeout("tcp", portalUrl[7 : len(portalUrl)] + ":80", timeOut)

	return err

}

func HttpRequest(apiPath, method string,  headers map[string][]string, data []byte, portalClient http.Client)(*http.Response, int,  cb.ErrMessage){

	var request *http.Request
	var err error

	portalUrl := model.PortalUrl
	u, err := url.ParseRequestURI(portalUrl+ apiPath)
	u.Path = apiPath
	urlStr := fmt.Sprintf("%v", u)

	if method == "GET"{
		request, err = http.NewRequest(method, urlStr, nil) // <-- URL-encoded payload
		if headers != nil{
			request.Header = headers
		}

	}else if method == "POST"{
		request, err = http.NewRequest(method,  portalUrl+ apiPath,  strings.NewReader(string(data)))
	}

	fmt.Println("request==================", request)
	request.Header.Set("Content-type", "application/json")

	resp, err := portalClient.Do(request)

	if err != nil {
		errMessage := cb.ErrMessage{
			"Message": err.Error() ,
		}
		return resp, http.StatusInternalServerError, errMessage
	}else{
		//defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK{

			errMessage := cb.ErrMessage{
				"Message": resp.Status ,
			}
			return resp, resp.StatusCode,  errMessage
		}
		return resp, http.StatusOK, nil;
	}
}

func RemoveDuplicates(elements []int64) []int64 {

	encountered := make(map[int64]bool)
	var result []int64

	for v := range elements {
		if encountered[elements[v]] == true {
		} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}

	return result
}

func ReadConfig(filename string) (map[string]string, error) {
	// init with some bogus data
	config := make(map[string]string, 0)
	config["server.port"] = "9999"

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
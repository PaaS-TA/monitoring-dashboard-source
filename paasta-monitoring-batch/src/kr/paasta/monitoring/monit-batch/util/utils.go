package util

import (
	"fmt"
	"math"
	"strconv"
	cb "kr/paasta/monitoring/monit-batch/models/base"
	"time"
	"net/url"
	"net/http"
	"kr/paasta/monitoring/monit-batch/models"
	"strings"
	"net"
)


func GetInsertCurrentTime() time.Time{
	now := time.Now()
	t := now.Format(cb.INSERT_DATE_FORMAT)
	currentTime, _ := time.Parse(time.RFC3339,t)
	return currentTime
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

func floattostrwithprec(fv float64, prec int) string {
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
	portalUrl := models.PortalUrl

	fmt.Println("portalUrl:::", portalUrl[7 : len(portalUrl)])
	_, err := net.DialTimeout("tcp", portalUrl[7 : len(portalUrl)] + ":80", timeOut)

	return err

}

func HttpRequest(apiPath, method string,  headers map[string][]string, data []byte, portalClient http.Client)(*http.Response, int,  cb.ErrMessage){

	var request *http.Request
	var err error

	portalUrl := models.PortalUrl
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

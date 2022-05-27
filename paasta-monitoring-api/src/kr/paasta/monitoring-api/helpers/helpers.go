package helpers

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"strconv"
)

//Int64ToString function convert a float number to a string
func Int64ToString(inputNum int64) string {
	return strconv.FormatInt(inputNum, 10)
}

func GetDBConnectionString(dbtype, user, password, protocol, host, port, dbname, charset, parseTime string) (string, string) {
	return dbtype, fmt.Sprintf("%s:%s@%s([%s]:%s)/%s?charset=%s&parseTime=%s",
		user, password, protocol, host, port, dbname, charset, parseTime)
}

// BindRequestAndCheckValid :: 클라이언트의 Requset Body를
// userRequest로 받은 구조체에 바인딩한 결과와 구조체 값의 유효성을 검사한 결과를 반환한다.
func BindRequestAndCheckValid(c echo.Context, request interface{}) error {
	bindErr := c.Bind(&request)
	if bindErr != nil {
		return bindErr
	}

	validator := validator.New()
	validErr := validator.Var(request, "dive")
	if validErr != nil {
		return validErr
	}

	return nil
}

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
// UserInfo 구조체에 바인딩한 결과와 구조체 값의 유효성을 검사한 결과를 반환한다.
func BindRequestAndCheckValid(c echo.Context, userRequest interface{}) error {
    /* Request Body를 userInfo와 바인드한 결과 저장 */
    bindErr := c.Bind(&userRequest) // 유효한 Request Body가 아니라면 error 반환 (ex: JSON 문법 에러)
    if bindErr != nil {
        return bindErr
    }

    /* userInfo의 유효성 검사를 담당하는 객체 생성 */
    validator := validator.New()
    validErr := validator.Struct(userRequest) // models 패키지에 정의된 태그에 따라 유효성 검사, 문제시 error 반환
    if validErr != nil {
        return validErr
    }

    return nil
}

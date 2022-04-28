package apiHelpers

import (
    "github.com/labstack/echo/v4"
    "time"
)

//ResponseData structure
type ResponseData struct {
    Data interface{} `json:"data"`
    Meta interface{} `json:"meta"`
}

type BasicResponseForm struct {
    StatusCode   int         `json:"statusCode"`
    Message      string      `json:"message"`
    ResponseTime time.Time   `json:"responseTime"`
    ResponseInfo interface{} `json:"responseInfo"`
}

// Internal Server Error Format
func InternalErrMessage(err error) {
    panic(err)
}

// Request & Response Error Format
func ExternalErrMessage(status int, message string) map[string]interface{} {
    return map[string]interface{}{"status": status, "message": message}
}

func SetBasicResponseForm(statusCode int, message string, responseInfo interface{}) BasicResponseForm {
    responseForm := BasicResponseForm{
        StatusCode:   statusCode,
        Message:      message,
        ResponseTime: time.Now().Local(),
        ResponseInfo: responseInfo,
    }
    return responseForm
}

//Respond returns JSON format message with basic response form
func Respond(c echo.Context, statusCode int, message string, responseInfo ...interface{}) {
    response := SetBasicResponseForm(statusCode, message, responseInfo)
    c.JSON(statusCode, response)
}

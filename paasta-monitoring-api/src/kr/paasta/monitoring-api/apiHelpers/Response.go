package apiHelpers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

//ResponseData structure
type ResponseData struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta"`
}

type BasicResponseForm struct {
	StatusText   string      `json:"statusText"`
	StatusCode   int         `json:"statusCode"`
	Message      string      `json:"message"`
	ResponseTime time.Time   `json:"responseTime"`
	ResponseInfo interface{} `json:"responseInfo"`
}

func SetBasicResponseForm(statusCode int, message string, responseInfo interface{}) BasicResponseForm {
	responseForm := BasicResponseForm{
		StatusText:   http.StatusText(statusCode),
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

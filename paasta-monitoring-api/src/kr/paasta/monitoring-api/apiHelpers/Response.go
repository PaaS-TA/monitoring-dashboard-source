package apiHelpers

import (
	"encoding/json"
	"net/http"
)

//ResponseData structure
type ResponseData struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta"`
}

// Internal Server Error Format
func InternalErrMessage(err error) {
	panic(err)
}

// Request & Response Error Format
func ExternalErrMessage(status int, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

//Respond returns basic response structure
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

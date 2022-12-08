package model

type SocketParam struct {
	ClientId string `json:"client_id"`
	Category string `json:"category"`
	Command string `json:"command"`
	ExtraParam map[string]interface{} `json:"extra_param"`
}
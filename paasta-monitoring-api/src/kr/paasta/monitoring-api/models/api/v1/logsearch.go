package v1

type (
	Logs struct {
		UUID       string `json:"uuid"`
		LogType    string `json:"logType"`
		Keyword    string `json:"keyword"`
		TargetDate string `json:"targetDate"`
		Period     string `json:"period"`
		StartTime  string `json:"startTime"`
		EndTime    string `json:"endTime"`
		Messages   interface{}
	}
)

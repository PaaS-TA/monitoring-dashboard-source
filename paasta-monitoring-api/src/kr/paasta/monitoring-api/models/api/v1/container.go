package v1

type (
	CellInfo struct {
		ZoneName string `json:"zoneName"`
		CellName string `json:"cellName"`
		Ip       string `json:"ip"`
		Id       uint   `json:"id"`
	}
)

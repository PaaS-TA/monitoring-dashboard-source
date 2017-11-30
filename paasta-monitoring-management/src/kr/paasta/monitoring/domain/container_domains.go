package domain

type (

	ZoneCellInfo struct {
		ZoneName string
		CellName string
		Ip       string
		Id       uint
	}

	CellTileView struct{
		CellName   string 		   	 `json:"cellName"`
		Ip         string 		   	 `json:"ip"`
		ContainerTileView   []ContainerTileView  `json:"containers"`
	}

	ContainerTileView struct {
		AppName  string                    `json:"appName"`
		AppIndex string                    `json:"appIndex"`
	}
)
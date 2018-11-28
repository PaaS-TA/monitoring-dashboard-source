package model


const (
	CON_MTR_CPU_USAGE = "cpu_usage_total"
	CON_MTR_LOAD_AVG = "load_average"
	CON_MTR_MEM_USAGE = "memory_usage"
	CON_MTR_DISK_USAGE = "disk_usage"
	CON_MTR_RX_BYTES = "rx_bytes"
	CON_MTR_TX_BYTES = "tx_bytes"
	CON_MTR_RX_DROPPED = "rx_dropped"
	CON_MTR_TX_DROPPED = "tx_dropped"
	CON_MTR_RX_ERRORS = "rx_errors"
	CON_MTR_TX_ERRORS = "tx_errors"

	CON_MTR_ID_PREFIX = "/garden/"
)


type (

	ZoneCellInfo struct {
		ZoneName string
		CellName string
		Ip       string
		Id       uint
	}

	CellTileView struct{
		ZoneName   string			 `json:"zoneName"`
		CellName   string 		   	 `json:"cellName"`
		Ip         string 		   	 `json:"ip"`
		ContainerTileView   []ContainerTileView  `json:"containers"`
	}

	ContainerTileView struct {
		CellName		string		`json:"-"`
		AppName  		string		`json:"appName"`
		AppIndex 		string		`json:"appIndex"`
		Ip         		string		`json:"-"`
		ContainerId		string		`json:"containerId"`
	}

	ContainerSummary struct {
		ZoneName   		string		`json:"zoneName"`
		CellName  		string      `json:"cellName"`
		CellCnt   		string      `json:"cellCnt"`
		ContainerCnt   	string      `json:"containerCnt"`
		AppCnt    		string    	`json:"appCnt"`
		FailCnt   		string		`json:"failedCnt"`
		CriticalCnt 	string		`json:"criticalCnt"`
		WarningCnt 		string		`json:"warningCnt"`
		RunningCnt  	string		`json:"runningCnt"`
	}

	ContainerSummaryPagingRes struct {
		TotalCount			int					`json:"totalCount"`
		PageItems			int					`json:"pageItems"`
		Overview			OverviewCntRes 		`json:"overview"`
		CotainerSummaryList []ContainerSummary	`json:"data"`
	}

	ContainerRelationshipRes struct {
		ContainerSummary
		Ip			string			`json:"cellIp"`
		AppInfoList	[]AppStatusInfo `json:"containers"`
	}

	AppStatusInfo struct {
		AppName   		string        	`json:"appName"`
		AppIndex  		string          `json:"appIndex"`
		Status			string			`json:"status"`
		ContainerId		string			`json:"containerId"`
	}

	CellOverviewRes struct {
		CellName			string			`json:"cellName"`
		CellId				string			`json:"cellId"`
		Ip					string			`json:"cellIp"`
		State				string			`json:"status"`
		Core				string			`json:"core"`
		CpuUsage			float64			`json:"cpuUsage"`
		CpuState			string			`json:"cpuErrStat"`
		TotalMemory			float64			`json:"totalMemory"`
		MemoryUsage			float64			`json:"memoryUsage"`
		MemoryState			string			`json:"memErrStat"`
		TotalDisk			float64			`json:"totalDisk"`
		TotalDiskUsage		string			`json:"-"`
		DiskStatus			string			`json:"diskStatus"`
		CellState			string			`json:"-"`
		TotalDiskState		string			`json:"diskErrStat"`
	}

	ContainerOverviewRes struct {
		ZoneName  string            `json:"zoneName"`
		CellName  string            `json:"cellName"`
		Ip		  string            `json:"-"`
		AppName   string            `json:"appName"`
		AppIndex  string            `json:"appIndex"`
		ContainerName   string      `json:"containerId"`
		CpuUsage    float64    		`json:"cpuUsage"`
		CpuState	string			`json:"cpuErrStat"`
		MemoryUsage float64			`json:"memoryUsage"`
		MemoryState	string			`json:"memErrStat"`
		DiskUsage   float64			`json:"diskUsage"`
		DiskState	string			`json:"diskErrStat"`
		Status      string			`json:"status"`
	}

	OverviewCntRes struct {
		Running		string        		`json:"running"`
		Failed		string        		`json:"failed"`
		Critical	string        		`json:"critical"`
		Warning		string              `json:"warning"`
		Total		string              `json:"total"`
	}

	ContainerReq struct {
		BaseModel
		PageIndex 		  int 				`json:"pageIndex"`
		PageItems 		  int 				`json:"pageItems"`
		AppName   		  string        	`json:"appName"`
		AppIndex  		  string            `json:"appIndex"`
		ContainerName     string            `json:"containerName"`
		Name              string            `json:"name"`
		CellIp            string        	`json:"cellIp"`
		Status	          string
		SqlQuery          string
		Time			  string
		Item			  []ContainerDetailReq
	}

	ContainerDetailReq struct {
		Name		string
		ResName		string
	}

)
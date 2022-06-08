package ap

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/jinzhu/gorm"
	dao "paasta-monitoring-api/dao/api/v1/ap"
	"paasta-monitoring-api/helpers"
	models "paasta-monitoring-api/models/api/v1"
	"sort"
	"strconv"
	"time"
)

type ApBoshService struct {
	DbInfo         *gorm.DB
	InfluxDBClient client.Client
	BoshInfoList   []models.Bosh
	Env            map[string]interface{}
}

func GetApBoshService(DbInfo *gorm.DB, InfluxDBClient client.Client, BoshInfoList []models.Bosh, Env map[string]interface{}) *ApBoshService {
	return &ApBoshService{
		DbInfo:         DbInfo,
		InfluxDBClient: InfluxDBClient,
		BoshInfoList:   BoshInfoList,
		Env:            Env,
	}
}

func (b *ApBoshService) GetBoshInfoList() ([]models.Bosh, error) {
	results := b.BoshInfoList
	return results, nil
}

func (b *ApBoshService) GetBoshOverview() ([]models.BoshOverview, error) {
	var results []models.BoshOverview

	return results, nil
}

func (b *ApBoshService) GetBoshSummary() ([]models.BoshSummary, error) {
	var results []models.BoshSummary
	return results, nil
}

func (b *ApBoshService) GetBoshProcessByMemory() ([]models.BoshProcess, error) {
	var results []models.BoshProcess

	for _, BoshInfo := range b.BoshInfoList {
		resp, err := dao.GetBoshDao(b.DbInfo, b.InfluxDBClient, b.BoshInfoList, b.Env).GetBoshProcessByMemory(BoshInfo.UUID)

		if err != nil {
			fmt.Println(err.Error())
			return results, err
		} else {
			valueList, _ := helpers.InfluxConverterToMap(resp)

			var resList []map[string]interface{}

			for z := 0; z < len(valueList); z++ {
				if len(resList) > 0 {
					chk := false
					for y := 0; y < len(resList); y++ {
						if resList[y][models.IFX_MTR_PROC_NAME] == valueList[z][models.IFX_MTR_PROC_NAME] && resList[y][models.IFX_MTR_PROC_PID] == valueList[z][models.IFX_MTR_PROC_PID] {
							chk = true
						}
					}
					if !chk {
						resList = append(resList, valueList[z])
					}
				} else {
					resList = append(resList, valueList[z])
				}
			}

			// mem sort
			sort.Slice(resList, func(i, j int) bool {
				return helpers.TypeChecker_float64(resList[j][models.IFX_MTR_MEM_USAGE]).(float64) < helpers.TypeChecker_float64(resList[i][models.IFX_MTR_MEM_USAGE]).(float64)
			})

			var idx int

			for _, vl := range resList {
				var BoshProcess models.BoshProcess

				BoshProcess.Index = strconv.Itoa(idx + 1)
				BoshProcess.Process = helpers.TypeChecker_string(vl[models.IFX_MTR_PROC_NAME]).(string)
				BoshProcess.Memory = helpers.TypeChecker_float64(vl[models.IFX_MTR_MEM_USAGE]).(float64) / models.MB
				BoshProcess.Pid = strconv.FormatFloat(helpers.TypeChecker_float64(vl[models.IFX_MTR_PROC_PID]).(float64), 'f', 0, 64)
				BoshProcess.Time = time.Unix(vl[models.IFX_MTR_TIME].(int64), 0).Format(time.RFC3339)[0:19]
				BoshProcess.UUID = BoshInfo.UUID
				results = append(results, BoshProcess)
				idx++
				if idx == 5 {
					break
				} //fixed 5row
			}
		}
	}

	return results, nil
}

func (b *ApBoshService) GetBoshChart() ([]models.BoshChart, error) {
	var results []models.BoshChart
	return results, nil
}

func (b *ApBoshService) GetBoshLog() ([]models.BoshLog, error) {
	var results []models.BoshLog
	return results, nil
}

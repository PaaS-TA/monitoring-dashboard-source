package caas

import (
	models "paasta-monitoring-api/models/api/v1"
)

type WorkloadService struct {
	CaasConfig models.CaasConfig
}

func GetWorkloadService(config models.CaasConfig) *WorkloadService{
	return &WorkloadService{
		CaasConfig: config,
	}
}




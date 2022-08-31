package cp

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	models "paasta-monitoring-api/models/api/v1"
	service "paasta-monitoring-api/services/api/v1/cp"
)

type ClusterController struct {
	CaaS models.CP
}

func GetClusterController(config models.CP) *ClusterController {
	return &ClusterController{
		CaaS: config,
	}
}

// GetClusterAverage
//  @tags         CP
//  @Summary      쿠버네티스 클러스터 리소스 사용 평균 가져오기
//  @Description  쿠버네티스 클러스터 리소스 사용 평균을 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        type  path      string  true  "Type 정보를 주입한다."  enums(pod, cpu, memory, disk)
//  @Success      200   {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/cp/cluster/average/{type} [get]
func (controller *ClusterController) GetClusterAverage(ctx echo.Context) error {
	results, err := service.GetClusterService(controller.CaaS).GetClusterAverage(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get cluster average.", results)
	}
	return nil
}

// GetWorkNodeList
//  @tags         CP
//  @Summary      쿠버네티스 클러스터 전체 노드 정보 가져오기
//  @Description  쿠버네티스 클러스터 전체 노드 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/cp/cluster/worknodes [get]
func (controller *ClusterController) GetWorkNodeList(ctx echo.Context) error {
	results, err := service.GetClusterService(controller.CaaS).GetWorkNodeList(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Hypervisor statistics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get work node list.", results)
	}
	return nil
}

// GetWorkNode
//  @tags         CP
//  @Summary      쿠버네티스 클러스터 노드 정보 가져오기
//  @Description  쿠버네티스 클러스터 노드 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        nodename  query     string  false  "파라미터에 대한 Pod 관련 메트릭 정보를 반환한다."                example(ip-10-0-0-20)
//  @Param        instance  query     string  false  "파라미터에 대햔 CPU, Memory, Disk 관련 메트릭 정보를 반환한다."  example(10.0.0.20:9100)
//  @Success      200       {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/cp/cluster/worknode [get]
func (controller *ClusterController) GetWorkNode(ctx echo.Context) error {
	results, err := service.GetClusterService(controller.CaaS).GetWorkNode(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Worker Node data.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get work node.", results)
	}
	return nil
}

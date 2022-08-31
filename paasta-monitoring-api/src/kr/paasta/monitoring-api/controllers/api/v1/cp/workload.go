package cp

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	models "paasta-monitoring-api/models/api/v1"
	service "paasta-monitoring-api/services/api/v1/cp"
)

type WorkloadController struct {
	CaaS models.CP
}

func GetWorkloadController(config models.CP) *WorkloadController {
	return &WorkloadController{
		CaaS: config,
	}
}

// GetWorkloadStatus
//  @tags         CP
//  @Summary      쿠버네티스 워크로드 상태 정보 가져오기
//  @Description  쿠버네티스 워크로드 상태 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200        {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/cp/workload/status [get]
func (controller *WorkloadController) GetWorkloadStatus(ctx echo.Context) error {
	results, err := service.GetWorkloadService(controller.CaaS).GetWorkloadStatus(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Workload Status.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Workload Status.", results)
	}
	return nil
}

// GetWorkloadList
//  @tags         CP
//  @Summary      쿠버네티스 워크로드 리스트 가져오기
//  @Description  쿠버네티스 워크로드 리스트를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200        {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/cp/workload/list [get]
func (controller *WorkloadController) GetWorkloadList(ctx echo.Context) error {
	results, err := service.GetWorkloadService(controller.CaaS).GetWorkloadList()
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Workload List.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Workload List.", results)
	}
	return nil
}

// GetWorkloadDetailMetrics
//  @tags         CP
//  @Summary      쿠버네티스 워크로드 상세 메트릭 정보 가져오기
//  @Description  쿠버네티스 워크로드 상세 메트릭 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        workload  query     string  true  "워크로드명을 주입한다."  enums(deployment, statefulset, daemonset)
//  @Success      200       {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/cp/workload/metrics [get]
func (controller *WorkloadController) GetWorkloadDetailMetrics(ctx echo.Context) error {
	results, err := service.GetWorkloadService(controller.CaaS).GetWorkloadDetailMetrics(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed tto get Workload Detail Metrics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Workload Detail Metrics.", results)
	}
	return nil
}

// GetWorkloadContainerList
//  @tags         CP
//  @Summary      쿠버네티스 워크로드 컨테이너 리스트 가져오기
//  @Description  쿠버네티스 워크로드 컨테이너 리스트를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        workload  query     string  false  "워크로드명을 주입한다."  enums(deployment, statefulset, daemonset)
//  @Success      200       {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/cp/workload/container/list [get]
func (controller *WorkloadController) GetWorkloadContainerList(ctx echo.Context) error {
	results, err := service.GetWorkloadService(controller.CaaS).GetWorkloadContainerList(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Workload Container List.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Workload Container List.", results)
	}
	return nil
}

// GetContainerMetrics
//  @tags         CP
//  @Summary      쿠버네티스 컨테이너 메트릭 정보 가져오기
//  @Description  쿠버네티스 컨테이너 메트릭 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        namespace  query     string  true  "네임스페이스명을 주입한다."  example(kube-system)
//  @Param        container  query     string  true  "컨테이너명을 주입한다."    example(kube-proxy)
//  @Param        pod        query     string  true  "파드명을 주입한다."      example(kube-proxy-v47r9)
//  @Success      200  {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/cp/workload/container/metrics [get]
func (controller *WorkloadController) GetContainerMetrics(ctx echo.Context) error {
	results, err := service.GetWorkloadService(controller.CaaS).GetContainerMetrics(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Workload Container Metrics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Workload Container Metrics.", results)
	}
	return nil
}

// GetContainerLog
//  @tags         CP
//  @Summary      쿠버네티스 컨테이너 로그 가져오기
//  @Description  쿠버네티스 컨테이너 로그를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        namespace  query     string  true  "네임스페이스명을 주입한다."  example(kube-system)
//  @Param        container  query     string  true  "컨테이너명을 주입한다."    example(kube-proxy)
//  @Param        pod        query     string  true  "파드명을 주입한다."      example(kube-proxy-v47r9)
//  @Success      200  {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/cp/workload/container/log [get]
func (controller *WorkloadController) GetContainerLog(ctx echo.Context) error {
	results, err := service.GetWorkloadService(controller.CaaS).GetContainerLog(ctx)
	if err != nil {
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Workload Container Log.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Workload Container Log.", results)
	}
	return nil
}

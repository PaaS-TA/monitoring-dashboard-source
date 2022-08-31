package cp

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"paasta-monitoring-api/apiHelpers"
	models "paasta-monitoring-api/models/api/v1"
	service "paasta-monitoring-api/services/api/v1/cp"
)

type PodController struct {
	CaaS models.CP
}

func GetPodController(config models.CP) *PodController {
	return &PodController{
		CaaS: config,
	}
}

// GetPodStatus
//  @tags         CP
//  @Summary      쿠버네티스 파드 상태 정보 가져오기
//  @Description  쿠버네티스 파드 상태 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/cp/pod/status [get]
func (controller *PodController) GetPodStatus(ctx echo.Context) error {
	results, err := service.GetPodService(controller.CaaS).GetPodStatus(ctx)
	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Kubernetes Pod Status.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Kubernetes Pod Status.", results)
	}
	return nil
}

// GetPodList
//  @tags         CP
//  @Summary      쿠버네티스 파드 리스트 가져오기
//  @Description  쿠버네티스 파드 리스트를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Success      200  {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/cp/pod/list [get]
func (controller *PodController) GetPodList(ctx echo.Context) error {
	results, err := service.GetPodService(controller.CaaS).GetPodList()
	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed Succeeded to get Kubernetes Pod List.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Kubernetes Pod List.", results)
	}
	return nil
}

// GetPodDetailMetrics
//  @tags         CP
//  @Summary      쿠버네티스 파드 상세 메트릭 정보 가져오기
//  @Description  쿠버네티스 파드 상세 메트릭 정보를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        pod  query     string  true  "파드명을 주입한다."  example(prometheus-kube-prometheus-stack-1648-prometheus-0)
//  @Success      200  {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/cp/pod/metrics [get]
func (controller *PodController) GetPodDetailMetrics(ctx echo.Context) error {
	pod := ctx.QueryParam("pod")
	results, err := service.GetPodService(controller.CaaS).GetPodDetailMetrics(pod)
	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Kubernetes Pod Detail Metrics.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Kubernetes Pod Detail Metrics.", results)
	}
	return nil
}

// GetPodContainerList
//  @tags         CP
//  @Summary      쿠버네티스 파드 컨테이너 리스트 가져오기
//  @Description  쿠버네티스 파드 컨테이너 리스트를 가져온다.
//  @Accept       json
//  @Produce      json
//  @Param        pod  query     string  true  "파드명을 주입한다."  example(prometheus-kube-prometheus-stack-1648-prometheus-0)
//  @Success      200  {object}  apiHelpers.BasicResponseForm
//  @Router       /api/v1/cp/pod/container/list [get]
func (controller *PodController) GetPodContainerList(ctx echo.Context) error {
	pod := ctx.QueryParam("pod")
	results, err := service.GetPodService(controller.CaaS).GetPodContainerList(pod)
	if err != nil {
		log.Println(err.Error())
		apiHelpers.Respond(ctx, http.StatusBadRequest, "Failed to get Kubernetes Pod Container List.", err.Error())
		return err
	} else {
		apiHelpers.Respond(ctx, http.StatusOK, "Succeeded to get Kubernetes Pod Container List.", results)
	}
	return nil
}

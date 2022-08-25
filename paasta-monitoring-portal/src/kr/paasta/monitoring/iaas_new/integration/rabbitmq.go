package integration

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"monitoring-portal/iaas_new/model"
	"monitoring-portal/utils"
	"net/http"
	"time"
)

type RabbitMq struct {
	OpenstackProvider model.OpenstackProvider
}

func GetRabbitMq(openstack_provider model.OpenstackProvider) *RabbitMq {
	return &RabbitMq{
		OpenstackProvider: openstack_provider,
	}
}

/**
Description : Get project Storage Max & Used information
*/
//func (r *RabbitMq) GetRabbitMQOverview() (result map[string]interface{}, err error) {
func (r *RabbitMq) GetRabbitMQOverview() (rabbitmq_overview model.RabbitMQGlobalResource, err error) {
	//var rabbitmq_overview models.RabbitMQGlobalCounts
	client := &http.Client{
		CheckRedirect: func(req *http.Request, _ []*http.Request) error {
			return errors.New("No redirects")
		},
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives:   true,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	overview_url := fmt.Sprintf("http://%s:%s@%s:%s/api/overview", r.OpenstackProvider.RabbitmqUser, r.OpenstackProvider.RabbitmqPass, model.RabbitMqIp, model.RabbitMqPort)

	req, err := http.NewRequest("GET", overview_url, nil)
	if err != nil {
		return rabbitmq_overview, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return rabbitmq_overview, err
	}

	var data interface{}
	rawdata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return rabbitmq_overview, err
	}
	json.Unmarshal(rawdata, &data)
	msg := data.(map[string]interface{})
	resources := msg["object_totals"].(map[string]interface{})
	if len(resources) > 0 {
		rabbitmq_overview.Connections = utils.TypeChecker_int(resources["connections"]).(int)
		rabbitmq_overview.Channels = utils.TypeChecker_int(resources["channels"]).(int)
		rabbitmq_overview.Queues = utils.TypeChecker_int(resources["queues"]).(int)
		rabbitmq_overview.Consumers = utils.TypeChecker_int(resources["consumers"]).(int)
		rabbitmq_overview.Exchanges = utils.TypeChecker_int(resources["exchanges"]).(int)
	}

	var rabbitmq_node_resources model.RabbitMQNodeResources
	overview_url = fmt.Sprintf("http://%s:%s@%s:%s/api/nodes/%s", r.OpenstackProvider.RabbitmqUser, r.OpenstackProvider.RabbitmqPass, model.RabbitMqIp, model.RabbitMqPort, r.OpenstackProvider.RabbitmqTargetNode)

	req, err = http.NewRequest("GET", overview_url, nil)
	if err != nil {
		fmt.Println("Err Rabbit:", err.Error())
		return rabbitmq_overview, err
	}
	resp, err = client.Do(req)
	if err != nil {
		return rabbitmq_overview, err
	}
	rawdata, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return rabbitmq_overview, err
	}

	var rabbitResource interface{}
	json.Unmarshal(rawdata, &rabbitResource)
	rabbitResourceMap := rabbitResource.(map[string]interface{})

	rabbitmq_node_resources.FileDescriptorLimit = utils.TypeChecker_float64(rabbitResourceMap["fd_total"]).(float64)
	rabbitmq_node_resources.FileDescriptorUsed = utils.TypeChecker_float64(rabbitResourceMap["fd_used"]).(float64)
	rabbitmq_node_resources.SocketLimit = utils.TypeChecker_float64(rabbitResourceMap["sockets_total"]).(float64)
	rabbitmq_node_resources.SocketUsed = utils.TypeChecker_float64(rabbitResourceMap["sockets_used"]).(float64)
	rabbitmq_node_resources.ErlangProcLimit = utils.TypeChecker_int(rabbitResourceMap["proc_total"]).(int)
	rabbitmq_node_resources.ErlangProcUsed = utils.TypeChecker_int(rabbitResourceMap["proc_used"]).(int)
	rabbitmq_node_resources.MemoryMbLimit = utils.TypeChecker_float64(rabbitResourceMap["mem_limit"]).(float64) / 1024 / 1024
	rabbitmq_node_resources.MemoryMbUsed = utils.TypeChecker_float64(rabbitResourceMap["mem_used"]).(float64) / 1024 / 1024
	rabbitmq_node_resources.DiskMbLimit = utils.TypeChecker_float64(rabbitResourceMap["disk_free_limit"]).(float64) / 1024 / 1024
	rabbitmq_node_resources.DiskMbFree = utils.TypeChecker_float64(rabbitResourceMap["disk_free"]).(float64) / 1024 / 1024

	rabbitmq_overview.NodeResources = rabbitmq_node_resources

	return rabbitmq_overview, nil
}

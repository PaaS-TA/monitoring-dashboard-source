package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gorilla/websocket"
	influxClient "github.com/influxdata/influxdb1-client/v2"
	"log"
	"net/http"
	"paasta-websocket/model"
	"paasta-websocket/service"
	"paasta-websocket/utils"
	"paasta-websocket/zabbix-client/lib/go-zabbix"
	"sync"
)


var upgrader = websocket.Upgrader {
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {   // Allow all origin
		return true
	},
}

type Client struct {
	WebsocketClint *websocket.Conn
	ZabbixSession *zabbix.Session
	OpenstackProvider *gophercloud.ProviderClient
	InfluxDBClient influxClient.Client
	MessageType int
	Message chan interface{}
	Command chan string
}


var clientQueue *model.SessionQueue

func main() {
	log.Println("[system] main() called..")
	clientQueue = model.NewQueue()
	client := &Client{WebsocketClint: nil, MessageType: 0, Message: make(chan interface{}, 1024)}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request){
		handler(client, w, r)
	})
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}


func (c *Client) initThirdPartySession() {
	log.Println("[system] initThirdPartySession() called..")

	config, _ := utils.ReadConfig("config.ini")

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: config["identity.endpoint"],
		Username:         config["default.username"],
		Password: config["default.password"],
		TenantID:    config["default.tenant_id"],
		DomainName: config["default.domain"],
	}

	//Provider is the top-level client that all of your OpenStack services
	providerClient, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		fmt.Println(err.Error())
	}
	log.Println("OpenStack TokenID : " + providerClient.TokenID)
	c.OpenstackProvider = providerClient

	zabbixHost := config["zabbix.host"]
	zabbixAdminId := config["zabbix.admin.id"]
	zabbixAdminPw := config["zabbix.admin.pw"]
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	cache := zabbix.NewSessionFileCache().SetFilePath("./zabbix_session")
	zabbixSession, zabbixErr := zabbix.CreateClient(zabbixHost).
		WithCache(cache).
		WithHTTPClient(client).
		WithCredentials(zabbixAdminId, zabbixAdminPw).Connect()
	if zabbixErr != nil {
		log.Println(zabbixErr)
	}
	c.ZabbixSession = zabbixSession

	paasUrl, _ := config["paas.metric.db.url"]
	paasuserName, _ := config["paas.metric.db.username"]
	paasPassword, _ := config["paas.metric.db.password"]

	influxDbClient, err := influxClient.NewHTTPClient(influxClient.HTTPConfig{
		Addr:               paasUrl,
		Username:           paasuserName,
		Password:           paasPassword,
		InsecureSkipVerify: true,
	})
	c.InfluxDBClient = influxDbClient

}


var wg sync.WaitGroup  // WaitGroup for goroutine
func handler(client *Client, w http.ResponseWriter, r *http.Request) {
	log.Println("[system] handler() called..")

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client.WebsocketClint = ws
	client.initThirdPartySession()

	wg.Add(2)  // Set sub goroutine count
	go client.readMessage(r)
	go client.writeMessage()
	wg.Wait()
}


func (c *Client) readMessage(r *http.Request) {
	log.Println("[system] readMessage() called..")
	defer c.WebsocketClint.Close()
	for {
		messageType, message, err := c.WebsocketClint.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		c.MessageType = messageType
		log.Printf("Received message : %v / %v\n", c.MessageType, string(message))

		param := model.SocketParam{}
		json.Unmarshal(message, &param)


		if c.MessageType == websocket.CloseMessage {
			clientQueue.Pop(param.ClientId)
			c.WebsocketClint.Close()
		} else {
			c.Command <- param.Command
			switch param.Command {
			case "initialize":
				log.Println("[system] readMessage() initialize command")
				clientId := uuid.New().String()
				clientQueue.Push(clientId, c)
				log.Printf("[system] initialize new client : %s\n", clientId)
				c.Message <- clientId
				clientQueue.PrintQueueList()

			case "cpuUsage":
				log.Println("[system] readMessage() cpuUsage command")
				if param.Category == "iaas" {
					instanceId := param.ExtraParam["instance_id"].(string)
					hypervisorName := param.ExtraParam["host"].(string)
					result, _ := service.GetZabbixService(c.ZabbixSession, c.OpenstackProvider).GetCpuUsage(instanceId, hypervisorName, r)

					resultMap := make(map[string]interface{})
					resultMap["label"] = "CPU"
					resultMap["data"] = result

					resultList := make([]interface{}, 1)
					resultList[0] = resultMap

					log.Printf("result : %v\n", resultList)
					c.Message <- resultList
				} else if param.Category == "paas" {
					var params model.MetricsRequest
					params.ServiceName = param.ExtraParam["guid"].(string)
					params.Index = param.ExtraParam["idx"].(string)
					params.DefaultTimeRange = param.ExtraParam["time_range"].(string)
					params.TimeRangeFrom = param.ExtraParam["from_time"].(string)
					params.TimeRangeTo = param.ExtraParam["to_time"].(string)
					params.GroupBy = param.ExtraParam["group_by"].(string)

					result, _ := service.GetMetricsService(c.InfluxDBClient).GetAppCpuUsage(params)
					c.Message <- result
				}

			case "memoryUsage":
				log.Println("[system] readMessage() cpuUsage command")
				if param.Category == "iaas" {
					instanceId := param.ExtraParam["instance_id"].(string)
					hypervisorName := param.ExtraParam["host"].(string)
					result, _ := service.GetZabbixService(c.ZabbixSession, c.OpenstackProvider).GetMemoryUsage(instanceId, hypervisorName, r)

					resultMap := make(map[string]interface{})
					resultMap["label"] = "CPU"
					resultMap["data"] = result

					resultList := make([]interface{}, 1)
					resultList[0] = resultMap

					log.Printf("result : %v\n", resultList)
					c.Message <- resultList
				} else if param.Category == "paas" {
					var params model.MetricsRequest
					params.ServiceName = param.ExtraParam["guid"].(string)
					params.Index = param.ExtraParam["idx"].(string)
					params.DefaultTimeRange = param.ExtraParam["time_range"].(string)
					params.TimeRangeFrom = param.ExtraParam["from_time"].(string)
					params.TimeRangeTo = param.ExtraParam["to_time"].(string)
					params.GroupBy = param.ExtraParam["group_by"].(string)

					result, _ := service.GetMetricsService(c.InfluxDBClient).GetAppMemoryUsage(params)
					c.Message <- result
				}

			case "diskUsage":
				if param.Category == "iaas" {
					instanceId := param.ExtraParam["instance_id"].(string)
					hypervisorName := param.ExtraParam["host"].(string)
					result, _ := service.GetZabbixService(c.ZabbixSession, c.OpenstackProvider).GetDiskUsage(instanceId, hypervisorName, r)

					resultMap := make(map[string]interface{})
					resultMap["label"] = "Disk"
					resultMap["data"] = result

					resultList := make([]interface{}, 1)
					resultList[0] = resultMap

					log.Printf("result : %v\n", resultList)
					c.Message <- resultList
				} else if param.Category == "paas" {
					var params model.MetricsRequest
					params.ServiceName = param.ExtraParam["guid"].(string)
					params.Index = param.ExtraParam["idx"].(string)
					params.DefaultTimeRange = param.ExtraParam["time_range"].(string)
					params.TimeRangeFrom = param.ExtraParam["from_time"].(string)
					params.TimeRangeTo = param.ExtraParam["to_time"].(string)
					params.GroupBy = param.ExtraParam["group_by"].(string)

					result, _ := service.GetMetricsService(c.InfluxDBClient).GetDiskUsage(params)
					c.Message <- result
				}
			}
		}
	}
	wg.Done()
}


func (c *Client) writeMessage() {
	log.Println("[system] writeMessage() called..")
	//defer c.WebsocketClint.Close()
	for {
		if c.MessageType > 0 {
			message := <- c.Message  // Get data from Channel
			log.Printf("Send message : %v / %v\n", c.MessageType, fmt.Sprintf("%v", message))

			var responseBytes []byte
			if <-c.Command == "initialize" {
				responseBytes = []byte(message.(string))
			} else {
				responseBytes, _ = json.Marshal(message)
			}

			if err := c.WebsocketClint.WriteMessage(c.MessageType, responseBytes); err != nil {
				log.Println(err)
				return
			}
		}
	}
	wg.Done()
}
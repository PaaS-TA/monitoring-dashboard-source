// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package influxdb

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	/*"k8s.io/klog/v2"*/
	"net/http"

	/*"errors"*/
	"flag"
	"fmt"
	/*"io/ioutil"*/
	"net"
	/*"net/http"*/

	/*"net/url"*/
	"os"
	/*"os/exec"*/
	"strings"
	"sync"
	"time"

	info "github.com/google/cadvisor/info/v1"
	"github.com/google/cadvisor/storage"
	/*"github.com/google/cadvisor/version"*/

	/*influxdb "github.com/influxdb/influxdb/client"*/
	influxdb "github.com/influxdata/influxdb/client/v2"
)

func init() {
	storage.RegisterStorageDriver("influxdb", new)
}

var argDbRetentionPolicy = flag.String("storage_driver_influxdb_retention_policy", "", "retention policy")

type influxdbStorage struct {
	client          influxdb.Client
	cellIp          string
	machineName     string
	database        string
	retentionPolicy string
	bufferDuration  time.Duration
	lastWrite       time.Time
	points          []*influxdb.Point
	lock            sync.Mutex
	readyToFlush    func() bool
}

//====================================================================================
// Container Metrics Metadata from REP (127.0.0.1:1800/v1/containers)
type ContainerMetricsMetadata struct {
	Limits            Limits       `json:"limits,omitempty"`
	UsageMetrics      UsageMetrics `json:"usage_metrics,omitempty"`
	Container_Id      string       `json:"container_id,omitempty"`
	Interface_Id      string       `json:"interface_id,omitempty"`
	Application_Id    string       `json:"application_id,omitempty"`
	Application_Index string       `json:"application_index,omitempty"`
	Application_Name  string       `json:"application_name,omitempty"`
	Application_Urls  []string     `json:"application_uris,omitempty"`
}

type Limits struct {
	Fds    int32 `json:"fds,omitempty"`
	Memory int32 `json:"mem,omitempty"`
	Disk   int32 `json:"disk,omitempty"`
}

type UsageMetrics struct {
	MemoryUsageInBytes uint64        `json:"memory_usage_in_bytes"`
	DiskUsageInBytes   uint64        `json:"disk_usage_in_bytes"`
	TimeSpentInCPU     time.Duration `json:"time_spent_in_cpu"`
}

//====================================================================================

// Series names
const (
	// Cumulative CPU usage
	serCpuUsageTotal  string = "cpu_usage_total"
	serCpuUsageSystem string = "cpu_usage_system"
	serCpuUsageUser   string = "cpu_usage_user"
	serCpuUsagePerCpu string = "cpu_usage_per_cpu"
	// Smoothed average of number of runnable threads x 1000.
	serLoadAverage string = "load_average"
	// Memory Usage
	serMemoryUsage string = "memory_usage"
	// Maximum memory usage recorded
	serMemoryMaxUsage string = "memory_max_usage"
	// //Number of bytes of page cache memory
	serMemoryCache string = "memory_cache"
	// Size of RSS
	serMemoryRss string = "memory_rss"
	// Container swap usage
	serMemorySwap string = "memory_swap"
	// Size of memory mapped files in bytes
	serMemoryMappedFile string = "memory_mapped_file"
	// Working set size
	serMemoryWorkingSet string = "memory_working_set"
	// Number of memory usage hits limits
	serMemoryFailcnt string = "memory_failcnt"
	// Cumulative count of memory allocation failures
	serMemoryFailure string = "memory_failure"
	// Cumulative count of bytes received.
	serRxBytes string = "rx_bytes"
	// Cumulative count of receive errors encountered.
	serRxErrors string = "rx_errors"
	// Cumulative count of bytes transmitted.
	serTxBytes string = "tx_bytes"
	// Cumulative count of transmit errors encountered.
	serTxErrors string = "tx_errors"
	// Filesystem limit.
	serFsLimit string = "fs_limit"
	// Filesystem usage.
	serFsUsage string = "fs_usage"
	// Hugetlb stat - current res_counter usage for hugetlb
	setHugetlbUsage = "hugetlb_usage"
	// Hugetlb stat - maximum usage ever recorded
	setHugetlbMaxUsage = "hugetlb_max_usage"
	// Hugetlb stat - number of times hugetlb usage allocation failure
	setHugetlbFailcnt = "hugetlb_failcnt"
	// Perf statistics
	serPerfStat = "perf_stat"
	// Referenced memory
	serReferencedMemory = "referenced_memory"
	// Resctrl - Total memory bandwidth
	serResctrlMemoryBandwidthTotal = "resctrl_memory_bandwidth_total"
	// Resctrl - Local memory bandwidth
	serResctrlMemoryBandwidthLocal = "resctrl_memory_bandwidth_local"
	// Resctrl - Last level cache usage
	serResctrlLLCOccupancy = "resctrl_llc_occupancy"

	serRxDropped string = "rx_dropped"
	serTxDropped string = "tx_dropped"
	// Filesystem device.
	//serFsDevice string = "fs_device"
	// Filesystem limit.
	//serFsLimit string = "fs_limit"
	// Filesystem usage.
	//serFsUsage string = "fs_usage"

	// Disk Usage
	serDiskUsage string = "disk_usage"

	// Container Measurement
	serContainerMeausement string = "container_metrics"
)

func new() (storage.StorageDriver, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	//=================================================================================
	// VM IP address
	var instanceIp string
	ifaces, err := net.Interfaces()
	if err != nil {
		instanceIp = ""
	}
	for _, iface := range ifaces {
		//fmt.Println("##### local network interface name :", iface.Name)
		if strings.HasPrefix(iface.Name, "en") || strings.HasPrefix(iface.Name, "eth") {
			addrs, _ := iface.Addrs()
			for _, addr := range addrs {
				//Check whether addr is  IP adress or Mac address.
				ip_array := strings.Split(addr.String(), ".")
				if len(ip_array) >= 4 {
					//var ip net.IP
					switch v := addr.(type) {
					case *net.IPNet:
						instanceIp = v.IP.String()
					}
				}
			}
		}
	}
	//fmt.Println("##### local network address 1 :", instanceIp)
	if instanceIp == "" {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			instanceIp = ""
		}
		for _, address := range addrs {
			// check the address type and if it is not a loopback the display it
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					instanceIp = ipnet.IP.String() //fmt.Println(ipnet.IP.String())
				}

			}
		}
	}
	fmt.Println("##### local network address 2 :", instanceIp)
	//=================================================================================
	/*
		var cellIp string
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			return nil, err
		}

		for _, address := range addrs {
			// check the address type and if it is not a loopback the display it
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {

				fmt.Println("address :", address.Network())
				fmt.Println("ipnet :", ipnet.IP)

				if ipnet.IP.To4() != nil {
					cellIp = ipnet.IP.String()
					fmt.Println("cellIp:", ipnet.IP.String())
				}

			}
		}*/

	return newStorage(
		hostname,
		instanceIp,
		*storage.ArgDbTable,
		*storage.ArgDbName,
		*argDbRetentionPolicy,
		*storage.ArgDbUsername,
		*storage.ArgDbPassword,
		*storage.ArgDbHost,
		*storage.ArgDbIsSecure,
		*storage.ArgDbBufferDuration,
	)
}

// machineName: A unique identifier to identify the host that current cAdvisor
// instance is running on.
// influxdbHost: The host which runs influxdb (host:port)
func newStorage(
	machineName,
	cellIp,
	tablename,
	database,
	retentionPolicy,
	username,
	password,
	influxdbHost string,
	isSecure bool,
	bufferDuration time.Duration,
) (*influxdbStorage, error) {
	/*url := &url.URL{
		Scheme: "http",
		Host:   influxdbHost,
	}
	if isSecure {
		url.Scheme = "https"
	}

	config := &influxdb.Config{
		URL:       *url,
		Username:  username,
		Password:  password,
		UserAgent: fmt.Sprintf("%v/%v", "cAdvisor", version.Info["version"]),
	}
	client, err := influxdb.NewClient(*config)
	if err != nil {
		return nil, err
	}

	ret := &influxdbStorage{
		client:          client,
		machineName:     machineName,
		cellIp:          cellIp,
		database:        database,
		retentionPolicy: retentionPolicy,
		bufferDuration:  bufferDuration,
		lastWrite:       time.Now(),
		points:          make([]*influxdb.Point, 0),
	}
	ret.readyToFlush = ret.defaultReadyToFlush
	return ret, nil*/
	// Make client
	client, err := influxdb.NewUDPClient(influxdb.UDPConfig{
		Addr: influxdbHost,
		//PayloadSize: 4096,
	})

	if err != nil {
		return nil, err
	}

	ret := &influxdbStorage{
		client:      client,
		machineName: machineName,
		cellIp:      cellIp,
		database:    database,
		lastWrite:   time.Now(),
		points:      make([]*influxdb.Point, 0),
	}
	ret.readyToFlush = ret.defaultReadyToFlush
	return ret, nil
}

// Field names
const (
	fieldAppDisk string = "app_disk"
	fieldAppMem  string = "app_mem"
	fieldValue   string = "value"
	fieldType    string = "type"
	fieldDevice  string = "device"
)

// Tag names
const (
	tagName               string = "name"
	tagMachineName        string = "machine"
	tagContainerName      string = "container_id"
	tagContainerInterface string = "container_interface"
	tagCellIp             string = "cell_ip"
	tagApplicationId      string = "application_id"
	tagApplicationIndex   string = "application_index"
	tagApplicationName    string = "application_name"
	tagApplicationUrl     string = "application_url"
)

// Container Metrics Metadata from REP (127.0.0.1:1800/v1/containers)
func (self *influxdbStorage) containerMetricsMedataData() []ContainerMetricsMetadata {
	cafile := "/var/vcap/jobs/cadvisor/config/certs/rep/client.crt"
	caCert, caCertErr := ioutil.ReadFile(cafile)
	if caCertErr != nil {
		fmt.Println("##### get Container Metrics Metadata caCertErr:", caCertErr)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	certfile := "/var/vcap/jobs/cadvisor/config/certs/rep/client.crt"
	keyfile := "/var/vcap/jobs/cadvisor/config/certs/rep/client.key"

	cert, ceatErr := tls.LoadX509KeyPair(certfile, keyfile)
	if ceatErr != nil {
		fmt.Println("##### Failed to use the collector certificate and key:", ceatErr)
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}
	tlsConfig.BuildNameToCertificate()

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := http.Client{Transport: transport}
	resp, respErr := client.Get("https://127.0.0.1:1800/v1/containers")
	if respErr != nil {
		fmt.Println("##### get Container Metrics Metadata respErr:", respErr)
		//glog.Error("##### get Container Metrics Metadata request err:", err)
	}
	defer resp.Body.Close()
	//fmt.Println("##### get Container Metrics Metadata resp:", resp)
	//fmt.Println("##### get Container Metrics Metadata resp.Body:", resp.Body)

	bytes, _ := ioutil.ReadAll(resp.Body)
	//str := string(bytes) // bytes -> string
	//fmt.Println(str)

	/*if err != nil {
		fmt.Println("##### get Container Metrics Metadata request err:", err)
		//glog.Error("##### get Container Metrics Metadata request err:", err)
	}*/
	if resp != nil {
		//rawdata, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println("##### Response Data :", string(rawdata))

		var containermetrics []ContainerMetricsMetadata
		//rawdata := []byte("[{\"limits\":{\"fds\":16384,\"mem\":1024,\"disk\":1024},\"usage_metrics\":{\"memory_usage_in_bytes\":230674432,\"disk_usage_in_bytes\":173891584,\"time_spent_in_cpu\":29331153300},\"container_id\":\"770e2059-b934-4dfe-7871-e4f9\",\"interface_id\":\"wggd5cu4rlph-0\",\"application_id\":\"21ed5b0e-2cda-4b0f-8e9d-4fa5fbb80088\",\"application_index\":\"0\",\"application_name\":\"spring-music-pinpoint-1\",\"application_uris\":[\"spring-music-pinpoint-1-unexpected-tasmaniandevil-sf.182.252.135.97.xip.io\"]}]")
		json.Unmarshal(bytes, &containermetrics)

		/*fmt.Println("##### Container Metrics Metadata :", containermetrics, len(containermetrics))
		for _, metrics :=range containermetrics{
			fmt.Println("##### Container Metrics container id :", metrics.Container_Id)
			fmt.Println("##### Container Metrics app id :", metrics.Application_Id)
			fmt.Println("##### Container Metrics app name :", metrics.Application_Name)
			fmt.Println("##### Container Metrics app urls :", metrics.Application_Urls)
			fmt.Println("##### Container Metrics app limits :", metrics.Limits)
			fmt.Println("##### Container Metrics app usage-memory :", metrics.UsageMetrics.MemoryUsageInBytes)
			fmt.Println("##### Container Metrics app usage-disk :", metrics.UsageMetrics.DiskUsageInBytes)
			fmt.Println("##### Container Metrics app usage-cpu(second) :", metrics.UsageMetrics.TimeSpentInCPU.Seconds())
		}*/
		fmt.Println("!@#$!%!@#%", containermetrics)
		return containermetrics
	}
	return nil
}

//====================================================================================

func (s *influxdbStorage) containerFilesystemStatsToPoints(
	cInfo *info.ContainerInfo,
	stats *info.ContainerStats) (points []*influxdb.Point) {
	if len(stats.Filesystem) == 0 {
		return points
	}
	for _, fsStat := range stats.Filesystem {
		tagsFsUsage := map[string]string{
			tagMachineName:        s.machineName,
			tagContainerName:      cInfo.Name,
			tagContainerInterface: cInfo.Namespace,
			fieldDevice:           fsStat.Device,
			fieldType:             "usage",
		}
		fieldsFsUsage := map[string]interface{}{
			fieldValue: float64(fsStat.Usage),
		}
		/*pointFsUsage := &influxdb.Point{
			Measurement: serContainerMeausement,
			Tags:        tagsFsUsage,
			Fields:      fieldsFsUsage,
		}*/
		pointFsUsage, err := influxdb.NewPoint(serContainerMeausement, tagsFsUsage, fieldsFsUsage)
		if err != nil {
			fmt.Println(err)
			/*glog.Fatalf("Failed to create NewPoint for FieldsFsUsage: %v", err)*/
		}

		tagsFsLimit := map[string]string{
			tagMachineName:        s.machineName,
			tagContainerName:      cInfo.Name,
			tagContainerInterface: cInfo.Namespace,
			fieldDevice:           fsStat.Device,
			fieldType:             "limit",
		}
		fieldsFsLimit := map[string]interface{}{
			fieldValue: float64(fsStat.Limit),
		}
		/*pointFsLimit := &influxdb.Point{
			Measurement: serContainerMeausement,
			Tags:        tagsFsLimit,
			Fields:      fieldsFsLimit,
		}*/
		pointFsLimit, err := influxdb.NewPoint(serContainerMeausement, tagsFsLimit, fieldsFsLimit)
		if err != nil {
			fmt.Println(err)
			/*glog.Fatalf("Failed to create NewPoint for FieldsFsLimit: %v", err)*/
		}

		points = append(points, pointFsUsage, pointFsLimit)
	}

	//s.tagPoints(cInfo, stats, points)

	return points
}

// Set tags and timestamp for all points of the batch.
// Points should inherit the tags that are set for BatchPoints, but that does not seem to work.
/*func (s *influxdbStorage) tagPoints(cInfo *info.ContainerInfo, stats *info.ContainerStats, points []*influxdb.Point) {
	// Use container alias if possible
	var containerName string
	if len(cInfo.ContainerReference.Aliases) > 0 {
		containerName = cInfo.ContainerReference.Aliases[0]
	} else {
		containerName = cInfo.ContainerReference.Name
	}

	commonTags := map[string]string{
		tagMachineName:   s.machineName,
		tagContainerName: containerName,
	}
	for i := 0; i < len(points); i++ {
		// merge with existing tags if any
		addTagsToPoint(points[i], commonTags)
		addTagsToPoint(points[i], cInfo.Spec.Labels)
		points[i].Time = stats.Timestamp
	}
}*/

func (s *influxdbStorage) containerStatsToPoints(
	cInfo *info.ContainerInfo,
	containerMetric ContainerMetricsMetadata,
	stats *info.ContainerStats,
) (points []*influxdb.Point) {
	/*
		// CPU usage: Total usage in nanoseconds
		points = append(points, makePoint(serCpuUsageTotal, stats.Cpu.Usage.Total))

		// CPU usage: Time spend in system space (in nanoseconds)
		points = append(points, makePoint(serCpuUsageSystem, stats.Cpu.Usage.System))

		// CPU usage: Time spent in user space (in nanoseconds)
		points = append(points, makePoint(serCpuUsageUser, stats.Cpu.Usage.User))

		// CPU usage per CPU
		for i := 0; i < len(stats.Cpu.Usage.PerCpu); i++ {
			point := makePoint(serCpuUsagePerCpu, stats.Cpu.Usage.PerCpu[i])
			tags := map[string]string{"instance": fmt.Sprintf("%v", i)}
			addTagsToPoint(point, tags)

			points = append(points, point)
		}

		// Load Average
		points = append(points, makePoint(serLoadAverage, stats.Cpu.LoadAverage))

		// Network Stats
		points = append(points, makePoint(serRxBytes, stats.Network.RxBytes))
		points = append(points, makePoint(serRxErrors, stats.Network.RxErrors))
		points = append(points, makePoint(serTxBytes, stats.Network.TxBytes))
		points = append(points, makePoint(serTxErrors, stats.Network.TxErrors))

		// Referenced Memory
		points = append(points, makePoint(serReferencedMemory, stats.ReferencedMemory))*/

	// CPU Usage
	points = append(points, makePoint(s.machineName, s.cellIp, cInfo.Name, serCpuUsageTotal, containerMetric, containerMetric.UsageMetrics.TimeSpentInCPU.Seconds()))
	// Load Average
	points = append(points, makePoint(s.machineName, s.cellIp, cInfo.Name, serLoadAverage, containerMetric, float64(stats.Cpu.LoadAverage)))
	// Memory Usage
	points = append(points, makePoint(s.machineName, s.cellIp, cInfo.Name, serMemoryUsage, containerMetric, float64(containerMetric.UsageMetrics.MemoryUsageInBytes)))
	// Disk Usage
	points = append(points, makePoint(s.machineName, s.cellIp, cInfo.Name, serDiskUsage, containerMetric, float64(containerMetric.UsageMetrics.DiskUsageInBytes)))

	// Network Stats
	for i := 0; i < len(stats.Network.Interfaces); i++ {
		/*fmt.Println("interface name :", stats.Network.Interfaces[i].Name)
		fmt.Println("rxbytes :", stats.Network.Interfaces[i].RxBytes)
		fmt.Println("rxerror :", stats.Network.Interfaces[i].RxErrors)
		fmt.Println("rxdropped :", stats.Network.Interfaces[i].RxDropped)
		fmt.Println("txbytes :", stats.Network.Interfaces[i].TxBytes)
		fmt.Println("txerror :", stats.Network.Interfaces[i].TxErrors)
		fmt.Println("txdropped :", stats.Network.Interfaces[i].TxDropped)*/
		points = append(points, makePoint(s.machineName, s.cellIp, stats.Network.Interfaces[i].Name, serRxBytes, containerMetric, float64(stats.Network.Interfaces[i].RxBytes)))
		points = append(points, makePoint(s.machineName, s.cellIp, stats.Network.Interfaces[i].Name, serRxErrors, containerMetric, float64(stats.Network.Interfaces[i].RxErrors)))
		points = append(points, makePoint(s.machineName, s.cellIp, stats.Network.Interfaces[i].Name, serRxDropped, containerMetric, float64(stats.Network.Interfaces[i].RxDropped)))
		points = append(points, makePoint(s.machineName, s.cellIp, stats.Network.Interfaces[i].Name, serTxBytes, containerMetric, float64(stats.Network.Interfaces[i].TxBytes)))
		points = append(points, makePoint(s.machineName, s.cellIp, stats.Network.Interfaces[i].Name, serTxErrors, containerMetric, float64(stats.Network.Interfaces[i].TxErrors)))
		points = append(points, makePoint(s.machineName, s.cellIp, stats.Network.Interfaces[i].Name, serTxDropped, containerMetric, float64(stats.Network.Interfaces[i].TxDropped)))
	}

	//s.tagPoints(cInfo, stats, points)

	return points
}

/*func (s *influxdbStorage) memoryStatsToPoints(
	cInfo *info.ContainerInfo,
	stats *info.ContainerStats,
) (points []*influxdb.Point) {
	// Memory Usage
	points = append(points, makePoint(serMemoryUsage, stats.Memory.Usage))
	// Maximum memory usage recorded
	points = append(points, makePoint(serMemoryMaxUsage, stats.Memory.MaxUsage))
	//Number of bytes of page cache memory
	points = append(points, makePoint(serMemoryCache, stats.Memory.Cache))
	// Size of RSS
	points = append(points, makePoint(serMemoryRss, stats.Memory.RSS))
	// Container swap usage
	points = append(points, makePoint(serMemorySwap, stats.Memory.Swap))
	// Size of memory mapped files in bytes
	points = append(points, makePoint(serMemoryMappedFile, stats.Memory.MappedFile))
	// Working Set Size
	points = append(points, makePoint(serMemoryWorkingSet, stats.Memory.WorkingSet))
	// Number of memory usage hits limits
	points = append(points, makePoint(serMemoryFailcnt, stats.Memory.Failcnt))

	// Cumulative count of memory allocation failures
	memoryFailuresTags := map[string]string{
		"failure_type": "pgfault",
		"scope":        "container",
	}
	memoryFailurePoint := makePoint(serMemoryFailure, stats.Memory.ContainerData.Pgfault)
	addTagsToPoint(memoryFailurePoint, memoryFailuresTags)
	points = append(points, memoryFailurePoint)

	memoryFailuresTags["failure_type"] = "pgmajfault"
	memoryFailurePoint = makePoint(serMemoryFailure, stats.Memory.ContainerData.Pgmajfault)
	addTagsToPoint(memoryFailurePoint, memoryFailuresTags)
	points = append(points, memoryFailurePoint)

	memoryFailuresTags["failure_type"] = "pgfault"
	memoryFailuresTags["scope"] = "hierarchical"
	memoryFailurePoint = makePoint(serMemoryFailure, stats.Memory.HierarchicalData.Pgfault)
	addTagsToPoint(memoryFailurePoint, memoryFailuresTags)
	points = append(points, memoryFailurePoint)

	memoryFailuresTags["failure_type"] = "pgmajfault"
	memoryFailurePoint = makePoint(serMemoryFailure, stats.Memory.HierarchicalData.Pgmajfault)
	addTagsToPoint(memoryFailurePoint, memoryFailuresTags)
	points = append(points, memoryFailurePoint)

	//s.tagPoints(cInfo, stats, points)

	return points
}*/

/*func (s *influxdbStorage) hugetlbStatsToPoints(
	cInfo *info.ContainerInfo,
	stats *info.ContainerStats,
) (points []*influxdb.Point) {

	for pageSize, hugetlbStat := range stats.Hugetlb {
		tags := map[string]string{
			"page_size": pageSize,
		}

		// Hugepage usage
		point := makePoint(setHugetlbUsage, hugetlbStat.Usage)
		addTagsToPoint(point, tags)
		points = append(points, point)

		//Maximum hugepage usage recorded
		point = makePoint(setHugetlbMaxUsage, hugetlbStat.MaxUsage)
		addTagsToPoint(point, tags)
		points = append(points, point)

		// Number of hugepage usage hits limits
		point = makePoint(setHugetlbFailcnt, hugetlbStat.Failcnt)
		addTagsToPoint(point, tags)
		points = append(points, point)
	}

	//s.tagPoints(cInfo, stats, points)

	return points
}*/

/*func (s *influxdbStorage) perfStatsToPoints(
	cInfo *info.ContainerInfo,
	stats *info.ContainerStats,
) (points []*influxdb.Point) {

	for _, perfStat := range stats.PerfStats {
		point := makePoint(serPerfStat, perfStat.Value)
		tags := map[string]string{
			"cpu":           fmt.Sprintf("%v", perfStat.Cpu),
			"name":          perfStat.Name,
			"scaling_ratio": fmt.Sprintf("%v", perfStat.ScalingRatio),
		}
		addTagsToPoint(point, tags)
		points = append(points, point)
	}

	//s.tagPoints(cInfo, stats, points)

	return points
}*/

/*func (s *influxdbStorage) resctrlStatsToPoints(
	cInfo *info.ContainerInfo,
	stats *info.ContainerStats,
) (points []*influxdb.Point) {

	// Memory bandwidth
	for nodeID, rdtMemoryBandwidth := range stats.Resctrl.MemoryBandwidth {
		tags := map[string]string{
			"node_id": fmt.Sprintf("%v", nodeID),
		}
		point := makePoint(serResctrlMemoryBandwidthTotal, rdtMemoryBandwidth.TotalBytes)
		addTagsToPoint(point, tags)
		points = append(points, point)

		point = makePoint(serResctrlMemoryBandwidthLocal, rdtMemoryBandwidth.LocalBytes)
		addTagsToPoint(point, tags)
		points = append(points, point)
	}

	// Cache
	for nodeID, rdtCache := range stats.Resctrl.Cache {
		tags := map[string]string{
			"node_id": fmt.Sprintf("%v", nodeID),
		}
		point := makePoint(serResctrlLLCOccupancy, rdtCache.LLCOccupancy)
		addTagsToPoint(point, tags)
		points = append(points, point)
	}

	//s.tagPoints(cInfo, stats, points)

	return points
}*/

func (s *influxdbStorage) OverrideReadyToFlush(readyToFlush func() bool) {
	s.readyToFlush = readyToFlush
}

func (s *influxdbStorage) defaultReadyToFlush() bool {
	return time.Since(s.lastWrite) >= s.bufferDuration
}

func (s *influxdbStorage) AddStats(cInfo *info.ContainerInfo, stats *info.ContainerStats) error {
	if stats == nil {
		return nil
	}
	var pointsToFlush []*influxdb.Point
	func() {
		// AddStats will be invoked simultaneously from multiple threads and only one of them will perform a write.
		s.lock.Lock()
		defer s.lock.Unlock()

		var containerName string
		var containerMetric ContainerMetricsMetadata
		if len(cInfo.Aliases) > 0 {
			containerName = cInfo.Aliases[0]
		} else {
			containerName = cInfo.Name
		}
		//fmt.Println("================ containerName :", containerName)

		// here, container id is seperation process, because need to containerMetricsMetadata function call control

		//===================================================================
		// Container Metrics Metadata from REP (127.0.0.1:1800/v1/containers)
		containerMetrics := s.containerMetricsMedataData()
		containerNames := strings.Split(containerName, "-")
		containerMetric.Container_Id = containerNames[len(containerNames)-1]
		for _, metrics := range containerMetrics {
			//fmt.Println("================ metrics:", metrics)
			//fmt.Println("================ metrics-contaier-id:", metrics.Container_Id)
			//fmt.Println("================ metrics-contaier-interface-id:", metrics.Interface_Id)

			//CAdvisor 0.23 버전과의 차이
			metric_container_id_array := strings.Split(metrics.Container_Id, "-")
			if len(metric_container_id_array) > 4 {
				//if metrics.Container_Id == containerMetric.Container_Id {
				if metric_container_id_array[len(metric_container_id_array)-1] == containerMetric.Container_Id {
					//fmt.Println("================ container id:", containerMetric.Container_Id)
					containerMetric.Interface_Id = metrics.Interface_Id
					containerMetric.Application_Id = metrics.Application_Id
					containerMetric.Application_Name = metrics.Application_Name
					containerMetric.Application_Urls = metrics.Application_Urls
					containerMetric.Application_Index = metrics.Application_Index
					containerMetric.Limits.Disk = metrics.Limits.Disk
					containerMetric.Limits.Memory = metrics.Limits.Memory
					containerMetric.UsageMetrics.MemoryUsageInBytes = metrics.UsageMetrics.MemoryUsageInBytes
					containerMetric.UsageMetrics.DiskUsageInBytes = metrics.UsageMetrics.DiskUsageInBytes
					containerMetric.UsageMetrics.TimeSpentInCPU = metrics.UsageMetrics.TimeSpentInCPU
				}
			} else {
				fmt.Println("================ else container id:", containerMetric.Container_Id, metrics.Container_Id)
			}
		}
		//===================================================================

		s.points = append(s.points, s.containerStatsToPoints(cInfo, containerMetric, stats)...)
		/*s.points = append(s.points, s.memoryStatsToPoints(cInfo, stats)...)
		s.points = append(s.points, s.hugetlbStatsToPoints(cInfo, stats)...)
		s.points = append(s.points, s.perfStatsToPoints(cInfo, stats)...)
		s.points = append(s.points, s.resctrlStatsToPoints(cInfo, stats)...)*/
		s.points = append(s.points, s.containerFilesystemStatsToPoints(cInfo, stats)...)
		if s.readyToFlush() {
			pointsToFlush = s.points
			s.points = make([]*influxdb.Point, 0)
			s.lastWrite = time.Now()
		}
	}()
	if len(pointsToFlush) > 0 {
		/*points := make([]influxdb.Point, len(pointsToFlush))
		for i, p := range pointsToFlush {
			points[i] = *p
		}

		batchTags := map[string]string{tagMachineName: s.machineName}
		bp := influxdb.BatchPoints{
			Points:          points,
			Database:        s.database,
			RetentionPolicy: s.retentionPolicy,
			Tags:            batchTags,
			Time:            stats.Timestamp,
		}
		response, err := s.client.Write(bp)
		if err != nil || checkResponseForErrors(response) != nil {
			return fmt.Errorf("failed to write stats to influxDb - %s", err)
		}*/

		/*ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		var stopChan chan bool = nil*/

		bp, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
			Database:  s.database,
			Precision: "s",
		})
		if err != nil {
			fmt.Println(err)
			/*glog.Fatalf("Failed to create NewBatchPoint: %v", err)*/
		}

		//points := make([]influxdb.Point, len(pointsToFlush))
		for _, p := range pointsToFlush {
			//points[i] = *p
			//fmt.Println("point to save at database ",self.database, i, p)
			bp.AddPoint(p)
		}
		err = s.client.Write(bp)
		if err != nil {
			fmt.Println(err)
			/*glog.Fatalf("Failed to send point to influxdb: %v", err)*/
		}

		/*select {
		case <-ticker.C:
		case <-stopChan:
			return nil
		}*/ //end select
	}
	return nil
}

func (s *influxdbStorage) Close() error {
	s.client = nil
	return nil
}

// Creates a measurement point with a single value field
func makePoint(machineName, cellIp, containerName, name string, containerMetric ContainerMetricsMetadata, value float64) *influxdb.Point {
	var tags map[string]string
	var fields map[string]interface{}
	if containerMetric.Application_Id != "" {
		tags = map[string]string{
			tagName:               name,
			tagMachineName:        machineName,
			tagCellIp:             cellIp,
			tagContainerName:      containerName,
			tagContainerInterface: containerMetric.Interface_Id,
			tagApplicationId:      containerMetric.Application_Id,
			tagApplicationIndex:   containerMetric.Application_Index,
			tagApplicationName:    containerMetric.Application_Name,
			tagApplicationUrl:     containerMetric.Application_Urls[0],
		}
	} else {
		tags = map[string]string{
			tagName:               name,
			tagMachineName:        machineName,
			tagCellIp:             cellIp,
			tagContainerName:      containerName,
			tagContainerInterface: containerMetric.Interface_Id,
		}

	}
	if containerMetric.Application_Id != "" {
		fields = map[string]interface{}{
			fieldValue:   value,
			fieldAppDisk: float64(containerMetric.Limits.Disk * 1024 * 1024),
			fieldAppMem:  float64(containerMetric.Limits.Memory * 1024 * 1024),
		}
	} else {
		fields = map[string]interface{}{
			fieldValue: value,
		}
	}

	/*return &influxdb.Point{
		//Measurement: name,
		Measurement: serContainerMeausement,
		Tags:        tags,
		Fields:      fields,
	}*/
	mkPoint, err := influxdb.NewPoint(serContainerMeausement, tags, fields)
	if err != nil {
		fmt.Println(err)
		/*glog.Fatalf("Failed to create NewPoint for FieldsFsLimit: %v", err)*/
	}

	return mkPoint
}

// Adds additional tags to the existing tags of a point
/*func addTagsToPoint(point *influxdb.Point, tags map[string]string) {
	if point.Tags == nil {
		point.Tags = tags
	} else {
		for k, v := range tags {
			point.Tags[k] = v
		}
	}
}*/

// Checks response for possible errors
/*func checkResponseForErrors(response *influxdb.Response) error {
	const msg = "failed to write stats to influxDb - %s"

	if response != nil && response.Err != nil {
		return fmt.Errorf(msg, response.Err)
	}
	if response != nil && response.Results != nil {
		for _, result := range response.Results {
			if result.Err != nil {
				return fmt.Errorf(msg, result.Err)
			}
			if result.Series != nil {
				for _, row := range result.Series {
					if row.Err != nil {
						return fmt.Errorf(msg, row.Err)
					}
				}
			}
		}
	}
	return nil
}*/

// Some stats have type unsigned integer, but the InfluxDB client accepts only signed integers.
func toSignedIfUnsigned(value interface{}) interface{} {
	switch v := value.(type) {
	case uint64:
		return int64(v)
	case uint32:
		return int32(v)
	case uint16:
		return int16(v)
	case uint8:
		return int8(v)
	case uint:
		return int(v)
	}
	return value
}

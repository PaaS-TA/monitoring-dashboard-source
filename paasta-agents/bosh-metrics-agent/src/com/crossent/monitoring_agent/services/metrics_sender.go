package services

import (
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"code.cloudfoundry.org/lager"
	influxdb "github.com/influxdata/influxdb/client/v2"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

type processStat struct {
	pid       int32
	ppid      int32
	startTime string
	memUsage  uint64
	name      string
}

type MetricSender struct {
	logger   lager.Logger
	influx   *InfluxConfig
	origin   *OriginConfig
	interval time.Duration
	Success  bool
}

type OriginConfig struct {
	Origin         string
	OriginId       string
	OriginIndex    string
	OriginIp       string
	DeployemntName string
}

type InfluxConfig struct {
	InfluxUrl          string
	InfluxUser         string
	InfluxPass         string
	InfluxDatabase     string
	Measurement        string
	ProcessMeasurement string
}

type CpuStats struct {
	CoreNum int
	Status  float64
}

type MemStatus struct {
	Total       float64
	Available   float64
	Used        float64
	UsedPercent float64
}

func NewMetricSender(logger lager.Logger, influx *InfluxConfig, origin *OriginConfig, interval time.Duration) *MetricSender {
	return &MetricSender{
		logger:   logger,
		influx:   influx,
		origin:   origin,
		interval: interval,
		Success:  false,
	}
}

func (f *MetricSender) SendMetricsToInfluxDb(stopChan <-chan bool) error {
	f.logger.Info("influx :", map[string]interface{}{"influxUrl": f.influx.InfluxUrl, "influxdatabase": f.influx.InfluxDatabase})

	var err error
	var c influxdb.Client
	ticker := time.NewTicker(f.interval)

	defer func() {
		ticker.Stop()
		_ = c.Close()
	}()

	for {
		// Make influxDB client
		//c, err := influxdb.NewUDPClient(influxdb.UDPConfig{
		c, err = influxdb.NewUDPClient(influxdb.UDPConfig{
			Addr: f.influx.InfluxUrl,
			//PayloadSize: 4096,
		})
		if err != nil {
			f.logger.Error("#metrics_sender.SendMetricsToInfluxDb  : There is an error during creating influxdb client:", err)
		} else {
			err = f.CollectSave(c)
			if err != nil {
				f.logger.Error("#metrics_sender.SendMetricsToInfluxDb  : There is an error during sending metrics to influxdb:", err)
			}
		}
		if c != nil {
			_ = c.Close()
			c = nil
		}
		select {
		case <-ticker.C:
		case <-stopChan:
			return nil
		} //end select
	}

	//return err
}

func (f *MetricSender) CollectSave(c influxdb.Client) error {
	//===============================================================
	//catch or finally
	/*defer func() {
		if err := recover(); err != nil { //catch
			//fmt.Fprintf(os.Stderr, "Exception at metrics_sender.Collect_Save(): %v\n", err)
			//os.Exit(1)
			return err
		}
	}()*/
	defer func() {
		_ = c.Close()
		c = nil
	}()

	f.Success = false
	//===============================================================
	//fmt.Println("##### Collect_Save Called #####")
	// Create a new point batch
	bp, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Database:  f.influx.InfluxDatabase,
		Precision: "s",
	})

	if err != nil {
		return err
	}

	//======================================================================
	// ========= CPU Status ========
	//======================================================================
	f.emitCpuMetrics(bp)

	//======================================================================
	// ========= Memory Status ========
	//======================================================================
	f.emitMemMetrics(bp)

	//======================================================================
	// ========= Disk Status ========
	//======================================================================
	f.emitDiskMetrics(bp)

	//======================================================================
	// ========= Network Status ========
	//======================================================================
	f.emitNetworkMetrics(bp)

	//======================================================================
	// ========= Processes Status ========
	//======================================================================
	f.emitProcessMetrics(bp)

	// Write the batch
	err = c.Write(bp)
	if err != nil {
		return err
	}

	f.Success = true
	return nil
}

// Creates a measurement point with a single value field
func (f *MetricSender) makePoint(name string, value float64) *influxdb.Point {
	var tags map[string]string
	var fields map[string]interface{}
	tags = map[string]string{
		"origin":     f.origin.Origin,
		"id":         f.origin.OriginId,
		"ip":         f.origin.OriginIp,
		"index":      f.origin.OriginIndex,
		"deployment": f.origin.DeployemntName,
		"metricname": name,
	}
	fields = map[string]interface{}{
		"value": value,
	}

	mkPoint, err := influxdb.NewPoint(f.influx.Measurement, tags, fields)
	if err != nil {
		f.logger.Error("#metrics_collector.makePoint : error caused during making a new point", err)
	}

	return mkPoint
}

/*
 Description: VM - CPU Info metrics
*/
func (f *MetricSender) emitCpuMetrics(bp influxdb.BatchPoints) {
	numCpu := runtime.NumCPU()
	duration := time.Duration(1) * time.Second
	c, err := cpu.Percent(duration, true)
	if err != nil {
		f.logger.Error("getting cpu metrics error %v", err)
	}

	//var cpuStatusArray []CpuStats
	//var load1AvgStat, load5AvgStat, load15AvgStat float64
	for k, percent := range c {
		//var cpuStatus CpuStats
		// Check for slightly greater then 100% to account for any rounding issues.
		if percent < 0.0 || percent > 100.0001*float64(numCpu) {
			f.logger.Info("CPUPercent value is invalid: %f", lager.Data{"percent": percent})
		} else {
			//cpuStatus.CoreNum = k
			//cpuStatus.Status = float64(percent)
			bp.AddPoint(f.makePoint(fmt.Sprintf("cpuStats.core.%d", k), percent))

		}
		//cpuStatusArray = append(cpuStatusArray, cpuStatus)
	}

	//============ CPU Load Average : Only support linux & freebsd ==============
	h, err := host.Info()
	if err != nil {
		return
	}

	if h.OS == "linux" || h.OS == "freebsd" {
		loadAvgStat, err := load.Avg()
		if err != nil {
			f.logger.Error("LoadAvgStats: failed to get LoadAvg information: %v", err)
		}
		bp.AddPoint(f.makePoint("cpuStats.LoadAvg1Stats", loadAvgStat.Load1))
		bp.AddPoint(f.makePoint("cpuStats.LoadAvg5Stats", loadAvgStat.Load5))
		bp.AddPoint(f.makePoint("cpuStats.LoadAvg15Stats", loadAvgStat.Load15))
		/*load1AvgStat = float64(loadAvgStat.Load1)
		load5AvgStat = float64(loadAvgStat.Load5)
		load15AvgStat = float64(loadAvgStat.Load15)*/
	}
	//===========================================================================
	/*cpuStatus, load1avg, load5avg, load15avg := f.emitCpuMetrics(bp)
	for _, v := range cpuStatus{
		//fmt.Println(fmt.Sprintf("cpuStats.core.%d", v.CoreNum), float64(v.Status))
		bp.AddPoint(f.makePoint(f.origin, f.cellIp, fmt.Sprintf("cpuStats.core.%d", v.CoreNum), float64(v.Status)))
	}

	fmt.Println("load 1 avg : ", load1avg)
	fmt.Println("load 5 avg : ", load5avg)
	fmt.Println("load 15 avg : ", load15avg)

	bp.AddPoint(f.makePoint(f.origin, f.cellIp, "cpuStats.LoadAvg1Stats", float64(load1avg)))
	bp.AddPoint(f.makePoint(f.origin, f.cellIp, "cpuStats.LoadAvg5Stats", float64(load5avg)))
	bp.AddPoint(f.makePoint(f.origin, f.cellIp, "cpuStats.LoadAvg15Stats", float64(load15avg)))
	//======================================================================

	return cpuStatusArray, load1AvgStat, load5AvgStat, load15AvgStat*/
}

/*
 Description: VM - Memory Info metrics
*/
func (f *MetricSender) emitMemMetrics(bp influxdb.BatchPoints) {
	m, err := mem.VirtualMemory()
	if err != nil {
		f.logger.Error("MemStats: failed to get Memory Info: %v", err)
	} else {
		/*fmt.Println("memoryStats.TotalMemory", float64(m.Total))
		fmt.Println("memoryStats.AvailableMemory", float64(m.Available))
		fmt.Println("memoryStats.UsedMemory", float64(m.Used))
		fmt.Println("memoryStats.UsedPercent", float64(m.UsedPercent))*/
		bp.AddPoint(f.makePoint("memoryStats.TotalMemory", float64(m.Total)))
		bp.AddPoint(f.makePoint("memoryStats.FreeMemory", float64(m.Free)))
		bp.AddPoint(f.makePoint("memoryStats.BuffersMemory", float64(m.Buffers)))
		bp.AddPoint(f.makePoint("memoryStats.CachedMemory", float64(m.Cached)))
		bp.AddPoint(f.makePoint("memoryStats.SwapCachedMemory", float64(m.SwapCached)))
		bp.AddPoint(f.makePoint("memoryStats.AvailableMemory", float64(m.Available)))
		bp.AddPoint(f.makePoint("memoryStats.UsedMemory", float64(m.Used)))
		bp.AddPoint(f.makePoint("memoryStats.UsedPercent", m.UsedPercent))
	}
}

/*
 Description: VM - Disk Info metrics
*/
func (f *MetricSender) emitDiskMetrics(bp influxdb.BatchPoints) {
	if runtime.GOOS == "windows" {
		var pathKey []string
		diskios, _ := disk.IOCounters()
		for key, value := range diskios {
			pathKey = append(pathKey, key)

			/*fmt.Println(fmt.Sprintf("diskIOStats.%s.readCount", key), float64(value.ReadCount))
			fmt.Println(fmt.Sprintf("diskIOStats.%s.writeCount", key), float64(value.WriteCount))
			fmt.Println(fmt.Sprintf("diskIOStats.%s.readBytes", key), float64(value.ReadBytes))
			fmt.Println(fmt.Sprintf("diskIOStats.%s.writeBytes", key), float64(value.WriteBytes))
			fmt.Println(fmt.Sprintf("diskIOStats.%s.readTime", key), float64(value.ReadTime))
			fmt.Println(fmt.Sprintf("diskIOStats.%s.writeTime", key), float64(value.WriteTime))
			fmt.Println(fmt.Sprintf("diskIOStats.%s.ioTime", key), float64(value.IoTime))*/

			bp.AddPoint(f.makePoint(fmt.Sprintf("diskIOStats.%s.readCount", key), float64(value.ReadCount)))
			bp.AddPoint(f.makePoint(fmt.Sprintf("diskIOStats.%s.writeCount", key), float64(value.WriteCount)))
			bp.AddPoint(f.makePoint(fmt.Sprintf("diskIOStats.%s.readBytes", key), float64(value.ReadBytes)))
			bp.AddPoint(f.makePoint(fmt.Sprintf("diskIOStats.%s.writeBytes", key), float64(value.WriteBytes)))
			bp.AddPoint(f.makePoint(fmt.Sprintf("diskIOStats.%s.readTime", key), float64(value.ReadTime)))
			bp.AddPoint(f.makePoint(fmt.Sprintf("diskIOStats.%s.writeTime", key), float64(value.WriteTime)))
			bp.AddPoint(f.makePoint(fmt.Sprintf("diskIOStats.%s.ioTime", key), float64(value.IoTime)))
		}
		for _, value := range pathKey {
			d, err := disk.Usage(value)
			if err != nil {
				f.logger.Error("getting disk info error %v", err)
			}
			/*fmt.Println(fmt.Sprintf("diskStats.windows.%s.Total", d.Path), float64(d.Total))
			fmt.Println(fmt.Sprintf("diskStats.windows.%s.Used", d.Path), float64(d.Used))
			fmt.Println(fmt.Sprintf("diskStats.windows.%s.Available", d.Path), float64(d.Free))
			fmt.Println(fmt.Sprintf("diskStats.windows.%s.Usage", d.Path), float64(d.UsedPercent))*/

			bp.AddPoint(f.makePoint(fmt.Sprintf("diskStats.windows.%s.Total", d.Path), float64(d.Total)))
			bp.AddPoint(f.makePoint(fmt.Sprintf("diskStats.windows.%s.Used", d.Path), float64(d.Used)))
			bp.AddPoint(f.makePoint(fmt.Sprintf("diskStats.windows.%s.Available", d.Path), float64(d.Free)))
			bp.AddPoint(f.makePoint(fmt.Sprintf("diskStats.windows.%s.Usage", d.Path), d.UsedPercent))
		}
	} else {

		diskparts, err := disk.Partitions(false)
		if err != nil {
			f.logger.Error("get disk partitions error: %v", err)
		}
		for _, partition := range diskparts {
			//fmt.Println("partition KEY:", key, "value:", partition)
			mountpoints := strings.Split(partition.Device, "/")
			d, err := disk.Usage(partition.Mountpoint)
			if err != nil {
				f.logger.Error("getting disk info error %v", err)
				//return err
			}
			/*fmt.Println(fmt.Sprintf("diskStats.%s.Total", partition.Mountpoint, float64(d.Total)))
			fmt.Println(fmt.Sprintf("diskStats.%s.Used", partition.Mountpoint, float64(d.Used)))
			fmt.Println(fmt.Sprintf("diskStats.%s.Available", partition.Mountpoint, float64(d.Free)))
			fmt.Println(fmt.Sprintf("diskStats.%s.Usage", partition.Mountpoint, float64(d.UsedPercent)))*/

			bp.AddPoint(f.makePoint(fmt.Sprintf("diskStats.%s.Total", partition.Mountpoint), float64(d.Total)))
			bp.AddPoint(f.makePoint(fmt.Sprintf("diskStats.%s.Used", partition.Mountpoint), float64(d.Used)))
			bp.AddPoint(f.makePoint(fmt.Sprintf("diskStats.%s.Available", partition.Mountpoint), float64(d.Free)))
			bp.AddPoint(f.makePoint(fmt.Sprintf("diskStats.%s.Usage", partition.Mountpoint), d.UsedPercent))

			//Newly Added - Disk I/O (2017.04)
			diskios, _ := disk.IOCounters()
			for key, value := range diskios {
				if mountpoints[len(mountpoints)-1] == key {
					/*fmt.Println(fmt.Sprintf("diskIOStats.%s.%s.readCount", partition.Mountpoint, key), float64(value.ReadCount))
					fmt.Println(fmt.Sprintf("diskIOStats.%s.%s.writeCount", partition.Mountpoint, key), float64(value.WriteCount))
					fmt.Println(fmt.Sprintf("diskIOStats.%s.%s.readBytes", partition.Mountpoint, key), float64(value.ReadBytes))
					fmt.Println(fmt.Sprintf("diskIOStats.%s.%s.writeBytes", partition.Mountpoint, key), float64(value.WriteBytes))
					fmt.Println(fmt.Sprintf("diskIOStats.%s.%s.readTime", partition.Mountpoint, key), float64(value.ReadTime))
					fmt.Println(fmt.Sprintf("diskIOStats.%s.%s.writeTime", partition.Mountpoint, key), float64(value.WriteTime))
					fmt.Println(fmt.Sprintf("diskIOStats.%s.%s.ioTime", partition.Mountpoint, key), float64(value.IoTime))*/

					bp.AddPoint(f.makePoint(fmt.Sprintf("diskIOStats.%s.%s.readCount", partition.Mountpoint, key), float64(value.ReadCount)))
					bp.AddPoint(f.makePoint(fmt.Sprintf("diskIOStats.%s.%s.writeCount", partition.Mountpoint, key), float64(value.WriteCount)))
					bp.AddPoint(f.makePoint(fmt.Sprintf("diskIOStats.%s.%s.readBytes", partition.Mountpoint, key), float64(value.ReadBytes)))
					bp.AddPoint(f.makePoint(fmt.Sprintf("diskIOStats.%s.%s.writeBytes", partition.Mountpoint, key), float64(value.WriteBytes)))
					bp.AddPoint(f.makePoint(fmt.Sprintf("diskIOStats.%s.%s.readTime", partition.Mountpoint, key), float64(value.ReadTime)))
					bp.AddPoint(f.makePoint(fmt.Sprintf("diskIOStats.%s.%s.writeTime", partition.Mountpoint, key), float64(value.WriteTime)))
					bp.AddPoint(f.makePoint(fmt.Sprintf("diskIOStats.%s.%s.ioTime", partition.Mountpoint, key), float64(value.IoTime)))
				}
			}
		}
	}
}

/*
 Description: VM - Network Info metrics
*/
func (f *MetricSender) emitNetworkMetrics(bp influxdb.BatchPoints) {
	nifs, err := net.Interfaces()
	if err != nil {
		f.logger.Error("getting network interface info error %v", err)
		//return err
	}

	for _, intf := range nifs {
		//fmt.Println(fmt.Sprintf("networkInterface.%s.%s", intf.Name, "MTU"), float64(intf.MTU))
		bp.AddPoint(f.makePoint(fmt.Sprintf("networkInterface.%s.%s", intf.Name, "MTU"), float64(intf.MTU)))
	}

	ios, err := net.IOCounters(true)
	for _, value := range ios {
		/*fmt.Println(fmt.Sprintf("networkIOStats.%s.bytesRecv", value.Name), float64(value.BytesRecv))
		fmt.Println(fmt.Sprintf("networkIOStats.%s.bytesSent", value.Name), float64(value.BytesSent))
		fmt.Println(fmt.Sprintf("networkIOStats.%s.packetRecv", value.Name), float64(value.PacketsRecv))
		fmt.Println(fmt.Sprintf("networkIOStats.%s.packetSent", value.Name), float64(value.PacketsSent))
		fmt.Println(fmt.Sprintf("networkIOStats.%s.dropIn", value.Name), float64(value.Dropin))
		fmt.Println(fmt.Sprintf("networkIOStats.%s.dropOut", value.Name), float64(value.Dropout))
		fmt.Println(fmt.Sprintf("networkIOStats.%s.errIn", value.Name), float64(value.Errin))
		fmt.Println(fmt.Sprintf("networkIOStats.%s.errOut", value.Name), float64(value.Errout))*/
		bp.AddPoint(f.makePoint(fmt.Sprintf("networkIOStats.%s.bytesRecv", value.Name), float64(value.BytesRecv)))
		bp.AddPoint(f.makePoint(fmt.Sprintf("networkIOStats.%s.bytesSent", value.Name), float64(value.BytesSent)))
		bp.AddPoint(f.makePoint(fmt.Sprintf("networkIOStats.%s.packetRecv", value.Name), float64(value.PacketsRecv)))
		bp.AddPoint(f.makePoint(fmt.Sprintf("networkIOStats.%s.packetSent", value.Name), float64(value.PacketsSent)))
		bp.AddPoint(f.makePoint(fmt.Sprintf("networkIOStats.%s.dropIn", value.Name), float64(value.Dropin)))
		bp.AddPoint(f.makePoint(fmt.Sprintf("networkIOStats.%s.dropOut", value.Name), float64(value.Dropout)))
		bp.AddPoint(f.makePoint(fmt.Sprintf("networkIOStats.%s.errIn", value.Name), float64(value.Errin)))
		bp.AddPoint(f.makePoint(fmt.Sprintf("networkIOStats.%s.errOut", value.Name), float64(value.Errout)))
	}
}

/*
 Description: VM - Process Info metrics
*/
func (f *MetricSender) emitProcessMetrics(bp influxdb.BatchPoints) {
	procs, err := process.Pids()
	if err != nil {
		f.logger.Error("getting processes error %v", err)
		//return err
	}

	pStatArray := make([]processStat, 0)

	for _, value := range procs {
		p, err := process.NewProcess(value)
		if err != nil {
			f.logger.Error("getting single process info error %v", err)
			//return err
		} else {
			var pStat processStat
			ct, _ := p.CreateTime()
			sTimestamp := strconv.FormatInt(ct, 10)

			//window 환경에서는 모든 프로세스를 조회하기 때문에 많은 시간이 소요된다.
			//이를 방지하기 위해 프로세스시작 시간을 통해 제어한다.
			if len(sTimestamp) >= 10 {
				pStat.startTime = sTimestamp[:10]
				pStat.pid = p.Pid
				pp, _ := p.Ppid()
				pStat.ppid = pp

				pname, _ := p.Name()
				pStat.name = pname

				m, err := p.MemoryInfo()
				if err == nil {
					pStat.memUsage = m.RSS
				}
				pStatArray = append(pStatArray, pStat)
			}
		}
	}

	//Memroy 점유 크기별로 Sorting
	//slice.Sort(pStatArray[:], func(i, j int) bool {
	sort.SliceStable(pStatArray[:], func(i, j int) bool {
		return pStatArray[i].memUsage > pStatArray[j].memUsage
	})

	var tags map[string]string
	var fields map[string]interface{}
	var index, startTime int64
	for _, ps := range pStatArray {
		if index > 20 {
			break
		}
		/*fmt.Println(fmt.Sprintf("processStats.%d.%s.pid",index, ps.name), float64(ps.pid))
		fmt.Println(fmt.Sprintf("processStats.%d.%s.ppid",index, ps.name), float64(ps.ppid))
		fmt.Println(fmt.Sprintf("processStats.%d.%s.memUsage",index, ps.name), float64(ps.memUsage))*/

		if startTime, err = strconv.ParseInt(ps.startTime, 10, 0); err != nil {
			startTime = 0
		}

		tags = map[string]string{
			"origin":     f.origin.Origin,
			"ip":         f.origin.OriginIp,
			"id":         f.origin.OriginId,
			"index":      f.origin.OriginIndex,
			"deployment": f.origin.DeployemntName,
			"metricname": "processStats",
			"proc_name":  ps.name,
		}
		fields = map[string]interface{}{
			"proc_index": index,
			"proc_pid":   ps.pid,
			"proc_ppid":  ps.ppid,
			"mem_usage":  ps.memUsage,
			"start_time": startTime,
		}

		mkPoint, err := influxdb.NewPoint(f.influx.ProcessMeasurement, tags, fields)
		if err != nil {
			f.logger.Error("makePoint for process: error caused during making a new point", err)
		}
		bp.AddPoint(mkPoint)
		/*bp.AddPoint(f.makePoint(f.origin, f.cellIp, fmt.Sprintf("processStats.%d.%s.pid",index, ps.name), float64(ps.pid)))
		bp.AddPoint(f.makePoint(f.origin, f.cellIp, fmt.Sprintf("processStats.%d.%s.ppid",index, ps.name), float64(ps.ppid)))
		if flt, err := strconv.ParseFloat(ps.startTime, 64); err == nil {
			//fmt.Println(fmt.Sprintf("processStats.%d.%s.startTime", index, ps.name),  flt)
			bp.AddPoint(f.makePoint(f.origin, f.cellIp, fmt.Sprintf("processStats.%d.%s.startTime", index, ps.name),  flt))
		}
		bp.AddPoint(f.makePoint(f.origin, f.cellIp, fmt.Sprintf("processStats.%d.%s.memUsage",index, ps.name), float64(ps.memUsage)))*/
		index++
	}

}

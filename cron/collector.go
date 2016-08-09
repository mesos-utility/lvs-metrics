package cron

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/google/seesaw/ipvs"
	"github.com/mesos-utility/lvs-metrics/g"
	"github.com/open-falcon/common/model"
)

var STATS = []string{
	"total.conns",
	"in.packets",
	"out.packets",
	"in.bytes",
	"out.bytes"}

//var IPVSFILE string
var IPVSSTATSFILE string

func init() {
	//IPVSFILE = "/proc/net/ip_vs"
	IPVSSTATSFILE = "/proc/net/ip_vs_stats"

	if os.Getenv("LVSDEV") != "" {
		pwd, _ := os.Getwd()

		//IPVSFILE = fmt.Sprintf("%s/resource/ip_vs", pwd)
		IPVSSTATSFILE = fmt.Sprintf("%s/resource/ip_vs_stats", pwd)
	}
}

func Collect() {
	if !g.Config().Transfer.Enable {
		glog.Warningf("Open falcon transfer is not enabled!!!")
		return
	}

	if g.Config().Transfer.Addr == "" {
		glog.Warningf("Open falcon transfer addr is null!!!")
		return
	}

	go collect()
}

func collect() {
	// start collect data for lvs cluster.
	var interval int64 = g.Config().Transfer.Interval
	var ticker = time.NewTicker(time.Duration(interval) * time.Second)

	ipvs.Init()
	for {
		<-ticker.C

		mvs := []*model.MetricValue{}

		// Collect metrics from /proc/net/ip_vs
		vips, err := GetIPVSStats()
		if err != nil {
			glog.Errorf("%s", err.Error())
		}
		mvs, _ = ConvertVIPs2Metrics(vips)
		g.SendMetrics(mvs)

		// Collect metrics from /proc/net/ip_vs_stats
		mvs, err = ParseIPVSStats(IPVSSTATSFILE)
		if os.IsNotExist(err) {
			glog.Errorf("%s", err.Error())
		}
		g.SendMetrics(mvs)
	}
}

func ConvertVIPs2Metrics(vips []*VirtualIPPoint) (metrics []*model.MetricValue, err error) {
	if len(vips) <= 0 {
		return nil, nil
	}

	var tags string
	var attachtags = g.Config().AttachTags
	var interval int64 = g.Config().Transfer.Interval
	if attachtags != "" {
		tags = attachtags
	}

	hostname, _ := g.Hostname()
	now := time.Now().Unix()
	var metric *model.MetricValue
	for _, vip := range vips {
		var tag string
		if tags != "" {
			tag = fmt.Sprintf("%s,vip=%s,port=%d", tags, vip.IP, vip.Port)
		} else {
			tag = fmt.Sprintf("vip=%s,port=%d", vip.IP, vip.Port)
		}
		metric = &model.MetricValue{
			Endpoint:  hostname,
			Metric:    "lvs.vip.active_conn",
			Value:     vip.ActiveConns,
			Timestamp: now,
			Step:      interval,
			Type:      "GAUGE",
			Tags:      tag,
		}
		metrics = append(metrics, metric)

		metric = &model.MetricValue{
			Endpoint:  hostname,
			Metric:    "lvs.vip.inact_conn",
			Value:     vip.InactiveConns,
			Timestamp: now,
			Step:      interval,
			Type:      "GAUGE",
			Tags:      tag,
		}
		metrics = append(metrics, metric)

		metric = &model.MetricValue{
			Endpoint:  hostname,
			Metric:    "lvs.vip.rs_num",
			Value:     vip.RealServerNum,
			Timestamp: now,
			Step:      interval,
			Type:      "GAUGE",
			Tags:      tag,
		}
		metrics = append(metrics, metric)

		metric = &model.MetricValue{
			Endpoint:  hostname,
			Metric:    "lvs.vip.conns",
			Value:     vip.Connections,
			Timestamp: now,
			Step:      interval,
			Type:      "COUNTER",
			Tags:      tag,
		}
		metrics = append(metrics, metric)

		metric = &model.MetricValue{
			Endpoint:  hostname,
			Metric:    "lvs.vip.inpkts",
			Value:     vip.PacketsIn,
			Timestamp: now,
			Step:      interval,
			Type:      "COUNTER",
			Tags:      tag,
		}
		metrics = append(metrics, metric)

		metric = &model.MetricValue{
			Endpoint:  hostname,
			Metric:    "lvs.vip.outpkts",
			Value:     vip.PacketsOut,
			Timestamp: now,
			Step:      interval,
			Type:      "COUNTER",
			Tags:      tag,
		}
		metrics = append(metrics, metric)

		metric = &model.MetricValue{
			Endpoint:  hostname,
			Metric:    "lvs.vip.inbytes",
			Value:     vip.BytesIn,
			Timestamp: now,
			Step:      interval,
			Type:      "COUNTER",
			Tags:      tag,
		}
		metrics = append(metrics, metric)

		metric = &model.MetricValue{
			Endpoint:  hostname,
			Metric:    "lvs.vip.outbytes",
			Value:     vip.BytesOut,
			Timestamp: now,
			Step:      interval,
			Type:      "COUNTER",
			Tags:      tag,
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

// Parse /proc/net/ip_vs_stats
func ParseIPVSStats(file string) (metrics []*model.MetricValue, err error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	strcon := string(content)
	lines := strings.Split(strcon, "\n")
	if len(lines) < 6 {
		return nil, fmt.Errorf("ip_vs_stats content invalid")
	}
	array := strings.Fields(lines[5])

	var attachtags = g.Config().AttachTags
	var interval int64 = g.Config().Transfer.Interval
	now := time.Now().Unix()
	hostname, _ := g.Hostname()
	for i, v := range STATS {
		value, _ := strconv.ParseInt(array[i], 16, 0)
		metricName := fmt.Sprintf("lvs.%s", v)
		metric := &model.MetricValue{
			Endpoint:  hostname,
			Metric:    metricName,
			Value:     value,
			Timestamp: now,
			Step:      interval,
			Type:      "GAUGE",
			Tags:      attachtags,
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

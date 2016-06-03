package cron

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/mesos-utility/lvs-metrics/g"
	"github.com/open-falcon/common/model"
)

var STATS = []string{
	"total.conns",
	"in.packets",
	"out.packets",
	"in.bytes",
	"out.bytes"}

const (
	IPVSFILE      = "/proc/net/ip_vs"
	IPVSSTATSFILE = "/proc/net/ip_vs_stats"
)

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

	for {
		<-ticker.C

		mvs := []*model.MetricValue{}

		// Collect metrics from /proc/net/ip_vs
		vips, err := ParseIPVS(IPVSFILE)
		if os.IsNotExist(err) {
			glog.Fatalf("%s", err.Error())
		}
		mvs, _ = ConvertVIPs2Metrics(vips)
		g.SendMetrics(mvs)

		// Collect metrics from /proc/net/ip_vs_stats
		mvs, err = ParseIPVSStats(IPVSSTATSFILE)
		if os.IsNotExist(err) {
			glog.Fatalf("%s", err.Error())
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
	for _, vip := range vips {
		var tag string
		if tags != "" {
			tag = fmt.Sprintf("%s,vip=%s,port=%d", tags, vip.IP, vip.Port)
		} else {
			tag = fmt.Sprintf("vip=%s,port=%d", vip.IP, vip.Port)
		}
		metric := &model.MetricValue{
			Endpoint:  hostname,
			Metric:    "lvs.active.conn",
			Value:     vip.TotalActiveConn,
			Timestamp: now,
			Step:      interval,
			Type:      "USAGE",
			Tags:      tag,
		}
		metrics = append(metrics, metric)

		metric = &model.MetricValue{
			Endpoint:  hostname,
			Metric:    "lvs.inact.conn",
			Value:     vip.TotalInActConn,
			Timestamp: now,
			Step:      interval,
			Type:      "USAGE",
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
			Type:      "USAGE",
			Tags:      attachtags,
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}
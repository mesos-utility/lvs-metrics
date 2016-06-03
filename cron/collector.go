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
	var attachtags = g.Config().AttachTags
	var interval int64 = g.Config().Transfer.Interval
	//var stats = make(map[string]string)
	var ticker = time.NewTicker(time.Duration(interval) * time.Second)

	for {
	REST:
		<-ticker.C
		hostname, err := g.Hostname()
		if err != nil {
			goto REST
		}

		mvs := []*model.MetricValue{}
		var tags string
		if attachtags != "" {
			tags = attachtags
		}

		now := time.Now().Unix()
		vips, err := ParseIPVS(IPVSFILE)
		if os.IsNotExist(err) {
			glog.Fatalf("%s", err.Error())
		}
		//glog.Infof("%v\n", len(vips))
		for _, vip := range vips {
			tag := fmt.Sprintf("%s,vip=%s:%d", tags, vip.IP, vip.Port)
			metric := &model.MetricValue{
				Endpoint:  hostname,
				Metric:    "lvs.ActiveConn",
				Value:     vip.TotalActiveConn,
				Timestamp: now,
				Step:      interval,
				Type:      "USAGE",
				Tags:      tag,
			}
			mvs = append(mvs, metric)
			//glog.Infof("%v\n", metric)

			metric = &model.MetricValue{
				Endpoint:  hostname,
				Metric:    "lvs.InActConn",
				Value:     vip.TotalInActConn,
				Timestamp: now,
				Step:      interval,
				Type:      "USAGE",
				Tags:      tag,
			}
			mvs = append(mvs, metric)
			//glog.Infof("%v\n", metric)
		}

		metrics, err := ParseIPVSStats(IPVSSTATSFILE)
		if os.IsNotExist(err) {
			glog.Fatalf("%s", err.Error())
		}
		for _, metric := range metrics {
			mvs = append(mvs, metric)
		}

		g.SendMetrics(mvs)
		mvs = nil
	}
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
		//glog.Infof("%v\n", metric)
	}

	return metrics, nil
}

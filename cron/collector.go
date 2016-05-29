package cron

import (
	//"fmt"
	//"net"
	//"regexp"
	"time"

	"github.com/golang/glog"
	"github.com/mesos-utility/lvs-metrics/g"
	"github.com/open-falcon/common/model"
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

	addrs := g.Config().Daemon.Addrs
	if !g.Config().Daemon.Enable {
		glog.Warningf("Daemon collect not enabled in cfg.json!!!")

		if len(addrs) < 1 {
			glog.Warningf("Not set addrs of daemon in cfg.json!!!")
		}
		return
	}

	go collect(addrs)
}

func collect(addrs []string) {
	// start collect data for lvs cluster.
	//var attachtags = g.Config().AttachTags
	var interval int64 = g.Config().Transfer.Interval
	//var stats = make(map[string]string)
	var ticker = time.NewTicker(time.Duration(interval) * time.Second)

	for {
		//REST:
		<-ticker.C
		//hostname, err := g.Hostname()
		//if err != nil {
		//	goto REST
		//}

		mvs := []*model.MetricValue{}
		//for _, addr := range addrs {
		//	var tags string
		//	if attachtags != "" {
		//		tags = fmt.Sprintf("%s,%s", attachtags)
		//	}

		//	now := time.Now().Unix()
		//	var suffix, vtype string
		//	for k, v := range stats {
		//		if _, ok := gaugess[k]; ok {
		//			suffix = ""
		//			vtype = "GAUGE"
		//		} else {
		//			suffix = "_cps"
		//			vtype = "COUNTER"
		//		}

		//		key := fmt.Sprintf("lvs.%s%s", k, suffix)

		//		metric := &model.MetricValue{
		//			Endpoint:  hostname,
		//			Metric:    key,
		//			Value:     v,
		//			Timestamp: now,
		//			Step:      interval,
		//			Type:      vtype,
		//			Tags:      tags,
		//		}

		//		mvs = append(mvs, metric)
		//		//glog.Infof("%v\n", metric)
		//	}
		//}
		g.SendMetrics(mvs)
		mvs = nil
	}
}

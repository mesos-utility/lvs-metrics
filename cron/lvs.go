package cron

import (
	"fmt"

	"github.com/google/seesaw/ipvs"
)

type RealServer struct {
	IPPort     string
	ActiveConn int
	InActConn  int
}

func NewRealServer(end string, actconn, inactconn int) *RealServer {
	return &RealServer{
		IPPort:     end,
		ActiveConn: actconn,
		InActConn:  inactconn,
	}
}

type VirtualIPPoint struct {
	IP            string
	Port          int
	ActiveConns   uint32
	InactiveConns uint32
	RealServerNum int

	// stats
	Connections uint32
	PacketsIn   uint32
	PacketsOut  uint32
	BytesIn     uint64
	BytesOut    uint64
	// Realservers     [](*RealServer)
}

func NewVirtualIPPoint(ip string, port int, actconn, inactconn uint32) *VirtualIPPoint {
	return &VirtualIPPoint{
		IP:            ip,
		Port:          port,
		ActiveConns:   actconn,
		InactiveConns: inactconn,
	}
}

func GetIPVSStats() (vips []*VirtualIPPoint, err error) {
	svcs, err := ipvs.GetServices()
	if err != nil {
		return nil, err
	}

	var vip *VirtualIPPoint
	var ActiveConns uint32
	var InactiveConns uint32
	var RsNum int
	//var PersistConns uint32
	for _, svc := range svcs {
		ActiveConns = 0
		InactiveConns = 0
		RsNum = len(svc.Destinations)

		for _, dest := range svc.Destinations {
			ActiveConns += dest.Statistics.ActiveConns
			InactiveConns += dest.Statistics.InactiveConns
			//PersistConns += dest.Statistics.PersistConns
		}

		ipstr := fmt.Sprintf("%v", svc.Address)
		vip = &VirtualIPPoint{
			IP:            ipstr,
			Port:          int(svc.Port),
			ActiveConns:   ActiveConns,
			InactiveConns: InactiveConns,
			RealServerNum: RsNum,

			Connections: svc.Statistics.Connections,
			PacketsIn:   svc.Statistics.PacketsIn,
			PacketsOut:  svc.Statistics.PacketsOut,
			BytesIn:     svc.Statistics.BytesIn,
			BytesOut:    svc.Statistics.BytesOut,
		}

		vips = append(vips, vip)
	}

	return vips, nil
}

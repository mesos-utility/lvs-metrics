package cron

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var RE = regexp.MustCompile(`^TCP|^UDP`)

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
	IP              string
	Port            int
	TotalActiveConn int
	TotalInActConn  int
	// Realservers     [](*RealServer)
}

func NewVirtualIPPoint(ip string, port, actconn, inactconn int) *VirtualIPPoint {
	return &VirtualIPPoint{
		IP:              ip,
		Port:            port,
		TotalActiveConn: actconn,
		TotalInActConn:  inactconn,
	}
}

// Parse /proc/net/ip_vs
func ParseIPVS(file string) (vips []*VirtualIPPoint, err error) {
	var line string
	var vip *VirtualIPPoint
	var totalAct, totalInact int

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// read line by line for parse.
	r := bufio.NewReader(f)
	for {
		line, err = r.ReadString('\n')
		if err != nil {
			break
		}

	CONT:
		totalAct = 0
		totalInact = 0

		if find := RE.MatchString(line); find {
			array := strings.Fields(line)
			pair := strings.SplitN(array[1], ":", 2)
			ipstr, _ := Hex2IPV4(pair[0])
			port, _ := strconv.ParseInt(pair[1], 16, 0)

			for {
				line, err = r.ReadString('\n')
				if err != nil {
					vip = &VirtualIPPoint{
						IP:              ipstr,
						Port:            int(port),
						TotalActiveConn: totalAct,
						TotalInActConn:  totalInact,
					}
					vips = append(vips, vip)
					break
				}

				if strings.Contains(line, "->") {
					array := strings.Fields(line)
					act, _ := strconv.ParseInt(array[4], 16, 0)
					inact, _ := strconv.ParseInt(array[5], 16, 0)
					totalAct += int(act)
					totalInact += int(inact)
				} else {
					vip = &VirtualIPPoint{
						IP:              ipstr,
						Port:            int(port),
						TotalActiveConn: totalAct,
						TotalInActConn:  totalInact,
					}
					vips = append(vips, vip)
					goto CONT
				}
			}
		}
	}

	return vips, nil
}

func Hex2IPV4(hexip string) (ip string, err error) {
	if len(hexip) != 8 {
		return "", fmt.Errorf("Invalid format of ipv4.")
	}

	var array []int
	for i := 0; i < 4; i++ {
		pos := 2 * (i + 1)
		num, _ := strconv.ParseInt(hexip[2*i:pos], 16, 0)
		array = append(array, int(num))
	}

	if len(array) != 4 {
		return "", fmt.Errorf("Invalid format of ipv4.")
	} else {
		return fmt.Sprintf("%v.%v.%v.%v", array[0], array[1], array[2], array[3]), nil
	}
}

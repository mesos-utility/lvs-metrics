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
	var txt string
	var totalAct, totalInact int

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
	CONT:
		totalAct = 0
		totalInact = 0
		txt = scanner.Text()
		if find := RE.MatchString(txt); find {
			array := strings.Fields(txt)
			pair := strings.SplitN(array[1], ":", 2)
			ipstr, _ := Hex2IPV4(pair[0])
			port, _ := strconv.ParseInt(pair[1], 16, 0)

			for scanner.Scan() {
				txt = scanner.Text()
				if strings.Contains(txt, "->") {
					array := strings.Fields(txt)
					act, _ := strconv.ParseInt(array[4], 16, 0)
					inact, _ := strconv.ParseInt(array[5], 16, 0)
					totalAct += int(act)
					totalInact += int(inact)

					vip := &VirtualIPPoint{
						IP:              ipstr,
						Port:            int(port),
						TotalActiveConn: totalAct,
						TotalInActConn:  totalInact,
					}
					vips = append(vips, vip)
				} else {
					goto CONT
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return vips, err
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

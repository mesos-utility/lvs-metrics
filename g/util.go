package g

import (
	"fmt"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

// get hostname
func Hostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		glog.Warningf("ERROR: os.Hostname() fail %v", err)
	}
	return hostname, err
}

// calculate metric ratio:
//   string(100 * metrics[0] / sum(metrics))
func CalculateMetricRatio(metrics ...string) string {
	if len(metrics) < 1 {
		return "0"
	} else if len(metrics) == 1 {
		return metrics[0]
	}

	first, err := strconv.ParseFloat(strings.TrimSpace(metrics[0]), 64)
	if err != nil {
		return "0"
	}
	var total float64 = 0

	for _, metric := range metrics {
		fmetric, err := strconv.ParseFloat(strings.TrimSpace(metric), 64)
		if err != nil {
			total = total + 0
		} else {
			total = total + fmetric
		}
	}

	if total == 0 {
		return "0"
	} else {
		return fmt.Sprintf("%.2f", 100*(first/total))
	}
}

// display version info.
func HandleVersion(displayVersion bool) {
	if displayVersion {
		fmt.Println(Version)
		os.Exit(0)
	}
}

// set memprofile
func HandleMemProfile(memprofile string) (file *os.File, err error) {
	if memprofile != "" {
		var err error
		memFile, err := os.Create(memprofile)
		if err != nil {
			glog.Warningf("Start write heap profile....")
			return nil, err
		} else {
			glog.Infoln("Start write heap profile....")
			pprof.WriteHeapProfile(memFile)
		}

		return memFile, nil
	}

	return nil, nil
}

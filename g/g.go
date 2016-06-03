package g

import (
	"runtime"
)

const (
	VERSION = "0.2.0"
	Commit  = ""
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

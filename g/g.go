package g

import (
	"runtime"
)

const (
	VERSION = "0.2.0"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

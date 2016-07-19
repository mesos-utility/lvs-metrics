package g

import (
	"runtime"
)

// version will be populated by the Makefile, read from
// VERSION file of the source code.
var Version = ""

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

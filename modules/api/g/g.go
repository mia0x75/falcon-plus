package g

import (
	"runtime"
)

// change log:
const (
	VERSION = "0.0.1"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

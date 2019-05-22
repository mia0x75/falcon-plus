package g

import (
	nproc "github.com/toolkits/proc"
)

// 变量定义
var (
	RecvDataTrace  = nproc.NewDataTrace("RecvDataTrace", 5)
	RecvDataFilter = nproc.NewDataFilter("RecvDataFilter", 5)
)

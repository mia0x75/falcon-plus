package receiver

import (
	"github.com/open-falcon/falcon-plus/modules/transfer/receiver/rpc"
	"github.com/open-falcon/falcon-plus/modules/transfer/receiver/socket"
)

// Start 启动服务
func Start() {
	rpc.Start()
	socket.Start()
}

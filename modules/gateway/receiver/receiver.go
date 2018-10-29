package receiver

import (
	"github.com/open-falcon/falcon-plus/modules/gateway/receiver/rpc"
	"github.com/open-falcon/falcon-plus/modules/gateway/receiver/socket"
)

func Start() {
	rpc.Start()
	socket.Start()
}

package rpc

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/transfer/g"
)

// Start 启动服务
func Start() {
	go start()
}

func start() {
	if !g.Config().RPC.Enabled {
		return
	}

	rpc.Register(new(Transfer))

	addr := g.Config().RPC.Listen
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatalf("[F] rpc.Start error, net.ResolveTCPAddr fail, %s", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("[F] rpc.Start error, listen %s fail, %s", addr, err)
	} else {
		log.Infof("[I] rpc listening %s", addr)
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Errorf("[E] listener.Accept occur error: %v", err)
			continue
		}
		conn.SetKeepAlive(true)
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

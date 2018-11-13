package rpc

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/api/g"
)

func Start() {
	go start()
}

func start() {
	if !g.Config().Rpc.Enabled {
		return
	}

	rpc.Register(new(Graph))

	addr := g.Config().Rpc.Listen
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

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Infof("[I] listener accept fail: %v", err)
				time.Sleep(time.Duration(100) * time.Millisecond)
				continue
			}
			go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}()
}

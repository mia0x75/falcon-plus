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
		log.Fatalf("rpc.Start error, net.ResolveTCPAddr fail, %s", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("rpc.Start error, listen %s fail, %s", addr, err)
	} else {
		log.Printf("rpc listening %s", addr)
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("listener accept fail:", err)
				time.Sleep(time.Duration(100) * time.Millisecond)
				continue
			}
			go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}()
}

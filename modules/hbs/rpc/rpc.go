package rpc

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/hbs/g"
)

func Start() {
	go start()
}

func start() {
	addr := g.Config().Listen

	rpc.Register(new(Agent))
	rpc.Register(new(Hbs))

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatalf("rpc.Start error, net.ResolveTCPAddr fail, %s", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("rpc.Start error, listen %s fail, %s", addr, err)
	} else {
		log.Println("rpc listening", addr)
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

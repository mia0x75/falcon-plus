package rpc

import (
	"net"
	"net/rpc"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/judge/g"
)

func Start() {
	go start()
}

func start() {
	if !g.Config().Rpc.Enabled {
		return
	}

	rpc.Register(new(Judge))

	addr := g.Config().Rpc.Listen
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
				log.Printf("listener.Accept occur error: %s", err)
				continue
			}
			go rpc.ServeConn(conn)
		}
	}()
}

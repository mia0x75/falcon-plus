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
	if !g.Config().Rpc.Enabled {
		return
	}

	addr := g.Config().Rpc.Listen
	server := rpc.NewServer()
	server.Register(new(GraphRpc))
	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatalln("listen error:", e)
	} else {
		log.Println("listening", addr)
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Println("listener accept fail:", err)
				time.Sleep(time.Duration(100) * time.Millisecond)
				continue
			}
			go server.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}()
}

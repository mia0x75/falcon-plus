package rpc

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/hbs/g"
)

type Hbs int
type Agent int

func Start() {
	addr := g.Config().Listen

	server := rpc.NewServer()
	server.Register(new(Agent))
	server.Register(new(Hbs))

	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatalf("rpc listen error:%v", e)
	} else {
		log.Printf("rpc listening %s", addr)
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

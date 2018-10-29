package socket

import (
	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/transfer/g"
)

func Start() {
	go start()
}

func start() {
	if !g.Config().Socket.Enabled {
		return
	}

	addr := g.Config().Socket.Listen
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatalf("socket.Start error, net.ResolveTCPAddr fail, %s", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("socket.Start error, listen %s fail, %s", addr, err)
	} else {
		log.Println("socket listening", addr)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept occur error:", err)
			continue
		}

		go socketTelnetHandle(conn)
	}
}

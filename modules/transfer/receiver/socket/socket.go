package socket

import (
	"net"

	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/transfer/g"
)

// Start 启动服务
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
		log.Fatalf("[F] socket.Start error, net.ResolveTCPAddr fail, %s", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("[F] socket.Start error, listen %s fail, %s", addr, err)
	} else {
		log.Infof("[I] socket listening %s", addr)
	}

	defer listener.Close()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Errorf("[E] listener.Accept occur error: %v", err)
			continue
		}
		conn.SetKeepAlive(true)
		go socketTelnetHandle(conn)
	}
}

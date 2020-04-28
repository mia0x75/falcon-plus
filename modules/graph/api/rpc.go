package api

import (
	"container/list"
	"net"
	"net/rpc"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/open-falcon/falcon-plus/modules/graph/g"
)

type conn_list struct {
	sync.RWMutex
	list *list.List
}

func (l *conn_list) insert(c net.Conn) *list.Element {
	l.Lock()
	defer l.Unlock()
	return l.list.PushBack(c)
}

func (l *conn_list) remove(e *list.Element) net.Conn {
	l.Lock()
	defer l.Unlock()
	return l.list.Remove(e).(net.Conn)
}

var Close_chan, Close_done_chan chan int
var connects conn_list

func init() {
	Close_chan = make(chan int, 1)
	Close_done_chan = make(chan int, 1)
	connects = conn_list{list: list.New()}
}

func Start() {
	go start()
}

func start() {
	if !g.Config().RPC.Enabled {
		log.Info("[I] rpc.Start warning, not enabled")
		return
	}

	rpc.Register(new(Graph))

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

	go func() {
		var tempDelay time.Duration // how long to sleep on accept failure
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			tempDelay = 0
			conn.SetKeepAlive(true)
			go func() {
				e := connects.insert(conn)
				defer connects.remove(e)
				rpc.ServeConn(conn)
			}()
		}
	}()

	select {
	case <-Close_chan:
		log.Info("[I] rpc, recv sigout and exiting...")
		listener.Close()
		Close_done_chan <- 1

		connects.Lock()
		for e := connects.list.Front(); e != nil; e = e.Next() {
			e.Value.(net.Conn).Close()
		}
		connects.Unlock()

		return
	}
}

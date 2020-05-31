package pool

import (
	"fmt"
	"net"
	"time"

	connp "github.com/toolkits/conn_pool"
)

// TSDBClient TODO:
type TSDBClient struct {
	cli  *struct{ net.Conn }
	name string
}

// Name TODO:
func (t TSDBClient) Name() string {
	return t.name
}

// Closed TODO:
func (t TSDBClient) Closed() bool {
	return t.cli.Conn == nil
}

// Close TODO:
func (t TSDBClient) Close() error {
	if t.cli != nil {
		err := t.cli.Close()
		t.cli.Conn = nil
		return err
	}
	return nil
}

func newTSDBConnPool(address string, maxConns int, maxIdle int, connTimeout int) *connp.ConnPool {
	pool := connp.NewConnPool("tsdb", address, int32(maxConns), int32(maxIdle))

	pool.New = func(name string) (connp.NConn, error) {
		_, err := net.ResolveTCPAddr("tcp", address)
		if err != nil {
			return nil, err
		}

		conn, err := net.DialTimeout("tcp", address, time.Duration(connTimeout)*time.Millisecond)
		if err != nil {
			return nil, err
		}

		return TSDBClient{
			cli:  &struct{ net.Conn }{conn},
			name: name,
		}, nil
	}

	return pool
}

// TSDBConnPoolHelper TODO:
type TSDBConnPoolHelper struct {
	p           *connp.ConnPool
	maxConns    int
	maxIdle     int
	connTimeout int
	callTimeout int
	address     string
}

// NewTSDBConnPoolHelper TODO:
func NewTSDBConnPoolHelper(address string, maxConns, maxIdle, connTimeout, callTimeout int) *TSDBConnPoolHelper {
	return &TSDBConnPoolHelper{
		p:           newTSDBConnPool(address, maxConns, maxIdle, connTimeout),
		maxConns:    maxConns,
		maxIdle:     maxIdle,
		connTimeout: connTimeout,
		callTimeout: callTimeout,
		address:     address,
	}
}

// Send TODO:
func (t *TSDBConnPoolHelper) Send(data []byte) (err error) {
	conn, err := t.p.Fetch()
	if err != nil {
		return fmt.Errorf("get connection fail: err %v. proc: %s", err, t.p.Proc())
	}

	cli := conn.(TSDBClient).cli

	done := make(chan error, 1)
	go func() {
		_, err = cli.Write(data)
		done <- err
	}()

	select {
	case <-time.After(time.Duration(t.callTimeout) * time.Millisecond):
		t.p.ForceClose(conn)
		return fmt.Errorf("%s, call timeout", t.address)
	case err = <-done:
		if err != nil {
			t.p.ForceClose(conn)
			err = fmt.Errorf("%s, call failed, err %v. proc: %s", t.address, err, t.p.Proc())
		} else {
			t.p.Release(conn)
		}
		return err
	}
}

// Destroy TODO:
func (t *TSDBConnPoolHelper) Destroy() {
	if t.p != nil {
		t.p.Destroy()
	}
}

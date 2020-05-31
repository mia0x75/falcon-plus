package g

import (
	"bytes"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	cm "github.com/open-falcon/falcon-plus/common/model"
)

var Root string

// InitRootDir
func InitRootDir() {
	defer func() {
		log.Debugf("[D] Root dir: %s", Root)
	}()

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("[F] getwd fail: %v", err)
	}
	if root := lookupRootDir(pwd); root != "" {
		Root = root
		return
	}
	file, _ := exec.LookPath(os.Args[0])
	execDir, _ := filepath.Abs(file)
	if root := lookupRootDir(path.Dir(execDir)); root != "" {
		Root = root
		return
	}
	Root = pwd
}

// lookupRootDir path目录下找 public，判断是否存在，不存在上级目录找
func lookupRootDir(dir string) string {
	index := filepath.Join("public", "index.html")
	public := filepath.Join(dir, index)
	if _, err := os.Stat(public); os.IsNotExist(err) {
		for i := 0; i < 2; i++ {
			parent := path.Dir(dir)
			public = filepath.Join(parent, index)
			if _, err := os.Stat(public); err == nil {
				return parent
			}
			dir = parent
		}
	} else if err == nil {
		return dir
	}
	return ""
}

// LocalIP TODO:
var LocalIP string

// InitLocalIP TODO:
func InitLocalIP() {
	for _, addr := range Config().Heartbeat.Addrs {
		conn, err := net.DialTimeout("tcp", addr, time.Second*10)
		if err != nil {
			log.Errorf("[E] connect to heartbeat server %s failed", addr)
		} else {
			defer conn.Close()
			LocalIP = strings.Split(conn.LocalAddr().String(), ":")[0]
			break
		}
	}
}

// TODO:
var (
	HbsClient *SingleConnRPCClient
)

// InitRPCClients TODO:
func InitRPCClients() {
	if len(Config().Heartbeat.Addrs) > 0 {
		HbsClient = &SingleConnRPCClient{
			RPCServers: Config().Heartbeat.Addrs,
			Timeout:    time.Duration(Config().Heartbeat.Timeout) * time.Millisecond,
		}
	} else {
		// TODO: panic
	}
}

// SendToTransfer TODO:
func SendToTransfer(metrics []*cm.MetricValue) {
	if len(metrics) == 0 {
		return
	}

	dt := Config().DefaultTags
	if len(dt) > 0 {
		var buf bytes.Buffer
		list := []string{}
		for k, v := range dt {
			buf.Reset()
			buf.WriteString(k)
			buf.WriteString("=")
			buf.WriteString(v)
			list = append(list, buf.String())
		}
		defaultTags := strings.Join(list, ",")

		for i, x := range metrics {
			buf.Reset()
			if x.Tags == "" {
				metrics[i].Tags = defaultTags
			} else {
				buf.WriteString(metrics[i].Tags)
				buf.WriteString(",")
				buf.WriteString(defaultTags)
				metrics[i].Tags = buf.String()
			}
		}
	}
	for _, m := range metrics {
		log.Debugf("[D] => Metric %v", m)
	}

	log.Debugf("[D] => <Total=%d> %v", len(metrics), metrics[0])

	var resp cm.TransferResponse
	SendMetrics(metrics, &resp)

	log.Debugf("[D] <= %v", &resp)
}

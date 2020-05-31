package g

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"unsafe"

	pfcg "github.com/mia0x75/gopfc/g"
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/file"
)

// File TODO:
type File struct {
	Filename string
	Body     []byte
}

// HTTPConfig HTTP配置
type HTTPConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

// RPCConfig RPC配置
type RPCConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

// RRDConfig RRD配置
type RRDConfig struct {
	Storage string `json:"storage"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Addr           string `json:"addr"`
	MaxIdle        int    `json:"max_idle"`
	MaxConnections int    `json:"max_connections"`
	WaitTimeout    int    `json:"wait_timeout"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `json:"level"`
}

// GlobalConfig Graph模块配置
type GlobalConfig struct {
	Log            *LogConfig      `json:"log"`
	Pid            string          `json:"pid"`
	HTTP           *HTTPConfig     `json:"http"`
	RPC            *RPCConfig      `json:"rpc"`
	RRD            *RRDConfig      `json:"rrd"`
	Database       *DatabaseConfig `json:"database"`
	ExecuteTimeout int32           `json:"execute_timeout"`
	IOWorkerNum    int             `json:"io_workers"`
	FirstBytesSize int
	Migrate        struct {
		Concurrency int               `json:"concurrency"` // number of multiple worker per node
		Enabled     bool              `json:"enabled"`
		Replicas    int               `json:"replicas"`
		Cluster     map[string]string `json:"cluster"`
	} `json:"migrate"`
	PerfCounter *pfcg.GlobalConfig `json:"pfc"`
}

// 变量定义
var (
	ConfigFile string
	ptr        unsafe.Pointer
)

// Config 获取Graph模块配置
func Config() *GlobalConfig {
	return (*GlobalConfig)(atomic.LoadPointer(&ptr))
}

// ParseConfig 解析配置文件
func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatal("[F] config file not specified: use -c $filename")
	}

	if !file.IsExist(cfg) {
		log.Fatalf("[F] config file specified not found: %v", cfg)
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalf("[F] read config file %s error: %v", cfg, err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalf("[F] parse config file %s error: %v", cfg, err)
	}

	if c.Migrate.Enabled && len(c.Migrate.Cluster) == 0 {
		c.Migrate.Enabled = false
	}

	// 确保ioWorkerNum是2^N
	if c.IOWorkerNum == 0 || (c.IOWorkerNum&(c.IOWorkerNum-1) != 0) {
		log.Fatalf("[F] IOWorkerNum must be 2^N, current IOWorkerNum is %v", c.IOWorkerNum)
	}

	// 需要md5的前多少位参与ioWorker的分片计算
	c.FirstBytesSize = len(strconv.FormatInt(int64(c.IOWorkerNum), 16))

	if c.PerfCounter != nil {
		c.PerfCounter.Debug = c.Log.Level == "debug"
		port := strings.Split(c.HTTP.Listen, ":")[1]
		if c.PerfCounter.Tags == "" {
			c.PerfCounter.Tags = fmt.Sprintf("port=%s", port)
		} else {
			c.PerfCounter.Tags = fmt.Sprintf("%s,port=%s", c.PerfCounter.Tags, port)
		}
	}

	// set config
	atomic.StorePointer(&ptr, unsafe.Pointer(&c))

	log.Debugf("[D] read config file: %s successfully", cfg)
}

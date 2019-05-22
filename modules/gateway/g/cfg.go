package g

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	pfcg "github.com/mia0x75/gopfc/g"
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/file"
)

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

// SocketConfig TCP配置
type SocketConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
	Timeout int32  `json:"timeout"`
}

// TransferConfig Transfer配置
type TransferConfig struct {
	Enabled        bool              `json:"enabled"`
	Batch          int32             `json:"batch"`
	Retry          int32             `json:"retry"`
	ConnectTimeout int32             `json:"connect_timeout"`
	ExecuteTimeout int32             `json:"execute_timeout"`
	MaxConnections int32             `json:"max_connections"`
	MaxIdle        int32             `json:"max_idle"`
	Cluster        map[string]string `json:"cluster"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `json:"level"`
}

// GlobalConfig 网关配置
type GlobalConfig struct {
	Log         *LogConfig         `json:"log"`
	HTTP        *HTTPConfig        `json:"http"`
	RPC         *RPCConfig         `json:"rpc"`
	Socket      *SocketConfig      `json:"socket"`
	Transfer    *TransferConfig    `json:"transfer"`
	Host        string             `json:"host"`
	PerfCounter *pfcg.GlobalConfig `json:"pfc"`
}

// 变量定义
var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

// Config 获取网关配置
func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

// ParseConfig 解析配置文件
func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatal("[F] use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalf("[F] config file: %s is not existent. maybe you need `mv cfg.example.json cfg.json`", cfg)
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalf("[F] read config file: %s fail: %v", cfg, err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalf("[F] parse config file: %s fail: %v", cfg, err)
	}

	// 配置文件正确性 校验, 不合法则直接 Exit(1)
	// TODO

	configLock.Lock()
	defer configLock.Unlock()

	if c.PerfCounter != nil {
		c.PerfCounter.Debug = c.Log.Level == "debug"
		port := strings.Split(c.HTTP.Listen, ":")[1]
		if c.PerfCounter.Tags == "" {
			c.PerfCounter.Tags = fmt.Sprintf("port=%s", port)
		} else {
			c.PerfCounter.Tags = fmt.Sprintf("%s,port=%s", c.PerfCounter.Tags, port)
		}
	}
	config = &c

	log.Debugf("[D] read config file: %s successfully", cfg)
}

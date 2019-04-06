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

type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type RpcConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type SocketConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
	Timeout int32  `json:"timeout"`
}

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

type LogConfig struct {
	Level string `json:"level"`
}

type GlobalConfig struct {
	Log         *LogConfig         `json:"log"`
	Http        *HttpConfig        `json:"http"`
	Rpc         *RpcConfig         `json:"rpc"`
	Socket      *SocketConfig      `json:"socket"`
	Transfer    *TransferConfig    `json:"transfer"`
	Host        string             `json:"host"`
	PerfCounter *pfcg.GlobalConfig `json:"pfc"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

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
		port := strings.Split(c.Http.Listen, ":")[1]
		if c.PerfCounter.Tags == "" {
			c.PerfCounter.Tags = fmt.Sprintf("port=%s", port)
		} else {
			c.PerfCounter.Tags = fmt.Sprintf("%s,port=%s", c.PerfCounter.Tags, port)
		}
	}
	config = &c

	log.Debugf("[D] read config file: %s successfully", cfg)
}

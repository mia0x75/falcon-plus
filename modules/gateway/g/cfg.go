package g

import (
	"encoding/json"
	"sync"

	log "github.com/Sirupsen/logrus"
	pfcg "github.com/mia0x75/gopfc/g"
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
	Enabled     bool              `json:"enabled"`
	Batch       int32             `json:"batch"`
	Retry       int32             `json:"retry"`
	ConnTimeout int32             `json:"connect_timeout"`
	CallTimeout int32             `json:"execute_timeout"`
	MaxConns    int32             `json:"max_connections"`
	MaxIdle     int32             `json:"max_idle"`
	Cluster     map[string]string `json:"cluster"`
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
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	// 配置文件正确性 校验, 不合法则直接 Exit(1)
	// TODO

	configLock.Lock()
	defer configLock.Unlock()
	c.PerfCounter.Debug = c.Log.Level == "debug"
	config = &c

	log.Debugln("read config file:", cfg, "successfully")
}

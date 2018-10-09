package g

import (
	"encoding/json"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/toolkits/file"
)

type HttpConfig struct {
	Enable bool   `json:"enable"`
	Listen string `json:"listen"`
}

type IndexConfig struct {
	Enable         bool              `json:"enable"`
	Addr           string            `json:"addr"`
	MaxIdle        int               `json:"max_idle"`
	MaxConnections int               `json:"max_connections"`
	AutoDelete     bool              `json:"auto_delete"`
	Cluster        map[string]string `json:"cluster"`
}

type CollectorConfig struct {
	Enable    bool     `json:"enable"`
	DestUrl   string   `json:"agent"`
	SrcUrlFmt string   `json:"url_pattern"`
	Cluster   []string `json:"cluster"`
}

type PluginConfig struct {
	ApiUrlFmt      string `json:"url_pattern"`
	Interval       int32  `json:"interval"`
	Concurrent     int32  `json:"concurrent"`
	ConnectTimeout int32  `json:"connect_timeout"`
	RequestTimeout int32  `json:"request_timeout"`
}

type CleanerConfig struct {
	Interval int32 `json:"interval"`
}

type AgentConfig struct {
	Enable  bool           `json:"enable"`
	Dsn     string         `json:"dsn"`
	MaxIdle int32          `json:"max_idle"`
	Plugin  *PluginConfig  `json:"plugin"`
	Cleaner *CleanerConfig `json:"cleaner"`
}

type LogConfig struct {
	Level string `json:"level"`
}

type GlobalConfig struct {
	Log       *LogConfig       `json:"log"`
	Http      *HttpConfig      `json:"http"`
	Index     *IndexConfig     `json:"index"`
	Collector *CollectorConfig `json:"collector"`
	Agent     *AgentConfig     `json:"agent"`
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

	configLock.Lock()
	defer configLock.Unlock()
	config = &c

	log.Debugln("read config file:", cfg, "successfully")
}

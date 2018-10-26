package g

import (
	"encoding/json"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/toolkits/file"
)

type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type PlusAPIConfig struct {
	Addr           string `json:"addr"`
	Token          string `json:"token"`
	ConnectTimeout int32  `json:"connect_timeout"`
	RequestTimeout int32  `json:"request_timeout"`
}

type DatabaseConfig struct {
	Enabled        bool   `json:"enabled"`
	Addr           string `json:"addr"`
	MaxIdle        int    `json:"max_idle"`
	MaxConnections int    `json:"max_connections"`
	WaitTimeout    int    `json:"wait_timeout"`
}

type CollectorConfig struct {
	Enabled    bool  `json:"enabled"`
	Batch      int32 `json:"batch"`
	Concurrent int32 `json:"concurrent"`
}

type TransferConfig struct {
	Enabled        bool   `json:"enabled"`
	Addr           string `json:"addr"`
	ConnectTimeout int32  `json:"connect_timeout"`
	RequestTimeout int32  `json:"request_timeout"`
	Batch          int32  `json:"batch"`
}

type LogConfig struct {
	Level string `json:"level"`
}

type GlobalConfig struct {
	Log       *LogConfig       `json:"log"`
	Http      *HttpConfig      `json:"http"`
	Api       *PlusAPIConfig   `json:"api"`
	Database  *DatabaseConfig  `json:"database"`
	Collector *CollectorConfig `json:"collector"`
	Transfer  *TransferConfig  `json:"transfer"`
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
		log.Fatalln("config file:", cfg, "is not existent")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("g.ParseConfig error, parse config file", cfg, "fail,", err)
	}

	configLock.Lock()
	defer configLock.Unlock()
	config = &c

	log.Debugln("read config file:", cfg, "successfully")
}

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

type DatabaseConfig struct {
	Addr     string `json:"addr"`
	Idle     int    `json:"idle"`
	Ids      []int  `json:"ids"`
	Interval int64  `json:"interval"`
}

type ApiConfig struct {
	ConnectTimeout int32  `json:"connect_timeout"`
	RequestTimeout int32  `json:"request_timeout"`
	Api            string `json:"api"`
	Token          string `json:"token"`
	Agent          string `json:"agent"`
}

type LogConfig struct {
	Level string `json:"level"`
}

type GlobalConfig struct {
	Log      *LogConfig      `json:"log"`
	Http     *HttpConfig     `json:"http"`
	Database *DatabaseConfig `json:"database"`
	Api      *ApiConfig      `json:"api"`
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
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	configLock.Lock()
	defer configLock.Unlock()
	config = &c
	log.Println("read config file:", cfg, "successfully")
}

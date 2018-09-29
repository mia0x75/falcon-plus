package g

import (
	"encoding/json"
	"log"
	"sync"

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

type NdConfig struct {
	Enabled bool   `json:"enabled"`
	Dsn     string `json:"dsn"`
	MaxIdle int32  `json:"max_idle"`
}

type CollectorConfig struct {
	Enabled    bool  `json:"enabled"`
	Batch      int32 `json:"batch"`
	Concurrent int32 `json:"concurrent"`
}

type SenderConfig struct {
	Enabled        bool   `json:"enabled"`
	TransferAddr   string `json:"transfer_addr"`
	ConnectTimeout int32  `json:"connect_timeout"`
	RequestTimeout int32  `json:"request_timeout"`
	Batch          int32  `json:"batch"`
}

type GlobalConfig struct {
	Debug     bool             `json:"debug"`
	Http      *HttpConfig      `json:"http"`
	Api       *PlusAPIConfig   `json:"api"`
	Config    *NdConfig        `json:"config"`
	Collector *CollectorConfig `json:"collector"`
	Sender    *SenderConfig    `json:"sender"`
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

	log.Println("g.ParseConfig ok, file ", cfg)
}

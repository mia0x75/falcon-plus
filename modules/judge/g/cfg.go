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

type RpcConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type HbsConfig struct {
	Servers  []string `json:"servers"`
	Timeout  int64    `json:"timeout"`
	Interval int64    `json:"interval"`
}

type RedisConfig struct {
	Addr         string `json:"addr"`
	MaxIdle      int    `json:"max_idle"`
	ConnTimeout  int    `json:"connect_timeout"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
	WaitTimeout  int    `json:"wait_timeout"`
}

type AlarmConfig struct {
	MinInterval  int64  `json:"min_interval"`
	QueuePattern string `json:"queue_pattern"`
}

type LogConfig struct {
	Level string `json:"level"`
}

type GlobalConfig struct {
	Log       *LogConfig   `json:"log"`
	DebugHost string       `json:"debug_host"`
	Remain    int          `json:"remain"`
	Http      *HttpConfig  `json:"http"`
	Rpc       *RpcConfig   `json:"rpc"`
	Hbs       *HbsConfig   `json:"hbs"`
	Alarm     *AlarmConfig `json:"alarm"`
	Redis     *RedisConfig `json:"redis"`
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
		log.Fatalf("[F] config file: %s is not existent", cfg)
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

	configLock.Lock()
	defer configLock.Unlock()

	config = &c

	log.Debugf("[D] read config file: %s successfully", cfg)
}

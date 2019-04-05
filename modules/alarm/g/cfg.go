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

type RedisConfig struct {
	Addr        string `json:"addr"`
	MaxIdle     int    `json:"max_idle"`
	WaitTimeout int    `json:"wait_timeout"`
}

type QueueConfig struct {
	HighQueues    []string       `json:"high_queues"`
	LowQueues     []string       `json:"low_queues"`
	InstantQueues *ChannelConfig `json:"instant_queues"`
	LatentQueues  *ChannelConfig `json:"latent_queues"`
}

type ChannelConfig struct {
	IMQueue   string `json:"im"`
	SmsQueue  string `json:"sms"`
	MailQueue string `json:"mail"`
}

type ApiConfig struct {
	Sms       string `json:"sms"`
	Mail      string `json:"mail"`
	Dashboard string `json:"dashboard"`
	Api       string `json:"api"`
	Token     string `json:"token"`
	IM        string `json:"im"`
}

type DatabaseConfig struct {
	Addr           string `json:"addr"`
	MaxIdle        int    `json:"max_idle"`
	MaxConnections int    `json:"max_connections"`
	WaitTimeout    int    `json:"wait_timeout"`
}

type WorkerConfig struct {
	IM   int `json:"im"`
	Sms  int `json:"sms"`
	Mail int `json:"mail"`
}

type HousekeeperConfig struct {
	EventRetentionDays int `json:"event_retention_days"`
	EventDeleteBatch   int `json:"event_delete_batch"`
}

type LogConfig struct {
	Level string `json:"level"`
}

type GlobalConfig struct {
	Log         *LogConfig         `json:"log"`
	Database    *DatabaseConfig    `json:"database"`
	Http        *HttpConfig        `json:"http"`
	Redis       *RedisConfig       `json:"redis"`
	Queue       *QueueConfig       `json:"queue"`
	Api         *ApiConfig         `json:"api"`
	Worker      *WorkerConfig      `json:"worker"`
	Housekeeper *HousekeeperConfig `json:"housekeeper"`
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

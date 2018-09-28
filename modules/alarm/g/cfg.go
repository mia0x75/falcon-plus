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

type RedisConfig struct {
	Addr          string   `json:"addr"`
	MaxIdle       int      `json:"max_idle"`
	HighQueues    []string `json:"high_queues"`
	LowQueues     []string `json:"low_queues"`
	UserIMQueue   string   `json:"user_im_queue"`
	UserSmsQueue  string   `json:"user_sms_queue"`
	UserMailQueue string   `json:"user_mail_queue"`
}

type ApiConfig struct {
	Sms          string `json:"sms"`
	Mail         string `json:"mail"`
	Dashboard    string `json:"dashboard"`
	PlusApi      string `json:"plus_api"`
	PlusApiToken string `json:"plus_api_token"`
	IM           string `json:"im"`
}

type PortalConfig struct {
	Addr string `json:"addr"`
	Idle int    `json:"idle"`
	Max  int    `json:"max"`
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

type GlobalConfig struct {
	LogLevel    string             `json:"log_level"`
	Portal      *PortalConfig      `json:"portal"`
	Http        *HttpConfig        `json:"http"`
	Redis       *RedisConfig       `json:"redis"`
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

package g

import (
	"encoding/json"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/toolkits/file"
)

// HTTPConfig HTTP配置
type HTTPConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

// RedisConfig 缓存服务器配置
type RedisConfig struct {
	Addr        string `json:"addr"`
	MaxIdle     int    `json:"max_idle"`
	WaitTimeout int    `json:"wait_timeout"`
}

// QueueConfig 消息队列配置
type QueueConfig struct {
	HighQueues    []string       `json:"high_queues"`
	LowQueues     []string       `json:"low_queues"`
	InstantQueues *ChannelConfig `json:"instant_queues"`
	LatentQueues  *ChannelConfig `json:"latent_queues"`
}

// ChannelConfig 消息队列配置
type ChannelConfig struct {
	IMQueue   string `json:"im"`
	SmsQueue  string `json:"sms"`
	MailQueue string `json:"mail"`
}

// APIConfig API配置
type APIConfig struct {
	Sms       string `json:"sms"`
	Mail      string `json:"mail"`
	Dashboard string `json:"dashboard"`
	API       string `json:"api"`
	Token     string `json:"token"`
	IM        string `json:"im"`
}

// DatabaseConfig 数据库连接配置
type DatabaseConfig struct {
	Addr           string `json:"addr"`
	MaxIdle        int    `json:"max_idle"`
	MaxConnections int    `json:"max_connections"`
	WaitTimeout    int    `json:"wait_timeout"`
}

// WorkerConfig Worker配置
type WorkerConfig struct {
	IM   int `json:"im"`
	Sms  int `json:"sms"`
	Mail int `json:"mail"`
}

// HousekeeperConfig 数据清理配置
type HousekeeperConfig struct {
	EventRetentionDays int `json:"event_retention_days"`
	EventDeleteBatch   int `json:"event_delete_batch"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `json:"level"`
}

// GlobalConfig 告警配置
type GlobalConfig struct {
	Log         *LogConfig         `json:"log"`
	Database    *DatabaseConfig    `json:"database"`
	HTTP        *HTTPConfig        `json:"http"`
	Redis       *RedisConfig       `json:"redis"`
	Queue       *QueueConfig       `json:"queue"`
	API         *APIConfig         `json:"api"`
	Worker      *WorkerConfig      `json:"worker"`
	Housekeeper *HousekeeperConfig `json:"housekeeper"`
}

// 变量定义
var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

// Config 获取告警模块配置
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

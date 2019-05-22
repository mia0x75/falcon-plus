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

// DatabaseConfig 数据库连接配置
type DatabaseConfig struct {
	Addr           string `json:"addr"`
	MaxIdle        int    `json:"max_idle"`
	MaxConnections int    `json:"max_connections"`
	WaitTimeout    int    `json:"wait_timeout"`
	Ids            []int  `json:"ids"`
	Interval       int64  `json:"interval"`
}

// APIConfig API配置
type APIConfig struct {
	ConnectTimeout int32  `json:"connect_timeout"`
	RequestTimeout int32  `json:"request_timeout"`
	API            string `json:"api"`
	Token          string `json:"token"`
	Agent          string `json:"agent"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `json:"level"`
}

// GlobalConfig 聚合组件配置
type GlobalConfig struct {
	Log      *LogConfig      `json:"log"`
	HTTP     *HTTPConfig     `json:"http"`
	Database *DatabaseConfig `json:"database"`
	API      *APIConfig      `json:"api"`
}

// 变量定义
var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

// Config 获取聚合配置
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

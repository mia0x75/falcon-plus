// change log:
// 1.0.7: code refactor for open source
// 1.0.8: bugfix loop init cache
// 1.0.9: update host table anyway
// 1.1.0: remove Checksum when query plugins
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

// LogConfig 日志配置
type LogConfig struct {
	Level string `json:"level"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Addr           string `json:"addr"`
	MaxIdle        int    `json:"max_idle"`
	MaxConnections int    `json:"max_connections"`
	WaitTimeout    int    `json:"wait_timeout"`
}

// GlobalConfig HBS模块配置
type GlobalConfig struct {
	Log            *LogConfig      `json:"log"`
	Hosts          string          `json:"hosts"`
	Database       *DatabaseConfig `json:"database"`
	MaxConnections int             `json:"max_connections"`
	MaxIdle        int             `json:"max_idle"`
	Listen         string          `json:"listen"`
	Trustable      []string        `json:"trustable"`
	HTTP           *HTTPConfig     `json:"http"`
}

// 变量定义
var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

// Config 获取HBS模块配置
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

package g

import (
	"encoding/json"
	"fmt"
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

// GlobalConfig Updater模块配置
type GlobalConfig struct {
	Log          *LogConfig  `json:"log"`
	Hostname     string      `json:"hostname"`
	DesiredAgent string      `json:"desired_agent"`
	Server       string      `json:"server"`
	Interval     int         `json:"interval"`
	HTTP         *HTTPConfig `json:"http"`
}

// 变量定义
var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

// Config 获取Updater模块的配置
func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

// ParseConfig 解析配置文件
func ParseConfig(cfg string) error {
	if cfg == "" {
		return fmt.Errorf("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		return fmt.Errorf("config file %s is nonexistent", cfg)
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		return fmt.Errorf("read config file %s fail %s", cfg, err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		return fmt.Errorf("parse config file %s fail %s", cfg, err)
	}

	configLock.Lock()
	defer configLock.Unlock()

	config = &c

	log.Debugf("[D] read config file: %s successfully", cfg)
	return nil
}

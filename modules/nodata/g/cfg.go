// changelog:
// 0.0.1 init project
// 0.0.2 make mock item.Ts one step after now(), rm sending log, add flood in proc
// 0.0.3 mv common to falcon.common, simplify nodata's codes, mv cfgcenter to nodata
// 0.0.4 fix bug of nil response on collecting from query
// 0.0.5 collect items concurrently from query
// 0.0.6 clear send buffer when blocking
// 0.0.7 use gauss distribution to get threshold, sync judge and sender, fix bug of collector's cache
// 0.0.8 simplify project
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

// PlusAPIConfig PlusAPI配置
type PlusAPIConfig struct {
	Addr           string `json:"addr"`
	Token          string `json:"token"`
	ConnectTimeout int32  `json:"connect_timeout"`
	RequestTimeout int32  `json:"request_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Enabled        bool   `json:"enabled"`
	Addr           string `json:"addr"`
	MaxIdle        int    `json:"max_idle"`
	MaxConnections int    `json:"max_connections"`
	WaitTimeout    int    `json:"wait_timeout"`
}

// CollectorConfig Collector配置
type CollectorConfig struct {
	Enabled    bool  `json:"enabled"`
	Batch      int32 `json:"batch"`
	Concurrent int32 `json:"concurrent"`
}

// TransferConfig Transfer配置
type TransferConfig struct {
	Enabled        bool   `json:"enabled"`
	Addr           string `json:"addr"`
	ConnectTimeout int32  `json:"connect_timeout"`
	RequestTimeout int32  `json:"request_timeout"`
	Batch          int32  `json:"batch"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `json:"level"`
}

// GlobalConfig Nodata模块配置
type GlobalConfig struct {
	Log       *LogConfig       `json:"log"`
	HTTP      *HTTPConfig      `json:"http"`
	API       *PlusAPIConfig   `json:"api"`
	Database  *DatabaseConfig  `json:"database"`
	Collector *CollectorConfig `json:"collector"`
	Transfer  *TransferConfig  `json:"transfer"`
}

// 变量定义
var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

// Config 获取Nodata模块的配置
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
		log.Fatalf("[F] g.ParseConfig error, parse config file %s fail: %v", cfg, err)
	}

	configLock.Lock()
	defer configLock.Unlock()
	config = &c

	log.Debugf("[D] read config file: %s successfully", cfg)
}

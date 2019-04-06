// changelog:
// 0.0.1: init project
// 0.0.3: add readme, add gitversion, modify proc, add config reload
// 0.0.4: make collector configurable, add monitor cron, adjust index db
// Changes: send turning-ok only after alarm happens, add conn timeout for http
//			maybe fix bug of 'too many open files', rollback to central lib
// 0.0.5: move self.monitor to anteye
// 0.0.6: make index update configurable, use global time formater
// 0.0.7: fix bug of index_update_all
// 0.0.8: add agents' house_keeper, use relative paths in 'import'
// 0.0.9: gen exporter.alive, use common module, use absolute paths in import
// 0.0.10: rm monitor, add controller for index cleaner
package g

import (
	"encoding/json"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/toolkits/file"
)

type AlarmConfig struct {
	Enabled bool   `json:"enabled"`
	Url     string `json:"url"`
}

type IndexConfig struct {
	Enabled        bool              `json:"enabled"`
	Addr           string            `json:"addr"`
	MaxIdle        int               `json:"max_idle"`
	MaxConnections int               `json:"max_connections"`
	WaitTimeout    int               `json:"wait_timeout"`
	AutoDelete     bool              `json:"auto_delete"`
	Cluster        map[string]string `json:"cluster"`
}

type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type MonitorConfig struct {
	Enabled bool         `json:"enabled"`
	Alarm   *AlarmConfig `json:"alarm"`
	Pattern string       `json:"pattern"`
	Hosts   *HostsConfig `json:"hosts"`
}

type HostsConfig struct {
	Agents  []string          `json:"agents"`
	Modules map[string]string `json:"modules"`
}

type CollectorConfig struct {
	Enabled bool     `json:"enabled"`
	Agent   string   `json:"agent"`
	Pattern string   `json:"pattern"`
	Cluster []string `json:"cluster"`
}

type PluginConfig struct {
	Pattern        string `json:"pattern"`
	Interval       int32  `json:"interval"`
	Concurrent     int32  `json:"concurrent"`
	ConnectTimeout int32  `json:"connect_timeout"`
	RequestTimeout int32  `json:"request_timeout"`
}

type CleanerConfig struct {
	Interval int32 `json:"interval"`
}

type AgentConfig struct {
	Enabled bool           `json:"enabled"`
	Dsn     string         `json:"dsn"`
	MaxIdle int32          `json:"max_idle"`
	Plugin  *PluginConfig  `json:"plugin"`
	Cleaner *CleanerConfig `json:"cleaner"`
}

type LogConfig struct {
	Level string `json:"level"`
}

type GlobalConfig struct {
	Log       *LogConfig       `json:"log"`
	Http      *HttpConfig      `json:"http"`
	Index     *IndexConfig     `json:"index"`
	Collector *CollectorConfig `json:"collector"`
	Monitor   *MonitorConfig   `json:"monitor"`
	Agent     *AgentConfig     `json:"agent"`
	Host      string           `json:"host"`
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
		log.Fatalf("[F] config file: %s is not existent. maybe you need `mv cfg.example.json cfg.json`", cfg)
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

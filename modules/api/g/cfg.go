package g

import (
	"encoding/json"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/toolkits/file"
)

type DatabaseConfig struct {
	Addr           string `json:"addr"`
	MaxIdle        int    `json:"max_idle"`
	MaxConnections int    `json:"max_connections"`
	WaitTimeout    int    `json:"wait_timeout"`
}

type GraphsConfig struct {
	MaxConnections int               `json:"max_connections"`
	MaxIdle        int               `json:"max_idle"`
	ConnectTimeout int               `json:"connect_timeout"`
	ExecuteTimeout int               `json:"execute_timeout"`
	Replicas       int               `json:"replicas"`
	Cluster        map[string]string `json:"cluster"`
}

type RPCConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

type StatsdConfig struct {
	Enabled bool `json:"enabled"`
	Port    int  `json:"port"`
}

type LogConfig struct {
	Level string `json:"level"`
}

type GlobalConfig struct {
	Log            *LogConfig      `json:"log"`
	Listen         string          `json:"listen"`
	AccessControl  bool            `json:"access_control"`
	SignupDisable  bool            `json:"signup_disable"`
	SkipAuth       bool            `json:"skip_auth"`
	DefaultToken   string          `json:"default_token"`
	GenDoc         bool            `json:"gen_doc"`
	GenDocPath     string          `json:"gen_doc_path"`
	MetricListFile string          `json:"metric_list_file"`
	RPC            *RPCConfig      `json:"rpc"`
	Statsd         *StatsdConfig   `json:"statsd"`
	Database       *DatabaseConfig `json:"database"`
	Graphs         *GraphsConfig   `json:"graphs"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
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

	lock.Lock()
	defer lock.Unlock()

	config = &c

	log.Debugf("[D] read config file: %s successfully", cfg)
}
